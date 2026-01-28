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
    :title="props.title"
    flat
    bordered
    binary-state-sort
  >
    <template v-slot:top-right>
      <div class="row q-gutter-sm">
        <q-btn @click="filterDialog = true" color="primary" label="상세 검색" icon="search" />

        <ReservationFilterDialog
          v-model="filterDialog"
          :filter="filter"
          :due-options="dueOptions"
          :status-options="statusOptions"
          :default-stay-start-at="defaultConfig.filter.stayStartAt"
          @apply="handleFilterApply"
        />

        <q-btn :to="createPageLink()" icon="add" color="grey" dense round flat />
      </div>
    </template>

    <template #body-cell-rooms="props">
      <q-td key="rooms" :props="props">
        <div v-if="props.row.rooms.length !== 0">
          <span v-for="room in props.row.rooms" :key="room.id">
            <q-btn :to="{ name: 'Room', params: { id: room.id } }" align="left" color="primary" dense flat>
              {{ room.number }}
            </q-btn>
          </span>
        </div>
        <div v-else class="text-grey">미배정</div>
      </q-td>
    </template>

    <template #body-cell-name="props">
      <q-td key="name" :props="props">
        <q-btn :to="objectPageLink(props.row.id)" class="full-width" align="left" color="primary" dense flat
          >{{ props.row.name }}
        </q-btn>
      </q-td>
    </template>

    <template #body-cell-phone="props">
      <q-td :props="props">
        <a v-if="props.row.phone" :href="`tel:${props.row.phone}`" class="text-primary">
          <q-icon name="phone" size="xs" class="q-mr-xs" />
          {{ props.row.phone }}
        </a>
        <span v-else>-</span>
      </q-td>
    </template>

    <template #body-cell-note="props">
      <q-td :props="props" key="note">
        <q-btn
          v-if="props.row.note"
          @click="
            $q.dialog({
              title: `${props.row.name}님 예약 메모`,
              message: props.row.note,
            })
          "
          color="primary"
        >
          메모 확인
        </q-btn>
      </q-td>
    </template>

    <template #body-cell-actions="props">
      <q-td key="actions" :props="props">
        <q-btn dense round flat color="grey" icon="edit" :to="editPageLink(props.row.id)"></q-btn>
        <q-btn dense round flat color="grey" icon="delete" @click="deleteItem(props.row)"></q-btn>
      </q-td>
    </template>
  </q-table>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import dayjs from "dayjs";
import { useQuasar } from "quasar";
import { formatDate } from "src/util/format-util";
import { getReservationFieldDetail, Reservation, ReservationType } from "src/schema/reservation";
import { convertTableColumnDef } from "src/util/table-util";
import { deleteReservation, fetchReservations } from "src/api/v1/reservation";
import { getErrorMessage } from "src/util/errorHandler";
import { useRoute, useRouter } from "vue-router";
import { useTable } from "src/composables/useTable";
import ReservationFilterDialog from "./ReservationFilterDialog.vue";

const $q = useQuasar();
const props = withDefaults(
  defineProps<{
    title: string;
    reservationType: ReservationType;
  }>(),
  {
    title: "다가오는 예약",
  },
);
const tableRef = ref();
const route = useRoute();
const router = useRouter();

const defaultConfig = {
  pagination: {
    sortBy: "stayStartAt",
    descending: false,
    page: 1,
    rowsPerPage: 15,
  },
  filter: {
    peopleInfo: "",
    dueOption: "6M",
    stayStartAt: formatDate(),
    stayEndAt: dayjs().add(6, "M").format("YYYY-MM-DD"),
    status: "NORMAL",
  },
};
// Load filter from URL query string
const filter = ref({
  peopleInfo: route.query.peopleInfo?.toString() ?? defaultConfig.filter.peopleInfo,
  dueOption: defaultConfig.filter.dueOption,
  stayStartAt: route.query.stayStartAt?.toString() ?? defaultConfig.filter.stayStartAt,
  stayEndAt: route.query.stayEndAt?.toString() ?? defaultConfig.filter.stayEndAt,
  status: route.query.status?.toString().toUpperCase() ?? defaultConfig.filter.status,
});

