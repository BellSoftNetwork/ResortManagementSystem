import { useAuthStore } from "@/store/auth.js"

export default (router) => {
  router.beforeEach((to, from, next) => {
    const authStore = useAuthStore()

    if (!(to.meta.isAuthenticated === true || to.meta.roles))
      return next()

    if (!authStore.isLoggedIn)
      return next({ name: "Login" })

    return next()
  })
}
