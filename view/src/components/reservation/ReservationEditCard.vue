<template>
  <q-stepper
    v-model="formStatus.step"
    ref="stepper"
    color="primary"
    animated
    vertical
  >
    <q-step
      :name="1"
      :caption="
        formatStayCaption(formModel.stayDate?.from, formModel.stayDate?.to)
      "
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
      :caption="
        selectedRoom[0] && Object.keys(selectedRoom[0]).includes('number')
          ? selectedRoom[0].number
          : '추후 배정'
      "
      icon="create_new_folder"
    >
      <RoomSelectTable
        v-model:selected="selectedRoom"
        :first-value="entity.room"
        :stay-start-at="formModel.stayDate.from"
        :stay-end-at="formModel.stayDate.to"
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
          v-model="formModel.reservationMethod"
          @update:model-value="changePrice()"
          :loading="reservationMethods.status.isLoading"
          :disable="!reservationMethods.status.isLoaded"
          :options="reservationMethods.values"
          option-label="name"
          label="예약 수단"
          required
          map-options
        ></q-select>

        <q-input
          v-model.number="formModel.brokerFee"
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

    <q-step :name="5" title="수정" icon="add_comment">
      <div>
        {{ Object.keys(patchedData()).length }}개 항목이 변경되었습니다.
      </div>

      <q-stepper-navigation>
        <q-btn @click="update" label="수정" color="primary" />
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
import { computed, onBeforeMount, ref } from "vue";
import { useRouter } from "vue-router";
import { useQuasar } from "quasar";
import RoomSelectTable from "components/room/RoomSelectTable.vue";
import {
  formatDate,
  formatDiffDays,
  formatPrice,
  formatStayCaption,
  formatStayTitle,
} from "src/util/format-util";
import {
  Reservation,
  reservationDynamicRules,
  reservationStaticRules,
} from "src/schema/reservation";
import {
  fetchReservation,
  patchReservation,
  ReservationPatchParams,
} from "src/api/v1/reservation";
import { fetchReservationMethods } from "src/api/v1/reservation-method";
import { formatSortParam } from "src/util/query-string-util";
import { Room } from "src/schema/room";
import { getPatchedFormData } from "src/util/data-util";

const router = useRouter();
const $q = useQuasar();

const props = defineProps<{
  id: number;
}>();
const id = props.id;
const entity = ref<
  {
    reservationMethodId: number;
    roomId: number | null;
  } & Reservation
>();
const formStatus = ref({
  step: 1,
  isProgress: false,
});
const formModel = ref({
  id: props.id,
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
});
const selectedRoom = ref<Room[]>([]);
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
  formatDiffDays(formModel.value.stayDate?.from, formModel.value.stayDate?.to),
);
const reservationMethods = ref({
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

function setReservation(reservation: Reservation) {
  entity.value = {
    ...reservation,
    reservationMethodId: reservation.reservationMethod.id,
    roomId: reservation.room?.id,
  };
}

function fetchData() {
  status.value.isProgress = true;

  return fetchReservation(id)
    .then((response) => {
      setReservation(response.value);
    })
    .catch((error) => {
      if (error.response.status === 404) router.push({ name: "ErrorNotFound" });

      console.log(error);
    })
    .finally(() => {
      status.value.isProgress = false;
    });
}

function loadReservationMethods() {
  reservationMethods.value.status.isLoading = true;
  reservationMethods.value.status.isLoaded = false;
  reservationMethods.value.values = [];

  return fetchReservationMethods({
    sort: formatSortParam({ field: "name" }),
  })
    .then((response) => {
      reservationMethods.value.values = response.values;

      reservationMethods.value.status.isLoaded = true;
    })
    .finally(() => {
      reservationMethods.value.status.isLoading = false;
    });
}

function update() {
  status.value.isProgress = true;

  patchReservation(id, patchedData())
    .then(() => {
      router.push({ name: "Reservation", params: { id: id } });

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

function getFormData(): ReservationPatchParams {
  const formData: Partial<typeof formModel.value> & ReservationPatchParams = {
    ...formModel.value,
    reservationMethodId: formModel.value.reservationMethod.id,
    roomId:
      selectedRoom.value[0] && Object.keys(selectedRoom.value[0]).includes("id")
        ? selectedRoom.value[0].id
        : null,
    stayStartAt: formModel.value.stayDate?.from,
    stayEndAt: formModel.value.stayDate?.to,
  };

  delete formData.id;
  delete formData.reservationMethod;
  delete formData.stayDate;

  return formData;
}

function patchedData(): ReservationPatchParams {
  return getPatchedFormData(entity.value, getFormData());
}

function changePrice() {
  formModel.value.brokerFee =
    formModel.value.price * formModel.value.reservationMethod.commissionRate;
}

function resetForm() {
  Object.assign(formModel.value, entity.value);
  formModel.value.stayDate.from = formatDate(entity.value?.stayStartAt);
  formModel.value.stayDate.to = formatDate(entity.value?.stayEndAt);
}

onBeforeMount(() => {
  fetchData().then(() => {
    if (entity.value === undefined) return;

    resetForm();
    if (entity.value.room) selectedRoom.value = [entity.value.room];
    loadReservationMethods().then(() => {
      formModel.value.reservationMethod = reservationMethods.value.values.find(
        (item) => item.id === entity.value.reservationMethod.id,
      );
    });
  });
});
</script>
