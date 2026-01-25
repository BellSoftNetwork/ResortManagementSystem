<template>
  <q-card>
    <q-card-section class="text-h6">객실 추가</q-card-section>

    <q-form @submit="create">
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
          :rules="roomStaticRules.number"
          label="번호"
          placeholder="101호"
          required
        ></q-input>

        <q-input
          v-model="formData.note"
          :rules="roomStaticRules.note"
          type="textarea"
          label="메모 (관리용)"
          placeholder="문고리 고장"
        ></q-input>

        <q-select
          v-model="formData.status"
          :options="options.status"
          label="상태"
          required
          emit-value
          map-options
        ></q-select>
      </q-card-section>

      <q-card-actions align="right">
        <q-btn :disable="status.isProgress" :to="{ name: 'Rooms' }" color="primary" label="취소" flat />
        <q-btn :loading="status.isProgress" type="submit" color="red" label="추가" flat />
      </q-card-actions>
    </q-form>
  </q-card>
</template>

<script setup lang="ts">
import { onBeforeMount, ref } from "vue";
import { useRouter } from "vue-router";
import { useQuasar } from "quasar";
import { Room, roomStaticRules } from "src/schema/room";
import { createRoom } from "src/api/v1/room";
import { fetchRoomGroups } from "src/api/v1/room-group";
import { RoomGroup } from "src/schema/room-group";
import { getErrorMessage } from "src/util/errorHandler";

const router = useRouter();
const $q = useQuasar();

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

function create() {
  status.value.isProgress = true;

  createRoom(formData.value)
    .then(() => {
      router.push({ name: "Rooms" });

      resetForm();
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

function loadRoomGroups() {
  roomGroups.value = null;

  return fetchRoomGroups().then((response) => {
    roomGroups.value = response.values;
  });
}

function resetForm() {
  formData.value.number = "";
  formData.value.status = "NORMAL";
}

onBeforeMount(() => {
  resetForm();
  loadRoomGroups().then(() => {
    formData.value.roomGroup = roomGroups.value?.[0];
  });
});
</script>
