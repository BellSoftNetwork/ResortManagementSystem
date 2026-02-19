<template>
  <q-table
    @request="onRequest"
    ref="tableRef"
    v-model:pagination="pagination"
    :loading="loading"
    :columns="columns"
    :rows="rows"
    :filter="filter"
    style="height: 90vh"
    row-key="id"
    title="감사 로그"
    flat
    bordered
    binary-state-sort
  >
    <template v-slot:top-right>
      <div class="row q-gutter-sm">
        <q-btn @click="filterDialog = true" color="primary" label="상세 검색" icon="search" />

        <q-dialog v-model="filterDialog" position="right" full-height>
          <q-card class="column full-height" style="width: 400px">
            <q-card-section class="row items-center q-pb-none">
              <div class="text-h6">상세 검색</div>
              <q-space />
              <q-btn icon="close" flat round dense v-close-popup />
            </q-card-section>

            <q-card-section class="col q-pa-md scroll">
              <div class="q-gutter-y-md">
                <!-- Date Range -->
                <div>
                  <div class="text-subtitle2 q-mb-sm">조회 기간</div>
                  <div class="row q-gutter-sm q-mb-sm">
                    <q-btn
                      v-for="opt in datePresets"
                      :key="opt.value"
                      :label="opt.label"
                      :color="filterBuffer.datePreset === opt.value ? 'primary' : 'grey-3'"
                      :text-color="filterBuffer.datePreset === opt.value ? 'white' : 'black'"
                      unelevated
                      size="sm"
                      @click="applyDatePreset(opt.value)"
                    />
                  </div>
                  <div class="row items-center q-gutter-x-sm">
                    <q-input
                      v-model="filterBuffer.startDate"
                      mask="####-##-##"
                      label="시작일"
                      outlined
                      dense
                      class="col"
                    >
                      <template v-slot:append>
                        <q-icon name="event" class="cursor-pointer">
                          <q-popup-proxy cover transition-show="scale" transition-hide="scale">
                            <q-date v-model="filterBuffer.startDate" mask="YYYY-MM-DD">
                              <div class="row items-center justify-end">
                                <q-btn v-close-popup label="Close" color="primary" flat />
                              </div>
                            </q-date>
                          </q-popup-proxy>
                        </q-icon>
                      </template>
                    </q-input>
                    <span>~</span>
                    <q-input v-model="filterBuffer.endDate" mask="####-##-##" label="종료일" outlined dense class="col">
                      <template v-slot:append>
                        <q-icon name="event" class="cursor-pointer">
                          <q-popup-proxy cover transition-show="scale" transition-hide="scale">
                            <q-date v-model="filterBuffer.endDate" mask="YYYY-MM-DD">
                              <div class="row items-center justify-end">
                                <q-btn v-close-popup label="Close" color="primary" flat />
                              </div>
                            </q-date>
                          </q-popup-proxy>
                        </q-icon>
                      </template>
                    </q-input>
                  </div>
                </div>

                <!-- Action Type -->
                <q-select
                  v-model="filterBuffer.action"
                  :options="actionOptions"
                  label="액션 타입"
                  outlined
                  emit-value
                  map-options
                  clearable
                />

                <!-- Entity Type -->
                <q-select
                  v-model="filterBuffer.entityType"
                  :options="entityTypeOptions"
                  label="대상 유형"
                  outlined
                  emit-value
                  map-options
                  clearable
                />

                <!-- Entity ID -->
                <q-input v-model="filterBuffer.entityId" label="대상 ID" outlined type="number" clearable />

                <!-- User ID -->
                <q-input v-model="filterBuffer.userId" label="변경자 ID" outlined type="number" clearable />
              </div>
            </q-card-section>

            <q-card-actions align="right" class="bg-grey-1">
              <q-btn flat label="초기화" color="primary" @click="resetFilter" />
              <q-btn label="검색" color="primary" @click="handleFilterApply" />
            </q-card-actions>
          </q-card>
        </q-dialog>
      </div>
    </template>

    <template v-slot:body="props">
      <q-tr :props="props">
        <!-- Expand Button -->
        <q-td auto-width>
          <q-btn
            size="sm"
            color="primary"
            round
            dense
            flat
            @click="toggleExpand(props.row)"
            :icon="expanded.has(props.row.id) ? 'keyboard_arrow_up' : 'keyboard_arrow_down'"
          />
        </q-td>

        <!-- Created At -->
        <q-td key="createdAt" :props="props">
          {{ formatDateTime(props.row.createdAt) }}
        </q-td>

        <!-- Action -->
        <q-td key="action" :props="props">
          <q-chip
            :color="REVISION_TYPE_MAP[props.row.action]?.color || 'grey'"
            text-color="white"
            size="sm"
            :icon="REVISION_TYPE_MAP[props.row.action]?.icon"
          >
            {{ REVISION_TYPE_MAP[props.row.action]?.name || props.row.action }}
          </q-chip>
        </q-td>

        <!-- Entity Type -->
        <q-td key="entityType" :props="props">
          {{ ENTITY_TYPE_LABELS[props.row.entityType] || props.row.entityType }}
        </q-td>

        <!-- Entity ID -->
        <q-td key="entityId" :props="props">
          <template v-if="props.row.action.includes('DELETE')">
            <span class="text-grey-7">{{ props.row.entityId }}</span>
            <q-badge color="grey" class="q-ml-xs">삭제됨</q-badge>
          </template>
          <q-btn
            v-else-if="props.row.entityType === 'reservation'"
            flat
            dense
            color="primary"
            :label="props.row.entityId"
            :to="{ name: 'Reservation', params: { id: props.row.entityId } }"
          />
          <q-btn
            v-else-if="props.row.entityType === 'room'"
            flat
            dense
            color="primary"
            :label="props.row.entityId"
            :to="{ name: 'Room', params: { id: props.row.entityId } }"
          />
          <span v-else>{{ props.row.entityId }}</span>
        </q-td>

        <!-- Username -->
        <q-td key="username" :props="props">
          {{ props.row.username }}
        </q-td>

        <!-- Changed Fields -->
        <q-td key="changedFields" :props="props">
          <div v-if="props.row.changedFields && props.row.changedFields.length > 0">
            <q-badge
              v-for="field in props.row.changedFields"
              :key="field"
              color="grey-2"
              text-color="black"
              class="q-mr-xs"
            >
              {{ field }}
            </q-badge>
          </div>
          <span v-else class="text-grey">-</span>
        </q-td>
      </q-tr>

      <!-- Expandable Detail Row -->
      <q-tr v-if="expanded.has(props.row.id)" :props="props">
        <q-td colspan="100%">
          <div class="q-pa-md bg-grey-1" style="max-height: 300px; overflow-y: auto">
            <div v-if="loadingDetail[props.row.id]" class="row justify-center q-pa-md">
              <q-spinner color="primary" size="2em" />
            </div>
            <div v-else-if="detailCache[props.row.id]">
              <div
                v-if="!detailCache[props.row.id].action.includes('CREATE') && detailCache[props.row.id].oldValues"
                class="q-mb-md"
              >
                <div class="text-subtitle2 text-grey-7 q-mb-xs">변경 전</div>
                <pre class="bg-white q-pa-sm rounded-borders shadow-1" style="white-space: pre-wrap">{{
                  JSON.stringify(detailCache[props.row.id].oldValues, null, 2)
                }}</pre>
              </div>
              <div v-if="!detailCache[props.row.id].action.includes('DELETE') && detailCache[props.row.id].newValues">
                <div class="text-subtitle2 text-grey-7 q-mb-xs">
                  {{ detailCache[props.row.id].action.includes("CREATE") ? "신규 생성" : "변경 후" }}
                </div>
                <pre class="bg-white q-pa-sm rounded-borders shadow-1" style="white-space: pre-wrap">{{
                  JSON.stringify(detailCache[props.row.id].newValues, null, 2)
                }}</pre>
              </div>
            </div>
            <div v-else class="text-center text-grey">상세 정보를 불러올 수 없습니다.</div>
          </div>
        </q-td>
      </q-tr>
    </template>

    <template #body-cell-action="props">
      <q-td :props="props">
        <q-chip
          :color="REVISION_TYPE_MAP[props.row.action]?.color || 'grey'"
          text-color="white"
          size="sm"
          :icon="REVISION_TYPE_MAP[props.row.action]?.icon"
        >
          {{ REVISION_TYPE_MAP[props.row.action]?.name || props.row.action }}
        </q-chip>
      </q-td>
    </template>

    <template #body-cell-changedFields="props">
      <q-td :props="props">
        <div v-if="props.row.changedFields && props.row.changedFields.length > 0">
          <q-badge
            v-for="field in props.row.changedFields"
            :key="field"
            color="grey-2"
            text-color="black"
            class="q-mr-xs"
          >
            {{ field }}
          </q-badge>
        </div>
        <span v-else class="text-grey">-</span>
      </q-td>
    </template>

    <template #body-cell-entityType="props">
      <q-td :props="props">
        {{ ENTITY_TYPE_LABELS[props.row.entityType] || props.row.entityType }}
      </q-td>
    </template>

    <template #body-cell-entityId="props">
      <q-td :props="props">
        <template v-if="props.row.action.includes('DELETE')">
          <span class="text-grey-7">{{ props.row.entityId }}</span>
          <q-badge color="grey" class="q-ml-xs">삭제됨</q-badge>
        </template>
        <q-btn
          v-else-if="props.row.entityType === 'reservation'"
          flat
          dense
          color="primary"
          :label="props.row.entityId"
          :to="{ name: 'Reservation', params: { id: props.row.entityId } }"
        />
        <q-btn
          v-else-if="props.row.entityType === 'room'"
          flat
          dense
          color="primary"
          :label="props.row.entityId"
          :to="{ name: 'Room', params: { id: props.row.entityId } }"
        />
        <span v-else>{{ props.row.entityId }}</span>
      </q-td>
    </template>
  </q-table>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import dayjs from "dayjs";
