import { defineStore } from "pinia";
import {
  fasBook,
  fasCalendar,
  fasChartLine,
  fasCommentDollar,
  fasHotel,
  fasPersonShelter,
  fasSignature,
  fasTableColumns,
  fasUser,
  fasFlask,
  fasClockRotateLeft as fasHistory,
  fasCalendarXmark,
} from "@quasar/extras/fontawesome-v6";
import { useAuthStore } from "stores/auth";

interface Link {
  icon: string;
  text: string;
  to: string;
  gnb?: boolean;
}

const normalLinks: Link[] = [{ icon: fasTableColumns, text: "대시보드", to: "Home", gnb: true }] as const;
const adminLinks: Link[] = [
  { icon: fasBook, text: "예약", to: "Reservations", gnb: true },
  { icon: fasSignature, text: "달방", to: "MonthlyRents", gnb: true },
  { icon: fasCalendar, text: "객실 현황", to: "RoomStatus", gnb: true },
  { icon: fasChartLine, text: "통계", to: "Stats", gnb: true },
  { icon: fasPersonShelter, text: "객실", to: "Rooms" },
  { icon: fasHotel, text: "객실 그룹", to: "RoomGroups" },
  { icon: fasCommentDollar, text: "결제 수단", to: "PaymentMethods" },
  { icon: fasUser, text: "계정 관리", to: "AdminAccounts" },
  { icon: fasHistory, text: "감사 로그", to: "AuditLogs" },
  { icon: fasCalendarXmark, text: "예약 마감", to: "DateBlocks" },
] as const;

const superAdminLinks: Link[] = [{ icon: fasFlask, text: "개발 테스트", to: "DevTest" }] as const;

export const useMenuStore = defineStore("menu", {
  getters: {
    allLinks: () => {
      const authStore = useAuthStore();
      const links = [normalLinks];

      if (authStore.isAdminRole) links.push(adminLinks);
      if (authStore.isSuperAdminRole) links.push(superAdminLinks);

      return links;
    },
    tabLinks: () => {
      const authStore = useAuthStore();
      const links = [...normalLinks];

      if (authStore.isAdminRole) links.push(...adminLinks);

      return links.filter((link) => link.gnb);
    },
  },
});
