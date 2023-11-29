<template>
  <q-page padding>
    <ReservationEditCard :id="id" />
  </q-page>
</template>

<script setup lang="ts">
import ReservationEditCard from "components/reservation/ReservationEditCard.vue";
import { useRoute } from "vue-router";

const route = useRoute();

const id = Number.parseInt(route.params.id as string);
</script>
