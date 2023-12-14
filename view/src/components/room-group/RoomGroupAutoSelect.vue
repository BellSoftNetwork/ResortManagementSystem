<template>
  <q-card>
    <q-toolbar class="bg-primary text-white shadow-2">
      <q-toolbar-title>객실 배정</q-toolbar-title>
      <q-select
        @update:model-value="loadRoomGroup"
        v-model="selectedRoomGroup"
        :loading="roomGroups === null"
        :disable="roomGroups === null"
        :options="roomGroups"
        option-label="name"
        label="객실 그룹"
        bg-color="white"
        dense
        required
        map-options
        filled
      ></q-select>
    </q-toolbar>

    <div>
      <q-card-section v-if="roomGroupDetail === null || selectedRoomGroup === null" style="height: 200px">
        <q-inner-loading :showing="true">
          <q-spinner-gears size="100px" color="primary"></q-spinner-gears>
        </q-inner-loading>
      </q-card-section>
      <q-card-section v-else style="height: 200px">
        <p>
          객실 단가: {{ formatPrice(selectedRoomGroup.peekPrice) }} / {{ formatPrice(selectedRoomGroup.peekPrice) }}
        </p>
        <q-input
          v-model.number="addRoomCount"
          type="number"
          min="1"
          :max="filteredRooms.length"
          label="객실 개수"
          hint="희망 기간 전 퇴실일이 가장 먼 순서로 자동 배정됩니다."
        >
          <template v-slot:append>/ {{ filteredRooms.length }}개</template>
          <template v-slot:after>
            <q-btn @click="addRooms()" :disable="filteredRooms.length <= 0" color="primary" class="full-height"
              >추가 배정
            </q-btn>
          </template>
        </q-input>
      </q-card-section>
    </div>
  </q-card>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { Room } from "src/schema/room";
import { Reservation } from "src/schema/reservation";
import { fetchRoomGroup, fetchRoomGroups, RoomGroupDetailResponse } from "src/api/v1/room-group";
import { RoomGroup } from "src/schema/room-group";
import { formatPrice } from "src/util/format-util";

const props = defineProps<{
  selected: Room[];
  parentReservation?: Reservation;
  stayStartAt: string;
  stayEndAt: string;
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

const selectedRoomGroup = ref<RoomGroup | null>(null);
const roomGroups = ref<RoomGroup[] | null>(null);
const roomGroupDetail = ref<RoomGroupDetailResponse | null>(null);
const filteredRooms = computed(() => {
  const selectedIds = selected.value.map((selectedRoom) => selectedRoom.id);

  return roomGroupDetail.value.rooms.filter(
    (roomWithReservation) => !selectedIds.includes(roomWithReservation.room.id),
  );
});
const addRoomCount = ref(1);

function loadRoomGroups() {
  roomGroups.value = null;

  return fetchRoomGroups().then((response) => {
    roomGroups.value = response.values;
  });
}

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

function loadRoomGroup(roomGroup: RoomGroup) {
  roomGroupDetail.value = null;

  fetchRoomGroup(roomGroup.id, {
    stayStartAt: props.stayStartAt,
    stayEndAt: props.stayEndAt,
    status: "NORMAL",
    excludeReservationId: props.parentReservation?.id ?? undefined,
  }).then((response) => {
    roomGroupDetail.value = response.value;
  });
}

onMounted(() => {
  loadRoomGroups().then(() => {
    if (roomGroups.value.length > 0) {
      selectedRoomGroup.value = roomGroups.value[0];
      loadRoomGroup(roomGroups.value[0]);
    }
  });
});
</script>
