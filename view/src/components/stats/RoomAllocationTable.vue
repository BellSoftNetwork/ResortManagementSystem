<template>
  <q-card flat bordered>
    <q-card-section>
      <div class="text-h6">객실별 예약 배정 건수</div>
      <q-inner-loading :showing="isLoading">
        <q-spinner-dots size="50px" color="primary" />
      </q-inner-loading>
      <div :class="{ invisible: isLoading }">
        <q-table :rows="roomAllocationStats" :columns="columns" row-key="roomGroup" flat bordered />
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
    name: "count",
    label: "예약 건수",
    field: "count",
    align: "left",
  },
  {
    name: "rooms",
    label: "포함 객실",
    field: "rooms",
    align: "left",
    format: (val) => val.join(", "),
  },
];

defineProps({
  roomAllocationStats: {
    type: Array,
    required: true,
  },
  isLoading: {
    type: Boolean,
    default: false,
  },
});
</script>
