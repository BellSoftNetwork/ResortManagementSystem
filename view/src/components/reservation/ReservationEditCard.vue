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
        <q-btn @click="$refs.stepper.next() || checkRooms()" label="다음" color="primary" />
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
        :parent-reservation="props.reservation"
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
          :loading="paymentMethods.status.isLoading"
          :disable="!paymentMethods.status.isLoaded"
          :options="paymentMethods.values"
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

    <q-step :name="5" title="수정" icon="add_comment">
      <div>{{ Object.keys(patchedData()).length }}개 항목이 변경되었습니다.</div>

      <q-stepper-navigation>
        <q-btn @click="update" label="수정" color="primary" />
        <q-btn @click="$refs.stepper.previous()" color="primary" label="이전" class="q-ml-sm" flat />
      </q-stepper-navigation>
    </q-step>
  </q-stepper>
</template>

<script setup lang="ts">
import { computed, onBeforeMount, ref } from "vue";
import { useRouter } from "vue-router";
import { useQuasar } from "quasar";
import { formatDate, formatDiffDays, formatPrice, formatStayCaption, formatStayTitle } from "src/util/format-util";
import { Reservation, reservationDynamicRules, reservationStaticRules } from "src/schema/reservation";
import { patchReservation, ReservationPatchParams } from "src/api/v1/reservation";
import { fetchPaymentMethods } from "src/api/v1/payment-method";
import { formatSortParam } from "src/util/query-string-util";
import { Room } from "src/schema/room";
import { getPatchedFormData } from "src/util/data-util";
import { fetchRooms } from "src/api/v1/room";
import RoomGroupSelector from "components/room-group/RoomGroupSelector.vue";

const router = useRouter();
const $q = useQuasar();

const props = defineProps<{
  reservation: Reservation;
}>();
const formStatus = ref({
  step: 1,
  isProgress: false,
});
const formModel = ref({
  id: props.reservation.id,
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
const selectedRooms = ref<Room[]>(props.reservation.rooms);
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
const stayDateDiff = computed(() => formatDiffDays(formModel.value.stayDate?.from, formModel.value.stayDate?.to));
const paymentMethods = ref({
  status: {
    isLoading: false,
    isLoaded: false,
  },
  values: [
    {
      id: -1,
      name: "네이버",
      commissionRate: 0.1,
      createdAt: "2021-01-01T00:00:00.000Z",
      updatedAt: "2021-01-01T00:00:00.000Z",
    },
  ],
});

function loadPaymentMethods() {
  paymentMethods.value.status.isLoading = true;
  paymentMethods.value.status.isLoaded = false;
  paymentMethods.value.values = [];

  return fetchPaymentMethods({
    sort: formatSortParam({ field: "name" }),
  })
    .then((response) => {
      paymentMethods.value.values = response.values;

      paymentMethods.value.status.isLoaded = true;
    })
    .finally(() => {
      paymentMethods.value.status.isLoading = false;
    });
}

function update() {
  const patchParams = patchedData();

  if (Object.keys(patchParams).length === 0) {
    $q.notify({
      message: "수정된 항목이 없습니다.",
      type: "info",
      actions: [
        {
          icon: "close",
          color: "white",
          round: true,
        },
      ],
    });

    return;
  }

  status.value.isProgress = true;

  patchReservation(props.reservation.id, patchParams)
    .then(() => {
      router.push({ name: "Reservation", params: { id: props.reservation.id } });

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

function checkRooms() {
  fetchRooms({
    stayStartAt: formModel.value.stayDate.from,
    stayEndAt: formModel.value.stayDate.to,
    status: "NORMAL",
    excludeReservationId: props.reservation.id,
  }).then((response) => {
    const availableRoomIds = response.values.map((room) => room.id);
    const unavailableRooms = props.reservation.rooms.filter((room) => !availableRoomIds.includes(room.id));

    selectedRooms.value = props.reservation.rooms.filter((room) => availableRoomIds.includes(room.id));

    if (unavailableRooms.length > 0) {
      $q.notify({
        message:
          `기존에 배정된 ${unavailableRooms.length}개의 객실이 해당 기간에 이용할 수 없어 제외되었습니다.<br />` +
          `이용 불가 객실: ${unavailableRooms.map((room) => room.number).join(", ")}`,
        type: "warning",
        html: true,
        actions: [
          {
            icon: "close",
            color: "white",
            round: true,
          },
        ],
      });
    }
  });
}

function getFormData(): ReservationPatchParams {
  const formData: Partial<typeof formModel.value> & ReservationPatchParams = {
    ...formModel.value,
    rooms: selectedRooms.value,
    stayStartAt: formModel.value.stayDate?.from,
    stayEndAt: formModel.value.stayDate?.to,
  };

  delete formData.id;
  delete formData.stayDate;

  return formData;
}

function patchedData(): ReservationPatchParams {
  return getPatchedFormData(props.reservation, getFormData());
}

function changePrice() {
  formModel.value.brokerFee = formModel.value.price * formModel.value.paymentMethod.commissionRate;
}

function resetForm() {
  Object.assign(formModel.value, props.reservation);
  formModel.value.stayDate.from = formatDate(props.reservation.stayStartAt);
  formModel.value.stayDate.to = formatDate(props.reservation.stayEndAt);
}

onBeforeMount(() => {
  resetForm();
  loadPaymentMethods().then(() => {
    formModel.value.paymentMethod = paymentMethods.value.values.find(
      (item) => item.id === props.reservation.paymentMethod.id,
    );
  });
});
</script>
