import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { calculateDateRange, type DueOption } from "../date-preset-util";
import * as formatUtil from "../format-util";

vi.mock("../format-util", () => ({
  formatDate: vi.fn(() => "2026-01-26"),
}));

describe("calculateDateRange", () => {
  const mockToday = "2026-01-26";

  beforeEach(() => {
    const mockDate = new Date("2026-01-26T00:00:00.000Z");
    vi.setSystemTime(mockDate);
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  describe("ALL 옵션", () => {
    it("ALL 옵션은 빈 날짜 범위를 반환한다", () => {
      const result = calculateDateRange("ALL");

      expect(result).toEqual({
        startAt: "",
        endAt: "",
      });
    });

    it("ALL 옵션은 startDate 파라미터를 무시하고 빈 날짜 범위를 반환한다", () => {
      const result = calculateDateRange("ALL", "2026-03-15");

      expect(result).toEqual({
        startAt: "",
        endAt: "",
      });
    });
  });

  describe("CUSTOM 옵션", () => {
    it("CUSTOM 옵션은 오늘 날짜를 startAt으로, 빈 문자열을 endAt으로 반환한다", () => {
      const result = calculateDateRange("CUSTOM");

      expect(result.startAt).toBe(mockToday);
      expect(result.endAt).toBe("");
    });

    it("CUSTOM 옵션은 제공된 startDate를 startAt으로, 빈 문자열을 endAt으로 반환한다", () => {
      const customDate = "2026-03-15";
      const result = calculateDateRange("CUSTOM", customDate);

      expect(result).toEqual({
        startAt: customDate,
        endAt: "",
      });
    });

    it("CUSTOM 옵션은 과거 날짜를 startDate로 받을 수 있다", () => {
      const pastDate = "2025-12-01";
      const result = calculateDateRange("CUSTOM", pastDate);

      expect(result).toEqual({
        startAt: pastDate,
        endAt: "",
      });
    });

    it("CUSTOM 옵션은 미래 날짜를 startDate로 받을 수 있다", () => {
      const futureDate = "2027-06-30";
      const result = calculateDateRange("CUSTOM", futureDate);

      expect(result).toEqual({
        startAt: futureDate,
        endAt: "",
      });
    });
  });

  describe("1M 옵션", () => {
    it("1M 옵션은 오늘부터 1개월 후까지의 날짜 범위를 반환한다", () => {
      const result = calculateDateRange("1M");

      expect(result.startAt).toBe(mockToday);
      expect(result.endAt).toBe("2026-02-26");
    });

    it("1M 옵션은 제공된 startDate를 startAt으로 사용한다", () => {
      const customStart = "2026-02-01";
      const result = calculateDateRange("1M", customStart);

      expect(result.startAt).toBe(customStart);
      expect(result.endAt).toBe("2026-02-26"); // endAt은 항상 오늘 기준
    });
  });

  describe("2M 옵션", () => {
    it("2M 옵션은 오늘부터 2개월 후까지의 날짜 범위를 반환한다", () => {
      const result = calculateDateRange("2M");

      expect(result.startAt).toBe(mockToday);
      expect(result.endAt).toBe("2026-03-26");
    });

    it("2M 옵션은 제공된 startDate를 startAt으로 사용한다", () => {
      const customStart = "2026-01-01";
      const result = calculateDateRange("2M", customStart);

      expect(result.startAt).toBe(customStart);
      expect(result.endAt).toBe("2026-03-26");
    });
  });

  describe("3M 옵션", () => {
    it("3M 옵션은 오늘부터 3개월 후까지의 날짜 범위를 반환한다", () => {
      const result = calculateDateRange("3M");

      expect(result.startAt).toBe(mockToday);
      expect(result.endAt).toBe("2026-04-26");
    });

    it("3M 옵션은 제공된 startDate를 startAt으로 사용한다", () => {
      const customStart = "2025-12-15";
      const result = calculateDateRange("3M", customStart);

      expect(result.startAt).toBe(customStart);
      expect(result.endAt).toBe("2026-04-26");
    });
  });

  describe("6M 옵션", () => {
    it("6M 옵션은 오늘부터 6개월 후까지의 날짜 범위를 반환한다", () => {
      const result = calculateDateRange("6M");

      expect(result.startAt).toBe(mockToday);
      expect(result.endAt).toBe("2026-07-26");
    });

    it("6M 옵션은 제공된 startDate를 startAt으로 사용한다", () => {
      const customStart = "2026-01-15";
      const result = calculateDateRange("6M", customStart);

      expect(result.startAt).toBe(customStart);
      expect(result.endAt).toBe("2026-07-26");
    });
  });

  describe("엣지 케이스", () => {
    it("모든 DueOption 타입이 올바른 DateRange를 반환한다", () => {
      const options: DueOption[] = ["ALL", "1M", "2M", "3M", "6M", "CUSTOM"];

      options.forEach((option) => {
        const result = calculateDateRange(option);

        expect(result).toHaveProperty("startAt");
        expect(result).toHaveProperty("endAt");
        expect(typeof result.startAt).toBe("string");
        expect(typeof result.endAt).toBe("string");
      });
    });

    it("startDate가 빈 문자열일 때 formatDate()를 사용한다", () => {
      const result = calculateDateRange("1M", "");

      expect(result.startAt).toBe(mockToday);
      expect(result.endAt).toBe("2026-02-26");
    });

    it("월 경계를 넘어가는 날짜 계산이 올바르게 작동한다", () => {
      const endOfMonth = "2026-01-31";
      vi.setSystemTime(new Date("2026-01-31T00:00:00.000Z"));
      vi.mocked(formatUtil.formatDate).mockReturnValue(endOfMonth);

      const result = calculateDateRange("1M");

      expect(result.startAt).toBe(endOfMonth);
      expect(result.endAt).toBe("2026-02-28");
    });

    it("연도를 넘어가는 날짜 계산이 올바르게 작동한다", () => {
      const endOfYear = "2025-12-15";
      vi.setSystemTime(new Date("2025-12-15T00:00:00.000Z"));
      vi.mocked(formatUtil.formatDate).mockReturnValue(endOfYear);

      const result = calculateDateRange("3M");

      expect(result.startAt).toBe(endOfYear);
      expect(result.endAt).toBe("2026-03-15");
    });

    it("날짜 형식이 YYYY-MM-DD 형식을 유지한다", () => {
      const result = calculateDateRange("6M");

      expect(result.startAt).toMatch(/^\d{4}-\d{2}-\d{2}$/);
      expect(result.endAt).toMatch(/^\d{4}-\d{2}-\d{2}$/);
    });

    it("startDate가 유효하지 않은 날짜 형식이어도 그대로 반환한다", () => {
      const invalidDate = "invalid-date";
      const result = calculateDateRange("CUSTOM", invalidDate);

      expect(result.startAt).toBe(invalidDate);
      expect(result.endAt).toBe("");
    });
  });

  describe("타입 안정성", () => {
    it("반환값이 DateRange 인터페이스를 만족한다", () => {
      const result = calculateDateRange("1M");

      expect(result).toMatchObject({
        startAt: expect.any(String),
        endAt: expect.any(String),
      });
    });

    it("모든 DueOption 값에 대해 일관된 반환 타입을 제공한다", () => {
      const options: DueOption[] = ["ALL", "1M", "2M", "3M", "6M", "CUSTOM"];

      options.forEach((option) => {
        const result = calculateDateRange(option);
        expect(Object.keys(result)).toEqual(["startAt", "endAt"]);
      });
    });
  });
});
