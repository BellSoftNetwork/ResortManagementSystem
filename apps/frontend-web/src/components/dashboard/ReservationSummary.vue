<template>
  <q-card flat bordered>
    <q-card-section class="text-h6 bg-primary text-white">
      <div class="row items-center justify-between">
        <div class="row items-center">
          <q-icon name="hotel" size="md" class="q-mr-sm" />
          <div>입실 정보 요약</div>
        </div>
        <q-spinner v-if="status.isLoading" color="white" size="1.5em" />
      </div>
    </q-card-section>

    <q-card-section>
      <q-inner-loading :showing="status.isLoading">
        <q-spinner size="50px" color="primary" />
      </q-inner-loading>

      <!-- 달력 네비게이션 컨트롤 -->
      <div class="row q-mb-md justify-between items-center">
        <q-btn icon="chevron_left" color="primary" flat round dense @click="calendar?.prev()" />
        <div class="text-h6">{{ currentMonthLabel }}</div>
        <q-btn icon="chevron_right" color="primary" flat round dense @click="calendar?.next()" />
      </div>

      <!-- 새로고침 상태 표시 -->
      <div class="row q-mb-md justify-between items-center">
        <div class="text-caption text-grey">
          마지막 갱신: {{ lastRefreshTime }}
          <q-tooltip>페이지 포커스 복귀 시 자동 갱신됩니다.</q-tooltip>
        </div>
        <q-btn
          icon="refresh"
          color="primary"
          flat
          dense
          :loading="status.isLoading"
          @click="fetchData()"
          label="새로고침"
        />
      </div>

      <ReservationCalendar
        ref="calendar"
        :calendar-events="calendarEvents"
        :selected-date="selectedDate"
        :calendar-year="calendarYear"
        :calendar-month="calendarMonth"
        @navigation="changeView"
        @date-select="onDateSelect"
      />

      <ReservationDayTable
        :reservations="currentDateReservations"
        :selected-date="selectedDate"
        :check-in-out-counts="checkInOutCounts"
        :columns="columns"
        @prev-date="changeDatePrev"
        @next-date="changeDateNext"
      />
    </q-card-section>
  </q-card>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from "vue";
import dayjs from "dayjs";
import { formatDate } from "src/util/format-util";
import { convertTableColumnDef } from "src/util/table-util";
import { getReservationFieldDetail, Reservation } from "src/schema/reservation";
import { fetchReservations } from "src/api/v1/reservation";
import { useReservationCalendar } from "src/composables/useReservationCalendar";
import ReservationCalendar from "src/components/dashboard/ReservationCalendar.vue";
import ReservationDayTable from "src/components/dashboard/ReservationDayTable.vue";

const { formatReservations, calculateExtendedDateRange, getRoomGroupColor } = useReservationCalendar();
const status = ref({
  isLoading: false,
  isLoaded: false,
  isPatching: false,
});
const filter = ref({
  sort: "stayStartAt",
  // 초기값은 나중에 calculateExtendedDateRange 함수로 계산된 값으로 대체됩니다
  stayStartAt: dayjs().startOf("month").format("YYYY-MM-DD"),
  stayEndAt: dayjs().endOf("month").format("YYYY-MM-DD"),
});
const columns = [
  {
    name: "type",
    field: "type",
    label: "구분",
    align: "left",
    required: true,
    sortable: true,
  },
  {
    ...getColumnDef("name"),
    align: "left",
    required: true,
    sortable: true,
  },
  {
    name: "peopleCount",
    field: "peopleCount",
    label: "인원수",
    align: "center",
    required: true,
    sortable: true,
  },
  {
    name: "missPrice",
    field: "missPrice",
    label: "미수금",
    align: "left",
    required: true,
  },
  {
    ...getColumnDef("paymentMethod"),
    align: "left",
    required: true,
    sortable: true,
  },
  {
    ...getColumnDef("rooms"),
    align: "left",
    required: true,
    sortable: true,
  },
  {
    ...getColumnDef("note"),
    align: "left",
    headerStyle: "width: 10%",
  },
];
const calendar = ref<InstanceType<typeof ReservationCalendar>>();
const selectedDate = ref(formatDate());
const currentMonth = ref(dayjs());
const reservationsOfDay = ref({});
const lastRefreshTime = ref("로딩 중...");
const refreshTimer = ref<number | null>(null);
const refreshInterval = ref(5 * 60 * 1000); // 기본값: 5분마다 갱신

// 현재 선택된 날짜의 예약 정보를 계산하는 속성
const currentDateReservations = computed((): Reservation[] => {
  return reservationsOfDay.value[selectedDate.value] || [];
});

const calendarEvents = computed(() => {
  const events = [];

  Object.entries(reservationsOfDay.value).forEach(([date, reservations]) => {
    (reservations as Reservation[]).forEach((reservation) => {
      if (reservation.stayStartAt === date) {
        events.push({
          id: reservation.id,
          name: reservation.name,
          peopleCount: reservation.peopleCount,
          rooms: reservation.rooms.map((room) => room.number),
          start: reservation.stayStartAt,
          end: reservation.stayEndAt,
          bgcolor: getRoomGroupColor(reservation),
          details: `${reservation.name} - ${reservation.stayStartAt} ~ ${reservation.stayEndAt}`,
        });
      }
    });
  });

  return events;
});

