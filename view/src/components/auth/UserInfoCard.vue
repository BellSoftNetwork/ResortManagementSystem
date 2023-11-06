<template>
  <v-container>
    <v-card
      title="Current User Info"
      color="primary"
      variant="outlined"
    >
      <v-card-text>
        <p>Email: <span>{{ authStore.user.email }}</span></p>
        <p>Name: <span>{{ authStore.user.name }}</span></p>
        <p>Role: <span>{{ authStore.user.role }}</span></p>
        <p>CreatedAt: <span>{{ dayjs(authStore.user.createdAt).format("YYYY-MM-DD HH:mm:ss") }}</span></p>
      </v-card-text>
    </v-card>
  </v-container>
</template>

<script setup>
import dayjs from "dayjs"
import { useAuthStore } from "@/store/auth.js"

const authStore = useAuthStore()
</script>