// Create a computed filter for the API call
const apiFilter = computed(() => ({
  stayStartAt: filter.value.stayStartAt || undefined,
  stayEndAt: filter.value.stayEndAt || undefined,
  searchText: filter.value.peopleInfo || undefined,
  status: filter.value.status === "ALL" ? undefined : filter.value.status,
  type: props.reservationType,
}));
const statusOptions = [
  { label: "전체", value: "ALL" },
  { label: "예약 대기", value: "PENDING" },
  { label: "예약 확정", value: "NORMAL" },
  { label: "예약 취소", value: "CANCEL" },
  { label: "환불 완료", value: "REFUND" },
];
const dueOptions = [
  { label: "전체", value: "ALL" },
  { label: "1개월", value: "1M" },
  { label: "2개월", value: "2M" },
  { label: "3개월", value: "3M" },
  { label: "6개월", value: "6M" },
  { label: "직접 선택", value: "CUSTOM" },
];
// Use the useTable composable
const { pagination, loading, rows, onRequest } = useTable<Reservation>({
  fetchFn: fetchReservations,
  defaultPagination: {
    sortBy: "stayStartAt",
    descending: false,
    page: 1,
    rowsPerPage: 15,
  },
  filter: apiFilter,
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

const filterDialog = ref(false);
const columns = [
  {
    ...getColumnDef("stayStartAt"),
    align: "left",
    headerStyle: "width: 10%",
    required: true,
    sortable: true,
  },
  {
    ...getColumnDef("stayEndAt"),
    align: "left",
    headerStyle: "width: 10%",
    required: true,
    sortable: true,
  },
  {
    ...getColumnDef("name"),
    align: "left",
    headerStyle: "width: 10%",
    required: true,
    sortable: true,
  },
  {
    ...getColumnDef("phone"),
    align: "left",
    headerStyle: "width: 10%",
    required: true,
    sortable: true,
  },
  {
    ...getColumnDef("peopleCount"),
    align: "left",
    headerStyle: "width: 10%",
    required: true,
    sortable: true,
  },
  {
    ...getColumnDef("price"),
    align: "left",
    headerStyle: "width: 10%",
    required: true,
    sortable: true,
  },
  {
    ...getColumnDef("paymentAmount"),
    align: "left",
    headerStyle: "width: 10%",
    required: true,
    sortable: true,
  },
  {
    ...getColumnDef("rooms"),
    align: "left",
    headerStyle: "width: 10%",
    required: true,
    sortable: true,
  },
  {
    ...getColumnDef("paymentMethod"),
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
  ...(props.reservationType === "STAY"
    ? [
        {
          ...getColumnDef("note"),
          align: "left",
          headerStyle: "width: 10%",
        },
      ]
    : []),
  ...(props.reservationType === "MONTHLY_RENT"
    ? [
        {
          ...getColumnDef("deposit"),
          align: "left",
          headerStyle: "width: 15%",
          required: true,
          sortable: true,
        },
      ]
    : []),
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

function createPageLink() {
  if (props.reservationType === "MONTHLY_RENT") return { name: "CreateMonthlyRent" };

  return { name: "CreateReservation" };
}

function editPageLink(reservationId: number) {
  if (props.reservationType === "MONTHLY_RENT") return { name: "EditMonthlyRent", params: { id: reservationId } };

  return { name: "EditReservation", params: { id: reservationId } };
}

function objectPageLink(reservationId: number) {
  if (props.reservationType === "MONTHLY_RENT") return { name: "MonthlyRent", params: { id: reservationId } };

  return { name: "Reservation", params: { id: reservationId } };
}

function handleFilterApply(newFilter: typeof filter.value) {
  Object.assign(filter.value, newFilter);

  router.push({
    query: {
      ...route.query,
      peopleInfo: filter.value.peopleInfo !== defaultConfig.filter.peopleInfo ? filter.value.peopleInfo : undefined,
      stayStartAt: filter.value.stayStartAt !== defaultConfig.filter.stayStartAt ? filter.value.stayStartAt : undefined,
      stayEndAt: filter.value.stayEndAt !== defaultConfig.filter.stayEndAt ? filter.value.stayEndAt : undefined,
      status: filter.value.status !== defaultConfig.filter.status ? filter.value.status.toLowerCase() : undefined,
    },
  });

  filterDialog.value = false;
}

function getColumnDef(field: string) {
  return convertTableColumnDef(getReservationFieldDetail(field));
}

function reloadData() {
  tableRef.value.requestServerInteraction();
}

function deleteItem(row: Reservation) {
  const itemId = row.id;
  const itemName = row.name;

  $q.dialog({
    title: "삭제",
    message: `정말로 ${itemName}님의 예약 정보를 삭제하시겠습니까?`,
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
    deleteReservation(itemId)
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
