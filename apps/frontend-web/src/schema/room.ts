import { EnumMap, FieldDetail, FieldMap, StaticRuleMap } from "src/types/map";
import { BASE_AUDIT_FIELD_MAP, BASE_TIME_FIELD_MAP, BaseAudit, BaseTime } from "src/schema/base";
import { RoomGroup } from "src/schema/room-group";

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
  roomGroup: RoomGroup;
  note: string;
  status: RoomStatus;
} & BaseTime &
  BaseAudit;

const RoomFieldMap: FieldMap = {
  id: { label: "ID" } as const,
  number: { label: "객실 번호" } as const,
  roomGroup: { label: "객실 그룹", format: (roomGroup: RoomGroup) => roomGroup.name } as const,
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
  note: [(value: string) => (value.length >= 0 && value.length <= 200) || "200 글자까지 입력 가능합니다"] as const,
} as const;
