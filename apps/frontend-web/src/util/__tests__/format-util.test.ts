import { describe, it, expect } from "vitest";
import {
  formatPrice,
  formatCommissionRate,
  formatDateTime,
  formatDate,
  formatSimpleDate,
  formatTime,
  formatDiffDays,
  formatStayTitle,
  formatStayCaption,
} from "../format-util";

describe("format-util", () => {
  describe("formatPrice", () => {
    it("한국 원화 형식으로 가격을 포맷한다", () => {
      expect(formatPrice(10000)).toBe("₩10,000");
    });

    it("0원을 올바르게 포맷한다", () => {
      expect(formatPrice(0)).toBe("₩0");
    });

    it("큰 금액을 올바르게 포맷한다", () => {
      expect(formatPrice(1000000)).toBe("₩1,000,000");
    });
  });

  describe("formatCommissionRate", () => {
    it("소수점 비율을 백분율로 변환한다", () => {
      expect(formatCommissionRate(0.1)).toBe("10%");
    });

    it("0 비율을 처리한다", () => {
      expect(formatCommissionRate(0)).toBe("0%");
    });

    it("1(100%) 비율을 처리한다", () => {
      expect(formatCommissionRate(1)).toBe("100%");
    });
  });

  describe("formatDateTime", () => {
    it("날짜와 시간을 포맷한다", () => {
      const result = formatDateTime("2026-01-15T10:30:45");
      expect(result).toBe("2026-01-15 10:30:45");
    });
  });

  describe("formatDate", () => {
    it("날짜를 YYYY-MM-DD 형식으로 포맷한다", () => {
      const result = formatDate("2026-01-15T10:30:45");
      expect(result).toBe("2026-01-15");
    });
  });

  describe("formatSimpleDate", () => {
    it("날짜를 YY/MM/DD 형식으로 포맷한다", () => {
      const result = formatSimpleDate("2026-01-15T10:30:45");
      expect(result).toBe("26/01/15");
    });
  });

  describe("formatTime", () => {
    it("시간을 HH:mm:ss 형식으로 포맷한다", () => {
      const result = formatTime("2026-01-15T10:30:45");
      expect(result).toBe("10:30:45");
    });
  });

  describe("formatDiffDays", () => {
    it("두 날짜 사이의 일수 차이를 계산한다", () => {
      expect(formatDiffDays("2026-01-01", "2026-01-03")).toBe(2);
    });

    it("같은 날짜면 0을 반환한다", () => {
      expect(formatDiffDays("2026-01-01", "2026-01-01")).toBe(0);
    });

    it("시작일이 null이면 0을 반환한다", () => {
      expect(formatDiffDays(null, "2026-01-03")).toBe(0);
    });

    it("종료일이 null이면 0을 반환한다", () => {
      expect(formatDiffDays("2026-01-01", null)).toBe(0);
    });

    it("둘 다 null이면 0을 반환한다", () => {
      expect(formatDiffDays(null, null)).toBe(0);
    });
  });

  describe("formatStayTitle", () => {
    it("1박을 올바르게 표시한다", () => {
      expect(formatStayTitle(1)).toBe("1박 2일");
    });

    it("2박을 올바르게 표시한다", () => {
      expect(formatStayTitle(2)).toBe("2박 3일");
    });

    it("0박을 처리한다", () => {
      expect(formatStayTitle(0)).toBe("0박 1일");
    });
  });

  describe("formatStayCaption", () => {
    it("숙박 기간 캡션을 생성한다", () => {
      const result = formatStayCaption("2026-01-01", "2026-01-03");
      expect(result).toBe("2026-01-01 ~ 2026-01-03 (2박 3일)");
    });

    it("시작일이 null이면 null을 반환한다", () => {
      expect(formatStayCaption(null, "2026-01-03")).toBeNull();
    });

    it("종료일이 null이면 null을 반환한다", () => {
      expect(formatStayCaption("2026-01-01", null)).toBeNull();
    });

    it("같은 날짜면 null을 반환한다", () => {
      expect(formatStayCaption("2026-01-01", "2026-01-01")).toBeNull();
    });
  });
});
