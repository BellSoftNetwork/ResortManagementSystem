import { defineStore } from "pinia";
import { ref } from "vue";
import { DateBlock } from "src/schema/date-block";
import {
  fetchDateBlocks,
  fetchDateBlock,
  createDateBlock,
  updateDateBlock,
  deleteDateBlock,
  fetchDateBlockHistories,
} from "src/api/v1/date-block";
import type { PageRequestParams, SortRequestParams } from "src/schema/response";
import type { DateBlockFilter, CreateDateBlockRequest, UpdateDateBlockRequest } from "src/schema/date-block";
import type { Revision } from "src/schema/revision";

export const useDateBlockStore = defineStore("dateBlock", () => {
  // State
  const dateBlocks = ref<DateBlock[]>([]);
  const currentDateBlock = ref<DateBlock | null>(null);
  const dateBlockHistories = ref<Revision<DateBlock>[]>([]);
  const loading = ref(false);
  const error = ref<string | null>(null);
  const totalCount = ref(0);

  // Actions
  async function fetchDateBlockList(filter?: DateBlockFilter) {
    loading.value = true;
    error.value = null;

    try {
      const response = await fetchDateBlocks(filter || {});
      dateBlocks.value = response.values;
      totalCount.value = response.page?.totalElements || 0;
      return response;
    } catch (err) {
      error.value = err instanceof Error ? err.message : "Failed to fetch date blocks";
      throw err;
    } finally {
      loading.value = false;
    }
  }

  async function fetchDateBlockById(id: number) {
    loading.value = true;
    error.value = null;

    try {
      const response = await fetchDateBlock(id);
      currentDateBlock.value = response.value;
      return response;
    } catch (err) {
      error.value = err instanceof Error ? err.message : "Failed to fetch date block";
      throw err;
    } finally {
      loading.value = false;
    }
  }

  async function createNewDateBlock(data: CreateDateBlockRequest) {
    loading.value = true;
    error.value = null;

    try {
      const response = await createDateBlock(data);
      dateBlocks.value.push(response.value);
      totalCount.value += 1;
      return response;
    } catch (err) {
      error.value = err instanceof Error ? err.message : "Failed to create date block";
      throw err;
    } finally {
      loading.value = false;
    }
  }

  async function updateDateBlockById(id: number, data: UpdateDateBlockRequest) {
    loading.value = true;
    error.value = null;

    try {
      const response = await updateDateBlock(id, data);
      const index = dateBlocks.value.findIndex((block) => block.id === id);
      if (index !== -1) {
        dateBlocks.value[index] = response.value;
      }
      if (currentDateBlock.value?.id === id) {
        currentDateBlock.value = response.value;
      }
      return response;
    } catch (err) {
      error.value = err instanceof Error ? err.message : "Failed to update date block";
      throw err;
    } finally {
      loading.value = false;
    }
  }

  async function removeDateBlock(id: number) {
    loading.value = true;
    error.value = null;

    try {
      await deleteDateBlock(id);
      dateBlocks.value = dateBlocks.value.filter((block) => block.id !== id);
      totalCount.value -= 1;
      if (currentDateBlock.value?.id === id) {
        currentDateBlock.value = null;
      }
    } catch (err) {
      error.value = err instanceof Error ? err.message : "Failed to delete date block";
      throw err;
    } finally {
      loading.value = false;
    }
  }

  async function fetchDateBlockHistoryList(id: number, params?: Partial<PageRequestParams & SortRequestParams>) {
    loading.value = true;
    error.value = null;

    try {
      const response = await fetchDateBlockHistories(id, params || {});
      dateBlockHistories.value = response.values;
      return response;
    } catch (err) {
      error.value = err instanceof Error ? err.message : "Failed to fetch date block histories";
      throw err;
    } finally {
      loading.value = false;
    }
  }

  function clearError() {
    error.value = null;
  }

  function clearCurrentDateBlock() {
    currentDateBlock.value = null;
  }

  return {
    // State
    dateBlocks,
    currentDateBlock,
    dateBlockHistories,
    loading,
    error,
    totalCount,
    // Actions
    fetchDateBlockList,
    fetchDateBlockById,
    createNewDateBlock,
    updateDateBlockById,
    removeDateBlock,
    fetchDateBlockHistoryList,
    clearError,
    clearCurrentDateBlock,
  };
});
