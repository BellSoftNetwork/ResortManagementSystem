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
    selection="single"
    v-model:selected="selected"
    flat
    bordered
    binary-state-sort
  ></q-table>
</template>

<script setup>
import { computed, onBeforeMount, onMounted, ref } from "vue"
import { useQuasar } from "quasar"
import { api } from "boot/axios"

const props = defineProps([
  "selected",
  "firstValue",
  "stayStartAt",
  "stayEndAt",
])
const emit = defineEmits(["update:selected"])

const $q = useQuasar()
const status = ref({
  isLoading: false,
  isLoaded: false,
  isPatching: false,
})

const selected = computed({
  get() {
    return props.selected
  },
  set(value) {
    emit("update:selected", value)
  },
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

function onRequest(tableProps) {
  const { page, rowsPerPage, sortBy, descending } = tableProps.pagination

  status.value.isLoading = true
  status.value.isLoaded = false

  const queryParams = {
    size: rowsPerPage,
    page: page - 1,
    sort: `${sortBy},${descending ? "desc" : "asc"}`,
    stayStartAt: props.stayStartAt,
    stayEndAt: props.stayEndAt,
    status: "NORMAL",
  }

  const queryString = Object.keys(queryParams).map(key => `${key}=${queryParams[key]}`)

  api.get(`/api/v1/rooms?${queryString.join("&")}`)
    .then(response => {
      responseData.value = response.data
      if (props.firstValue)
        responseData.value.values.push(props.firstValue)

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
