<template>
  <q-card flat bordered>
    <q-card-section class="text-h6">입실 정보 요약</q-card-section>

    <q-card-section>
      <div class="row">
        <div class="col-12 col-md-auto q-pa-sm">
          <q-date v-model="date" @navigation="changeView" :events="events" mask="YYYY-MM-DD" />
        </div>

        <div class="col-12 col-md q-pa-md-sm">
          <q-tab-panels v-model="date">
            <q-tab-panel
              v-for="(reservations, stayStartAt) in reservationsOfDay"
              :key="stayStartAt"
              :name="stayStartAt"
              class="q-px-none"
            >
              <q-table
                :columns="columns"
                :rows="reservations"
                row-key="id"
                :title="stayStartAt"
                flat
                bordered
                :pagination="{ rowsPerPage: 20 }"
              >
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
              </q-table>
            </q-tab-panel>
          </q-tab-panels>
        </div>
      </div>
    </q-card-section>
  </q-card>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import dayjs from "dayjs";
import { useQuasar } from "quasar";
import { useAuthStore } from "stores/auth";
import { formatDate, formatPrice } from "src/util/format-util";
import { convertTableColumnDef } from "src/util/table-util";
import { getReservationFieldDetail, Reservation } from "src/schema/reservation";
import { fetchReservations } from "src/api/v1/reservation";

const $q = useQuasar();
const authStore = useAuthStore();
const status = ref({
  isLoading: false,
  isLoaded: false,
  isPatching: false,
});
const filter = ref({
  sort: "stayStartAt",
  stayStartAt: dayjs().startOf("month").format("YYYY-MM-DD"),
  stayEndAt: dayjs().endOf("month").format("YYYY-MM-DD"),
});
const columns = [
  {
    name: "type",
    field: "type",
    label: "구분",
    align: "left",
    required: true,
    sortable: true,
  },
  {
    ...getColumnDef("name"),
    align: "left",
    required: true,
    sortable: true,
  },
  {
    name: "missPrice",
    field: "missPrice",
    label: "미수금",
    align: "left",
    required: true,
  },
  {
    ...getColumnDef("paymentMethod"),
    align: "left",
    required: true,
    sortable: true,
  },
  {
    ...getColumnDef("rooms"),
    align: "left",
    required: true,
    sortable: true,
  },
  {
    ...getColumnDef("note"),
    align: "left",
    headerStyle: "width: 10%",
  },
];
const date = ref(formatDate());
const reservationsOfDay = ref({});
const events = computed(() => Object.keys(reservationsOfDay.value).map((date) => dayjs(date).format("YYYY/MM/DD")));

function getColumnDef(field: string) {
  return convertTableColumnDef(getReservationFieldDetail(field));
}

function fetchData() {
  status.value.isLoading = true;
  status.value.isLoaded = false;

  fetchReservations({
    stayStartAt: filter.value.stayStartAt,
    stayEndAt: filter.value.stayEndAt,
    status: "NORMAL",
    type: "STAY",
    size: 200, // TODO: 임시 수정
  })
    .then((response) => {
      reservationsOfDay.value = formatReservations(response.values);

      status.value.isLoaded = true;
    })
    .finally(() => {
      status.value.isLoading = false;
    });
}

function formatReservations(reservations: Reservation[]) {
  const reservationMap: {
    [date: string]: Reservation[];
  } = {};

  reservations.forEach((reservation) => {
    reservation.missPrice = reservation.price - reservation.paymentAmount;

    for (const [index, date] of getDateArray(reservation.stayStartAt, reservation.stayEndAt).entries()) {
      if (!Object.keys(reservationMap).includes(date)) reservationMap[date] = [];

      const reservationCopy = { ...reservation, type: "N/A" };

      if (index === 0) reservationCopy.type = "입실";
      else if (date === reservation.stayEndAt) reservationCopy.type = "퇴실";
      else reservationCopy.type = "연박";

      reservationMap[date].push(reservationCopy);
    }
  });

  return reservationMap;
}

function getDateArray(startDate: string, endDate: string) {
  const stayStartDate = formatDate(startDate);
  const stayEndAt = formatDate(endDate);
  const dateArray = [];

  for (let date = stayStartDate; date <= stayEndAt; date = dayjs(date).add(1, "day").format("YYYY-MM-DD")) {
    dateArray.push(date);
  }

  return dateArray;
}

function changeView(view) {
  const year = view.year;
  const month = view.month;

  filter.value.stayStartAt = dayjs(`${year}-${month}-01`).startOf("month").format("YYYY-MM-DD");
  filter.value.stayEndAt = dayjs(`${year}-${month}-01`).endOf("month").format("YYYY-MM-DD");

  fetchData();
}

function missPriceBackgroundColor(value) {
  if (value.paymentMethod.requireUnpaidAmountCheck === false) return "";

  if (value.missPrice > 0) return "bg-warning";

  return "";
}

onMounted(() => {
  fetchData();
});
</script>
