<template>
  <q-page class="q-pa-md">
    <div v-if="isLoading" class="text-center q-pa-xl">
      <q-spinner color="primary" size="3em" />
      <div class="q-mt-md">Loading file...</div>
    </div>

    <div v-else-if="error" class="text-center q-pa-xl text-negative">
      <q-icon name="error" size="3em" />
      <div class="q-mt-md">{{ error }}</div>
    </div>

    <div v-else-if="file">
      <div class="row items-center q-mb-md">
        <div class="text-h5 col">{{ file.name }}</div>
        <q-btn color="primary" icon="save" label="Save" :loading="isSaving" @click="saveFile" class="q-ml-auto" />
      </div>

      <q-card class="std-card">
        <q-card-section>
          <q-input v-model="content" type="textarea" filled autogrow :readonly="isSaving" placeholder="Enter text content..." style="min-height: 400px; font-family: monospace" />
        </q-card-section>
      </q-card>
    </div>
  </q-page>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useRoute } from "vue-router";
import { useContentStore } from "stores/content";
import { callBlobGetApi, callBlobSetApi, callFileGetApi, callBucketListApi, BlobApiError } from "integration/content-apis";
import { dialogService } from "services/dialog-service";
import { errorService } from "services/error-service";
import { getOrCollectPasswordForBucket } from "lib/password-provider";
import { createEncryptionKeyFromPassword, makeRandomIv, makeRandomSalt, buildCryptoHeader, unbuildCryptoHeader, convertSmallUint8ArrayToString } from "nkrypt-xyz-core-nodejs";
import { File } from "models/common";

const route = useRoute();
const contentStore = useContentStore();

const isLoading = ref(false);
const isSaving = ref(false);
const error = ref<string | null>(null);
const file = ref<File | null>(null);
const content = ref("");
const originalContent = ref("");

onMounted(async () => {
  await loadFile();
});

async function loadFile() {
  isLoading.value = true;
  error.value = null;

  try {
    const bucketId = route.params.bucketId as string;
    const fileId = route.params.fileId as string;

    // Ensure bucket list is loaded
    if (contentStore.bucketList.length === 0) {
      console.debug("Loading bucket list in text editor...");
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

    let decryptedText: string;

    try {
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

      const decoder = new TextDecoder();
      decryptedText = decoder.decode(decryptedData);
    } catch (err) {
      // New/empty files have no blob yet - treat as empty content
      if (err instanceof BlobApiError && err.code === "BLOB_NOT_FOUND") {
        decryptedText = "";
      } else {
        throw err;
      }
    }

    content.value = decryptedText;
    originalContent.value = content.value;
  } catch (err) {
    error.value = "Failed to load file";
    console.error(err);
  } finally {
    isLoading.value = false;
  }
}

async function saveFile() {
  const currentFile = file.value;
  if (!currentFile) {
    await dialogService.alert("Error", "No file loaded");
    return;
  }

  if (content.value === originalContent.value) {
    dialogService.notify("info", "No changes to save");
    return;
  }

  const bucket = contentStore.getBucketById(currentFile.bucketId);
  if (!bucket) {
    await dialogService.alert("Error", "Bucket not found");
    return;
  }
  const bucketPassword = await getOrCollectPasswordForBucket(bucket);
  if (!bucketPassword) {
    await dialogService.alert("Error", "Bucket password not found");
    return;
  }

  isSaving.value = true;
  try {
    const encoder = new TextEncoder();
    const plainData = encoder.encode(content.value);

    // Generate random IV and salt for this encryption
    const { iv } = await makeRandomIv();
    const { salt } = await makeRandomSalt();

    // Create encryption key from password and salt
    const { key: encryptionKey } = await createEncryptionKeyFromPassword(bucketPassword, salt);

    // Encrypt the data
    const encryptedData = await crypto.subtle.encrypt(
      {
        name: "AES-GCM",
        iv,
      },
      encryptionKey,
      plainData
    );

    // Build crypto header for HTTP header (not embedded in blob)
    const ivStr = convertSmallUint8ArrayToString(iv);
    const saltStr = convertSmallUint8ArrayToString(salt);
    const cryptoMetaHeader = buildCryptoHeader(ivStr, saltStr);

    const blob = new Blob([encryptedData]);

    await callBlobSetApi({
      bucketId: currentFile.bucketId,
      fileId: currentFile._id,
      blob,
      cryptoMeta: cryptoMetaHeader,
    });

    originalContent.value = content.value;
    dialogService.notify("success", "File saved successfully!");
  } catch (error) {
    await errorService.handleUnexpectedError(error);
  } finally {
    isSaving.value = false;
  }
}
</script>
