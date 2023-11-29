<template>
  <slot :dialog="dialog">
    <q-btn @click="dialog.isOpen = true"> 로그아웃</q-btn>
  </slot>

  <q-dialog v-model="dialog.isOpen" :persistent="status.isProgress">
    <q-card>
      <q-card-section class="text-h6"> 로그아웃</q-card-section>

      <q-card-section>
        로그아웃 시 모든 기능을 이용하실 수 없습니다.
      </q-card-section>

      <q-card-actions align="right">
        <q-btn
          flat
          label="로그인 유지"
          color="primary"
          v-close-popup
          :disable="status.isProgress"
        />
        <q-btn
          flat
          label="로그아웃"
          color="red"
          v-close-popup
          @click="logout"
          :loading="status.isProgress"
        />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { useRouter } from "vue-router";
import { useAuthStore } from "stores/auth";
import { useQuasar } from "quasar";

defineExpose({
  openDialog,
});

const router = useRouter();
const authStore = useAuthStore();
const $q = useQuasar();

const dialog = ref({
  isOpen: false,
});
const status = ref({
  isProgress: false,
});

function openDialog() {
  dialog.value.isOpen = true;
}

function logout() {
  status.value.isProgress = true;

  authStore
    .logout()
    .then(() => {
      authStore.loadAccountInfo().finally(() => {
        router.push({ name: "Login" });
      });
    })
    .catch((error) => {
      $q.notify({
        message: error.response.data.message,
        type: "negative",
        actions: [
          {
            icon: "close",
            color: "white",
            round: true,
          },
        ],
      });
    })
    .finally(() => {
      status.value.isProgress = false;
    });
}
</script>
