<template>
  <q-dialog :model-value="modelValue" @update:model-value="emit('update:modelValue', $event)" maximized>
    <q-card>
      <q-card-section class="row items-center q-pb-none">
        <div class="text-h6">예약 마감 변경 이력</div>
        <q-space />
        <q-btn icon="close" flat round dense v-close-popup />
      </q-card-section>
      <q-card-section>
        <q-table
          @request="onRequest"
          ref="tableRef"
          v-model:pagination="pagination"
          :loading="status.isLoading"
          :columns="columns"
          :rows="histories"
          row-key="historyCreatedAt"
          flat
          bordered
          binary-state-sort
        >
          <template #body-cell-historyType="props">
            <q-td key="historyType" :props="props">
              <q-chip
                :icon="REVISION_TYPE_MAP[props.row.historyType]?.icon"
                :color="REVISION_TYPE_MAP[props.row.historyType]?.color"
                outline
              >
                {{ REVISION_TYPE_MAP[props.row.historyType]?.name }}
              </q-chip>
            </q-td>
          </template>

          <template #body-cell-updatedValue="props">
            <q-td key="updatedValue" :props="props">
              <div class="row q-gutter-sm">
                <template v-if="props.row.updatedFields && props.row.updatedFields.length > 0">
                  <q-card v-for="field in props.row.updatedFields" :key="field" bordered flat>
                    <q-card-section horizontal>
                      <q-card-section class="bg-blue-1 q-pa-xs text-caption">
                        {{ formatDateBlockFieldToLabel(field) }}
                      </q-card-section>

                      <q-separator vertical />

                      <q-card-section class="q-pa-xs text-caption">
                        {{ formatDateBlockValue(field, props.row.entity[field]) }}
                      </q-card-section>
                    </q-card-section>
                  </q-card>
                </template>
                <template v-else>
                  <span class="text-grey-6">변경 내역 없음 (전체 스냅샷)</span>
                </template>
              </div>
            </q-td>
          </template>
        </q-table>
      </q-card-section>
    </q-card>
  </q-dialog>
</template>

<script setup lang="ts">
import { ref, watch } from "vue";
import { formatDateTime } from "src/util/format-util";
import { formatDateBlockFieldToLabel, formatDateBlockValue, DateBlock } from "src/schema/date-block";
import { Revision, REVISION_TYPE_MAP } from "src/schema/revision";
import { fetchDateBlockHistories } from "src/api/v1/date-block";
import { formatSortParam } from "src/util/query-string-util";

const props = defineProps<{
  modelValue: boolean;
  dateBlockId: number;
}>();

const emit = defineEmits<{
  (e: "update:modelValue", value: boolean): void;
}>();

const status = ref({
  isLoading: false,
  isLoaded: false,
});

const tableRef = ref();
const pagination = ref({
  sortBy: "historyCreatedAt",
  descending: true,
  page: 1,
  rowsPerPage: 15,
  rowsNumber: 0,
});

const columns = [
  {
    name: "historyType",
    field: "historyType",
    label: "타입",
    align: "left",
    headerStyle: "width: 100px",
  },
  {
    name: "updatedValue",
    field: "updatedValue",
    label: "변경 내역",
    align: "left",
  },
  {
    name: "historyUsername",
    field: "historyUsername",
    label: "변경자",
    align: "left",
    headerStyle: "width: 140px",
    format: (value: string | undefined) => value || "-",
  },
  {
    name: "updatedAt",
    field: "historyCreatedAt",
    label: "변경 시각",
    align: "left",
    headerStyle: "width: 180px",
    sortable: true,
    format: formatDateTime,
  },
];

type DateBlockHistoryRow = Revision<DateBlock> & {
  historyUsername?: string;
};

const histories = ref<DateBlockHistoryRow[]>([]);

async function onRequest(tableProps: any) {
  const { page, rowsPerPage, sortBy, descending } = tableProps.pagination;

  status.value.isLoading = true;

  try {
    const response = await fetchDateBlockHistories(props.dateBlockId, {
      page: page - 1,
      size: rowsPerPage,
      sort: formatSortParam({ field: sortBy, isDescending: descending }),
    });

    histories.value = response.values;
    const pageInfo = response.page;

    pagination.value.rowsNumber = pageInfo.totalElements;
    pagination.value.page = pageInfo.index + 1;
    pagination.value.rowsPerPage = pageInfo.size;
    pagination.value.sortBy = sortBy;
    pagination.value.descending = descending;

    status.value.isLoaded = true;
  } catch (error) {
    console.error("Failed to fetch history", error);
  } finally {
    status.value.isLoading = false;
  }
}

function reloadData() {
  // Reset pagination to first page on new open if desired, or just refresh current
  if (tableRef.value) {
    tableRef.value.requestServerInteraction();
  } else {
    // If table not mounted yet, it will trigger onRequest on mount if we don't interfere,
    // but since it's in a dialog, v-if might be involved or just keep-alive.
    // simpler: when dialog opens, we want to load.
  }
}

watch(
  () => props.modelValue,
  (newVal) => {
    if (newVal && props.dateBlockId) {
      // Small delay to ensure table is mounted in dialog
      setTimeout(() => {
        reloadData();
      }, 100);
    }
  },
);
</script>
