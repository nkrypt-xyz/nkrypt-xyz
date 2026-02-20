<template>
  <q-dialog v-model="isOpen" persistent>
    <q-card style="min-width: 400px">
      <q-card-section>
        <div class="text-h6">{{ isDirectory ? "Directory" : "File" }} Properties</div>
      </q-card-section>

      <q-card-section>
        <q-markup-table class="base-table">
          <tbody>
            <tr>
              <td class="text-bold">Name</td>
              <td>{{ entity?.name }}</td>
            </tr>
            <tr>
              <td class="text-bold">ID</td>
              <td class="text-caption">{{ entity?._id }}</td>
            </tr>
            <tr v-if="!isDirectory">
              <td class="text-bold">Blob ID</td>
              <td class="text-caption">{{ (entity as any)?.blobId }}</td>
            </tr>
            <tr v-if="entity?.metaData?.createdAt">
              <td class="text-bold">Created</td>
              <td>{{ formatDate(entity.metaData.createdAt) }}</td>
            </tr>
            <tr v-if="decryptedMetaData">
              <td class="text-bold">Content Type</td>
              <td>{{ decryptedMetaData.contentType || "N/A" }}</td>
            </tr>
          </tbody>
        </q-markup-table>
      </q-card-section>

      <q-card-actions align="right">
        <q-btn flat label="Close" color="primary" v-close-popup />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { decryptToObject } from "utils/crypto-utils";
import { Directory, File } from "models/common";

const isOpen = ref(false);
const entity = ref<Directory | File | null>(null);
const isDirectory = ref(false);
const bucketPassword = ref<string | null>(null);
const decryptedMetaData = ref<any>(null);

function formatDate(timestamp: number): string {
  return new Date(timestamp).toLocaleString();
}

async function show(options: { entity: Directory | File; isDirectory: boolean; bucketPassword: string | null }) {
  entity.value = options.entity;
  isDirectory.value = options.isDirectory;
  bucketPassword.value = options.bucketPassword;
  decryptedMetaData.value = null;

  if (options.bucketPassword && options.entity.encryptedMetaData) {
    try {
      decryptedMetaData.value = await decryptToObject(options.entity.encryptedMetaData, options.bucketPassword);
    } catch (error) {
      console.error("Failed to decrypt metadata:", error);
    }
  }

  isOpen.value = true;
}

defineExpose({ show });
</script>

<style scoped lang="scss">
.base-table {
  td {
    padding: 12px;
  }
}
</style>
