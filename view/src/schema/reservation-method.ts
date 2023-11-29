import { FieldDetail, FieldMap, StaticRuleMap } from "src/types/map";
import { formatCommissionRate } from "src/util/format-util";
import { BASE_TIME_FIELD_MAP, BaseTime } from "src/schema/base";

export type ReservationMethod = {
  id: number;
  name: string;
  commissionRate: number;
  requireUnpaidAmountCheck: boolean;
} & BaseTime;

const ReservationMethodFieldMap: FieldMap = {
  id: { label: "ID" } as const,
  name: { label: "이름" } as const,
  commissionRate: { label: "수수료율", format: formatCommissionRate } as const,
  requireUnpaidAmountCheck: {
    label: "미수금 금액 알림",
    format: (value: boolean) => (value ? "활성" : "비활성"),
  } as const,
  ...BASE_TIME_FIELD_MAP,
} as const;

export function getReservationMethodFieldDetail(
  field: string,
): FieldDetail | null {
  const fieldDetail = ReservationMethodFieldMap[field];

  return fieldDetail
    ? {
        field: field,
        ...fieldDetail,
      }
    : null;
}

export function formatReservationMethodFieldToLabel(field: string) {
  return getReservationMethodFieldDetail(field)?.label ?? field;
}

export function formatReservationMethodValue(
  field: string,
  value: string | number | null,
) {
  const fieldDetail = getReservationMethodFieldDetail(field);
  return fieldDetail?.format ? fieldDetail.format(value) : value;
}

export const reservationMethodStaticRules: StaticRuleMap = {
  name: [
    (value: string) =>
      (value.length >= 2 && value.length <= 20) || "2~20 글자가 필요합니다",
  ] as const,
  commissionRatePercent: [
    (value: number) =>
      (value >= 0 && value <= 100) || "수수료율이 유효하지 않습니다.",
  ] as const,
} as const;
