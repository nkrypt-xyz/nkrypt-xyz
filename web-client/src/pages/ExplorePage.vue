<template>
  <q-page class="q-pa-md">
    <div v-if="isLoading" class="text-center q-pa-xl">
      <q-spinner color="primary" size="3em" />
      <div class="q-mt-md">Loading...</div>
    </div>

    <div v-else-if="error" class="text-center q-pa-xl text-negative">
      <q-icon name="error" size="3em" />
      <div class="q-mt-md">{{ error }}</div>
    </div>

    <div v-else-if="currentBucket && entityStack.length > 0">
      <!-- Breadcrumbs -->
      <Breadcrumbs :entity-stack="entityStack" :current-bucket="currentBucket" @breadcrumb-clicked="handleBreadcrumbClick" />

      <!-- Clipboard Banner -->
      <Clipboard v-if="clipboard" :clipboard="clipboard" @perform-action="handleClipboardAction" />

      <!-- Action Buttons -->
      <div class="q-mb-md row q-gutter-sm">
        <q-btn color="primary" icon="create_new_folder" label="New Directory" @click="createDirectory" />
        <q-btn color="primary" icon="note_add" label="New Text File" @click="createNewTextFile" />
        <q-btn color="primary" icon="upload" label="Upload File" @click="uploadFile" />
      </div>

      <!-- Directory Section -->
      <DirectorySection :child-directory-list="childDirectoryList" @directory-clicked="handleDirectoryClick" @initiate-move="initiateMoveOperation" @view-properties="viewProperties" @refresh-needed="loadDirectoryContents" />

      <!-- File Section -->
      <FileSection :child-file-list="childFileList" @file-clicked="handleFileClick" @download-file="downloadFile" @initiate-move="initiateMoveOperation" @view-properties="viewProperties" @refresh-needed="loadDirectoryContents" />

      <!-- Empty State -->
      <q-card v-if="childDirectoryList.length === 0 && childFileList.length === 0" class="std-card text-center q-pa-xl">
        <q-icon name="inbox" size="4em" color="grey-5" />
        <div class="text-h6 text-grey-7 q-mt-md">This directory is empty</div>
        <div class="text-grey-6">Create a new directory or upload a file to get started.</div>
      </q-card>
    </div>

    <!-- Properties Modal -->
    <PropertiesModal ref="propertiesModalRef" />

    <!-- Upload Progress Dialog -->
    <UploadProgressDialog ref="uploadProgressDialogRef" />

    <!-- Download Progress Dialog -->
    <DownloadProgressDialog ref="downloadProgressDialogRef" />

    <!-- Hidden file input for upload -->
    <input ref="fileInputRef" type="file" style="display: none" @change="handleFileSelected" />
  </q-page>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useContentStore } from "stores/content";
import { usePasswordStore } from "stores/password";
import { useUIStore } from "stores/ui";
import { callDirectoryGetApi, callDirectoryCreateApi, callDirectoryMoveApi, callFileMoveApi, callBucketListApi } from "integration/content-apis";
import { dialogService } from "services/dialog-service";
import { errorService } from "services/error-service";
import { navigationService } from "services/navigation-service";
import { getOrCollectPasswordForBucket } from "lib/password-provider";
import { createTextFile } from "lib/create-text-file";
import { ClipboardAction } from "lib/clipboard-helper";
import { uploadFileWithEncryption } from "lib/file-upload";
import { downloadFileWithDecryption } from "lib/file-download";
import { MetaDataConstant } from "constants/meta-data-constants";
import { MiscConstant } from "constants/misc-constants";
import { CommonConstant } from "constants/common-constants";
import Breadcrumbs from "components/explore/Breadcrumbs.vue";
import Clipboard from "components/explore/Clipboard.vue";
import DirectorySection from "components/explore/DirectorySection.vue";
import FileSection from "components/explore/FileSection.vue";
import PropertiesModal from "components/explore/PropertiesModal.vue";
import UploadProgressDialog from "components/explore/UploadProgressDialog.vue";
import DownloadProgressDialog from "components/explore/DownloadProgressDialog.vue";
import { Bucket, Directory, File, EntityStackItem, ClipboardData } from "models/common";

const route = useRoute();
const router = useRouter();
const contentStore = useContentStore();
const passwordStore = usePasswordStore();
const uiStore = useUIStore();

