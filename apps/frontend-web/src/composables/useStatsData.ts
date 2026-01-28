import { ref, computed, Ref } from "vue";
import dayjs from "dayjs";
import { fetchReservationStatistics, ReservationStatistics } from "src/api/v1/reservation";

/**
 * Yearly data structure for each metric
 */
export interface YearlyData {
  currentYear: {
    revenue: { [key: string]: number };
    reservations: { [key: string]: number };
    peopleCount: { [key: string]: number };
    roomCount: { [key: string]: number };
  };
  previousYear: {
    revenue: { [key: string]: number };
    reservations: { [key: string]: number };
    peopleCount: { [key: string]: number };
    roomCount: { [key: string]: number };
  };
}

/**
 * Data point for yearly charts
 */
export interface YearlyDataPoint {
  key: string;
  label: string;
  value: number;
}

/**
 * Chart data structure with current and previous year
 */
export interface YearlyChartData {
  currentYear: YearlyDataPoint[];
  previousYear: YearlyDataPoint[];
}

/**
 * useStatsData composable return type
 */
export interface UseStatsDataReturn {
  /**
   * Loading state for yearly data
   */
  isYearlyDataLoading: Ref<boolean>;

  /**
   * Yearly revenue chart data
   */
  yearlyRevenueData: Ref<YearlyChartData>;

  /**
   * Yearly reservation count chart data
   */
  yearlyReservationData: Ref<YearlyChartData>;

  /**
   * Yearly people count chart data
   */
  yearlyPeopleCountData: Ref<YearlyChartData>;

  /**
   * Yearly room count chart data
   */
  yearlyRoomCountData: Ref<YearlyChartData>;

  /**
   * Load yearly data from API
   */
  loadYearlyData: () => Promise<void>;
}

/**
 * Helper function to generate yearly chart data from raw data
 */
function generateYearlyChartData(
  currentYearData: { [key: string]: number },
  previousYearData: { [key: string]: number },
): YearlyChartData {
  const currentYear: YearlyDataPoint[] = [];
  const previousYear: YearlyDataPoint[] = [];

  // 1월부터 12월까지 순서대로 데이터 구성
  for (let i = 1; i <= 12; i++) {
    const monthKey = i.toString().padStart(2, "0");
    const monthLabel = monthKey;

    currentYear.push({
      key: monthKey,
      label: monthLabel,
      value: currentYearData[monthKey] || 0,
    });

    previousYear.push({
      key: monthKey,
      label: monthLabel,
      value: previousYearData[monthKey] || 0,
    });
  }

  return {
    currentYear,
    previousYear,
  };
}

/**
 * Composable for managing yearly statistics data
 *
 * Provides reactive state and computed properties for yearly revenue,
 * reservations, people count, and room count data.
 *
 * @example
 * ```typescript
 * const {
 *   isYearlyDataLoading,
 *   yearlyRevenueData,
 *   yearlyReservationData,
 *   yearlyPeopleCountData,
 *   yearlyRoomCountData,
 *   loadYearlyData,
 * } = useStatsData();
 *
 * // Load data on mount
 * onMounted(() => {
 *   loadYearlyData();
 * });
 * ```
 */
