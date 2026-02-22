<template>
  <q-page class="q-pa-md">
    <div v-if="isLoading" class="text-center q-pa-xl">
      <q-spinner color="primary" size="3em" />
      <div class="q-mt-md">Loading image...</div>
    </div>

    <div v-else-if="error" class="text-center q-pa-xl text-negative">
      <q-icon name="error" size="3em" />
      <div class="q-mt-md">{{ error }}</div>
    </div>

    <div v-else-if="file && imageUrl">
      <div class="row items-center q-mb-md">
        <div class="text-h5 col">{{ file.name }}</div>
        <q-btn color="primary" icon="download" label="Download" @click="downloadImage" class="q-ml-auto" />
      </div>

      <q-card class="std-card">
        <q-card-section>
          <div class="text-center">
            <img :src="imageUrl" :alt="file.name" style="max-width: 100%; height: auto" />
          </div>
        </q-card-section>
      </q-card>
    </div>
  </q-page>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from "vue";
import { useRoute } from "vue-router";
import { useContentStore } from "stores/content";
import { callBlobGetApi, callFileGetApi, callBucketListApi } from "integration/content-apis";
import { errorService } from "services/error-service";
import { getOrCollectPasswordForBucket } from "lib/password-provider";
import { downloadFileWithDecryption } from "lib/file-download";
import { createEncryptionKeyFromPassword, unbuildCryptoHeader } from "nkrypt-xyz-core-nodejs";
import { File } from "models/common";

const route = useRoute();
const contentStore = useContentStore();

const isLoading = ref(false);
const error = ref<string | null>(null);
const file = ref<File | null>(null);
const imageUrl = ref<string | null>(null);

onMounted(async () => {
  await loadImage();
});

onBeforeUnmount(() => {
  if (imageUrl.value) {
    URL.revokeObjectURL(imageUrl.value);
  }
});

async function loadImage() {
  isLoading.value = true;
  error.value = null;

  try {
    const bucketId = route.params.bucketId as string;
    const fileId = route.params.fileId as string;

    // Ensure bucket list is loaded
    if (contentStore.bucketList.length === 0) {
      console.debug("Loading bucket list in image viewer...");
      const response = await callBucketListApi({});
      contentStore.setBucketList(response.bucketList);
    }

    const bucket = contentStore.bucketList.find((b) => b._id === bucketId);
    if (!bucket) {
      error.value = "Bucket not found";
      return;
    }

    const bucketPassword = await getOrCollectPasswordForBucket(bucket);
    if (!bucketPassword) {
      error.value = "Password required to access file";
      return;
    }

    const response = await callFileGetApi({ bucketId, fileId });
    const fileData = response.file;
    if (!fileData) {
      throw new Error("Invalid file response from server");
    }
    file.value = fileData;

    const { blob, cryptoMeta } = await callBlobGetApi({
      bucketId: fileData.bucketId,
      fileId: fileData._id,
    });

    if (!cryptoMeta) {
      throw new Error("Missing crypto metadata from server");
    }

    // Parse crypto header to get IV and salt
    const [ivBuffer, saltBuffer] = unbuildCryptoHeader(cryptoMeta);
    const iv = new Uint8Array(ivBuffer);
    const salt = new Uint8Array(saltBuffer);

    // Create encryption key from password and salt
    const { key: encryptionKey } = await createEncryptionKeyFromPassword(bucketPassword, salt);

    // Decrypt the blob data (no embedded header, just encrypted content)
    const arrayBuffer = await blob.arrayBuffer();
    const decryptedData = await crypto.subtle.decrypt(
      {
        name: "AES-GCM",
        iv,
      },
      encryptionKey,
      arrayBuffer
    );

    const imageBlob = new Blob([decryptedData], { type: "image/*" });
    imageUrl.value = URL.createObjectURL(imageBlob);
  } catch (err) {
    error.value = "Failed to load image";
    console.error(err);
  } finally {
    isLoading.value = false;
  }
}

async function downloadImage() {
  const currentFile = file.value;
  if (!currentFile) {
    await errorService.handleUnexpectedError(new Error("No file loaded"));
    return;
  }

  const bucket = contentStore.getBucketById(currentFile.bucketId);
  if (!bucket) {
    await errorService.handleUnexpectedError(new Error("Bucket not found"));
    return;
  }
  const bucketPassword = await getOrCollectPasswordForBucket(bucket);
  if (!bucketPassword) {
    await errorService.handleUnexpectedError(new Error("Bucket password not found"));
    return;
  }

  try {
    await downloadFileWithDecryption(currentFile, bucketPassword);
  } catch (error) {
    await errorService.handleUnexpectedError(error);
  }
}
</script>
