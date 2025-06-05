<template>
  <slot :dialog="dialog">
    <q-btn @click="dialog.isOpen = true">추가</q-btn>
  </slot>

  <q-dialog v-model="dialog.isOpen" :persistent="status.isProgress" @beforeShow="resetForm">
    <q-card style="width: 500px">
      <q-card-section class="text-h6">계정 추가</q-card-section>

      <q-form @submit="create">
        <q-card-section>
          <q-input v-model="formData.name" :rules="userStaticRules.name" label="이름" required></q-input>

          <q-input v-model="formData.userId" :rules="userStaticRules.userId" label="계정 ID" required></q-input>

          <q-input v-model="formData.email" :rules="userStaticRules.email" label="이메일"></q-input>

          <q-input
            v-model="formData.password"
            :rules="userStaticRules.password"
            type="password"
            label="비밀번호"
            required
          ></q-input>

          <q-select
            v-model="formData.role"
            :options="options.role"
            label="권한"
            required
            emit-value
            map-options
          ></q-select>
        </q-card-section>

        <q-card-actions align="right">
          <q-btn v-close-popup :disable="status.isProgress" color="primary" label="취소" flat />
          <q-btn :loading="status.isProgress" type="submit" color="red" label="추가" flat />
        </q-card-actions>
      </q-form>
    </q-card>
  </q-dialog>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { useAuthStore } from "stores/auth";
import { useQuasar } from "quasar";
import { userStaticRules } from "src/schema/user";
import { createAdminAccount } from "src/api/v1/admin/account";

const authStore = useAuthStore();
const $q = useQuasar();

const emit = defineEmits(["complete"]);
const dialog = ref({
  isOpen: false,
});
const status = ref({
  isProgress: false,
});
const formData = ref({
  name: "",
  userId: "",
  email: "",
  password: "",
  role: "NORMAL",
});
const options = {
  role: [
    { label: "일반", value: "NORMAL" },
    { label: "관리자", value: "ADMIN", disable: !authStore.isSuperAdminRole },
    { label: "최고 관리자", value: "SUPER_ADMIN", disable: !authStore.isSuperAdminRole },
  ],
};

function create() {
  status.value.isProgress = true;

  createAdminAccount(formData.value)
    .then(() => {
      emit("complete");
      dialog.value.isOpen = false;

      resetForm();
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
      status.value.isProgress = false;
    });
}

function resetForm() {
  formData.value.name = "";
  formData.value.userId = "";
  formData.value.email = "";
  formData.value.password = "";
  formData.value.role = "NORMAL";
}
</script>
