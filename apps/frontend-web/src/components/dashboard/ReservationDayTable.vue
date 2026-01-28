<template>
  <div class="row q-pt-sm">
    <div class="col q-pa-md-sm">
      <div style="min-height: 500px">
        <q-table
          :columns="columns"
          :rows="reservations"
          row-key="id"
          flat
          bordered
          :pagination="{ rowsPerPage: 20, sortBy: 'type', descending: false }"
        >
          <template #top>
            <div class="row justify-between items-center q-pa-sm full-width">
              <div class="row items-center">
                <q-btn icon="chevron_left" color="primary" flat round dense @click="emit('prevDate')" />
                <div class="text-h6">{{ selectedDate }} 예약 현황</div>
                <q-btn icon="chevron_right" color="primary" flat round dense @click="emit('nextDate')" />
              </div>
              <div class="row q-gutter-sm">
                <q-badge v-if="checkInOutCounts[selectedDate]?.checkIn > 0" color="positive" class="q-px-sm">
                  입실 {{ checkInOutCounts[selectedDate]?.checkIn }}
                </q-badge>
                <q-badge v-if="checkInOutCounts[selectedDate]?.checkOut > 0" color="negative" class="q-px-sm">
                  퇴실 {{ checkInOutCounts[selectedDate]?.checkOut }}
                </q-badge>
                <q-badge v-if="!reservations.length" color="grey" class="q-px-sm"> 예약 없음</q-badge>
              </div>
            </div>
          </template>

          <template #header="props">
            <q-tr :props="props">
              <q-th v-for="col in props.cols" :key="col.name" :props="props" class="bg-blue-1">
                {{ col.label }}
              </q-th>
            </q-tr>
          </template>

          <template #body-cell-type="props">
            <q-td key="type" :props="props">
              <q-badge :color="getTypeColor(props.row.type)" text-color="white" class="q-px-sm">
                {{ props.row.type }}
              </q-badge>
            </q-td>
          </template>

          <template #body-cell-rooms="props">
            <q-td key="rooms" :props="props">
              <div v-if="props.row.rooms.length !== 0">
                <span v-for="room in props.row.rooms" :key="room.id">
                  <div v-if="authStore.isAdminRole">
                    <q-btn
                      :to="{
                        name: 'Room',
                        params: { id: room.id },
                      }"
                      align="left"
                      color="primary"
                      dense
                      flat
                      >{{ room.number }}
                    </q-btn>
                  </div>
                  <div v-else>
                    {{ room.number }}
                  </div>
                </span>
              </div>
              <div v-else class="text-grey">미배정</div>
            </q-td>
          </template>

          <template #body-cell-name="props">
            <q-td key="name" :props="props">
              <div v-if="authStore.isAdminRole">
                <q-btn
                  :to="{
                    name: 'Reservation',
                    params: { id: props.row.id },
                  }"
                  class="full-width"
                  align="left"
                  color="primary"
                  dense
                  flat
                  >{{ props.row.name }}
                </q-btn>
              </div>
              <div v-else>
                {{ props.row.name }}
              </div>
            </q-td>
          </template>

          <template #body-cell-missPrice="props">
            <q-td :props="props" key="missPrice" :class="missPriceBackgroundColor(props.row)">
              {{ formatPrice(props.row.missPrice) }}
            </q-td>
          </template>

          <template #body-cell-note="props">
            <q-td :props="props" key="note">
              <q-btn
                v-if="props.row.note"
                @click="
                  $q.dialog({
                    title: `${props.row.name}님 예약 메모`,
                    message: props.row.note,
                  })
                "
                color="primary"
              >
                메모 확인
              </q-btn>
            </q-td>
          </template>

          <template #no-data>
            <div class="full-width row justify-center q-py-md">
              <div class="text-grey text-body1">이 날짜에 등록된 예약이 없습니다.</div>
            </div>
          </template>
        </q-table>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useQuasar } from "quasar";
import { useAuthStore } from "stores/auth";
import { formatPrice } from "src/util/format-util";
import { Reservation } from "src/schema/reservation";
import { useReservationCalendar } from "src/composables/useReservationCalendar";

interface Props {
  reservations: Reservation[];
  selectedDate: string;
  checkInOutCounts: { [date: string]: { checkIn: number; checkOut: number } };
  columns: any[];
}

interface Emits {
  (e: "prevDate"): void;
  (e: "nextDate"): void;
}

defineProps<Props>();
const emit = defineEmits<Emits>();

const $q = useQuasar();
const authStore = useAuthStore();
const { getTypeColor } = useReservationCalendar();

function missPriceBackgroundColor(value: Reservation) {
  if (value.paymentMethod.requireUnpaidAmountCheck === false) return "";

  if (value.missPrice > 0) return "bg-warning";

  return "";
}
</script>
