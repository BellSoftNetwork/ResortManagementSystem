<template>
  <div class="row">
    <div class="col q-pa-md-sm">
      <QCalendarMonth
        ref="calendar"
        :model-value="selectedDate"
        @navigation="handleNavigation"
        @update:model-value="handleDateSelect"
        @click-day="handleDayClick"
        @click-date="handleDateClick"
        :selected-dates="selectedDatesArray"
        :year="calendarYear"
        :month="calendarMonth"
        :hoverable="true"
        :focusable="false"
        :focus-type="['day', 'date']"
        :now="formatDate()"
        mask="YYYY-MM-DD"
        locale="ko-KR"
        :weekdays="[1, 2, 3, 4, 5, 6, 0]"
        day-min-height="150"
        :day-height="0"
        animated
        bordered
        :day-class="getDayClass"
      >
        <template #head-day-label="{ scope: { timestamp, dayLabel } }">
          <span>{{ dayLabel }}</span>
          <q-icon v-if="blockedDates?.has(timestamp.date)" name="event_busy" color="red-4" size="xs" class="q-ml-xs" />
        </template>
        <template #week="{ scope: { week } }">
          <template v-for="(displayedEvent, index) in getWeekEvents(week, calendarEvents)" :key="index">
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
</template>

<script setup lang="ts">
import { ref, computed } from "vue";
import { QCalendarMonth } from "@quasar/quasar-ui-qcalendar";
import "@quasar/quasar-ui-qcalendar/dist/index.css";
import { formatDate } from "src/util/format-util";
import { useReservationCalendar, CalendarEvent } from "src/composables/useReservationCalendar";

interface Props {
  calendarEvents: CalendarEvent[];
  selectedDate: string;
  calendarYear: number;
  calendarMonth: number;
  blockedDates?: Set<string>;
}

interface Emits {
  (e: "navigation", view: { year: number; month: number }): void;
  (e: "dateSelect", date: string): void;
}

const props = defineProps<Props>();
const emit = defineEmits<Emits>();

const calendar = ref<QCalendarMonth>();
const { getWeekEvents, badgeClasses, badgeStyles } = useReservationCalendar();

// 선택된 날짜 배열 (QCalendarMonth selected-dates prop용)
const selectedDatesArray = computed(() => [props.selectedDate]);

// 날짜 셀(day 영역) 클릭 핸들러
// @click-day payload: { scope: { timestamp: { date: string } }, e: MouseEvent }
function handleDayClick({ scope }: { scope: { timestamp: { date: string } } }) {
  emit("dateSelect", scope.timestamp.date);
}

// 날짜 번호(head 영역) 클릭 핸들러
// @click-date payload: { scope: { timestamp: { date: string } }, event: MouseEvent }
function handleDateClick({ scope }: { scope: { timestamp: { date: string } } }) {
  emit("dateSelect", scope.timestamp.date);
}

function getDayClass({ scope }: { scope: { timestamp: { date: string } } }) {
  return {
    "blocked-date": !!props.blockedDates?.has(scope.timestamp.date),
    "selected-date": scope.timestamp.date === props.selectedDate,
  };
}

function handleNavigation(view: { year: number; month: number }) {
  emit("navigation", view);
}

function handleDateSelect(date: string) {
  emit("dateSelect", date);
}

// Expose calendar ref for parent component to call prev/next
defineExpose({
  prev: () => calendar.value?.prev(),
  next: () => calendar.value?.next(),
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
  pointer-events: auto;
}

.my-void-event {
  display: inline-flex;
  white-space: nowrap;
  height: 1px;
}

.rounded-border {
  border-radius: 6px;
}

:deep(.blocked-date) {
  background-color: rgba(239, 83, 80, 0.12) !important;
}

:deep(.selected-date) {
  background-color: rgba(25, 118, 210, 0.15) !important;
  border-left: 3px solid #1976d2 !important;
}
</style>

<!-- QCalendar 내부 DOM 요소에 대한 CSS 오버라이드.
  .q-calendar-month__week--events는 QCalendarMonth 자식 컴포넌트가 렌더링하는 요소이므로
  Vue 3 scoped CSS의 :deep()보다 unscoped 블록으로 선언하는 것이 빌드 환경 안정성을 보장합니다.
  See: https://vuejs.org/api/sfc-css-features#deep-selectors -->
<style lang="scss">
/* 이벤트 오버레이 레이어의 클릭 투과 처리.
   이 레이어(position: absolute)가 pointer-events: auto이면 하위 날짜 셀 클릭이 차단됩니다. */
.q-calendar-month__week--events {
  pointer-events: none;
}
</style>
