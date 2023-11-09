<template>
  <v-menu v-model="status.isDialogActive">
    <template v-slot:activator="{ props }">
      <v-btn
        v-bind="props"
        variant="text"
      >
        <v-icon>fa-solid fa-ellipsis-vertical</v-icon>
      </v-btn>
    </template>

    <v-list>
      <v-list-item>
        <AccountEditDialog @complete="emit('complete')" :account="props.account">
          <v-list-item-title>수정</v-list-item-title>
        </AccountEditDialog>
      </v-list-item>

      <v-divider />

      <v-list-item>
        <!-- TODO: 계정 잠금 확인 dialog 컴포넌트 구현 필요 -->
        <v-tooltip text="구현 예정">
          <template v-slot:activator="{ props }">
            <div v-bind="props">
              <v-btn variant="text" color="red" disabled>
                <v-list-item-title>잠금</v-list-item-title>
              </v-btn>
            </div>
          </template>
        </v-tooltip>
      </v-list-item>

      <v-list-item>
        <!-- TODO: 계정 삭제 확인 dialog 컴포넌트 구현 필요 -->
        <v-tooltip text="구현 예정">
          <template v-slot:activator="{ props }">
            <div v-bind="props">
              <v-btn variant="text" color="red" disabled>
                <v-list-item-title>삭제</v-list-item-title>
              </v-btn>
            </div>
          </template>
        </v-tooltip>
      </v-list-item>
    </v-list>
  </v-menu>
</template>

<script setup>
import { ref } from "vue"
import AccountEditDialog from "@/components/account/AccountEditDialog.vue"

const emit = defineEmits(["complete"])
const props = defineProps({
  account: Object,
})
const status = ref({
  isDialogActive: false,
})
</script>
