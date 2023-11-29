<template>
  <q-card flat bordered>
    <q-inner-loading :showing="status.isProgress">
      <q-spinner-gears size="50px" color="primary" />
    </q-inner-loading>

    <q-card-section class="text-h6">
      {{ entity?.number }}
    </q-card-section>

    <q-card-section>
      <q-input
        :model-value="entity?.peekPrice"
        :loading="status.isProgress"
        :readonly="true"
        label="성수기 예약금"
        type="number"
      ></q-input>

      <q-input
        :model-value="entity?.offPeekPrice"
        :loading="status.isProgress"
        :readonly="true"
        label="비성수기 예약금"
        type="number"
      ></q-input>

      <q-input
        :model-value="entity?.description"
        :loading="status.isProgress"
        :readonly="true"
        type="textarea"
        label="설명"
      ></q-input>

      <q-input
        :model-value="entity?.note"
        :loading="status.isProgress"
        :readonly="true"
        type="textarea"
        label="메모 (관리용)"
      ></q-input>

      <q-select
        :model-value="entity?.status"
        :loading="status.isProgress"
        :readonly="true"
        :options="options.status"
        label="상태"
        emit-value
        map-options
      ></q-select>
    </q-card-section>

    <q-card-actions align="right">
      <q-btn @click="deleteItem()" color="red" label="삭제" dense flat></q-btn>
      <q-btn
        :disable="status.isProgress"
        :to="{ name: 'EditRoom', params: { id: entity?.id } }"
        color="primary"
        label="수정"
        flat
      />
    </q-card-actions>
  </q-card>
</template>

<script setup lang="ts">
import { onBeforeMount, ref } from "vue";
import { useRouter } from "vue-router";
import { useQuasar } from "quasar";
import { deleteRoom, fetchRoom } from "src/api/v1/room";
import { Room } from "src/schema/room";

const router = useRouter();
const $q = useQuasar();
const props = defineProps<{
  id: number;
}>();
const id = props.id;
const status = ref({
  isProgress: false,
});
const entity = ref<Room | null>(null);
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

function deleteItem() {
  if (entity.value === null) return;

  const itemId = entity.value.id;
  const itemName = entity.value.number;

  $q.dialog({
    title: "삭제",
    message: `정말로 '${itemName}'을 삭제하시겠습니까?`,
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
    deleteRoom(itemId)
      .then(() => {
        router.push({ name: "Rooms" });
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

onBeforeMount(() => {
  fetchData();
});
</script>
