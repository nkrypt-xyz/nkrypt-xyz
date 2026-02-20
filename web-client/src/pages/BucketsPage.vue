<template>
  <q-page padding class="page">
    <div class="row">
      <div class="col-12 col-lg-8 offset-lg-2">
        <div class="flex items-center justify-between q-mb-md">
          <div class="text-h5">Buckets</div>
          <q-btn v-if="canCreateBucket()" icon="add" label="Create Bucket" color="primary" @click="createBucketClicked" />
        </div>

        <q-card class="std-card" v-if="isLoading">
          <q-card-section class="text-center">
            <q-spinner color="primary" size="50px" />
          </q-card-section>
        </q-card>

        <q-card class="std-card" v-if="!isLoading">
          <q-markup-table class="base-table">
            <thead>
              <tr class="bg-primary text-white">
                <th class="text-left">Name</th>
                <th class="text-left">Created</th>
                <th class="text-right">Actions</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="bucket in contentStore.bucketList" :key="bucket._id">
                <td>{{ bucket.name }}</td>
                <td>{{ formatDate(bucket.metaData?.createdAt) }}</td>
                <td class="text-right">
                  <q-btn icon="folder_open" flat dense @click="exploreBucket(bucket)">
                    <q-tooltip>Explore</q-tooltip>
                  </q-btn>
                  <q-btn icon="edit" flat dense @click="editBucket(bucket)">
                    <q-tooltip>Edit</q-tooltip>
                  </q-btn>
                  <q-btn icon="delete" flat dense color="negative" @click="deleteBucket(bucket)">
                    <q-tooltip>Delete</q-tooltip>
                  </q-btn>
                </td>
              </tr>
            </tbody>
          </q-markup-table>

          <q-card-section v-if="contentStore.bucketList.length === 0">
            <div class="text-center text-grey-7">No buckets yet. Create one to get started!</div>
          </q-card-section>
        </q-card>
      </div>
    </div>
  </q-page>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useRouter } from "vue-router";
import { useContentStore } from "stores/content";
import { useUIStore } from "stores/ui";
import { useUserStore } from "stores/user";
import { callBucketListApi, callBucketDestroyApi } from "integration/content-apis";
import { GLOBAL_PERMISSION } from "constants/permissions";
import { dialogService } from "services/dialog-service";
import { errorService } from "services/error-service";
import { Bucket } from "models/common";

const router = useRouter();
const contentStore = useContentStore();
const uiStore = useUIStore();
const userStore = useUserStore();

const isLoading = ref(true);
const canCreateBucket = () => userStore.hasPermission(GLOBAL_PERMISSION.CREATE_BUCKET);

function formatDate(timestamp?: number): string {
  if (!timestamp) return "N/A";
  return new Date(timestamp).toLocaleDateString();
}

async function loadBuckets() {
  try {
    isLoading.value = true;
    const response = await callBucketListApi({});
    contentStore.setBucketList(response.bucketList);
  } catch (error) {
    await errorService.handleUnexpectedError(error);
  } finally {
    isLoading.value = false;
  }
}

function createBucketClicked() {
  router.push({ name: "bucket-create" });
}

function exploreBucket(bucket: Bucket) {
  router.push({ name: "explore", params: { bucketId: bucket._id, pathMatch: [] } });
}

function editBucket(bucket: Bucket) {
  router.push({ name: "bucket-edit", params: { bucketId: bucket._id } });
}

async function deleteBucket(bucket: Bucket) {
  const confirmed = await dialogService.confirm("Delete Bucket", `Are you sure you want to delete "${bucket.name}"? This action cannot be undone.`);

  if (!confirmed) return;

  uiStore.incrementActiveGlobalObtrusiveTaskCount();
  try {
    await callBucketDestroyApi({
      bucketId: bucket._id,
      name: bucket.name,
    });
    dialogService.notify("success", "Bucket deleted successfully");
    await loadBuckets();
  } catch (error) {
    await errorService.handleUnexpectedError(error);
  } finally {
    uiStore.decrementActiveGlobalObtrusiveTaskCount();
  }
}

onMounted(() => {
  loadBuckets();
});
</script>

<style scoped lang="scss">
.page {
  display: flex;
  flex-direction: column;
}

.base-table {
  th {
    font-weight: 600;
  }

  td {
    padding: 12px;
  }
}
</style>
