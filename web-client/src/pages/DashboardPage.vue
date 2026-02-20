<template>
  <q-page padding class="page">
    <div class="row">
      <div class="col-12 col-md-8 offset-md-2">
        <q-card class="std-card q-mb-md">
          <q-card-section>
            <div class="text-h5">Welcome, {{ userStore.displayName }}</div>
          </q-card-section>
        </q-card>

        <q-card class="std-card" v-if="isLoading">
          <q-card-section class="text-center">
            <q-spinner color="primary" size="50px" />
            <div class="q-mt-md">Loading metrics...</div>
          </q-card-section>
        </q-card>

        <q-card class="std-card" v-if="!isLoading && metrics">
          <q-card-section>
            <div class="text-h6">Disk Usage</div>
          </q-card-section>

          <q-markup-table class="base-table">
            <tbody>
              <tr>
                <td class="text-bold">Used</td>
                <td>
                  {{ formatSize(metrics.disk.usedBytes) }}
                  ({{ formatPercentage(metrics.disk.usedBytes, metrics.disk.totalBytes) }})
                </td>
              </tr>
              <tr>
                <td class="text-bold">Total</td>
                <td>{{ formatSize(metrics.disk.totalBytes) }}</td>
              </tr>
            </tbody>
          </q-markup-table>
        </q-card>
      </div>
    </div>
  </q-page>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useUserStore } from "stores/user";
import { callMetricsGetSummaryApi } from "integration/metrics-apis";
import { errorService } from "services/error-service";

const userStore = useUserStore();

const isLoading = ref(true);
const metrics = ref<any>(null);

function formatSize(bytes: number): string {
  const gb = Math.round((bytes / 1000 / 1000 / 1000) * 100) / 100;
  return `${gb} GB`;
}

function formatPercentage(value: number, total: number): string {
  const percentage = Math.round((value / total) * 100 * 100) / 100;
  return `${percentage}%`;
}

async function loadMetrics() {
  try {
    isLoading.value = true;
    const response = await callMetricsGetSummaryApi({});
    metrics.value = response;
  } catch (error) {
    await errorService.handleUnexpectedError(error);
  } finally {
    isLoading.value = false;
  }
}

onMounted(() => {
  loadMetrics();
});
</script>

<style scoped lang="scss">
.page {
  display: flex;
  flex-direction: column;
}

.base-table {
  td {
    padding: 12px;
  }
}
</style>
