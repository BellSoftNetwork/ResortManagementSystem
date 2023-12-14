<template>
  <q-stepper v-model="formStatus.step" ref="stepper" color="primary" animated vertical>
    <q-step
      :name="1"
      :caption="formatStayCaption(formModel.stayDate?.from, formModel.stayDate?.to)"
      :error="stayDateDiff <= 0"
      :done="formStatus.step > 1"
      title="숙박 기간"
      icon="settings"
    >
      <q-date
        v-model="formModel.stayDate"
        :color="stayDateDiff > 0 ? 'primary' : 'red'"
        :title="formatStayTitle(stayDateDiff)"
        subtitle="숙박 기간"
        range
        mask="YYYY-MM-DD"
      />

      <q-stepper-navigation>
        <q-btn @click="$refs.stepper.next()" label="다음" color="primary" />
      </q-stepper-navigation>
    </q-step>

    <q-step
      :name="2"
      :done="formStatus.step > 2"
      title="객실 배정"
      :caption="selectedRooms.length !== 0 ? selectedRooms.map((room) => room.number).join(', ') : '추후 배정'"
      icon="create_new_folder"
    >
      <RoomGroupSelector
        v-model:selected="selectedRooms"
        :stay-start-at="formModel.stayDate.from"
        :stay-end-at="formModel.stayDate.to"
      />

      <q-stepper-navigation>
        <q-btn @click="$refs.stepper.next()" label="다음" color="primary" />
        <q-btn @click="$refs.stepper.previous()" color="primary" label="이전" class="q-ml-sm" flat />
      </q-stepper-navigation>
    </q-step>

    <q-step
      :name="3"
      :done="formStatus.step > 3"
      title="금액 확인"
      :caption="formatPrice(formModel.price)"
      icon="create_new_folder"
    >
      <q-form @submit="formStatus.step = 4">
        <q-input
          v-model.number="formModel.price"
          :rules="reservationStaticRules.price"
          @update:model-value="changePrice()"
          label="판매 금액"
          placeholder="100000"
          type="number"
          min="0"
          max="10000000"
          required
        ></q-input>

        <q-input
          v-model.number="formModel.paymentAmount"
          :rules="reservationDynamicRules.paymentAmount(formModel.price)"
          label="누적 결제 금액"
          placeholder="80000"
          type="number"
          min="0"
          max="10000000"
          required
        ></q-input>

        <q-select
          v-model="formModel.paymentMethod"
          @update:model-value="changePrice()"
          :loading="paymentMethodStatus.isLoading"
          :disable="!paymentMethodStatus.isLoaded"
          :options="paymentMethods"
          option-label="name"
          label="결제 수단"
          required
          map-options
        ></q-select>

        <q-input
          v-model.number="formModel.brokerFee"
          :rules="reservationStaticRules.brokerFee"
          :readonly="true"
          label="결제 수단 수수료"
          placeholder="5000"
          type="number"
          min="0"
          max="10000000"
          required
        ></q-input>

        <q-stepper-navigation>
          <q-btn type="submit" label="다음" color="primary" />
          <q-btn @click="$refs.stepper.previous()" color="primary" label="이전" class="q-ml-sm" flat />
        </q-stepper-navigation>
      </q-form>
    </q-step>

    <q-step
      :name="4"
      :done="formStatus.step > 4"
      title="예약자 정보"
      :caption="`${formModel.name} (${formModel.phone}) / ${formModel.peopleCount}명`"
      icon="create_new_folder"
    >
      <q-form @submit="formStatus.step = 5">
        <q-input
          v-model="formModel.name"
          :rules="reservationStaticRules.name"
          label="예약자명"
          placeholder="홍길동"
          required
        ></q-input>

        <q-input
          v-model="formModel.phone"
          :rules="reservationStaticRules.phone"
          label="예약자 연락처"
          placeholder="010-0000-0000"
        ></q-input>

        <q-input
          v-model.number="formModel.peopleCount"
          :rules="reservationStaticRules.peopleCount"
          label="예약인원"
          placeholder="4"
          type="number"
          min="0"
          max="1000"
          required
        ></q-input>

        <q-select
          v-model="formModel.status"
          :options="options.status"
          label="예약 상태"
          required
          emit-value
          map-options
        ></q-select>

        <q-input
          v-model="formModel.note"
          :rules="reservationStaticRules.note"
          type="textarea"
          label="메모"
          placeholder="밤 늦게 입실 예정"
        ></q-input>

        <q-stepper-navigation>
          <q-btn label="다음" type="submit" color="primary" />
          <q-btn @click="$refs.stepper.previous()" color="primary" label="이전" class="q-ml-sm" flat />
        </q-stepper-navigation>
      </q-form>
    </q-step>

    <q-step :name="5" title="등록" icon="add_comment">
      <q-stepper-navigation>
        <q-btn @click="create" label="등록" color="primary" />
        <q-btn @click="$refs.stepper.previous()" color="primary" label="이전" class="q-ml-sm" flat />
      </q-stepper-navigation>
    </q-step>
  </q-stepper>
