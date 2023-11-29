import { FieldMap } from "src/types/map"
import { formatDateTime } from "src/util/format-util"
import { User } from "src/schema/user"

export type BaseTime = {
  createdAt: string;
  updatedAt: string;
};

export const BASE_TIME_FIELD_MAP: FieldMap = {
  createdAt: { label: "등록 시각", format: formatDateTime } as const,
  updatedAt: { label: "수정 시각", format: formatDateTime } as const,
};

export type BaseAudit = {
  createdBy: string; // TODO: API 응답 개선하면서 User 로 변경 필요
  updatedBy: string; // TODO API 응답 개선하면서 User 로 변경 필요
};

export const BASE_AUDIT_FIELD_MAP: FieldMap = {
  createdBy: { label: "등록자", format: (value: User) => value.name } as const,
  updatedBy: { label: "수정자", format: (value: User) => value.name } as const,
};
