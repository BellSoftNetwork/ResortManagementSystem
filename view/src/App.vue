<template>
  <router-view />
</template>

<script setup>
import { useRouter } from "vue-router"
import { useAuthStore } from "@/store/auth.js"

const router = useRouter()
const authStore = useAuthStore()

authStore.loadAccountInfo().finally(() => {
  if (!authStore.isLoggedIn)
    router.push({ name: "Login" })
})
</script>
