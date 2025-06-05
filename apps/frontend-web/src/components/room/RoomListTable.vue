<template>
  <q-table
    :columns="columns"
    :rows="props.roomLastStayDetails"
    row-key="id"
    title="객실"
    flat
    bordered
    binary-state-sort
  >
    <template v-slot:top-right>
      <div class="row q-gutter-sm">
        <q-btn :to="{ name: 'CreateRoom' }" icon="add" color="grey" dense round flat />
      </div>
    </template>

    <template #body-cell-number="props">
      <q-td key="number" :props="props">
        <q-btn
          :to="{ name: 'Room', params: { id: props.row.room.id } }"
          class="full-width"
          align="left"
          color="primary"
          dense
          flat
          >{{ props.row.room.number }}
        </q-btn>
      </q-td>
    </template>

    <template #body-cell-lastReservation="props">
      <q-td key="lastReservation" :props="props">
        <q-btn
          v-if="props.row.lastReservation"
          :to="{ name: 'Reservation', params: { id: props.row.lastReservation.id } }"
          class="full-width"
          align="left"
          color="primary"
          dense
          flat
          >{{ props.row.lastReservation.name }}
        </q-btn>
      </q-td>
    </template>
  </q-table>
</template>

<script setup lang="ts">
import { getRoomFieldDetail } from "src/schema/room";
import { convertTableColumnDef } from "src/util/table-util";
import { RoomLastStayDetail } from "src/api/v1/room-group";

const props = defineProps<{
  roomLastStayDetails: RoomLastStayDetail[];
}>();
const columns = [
  {
    ...getColumnDef("number"),
    align: "left",
    required: true,
    sortable: true,
    field: (row: RoomLastStayDetail) => row.room.number,
  },
  {
    ...getColumnDef("status"),
    align: "left",
    headerStyle: "width: 10%",
    required: true,
    sortable: true,
    field: (row: RoomLastStayDetail) => row.room.status,
  },
  {
    label: "최근 예약 정보",
    name: "lastReservation",
    align: "left",
    required: true,
    sortable: true,
  },
  {
    label: "최근 퇴실일",
    name: "lastStayAt",
    align: "left",
    required: true,
    sortable: true,
    field: (row: RoomLastStayDetail) => row.lastReservation?.stayEndAt,
  },
];

function getColumnDef(field: string) {
  return convertTableColumnDef(getRoomFieldDetail(field));
}
</script>
