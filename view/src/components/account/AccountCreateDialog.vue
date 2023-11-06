<template>
  <v-dialog
    v-model="status.isDialogActive"
    :persistent="status.isProgress"
    width="500"
  >
    <template v-slot:activator="{ props }">
      <v-btn
        v-bind="props"
        color="primary"
        text="계정 추가"
        block
      ></v-btn>
    </template>

    <template v-slot:default>
      <v-card
        title="계정 추가"
        :loading="status.isProgress"
        :disabled="status.isProgress"
      >
        <v-card-text>
          <v-form
            v-model="status.isValid"
            @submit.prevent
            fast-fail
            ref="form"
          >
            <v-text-field
              v-model="account.name"
              label="이름"
              :rules="rules.name"
              required
            ></v-text-field>

            <v-text-field
              v-model="account.email"
              label="이메일"
              :rules="rules.email"
              required
            ></v-text-field>

            <v-text-field
              v-model="account.password"
              type="password"
              label="비밀번호"
              :rules="rules.password"
              required
            ></v-text-field>

            <v-select
              v-model="account.role"
              label="권한"
              :items="rules.role"
              required
            ></v-select>
          </v-form>
        </v-card-text>

        <v-card-actions>
          <v-spacer></v-spacer>

          <v-btn
            text="취소"
            color="primary"
            @click="status.isDialogActive = false"
          ></v-btn>

          <v-btn
            text="계정 추가"
            color="red"
            @click="createAccount"
          ></v-btn>
        </v-card-actions>
      </v-card>
    </template>
  </v-dialog>

  <v-snackbar
    v-model="status.isError"
  >
    계정 추가 실패 ({{ status.errorMessage }})

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
import axios from "@/modules/axios-wrapper"

const router = useRouter()
const authStore = useAuthStore()

const emit = defineEmits(["created"])
const status = ref({
  isDialogActive: false,
  isValid: false,
  isProgress: false,
  isError: false,
  errorMessage: null,
})
const account = ref({
  name: "",
  email: "",
  password: "",
  role: "NORMAL",
})
const rules = {
  name: [value => (value.length >= 2 && value.length <= 20) || "2~20 글자가 필요합니다"],
  email: [value => /^\w+([.-]?\w+)*@\w+([.-]?\w+)*(\.\w{2,3})+$/.test(value) || "이메일이 유효하지 않습니다."],
  password: [value => (value.length >= 8 && value.length <= 20) || "비밀번호는 8~20 글자가 필요합니다."],
  role: [
    { title: "일반", value: "NORMAL" },
  ],
}

if (authStore.isSuperAdminRole)
  rules.role.push({ title: "관리자", value: "ADMIN" })


function createAccount() {
  status.value.isProgress = true

  axios.post("/api/v1/admin/accounts", account.value)
    .then(() => {
      emit("created")
      status.value.isDialogActive = false

      resetForm()
    })
    .catch((error) => {
      status.value.errorMessage = error.response.data.message
      status.value.isError = true
    })
    .finally(() => {
      status.value.isProgress = false
    })
}

function resetForm() {
  account.value.name = ""
  account.value.email = ""
  account.value.password = ""
  account.value.role = "NORMAL"
}
</script>
