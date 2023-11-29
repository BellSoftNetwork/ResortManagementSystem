import { api } from "boot/axios"
import {
  ListResponse,
  PageRequestParams,
  SingleResponse,
  SortRequestParams,
} from "src/schema/response";
import { ReservationMethod } from "src/schema/reservation-method"

type FetchReservationMethodsRequestParams = Partial<
  PageRequestParams & SortRequestParams
>;

export async function fetchReservationMethods(
  params: FetchReservationMethodsRequestParams,
) {
  const result = await api.get<ListResponse<ReservationMethod>>(
    "/api/v1/reservation-methods",
    {
      params,
    },
  );
  return result.data
}

export async function fetchReservationMethod(id: number) {
  const result = await api.get<SingleResponse<ReservationMethod>>(
    `/api/v1/reservation-methods/${id}`,
  );
  return result.data
}

type ReservationMethodParams = {
  name: string;
  commissionRate: number;
  requireUnpaidAmountCheck: boolean;
};

export async function createReservationMethod(params: ReservationMethodParams) {
  const result = await api.post<SingleResponse<ReservationMethod>>(
    "/api/v1/reservation-methods",
    params,
  );
  return result.data
}

export async function patchReservationMethod(
  id: number,
  params: Partial<ReservationMethodParams>,
) {
  const result = await api.patch<SingleResponse<ReservationMethod>>(
    `/api/v1/reservation-methods/${id}`,
    params,
  );
  return result.data
}

export async function deleteReservationMethod(id: number) {
  const result = await api.delete(`/api/v1/reservation-methods/${id}`)
  return result.data
}
