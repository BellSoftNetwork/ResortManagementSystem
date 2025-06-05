<template>
  <q-page padding>
    <div class="q-pa-md">
      <div class="text-h4 q-mb-md">통계</div>

      <!-- 매출액 그래프 (전년 대비) -->
      <div class="row q-col-gutter-md q-mb-md">
        <div class="col-12">
          <YearlyRevenueChart :chartData="yearlyRevenueData" :isLoading="isYearlyDataLoading" />
        </div>
      </div>

      <!-- 예약 건수 그래프 (전년 대비) -->
      <div class="row q-col-gutter-md q-mb-md">
        <div class="col-12">
          <YearlyReservationCountChart :chartData="yearlyReservationData" :isLoading="isYearlyDataLoading" />
        </div>
      </div>

      <!-- 다녀간 인원 수 그래프 (전년 대비) -->
      <div class="row q-col-gutter-md q-mb-md">
        <div class="col-12">
          <YearlyPeopleCountChart :chartData="yearlyPeopleCountData" :isLoading="isYearlyDataLoading" />
        </div>
      </div>

      <!-- 예약된 객실 수 그래프 (전년 대비) -->
      <div class="row q-col-gutter-md q-mb-md">
        <div class="col-12">
          <YearlyRoomCountChart :chartData="yearlyRoomCountData" :isLoading="isYearlyDataLoading" />
        </div>
      </div>

      <div class="row q-col-gutter-md">
        <!-- 월 선택 -->
        <div class="col-12">
          <MonthSelector :selectedMonth="selectedMonth" @update:selectedMonth="handleMonthChange" />
        </div>

        <!-- 월별 통계 데이터 -->
        <div class="col-12">
          <div class="row q-col-gutter-md q-mb-md">
            <!-- 매출액 -->
            <div class="col-12 col-md-6">
              <MonthlySalesCard :totalSales="totalSales" :isLoading="isMonthlyDataLoading" />
            </div>

            <!-- 예약 건수 -->
            <div class="col-12 col-md-6">
              <MonthlyReservationCountCard :totalReservations="totalReservations" :isLoading="isMonthlyDataLoading" />
            </div>
          </div>

          <div class="row q-col-gutter-md q-mb-md">
            <!-- 객실별 예약 배정 건수 -->
            <div class="col-12">
              <RoomAllocationTable :roomAllocationStats="roomAllocationStats" :isLoading="isMonthlyDataLoading" />
            </div>
          </div>

          <div class="row q-col-gutter-md q-mb-md">
            <!-- 객실 그룹별 평균 점유율 -->
            <div class="col-12">
              <RoomGroupOccupancyTable
                :roomGroupOccupancyStats="roomGroupOccupancyStats"
                :isLoading="isMonthlyDataLoading"
              />
            </div>
          </div>

          <div class="row q-col-gutter-md">
            <!-- 총 다녀간 인원 -->
            <div class="col-12 col-md-6">
              <MonthlyGuestsCard :totalGuests="totalGuests" :isLoading="isMonthlyDataLoading" />
            </div>
          </div>
        </div>
      </div>
    </div>
  </q-page>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from "vue";
import dayjs from "dayjs";
import { useQuasar } from "quasar";
import {
  fetchReservations,
  FetchReservationsRequestParams,
  fetchReservationStatistics,
  ReservationStatistics,
} from "src/api/v1/reservation";
import { fetchRooms } from "src/api/v1/room";
// formatPrice is now used in the YearlyRevenueChart component
import { Reservation } from "src/schema/reservation";
import { Room } from "src/schema/room";

// Import components
import YearlyRevenueChart from "src/components/stats/YearlyRevenueChart.vue";
import YearlyReservationCountChart from "src/components/stats/YearlyReservationCountChart.vue";
import YearlyPeopleCountChart from "src/components/stats/YearlyPeopleCountChart.vue";
import YearlyRoomCountChart from "src/components/stats/YearlyRoomCountChart.vue";
import MonthSelector from "src/components/stats/MonthSelector.vue";
import MonthlySalesCard from "src/components/stats/MonthlySalesCard.vue";
import MonthlyReservationCountCard from "src/components/stats/MonthlyReservationCountCard.vue";
import RoomAllocationTable from "src/components/stats/RoomAllocationTable.vue";
import RoomGroupOccupancyTable from "src/components/stats/RoomGroupOccupancyTable.vue";
import MonthlyGuestsCard from "src/components/stats/MonthlyGuestsCard.vue";