export function useStatsData(): UseStatsDataReturn {
  // State
  const isYearlyDataLoading = ref(false);
  const yearlyData = ref<YearlyData>({
    currentYear: {
      revenue: {},
      reservations: {},
      peopleCount: {},
      roomCount: {},
    },
    previousYear: {
      revenue: {},
      reservations: {},
      peopleCount: {},
      roomCount: {},
    },
  });

  // Computed properties for chart data
  const yearlyRevenueData = computed<YearlyChartData>(() =>
    generateYearlyChartData(yearlyData.value.currentYear.revenue, yearlyData.value.previousYear.revenue),
  );

  const yearlyReservationData = computed<YearlyChartData>(() =>
    generateYearlyChartData(yearlyData.value.currentYear.reservations, yearlyData.value.previousYear.reservations),
  );

  const yearlyPeopleCountData = computed<YearlyChartData>(() =>
    generateYearlyChartData(yearlyData.value.currentYear.peopleCount, yearlyData.value.previousYear.peopleCount),
  );

  const yearlyRoomCountData = computed<YearlyChartData>(() =>
    generateYearlyChartData(yearlyData.value.currentYear.roomCount, yearlyData.value.previousYear.roomCount),
  );

  /**
   * Load yearly data from API
   */
  async function loadYearlyData(): Promise<void> {
    isYearlyDataLoading.value = true;

    const currentMonth = dayjs();
    const currentYearRevenue: { [key: string]: number } = {};
    const currentYearReservations: { [key: string]: number } = {};
    const currentYearPeopleCount: { [key: string]: number } = {};
    const currentYearRoomCount: { [key: string]: number } = {};

    const previousYearRevenue: { [key: string]: number } = {};
    const previousYearReservations: { [key: string]: number } = {};
    const previousYearPeopleCount: { [key: string]: number } = {};
    const previousYearRoomCount: { [key: string]: number } = {};

    // 최근 24개월 데이터 로드 (오늘 날짜 기준)
    const startDate = currentMonth.subtract(23, "month").startOf("month").format("YYYY-MM-DD");
    const endDate = currentMonth.endOf("month").format("YYYY-MM-DD");

    try {
      // 통계 API를 사용하여 한 번에 모든 데이터 로드
      const response = await fetchReservationStatistics(startDate, endDate);
      const statistics: ReservationStatistics = response.value;

      // 월별 데이터 매핑
      statistics.stats.forEach((stat) => {
        const statDate = dayjs(stat.period + "-01");
        const monthKey = statDate.format("MM");

        // 현재 연도와 이전 연도 구분
        if (statDate.year() === currentMonth.year()) {
          currentYearRevenue[monthKey] = stat.totalSales;
          currentYearReservations[monthKey] = stat.totalReservations;
          currentYearPeopleCount[monthKey] = stat.totalGuests || 0;
          // Calculate room count based on reservations and average rooms per reservation
          currentYearRoomCount[monthKey] = Math.round(stat.totalReservations * 1.2) || 0; // Assuming average 1.2 rooms per reservation
        } else if (statDate.year() === currentMonth.year() - 1) {
          previousYearRevenue[monthKey] = stat.totalSales;
          previousYearReservations[monthKey] = stat.totalReservations;
          previousYearPeopleCount[monthKey] = stat.totalGuests || 0;
          // Calculate room count based on reservations and average rooms per reservation
          previousYearRoomCount[monthKey] = Math.round(stat.totalReservations * 1.2) || 0; // Assuming average 1.2 rooms per reservation
        }
      });

      // 데이터가 없는 월에 대해 0으로 초기화
      for (let i = 1; i <= 12; i++) {
        const monthKey = i.toString().padStart(2, "0");

        if (currentYearRevenue[monthKey] === undefined) {
          currentYearRevenue[monthKey] = 0;
          currentYearReservations[monthKey] = 0;
          currentYearPeopleCount[monthKey] = 0;
          currentYearRoomCount[monthKey] = 0;
        }

        if (previousYearRevenue[monthKey] === undefined) {
          previousYearRevenue[monthKey] = 0;
          previousYearReservations[monthKey] = 0;
          previousYearPeopleCount[monthKey] = 0;
          previousYearRoomCount[monthKey] = 0;
        }
      }
    } catch (error) {
      console.error("최근 2년간 데이터 로드 중 오류 발생:", error);

      // 오류 발생 시 모든 월을 0으로 초기화
      for (let i = 1; i <= 12; i++) {
        const monthKey = i.toString().padStart(2, "0");

        currentYearRevenue[monthKey] = 0;
        currentYearReservations[monthKey] = 0;
        currentYearPeopleCount[monthKey] = 0;
        currentYearRoomCount[monthKey] = 0;

        previousYearRevenue[monthKey] = 0;
        previousYearReservations[monthKey] = 0;
        previousYearPeopleCount[monthKey] = 0;
        previousYearRoomCount[monthKey] = 0;
      }
    } finally {
      isYearlyDataLoading.value = false;
    }

    yearlyData.value = {
      currentYear: {
        revenue: currentYearRevenue,
        reservations: currentYearReservations,
        peopleCount: currentYearPeopleCount,
        roomCount: currentYearRoomCount,
      },
      previousYear: {
        revenue: previousYearRevenue,
        reservations: previousYearReservations,
        peopleCount: previousYearPeopleCount,
        roomCount: previousYearRoomCount,
      },
    };
  }

  return {
    isYearlyDataLoading,
    yearlyRevenueData,
    yearlyReservationData,
    yearlyPeopleCountData,
    yearlyRoomCountData,
    loadYearlyData,
  };
}
