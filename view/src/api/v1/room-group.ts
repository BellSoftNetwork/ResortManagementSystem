import { api } from "boot/axios";
import { ListResponse, PageRequestParams, SingleResponse, SortRequestParams } from "src/schema/response";
import { Room, RoomStatus } from "src/schema/room";
import { RoomGroup } from "src/schema/room-group";
import { Reservation } from "src/schema/reservation";

type RoomGroupSummary = Omit<RoomGroup, "rooms">;

export async function fetchRoomGroups() {
  const result = await api.get<ListResponse<RoomGroupSummary>>("/api/v1/room-groups");
  return result.data;
}

type FetchRoomFilterParams = Partial<
  {
    stayStartAt: string;
    stayEndAt: string;
    status: RoomStatus;
    excludeReservationId: number;
  } & PageRequestParams &
    SortRequestParams
>;

export type RoomLastStayDetail = {
  room: Room;
  lastReservation: Reservation;
};
export type RoomGroupDetailResponse = RoomGroupSummary & {
  rooms: RoomLastStayDetail[];
};

export async function fetchRoomGroup(id: number, params: FetchRoomFilterParams) {
  const result = await api.get<SingleResponse<RoomGroupDetailResponse>>(`/api/v1/room-groups/${id}`, { params });
  return result.data;
}

type RoomGroupParams = {
  name: string;
  peekPrice: number;
  offPeekPrice: number;
  description: string;
};

export async function createRoomGroup(params: RoomGroupParams) {
  const result = await api.post<SingleResponse<RoomGroupSummary>>("/api/v1/room-groups", params);
  return result.data;
}

export async function patchRoomGroup(id: number, params: Partial<RoomGroupParams>) {
  const result = await api.patch<SingleResponse<RoomGroupSummary>>(`/api/v1/room-groups/${id}`, params);
  return result.data;
}

export async function deleteRoomGroup(id: number) {
  const result = await api.delete(`/api/v1/room-groups/${id}`);
  return result.data;
}
