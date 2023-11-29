<template>
  <q-table
    @request="onRequest"
    ref="tableRef"
    v-model:pagination="pagination"
    :loading="status.isLoading"
    :columns="columns"
    :rows="reservationMethods"
    :filter="filter"
    style="height: 90vh"
    row-key="id"
    title="예약 수단"
    flat
    bordered
    binary-state-sort
  >
    <template v-slot:top-right>
      <div class="row q-gutter-sm">
        <ReservationMethodCreateDialog
          v-slot="{ dialog }"
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
            :rules="reservationMethodStaticRules.name"
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
        {{ formatCommissionRate(props.row.commissionRate) }}
        <q-popup-edit
          v-slot="scope"
          :model-value="props.row.commissionRate * 100"
          :persistent="status.isPatching"
        >
          <q-input
            v-model.number="scope.value"
            @keyup.enter="
              updateScope(
                props.row,
                scope,
                'commissionRate',
                (value) => value / 100,
              )
            "
            :loading="status.isPatching"
            :disable="status.isPatching"
            :rules="reservationMethodStaticRules.commissionRatePercent"
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

    <template #body-cell-requireUnpaidAmountCheck="props">
      <q-td key="requireUnpaidAmountCheck" :props="props">
        {{ props.row.requireUnpaidAmountCheck ? "활성" : "비활성" }}
        <q-popup-edit
          v-slot="scope"
          :model-value="props.row.requireUnpaidAmountCheck"
          :persistent="status.isPatching"
        >
          <q-checkbox
            v-model="scope.value"
            @update:model-value="
              updateScope(props.row, scope, 'requireUnpaidAmountCheck')
            "
            :loading="status.isPatching"
            :disable="status.isPatching"
            label="미수금 금액 알림"
          >
          </q-checkbox>
        </q-popup-edit>
      </q-td>
    </template>

    <template #body-cell-actions="props">
      <q-td key="actions" :props="props">
        <q-btn
          dense
          round
          flat
          color="grey"
          icon="delete"
          @click="deleteItem(props.row)"
        ></q-btn>
      </q-td>
    </template>
  </q-table>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useQuasar } from "quasar";
import ReservationMethodCreateDialog from "components/reservation-method/ReservationMethodCreateDialog.vue";
import {
  getReservationMethodFieldDetail,
  ReservationMethod,
  reservationMethodStaticRules,
} from "src/schema/reservation-method";
import { convertTableColumnDef } from "src/util/table-util";
import { formatCommissionRate } from "src/util/format-util";
import {
  deleteReservationMethod,
  fetchReservationMethods,
  patchReservationMethod,
} from "src/api/v1/reservation-method";
import { formatSortParam } from "src/util/query-string-util";

const $q = useQuasar();
const status = ref({
  isLoading: false,
  isLoaded: false,
  isPatching: false,
});
const tableRef = ref();
const inputRef = ref(null);
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
    required: true,
    sortable: true,
  },
  {
    ...getColumnDef("commissionRate"),
    align: "left",
    headerStyle: "width: 10%",
    required: true,
    sortable: true,
  },
  {
    ...getColumnDef("requireUnpaidAmountCheck"),
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
const reservationMethods = ref<ReservationMethod[]>();

function getColumnDef(field: string) {
  return convertTableColumnDef(getReservationMethodFieldDetail(field));
}

function onRequest(props) {
  const { page, rowsPerPage, sortBy, descending } = props.pagination;

  status.value.isLoading = true;
  status.value.isLoaded = false;

  fetchReservationMethods({
    page: page - 1,
    size: rowsPerPage,
    sort: formatSortParam({ field: sortBy, isDescending: descending }),
  })
    .then((response) => {
      reservationMethods.value = response.values;
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

function updateScope(row, scope, key, formatter) {
  if (
    (inputRef.value && !inputRef.value.validate()) ||
    row[key] === scope.value
  )
    return;

  const patchData = {};
  patchData[key] = formatter ? formatter(scope.value) : scope.value;

  status.value.isPatching = true;
  patchReservationMethod(row.id, patchData)
    .then((response) => {
      scope.set();
      row[key] = response.value[key];
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
    })
    .finally(() => {
      status.value.isPatching = false;
    });
}

function deleteItem(row: ReservationMethod) {
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
    deleteReservationMethod(itemId)
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
