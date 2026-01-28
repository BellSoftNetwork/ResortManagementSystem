package services

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/sirupsen/logrus"
	"gitlab.bellsoft.net/rms/api-core/internal/models"
	"gorm.io/gorm"
)

// ReservationGenerationOptions represents options for generating reservation data
type ReservationGenerationOptions struct {
	StartDate           *time.Time
	EndDate             *time.Time
	RegularReservations *int
	MonthlyReservations *int
}

type DevelopmentService interface {
	GenerateTestData(dataType string, reservationOptions *ReservationGenerationOptions) (map[string]interface{}, error)
}

type developmentService struct {
	db *gorm.DB
}

func NewDevelopmentServiceV2(db *gorm.DB) DevelopmentService {
	return &developmentService{
		db: db,
	}
}

func (s *developmentService) resetData() error {
	logrus.Info("=== resetData called ===")

	tables := []string{
		"reservation_room",
		"reservation",
		"room",
		"room_group",
		"payment_method",
	}

	activeTime := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)

	for _, table := range tables {
		logrus.Infof("Soft deleting records from table: %s", table)
		if err := s.db.Exec(fmt.Sprintf("UPDATE %s SET deleted_at = NOW() WHERE deleted_at = ?", table), activeTime).Error; err != nil {
			logrus.Errorf("Failed to reset table %s: %v", table, err)
			return fmt.Errorf("failed to reset table %s: %w", table, err)
		}
	}

	logrus.Info("=== resetData completed ===")
	return nil
}

func (s *developmentService) generateEssentialData() (map[string]interface{}, error) {
	logrus.Info("=== GenerateEssentialData called ===")
	result := make(map[string]interface{})

	// Get the super admin user ID
	var adminUser models.User
	if err := s.db.Where("role = ?", models.UserRoleSuperAdmin).First(&adminUser).Error; err != nil {
		logrus.Errorf("Failed to find super admin user: %v", err)
		return nil, fmt.Errorf("failed to find super admin user: %w", err)
	}
	adminUserID := adminUser.ID
	logrus.Infof("Found admin user with ID: %d", adminUserID)

	// Create payment methods only for now
	paymentMethods := []models.PaymentMethod{
		{
			Name:                     "야놀자",
			CommissionRate:           0.15,
			RequireUnpaidAmountCheck: false,
			IsDefaultSelect:          false,
			Status:                   models.PaymentMethodStatusActive,
		},
		{
			Name:                     "통장입금",
			CommissionRate:           0.0,
			RequireUnpaidAmountCheck: true,
			IsDefaultSelect:          true,
			Status:                   models.PaymentMethodStatusActive,
		},
		{
			Name:                     "펜션다나와",
			CommissionRate:           0.13,
			RequireUnpaidAmountCheck: false,
			IsDefaultSelect:          false,
			Status:                   models.PaymentMethodStatusActive,
		},
		{
			Name:                     "현장결제",
			CommissionRate:           0.0,
			RequireUnpaidAmountCheck: true,
			IsDefaultSelect:          false,
			Status:                   models.PaymentMethodStatusActive,
		},
	}

	createdPaymentMethods := make([]*models.PaymentMethod, 0)
	for _, pm := range paymentMethods {
		// Check if already exists
		var count int64
		s.db.Model(&models.PaymentMethod{}).Where("name = ? AND deleted_at = ?", pm.Name, time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)).Count(&count)
		if count > 0 {
			logrus.Infof("Payment method '%s' already exists, skipping", pm.Name)
			continue
		}

		if err := s.db.Session(&gorm.Session{SkipHooks: true}).Omit("CreatedByUser", "UpdatedByUser").Create(&pm).Error; err != nil {
			logrus.Errorf("Failed to create payment method '%s': %v", pm.Name, err)
			// Continue processing other items instead of returning error
			continue
		}
		logrus.Infof("Created payment method '%s'", pm.Name)
		createdPaymentMethods = append(createdPaymentMethods, &pm)
	}

	result["paymentMethods"] = map[string]interface{}{
		"created": len(createdPaymentMethods),
		"items":   createdPaymentMethods,
	}

	// Create room groups using raw SQL to avoid GORM hooks and relation issues
	logrus.Info("Creating room groups")
	roomGroups := []struct {
		Name         string
		PeekPrice    int
		OffPeekPrice int
		Description  string
	}{
		{"5층 투룸", 150000, 120000, "5층 투룸 객실"},
		{"4층 투룸", 140000, 110000, "4층 투룸 객실"},
		{"3층 투룸", 130000, 100000, "3층 투룸 객실"},
		{"2층 투룸", 120000, 90000, "2층 투룸 객실"},
		{"원룸(20평형)", 100000, 80000, "원룸 20평형"},
		{"부대시설", 50000, 50000, "부대시설"},
		{"대형룸", 200000, 160000, "대형룸"},
	}

	createdRoomGroups := 0
	for _, rg := range roomGroups {
		// Check if already exists
		var count int64
		s.db.Table("room_group").Where("name = ? AND deleted_at = ?", rg.Name, time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)).Count(&count)
		if count > 0 {
			logrus.Infof("Room group '%s' already exists, skipping", rg.Name)
			continue
		}

		// Use raw SQL to insert
		if err := s.db.Exec(`INSERT INTO room_group (name, peek_price, off_peek_price, description, created_by, updated_by, created_at, updated_at, deleted_at)
			VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW(), '1970-01-01 00:00:00')`,
			rg.Name, rg.PeekPrice, rg.OffPeekPrice, rg.Description, adminUserID, adminUserID).Error; err != nil {
			logrus.Errorf("Failed to create room group '%s': %v", rg.Name, err)
			// Continue processing other items instead of returning error
			continue
		}
		logrus.Infof("Created room group '%s'", rg.Name)
		createdRoomGroups++
	}

	result["roomGroups"] = map[string]interface{}{
		"created": createdRoomGroups,
		"items":   roomGroups,
	}

	// Get room groups with their IDs
	roomGroupMap := make(map[string]uint)
	rows, err := s.db.Table("room_group").Select("id, name").Where("deleted_at = ?", time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)).Rows()
	if err != nil {
		logrus.Errorf("Failed to get room groups: %v", err)
		return nil, fmt.Errorf("failed to get room groups: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id uint
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			continue
		}
		roomGroupMap[name] = id
	}

	// Create rooms
	logrus.Info("Creating rooms")
	roomsData := []struct {
		GroupName string
		Rooms     []struct {
			Number string
			Status int // 1: Normal, 2: Inactive, 3: Construction
		}
	}{
		{
			GroupName: "5층 투룸",
			Rooms: func() []struct {
				Number string
				Status int
			} {
				rooms := make([]struct {
					Number string
					Status int
				}, 10)
				for i := 0; i < 10; i++ {
					rooms[i] = struct {
						Number string
						Status int
					}{
						Number: fmt.Sprintf("%d호", 501+i),
						Status: 1, // Normal
					}
				}
				return rooms
			}(),
		},
		{
			GroupName: "4층 투룸",
			Rooms: func() []struct {
				Number string
				Status int
			} {
				rooms := make([]struct {
					Number string
					Status int
				}, 12)
				for i := 0; i < 12; i++ {
					rooms[i] = struct {
						Number string
						Status int
					}{
						Number: fmt.Sprintf("%d호", 401+i),
						Status: 1, // Normal
					}
				}
				return rooms
			}(),
		},
		{
			GroupName: "3층 투룸",
			Rooms: func() []struct {
				Number string
				Status int
			} {
				rooms := make([]struct {
					Number string
					Status int
				}, 14)
				for i := 0; i < 14; i++ {
					rooms[i] = struct {
						Number string
						Status int
					}{
						Number: fmt.Sprintf("%d호", 301+i),
						Status: 1, // Normal
					}
				}
				return rooms
			}(),
		},
		{
			GroupName: "2층 투룸",
			Rooms: func() []struct {
				Number string
				Status int
			} {
				rooms := make([]struct {
					Number string
					Status int
				}, 14)
				for i := 0; i < 14; i++ {
					status := 1   // Normal
					if i+1 == 8 { // 208호는 이용 불가
						status = 2 // Inactive
					}
					rooms[i] = struct {
						Number string
						Status int
					}{
						Number: fmt.Sprintf("%d호", 201+i),
						Status: status,
					}
				}
				return rooms
			}(),
		},
		{
			GroupName: "원룸(20평형)",
			Rooms: []struct {
				Number string
				Status int
			}{
				{Number: "403호", Status: 1}, // Normal
				{Number: "410호", Status: 1}, // Normal
				{Number: "303호", Status: 1}, // Normal
				{Number: "304호", Status: 1}, // Normal
				{Number: "311호", Status: 1}, // Normal
				{Number: "312호", Status: 1}, // Normal
				{Number: "203호", Status: 3}, // Construction
				{Number: "204호", Status: 1}, // Normal
				{Number: "211호", Status: 1}, // Normal
				{Number: "212호", Status: 3}, // Construction
			},
		},
		{
			GroupName: "부대시설",
			Rooms: []struct {
				Number string
				Status int
			}{
				{Number: "1층세미나실", Status: 1},
				{Number: "카페(내부)", Status: 1},
				{Number: "카페(외부)", Status: 1},
				{Number: "음향시설", Status: 1},
				{Number: "바베큐그릴 고급형", Status: 1},
				{Number: "소강당", Status: 1},
			},
		},
		{
			GroupName: "대형룸",
			Rooms: []struct {
				Number string
				Status int
			}{
				{Number: "207호208호", Status: 1}, // Normal
			},
		},
	}

	createdRooms := 0
	for _, groupData := range roomsData {
		roomGroupID, ok := roomGroupMap[groupData.GroupName]
		if !ok {
			logrus.Warnf("Room group '%s' not found, skipping rooms", groupData.GroupName)
			continue
		}

		for _, roomData := range groupData.Rooms {
			// Check if already exists
			var count int64
			s.db.Table("room").Where("number = ? AND deleted_at = ?", roomData.Number, time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)).Count(&count)
			if count > 0 {
				logrus.Infof("Room '%s' already exists, skipping", roomData.Number)
				continue
			}

			// Use raw SQL to insert
			if err := s.db.Exec(`INSERT INTO room (number, room_group_id, status, note, created_by, updated_by, created_at, updated_at, deleted_at)
				VALUES (?, ?, ?, '', ?, ?, NOW(), NOW(), '1970-01-01 00:00:00')`,
				roomData.Number, roomGroupID, roomData.Status, adminUserID, adminUserID).Error; err != nil {
				logrus.Errorf("Failed to create room '%s': %v", roomData.Number, err)
				continue
			}
			logrus.Debugf("Created room '%s'", roomData.Number)
			createdRooms++
		}
	}

	result["rooms"] = map[string]interface{}{
		"created": createdRooms,
		"count":   createdRooms,
	}

	logrus.Info("=== GenerateEssentialData completed ===")
	return result, nil
}

