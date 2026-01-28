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

          <div class="row q-col-gutter-md q-mt-md">
            <!-- 수수료 -->
            <div class="col-12 col-md-6">
              <q-card class="bg-orange-1">
                <q-card-section>
                  <div class="text-subtitle2 text-grey-8">월별 총 수수료</div>
                  <div class="text-h5 text-orange-8">
                    <q-skeleton v-if="isMonthlyDataLoading" type="text" />
                    <span v-else>{{ formatPrice(totalBrokerFee) }}</span>
                  </div>
                </q-card-section>
              </q-card>
            </div>

            <!-- 보증금 -->
            <div class="col-12 col-md-6">
              <q-card class="bg-cyan-1">
                <q-card-section>
                  <div class="text-subtitle2 text-grey-8">월별 총 보증금</div>
                  <div class="text-h5 text-cyan-8">
                    <q-skeleton v-if="isMonthlyDataLoading" type="text" />
                    <span v-else>{{ formatPrice(totalDeposit) }}</span>
                  </div>
                </q-card-section>
              </q-card>
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
import { fetchReservations, FetchReservationsRequestParams } from "src/api/v1/reservation";
import { fetchRooms } from "src/api/v1/room";
import { Reservation } from "src/schema/reservation";
import { Room } from "src/schema/room";
import { useStatsData } from "src/composables/useStatsData";

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

import { formatPrice } from "src/util/format-util";

const $q = useQuasar();
const selectedMonth = ref(dayjs().format("YYYY-MM"));
const reservations = ref<Reservation[]>([]);
const rooms = ref<Room[]>([]);
const isMonthlyDataLoading = ref(false);

// Use composable for yearly data
const {
  isYearlyDataLoading,
  yearlyRevenueData,
  yearlyReservationData,
  yearlyPeopleCountData,
  yearlyRoomCountData,
  loadYearlyData,
} = useStatsData();

// 부대시설인지 확인 (이름에 "부대시설" 포함 여부)
function isFacilityRoom(room: Room): boolean {
  return room.roomGroup.name.includes("부대시설");
}

// 부대시설만 있는 예약 제외
const filteredReservations = computed(() => {
  return reservations.value.filter((reservation) => {
    // 객실이 없는 예약은 포함
    if (reservation.rooms.length === 0) return true;
    // 하나라도 일반 객실이 있으면 포함
    return reservation.rooms.some((room) => !room.roomGroup.name.includes("부대시설"));
  });
});

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
  return filteredReservations.value.reduce((total, reservation) => {
    return total + reservation.price;
  }, 0);
});

// 월별 총 예약 건수
const totalReservations = computed(() => {
  return filteredReservations.value.length;
});

// 월별 객실별 예약 배정 건수
const roomAllocationStats = computed(() => {
  const roomGroups = new Map();

  // 부대시설 제외한 객실만 처리
  const filteredRooms = rooms.value.filter((room) => !isFacilityRoom(room));

  // 객실 그룹별로 초기화
  filteredRooms.forEach((room) => {
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
  filteredReservations.value.forEach((reservation) => {
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
  // 부대시설 제외한 객실만 처리
  const filteredRooms = rooms.value.filter((room) => !isFacilityRoom(room));
  if (filteredRooms.length === 0) return [];

  const daysInMonth = dayjs(selectedMonth.value).daysInMonth();
  const roomGroups = new Map();

  // 객실 그룹별로 초기화
  filteredRooms.forEach((room) => {
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
  filteredReservations.value.forEach((reservation) => {
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
  return filteredReservations.value.reduce((total, reservation) => {
    return total + (reservation.peopleCount || 0);
  }, 0);
});

// 월별 총 수수료 (brokerFee)
const totalBrokerFee = computed(() => {
  return reservations.value.reduce((total, reservation) => {
    return total + (reservation.brokerFee || 0);
  }, 0);
});

// 월별 총 보증금 (deposit)
const totalDeposit = computed(() => {
  return reservations.value.reduce((total, reservation) => {
    return total + (reservation.deposit || 0);
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
    if (yearlyRevenueData.value.currentYear.length === 0) {
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

// 선택된 월이 변경될 때 데이터 다시 로드
watch(selectedMonth, () => {
  loadData();
});

onMounted(() => {
  loadData();
});
</script>
