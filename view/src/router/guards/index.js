import accountInfoLoader from "@/router/guards/account-info-loader"
import authGuard from "@/router/guards/authenticate-guard"
import roleGuard from "@/router/guards/role-guard"

export default function registerRouterGuards(router) {
  accountInfoLoader(router)
  authGuard(router)
  roleGuard(router)
}
