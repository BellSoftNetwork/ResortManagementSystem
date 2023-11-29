<template>
  <q-stepper
    v-model="formModel.status.step"
    ref="stepper"
    color="primary"
    animated
    vertical
  >
    <q-step
      :name="1"
      :caption="
        formatStayCaption(
          formModel.value.stayDate.from,
          formModel.value.stayDate.to,
        )
      "
      :error="stayDateDiff <= 0"
      :done="formModel.status.step > 1"
      title="숙박 기간"
      icon="settings"
    >
      <q-date
        v-model="formModel.value.stayDate"
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
      :done="formModel.status.step > 2"
      title="객실 배정"
      :caption="
        selectedRoom[0] && Object.keys(selectedRoom[0]).includes('number')
          ? selectedRoom[0].number
          : '추후 배정'
      "
      icon="create_new_folder"
    >
      <RoomSelectTable
        v-model:selected="selectedRoom"
        :stay-start-at="formModel.value.stayDate.from"
        :stay-end-at="formModel.value.stayDate.to"
      />

      <q-stepper-navigation>
        <q-btn @click="$refs.stepper.next()" label="다음" color="primary" />
        <q-btn
          @click="$refs.stepper.previous()"
          color="primary"
          label="이전"
          class="q-ml-sm"
          flat
        />
      </q-stepper-navigation>
    </q-step>

    <q-step
      :name="3"
      :done="formModel.status.step > 3"
      title="금액 확인"
      :caption="formatPrice(formModel.value.price)"
      icon="create_new_folder"
    >
      <q-form @submit="formModel.status.step = 4">
        <q-input
          v-model.number="formModel.value.price"
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
          v-model.number="formModel.value.paymentAmount"
          :rules="reservationDynamicRules.paymentAmount(formModel.value.price)"
          label="누적 결제 금액"
          placeholder="80000"
          type="number"
          min="0"
          max="10000000"
          required
        ></q-input>

        <q-select
          v-model="formModel.value.reservationMethod"
          @update:model-value="changePrice()"
          :loading="reservationMethodStatus.isLoading"
          :disable="!reservationMethodStatus.isLoaded"
          :options="reservationMethods"
          option-label="name"
          label="예약 수단"
          required
          map-options
        ></q-select>

        <q-input
          v-model.number="formModel.value.brokerFee"
          :rules="reservationStaticRules.brokerFee"
          :readonly="true"
          label="예약 수단 수수료"
          placeholder="5000"
          type="number"
          min="0"
          max="10000000"
          required
        ></q-input>

        <q-stepper-navigation>
          <q-btn type="submit" label="다음" color="primary" />
          <q-btn
            @click="$refs.stepper.previous()"
            color="primary"
            label="이전"
            class="q-ml-sm"
            flat
          />
        </q-stepper-navigation>
      </q-form>
    </q-step>

    <q-step
      :name="4"
      :done="formModel.status.step > 4"
      title="예약자 정보"
      :caption="`${formModel.value.name} (${formModel.value.phone}) / ${formModel.value.peopleCount}명`"
      icon="create_new_folder"
    >
      <q-form @submit="formModel.status.step = 5">
        <q-input
          v-model="formModel.value.name"
          :rules="reservationStaticRules.name"
          label="예약자명"
          placeholder="홍길동"
          required
        ></q-input>

        <q-input
          v-model="formModel.value.phone"
          :rules="reservationStaticRules.phone"
          label="예약자 연락처"
          placeholder="010-0000-0000"
        ></q-input>

        <q-input
          v-model.number="formModel.value.peopleCount"
          :rules="reservationStaticRules.peopleCount"
          label="예약인원"
          placeholder="4"
          type="number"
          min="0"
          max="1000"
          required
        ></q-input>

        <q-select
          v-model="formModel.value.status"
          :options="options.status"
          label="예약 상태"
          required
          emit-value
          map-options
        ></q-select>

        <q-input
          v-model="formModel.value.note"
          :rules="reservationStaticRules.note"
          type="textarea"
          label="메모"
          placeholder="밤 늦게 입실 예정"
        ></q-input>

        <q-stepper-navigation>
          <q-btn label="다음" type="submit" color="primary" />
          <q-btn
            @click="$refs.stepper.previous()"
            color="primary"
            label="이전"
            class="q-ml-sm"
            flat
          />
        </q-stepper-navigation>
      </q-form>
    </q-step>

    <q-step :name="5" title="등록" icon="add_comment">
      <q-stepper-navigation>
        <q-btn @click="create" label="등록" color="primary" />
        <q-btn
          @click="$refs.stepper.previous()"
          color="primary"
          label="이전"
          class="q-ml-sm"
          flat
        />
      </q-stepper-navigation>
    </q-step>
  </q-stepper>
