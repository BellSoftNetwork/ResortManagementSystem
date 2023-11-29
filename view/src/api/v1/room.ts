import { api } from "boot/axios"
import {
  ListResponse,
  PageRequestParams,
  SingleResponse,
  SortRequestParams,
} from "src/schema/response";
import { Room, RoomStatus } from "src/schema/room"
import { Revision } from "src/schema/revision"

type FetchRoomsRequestParams = Partial<
  {
    stayStartAt: string;
    stayEndAt: string;
    status: RoomStatus;
  } & PageRequestParams &
  SortRequestParams
>;

export async function fetchRooms(params: FetchRoomsRequestParams) {
  const result = await api.get<ListResponse<Room>>("/api/v1/rooms", {
    params,
  });
  return result.data
}

export async function fetchRoom(id: number) {
  const result = await api.get<SingleResponse<Room>>(`/api/v1/rooms/${id}`)
  return result.data
}

type RoomParams = {
  number: string;
  peekPrice: number;
  offPeekPrice: number;
  description: string;
  note: string;
  status: RoomStatus;
};

export async function createRoom(params: RoomParams) {
  const result = await api.post<SingleResponse<Room>>("/api/v1/rooms", params)
  return result.data
}

export async function patchRoom(id: number, params: Partial<RoomParams>) {
  const result = await api.patch<SingleResponse<Room>>(
    `/api/v1/rooms/${id}`,
    params,
  );
  return result.data
}

export async function deleteRoom(id: number) {
  const result = await api.delete(`/api/v1/rooms/${id}`)
  return result.data
}

type FetchRoomHistoriesRequestParams = Partial<
  PageRequestParams & SortRequestParams
>;

export async function fetchRoomHistories(
  id: number,
  params: FetchRoomHistoriesRequestParams,
) {
  const result = await api.get<ListResponse<Revision<Room>>>(
    `/api/v1/rooms/${id}/histories`,
    {
      params,
    },
  );
  return result.data
}