func (s *developmentService) generateEssentialData_OLD() (map[string]interface{}, error) {
	result := make(map[string]interface{})

	logrus.Info("GenerateEssentialData called")

	// Get the super admin user ID
	var adminUser models.User
	if err := s.db.Where("role = ?", models.UserRoleSuperAdmin).First(&adminUser).Error; err != nil {
		logrus.Errorf("Failed to find super admin user: %v", err)
		return nil, fmt.Errorf("failed to find super admin user: %w", err)
	}
	adminUserID := adminUser.ID
	logrus.Infof("Found admin user with ID: %d", adminUserID)

	// Use direct DB connection instead of transaction to avoid issues
	db := s.db

	// Create payment methods
	paymentMethods := []models.PaymentMethod{
		{
			Name:                     "야놀자",
			CommissionRate:           0.15,
			RequireUnpaidAmountCheck: false,
			IsDefaultSelect:          false,
			Status:                   models.PaymentMethodStatusActive,
		},
		{
			Name:                     "통장입금",
			CommissionRate:           0.0,
			RequireUnpaidAmountCheck: true,
			IsDefaultSelect:          true,
			Status:                   models.PaymentMethodStatusActive,
		},
		{
			Name:                     "펜션다나와",
			CommissionRate:           0.13,
			RequireUnpaidAmountCheck: false,
			IsDefaultSelect:          false,
			Status:                   models.PaymentMethodStatusActive,
		},
		{
			Name:                     "현장결제",
			CommissionRate:           0.0,
			RequireUnpaidAmountCheck: true,
			IsDefaultSelect:          false,
			Status:                   models.PaymentMethodStatusActive,
		},
	}

	createdPaymentMethods := make([]*models.PaymentMethod, 0)
	for _, pm := range paymentMethods {
		// Check if already exists
		var count int64
		db.Model(&models.PaymentMethod{}).Where("name = ? AND deleted_at = ?", pm.Name, time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)).Count(&count)
		if count > 0 {
			logrus.Infof("Payment method '%s' already exists, skipping", pm.Name)
			continue
		}

		if err := db.Session(&gorm.Session{SkipHooks: true}).Omit("CreatedByUser", "UpdatedByUser").Create(&pm).Error; err != nil {
			return nil, fmt.Errorf("failed to create payment method '%s': %w", pm.Name, err)
		}
		createdPaymentMethods = append(createdPaymentMethods, &pm)
	}

	// Create room groups
	logrus.Info("Starting room group creation")
	// Temporarily comment out to test
	/*
		roomGroups := []models.RoomGroup{
			{
				BaseMustAuditEntity: models.BaseMustAuditEntity{
					CreatedBy: adminUserID,
					UpdatedBy: adminUserID,
				},
				Name:         "5층 투룸",
				PeekPrice:    150000,
				OffPeekPrice: 120000,
				Description:  "5층 투룸 객실",
			},
			{
				BaseMustAuditEntity: models.BaseMustAuditEntity{
					CreatedBy: adminUserID,
					UpdatedBy: adminUserID,
				},
				Name:         "4층 투룸",
				PeekPrice:    140000,
				OffPeekPrice: 110000,
				Description:  "4층 투룸 객실",
			},
			{
				BaseMustAuditEntity: models.BaseMustAuditEntity{
					CreatedBy: adminUserID,
					UpdatedBy: adminUserID,
				},
				Name:         "3층 투룸",
				PeekPrice:    130000,
				OffPeekPrice: 100000,
				Description:  "3층 투룸 객실",
			},
			{
				BaseMustAuditEntity: models.BaseMustAuditEntity{
					CreatedBy: adminUserID,
					UpdatedBy: adminUserID,
				},
				Name:         "2층 투룸",
				PeekPrice:    120000,
				OffPeekPrice: 90000,
				Description:  "2층 투룸 객실",
			},
			{
				BaseMustAuditEntity: models.BaseMustAuditEntity{
					CreatedBy: adminUserID,
					UpdatedBy: adminUserID,
				},
				Name:         "원룸(20평형)",
				PeekPrice:    100000,
				OffPeekPrice: 80000,
				Description:  "원룸 20평형",
			},
			{
				BaseMustAuditEntity: models.BaseMustAuditEntity{
					CreatedBy: adminUserID,
					UpdatedBy: adminUserID,
				},
				Name:         "부대시설",
				PeekPrice:    50000,
				OffPeekPrice: 50000,
				Description:  "부대시설",
			},
			{
				BaseMustAuditEntity: models.BaseMustAuditEntity{
					CreatedBy: adminUserID,
					UpdatedBy: adminUserID,
				},
				Name:         "대형룸",
				PeekPrice:    200000,
				OffPeekPrice: 160000,
				Description:  "대형룸",
			},
		}
	*/

	createdRoomGroups := make([]*models.RoomGroup, 0)
	// roomGroupMap := make(map[string]*models.RoomGroup) // commented out
	createdRooms := make([]*models.Room, 0) // moved here to fix compilation

	/*
		logrus.Infof("Processing %d room groups", len(roomGroups))
		for _, rg := range roomGroups {
			logrus.Infof("Processing room group: %s", rg.Name)
			// Check if already exists
			var count int64
			db.Model(&models.RoomGroup{}).Where("name = ? AND deleted_at = ?", rg.Name, time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)).Count(&count)
			if count > 0 {
				logrus.Infof("Room group '%s' already exists, skipping", rg.Name)
				// Get just the ID
				var existingID uint
				db.Model(&models.RoomGroup{}).Select("id").Where("name = ? AND deleted_at = ?", rg.Name, time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)).Scan(&existingID)

				existing := models.RoomGroup{
					BaseMustAuditEntity: models.BaseMustAuditEntity{
						BaseTimeEntity: models.BaseTimeEntity{
							BaseEntity: models.BaseEntity{
								ID: existingID,
							},
						},
					},
					Name: rg.Name,
				}
				createdRoomGroups = append(createdRoomGroups, &existing)
				roomGroupMap[rg.Name] = &existing
				continue
			}

			// Use raw SQL to avoid hooks and relation loading issues
			logrus.Infof("Creating room group '%s' with raw SQL", rg.Name)
			result := db.Exec(`INSERT INTO room_group (name, peek_price, off_peek_price, description, created_by, updated_by, created_at, updated_at, deleted_at)
				VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW(), '1970-01-01 00:00:00')`,
				rg.Name, rg.PeekPrice, rg.OffPeekPrice, rg.Description, adminUserID, adminUserID)
			if result.Error != nil {
				logrus.Errorf("Failed to create room group '%s' with raw SQL: %v", rg.Name, result.Error)
				return nil, fmt.Errorf("failed to create room group '%s': %w", rg.Name, result.Error)
			}
			logrus.Infof("Successfully created room group '%s' with raw SQL", rg.Name)

			// Get the last insert ID manually
			var lastInsertID uint
			if err := db.Raw("SELECT LAST_INSERT_ID()").Scan(&lastInsertID).Error; err != nil {
				return nil, fmt.Errorf("failed to get last insert ID for room group '%s': %w", rg.Name, err)
			}

			// Create a minimal room group object with just the ID
			createdRG := models.RoomGroup{
				BaseMustAuditEntity: models.BaseMustAuditEntity{
					BaseTimeEntity: models.BaseTimeEntity{
						BaseEntity: models.BaseEntity{
							ID: lastInsertID,
						},
					},
				},
				Name:         rg.Name,
				PeekPrice:    rg.PeekPrice,
				OffPeekPrice: rg.OffPeekPrice,
				Description:  rg.Description,
			}

			createdRoomGroups = append(createdRoomGroups, &createdRG)
			roomGroupMap[rg.Name] = &createdRG
		}
	*/

	// Create rooms
	// Temporarily comment out
	/*
		roomsData := []struct {
			GroupName string
			Rooms     []struct {
				Number string
				Status models.RoomStatus
			}
		}{
			{
				GroupName: "5층 투룸",
				Rooms: func() []struct {
					Number string
					Status models.RoomStatus
				} {
					rooms := make([]struct {
						Number string
						Status models.RoomStatus
					}, 10)
					for i := 0; i < 10; i++ {
						rooms[i] = struct {
							Number string
							Status models.RoomStatus
						}{
							Number: fmt.Sprintf("50%d호", i+1),
							Status: models.RoomStatusNormal,
						}
					}
					return rooms
				}(),
			},
			{
				GroupName: "4층 투룸",
				Rooms: func() []struct {
					Number string
					Status models.RoomStatus
				} {
					rooms := make([]struct {
						Number string
						Status models.RoomStatus
					}, 12)
					for i := 0; i < 12; i++ {
						rooms[i] = struct {
							Number string
							Status models.RoomStatus
						}{
							Number: fmt.Sprintf("40%d호", i+1),
							Status: models.RoomStatusNormal,
						}
					}
					return rooms
				}(),
			},
			{
				GroupName: "3층 투룸",
				Rooms: func() []struct {
					Number string
					Status models.RoomStatus
				} {
					rooms := make([]struct {
						Number string
						Status models.RoomStatus
					}, 14)
					for i := 0; i < 14; i++ {
						rooms[i] = struct {
							Number string
							Status models.RoomStatus
						}{
							Number: fmt.Sprintf("30%d호", i+1),
							Status: models.RoomStatusNormal,
						}
					}
					return rooms
				}(),
			},
			{
				GroupName: "2층 투룸",
				Rooms: func() []struct {
					Number string
					Status models.RoomStatus
				} {
					rooms := make([]struct {
						Number string
						Status models.RoomStatus
					}, 14)
					for i := 0; i < 14; i++ {
						rooms[i] = struct {
							Number string
							Status models.RoomStatus
						}{
							Number: fmt.Sprintf("20%d호", i+1),
							Status: models.RoomStatusNormal,
						}
						// 208호는 이용 불가 상태
						if i == 7 {
							rooms[i].Status = models.RoomStatusInactive
						}
					}
					return rooms
				}(),
			},
			{
				GroupName: "원룸(20평형)",
				Rooms: []struct {
					Number string
					Status models.RoomStatus
				}{
					{Number: "403호", Status: models.RoomStatusNormal},
					{Number: "410호", Status: models.RoomStatusNormal},
					{Number: "303호", Status: models.RoomStatusNormal},
					{Number: "304호", Status: models.RoomStatusNormal},
					{Number: "311호", Status: models.RoomStatusNormal},
					{Number: "312호", Status: models.RoomStatusNormal},
					{Number: "203호", Status: models.RoomStatusConstruction}, // 공사 중
					{Number: "204호", Status: models.RoomStatusNormal},
					{Number: "211호", Status: models.RoomStatusNormal},
					{Number: "212호", Status: models.RoomStatusConstruction}, // 공사 중
				},
			},
			{
				GroupName: "부대시설",
				Rooms: []struct {
					Number string
					Status models.RoomStatus
				}{
					{Number: "1층세미나실", Status: models.RoomStatusNormal},
					{Number: "카페(내부)", Status: models.RoomStatusNormal},
					{Number: "카페(외부)", Status: models.RoomStatusNormal},
					{Number: "음향시설", Status: models.RoomStatusNormal},
					{Number: "바베큐그릴 고급형", Status: models.RoomStatusNormal},
					{Number: "소강당", Status: models.RoomStatusNormal},
				},
			},
			{
				GroupName: "대형룸",
				Rooms: []struct {
					Number string
					Status models.RoomStatus
				}{
					{Number: "207호208호", Status: models.RoomStatusNormal},
				},
			},
		}

		// createdRooms := make([]*models.Room, 0) - moved above
		// _ = roomGroupMap // avoid unused variable error
		/*
		for _, groupData := range roomsData {
			roomGroup, ok := roomGroupMap[groupData.GroupName]
			if !ok {
				continue
			}

			for _, roomData := range groupData.Rooms {
				// Check if already exists
				var count int64
				db.Model(&models.Room{}).Where("number = ? AND deleted_at = ?", roomData.Number, time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)).Count(&count)
				if count > 0 {
					logrus.Infof("Room '%s' already exists, skipping", roomData.Number)
					continue
				}

				room := models.Room{
					BaseMustAuditEntity: models.BaseMustAuditEntity{
						CreatedBy: adminUserID,
						UpdatedBy: adminUserID,
					},
					Number:      roomData.Number,
					RoomGroupID: roomGroup.ID,
					Status:      roomData.Status,
					Note:        "",
				}

				if err := db.Session(&gorm.Session{SkipHooks: true}).Omit("CreatedByUser", "UpdatedByUser", "RoomGroup").Create(&room).Error; err != nil {
						return nil, fmt.Errorf("failed to create room '%s': %w", room.Number, err)
				}
				createdRooms = append(createdRooms, &room)
			}
		}
	*/

	result["paymentMethods"] = map[string]interface{}{
		"created": len(createdPaymentMethods),
		"items":   createdPaymentMethods,
	}
	result["roomGroups"] = map[string]interface{}{
		"created": len(createdRoomGroups),
		"items":   createdRoomGroups,
	}
	result["rooms"] = map[string]interface{}{
		"created": len(createdRooms),
		"count":   len(createdRooms),
	}

	return result, nil
}

