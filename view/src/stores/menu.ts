import { defineStore } from "pinia";
import {
  fasBook,
  fasCommentDollar,
  fasHotel,
  fasPersonShelter,
  fasSignature,
  fasTableColumns,
} from "@quasar/extras/fontawesome-v6";
import { useAuthStore } from "stores/auth";
import { matCalendarViewMonth } from "@quasar/extras/material-icons";

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
  { icon: fasPersonShelter, text: "객실", to: "Rooms", gnb: true },
  { icon: matCalendarViewMonth, text: "객실 현황", to: "RoomStatus", gnb: true },
  { icon: fasHotel, text: "객실 그룹", to: "RoomGroups" },
  { icon: fasCommentDollar, text: "결제 수단", to: "PaymentMethods" },
  { icon: "person", text: "계정 관리", to: "AdminAccounts" },
] as const;

export const useMenuStore = defineStore("menu", {
  getters: {
    allLinks: () => {
      const authStore = useAuthStore();
      const links = [normalLinks];

      if (authStore.isAdminRole) links.push(adminLinks);

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
