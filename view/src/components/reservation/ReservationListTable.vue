<template>
  <q-table
    @request="onRequest"
    ref="tableRef"
    v-model:pagination="pagination"
    :loading="status.isLoading"
    :columns="columns"
    :rows="responseData.values"
    :filter="filter"
    style="height: 90vh;"
    row-key="id"
    title="다가오는 예약"
    flat
    bordered
    binary-state-sort
  >
    <template v-slot:top-right>
      <div class="row q-gutter-sm">
        <q-btn
          color="primary"
          label="상세 검색"
          icon="search"
        >
          <q-menu
            anchor="bottom end"
            self="top end"
          >
            <div class="row no-wrap q-pa-md">
              <q-input
                v-model="filter.peopleInfo"
                debounce="300"
                placeholder="홍길동"
                label="예약자 정보"
                class="fit"
              />
            </div>

            <div class="row no-wrap q-pa-md">
              <q-input v-model="filter.stayStartAt" mask="####-##-##" :readonly="true" outlined>
                <template v-slot:append>
                  <q-icon name="event" class="cursor-pointer">
                    <q-popup-proxy cover transition-show="scale" transition-hide="scale">
                      <q-date v-model="filter.stayStartAt" mask="YYYY-MM-DD">
                        <div class="row items-center justify-end">
                          <q-btn v-close-popup label="Close" color="primary" flat />
                        </div>
                      </q-date>
                    </q-popup-proxy>
                  </q-icon>
                </template>
              </q-input>
              <span class="self-center q-mx-sm">~</span>
              <q-input v-model="filter.stayEndAt" mask="####-##-##" :readonly="true" outlined>
                <template v-slot:append>
                  <q-icon name="event" class="cursor-pointer">
                    <q-popup-proxy cover transition-show="scale" transition-hide="scale">
                      <q-date v-model="filter.stayEndAt" mask="YYYY-MM-DD">
                        <div class="row items-center justify-end">
                          <q-btn v-close-popup label="Close" color="primary" flat />
                        </div>
                      </q-date>
                    </q-popup-proxy>
                  </q-icon>
                </template>
              </q-input>
            </div>
          </q-menu>
        </q-btn>

        <q-btn
          :to="{ name: 'CreateReservation' }"
          icon="add"
          color="grey"
          dense
          round
          flat
        />
      </div>
    </template>

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

    <template #body-cell-actions="props">
      <q-td key="actions" :props="props">
        <q-btn dense round flat color="grey" icon="edit"
               :to="{ name: 'EditReservation', params: { id: props.row.id } }"></q-btn>
        <q-btn dense round flat color="grey" icon="delete" @click="deleteItem(props.row)"></q-btn>
      </q-td>
    </template>
  </q-table>
</template>

<script setup>
import { onBeforeMount, onMounted, ref } from "vue"
import dayjs from "dayjs"
import { useQuasar } from "quasar"
import { api } from "boot/axios"

