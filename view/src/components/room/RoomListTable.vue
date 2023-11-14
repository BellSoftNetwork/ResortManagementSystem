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
    title="객실"
    flat
    bordered
    binary-state-sort
  >
    <template v-slot:top-right>
      <div class="row q-gutter-sm">
        <q-btn
          :to="{ name: 'CreateRoom' }"
          icon="add"
          color="grey"
          dense
          round
          flat
        />
      </div>
    </template>

    <template #body-cell-number="props">
      <q-td key="number" :props="props">
        <q-btn
          :to="{ name: 'Room', params: { id: props.row.id } }"
          class="full-width"
          align="left"
          color="primary"
          dense
          flat
        >{{ props.row.number }}
        </q-btn>
      </q-td>
    </template>

    <template #body-cell-actions="props">
      <q-td key="actions" :props="props">
        <q-btn dense round flat color="grey" icon="edit"
               :to="{ name: 'EditRoom', params: { id: props.row.id } }"></q-btn>
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
const filter = ref("")
const pagination = ref({
  sortBy: "number",
  descending: false,
  page: 1,
  rowsPerPage: 15,
  rowsNumber: 10,
})
const columns = [
  {
    name: "number",
    field: "number",
    label: "객실 번호",
    align: "left",
    required: true,
    sortable: true,
  },
  {
    name: "peekPrice",
    field: "peekPrice",
    label: "성수기 예약금",
    align: "left",
    headerStyle: "width: 10%",
    required: true,
    sortable: true,
    format: (value) => formatPrice(value),
  },
  {
    name: "offPeekPrice",
    field: "offPeekPrice",
    label: "비성수기 예약금",
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
    name: "createdAt",
    field: "createdAt",
    label: "생성 시각",
    align: "left",
    headerStyle: "width: 15%",
    required: true,
    sortable: true,
    format: val => dayjs(val).format("YYYY-MM-DD HH:mm:ss"),
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

  api.get(`/api/v1/rooms?size=${rowsPerPage}&page=${page - 1}&sort=${sortBy},${descending ? "desc" : "asc"}`)
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
  const itemName = row.number

  $q.dialog({
    title: "삭제",
    message: `정말로 '${itemName}'을 삭제하시겠습니까?`,
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
    api.delete(`/api/v1/rooms/${itemId}`)
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
      return "정상"
    case "INACTIVE":
      return "이용불가"
    case "DAMAGED":
      return "파손"
    case "CONSTRUCTION":
      return "공사 중"
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
