<template>
  <q-table
    @request="onRequest"
    ref="tableRef"
    v-model:pagination="pagination"
    :loading="status.isLoading"
    :columns="columns"
    :rows="rooms"
    :filter="filter"
    row-key="id"
    selection="single"
    v-model:selected="selected"
    flat
    bordered
    binary-state-sort
  ></q-table>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from "vue"
import { getRoomFieldDetail, Room } from "src/schema/room"
import { convertTableColumnDef } from "src/util/table-util"
import { fetchRooms } from "src/api/v1/room"
import { formatSortParam } from "src/util/query-string-util"

const props = defineProps<{
  selected: Room[];
  firstValue?: Room;
  stayStartAt: string;
  stayEndAt: string;
}>();
const emit = defineEmits(["update:selected"])

const status = ref({
  isLoading: false,
  isLoaded: false,
  isPatching: false,
});

const selected = computed({
  get() {
    return props.selected
  },
  set(value) {
    emit("update:selected", value)
  },
});

const tableRef = ref()
const filter = ref("")
const pagination = ref({
  sortBy: "number",
  descending: false,
  page: 1,
  rowsPerPage: 15,
  rowsNumber: 10,
});
const columns = [
  {
    ...getColumnDef("number"),
    align: "left",
    required: true,
    sortable: true,
  },
  {
    ...getColumnDef("peekPrice"),
    align: "left",
    headerStyle: "width: 10%",
    required: true,
    sortable: true,
  },
  {
    ...getColumnDef("offPeekPrice"),
    align: "left",
    headerStyle: "width: 10%",
    required: true,
    sortable: true,
  },
  {
    ...getColumnDef("status"),
    align: "left",
    headerStyle: "width: 10%",
    required: true,
    sortable: true,
  },
];
const rooms = ref<Room[]>()

function getColumnDef(field: string) {
  return convertTableColumnDef(getRoomFieldDetail(field))
}

function onRequest(tableProps) {
  const { page, rowsPerPage, sortBy, descending } = tableProps.pagination

  status.value.isLoading = true
  status.value.isLoaded = false

  fetchRooms({
    page: page - 1,
    size: rowsPerPage,
    sort: formatSortParam({ field: sortBy, isDescending: descending }),
    stayStartAt: props.stayStartAt,
    stayEndAt: props.stayEndAt,
    status: "NORMAL",
  })
    .then((response) => {
      rooms.value = response.values
      if (props.firstValue) rooms.value.push(props.firstValue)

      const page = response.page

      pagination.value.rowsNumber = page.totalElements
      pagination.value.page = page.index + 1
      pagination.value.rowsPerPage = page.size
      pagination.value.sortBy = sortBy
      pagination.value.descending = descending

      status.value.isLoaded = true
    })
    .finally(() => {
      status.value.isLoading = false
    });
}

function reloadData() {
  tableRef.value.requestServerInteraction()
}

onMounted(() => {
  reloadData()
});
</script>
