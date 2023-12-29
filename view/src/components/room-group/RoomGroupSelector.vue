<template>
  <q-btn-toggle
    v-model="selectMode"
    spread
    no-caps
    toggle-color="primary"
    :options="[
      { label: '자동 배정', value: 'auto' },
      { label: '수동 배정', value: 'manual' },
    ]"
  />

  <div class="row">
    <div :class="hideSelected ? '' : 'col-12 col-md-6 q-pa-sm'">
      <q-card>
        <q-toolbar class="bg-primary text-white shadow-2">
          <q-toolbar-title>객실 배정</q-toolbar-title>

          &nbsp;
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

        <q-card-section v-if="roomGroupDetail === null" style="min-height: 200px">
          <q-inner-loading :showing="true">
            <div v-if="stayStartAt && stayEndAt">
              <q-spinner-gears size="100px" color="primary"></q-spinner-gears>
            </div>
            <div v-else>
              <div class="text-center">
                <q-icon name="warning" color="red" size="100px"></q-icon>
              </div>
            </div>
          </q-inner-loading>
        </q-card-section>

        <div v-else>
          <RoomGroupAutoSelect
            v-if="selectMode === 'auto'"
            v-model:selected="selected"
            :room-group-detail="roomGroupDetail"
            :stay-start-at="props.stayStartAt"
            :stay-end-at="props.stayEndAt"
          />
          <RoomGroupManualSelect
            v-else
            v-model:selected="selected"
            :room-group-detail="roomGroupDetail"
            :stay-start-at="props.stayStartAt"
            :stay-end-at="props.stayEndAt"
          />
        </div>
      </q-card>
    </div>

    <div v-if="!hideSelected" class="col-12 col-md-6 q-pa-sm">
      <q-card>
        <q-toolbar class="bg-accent text-white shadow-2">
          <q-toolbar-title>배정된 객실</q-toolbar-title>
        </q-toolbar>
        <div style="height: 200px">
          <q-virtual-scroll :items="selected" v-slot="{ item }" style="height: 200px" bordered separator>
            <q-item :key="item.id" @click="selected = selected.filter((r) => r.id !== item.id)" clickable v-ripple>
              <q-item-section>
                <q-item-label>{{ item.number }}</q-item-label>
              </q-item-section>
            </q-item>
          </q-virtual-scroll>
        </div>
      </q-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from "vue";
import { Room } from "src/schema/room";
import { Reservation } from "src/schema/reservation";
import RoomGroupManualSelect from "components/room-group/RoomGroupManualSelect.vue";
import RoomGroupAutoSelect from "components/room-group/RoomGroupAutoSelect.vue";
import { RoomGroup } from "src/schema/room-group";
import { fetchRoomGroup, fetchRoomGroups, RoomGroupDetailResponse } from "src/api/v1/room-group";

const props = defineProps<{
  selected: Room[];
  parentReservation?: Reservation;
  stayStartAt: string | null;
  stayEndAt: string | null;
  hideSelected?: boolean;
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
const roomGroups = ref<RoomGroup[] | null>(null);
const selectedRoomGroup = ref<RoomGroup | null>(null);
const roomGroupDetail = ref<RoomGroupDetailResponse | null>(null);

const selectMode = ref("auto");

function loadRoomGroups() {
  roomGroups.value = null;

  return fetchRoomGroups().then((response) => {
    roomGroups.value = response.values;
  });
}

function loadRoomGroup(roomGroup: RoomGroup) {
  roomGroupDetail.value = null;

  if (!(props.stayStartAt && props.stayEndAt)) return;

  fetchRoomGroup(roomGroup.id, {
    stayStartAt: props.stayStartAt,
    stayEndAt: props.stayEndAt,
    status: "NORMAL",
    excludeReservationId: props.parentReservation?.id ?? undefined,
  }).then((response) => {
    roomGroupDetail.value = response.value;
  });
}

watch(props, (value, oldValue) => {
  if (value.stayStartAt !== oldValue.stayStartAt || value.stayEndAt !== oldValue.stayEndAt) {
    if (selectedRoomGroup.value) loadRoomGroup(selectedRoomGroup.value);
  }
});

onMounted(() => {
  loadRoomGroups().then(() => {
    if (roomGroups.value.length > 0) {
      selectedRoomGroup.value = roomGroups.value[0];
      loadRoomGroup(roomGroups.value[0]);
    }
  });
});
</script>
