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
    title="예약 수단"
    flat
    bordered
    binary-state-sort
  >
    <template v-slot:top-right>
      <div class="row q-gutter-sm">
        <ReservationMethodCreateDialog
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
        </ReservationMethodCreateDialog>
      </div>
    </template>

    <template #body-cell-name="props">
      <q-td key="name" :props="props">
        {{ props.row.name }}
        <q-popup-edit
          v-slot="scope"
          :model-value="props.row.name"
          :persistent="status.isPatching"
        >
          <q-input
            v-model="scope.value"
            @keyup.enter="updateScope(props.row, scope, 'name')"
            :loading="status.isPatching"
            :disable="status.isPatching"
            :rules="rules.name"
            ref="inputRef"
            dense
            autofocus
            counter
          />
        </q-popup-edit>
      </q-td>
    </template>

    <template #body-cell-commissionRate="props">
      <q-td key="commissionRate" :props="props">
        {{ commissionRateFormatter(props.row.commissionRate) }}
        <q-popup-edit
          v-slot="scope"
          :model-value="props.row.commissionRate * 100"
          :persistent="status.isPatching"
        >
          <q-input
            v-model="scope.value"
            @keyup.enter="updateScope(props.row, scope, 'commissionRate', (value) => value / 100)"
            :loading="status.isPatching"
            :disable="status.isPatching"
            :rules="rules.commissionRatePercent"
            ref="inputRef"
            type="number"
            dense
            autofocus
          >
            <template v-slot:after>
              <q-icon name="percent" />
            </template>
          </q-input>
        </q-popup-edit>
      </q-td>
    </template>

    <template #body-cell-actions="props">
      <q-td key="actions" :props="props">
        <q-btn dense round flat color="grey" icon="delete" @click="deleteItem(props.row)"></q-btn>
      </q-td>
    </template>
  </q-table>
</template>

<script setup>
import { onMounted, ref } from "vue"
import dayjs from "dayjs"
import { useQuasar } from "quasar"
import ReservationMethodCreateDialog from "components/reservation-method/ReservationMethodCreateDialog.vue"
import { api } from "boot/axios"

const $q = useQuasar()
const status = ref({
  isLoading: false,
  isLoaded: false,
  isPatching: false,
})
const tableRef = ref()
const inputRef = ref(null)
const filter = ref("")
const pagination = ref({
  sortBy: "name",
  descending: false,
  page: 1,
  rowsPerPage: 15,
  rowsNumber: 10,
})
const rules = {
  name: [value => (value.length >= 2 && value.length <= 20) || "2~20 글자가 필요합니다"],
  commissionRatePercent: [value => (value >= 0 && value <= 100) || "수수료율이 유효하지 않습니다."],
}
const columns = [
  {
    name: "name",
    field: "name",
    label: "이름",
    align: "left",
    required: true,
    sortable: true,
  },
  {
    name: "commissionRate",
    field: "commissionRate",
    label: "수수료율",
    align: "left",
    headerStyle: "width: 10%",
    required: true,
    sortable: true,
    format: val => commissionRateFormatter(val),
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
      name: "네이버",
      commissionRate: 0.1,
      createdAt: "2021-01-01T00:00:00.000Z",
      updatedAt: "2021-01-01T00:00:00.000Z",
    },
  ],
})

function onRequest(props) {
  const { page, rowsPerPage, sortBy, descending } = props.pagination

  status.value.isLoading = true
  status.value.isLoaded = false

  api.get(`/api/v1/reservation-methods?size=${rowsPerPage}&page=${page - 1}&sort=${sortBy},${descending ? "desc" : "asc"}`)
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

function updateScope(row, scope, key, formatter) {
  if (!inputRef.value.validate() || row[key] === scope.value)
    return

  const patchData = {}
  patchData[key] = formatter ? formatter(scope.value) : scope.value

  status.value.isPatching = true
  api.patch(`/api/v1/reservation-methods/${row.id}`, patchData).then((response) => {
    scope.set()
    row[key] = response.data.value[key]
  }).catch((error) => {
    $q.notify({
      message: error.response.data.message,
      type: "negative",
      actions: [
        {
          icon: "close", color: "white", round: true,
        },
      ],
    })
  }).finally(() => {
    status.value.isPatching = false
  })
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
    api.delete(`/api/v1/reservation-methods/${itemId}`)
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

function commissionRateFormatter(value) {
  return value * 100 + "%"
}

function resetData() {
  responseData.value.values = []
}

onMounted(() => {
  resetData()
  reloadData()
})
</script>
