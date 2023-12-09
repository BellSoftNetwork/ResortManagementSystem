<template>
  <q-drawer v-model="leftDrawerOpen" show-if-above bordered class="bg-white" :width="280">
    <q-scroll-area class="fit">
      <q-list padding class="text-grey-8">
        <div v-for="links in allLinks" :key="links">
          <q-item
            v-for="link in links"
            :key="link.text"
            :to="{ name: link.to }"
            class="GNL__drawer-item"
            v-ripple
            clickable
          >
            <q-item-section avatar>
              <q-icon :name="link.icon" />
            </q-item-section>
            <q-item-section>
              <q-item-label>{{ link.text }}</q-item-label>
            </q-item-section>
          </q-item>

          <q-separator inset class="q-my-sm" />
        </div>
      </q-list>
    </q-scroll-area>
  </q-drawer>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { useAuthStore } from "stores/auth";
import { fasBook, fasCommentDollar, fasPersonShelter, fasTableColumns } from "@quasar/extras/fontawesome-v6";

defineExpose({
  toggleLeftDrawer,
});

const authStore = useAuthStore();

const leftDrawerOpen = ref(false);

const normalLinks = [{ icon: fasTableColumns, text: "대시보드", to: "Home" }];
const adminLinks = [
  { icon: fasBook, text: "예약", to: "Reservations" },
  { icon: fasPersonShelter, text: "객실", to: "Rooms" },
  { icon: fasCommentDollar, text: "결제 수단", to: "PaymentMethods" },
  { icon: "person", text: "계정 관리", to: "AdminAccounts" },
];

const allLinks = [normalLinks];
if (authStore.isAdminRole) allLinks.push(adminLinks);

function toggleLeftDrawer() {
  leftDrawerOpen.value = !leftDrawerOpen.value;
}
</script>
