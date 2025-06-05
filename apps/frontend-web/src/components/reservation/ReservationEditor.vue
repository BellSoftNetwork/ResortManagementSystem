<template>
  <q-card>
    <q-card-section class="text-h6">{{ typeName }} {{ modeTitle }}</q-card-section>

    <form @submit.prevent="submit">
      <q-card-section>
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
                  @update:model-value="recalculateStayEndAt"
                  v-model="formModel.stayStartAt"
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
              v-model="formModel.name"
              :rules="reservationStaticRules.name"
              :readonly="mode === 'view'"
              label="예약자명"
              placeholder="홍길동"
              required
            ></q-input>

            <q-input
              v-model.number="formModel.peopleCount"
              :rules="reservationStaticRules.peopleCount"
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
                  @update:model-value="recalculateStayStartAt"
                  v-model="formModel.stayEndAt"
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
              v-model="formModel.phone"
              :rules="reservationStaticRules.phone"
              :readonly="mode === 'view'"
              label="예약자 연락처"
              placeholder="010-0000-0000"
            ></q-input>

            <q-select
              v-model="formModel.status"
              :options="options.status"
              :readonly="mode === 'view'"
              label="예약 상태"
              required
              emit-value
              map-options
            ></q-select>
          </div>
        </div>

        <div class="row q-gutter-md q-pa-sm">
          <div v-if="mode !== 'view'" class="col q-ma-none">
            <RoomGroupSelector
              v-model:selected="selectedRooms"
              :stay-start-at="formModel.stayStartAt"
              :stay-end-at="formModel.stayEndAt"
            />
          </div>
          <div v-else class="col">
            <div class="q-py-sm">
              <div class="text-caption">객실</div>
              <div class="text-body1">
                <div v-if="formModel.rooms.length !== 0">
                  <span v-for="room in formModel.rooms" :key="room.id">
                    <q-btn :to="{ name: 'Room', params: { id: room.id } }" color="primary" flat dense>
                      {{ room.number }}
                    </q-btn>
                  </span>
                </div>
                <div v-else class="text-grey">미배정</div>
              </div>
            </div>
          </div>
        </div>

        <div class="row q-gutter-md q-pa-sm">
          <div class="col">
            <q-input
              v-model.number="formModel.price"
              :rules="reservationStaticRules.price"
              @update:model-value="changePrice()"
              :readonly="mode === 'view'"
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
              :readonly="mode === 'view'"
              label="누적 결제 금액"
              placeholder="80000"
              type="number"
              min="0"
              max="10000000"
              required
            ></q-input>

            <q-input
              v-model.number="formModel.deposit"
              :rules="reservationDynamicRules.paymentAmount(formModel.deposit)"
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
              v-model="formModel.paymentMethod"
              @update:model-value="changePrice()"
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
          </div>
        </div>

        <div class="row q-gutter-md q-pa-sm">
          <div class="col">
            <q-input
              v-model="formModel.note"
              :rules="reservationStaticRules.note"
              :readonly="mode === 'view'"
              type="textarea"
              label="메모"
              placeholder="밤 늦게 입실 예정"
            ></q-input>
          </div>
        </div>
      </q-card-section>

      <q-card-actions v-if="mode === 'view'" align="right">
        <q-btn @click="deleteItem()" color="red" label="삭제" dense flat></q-btn>
        <q-btn @click="mode = 'update'" color="primary" label="수정" flat />
      </q-card-actions>
      <q-card-actions v-else>
        <q-btn :label="modeTitle" color="primary" class="full-width" type="submit" />
      </q-card-actions>
    </form>
  </q-card>
</template>

<script setup lang="ts">
import { computed, onBeforeMount, ref } from "vue";
import { useRouter } from "vue-router";
import { useQuasar } from "quasar";
import dayjs from "dayjs";
import { formatDate } from "src/util/format-util";
import { Reservation, reservationDynamicRules, reservationStaticRules, ReservationType } from "src/schema/reservation";
import {
  createReservation,
  deleteReservation,
  patchReservation,
  ReservationCreateParams,
} from "src/api/v1/reservation";
import { fetchPaymentMethods } from "src/api/v1/payment-method";
import { formatSortParam } from "src/util/query-string-util";
import { PaymentMethod } from "src/schema/payment-method";
import { Room } from "src/schema/room";
import RoomGroupSelector from "components/room-group/RoomGroupSelector.vue";
import { getPatchedFormData } from "src/util/data-util";
import _ from "lodash";

const router = useRouter();
const $q = useQuasar();

const props = withDefaults(
  defineProps<{
    mode: "create" | "update" | "view";
    reservation?: Reservation;
    reservationType: ReservationType;
  }>(),
  {
    mode: "create",
  },
);
const typeName = computed(() => {
  if (props.reservationType === "MONTHLY_RENT") return "달방";
  else return "예약";
});
const mode = ref(props.mode);
const modeTitle = computed(() => {
  switch (mode.value) {
    case "create":
      return "등록";
    case "update":
      return "수정";
    case "view":
      return "정보";
    default:
      return "";
  }
});
let defaultReservationValue = <
  {
    paymentMethod: Partial<PaymentMethod>;
  } & ReservationCreateParams
