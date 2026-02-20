<template>
  <q-page padding class="page">
    <div class="row">
      <div class="col-12 col-md-6 offset-md-3">
        <q-card class="std-card" v-if="isLoading">
          <q-card-section class="text-center">
            <q-spinner color="primary" size="50px" />
          </q-card-section>
        </q-card>

        <q-card class="std-card" v-if="!isLoading && bucket">
          <q-card-section>
            <div class="text-h6">Edit Bucket</div>
          </q-card-section>

          <q-card-section>
            <q-form ref="formRef" @submit="onSubmit" class="q-gutter-md">
              <q-input standout="bg-primary text-white" v-model="bucketName" label="Bucket Name" lazy-rules :rules="[(val) => !!val || 'Bucket name is required', (val) => val.length >= 3 || 'Name must be at least 3 characters']">
                <template v-slot:prepend>
                  <q-icon name="folder_zip" />
                </template>
              </q-input>

              <div class="row justify-end q-gutter-sm q-mt-lg">
                <q-btn label="Cancel" outline color="grey" @click="cancelClicked" />
                <q-btn label="Save Changes" type="submit" color="primary" unelevated :loading="isSaving" />
              </div>
            </q-form>
          </q-card-section>

          <q-separator />

          <q-card-section>
            <div class="text-h6 text-negative q-mb-md">Danger Zone</div>
            <q-btn label="Delete Bucket" color="negative" outline @click="deleteBucketClicked" />
          </q-card-section>
        </q-card>
      </div>
    </div>
  </q-page>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { QForm } from "quasar";
import { useRouter, useRoute } from "vue-router";
import { useContentStore } from "stores/content";
import { useUIStore } from "stores/ui";
import { callBucketUpdateApi, callBucketDestroyApi } from "integration/content-apis";
import { dialogService } from "services/dialog-service";
import { errorService } from "services/error-service";
import { Bucket } from "models/common";

const router = useRouter();
const route = useRoute();
const contentStore = useContentStore();
const uiStore = useUIStore();

const formRef = ref<QForm | null>(null);

const bucket = ref<Bucket | null>(null);
const bucketName = ref("");

const isLoading = ref(true);
const isSaving = ref(false);

async function loadBucket() {
  const bucketId = route.params.bucketId as string;
  const foundBucket = contentStore.getBucketById(bucketId);

  if (!foundBucket) {
    await dialogService.alert("Error", "Bucket not found");
    router.push({ name: "buckets" });
    return;
  }

  bucket.value = foundBucket;
  bucketName.value = foundBucket.name;
  isLoading.value = false;
}

async function onSubmit() {
  const isValid = await formRef.value!.validate();
  if (!isValid) return;

  isSaving.value = true;
  uiStore.incrementActiveGlobalObtrusiveTaskCount();

  try {
    await callBucketUpdateApi({
      bucketId: bucket.value!._id,
      name: bucketName.value,
    });

    contentStore.updateBucket(bucket.value!._id, { name: bucketName.value });

    dialogService.notify("success", "Bucket updated successfully!");
    router.push({ name: "buckets" });
  } catch (error) {
    await errorService.handleUnexpectedError(error);
  } finally {
    isSaving.value = false;
    uiStore.decrementActiveGlobalObtrusiveTaskCount();
  }
}

async function deleteBucketClicked() {
  const confirmed = await dialogService.confirm("Delete Bucket", `Are you sure you want to delete "${bucket.value!.name}"? This action cannot be undone and will delete all contents.`);

  if (!confirmed) return;

  uiStore.incrementActiveGlobalObtrusiveTaskCount();
  try {
    await callBucketDestroyApi({
      bucketId: bucket.value!._id,
      name: bucket.value!.name,
    });

    contentStore.removeBucket(bucket.value!._id);

    dialogService.notify("success", "Bucket deleted successfully");
    router.push({ name: "buckets" });
  } catch (error) {
    await errorService.handleUnexpectedError(error);
  } finally {
    uiStore.decrementActiveGlobalObtrusiveTaskCount();
  }
}

function cancelClicked() {
  router.back();
}

onMounted(() => {
  loadBucket();
});
</script>

<style scoped lang="scss">
.page {
  display: flex;
  flex-direction: column;
}
</style>
