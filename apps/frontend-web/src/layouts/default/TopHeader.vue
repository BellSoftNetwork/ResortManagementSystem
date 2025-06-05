<template>
  <q-header elevated class="bg-white text-grey-8" height-hint="64">
    <q-toolbar>
      <q-btn @click="emit('toggleLeftDrawer')" icon="menu" class="q-mr-sm" flat dense round aria-label="Menu" />

      <q-toolbar-title shrink class="row items-center no-wrap">
        <span v-if="$q.screen.gt.xs" class="q-ml-sm">Resort Management System</span>
        <span v-else class="q-ml-sm">RMS</span>
        <div v-if="env !== 'prod'">
          &nbsp;
          <q-badge v-if="env === 'local'" align="middle" color="primary">Local</q-badge>
          <q-badge v-else-if="env === 'dev'" align="middle" color="secondary">Dev</q-badge>
          <q-badge v-else-if="env === 'staging'" e align="middle" color="warning">Staging</q-badge>
        </div>
      </q-toolbar-title>

      <q-space />

      <div class="q-gutter-sm row items-center no-wrap">
        <q-btn round flat>
          <q-avatar size="26px">
            <img :src="authStore.user?.profileImageUrl" alt="Profile Image" />
          </q-avatar>
          <q-menu auto-close>
            <q-list>
              <q-item clickable :to="{ name: 'MyDetail' }">
                <q-item-section>내 정보</q-item-section>
              </q-item>
              <q-separator />
              <q-item clickable @click="openDialog">
                <q-item-section>로그아웃</q-item-section>
              </q-item>
            </q-list>
          </q-menu>

          <q-tooltip>{{ authStore.user?.name }}</q-tooltip>
        </q-btn>
      </div>
    </q-toolbar>

    <TopTab v-if="$q.screen.lt.md" />

    <LogoutDialog ref="logoutDialogRef">
      <template></template>
    </LogoutDialog>
  </q-header>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useAuthStore } from "stores/auth";
import LogoutDialog from "components/auth/LogoutDialog.vue";
import TopTab from "layouts/default/TopTab.vue";

const emit = defineEmits(["toggleLeftDrawer"]);
const authStore = useAuthStore();
const logoutDialogRef = ref();
const env = ref<"local" | "dev" | "staging" | "prod">("prod");

function openDialog() {
  logoutDialogRef.value.openDialog();
}

function getEnv() {
  const hostname = window.location.hostname;

  switch (hostname) {
    case "rms.bellsoft.net":
      return "prod";
    case "staging.rms.bellsoft.net":
      return "staging";
    case "development.rms.bellsoft.net":
      return "dev";

    case "localhost":
    case "127.0.0.1":
    case "[::1]":
    default:
      return "local";
  }
}

onMounted(() => {
  env.value = getEnv();
});
</script>
