<template>
  <q-table
    @request="onRequest"
    ref="tableRef"
    v-model:pagination="pagination"
    :loading="loading"
    :columns="columns"
    :rows="rows"
    style="height: 90vh"
    row-key="id"
    title="객실 그룹"
    flat
    bordered
    binary-state-sort
  >
    <template v-slot:top-right>
      <div class="row q-gutter-sm">
        <q-btn :to="{ name: 'CreateRoomGroup' }" icon="add" color="grey" dense round flat />
      </div>
    </template>

    <template #body-cell-name="props">
      <q-td key="name" :props="props">
        <q-btn
          :to="{ name: 'RoomGroup', params: { id: props.row.id } }"
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
        <q-btn
          dense
          round
          flat
          color="grey"
          icon="edit"
          :to="{ name: 'EditRoomGroup', params: { id: props.row.id } }"
        ></q-btn>
        <q-btn dense round flat color="grey" icon="delete" @click="deleteItem(props.row)"></q-btn>
      </q-td>
    </template>
  </q-table>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useQuasar } from "quasar";
import { convertTableColumnDef } from "src/util/table-util";
import { getErrorMessage } from "src/util/errorHandler";
import { getRoomGroupFieldDetail, RoomGroup } from "src/schema/room-group";
import { deleteRoomGroup, fetchRoomGroups } from "src/api/v1/room-group";
import { useTable } from "src/composables/useTable";

const $q = useQuasar();
const tableRef = ref();

const { pagination, loading, rows, onRequest } = useTable<RoomGroup>({
  fetchFn: fetchRoomGroups,
  defaultPagination: {
    sortBy: "name",
    descending: false,
    page: 1,
    rowsPerPage: 15,
  },
  onError: (error) => {
    $q.notify({
      message: getErrorMessage(error),
      type: "negative",
      actions: [
        {
          icon: "close",
          color: "white",
          round: true,
        },
      ],
    });
  },
});

const columns = [
  {
    ...getColumnDef("name"),
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
    ...getColumnDef("createdAt"),
    align: "left",
    headerStyle: "width: 15%",
    required: true,
    sortable: true,
  },
  {
    ...getColumnDef("updatedAt"),
    align: "left",
    headerStyle: "width: 15%",
    required: true,
    sortable: true,
  },
  {
    name: "actions",
    label: "액션",
    align: "center",
    headerStyle: "width: 5%",
  },
];

function getColumnDef(field: string) {
  return convertTableColumnDef(getRoomGroupFieldDetail(field));
}

function reloadData() {
  tableRef.value.requestServerInteraction();
}

function deleteItem(row: RoomGroup) {
  const itemId = row.id;
  const itemName = row.name;

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
    deleteRoomGroup(itemId)
      .then(() => {
        reloadData();
      })
      .catch((error) => {
        $q.notify({
          message: getErrorMessage(error),
          type: "negative",
          actions: [
            {
              icon: "close",
              color: "white",
              round: true,
            },
          ],
        });
      });
  });
}

onMounted(() => {
  reloadData();
});
</script>
