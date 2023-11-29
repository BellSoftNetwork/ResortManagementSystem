<template>
  <q-page padding class="q-gutter-sm">
    <div class="row">
      <div class="col">
        <RoomDetailCard :id="id" />
      </div>
    </div>

    <div class="row">
      <div class="col">
        <RoomHistoryTable :id="id" />
      </div>
    </div>
  </q-page>
</template>

<script setup lang="ts">
import RoomDetailCard from "components/room/RoomDetailCard.vue";
import RoomHistoryTable from "components/room/RoomHistoryTable.vue";
import { useRoute } from "vue-router";

const route = useRoute();

const id = Number.parseInt(route.params.id as string);
</script>
