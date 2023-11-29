<template>
  <q-layout view="hHh lpR fFf" class="bg-grey-1">
    <TopHeader @toggle-left-drawer="toggleLeftDrawer" />

    <LeftDrawer ref="leftDrawerRef" />

    <q-page-container>
      <router-view />
    </q-page-container>
  </q-layout>
</template>

<script setup lang="ts">
import { ref } from "vue"
import TopHeader from "./TopHeader.vue"
import LeftDrawer from "./LeftDrawer.vue"

const leftDrawerRef = ref()

function toggleLeftDrawer() {
  leftDrawerRef.value.toggleLeftDrawer()
}
</script>
