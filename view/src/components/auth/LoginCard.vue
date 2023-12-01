<template>
  <q-card
    :loading="status.isProgress"
    class="q-pa-lg shadow-1"
    style="min-width: 400px"
    square
    bordered
  >
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

      <q-card-section
        v-if="appConfigStore.config.isAvailableRegistration"
        class="text-center q-py-none"
      >
        <q-btn
          :to="{ name: 'Register' }"
          class="text-grey-6 mt-2 full-width"
          flat
        >
          회원 가입
        </q-btn>
      </q-card-section>
    </q-form>
  </q-card>
</template>

<script setup lang="ts">
import { onBeforeMount, ref } from "vue";
import { useRouter } from "vue-router";
import { useAuthStore } from "stores/auth";
import { useQuasar } from "quasar";
import { useAppConfigStore } from "stores/app-config";
import { userStaticRules } from "src/schema/user";

const router = useRouter();
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
      router.push({ name: "Home" });
    })
    .catch((error) => {
      $q.notify({
        message: `로그인 실패 (${error.response.data.message})`,
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
