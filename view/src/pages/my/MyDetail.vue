<template>
  <q-page padding>
    <div v-if="authStore.user === null">
      <p>알 수 없는 문제</p>
    </div>
    <div v-else>
      <MyDetailCard :user="authStore.user" />
    </div>
  </q-page>
</template>

<script setup lang="ts">
import MyDetailCard from "components/my/MyDetailCard.vue";
import { useAuthStore } from "stores/auth";

const authStore = useAuthStore();
</script>
