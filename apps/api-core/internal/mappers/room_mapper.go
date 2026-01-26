package mappers

import (
	"context"

	"gitlab.bellsoft.net/rms/api-core/internal/dto"
	"gitlab.bellsoft.net/rms/api-core/internal/models"
)

// GetUserSummaryFunc is a function type for retrieving user summary
type GetUserSummaryFunc func(ctx context.Context, userID uint) *dto.UserSummaryResponse

// ToRoomResponse converts a Room model to RoomResponse DTO
func ToRoomResponse(ctx context.Context, room *models.Room, getUserSummary GetUserSummaryFunc) dto.RoomResponse {
	resp := dto.RoomResponse{
		ID:          room.ID,
		Number:      room.Number,
		RoomGroupID: room.RoomGroupID,
		Note:        room.Note,
		Status:      room.Status.String(),
		CreatedAt:   dto.CustomTime{Time: room.CreatedAt},
		UpdatedAt:   dto.CustomTime{Time: room.UpdatedAt},
		CreatedBy:   getUserSummary(ctx, room.CreatedBy),
		UpdatedBy:   getUserSummary(ctx, room.UpdatedBy),
	}

	if room.RoomGroup != nil {
		resp.RoomGroup = &dto.RoomGroupResponse{
			ID:           room.RoomGroup.ID,
			Name:         room.RoomGroup.Name,
			PeekPrice:    room.RoomGroup.PeekPrice,
			OffPeekPrice: room.RoomGroup.OffPeekPrice,
			Description:  room.RoomGroup.Description,
			CreatedAt:    dto.CustomTime{Time: room.RoomGroup.CreatedAt},
			UpdatedAt:    dto.CustomTime{Time: room.RoomGroup.UpdatedAt},
		}
	}

	return resp
}
