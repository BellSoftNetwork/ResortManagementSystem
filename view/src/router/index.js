// Composable
import { createRouter, createWebHistory } from "vue-router"
import registerRouterGuards from "@/router/guards"
import DefaultLayout from "@/layouts/default/DefaultLayout.vue"
import NotFound from "@/views/error/NotFound.vue"
import Forbidden from "@/views/error/Forbidden.vue"

const routes = [
  {
    path: "/",
    component: DefaultLayout,
    meta: {
      isAuthenticated: true,
    },
    children: [
      {
        path: "",
        name: "Home",
        component: () => import(/* webpackChunkName: "home" */ "@/views/DefaultHome.vue"),
      },
    ],
  },
  {
    path: "/reservation-methods",
    component: DefaultLayout,
    meta: {
      isAuthenticated: true,
    },
    children: [
      {
        path: "",
        name: "ReservationMethods",
        component: () => import(/* webpackChunkName: "reservationMethods" */ "@/views/reservation-method/ReservationMethodList.vue"),
      },
    ],
  },
  {
    path: "/admin",
    component: DefaultLayout,
    meta: {
      roles: ["ADMIN", "SUPER_ADMIN"],
      isAuthenticated: true,
    },
    children: [
      {
        path: "accounts",
        name: "AdminAccounts",
        component: () => import(/* webpackChunkName: "adminAccounts" */ "@/views/admin/AccountList.vue"),
      },
    ],
  },
  {
    path: "/login",
    meta: {
      isAuthenticated: false,
    },
    children: [
      {
        path: "",
        name: "Login",
        component: () => import(/* webpackChunkName: "login" */ "@/views/auth/MainLogin.vue"),
      },
    ],
  },
  {
    path: "/error/403",
    component: DefaultLayout,
    children: [
      {
        path: "",
        name: "forbidden",
        component: Forbidden,
      },
    ],
  },
  {
    path: "/:pathMatch(.*)*",
    component: DefaultLayout,
    children: [
      {
        path: "",
        name: "notFound",
        component: NotFound,
      },
    ],
  },
]

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes,
})

registerRouterGuards(router)

export default router
