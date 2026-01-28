<template>
  <div class="row q-gutter-md q-pa-sm">
    <div class="col">
      <q-input
        :model-value="formModel.stayStartAt"
        :readonly="true"
        :rules="[() => true]"
        style="min-width: 60px"
        label="입실일"
        dense
      >
        <q-popup-proxy v-if="mode !== 'view'" transition-show="scale" transition-hide="scale">
          <q-date
            @update:model-value="handleStayStartAtChange"
            :model-value="formModel.stayStartAt"
            mask="YYYY-MM-DD"
            no-unset
          >
            <div class="row items-center justify-end">
              <q-btn v-close-popup label="Close" color="primary" flat />
            </div>
          </q-date>
        </q-popup-proxy>
      </q-input>

      <q-input
        :model-value="formModel.name"
        @update:model-value="updateField('name', $event)"
        :rules="rules.name"
        :readonly="mode === 'view'"
        label="예약자명"
        placeholder="홍길동"
        required
      ></q-input>

      <q-input
        :model-value="formModel.peopleCount"
        @update:model-value="updateField('peopleCount', Number($event))"
        :rules="rules.peopleCount"
        :readonly="mode === 'view'"
        label="예약인원"
        placeholder="4"
        type="number"
        min="0"
        max="1000"
        required
      ></q-input>
    </div>

    <div class="col">
      <q-input
        :model-value="formModel.stayEndAt"
        :readonly="true"
        :rules="[() => true]"
        label="퇴실일"
        style="min-width: 60px"
        dense
      >
        <q-popup-proxy v-if="mode !== 'view'" transition-show="scale" transition-hide="scale">
          <q-date
            @update:model-value="handleStayEndAtChange"
            :model-value="formModel.stayEndAt"
            mask="YYYY-MM-DD"
            no-unset
          >
            <div class="row items-center justify-end">
              <q-btn v-close-popup label="Close" color="primary" flat />
            </div>
          </q-date>
        </q-popup-proxy>
      </q-input>

      <!-- view 모드: tel: 링크 -->
      <div v-if="mode === 'view'" class="q-field q-field--readonly">
        <div class="q-field__label">예약자 연락처</div>
        <div class="q-field__native">
          <a v-if="formModel.phone" :href="`tel:${formModel.phone}`" class="text-primary">
            <q-icon name="phone" size="xs" class="q-mr-xs" />
            {{ formModel.phone }}
          </a>
          <span v-else>-</span>
        </div>
      </div>

      <!-- edit/create 모드: 기존 q-input 유지 -->
      <q-input
        v-else
        :model-value="formModel.phone"
        @update:model-value="updateField('phone', $event)"
        :rules="rules.phone"
        label="예약자 연락처"
        placeholder="010-0000-0000"
      ></q-input>

      <q-select
        :model-value="formModel.status"
        @update:model-value="updateField('status', $event)"
        :options="statusOptions"
        :readonly="mode === 'view'"
        label="예약 상태"
        required
        emit-value
        map-options
      ></q-select>
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
  rules: {
    name: Array<(val: string) => boolean | string>;
    phone: Array<(val: string) => boolean | string>;
    peopleCount: Array<(val: number) => boolean | string>;
  };
}>();

const emit = defineEmits<{
  (e: "update:formModel", value: typeof props.formModel): void;
  (e: "recalculateStayStartAt", stayEndAt: string): void;
  (e: "recalculateStayEndAt", stayStartAt: string): void;
}>();

const statusOptions = [
  { label: "예약 대기", value: "PENDING" },
  { label: "예약 확정", value: "NORMAL" },
  { label: "예약 취소", value: "CANCEL" },
  { label: "환불 완료", value: "REFUND" },
];

function updateField(field: string, value: unknown) {
  emit("update:formModel", {
    ...props.formModel,
    [field]: value,
  });
}

function handleStayStartAtChange(stayStartAt: string) {
  updateField("stayStartAt", stayStartAt);
  emit("recalculateStayEndAt", stayStartAt);
}

function handleStayEndAtChange(stayEndAt: string) {
  updateField("stayEndAt", stayEndAt);
  emit("recalculateStayStartAt", stayEndAt);
}
</script>
