<template>
  <div class="row q-gutter-md q-pa-sm">
    <div class="col">
      <q-input
        :model-value="formModel.price"
        @update:model-value="handlePriceChange"
        :rules="rules.price"
        :readonly="mode === 'view'"
        label="판매 금액"
        placeholder="100000"
        type="number"
        min="0"
        max="10000000"
        required
      ></q-input>

      <q-input
        :model-value="formModel.paymentAmount"
        @update:model-value="updateField('paymentAmount', Number($event))"
        :rules="rules.paymentAmount"
        :readonly="mode === 'view'"
        label="누적 결제 금액"
        placeholder="80000"
        type="number"
        min="0"
        max="10000000"
        required
      ></q-input>

      <q-input
        :model-value="formModel.deposit"
        @update:model-value="updateField('deposit', Number($event))"
        :rules="rules.deposit"
        :readonly="mode === 'view'"
        label="보증금"
        placeholder="100000"
        type="number"
        min="0"
        max="10000000"
        required
      ></q-input>
    </div>

    <div class="col">
      <q-select
        :model-value="formModel.paymentMethod"
        @update:model-value="handlePaymentMethodChange"
        :loading="paymentMethodStatus.isLoading"
        :disable="!paymentMethodStatus.isLoaded"
        :options="paymentMethods"
        :readonly="mode === 'view'"
        option-label="name"
        label="결제 수단"
        required
        map-options
      ></q-select>

      <q-input
        :model-value="formModel.brokerFee"
        @update:model-value="updateField('brokerFee', Number($event))"
        :rules="rules.brokerFee"
        :readonly="true"
        label="결제 수단 수수료"
        placeholder="5000"
        type="number"
        min="0"
        max="10000000"
        required
      ></q-input>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ReservationCreateParams } from "src/api/v1/reservation";
import { PaymentMethod } from "src/schema/payment-method";

const props = defineProps<{
  formModel: {
    paymentMethod: Partial<PaymentMethod>;
  } & ReservationCreateParams;
  mode: "create" | "update" | "view";
  paymentMethods: PaymentMethod[] | undefined;
  paymentMethodStatus: {
    isLoading: boolean;
    isLoaded: boolean;
  };
  rules: {
    price: Array<(val: number) => boolean | string>;
    paymentAmount: Array<(val: number) => boolean | string>;
    deposit: Array<(val: number) => boolean | string>;
    brokerFee: Array<(val: number) => boolean | string>;
  };
}>();

const emit = defineEmits<{
  (e: "update:formModel", value: typeof props.formModel): void;
  (e: "changePrice"): void;
}>();

function updateField(field: string, value: unknown) {
  emit("update:formModel", {
    ...props.formModel,
    [field]: value,
  });
}

function handlePriceChange(value: number) {
  updateField("price", value);
  emit("changePrice");
}

function handlePaymentMethodChange(value: Partial<PaymentMethod>) {
  updateField("paymentMethod", value);
  emit("changePrice");
}
</script>
