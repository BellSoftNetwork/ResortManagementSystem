<template>
  <q-page padding>
    <div v-if="roomGroup === null">
      <q-inner-loading :showing="true">
        <q-spinner-gears size="20vh" color="primary" />
      </q-inner-loading>
    </div>
    <div v-else>
      <RoomGroupEditCard :room-group="roomGroup" />
    </div>
  </q-page>
</template>

<script setup lang="ts">
import RoomGroupEditCard from "components/room-group/RoomGroupEditCard.vue";
import { useRoute, useRouter } from "vue-router";
import { onBeforeMount, ref } from "vue";
import RoomGroupDetail from "pages/room-group/RoomGroupDetail.vue";
import { fetchRoomGroup } from "src/api/v1/room-group";

const route = useRoute();
const router = useRouter();

const id = Number.parseInt(route.params.id as string);
const roomGroup = ref<RoomGroupDetail | null>(null);

function fetchData() {
  roomGroup.value = null;

  return fetchRoomGroup(id)
    .then((response) => {
      roomGroup.value = response.value;
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
