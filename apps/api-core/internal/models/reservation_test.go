package models_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.bellsoft.net/rms/api-core/internal/models"
)

func TestGetAuditFields_날짜는_YYYY_MM_DD_형식이다(t *testing.T) {
	t.Run("stayStartAt과 stayEndAt은 YYYY-MM-DD 형식", func(t *testing.T) {
		reservation := &models.Reservation{
			StayStartAt: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			StayEndAt:   time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
		}
		fields := reservation.GetAuditFields()
		assert.Equal(t, "2024-01-15", fields["stayStartAt"])
		assert.Equal(t, "2024-01-20", fields["stayEndAt"])
	})

	t.Run("checkInAt이 nil이면 nil 반환", func(t *testing.T) {
		reservation := &models.Reservation{}
		fields := reservation.GetAuditFields()
		assert.Nil(t, fields["checkInAt"])
	})

	t.Run("checkInAt이 있으면 RFC3339 문자열 반환", func(t *testing.T) {
		checkIn := time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC)
		reservation := &models.Reservation{CheckInAt: &checkIn}
		fields := reservation.GetAuditFields()
		assert.Equal(t, "2024-01-15T14:30:00Z", fields["checkInAt"])
	})
}

func TestReservation_GetAuditFields_rooms가_포함된다(t *testing.T) {
	// Given - 객실이 포함된 예약이 있을 때
	reservation := &models.Reservation{
		Rooms: []models.ReservationRoom{
			{RoomID: 1, Room: &models.Room{Number: "101호"}},
			{RoomID: 2, Room: &models.Room{Number: "102호"}},
		},
	}
	reservation.ID = 1

	// When - GetAuditFields를 호출하면
	fields := reservation.GetAuditFields()

	// Then - rooms 필드가 [{id, number}] 형태로 포함되어야 한다
	rooms, ok := fields["rooms"]
	assert.True(t, ok, "audit fields에 rooms가 포함되어야 함")

	roomsSlice, ok := rooms.([]map[string]interface{})
	assert.True(t, ok, "rooms는 []map[string]interface{} 타입이어야 함")
	assert.Len(t, roomsSlice, 2)

	// ID 기준 정렬되어 있어야 함
	assert.Equal(t, uint(1), roomsSlice[0]["id"])
	assert.Equal(t, "101호", roomsSlice[0]["number"])
	assert.Equal(t, uint(2), roomsSlice[1]["id"])
	assert.Equal(t, "102호", roomsSlice[1]["number"])
}

func TestReservation_GetAuditFields_paymentMethod가_포함된다(t *testing.T) {
	// Given - 결제수단이 포함된 예약이 있을 때
	reservation := &models.Reservation{
		PaymentMethodID: 3,
		PaymentMethod:   &models.PaymentMethod{Name: "펜션다나와"},
	}
	reservation.ID = 1

	// When - GetAuditFields를 호출하면
	fields := reservation.GetAuditFields()

	// Then - paymentMethod 필드가 {id, name} 형태로 포함되어야 한다
	paymentMethod, ok := fields["paymentMethod"]
	assert.True(t, ok, "audit fields에 paymentMethod가 포함되어야 함")

	pmMap, ok := paymentMethod.(map[string]interface{})
	assert.True(t, ok, "paymentMethod는 map[string]interface{} 타입이어야 함")
	assert.Equal(t, uint(3), pmMap["id"])
	assert.Equal(t, "펜션다나와", pmMap["name"])
}

func TestReservation_GetAuditFields_Room이_nil이면_number가_빈문자열(t *testing.T) {
	// Given - Room이 preload되지 않은 예약이 있을 때
	reservation := &models.Reservation{
		Rooms: []models.ReservationRoom{
			{RoomID: 1, Room: nil},
		},
	}
	reservation.ID = 1

	// When - GetAuditFields를 호출하면
	fields := reservation.GetAuditFields()

	// Then - rooms의 number는 빈 문자열이어야 한다
	rooms := fields["rooms"].([]map[string]interface{})
	assert.Equal(t, "", rooms[0]["number"])
}
