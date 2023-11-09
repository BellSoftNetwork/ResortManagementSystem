<template>
  <v-table height="80vh">
    <thead>
    <tr>
      <th scole="col" class="text-left" style="width: 20%">
        이름
      </th>
      <th scole="col" class="text-left" style="width: 35%">
        수수료율
      </th>
      <th scole="col" class="text-left" style="width: 20%">
        생성일
      </th>
      <th scole="col" class="text-left" style="width: 20%">
        수정일
      </th>
      <th scole="col" style="width: 5%"></th>
    </tr>
    </thead>
    <tbody>
    <tr
      v-if="status.isLoaded"
      v-for="reservationMethod in responseData.values"
      :key="reservationMethod.id"
    >
      <td>{{ reservationMethod.name }}</td>
      <td>{{ commissionRateFormatter(reservationMethod.commissionRate) }}</td>
      <td>{{ dayjs(reservationMethod.createdAt).format("YYYY-MM-DD HH:mm:ss") }}</td>
      <td>{{ dayjs(reservationMethod.updatedAt).format("YYYY-MM-DD HH:mm:ss") }}</td>
      <td>
        <ReservationMethodControlMenu @complete="fetchReservationMethods" :reservationMethod="reservationMethod" />
      </td>
    </tr>
    </tbody>
  </v-table>

  <v-pagination
    v-model="requestPage"
    :length="responseData.totalPages"
  ></v-pagination>

  <ReservationMethodCreateDialog @created="fetchReservationMethods" />
</template>

<script setup>
import { ref, watch } from "vue"
import axios from "@/modules/axios-wrapper"
import dayjs from "dayjs"

import ReservationMethodCreateDialog from "@/components/reservation-method/ReservationMethodCreateDialog.vue"
import ReservationMethodControlMenu from "@/components/reservation-method/ReservationMethodControlMenu.vue"

const status = ref({
  isLoading: false,
  isLoaded: false,
})
const requestPage = ref(1)
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

function fetchReservationMethods() {
  status.value.isLoading = true
  status.value.isLoaded = false
  responseData.value.values = []

  axios.get(`/api/v1/reservation-methods?size=14&page=${requestPage.value - 1}`).then(response => {
    responseData.value = response.data
    status.value.isLoaded = true
  }).finally(() => {
    status.value.isLoading = false
  })
}

function commissionRateFormatter(value) {
  return value * 100 + "%"
}

fetchReservationMethods()

watch(requestPage, fetchReservationMethods)
</script>
