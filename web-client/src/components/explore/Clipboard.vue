<template>
  <q-banner class="bg-orange text-white q-mb-md">
    <template v-slot:avatar>
      <q-icon name="content_cut" />
    </template>

    <div>
      <strong>{{ clipboard.isDirectory ? "Directory" : "File" }} cut:</strong>
      {{ clipboard.entity.name }}
    </div>

    <template v-slot:action>
      <q-btn flat label="Paste Here" @click="emit('perform-action', true)" />
      <q-btn flat label="Cancel" @click="emit('perform-action', false)" />
    </template>
  </q-banner>
</template>

<script setup lang="ts">
/* eslint-disable vue/multi-word-component-names */
import { ClipboardData } from "models/common";

defineProps<{
  clipboard: ClipboardData;
}>();

const emit = defineEmits<{
  "perform-action": [affirmative: boolean];
}>();
</script>
