<template>
  <slot :dialog="dialog">
    <q-btn @click="dialog.isOpen = true">추가</q-btn>
  </slot>

  <q-dialog v-model="dialog.isOpen" :persistent="status.isProgress" @beforeShow="resetForm">
    <q-card style="width: 500px">
      <q-card-section class="text-h6">결제 수단 추가</q-card-section>

      <q-form @submit="create">
        <q-card-section>
          <q-input v-model="formData.name" :rules="paymentMethodStaticRules.name" label="이름" required></q-input>

          <q-input
            v-model.number="formData.commissionRatePercent"
            :rules="paymentMethodStaticRules.commissionRatePercent"
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

          <q-checkbox v-model="formData.requireUnpaidAmountCheck" label="미수금 금액 알림"></q-checkbox>
        </q-card-section>

        <q-card-actions align="right">
          <q-btn v-close-popup :disable="status.isProgress" color="primary" label="취소" flat />
          <q-btn :loading="status.isProgress" type="submit" color="red" label="추가" flat />
        </q-card-actions>
      </q-form>
    </q-card>
  </q-dialog>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { useQuasar } from "quasar";
import { paymentMethodStaticRules } from "src/schema/payment-method";
import { createPaymentMethod } from "src/api/v1/payment-method";
import { getErrorMessage } from "src/util/errorHandler";

const $q = useQuasar();

const emit = defineEmits(["complete"]);
const dialog = ref({
  isOpen: false,
});
const status = ref({
  isProgress: false,
});
const formData = ref({
  name: "",
  commissionRatePercent: 0,
  requireUnpaidAmountCheck: false,
});

function create() {
  status.value.isProgress = true;

  createPaymentMethod({
    name: formData.value.name,
    commissionRate: formData.value.commissionRatePercent / 100,
    requireUnpaidAmountCheck: formData.value.requireUnpaidAmountCheck,
  })
    .then(() => {
      emit("complete");
      dialog.value.isOpen = false;

      resetForm();
    })
    .catch((error) => {
      $q.notify({
        message: getErrorMessage(error),
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

function resetForm() {
  formData.value.name = "";
  formData.value.commissionRatePercent = 0;
  formData.value.requireUnpaidAmountCheck = false;
}
</script>
