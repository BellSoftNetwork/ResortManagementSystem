package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gitlab.bellsoft.net/rms/api-core/internal/audit"
	"gitlab.bellsoft.net/rms/api-core/internal/config"
	"gitlab.bellsoft.net/rms/api-core/internal/database"
	"gitlab.bellsoft.net/rms/api-core/internal/handlers"
	"gitlab.bellsoft.net/rms/api-core/internal/middleware"
	"gitlab.bellsoft.net/rms/api-core/internal/repositories"
	"gitlab.bellsoft.net/rms/api-core/internal/services"
	"gitlab.bellsoft.net/rms/api-core/pkg/auth"
)

func main() {
	// Parse command line flags
	var (
		runMigrations = flag.Bool("migrate", false, "Run database migrations before starting the server")
		migrateOnly   = flag.Bool("migrate-only", false, "Run database migrations and exit")
	)
	flag.Parse()

	cfg := config.Load()

	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-run migrations in development mode or if explicitly requested
	autoMigrate := cfg.Environment == "local" || cfg.Environment == "development"
	if *runMigrations || *migrateOnly || autoMigrate {
		log.Println("Running database migrations...")
		if err := database.Migrate(db); err != nil {
			log.Fatal("Failed to run migrations:", err)
		}
		log.Println("Migrations completed successfully")

		if *migrateOnly {
			log.Println("Migration-only mode, exiting...")
			os.Exit(0)
		}
	}

	redis, err := database.InitRedis(cfg)
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}

	jwtService := auth.NewJWTService(cfg.JWT.Secret, cfg.JWT.AccessTokenExpiry, cfg.JWT.RefreshTokenExpiry)

	userRepo := repositories.NewUserRepository(db)
	roomRepo := repositories.NewRoomRepository(db)
	roomGroupRepo := repositories.NewRoomGroupRepository(db)
	reservationRepo := repositories.NewReservationRepository(db)
	paymentMethodRepo := repositories.NewPaymentMethodRepository(db)
	loginAttemptRepo := repositories.NewLoginAttemptRepository(db)
	// reservationRoomRepo := repositories.NewReservationRoomRepository(db) // Not used

	// Initialize audit service first
	auditService := audit.NewService(db)
	// Register audit hooks for GORM - disabled due to JSON depth issue
	audit.RegisterHooks(db, auditService)

	authService := services.NewAuthService(userRepo, loginAttemptRepo, jwtService, cfg)
	userService := services.NewUserService(userRepo)
	roomService := services.NewRoomService(roomRepo, roomGroupRepo, auditService)
	roomGroupService := services.NewRoomGroupService(roomGroupRepo)
	reservationService := services.NewReservationService(reservationRepo, roomRepo, paymentMethodRepo, auditService)
	paymentMethodService := services.NewPaymentMethodService(paymentMethodRepo)
	configService := services.NewConfigService(cfg)
	developmentService := services.NewDevelopmentServiceV2(db)
	historyService := services.NewHistoryService(auditService, userService)

	authHandler := handlers.NewAuthHandler(authService)
	mainHandler := handlers.NewMainHandler(configService, userRepo)
	userHandler := handlers.NewUserHandler(userService)
	roomHandler := handlers.NewRoomHandler(roomService, userService, historyService)
	roomGroupHandler := handlers.NewRoomGroupHandler(roomGroupService, reservationService, userService)
	reservationHandler := handlers.NewReservationHandler(reservationService, userService, historyService)
	paymentMethodHandler := handlers.NewPaymentMethodHandler(paymentMethodService)
	developmentHandler := handlers.NewDevelopmentHandler(developmentService)
	healthHandler := handlers.NewHealthHandler(db, redis)
	docsHandler := handlers.NewDocsHandler()
	auditHandler := handlers.NewAuditHandler(auditService)

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.ErrorHandler())

	corsConfig := cors.Config{
		AllowOrigins:     cfg.CORS.AllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	router.Use(cors.New(corsConfig))

	router.NoRoute(func(c *gin.Context) {
		c.File("./public/index.html")
	})

	setupRoutes(router, authHandler, mainHandler, userHandler, roomHandler, roomGroupHandler, reservationHandler, paymentMethodHandler, developmentHandler, healthHandler, docsHandler, auditHandler, jwtService, cfg)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	log.Printf("Server started on port %d", cfg.Server.Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}

func setupRoutes(r *gin.Engine, authHandler *handlers.AuthHandler, mainHandler *handlers.MainHandler,
	userHandler *handlers.UserHandler, roomHandler *handlers.RoomHandler,
	roomGroupHandler *handlers.RoomGroupHandler, reservationHandler *handlers.ReservationHandler,
	paymentMethodHandler *handlers.PaymentMethodHandler, developmentHandler *handlers.DevelopmentHandler,
	healthHandler *handlers.HealthHandler, docsHandler *handlers.DocsHandler, auditHandler *handlers.AuditHandler,
	jwtService *auth.JWTService, cfg *config.Config) {

	// Health check endpoints (Spring Boot Actuator compatible)
	actuator := r.Group("/actuator")
	{
		actuator.GET("/health", healthHandler.Health)
		actuator.GET("/health/liveness", healthHandler.Liveness)
		actuator.GET("/health/readiness", healthHandler.Readiness)
	}

	// Documentation endpoints
	docs := r.Group("/docs")
	{
		docs.GET("/schema", docsHandler.GetOpenAPISchema)
		docs.GET("/swagger-ui", docsHandler.GetSwaggerUI)
	}

	// API v1 for api-core
	api := r.Group("/api/v1")
	{
		api.GET("/env", mainHandler.GetEnvironment)
		api.GET("/config", mainHandler.GetConfig)

		authRoutes := api.Group("/auth")
		{
			authRoutes.POST("/register", authHandler.Register)
			authRoutes.POST("/login", authHandler.Login)
			authRoutes.POST("/refresh", authHandler.RefreshToken)
		}

		authenticated := api.Group("")
		authenticated.Use(middleware.AuthMiddleware(jwtService))
		authenticated.Use(middleware.AuditMiddleware())
		{
			myRoutes := authenticated.Group("/my")
			{
				myRoutes.GET("", userHandler.GetCurrentUser)
				myRoutes.POST("", userHandler.GetCurrentUser)
				myRoutes.PATCH("", userHandler.UpdateCurrentUser)
			}

			adminRoutes := authenticated.Group("/admin")
			adminRoutes.Use(middleware.RoleMiddleware("ADMIN", "SUPER_ADMIN"))
			{
				accountRoutes := adminRoutes.Group("/accounts")
				{
					accountRoutes.GET("", userHandler.ListUsers)
					accountRoutes.POST("", userHandler.CreateUser)
					accountRoutes.PATCH("/:id", userHandler.UpdateUser)
				}
				adminRoutes.GET("/audit-logs", auditHandler.ListAuditLogs)
				adminRoutes.GET("/audit-logs/:id", auditHandler.GetAuditLog)
			}

			roomRoutes := authenticated.Group("/rooms")
			{
				roomRoutes.GET("", roomHandler.ListRooms)
				roomRoutes.GET("/:id", roomHandler.GetRoom)
				roomRoutes.POST("", middleware.RoleMiddleware("ADMIN", "SUPER_ADMIN"), roomHandler.CreateRoom)
				roomRoutes.PATCH("/:id", middleware.RoleMiddleware("ADMIN", "SUPER_ADMIN"), roomHandler.UpdateRoom)
				roomRoutes.DELETE("/:id", middleware.RoleMiddleware("ADMIN", "SUPER_ADMIN"), roomHandler.DeleteRoom)
				roomRoutes.GET("/:id/histories", middleware.RoleMiddleware("ADMIN", "SUPER_ADMIN"), roomHandler.GetRoomHistories)
			}

			roomGroupRoutes := authenticated.Group("/room-groups")
			{
				roomGroupRoutes.GET("", roomGroupHandler.ListRoomGroups)
				roomGroupRoutes.GET("/:id", roomGroupHandler.GetRoomGroup)
				roomGroupRoutes.POST("", middleware.RoleMiddleware("ADMIN", "SUPER_ADMIN"), roomGroupHandler.CreateRoomGroup)
				roomGroupRoutes.PATCH("/:id", middleware.RoleMiddleware("ADMIN", "SUPER_ADMIN"), roomGroupHandler.UpdateRoomGroup)
				roomGroupRoutes.DELETE("/:id", middleware.RoleMiddleware("ADMIN", "SUPER_ADMIN"), roomGroupHandler.DeleteRoomGroup)
			}

			reservationRoutes := authenticated.Group("/reservations")
			{
				reservationRoutes.GET("", reservationHandler.ListReservations)
				reservationRoutes.GET("/:id", reservationHandler.GetReservation)
				reservationRoutes.POST("", middleware.RoleMiddleware("ADMIN", "SUPER_ADMIN"), reservationHandler.CreateReservation)
				reservationRoutes.PATCH("/:id", middleware.RoleMiddleware("ADMIN", "SUPER_ADMIN"), reservationHandler.UpdateReservation)
				reservationRoutes.DELETE("/:id", middleware.RoleMiddleware("ADMIN", "SUPER_ADMIN"), reservationHandler.DeleteReservation)
				reservationRoutes.GET("/:id/histories", middleware.RoleMiddleware("ADMIN", "SUPER_ADMIN"), reservationHandler.GetReservationHistories)
			}

			reservationStatsRoutes := authenticated.Group("/reservation-statistics")
			{
				reservationStatsRoutes.GET("", reservationHandler.GetReservationStatistics)
			}

			paymentMethodRoutes := authenticated.Group("/payment-methods")
			{
				paymentMethodRoutes.GET("", paymentMethodHandler.ListPaymentMethods)
				paymentMethodRoutes.GET("/:id", paymentMethodHandler.GetPaymentMethod)
				paymentMethodRoutes.POST("", middleware.RoleMiddleware("ADMIN", "SUPER_ADMIN"), paymentMethodHandler.CreatePaymentMethod)
				paymentMethodRoutes.PATCH("/:id", middleware.RoleMiddleware("ADMIN", "SUPER_ADMIN"), paymentMethodHandler.UpdatePaymentMethod)
				paymentMethodRoutes.DELETE("/:id", middleware.RoleMiddleware("ADMIN", "SUPER_ADMIN"), paymentMethodHandler.DeletePaymentMethod)
			}

			// Development endpoints (only available in non-production environments)
			if cfg.Environment != "production" {
				devRoutes := authenticated.Group("/dev")
				devRoutes.Use(middleware.DevelopmentOnlyMiddleware())
				devRoutes.Use(middleware.RoleMiddleware("SUPER_ADMIN"))
				{
					devRoutes.POST("/test-data", developmentHandler.GenerateTestData)
				}
			}
		}
	}
}
