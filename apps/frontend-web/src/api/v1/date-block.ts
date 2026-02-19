import { api } from "boot/axios";
import { ListResponse, SingleResponse, PageRequestParams, SortRequestParams } from "src/schema/response";
import { DateBlock, CreateDateBlockRequest, UpdateDateBlockRequest } from "src/schema/date-block";
import { Revision } from "src/schema/revision";

type FetchDateBlocksRequestParams = Partial<
  { startDate?: string; endDate?: string } & PageRequestParams & SortRequestParams
>;

export async function fetchDateBlocks(params: FetchDateBlocksRequestParams) {
  const result = await api.get<ListResponse<DateBlock>>("/api/v1/date-blocks", {
    params,
  });
  return result.data;
}

export async function fetchDateBlock(id: number) {
  const result = await api.get<SingleResponse<DateBlock>>(`/api/v1/date-blocks/${id}`);
  return result.data;
}

export async function createDateBlock(data: CreateDateBlockRequest) {
  const result = await api.post<SingleResponse<DateBlock>>("/api/v1/date-blocks", data);
  return result.data;
}

export async function updateDateBlock(id: number, params: UpdateDateBlockRequest) {
  const result = await api.patch<SingleResponse<DateBlock>>(`/api/v1/date-blocks/${id}`, params);
  return result.data;
}

export async function deleteDateBlock(id: number) {
  const result = await api.delete(`/api/v1/date-blocks/${id}`);
  return result.data;
}

type FetchDateBlockHistoriesRequestParams = Partial<PageRequestParams & SortRequestParams>;

export async function fetchDateBlockHistories(id: number, params: FetchDateBlockHistoriesRequestParams) {
  const result = await api.get<ListResponse<Revision<DateBlock>>>(`/api/v1/date-blocks/${id}/histories`, {
    params,
  });
  return result.data;
}
