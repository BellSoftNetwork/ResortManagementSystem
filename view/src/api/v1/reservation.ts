import { api } from "boot/axios";
import {
  ListResponse,
  PageRequestParams,
  SingleResponse,
  SortRequestParams,
} from "src/schema/response";
import { Reservation, ReservationStatus } from "src/schema/reservation";
import { Revision } from "src/schema/revision";
import { EntityReference } from "src/schema/base";

export type FetchReservationsRequestParams = Partial<
  {
    stayStartAt: string;
    stayEndAt: string;
    searchText: string;
    status: ReservationStatus;
  } & PageRequestParams &
    SortRequestParams
>;

export async function fetchReservations(
  params: FetchReservationsRequestParams,
) {
  const result = await api.get<ListResponse<Reservation>>(
    "/api/v1/reservations",
    {
      params,
    },
  );
  return result.data;
}

export async function fetchReservation(id: number) {
  const result = await api.get<SingleResponse<Reservation>>(
    `/api/v1/reservations/${id}`,
  );
  return result.data;
}

type ReservationParams = {
  paymentMethodId: number;
  rooms: EntityReference[];
  name: string;
  phone: string;
  peopleCount: number;
  stayStartAt: string;
  stayEndAt: string;
  price: number;
  paymentAmount: number;
  brokerFee: number;
  note: string;
  status: ReservationStatus;
};

export type ReservationCreateParams = ReservationParams;

export async function createReservation(params: ReservationCreateParams) {
  const result = await api.post<SingleResponse<Reservation>>(
    "/api/v1/reservations",
    params,
  );
  return result.data;
}

export type ReservationPatchParams = Partial<ReservationParams>;

export async function patchReservation(
  id: number,
  params: Partial<ReservationPatchParams>,
) {
  const result = await api.patch<SingleResponse<Reservation>>(
    `/api/v1/reservations/${id}`,
    params,
  );
  return result.data;
}

export async function deleteReservation(id: number) {
  const result = await api.delete(`/api/v1/reservations/${id}`);
  return result.data;
}

type FetchReservationHistoriesRequestParams = Partial<
  PageRequestParams & SortRequestParams
>;

export async function fetchReservationHistories(
  id: number,
  params: FetchReservationHistoriesRequestParams,
) {
  const result = await api.get<ListResponse<Revision<Reservation>>>(
    `/api/v1/reservations/${id}/histories`,
    {
      params,
    },
  );
  return result.data;
}
