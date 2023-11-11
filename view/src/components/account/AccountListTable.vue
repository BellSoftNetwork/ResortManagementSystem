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
    title="사용자 계정"
    flat
    bordered
    binary-state-sort
  >
    <template v-slot:top-right>
      <div class="row q-gutter-sm">
        <AccountCreateDialog
          v-slot="{dialog}"
          @complete="reloadData"
        >
          <q-btn
            @click="dialog.isOpen = true"
            icon="add"
            color="grey"
            dense
            round
            flat
          />
        </AccountCreateDialog>
      </div>
    </template>

    <template #body-cell-actions="props">
      <q-td key="actions" :props="props">
        <AccountEditDialog
          v-slot="{dialog}"
          @complete="reloadData"
          :entity="props.row"
        >
          <q-btn
            @click="dialog.isOpen = true"
            icon="edit"
            color="grey"
            dense
            round
            flat
          />
        </AccountEditDialog>
        <q-btn
          @click="deleteItem(props.row)"
          :disable="true"
          icon="delete"
          color="grey"
          dense
          round
          flat
        >
          <q-tooltip>구현 예정</q-tooltip>
        </q-btn>
      </q-td>
    </template>
  </q-table>
</template>

<script setup>
import { onMounted, ref } from "vue"
import dayjs from "dayjs"
import { useQuasar } from "quasar"
import AccountCreateDialog from "components/account/AccountCreateDialog.vue"
import AccountEditDialog from "components/account/AccountEditDialog.vue"
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
  sortBy: "name",
  descending: false,
  page: 1,
  rowsPerPage: 15,
  rowsNumber: 10,
})
const columns = [
  {
    name: "name",
    field: "name",
    label: "이름",
    align: "left",
    headerStyle: "width: 10%",
    required: true,
    sortable: true,
  },
  {
    name: "email",
    field: "email",
    label: "이메일",
    align: "left",
    required: true,
    sortable: true,
  },
  {
    name: "role",
    field: "role",
    label: "권한",
    align: "left",
    headerStyle: "width: 10%",
    required: true,
    sortable: true,
    format: val => roleLabelConvert(val),
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
      name: "방울",
      email: "bell04204@gmail.com",
      role: "최고 관리자",
    },
  ],
})
const roleMap = {
  "NORMAL": "일반",
  "ADMIN": "관리자",
  "SUPER_ADMIN": "최고 관리자",
}

function onRequest(props) {
  const { page, rowsPerPage, sortBy, descending } = props.pagination

  status.value.isLoading = true
  status.value.isLoaded = false

  api.get(`/api/v1/admin/accounts?size=${rowsPerPage}&page=${page - 1}&sort=${sortBy},${descending ? "desc" : "asc"}`)
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
  const itemName = row.name

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
    api.delete(`/api/v1/admin/accounts/${itemId}`)
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

function roleLabelConvert(role) {
  return roleMap[role] || role
}

function resetData() {
  responseData.value.values = []
}

onMounted(() => {
  resetData()
  reloadData()
})
</script>
