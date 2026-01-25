<template>
  <slot :dialog="dialog">
    <q-btn @click="dialog.isOpen = true">추가</q-btn>
  </slot>

  <q-dialog v-model="dialog.isOpen" :persistent="status.isProgress" @beforeShow="resetForm">
    <q-card style="width: 500px">
      <q-card-section class="text-h6">계정 수정</q-card-section>

      <q-form @submit="update">
        <q-card-section>
          <q-input v-model="formData.userId" :rules="userStaticRules.userId" label="계정 ID" required></q-input>

          <q-input v-model="formData.email" :rules="userStaticRules.email" label="이메일"></q-input>

          <q-input v-model="formData.name" :rules="userStaticRules.name" label="이름" required></q-input>

          <q-input
            v-model="formData.password"
            :rules="userStaticRules.password"
            type="password"
            label="비밀번호"
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
          <q-btn :loading="status.isProgress" type="submit" color="red" label="수정" flat />
        </q-card-actions>
      </q-form>
    </q-card>
  </q-dialog>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { useAuthStore } from "stores/auth";
import { useQuasar } from "quasar";
import { User, userStaticRules } from "src/schema/user";
import { AdminAccountPatchParams, patchAdminAccount } from "src/api/v1/admin/account";
import { getErrorMessage } from "src/util/errorHandler";

const authStore = useAuthStore();
const $q = useQuasar();

const emit = defineEmits(["complete"]);
const props = defineProps<{
  user: User;
}>();
const dialog = ref({
  isOpen: false,
});
const status = ref({
  isProgress: false,
});
const formData = ref({
  userId: props.user.userId,
  email: props.user.email,
  name: props.user.name,
  password: "",
  role: props.user.role,
});
const options = {
  role: [
    { label: "일반", value: "NORMAL" },
    { label: "관리자", value: "ADMIN", disable: !authStore.isSuperAdminRole },
    { label: "최고 관리자", value: "SUPER_ADMIN", disable: !authStore.isSuperAdminRole },
  ],
};

function update() {
  if (!isChanged()) {
    $q.notify({
      message: "수정된 항목이 없습니다.",
      type: "info",
      actions: [
        {
          icon: "close",
          color: "white",
          round: true,
        },
      ],
    });

    return;
  }

  status.value.isProgress = true;

  patchAdminAccount(props.user.id, patchedData())
    .then(() => {
      emit("complete");
      dialog.value.isOpen = false;
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
    })
    .finally(() => {
      status.value.isProgress = false;
    });
}

function isChanged() {
  return (
    formData.value.userId !== props.user.userId ||
    formData.value.email !== props.user.email ||
    formData.value.name !== props.user.name ||
    formData.value.role !== props.user.role ||
    formData.value.password.length > 0
  );
}

function patchedData() {
  const patchData: AdminAccountPatchParams = {};

  if (props.user.userId !== formData.value.userId) patchData.userId = formData.value.userId;
  if (props.user.email !== formData.value.email) {
    const email = formData.value.email ? formData.value.email : null;
    patchData.email = email;
  }
  if (props.user.name !== formData.value.name) patchData.name = formData.value.name;
  if (props.user.role !== formData.value.role) patchData.role = formData.value.role;
  if (formData.value.password.length > 0) patchData.password = formData.value.password;

  return patchData;
}

function resetForm() {
  formData.value.userId = props.user.userId;
  formData.value.email = props.user.email;
  formData.value.name = props.user.name;
  formData.value.password = "";
  formData.value.role = props.user.role;
}
</script>
