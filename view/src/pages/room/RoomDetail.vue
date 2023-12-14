<template>
  <q-page padding>
    <div v-if="room === null">
      <q-inner-loading :showing="true">
        <q-spinner-gears size="20vh" color="primary" />
      </q-inner-loading>
    </div>
    <div v-else class="q-gutter-sm">
      <div class="row">
        <div class="col">
          <RoomDetailCard :room="room" />
        </div>
      </div>

      <div class="row">
        <div class="col">
          <RoomHistoryDynamicTable :id="id" />
        </div>
      </div>
    </div>
  </q-page>
</template>

<script setup lang="ts">
import RoomDetailCard from "components/room/RoomDetailCard.vue";
import RoomHistoryDynamicTable from "components/room/RoomHistoryDynamicTable.vue";
import { useRoute, useRouter } from "vue-router";
import { fetchRoom } from "src/api/v1/room";
import { onBeforeMount, ref } from "vue";
import { Room } from "src/schema/room";

const route = useRoute();
const router = useRouter();

const id = Number.parseInt(route.params.id as string);
const room = ref<Room | null>(null);

function fetchData() {
  room.value = null;

  return fetchRoom(id)
    .then((response) => {
      room.value = response.value;
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
