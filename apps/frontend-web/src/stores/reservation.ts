import { defineStore } from "pinia";
import { ref } from "vue";
import { Reservation, ReservationStatus, ReservationType } from "src/schema/reservation";
import {
  fetchReservations,
  fetchReservation,
  createReservation,
  patchReservation,
  deleteReservation,
  fetchReservationHistories,
  fetchReservationStatistics,
  ReservationCreateParams,
  ReservationPatchParams,
  StatisticsPeriodType,
  ReservationStatistics,
} from "src/api/v1/reservation";
import type { PageRequestParams, SortRequestParams } from "src/schema/response";
import type { Revision } from "src/schema/revision";

export type ReservationFilter = Partial<{
  stayStartAt: string;
  stayEndAt: string;
  searchText: string;
  status: ReservationStatus;
  type: ReservationType;
}> &
  Partial<PageRequestParams> &
  Partial<SortRequestParams>;

export type CreateReservationRequest = ReservationCreateParams;
export type UpdateReservationRequest = ReservationPatchParams;

export const useReservationStore = defineStore("reservation", () => {
  // State
  const reservations = ref<Reservation[]>([]);
  const currentReservation = ref<Reservation | null>(null);
  const reservationHistories = ref<Revision<Reservation>[]>([]);
  const statistics = ref<ReservationStatistics | null>(null);
  const loading = ref(false);
  const error = ref<string | null>(null);
  const totalCount = ref(0);

  // Actions
  async function fetchReservationList(filter?: ReservationFilter) {
    loading.value = true;
    error.value = null;

    try {
      const response = await fetchReservations(filter || {});
      reservations.value = response.value;
      totalCount.value = response.totalCount || 0;
      return response;
    } catch (err) {
      error.value = err instanceof Error ? err.message : "Failed to fetch reservations";
      throw err;
    } finally {
      loading.value = false;
    }
  }

  async function fetchReservationById(id: number) {
    loading.value = true;
    error.value = null;

    try {
      const response = await fetchReservation(id);
      currentReservation.value = response.value;
      return response;
    } catch (err) {
      error.value = err instanceof Error ? err.message : "Failed to fetch reservation";
      throw err;
    } finally {
      loading.value = false;
    }
  }

  async function createNewReservation(data: CreateReservationRequest) {
    loading.value = true;
    error.value = null;

    try {
      const response = await createReservation(data);
      // Add the new reservation to the list
      reservations.value.push(response.value);
      totalCount.value += 1;
      return response;
    } catch (err) {
      error.value = err instanceof Error ? err.message : "Failed to create reservation";
      throw err;
    } finally {
      loading.value = false;
    }
  }

  async function updateReservation(id: number, data: UpdateReservationRequest) {
    loading.value = true;
    error.value = null;

    try {
      const response = await patchReservation(id, data);
      // Update the reservation in the list
      const index = reservations.value.findIndex((reservation) => reservation.id === id);
      if (index !== -1) {
        reservations.value[index] = response.value;
      }
      // Update current reservation if it's the same
      if (currentReservation.value?.id === id) {
        currentReservation.value = response.value;
      }
      return response;
    } catch (err) {
      error.value = err instanceof Error ? err.message : "Failed to update reservation";
      throw err;
    } finally {
      loading.value = false;
    }
  }

  async function removeReservation(id: number) {
    loading.value = true;
    error.value = null;

    try {
      await deleteReservation(id);
      // Remove the reservation from the list
      reservations.value = reservations.value.filter((reservation) => reservation.id !== id);
      totalCount.value -= 1;
      // Clear current reservation if it's the same
      if (currentReservation.value?.id === id) {
        currentReservation.value = null;
      }
    } catch (err) {
      error.value = err instanceof Error ? err.message : "Failed to delete reservation";
      throw err;
    } finally {
      loading.value = false;
    }
  }

  async function fetchReservationHistoryList(id: number, params?: Partial<PageRequestParams & SortRequestParams>) {
    loading.value = true;
    error.value = null;

    try {
      const response = await fetchReservationHistories(id, params || {});
      reservationHistories.value = response.value;
      return response;
    } catch (err) {
      error.value = err instanceof Error ? err.message : "Failed to fetch reservation histories";
      throw err;
    } finally {
      loading.value = false;
    }
  }

  async function fetchStatistics(
    startDate: string,
    endDate: string,
    periodType: StatisticsPeriodType = StatisticsPeriodType.MONTHLY,
  ) {
    loading.value = true;
    error.value = null;

    try {
      const response = await fetchReservationStatistics(startDate, endDate, periodType);
      statistics.value = response.value;
      return response;
    } catch (err) {
      error.value = err instanceof Error ? err.message : "Failed to fetch reservation statistics";
      throw err;
    } finally {
      loading.value = false;
    }
  }

  function clearError() {
    error.value = null;
  }

  function clearCurrentReservation() {
    currentReservation.value = null;
  }

  function clearStatistics() {
    statistics.value = null;
  }

  return {
    // State
    reservations,
    currentReservation,
    reservationHistories,
    statistics,
    loading,
    error,
    totalCount,
    // Actions
    fetchReservationList,
    fetchReservationById,
    createNewReservation,
    updateReservation,
    removeReservation,
    fetchReservationHistoryList,
    fetchStatistics,
    clearError,
    clearCurrentReservation,
    clearStatistics,
  };
});
