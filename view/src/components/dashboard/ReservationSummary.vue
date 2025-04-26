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

      <div class="row">
        <div class="col q-pa-md-sm">
          <QCalendarMonth
            ref="calendar"
            v-model="selectedDate"
            @navigation="changeView"
            @update:model-value="onDateSelect"
            :year="calendarYear"
            :month="calendarMonth"
            :hoverable="true"
            :focusable="false"
            :focus-type="['day', 'date']"
            :now="formatDate"
            mask="YYYY-MM-DD"
            locale="ko-KR"
            day-min-height="150"
            :day-height="0"
            animated
            bordered
          >
            <template #week="{ scope: { week } }">
              <template v-for="(displayedEvent, index) in getWeekEvents(week)" :key="index">
                <div :class="badgeClasses(displayedEvent)" :style="badgeStyles(displayedEvent, week.length)">
                  <q-btn
                    v-if="displayedEvent.event && displayedEvent.event.name"
                    :to="{ name: 'Reservation', params: { id: displayedEvent.event.id } }"
                    class="full-width full-height q-pa-none block"
                    size="sm"
                    dense
                    flat
                  >
                    {{ displayedEvent.event.name }}님&nbsp;
                    <span v-if="displayedEvent.event.peopleCount > 1">
                      외 {{ displayedEvent.event.peopleCount - 1 }}명 </span
                    >&nbsp;
                    <span v-if="displayedEvent.event.rooms.length > 0">
                      ({{ displayedEvent.event.rooms.join(", ") }})
                    </span>
                    <span v-else>(미배정)</span>
                  </q-btn>
                </div>
              </template>
            </template>
          </QCalendarMonth>
        </div>
      </div>

      <div class="row q-pt-sm">
        <div class="col q-pa-md-sm">
          <div style="min-height: 500px">
            <q-table
              :columns="columns"
              :rows="currentDateReservations"
              row-key="id"
              flat
              bordered
              :pagination="{ rowsPerPage: 20, sortBy: 'type', descending: false }"
            >
              <template #top>
                <div class="row justify-between items-center q-pa-sm full-width">
                  <div class="row items-center">
                    <q-btn icon="chevron_left" color="primary" flat round dense @click="changeDatePrev()" />
                    <div class="text-h6">{{ selectedDate }} 예약 현황</div>
                    <q-btn icon="chevron_right" color="primary" flat round dense @click="changeDateNext()" />
                  </div>
                  <div class="row q-gutter-sm">
                    <q-badge v-if="checkInOutCounts[selectedDate]?.checkIn > 0" color="positive" class="q-px-sm">
                      입실 {{ checkInOutCounts[selectedDate]?.checkIn }}
                    </q-badge>
                    <q-badge v-if="checkInOutCounts[selectedDate]?.checkOut > 0" color="negative" class="q-px-sm">
                      퇴실 {{ checkInOutCounts[selectedDate]?.checkOut }}
                    </q-badge>
                    <q-badge v-if="!currentDateReservations.length" color="grey" class="q-px-sm"> 예약 없음</q-badge>
                  </div>
                </div>
              </template>

              <template #header="props">
                <q-tr :props="props">
                  <q-th v-for="col in props.cols" :key="col.name" :props="props" class="bg-blue-1">
                    {{ col.label }}
                  </q-th>
                </q-tr>
              </template>

              <template #body-cell-type="props">
                <q-td key="type" :props="props">
                  <q-badge :color="getTypeColor(props.row.type)" text-color="white" class="q-px-sm">
                    {{ props.row.type }}
                  </q-badge>
                </q-td>
              </template>

              <template #body-cell-rooms="props">
                <q-td key="rooms" :props="props">
                  <div v-if="props.row.rooms.length !== 0">
                    <span v-for="room in props.row.rooms" :key="room.id">
                      <div v-if="authStore.isAdminRole">
                        <q-btn
                          :to="{
                            name: 'Room',
                            params: { id: room.id },
                          }"
                          align="left"
                          color="primary"
                          dense
                          flat
                          >{{ room.number }}
                        </q-btn>
                      </div>
                      <div v-else>
                        {{ room.number }}
                      </div>
                    </span>
                  </div>
                  <div v-else class="text-grey">미배정</div>
                </q-td>
              </template>

              <template #body-cell-name="props">
                <q-td key="name" :props="props">
                  <div v-if="authStore.isAdminRole">
                    <q-btn
                      :to="{
                        name: 'Reservation',
                        params: { id: props.row.id },
                      }"
                      class="full-width"
                      align="left"
                      color="primary"
                      dense
                      flat
                      >{{ props.row.name }}
                    </q-btn>
                  </div>
                  <div v-else>
                    {{ props.row.name }}
                  </div>
                </q-td>
              </template>

              <template #body-cell-missPrice="props">
                <q-td :props="props" key="missPrice" :class="missPriceBackgroundColor(props.row)">
                  {{ formatPrice(props.row.missPrice) }}
                </q-td>
              </template>

              <template #body-cell-note="props">
                <q-td :props="props" key="note">
                  <q-btn
                    v-if="props.row.note"
                    @click="
                      $q.dialog({
                        title: `${props.row.name}님 예약 메모`,
                        message: props.row.note,
                      })
                    "
                    color="primary"
                  >
                    메모 확인
                  </q-btn>
                </q-td>
              </template>

              <template #no-data>
                <div class="full-width row justify-center q-py-md">
                  <div class="text-grey text-body1">이 날짜에 등록된 예약이 없습니다.</div>
                </div>
              </template>
            </q-table>
          </div>
        </div>
      </div>
    </q-card-section>
  </q-card>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from "vue";