func (s *developmentService) generateReservationData() (map[string]interface{}, error) {
	logrus.Info("=== GenerateReservationData called ===")
	result := make(map[string]interface{})

	// Get the super admin user ID
	var adminUser models.User
	if err := s.db.Where("role = ?", models.UserRoleSuperAdmin).First(&adminUser).Error; err != nil {
		logrus.Errorf("Failed to find super admin user: %v", err)
		return nil, fmt.Errorf("failed to find super admin user: %w", err)
	}
	adminUserID := adminUser.ID
	logrus.Infof("Found admin user with ID: %d", adminUserID)

	// Get available rooms (excluding construction and facility rooms)
	var availableRooms []struct {
		ID     uint
		Number string
	}
	if err := s.db.Table("room r").
		Select("r.id, r.number").
		Joins("JOIN room_group rg ON r.room_group_id = rg.id").
		Where("r.status = ? AND r.deleted_at = ? AND rg.name NOT IN (?) AND rg.deleted_at = ?",
			1, // Normal status
			time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
			[]string{"부대시설"},
			time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)).
		Order("r.number").
		Scan(&availableRooms).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch available rooms: %w", err)
	}

	if len(availableRooms) == 0 {
		return nil, fmt.Errorf("no available rooms found")
	}

	logrus.Infof("Found %d available rooms", len(availableRooms))

	// Get payment methods
	var paymentMethodID uint
	if err := s.db.Table("payment_method").
		Select("id").
		Where("name = ? AND deleted_at = ?", "통장입금", time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)).
		Scan(&paymentMethodID).Error; err != nil {
		return nil, fmt.Errorf("failed to find payment method: %w", err)
	}

	today := time.Now().Truncate(24 * time.Hour)
	createdReservations := 0
	oneDayStay := 0
	twoDaysStay := 0
	monthlyRent := 0

	// 1. 현재 날짜를 입실일로 1박 2일 예약 10건
	for i := 0; i < 10 && i < len(availableRooms); i++ {
		checkIn := today
		checkOut := checkIn.AddDate(0, 0, 1)
		guestName := fmt.Sprintf("테스트고객%d", i+1)

		// Check if already exists
		var count int64
		s.db.Table("reservation").
			Where("name = ? AND stay_start_at = ? AND deleted_at = ?",
				guestName, checkIn, time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)).
			Count(&count)
		if count > 0 {
			logrus.Infof("Reservation for guest '%s' on %s already exists, skipping", guestName, checkIn.Format("2006-01-02"))
			continue
		}

		// Insert reservation using raw SQL (using legacy schema)
		var reservationID uint
		result := s.db.Exec(`INSERT INTO reservation (name, phone, people_count, payment_method_id,
			stay_start_at, stay_end_at, check_in_at, check_out_at,
			price, deposit, payment_amount, refund_amount, broker_fee,
			note, status, type, created_by, updated_by, created_at, updated_at, deleted_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW(), '1970-01-01 00:00:00')`,
			guestName, fmt.Sprintf("010-1234-%04d", 1000+i), 2, paymentMethodID,
			checkIn, checkOut, checkIn, checkOut,
			100000, 50000, 50000, 0, 0,
			"1박 2일 테스트 예약", 1, 0, adminUserID, adminUserID)
		if result.Error != nil {
			logrus.Errorf("Failed to create 1-night reservation %d: %v", i+1, result.Error)
			continue
		}
		// Get the last insert ID
		s.db.Raw("SELECT LAST_INSERT_ID()").Scan(&reservationID)

		// Add room to reservation
		if err := s.db.Exec(`INSERT INTO reservation_room (reservation_id, room_id, created_by, updated_by, created_at, updated_at, deleted_at)
			VALUES (?, ?, ?, ?, NOW(), NOW(), '1970-01-01 00:00:00')`,
			reservationID, availableRooms[i].ID, adminUserID, adminUserID).Error; err != nil {
			logrus.Errorf("Failed to create reservation room: %v", err)
			continue
		}

		createdReservations++
		oneDayStay++
	}

	// 2. 현재 날짜를 입실일로 2박 3일 예약 4건
	for i := 0; i < 4 && (10+i) < len(availableRooms); i++ {
		checkIn := today
		checkOut := checkIn.AddDate(0, 0, 2)
		guestName := fmt.Sprintf("장기고객%d", i+1)
		roomIndex := 10 + i

		// Check if already exists
		var count int64
		s.db.Table("reservation").
			Where("name = ? AND stay_start_at = ? AND deleted_at = ?",
				guestName, checkIn, time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)).
			Count(&count)
		if count > 0 {
			logrus.Infof("Reservation for guest '%s' on %s already exists, skipping", guestName, checkIn.Format("2006-01-02"))
			continue
		}

		// Insert reservation
		var reservationID uint
		result := s.db.Exec(`INSERT INTO reservation (name, phone, people_count, payment_method_id,
			stay_start_at, stay_end_at, check_in_at, check_out_at,
			price, deposit, payment_amount, refund_amount, broker_fee,
			note, status, type, created_by, updated_by, created_at, updated_at, deleted_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW(), '1970-01-01 00:00:00')`,
			guestName, fmt.Sprintf("010-5678-%04d", 1000+i), 4, paymentMethodID,
			checkIn, checkOut, checkIn, checkOut,
			200000, 100000, 100000, 0, 0,
			"2박 3일 테스트 예약", 1, 0, adminUserID, adminUserID)

		if result.Error != nil {
			logrus.Errorf("Failed to create 2-night reservation %d: %v", i+1, result.Error)
			continue
		}

		// Get the last insert ID
		s.db.Raw("SELECT LAST_INSERT_ID()").Scan(&reservationID)

		// Add room to reservation
		if err := s.db.Exec(`INSERT INTO reservation_room (reservation_id, room_id, created_by, updated_by, created_at, updated_at, deleted_at)
			VALUES (?, ?, ?, ?, NOW(), NOW(), '1970-01-01 00:00:00')`,
			reservationID, availableRooms[roomIndex].ID, adminUserID, adminUserID).Error; err != nil {
			logrus.Errorf("Failed to create reservation room: %v", err)
			continue
		}

		createdReservations++
		twoDaysStay++
	}

	// 3. 현재 날짜 기준 1달 전 1일 ~ 3달 뒤 마지막날까지 달방 예약 정보 4건
	oneMonthAgo := today.AddDate(0, -1, 0)
	threeMonthsLater := today.AddDate(0, 3, 0)
	lastDayOfThirdMonth := time.Date(threeMonthsLater.Year(), threeMonthsLater.Month()+1, 0, 0, 0, 0, 0, time.Local)

	for i := 0; i < 4 && (14+i) < len(availableRooms); i++ {
		roomIndex := 14 + i

		// 월 단위로 분산
		startDate := oneMonthAgo.AddDate(0, i, 0)
		endDate := startDate.AddDate(0, 1, 0)
		if endDate.After(lastDayOfThirdMonth) {
			endDate = lastDayOfThirdMonth
		}

		guestName := fmt.Sprintf("월세고객%d", i+1)

		// Check if already exists
		var count int64
		s.db.Table("reservation").
			Where("name = ? AND stay_start_at = ? AND deleted_at = ?",
				guestName, startDate, time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)).
			Count(&count)
		if count > 0 {
			logrus.Infof("Monthly reservation for guest '%s' on %s already exists, skipping", guestName, startDate.Format("2006-01-02"))
			continue
		}

		// Insert reservation
		var reservationID uint
		result := s.db.Exec(`INSERT INTO reservation (name, phone, people_count, payment_method_id,
			stay_start_at, stay_end_at, check_in_at, check_out_at,
			price, deposit, payment_amount, refund_amount, broker_fee,
			note, status, type, created_by, updated_by, created_at, updated_at, deleted_at)
			VALUES (?, ?, ?, ?, ?, ?, NULL, NULL, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW(), '1970-01-01 00:00:00')`,
			guestName, fmt.Sprintf("010-9999-%04d", 1000+i), 2, paymentMethodID,
			startDate, endDate,
			3000000, 3000000, 3000000, 0, 0,
			"월세 예약", 1, 10, adminUserID, adminUserID)

		if result.Error != nil {
			logrus.Errorf("Failed to create monthly reservation %d: %v", i+1, result.Error)
			continue
		}

		// Get the last insert ID
		s.db.Raw("SELECT LAST_INSERT_ID()").Scan(&reservationID)

		// Add room to reservation
		if err := s.db.Exec(`INSERT INTO reservation_room (reservation_id, room_id, created_by, updated_by, created_at, updated_at, deleted_at)
			VALUES (?, ?, ?, ?, NOW(), NOW(), '1970-01-01 00:00:00')`,
			reservationID, availableRooms[roomIndex].ID, adminUserID, adminUserID).Error; err != nil {
			logrus.Errorf("Failed to create reservation room: %v", err)
			continue
		}

		createdReservations++
		monthlyRent++
	}

	result["reservations"] = map[string]interface{}{
		"created":      createdReservations,
		"oneDayStay":   oneDayStay,
		"twoDaysStay":  twoDaysStay,
		"monthlyRent":  monthlyRent,
		"totalCreated": createdReservations,
	}

	logrus.Info("=== GenerateReservationData completed ===")
	return result, nil
}

