<template>
  <q-page padding>
    <ReservationListDynamicTable reservation-type="STAY" />
  </q-page>
</template>

<script setup lang="ts">
import ReservationListDynamicTable from "components/reservation/ReservationListDynamicTable.vue";
</script>
