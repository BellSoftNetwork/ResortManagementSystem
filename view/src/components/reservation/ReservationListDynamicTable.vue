<template>
  <q-table
    @request="onRequest"
    ref="tableRef"
    v-model:pagination="pagination"
    :loading="status.isLoading"
    :columns="columns"
    :rows="reservations"
    :filter="filter"
    style="height: 90vh"
    row-key="id"
    title="다가오는 예약"
    flat
    bordered
    binary-state-sort
  >
    <template v-slot:top-right>
      <div class="row q-gutter-sm">
        <q-btn color="primary" label="상세 검색" icon="search">
          <q-menu anchor="bottom end" self="top end" @hide="setFilterQuery">
            <div class="row no-wrap q-pa-md">
              <q-input
                v-model="filter.peopleInfo"
                debounce="300"
                placeholder="홍길동"
                label="예약자 정보"
                class="fit"
              />
            </div>

            <div class="row no-wrap q-pa-md">
              <q-input v-model="filter.stayStartAt" mask="####-##-##" :readonly="true" outlined>
                <template v-slot:append>
                  <q-icon name="event" class="cursor-pointer">
                    <q-popup-proxy cover transition-show="scale" transition-hide="scale">
                      <q-date v-model="filter.stayStartAt" mask="YYYY-MM-DD">
                        <div class="row items-center justify-end">
                          <q-btn v-close-popup label="Close" color="primary" flat />
                        </div>
                      </q-date>
                    </q-popup-proxy>
                  </q-icon>
                </template>
              </q-input>
              <span class="self-center q-mx-sm">~</span>
              <q-input v-model="filter.stayEndAt" mask="####-##-##" :readonly="true" outlined>
                <template v-slot:append>
                  <q-icon name="event" class="cursor-pointer">
                    <q-popup-proxy cover transition-show="scale" transition-hide="scale">
                      <q-date v-model="filter.stayEndAt" mask="YYYY-MM-DD">
                        <div class="row items-center justify-end">
                          <q-btn v-close-popup label="Close" color="primary" flat />
                        </div>
                      </q-date>
                    </q-popup-proxy>
                  </q-icon>
                </template>
              </q-input>
            </div>
          </q-menu>
        </q-btn>

        <q-btn :to="{ name: 'CreateReservation' }" icon="add" color="grey" dense round flat />
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
        <q-btn
          :to="{ name: 'Reservation', params: { id: props.row.id } }"
          class="full-width"
          align="left"
          color="primary"
          dense
          flat
          >{{ props.row.name }}
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
          :to="{ name: 'EditReservation', params: { id: props.row.id } }"
        ></q-btn>
        <q-btn dense round flat color="grey" icon="delete" @click="deleteItem(props.row)"></q-btn>
      </q-td>
    </template>
  </q-table>
</template>

<script setup lang="ts">
import { onMounted, ref, watch } from "vue";
import dayjs from "dayjs";
import { useQuasar } from "quasar";
import { formatDate } from "src/util/format-util";
import { getReservationFieldDetail, Reservation } from "src/schema/reservation";
import { convertTableColumnDef } from "src/util/table-util";
import { deleteReservation, fetchReservations } from "src/api/v1/reservation";
import { formatSortParam } from "src/util/query-string-util";
import { useRoute, useRouter } from "vue-router";

const $q = useQuasar();
const status = ref({
  isLoading: false,
});
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
    stayStartAt: formatDate(),
    stayEndAt: dayjs().add(3, "M").format("YYYY-MM-DD"),
  },
};
const filter = ref({
  peopleInfo: defaultConfig.filter.peopleInfo,
  stayStartAt: defaultConfig.filter.stayStartAt,
  stayEndAt: defaultConfig.filter.stayEndAt,
});
const pagination = ref({
  sortBy: defaultConfig.pagination.sortBy,
  descending: defaultConfig.pagination.descending,
  page: defaultConfig.pagination.page,
  rowsPerPage: defaultConfig.pagination.rowsPerPage,
  rowsNumber: 0,
});
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
  {
    ...getColumnDef("checkInAt"),
    align: "left",
    headerStyle: "width: 15%",
    required: true,
    sortable: true,
  },
  {
    ...getColumnDef("checkOutAt"),
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
const reservations = ref<Reservation[]>();

loadQueryString();

watch(route, () => {
  loadQueryString();
});

function loadQueryString() {
  filter.value.peopleInfo = route.query.peopleInfo?.toString() ?? defaultConfig.filter.peopleInfo;
  filter.value.stayStartAt = route.query.stayStartAt?.toString() ?? defaultConfig.filter.stayStartAt;
  filter.value.stayEndAt = route.query.stayEndAt?.toString() ?? defaultConfig.filter.stayEndAt;

  pagination.value.sortBy = route.query.sortBy?.toString() ?? defaultConfig.pagination.sortBy;
  pagination.value.descending = Boolean(route.query.descending ?? defaultConfig.pagination.descending);
  pagination.value.page = Number(route.query.page ?? defaultConfig.pagination.page);
  pagination.value.rowsPerPage = Number(route.query.rowsPerPage ?? defaultConfig.pagination.rowsPerPage);
}

function setFilterQuery() {
  router.push({
    query: {
      ...route.query,
      peopleInfo: filter.value.peopleInfo !== defaultConfig.filter.peopleInfo ? filter.value.peopleInfo : undefined,
      stayStartAt: filter.value.stayStartAt !== defaultConfig.filter.stayStartAt ? filter.value.stayStartAt : undefined,
      stayEndAt: filter.value.stayEndAt !== defaultConfig.filter.stayEndAt ? filter.value.stayEndAt : undefined,
    },
  });
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
  return convertTableColumnDef(getReservationFieldDetail(field));
}

function onRequest(props) {
  const { page, rowsPerPage, sortBy, descending } = props.pagination;

  status.value.isLoading = true;

  fetchReservations({
    page: page - 1,
    size: rowsPerPage,
    sort: formatSortParam({ field: sortBy, isDescending: descending }),
    stayStartAt: filter.value.stayStartAt,
    stayEndAt: filter.value.stayEndAt,
    searchText: filter.value.peopleInfo || undefined,
  })
    .then((response) => {
      reservations.value = response.values;
      const pageInfo = response.page;

      pagination.value.rowsNumber = pageInfo.totalElements;
      pagination.value.page = pageInfo.index + 1;
      pagination.value.rowsPerPage = pageInfo.size;
      pagination.value.sortBy = sortBy;
      pagination.value.descending = descending;

      setPaginationQuery();
    })
    .catch((error) => {
      reservations.value = [];

      console.error(error);
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
      status.value.isLoading = false;
    });
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
