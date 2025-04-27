<template>
  <q-card flat bordered>
    <q-card-section>
      <div class="text-h6">매출액</div>
      <q-inner-loading :showing="isLoading">
        <q-spinner-dots size="50px" color="primary" />
      </q-inner-loading>
      <div class="text-h4 q-mt-sm text-primary" :class="{ invisible: isLoading }">
        {{ formatPrice(totalSales) }}
      </div>
    </q-card-section>
  </q-card>
</template>

<script setup lang="ts">
import { defineProps } from "vue";
import { formatPrice } from "src/util/format-util";

defineProps({
  totalSales: {
    type: Number,
    required: true,
  },
  isLoading: {
    type: Boolean,
    default: false,
  },
});
</script>