const isLoading = ref(false);
const error = ref<string | null>(null);
const entityStack = ref<EntityStackItem[]>([]);
const childDirectoryList = ref<Directory[]>([]);
const childFileList = ref<File[]>([]);
const clipboard = ref<ClipboardData | null>(null);
const propertiesModalRef = ref<InstanceType<typeof PropertiesModal> | null>(null);
const uploadProgressDialogRef = ref<InstanceType<typeof UploadProgressDialog> | null>(null);
const downloadProgressDialogRef = ref<InstanceType<typeof DownloadProgressDialog> | null>(null);
const fileInputRef = ref<HTMLInputElement | null>(null);

const currentBucket = computed(() => {
  if (entityStack.value.length === 0) return null;
  return entityStack.value[0].bucket;
});

const currentDirectory = computed(() => {
  if (entityStack.value.length === 0) return null;
  return entityStack.value[entityStack.value.length - 1].directory;
});

onMounted(async () => {
  // Ensure bucket list is loaded
  if (contentStore.bucketList.length === 0) {
    console.debug("Loading bucket list...");
    try {
      const response = await callBucketListApi({});
      contentStore.setBucketList(response.bucketList);
    } catch (err) {
      console.error("Failed to load bucket list:", err);
      error.value = "Failed to load buckets";
      isLoading.value = false;
      return;
    }
  }
  await initializePage();
});

watch(
  () => [route.params.bucketId, route.params.pathMatch, contentStore.bucketList.length],
  async () => {
    // Only initialize if we have buckets loaded
    if (contentStore.bucketList.length > 0) {
      await initializePage();
    }
  }
);

async function initializePage() {
  isLoading.value = true;
  error.value = null;

  try {
    const bucketId = route.params.bucketId as string;
    const pathMatch = route.params.pathMatch as string[] | undefined;
    
    // Filter out empty strings from pathMatch
    let pathParts: string[] = [];
    if (pathMatch && Array.isArray(pathMatch)) {
      pathParts = pathMatch.filter((part) => part && part.length > 0);
    }

    console.log("=== EXPLORE PAGE INIT ===");
    console.log("bucketId:", bucketId);
    console.log("pathParts:", pathParts);
    console.log("bucketList length:", contentStore.bucketList.length);

    const bucket = contentStore.bucketList.find((b) => b._id === bucketId);
    if (!bucket) {
      console.error("❌ Bucket not found in store");
      console.error("Available buckets:", contentStore.bucketList);
      error.value = "Bucket not found";
      return;
    }

    console.log("Found bucket:", bucket._id);
    console.log("Bucket name:", bucket.name);
    console.log("Bucket rootDirectoryId:", bucket.rootDirectoryId);
    console.log("Full bucket object:", bucket);

    // Check if rootDirectoryId exists
    if (!bucket.rootDirectoryId) {
      console.error("CRITICAL: Bucket has no rootDirectoryId!", {
        bucket,
        keys: Object.keys(bucket),
      });
      error.value = "Bucket configuration error: missing rootDirectoryId";
      return;
    }

    // Don't prompt for password yet - only when needed for encryption/decryption
    // Password will be requested lazily when user creates/uploads/downloads files

    await buildEntityStackFromPath(bucket, pathParts);
    await loadDirectoryContents();
  } catch (err) {
    error.value = "Failed to load bucket contents";
    console.error(err);
  } finally {
    isLoading.value = false;
  }
}

