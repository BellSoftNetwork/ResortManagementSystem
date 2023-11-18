<template>
  <q-page padding>
    <ReservationEditCard :id="$route.params.id" />
  </q-page>
</template>

<script setup>
import ReservationEditCard from "components/reservation/ReservationEditCard.vue"
</script>