const $q = useQuasar()
const status = ref({
  isLoading: false,
  isLoaded: false,
  isPatching: false,
})
const tableRef = ref()
const filter = ref({
  peopleInfo: "",
  stayStartAt: dayjs().format("YYYY-MM-DD"),
  stayEndAt: dayjs().add(3, "M").format("YYYY-MM-DD"),
})
const pagination = ref({
  sortBy: "stayStartAt",
  descending: true,
  page: 1,
  rowsPerPage: 15,
  rowsNumber: 10,
})
const columns = [
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
    headerStyle: "width: 10%",
    required: true,
    sortable: true,
  },
  {
    name: "name",
    field: "name",
    label: "예약자명",
    align: "left",
    headerStyle: "width: 10%",
    required: true,
    sortable: true,
  },
  {
    name: "phone",
    field: "phone",
    label: "예약자 전화번호",
    align: "left",
    headerStyle: "width: 10%",
    required: true,
    sortable: true,
  },
  {
    name: "peopleCount",
    field: "peopleCount",
    label: "숙박 인원",
    align: "left",
    headerStyle: "width: 10%",
    required: true,
    sortable: true,
  },
  {
    name: "stayStartAt",
    field: "stayStartAt",
    label: "입실일",
    align: "left",
    headerStyle: "width: 10%",
    required: true,
    sortable: true,
    format: val => dayjs(val).format("YYYY-MM-DD"),
  },
  {
    name: "stayEndAt",
    field: "stayEndAt",
    label: "퇴실일",
    align: "left",
    headerStyle: "width: 10%",
    required: true,
    sortable: true,
    format: val => dayjs(val).format("YYYY-MM-DD"),
  },
  {
    name: "checkInAt",
    field: "checkInAt",
    label: "체크인 시각",
    align: "left",
    headerStyle: "width: 15%",
    required: true,
    sortable: true,
    format: val => val ? dayjs(val).format("YYYY-MM-DD HH:mm:ss") : "",
  },
  {
    name: "checkOutAt",
    field: "checkOutAt",
    label: "체크아웃 시각",
    align: "left",
    headerStyle: "width: 15%",
    required: true,
    sortable: true,
    format: val => val ? dayjs(val).format("YYYY-MM-DD HH:mm:ss") : "",
  },
  {
    name: "price",
    field: "price",
    label: "총 금액",
    align: "left",
    headerStyle: "width: 10%",
    required: true,
    sortable: true,
    format: (value) => formatPrice(value),
  },
  {
    name: "paymentAmount",
    field: "paymentAmount",
    label: "지불 금액",
    align: "left",
    headerStyle: "width: 10%",
    required: true,
    sortable: true,
    format: (value) => formatPrice(value),
  },
  {
    name: "status",
    field: "status",
    label: "상태",
    align: "left",
    headerStyle: "width: 10%",
    required: true,
    sortable: true,
    format: (value) => formatStatus(value),
  },
  {
    name: "updatedAt",
    field: "updatedAt",
    label: "수정 시각",
    align: "left",
    headerStyle: "width: 15%",
    required: true,
    sortable: true,
    format: val => dayjs(val).format("YYYY-MM-DD HH:mm:ss"),
  },
  {
    name: "actions",
    label: "액션",
    align: "center",
    headerStyle: "width: 5%",
  },
]
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

function onRequest(props) {
  const { page, rowsPerPage, sortBy, descending } = props.pagination

  status.value.isLoading = true
  status.value.isLoaded = false

  const queryParams = {
    size: rowsPerPage,
    page: page - 1,
    sort: `${sortBy},${descending ? "desc" : "asc"}`,
    stayStartAt: filter.value.stayStartAt,
    stayEndAt: filter.value.stayEndAt,
  }

  if (filter.value.peopleInfo)
    queryParams.searchText = encodeURIComponent(filter.value.peopleInfo)

  const queryString = Object.keys(queryParams).map(key => `${key}=${queryParams[key]}`)

  api.get(`/api/v1/reservations?${queryString.join("&")}`)
    .then(response => {
      responseData.value = response.data
      const page = responseData.value.page

      pagination.value.rowsNumber = page.totalElements
      pagination.value.page = page.index + 1
      pagination.value.rowsPerPage = page.size
      pagination.value.sortBy = sortBy
      pagination.value.descending = descending

      status.value.isLoaded = true
    }).finally(() => {
    status.value.isLoading = false
  })
}

function reloadData() {
  tableRef.value.requestServerInteraction()
}

function deleteItem(row) {
  const itemId = row.id
  const itemName = row.name

  $q.dialog({
    title: "삭제",
    message: `정말로 ${itemName}님의 예약 정보를 삭제하시겠습니까?`,
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
    api.delete(`/api/v1/reservations/${itemId}`)
      .then(() => {
        reloadData()
      })
      .catch((error) => {
        $q.notify({
          message: error.response.data.message,
          type: "negative",
          actions: [
            {
              icon: "close", color: "white", round: true,
            },
          ],
        })
      })
  })
}

function formatPrice(value) {
  return new Intl.NumberFormat("ko-KR", {
    style: "currency",
    currency: "KRW",
  }).format(value)
}

function formatStatus(value) {
  switch (value) {
    case "NORMAL":
      return "예약 확정"
    case "PENDING":
      return "예약 대기"
    case "CANCEL":
      return "취소 요청"
    case "REFUND":
      return "환불 완료"
    default:
      return value
  }
}

function resetData() {
  responseData.value.values = []
}

onBeforeMount(() => {
  resetData()
})

onMounted(() => {
  reloadData()
})
</script>
