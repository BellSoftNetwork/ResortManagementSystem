<template>
  <q-card flat bordered>
    <q-card-section>
      <div class="text-h6">예약 건수</div>
      <q-inner-loading :showing="isLoading">
        <q-spinner-dots size="50px" color="primary" />
      </q-inner-loading>
      <div class="text-h4 q-mt-sm text-primary" :class="{ invisible: isLoading }">{{ totalReservations }}건</div>
    </q-card-section>
  </q-card>
</template>

<script setup lang="ts">
import { defineProps } from "vue";

defineProps({
  totalReservations: {
    type: Number,
    required: true,
  },
  isLoading: {
    type: Boolean,
    default: false,
  },
});
</script>
