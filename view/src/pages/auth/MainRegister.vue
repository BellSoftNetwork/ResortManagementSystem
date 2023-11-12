<template>
  <q-layout>
    <q-page-container>
      <q-page class="bg-primary window-height window-width row justify-center items-center">
        <div class="column">
          <div class="row">
            <h5 class="text-h5 text-white q-my-md">Resort Management System</h5>
          </div>
          <div class="row">
            <RegisterCard />
          </div>
        </div>
      </q-page>
    </q-page-container>
  </q-layout>
</template>

<script setup>
import { useRouter } from "vue-router"
import { useAuthStore } from "stores/auth.js"

import RegisterCard from "components/auth/RegisterCard.vue"
import { useAppConfigStore } from "stores/appConfig"
import { onBeforeMount } from "vue"

const router = useRouter()
const authStore = useAuthStore()
const appConfigStore = useAppConfigStore()


onBeforeMount(() => {
  if (authStore.isLoggedIn)
    router.push({ name: "Home" })

  appConfigStore.loadAppConfig(true).finally(() => {
    if (!appConfigStore.config.isAvailableRegistration)
      router.push({ name: "Login" })
  })
})
</script>
