import { useAuthStore } from "stores/auth"
import { Router } from "vue-router"

export default (router: Router) => {
  router.beforeEach((to, from, next) => {
    const authStore = useAuthStore()

    if (!authStore.isFirstRequest) return next()

    authStore.loadAccountInfo().finally(() => next())
  });
};
