<template>
  <slot :dialog="dialog">
    <q-btn @click="dialog.isOpen = true">
      추가
    </q-btn>
  </slot>

  <q-dialog
    v-model="dialog.isOpen"
    :persistent="status.isProgress"
    @beforeShow="resetForm"
  >
    <q-card style="width: 500px">
      <q-card-section class="text-h6">
        계정 수정
      </q-card-section>

      <q-form @submit="update">
        <q-card-section>
          <q-input
            v-model="formData.email"
            :rules="rules.email"
            :readonly="true"
            :disable="true"
            label="이메일"
          ></q-input>

          <q-input
            v-model="formData.name"
            :rules="rules.name"
            label="이름"
            required
          ></q-input>

          <q-input
            v-model="formData.password"
            :rules="rules.password"
            type="password"
            label="비밀번호"
          ></q-input>

          <q-select
            v-model="formData.role"
            :options="options.role"
            label="권한"
            required
            emit-value
            map-options
          ></q-select>
        </q-card-section>

        <q-card-actions align="right">
          <q-btn
            v-close-popup
            :disable="status.isProgress"
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
  </q-dialog>
</template>

<script setup>
import { ref } from "vue"
import { useRouter } from "vue-router"
import { useAuthStore } from "stores/auth.js"
import { useQuasar } from "quasar"
import { api } from "boot/axios"

const router = useRouter()
const authStore = useAuthStore()
const $q = useQuasar()

const emit = defineEmits(["complete"])
const props = defineProps({
  entity: Object,
})
const dialog = ref({
  isOpen: false,
})
const status = ref({
  isProgress: false,
})
const formData = ref({
  email: props.entity.email,
  name: props.entity.name,
  password: "",
  role: props.entity.role,
})
const rules = {
  name: [value => (value.length >= 2 && value.length <= 20) || "2~20 글자가 필요합니다"],
  email: [value => /^\w+([.-]?\w+)*@\w+([.-]?\w+)*(\.\w{2,3})+$/.test(value) || "이메일이 유효하지 않습니다."],
  password: [value => (value.length === 0 || (value.length >= 8 && value.length <= 20)) || "비밀번호는 8~20 글자가 필요합니다."],
}
const options = {
  role: [
    { label: "일반", value: "NORMAL" },
    { label: "관리자", value: "ADMIN", disable: !authStore.isSuperAdminRole },
  ],
}

function update() {
  if (!isChanged()) {
    $q.notify({
      message: "수정된 항목이 없습니다.",
      type: "info",
      actions: [
        {
          icon: "close", color: "white", round: true,
        },
      ],
    })

    return
  }

  status.value.isProgress = true

  api.patch(`/api/v1/admin/accounts/${props.entity.id}`, patchedData())
    .then(() => {
      emit("complete")
      dialog.value.isOpen = false
    })
    .catch((error) => {
      $q.notify({
        message: error.response.data.message,
        type: "negative",
        actions: [
          {
            icon: "close", color: "white", round: true,
          },
        ],
      })
    })
    .finally(() => {
      status.value.isProgress = false
    })
}

function isChanged() {
  return formData.value.name !== props.entity.name ||
    formData.value.role !== props.entity.role ||
    formData.value.password.length > 0
}

function patchedData() {
  const patchData = {}

  if (props.entity.name !== formData.value.name)
    patchData.name = formData.value.name
  if (props.entity.role !== formData.value.role)
    patchData.role = formData.value.role
  if (formData.value.password.length > 0)
    patchData.password = formData.value.password

  return patchData
}

function resetForm() {
  formData.value.email = props.entity.email
  formData.value.name = props.entity.name
  formData.value.password = ""
  formData.value.role = props.entity.role
}
</script>