const $q = useQuasar();
const selectedMonth = ref(dayjs().format("YYYY-MM"));
const reservations = ref<Reservation[]>([]);
const rooms = ref<Room[]>([]);
const isMonthlyDataLoading = ref(false);
const isYearlyDataLoading = ref(false);
const yearlyData = ref<{
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
}>({
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

// 최근 1년간 매출액 데이터
const yearlyRevenueData = computed(() => {
  const currentYearData = [];
  const previousYearData = [];

  // 1월부터 12월까지 순서대로 데이터 구성
  for (let i = 1; i <= 12; i++) {
    const monthKey = i.toString().padStart(2, "0");
    const monthLabel = monthKey;

    currentYearData.push({
      key: monthKey,
      label: monthLabel,
      value: yearlyData.value.currentYear.revenue[monthKey] || 0,
    });

    previousYearData.push({
      key: monthKey,
      label: monthLabel,
      value: yearlyData.value.previousYear.revenue[monthKey] || 0,
    });
  }

  return {
    currentYear: currentYearData,
    previousYear: previousYearData,
  };
});

// 최근 1년간 예약 건수 데이터
const yearlyReservationData = computed(() => {
  const currentYearData = [];
  const previousYearData = [];

  // 1월부터 12월까지 순서대로 데이터 구성
  for (let i = 1; i <= 12; i++) {
    const monthKey = i.toString().padStart(2, "0");
    const monthLabel = monthKey;

    currentYearData.push({
      key: monthKey,
      label: monthLabel,
      value: yearlyData.value.currentYear.reservations[monthKey] || 0,
    });

    previousYearData.push({
      key: monthKey,
      label: monthLabel,
      value: yearlyData.value.previousYear.reservations[monthKey] || 0,
    });
  }

  return {
    currentYear: currentYearData,
    previousYear: previousYearData,
  };
});

// 최근 1년간 다녀간 인원 수 데이터
const yearlyPeopleCountData = computed(() => {
  const currentYearData = [];
  const previousYearData = [];

  // 1월부터 12월까지 순서대로 데이터 구성
  for (let i = 1; i <= 12; i++) {
    const monthKey = i.toString().padStart(2, "0");
    const monthLabel = monthKey;

    currentYearData.push({
      key: monthKey,
      label: monthLabel,
      value: yearlyData.value.currentYear.peopleCount[monthKey] || 0,
    });

    previousYearData.push({
      key: monthKey,
      label: monthLabel,
      value: yearlyData.value.previousYear.peopleCount[monthKey] || 0,
    });
  }

  return {
    currentYear: currentYearData,
    previousYear: previousYearData,
  };
});

// 최근 1년간 예약된 객실 수 데이터
const yearlyRoomCountData = computed(() => {
  const currentYearData = [];
  const previousYearData = [];

  // 1월부터 12월까지 순서대로 데이터 구성
  for (let i = 1; i <= 12; i++) {
    const monthKey = i.toString().padStart(2, "0");
    const monthLabel = monthKey;

    currentYearData.push({
      key: monthKey,
      label: monthLabel,
      value: yearlyData.value.currentYear.roomCount[monthKey] || 0,
    });

    previousYearData.push({
      key: monthKey,
      label: monthLabel,
      value: yearlyData.value.previousYear.roomCount[monthKey] || 0,
    });
  }

  return {
    currentYear: currentYearData,
    previousYear: previousYearData,
  };
});

// 최대값은 이제 각 컴포넌트 내부에서 계산됩니다.

// 월 이동 함수
function handleMonthChange(direction: string) {
  if (direction === "previous") {
    selectedMonth.value = dayjs(selectedMonth.value).subtract(1, "month").format("YYYY-MM");
  } else if (direction === "next") {
    selectedMonth.value = dayjs(selectedMonth.value).add(1, "month").format("YYYY-MM");
  }
}

// 월별 총 매출액
const totalSales = computed(() => {
  return reservations.value.reduce((total, reservation) => {
    return total + reservation.price;
  }, 0);
});

// 월별 총 예약 건수
const totalReservations = computed(() => {
  return reservations.value.length;
});

// 월별 객실별 예약 배정 건수
const roomAllocationStats = computed(() => {
  const roomGroups = new Map();

  // 객실 그룹별로 초기화
  rooms.value.forEach((room) => {
    const groupName = room.roomGroup.name;
    if (!roomGroups.has(groupName)) {
      roomGroups.set(groupName, {
        roomGroup: groupName,
        count: 0,
        rooms: [],
      });
    }

    roomGroups.get(groupName).rooms.push(room.number);
  });

  // 예약 데이터로 카운트
  reservations.value.forEach((reservation) => {
    reservation.rooms.forEach((room) => {
      const groupName = room.roomGroup.name;
      if (roomGroups.has(groupName)) {
        roomGroups.get(groupName).count++;
      }
    });
  });

  return Array.from(roomGroups.values());
});

// 테이블 컬럼 정의는 이제 각 테이블 컴포넌트 내부에서 정의됩니다.

// 월별 객실 그룹별 평균 점유율
const roomGroupOccupancyStats = computed(() => {
  if (rooms.value.length === 0) return [];

  const daysInMonth = dayjs(selectedMonth.value).daysInMonth();
  const roomGroups = new Map();

  // 객실 그룹별로 초기화
  rooms.value.forEach((room) => {
    const groupName = room.roomGroup.name;
    const groupId = room.roomGroup.id;

    if (!roomGroups.has(groupId)) {
      roomGroups.set(groupId, {
        roomGroup: groupName,
        roomCount: 0,
        occupiedDays: 0,
        occupancyRate: 0,
      });
    }

    roomGroups.get(groupId).roomCount++;
  });

  // 예약 데이터로 점유일 계산
  reservations.value.forEach((reservation) => {
    const startDate = dayjs(reservation.stayStartAt);
    const endDate = dayjs(reservation.stayEndAt);
    const stayDuration = endDate.diff(startDate, "day");

    reservation.rooms.forEach((room) => {
      const groupId = room.roomGroup.id;
      if (roomGroups.has(groupId)) {
        roomGroups.get(groupId).occupiedDays += stayDuration;
      }
    });
  });

  // 점유율 계산
  roomGroups.forEach((group) => {
    const totalPossibleDays = group.roomCount * daysInMonth;
    group.occupancyRate = totalPossibleDays > 0 ? ((group.occupiedDays / totalPossibleDays) * 100).toFixed(2) : 0;
  });

  return Array.from(roomGroups.values());
});

// 테이블 컬럼 정의는 이제 각 테이블 컴포넌트 내부에서 정의됩니다.

// 월별 총 다녀간 인원
const totalGuests = computed(() => {
  return reservations.value.reduce((total, reservation) => {
    return total + (reservation.peopleCount || 0);
  }, 0);
});

// 데이터 로드 함수
async function loadData() {
  isMonthlyDataLoading.value = true;

  try {
    // 현재 선택된 월의 데이터 로드
    const startDate = dayjs(selectedMonth.value).startOf("month").format("YYYY-MM-DD");
    const endDate = dayjs(selectedMonth.value).endOf("month").format("YYYY-MM-DD");

    const params: FetchReservationsRequestParams = {
      stayStartAt: startDate,
      stayEndAt: endDate,
      size: 1000, // 충분히 큰 값으로 설정
    };

    const [reservationsResponse, roomsResponse] = await Promise.all([fetchReservations(params), fetchRooms({})]);

    reservations.value = reservationsResponse.values;
    rooms.value = roomsResponse.values;

    // 최근 1년간 데이터 로드 (첫 로드 시에만)
    if (!yearlyData.value.revenue || Object.keys(yearlyData.value.revenue).length === 0) {
      await loadYearlyData();
    }
  } catch (error) {
    console.error("데이터 로드 중 오류 발생:", error);
    $q.notify({
      message: "통계 데이터를 불러오는 중 오류가 발생했습니다.",
      type: "negative",
      actions: [{ icon: "close", color: "white", round: true }],
    });
  } finally {
    isMonthlyDataLoading.value = false;
  }
}

// 최근 2년간 데이터 로드
async function loadYearlyData() {
  isYearlyDataLoading.value = true;

  const currentMonth = dayjs();
  const currentYearRevenue = {};
  const currentYearReservations = {};
  const currentYearPeopleCount = {};
  const currentYearRoomCount = {};

  const previousYearRevenue = {};
  const previousYearReservations = {};
  const previousYearPeopleCount = {};
  const previousYearRoomCount = {};

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

// 선택된 월이 변경될 때 데이터 다시 로드
watch(selectedMonth, () => {
  loadData();
});

onMounted(() => {
  loadData();
});
</script>
