import { formatPrice } from "src/util/format-util";
import { EnumMap, FieldDetail, FieldMap, StaticRuleMap } from "src/types/map";
import { BASE_AUDIT_FIELD_MAP, BASE_TIME_FIELD_MAP, BaseAudit, BaseTime } from "src/schema/base";

const ROOM_STATUS_MAP: EnumMap = {
  NORMAL: "정상",
  INACTIVE: "이용 불가",
  CONSTRUCTION: "공사 중",
  DAMAGED: "파손",
} as const;
export type RoomStatus = keyof typeof ROOM_STATUS_MAP;

export function roomStatusValueToName(role: RoomStatus) {
  return ROOM_STATUS_MAP[role] || role;
}

export type Room = {
  id: number;
  number: string;
  peekPrice: number;
  offPeekPrice: number;
  description: string;
  note: string;
  status: RoomStatus;
} & BaseTime &
  BaseAudit;

const RoomFieldMap: FieldMap = {
  id: { label: "ID" } as const,
  number: { label: "객실 번호" } as const,
  peekPrice: { label: "성수기 예약금", format: formatPrice } as const,
  offPeekPrice: { label: "비성수기 예약금", format: formatPrice } as const,
  description: { label: "설명" } as const,
  note: { label: "메모" } as const,
  status: { label: "상태", format: roomStatusValueToName } as const,
  ...BASE_TIME_FIELD_MAP,
  ...BASE_AUDIT_FIELD_MAP,
} as const;

export function getRoomFieldDetail(field: string): FieldDetail | null {
  const fieldDetail = RoomFieldMap[field];

  return fieldDetail
    ? {
        field: field,
        ...fieldDetail,
      }
    : null;
}

export function formatRoomFieldToLabel(field: string) {
  return getRoomFieldDetail(field)?.label ?? field;
}

export function formatRoomValue(field: string, value: string | number | null) {
  const fieldDetail = getRoomFieldDetail(field);
  return fieldDetail?.format ? fieldDetail.format(value) : value;
}

export const roomStaticRules: StaticRuleMap = {
  number: [(value: string) => (value.length >= 2 && value.length <= 20) || "2~20 글자가 필요합니다"] as const,
  peekPrice: [(value: number) => (value >= 0 && value <= 100000000) || "금액은 1억 미만 양수만 가능합니다"] as const,
  offPeekPrice: [(value: number) => (value >= 0 && value <= 100000000) || "금액은 1억 미만 양수만 가능합니다"] as const,
  description: [
    (value: string) => (value.length >= 0 && value.length <= 200) || "200 글자까지 입력 가능합니다",
  ] as const,
  note: [(value: string) => (value.length >= 0 && value.length <= 200) || "200 글자까지 입력 가능합니다"] as const,
} as const;