import { useQuasar } from "quasar";
import { useRoute, useRouter } from "vue-router";
import { useTable } from "src/composables/useTable";
import { formatDateTime, formatDate } from "src/util/format-util";
import { REVISION_TYPE_MAP, HistoryType } from "src/schema/revision";
import { AuditLog, AuditLogDetail, fetchAuditLogs, fetchAuditLogDetail } from "src/api/v1/audit";
import { getErrorMessage } from "src/util/errorHandler";

const $q = useQuasar();
const route = useRoute();
const router = useRouter();
const tableRef = ref();
const filterDialog = ref(false);
const expanded = ref(new Set<number>());
const detailCache = ref<Record<number, AuditLogDetail>>({});
const loadingDetail = ref<Record<number, boolean>>({});

async function toggleExpand(row: AuditLog) {
  if (expanded.value.has(row.id)) {
    expanded.value.delete(row.id);
  } else {
    expanded.value.add(row.id);
    if (!detailCache.value[row.id]) {
      loadingDetail.value[row.id] = true;
      try {
        const detail = await fetchAuditLogDetail(row.id);
        detailCache.value[row.id] = detail;
      } catch (e) {
        console.error(e);
        $q.notify({
          message: "상세 정보를 불러오는데 실패했습니다.",
          type: "negative",
        });
      } finally {
        loadingDetail.value[row.id] = false;
      }
    }
  }
}

