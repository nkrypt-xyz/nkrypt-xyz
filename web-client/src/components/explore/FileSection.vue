<template>
  <q-card class="std-card" v-if="childFileList.length > 0">
    <q-card-section>
      <div class="text-h6">Files</div>
    </q-card-section>

    <q-list separator>
      <q-item v-for="file in childFileList" :key="file._id" clickable @click="emit('file-clicked', file)">
        <q-item-section avatar>
          <q-icon :name="getFileIcon(file)" color="secondary" size="32px" />
        </q-item-section>

        <q-item-section>
          <q-item-label>{{ file.name }}</q-item-label>
          <q-item-label caption> Created: {{ formatDate(file.metaData?.createdAt) }} </q-item-label>
        </q-item-section>

        <q-item-section side>
          <q-btn icon="more_vert" flat round dense @click.stop>
            <q-menu>
              <q-list style="min-width: 150px">
                <q-item clickable v-close-popup @click="emit('download-file', file)">
                  <q-item-section avatar>
                    <q-icon name="download" />
                  </q-item-section>
                  <q-item-section>Download</q-item-section>
                </q-item>
                <q-separator />
                <q-item clickable v-close-popup @click="renameFile(file)">
                  <q-item-section avatar>
                    <q-icon name="edit" />
                  </q-item-section>
                  <q-item-section>Rename</q-item-section>
                </q-item>
                <q-item clickable v-close-popup @click="emit('initiate-move', file, false)">
                  <q-item-section avatar>
                    <q-icon name="content_cut" />
                  </q-item-section>
                  <q-item-section>Cut</q-item-section>
                </q-item>
                <q-item clickable v-close-popup @click="emit('view-properties', file)">
                  <q-item-section avatar>
                    <q-icon name="info" />
                  </q-item-section>
                  <q-item-section>Properties</q-item-section>
                </q-item>
                <q-separator />
                <q-item clickable v-close-popup @click="deleteFile(file)">
                  <q-item-section avatar>
                    <q-icon name="delete" color="negative" />
                  </q-item-section>
                  <q-item-section>Delete</q-item-section>
                </q-item>
              </q-list>
            </q-menu>
          </q-btn>
        </q-item-section>
      </q-item>
    </q-list>
  </q-card>
</template>

<script setup lang="ts">
import { useUIStore } from "stores/ui";
import { callFileRenameApi, callFileDeleteApi } from "integration/content-apis";
import { dialogService } from "services/dialog-service";
import { errorService } from "services/error-service";
import { File } from "models/common";

const props = defineProps<{
  childFileList: File[];
}>();

const emit = defineEmits<{
  "file-clicked": [file: File];
  "download-file": [file: File];
  "initiate-move": [file: File, isDirectory: boolean];
  "view-properties": [file: File];
  "refresh-needed": [];
}>();

const uiStore = useUIStore();

function formatDate(timestamp?: number): string {
  if (!timestamp) return "N/A";
  return new Date(timestamp).toLocaleDateString();
}

function getFileIcon(file: File): string {
  if (file.name.endsWith(".txt")) return "description";
  if (file.name.match(/\.(jpg|jpeg|png|gif|webp)$/i)) return "image";
  if (file.name.match(/\.(pdf)$/i)) return "picture_as_pdf";
  if (file.name.match(/\.(zip|tar|gz)$/i)) return "archive";
  return "insert_drive_file";
}

async function renameFile(file: File) {
  const newName = await dialogService.prompt("Rename File", "Enter new name:", file.name);

  if (!newName || newName === file.name) return;

  uiStore.incrementActiveGlobalObtrusiveTaskCount();
  try {
    await callFileRenameApi({
      bucketId: file.bucketId,
      fileId: file._id,
      name: newName,
    });
    dialogService.notify("success", "File renamed!");
    emit("refresh-needed");
  } catch (error) {
    await errorService.handleUnexpectedError(error);
  } finally {
    uiStore.decrementActiveGlobalObtrusiveTaskCount();
  }
}

async function deleteFile(file: File) {
  const confirmed = await dialogService.confirm("Delete File", `Are you sure you want to delete "${file.name}"?`);

  if (!confirmed) return;

  uiStore.incrementActiveGlobalObtrusiveTaskCount();
  try {
    await callFileDeleteApi({
      bucketId: file.bucketId,
      fileId: file._id,
    });
    dialogService.notify("success", "File deleted!");
    emit("refresh-needed");
  } catch (error) {
    await errorService.handleUnexpectedError(error);
  } finally {
    uiStore.decrementActiveGlobalObtrusiveTaskCount();
  }
}
</script>
