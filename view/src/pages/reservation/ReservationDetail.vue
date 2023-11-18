<template>
  <q-page padding class="q-gutter-sm">
    <div class="row">
      <div class="col">
        <ReservationDetailCard :id="$route.params.id" />
      </div>
    </div>

    <div class="row">
      <div class="col">
        <ReservationHistoryTable :id="$route.params.id" />
      </div>
    </div>
  </q-page>
</template>

<script setup>
import ReservationDetailCard from "components/reservation/ReservationDetailCard.vue"
import ReservationHistoryTable from "components/reservation/ReservationHistoryTable.vue"
</script>
