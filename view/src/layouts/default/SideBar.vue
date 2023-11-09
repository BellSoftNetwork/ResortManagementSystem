<template>
  <v-navigation-drawer v-model="drawer">
    <v-sheet
      color="grey-lighten-4"
      class="pa-4"
    >
      <v-toolbar-title>
        <v-icon icon="fas fa-hotel" />&nbsp;리조트 관리 시스템
      </v-toolbar-title>
      <br />
      <div>{{ authStore.user.email }}</div>
    </v-sheet>

    <v-divider></v-divider>

    <v-list>
      <v-list-item
        v-for="[icon, text, routeName] in links"
        :key="icon"
        link
        :to="{name: routeName}"
      >
        <template v-slot:prepend>
          <v-icon>{{ icon }}</v-icon>
        </template>

        <v-list-item-title>{{ text }}</v-list-item-title>
      </v-list-item>
    </v-list>

    <template v-slot:append>
      <div class="pa-2">
        <LogoutDialog />
      </div>
    </template>
  </v-navigation-drawer>
</template>

<script setup>
import { ref } from "vue"
import { useAuthStore } from "@/store/auth.js"

import LogoutDialog from "@/components/auth/LogoutDialog.vue"

const authStore = useAuthStore()

const drawer = ref(null)
const links = [
  ["fa-solid fa-table-columns", "대시보드", "Home"],
  ["fa-solid fa-comment-dollar", "예약 수단", "ReservationMethods"],
]

if (authStore.isAdminRole)
  links.push(["fa-solid fa-id-card", "계정 관리", "AdminAccounts"])
</script>
