<template>
  <q-page padding>
    <div v-if="roomGroup === null">
      <q-inner-loading :showing="true">
        <q-spinner-gears size="20vh" color="primary" />
      </q-inner-loading>
    </div>
    <div v-else class="q-gutter-sm">
      <div class="row">
        <div class="col">
          <RoomGroupDetailCard :room-group="roomGroup" />
        </div>
      </div>

      <div class="row">
        <div class="col">
          <RoomListTable :room-last-stay-details="roomGroup.rooms" />
        </div>
      </div>
    </div>
  </q-page>
</template>

<script setup lang="ts">
import RoomGroupDetailCard from "components/room-group/RoomGroupDetailCard.vue";
import { useRoute, useRouter } from "vue-router";
import { onBeforeMount, ref } from "vue";
import { fetchRoomGroup } from "src/api/v1/room-group";
import RoomGroupDetail from "pages/room-group/RoomGroupDetail.vue";
import RoomListTable from "components/room/RoomListTable.vue";

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
