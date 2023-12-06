import { Room } from "src/schema/room";
import { ReservationMethod } from "src/schema/reservation-method";
import {
  DynamicRuleMap,
  EnumMap,
  FieldDetail,
  FieldMap,
  StaticRuleMap,
} from "src/types/map";
import { formatDate, formatDateTime, formatPrice } from "src/util/format-util";
import {
  BASE_AUDIT_FIELD_MAP,
  BASE_TIME_FIELD_MAP,
  BaseAudit,
  BaseTime,
} from "src/schema/base";

const RESERVATION_STATUS_MAP: EnumMap = {
  NORMAL: "예약 확정",
  PENDING: "예약 대기",
  CANCEL: "취소 요청",
  REFUND: "환불 완료",
} as const;
export type ReservationStatus = keyof typeof RESERVATION_STATUS_MAP;

export function reservationStatusValueToName(role: ReservationStatus) {
  return RESERVATION_STATUS_MAP[role] || role;
}

export type Reservation = {
  id: number;
  reservationMethod: ReservationMethod;
  rooms: Room[];
  name: string;
  phone: string;
  peopleCount: number;
  stayStartAt: string;
  stayEndAt: string;
  checkInAt: string;
  checkOutAt: string;
  price: number;
  paymentAmount: number;
  refundAmount: number;
  brokerFee: number;
  note: string;
  canceledAt: string;
  status: ReservationStatus;
} & BaseTime &
  BaseAudit;

const ReservationFieldMap: FieldMap = {
  id: { label: "ID" } as const,
  reservationMethod: {
    label: "예약 수단",
    format: (value) => value.name,
  } as const,
  // NOTE: 히스토리 조회 시 하위 호환을 위해 유지 필요
  room: {
    label: "객실",
    format: (value: Room) => (value ? value.number : "미배정"),
  } as const,
  rooms: {
    label: "객실",
    format: (value: Room[]) =>
      value.length !== 0
        ? value.map((room) => room.number).join(", ")
        : "미배정",
  } as const,
  name: { label: "예약자명" } as const,
  phone: {
    label: "예약자 연락처",
    format: (value: string) => value || "-",
  } as const,
  peopleCount: { label: "예약인원" } as const,
  stayStartAt: { label: "입실일", format: formatDate } as const,
  stayEndAt: { label: "퇴실일", format: formatDate } as const,
  checkInAt: {
    label: "체크인 시각",
    format: (value: string | null) => (value ? formatDateTime(value) : ""),
  } as const,
  checkOutAt: {
    label: "체크아웃 시각",
    format: (value: string | null) => (value ? formatDateTime(value) : ""),
  } as const,
  price: { label: "판매 금액", format: formatPrice } as const,
  paymentAmount: { label: "누적 결제 금액", format: formatPrice } as const,
  refundAmount: { label: "환불 금액", format: formatPrice } as const,
  brokerFee: { label: "예약 수단 수수료", format: formatPrice } as const,
  note: { label: "메모" } as const,
  canceledAt: { label: "취소 시각", format: formatDateTime } as const,
  status: { label: "상태", format: reservationStatusValueToName } as const,
  ...BASE_TIME_FIELD_MAP,
  ...BASE_AUDIT_FIELD_MAP,
} as const;

export function getReservationFieldDetail(field: string): FieldDetail | null {
  const fieldDetail = ReservationFieldMap[field];

  return fieldDetail
    ? {
        field: field,
        ...fieldDetail,
      }
    : null;
}

export function formatReservationFieldToLabel(field: string) {
  return getReservationFieldDetail(field)?.label ?? field;
}

export function formatReservationValue(
  field: string,
  value: string | number | null,
) {
  const fieldDetail = getReservationFieldDetail(field);
  return fieldDetail?.format ? fieldDetail.format(value) : value;
}

export const reservationStaticRules: StaticRuleMap = {
  name: [
    (value: string) =>
      (value.length >= 2 && value.length <= 30) || "2~30 글자가 필요합니다",
  ] as const,
  phone: [
    (value: string) => value.length <= 20 || "20 글자 이내로 입력 가능합니다",
  ] as const,
  peopleCount: [
    (value: number) =>
      (value >= 0 && value <= 1000) || "1000 명 이하만 입실 가능합니다",
  ] as const,
  stayStartAt: [
    (value: string) =>
      /^-?[\d]+-[0-1]\d-[0-3]\d$/.test(value) ||
      "####-##-## 형태의 날짜만 입력 가능합니다.",
  ] as const,
  stayEndAt: [
    (value: string) =>
      /^-?[\d]+-[0-1]\d-[0-3]\d$/.test(value) ||
      "####-##-## 형태의 날짜만 입력 가능합니다.",
  ] as const,
  price: [
    (value: number) =>
      (value >= 0 && value <= 100000000) || "금액은 1억 미만 양수만 가능합니다",
  ] as const,
  brokerFee: [
    (value: number) =>
      (value >= 0 && value <= 100000000) || "금액은 1억 미만 양수만 가능합니다",
  ] as const,
  note: [
    (value: string) =>
      (value.length >= 0 && value.length <= 200) ||
      "200 글자까지 입력 가능합니다",
  ] as const,
} as const;

export const reservationDynamicRules: DynamicRuleMap = {
  paymentAmount: (price: number) =>
    [
      (value: number) =>
        (value >= 0 && value <= 100000000) ||
        "금액은 1억 미만 양수만 가능합니다",
      (value: number) => value <= price || "판매 금액보다 클 수 없습니다",
    ] as const,
} as const;