func (s *developmentService) generateReservationData_OLD() (map[string]interface{}, error) {
	result := make(map[string]interface{})

	// Get the super admin user ID
	var adminUser models.User
	if err := s.db.Where("role = ?", models.UserRoleSuperAdmin).First(&adminUser).Error; err != nil {
		return nil, fmt.Errorf("failed to find super admin user: %w", err)
	}
	adminUserID := adminUser.ID

	// Use direct DB connection instead of transaction to avoid issues
	db := s.db

	// Get payment methods
	var paymentMethods []models.PaymentMethod
	if err := db.Where("deleted_at = ?", time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)).Find(&paymentMethods).Error; err != nil || len(paymentMethods) == 0 {
		return nil, fmt.Errorf("no payment methods found, please generate essential data first")
	}

	// Get available rooms
	var availableRooms []models.Room
	if err := db.Where("status = ? AND deleted_at = ?", models.RoomStatusNormal, time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)).Order("room_group_id, number").Find(&availableRooms).Error; err != nil || len(availableRooms) == 0 {
		return nil, fmt.Errorf("no available rooms found, please generate essential data first")
	}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)

	createdReservations := make([]*models.Reservation, 0)

	// 1. 현재 날짜를 입실일로 1박 2일 예약 10건
	for i := 0; i < 10; i++ {
		if i >= len(availableRooms) {
			break
		}

		// Check if reservation already exists
		guestName := fmt.Sprintf("테스트고객%d", i+1)
		var existingCount int64
		db.Model(&models.Reservation{}).Where("name = ? AND stay_start_at = ? AND deleted_at = ?", guestName, today, time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)).Count(&existingCount)
		if existingCount > 0 {
			logrus.Infof("Reservation for guest '%s' on %s already exists, skipping", guestName, today.Format("2006-01-02"))
			continue
		}

		reservation := models.Reservation{
			BaseMustAuditEntity: models.BaseMustAuditEntity{
				CreatedBy: adminUserID,
				UpdatedBy: adminUserID,
			},
			PaymentMethodID: paymentMethods[i%len(paymentMethods)].ID,
			Name:            fmt.Sprintf("테스트고객%d", i+1),
			Phone:           fmt.Sprintf("010-1234-%04d", i+1),
			PeopleCount:     2 + (i % 3),
			StayStartAt:     today,
			StayEndAt:       today.AddDate(0, 0, 1), // 1박 2일
			Price:           100000 + (i * 10000),
			Deposit:         50000,
			PaymentAmount:   50000,
			RefundAmount:    0,
			BrokerFee:       0,
			Note:            fmt.Sprintf("1박2일 테스트 예약 %d", i+1),
			Status:          models.ReservationStatusNormal,
			Type:            models.ReservationTypeStay,
		}

		if err := db.Session(&gorm.Session{SkipHooks: true}).Omit("CreatedByUser", "UpdatedByUser", "PaymentMethod").Create(&reservation).Error; err != nil {
			return nil, fmt.Errorf("failed to create 1-night reservation %d: %w", i+1, err)
		}

		// Add room to reservation
		reservationRoom := models.ReservationRoom{
			BaseMustAuditEntity: models.BaseMustAuditEntity{
				CreatedBy: adminUserID,
				UpdatedBy: adminUserID,
			},
			ReservationID: reservation.ID,
			RoomID:        availableRooms[i].ID,
		}

		if err := db.Session(&gorm.Session{SkipHooks: true}).Omit("CreatedByUser", "UpdatedByUser", "Reservation", "Room").Create(&reservationRoom).Error; err != nil {
			return nil, fmt.Errorf("failed to create reservation room for reservation %d: %w", reservation.ID, err)
		}

		createdReservations = append(createdReservations, &reservation)
	}

	// 2. 현재 날짜 기준 2박 3일 예약 4건
	for i := 0; i < 4; i++ {
		roomIndex := 10 + i
		if roomIndex >= len(availableRooms) {
			break
		}

		// Check if reservation already exists
		guestName := fmt.Sprintf("장기고객%d", i+1)
		stayStart := today.AddDate(0, 0, i+1)
		var existingCount int64
		db.Model(&models.Reservation{}).Where("name = ? AND stay_start_at = ? AND deleted_at = ?", guestName, stayStart, time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)).Count(&existingCount)
		if existingCount > 0 {
			logrus.Infof("Reservation for guest '%s' on %s already exists, skipping", guestName, stayStart.Format("2006-01-02"))
			continue
		}

		reservation := models.Reservation{
			BaseMustAuditEntity: models.BaseMustAuditEntity{
				CreatedBy: adminUserID,
				UpdatedBy: adminUserID,
			},
			PaymentMethodID: paymentMethods[i%len(paymentMethods)].ID,
			Name:            fmt.Sprintf("장기고객%d", i+1),
			Phone:           fmt.Sprintf("010-5678-%04d", i+1),
			PeopleCount:     3 + (i % 2),
			StayStartAt:     today.AddDate(0, 0, i+1), // 내일부터 시작
			StayEndAt:       today.AddDate(0, 0, i+3), // 2박 3일
			Price:           200000 + (i * 20000),
			Deposit:         100000,
			PaymentAmount:   100000,
			RefundAmount:    0,
			BrokerFee:       0,
			Note:            fmt.Sprintf("2박3일 테스트 예약 %d", i+1),
			Status:          models.ReservationStatusNormal,
			Type:            models.ReservationTypeStay,
		}

		if err := db.Session(&gorm.Session{SkipHooks: true}).Omit("CreatedByUser", "UpdatedByUser", "PaymentMethod").Create(&reservation).Error; err != nil {
			return nil, fmt.Errorf("failed to create 2-night reservation %d: %w", i+1, err)
		}

		// Add room to reservation
		reservationRoom := models.ReservationRoom{
			BaseMustAuditEntity: models.BaseMustAuditEntity{
				CreatedBy: adminUserID,
				UpdatedBy: adminUserID,
			},
			ReservationID: reservation.ID,
			RoomID:        availableRooms[roomIndex].ID,
		}

		if err := db.Session(&gorm.Session{SkipHooks: true}).Omit("CreatedByUser", "UpdatedByUser", "Reservation", "Room").Create(&reservationRoom).Error; err != nil {
			return nil, fmt.Errorf("failed to create reservation room for reservation %d: %w", reservation.ID, err)
		}

		createdReservations = append(createdReservations, &reservation)
	}

	// 3. 현재 날짜 기준 1달 전 1일 ~ 3달 뒤 마지막날까지 달방 예약 정보 4건
	oneMonthAgo := today.AddDate(0, -1, 0)
	threeMonthsLater := today.AddDate(0, 3, 0)
	lastDayOfThirdMonth := time.Date(threeMonthsLater.Year(), threeMonthsLater.Month()+1, 0, 0, 0, 0, 0, time.Local)

	for i := 0; i < 4; i++ {
		roomIndex := 14 + i
		if roomIndex >= len(availableRooms) {
			break
		}

		// 월 단위로 분산
		startDate := oneMonthAgo.AddDate(0, i, 0)
		endDate := startDate.AddDate(0, 1, 0)
		if endDate.After(lastDayOfThirdMonth) {
			endDate = lastDayOfThirdMonth
		}

		// Check if reservation already exists
		guestName := fmt.Sprintf("월세고객%d", i+1)
		var existingCount int64
		db.Model(&models.Reservation{}).Where("name = ? AND stay_start_at = ? AND deleted_at = ?", guestName, startDate, time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)).Count(&existingCount)
		if existingCount > 0 {
			logrus.Infof("Monthly reservation for guest '%s' on %s already exists, skipping", guestName, startDate.Format("2006-01-02"))
			continue
		}

		reservation := models.Reservation{
			BaseMustAuditEntity: models.BaseMustAuditEntity{
				CreatedBy: adminUserID,
				UpdatedBy: adminUserID,
			},
			PaymentMethodID: paymentMethods[i%len(paymentMethods)].ID,
			Name:            fmt.Sprintf("월세고객%d", i+1),
			Phone:           fmt.Sprintf("010-9999-%04d", i+1),
			PeopleCount:     2,
			StayStartAt:     startDate,
			StayEndAt:       endDate,
			Price:           1500000, // 월세
			Deposit:         1500000,
			PaymentAmount:   1500000,
			RefundAmount:    0,
			BrokerFee:       0,
			Note:            fmt.Sprintf("월세 테스트 예약 %d", i+1),
			Status:          models.ReservationStatusNormal,
			Type:            models.ReservationTypeMonthlyRent,
		}

		if err := db.Session(&gorm.Session{SkipHooks: true}).Omit("CreatedByUser", "UpdatedByUser", "PaymentMethod").Create(&reservation).Error; err != nil {
			return nil, fmt.Errorf("failed to create monthly reservation %d: %w", i+1, err)
		}

		// Add room to reservation
		reservationRoom := models.ReservationRoom{
			BaseMustAuditEntity: models.BaseMustAuditEntity{
				CreatedBy: adminUserID,
				UpdatedBy: adminUserID,
			},
			ReservationID: reservation.ID,
			RoomID:        availableRooms[roomIndex].ID,
		}

		if err := db.Session(&gorm.Session{SkipHooks: true}).Omit("CreatedByUser", "UpdatedByUser", "Reservation", "Room").Create(&reservationRoom).Error; err != nil {
			return nil, fmt.Errorf("failed to create reservation room for reservation %d: %w", reservation.ID, err)
		}

		createdReservations = append(createdReservations, &reservation)
	}

	result["reservations"] = map[string]interface{}{
		"created":      len(createdReservations),
		"oneDayStay":   10,
		"twoDaysStay":  4,
		"monthlyRent":  4,
		"totalCreated": len(createdReservations),
	}

	return result, nil
}

