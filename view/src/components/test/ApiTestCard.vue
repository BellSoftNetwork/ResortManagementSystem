<template>
  <q-card bordered>
    <q-card-section>
      <div class="text-h6">API Test Result</div>
      <div class="text-subtitle1">
        <div v-if="env">
          {{ env.applicationFullName }}
        </div>
        <div v-else>Loading...</div>
      </div>
    </q-card-section>

    <q-card-section>
      <div v-if="env">
        <div>Commit SHA: {{ env.commitSha }}</div>
        <div>commit Short SHA: {{ env.commitShortSha }}</div>
        <div>Commit Title: {{ env.commitTitle }}</div>
        <div>Commit Timestamp: {{ env.commitTimestamp }}</div>
      </div>
    </q-card-section>
  </q-card>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue";
import { getServerEnv } from "src/api/v1/main";
import { ServerEnv } from "src/schema/server-config";

let env = ref<ServerEnv | null>(null);

onMounted(() => {
  getServerEnv().then((response) => (env.value = response.value));
});
</script>
