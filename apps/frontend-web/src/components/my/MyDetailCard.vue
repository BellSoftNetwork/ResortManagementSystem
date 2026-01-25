<template>
  <q-card>
    <q-card-section class="text-h6">내 정보</q-card-section>

    <q-form @submit="update">
      <q-card-section>
        <q-input :model-value="props.user.name" label="이름" :readonly="true"></q-input>

        <q-input :model-value="props.user.userId" :readonly="true" label="계정 ID"></q-input>

        <q-input v-model="formData.email" :rules="userStaticRules.email" label="이메일" required></q-input>

        <q-input
          v-model="formData.password"
          :rules="userStaticRules.password"
          type="password"
          label="비밀번호"
        ></q-input>

        <q-input
          v-model="formData.passwordConfirm"
          :rules="userDynamicRules.passwordConfirm(formData.password)"
          type="password"
          label="비밀번호 확인"
        ></q-input>
      </q-card-section>

      <q-card-actions align="right">
        <q-btn :loading="status.isProgress" type="submit" color="red" label="수정" flat />
      </q-card-actions>
    </q-form>
  </q-card>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { useQuasar } from "quasar";
import { User, userDynamicRules, userStaticRules } from "src/schema/user";
import { MyPatchParams, patchMy } from "src/api/v1/main";
import { getErrorMessage } from "src/util/errorHandler";

const $q = useQuasar();

const props = defineProps<{
  user: User;
}>();
const status = ref({
  isProgress: false,
});
const formData = ref({
  email: props.user.email,
  password: "",
  passwordConfirm: "",
});

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

  patchMy(patchedData())
    .then(() => {
      $q.notify({
        message: "정보가 정상적으로 변경되었습니다",
        type: "positive",
        actions: [
          {
            icon: "close",
            color: "white",
            round: true,
          },
        ],
      });
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
  return formData.value.email !== props.user.email || formData.value.password.length > 0;
}

function patchedData() {
  const patchData: MyPatchParams = {};

  if (formData.value.email !== props.user.email) patchData.email = formData.value.email;

  if (formData.value.password.length > 0) patchData.password = formData.value.password;

  return patchData;
}
</script>
