<template>
  <q-card
    :loading="status.isProgress"
    class="q-pa-lg shadow-1"
    style="min-width: 400px"
    square
    bordered
  >
    <q-card-section>
      <div class="text-h5 text-center">Register</div>
    </q-card-section>

    <q-form @submit="register">
      <q-card-section>
        <q-input
          v-model="formData.userId"
          label="계정 ID"
          :rules="userStaticRules.userId"
          :disabled="status.isProgress"
          required
          autofocus
        ></q-input>

        <q-input
          v-model="formData.email"
          label="이메일"
          :rules="userStaticRules.email"
          :disabled="status.isProgress"
          required
        ></q-input>

        <q-input
          v-model="formData.name"
          label="이름"
          :rules="userStaticRules.name"
          :disabled="status.isProgress"
          required
        ></q-input>

        <q-input
          v-model="formData.password"
          type="password"
          label="비밀번호"
          :rules="userStaticRules.password"
          :disabled="status.isProgress"
          required
        ></q-input>

        <q-input
          v-model="formData.passwordConfirm"
          type="password"
          label="비밀번호 확인"
          :rules="userDynamicRules.passwordConfirm(formData.password)"
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
          회원가입
        </q-btn>
      </q-card-section>

      <q-card-section class="text-center q-py-none">
        <q-btn :to="{ name: 'Login' }" class="text-grey-6 mt-2 full-width" flat>
          로그인
        </q-btn>
      </q-card-section>
    </q-form>
  </q-card>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { useRouter } from "vue-router";
import { useAuthStore } from "stores/auth";
import { useQuasar } from "quasar";
import { userDynamicRules, userStaticRules } from "src/schema/user";
import { postRegister } from "src/api/v1/auth";

const router = useRouter();
const $q = useQuasar();
const authStore = useAuthStore();

const status = ref({
  isProgress: false,
});
const formData = ref({
  userId: "",
  email: "",
  name: "",
  password: "",
  passwordConfirm: "",
});

function register() {
  status.value.isProgress = true;

  postRegister(formData.value)
    .then(() => {
      authStore
        .login({
          username: formData.value.email,
          password: formData.value.password,
        })
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

          router.push({ name: "Login" });
        })
        .finally(() => {
          status.value.isProgress = false;
        });
    })
    .catch((error) => {
      $q.notify({
        message: `회원 가입 실패 (${error.response.data.message})`,
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
