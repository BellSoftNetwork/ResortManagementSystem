<template>
  <q-page class="row q-pa-lg">
    <div class="col-12">
      <div class="text-h4 q-mb-lg">개발 테스트 도구</div>

      <q-card class="q-mb-md">
        <q-card-section>
          <div class="text-h6 q-mb-md">더미 데이터 생성</div>
          <div class="text-caption text-grey q-mb-md">개발 및 테스트 환경에서만 사용 가능합니다.</div>
        </q-card-section>

        <q-separator />

        <q-card-section>
          <div class="row q-gutter-md">
            <div class="col-12">
              <q-btn
                @click="generateAllData"
                :loading="allDataLoading"
                color="positive"
                class="full-width"
                size="lg"
                unelevated
              >
                <q-icon left name="fas fa-magic" />
                전체 데이터 생성 (필수 + 예약)
              </q-btn>
              <div class="text-caption text-grey q-mt-sm">필수 데이터와 예약 데이터를 모두 생성합니다.</div>
            </div>

            <div class="col-12 col-md-6">
              <q-btn
                @click="generateEssentialData"
                :loading="essentialDataLoading"
                color="primary"
                class="full-width"
                size="lg"
                unelevated
              >
                <q-icon left name="fas fa-database" />
                필수 데이터만 생성
              </q-btn>
              <div class="text-caption text-grey q-mt-sm">결제 수단, 객실 그룹, 객실 데이터를 생성합니다.</div>
            </div>

            <div class="col-12 col-md-6">
              <q-btn
                @click="showReservationOptions = true"
                :loading="reservationDataLoading"
                color="secondary"
                class="full-width"
                size="lg"
                unelevated
              >
                <q-icon left name="fas fa-calendar-check" />
                예약 데이터만 생성
              </q-btn>
              <div class="text-caption text-grey q-mt-sm">다양한 예약 시나리오 데이터를 생성합니다.</div>
            </div>
          </div>
        </q-card-section>
      </q-card>

      <q-card v-if="lastResult">
        <q-card-section>
          <div class="text-h6 q-mb-md">실행 결과</div>
        </q-card-section>

        <q-separator />

        <q-card-section>
          <pre class="text-caption">{{ JSON.stringify(lastResult, null, 2) }}</pre>
        </q-card-section>
      </q-card>
    </div>

    <!-- 예약 옵션 다이얼로그 -->
    <q-dialog v-model="showReservationOptions">
      <q-card style="min-width: 400px">
        <q-card-section>
          <div class="text-h6">예약 데이터 생성 옵션</div>
        </q-card-section>

        <q-separator />

        <q-card-section>
          <div class="row q-gutter-md">
            <div class="col-12">
              <q-input
                v-model="reservationOptions.startDate"
                label="시작일"
                type="date"
                filled
                hint="예약 생성 기간의 시작일"
              />
            </div>

            <div class="col-12">
              <q-input
                v-model="reservationOptions.endDate"
                label="종료일"
                type="date"
                filled
                hint="예약 생성 기간의 종료일"
              />
            </div>

            <div class="col-12">
              <q-input
                v-model.number="reservationOptions.regularReservations"
                label="일반 예약 건수"
                type="number"
                filled
                hint="생성할 일반 예약 건수"
              />
            </div>

            <div class="col-12">
              <q-input
                v-model.number="reservationOptions.monthlyReservations"
                label="달방 예약 건수"
                type="number"
                filled
                hint="생성할 월세 예약 건수"
              />
            </div>
          </div>
        </q-card-section>

        <q-separator />

        <q-card-actions align="right">
          <q-btn flat label="취소" v-close-popup />
          <q-btn flat label="기본값으로 생성" color="primary" @click="generateReservationDataWithDefaults" />
          <q-btn flat label="옵션 적용하여 생성" color="primary" @click="generateReservationDataWithOptions" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup lang="ts">
