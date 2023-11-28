<template>
  <q-table
    @request="onRequest"
    ref="tableRef"
    v-model:pagination="pagination"
    :loading="status.isLoading"
    :columns="columns"
    :rows="responseData.values"
    :filter="filter"
    row-key="id"
    title="변경 이력"
    flat
    bordered
    binary-state-sort
  >
    <template #body-cell-historyType="props">
      <q-td key="historyType" :props="props">
        <q-chip
          :icon="historyTypeMap[props.row.historyType].icon"
          :color="historyTypeMap[props.row.historyType].color"
          outline
        >
          {{ historyTypeMap[props.row.historyType].name }}
        </q-chip>
      </q-td>
    </template>

    <template #body-cell-updatedValue="props">
      <q-td key="updatedValue" :props="props">
        <div class="row q-gutter-sm">
          <q-card
            v-for="field in props.row.updatedFields"
            v-bind:key="field"
            bordered
          >
            <q-card-section horizontal>
              <q-card-section class="bg-blue-3 q-pa-xs">
                {{ columnMap[field].name }}
              </q-card-section>

              <q-card-section class="bg-grey-4 q-pa-xs">
                {{ formatValue(field, props.row.entity[field]) }}
              </q-card-section>
            </q-card-section>
          </q-card>
        </div>
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
const props = defineProps({
  id: Number,
})
const status = ref({
  isLoading: false,
  isLoaded: false,
  isPatching: false,
})
const tableRef = ref()
const filter = ref("")
const pagination = ref({
  sortBy: "updatedAt",
  descending: true,
  page: 1,
  rowsPerPage: 15,
  rowsNumber: 10,
})
const columns = [
  {
    name: "historyType",
    field: "historyType",
    label: "타입",
    align: "left",
    headerStyle: "width: 10%",
    required: true,
  },
  {
    name: "updatedValue",
    field: "updatedValue",
    label: "변경 내역",
    align: "left",
    required: true,
  },
  {
    name: "updatedAt",
    field: "historyCreatedAt",
    label: "변경 시각",
    align: "left",
    headerStyle: "width: 15%",
    required: true,
    sortable: true,
    format: val => dayjs(val).format("YYYY-MM-DD HH:mm:ss"),
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
      entity: {
        name: "101",
      },
      historyCreatedAt: "2021-01-01T00:00:00.000Z",
      historyType: "CREATED",
      updatedFields: ["number"],
    },
  ],
})
const columnMap = {
  number: { name: "객실 번호" },
  peekPrice: { name: "성수기 예약금", format: (value) => formatPrice(value) },
  offPeekPrice: { name: "비성수기 예약금", format: (value) => formatPrice(value) },
  description: { name: "설명" },
  note: { name: "메모" },
  status: { name: "상태", format: (value) => formatStatus(value) },
}
const historyTypeMap = {
  CREATED: { name: "생성", color: "primary", icon: "add" },
  UPDATED: { name: "변경", color: "warning", icon: "edit" },
  DELETED: { name: "삭제", color: "red", icon: "remove" },
}

function onRequest(tableProps) {
  const { page, rowsPerPage, sortBy, descending } = tableProps.pagination

  status.value.isLoading = true
  status.value.isLoaded = false

  api.get(`/api/v1/rooms/${props.id}/histories?size=${rowsPerPage}&page=${page - 1}&sort=${sortBy},${descending ? "desc" : "asc"}`)
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

function formatValue(field, value) {
  return columnMap[field].format ? columnMap[field].format(value) : value
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
