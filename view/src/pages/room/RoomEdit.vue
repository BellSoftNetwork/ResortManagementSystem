<template>
  <q-page padding>
    <div v-if="room === null">
      <q-inner-loading :showing="true">
        <q-spinner-gears size="20vh" color="primary" />
      </q-inner-loading>
    </div>
    <div v-else>
      <RoomEditCard :room="room" />
    </div>
  </q-page>
</template>

<script setup lang="ts">
import RoomEditCard from "components/room/RoomEditCard.vue";
import { useRoute, useRouter } from "vue-router";
import { onBeforeMount, ref } from "vue";
import { Room } from "src/schema/room";
import { fetchRoom } from "src/api/v1/room";

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
