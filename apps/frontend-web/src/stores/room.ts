import { defineStore } from "pinia";
import { ref } from "vue";
import { Room } from "src/schema/room";
import { fetchRooms, fetchRoom, createRoom, patchRoom, deleteRoom, fetchRoomHistories } from "src/api/v1/room";
import type { PageRequestParams, SortRequestParams } from "src/schema/response";
import type { RoomStatus } from "src/schema/room";
import type { Revision } from "src/schema/revision";
import type { EntityReference } from "src/schema/base";

export type RoomFilter = Partial<{
  stayStartAt: string;
  stayEndAt: string;
  status: RoomStatus;
  excludeReservationId: number;
}> &
  Partial<PageRequestParams> &
  Partial<SortRequestParams>;

export type CreateRoomRequest = {
  roomGroup: EntityReference;
  number: string;
  note: string;
  status: RoomStatus;
};

export type UpdateRoomRequest = Partial<CreateRoomRequest>;

export const useRoomStore = defineStore("room", () => {
  // State
  const rooms = ref<Room[]>([]);
  const currentRoom = ref<Room | null>(null);
  const roomHistories = ref<Revision<Room>[]>([]);
  const loading = ref(false);
  const error = ref<string | null>(null);
  const totalCount = ref(0);

  // Actions
  async function fetchRoomList(filter?: RoomFilter) {
    loading.value = true;
    error.value = null;

    try {
      const response = await fetchRooms(filter || {});
      rooms.value = response.value;
      totalCount.value = response.totalCount || 0;
      return response;
    } catch (err) {
      error.value = err instanceof Error ? err.message : "Failed to fetch rooms";
      throw err;
    } finally {
      loading.value = false;
    }
  }

  async function fetchRoomById(id: number) {
    loading.value = true;
    error.value = null;

    try {
      const response = await fetchRoom(id);
      currentRoom.value = response.value;
      return response;
    } catch (err) {
      error.value = err instanceof Error ? err.message : "Failed to fetch room";
      throw err;
    } finally {
      loading.value = false;
    }
  }

  async function createNewRoom(data: CreateRoomRequest) {
    loading.value = true;
    error.value = null;

    try {
      const response = await createRoom(data);
      // Add the new room to the list
      rooms.value.push(response.value);
      totalCount.value += 1;
      return response;
    } catch (err) {
      error.value = err instanceof Error ? err.message : "Failed to create room";
      throw err;
    } finally {
      loading.value = false;
    }
  }

  async function updateRoom(id: number, data: UpdateRoomRequest) {
    loading.value = true;
    error.value = null;

    try {
      const response = await patchRoom(id, data);
      // Update the room in the list
      const index = rooms.value.findIndex((room) => room.id === id);
      if (index !== -1) {
        rooms.value[index] = response.value;
      }
      // Update current room if it's the same
      if (currentRoom.value?.id === id) {
        currentRoom.value = response.value;
      }
      return response;
    } catch (err) {
      error.value = err instanceof Error ? err.message : "Failed to update room";
      throw err;
    } finally {
      loading.value = false;
    }
  }

  async function removeRoom(id: number) {
    loading.value = true;
    error.value = null;

    try {
      await deleteRoom(id);
      // Remove the room from the list
      rooms.value = rooms.value.filter((room) => room.id !== id);
      totalCount.value -= 1;
      // Clear current room if it's the same
      if (currentRoom.value?.id === id) {
        currentRoom.value = null;
      }
    } catch (err) {
      error.value = err instanceof Error ? err.message : "Failed to delete room";
      throw err;
    } finally {
      loading.value = false;
    }
  }

  async function fetchRoomHistoryList(id: number, params?: Partial<PageRequestParams & SortRequestParams>) {
    loading.value = true;
    error.value = null;

    try {
      const response = await fetchRoomHistories(id, params || {});
      roomHistories.value = response.value;
      return response;
    } catch (err) {
      error.value = err instanceof Error ? err.message : "Failed to fetch room histories";
      throw err;
    } finally {
      loading.value = false;
    }
  }

  function clearError() {
    error.value = null;
  }

  function clearCurrentRoom() {
    currentRoom.value = null;
  }

  return {
    // State
    rooms,
    currentRoom,
    roomHistories,
    loading,
    error,
    totalCount,
    // Actions
    fetchRoomList,
    fetchRoomById,
    createNewRoom,
    updateRoom,
    removeRoom,
    fetchRoomHistoryList,
    clearError,
    clearCurrentRoom,
  };
});
