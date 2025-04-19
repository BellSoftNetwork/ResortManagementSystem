<template>
  <q-page padding>
    <div class="q-pa-md">
      <div class="row q-col-gutter-md q-mb-md">
        <div class="col-12 col-md-6">
          <h5 class="q-mt-none q-mb-md">객실 현황</h5>
        </div>
        <div class="col-12 col-md-6">
          <div class="row q-col-gutter-sm">
            <div class="col">
              <q-input
                v-model="dateRange.from"
                label="시작일"
                type="date"
                outlined
                dense
                @update:model-value="loadRoomStatus"
              />
            </div>
            <div class="col">
              <q-input
                v-model="dateRange.to"
                label="종료일"
                type="date"
                outlined
                dense
                @update:model-value="loadRoomStatus"
              />
            </div>
            <div class="col-12 text-right q-mt-sm">
              <q-badge color="blue" class="q-pa-xs">
                기준일: {{ formatSimpleDate(today.toISOString().split("T")[0]) }}
              </q-badge>
            </div>
          </div>
        </div>
      </div>

      <div class="row q-col-gutter-md">
        <div v-if="loading" class="col-12 text-center q-pa-xl">
          <q-spinner-dots size="80px" color="primary" />
          <div class="text-body1 q-mt-md">객실 정보를 불러오는 중입니다...</div>
        </div>
        <template v-else>
          <div v-for="room in rooms" :key="room.id" class="col-12 col-md-6 col-lg-4">
            <q-card class="room-card">
              <q-card-section
                :class="isRoomOccupiedToday(room.id) ? 'bg-orange-8 text-white' : 'bg-primary text-white'"
              >
                <div class="row items-center">
                  <div class="col">
                    <div
                      class="text-h6 cursor-pointer"
                      @click="$router.push({ name: 'Room', params: { id: room.id } })"
                    >
                      {{ room.number }}
                      <q-badge :color="isRoomOccupiedToday(room.id) ? 'red' : 'green'" class="q-ml-sm">
                        {{ isRoomOccupiedToday(room.id) ? "입실 중" : "빈 방" }}
                      </q-badge>
                    </div>
                    <div class="text-subtitle2">{{ room.roomGroup.name }}</div>
                  </div>
                  <div class="col-auto">
                    <q-badge :color="getStatusBadgeColor(room.status)">
                      {{ roomStatusValueToName(room.status) }}
                    </q-badge>
                  </div>
                </div>
              </q-card-section>
              <q-list bordered separator v-if="roomReservations[room.id] && roomReservations[room.id].length > 0">
                <q-item
                  v-for="reservation in roomReservations[room.id]"
                  :key="reservation.id"
                  clickable
                  :to="{ name: 'Reservation', params: { id: reservation.id } }"
                  :class="{ 'bg-yellow-1': isCheckInOrOutToday(reservation) }"
                >
                  <q-item-section>
                    <q-item-label>{{ reservation.name }} ({{ reservation.phone }})</q-item-label>
                    <q-item-label caption>
                      {{ formatSimpleDate(reservation.stayStartAt) }} ~
                      {{ formatSimpleDate(reservation.stayEndAt) }}
                    </q-item-label>
                    <q-item-label caption>
                      <q-badge :color="getStatusColor(reservation.status)">
                        {{ reservationStatusValueToName(reservation.status) }}
                      </q-badge>
                      <q-badge v-if="isCheckInToday(reservation)" color="green" class="q-ml-xs">오늘 입실</q-badge>
                      <q-badge v-if="isCheckOutToday(reservation)" color="blue" class="q-ml-xs">오늘 퇴실</q-badge>
                    </q-item-label>
                  </q-item-section>
                  <q-item-section side>
                    <q-icon name="arrow_forward_ios" color="primary" size="xs" />
                  </q-item-section>
                </q-item>
              </q-list>
              <q-card-section v-else>
                <div class="text-grey text-center q-pa-md">해당 기간에 예약이 없습니다.</div>
              </q-card-section>
            </q-card>
          </div>
        </template>
      </div>
    </div>
  </q-page>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { fetchRooms } from "src/api/v1/room";
import { fetchReservations } from "src/api/v1/reservation";
import { Room, roomStatusValueToName } from "src/schema/room";
import { Reservation, ReservationStatus, reservationStatusValueToName } from "src/schema/reservation";
import { formatSimpleDate } from "src/util/format-util";

const rooms = ref<Room[]>([]);
const reservations = ref<Reservation[]>([]);
const loading = ref(true);

const today = new Date();
const dateRange = ref({
  from: today.toISOString().split("T")[0],
  to: new Date(today.getTime() + 7 * 24 * 60 * 60 * 1000).toISOString().split("T")[0],
});

const roomReservations = computed(() => {
  const result: Record<number, Reservation[]> = {};

  // Initialize with empty arrays for all rooms
  rooms.value.forEach((room) => {
    result[room.id] = [];
  });

  // Populate with reservations
  reservations.value.forEach((reservation) => {
    reservation.rooms.forEach((room) => {
      if (result[room.id]) {
        result[room.id].push(reservation);
      }
    });
  });

  return result;
});

function getStatusColor(status: ReservationStatus): string {
  switch (status) {
    case "NORMAL":
      return "positive";
    case "PENDING":
      return "warning";
    case "CANCEL":
      return "negative";
    case "REFUND":
      return "grey";
    default:
      return "grey";
  }
}

function getStatusBadgeColor(status: string): string {
  switch (status) {
    case "NORMAL":
      return "positive";
    case "INACTIVE":
      return "grey";
    case "CONSTRUCTION":
      return "orange";
    case "DAMAGED":
      return "negative";
    default:
      return "blue";
  }
}

// Check if a reservation has check-in today
function isCheckInToday(reservation: Reservation): boolean {
  const today = new Date().toISOString().split("T")[0];
  return reservation.stayStartAt === today;
}

// Check if a reservation has check-out today
function isCheckOutToday(reservation: Reservation): boolean {
  const today = new Date().toISOString().split("T")[0];
  return reservation.stayEndAt === today;
}

// Check if a reservation has either check-in or check-out today
function isCheckInOrOutToday(reservation: Reservation): boolean {
  return isCheckInToday(reservation) || isCheckOutToday(reservation);
}

// Check if a room is currently occupied (has an active reservation that includes today)
function isRoomOccupiedToday(roomId: number): boolean {
  const today = new Date().toISOString().split("T")[0];

  if (!roomReservations.value[roomId]) return false;

  return roomReservations.value[roomId].some((reservation) => {
    return reservation.status === "NORMAL" && reservation.stayStartAt <= today && reservation.stayEndAt > today;
  });
}

async function loadRoomStatus() {
  loading.value = true;

  try {
    // Fetch all rooms
    const roomsResponse = await fetchRooms({});
    rooms.value = roomsResponse.values;

    // Fetch reservations for the selected date range
    const reservationsResponse = await fetchReservations({
      stayStartAt: dateRange.value.from,
      stayEndAt: dateRange.value.to,
    });
    reservations.value = reservationsResponse.values;
  } catch (error) {
    console.error("Error loading room status:", error);
  } finally {
    loading.value = false;
  }
}

onMounted(() => {
  loadRoomStatus();
});
</script>

<style scoped>
.room-card {
  height: 100%;
  display: flex;
  flex-direction: column;
}
</style>