import dayjs from "dayjs";
import { useQuasar } from "quasar";
import { daysBetween, isOverlappingDates, parsed, QCalendarMonth, Timestamp } from "@quasar/quasar-ui-qcalendar";
import { indexOf } from "@quasar/quasar-ui-qcalendar/src/utils/helpers.js";
import "@quasar/quasar-ui-qcalendar/dist/index.css";
import { useAuthStore } from "stores/auth";
import { formatDate, formatPrice } from "src/util/format-util";
import { convertTableColumnDef } from "src/util/table-util";
import { getReservationFieldDetail, Reservation } from "src/schema/reservation";
import { fetchReservations } from "src/api/v1/reservation";

const $q = useQuasar();
const authStore = useAuthStore();
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
const calendar = ref<QCalendarMonth>();
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

function formatReservations(reservations: Reservation[]) {
  const reservationMap: {
    [date: string]: Reservation[];
  } = {};

  reservations.forEach((reservation) => {
    reservation.missPrice = reservation.price - reservation.paymentAmount;

    for (const [index, date] of getDateArray(reservation.stayStartAt, reservation.stayEndAt).entries()) {
      if (!Object.keys(reservationMap).includes(date)) reservationMap[date] = [];

      const reservationCopy = { ...reservation, type: "N/A" };

      if (index === 0) reservationCopy.type = "입실";
      else if (date === reservation.stayEndAt) reservationCopy.type = "퇴실";
      else reservationCopy.type = "연박";

      reservationMap[date].push(reservationCopy);
    }
  });

  return reservationMap;
}

function getDateArray(startDate: string, endDate: string) {
  const stayStartDate = formatDate(startDate);
  const stayEndAt = formatDate(endDate);
  const dateArray = [];

  for (let date = stayStartDate; date <= stayEndAt; date = dayjs(date).add(1, "day").format("YYYY-MM-DD")) {
    dateArray.push(date);
  }

  return dateArray;
}

