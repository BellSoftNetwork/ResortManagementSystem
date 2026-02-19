<template>
  <q-dialog v-model="dialog.isOpen" :persistent="status.isProgress" @before-show="resetForm">
    <q-card style="width: 500px">
      <q-card-section class="text-h6">예약 마감</q-card-section>

      <q-form @submit="create">
        <q-card-section>
          <q-input
            v-model="formData.startDate"
            :rules="[(v) => !!v || '시작일을 입력해주세요']"
            type="date"
            label="시작일"
            required
          />
          <q-input
            v-model="formData.endDate"
            :rules="[
              (v) => !!v || '종료일을 입력해주세요',
              (v) => !v || v >= formData.startDate || '종료일은 시작일 이후여야 합니다',
            ]"
            type="date"
            label="종료일"
            required
          />
          <q-input
            v-model="formData.reason"
            :rules="[(v) => !!v || '사유를 입력해주세요', (v) => v.length <= 200 || '200자 이내로 입력해주세요']"
            label="사유"
            maxlength="200"
            required
          />
        </q-card-section>

        <q-card-actions align="right">
          <q-btn v-close-popup :disable="status.isProgress" color="primary" label="취소" flat />
          <q-btn :loading="status.isProgress" type="submit" color="red" label="마감" flat />
        </q-card-actions>
      </q-form>
    </q-card>
  </q-dialog>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { useQuasar } from "quasar";
import { createDateBlock } from "src/api/v1/date-block";
import { getErrorMessage } from "src/util/errorHandler";

const $q = useQuasar();
const emit = defineEmits(["created"]);

const dialog = ref({ isOpen: false });
const status = ref({ isProgress: false });
const formData = ref({ startDate: "", endDate: "", reason: "" });

function open() {
  dialog.value.isOpen = true;
}

function resetForm() {
  formData.value.startDate = "";
  formData.value.endDate = "";
  formData.value.reason = "";
}

function create() {
  status.value.isProgress = true;

  createDateBlock(formData.value)
    .then(() => {
      emit("created");
      dialog.value.isOpen = false;
      resetForm();
    })
    .catch((error) => {
      $q.notify({
        message: getErrorMessage(error),
        type: "negative",
        actions: [{ icon: "close", color: "white", round: true }],
      });
    })
    .finally(() => {
      status.value.isProgress = false;
    });
}

defineExpose({ open });
</script>
