<template>
  <q-card bordered>
    <q-card-section>
      <div class="text-h6">Current User Info</div>
    </q-card-section>

    <q-card-section>
      <div>
        Email: <span>{{ authStore.user?.email }}</span>
      </div>
      <div>
        Name: <span>{{ authStore.user?.name }}</span>
      </div>
      <div>
        Role: <span>{{ authStore.user?.role }}</span>
      </div>
      <div>
        CreatedAt: <span>{{ formatDateTime(authStore.user?.createdAt) }}</span>
      </div>
    </q-card-section>
  </q-card>
</template>

<script setup lang="ts">
import { useAuthStore } from "stores/auth"
import { formatDateTime } from "src/util/format-util"

const authStore = useAuthStore()
</script>
