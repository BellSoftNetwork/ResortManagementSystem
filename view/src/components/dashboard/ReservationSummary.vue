<template>
  <q-card flat bordered>
    <q-card-section class="text-h6">
      입실 정보 요약
    </q-card-section>

    <q-card-section>
      <div class="row">
        <div class="col-12 col-md-auto q-pa-sm">
          <q-date
            v-model="date"
            @navigation="changeView"
            :events="events"
            mask="YYYY-MM-DD"
          />
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
              >
                <template #body-cell-reservationMethod="props">
                  <q-td key="reservationMethod" :props="props">{{ props.row.reservationMethod.name }}</q-td>
                </template>

                <template #body-cell-room="props">
                  <q-td key="room" :props="props">
                    <div v-if="props.row.room">
                      <q-btn
                        :to="{ name: 'Room', params: { id: props.row.room.id } }"
                        class="full-width"
                        align="left"
                        color="primary"
                        dense
                        flat
                      >{{ props.row.room.number }}
                      </q-btn>
                    </div>
                    <div v-else class="text-grey">미배정</div>
                  </q-td>
                </template>

                <template #body-cell-name="props">
                  <q-td key="name" :props="props">
                    <q-btn
                      :to="{ name: 'Reservation', params: { id: props.row.id } }"
                      class="full-width"
                      align="left"
                      color="primary"
                      dense
                      flat
                    >{{ props.row.name }}
                    </q-btn>
                  </q-td>
                </template>

                <template #body-cell-missPrice="props">
                  <q-td
                    :props="props"
                    key="missPrice"
                    :class="missPriceBackgroundColor(props.row)"
                  >
                    {{ formatPrice(props.row.missPrice) }}
                  </q-td>
                </template>

                <template #body-cell-note="props">
                  <q-td
                    :props="props"
                    key="note"
                  >
                    <q-btn
                      v-if="props.row.note"
                      @click="$q.dialog({ title: `${props.row.name}님 예약 메모`, message: props.row.note })"
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

<script setup>
import { computed, onBeforeMount, onMounted, ref } from "vue"
import dayjs from "dayjs"
import { useQuasar } from "quasar"
import { api } from "boot/axios"

const $q = useQuasar()
const status = ref({
  isLoading: false,
  isLoaded: false,
  isPatching: false,
})
const filter = ref({
  sort: "stayStartAt",
  stayStartAt: dayjs().startOf("month").format("YYYY-MM-DD"),
  stayEndAt: dayjs().endOf("month").format("YYYY-MM-DD"),
})
const columns = [
  {
    name: "name",
    field: "name",
    label: "예약자명",
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
    name: "reservationMethod",
    field: "reservationMethod",
    label: "예약 수단",
    align: "left",
    required: true,
    sortable: true,
  },
  {
    name: "room",
    field: "room",
    label: "객실",
    align: "left",
    required: true,
    sortable: true,
  },
  {
    name: "note",
    field: "note",
    label: "메모",
    align: "left",
    headerStyle: "width: 10%",
  },
]
const date = ref(dayjs().format("YYYY-MM-DD"))
const reservationsOfDay = computed(() => {
  const reservationMap = {}

  responseData.value.values.forEach((reservation) => {
    if (!Object.keys(reservationMap).includes(reservation.stayStartAt))
      reservationMap[reservation.stayStartAt] = []

    reservation.missPrice = reservation.price - reservation.paymentAmount
    reservationMap[reservation.stayStartAt].push(reservation)
  })

  return reservationMap
})
const events = computed(() => Object.keys(reservationsOfDay.value).map((date) => dayjs(date).format("YYYY/MM/DD")))
const responseData = ref({
  page: {
    index: 0,
    size: 0,
    totalPages: 0,
    totalElements: 0,
  },
  values: [
    {
      id: 1,
      number: "101",
      peekPrice: 100000,
      offPeekPrice: 80000,
      status: "ACTIVE",
      createdAt: "2021-01-01T00:00:00.000Z",
      updatedAt: "2021-01-01T00:00:00.000Z",
    },
  ],
})

function fetchData() {
  status.value.isLoading = true
  status.value.isLoaded = false

  const queryParams = {
    stayStartAt: filter.value.stayStartAt,
    stayEndAt: filter.value.stayEndAt,
  }

  const queryString = Object.keys(queryParams).map(key => `${key}=${queryParams[key]}`)

  api.get(`/api/v1/reservations?${queryString.join("&")}`)
    .then(response => {
      responseData.value = response.data

      status.value.isLoaded = true
    }).finally(() => {
    status.value.isLoading = false
  })
}

function changeView(view) {
  const year = view.year
  const month = view.month

  filter.value.stayStartAt = dayjs(`${year}-${month}-01`).startOf("month").format("YYYY-MM-DD")
  filter.value.stayEndAt = dayjs(`${year}-${month}-01`).endOf("month").format("YYYY-MM-DD")

  fetchData()
}

function formatPrice(value) {
  return new Intl.NumberFormat("ko-KR", {
    style: "currency",
    currency: "KRW",
  }).format(value)
}

function missPriceBackgroundColor(value) {
  if (value.reservationMethod.requireUnpaidAmountCheck === false)
    return ""

  if (value.missPrice > 0)
    return "bg-warning"

  return ""
}

function resetData() {
  responseData.value.values = []
}

onBeforeMount(() => {
  resetData()
})

onMounted(() => {
  fetchData()
})
</script>
