<template>
  <q-card flat bordered>
    <q-card-section>
      <div class="row items-center">
        <q-btn icon="chevron_left" flat round @click="onPreviousMonth" />
        <div class="text-h6">{{ formattedMonthTitle }}</div>
        <q-btn icon="chevron_right" flat round @click="onNextMonth" />
      </div>
    </q-card-section>
  </q-card>
</template>

<script setup lang="ts">
import { computed, defineEmits, defineProps } from "vue";

const props = defineProps({
  selectedMonth: {
    type: String,
    required: true,
  },
});

const emit = defineEmits(["update:selectedMonth"]);

// 월별 통계 제목 포맷
const formattedMonthTitle = computed(() => {
  const year = props.selectedMonth.split("-")[0];
  const month = props.selectedMonth.split("-")[1];
  return `${year}년 ${month}월 통계`;
});

// 월 이동 함수
function onPreviousMonth() {
  emit("update:selectedMonth", "previous");
}

function onNextMonth() {
  emit("update:selectedMonth", "next");
}
</script>
