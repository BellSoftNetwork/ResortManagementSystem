package handlers

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

// HealthHandler 구조체
type HealthHandler struct {
	db    *gorm.DB
	redis *redis.Client
}

// NewHealthHandler 생성자
func NewHealthHandler(db *gorm.DB, redis *redis.Client) *HealthHandler {
	return &HealthHandler{
		db:    db,
		redis: redis,
	}
}

// HealthStatus 응답 구조체
type HealthStatus struct {
	Status     string               `json:"status"`
	Components map[string]Component `json:"components,omitempty"`
}

// Component 상태 구조체
type Component struct {
	Status  string                 `json:"status"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// Health 기본 헬스체크
// @Summary 헬스체크
// @Description 애플리케이션의 전반적인 상태를 확인합니다
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} HealthStatus
// @Success 503 {object} HealthStatus
// @Router /actuator/health [get]
func (h *HealthHandler) Health(c *gin.Context) {
	status := "UP"
	components := make(map[string]Component)

	// MySQL 체크
	mysqlStatus := h.checkMySQL()
	components["db"] = mysqlStatus
	if mysqlStatus.Status != "UP" {
		status = "DOWN"
	}

	// Redis 체크
	redisStatus := h.checkRedis()
	components["redis"] = redisStatus
	if redisStatus.Status != "UP" {
		status = "DOWN"
	}

	// 응답 생성
	response := HealthStatus{
		Status:     status,
		Components: components,
	}

	// 상태 코드 결정
	statusCode := http.StatusOK
	if status == "DOWN" {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, response)
}

// Liveness 라이브니스 프로브
// @Summary 라이브니스 프로브
// @Description Kubernetes 라이브니스 프로브를 위한 엔드포인트
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} HealthStatus
// @Router /actuator/health/liveness [get]
func (h *HealthHandler) Liveness(c *gin.Context) {
	// 라이브니스는 애플리케이션이 살아있는지만 확인
	// 데드락이나 무한루프 등의 상태가 아니면 항상 UP
	response := HealthStatus{
		Status: "UP",
	}
	c.JSON(http.StatusOK, response)
}

// Readiness 레디니스 프로브
// @Summary 레디니스 프로브
// @Description Kubernetes 레디니스 프로브를 위한 엔드포인트
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} HealthStatus
// @Success 503 {object} HealthStatus
// @Router /actuator/health/readiness [get]
func (h *HealthHandler) Readiness(c *gin.Context) {
	status := "UP"
	components := make(map[string]Component)

	// MySQL 체크 (필수)
	mysqlStatus := h.checkMySQL()
	components["db"] = mysqlStatus
	if mysqlStatus.Status != "UP" {
		status = "DOWN"
	}

	// Redis 체크 (옵션)
	redisStatus := h.checkRedis()
	components["redis"] = redisStatus
	// Redis는 다운되어도 서비스는 가능하므로 readiness에는 영향 없음

	// 응답 생성
	response := HealthStatus{
		Status:     status,
		Components: components,
	}

	// 상태 코드 결정
	statusCode := http.StatusOK
	if status == "DOWN" {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, response)
}

// checkMySQL MySQL 연결 상태 확인
func (h *HealthHandler) checkMySQL() Component {
	component := Component{
		Status:  "UP",
		Details: make(map[string]interface{}),
	}

	// DB가 설정되지 않았으면 스킵
	if h.db == nil {
		component.Status = "UNKNOWN"
		component.Details["message"] = "Database not configured"
		return component
	}

	// 연결 테스트
	start := time.Now()

	sqlDB, err := h.db.DB()
	if err != nil {
		component.Status = "DOWN"
		component.Details["error"] = err.Error()
		return component
	}

	// Ping 테스트
	err = sqlDB.Ping()
	if err != nil {
		component.Status = "DOWN"
		component.Details["error"] = err.Error()
		return component
	}

	// 응답 시간
	component.Details["responseTime"] = time.Since(start).Milliseconds()

	// 연결 통계
	stats := sqlDB.Stats()
	component.Details["openConnections"] = stats.OpenConnections
	component.Details["inUse"] = stats.InUse
	component.Details["idle"] = stats.Idle

	return component
}

// checkRedis Redis 연결 상태 확인
func (h *HealthHandler) checkRedis() Component {
	component := Component{
		Status:  "UP",
		Details: make(map[string]interface{}),
	}

	// Redis가 설정되지 않았으면 스킵
	if h.redis == nil {
		component.Status = "UNKNOWN"
		component.Details["message"] = "Redis not configured"
		return component
	}

	// Ping 테스트
	start := time.Now()

	pong, err := h.redis.Ping(context.Background()).Result()
	if err != nil {
		component.Status = "DOWN"
		component.Details["error"] = err.Error()
		return component
	}

	// 응답 시간
	component.Details["responseTime"] = time.Since(start).Milliseconds()
	component.Details["response"] = pong

	// Redis 정보
	info, err := h.redis.Info(context.Background(), "server").Result()
	if err == nil {
		// 간단한 버전 정보만 추출
		for _, line := range strings.Split(info, "\n") {
			if strings.HasPrefix(line, "redis_version:") {
				component.Details["version"] = strings.TrimSpace(strings.TrimPrefix(line, "redis_version:"))
				break
			}
		}
	}

	return component
}