</template>

<script setup lang="ts">
import { computed, onBeforeMount, ref } from "vue";
import { useRouter } from "vue-router";
import { useQuasar } from "quasar";
import dayjs from "dayjs";
import { formatDate, formatDiffDays, formatPrice, formatStayCaption, formatStayTitle } from "src/util/format-util";
import { reservationDynamicRules, reservationStaticRules } from "src/schema/reservation";
import { createReservation } from "src/api/v1/reservation";
import { fetchPaymentMethods } from "src/api/v1/payment-method";
import { formatSortParam } from "src/util/query-string-util";
import { PaymentMethod } from "src/schema/payment-method";
import { Room } from "src/schema/room";
import RoomGroupSelector from "components/room-group/RoomGroupSelector.vue";

const router = useRouter();
const $q = useQuasar();

const formStatus = ref({
  step: 1,
  isProgress: false,
});
const formModel = ref({
  paymentMethod: {
    id: -1,
    name: "네이버",
    commissionRate: 0.1,
    createdAt: "2021-01-01T00:00:00.000Z",
    updatedAt: "2021-01-01T00:00:00.000Z",
  },
  name: "",
  phone: "",
  peopleCount: 4,
  stayDate: { from: "", to: "" },
  price: 0,
  paymentAmount: 0,
  brokerFee: 0,
  note: "",
  status: "NORMAL",
});
const selectedRooms = ref<Room[]>([]);
const status = ref({
  isProgress: false,
});
const options = {
  status: [
    { label: "예약 대기", value: "PENDING" },
    { label: "예약 확정", value: "NORMAL" },
    { label: "예약 취소", value: "CANCEL" },
    { label: "환불 완료", value: "REFUND" },
  ],
};
const stayDateDiff = computed(() => formatDiffDays(formModel.value.stayDate.from, formModel.value.stayDate.to));
const paymentMethodStatus = ref({
  isLoading: false,
  isLoaded: false,
});
const paymentMethods = ref<PaymentMethod[]>();

function loadPaymentMethods() {
  paymentMethodStatus.value.isLoading = true;
  paymentMethodStatus.value.isLoaded = false;
  paymentMethods.value = [];

  fetchPaymentMethods({
    sort: formatSortParam({ field: "name" }),
  })
    .then((response) => {
      paymentMethods.value = response.values;
      formModel.value.paymentMethod = response.values[0];

      paymentMethodStatus.value.isLoaded = true;
    })
    .finally(() => {
      paymentMethodStatus.value.isLoading = false;
    });
}

function create() {
  status.value.isProgress = true;

  createReservation(formData())
    .then(() => {
      router.push({ name: "Reservations" });

      resetForm();
    })
    .catch((error) => {
      $q.notify({
        message: error.response.data.message,
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

function formData() {
  return {
    ...formModel.value,
    rooms: selectedRooms.value,
    stayStartAt: formModel.value.stayDate?.from,
    stayEndAt: formModel.value.stayDate?.to,
  };
}

function changePrice() {
  formModel.value.brokerFee = formModel.value.price * formModel.value.paymentMethod.commissionRate;
}

function resetForm() {
  formModel.value.name = "";
  formModel.value.phone = "";
  formModel.value.peopleCount = 4;
  formModel.value.stayDate.from = formatDate();
  formModel.value.stayDate.to = dayjs().add(1, "d").format("YYYY-MM-DD");
  formModel.value.price = 0;
  formModel.value.paymentAmount = 0;
  formModel.value.brokerFee = 0;
  formModel.value.note = "";
  formModel.value.status = "NORMAL";
}

onBeforeMount(() => {
  resetForm();
  loadPaymentMethods();
});
</script>
