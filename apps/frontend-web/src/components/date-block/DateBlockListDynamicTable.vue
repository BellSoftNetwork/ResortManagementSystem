<template>
  <q-table
    ref="tableRef"
    v-model:pagination="pagination"
    :loading="loading"
    :columns="columns"
    :rows="rows"
    style="height: 90vh"
    row-key="id"
    title="예약 마감"
    flat
    bordered
    binary-state-sort
    @request="onRequest"
  >
    <template v-slot:top-right>
      <div class="row q-gutter-sm">
        <q-btn icon="add" color="grey" dense round flat @click="openCreateDialog" />
      </div>
    </template>

    <template #body-cell-actions="props">
      <q-td key="actions" :props="props">
        <q-btn dense round flat color="grey" icon="edit" @click="openEditDialog(props.row)"></q-btn>
        <q-btn dense round flat color="grey" icon="delete" @click="deleteItem(props.row)"></q-btn>
        <q-btn dense round flat color="grey" icon="history" @click="openHistoryDialog(props.row)"></q-btn>
      </q-td>
    </template>
  </q-table>

  <DateBlockCreateDialog ref="createDialogRef" @created="reloadData" />
  <DateBlockEditDialog ref="editDialogRef" @updated="reloadData" />
  <DateBlockHistoryDialog v-model="historyDialogOpen" :date-block-id="historyDateBlockId" />
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useQuasar } from "quasar";
import { DateBlock, getDateBlockFieldDetail } from "src/schema/date-block";
import { convertTableColumnDef } from "src/util/table-util";
import { deleteDateBlock, fetchDateBlocks } from "src/api/v1/date-block";
import { getErrorMessage } from "src/util/errorHandler";
import { useTable } from "src/composables/useTable";
import DateBlockCreateDialog from "components/dashboard/DateBlockCreateDialog.vue";
import DateBlockEditDialog from "components/date-block/DateBlockEditDialog.vue";
import DateBlockHistoryDialog from "components/dashboard/DateBlockHistoryDialog.vue";

const $q = useQuasar();
const tableRef = ref();
const createDialogRef = ref();
const editDialogRef = ref();
const historyDialogOpen = ref(false);
const historyDateBlockId = ref(0);

const { pagination, loading, rows, onRequest } = useTable<DateBlock>({
  fetchFn: fetchDateBlocks,
  defaultPagination: {
    sortBy: "startDate",
    descending: true,
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
    ...getColumnDef("startDate"),
    align: "left",
    required: true,
    sortable: true,
  },
  {
    ...getColumnDef("endDate"),
    align: "left",
    required: true,
    sortable: true,
  },
  {
    ...getColumnDef("reason"),
    align: "left",
    required: true,
    sortable: true,
  },
  {
    ...getColumnDef("createdBy"),
    align: "left",
    required: true,
    sortable: true,
  },
  {
    ...getColumnDef("createdAt"),
    align: "left",
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
  return convertTableColumnDef(getDateBlockFieldDetail(field));
}

function reloadData() {
  tableRef.value.requestServerInteraction();
}

function openCreateDialog() {
  createDialogRef.value.open();
}

function openEditDialog(row: DateBlock) {
  editDialogRef.value.open(row);
}

function openHistoryDialog(row: DateBlock) {
  historyDateBlockId.value = row.id;
  historyDialogOpen.value = true;
}

function deleteItem(row: DateBlock) {
  const itemId = row.id;
  const itemName = `${row.startDate} ~ ${row.endDate}`;

  $q.dialog({
    title: "삭제",
    message: `정말로 '${itemName}' 기간의 예약을 삭제하시겠습니까?`,
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
    deleteDateBlock(itemId)
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
