<template>
  <v-dialog
    v-model="status.isDialogActive"
    :persistent="status.isProgress"
    width="500"
  >
    <template v-slot:activator="{ props }">
      <v-btn
        v-bind="props"
        color="red"
        variant="text"
      >
        <slot>예약 수단 삭제</slot>
      </v-btn>
    </template>

    <template v-slot:default>
      <v-card
        title="예약 수단 삭제"
        :loading="status.isProgress"
        :disabled="status.isProgress"
      >
        <v-card-text>
          <v-form
            @submit.prevent
            fast-fail
            ref="form"
          >
            <v-text-field
              v-model="reservationMethod.name"
              label="이름"
              :readonly="true"
            ></v-text-field>

            <v-text-field
              v-model="reservationMethod.commissionRatePercent"
              label="수수료율"
              append-inner-icon="fa-solid fa-percent"
              :readonly="true"
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
            text="삭제"
            color="red"
            @click="updateReservationMethod"
          ></v-btn>
        </v-card-actions>
      </v-card>
    </template>
  </v-dialog>

  <v-snackbar
    v-model="status.isError"
  >
    예약 수단 삭제 실패 ({{ status.errorMessage }})

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

const emit = defineEmits(["complete"])
const props = defineProps({
  reservationMethod: Object,
})
const status = ref({
  isDialogActive: false,
  isValid: false,
  isProgress: false,
  isError: false,
  errorMessage: null,
})
const reservationMethod = ref({
  name: props.reservationMethod.name,
  commissionRatePercent: props.reservationMethod.commissionRate * 100,
})

function updateReservationMethod() {
  status.value.isProgress = true

  axios.delete(`/api/v1/reservation-methods/${props.reservationMethod.id}`)
    .then(() => {
      emit("complete")
      status.value.isDialogActive = false
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