</template>

<script setup lang="ts">
import { computed, onBeforeMount, ref } from "vue"
import { useRouter } from "vue-router"
import { useQuasar } from "quasar"
import dayjs from "dayjs"
import RoomSelectTable from "components/room/RoomSelectTable.vue"
import {
  formatDate,
  formatDiffDays,
  formatPrice,
  formatStayCaption,
  formatStayTitle,
} from "src/util/format-util"
import {
  reservationDynamicRules,
  reservationStaticRules,
} from "src/schema/reservation"
import { createReservation } from "src/api/v1/reservation"
import { fetchReservationMethods } from "src/api/v1/reservation-method"
import { formatSortParam } from "src/util/query-string-util"
import { ReservationMethod } from "src/schema/reservation-method"
import { Room } from "src/schema/room"

const router = useRouter()
const $q = useQuasar()

const formModel = ref({
  status: {
    step: 1,
    isProgress: false,
  },
  value: {
    reservationMethod: {
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
    status: "PENDING",
  },
});
const selectedRoom = ref<Room[]>([])
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
const stayDateDiff = computed(() =>
  formatDiffDays(
    formModel.value.value.stayDate.from,
    formModel.value.value.stayDate.to,
  ),
);
const reservationMethodStatus = ref({
  isLoading: false,
  isLoaded: false,
});
const reservationMethods = ref<ReservationMethod[]>()

function loadReservationMethods() {
  reservationMethodStatus.value.isLoading = true
  reservationMethodStatus.value.isLoaded = false
  reservationMethods.value = []

  fetchReservationMethods({
    sort: formatSortParam({ field: "name" }),
  })
    .then((response) => {
      reservationMethods.value = response.values
      formModel.value.value.reservationMethod = response.values[0]

      reservationMethodStatus.value.isLoaded = true
    })
    .finally(() => {
      reservationMethodStatus.value.isLoading = false
    });
}

function create() {
  status.value.isProgress = true

  createReservation(formData())
    .then(() => {
      router.push({ name: "Reservations" })

      resetForm()
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
      status.value.isProgress = false
    });
}

function formData() {
  return {
    reservationMethodId: formModel.value.value.reservationMethod.id,
    roomId:
      selectedRoom.value[0] && Object.keys(selectedRoom.value[0]).includes("id")
        ? selectedRoom.value[0].id
        : null,
    name: formModel.value.value.name,
    phone: formModel.value.value.phone,
    peopleCount: formModel.value.value.peopleCount,
    stayStartAt: formModel.value.value.stayDate.from,
    stayEndAt: formModel.value.value.stayDate.to,
    price: formModel.value.value.price,
    paymentAmount: formModel.value.value.paymentAmount,
    brokerFee: formModel.value.value.brokerFee,
    note: formModel.value.value.note,
    status: formModel.value.value.status,
  };
}

function changePrice() {
  formModel.value.value.brokerFee =
    formModel.value.value.price *
    formModel.value.value.reservationMethod.commissionRate;
}

function resetForm() {
  formModel.value.value.name = ""
  formModel.value.value.phone = ""
  formModel.value.value.peopleCount = 4
  formModel.value.value.stayDate.from = formatDate()
  formModel.value.value.stayDate.to = dayjs().add(1, "d").format("YYYY-MM-DD")
  formModel.value.value.price = 0
  formModel.value.value.paymentAmount = 0
  formModel.value.value.brokerFee = 0
  formModel.value.value.note = ""
  formModel.value.value.status = "PENDING"
}

onBeforeMount(() => {
  resetForm()
  loadReservationMethods()
});
</script>
