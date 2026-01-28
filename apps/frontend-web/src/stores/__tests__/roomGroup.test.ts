import { describe, it, expect, vi, beforeEach } from "vitest";
import { setActivePinia, createPinia } from "pinia";
import { useRoomGroupStore } from "../roomGroup";
import * as roomGroupApi from "src/api/v1/room-group";
import { createMockApiResponse, createMockApiError } from "test/vitest/helpers";
import type { RoomGroup } from "src/schema/room-group";
import type { RoomGroupDetailResponse } from "src/api/v1/room-group";

vi.mock("src/api/v1/room-group");

describe("useRoomGroupStore", () => {
  beforeEach(() => {
    setActivePinia(createPinia());
    vi.clearAllMocks();
  });

  // Mock 데이터 생성 헬퍼
  const createMockRoomGroup = (id: number, name: string): Omit<RoomGroup, "rooms"> => ({
    id,
    name,
    description: "테스트 객실 그룹",
    peekPrice: 100000,
    offPeekPrice: 80000,
    createdAt: "2026-01-01T00:00:00Z",
    updatedAt: "2026-01-01T00:00:00Z",
    createdBy: "admin",
    updatedBy: "admin",
  });

  const createMockRoomGroupDetail = (id: number, name: string): RoomGroupDetailResponse => ({
    id,
    name,
    description: "테스트 객실 그룹",
    peekPrice: 100000,
    offPeekPrice: 80000,
    createdAt: "2026-01-01T00:00:00Z",
    updatedAt: "2026-01-01T00:00:00Z",
    createdBy: "admin",
    updatedBy: "admin",
    rooms: [],
  });

  describe("fetchRoomGroupList", () => {
    it("객실 그룹 목록을 성공적으로 조회한다", async () => {
      const mockRoomGroups = [createMockRoomGroup(1, "스탠다드"), createMockRoomGroup(2, "디럭스")];
      vi.mocked(roomGroupApi.fetchRoomGroups).mockResolvedValue(createMockApiResponse(mockRoomGroups));

      const store = useRoomGroupStore();
      await store.fetchRoomGroupList();

      expect(store.roomGroups).toEqual(mockRoomGroups);
      expect(store.loading).toBe(false);
      expect(store.error).toBeNull();
    });

    it("객실 그룹 목록 조회 중 로딩 상태를 관리한다", async () => {
      const mockRoomGroups = [createMockRoomGroup(1, "스탠다드")];
      vi.mocked(roomGroupApi.fetchRoomGroups).mockImplementation(
        () =>
          new Promise((resolve) => {
            setTimeout(() => resolve(createMockApiResponse(mockRoomGroups)), 10);
          })
      );

      const store = useRoomGroupStore();
      const promise = store.fetchRoomGroupList();

      expect(store.loading).toBe(true);

      await promise;

      expect(store.loading).toBe(false);
    });

    it("객실 그룹 목록 조회 실패 시 에러를 처리한다", async () => {
      const errorMessage = "Failed to fetch room groups";
      vi.mocked(roomGroupApi.fetchRoomGroups).mockRejectedValue(createMockApiError(errorMessage));

      const store = useRoomGroupStore();

      await expect(store.fetchRoomGroupList()).rejects.toThrow(errorMessage);
      expect(store.error).toBe(errorMessage);
      expect(store.loading).toBe(false);
    });
  });

  describe("fetchRoomGroupById", () => {
    it("ID로 객실 그룹을 성공적으로 조회한다", async () => {
      const mockRoomGroup = createMockRoomGroupDetail(1, "스탠다드");
      vi.mocked(roomGroupApi.fetchRoomGroup).mockResolvedValue(createMockApiResponse(mockRoomGroup));

      const store = useRoomGroupStore();
      await store.fetchRoomGroupById(1);

      expect(store.currentRoomGroup).toEqual(mockRoomGroup);
      expect(store.loading).toBe(false);
      expect(store.error).toBeNull();
    });

    it("필터 파라미터와 함께 객실 그룹을 조회한다", async () => {
      const mockRoomGroup = createMockRoomGroupDetail(1, "스탠다드");
      vi.mocked(roomGroupApi.fetchRoomGroup).mockResolvedValue(createMockApiResponse(mockRoomGroup));

      const store = useRoomGroupStore();
      const filter = {
        stayStartAt: "2026-01-01",
        stayEndAt: "2026-01-31",
        status: "NORMAL" as const,
      };
      await store.fetchRoomGroupById(1, filter);

      expect(roomGroupApi.fetchRoomGroup).toHaveBeenCalledWith(1, filter);
      expect(store.currentRoomGroup).toEqual(mockRoomGroup);
    });

    it("객실 그룹 조회 중 로딩 상태를 관리한다", async () => {
      const mockRoomGroup = createMockRoomGroupDetail(1, "스탠다드");
      vi.mocked(roomGroupApi.fetchRoomGroup).mockImplementation(
        () =>
          new Promise((resolve) => {
            setTimeout(() => resolve(createMockApiResponse(mockRoomGroup)), 10);
          })
      );

      const store = useRoomGroupStore();
      const promise = store.fetchRoomGroupById(1);

      expect(store.loading).toBe(true);

      await promise;

      expect(store.loading).toBe(false);
    });

    it("객실 그룹 조회 실패 시 에러를 처리한다", async () => {
      const errorMessage = "Failed to fetch room group";
      vi.mocked(roomGroupApi.fetchRoomGroup).mockRejectedValue(createMockApiError(errorMessage));

      const store = useRoomGroupStore();

      await expect(store.fetchRoomGroupById(1)).rejects.toThrow(errorMessage);
      expect(store.error).toBe(errorMessage);
      expect(store.loading).toBe(false);
    });
  });

  describe("createNewRoomGroup", () => {
    it("새 객실 그룹을 성공적으로 생성한다", async () => {
      const newRoomGroupData = {
        name: "프리미엄",
        peekPrice: 150000,
        offPeekPrice: 120000,
        description: "프리미엄 객실",
      };
      const createdRoomGroup = createMockRoomGroup(3, "프리미엄");
      vi.mocked(roomGroupApi.createRoomGroup).mockResolvedValue(createMockApiResponse(createdRoomGroup));

      const store = useRoomGroupStore();
      store.roomGroups = [createMockRoomGroup(1, "스탠다드"), createMockRoomGroup(2, "디럭스")];

      await store.createNewRoomGroup(newRoomGroupData);

      expect(store.roomGroups).toHaveLength(3);
      expect(store.roomGroups[2]).toEqual(createdRoomGroup);
      expect(store.loading).toBe(false);
      expect(store.error).toBeNull();
    });

    it("객실 그룹 생성 중 로딩 상태를 관리한다", async () => {
      const newRoomGroupData = {
        name: "프리미엄",
        peekPrice: 150000,
        offPeekPrice: 120000,
        description: "프리미엄 객실",
      };
      const createdRoomGroup = createMockRoomGroup(3, "프리미엄");
      vi.mocked(roomGroupApi.createRoomGroup).mockImplementation(
        () =>
          new Promise((resolve) => {
            setTimeout(() => resolve(createMockApiResponse(createdRoomGroup)), 10);
          })
      );

      const store = useRoomGroupStore();
      const promise = store.createNewRoomGroup(newRoomGroupData);

      expect(store.loading).toBe(true);

      await promise;

      expect(store.loading).toBe(false);
    });

    it("객실 그룹 생성 실패 시 에러를 처리한다", async () => {
      const newRoomGroupData = {
        name: "프리미엄",
        peekPrice: 150000,
        offPeekPrice: 120000,
        description: "프리미엄 객실",
      };
      const errorMessage = "Failed to create room group";
      vi.mocked(roomGroupApi.createRoomGroup).mockRejectedValue(createMockApiError(errorMessage));

      const store = useRoomGroupStore();

      await expect(store.createNewRoomGroup(newRoomGroupData)).rejects.toThrow(errorMessage);
      expect(store.error).toBe(errorMessage);
      expect(store.loading).toBe(false);
    });
  });

  describe("updateRoomGroup", () => {
    it("객실 그룹을 성공적으로 수정한다", async () => {
      const updateData = {
        name: "스탠다드 플러스",
        peekPrice: 110000,
      };
      const updatedRoomGroup = { ...createMockRoomGroup(1, "스탠다드 플러스"), peekPrice: 110000 };
      vi.mocked(roomGroupApi.patchRoomGroup).mockResolvedValue(createMockApiResponse(updatedRoomGroup));

      const store = useRoomGroupStore();
      store.roomGroups = [createMockRoomGroup(1, "스탠다드"), createMockRoomGroup(2, "디럭스")];

      await store.updateRoomGroup(1, updateData);

      expect(store.roomGroups[0]).toEqual(updatedRoomGroup);
      expect(store.loading).toBe(false);
      expect(store.error).toBeNull();
    });

    it("객실 그룹 수정 시 currentRoomGroup도 함께 업데이트한다", async () => {
      const updateData = {
        name: "스탠다드 플러스",
      };
      const updatedRoomGroup = createMockRoomGroup(1, "스탠다드 플러스");
      vi.mocked(roomGroupApi.patchRoomGroup).mockResolvedValue(createMockApiResponse(updatedRoomGroup));

      const store = useRoomGroupStore();
      store.roomGroups = [createMockRoomGroup(1, "스탠다드"), createMockRoomGroup(2, "디럭스")];
      store.currentRoomGroup = createMockRoomGroupDetail(1, "스탠다드");

      await store.updateRoomGroup(1, updateData);

      expect(store.currentRoomGroup).toMatchObject(updatedRoomGroup);
    });

    it("목록에 없는 객실 그룹을 수정해도 에러가 발생하지 않는다", async () => {
      const updateData = {
        name: "프리미엄",
      };
      const updatedRoomGroup = createMockRoomGroup(3, "프리미엄");
      vi.mocked(roomGroupApi.patchRoomGroup).mockResolvedValue(createMockApiResponse(updatedRoomGroup));

      const store = useRoomGroupStore();
      store.roomGroups = [createMockRoomGroup(1, "스탠다드"), createMockRoomGroup(2, "디럭스")];

      await store.updateRoomGroup(3, updateData);

      expect(store.roomGroups).toHaveLength(2);
      expect(store.loading).toBe(false);
    });

    it("객실 그룹 수정 중 로딩 상태를 관리한다", async () => {
      const updateData = {
        name: "스탠다드 플러스",
      };
      const updatedRoomGroup = createMockRoomGroup(1, "스탠다드 플러스");
      vi.mocked(roomGroupApi.patchRoomGroup).mockImplementation(
        () =>
          new Promise((resolve) => {
            setTimeout(() => resolve(createMockApiResponse(updatedRoomGroup)), 10);
          })
      );

      const store = useRoomGroupStore();
      store.roomGroups = [createMockRoomGroup(1, "스탠다드")];
      const promise = store.updateRoomGroup(1, updateData);

      expect(store.loading).toBe(true);

      await promise;

      expect(store.loading).toBe(false);
    });

    it("객실 그룹 수정 실패 시 에러를 처리한다", async () => {
      const updateData = {
        name: "스탠다드 플러스",
      };
      const errorMessage = "Failed to update room group";
      vi.mocked(roomGroupApi.patchRoomGroup).mockRejectedValue(createMockApiError(errorMessage));

      const store = useRoomGroupStore();

      await expect(store.updateRoomGroup(1, updateData)).rejects.toThrow(errorMessage);
      expect(store.error).toBe(errorMessage);
      expect(store.loading).toBe(false);
    });
  });

  describe("removeRoomGroup", () => {
    it("객실 그룹을 성공적으로 삭제한다", async () => {
      vi.mocked(roomGroupApi.deleteRoomGroup).mockResolvedValue(undefined);

      const store = useRoomGroupStore();
      store.roomGroups = [createMockRoomGroup(1, "스탠다드"), createMockRoomGroup(2, "디럭스")];

      await store.removeRoomGroup(1);

      expect(store.roomGroups).toHaveLength(1);
      expect(store.roomGroups[0].id).toBe(2);
      expect(store.loading).toBe(false);
      expect(store.error).toBeNull();
    });

    it("객실 그룹 삭제 시 currentRoomGroup도 함께 초기화한다", async () => {
      vi.mocked(roomGroupApi.deleteRoomGroup).mockResolvedValue(undefined);

      const store = useRoomGroupStore();
      store.roomGroups = [createMockRoomGroup(1, "스탠다드"), createMockRoomGroup(2, "디럭스")];
      store.currentRoomGroup = createMockRoomGroupDetail(1, "스탠다드");

      await store.removeRoomGroup(1);

      expect(store.currentRoomGroup).toBeNull();
    });

    it("다른 객실 그룹을 삭제해도 currentRoomGroup은 유지된다", async () => {
      vi.mocked(roomGroupApi.deleteRoomGroup).mockResolvedValue(undefined);

      const store = useRoomGroupStore();
      store.roomGroups = [createMockRoomGroup(1, "스탠다드"), createMockRoomGroup(2, "디럭스")];
      store.currentRoomGroup = createMockRoomGroupDetail(1, "스탠다드");

      await store.removeRoomGroup(2);

      expect(store.currentRoomGroup).toEqual(createMockRoomGroupDetail(1, "스탠다드"));
    });

    it("객실 그룹 삭제 중 로딩 상태를 관리한다", async () => {
      vi.mocked(roomGroupApi.deleteRoomGroup).mockImplementation(
        () =>
          new Promise((resolve) => {
            setTimeout(() => resolve(undefined), 10);
          })
      );

      const store = useRoomGroupStore();
      store.roomGroups = [createMockRoomGroup(1, "스탠다드")];
      const promise = store.removeRoomGroup(1);

      expect(store.loading).toBe(true);

      await promise;

      expect(store.loading).toBe(false);
    });

    it("객실 그룹 삭제 실패 시 에러를 처리한다", async () => {
      const errorMessage = "Failed to delete room group";
      vi.mocked(roomGroupApi.deleteRoomGroup).mockRejectedValue(createMockApiError(errorMessage));

      const store = useRoomGroupStore();

      await expect(store.removeRoomGroup(1)).rejects.toThrow(errorMessage);
      expect(store.error).toBe(errorMessage);
      expect(store.loading).toBe(false);
    });
  });

  describe("clearError", () => {
    it("에러 상태를 초기화한다", () => {
      const store = useRoomGroupStore();
      store.error = "Some error";

      store.clearError();

      expect(store.error).toBeNull();
    });
  });

  describe("clearCurrentRoomGroup", () => {
    it("현재 객실 그룹 상태를 초기화한다", () => {
      const store = useRoomGroupStore();
      store.currentRoomGroup = createMockRoomGroupDetail(1, "스탠다드");

      store.clearCurrentRoomGroup();

      expect(store.currentRoomGroup).toBeNull();
    });
  });

  describe("초기 상태", () => {
    it("스토어가 올바른 초기 상태를 가진다", () => {
      const store = useRoomGroupStore();

      expect(store.roomGroups).toEqual([]);
      expect(store.currentRoomGroup).toBeNull();
      expect(store.loading).toBe(false);
      expect(store.error).toBeNull();
    });
  });
});
