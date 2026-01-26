import { defineStore } from "pinia";
import { ref } from "vue";
import { RoomGroup } from "src/schema/room-group";
import {
  fetchRoomGroups,
  fetchRoomGroup,
  createRoomGroup,
  patchRoomGroup,
  deleteRoomGroup,
  RoomGroupDetailResponse,
} from "src/api/v1/room-group";
import type { RoomStatus } from "src/schema/room";

export type RoomGroupFilter = Partial<{
  stayStartAt: string;
  stayEndAt: string;
  status: RoomStatus;
  excludeReservationId: number;
}>;

export type CreateRoomGroupRequest = {
  name: string;
  peekPrice: number;
  offPeekPrice: number;
  description: string;
};

export type UpdateRoomGroupRequest = Partial<CreateRoomGroupRequest>;

// RoomGroupSummary type (without rooms array)
type RoomGroupSummary = Omit<RoomGroup, "rooms">;

export const useRoomGroupStore = defineStore("roomGroup", () => {
  // State
  const roomGroups = ref<RoomGroupSummary[]>([]);
  const currentRoomGroup = ref<RoomGroupDetailResponse | null>(null);
  const loading = ref(false);
  const error = ref<string | null>(null);

  // Actions
  async function fetchRoomGroupList() {
    loading.value = true;
    error.value = null;

    try {
      const response = await fetchRoomGroups();
      roomGroups.value = response.value;
      return response;
    } catch (err) {
      error.value = err instanceof Error ? err.message : "Failed to fetch room groups";
      throw err;
    } finally {
      loading.value = false;
    }
  }

  async function fetchRoomGroupById(id: number, filter?: RoomGroupFilter) {
    loading.value = true;
    error.value = null;

    try {
      const response = await fetchRoomGroup(id, filter || {});
      currentRoomGroup.value = response.value;
      return response;
    } catch (err) {
      error.value = err instanceof Error ? err.message : "Failed to fetch room group";
      throw err;
    } finally {
      loading.value = false;
    }
  }

  async function createNewRoomGroup(data: CreateRoomGroupRequest) {
    loading.value = true;
    error.value = null;

    try {
      const response = await createRoomGroup(data);
      // Add the new room group to the list
      roomGroups.value.push(response.value);
      return response;
    } catch (err) {
      error.value = err instanceof Error ? err.message : "Failed to create room group";
      throw err;
    } finally {
      loading.value = false;
    }
  }

  async function updateRoomGroup(id: number, data: UpdateRoomGroupRequest) {
    loading.value = true;
    error.value = null;

    try {
      const response = await patchRoomGroup(id, data);
      // Update the room group in the list
      const index = roomGroups.value.findIndex((roomGroup) => roomGroup.id === id);
      if (index !== -1) {
        roomGroups.value[index] = response.value;
      }
      // Update current room group if it's the same
      if (currentRoomGroup.value?.id === id) {
        currentRoomGroup.value = {
          ...currentRoomGroup.value,
          ...response.value,
        };
      }
      return response;
    } catch (err) {
      error.value = err instanceof Error ? err.message : "Failed to update room group";
      throw err;
    } finally {
      loading.value = false;
    }
  }

  async function removeRoomGroup(id: number) {
    loading.value = true;
    error.value = null;

    try {
      await deleteRoomGroup(id);
      // Remove the room group from the list
      roomGroups.value = roomGroups.value.filter((roomGroup) => roomGroup.id !== id);
      // Clear current room group if it's the same
      if (currentRoomGroup.value?.id === id) {
        currentRoomGroup.value = null;
      }
    } catch (err) {
      error.value = err instanceof Error ? err.message : "Failed to delete room group";
      throw err;
    } finally {
      loading.value = false;
    }
  }

  function clearError() {
    error.value = null;
  }

  function clearCurrentRoomGroup() {
    currentRoomGroup.value = null;
  }

  return {
    // State
    roomGroups,
    currentRoomGroup,
    loading,
    error,
    // Actions
    fetchRoomGroupList,
    fetchRoomGroupById,
    createNewRoomGroup,
    updateRoomGroup,
    removeRoomGroup,
    clearError,
    clearCurrentRoomGroup,
  };
});