// GenerateTestData generates test data based on the requested type
func (s *developmentService) GenerateTestData(dataType string, reservationOptions *ReservationGenerationOptions) (map[string]interface{}, error) {
	logrus.Infof("=== GenerateTestData called with type: %s ===", dataType)
	result := make(map[string]interface{})

	switch dataType {
	case "essential":
		essentialResult, err := s.generateEssentialData()
		if err != nil {
			return nil, err
		}
		result = essentialResult

	case "reset":
		if err := s.resetData(); err != nil {
			return nil, err
		}
		essentialResult, err := s.generateEssentialData()
		if err != nil {
			return nil, err
		}
		result = essentialResult

	case "reservation":
		if reservationOptions != nil {
			reservationResult, err := s.generateReservationDataWithOptions(reservationOptions)
			if err != nil {
				return nil, err
			}
			result = reservationResult
		} else {
			reservationResult, err := s.generateReservationData()
			if err != nil {
				return nil, err
			}
			result = reservationResult
		}

	case "all":
		// Generate essential data first
		essentialResult, err := s.generateEssentialData()
		if err != nil {
			return nil, err
		}
		for k, v := range essentialResult {
			result[k] = v
		}

		// Then generate reservation data
		if reservationOptions != nil {
			reservationResult, err := s.generateReservationDataWithOptions(reservationOptions)
			if err != nil {
				return nil, err
			}
			for k, v := range reservationResult {
				result[k] = v
			}
		} else {
			reservationResult, err := s.generateReservationData()
			if err != nil {
				return nil, err
			}
			for k, v := range reservationResult {
				result[k] = v
			}
		}

	default:
		return nil, fmt.Errorf("invalid data type: %s", dataType)
	}

	logrus.Info("=== GenerateTestData completed ===")
	return result, nil
}

