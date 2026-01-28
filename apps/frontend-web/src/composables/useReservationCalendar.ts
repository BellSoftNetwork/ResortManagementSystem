import dayjs, { Dayjs } from "dayjs";
import { parsed, Timestamp, daysBetween, isOverlappingDates } from "@quasar/quasar-ui-qcalendar";
import { indexOf } from "@quasar/quasar-ui-qcalendar/src/utils/helpers.js";
import { formatDate } from "src/util/format-util";
import { Reservation } from "src/schema/reservation";

export interface CalendarEvent {
  id: number;
  name: string;
  peopleCount: number;
  rooms: string[];
  start: string;
  end: string;
  bgcolor: string;
  details: string;
}

export interface DisplayedEvent {
  id?: number;
  left?: number;
  right?: number;
  size: number;
  event?: CalendarEvent;
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

export function useReservationCalendar() {
  /**
   * 예약 목록을 날짜별로 그룹화하고 타입(입실/퇴실/연박)을 설정
   */
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

  /**
   * 시작일과 종료일 사이의 모든 날짜 배열 생성
   */
  function getDateArray(startDate: string, endDate: string) {
    const stayStartDate = formatDate(startDate);
    const stayEndAt = formatDate(endDate);
    const dateArray = [];

    for (let date = stayStartDate; date <= stayEndAt; date = dayjs(date).add(1, "day").format("YYYY-MM-DD")) {
      dateArray.push(date);
    }

    return dateArray;
  }

  /**
   * 달력에 표시되는 확장된 날짜 범위를 계산하는 함수
   * 월요일 시작 ~ 일요일 종료로 확장
   */
  function calculateExtendedDateRange(month: Dayjs) {
    // 해당 월의 시작일과 마지막일
    const monthStart = month.startOf("month");
    const monthEnd = month.endOf("month");

    // 월요일(1)을 기준으로 시작 날짜 계산
    // day()가 1(월요일)이면 0, 0(일요일)이면 -6, 나머지는 해당 요일만큼 빼기
    let daysToSubtract = 0;
    const firstDayOfMonth = monthStart.day();

    if (firstDayOfMonth === 0) {
      // 일요일인 경우 6일 전으로
      daysToSubtract = 6;
    } else if (firstDayOfMonth !== 1) {
      // 월요일이 아닌 경우 (화~토), 해당 요일에서 1 빼기
      daysToSubtract = firstDayOfMonth - 1;
    }
    // 월요일인 경우 daysToSubtract는 0 유지

    const startOfFirstWeek = monthStart.subtract(daysToSubtract, "day");

    // 일요일(0)을 기준으로 종료 날짜 계산
    // day()가 0(일요일)이면 0, 나머지는 (7 - 해당 요일) 더하기
    const daysToAdd = monthEnd.day() === 0 ? 0 : 7 - monthEnd.day();
    const endOfLastWeek = monthEnd.add(daysToAdd, "day");

    return {
      startAt: startOfFirstWeek.format("YYYY-MM-DD"),
      endAt: endOfLastWeek.format("YYYY-MM-DD"),
    };
  }

  /**
   * 주어진 주의 모든 이벤트를 계산하여 표시 가능한 형태로 반환
   */
  function getWeekEvents(week: Timestamp[], calendarEvents: CalendarEvent[]): DisplayedEvent[] {
    if (!week || week.length === 0) return [];

    // Define week range
    const firstDay = parsed(`${week[0]!.date} 00:00`);
    const lastDay = parsed(`${week[week.length - 1]?.date} 23:59`);
    if (!firstDay || !lastDay) return [];

    // Filter and process events
    const eventsWeek = calendarEvents
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

  /**
   * 이벤트를 주 구조에 재귀적으로 삽입
   */
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

  /**
   * 표시된 이벤트의 CSS 클래스 반환
   */
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

  /**
   * 표시된 이벤트의 CSS 스타일 반환
   */
  function badgeStyles(displayedEvent: DisplayedEvent, weekLength: number) {
    const s: Record<string, string | number> = {};
    if (displayedEvent.size !== undefined) {
      s.width = (100 / weekLength) * displayedEvent.size + "%";
    }
    return s;
  }

  /**
   * 객실 그룹 ID에 기반한 색상 반환
   */
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

  /**
   * 예약 타입에 따른 색상 반환
   */
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

  return {
    formatReservations,
    getDateArray,
    calculateExtendedDateRange,
    getWeekEvents,
    insertEvent,
    badgeClasses,
    badgeStyles,
    getRoomGroupColor,
    getTypeColor,
  };
}
