<template>
  <q-card flat bordered>
    <q-card-section>
      <div class="text-h6">총 다녀간 인원</div>
      <q-inner-loading :showing="isLoading">
        <q-spinner-dots size="50px" color="primary" />
      </q-inner-loading>
      <div class="text-h4 q-mt-sm text-primary" :class="{ invisible: isLoading }">{{ totalGuests }}명</div>
    </q-card-section>
  </q-card>
</template>

<script setup lang="ts">
import { defineProps } from "vue";

defineProps({
  totalGuests: {
    type: Number,
    required: true,
  },
  isLoading: {
    type: Boolean,
    default: false,
  },
});
</script>
