<template>
  <q-card flat bordered>
    <q-card-section>
      <div class="text-h6">예약 건수 (전년 대비)</div>
      <q-inner-loading :showing="isLoading">
        <q-spinner-dots size="50px" color="primary" />
      </q-inner-loading>
      <div class="q-mt-md" :class="{ invisible: isLoading }">
        <apexchart type="line" height="250" :options="chartOptions" :series="series"></apexchart>
      </div>
    </q-card-section>
  </q-card>
</template>

<script setup lang="ts">
import { computed, defineProps } from "vue";
import dayjs from "dayjs";

const props = defineProps({
  chartData: {
    type: Object,
    required: true,
    validator: (value) => {
      return value.currentYear && value.previousYear;
    },
  },
  isLoading: {
    type: Boolean,
    default: false,
  },
});

const series = computed(() => {
  const currentYear = dayjs().year();
  const previousYear = currentYear - 1;

  return [
    {
      name: `${currentYear}년`,
      data: props.chartData.currentYear.map((item) => item.value),
    },
    {
      name: `${previousYear}년`,
      data: props.chartData.previousYear.map((item) => item.value),
    },
  ];
});

const chartOptions = computed(() => {
  return {
    chart: {
      type: "line",
      toolbar: {
        show: false,
      },
      zoom: {
        enabled: false,
      },
    },
    stroke: {
      curve: "smooth",
      width: [3, 2],
      dashArray: [0, 5],
    },
    dataLabels: {
      enabled: false,
    },
    colors: ["#00C853", "#FF9800"],
    xaxis: {
      categories: props.chartData.currentYear.map((item) => item.label + "월"),
      labels: {
        style: {
          fontSize: "12px",
        },
      },
    },
    tooltip: {
      y: {
        formatter: function (val) {
          return `${val}건`;
        },
      },
    },
    legend: {
      position: "top",
    },
  };
});
</script>
