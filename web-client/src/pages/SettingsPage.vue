<template>
  <q-page padding class="page">
    <div class="row">
      <div class="col-12 col-md-6 offset-md-3">
        <q-card class="std-card">
          <q-card-section>
            <div class="text-h6">Settings</div>
          </q-card-section>

          <q-card-section>
            <div class="q-gutter-md">
              <q-select
                standout="bg-primary text-white"
                v-model="settingsStore.uploadMechanism"
                :options="uploadMechanismOptions"
                label="Upload Mechanism"
                emit-value
                map-options
                @update:model-value="settingsStore.setUploadMechanism"
              >
                <template v-slot:prepend>
                  <q-icon name="upload" />
                </template>
              </q-select>

              <q-select
                standout="bg-primary text-white"
                v-model="settingsStore.downloadMechanism"
                :options="downloadMechanismOptions"
                label="Download Mechanism"
                emit-value
                map-options
                @update:model-value="settingsStore.setDownloadMechanism"
              >
                <template v-slot:prepend>
                  <q-icon name="download" />
                </template>
              </q-select>

              <q-toggle v-model="settingsStore.plainTextEditorNoRestrictions" label="Plain Text Editor - No Restrictions" @update:model-value="settingsStore.setPlainTextEditorNoRestrictions" />

              <div class="text-caption text-grey-7 q-mt-sm">When enabled, the plain text editor can open files of any size. Use with caution as this may cause performance issues with large files.</div>
            </div>
          </q-card-section>
        </q-card>
      </div>
    </div>
  </q-page>
</template>

<script setup lang="ts">
import { useSettingsStore } from "stores/settings";

const settingsStore = useSettingsStore();

const uploadMechanismOptions = [
  { label: "Basic", value: "basic" },
  { label: "Stream", value: "stream" },
  { label: "Chunked Stream", value: "chunked-stream" },
];

const downloadMechanismOptions = [
  { label: "Basic", value: "basic" },
  { label: "Stream", value: "stream" },
  { label: "FileSystem", value: "fs" },
];
</script>

<style scoped lang="scss">
.page {
  display: flex;
  flex-direction: column;
}
</style>
