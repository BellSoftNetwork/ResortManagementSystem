<template>
  <q-layout>
    <q-page-container>
      <q-page class="bg-primary window-height window-width row justify-center items-center">
        <div class="column q-pa-md" style="width: 100%; max-width: 450px">
          <div class="row justify-center">
            <h5 class="text-h5 text-white q-my-md">Resort Management System</h5>
          </div>
          <div class="row justify-center">
            <LoginCard />
          </div>
        </div>
      </q-page>
    </q-page-container>
  </q-layout>
</template>

<script setup lang="ts">
import { useRouter } from "vue-router";
import { useAuthStore } from "stores/auth";

import LoginCard from "components/auth/LoginCard.vue";
import { onBeforeMount } from "vue";

const router = useRouter();
const authStore = useAuthStore();

onBeforeMount(() => {
  if (authStore.isLoggedIn) router.push({ name: "Home" });
});
</script>
