<template>
  <q-card bordered>
    <q-card-section>
      <div class="text-h6">API Test Result</div>
      <div class="text-subtitle1">{{ env.applicationFullName }}</div>
    </q-card-section>

    <q-card-section>
      <div>Version: {{ env.version }}</div>
    </q-card-section>
  </q-card>
</template>

<script setup>
import { ref } from "vue"

let env = ref({
  applicationFullName: "",
  applicationShortName: "",
  version: "",
})

fetch("/api/v1/env")
  .then((response) => response.json())
  .then((data) => (env.value = data))
</script>
