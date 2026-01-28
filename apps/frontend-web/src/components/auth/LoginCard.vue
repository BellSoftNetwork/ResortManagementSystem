<template>
  <q-card :loading="status.isProgress" class="q-pa-lg shadow-1" style="width: 100%; max-width: 400px" square bordered>
    <q-card-section>
      <div class="text-h5 text-center">Login</div>
    </q-card-section>

    <q-form @submit="login">
      <q-card-section>
        <q-input
          v-model="formData.username"
          label="ID 또는 이메일"
          :rules="userStaticRules.username"
          :disabled="status.isProgress"
          required
          autofocus
        ></q-input>

        <q-input
          v-model="formData.password"
          type="password"
          label="비밀번호"
          :rules="userStaticRules.password"
          :disabled="status.isProgress"
          required
        ></q-input>
      </q-card-section>

      <q-card-section>
        <q-btn
          type="submit"
          color="primary"
          class="mt-2 full-width"
          :loading="status.isProgress"
          :disabled="status.isProgress"
        >
          로그인
        </q-btn>
      </q-card-section>

      <q-card-section v-if="appConfigStore.config?.isAvailableRegistration" class="text-center q-py-none">
        <q-btn :to="{ name: 'Register' }" class="text-grey-6 mt-2 full-width" flat>회원 가입</q-btn>
      </q-card-section>
    </q-form>
  </q-card>
</template>

<script setup lang="ts">
import { onBeforeMount, ref } from "vue";
import { useRouter, useRoute } from "vue-router";
import { useAuthStore } from "stores/auth";
import { useQuasar } from "quasar";
import { useAppConfigStore } from "stores/app-config";
import { userStaticRules } from "src/schema/user";

const router = useRouter();
const route = useRoute();
const $q = useQuasar();
const authStore = useAuthStore();
const appConfigStore = useAppConfigStore();

const status = ref({
  isProgress: false,
});
const formData = ref({
  username: "",
  password: "",
});

function login() {
  status.value.isProgress = true;

  authStore
    .login(formData.value)
    .then(() => {
      const redirectPath = route.query.redirect as string | undefined;
      if (redirectPath) {
        router.push(redirectPath);
      } else {
        router.push({ name: "Home" });
      }
    })
    .catch((error) => {
      // 서버 오류 (5xx) 또는 네트워크 오류 감지
      const status = error.response?.status;
      const isServerError = status && status >= 500;
      const isNetworkError = !error.response;

      let message: string;
      if (isServerError) {
        message = `서버에 문제가 발생했습니다 (오류 코드: ${status})`;
      } else if (isNetworkError) {
        message = "서버에 연결할 수 없습니다";
      } else {
        // 일반 인증 오류 (400, 401, 403 등)
        message = error.response?.data?.message || "로그인에 실패했습니다";
      }

      $q.notify({
        message: `로그인 실패: ${message}`,
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

onBeforeMount(() => {
  appConfigStore.loadAppConfig();
});
</script>
