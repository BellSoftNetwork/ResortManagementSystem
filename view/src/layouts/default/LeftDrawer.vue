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
import { useMenuStore } from "stores/menu";

defineExpose({
  toggleLeftDrawer,
});

const menuStore = useMenuStore();

const leftDrawerOpen = ref(false);
const allLinks = menuStore.allLinks;

function toggleLeftDrawer() {
  leftDrawerOpen.value = !leftDrawerOpen.value;
}
</script>
