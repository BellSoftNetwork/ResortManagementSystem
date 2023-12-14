<template>
  <q-card flat bordered>
    <q-card-section class="text-h6">예약 정보</q-card-section>

    <q-card-section>
      <div class="row">
        <div class="col-12 col-md-auto q-px-sm">
          <q-date
            :model-value="{ from: props.reservation.stayStartAt, to: props.reservation.stayEndAt }"
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
              <div v-if="props.reservation.rooms.length !== 0">
                <span v-for="room in props.reservation.rooms" :key="room.id">
                  <q-btn :to="{ name: 'Room', params: { id: room.id } }" color="primary" flat dense>
                    {{ room.number }}
                  </q-btn>
                </span>
              </div>
              <div v-else class="text-grey">미배정</div>
            </div>
          </div>

          <div class="q-py-sm">
            <div class="text-caption">판매 금액</div>
            <div class="text-body1">{{ formatPrice(props.reservation.price) }}</div>
          </div>

          <div class="q-py-sm">
            <div class="text-caption">결제 금액</div>
            <div class="text-body1">
              {{ formatPrice(props.reservation.paymentAmount) }}
            </div>
          </div>

          <div class="q-py-sm">
            <div class="text-caption">결제 수단</div>
            <div class="text-body1">{{ props.reservation.paymentMethod.name }}</div>
          </div>

          <div class="q-py-sm">
            <div class="text-caption">결제 수단 수수료</div>
            <div class="text-body1">{{ formatPrice(props.reservation.brokerFee) }}</div>
          </div>
        </div>

        <div class="col-12 col-md-4 q-px-sm">
          <div class="q-py-sm">
            <div class="text-caption">예약자명</div>
            <div class="text-body1">{{ props.reservation.name }}</div>
          </div>

          <div class="q-py-sm">
            <div class="text-caption">예약자 연락처</div>
            <div class="text-body1">
              <a :href="'tel:' + props.reservation.phone">{{ props.reservation.phone }}</a>
            </div>
          </div>

          <div class="q-py-sm">
            <div class="text-caption">예약인원</div>
            <div class="text-body1">{{ props.reservation.peopleCount }}</div>
          </div>

          <div class="q-py-sm">
            <div class="text-caption">예약 상태</div>
            <div class="text-body1">
              {{ reservationStatusValueToName(props.reservation.status) }}
            </div>
          </div>
        </div>
      </div>
      <br />
      <div class="row">
        <div class="col">
          <div class="q-py-sm">
            <div class="text-caption">메모</div>
            <div class="text-body1">{{ props.reservation.note }}</div>
          </div>
        </div>
      </div>
    </q-card-section>

    <q-card-actions align="right">
      <q-btn @click="deleteItem()" color="red" label="삭제" dense flat></q-btn>
      <q-btn
        :to="{ name: 'EditReservation', params: { id: props.reservation.id } }"
        color="primary"
        label="수정"
        flat
      />
    </q-card-actions>
  </q-card>
</template>

<script setup lang="ts">
import { computed } from "vue";
import { useRouter } from "vue-router";
import { useQuasar } from "quasar";
import { formatDiffDays, formatPrice, formatStayTitle } from "src/util/format-util";
import { deleteReservation } from "src/api/v1/reservation";
import { Reservation, reservationStatusValueToName } from "src/schema/reservation";

const router = useRouter();
const $q = useQuasar();
const props = defineProps<{
  reservation: Reservation;
}>();
const stayDateDiff = computed(() => formatDiffDays(props.reservation.stayStartAt, props.reservation.stayEndAt));

function deleteItem() {
  const itemId = props.reservation.id;
  const itemName = props.reservation.name;

  $q.dialog({
    title: "삭제",
    message: `정말로 ${itemName}님의 ${formatStayTitle(stayDateDiff.value)} 예약 정보를 삭제하시겠습니까?`,
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
</script>
