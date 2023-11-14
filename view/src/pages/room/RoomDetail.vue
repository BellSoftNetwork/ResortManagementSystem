<template>
  <q-page padding class="q-gutter-sm">
    <div class="row">
      <div class="col">
        <RoomDetailCard :id="$route.params.id" />
      </div>
    </div>

    <div class="row">
      <div class="col">
        <RoomHistoryTable :id="$route.params.id" />
      </div>
    </div>
  </q-page>
</template>

<script setup>
import RoomDetailCard from "components/room/RoomDetailCard.vue"
import RoomHistoryTable from "components/room/RoomHistoryTable.vue"
</script>
