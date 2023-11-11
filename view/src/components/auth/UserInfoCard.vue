<template>
  <q-card bordered>
    <q-card-section>
      <div class="text-h6">Current User Info</div>
    </q-card-section>

    <q-card-section>
      <div>Email: <span>{{ authStore.user.email }}</span></div>
      <div>Name: <span>{{ authStore.user.name }}</span></div>
      <div>Role: <span>{{ authStore.user.role }}</span></div>
      <div>CreatedAt: <span>{{ dayjs(authStore.user.createdAt).format("YYYY-MM-DD HH:mm:ss") }}</span></div>
    </q-card-section>
  </q-card>
</template>

<script setup>
import dayjs from "dayjs"
import { useAuthStore } from "stores/auth.js"

const authStore = useAuthStore()
</script>
