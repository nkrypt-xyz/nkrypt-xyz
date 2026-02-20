<template>
  <q-card class="std-card q-mb-md" v-if="childDirectoryList.length > 0">
    <q-card-section>
      <div class="text-h6">Directories</div>
    </q-card-section>

    <q-list separator>
      <q-item v-for="directory in childDirectoryList" :key="directory._id" clickable @click="emit('directory-clicked', directory)">
        <q-item-section avatar>
          <q-icon name="folder" color="primary" size="32px" />
        </q-item-section>

        <q-item-section>
          <q-item-label>{{ directory.name }}</q-item-label>
          <q-item-label caption> Created: {{ formatDate(directory.metaData?.createdAt) }} </q-item-label>
        </q-item-section>

        <q-item-section side>
          <q-btn icon="more_vert" flat round dense @click.stop>
            <q-menu>
              <q-list style="min-width: 150px">
                <q-item clickable v-close-popup @click="renameDirectory(directory)">
                  <q-item-section avatar>
                    <q-icon name="edit" />
                  </q-item-section>
                  <q-item-section>Rename</q-item-section>
                </q-item>
                <q-item clickable v-close-popup @click="emit('initiate-move', directory, true)">
                  <q-item-section avatar>
                    <q-icon name="content_cut" />
                  </q-item-section>
                  <q-item-section>Cut</q-item-section>
                </q-item>
                <q-item clickable v-close-popup @click="emit('view-properties', directory)">
                  <q-item-section avatar>
                    <q-icon name="info" />
                  </q-item-section>
                  <q-item-section>Properties</q-item-section>
                </q-item>
                <q-separator />
                <q-item clickable v-close-popup @click="deleteDirectory(directory)">
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
import { callDirectoryRenameApi, callDirectoryDeleteApi } from "integration/content-apis";
import { dialogService } from "services/dialog-service";
import { errorService } from "services/error-service";
import { Directory } from "models/common";

const props = defineProps<{
  childDirectoryList: Directory[];
}>();

const emit = defineEmits<{
  "directory-clicked": [directory: Directory];
  "initiate-move": [directory: Directory, isDirectory: boolean];
  "view-properties": [directory: Directory];
  "refresh-needed": [];
}>();

const uiStore = useUIStore();

function formatDate(timestamp?: number): string {
  if (!timestamp) return "N/A";
  return new Date(timestamp).toLocaleDateString();
}

async function renameDirectory(directory: Directory) {
  const newName = await dialogService.prompt("Rename Directory", "Enter new name:", directory.name);

  if (!newName || newName === directory.name) return;

  uiStore.incrementActiveGlobalObtrusiveTaskCount();
  try {
    await callDirectoryRenameApi({
      bucketId: directory.bucketId,
      directoryId: directory._id,
      name: newName,
    });
    dialogService.notify("success", "Directory renamed!");
    emit("refresh-needed");
  } catch (error) {
    await errorService.handleUnexpectedError(error);
  } finally {
    uiStore.decrementActiveGlobalObtrusiveTaskCount();
  }
}

async function deleteDirectory(directory: Directory) {
  const confirmed = await dialogService.confirm("Delete Directory", `Are you sure you want to delete "${directory.name}"? This will delete all contents.`);

  if (!confirmed) return;

  uiStore.incrementActiveGlobalObtrusiveTaskCount();
  try {
    await callDirectoryDeleteApi({
      bucketId: directory.bucketId,
      directoryId: directory._id,
    });
    dialogService.notify("success", "Directory deleted!");
    emit("refresh-needed");
  } catch (error) {
    await errorService.handleUnexpectedError(error);
  } finally {
    uiStore.decrementActiveGlobalObtrusiveTaskCount();
  }
}
</script>
