<template>
  <q-page padding>
    <ReservationListDynamicTable title="숙박 중인 달방" reservation-type="MONTHLY_RENT" />
  </q-page>
</template>

<script setup lang="ts">
import ReservationListDynamicTable from "components/reservation/ReservationListDynamicTable.vue";
</script>
