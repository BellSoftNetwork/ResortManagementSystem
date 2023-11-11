import { useAuthStore } from "stores/auth.js"

export default (router) => {
  router.beforeEach((to, from, next) => {
    const authStore = useAuthStore()

    if (!authStore.isFirstRequest)
      return next()

    authStore.loadAccountInfo().finally(() => next())
  })
}
