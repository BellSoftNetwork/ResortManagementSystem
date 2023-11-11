<template>
  <q-layout>
    <q-page-container>
      <q-page class="bg-primary window-height window-width row justify-center items-center">
        <div class="column">
          <div class="row">
            <h5 class="text-h5 text-white q-my-md">Resort Management System</h5>
          </div>
          <div class="row">
            <LoginCard />
          </div>
        </div>
      </q-page>
    </q-page-container>
  </q-layout>
</template>

<script setup>
import { useRouter } from "vue-router"
import { useAuthStore } from "stores/auth.js"

import LoginCard from "components/auth/LoginCard.vue"

const router = useRouter()
const authStore = useAuthStore()

if (authStore.isLoggedIn)
  router.push({ name: "Home" })
</script>
