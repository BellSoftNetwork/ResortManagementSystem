<template>
  <q-card-section style="min-height: 200px">
    <p>객실 단가: {{ formatPrice(roomGroupDetail.peekPrice) }} / {{ formatPrice(roomGroupDetail.peekPrice) }}</p>
    <q-input
      v-model.number="addRoomCount"
      type="number"
      min="1"
      :max="filteredRooms?.length"
      label="객실 개수"
      hint="희망 기간 전 퇴실일이 가장 먼 순서로 자동 배정됩니다."
      class="q-my-md"
    >
      <template v-slot:append>/ {{ filteredRooms.length }}개</template>
    </q-input>

    <q-btn @click="addRooms()" :disable="filteredRooms.length <= 0" color="primary" class="full-width">
      추가 배정
    </q-btn>
  </q-card-section>
</template>

<script setup lang="ts">
import { computed, ref } from "vue";
import { Room } from "src/schema/room";
import { Reservation } from "src/schema/reservation";
import { RoomGroupDetailResponse } from "src/api/v1/room-group";
import { formatPrice } from "src/util/format-util";

const props = defineProps<{
  selected: Room[];
  parentReservation?: Reservation;
  stayStartAt: string | null;
  stayEndAt: string | null;
  roomGroupDetail: RoomGroupDetailResponse;
}>();
const emit = defineEmits(["update:selected"]);

const selected = computed({
  get() {
    return props.selected;
  },
  set(value) {
    emit("update:selected", value);
  },
});

const filteredRooms = computed(() => {
  const selectedIds = selected.value.map((selectedRoom) => selectedRoom.id);

  return props.roomGroupDetail.rooms.filter(
    (roomWithReservation) => !selectedIds.includes(roomWithReservation.room.id),
  );
});
const addRoomCount = ref(1);

function addRooms() {
  const count = addRoomCount.value > filteredRooms.value.length ? filteredRooms.value.length : addRoomCount.value;

  filteredRooms.value.slice(0, count).forEach((roomWithReservation) => {
    addRoom(roomWithReservation.room);
  });
}

function addRoom(room: Room) {
  if (!selected.value.find((r) => r.id === room.id)) {
    selected.value.push(room);
    selected.value.sort((a, b) => a.id - b.id);
  }
}
</script>
