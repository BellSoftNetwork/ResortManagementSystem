import accountInfoLoader from "./account-info-loader"
import authGuard from "./authenticate-guard"
import roleGuard from "./role-guard"

export default function registerRouterGuards(router) {
  accountInfoLoader(router)
  authGuard(router)
  roleGuard(router)
}
