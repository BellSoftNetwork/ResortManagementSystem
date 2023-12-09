<template>
  <q-table
    @request="onRequest"
    ref="tableRef"
    v-model:pagination="pagination"
    :loading="status.isLoading"
    :columns="columns"
    :rows="users"
    :filter="filter"
    style="height: 90vh"
    row-key="id"
    title="사용자 계정"
    flat
    bordered
    binary-state-sort
  >
    <template v-slot:top-right>
      <div class="row q-gutter-sm">
        <AccountCreateDialog v-slot="{ dialog }" @complete="reloadData">
          <q-btn @click="dialog.isOpen = true" icon="add" color="grey" dense round flat />
        </AccountCreateDialog>
      </div>
    </template>

    <template #body-cell-actions="props">
      <q-td key="actions" :props="props">
        <AccountEditDialog v-slot="{ dialog }" @complete="reloadData" :entity="props.row">
          <q-btn @click="dialog.isOpen = true" icon="edit" color="grey" dense round flat />
        </AccountEditDialog>
        <q-btn @click="deleteItem(props.row)" :disable="true" icon="delete" color="grey" dense round flat>
          <q-tooltip>구현 예정</q-tooltip>
        </q-btn>
      </q-td>
    </template>
  </q-table>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useQuasar } from "quasar";
import AccountCreateDialog from "components/account/AccountCreateDialog.vue";
import AccountEditDialog from "components/account/AccountEditDialog.vue";
import { getUserFieldDetail, User } from "src/schema/user";
import { convertTableColumnDef } from "src/util/table-util";
import { deleteAdminAccount, fetchAdminAccounts } from "src/api/v1/admin/account";
import { formatSortParam } from "src/util/query-string-util";

const $q = useQuasar();
const status = ref({
  isLoading: false,
  isLoaded: false,
  isPatching: false,
});
const tableRef = ref();
const filter = ref("");
const pagination = ref({
  sortBy: "name",
  descending: false,
  page: 1,
  rowsPerPage: 15,
  rowsNumber: 10,
});
const columns = [
  {
    ...getColumnDef("name"),
    align: "left",
    headerStyle: "width: 10%",
    required: true,
    sortable: true,
  },
  {
    ...getColumnDef("userId"),
    align: "left",
    headerStyle: "width: 10%",
    required: true,
    sortable: true,
  },
  {
    ...getColumnDef("email"),
    align: "left",
    required: true,
    sortable: true,
  },
  {
    ...getColumnDef("role"),
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
const users = ref<User[]>();

function getColumnDef(field: string) {
  return convertTableColumnDef(getUserFieldDetail(field));
}

function onRequest(props) {
  const { page, rowsPerPage, sortBy, descending } = props.pagination;

  status.value.isLoading = true;
  status.value.isLoaded = false;

  fetchAdminAccounts({
    page: page - 1,
    size: rowsPerPage,
    sort: formatSortParam({ field: sortBy, isDescending: descending }),
  })
    .then((response) => {
      users.value = response.values;
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

function deleteItem(row: User) {
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
    deleteAdminAccount(itemId)
      .then(() => {
        reloadData();
      })
      .catch((error) => {
        $q.notify({
          message: error.response.data.message,
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
