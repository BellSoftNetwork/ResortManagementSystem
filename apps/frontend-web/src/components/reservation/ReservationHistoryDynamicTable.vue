<template>
  <q-table
    @request="onRequest"
    ref="tableRef"
    v-model:pagination="pagination"
    :loading="status.isLoading"
    :columns="columns"
    :rows="reservationHistories"
    :filter="filter"
    row-key="id"
    title="변경 이력"
    flat
    bordered
    binary-state-sort
  >
    <template #body-cell-historyType="props">
      <q-td key="historyType" :props="props">
        <q-chip
          :icon="REVISION_TYPE_MAP[props.row.historyType].icon"
          :color="REVISION_TYPE_MAP[props.row.historyType].color"
          outline
        >
          {{ REVISION_TYPE_MAP[props.row.historyType].name }}
        </q-chip>
      </q-td>
    </template>

    <template #body-cell-updatedValue="props">
      <q-td key="updatedValue" :props="props">
        <div class="row q-gutter-sm">
          <q-card v-for="field in props.row.updatedFields" v-bind:key="field" bordered>
            <q-card-section horizontal>
              <q-card-section class="bg-blue-3 q-pa-xs">
                {{ formatReservationFieldToLabel(field) }}
              </q-card-section>

              <q-card-section class="bg-grey-4 q-pa-xs">
                {{ formatReservationValue(field, props.row.entity[field]) }}
              </q-card-section>
            </q-card-section>
          </q-card>
        </div>
      </q-td>
    </template>

    <template #body-cell-modifier="props">
      <q-td key="modifier" :props="props">
        {{ getModifierName(props.row) }}
      </q-td>
    </template>
  </q-table>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue";
import { formatDateTime } from "src/util/format-util";
import { formatReservationFieldToLabel, formatReservationValue, Reservation } from "src/schema/reservation";
import { Revision, REVISION_TYPE_MAP } from "src/schema/revision";
import { formatSortParam } from "src/util/query-string-util";
import { fetchReservationHistories } from "src/api/v1/reservation";

const props = defineProps<{
  id: number;
}>();
const status = ref({
  isLoading: false,
  isLoaded: false,
  isPatching: false,
});
const tableRef = ref();
const filter = ref("");
const pagination = ref({
  sortBy: "updatedAt",
  descending: true,
  page: 1,
  rowsPerPage: 15,
  rowsNumber: 10,
});
const columns = [
  {
    name: "historyType",
    field: "historyType",
    label: "타입",
    align: "left",
    headerStyle: "width: 10%",
    required: true,
  },
  {
    name: "updatedValue",
    field: "updatedValue",
    label: "변경 내역",
    align: "left",
    required: true,
  },
  {
    name: "modifier",
    field: "modifier",
    label: "변경자",
    align: "left",
    headerStyle: "width: 10%",
    required: true,
  },
  {
    name: "updatedAt",
    field: "historyCreatedAt",
    label: "변경 시각",
    align: "left",
    headerStyle: "width: 15%",
    required: true,
    sortable: true,
    format: formatDateTime,
  },
];
const reservationHistories = ref<Revision<Reservation>[]>([]);

function getModifierName(row: Revision<Reservation>): string {
  if (row.historyType.includes("CREATE")) {
    return row.entity.createdBy?.name || "-";
  }
  return row.entity.updatedBy?.name || "-";
}

function onRequest(tableProps) {
  const { page, rowsPerPage, sortBy, descending } = tableProps.pagination;

  status.value.isLoading = true;
  status.value.isLoaded = false;

  fetchReservationHistories(props.id, {
    page: page - 1,
    size: rowsPerPage,
    sort: formatSortParam({ field: sortBy, isDescending: descending }),
  })
    .then((response) => {
      reservationHistories.value = response.values;
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

defineExpose({
  reloadData,
});

onMounted(() => {
  reloadData();
});
</script>
