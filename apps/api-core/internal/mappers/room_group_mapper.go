package mappers

import (
	"context"

	"gitlab.bellsoft.net/rms/api-core/internal/dto"
	"gitlab.bellsoft.net/rms/api-core/internal/models"
	"gitlab.bellsoft.net/rms/api-core/internal/utils"
	pkgutils "gitlab.bellsoft.net/rms/api-core/pkg/utils"
)

// ToRoomGroupResponse converts a RoomGroup model to RoomGroupResponse DTO
func ToRoomGroupResponse(roomGroup *models.RoomGroup) dto.RoomGroupResponse {
	return dto.RoomGroupResponse{
		ID:           roomGroup.ID,
		Name:         roomGroup.Name,
		PeekPrice:    roomGroup.PeekPrice,
		OffPeekPrice: roomGroup.OffPeekPrice,
		Description:  roomGroup.Description,
		Rooms:        make([]dto.RoomLastStayDetailResponse, 0),
		CreatedAt:    dto.CustomTime{Time: roomGroup.CreatedAt},
		UpdatedAt:    dto.CustomTime{Time: roomGroup.UpdatedAt},
	}
}

// ToRoomGroupResponseWithRooms converts a RoomGroup model with rooms to RoomGroupResponse DTO
func ToRoomGroupResponseWithRooms(
	ctx context.Context,
	roomGroup *models.RoomGroup,
	getUserSummary GetUserSummaryFunc,
	getLastReservation func(ctx context.Context, roomID uint) *dto.ReservationResponse,
) dto.RoomGroupResponse {
	resp := ToRoomGroupResponse(roomGroup)

	resp.Rooms = make([]dto.RoomLastStayDetailResponse, 0)

	if len(roomGroup.Rooms) > 0 {
		resp.Rooms = make([]dto.RoomLastStayDetailResponse, len(roomGroup.Rooms))
		for i, room := range roomGroup.Rooms {
			roomResponse := dto.RoomResponse{
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

			var lastReservation *dto.ReservationResponse
			if getLastReservation != nil {
				lastReservation = getLastReservation(ctx, room.ID)
			}

			resp.Rooms[i] = dto.RoomLastStayDetailResponse{
				Room:            roomResponse,
				LastReservation: lastReservation,
			}
		}
	}

	return resp
}

// ToRoomGroupResponseWithUsers converts a RoomGroup model with user info to RoomGroupResponse DTO
func ToRoomGroupResponseWithUsers(roomGroup *models.RoomGroup) dto.RoomGroupResponse {
	resp := dto.RoomGroupResponse{
		ID:           roomGroup.ID,
		Name:         roomGroup.Name,
		PeekPrice:    roomGroup.PeekPrice,
		OffPeekPrice: roomGroup.OffPeekPrice,
		Description:  roomGroup.Description,
		Rooms:        make([]dto.RoomLastStayDetailResponse, 0),
		CreatedAt:    dto.CustomTime{Time: roomGroup.CreatedAt},
		UpdatedAt:    dto.CustomTime{Time: roomGroup.UpdatedAt},
	}

	if roomGroup.CreatedByUser != nil {
		resp.CreatedBy = &dto.UserSummaryResponse{
			ID:              roomGroup.CreatedByUser.ID,
			Email:           utils.StringPtrToString(roomGroup.CreatedByUser.Email),
			Name:            roomGroup.CreatedByUser.Name,
			ProfileImageURL: pkgutils.GenerateGravatarURL(utils.StringPtrToString(roomGroup.CreatedByUser.Email)),
		}
	}

	if roomGroup.UpdatedByUser != nil {
		resp.UpdatedBy = &dto.UserSummaryResponse{
			ID:              roomGroup.UpdatedByUser.ID,
			Email:           utils.StringPtrToString(roomGroup.UpdatedByUser.Email),
			Name:            roomGroup.UpdatedByUser.Name,
			ProfileImageURL: pkgutils.GenerateGravatarURL(utils.StringPtrToString(roomGroup.UpdatedByUser.Email)),
		}
	}

	return resp
}
