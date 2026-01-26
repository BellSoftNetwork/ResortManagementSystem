import dayjs from "dayjs";
import { formatDate } from "./format-util";

export type DueOption = "ALL" | "1M" | "2M" | "3M" | "6M" | "CUSTOM";

export interface DateRange {
  startAt: string;
  endAt: string;
}

/**
 * Calculate date range based on due option preset
 * @param dueOption - The preset option
 * @param startDate - Optional start date (defaults to today)
 * @returns DateRange with startAt and endAt
 */
export function calculateDateRange(dueOption: DueOption, startDate?: string): DateRange {
  if (dueOption === "ALL") {
    return { startAt: "", endAt: "" };
  }

  const start = startDate || formatDate();

  if (dueOption === "CUSTOM") {
    return { startAt: start, endAt: "" };
  }

  const monthMap: Record<string, number> = {
    "1M": 1,
    "2M": 2,
    "3M": 3,
    "6M": 6,
  };

  const months = monthMap[dueOption] || 6;
  const endAt = dayjs().add(months, "M").format("YYYY-MM-DD");

  return { startAt: start, endAt };
}
