import { describe, it, expect } from "vitest";
import dayjs from "dayjs";
import { useReservationCalendar, CalendarEvent, DisplayedEvent } from "../useReservationCalendar";
import type { Reservation } from "src/schema/reservation";
import type { Timestamp } from "@quasar/quasar-ui-qcalendar";

const createMockReservation = (
  id: number,
  stayStartAt: string,
  stayEndAt: string,
  roomGroupId = 1
): Reservation =>
  ({
    id,
    stayStartAt,
    stayEndAt,
    price: 100000,
    paymentAmount: 80000,
    rooms: [
      {
        id: 1,
        number: "101",
        roomGroup: { id: roomGroupId, name: "Standard" },
      },
    ],
    peopleCount: 2,
    name: `Guest ${id}`,
  }) as unknown as Reservation;

describe("useReservationCalendar", () => {
  const {
    formatReservations,
    getDateArray,
    calculateExtendedDateRange,
    badgeClasses,
    badgeStyles,
    getRoomGroupColor,
    getTypeColor,
  } = useReservationCalendar();

  describe("getDateArray", () => {
    it("시작일과 종료일 사이의 모든 날짜를 반환한다", () => {
      const dates = getDateArray("2026-01-01", "2026-01-03");

      expect(dates).toEqual(["2026-01-01", "2026-01-02", "2026-01-03"]);
    });

    it("하루만 있는 경우 하나의 날짜만 반환한다", () => {
      const dates = getDateArray("2026-01-15", "2026-01-15");

      expect(dates).toEqual(["2026-01-15"]);
    });

    it("월을 넘어가는 날짜 범위를 올바르게 처리한다", () => {
      const dates = getDateArray("2026-01-30", "2026-02-02");

      expect(dates).toEqual([
        "2026-01-30",
        "2026-01-31",
        "2026-02-01",
        "2026-02-02",
      ]);
    });
  });

  describe("formatReservations", () => {
    it("예약을 날짜별로 그룹화한다", () => {
      const reservations = [
        createMockReservation(1, "2026-01-01", "2026-01-03"),
      ];

      const result = formatReservations(reservations);

      expect(Object.keys(result)).toContain("2026-01-01");
      expect(Object.keys(result)).toContain("2026-01-02");
      expect(Object.keys(result)).toContain("2026-01-03");
    });

    it("입실일에는 type을 '입실'로 설정한다", () => {
      const reservations = [
        createMockReservation(1, "2026-01-01", "2026-01-03"),
      ];

      const result = formatReservations(reservations);

      expect(result["2026-01-01"][0].type).toBe("입실");
    });

    it("퇴실일에는 type을 '퇴실'로 설정한다", () => {
      const reservations = [
        createMockReservation(1, "2026-01-01", "2026-01-03"),
      ];

      const result = formatReservations(reservations);

      expect(result["2026-01-03"][0].type).toBe("퇴실");
    });

    it("중간 날짜에는 type을 '연박'으로 설정한다", () => {
      const reservations = [
        createMockReservation(1, "2026-01-01", "2026-01-03"),
      ];

      const result = formatReservations(reservations);

      expect(result["2026-01-02"][0].type).toBe("연박");
    });

    it("missPrice를 계산한다", () => {
      const reservations = [
        createMockReservation(1, "2026-01-01", "2026-01-02"),
      ];

      const result = formatReservations(reservations);

      expect(result["2026-01-01"][0].missPrice).toBe(20000);
    });
  });

  describe("calculateExtendedDateRange", () => {
    it("월요일로 시작하는 달의 시작일을 반환한다", () => {
      const month = dayjs("2026-06-01");
      const result = calculateExtendedDateRange(month);

      expect(result.startAt).toBe("2026-06-01");
    });

    it("일요일로 끝나는 주의 종료일을 반환한다", () => {
      const month = dayjs("2026-06-01");
      const result = calculateExtendedDateRange(month);

      expect(dayjs(result.endAt).day()).toBe(0);
    });

    it("월이 화요일로 시작하면 이전 월요일을 반환한다", () => {
      const month = dayjs("2026-09-01");
      const result = calculateExtendedDateRange(month);

      expect(dayjs(result.startAt).day()).toBe(1);
      expect(result.startAt).toBe("2026-08-31");
    });

    it("월이 일요일로 시작하면 6일 전 월요일을 반환한다", () => {
      const month = dayjs("2026-03-01");
      const result = calculateExtendedDateRange(month);

      expect(dayjs(result.startAt).day()).toBe(1);
      expect(result.startAt).toBe("2026-02-23");
    });
  });

  describe("badgeClasses", () => {
    it("이벤트가 있으면 이벤트 관련 클래스를 반환한다", () => {
      const displayedEvent: DisplayedEvent = {
        size: 3,
        event: {
          id: 1,
          name: "Test",
          peopleCount: 2,
          rooms: ["101"],
          start: "2026-01-01",
          end: "2026-01-03",
          bgcolor: "blue",
          details: "",
        },
      };

      const classes = badgeClasses(displayedEvent);

      expect(classes["my-event"]).toBe(true);
      expect(classes["text-white"]).toBe(true);
      expect(classes["bg-blue"]).toBe(true);
    });

    it("이벤트가 없으면 void 클래스를 반환한다", () => {
      const displayedEvent: DisplayedEvent = { size: 1 };

      const classes = badgeClasses(displayedEvent);

      expect(classes["my-void-event"]).toBe(true);
    });
  });

  describe("badgeStyles", () => {
    it("size에 따라 width를 계산한다", () => {
      const displayedEvent: DisplayedEvent = { size: 3 };
      const weekLength = 7;

      const styles = badgeStyles(displayedEvent, weekLength);

      expect(styles.width).toBe(`${(100 / 7) * 3}%`);
    });
  });

  describe("getRoomGroupColor", () => {
    it("객실이 없으면 grey-7을 반환한다", () => {
      const reservation = { rooms: [] } as unknown as Reservation;

      const color = getRoomGroupColor(reservation);

      expect(color).toBe("grey-7");
    });

    it("같은 객실 그룹에는 동일한 색상을 반환한다", () => {
      const reservation1 = createMockReservation(1, "2026-01-01", "2026-01-02", 100);
      const reservation2 = createMockReservation(2, "2026-01-03", "2026-01-04", 100);

      const color1 = getRoomGroupColor(reservation1);
      const color2 = getRoomGroupColor(reservation2);

      expect(color1).toBe(color2);
    });
  });

  describe("getTypeColor", () => {
    it("입실 타입에 positive 색상을 반환한다", () => {
      expect(getTypeColor("입실")).toBe("positive");
    });

    it("퇴실 타입에 negative 색상을 반환한다", () => {
      expect(getTypeColor("퇴실")).toBe("negative");
    });

    it("연박 타입에 info 색상을 반환한다", () => {
      expect(getTypeColor("연박")).toBe("info");
    });

    it("알 수 없는 타입에 grey 색상을 반환한다", () => {
      expect(getTypeColor("unknown")).toBe("grey");
    });
  });
});
