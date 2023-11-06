<template>
  <v-container>
    <v-row
      align="center"
      justify="center"
      style="height:100vh"
      dense
    >
      <v-col
        cols="12"
        class="fill-height d-flex flex-column justify-center align-center"
      >
        <LoginCard />
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup>
import { useRouter } from "vue-router"
import { useAuthStore } from "@/store/auth.js"

import LoginCard from "@/components/auth/LoginCard.vue"

const router = useRouter()
const authStore = useAuthStore()

if (authStore.isLoggedIn)
  router.push({ name: "Home" })
</script>
