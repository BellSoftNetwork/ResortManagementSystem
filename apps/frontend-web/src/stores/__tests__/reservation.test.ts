import { describe, it, expect, vi, beforeEach } from "vitest";
import { setActivePinia, createPinia } from "pinia";
import { useReservationStore } from "../reservation";
import * as reservationApi from "src/api/v1/reservation";
import { createMockApiResponse, createMockApiError } from "test/vitest/helpers";
import type { Reservation, ReservationStatus, ReservationType } from "src/schema/reservation";
import type { Revision } from "src/schema/revision";
import type { ReservationStatistics, StatisticsPeriodType } from "src/api/v1/reservation";

vi.mock("src/api/v1/reservation");

describe("useReservationStore", () => {
  beforeEach(() => {
    setActivePinia(createPinia());
    vi.clearAllMocks();
  });

  // Mock 데이터 생성 헬퍼
  const createMockReservation = (
    id: number,
    guestName: string,
    status: ReservationStatus = "CONFIRMED",
    type: ReservationType = "ONLINE"
  ): Reservation => ({
    id,
    guestName,
    guestPhone: "010-1234-5678",
    guestEmail: "guest@example.com",
    roomId: 1,
    roomNumber: "101",
    stayStartAt: "2026-02-01",
    stayEndAt: "2026-02-03",
    status,
    type,
    totalPrice: 200000,
    note: "테스트 예약",
    createdAt: "2026-01-01T00:00:00Z",
    updatedAt: "2026-01-01T00:00:00Z",
    createdBy: "admin",
    updatedBy: "admin",
  });

  const createMockRevision = (reservation: Reservation): Revision<Reservation> => ({
    entity: reservation,
    historyType: "UPDATED",
    historyCreatedAt: "2026-01-26T00:00:00Z",
    updatedFields: ["status"],
  });

  const createMockStatistics = (): ReservationStatistics => ({
    totalReservations: 100,
    totalRevenue: 10000000,
    averageStayDuration: 2.5,
    periodData: [
      {
        period: "2026-01",
        reservationCount: 50,
        revenue: 5000000,
      },
      {
        period: "2026-02",
        reservationCount: 50,
        revenue: 5000000,
      },
    ],
  });

  describe("fetchReservationList", () => {
    it("예약 목록을 성공적으로 조회한다", async () => {
      const mockReservations = [
        createMockReservation(1, "홍길동"),
        createMockReservation(2, "김철수"),
      ];
      vi.mocked(reservationApi.fetchReservations).mockResolvedValue(
        createMockApiResponse(mockReservations, 2)
      );

      const store = useReservationStore();
      await store.fetchReservationList();

      expect(store.reservations).toEqual(mockReservations);
      expect(store.totalCount).toBe(2);
      expect(store.loading).toBe(false);
      expect(store.error).toBeNull();
    });

    it("필터 파라미터와 함께 예약 목록을 조회한다", async () => {
      const mockReservations = [createMockReservation(1, "홍길동", "CONFIRMED")];
      vi.mocked(reservationApi.fetchReservations).mockResolvedValue(
        createMockApiResponse(mockReservations, 1)
      );

      const store = useReservationStore();
      const filter = {
        status: "CONFIRMED" as ReservationStatus,
        stayStartAt: "2026-02-01",
        stayEndAt: "2026-02-28",
        page: 1,
        size: 10,
      };
      await store.fetchReservationList(filter);

      expect(reservationApi.fetchReservations).toHaveBeenCalledWith(filter);
      expect(store.reservations).toEqual(mockReservations);
    });

    it("예약 목록 조회 중 로딩 상태를 관리한다", async () => {
      const mockReservations = [createMockReservation(1, "홍길동")];
      vi.mocked(reservationApi.fetchReservations).mockImplementation(
        () =>
          new Promise((resolve) => {
            setTimeout(() => resolve(createMockApiResponse(mockReservations, 1)), 10);
          })
      );

      const store = useReservationStore();
      const promise = store.fetchReservationList();

      expect(store.loading).toBe(true);

      await promise;

      expect(store.loading).toBe(false);
    });

    it("예약 목록 조회 실패 시 에러를 처리한다", async () => {
      const errorMessage = "Failed to fetch reservations";
      vi.mocked(reservationApi.fetchReservations).mockRejectedValue(
        createMockApiError(errorMessage)
      );

      const store = useReservationStore();

      await expect(store.fetchReservationList()).rejects.toThrow(errorMessage);
      expect(store.error).toBe(errorMessage);
      expect(store.loading).toBe(false);
    });
  });

  describe("fetchReservationById", () => {
    it("ID로 예약을 성공적으로 조회한다", async () => {
      const mockReservation = createMockReservation(1, "홍길동");
      vi.mocked(reservationApi.fetchReservation).mockResolvedValue(
        createMockApiResponse(mockReservation)
      );

      const store = useReservationStore();
      await store.fetchReservationById(1);

      expect(store.currentReservation).toEqual(mockReservation);
      expect(store.loading).toBe(false);
      expect(store.error).toBeNull();
    });

    it("예약 조회 중 로딩 상태를 관리한다", async () => {
      const mockReservation = createMockReservation(1, "홍길동");
      vi.mocked(reservationApi.fetchReservation).mockImplementation(
        () =>
          new Promise((resolve) => {
            setTimeout(() => resolve(createMockApiResponse(mockReservation)), 10);
          })
      );

      const store = useReservationStore();
      const promise = store.fetchReservationById(1);

      expect(store.loading).toBe(true);

      await promise;

      expect(store.loading).toBe(false);
    });

    it("예약 조회 실패 시 에러를 처리한다", async () => {
      const errorMessage = "Failed to fetch reservation";
      vi.mocked(reservationApi.fetchReservation).mockRejectedValue(
        createMockApiError(errorMessage)
      );

      const store = useReservationStore();

      await expect(store.fetchReservationById(1)).rejects.toThrow(errorMessage);
      expect(store.error).toBe(errorMessage);
      expect(store.loading).toBe(false);
    });
  });

  describe("createNewReservation", () => {
    it("새 예약을 성공적으로 생성한다", async () => {
      const newReservationData = {
        guestName: "이영희",
        guestPhone: "010-9876-5432",
        guestEmail: "lee@example.com",
        roomId: 2,
        stayStartAt: "2026-03-01",
        stayEndAt: "2026-03-03",
        status: "CONFIRMED" as ReservationStatus,
        type: "ONLINE" as ReservationType,
        totalPrice: 300000,
        note: "새 예약",
      };
      const createdReservation = createMockReservation(3, "이영희");
      vi.mocked(reservationApi.createReservation).mockResolvedValue(
        createMockApiResponse(createdReservation)
      );

      const store = useReservationStore();
      store.reservations = [createMockReservation(1, "홍길동"), createMockReservation(2, "김철수")];
      store.totalCount = 2;

      await store.createNewReservation(newReservationData);

      expect(store.reservations).toHaveLength(3);
      expect(store.reservations[2]).toEqual(createdReservation);
      expect(store.totalCount).toBe(3);
      expect(store.loading).toBe(false);
      expect(store.error).toBeNull();
    });

    it("예약 생성 중 로딩 상태를 관리한다", async () => {
      const newReservationData = {
        guestName: "이영희",
        guestPhone: "010-9876-5432",
        guestEmail: "lee@example.com",
        roomId: 2,
        stayStartAt: "2026-03-01",
        stayEndAt: "2026-03-03",
        status: "CONFIRMED" as ReservationStatus,
        type: "ONLINE" as ReservationType,
        totalPrice: 300000,
        note: "새 예약",
      };
      const createdReservation = createMockReservation(3, "이영희");
      vi.mocked(reservationApi.createReservation).mockImplementation(
        () =>
          new Promise((resolve) => {
            setTimeout(() => resolve(createMockApiResponse(createdReservation)), 10);
          })
      );

      const store = useReservationStore();
      const promise = store.createNewReservation(newReservationData);

      expect(store.loading).toBe(true);

      await promise;

      expect(store.loading).toBe(false);
    });

    it("예약 생성 실패 시 에러를 처리한다", async () => {
      const newReservationData = {
        guestName: "이영희",
        guestPhone: "010-9876-5432",
        guestEmail: "lee@example.com",
        roomId: 2,
        stayStartAt: "2026-03-01",
        stayEndAt: "2026-03-03",
        status: "CONFIRMED" as ReservationStatus,
        type: "ONLINE" as ReservationType,
        totalPrice: 300000,
        note: "새 예약",
      };
      const errorMessage = "Failed to create reservation";
      vi.mocked(reservationApi.createReservation).mockRejectedValue(
        createMockApiError(errorMessage)
      );

      const store = useReservationStore();

      await expect(store.createNewReservation(newReservationData)).rejects.toThrow(errorMessage);
      expect(store.error).toBe(errorMessage);
      expect(store.loading).toBe(false);
    });
  });

  describe("updateReservation", () => {
    it("예약을 성공적으로 수정한다", async () => {
      const updateData = {
        status: "CANCELLED" as ReservationStatus,
      };
      const updatedReservation = createMockReservation(1, "홍길동", "CANCELLED");
      vi.mocked(reservationApi.patchReservation).mockResolvedValue(
        createMockApiResponse(updatedReservation)
      );

      const store = useReservationStore();
      store.reservations = [createMockReservation(1, "홍길동"), createMockReservation(2, "김철수")];

      await store.updateReservation(1, updateData);

      expect(store.reservations[0]).toEqual(updatedReservation);
      expect(store.loading).toBe(false);
      expect(store.error).toBeNull();
    });

    it("예약 수정 시 currentReservation도 함께 업데이트한다", async () => {
      const updateData = {
        status: "CANCELLED" as ReservationStatus,
      };
      const updatedReservation = createMockReservation(1, "홍길동", "CANCELLED");
      vi.mocked(reservationApi.patchReservation).mockResolvedValue(
        createMockApiResponse(updatedReservation)
      );

      const store = useReservationStore();
      store.reservations = [createMockReservation(1, "홍길동"), createMockReservation(2, "김철수")];
      store.currentReservation = createMockReservation(1, "홍길동");

      await store.updateReservation(1, updateData);

      expect(store.currentReservation).toEqual(updatedReservation);
    });

    it("목록에 없는 예약을 수정해도 에러가 발생하지 않는다", async () => {
      const updateData = {
        status: "CANCELLED" as ReservationStatus,
      };
      const updatedReservation = createMockReservation(3, "이영희", "CANCELLED");
      vi.mocked(reservationApi.patchReservation).mockResolvedValue(
        createMockApiResponse(updatedReservation)
      );

      const store = useReservationStore();
      store.reservations = [createMockReservation(1, "홍길동"), createMockReservation(2, "김철수")];

      await store.updateReservation(3, updateData);

      expect(store.reservations).toHaveLength(2);
      expect(store.loading).toBe(false);
    });

    it("예약 수정 중 로딩 상태를 관리한다", async () => {
      const updateData = {
        status: "CANCELLED" as ReservationStatus,
      };
      const updatedReservation = createMockReservation(1, "홍길동", "CANCELLED");
      vi.mocked(reservationApi.patchReservation).mockImplementation(
        () =>
          new Promise((resolve) => {
            setTimeout(() => resolve(createMockApiResponse(updatedReservation)), 10);
          })
      );

      const store = useReservationStore();
      store.reservations = [createMockReservation(1, "홍길동")];
      const promise = store.updateReservation(1, updateData);

      expect(store.loading).toBe(true);

      await promise;

      expect(store.loading).toBe(false);
    });

    it("예약 수정 실패 시 에러를 처리한다", async () => {
      const updateData = {
        status: "CANCELLED" as ReservationStatus,
      };
      const errorMessage = "Failed to update reservation";
      vi.mocked(reservationApi.patchReservation).mockRejectedValue(
        createMockApiError(errorMessage)
      );

      const store = useReservationStore();

      await expect(store.updateReservation(1, updateData)).rejects.toThrow(errorMessage);
      expect(store.error).toBe(errorMessage);
      expect(store.loading).toBe(false);
    });
  });

  describe("removeReservation", () => {
    it("예약을 성공적으로 삭제한다", async () => {
      vi.mocked(reservationApi.deleteReservation).mockResolvedValue(undefined);

      const store = useReservationStore();
      store.reservations = [createMockReservation(1, "홍길동"), createMockReservation(2, "김철수")];
      store.totalCount = 2;

      await store.removeReservation(1);

      expect(store.reservations).toHaveLength(1);
      expect(store.reservations[0].id).toBe(2);
      expect(store.totalCount).toBe(1);
      expect(store.loading).toBe(false);
      expect(store.error).toBeNull();
    });

    it("예약 삭제 시 currentReservation도 함께 초기화한다", async () => {
      vi.mocked(reservationApi.deleteReservation).mockResolvedValue(undefined);

      const store = useReservationStore();
      store.reservations = [createMockReservation(1, "홍길동"), createMockReservation(2, "김철수")];
      store.currentReservation = createMockReservation(1, "홍길동");
      store.totalCount = 2;

      await store.removeReservation(1);

      expect(store.currentReservation).toBeNull();
    });

    it("다른 예약을 삭제해도 currentReservation은 유지된다", async () => {
      vi.mocked(reservationApi.deleteReservation).mockResolvedValue(undefined);

      const store = useReservationStore();
      store.reservations = [createMockReservation(1, "홍길동"), createMockReservation(2, "김철수")];
      store.currentReservation = createMockReservation(1, "홍길동");
      store.totalCount = 2;

      await store.removeReservation(2);

      expect(store.currentReservation).toEqual(createMockReservation(1, "홍길동"));
    });

    it("예약 삭제 중 로딩 상태를 관리한다", async () => {
      vi.mocked(reservationApi.deleteReservation).mockImplementation(
        () =>
          new Promise((resolve) => {
            setTimeout(() => resolve(undefined), 10);
          })
      );

      const store = useReservationStore();
      store.reservations = [createMockReservation(1, "홍길동")];
      const promise = store.removeReservation(1);

      expect(store.loading).toBe(true);

      await promise;

      expect(store.loading).toBe(false);
    });

    it("예약 삭제 실패 시 에러를 처리한다", async () => {
      const errorMessage = "Failed to delete reservation";
      vi.mocked(reservationApi.deleteReservation).mockRejectedValue(
        createMockApiError(errorMessage)
      );

      const store = useReservationStore();

      await expect(store.removeReservation(1)).rejects.toThrow(errorMessage);
      expect(store.error).toBe(errorMessage);
      expect(store.loading).toBe(false);
    });
  });

  describe("fetchReservationHistoryList", () => {
    it("예약 히스토리를 성공적으로 조회한다", async () => {
      const mockReservation = createMockReservation(1, "홍길동");
      const mockHistories = [createMockRevision(mockReservation)];
      vi.mocked(reservationApi.fetchReservationHistories).mockResolvedValue(
        createMockApiResponse(mockHistories, 1)
      );

      const store = useReservationStore();
      await store.fetchReservationHistoryList(1);

      expect(store.reservationHistories).toEqual(mockHistories);
      expect(store.loading).toBe(false);
      expect(store.error).toBeNull();
    });

    it("페이지네이션 파라미터와 함께 히스토리를 조회한다", async () => {
      const mockReservation = createMockReservation(1, "홍길동");
      const mockHistories = [createMockRevision(mockReservation)];
      vi.mocked(reservationApi.fetchReservationHistories).mockResolvedValue(
        createMockApiResponse(mockHistories, 1)
      );

      const store = useReservationStore();
      const params = {
        page: 1,
        size: 10,
        sort: "historyCreatedAt,desc",
      };
      await store.fetchReservationHistoryList(1, params);

      expect(reservationApi.fetchReservationHistories).toHaveBeenCalledWith(1, params);
      expect(store.reservationHistories).toEqual(mockHistories);
    });

    it("히스토리 조회 중 로딩 상태를 관리한다", async () => {
      const mockReservation = createMockReservation(1, "홍길동");
      const mockHistories = [createMockRevision(mockReservation)];
      vi.mocked(reservationApi.fetchReservationHistories).mockImplementation(
        () =>
          new Promise((resolve) => {
            setTimeout(() => resolve(createMockApiResponse(mockHistories, 1)), 10);
          })
      );

      const store = useReservationStore();
      const promise = store.fetchReservationHistoryList(1);

      expect(store.loading).toBe(true);

      await promise;

      expect(store.loading).toBe(false);
    });

    it("히스토리 조회 실패 시 에러를 처리한다", async () => {
      const errorMessage = "Failed to fetch reservation histories";
      vi.mocked(reservationApi.fetchReservationHistories).mockRejectedValue(
        createMockApiError(errorMessage)
      );

      const store = useReservationStore();

      await expect(store.fetchReservationHistoryList(1)).rejects.toThrow(errorMessage);
      expect(store.error).toBe(errorMessage);
      expect(store.loading).toBe(false);
    });
  });

  describe("fetchStatistics", () => {
    it("통계를 성공적으로 조회한다", async () => {
      const mockStatistics = createMockStatistics();
      vi.mocked(reservationApi.fetchReservationStatistics).mockResolvedValue(
        createMockApiResponse(mockStatistics)
      );

      const store = useReservationStore();
      await store.fetchStatistics("2026-01-01", "2026-02-28", "MONTHLY" as StatisticsPeriodType);

      expect(store.statistics).toEqual(mockStatistics);
      expect(store.loading).toBe(false);
      expect(store.error).toBeNull();
    });

    it("기본 periodType으로 통계를 조회한다", async () => {
      const mockStatistics = createMockStatistics();
      vi.mocked(reservationApi.fetchReservationStatistics).mockResolvedValue(
        createMockApiResponse(mockStatistics)
      );

      const store = useReservationStore();
      await store.fetchStatistics("2026-01-01", "2026-02-28");

      expect(reservationApi.fetchReservationStatistics).toHaveBeenCalledWith(
        "2026-01-01",
        "2026-02-28",
        "MONTHLY"
      );
      expect(store.statistics).toEqual(mockStatistics);
    });

    it("통계 조회 중 로딩 상태를 관리한다", async () => {
      const mockStatistics = createMockStatistics();
      vi.mocked(reservationApi.fetchReservationStatistics).mockImplementation(
        () =>
          new Promise((resolve) => {
            setTimeout(() => resolve(createMockApiResponse(mockStatistics)), 10);
          })
      );

      const store = useReservationStore();
      const promise = store.fetchStatistics("2026-01-01", "2026-02-28");

      expect(store.loading).toBe(true);

      await promise;

      expect(store.loading).toBe(false);
    });

    it("통계 조회 실패 시 에러를 처리한다", async () => {
      const errorMessage = "Failed to fetch reservation statistics";
      vi.mocked(reservationApi.fetchReservationStatistics).mockRejectedValue(
        createMockApiError(errorMessage)
      );

      const store = useReservationStore();

      await expect(
        store.fetchStatistics("2026-01-01", "2026-02-28", "MONTHLY" as StatisticsPeriodType)
      ).rejects.toThrow(errorMessage);
      expect(store.error).toBe(errorMessage);
      expect(store.loading).toBe(false);
    });
  });

  describe("clearError", () => {
    it("에러 상태를 초기화한다", () => {
      const store = useReservationStore();
      store.error = "Some error";

      store.clearError();

      expect(store.error).toBeNull();
    });
  });

  describe("clearCurrentReservation", () => {
    it("현재 예약 상태를 초기화한다", () => {
      const store = useReservationStore();
      store.currentReservation = createMockReservation(1, "홍길동");

      store.clearCurrentReservation();

      expect(store.currentReservation).toBeNull();
    });
  });

  describe("clearStatistics", () => {
    it("통계 상태를 초기화한다", () => {
      const store = useReservationStore();
      store.statistics = createMockStatistics();

      store.clearStatistics();

      expect(store.statistics).toBeNull();
    });
  });

  describe("초기 상태", () => {
    it("스토어가 올바른 초기 상태를 가진다", () => {
      const store = useReservationStore();

      expect(store.reservations).toEqual([]);
      expect(store.currentReservation).toBeNull();
      expect(store.reservationHistories).toEqual([]);
      expect(store.statistics).toBeNull();
      expect(store.loading).toBe(false);
      expect(store.error).toBeNull();
      expect(store.totalCount).toBe(0);
    });
  });
});
