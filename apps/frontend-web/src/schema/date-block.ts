import { UserSummary } from "src/schema/user";
import { FieldDetail, FieldMap } from "src/types/map";
import { BASE_AUDIT_FIELD_MAP, BASE_TIME_FIELD_MAP } from "src/schema/base";

export interface DateBlock {
  id: number;
  startDate: string;
  endDate: string;
  reason: string;
  createdBy: UserSummary;
  createdAt: string;
}

export interface CreateDateBlockRequest {
  startDate: string;
  endDate: string;
  reason: string;
}

export interface UpdateDateBlockRequest {
  startDate?: string;
  endDate?: string;
  reason?: string;
}

export interface DateBlockFilter {
  startDate?: string;
  endDate?: string;
  page?: number;
  size?: number;
  sort?: string;
}

const DateBlockFieldMap: FieldMap = {
  id: { label: "ID" } as const,
  startDate: { label: "시작일" } as const,
  endDate: { label: "종료일" } as const,
  reason: { label: "사유" } as const,
  ...BASE_TIME_FIELD_MAP,
  ...BASE_AUDIT_FIELD_MAP,
} as const;

export function getDateBlockFieldDetail(field: string): FieldDetail | null {
  const fieldDetail = DateBlockFieldMap[field];

  return fieldDetail
    ? {
        field: field,
        ...fieldDetail,
      }
    : null;
}

export function formatDateBlockFieldToLabel(field: string): string {
  return getDateBlockFieldDetail(field)?.label ?? field;
}

export function formatDateBlockValue(field: string, value: unknown): string {
  const fieldDetail = getDateBlockFieldDetail(field);
  return fieldDetail?.format ? fieldDetail.format(value) : String(value ?? "");
}
