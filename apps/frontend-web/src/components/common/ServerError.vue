<template>
  <div class="fullscreen bg-grey-2 text-center flex flex-center column">
    <div class="q-pa-lg">
      <!-- 에러 타입별 아이콘 -->
      <q-icon
        :name="networkStore.isServerError ? 'error_outline' : 'cloud_off'"
        size="80px"
        :color="networkStore.isServerError ? 'negative' : 'grey-6'"
      />

      <!-- 에러 타입별 제목 -->
      <h5 class="q-mt-lg q-mb-sm text-grey-8">
        {{ networkStore.isServerError ? "서버 오류" : "연결 끊김" }}
      </h5>

      <!-- 동적 에러 메시지 -->
      <p class="text-grey-6 q-mb-sm">
        {{ networkStore.errorMessage }}
      </p>

      <!-- 상태 코드 표시 (서버 에러 시) -->
      <p v-if="networkStore.statusCode" class="text-grey-5 text-caption q-mb-md">
        오류 코드: {{ networkStore.statusCode }}
      </p>

      <!-- 서버 에러가 아닐 때 여백 유지 -->
      <div v-else class="q-mb-md"></div>

      <!-- 자동 재시도 진행 상태 -->
      <div v-if="networkStore.autoRetryEnabled" class="q-mb-lg">
        <div class="flex items-center justify-center q-mb-xs">
          <q-spinner-dots color="primary" size="24px" class="q-mr-sm" />
          <span class="text-grey-7 text-body2"> 재연결 시도 중... ({{ networkStore.retryProgress }}) </span>
        </div>
        <p class="text-grey-5 text-caption q-mb-none">{{ countdown }}초 후 다시 시도</p>
      </div>

      <!-- 자동 재시도 비활성화 시 여백 유지 -->
      <div v-else class="q-mb-lg"></div>

      <q-btn color="primary" label="지금 다시 시도" icon="refresh" class="q-px-lg" @click="retry" />
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, onMounted, onUnmounted } from "vue";
import { useNetworkStore, RETRY_INTERVAL_MS } from "stores/network";
import axios from "axios";

export default defineComponent({
  name: "ServerError",
  setup() {
    const networkStore = useNetworkStore();
    const countdown = ref(Math.floor(RETRY_INTERVAL_MS / 1000));

    let retryIntervalId: number | null = null;
    let countdownIntervalId: number | null = null;

    /**
     * 헬스체크 API 호출
     */
    const checkHealth = async () => {
      try {
        await axios.get("/api/v1/config");
        // 성공 시 온라인 복구 후 페이지 새로고침
        networkStore.setOnline();
        window.location.reload();
      } catch {
        // 실패 시 재시도 횟수 증가
        networkStore.incrementAutoRetry();

        // 최대 횟수 도달 시 자동 재시도 중단
        if (!networkStore.canAutoRetry) {
          networkStore.stopAutoRetry();
        }
      }
    };

    /**
     * 카운트다운 시작
     */
    const startCountdown = () => {
      countdown.value = Math.floor(RETRY_INTERVAL_MS / 1000);

      countdownIntervalId = window.setInterval(() => {
        countdown.value--;
        if (countdown.value <= 0) {
          countdown.value = Math.floor(RETRY_INTERVAL_MS / 1000);
        }
      }, 1000);
    };

    /**
     * 자동 재시도 시작
     */
    const startAutoRetry = () => {
      networkStore.startAutoRetry();
      startCountdown();

      // 첫 번째 재시도 즉시 실행
      checkHealth();

      // 이후 10초마다 재시도
      retryIntervalId = window.setInterval(() => {
        if (networkStore.canAutoRetry) {
          checkHealth();
        } else {
          // 자동 재시도 중단
          if (retryIntervalId !== null) {
            clearInterval(retryIntervalId);
            retryIntervalId = null;
          }
          if (countdownIntervalId !== null) {
            clearInterval(countdownIntervalId);
            countdownIntervalId = null;
          }
        }
      }, RETRY_INTERVAL_MS);
    };

    /**
     * 수동 재시도 (즉시 페이지 새로고침)
     */
    const retry = () => {
      window.location.reload();
    };

    onMounted(() => {
      startAutoRetry();
    });

    onUnmounted(() => {
      if (retryIntervalId !== null) {
        clearInterval(retryIntervalId);
      }
      if (countdownIntervalId !== null) {
        clearInterval(countdownIntervalId);
      }
    });

    return { networkStore, countdown, retry };
  },
});
</script>
