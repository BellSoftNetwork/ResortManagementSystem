import { api } from "boot/axios";
import { ListResponse, PageRequestParams, SortRequestParams } from "src/schema/response";
import { HistoryType } from "src/schema/revision";

export interface AuditLog {
  id: number;
  entityType: string;
  entityId: number;
  action: HistoryType;
  changedFields: string[];
  userId: number;
  username: string;
  createdAt: string;
}

export interface AuditLogDetail extends AuditLog {
  oldValues: Record<string, unknown> | null;
  newValues: Record<string, unknown> | null;
}

export type FetchAuditLogsRequestParams = Partial<
  {
    entityType: string;
    startDate: string;
    endDate: string;
    action: HistoryType;
    userId: number;
    entityId: number;
  } & PageRequestParams &
    SortRequestParams
>;

export async function fetchAuditLogs(params: FetchAuditLogsRequestParams) {
  const result = await api.get<ListResponse<AuditLog>>("/api/v1/admin/audit-logs", {
    params,
  });
  return result.data;
}

export async function fetchAuditLogDetail(id: number) {
  const result = await api.get<{ value: AuditLogDetail }>(`/api/v1/admin/audit-logs/${id}`);
  return result.data.value;
}
