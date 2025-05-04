import { useAuthStore } from "stores/auth";
import { Router } from "vue-router";

export default (router: Router) => {
  router.beforeEach((to, from, next) => {
    const authStore = useAuthStore();
    const roles = to.meta.roles;

    if (!Array.isArray(roles) || roles.length <= 0) return next();

    // 토큰 갱신 중이면 권한 검사를 건너뜁니다
    if (authStore.isRefreshingToken) {
      return next();
    }

    if (!roles.includes(authStore.user?.role)) return next({ name: "ErrorForbidden" });

    return next();
  });
};
