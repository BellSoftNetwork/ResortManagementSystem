<template>
  <q-card
    :loading="status.isProgress"
    class="q-pa-lg shadow-1" style="min-width: 400px;"
    square
    bordered
  >
    <q-card-section>
      <div class="text-h5 text-center">Login</div>
    </q-card-section>

    <q-card-section>
      <q-form @submit="login">
        <q-input
          v-model="account.email"
          label="이메일"
          :rules="rules.email"
          :disabled="status.isProgress"
          required
          autofocus
        ></q-input>

        <q-input
          v-model="account.password"
          type="password"
          label="비밀번호"
          :rules="rules.password"
          :disabled="status.isProgress"
          required
        ></q-input>

        <q-btn
          type="submit"
          color="primary"
          class="mt-2 full-width"
          :loading="status.isProgress"
          :disabled="status.isProgress"
        >
          로그인
        </q-btn>
      </q-form>
    </q-card-section>
  </q-card>
</template>

<script setup>
import { ref } from "vue"
import { useRouter } from "vue-router"
import { useAuthStore } from "stores/auth.js"
import { useQuasar } from "quasar"

const router = useRouter()
const authStore = useAuthStore()
const $q = useQuasar()

const status = ref({
  isValid: false,
  isProgress: false,
})
const account = ref({
  email: "",
  password: "",
})
const rules = {
  email: [value => /^\w+([.-]?\w+)*@\w+([.-]?\w+)*(\.\w{2,3})+$/.test(value) || "이메일이 유효하지 않습니다."],
  password: [value => (value?.length >= 8 && value?.length <= 20) || "비밀번호는 8~20 글자가 필요합니다."],
}

function login() {
  status.value.isProgress = true

  authStore.login(account.value.email, account.value.password)
    .then(() => {
      router.push({ name: "Home" })
    })
    .catch((error) => {
      $q.notify({
        message: error.response.data.message,
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
