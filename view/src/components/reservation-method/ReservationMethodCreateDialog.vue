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
        text="예약 수단 추가"
        block
      ></v-btn>
    </template>

    <template v-slot:default>
      <v-card
        title="예약 수단 추가"
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
              v-model="reservationMethod.name"
              label="이름"
              :rules="rules.name"
              required
            ></v-text-field>

            <v-text-field
              v-model="reservationMethod.commissionRatePercent"
              label="수수료율"
              :rules="rules.commissionRatePercent"
              append-inner-icon="fa-solid fa-percent"
              type="number"
              min="0"
              max="100"
              required
            ></v-text-field>
          </v-form>
        </v-card-text>

        <v-card-actions>
          <v-spacer></v-spacer>

          <v-btn
            text="취소"
            @click="status.isDialogActive = false"
          ></v-btn>

          <v-btn
            text="추가"
            color="primary"
            @click="createReservationMethod"
            :disabled="!status.isValid"
          ></v-btn>
        </v-card-actions>
      </v-card>
    </template>
  </v-dialog>

  <v-snackbar
    v-model="status.isError"
  >
    예약 수단 추가 실패 ({{ status.errorMessage }})

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
const reservationMethod = ref({
  name: "",
  commissionRatePercent: 0,
})
const rules = {
  name: [value => (value.length >= 2 && value.length <= 20) || "2~20 글자가 필요합니다"],
  commissionRatePercent: [value => (value >= 0 && value <= 100) || "수수료율이 유효하지 않습니다."],
}

function createReservationMethod() {
  status.value.isProgress = true

  const data = {
    name: reservationMethod.value.name,
    commissionRate: reservationMethod.value.commissionRatePercent / 100,
  }

  axios.post("/api/v1/reservation-methods", data)
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
  reservationMethod.value.name = ""
  reservationMethod.value.commissionRate = 0
}
</script>
