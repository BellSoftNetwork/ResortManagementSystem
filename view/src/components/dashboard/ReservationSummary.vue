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
        <q-btn icon="chevron_left" color="primary" flat round dense @click="navigatePreviousMonth" />
        <div class="text-h6">{{ currentMonthLabel }}</div>
        <q-btn icon="chevron_right" color="primary" flat round dense @click="navigateNextMonth" />
      </div>

      <div class="row">
        <div class="col q-pa-md-sm">
          <QCalendarMonth
            ref="calendar"
            v-model="selectedDate"
            @navigation="changeView"
            :hoverable="true"
            :focusable="true"
            :focus-type="['day', 'date']"
            :now="formatDate"
            mask="YYYY-MM-DD"
            locale="ko-KR"
            day-min-height="60"
            :day-height="0"
            animated
            bordered
          >
            <template #week="{ scope: { week, weekdays } }">
              <template v-for="(displayedEvent, index) in getWeekEvents(week, weekdays)" :key="index">
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

      <div class="row">
        <div class="col q-pa-md-sm">
          <q-tab-panels v-model="selectedDate">
            <q-tab-panel
              v-for="(reservations, stayStartAt) in reservationsOfDay"
              :key="stayStartAt"
              :name="stayStartAt"
              class="q-px-none"
            >
              <q-table
                :columns="columns"
                :rows="reservations"
                row-key="id"
                flat
                bordered
                :pagination="{ rowsPerPage: 20, sortBy: 'type', descending: false }"
              >
                <template #top>
                  <div class="row justify-between items-center q-pa-sm full-width">
                    <div class="text-h6">{{ stayStartAt }} 예약 현황</div>
                    <div class="row q-gutter-sm">
                      <q-badge v-if="checkInOutCounts[stayStartAt]?.checkIn > 0" color="positive" class="q-px-sm">
                        입실 {{ checkInOutCounts[stayStartAt]?.checkIn }}
                      </q-badge>
                      <q-badge v-if="checkInOutCounts[stayStartAt]?.checkOut > 0" color="negative" class="q-px-sm">
                        퇴실 {{ checkInOutCounts[stayStartAt]?.checkOut }}
                      </q-badge>
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
              </q-table>
            </q-tab-panel>
          </q-tab-panels>
        </div>
      </div>
    </q-card-section>
  </q-card>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
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

function changeView(view) {
  const year = view.year;
  const month = view.month;

  currentMonth.value = dayjs(`${year}-${month}-01`);
  filter.value.stayStartAt = currentMonth.value.startOf("month").format("YYYY-MM-DD");
  filter.value.stayEndAt = currentMonth.value.endOf("month").format("YYYY-MM-DD");

  fetchData();
}

function navigatePreviousMonth() {
  currentMonth.value = currentMonth.value.subtract(1, "month");
  filter.value.stayStartAt = currentMonth.value.startOf("month").format("YYYY-MM-DD");
  filter.value.stayEndAt = currentMonth.value.endOf("month").format("YYYY-MM-DD");
  fetchData();
}

function navigateNextMonth() {
  currentMonth.value = currentMonth.value.add(1, "month");
  filter.value.stayStartAt = currentMonth.value.startOf("month").format("YYYY-MM-DD");
  filter.value.stayEndAt = currentMonth.value.endOf("month").format("YYYY-MM-DD");
  fetchData();
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

onMounted(() => {
  fetchData();
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
