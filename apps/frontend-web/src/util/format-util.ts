import dayjs from "dayjs";

export function formatPrice(value: number) {
  return new Intl.NumberFormat("ko-KR", {
    style: "currency",
    currency: "KRW",
  }).format(value);
}

export function formatCommissionRate(value: number) {
  return value * 100 + "%";
}

export function formatDateTime(value: string | undefined = undefined) {
  return dayjs(value).format("YYYY-MM-DD HH:mm:ss");
}

export function formatDate(value: string | undefined = undefined) {
  return dayjs(value).format("YYYY-MM-DD");
}

export function formatSimpleDate(value: string | undefined = undefined) {
  return dayjs(value).format("YY/MM/DD");
}

export function formatTime(value: string | undefined = undefined) {
  return dayjs(value).format("HH:mm:ss");
}

export function formatDiffDays(startDate: string | null, endDate: string | null): number {
  if (startDate === null || endDate === null) return 0;

  try {
    return dayjs(endDate).diff(dayjs(startDate), "day");
  } catch {
    return 0;
  }
}

export function formatStayTitle(dateDiff: number) {
  return `${dateDiff}박 ${dateDiff + 1}일`;
}

export function formatStayCaption(startDate: string | null, endDate: string | null) {
  if (startDate === null || endDate === null) return null;

  const diffDays = formatDiffDays(startDate, endDate);
  return formatDiffDays(startDate, endDate) > 0 ? `${startDate} ~ ${endDate} (${formatStayTitle(diffDays)})` : null;
}
