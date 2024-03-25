<template>
  <q-page padding>
    <div v-if="reservation === null">
      <q-inner-loading :showing="true">
        <q-spinner-gears size="20vh" color="primary" />
      </q-inner-loading>
    </div>
    <div v-else class="q-gutter-sm">
      <div class="row">
        <div class="col">
          <ReservationEditor :reservation="reservation" mode="view" reservation-type="MONTHLY_RENT" />
        </div>
      </div>

      <div class="row">
        <div class="col">
          <ReservationHistoryDynamicTable :id="id" />
        </div>
      </div>
    </div>
  </q-page>
</template>

<script setup lang="ts">
import ReservationHistoryDynamicTable from "components/reservation/ReservationHistoryDynamicTable.vue";
import { useRoute, useRouter } from "vue-router";
import { onBeforeMount, ref } from "vue";
import { Reservation } from "src/schema/reservation";
import { fetchReservation } from "src/api/v1/reservation";
import ReservationEditor from "components/reservation/ReservationEditor.vue";

const route = useRoute();
const router = useRouter();

const id = Number.parseInt(route.params.id as string);
const reservation = ref<Reservation | null>(null);

function fetchData() {
  reservation.value = null;

  return fetchReservation(id)
    .then((response) => {
      reservation.value = response.value;
    })
    .catch((error) => {
      if (error.response.status === 404) router.push({ name: "ErrorNotFound" });

      console.log(error);
    });
}

onBeforeMount(() => {
  fetchData();
});
</script>
