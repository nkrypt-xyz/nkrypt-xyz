<template>
  <q-page padding class="page">
    <div class="row">
      <div class="col-12 col-md-6 offset-md-3">
        <q-card class="std-card">
          <q-card-section>
            <div class="text-h6">Create New Bucket</div>
            <div class="text-subtitle2 text-grey-7">Create an encrypted bucket to store your files securely</div>
          </q-card-section>

          <q-card-section>
            <q-form ref="formRef" @submit="onSubmit" class="q-gutter-md">
              <q-input standout="bg-primary text-white" v-model="bucketName" label="Bucket Name" lazy-rules :rules="[(val) => !!val || 'Bucket name is required', (val) => val.length >= 3 || 'Name must be at least 3 characters']">
                <template v-slot:prepend>
                  <q-icon name="folder_zip" />
                </template>
              </q-input>

              <q-input type="password" standout="bg-primary text-white" v-model="bucketPassword" label="Encryption Password" hint="This password will encrypt all content in the bucket" lazy-rules :rules="[(val) => !!val || 'Password is required', (val) => val.length >= 8 || 'Password must be at least 8 characters']">
                <template v-slot:prepend>
                  <q-icon name="lock" />
                </template>
              </q-input>

              <q-input type="password" standout="bg-primary text-white" v-model="confirmPassword" label="Confirm Password" lazy-rules :rules="[(val) => !!val || 'Please confirm password', (val) => val === bucketPassword || 'Passwords do not match']">
                <template v-slot:prepend>
                  <q-icon name="lock" />
                </template>
              </q-input>

              <div class="row justify-end q-gutter-sm q-mt-lg">
                <q-btn label="Cancel" outline color="grey" @click="cancelClicked" />
                <q-btn label="Create Bucket" type="submit" color="primary" unelevated :loading="isCreating" />
              </div>
            </q-form>
          </q-card-section>
        </q-card>
      </div>
    </div>
  </q-page>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { QForm } from "quasar";
import { useRouter } from "vue-router";
import { useUIStore } from "stores/ui";
import { callBucketCreateApi } from "integration/content-apis";
import { dialogService } from "services/dialog-service";
import { errorService } from "services/error-service";
import { generateCryptoSpec, encryptCryptoData } from "nkrypt-xyz-core-nodejs";
import { MetaDataConstant } from "constants/meta-data-constants";
import { MiscConstant } from "constants/misc-constants";
import { CommonConstant } from "constants/common-constants";

const router = useRouter();
const uiStore = useUIStore();

const formRef = ref<QForm | null>(null);

const bucketName = ref("");
const bucketPassword = ref("");
const confirmPassword = ref("");

const isCreating = ref(false);

async function onSubmit() {
  const isValid = await formRef.value!.validate();
  if (!isValid) return;

  isCreating.value = true;
  uiStore.incrementActiveGlobalObtrusiveTaskCount();

  try {
    const cryptSpec = await generateCryptoSpec();
    const cryptData = await encryptCryptoData(bucketPassword.value, cryptSpec);

    await callBucketCreateApi({
      name: bucketName.value,
      cryptSpec: JSON.stringify(cryptSpec),
      cryptData: cryptData,
      metaData: {
        [MetaDataConstant.ORIGIN_GROUP_NAME]: {
          [MetaDataConstant.ORIGIN.CLIENT_NAME]: CommonConstant.CLIENT_NAME,
          [MetaDataConstant.ORIGIN.ORIGINATION_SOURCE]: MiscConstant.ORIGINATION_SOURCE_CREATE_BUCKET,
          [MetaDataConstant.ORIGIN.ORIGINATION_DATE]: Date.now(),
        },
      },
    });

    dialogService.notify("success", "Bucket created successfully!");
    router.push({ name: "buckets" });
  } catch (error) {
    await errorService.handleUnexpectedError(error);
  } finally {
    isCreating.value = false;
    uiStore.decrementActiveGlobalObtrusiveTaskCount();
  }
}

function cancelClicked() {
  router.back();
}
</script>

<style scoped lang="scss">
.page {
  display: flex;
  flex-direction: column;
}
</style>
