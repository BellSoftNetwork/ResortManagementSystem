<template>
  <q-card>
    <q-card-section class="text-h6">
      내 정보
    </q-card-section>

    <q-form @submit="update">
      <q-card-section>
        <q-input
          v-model="authStore.user.name"
          label="이름"
          :readonly="true"
        ></q-input>

        <q-input
          v-model="authStore.user.email"
          :readonly="true"
          label="이메일"
        ></q-input>

        <q-input
          v-model="formData.password"
          :rules="rules.password"
          type="password"
          label="비밀번호"
        ></q-input>

        <q-input
          v-model="formData.passwordConfirm"
          :rules="rules.passwordConfirm"
          type="password"
          label="비밀번호 확인"
        ></q-input>
      </q-card-section>

      <q-card-actions align="right">
        <q-btn
          :loading="status.isProgress"
          type="submit"
          color="red"
          label="수정"
          flat
        />
      </q-card-actions>
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
const authStore = useAuthStore()
const $q = useQuasar()

const status = ref({
  isProgress: false,
})
const formData = ref({
  password: "",
})
const rules = {
  password: [value => (value.length === 0 || (value.length >= 8 && value.length <= 20)) || "비밀번호는 8~20 글자가 필요합니다."],
  passwordConfirm: [value => (formData.value.password === value) || "비밀번호가 일치하지 않습니다."],
}

function update() {
  if (!isChanged()) {
    $q.notify({
      message: "수정된 항목이 없습니다.",
      type: "info",
      actions: [
        {
          icon: "close", color: "white", round: true,
        },
      ],
    })

    return
  }

  status.value.isProgress = true

  api.patch(`/api/v1/my`, patchedData())
    .then(() => {
      $q.notify({
        message: "정보가 정상적으로 변경되었습니다",
        type: "positive",
        actions: [
          {
            icon: "close", color: "white", round: true,
          },
        ],
      })

      authStore.logout()
        .then(() => {
          authStore.loadAccountInfo().finally(() => {
            router.push({ name: "Login" })
          })
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

function isChanged() {
  return formData.value.password.length > 0
}

function patchedData() {
  const patchData = {}

  if (formData.value.password.length > 0)
    patchData.password = formData.value.password

  return patchData
}
</script>
