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
    path: "/room-status",
    component: DefaultLayout,
    meta: {
      isAuthenticated: true,
    },
    children: [
      {
        path: "",
        name: "RoomStatus",
        component: () => import("pages/room/RoomStatus.vue"),
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
    path: "/room-groups",
    component: DefaultLayout,
    meta: {
      isAuthenticated: true,
      roles: ["ADMIN", "SUPER_ADMIN"],
    },
    children: [
      {
        path: "",
        name: "RoomGroups",
        component: () => import("pages/room-group/RoomGroupList.vue"),
      },
      {
        path: ":id",
        name: "RoomGroup",
        component: () => import("pages/room-group/RoomGroupDetail.vue"),
      },
      {
        path: ":id/edit",
        name: "EditRoomGroup",
        component: () => import("pages/room-group/RoomGroupEdit.vue"),
      },
      {
        path: "create",
        name: "CreateRoomGroup",
        component: () => import("pages/room-group/RoomGroupCreate.vue"),
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
    path: "/monthly-rents",
    component: DefaultLayout,
    meta: {
      isAuthenticated: true,
      roles: ["ADMIN", "SUPER_ADMIN"],
    },
    children: [
      {
        path: "",
        name: "MonthlyRents",
        component: () => import("pages/monthly-rent/MonthlyRentList.vue"),
      },
      {
        path: ":id",
        name: "MonthlyRent",
        component: () => import("pages/monthly-rent/MonthlyRentDetail.vue"),
      },
      {
        path: ":id/edit",
        name: "EditMonthlyRent",
        component: () => import("pages/monthly-rent/MonthlyRentEdit.vue"),
      },
      {
        path: "create",
        name: "CreateMonthlyRent",
        component: () => import("pages/monthly-rent/MonthlyRentCreate.vue"),
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
      {
        path: "audit-logs",
        name: "AuditLogs",
        component: () => import("pages/admin/AuditLogList.vue"),
      },
      {
        path: "dev-test",
        name: "DevTest",
        component: () => import("pages/admin/DevTest.vue"),
        meta: {
          roles: ["SUPER_ADMIN"],
        },
      },
    ],
  },
  {
    path: "/stats",
    component: DefaultLayout,
    meta: {
      isAuthenticated: true,
      roles: ["ADMIN", "SUPER_ADMIN"],
    },
    children: [
      {
        path: "",
        name: "Stats",
        component: () => import("pages/stats/StatsPage.vue"),
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
