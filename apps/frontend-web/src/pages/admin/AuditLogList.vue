<template>
  <q-page padding>
    <div class="row items-center justify-between q-mb-md">
      <div class="text-h5">감사 로그</div>
    </div>

    <AuditLogListDynamicTable />
  </q-page>
</template>

<script setup lang="ts">
import AuditLogListDynamicTable from "src/components/admin/AuditLogListDynamicTable.vue";
</script>
