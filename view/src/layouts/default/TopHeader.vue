<template>
  <q-header elevated class="bg-white text-grey-8" height-hint="64">
    <q-toolbar>
      <q-btn
        @click="emit('toggleLeftDrawer')"
        icon="menu"
        class="q-mr-sm"
        flat
        dense
        round
        aria-label="Menu"
      />

      <q-toolbar-title v-if="$q.screen.gt.xs" shrink class="row items-center no-wrap">
        <span class="q-ml-sm">Resort Management System</span>
      </q-toolbar-title>

      <q-space />

      <div class="q-gutter-sm row items-center no-wrap">
        <q-btn round flat>
          <q-avatar size="26px">
            <img :src="authStore.user.profileImageUrl" alt="Profile Image">
          </q-avatar>
          <q-menu auto-close>
            <q-list>
              <q-item clickable @click="openDialog">
                <q-item-section>로그아웃</q-item-section>
              </q-item>
            </q-list>
          </q-menu>

          <q-tooltip>{{ authStore.user.name }}</q-tooltip>
        </q-btn>
      </div>
    </q-toolbar>

    <LogoutDialog ref="logoutDialogRef">
      <template></template>
    </LogoutDialog>
  </q-header>
</template>

<script setup>
import { ref } from "vue"
import { useAuthStore } from "stores/auth.js"
import LogoutDialog from "components/auth/LogoutDialog.vue"

const emit = defineEmits(["toggleLeftDrawer"])
const authStore = useAuthStore()
const logoutDialogRef = ref()

function openDialog() {
  logoutDialogRef.value.openDialog()
}
</script>
