import { api } from "boot/axios";
import { ListResponse, PageRequestParams, SingleResponse, SortRequestParams } from "src/schema/response";
import { PaymentMethod } from "src/schema/payment-method";

type FetchPaymentMethodsRequestParams = Partial<PageRequestParams & SortRequestParams>;

export async function fetchPaymentMethods(params: FetchPaymentMethodsRequestParams) {
  const result = await api.get<ListResponse<PaymentMethod>>("/api/v1/payment-methods", {
    params,
  });
  return result.data;
}

export async function fetchPaymentMethod(id: number) {
  const result = await api.get<SingleResponse<PaymentMethod>>(`/api/v1/payment-methods/${id}`);
  return result.data;
}

type PaymentMethodParams = {
  name: string;
  commissionRate: number;
  requireUnpaidAmountCheck: boolean;
};

export async function createPaymentMethod(params: PaymentMethodParams) {
  const result = await api.post<SingleResponse<PaymentMethod>>("/api/v1/payment-methods", params);
  return result.data;
}

export async function patchPaymentMethod(id: number, params: Partial<PaymentMethodParams>) {
  const result = await api.patch<SingleResponse<PaymentMethod>>(`/api/v1/payment-methods/${id}`, params);
  return result.data;
}

export async function deletePaymentMethod(id: number) {
  const result = await api.delete(`/api/v1/payment-methods/${id}`);
  return result.data;
}
