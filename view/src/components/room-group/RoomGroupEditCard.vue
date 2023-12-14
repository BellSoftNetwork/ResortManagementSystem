<template>
  <q-card>
    <q-card-section class="text-h6">객실 그룹 수정</q-card-section>

    <q-form @submit="update">
      <q-card-section>
        <q-input
          v-model="formData.name"
          :loading="status.isProgress"
          :disable="status.isProgress"
          :rules="roomGroupStaticRules.name"
          label="번호"
          placeholder="101호"
          required
        ></q-input>

        <q-input
          v-model.number="formData.peekPrice"
          :loading="status.isProgress"
          :disable="status.isProgress"
          :rules="roomGroupStaticRules.peekPrice"
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
          :rules="roomGroupStaticRules.offPeekPrice"
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
          :rules="roomGroupStaticRules.description"
          type="textarea"
          label="설명"
          placeholder="와이파이 사용 가능"
        ></q-input>
      </q-card-section>

      <q-card-actions align="right">
        <q-btn :disable="status.isProgress" :to="{ name: 'RoomGroups' }" color="primary" label="취소" flat />
        <q-btn :loading="status.isProgress" type="submit" color="red" label="수정" flat />
      </q-card-actions>
    </q-form>
  </q-card>
</template>

<script setup lang="ts">
import { onBeforeMount, ref } from "vue";
import { useRouter } from "vue-router";
import { useQuasar } from "quasar";
import { getPatchedFormData, isFormValueChanged } from "src/util/data-util";
import { RoomGroup, roomGroupStaticRules } from "src/schema/room-group";
import { patchRoomGroup } from "src/api/v1/room-group";

const router = useRouter();
const $q = useQuasar();
const props = defineProps<{
  roomGroup: RoomGroup;
}>();
const status = ref({
  isProgress: false,
});
const formData = ref<Partial<RoomGroup>>({
  name: "",
  peekPrice: 0,
  offPeekPrice: 0,
  description: "",
});

function update() {
  if (!isFormValueChanged(props.roomGroup, formData.value)) {
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

  patchRoomGroup(props.roomGroup.id, getPatchedFormData(props.roomGroup, formData.value))
    .then(() => {
      router.push({ name: "RoomGroup", params: { id: props.roomGroup.id } });
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
  Object.assign(formData.value, props.roomGroup);
}

onBeforeMount(() => {
  resetForm();
});
</script>
