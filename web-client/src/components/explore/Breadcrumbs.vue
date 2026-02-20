<template>
  <q-breadcrumbs class="q-mb-md">
    <q-breadcrumbs-el :label="currentBucket.name" icon="folder_zip" @click="emit('breadcrumb-clicked', null)" />
    <q-breadcrumbs-el v-for="entity in entityStack.slice(1)" :key="entity.directory._id" :label="entity.directory.name" icon="folder" @click="emit('breadcrumb-clicked', entity)" />
  </q-breadcrumbs>
</template>

<script setup lang="ts">
/* eslint-disable vue/multi-word-component-names */
import { Bucket, EntityStackItem } from "models/common";

defineProps<{
  entityStack: EntityStackItem[];
  currentBucket: Bucket;
}>();

const emit = defineEmits<{
  "breadcrumb-clicked": [entity: EntityStackItem | null];
}>();
</script>

<style scoped>
:deep(.q-breadcrumbs__el) {
  cursor: pointer;
}
</style>
