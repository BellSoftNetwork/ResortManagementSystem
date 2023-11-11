import DefaultLayout from "layouts/default/DefaultLayout.vue"

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
        component: () => import("pages/DefaultHome.vue"),
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
        component: () => import("pages/reservation-method/ReservationMethodList.vue"),
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
        component: () => import("pages/admin/AccountList.vue"),
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
        component: () => import("pages/auth/MainLogin.vue"),
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
        component: () => import("pages/error/ErrorForbidden.vue"),
      },
    ],
  },

  {
    path: "/:catchAll(.*)*",
    component: () => import("pages/error/ErrorNotFound.vue"),
  },
]

export default routes
