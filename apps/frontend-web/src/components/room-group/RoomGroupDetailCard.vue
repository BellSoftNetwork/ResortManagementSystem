<template>
  <q-card flat bordered>
    <q-card-section class="text-h6">
      {{ props.roomGroup.name }}
    </q-card-section>

    <q-card-section>
      <div class="row">
        <div class="col q-py-sm">
          <div class="text-caption">성수기 예약금</div>
          <div class="text-body1">{{ formatPrice(props.roomGroup.peekPrice) }}</div>
        </div>

        <div class="col q-py-sm">
          <div class="text-caption">비성수기 예약금</div>
          <div class="text-body1">{{ formatPrice(props.roomGroup.offPeekPrice) }}</div>
        </div>
      </div>

      <div class="row">
        <div class="col">
          <div class="q-py-sm">
            <div class="text-caption">메모</div>
            <div class="text-body1">{{ props.roomGroup.description }}</div>
          </div>
        </div>
      </div>
    </q-card-section>

    <q-card-actions align="right">
      <q-btn @click="deleteItem()" color="red" label="삭제" dense flat></q-btn>
      <q-btn :to="{ name: 'EditRoomGroup', params: { id: props.roomGroup.id } }" color="primary" label="수정" flat />
    </q-card-actions>
  </q-card>
</template>

<script setup lang="ts">
import { useRouter } from "vue-router";
import { useQuasar } from "quasar";
import { RoomGroup } from "src/schema/room-group";
import { deleteRoomGroup } from "src/api/v1/room-group";
import { formatPrice } from "src/util/format-util";

const router = useRouter();
const $q = useQuasar();
const props = defineProps<{
  roomGroup: RoomGroup;
}>();

function deleteItem() {
  const itemId = props.roomGroup.id;
  const itemName = props.roomGroup.name;

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
    deleteRoomGroup(itemId)
      .then(() => {
        router.push({ name: "RoomGroups" });
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
</script>
