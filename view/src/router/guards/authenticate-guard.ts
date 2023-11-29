import { useAuthStore } from "stores/auth";
import { Router } from "vue-router";

export default (router: Router) => {
  router.beforeEach((to, from, next) => {
    const authStore = useAuthStore();

    if (!(to.meta.isAuthenticated === true || to.meta.roles)) return next();

    if (!authStore.isLoggedIn) return next({ name: "Login" });

    return next();
  });
};