async function buildEntityStackFromPath(bucket: Bucket, pathParts: string[]) {
  entityStack.value = [];

  // Validate bucket has rootDirectoryId
  if (!bucket.rootDirectoryId) {
    console.error("Bucket missing rootDirectoryId. Bucket data:", JSON.stringify(bucket, null, 2));
    throw new Error("Invalid bucket configuration: missing rootDirectoryId");
  }

  console.log("=== LOADING ROOT DIRECTORY ===");
  console.log("bucketId:", bucket._id);
  console.log("rootDirectoryId:", bucket.rootDirectoryId);
  console.log("Calling API with:", { bucketId: bucket._id, directoryId: bucket.rootDirectoryId });

  // Load root directory
  const rootResponse = await callDirectoryGetApi({
    bucketId: bucket._id,
    directoryId: bucket.rootDirectoryId,
  });
  
  console.log("Root directory loaded successfully");

  const rootDirectory: Directory = {
    _id: rootResponse.directory._id,
    name: bucket.name,
    bucketId: bucket._id,
    parentDirectoryId: null,
    metaData: rootResponse.directory.metaData,
    encryptedMetaData: rootResponse.directory.encryptedMetaData,
  };

  entityStack.value.push({ bucket, directory: rootDirectory });

  // Filter out empty path parts and traverse
  const validPathParts = pathParts.filter((part) => part && part.length > 0);
  
  for (let i = 0; i < validPathParts.length; i++) {
    const directoryId = validPathParts[i];

    // Skip if directoryId is invalid
    if (!directoryId || directoryId.length !== 16) {
      console.warn(`Invalid directory ID in path: "${directoryId}"`);
      break;
    }

    try {
      const response = await callDirectoryGetApi({
        bucketId: bucket._id,
        directoryId,
      });

      const directory: Directory = {
        _id: response.directory._id,
        name: response.directory.name,
        bucketId: bucket._id,
        parentDirectoryId: response.directory.parentDirectoryId,
        metaData: response.directory.metaData,
        encryptedMetaData: response.directory.encryptedMetaData,
      };

      entityStack.value.push({ bucket, directory });
    } catch (err) {
      console.error(`Failed to load directory ${directoryId}:`, err);
      // Invalid path - stop here
      break;
    }
  }
}

async function loadDirectoryContents() {
  if (!currentBucket.value || !currentDirectory.value) return;

  isLoading.value = true;
  try {
    const response = await callDirectoryGetApi({
      bucketId: currentBucket.value._id,
      directoryId: currentDirectory.value._id,
    });

    childDirectoryList.value = response.childDirectoryList;
    childFileList.value = response.childFileList;
  } catch (err) {
    await errorService.handleUnexpectedError(err);
  } finally {
    isLoading.value = false;
  }
}

function handleBreadcrumbClick(entity: EntityStackItem | null) {
  const index = entity === null ? 0 : entityStack.value.findIndex((e) => e.directory._id === entity.directory._id);
  
  if (index === -1) return;

  // Build path from entity stack up to clicked item
  const pathParts = entityStack.value
    .slice(1, index + 1) // Skip root, go up to clicked item
    .map((item) => item.directory._id);

  // Navigate to path
  router.push({
    name: "explore",
    params: {
      bucketId: currentBucket.value!._id,
      pathMatch: pathParts,
    },
  });
}

function handleDirectoryClick(directory: Directory) {
  // Build path from current entity stack + new directory
  const pathParts = entityStack.value
    .slice(1) // Skip root
    .map((item) => item.directory._id)
    .concat([directory._id]);

  // Navigate to new path
  router.push({
    name: "explore",
    params: {
      bucketId: currentBucket.value!._id,
      pathMatch: pathParts,
    },
  });
}

function handleFileClick(file: File) {
  if (file.name.endsWith(".txt")) {
    navigationService.push(`/editor/text/${file.bucketId}/${file._id}`);
  } else if (file.name.match(/\.(jpg|jpeg|png|gif|webp)$/i)) {
    navigationService.push(`/viewer/image/${file.bucketId}/${file._id}`);
  } else {
    dialogService.notify("info", "Preview not available for this file type");
  }
}

async function createDirectory() {
  const name = await dialogService.prompt("Create Directory", "Enter directory name:");

  if (!name) return;

  uiStore.incrementActiveGlobalObtrusiveTaskCount();
  try {
    await callDirectoryCreateApi({
      bucketId: currentBucket.value!._id,
      name,
      parentDirectoryId: currentDirectory.value!._id,
      metaData: {
        [MetaDataConstant.ORIGIN_GROUP_NAME]: {
          [MetaDataConstant.ORIGIN.CLIENT_NAME]: CommonConstant.CLIENT_NAME,
          [MetaDataConstant.ORIGIN.ORIGINATION_SOURCE]: MiscConstant.ORIGINATION_SOURCE_CREATE_DIRECTORY,
          [MetaDataConstant.ORIGIN.ORIGINATION_DATE]: Date.now(),
        },
      },
      encryptedMetaData: "-",
    });
    dialogService.notify("success", "Directory created!");
    await loadDirectoryContents();
  } catch (error) {
    await errorService.handleUnexpectedError(error);
  } finally {
    uiStore.decrementActiveGlobalObtrusiveTaskCount();
  }
}

