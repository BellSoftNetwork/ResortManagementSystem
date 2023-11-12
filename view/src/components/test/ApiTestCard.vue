<template>
  <q-card bordered>
    <q-card-section>
      <div class="text-h6">API Test Result</div>
      <div class="text-subtitle1">{{ env.applicationFullName }}</div>
    </q-card-section>

    <q-card-section>
      <div>Commit SHA: {{ env.commitSha }}</div>
      <div>commit Short SHA: {{ env.commitShortSha }}</div>
      <div>Commit Title: {{ env.commitTitle }}</div>
      <div>Commit Timestamp: {{ env.commitTimestamp }}</div>
    </q-card-section>
  </q-card>
</template>

<script setup>
import { ref } from "vue"

let env = ref({
  applicationFullName: "",
  applicationShortName: "",
  commitSha: "",
  commitShortSha: "",
  commitTitle: "",
  commitTimestamp: "",
})

fetch("/api/v1/env")
  .then((response) => response.json())
  .then((data) => (env.value = data))
</script>
