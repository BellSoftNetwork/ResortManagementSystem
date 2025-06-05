import { useAuthStore } from "stores/auth";
import { Router } from "vue-router";

export default (router: Router) => {
  router.beforeEach((to, from, next) => {
    const authStore = useAuthStore();

    if (!(to.meta.isAuthenticated === true || to.meta.roles)) return next();

    // 토큰 갱신 중이면 대기 (다음 라우터 가드에서 처리)
    if (authStore.isRefreshingToken) {
      return next();
    }

    if (!authStore.isLoggedIn) return next({ name: "Login" });

    return next();
  });
};