const props = defineProps<{
  entityType?: string;
}>();

const datePresets = [
  { label: "오늘", value: "TODAY" },
  { label: "7일", value: "7D" },
  { label: "1개월", value: "1M" },
  { label: "3개월", value: "3M" },
];

const actionOptions = [
  { label: "전체", value: null },
  { label: "생성", value: "CREATE" },
  { label: "변경", value: "UPDATE" },
  { label: "삭제", value: "DELETE" },
];

const entityTypeOptions = [
  { label: "전체", value: null },
  { label: "예약", value: "reservation" },
  { label: "객실", value: "room" },
  { label: "객실 그룹", value: "room_group" },
  { label: "결제 수단", value: "payment_method" },
  { label: "계정", value: "user" },
];

const ENTITY_TYPE_LABELS: Record<string, string> = {
  reservation: "예약",
  room: "객실",
  room_group: "객실 그룹",
  payment_method: "결제 수단",
  user: "계정",
};

const defaultConfig = {
  pagination: {
    sortBy: "createdAt",
    descending: true,
    page: 1,
    rowsPerPage: 20,
  },
  filter: {
    startDate: dayjs().subtract(7, "day").format("YYYY-MM-DD"),
    endDate: formatDate(),
    action: null as HistoryType | null,
    entityType: null as string | null,
    entityId: null as number | null,
    userId: null as number | null,
    datePreset: "7D",
  },
};

