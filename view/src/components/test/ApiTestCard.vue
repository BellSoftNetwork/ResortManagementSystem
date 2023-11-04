<template>
  <v-container>
    <v-card
      title="API Test Result"
      color="primary"
      :subtitle="env.applicationFullName"
      :text="env.version"
      variant="outlined"
    ></v-card>
  </v-container>
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
