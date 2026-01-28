<template>
  <q-dialog
    :model-value="modelValue"
    @update:model-value="$emit('update:modelValue', $event)"
    @before-show="resetFilterBuffer"
    :maximized="$q.screen.lt.md"
    transition-show="slide-up"
    transition-hide="slide-down"
    persistent
  >
    <q-card flat>
      <q-card-section>
        <q-input v-model="filterBuffer.peopleInfo" placeholder="홍길동" label="예약자 정보" class="fit" />
      </q-card-section>

      <q-card-section>
        <div class="row q-col-gutter-sm">
          <div class="col-12 col-sm-3">
            <q-select
              v-model="filterBuffer.dueOption"
              @update:model-value="updateDueDate"
              :options="dueOptions"
              label="검색 기간"
              emit-value
              map-options
              outlined
            />
          </div>

          <div class="col-12 col-sm-9">
            <div class="row no-wrap">
              <q-input
                v-model="filterBuffer.stayStartAt"
                mask="####-##-##"
                :readonly="true"
                :bg-color="filterBuffer.dueOption !== 'CUSTOM' ? 'grey-4' : ''"
                class="due-date-text"
                outlined
              >
                <template v-slot:append>
                  <q-icon @click="filterBuffer.dueOption = 'CUSTOM'" name="event" class="cursor-pointer">
                    <q-popup-proxy cover transition-show="scale" transition-hide="scale">
                      <q-date v-model="filterBuffer.stayStartAt" mask="YYYY-MM-DD">
                        <div class="row items-center justify-end">
                          <q-btn v-close-popup label="Close" color="primary" flat />
                        </div>
                      </q-date>
                    </q-popup-proxy>
                  </q-icon>
                </template>
              </q-input>
              <span class="self-center q-mx-sm">~</span>
              <q-input
                v-model="filterBuffer.stayEndAt"
                mask="####-##-##"
                :readonly="true"
                :bg-color="filterBuffer.dueOption !== 'CUSTOM' ? 'grey-4' : ''"
                class="due-date-text"
                outlined
              >
                <template v-slot:append>
                  <q-icon @click="filterBuffer.dueOption = 'CUSTOM'" name="event" class="cursor-pointer">
                    <q-popup-proxy cover transition-show="scale" transition-hide="scale">
                      <q-date v-model="filterBuffer.stayEndAt" mask="YYYY-MM-DD">
                        <div class="row items-center justify-end">
                          <q-btn v-close-popup label="Close" color="primary" flat />
                        </div>
                      </q-date>
                    </q-popup-proxy>
                  </q-icon>
                </template>
              </q-input>
            </div>
          </div>
        </div>
      </q-card-section>

      <q-card-section>
        <q-select
          v-model="filterBuffer.status"
          :options="statusOptions"
          class="full-width"
          label="예약 상태"
          emit-value
          map-options
          outlined
        />
      </q-card-section>

      <q-card-actions align="right">
        <q-btn @click="applyFilter" color="primary">적용</q-btn>
        <q-btn @click="$emit('update:modelValue', false)">취소</q-btn>
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { useQuasar } from "quasar";
import { calculateDateRange, type DueOption } from "src/util/date-preset-util";

interface FilterType {
  peopleInfo: string;
  dueOption: string;
  stayStartAt: string;
  stayEndAt: string;
  status: string;
}

interface Props {
  modelValue: boolean;
  filter: FilterType;
  dueOptions: Array<{ label: string; value: string }>;
  statusOptions: Array<{ label: string; value: string }>;
  defaultStayStartAt: string;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "update:modelValue", value: boolean): void;
  (e: "apply", filter: FilterType): void;
}>();

const $q = useQuasar();

const filterBuffer = ref<FilterType>({
  peopleInfo: "",
  dueOption: "6M",
  stayStartAt: "",
  stayEndAt: "",
  status: "NORMAL",
});

function resetFilterBuffer() {
  Object.assign(filterBuffer.value, props.filter);
}

function updateDueDate(dueOption: string) {
  const range = calculateDateRange(dueOption as DueOption, props.defaultStayStartAt);
  filterBuffer.value.stayStartAt = range.startAt;
  filterBuffer.value.stayEndAt = range.endAt;
}

function applyFilter() {
  emit("apply", { ...filterBuffer.value });
}
</script>
