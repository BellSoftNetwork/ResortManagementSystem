<template>
  <div>
    <!-- 로그인 페이지 포함 모든 페이지에서 서버 오류 시 ServerError 표시 -->
    <ServerError v-if="networkStore.isOffline" />

    <!-- Normal App Content -->
    <template v-else>
      <!-- 현재 경로가 로그인 페이지인 경우 항상 라우터 뷰 표시 -->
      <router-view v-if="isLoginPage" />
      <!-- 로그인 페이지가 아닌 경우에만 토큰 갱신 중 로딩 표시 -->
      <template v-else>
        <div v-if="authStore.isRefreshingToken" class="fullscreen bg-white text-center flex flex-center column">
          <q-spinner-dots size="80px" color="primary" />
          <div class="q-mt-md">인증 정보를 확인하는 중입니다...</div>
        </div>
        <router-view v-else />
      </template>
    </template>
  </div>
</template>

<script>
import { computed, defineComponent } from "vue";
import { useAuthStore } from "stores/auth";
import { useNetworkStore } from "stores/network";
import { useRoute } from "vue-router";
import ServerError from "components/common/ServerError.vue";

export default defineComponent({
  name: "App",
  components: { ServerError },
  setup() {
    const authStore = useAuthStore();
    const networkStore = useNetworkStore();
    const route = useRoute();

    // 현재 경로가 로그인 페이지인지 확인
    const isLoginPage = computed(() => {
      return route.name === "Login" || route.path === "/login";
    });

    return { authStore, networkStore, isLoginPage };
  },
});
</script>
