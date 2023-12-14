<template>
  <q-btn-toggle
    v-model="selectMode"
    spread
    no-caps
    rounded
    toggle-color="primary"
    :options="[
      { label: '자동 배정', value: 'auto' },
      { label: '수동 배정', value: 'manual' },
    ]"
  />

  <div class="row">
    <div class="col-12 col-md-6 q-pa-sm">
      <div v-if="selectMode === 'auto'">
        <RoomGroupAutoSelect
          v-model:selected="selected"
          :stay-start-at="props.stayStartAt"
          :stay-end-at="props.stayEndAt"
        />
      </div>
      <div v-else>
        <RoomGroupManualSelect
          v-model:selected="selected"
          :stay-start-at="props.stayStartAt"
          :stay-end-at="props.stayEndAt"
        />
      </div>
    </div>
    <div class="col-12 col-md-6 q-pa-sm">
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
import { computed, ref } from "vue";
import { Room } from "src/schema/room";
import { Reservation } from "src/schema/reservation";
import RoomGroupManualSelect from "components/room-group/RoomGroupManualSelect.vue";
import RoomGroupAutoSelect from "components/room-group/RoomGroupAutoSelect.vue";

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

const selectMode = ref("auto");
</script>