// Load filter from URL or default
const filter = ref({
  startDate: route.query.startDate?.toString() ?? defaultConfig.filter.startDate,
  endDate: route.query.endDate?.toString() ?? defaultConfig.filter.endDate,
  action: (route.query.action as HistoryType) ?? defaultConfig.filter.action,
  entityType: route.query.entityType?.toString() ?? defaultConfig.filter.entityType,
  entityId: route.query.entityId ? Number(route.query.entityId) : defaultConfig.filter.entityId,
  userId: route.query.userId ? Number(route.query.userId) : defaultConfig.filter.userId,
  datePreset: route.query.datePreset?.toString() ?? defaultConfig.filter.datePreset,
});

const filterBuffer = ref({ ...filter.value });

// API Filter
const apiFilter = computed(() => ({
  entityType: filter.value.entityType || props.entityType || undefined,
  startDate: filter.value.startDate,
  endDate: filter.value.endDate,
  action: filter.value.action || undefined,
  entityId: filter.value.entityId || undefined,
  userId: filter.value.userId || undefined,
}));

const { pagination, loading, rows, onRequest } = useTable<AuditLog>({
  fetchFn: fetchAuditLogs,
  defaultPagination: defaultConfig.pagination,
  filter: apiFilter,
  onError: (error) => {
    $q.notify({
      message: getErrorMessage(error),
      type: "negative",
    });
  },
});

const columns = [
  { name: "expand", label: "", field: "id", align: "left", sortable: false },
  { name: "createdAt", label: "시간", field: "createdAt", align: "left", sortable: true },
  { name: "action", label: "액션", field: "action", align: "left", sortable: true },
  { name: "entityType", label: "대상 유형", field: "entityType", align: "left", sortable: true },
  { name: "entityId", label: "대상 ID", field: "entityId", align: "left", sortable: true },
  { name: "username", label: "변경자", field: "username", align: "left", sortable: true },
  { name: "changedFields", label: "변경 필드", field: "changedFields", align: "left", sortable: false },
];

function applyDatePreset(preset: string) {
  filterBuffer.value.datePreset = preset;
  const today = dayjs().format("YYYY-MM-DD");
  filterBuffer.value.endDate = today;

  if (preset === "TODAY") {
    filterBuffer.value.startDate = today;
  } else if (preset === "7D") {
    filterBuffer.value.startDate = dayjs().subtract(7, "day").format("YYYY-MM-DD");
  } else if (preset === "1M") {
    filterBuffer.value.startDate = dayjs().subtract(1, "month").format("YYYY-MM-DD");
  } else if (preset === "3M") {
    filterBuffer.value.startDate = dayjs().subtract(3, "month").format("YYYY-MM-DD");
  }
}

function handleFilterApply() {
  Object.assign(filter.value, filterBuffer.value);
  filterDialog.value = false;

  // Sync to URL
  const query: any = { ...route.query };
  if (filter.value.startDate !== defaultConfig.filter.startDate) query.startDate = filter.value.startDate;
  else delete query.startDate;

  if (filter.value.endDate !== defaultConfig.filter.endDate) query.endDate = filter.value.endDate;
  else delete query.endDate;

  if (filter.value.action) query.action = filter.value.action;
  else delete query.action;

  if (filter.value.entityType) query.entityType = filter.value.entityType;
  else delete query.entityType;

  if (filter.value.entityId) query.entityId = filter.value.entityId.toString();
  else delete query.entityId;

  if (filter.value.userId) query.userId = filter.value.userId.toString();
  else delete query.userId;

  if (filter.value.datePreset !== defaultConfig.filter.datePreset) query.datePreset = filter.value.datePreset;
  else delete query.datePreset;

  router.push({ query });

  // Trigger reload via useTable (it watches route if syncUrl is true, or we call onRequest)
  // useTable by default syncs URL, so modifying route.query might trigger it if useTable watches route.
  // But we are using internal filter ref passed to useTable.
  // useTable watches `filter`? No, it takes `filter` ref.
  // If `filter` ref changes, we usually need to call reload or it's reactive.
  // In ReservationListDynamicTable, they call requestServerInteraction.
  tableRef.value?.requestServerInteraction();
}

function resetFilter() {
  filterBuffer.value = { ...defaultConfig.filter };
}

onMounted(() => {
  tableRef.value?.requestServerInteraction();
});
</script>
