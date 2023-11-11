<template>
  <q-page padding>
    <div class="row items-start q-gutter-md">
      <ApiTestCard />
      <UserInfoCard />
    </div>
  </q-page>
</template>

<script setup>
import ApiTestCard from "components/test/ApiTestCard.vue"
import UserInfoCard from "components/auth/UserInfoCard.vue"
</script>
