<template>
  <q-dialog v-model="isOpen" persistent>
    <q-card style="min-width: 400px">
      <q-card-section>
        <div class="text-h6">Downloading File</div>
      </q-card-section>

      <q-card-section>
        <div class="text-body2 q-mb-sm">{{ fileName }}</div>
        <q-linear-progress :value="progress.percentage / 100" color="primary" size="20px" />
        <div class="text-caption text-center q-mt-xs">{{ progress.percentage }}%</div>
        <div class="text-caption text-center">{{ statusText }}</div>
      </q-card-section>

      <q-card-actions align="right" v-if="progress.status === 'complete' || progress.status === 'error'">
        <q-btn flat label="Close" color="primary" v-close-popup />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script setup lang="ts">
import { ref, computed } from "vue";
import { DownloadProgress } from "lib/file-download";

const isOpen = ref(false);
const fileName = ref("");
const progress = ref<DownloadProgress>({
  bytesDownloaded: 0,
  totalBytes: 0,
  percentage: 0,
  status: "downloading",
});

const statusText = computed(() => {
  switch (progress.value.status) {
    case "downloading":
      return "Downloading...";
    case "decrypting":
      return "Decrypting...";
    case "complete":
      return "Download complete!";
    case "error":
      return "Download failed";
    default:
      return "";
  }
});

function show(name: string) {
  fileName.value = name;
  progress.value = {
    bytesDownloaded: 0,
    totalBytes: 0,
    percentage: 0,
    status: "downloading",
  };
  isOpen.value = true;
}

function updateProgress(newProgress: DownloadProgress) {
  progress.value = newProgress;
}

function hide() {
  isOpen.value = false;
}

defineExpose({ show, updateProgress, hide });
</script>
