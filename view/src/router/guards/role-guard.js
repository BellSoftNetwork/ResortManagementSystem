import { useAuthStore } from "stores/auth.js"

export default (router) => {
  router.beforeEach((to, from, next) => {
    const authStore = useAuthStore()
    const roles = to.meta.roles

    if (!Array.isArray(roles) || roles.length <= 0)
      return next()

    if (!roles.includes(authStore.user.role))
      return next({ name: "ErrorForbidden" })

    return next()
  })
}
