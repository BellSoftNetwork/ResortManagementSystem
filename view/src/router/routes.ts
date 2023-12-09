import DefaultLayout from "layouts/default/DefaultLayout.vue";
import { RouteRecordRaw } from "vue-router";

const routes: RouteRecordRaw[] = [
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
    path: "/my",
    component: DefaultLayout,
    meta: {
      isAuthenticated: true,
    },
    children: [
      {
        path: "",
        name: "MyDetail",
        component: () => import("pages/my/MyDetail.vue"),
      },
    ],
  },
  {
    path: "/rooms",
    component: DefaultLayout,
    meta: {
      isAuthenticated: true,
      roles: ["ADMIN", "SUPER_ADMIN"],
    },
    children: [
      {
        path: "",
        name: "Rooms",
        component: () => import("pages/room/RoomList.vue"),
      },
      {
        path: ":id",
        name: "Room",
        component: () => import("pages/room/RoomDetail.vue"),
      },
      {
        path: ":id/edit",
        name: "EditRoom",
        component: () => import("pages/room/RoomEdit.vue"),
      },
      {
        path: "create",
        name: "CreateRoom",
        component: () => import("pages/room/RoomCreate.vue"),
      },
    ],
  },
  {
    path: "/reservations",
    component: DefaultLayout,
    meta: {
      isAuthenticated: true,
      roles: ["ADMIN", "SUPER_ADMIN"],
    },
    children: [
      {
        path: "",
        name: "Reservations",
        component: () => import("pages/reservation/ReservationList.vue"),
      },
      {
        path: ":id",
        name: "Reservation",
        component: () => import("pages/reservation/ReservationDetail.vue"),
      },
      {
        path: ":id/edit",
        name: "EditReservation",
        component: () => import("pages/reservation/ReservationEdit.vue"),
      },
      {
        path: "create",
        name: "CreateReservation",
        component: () => import("pages/reservation/ReservationCreate.vue"),
      },
    ],
  },
  {
    path: "/payment-methods",
    component: DefaultLayout,
    meta: {
      isAuthenticated: true,
      roles: ["ADMIN", "SUPER_ADMIN"],
    },
    children: [
      {
        path: "",
        name: "PaymentMethods",
        component: () => import("pages/payment-method/PaymentMethodList.vue"),
      },
    ],
  },
  {
    path: "/admin",
    component: DefaultLayout,
    meta: {
      isAuthenticated: true,
      roles: ["ADMIN", "SUPER_ADMIN"],
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
    path: "/debug",
    component: DefaultLayout,
    meta: {
      isAuthenticated: true,
      roles: ["SUPER_ADMIN"],
    },
    children: [
      {
        path: "",
        name: "DevDebug",
        component: () => import("pages/DevDebug.vue"),
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
    path: "/register",
    meta: {
      isAuthenticated: false,
    },
    children: [
      {
        path: "",
        name: "Register",
        component: () => import("pages/auth/MainRegister.vue"),
      },
    ],
  },

  {
    path: "/error/403",
    component: DefaultLayout,
    children: [
      {
        path: "",
        name: "ErrorForbidden",
        component: () => import("pages/error/ErrorForbidden.vue"),
      },
    ],
  },
  {
    path: "/error/404",
    name: "ErrorNotFound",
    component: () => import("pages/error/ErrorNotFound.vue"),
  },

  {
    path: "/:catchAll(.*)*",
    component: () => import("pages/error/ErrorNotFound.vue"),
  },
];

export default routes;
