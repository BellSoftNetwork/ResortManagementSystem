<template>
  <q-card flat bordered>
    <q-card-section>
      <div class="text-h6">최근 1년간 예약 통계</div>
      <q-inner-loading :showing="isLoading">
        <q-spinner-dots size="50px" color="primary" />
      </q-inner-loading>
      <div class="q-mt-md" :class="{ invisible: isLoading }">
        <apexchart type="area" height="250" :options="chartOptions" :series="series"></apexchart>
      </div>
    </q-card-section>
  </q-card>
</template>

<script setup lang="ts">
import { computed, defineProps } from "vue";

const props = defineProps({
  chartData: {
    type: Array,
    required: true,
  },
  peopleCountData: {
    type: Array,
    default: () => [],
  },
  roomCountData: {
    type: Array,
    default: () => [],
  },
  isLoading: {
    type: Boolean,
    default: false,
  },
});

const series = computed(() => {
  const seriesData = [
    {
      name: "예약 건수",
      data: props.chartData.map((item) => item.value),
    },
  ];

  // Add people count data if available
  if (props.peopleCountData.length > 0) {
    seriesData.push({
      name: "다녀간 인원 수",
      data: props.peopleCountData.map((item) => item.value),
    });
  }

  // Add room count data if available
  if (props.roomCountData.length > 0) {
    seriesData.push({
      name: "예약된 객실 수",
      data: props.roomCountData.map((item) => item.value),
    });
  }

  return seriesData;
});

const chartOptions = computed(() => {
  return {
    chart: {
      type: "area",
      toolbar: {
        show: false,
      },
      stacked: false,
    },
    stroke: {
      curve: "smooth",
      width: 2,
    },
    fill: {
      type: "gradient",
      gradient: {
        shadeIntensity: 1,
        opacityFrom: 0.7,
        opacityTo: 0.9,
        stops: [0, 90, 100],
      },
    },
    dataLabels: {
      enabled: false,
    },
    colors: ["#1976d2", "#00C853", "#FF6D00"],
    xaxis: {
      categories: props.chartData.map((item) => item.label),
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
  };
});
</script>
