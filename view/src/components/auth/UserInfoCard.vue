<template>
  <v-container>
    <v-card
      title="Current User Info"
      color="primary"
      variant="outlined"
    >
      <v-card-text>
        Email: <span>{{ authStore.user.email }}</span>
      </v-card-text>
    </v-card>
  </v-container>
</template>

<script setup>
import { useAuthStore } from "@/store/auth.js"

const authStore = useAuthStore()
</script>
