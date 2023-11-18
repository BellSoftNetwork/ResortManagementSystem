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
    path: "/rooms",
    component: DefaultLayout,
    meta: {
      isAuthenticated: true,
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
]

export default routes
