<template>
  <q-card flat bordered>
    <q-card-section class="text-h6">
      {{ room.number }}
    </q-card-section>

    <q-card-section>
      <div class="row">
        <div class="col q-py-sm">
          <div class="text-caption">객실 그룹</div>
          <div class="text-body1">{{ props.room.roomGroup.name }}</div>
        </div>

        <div class="col q-py-sm">
          <div class="text-caption">객실 상태</div>
          <div class="text-body1">
            {{ roomStatusValueToName(props.room.status) }}
          </div>
        </div>
      </div>

      <div class="row">
        <div class="col">
          <div class="q-py-sm">
            <div class="text-caption">메모 (관리용)</div>
            <div class="text-body1">{{ props.room.note }}</div>
          </div>
        </div>
      </div>
    </q-card-section>

    <q-card-actions align="right">
      <q-btn @click="deleteItem()" color="red" label="삭제" dense flat></q-btn>
      <q-btn :to="{ name: 'EditRoom', params: { id: room.id } }" color="primary" label="수정" flat />
    </q-card-actions>
  </q-card>
</template>

<script setup lang="ts">
import { useRouter } from "vue-router";
import { useQuasar } from "quasar";
import { deleteRoom } from "src/api/v1/room";
import { Room, roomStatusValueToName } from "src/schema/room";

const router = useRouter();
const $q = useQuasar();
const props = defineProps<{
  room: Room;
}>();

function deleteItem() {
  const itemId = props.room.id;
  const itemName = props.room.number;

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
</script>
