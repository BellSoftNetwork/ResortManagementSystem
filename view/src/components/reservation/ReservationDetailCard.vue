<template>
  <q-card flat bordered>
    <q-inner-loading :showing="status.isProgress">
      <q-spinner-gears size="50px" color="primary" />
    </q-inner-loading>

    <q-card-section class="text-h6"> 예약 정보</q-card-section>

    <q-card-section>
      <div class="row">
        <div class="col-12 col-md-auto q-px-sm">
          <q-date
            :model-value="{ from: entity?.stayStartAt, to: entity?.stayEndAt }"
            :title="formatStayTitle(stayDateDiff)"
            subtitle="숙박 기간"
            range
            mask="YYYY-MM-DD"
            :readonly="true"
          />
        </div>

        <div class="col-12 col-md-4 q-px-sm">
          <div class="q-py-sm">
            <div class="text-caption">객실 번호</div>
            <div class="text-body1">
              <div v-if="entity?.room">
                <q-btn
                  :to="{ name: 'Room', params: { id: entity?.room.id } }"
                  color="primary"
                  flat
                  dense
                >
                  {{ entity?.room.number }}
                </q-btn>
              </div>
              <div v-else class="text-grey">미배정</div>
            </div>
          </div>

          <div class="q-py-sm">
            <div class="text-caption">판매 금액</div>
            <div class="text-body1">{{ formatPrice(entity?.price) }}</div>
          </div>

          <div class="q-py-sm">
            <div class="text-caption">결제 금액</div>
            <div class="text-body1">
              {{ formatPrice(entity?.paymentAmount) }}
            </div>
          </div>

          <div class="q-py-sm">
            <div class="text-caption">예약 수단</div>
            <div class="text-body1">{{ entity?.reservationMethod.name }}</div>
          </div>

          <div class="q-py-sm">
            <div class="text-caption">예약 수단 수수료</div>
            <div class="text-body1">{{ formatPrice(entity?.brokerFee) }}</div>
          </div>
        </div>

        <div class="col-12 col-md-4 q-px-sm">
          <div class="q-py-sm">
            <div class="text-caption">예약자명</div>
            <div class="text-body1">{{ entity?.name }}</div>
          </div>

          <div class="q-py-sm">
            <div class="text-caption">예약자 연락처</div>
            <div class="text-body1">
              <a :href="'tel:' + entity?.phone">{{ entity?.phone }}</a>
            </div>
          </div>

          <div class="q-py-sm">
            <div class="text-caption">예약인원</div>
            <div class="text-body1">{{ entity?.peopleCount }}</div>
          </div>

          <div class="q-py-sm">
            <div class="text-caption">예약 상태</div>
            <div class="text-body1">
              {{ entity ? reservationStatusValueToName(entity.status) : "" }}
            </div>
          </div>
        </div>
      </div>
      <br />
      <div class="row">
        <div class="col">
          <div class="q-py-sm">
            <div class="text-caption">메모</div>
            <div class="text-body1">{{ entity?.note }}</div>
          </div>
        </div>
      </div>
    </q-card-section>

    <q-card-actions align="right">
      <q-btn @click="deleteItem()" color="red" label="삭제" dense flat></q-btn>
      <q-btn
        :disable="status.isProgress"
        :to="{ name: 'EditReservation', params: { id: id } }"
        color="primary"
        label="수정"
        flat
      />
    </q-card-actions>
  </q-card>
</template>

<script setup lang="ts">
import { computed, onBeforeMount, ref } from "vue";
import { useRouter } from "vue-router";
import { useQuasar } from "quasar";
import {
  formatDiffDays,
  formatPrice,
  formatStayTitle,
} from "src/util/format-util";
import { deleteReservation, fetchReservation } from "src/api/v1/reservation";
import {
  Reservation,
  reservationStatusValueToName,
} from "src/schema/reservation";

const router = useRouter();
const $q = useQuasar();
const props = defineProps<{
  id: number;
}>();
const id = props.id;
const status = ref({
  isProgress: false,
});
const entity = ref<Reservation>();
const stayDateDiff = computed(() =>
  formatDiffDays(entity.value?.stayStartAt, entity.value?.stayEndAt),
);

function fetchData() {
  status.value.isProgress = true;

  fetchReservation(id)
    .then((response) => {
      entity.value = response.value;
    })
    .catch((error) => {
      if (error.response.status === 404) router.push({ name: "ErrorNotFound" });

      console.log(error);
    })
    .finally(() => {
      status.value.isProgress = false;
    });
}

function deleteItem() {
  if (entity.value === undefined) return;

  const itemId = entity.value.id;
  const itemName = entity.value.name;

  $q.dialog({
    title: "삭제",
    message: `정말로 ${itemName}님의 ${formatStayTitle(
      stayDateDiff.value,
    )} 예약 정보를 삭제하시겠습니까?`,
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
        router.push({ name: "Reservations" });
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

onBeforeMount(() => {
  fetchData();
});
</script>
