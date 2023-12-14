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
    selection="multiple"
    v-model:selected="selected"
    flat
    bordered
    binary-state-sort
  ></q-table>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { getRoomFieldDetail, Room } from "src/schema/room";
import { convertTableColumnDef } from "src/util/table-util";
import { fetchRooms } from "src/api/v1/room";
import { formatSortParam } from "src/util/query-string-util";
import { Reservation } from "src/schema/reservation";
import { useQuasar } from "quasar";

const $q = useQuasar();
const props = defineProps<{
  selected: Room[];
  parentReservation?: Reservation;
  stayStartAt: string;
  stayEndAt: string;
}>();
const emit = defineEmits(["update:selected"]);

const status = ref({
  isLoading: false,
  isLoaded: false,
  isPatching: false,
});

const selected = computed({
  get() {
    return props.selected;
  },
  set(value) {
    emit("update:selected", value);
  },
});

const tableRef = ref();
const filter = ref("");
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
const rooms = ref<Room[]>();

function getColumnDef(field: string) {
  return convertTableColumnDef(getRoomFieldDetail(field));
}

function onRequest(tableProps) {
  const { page, rowsPerPage, sortBy, descending } = tableProps.pagination;

  status.value.isLoading = true;
  status.value.isLoaded = false;

  fetchRooms({
    page: page - 1,
    size: rowsPerPage,
    sort: formatSortParam({ field: sortBy, isDescending: descending }),
    stayStartAt: props.stayStartAt,
    stayEndAt: props.stayEndAt,
    status: "NORMAL",
    excludeReservationId: props.parentReservation?.id ?? undefined,
  })
    .then((response) => {
      rooms.value = response.values;

      const reservation = props.parentReservation;
      if (reservation) {
        const reservationRoomIds = reservation.rooms.map((room) => room.id);
        const selectedRooms = rooms.value.filter((room) => reservationRoomIds.includes(room.id));
        selected.value = selectedRooms;

        const selectedRoomIds = selectedRooms.map((room) => room.id);
        const unavailableRooms = reservation.rooms.filter((room) => !selectedRoomIds.includes(room.id));
        if (unavailableRooms.length > 0) {
          $q.notify({
            message:
              `기존에 배정된 ${unavailableRooms.length}개의 객실이 해당 기간에 이용할 수 없어 제외되었습니다.<br />` +
              `이용 불가 객실: ${unavailableRooms.map((room) => room.number).join(", ")}`,
            type: "warning",
            html: true,
            actions: [
              {
                icon: "close",
                color: "white",
                round: true,
              },
            ],
          });
        }
      }

      const page = response.page;

      pagination.value.rowsNumber = page.totalElements;
      pagination.value.page = page.index + 1;
      pagination.value.rowsPerPage = page.size;
      pagination.value.sortBy = sortBy;
      pagination.value.descending = descending;

      status.value.isLoaded = true;
    })
    .finally(() => {
      status.value.isLoading = false;
    });
}

function reloadData() {
  tableRef.value.requestServerInteraction();
}

onMounted(() => {
  reloadData();
});
</script>
