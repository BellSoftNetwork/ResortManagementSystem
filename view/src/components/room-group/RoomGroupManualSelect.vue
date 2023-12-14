<template>
  <q-card>
    <q-toolbar class="bg-primary text-white shadow-2">
      <q-toolbar-title>배정 가능 객실</q-toolbar-title>
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
      <q-card-section v-if="roomGroupDetail === null" style="height: 200px">
        <q-inner-loading :showing="true">
          <q-spinner-gears size="100px" color="primary"></q-spinner-gears>
        </q-inner-loading>
      </q-card-section>
      <div v-else>
        <q-virtual-scroll
          v-if="roomGroupDetail"
          :items="filteredRooms"
          v-slot="{ item }"
          style="height: 200px"
          bordered
          separator
        >
          <q-item :key="item.room.id" @click="addRoom(item.room)" clickable v-ripple>
            <q-item-section>
              <q-item-label>
                {{ item.room.number }} (최근 퇴실일:&nbsp;
                <span v-if="item.lastReservation">{{ item.lastReservation.stayEndAt }}</span>
                <span v-else>정보 없음</span>)
              </q-item-label>
            </q-item-section>
          </q-item>
        </q-virtual-scroll>
      </div>
    </div>
  </q-card>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { Room } from "src/schema/room";
import { Reservation } from "src/schema/reservation";
import { fetchRoomGroup, fetchRoomGroups, RoomGroupDetailResponse } from "src/api/v1/room-group";
import { RoomGroup } from "src/schema/room-group";

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

function loadRoomGroups() {
  roomGroups.value = null;

  return fetchRoomGroups().then((response) => {
    roomGroups.value = response.values;
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
