<template>
  <v-dialog
    width="500"
    :persistent="status.isProgress"
  >
    <template v-slot:activator="{ props }">
      <v-btn
        v-bind="props"
        color="primary"
        text="로그아웃"
        block
      ></v-btn>
    </template>

    <template v-slot:default="{ isActive }">
      <v-card
        title="로그아웃"
        :loading="status.isProgress"
        :disabled="status.isProgress"
      >
        <v-card-text>
          로그아웃 시 모든 기능일 이용하실 수 없습니다.
        </v-card-text>

        <v-card-actions>
          <v-spacer></v-spacer>

          <v-btn
            text="로그인 유지"
            color="primary"
            @click="isActive.value = false"
          ></v-btn>

          <v-btn
            text="로그아웃"
            color="red"
            @click="logout"
          ></v-btn>
        </v-card-actions>
      </v-card>
    </template>
  </v-dialog>

  <v-snackbar
    v-model="status.isError"
  >
    로그아웃 실패 ({{ status.errorMessage }})

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
  isProgress: false,
  isError: false,
  errorMessage: null,
})

function logout() {
  status.value.isProgress = true

  authStore.logout()
    .then(() => {
      authStore.loadAccountInfo().finally(() => {
        router.push({ name: "Login" })
      })
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
