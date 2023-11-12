<template>
  <q-card
    :loading="status.isProgress"
    class="q-pa-lg shadow-1" style="min-width: 400px;"
    square
    bordered
  >
    <q-card-section>
      <div class="text-h5 text-center">Register</div>
    </q-card-section>

    <q-form @submit="register">
      <q-card-section>
        <q-input
          v-model="formData.email"
          label="이메일"
          :rules="rules.email"
          :disabled="status.isProgress"
          required
          autofocus
        ></q-input>

        <q-input
          v-model="formData.name"
          label="이름"
          :rules="rules.name"
          :disabled="status.isProgress"
          required
        ></q-input>

        <q-input
          v-model="formData.password"
          type="password"
          label="비밀번호"
          :rules="rules.password"
          :disabled="status.isProgress"
          required
        ></q-input>

        <q-input
          v-model="formData.passwordConfirm"
          type="password"
          label="비밀번호 확인"
          :rules="rules.passwordConfirm"
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
        <q-btn
          :to="{ name: 'Login' }"
          class="text-grey-6 mt-2 full-width"
          flat
        >
          로그인
        </q-btn>
      </q-card-section>
    </q-form>
  </q-card>
</template>

<script setup>
import { ref } from "vue"
import { useRouter } from "vue-router"
import { useAuthStore } from "stores/auth.js"
import { useQuasar } from "quasar"
import { api } from "boot/axios"

const router = useRouter()
const $q = useQuasar()
const authStore = useAuthStore()

const status = ref({
  isProgress: false,
})
const formData = ref({
  email: "",
  name: "",
  password: "",
  passwordConfirm: "",
})
const rules = {
  name: [value => (value.length >= 2 && value.length <= 20) || "2~20 글자가 필요합니다"],
  email: [value => /^\w+([.-]?\w+)*@\w+([.-]?\w+)*(\.\w{2,3})+$/.test(value) || "이메일이 유효하지 않습니다."],
  password: [value => (value?.length >= 8 && value?.length <= 20) || "비밀번호는 8~20 글자가 필요합니다."],
  passwordConfirm: [value => (formData.value.password === value) || "비밀번호가 일치하지 않습니다."],
}

function register() {
  status.value.isProgress = true

  api.post(`/api/v1/auth/register`, formData.value)
    .then(() => {
      authStore.login(formData.value.email, formData.value.password)
        .then(() => {
          router.push({ name: "Home" })
        })
        .catch((error) => {
          $q.notify({
            message: `로그인 실패 (${error.response.data.message})`,
            type: "negative",
            actions: [
              {
                icon: "close", color: "white", round: true,
              },
            ],
          })

          router.push({ name: "Login" })
        })
        .finally(() => {
          status.value.isProgress = false
        })
    })
    .catch((error) => {
      $q.notify({
        message: `회원 가입 실패 (${error.response.data.message})`,
        type: "negative",
        actions: [
          {
            icon: "close", color: "white", round: true,
          },
        ],
      })
    })
    .finally(() => {
      status.value.isProgress = false
    })
}
</script>
