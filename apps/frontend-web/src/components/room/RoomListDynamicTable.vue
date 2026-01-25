<template>
  <q-table
    @request="onRequest"
    ref="tableRef"
    v-model:pagination="pagination"
    :loading="status.isLoading"
    :columns="columns"
    :rows="rooms"
    style="height: 90vh"
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
        <q-btn
          dense
          round
          flat
          color="grey"
          icon="edit"
          :to="{ name: 'EditRoom', params: { id: props.row.id } }"
        ></q-btn>
        <q-btn dense round flat color="grey" icon="delete" @click="deleteItem(props.row)"></q-btn>
      </q-td>
    </template>
  </q-table>
</template>

<script setup lang="ts">
import { onMounted, ref, watch } from "vue";
import { useQuasar } from "quasar";
import { getRoomFieldDetail, Room } from "src/schema/room";
import { convertTableColumnDef } from "src/util/table-util";
import { deleteRoom, fetchRooms } from "src/api/v1/room";
import { formatSortParam } from "src/util/query-string-util";
import { useRoute, useRouter } from "vue-router";
import { getErrorMessage } from "src/util/errorHandler";

const $q = useQuasar();
const status = ref({
  isLoading: false,
});
const tableRef = ref();
const route = useRoute();
const router = useRouter();

const defaultConfig = {
  pagination: {
    sortBy: "number",
    descending: false,
    page: 1,
    rowsPerPage: 15,
  },
};
const pagination = ref({
  sortBy: defaultConfig.pagination.sortBy,
  descending: defaultConfig.pagination.descending,
  page: defaultConfig.pagination.page,
  rowsPerPage: defaultConfig.pagination.rowsPerPage,
  rowsNumber: 0,
});
const columns = [
  {
    ...getColumnDef("roomGroup"),
    align: "left",
    required: true,
    sortable: true,
  },
  {
    ...getColumnDef("number"),
    align: "left",
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
const rooms = ref<Room[]>([]);

loadQueryString();

watch(route, () => {
  loadQueryString();
});

function loadQueryString() {
  pagination.value.sortBy = route.query.sortBy?.toString() ?? defaultConfig.pagination.sortBy;
  pagination.value.descending = Boolean(route.query.descending ?? defaultConfig.pagination.descending);
  pagination.value.page = Number(route.query.page ?? defaultConfig.pagination.page);
  pagination.value.rowsPerPage = Number(route.query.rowsPerPage ?? defaultConfig.pagination.rowsPerPage);
}

function setPaginationQuery() {
  router.push({
    query: {
      ...route.query,
      page: pagination.value.page !== defaultConfig.pagination.page ? pagination.value.page : undefined,
      rowsPerPage:
        pagination.value.rowsPerPage !== defaultConfig.pagination.rowsPerPage
          ? pagination.value.rowsPerPage
          : undefined,
      sortBy: pagination.value.sortBy !== defaultConfig.pagination.sortBy ? pagination.value.sortBy : undefined,
      descending:
        pagination.value.descending !== defaultConfig.pagination.descending ? pagination.value.descending : undefined,
    },
  });
}

function getColumnDef(field: string) {
  return convertTableColumnDef(getRoomFieldDetail(field));
}

function onRequest(props) {
  const { page, rowsPerPage, sortBy, descending } = props.pagination;

  status.value.isLoading = true;

  fetchRooms({
    page: page - 1,
    size: rowsPerPage,
    sort: formatSortParam({ field: sortBy, isDescending: descending }),
  })
    .then((response) => {
      rooms.value = response.values;

      const page = response.page;

      pagination.value.rowsNumber = page.totalElements;
      pagination.value.page = page.index + 1;
      pagination.value.rowsPerPage = page.size;
      pagination.value.sortBy = sortBy;
      pagination.value.descending = descending;

      setPaginationQuery();
    })
    .catch((error) => {
      rooms.value = [];

      console.error(error);
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
    })
    .finally(() => {
      status.value.isLoading = false;
    });
}

function reloadData() {
  tableRef.value.requestServerInteraction();
}

function deleteItem(row: Room) {
  const itemId = row.id;
  const itemName = row.number;

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
    deleteRoom(itemId)
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
