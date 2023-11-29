<template>
  <q-card>
    <q-inner-loading :showing="status.isProgress">
      <q-spinner-gears size="50px" color="primary" />
    </q-inner-loading>

    <q-card-section class="text-h6"> 객실 수정</q-card-section>

    <q-form @submit="update">
      <q-card-section>
        <q-input
          v-model="formData.number"
          :loading="status.isProgress"
          :disable="status.isProgress"
          :rules="roomStaticRules.number"
          label="번호"
          placeholder="101호"
          required
        ></q-input>

        <q-input
          v-model.number="formData.peekPrice"
          :loading="status.isProgress"
          :disable="status.isProgress"
          :rules="roomStaticRules.peekPrice"
          label="성수기 예약금"
          type="number"
          min="0"
          max="100000000"
          required
        ></q-input>

        <q-input
          v-model.number="formData.offPeekPrice"
          :loading="status.isProgress"
          :disable="status.isProgress"
          :rules="roomStaticRules.offPeekPrice"
          label="비성수기 예약금"
          type="number"
          min="0"
          max="100000000"
          required
        ></q-input>

        <q-input
          v-model="formData.description"
          :loading="status.isProgress"
          :disable="status.isProgress"
          :rules="roomStaticRules.description"
          type="textarea"
          label="설명"
          placeholder="와이파이 사용 가능"
        ></q-input>

        <q-input
          v-model="formData.note"
          :loading="status.isProgress"
          :disable="status.isProgress"
          :rules="roomStaticRules.note"
          type="textarea"
          label="메모 (관리용)"
          placeholder="문고리 고장"
        ></q-input>

        <q-select
          v-model="formData.status"
          :loading="status.isProgress"
          :disable="status.isProgress"
          :options="options.status"
          label="상태"
          required
          emit-value
          map-options
        ></q-select>
      </q-card-section>

      <q-card-actions align="right">
        <q-btn
          :disable="status.isProgress"
          :to="{ name: 'Rooms' }"
          color="primary"
          label="취소"
          flat
        />
        <q-btn
          :loading="status.isProgress"
          type="submit"
          color="red"
          label="수정"
          flat
        />
      </q-card-actions>
    </q-form>
  </q-card>
</template>

<script setup lang="ts">
import { onBeforeMount, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useQuasar } from "quasar";
import { Room, roomStaticRules } from "src/schema/room";
import { fetchRoom, patchRoom } from "src/api/v1/room";
import { getPatchedFormData, isFormValueChanged } from "src/util/data-util";

const router = useRouter();
const route = useRoute();
const $q = useQuasar();
const id = Number.parseInt(route.params.id as string);
const status = ref({
  isProgress: false,
});
const entity = ref<Room | null>(null);
const formData = ref<Partial<Room>>({
  number: "",
  peekPrice: 0,
  offPeekPrice: 0,
  description: "",
  note: "",
  status: "NORMAL",
});
const options = {
  status: [
    { label: "정상", value: "NORMAL" },
    { label: "이용불가", value: "INACTIVE" },
    { label: "파손", value: "DAMAGED" },
    { label: "공사 중", value: "CONSTRUCTION" },
  ],
};

function fetchData() {
  status.value.isProgress = true;

  return fetchRoom(id)
    .then((response) => {
      entity.value = response.value;
    })
    .catch((error) => {
      if (error.response.status === 404) router.push({ name: "ErrorNotFound" });

      console.log(error);
    })
    .finally(() => {
      status.value.isProgress = false;
    });
}

function update() {
  if (!isFormValueChanged(entity.value, formData.value)) {
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

  patchRoom(id, getPatchedFormData(entity.value, formData.value))
    .then(() => {
      router.push({ name: "Room", params: { id: id } });
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
  Object.assign(formData.value, entity.value);
}

onBeforeMount(() => {
  fetchData().then(() => {
    resetForm();
  });
});
</script>
