import { describe, it, expect, vi, beforeEach } from "vitest";
import { useStatsData } from "../useStatsData";
import * as reservationApi from "src/api/v1/reservation";
import { createMockApiResponse } from "test/vitest/helpers";
import type { ReservationStatistics } from "src/api/v1/reservation";
import dayjs from "dayjs";

vi.mock("src/api/v1/reservation");

describe("useStatsData", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe("초기 상태", () => {
    it("isYearlyDataLoading은 false로 초기화된다", () => {
      const { isYearlyDataLoading } = useStatsData();

      expect(isYearlyDataLoading.value).toBe(false);
    });

    it("yearlyRevenueData는 빈 데이터로 초기화된다", () => {
      const { yearlyRevenueData } = useStatsData();

      expect(yearlyRevenueData.value.currentYear).toHaveLength(12);
      expect(yearlyRevenueData.value.previousYear).toHaveLength(12);
      yearlyRevenueData.value.currentYear.forEach((point) => {
        expect(point.value).toBe(0);
      });
    });

    it("yearlyReservationData는 빈 데이터로 초기화된다", () => {
      const { yearlyReservationData } = useStatsData();

      expect(yearlyReservationData.value.currentYear).toHaveLength(12);
      expect(yearlyReservationData.value.previousYear).toHaveLength(12);
    });

    it("yearlyPeopleCountData는 빈 데이터로 초기화된다", () => {
      const { yearlyPeopleCountData } = useStatsData();

      expect(yearlyPeopleCountData.value.currentYear).toHaveLength(12);
      expect(yearlyPeopleCountData.value.previousYear).toHaveLength(12);
    });

    it("yearlyRoomCountData는 빈 데이터로 초기화된다", () => {
      const { yearlyRoomCountData } = useStatsData();

      expect(yearlyRoomCountData.value.currentYear).toHaveLength(12);
      expect(yearlyRoomCountData.value.previousYear).toHaveLength(12);
    });
  });

  describe("loadYearlyData", () => {
    it("API 호출 중 로딩 상태를 관리한다", async () => {
      const mockStats: ReservationStatistics = { stats: [] };
      vi.mocked(reservationApi.fetchReservationStatistics).mockResolvedValue(
        createMockApiResponse(mockStats)
      );

      const { isYearlyDataLoading, loadYearlyData } = useStatsData();

      expect(isYearlyDataLoading.value).toBe(false);

      const promise = loadYearlyData();

      expect(isYearlyDataLoading.value).toBe(true);

      await promise;

      expect(isYearlyDataLoading.value).toBe(false);
    });

    it("API 응답 데이터를 올바르게 매핑한다", async () => {
      const currentYear = dayjs().year();
      const mockStats: ReservationStatistics = {
        stats: [
          {
            period: `${currentYear}-01`,
            totalSales: 1000000,
            totalReservations: 10,
            totalGuests: 25,
          },
          {
            period: `${currentYear}-06`,
            totalSales: 2000000,
            totalReservations: 20,
            totalGuests: 50,
          },
          {
            period: `${currentYear - 1}-01`,
            totalSales: 800000,
            totalReservations: 8,
            totalGuests: 20,
          },
        ],
      };

      vi.mocked(reservationApi.fetchReservationStatistics).mockResolvedValue(
        createMockApiResponse(mockStats)
      );

      const { yearlyRevenueData, yearlyReservationData, loadYearlyData } = useStatsData();

      await loadYearlyData();

      expect(yearlyRevenueData.value.currentYear[0].value).toBe(1000000);
      expect(yearlyRevenueData.value.currentYear[5].value).toBe(2000000);
      expect(yearlyRevenueData.value.previousYear[0].value).toBe(800000);

      expect(yearlyReservationData.value.currentYear[0].value).toBe(10);
      expect(yearlyReservationData.value.currentYear[5].value).toBe(20);
    });

    it("데이터가 없는 월은 0으로 초기화된다", async () => {
      const mockStats: ReservationStatistics = { stats: [] };
      vi.mocked(reservationApi.fetchReservationStatistics).mockResolvedValue(
        createMockApiResponse(mockStats)
      );

      const { yearlyRevenueData, loadYearlyData } = useStatsData();

      await loadYearlyData();

      yearlyRevenueData.value.currentYear.forEach((point) => {
        expect(point.value).toBe(0);
      });
      yearlyRevenueData.value.previousYear.forEach((point) => {
        expect(point.value).toBe(0);
      });
    });

    it("API 에러 시 모든 데이터를 0으로 초기화한다", async () => {
      vi.mocked(reservationApi.fetchReservationStatistics).mockRejectedValue(
        new Error("API Error")
      );

      const consoleErrorSpy = vi.spyOn(console, "error").mockImplementation(() => {});

      const { yearlyRevenueData, isYearlyDataLoading, loadYearlyData } = useStatsData();

      await loadYearlyData();

      expect(isYearlyDataLoading.value).toBe(false);
      yearlyRevenueData.value.currentYear.forEach((point) => {
        expect(point.value).toBe(0);
      });

      consoleErrorSpy.mockRestore();
    });

    it("월별 라벨이 올바르게 설정된다", async () => {
      const mockStats: ReservationStatistics = { stats: [] };
      vi.mocked(reservationApi.fetchReservationStatistics).mockResolvedValue(
        createMockApiResponse(mockStats)
      );

      const { yearlyRevenueData, loadYearlyData } = useStatsData();

      await loadYearlyData();

      const labels = yearlyRevenueData.value.currentYear.map((p) => p.label);
      expect(labels).toEqual([
        "01", "02", "03", "04", "05", "06",
        "07", "08", "09", "10", "11", "12",
      ]);
    });
  });
});