>{
  paymentMethod: { id: 0 },
  rooms: [],
  name: "",
  phone: "",
  peopleCount: 4,
  stayStartAt: formatDate(),
  stayEndAt: dayjs().add(1, "d").format("YYYY-MM-DD"),
  price: 0,
  deposit: 0,
  paymentAmount: 0,
  brokerFee: 0,
  note: "",
  status: "NORMAL",
  type: props.reservationType,
};
const formModel = ref<
  {
    paymentMethod: Partial<PaymentMethod>;
  } & ReservationCreateParams
>(_.cloneDeep(defaultReservationValue));
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

      if (defaultReservationValue.paymentMethod.id === 0) {
        const defaultSelectPaymentMethod = response.values.find((paymentMethod) => paymentMethod.isDefaultSelect);
        defaultReservationValue.paymentMethod = defaultSelectPaymentMethod
          ? defaultSelectPaymentMethod
          : response.values[0];
      }
      formModel.value.paymentMethod = defaultReservationValue.paymentMethod;

      paymentMethodStatus.value.isLoaded = true;
    })
    .finally(() => {
      paymentMethodStatus.value.isLoading = false;
    });
}

function submit() {
  // Check if rooms are selected
  if (selectedRooms.value.length === 0) {
    $q.dialog({
      title: "경고",
      message: [
        "객실 미배정 상태로 먼저 예약을 등록하시겠습니까?",
        "객실을 배정하면 객실 현황 페이지에서 현재 객실 입실 상태를 한 눈에 확인할 수 있습니다.",
      ].join("<br />"),
      html: true,
      cancel: {
        label: "취소",
        flat: true,
        color: "negative",
      },
      ok: {
        label: "등록",
        flat: true,
        color: "primary",
      },
      persistent: true,
    })
      .onOk(() => {
        // Proceed with registration without room assignment
        switch (mode.value) {
          case "create":
            return create();
          case "update":
            return update();
        }
      })
      .onCancel(() => {
        // Scroll to room assignment section
        const roomGroupSelector = document.querySelector(".q-ma-none");
        if (roomGroupSelector) {
          roomGroupSelector.scrollIntoView({ behavior: "smooth" });
        }
      });
    return;
  }

  switch (mode.value) {
    case "create":
      return create();

    case "update":
      return update();
  }
}

function create() {
  status.value.isProgress = true;

  createReservation(formData())
    .then(() => {
      resetForm();

      if (props.reservationType === "MONTHLY_RENT") return router.push({ name: "MonthlyRents" });
      return router.push({ name: "Reservations" });
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

function update() {
  const patchParams = getPatchedFormData(defaultReservationValue, formData());

  let updatedKeyCount = Object.keys(patchParams).length;
  if (updatedKeyCount === 0) {
    mode.value = "view";
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

  patchReservation(defaultReservationValue.id, patchParams)
    .then((response) => {
      $q.notify({
        message: `${updatedKeyCount} 건의 항목이 정상적으로 수정되었습니다.`,
        type: "info",
        actions: [
          {
            icon: "close",
            color: "white",
            round: true,
          },
        ],
      });

      mode.value = "view";

      defaultReservationValue = _.cloneDeep(response.value);
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

function deleteItem() {
  const itemId = defaultReservationValue.id;
  const itemName = defaultReservationValue.name;

  $q.dialog({
    title: "삭제",
    message: `정말로 ${itemName}님의 ${typeName.value} 정보를 삭제하시겠습니까?`,
    ok: {
      label: "삭제",
      flat: true,
      color: "negative",
    },
    cancel: {
      label: "유지",
      flat: true,
    },
    focus: "cancel",
  }).onOk(() => {
    deleteReservation(itemId)
      .then(() => {
        if (props.reservationType === "MONTHLY_RENT") return router.push({ name: "MonthlyRents" });
        return router.push({ name: "Reservations" });
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
      });
  });
}

function formData() {
  return {
    ...formModel.value,
    rooms: selectedRooms.value,
  };
}

function recalculateStayStartAt(stayEndAt: string) {
  if (dayjs(formModel.value.stayStartAt).isBefore(stayEndAt)) return;

  formModel.value.stayStartAt = dayjs(stayEndAt).add(-1, "d").format("YYYY-MM-DD");
}

function recalculateStayEndAt(stayStartAt: string) {
  if (dayjs(formModel.value.stayEndAt).isAfter(stayStartAt)) return;

  formModel.value.stayEndAt = dayjs(stayStartAt).add(1, "d").format("YYYY-MM-DD");
}

function changePrice() {
  formModel.value.brokerFee = formModel.value.price * formModel.value.paymentMethod.commissionRate;
}

function resetForm() {
  formModel.value = _.cloneDeep(defaultReservationValue);
  selectedRooms.value = _.cloneDeep(defaultReservationValue.rooms);
}

onBeforeMount(() => {
  if (props.mode !== "create") defaultReservationValue = _.cloneDeep(props.reservation);

  resetForm();
  loadPaymentMethods();
});
</script>