import { ref, reactive } from "vue";
import { useQuasar } from "quasar";
import { generateTestData, type ReservationGenerationOptions } from "src/services/dev-test";

const $q = useQuasar();

const allDataLoading = ref(false);
const essentialDataLoading = ref(false);
const reservationDataLoading = ref(false);
const lastResult = ref<any>(null);
const showReservationOptions = ref(false);

// 예약 옵션 기본값 설정
const today = new Date();
const oneMonthLater = new Date(today);
oneMonthLater.setMonth(oneMonthLater.getMonth() + 1);

const reservationOptions = reactive<ReservationGenerationOptions>({
  startDate: today.toISOString().split("T")[0],
  endDate: oneMonthLater.toISOString().split("T")[0],
  regularReservations: 20,
  monthlyReservations: 5,
});

async function generateAllData() {
  allDataLoading.value = true;
  lastResult.value = null;

  try {
    const result = await generateTestData({
      type: "all",
      reservationOptions: {
        ...reservationOptions,
        startDate: reservationOptions.startDate ? `${reservationOptions.startDate}T00:00:00Z` : undefined,
        endDate: reservationOptions.endDate ? `${reservationOptions.endDate}T23:59:59Z` : undefined,
      },
    });
    lastResult.value = result.data;
    $q.notify({
      type: "positive",
      message: "전체 데이터가 성공적으로 생성되었습니다.",
    });
  } catch (error: any) {
    console.error("Failed to generate all data:", error);
    $q.notify({
      type: "negative",
      message: error.response?.data?.message || "전체 데이터 생성에 실패했습니다.",
    });
  } finally {
    allDataLoading.value = false;
  }
}

async function generateEssentialData() {
  essentialDataLoading.value = true;
  lastResult.value = null;

  try {
    const result = await generateTestData({ type: "essential" });
    lastResult.value = result.data;
    $q.notify({
      type: "positive",
      message: "필수 데이터가 성공적으로 생성되었습니다.",
    });
  } catch (error: any) {
    console.error("Failed to generate essential data:", error);
    $q.notify({
      type: "negative",
      message: error.response?.data?.message || "필수 데이터 생성에 실패했습니다.",
    });
  } finally {
    essentialDataLoading.value = false;
  }
}

async function generateReservationDataWithDefaults() {
  showReservationOptions.value = false;
  reservationDataLoading.value = true;
  lastResult.value = null;

  try {
    const result = await generateTestData({ type: "reservation" });
    lastResult.value = result.data;
    $q.notify({
      type: "positive",
      message: "예약 데이터가 성공적으로 생성되었습니다.",
    });
  } catch (error: any) {
    console.error("Failed to generate reservation data:", error);
    $q.notify({
      type: "negative",
      message: error.response?.data?.message || "예약 데이터 생성에 실패했습니다.",
    });
  } finally {
    reservationDataLoading.value = false;
  }
}

async function generateReservationDataWithOptions() {
  showReservationOptions.value = false;
  reservationDataLoading.value = true;
  lastResult.value = null;

  try {
    const result = await generateTestData({
      type: "reservation",
      reservationOptions: {
        ...reservationOptions,
        startDate: reservationOptions.startDate ? `${reservationOptions.startDate}T00:00:00Z` : undefined,
        endDate: reservationOptions.endDate ? `${reservationOptions.endDate}T23:59:59Z` : undefined,
      },
    });
    lastResult.value = result.data;
    $q.notify({
      type: "positive",
      message: "예약 데이터가 성공적으로 생성되었습니다.",
    });
  } catch (error: any) {
    console.error("Failed to generate reservation data:", error);
    $q.notify({
      type: "negative",
      message: error.response?.data?.message || "예약 데이터 생성에 실패했습니다.",
    });
  } finally {
    reservationDataLoading.value = false;
  }
}
</script>

<style scoped>
pre {
  background-color: #f5f5f5;
  padding: 10px;
  border-radius: 4px;
  overflow-x: auto;
}
</style>
