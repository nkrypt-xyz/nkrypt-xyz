<template>
  <q-page padding class="page">
    <div class="row">
      <div class="col-12 col-lg-8 offset-lg-2">
        <div class="text-h5 q-mb-md">Credentials</div>
        <p class="text-body2 text-grey-7 q-mb-lg">
          Cached bucket passwords are stored locally for convenience. Remove them here if you want to be prompted for a password again when accessing a bucket.
        </p>

        <q-card class="std-card">
          <q-card-section>
            <div class="text-h6 q-mb-md">Cached Bucket Passwords</div>

            <q-list v-if="cachedCredentials.length > 0" bordered separator>
              <q-item v-for="cred in cachedCredentials" :key="cred.bucketId" class="q-py-sm">
                <q-item-section avatar>
                  <q-icon name="vpn_key" color="primary" />
                </q-item-section>
                <q-item-section>
                  <q-item-label>{{ cred.displayName }}</q-item-label>
                  <q-item-label caption v-if="!cred.bucketName">Bucket ID: {{ cred.bucketId }}</q-item-label>
                </q-item-section>
                <q-item-section side>
                  <q-btn icon="delete" flat dense color="negative" round @click="deleteCredential(cred)">
                    <q-tooltip>Remove cached password</q-tooltip>
                  </q-btn>
                </q-item-section>
              </q-item>
            </q-list>

            <div v-else class="text-center text-grey-7 q-pa-lg">
              <q-icon name="check_circle" size="48px" color="positive" class="q-mb-sm" />
              <div>No cached credentials. Bucket passwords are prompted when you access encrypted content.</div>
            </div>
          </q-card-section>
        </q-card>
      </div>
    </div>
  </q-page>
</template>

<script setup lang="ts">
import { computed, onMounted } from "vue";
import { usePasswordStore } from "stores/password";
import { useContentStore } from "stores/content";
import { callBucketListApi } from "integration/content-apis";
import { dialogService } from "services/dialog-service";

const passwordStore = usePasswordStore();
const contentStore = useContentStore();

const cachedCredentials = computed(() => {
  const bucketIds = passwordStore.cachedBucketIds;
  return bucketIds.map((bucketId) => {
    const bucket = contentStore.getBucketById(bucketId);
    return {
      bucketId,
      bucketName: bucket?.name,
      displayName: bucket?.name ?? `Bucket ${bucketId.slice(0, 8)}…`,
    };
  });
});

async function deleteCredential(cred: { bucketId: string; displayName: string }) {
  const confirmed = await dialogService.confirm(
    "Remove Cached Password",
    `Remove the cached password for "${cred.displayName}"? You will be prompted for the password again when accessing this bucket.`
  );

  if (!confirmed) return;

  passwordStore.clearPasswordForBucket(cred.bucketId);
  dialogService.notify("success", "Cached password removed");
}

onMounted(async () => {
  if (contentStore.bucketList.length === 0) {
    try {
      const response = await callBucketListApi({});
      contentStore.setBucketList(response.bucketList);
    } catch {
      // Ignore - bucket names will show as IDs
    }
  }
});
</script>

<style scoped lang="scss">
.page {
  display: flex;
  flex-direction: column;
}
</style>
