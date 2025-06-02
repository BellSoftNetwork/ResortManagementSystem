<template>
  <q-page padding>
    <ReservationListDynamicTable title="다가오는 예약" reservation-type="STAY" />
  </q-page>
</template>

<script setup lang="ts">
import ReservationListDynamicTable from "components/reservation/ReservationListDynamicTable.vue";
</script>