// generateReservationDataWithOptions generates reservation data with custom options
func (s *developmentService) generateReservationDataWithOptions(options *ReservationGenerationOptions) (map[string]interface{}, error) {
	logrus.Info("=== GenerateReservationDataWithOptions called ===")
	result := make(map[string]interface{})

	// Get the super admin user ID
	var adminUser models.User
	if err := s.db.Where("role = ?", models.UserRoleSuperAdmin).First(&adminUser).Error; err != nil {
		logrus.Errorf("Failed to find super admin user: %v", err)
		return nil, fmt.Errorf("failed to find super admin user: %w", err)
	}
	adminUserID := adminUser.ID
	logrus.Infof("Found admin user with ID: %d", adminUserID)

	// Get available rooms (excluding construction and facility rooms)
	var availableRooms []struct {
		ID          uint
		Number      string
		RoomGroupID uint
	}
	if err := s.db.Table("room r").
		Select("r.id, r.number, r.room_group_id").
		Joins("JOIN room_group rg ON r.room_group_id = rg.id").
		Where("r.status = ? AND r.deleted_at = ? AND rg.name NOT IN (?) AND rg.deleted_at = ?",
			1, // Normal status
			time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
			[]string{"부대시설"},
			time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)).
		Order("r.number").
		Scan(&availableRooms).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch available rooms: %w", err)
	}

	if len(availableRooms) == 0 {
		return nil, fmt.Errorf("no available rooms found")
	}

	logrus.Infof("Found %d available rooms", len(availableRooms))

	// Get all payment methods
	var paymentMethods []struct {
		ID   uint
		Name string
	}
	if err := s.db.Table("payment_method").
		Select("id, name").
		Where("deleted_at = ?", time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)).
		Scan(&paymentMethods).Error; err != nil || len(paymentMethods) == 0 {
		return nil, fmt.Errorf("failed to find payment methods: %w", err)
	}

	// Set defaults
	startDate := time.Now().Truncate(24 * time.Hour)
	endDate := startDate.AddDate(0, 1, 0) // Default 1 month
	regularReservations := 20
	monthlyReservations := 5

	if options.StartDate != nil {
		startDate = options.StartDate.Truncate(24 * time.Hour)
	}
	if options.EndDate != nil {
		endDate = options.EndDate.Truncate(24 * time.Hour)
	}
	if options.RegularReservations != nil {
		regularReservations = *options.RegularReservations
	}
	if options.MonthlyReservations != nil {
		monthlyReservations = *options.MonthlyReservations
	}

	// Initialize counters
	createdReservations := 0
	roomAssignedCount := 0
	paymentMethodCounts := make(map[string]int)
	statusCounts := make(map[string]int)

	// Define reservation statuses with weights
	statuses := []struct {
		Status int
		Name   string
		Weight int
	}{
		{1, "정상", 70},   // 70% 확정
		{2, "예약대기", 15}, // 15% 예약대기
		{3, "예약취소", 10}, // 10% 취소
		{5, "체크아웃", 5},  // 5% 체크아웃
	}

	// Create regular reservations
	for i := 0; i < regularReservations; i++ {
		// Random date within range
		dayRange := int(endDate.Sub(startDate).Hours() / 24)
		if dayRange <= 0 {
			dayRange = 1
		}
		randomDays := rand.Intn(dayRange)
		checkIn := startDate.AddDate(0, 0, randomDays)

		// Random stay duration (1-7 days)
		stayDays := rand.Intn(7) + 1
		checkOut := checkIn.AddDate(0, 0, stayDays)

		// Random guest info
		guestName := fmt.Sprintf("테스트고객%d", i+1)
		guestPhone := fmt.Sprintf("010-%04d-%04d", rand.Intn(10000), rand.Intn(10000))
		adultCount := rand.Intn(4) + 1
		childCount := rand.Intn(3)

		// Random payment method
		paymentMethod := paymentMethods[rand.Intn(len(paymentMethods))]
		paymentMethodCounts[paymentMethod.Name]++

		// Random status (weighted)
		statusRoll := rand.Intn(100)
		var status int
		var statusName string
		cumWeight := 0
		for _, s := range statuses {
			cumWeight += s.Weight
			if statusRoll < cumWeight {
				status = s.Status
				statusName = s.Name
				break
			}
		}
		statusCounts[statusName]++

		// Random price
		basePrice := 100000 + (rand.Intn(10) * 10000)
		totalPrice := basePrice * stayDays
		deposit := totalPrice / 2
		paymentAmount := deposit
		if status == 5 { // 체크아웃인 경우 전액 결제
			paymentAmount = totalPrice
		}

		// Insert reservation
		var reservationID uint
		checkInAt := checkIn
		checkOutAt := checkOut
		if status == 2 || status == 3 { // 예약대기나 취소인 경우 체크인/아웃 시간 NULL
			result := s.db.Exec(`INSERT INTO reservation (name, phone, people_count, payment_method_id,
				stay_start_at, stay_end_at, check_in_at, check_out_at,
				price, deposit, payment_amount, refund_amount, broker_fee,
				note, status, type, created_by, updated_by, created_at, updated_at, deleted_at)
				VALUES (?, ?, ?, ?, ?, ?, NULL, NULL, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW(), '1970-01-01 00:00:00')`,
				guestName, guestPhone, adultCount+childCount, paymentMethod.ID,
				checkIn, checkOut,
				totalPrice, deposit, paymentAmount, 0, 0,
				fmt.Sprintf("테스트 예약 - %s", statusName), status, 0, adminUserID, adminUserID)
			if result.Error != nil {
				logrus.Errorf("Failed to create reservation %d: %v", i+1, result.Error)
				continue
			}
		} else {
			result := s.db.Exec(`INSERT INTO reservation (name, phone, people_count, payment_method_id,
				stay_start_at, stay_end_at, check_in_at, check_out_at,
				price, deposit, payment_amount, refund_amount, broker_fee,
				note, status, type, created_by, updated_by, created_at, updated_at, deleted_at)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW(), '1970-01-01 00:00:00')`,
				guestName, guestPhone, adultCount+childCount, paymentMethod.ID,
				checkIn, checkOut, checkInAt, checkOutAt,
				totalPrice, deposit, paymentAmount, 0, 0,
				fmt.Sprintf("테스트 예약 - %s", statusName), status, 0, adminUserID, adminUserID)
			if result.Error != nil {
				logrus.Errorf("Failed to create reservation %d: %v", i+1, result.Error)
				continue
			}
		}

		// Get the last insert ID
		s.db.Raw("SELECT LAST_INSERT_ID()").Scan(&reservationID)

		// Assign room (70% chance)
		if rand.Float32() < 0.7 && len(availableRooms) > 0 {
			room := availableRooms[rand.Intn(len(availableRooms))]
			if err := s.db.Exec(`INSERT INTO reservation_room (reservation_id, room_id, created_by, updated_by, created_at, updated_at, deleted_at)
				VALUES (?, ?, ?, ?, NOW(), NOW(), '1970-01-01 00:00:00')`,
				reservationID, room.ID, adminUserID, adminUserID).Error; err != nil {
				logrus.Errorf("Failed to create reservation room: %v", err)
			} else {
				roomAssignedCount++
			}
		}

		createdReservations++
	}

	// Create monthly reservations
	for i := 0; i < monthlyReservations; i++ {
		// Random start date within range
		dayRange := int(endDate.Sub(startDate).Hours() / 24)
		var monthStart, monthEnd time.Time

		if dayRange <= 30 {
			// If date range is too small for monthly, use the start date
			monthStart = startDate
			monthEnd = monthStart.AddDate(0, 1, 0) // 1 month later
			if monthEnd.After(endDate) {
				monthEnd = endDate
			}
		} else {
			randomDays := rand.Intn(dayRange - 30) // Ensure at least 30 days for monthly
			monthStart = startDate.AddDate(0, 0, randomDays)
			monthEnd = monthStart.AddDate(0, 1, 0) // 1 month later
		}

		guestName := fmt.Sprintf("월세고객%d", i+1)
		guestPhone := fmt.Sprintf("010-%04d-%04d", rand.Intn(10000), rand.Intn(10000))

		// Random payment method
		paymentMethod := paymentMethods[rand.Intn(len(paymentMethods))]
		paymentMethodCounts[paymentMethod.Name]++

		// Monthly reservations are usually confirmed
		status := 1
		statusCounts["정상"]++

		// Monthly price
		monthlyPrice := 1500000 + (rand.Intn(10) * 100000)

		// Insert reservation
		var reservationID uint
		result := s.db.Exec(`INSERT INTO reservation (name, phone, people_count, payment_method_id,
			stay_start_at, stay_end_at, check_in_at, check_out_at,
			price, deposit, payment_amount, refund_amount, broker_fee,
			note, status, type, created_by, updated_by, created_at, updated_at, deleted_at)
			VALUES (?, ?, ?, ?, ?, ?, NULL, NULL, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW(), '1970-01-01 00:00:00')`,
			guestName, guestPhone, 2, paymentMethod.ID,
			monthStart, monthEnd,
			monthlyPrice, monthlyPrice, monthlyPrice, 0, 0,
			"월세 예약", status, 10, adminUserID, adminUserID)

		if result.Error != nil {
			logrus.Errorf("Failed to create monthly reservation %d: %v", i+1, result.Error)
			continue
		}

		// Get the last insert ID
		s.db.Raw("SELECT LAST_INSERT_ID()").Scan(&reservationID)

		// Always assign room for monthly reservations
		if len(availableRooms) > 0 {
			room := availableRooms[rand.Intn(len(availableRooms))]
			if err := s.db.Exec(`INSERT INTO reservation_room (reservation_id, room_id, created_by, updated_by, created_at, updated_at, deleted_at)
				VALUES (?, ?, ?, ?, NOW(), NOW(), '1970-01-01 00:00:00')`,
				reservationID, room.ID, adminUserID, adminUserID).Error; err != nil {
				logrus.Errorf("Failed to create reservation room: %v", err)
			} else {
				roomAssignedCount++
			}
		}

		createdReservations++
	}

	result["reservations"] = map[string]interface{}{
		"created":             createdReservations,
		"regularReservations": regularReservations,
		"monthlyReservations": monthlyReservations,
		"roomAssigned":        roomAssignedCount,
		"paymentMethods":      paymentMethodCounts,
		"statuses":            statusCounts,
		"dateRange": map[string]string{
			"start": startDate.Format("2006-01-02"),
			"end":   endDate.Format("2006-01-02"),
		},
	}

	logrus.Info("=== GenerateReservationDataWithOptions completed ===")
	return result, nil
}