const currentMonthLabel = computed(() => {
  return currentMonth.value.format("YYYY년 M월");
});
const calendarYear = computed(() => {
  return currentMonth.value.year();
});
const calendarMonth = computed(() => {
  return currentMonth.value.month() + 1; // dayjs는 0-indexed 월을 반환하지만 QCalendar는 1-indexed 월을 사용
});
const checkInOutCounts = computed(() => {
  const counts: { [date: string]: { checkIn: number; checkOut: number } } = {};

  // Initialize counts for all days in the current month
  const startDate = dayjs(filter.value.stayStartAt);
  const endDate = dayjs(filter.value.stayEndAt);
  let currentDate = startDate;

  while (currentDate.isBefore(endDate) || currentDate.isSame(endDate, "day")) {
    const dateStr = currentDate.format("YYYY-MM-DD");
    counts[dateStr] = { checkIn: 0, checkOut: 0 };
    currentDate = currentDate.add(1, "day");
  }

  // Count check-ins and check-outs
  Object.entries(reservationsOfDay.value).forEach(([date, reservations]) => {
    if (!counts[date]) {
      counts[date] = { checkIn: 0, checkOut: 0 };
    }

    (reservations as Reservation[]).forEach((reservation) => {
      if (reservation.stayStartAt === date) {
        counts[date].checkIn++;
      } else if (reservation.stayEndAt === date) {
        counts[date].checkOut++;
      }
    });
  });

  return counts;
});

function getColumnDef(field: string) {
  return convertTableColumnDef(getReservationFieldDetail(field));
}

function fetchData() {
  status.value.isLoading = true;
  status.value.isLoaded = false;

  fetchReservations({
    stayStartAt: filter.value.stayStartAt,
    stayEndAt: filter.value.stayEndAt,
    status: "NORMAL",
    type: "STAY",
    size: 200, // TODO: 임시 수정
  })
    .then((response) => {
      reservationsOfDay.value = formatReservations(response.values);
      status.value.isLoaded = true;
      // 마지막 갱신 시간 업데이트
      lastRefreshTime.value = dayjs().format("YYYY-MM-DD HH:mm:ss");
    })
    .finally(() => {
      status.value.isLoading = false;
    });
}

function changeView(view) {
  const year = view.year;
  const month = view.month;

  currentMonth.value = dayjs(`${year}-${month}-01`);

  // 확장된 날짜 범위 계산
  const dateRange = calculateExtendedDateRange(currentMonth.value);
  filter.value.stayStartAt = dateRange.startAt;
  filter.value.stayEndAt = dateRange.endAt;

  // 날짜 선택값도 현재 월의 1일로 업데이트
  selectedDate.value = currentMonth.value.startOf("month").format("YYYY-MM-DD");

  fetchData();
}

function onDateSelect(date) {
  // 날짜 선택 시 해당 날짜로 selectedDate 업데이트
  selectedDate.value = date;

  const selectedMonth = dayjs(date);
  // 선택된 날짜의 월이 현재 표시 중인 월과 다른 경우에만 업데이트
  if (selectedMonth.month() !== currentMonth.value.month() || selectedMonth.year() !== currentMonth.value.year()) {
    currentMonth.value = selectedMonth;

    // 확장된 날짜 범위 계산
    const dateRange = calculateExtendedDateRange(currentMonth.value);
    filter.value.stayStartAt = dateRange.startAt;
    filter.value.stayEndAt = dateRange.endAt;

    fetchData();
  }
}

function changeDatePrev() {
  const prevDate = dayjs(selectedDate.value).subtract(1, "day").format("YYYY-MM-DD");

  if (dayjs(prevDate).month() !== dayjs(selectedDate.value).month()) {
    calendar.value?.prev();
  }

  selectedDate.value = prevDate;
}

function changeDateNext() {
  const nextDate = dayjs(selectedDate.value).add(1, "day").format("YYYY-MM-DD");

  if (dayjs(nextDate).month() !== dayjs(selectedDate.value).month()) {
    calendar.value?.next();
  }

  selectedDate.value = nextDate;
}

// 페이지 포커스 이벤트 핸들러
function handleVisibilityChange() {
  if (document.visibilityState === "visible") {
    // 페이지가 다시 보이게 되면 데이터 새로고침
    fetchData();
  }
}

// 키보드 이벤트 핸들러 - 화살표 키로 날짜 이동
function handleKeyDown(e: KeyboardEvent) {
  // 다른 입력 필드에 포커스가 있을 때는 이벤트를 무시
  if (e.target instanceof HTMLInputElement || e.target instanceof HTMLTextAreaElement) {
    return;
  }

  if (e.key === "ArrowLeft") {
    changeDatePrev();
  } else if (e.key === "ArrowRight") {
    changeDateNext();
  }
}

// 자동 갱신 타이머 설정
function setupRefreshTimer() {
  // 기존 타이머가 있으면 정리
  if (refreshTimer.value !== null) {
    clearInterval(refreshTimer.value);
  }

  // 새 타이머 설정 (기본값: 5분)
  refreshTimer.value = window.setInterval(() => {
    fetchData();
  }, refreshInterval.value);
}

onMounted(() => {
  // 초기 로딩 시 현재 월에 대한 확장된 날짜 범위 설정
  const dateRange = calculateExtendedDateRange(currentMonth.value);
  filter.value.stayStartAt = dateRange.startAt;
  filter.value.stayEndAt = dateRange.endAt;

  // 페이지 포커스 변경 이벤트 리스너 등록
  document.addEventListener("visibilitychange", handleVisibilityChange);

  // 키보드 이벤트 리스너 등록
  document.addEventListener("keydown", handleKeyDown);

  // 초기 데이터 로드 및 자동 갱신 타이머 설정
  fetchData();
  setupRefreshTimer();
});

onUnmounted(() => {
  // 컴포넌트 언마운트 시 이벤트 리스너 및 타이머 정리
  document.removeEventListener("visibilitychange", handleVisibilityChange);
  document.removeEventListener("keydown", handleKeyDown);

  if (refreshTimer.value !== null) {
    clearInterval(refreshTimer.value);
    refreshTimer.value = null;
  }
});
</script>
