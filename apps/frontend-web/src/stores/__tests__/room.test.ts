import { describe, it, expect, vi, beforeEach } from "vitest";
import { setActivePinia, createPinia } from "pinia";
import { useRoomStore } from "../room";
import * as roomApi from "src/api/v1/room";
import { createMockApiResponse, createMockApiError } from "test/vitest/helpers";
import type { Room, RoomStatus } from "src/schema/room";
import type { Revision } from "src/schema/revision";
import type { RoomGroup } from "src/schema/room-group";

vi.mock("src/api/v1/room");

describe("useRoomStore", () => {
  beforeEach(() => {
    setActivePinia(createPinia());
    vi.clearAllMocks();
  });

  // Mock 데이터 생성 헬퍼
  const createMockRoomGroup = (): RoomGroup => ({
    id: 1,
    name: "스탠다드",
    description: "일반 객실",
    createdAt: "2026-01-01T00:00:00Z",
    updatedAt: "2026-01-01T00:00:00Z",
    createdBy: "admin",
    updatedBy: "admin",
  });

  const createMockRoom = (id: number, number: string, status: RoomStatus = "NORMAL"): Room => ({
    id,
    number,
    roomGroup: createMockRoomGroup(),
    note: "테스트 객실",
    status,
    createdAt: "2026-01-01T00:00:00Z",
    updatedAt: "2026-01-01T00:00:00Z",
    createdBy: "admin",
    updatedBy: "admin",
  });

  const createMockRevision = (room: Room): Revision<Room> => ({
    entity: room,
    historyType: "UPDATED",
    historyCreatedAt: "2026-01-26T00:00:00Z",
    updatedFields: ["status"],
  });

  describe("fetchRoomList", () => {
    it("객실 목록을 성공적으로 조회한다", async () => {
      const mockRooms = [createMockRoom(1, "101"), createMockRoom(2, "102")];
      vi.mocked(roomApi.fetchRooms).mockResolvedValue(createMockApiResponse(mockRooms, 2));

      const store = useRoomStore();
      await store.fetchRoomList();

      expect(store.rooms).toEqual(mockRooms);
      expect(store.totalCount).toBe(2);
      expect(store.loading).toBe(false);
      expect(store.error).toBeNull();
    });

    it("필터 파라미터와 함께 객실 목록을 조회한다", async () => {
      const mockRooms = [createMockRoom(1, "101", "NORMAL")];
      vi.mocked(roomApi.fetchRooms).mockResolvedValue(createMockApiResponse(mockRooms, 1));

      const store = useRoomStore();
      const filter = {
        status: "NORMAL" as RoomStatus,
        page: 1,
        size: 10,
      };
      await store.fetchRoomList(filter);

      expect(roomApi.fetchRooms).toHaveBeenCalledWith(filter);
      expect(store.rooms).toEqual(mockRooms);
    });

    it("객실 목록 조회 중 로딩 상태를 관리한다", async () => {
      const mockRooms = [createMockRoom(1, "101")];
      vi.mocked(roomApi.fetchRooms).mockImplementation(
        () =>
          new Promise((resolve) => {
            setTimeout(() => resolve(createMockApiResponse(mockRooms, 1)), 10);
          })
      );

      const store = useRoomStore();
      const promise = store.fetchRoomList();

      expect(store.loading).toBe(true);

      await promise;

      expect(store.loading).toBe(false);
    });

    it("객실 목록 조회 실패 시 에러를 처리한다", async () => {
      const errorMessage = "Failed to fetch rooms";
      vi.mocked(roomApi.fetchRooms).mockRejectedValue(createMockApiError(errorMessage));

      const store = useRoomStore();

      await expect(store.fetchRoomList()).rejects.toThrow(errorMessage);
      expect(store.error).toBe(errorMessage);
      expect(store.loading).toBe(false);
    });
  });

  describe("fetchRoomById", () => {
    it("ID로 객실을 성공적으로 조회한다", async () => {
      const mockRoom = createMockRoom(1, "101");
      vi.mocked(roomApi.fetchRoom).mockResolvedValue(createMockApiResponse(mockRoom));

      const store = useRoomStore();
      await store.fetchRoomById(1);

      expect(store.currentRoom).toEqual(mockRoom);
      expect(store.loading).toBe(false);
      expect(store.error).toBeNull();
    });

    it("객실 조회 중 로딩 상태를 관리한다", async () => {
      const mockRoom = createMockRoom(1, "101");
      vi.mocked(roomApi.fetchRoom).mockImplementation(
        () =>
          new Promise((resolve) => {
            setTimeout(() => resolve(createMockApiResponse(mockRoom)), 10);
          })
      );

      const store = useRoomStore();
      const promise = store.fetchRoomById(1);

      expect(store.loading).toBe(true);

      await promise;

      expect(store.loading).toBe(false);
    });

    it("객실 조회 실패 시 에러를 처리한다", async () => {
      const errorMessage = "Failed to fetch room";
      vi.mocked(roomApi.fetchRoom).mockRejectedValue(createMockApiError(errorMessage));

      const store = useRoomStore();

      await expect(store.fetchRoomById(1)).rejects.toThrow(errorMessage);
      expect(store.error).toBe(errorMessage);
      expect(store.loading).toBe(false);
    });
  });

  describe("createNewRoom", () => {
    it("새 객실을 성공적으로 생성한다", async () => {
      const newRoomData = {
        roomGroup: { id: 1 },
        number: "103",
        note: "새 객실",
        status: "NORMAL" as RoomStatus,
      };
      const createdRoom = createMockRoom(3, "103");
      vi.mocked(roomApi.createRoom).mockResolvedValue(createMockApiResponse(createdRoom));

      const store = useRoomStore();
      store.rooms = [createMockRoom(1, "101"), createMockRoom(2, "102")];
      store.totalCount = 2;

      await store.createNewRoom(newRoomData);

      expect(store.rooms).toHaveLength(3);
      expect(store.rooms[2]).toEqual(createdRoom);
      expect(store.totalCount).toBe(3);
      expect(store.loading).toBe(false);
      expect(store.error).toBeNull();
    });

    it("객실 생성 중 로딩 상태를 관리한다", async () => {
      const newRoomData = {
        roomGroup: { id: 1 },
        number: "103",
        note: "새 객실",
        status: "NORMAL" as RoomStatus,
      };
      const createdRoom = createMockRoom(3, "103");
      vi.mocked(roomApi.createRoom).mockImplementation(
        () =>
          new Promise((resolve) => {
            setTimeout(() => resolve(createMockApiResponse(createdRoom)), 10);
          })
      );

      const store = useRoomStore();
      const promise = store.createNewRoom(newRoomData);

      expect(store.loading).toBe(true);

      await promise;

      expect(store.loading).toBe(false);
    });

    it("객실 생성 실패 시 에러를 처리한다", async () => {
      const newRoomData = {
        roomGroup: { id: 1 },
        number: "103",
        note: "새 객실",
        status: "NORMAL" as RoomStatus,
      };
      const errorMessage = "Failed to create room";
      vi.mocked(roomApi.createRoom).mockRejectedValue(createMockApiError(errorMessage));

      const store = useRoomStore();

      await expect(store.createNewRoom(newRoomData)).rejects.toThrow(errorMessage);
      expect(store.error).toBe(errorMessage);
      expect(store.loading).toBe(false);
    });
  });

  describe("updateRoom", () => {
    it("객실을 성공적으로 수정한다", async () => {
      const updateData = {
        status: "INACTIVE" as RoomStatus,
      };
      const updatedRoom = createMockRoom(1, "101", "INACTIVE");
      vi.mocked(roomApi.patchRoom).mockResolvedValue(createMockApiResponse(updatedRoom));

      const store = useRoomStore();
      store.rooms = [createMockRoom(1, "101"), createMockRoom(2, "102")];

      await store.updateRoom(1, updateData);

      expect(store.rooms[0]).toEqual(updatedRoom);
      expect(store.loading).toBe(false);
      expect(store.error).toBeNull();
    });

    it("객실 수정 시 currentRoom도 함께 업데이트한다", async () => {
      const updateData = {
        status: "INACTIVE" as RoomStatus,
      };
      const updatedRoom = createMockRoom(1, "101", "INACTIVE");
      vi.mocked(roomApi.patchRoom).mockResolvedValue(createMockApiResponse(updatedRoom));

      const store = useRoomStore();
      store.rooms = [createMockRoom(1, "101"), createMockRoom(2, "102")];
      store.currentRoom = createMockRoom(1, "101");

      await store.updateRoom(1, updateData);

      expect(store.currentRoom).toEqual(updatedRoom);
    });

    it("목록에 없는 객실을 수정해도 에러가 발생하지 않는다", async () => {
      const updateData = {
        status: "INACTIVE" as RoomStatus,
      };
      const updatedRoom = createMockRoom(3, "103", "INACTIVE");
      vi.mocked(roomApi.patchRoom).mockResolvedValue(createMockApiResponse(updatedRoom));

      const store = useRoomStore();
      store.rooms = [createMockRoom(1, "101"), createMockRoom(2, "102")];

      await store.updateRoom(3, updateData);

      expect(store.rooms).toHaveLength(2);
      expect(store.loading).toBe(false);
    });

    it("객실 수정 중 로딩 상태를 관리한다", async () => {
      const updateData = {
        status: "INACTIVE" as RoomStatus,
      };
      const updatedRoom = createMockRoom(1, "101", "INACTIVE");
      vi.mocked(roomApi.patchRoom).mockImplementation(
        () =>
          new Promise((resolve) => {
            setTimeout(() => resolve(createMockApiResponse(updatedRoom)), 10);
          })
      );

      const store = useRoomStore();
      store.rooms = [createMockRoom(1, "101")];
      const promise = store.updateRoom(1, updateData);

      expect(store.loading).toBe(true);

      await promise;

      expect(store.loading).toBe(false);
    });

    it("객실 수정 실패 시 에러를 처리한다", async () => {
      const updateData = {
        status: "INACTIVE" as RoomStatus,
      };
      const errorMessage = "Failed to update room";
      vi.mocked(roomApi.patchRoom).mockRejectedValue(createMockApiError(errorMessage));

      const store = useRoomStore();

      await expect(store.updateRoom(1, updateData)).rejects.toThrow(errorMessage);
      expect(store.error).toBe(errorMessage);
      expect(store.loading).toBe(false);
    });
  });

  describe("removeRoom", () => {
    it("객실을 성공적으로 삭제한다", async () => {
      vi.mocked(roomApi.deleteRoom).mockResolvedValue(undefined);

      const store = useRoomStore();
      store.rooms = [createMockRoom(1, "101"), createMockRoom(2, "102")];
      store.totalCount = 2;

      await store.removeRoom(1);

      expect(store.rooms).toHaveLength(1);
      expect(store.rooms[0].id).toBe(2);
      expect(store.totalCount).toBe(1);
      expect(store.loading).toBe(false);
      expect(store.error).toBeNull();
    });

    it("객실 삭제 시 currentRoom도 함께 초기화한다", async () => {
      vi.mocked(roomApi.deleteRoom).mockResolvedValue(undefined);

      const store = useRoomStore();
      store.rooms = [createMockRoom(1, "101"), createMockRoom(2, "102")];
      store.currentRoom = createMockRoom(1, "101");
      store.totalCount = 2;

      await store.removeRoom(1);

      expect(store.currentRoom).toBeNull();
    });

    it("다른 객실을 삭제해도 currentRoom은 유지된다", async () => {
      vi.mocked(roomApi.deleteRoom).mockResolvedValue(undefined);

      const store = useRoomStore();
      store.rooms = [createMockRoom(1, "101"), createMockRoom(2, "102")];
      store.currentRoom = createMockRoom(1, "101");
      store.totalCount = 2;

      await store.removeRoom(2);

      expect(store.currentRoom).toEqual(createMockRoom(1, "101"));
    });

    it("객실 삭제 중 로딩 상태를 관리한다", async () => {
      vi.mocked(roomApi.deleteRoom).mockImplementation(
        () =>
          new Promise((resolve) => {
            setTimeout(() => resolve(undefined), 10);
          })
      );

      const store = useRoomStore();
      store.rooms = [createMockRoom(1, "101")];
      const promise = store.removeRoom(1);

      expect(store.loading).toBe(true);

      await promise;

      expect(store.loading).toBe(false);
    });

    it("객실 삭제 실패 시 에러를 처리한다", async () => {
      const errorMessage = "Failed to delete room";
      vi.mocked(roomApi.deleteRoom).mockRejectedValue(createMockApiError(errorMessage));

      const store = useRoomStore();

      await expect(store.removeRoom(1)).rejects.toThrow(errorMessage);
      expect(store.error).toBe(errorMessage);
      expect(store.loading).toBe(false);
    });
  });

  describe("fetchRoomHistoryList", () => {
    it("객실 히스토리를 성공적으로 조회한다", async () => {
      const mockRoom = createMockRoom(1, "101");
      const mockHistories = [createMockRevision(mockRoom)];
      vi.mocked(roomApi.fetchRoomHistories).mockResolvedValue(createMockApiResponse(mockHistories, 1));

      const store = useRoomStore();
      await store.fetchRoomHistoryList(1);

      expect(store.roomHistories).toEqual(mockHistories);
      expect(store.loading).toBe(false);
      expect(store.error).toBeNull();
    });

    it("페이지네이션 파라미터와 함께 히스토리를 조회한다", async () => {
      const mockRoom = createMockRoom(1, "101");
      const mockHistories = [createMockRevision(mockRoom)];
      vi.mocked(roomApi.fetchRoomHistories).mockResolvedValue(createMockApiResponse(mockHistories, 1));

      const store = useRoomStore();
      const params = {
        page: 1,
        size: 10,
        sort: "historyCreatedAt,desc",
      };
      await store.fetchRoomHistoryList(1, params);

      expect(roomApi.fetchRoomHistories).toHaveBeenCalledWith(1, params);
      expect(store.roomHistories).toEqual(mockHistories);
    });

    it("히스토리 조회 중 로딩 상태를 관리한다", async () => {
      const mockRoom = createMockRoom(1, "101");
      const mockHistories = [createMockRevision(mockRoom)];
      vi.mocked(roomApi.fetchRoomHistories).mockImplementation(
        () =>
          new Promise((resolve) => {
            setTimeout(() => resolve(createMockApiResponse(mockHistories, 1)), 10);
          })
      );

      const store = useRoomStore();
      const promise = store.fetchRoomHistoryList(1);

      expect(store.loading).toBe(true);

      await promise;

      expect(store.loading).toBe(false);
    });

    it("히스토리 조회 실패 시 에러를 처리한다", async () => {
      const errorMessage = "Failed to fetch room histories";
      vi.mocked(roomApi.fetchRoomHistories).mockRejectedValue(createMockApiError(errorMessage));

      const store = useRoomStore();

      await expect(store.fetchRoomHistoryList(1)).rejects.toThrow(errorMessage);
      expect(store.error).toBe(errorMessage);
      expect(store.loading).toBe(false);
    });
  });

  describe("clearError", () => {
    it("에러 상태를 초기화한다", () => {
      const store = useRoomStore();
      store.error = "Some error";

      store.clearError();

      expect(store.error).toBeNull();
    });
  });

  describe("clearCurrentRoom", () => {
    it("현재 객실 상태를 초기화한다", () => {
      const store = useRoomStore();
      store.currentRoom = createMockRoom(1, "101");

      store.clearCurrentRoom();

      expect(store.currentRoom).toBeNull();
    });
  });

  describe("초기 상태", () => {
    it("스토어가 올바른 초기 상태를 가진다", () => {
      const store = useRoomStore();

      expect(store.rooms).toEqual([]);
      expect(store.currentRoom).toBeNull();
      expect(store.roomHistories).toEqual([]);
      expect(store.loading).toBe(false);
      expect(store.error).toBeNull();
      expect(store.totalCount).toBe(0);
    });
  });
});
