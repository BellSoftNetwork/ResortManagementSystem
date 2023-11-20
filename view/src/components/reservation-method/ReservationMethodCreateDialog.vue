<template>
  <slot :dialog="dialog">
    <q-btn @click="dialog.isOpen = true">
      추가
    </q-btn>
  </slot>

  <q-dialog
    v-model="dialog.isOpen"
    :persistent="status.isProgress"
    @beforeShow="resetForm"
  >
    <q-card style="width: 500px">
      <q-card-section class="text-h6">
        예약 수단 추가
      </q-card-section>

      <q-form @submit="create">
        <q-card-section>
          <q-input
            v-model="formData.name"
            :rules="rules.name"
            label="이름"
            required
          ></q-input>

          <q-input
            v-model="formData.commissionRatePercent"
            :rules="rules.commissionRatePercent"
            label="수수료율"
            type="number"
            min="0"
            max="100"
            required
          >
            <template v-slot:after>
              <q-icon name="percent" />
            </template>
          </q-input>

          <q-checkbox
            v-model="formData.requireUnpaidAmountCheck"
            label="미수금 금액 알림"
          ></q-checkbox>
        </q-card-section>

        <q-card-actions align="right">
          <q-btn
            v-close-popup
            :disable="status.isProgress"
            color="primary"
            label="취소"
            flat
          />
          <q-btn
            :loading="status.isProgress"
            type="submit"
            color="red"
            label="추가"
            flat
          />
        </q-card-actions>
      </q-form>
    </q-card>
  </q-dialog>
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

const emit = defineEmits(["complete"])
const dialog = ref({
  isOpen: false,
})
const status = ref({
  isProgress: false,
})
const formData = ref({
  name: "",
  commissionRatePercent: 0,
  requireUnpaidAmountCheck: false,
})
const rules = {
  name: [value => (value.length >= 2 && value.length <= 20) || "2~20 글자가 필요합니다"],
  commissionRatePercent: [value => (value >= 0 && value <= 100) || "수수료율이 유효하지 않습니다."],
}

function create() {
  status.value.isProgress = true

  const data = {
    name: formData.value.name,
    commissionRate: formData.value.commissionRatePercent / 100,
    requireUnpaidAmountCheck: formData.value.requireUnpaidAmountCheck,
  }

  api.post("/api/v1/reservation-methods", data)
    .then(() => {
      emit("complete")
      dialog.value.isOpen = false

      resetForm()
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

function resetForm() {
  formData.value.name = ""
  formData.value.commissionRate = 0
  formData.value.commissionRatePercent = 0
  formData.value.requireUnpaidAmountCheck = false
}
</script>