async function createNewTextFile() {
  const fileName = await dialogService.prompt("Create Text File", "Enter file name:");

  if (!fileName) return;

  // Automatically append .txt extension if not present
  const finalFileName = fileName.endsWith(".txt") ? fileName : `${fileName}.txt`;

  // Lazily request password when needed
  const bucketPassword = await getOrCollectPasswordForBucket(currentBucket.value!);
  if (!bucketPassword) {
    return;
  }

  uiStore.incrementActiveGlobalObtrusiveTaskCount();
  try {
    await createTextFile(currentBucket.value!, currentDirectory.value!, bucketPassword, finalFileName);
    dialogService.notify("success", "Text file created!");
    await loadDirectoryContents();
  } catch (error) {
    await errorService.handleUnexpectedError(error);
  } finally {
    uiStore.decrementActiveGlobalObtrusiveTaskCount();
  }
}

function uploadFile() {
  fileInputRef.value?.click();
}

async function handleFileSelected(event: Event) {
  const target = event.target as HTMLInputElement;
  const file = target.files?.[0];

  if (!file) return;

  // Lazily request password when needed
  const bucketPassword = await getOrCollectPasswordForBucket(currentBucket.value!);
  if (!bucketPassword) {
    target.value = "";
    return;
  }

  uploadProgressDialogRef.value?.show(file.name);

  try {
    await uploadFileWithEncryption(file, currentBucket.value!, currentDirectory.value!, bucketPassword, (progress) => {
      uploadProgressDialogRef.value?.updateProgress(progress);
    });

    dialogService.notify("success", "File uploaded successfully!");
    await loadDirectoryContents();
  } catch (error) {
    await errorService.handleUnexpectedError(error);
  } finally {
    uploadProgressDialogRef.value?.hide();
  }

  if (fileInputRef.value) {
    fileInputRef.value.value = "";
  }
}

async function downloadFile(file: File) {
  // Lazily request password when needed
  const bucket = contentStore.bucketList.find((b) => b._id === file.bucketId);
  if (!bucket) {
    await dialogService.alert("Error", "Bucket not found");
    return;
  }

  const bucketPassword = await getOrCollectPasswordForBucket(bucket);
  if (!bucketPassword) {
    return;
  }

  downloadProgressDialogRef.value?.show(file.name);

  try {
    await downloadFileWithDecryption(file, bucketPassword, (progress) => {
      downloadProgressDialogRef.value?.updateProgress(progress);
    });

    dialogService.notify("success", "File downloaded successfully!");
  } catch (error) {
    await errorService.handleUnexpectedError(error);
  } finally {
    downloadProgressDialogRef.value?.hide();
  }
}

function initiateMoveOperation(entity: Directory | File, isDirectory: boolean) {
  clipboard.value = {
    entity,
    isDirectory,
    action: ClipboardAction.CUT,
  };
}

async function handleClipboardAction(affirmative: boolean) {
  if (!affirmative || !clipboard.value) {
    clipboard.value = null;
    return;
  }

  uiStore.incrementActiveGlobalObtrusiveTaskCount();
  try {
    if (clipboard.value.isDirectory) {
      await callDirectoryMoveApi({
        bucketId: currentBucket.value!._id,
        directoryId: clipboard.value.entity._id,
        targetParentDirectoryId: currentDirectory.value!._id,
      });
    } else {
      await callFileMoveApi({
        bucketId: currentBucket.value!._id,
        fileId: clipboard.value.entity._id,
        targetParentDirectoryId: currentDirectory.value!._id,
      });
    }
    dialogService.notify("success", "Item moved successfully!");
    clipboard.value = null;
    await loadDirectoryContents();
  } catch (error) {
    await errorService.handleUnexpectedError(error);
  } finally {
    uiStore.decrementActiveGlobalObtrusiveTaskCount();
  }
}

async function viewProperties(entity: Directory | File) {
  // Lazily request password when needed (for viewing encrypted metadata)
  const bucketPassword = await getOrCollectPasswordForBucket(currentBucket.value!);
  if (!bucketPassword) {
    return;
  }

  const isDirectory = "parentDirectoryId" in entity && !("blobId" in entity);
  propertiesModalRef.value?.show({ entity, isDirectory, bucketPassword });
}
</script>
