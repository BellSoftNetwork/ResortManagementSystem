<template>
  <q-card>
    <q-card-section class="text-h6">객실 그룹 추가</q-card-section>

    <q-form @submit="create">
      <q-card-section>
        <q-input
          v-model="formData.name"
          :rules="roomGroupStaticRules.name"
          label="이름"
          placeholder="20평형"
          required
        ></q-input>

        <q-input
          v-model.number="formData.peekPrice"
          :rules="roomStaticRules.peekPrice"
          label="성수기 예약금"
          type="number"
          min="0"
          max="100000000"
          required
        ></q-input>

        <q-input
          v-model.number="formData.offPeekPrice"
          :rules="roomStaticRules.offPeekPrice"
          label="비성수기 예약금"
          type="number"
          min="0"
          max="100000000"
          required
        ></q-input>

        <q-input
          v-model="formData.description"
          :rules="roomGroupStaticRules.description"
          type="textarea"
          label="설명"
          placeholder="와이파이 사용 가능"
        ></q-input>
      </q-card-section>

      <q-card-actions align="right">
        <q-btn :disable="status.isProgress" :to="{ name: 'RoomGroups' }" color="primary" label="취소" flat />
        <q-btn :loading="status.isProgress" type="submit" color="red" label="추가" flat />
      </q-card-actions>
    </q-form>
  </q-card>
</template>

<script setup lang="ts">
import { onBeforeMount, ref } from "vue";
import { useRouter } from "vue-router";
import { useQuasar } from "quasar";
import { roomGroupStaticRules } from "src/schema/room-group";
import { createRoomGroup } from "src/api/v1/room-group";
import { roomStaticRules } from "src/schema/room";

const router = useRouter();
const $q = useQuasar();

const status = ref({
  isProgress: false,
});
const formData = ref({
  name: "",
  peekPrice: 0,
  offPeekPrice: 0,
  description: "",
});

function create() {
  status.value.isProgress = true;

  createRoomGroup(formData.value)
    .then(() => {
      router.push({ name: "RoomGroups" });

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
  formData.value.peekPrice = 0;
  formData.value.offPeekPrice = 0;
  formData.value.description = "";
}

onBeforeMount(() => {
  resetForm();
});
</script>
