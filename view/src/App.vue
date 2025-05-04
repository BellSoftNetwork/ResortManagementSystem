<template>
  <div>
    <!-- 토큰 갱신 중일 때 로딩 표시 -->
    <div v-if="authStore.isRefreshingToken" class="fullscreen bg-white text-center flex flex-center column">
      <q-spinner-dots size="80px" color="primary" />
      <div class="q-mt-md">인증 정보를 확인하는 중입니다...</div>
    </div>
    <!-- 토큰 갱신 중이 아닐 때만 라우터 뷰 표시 -->
    <router-view v-else />
  </div>
</template>

<script>
import { defineComponent } from "vue";
import { useAuthStore } from "stores/auth";

export default defineComponent({
  name: "App",
  setup() {
    const authStore = useAuthStore();
    return { authStore };
  },
});
</script>
