<template>
  <q-page padding class="q-gutter-sm">
    <div class="row">
      <div class="col">
        <ReservationDetailCard :id="id" />
      </div>
    </div>

    <div class="row">
      <div class="col">
        <ReservationHistoryTable :id="id" />
      </div>
    </div>
  </q-page>
</template>

<script setup lang="ts">
import ReservationDetailCard from "components/reservation/ReservationDetailCard.vue"
import ReservationHistoryTable from "components/reservation/ReservationHistoryTable.vue"
import { useRoute } from "vue-router"

const route = useRoute()

const id = Number.parseInt(route.params.id as string)
</script>
