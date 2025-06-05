<template>
  <q-card>
    <q-card-section class="text-h6">객실 수정</q-card-section>

    <q-form @submit="update">
      <q-card-section>
        <q-select
          v-model="formData.roomGroup"
          :loading="status.isProgress || roomGroups === null"
          :disable="status.isProgress || roomGroups === null"
          :options="roomGroups"
          option-label="name"
          label="객실 그룹"
          required
          map-options
        ></q-select>

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
        <q-btn :disable="status.isProgress" :to="{ name: 'Rooms' }" color="primary" label="취소" flat />
        <q-btn :loading="status.isProgress" type="submit" color="red" label="수정" flat />
      </q-card-actions>
    </q-form>
  </q-card>
</template>

<script setup lang="ts">
import { onBeforeMount, ref } from "vue";
import { useRouter } from "vue-router";
import { useQuasar } from "quasar";
import { Room, roomStaticRules } from "src/schema/room";
import { patchRoom } from "src/api/v1/room";
import { getPatchedFormData, isFormValueChanged } from "src/util/data-util";
import { RoomGroup } from "src/schema/room-group";
import { fetchRoomGroups } from "src/api/v1/room-group";

const router = useRouter();
const $q = useQuasar();
const props = defineProps<{
  room: Room;
}>();
const status = ref({
  isProgress: false,
});
const roomGroups = ref<RoomGroup[] | null>(null);
const formData = ref<Partial<Room>>({
  number: "",
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

function loadRoomGroups() {
  roomGroups.value = null;

  return fetchRoomGroups().then((response) => {
    roomGroups.value = response.values;
  });
}

function update() {
  if (!isFormValueChanged(props.room, formData.value)) {
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

  patchRoom(props.room.id, getPatchedFormData(props.room, formData.value))
    .then(() => {
      router.push({ name: "Room", params: { id: props.room.id } });
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
  Object.assign(formData.value, props.room);
}

onBeforeMount(() => {
  resetForm();
  loadRoomGroups().then(() => {
    formData.value.roomGroup = roomGroups.value?.find((item) => item.id === props.room.roomGroup.id);
  });
});
</script>
