<template>
  <v-table height="80vh">
    <thead>
    <tr>
      <th scole="col" class="text-left" style="width: 20%">
        이름
      </th>
      <th scole="col" class="text-left" style="width: 35%">
        이메일
      </th>
      <th scole="col" class="text-left" style="width: 20%">
        권한
      </th>
      <th scole="col" class="text-left" style="width: 20%">
        생성일
      </th>
      <th scole="col" style="width: 5%"></th>
    </tr>
    </thead>
    <tbody>
    <tr
      v-if="status.isLoaded"
      v-for="account in responseData.values"
      :key="account.id"
    >
      <td>{{ account.name }}</td>
      <td>{{ account.email }}</td>
      <td>{{ roleLabelConvert(account.role) }}</td>
      <td>{{ dayjs(account.createdAt).format("YYYY-MM-DD HH:mm:ss") }}</td>
      <td>
        <AccountControlMenu @complete="fetchAccounts" :account="account" />
      </td>
    </tr>
    </tbody>
  </v-table>

  <v-pagination
    v-model="requestPage"
    :length="responseData.totalPages"
  ></v-pagination>

  <AccountCreateDialog @created="fetchAccounts" />
</template>

<script setup>
import { ref, watch } from "vue"
import axios from "@/modules/axios-wrapper"
import dayjs from "dayjs"

import AccountCreateDialog from "@/components/account/AccountCreateDialog.vue"
import AccountControlMenu from "@/components/account/AccountControlMenu.vue"

const status = ref({
  isLoading: false,
  isLoaded: false,
})
const requestPage = ref(1)
const responseData = ref({
  page: 0,
  totalPages: 1,
  totalElements: 0,
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

function fetchAccounts() {
  status.value.isLoading = true
  status.value.isLoaded = false
  responseData.value.values = []

  axios.get(`/api/v1/admin/accounts?size=14&page=${requestPage.value - 1}`).then(response => {
    responseData.value = response.data
    status.value.isLoaded = true
  }).finally(() => {
    status.value.isLoading = false
  })
}

function roleLabelConvert(role) {
  return roleMap[role] || role
}

fetchAccounts()

watch(requestPage, fetchAccounts)
</script>
