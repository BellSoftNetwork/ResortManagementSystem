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
        <ReservationMethodEditDialog @complete="emit('complete')" :reservationMethod="props.reservationMethod">
          <v-list-item-title>수정</v-list-item-title>
        </ReservationMethodEditDialog>
      </v-list-item>

      <v-divider />

      <v-list-item>
        <ReservationMethodDeleteDialog @complete="emit('complete')" :reservationMethod="props.reservationMethod">
          <v-list-item-title>삭제</v-list-item-title>
        </ReservationMethodDeleteDialog>
      </v-list-item>
    </v-list>
  </v-menu>
</template>

<script setup>
import { ref } from "vue"
import ReservationMethodEditDialog from "@/components/reservation-method/ReservationMethodEditDialog.vue"
import ReservationMethodDeleteDialog from "@/components/reservation-method/ReservationMethodDeleteDialog.vue"

const emit = defineEmits(["complete"])
const props = defineProps({
  reservationMethod: Object,
})
const status = ref({
  isDialogActive: false,
})
</script>