// 달력에 표시되는 확장된 날짜 범위를 계산하는 함수
function calculateExtendedDateRange(month: dayjs.Dayjs) {
  // 해당 월의 시작일과 마지막일
  const monthStart = month.startOf("month");
  const monthEnd = month.endOf("month");

  // 해당 월의 첫 날짜가 속한 주의 첫째 날 (전월의 마지막 주 포함)
  const startOfFirstWeek = monthStart.startOf("week");

  // 해당 월의 마지막 날짜가 속한 주의 마지막 날 (다음 월의 첫째 주 포함)
  const endOfLastWeek = monthEnd.endOf("week");

  return {
    startAt: startOfFirstWeek.format("YYYY-MM-DD"),
    endAt: endOfLastWeek.format("YYYY-MM-DD"),
  };
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

function missPriceBackgroundColor(value) {
  if (value.paymentMethod.requireUnpaidAmountCheck === false) return "";

  if (value.missPrice > 0) return "bg-warning";

  return "";
}

// 객실 그룹 ID에 따라 사용할 색상 배열 정의
const colorSet = [
  "purple",
  "deep-purple",
  "indigo",
  "blue",
  "light-blue",
  "cyan",
  "teal",
  "green",
  "light-green",
  "lime",
  "yellow",
  "amber",
  "orange",
  "deep-orange",
  "brown",
  "blue-grey",
  "pink",
  "red",
];

// 색상 캐시 (객실 그룹 ID => 색상 인덱스)
const roomGroupColorMap = new Map<number, number>();
let colorIndex = 0;

// 객실 그룹 ID에 기반한 색상 반환
function getRoomGroupColor(reservation: Reservation) {
  // 객실이 배정되지 않은 경우
  if (!reservation.rooms || reservation.rooms.length === 0) {
    return "grey-7";
  }

  // 첫 번째 객실의 그룹 ID를 사용
  const roomGroupId = reservation.rooms[0].roomGroup.id;

  // 이미 배정된 색상이 있으면 해당 색상 사용
  if (roomGroupColorMap.has(roomGroupId)) {
    return colorSet[roomGroupColorMap.get(roomGroupId)!];
  }

  // 새로운 색상 배정
  roomGroupColorMap.set(roomGroupId, colorIndex % colorSet.length);
  return colorSet[colorIndex++ % colorSet.length];
}

function getTypeColor(type: string) {
  switch (type) {
    case "입실":
      return "positive";
    case "퇴실":
      return "negative";
    case "연박":
      return "info";
    default:
      return "grey";
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

interface CalendarEvent {
  id: number;
  name: string;
  peopleCount: number;
  rooms: string[];
  start: string;
  end: string;
  bgcolor: string;
  details: string;
}

interface DisplayedEvent {
  id?: number;
  left?: number;
  right?: number;
  size: number;
  event?: CalendarEvent;
}

function getWeekEvents(week: Timestamp[]): DisplayedEvent[] {
  if (!week || week.length === 0) return [];

  // Define week range
  const firstDay = parsed(`${week[0]!.date} 00:00`);
  const lastDay = parsed(`${week[week.length - 1]?.date} 23:59`);
  if (!firstDay || !lastDay) return [];

  // Filter and process events
  const eventsWeek = calendarEvents.value
    .map((event, id) => {
      const startDate = parsed(`${event.start} 00:00`);
      const endDate = parsed(`${event.end} 23:59`);

      if (startDate && endDate && isOverlappingDates(startDate, endDate, firstDay, lastDay)) {
        const left = daysBetween(firstDay, startDate);
        const right = daysBetween(endDate, lastDay);
        return {
          id,
          left,
          right,
          size: week.length - (left + right),
          event,
        };
      }
      return null;
    })
    .filter(Boolean) as DisplayedEvent[]; // Remove null values

  // Sort and insert events into week structure
  const evts: DisplayedEvent[] = [];
  if (eventsWeek.length > 0) {
    const sortedWeek = eventsWeek.sort((a, b) => (a.left ?? 0) - (b.left ?? 0));
    sortedWeek.forEach((_, i) => {
      insertEvent(evts, week.length, sortedWeek, i, 0, 0);
    });
  }

  return evts;
}

function insertEvent(
  events: DisplayedEvent[],
  weekLength: number,
  infoWeek: DisplayedEvent[],
  index: number,
  availableDays: number,
  level: number,
) {
  const iEvent = infoWeek[index];
  if (iEvent !== undefined && "left" in iEvent && iEvent.left >= availableDays) {
    // If you have space available, more events are placed
    if (iEvent.left - availableDays) {
      // It is filled with empty events
      events.push({ size: iEvent.left - availableDays });
    }
    // The event is built
    events.push({ size: iEvent.size, event: iEvent.event });

    if (level !== 0) {
      // If it goes into recursion, then the item is deleted
      infoWeek.splice(index, 1);
    }

    const currentAvailableDays = iEvent.left + iEvent.size;

    if (currentAvailableDays <= weekLength) {
      const indexNextEvent = indexOf(
        infoWeek,
        (e: DisplayedEvent) => e.id !== iEvent.id && e.left !== undefined && e.left >= currentAvailableDays,
      );

      insertEvent(
        events,
        weekLength,
        infoWeek,
        indexNextEvent !== -1 ? indexNextEvent : index,
        currentAvailableDays,
        level + 1,
      );
    } // else: There are no more days available, end of iteration
  } else {
    events.push({ size: weekLength - availableDays });
    // end of iteration
  }
}

function badgeClasses(displayedEvent: DisplayedEvent) {
  if (displayedEvent.event !== undefined) {
    return {
      "my-event": true,
      "text-white": true,
      [`bg-${displayedEvent.event.bgcolor}`]: true,
      [`text-${displayedEvent.event.bgcolor}`]: false,
      "rounded-border": true,
    };
  }
  return {
    "my-void-event": true,
  };
}

function badgeStyles(displayedEvent: DisplayedEvent, weekLength: number) {
  const s: Record<string, string | number> = {};
  if (displayedEvent.size !== undefined) {
    s.width = (100 / weekLength) * displayedEvent.size + "%";
  }
  return s;
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

<style lang="scss" scoped>
.my-event {
  position: relative;
  display: inline-flex;
  white-space: nowrap;
  font-size: 12px;
  height: 20px;
  margin: 1px 0 0 0;
  padding: 2px 2px;
  justify-content: start;
  text-overflow: ellipsis;
  overflow: hidden;
  cursor: pointer;
}

.my-void-event {
  display: inline-flex;
  white-space: nowrap;
  height: 1px;
}

.rounded-border {
  border-radius: 6px;
}
</style>
