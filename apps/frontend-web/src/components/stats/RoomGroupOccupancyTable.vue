<template>
  <q-card flat bordered>
    <q-card-section>
      <div class="text-h6">객실 그룹별 평균 점유율</div>
      <q-inner-loading :showing="isLoading">
        <q-spinner-dots size="50px" color="primary" />
      </q-inner-loading>
      <div :class="{ invisible: isLoading }">
        <q-table :rows="roomGroupOccupancyStats" :columns="columns" row-key="roomGroup" flat bordered />
      </div>
    </q-card-section>
  </q-card>
</template>

<script setup lang="ts">
import { defineProps } from "vue";

const columns = [
  {
    name: "roomGroup",
    label: "객실 그룹",
    field: "roomGroup",
    align: "left",
  },
  {
    name: "roomCount",
    label: "객실 수",
    field: "roomCount",
    align: "left",
  },
  {
    name: "occupancyRate",
    label: "점유율",
    field: "occupancyRate",
    align: "left",
    format: (val) => `${val}%`,
  },
];

defineProps({
  roomGroupOccupancyStats: {
    type: Array,
    required: true,
  },
  isLoading: {
    type: Boolean,
    default: false,
  },
});
</script>
