<template>
  <div v-if="filteredRooms.length > 0">
    <q-virtual-scroll :items="filteredRooms" v-slot="{ item }" style="height: 200px" bordered separator>
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

  <q-card-section v-else style="height: 200px">
    <div class="text-center">
      <p>배정 가능한 객실이 없습니다.</p>
    </div>
  </q-card-section>
</template>

<script setup lang="ts">
import { computed } from "vue";
import { Room } from "src/schema/room";
import { Reservation } from "src/schema/reservation";
import { RoomGroupDetailResponse } from "src/api/v1/room-group";

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

function addRoom(room: Room) {
  if (!selected.value.find((r) => r.id === room.id)) {
    selected.value.push(room);
    selected.value.sort((a, b) => a.id - b.id);
  }
}
</script>
