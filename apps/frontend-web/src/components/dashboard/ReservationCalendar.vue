<template>
  <div class="row">
    <div class="col q-pa-md-sm">
      <QCalendarMonth
        ref="calendar"
        :model-value="selectedDate"
        @navigation="handleNavigation"
        @update:model-value="handleDateSelect"
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
      >
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
import { ref } from "vue";
import { QCalendarMonth } from "@quasar/quasar-ui-qcalendar";
import "@quasar/quasar-ui-qcalendar/dist/index.css";
import { formatDate } from "src/util/format-util";
import { useReservationCalendar, CalendarEvent } from "src/composables/useReservationCalendar";

interface Props {
  calendarEvents: CalendarEvent[];
  selectedDate: string;
  calendarYear: number;
  calendarMonth: number;
}

interface Emits {
  (e: "navigation", view: { year: number; month: number }): void;
  (e: "dateSelect", date: string): void;
}

defineProps<Props>();
const emit = defineEmits<Emits>();

const calendar = ref<QCalendarMonth>();
const { getWeekEvents, badgeClasses, badgeStyles } = useReservationCalendar();

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
