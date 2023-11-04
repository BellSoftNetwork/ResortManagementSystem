<template>
  <v-card
    title="Login"
    color="primary"
    min-width="500"
    subtitle="Resort Management System"
    variant="outlined"
    :loading="status.isProgress"
  >
    <v-card-text>
      <v-form fast-fail @submit.prevent ref="form" v-model="status.isValid">
        <v-text-field
          v-model="account.email"
          label="이메일"
          :rules="rules.email"
          :disabled="status.isProgress"
          required
        ></v-text-field>

        <v-text-field
          v-model="account.password"
          type="password"
          label="비밀번호"
          :rules="rules.password"
          :disabled="status.isProgress"
          required
        ></v-text-field>

        <v-btn
          type="submit"
          color="primary"
          class="mt-2"
          block
          @click="login"
          :loading="status.isProgress"
          :disabled="status.isProgress || !status.isValid"
        >
          로그인
        </v-btn>
      </v-form>
    </v-card-text>
  </v-card>

  <v-snackbar
    v-model="status.isError"
  >
    로그인 실패 ({{ status.errorMessage }})

    <template v-slot:actions>
      <v-btn
        color="pink"
        variant="text"
        @click="status.isError = false"
      >
        닫기
      </v-btn>
    </template>
  </v-snackbar>
</template>

<script setup>
import { ref } from "vue"
import { useRouter } from "vue-router"
import { useAuthStore } from "@/store/auth.js"

const router = useRouter()
const authStore = useAuthStore()

const status = ref({
  isValid: false,
  isProgress: false,
  isError: false,
  errorMessage: null,
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
  status.value.isError = false
  status.value.errorMessage = null

  authStore.login(account.value.email, account.value.password)
    .then(() => {
      router.push({ name: "Home" })
    })
    .catch((error) => {
      status.value.errorMessage = error.response.data.message
      status.value.isError = true
    })
    .finally(() => {
      status.value.isProgress = false
    })
}
</script>
