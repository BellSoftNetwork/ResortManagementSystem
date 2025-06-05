<template>
  <q-tabs v-model="tab" inline-label>
    <q-route-tab
      v-for="link in allLinks"
      :key="link.text"
      :to="{ name: link.to }"
      :name="link.text"
      :icon="link.icon"
      :label="link.text"
    />
  </q-tabs>
</template>

<script setup lang="ts">
import { useMenuStore } from "stores/menu";
import { ref } from "vue";

const menuStore = useMenuStore();

const allLinks = menuStore.tabLinks;
const tab = ref();
</script>
