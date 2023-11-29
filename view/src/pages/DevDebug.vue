<template>
  <q-page padding>
    <div class="row items-start q-gutter-md">
      <ApiTestCard />
      <UserInfoCard />
    </div>
  </q-page>
</template>

<script setup lang="ts">
import ApiTestCard from "components/test/ApiTestCard.vue";
import UserInfoCard from "components/test/UserInfoCard.vue";
</script>
