import { useAuthStore } from "stores/auth";
import { Router } from "vue-router";

export default (router: Router) => {
  router.beforeEach((to, from, next) => {
    const authStore = useAuthStore();

    // 첫 요청이 아니고, 토큰이 있지만 사용자 정보가 없는 경우에도 로드
    if (!authStore.isFirstRequest && !(authStore.accessToken && !authStore.user)) return next();

    authStore.loadAccountInfo().finally(() => next());
  });
};
