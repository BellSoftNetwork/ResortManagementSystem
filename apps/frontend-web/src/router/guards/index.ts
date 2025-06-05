import accountInfoLoader from "./account-info-loader";
import authGuard from "./authenticate-guard";
import roleGuard from "./role-guard";
import { Router } from "vue-router";

export default function registerRouterGuards(router: Router) {
  accountInfoLoader(router);
  authGuard(router);
  roleGuard(router);
}
