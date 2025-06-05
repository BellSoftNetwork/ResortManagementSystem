import { FieldDetail, FieldMap, StaticRuleMap } from "src/types/map";
import { BASE_AUDIT_FIELD_MAP, BASE_TIME_FIELD_MAP, BaseAudit, BaseTime } from "src/schema/base";
import { formatPrice } from "src/util/format-util";

export type RoomGroup = {
  id: number;
  name: string;
  peekPrice: number;
  offPeekPrice: number;
  description: string;
} & BaseTime &
  BaseAudit;

const RoomGroupFieldMap: FieldMap = {
  id: { label: "ID" } as const,
  name: { label: "객실 그룹명" } as const,
  peekPrice: { label: "성수기 예약금", format: formatPrice } as const,
  offPeekPrice: { label: "비성수기 예약금", format: formatPrice } as const,
  description: { label: "설명" } as const,
  ...BASE_TIME_FIELD_MAP,
  ...BASE_AUDIT_FIELD_MAP,
} as const;

export function getRoomGroupFieldDetail(field: string): FieldDetail | null {
  const fieldDetail = RoomGroupFieldMap[field];

  return fieldDetail
    ? {
        field: field,
        ...fieldDetail,
      }
    : null;
}

export function formatRoomGroupFieldToLabel(field: string) {
  return getRoomGroupFieldDetail(field)?.label ?? field;
}

export function formatRoomGroupValue(field: string, value: string | number | null) {
  const fieldDetail = getRoomGroupFieldDetail(field);
  return fieldDetail?.format ? fieldDetail.format(value) : value;
}

export const roomGroupStaticRules: StaticRuleMap = {
  name: [(value: string) => (value.length >= 2 && value.length <= 20) || "2~20 글자가 필요합니다"] as const,
  peekPrice: [(value: number) => (value >= 0 && value <= 100000000) || "금액은 1억 미만 양수만 가능합니다"] as const,
  offPeekPrice: [(value: number) => (value >= 0 && value <= 100000000) || "금액은 1억 미만 양수만 가능합니다"] as const,
  description: [
    (value: string) => (value.length >= 0 && value.length <= 200) || "200 글자까지 입력 가능합니다",
  ] as const,
} as const;
