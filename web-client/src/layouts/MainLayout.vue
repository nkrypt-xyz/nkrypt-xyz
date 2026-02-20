<template>
  <q-layout view="lHh Lpr lFf">
    <q-header v-if="userStore.isUserLoggedIn">
      <q-toolbar>
        <q-btn v-if="!$route.meta.backButton" flat dense round icon="menu" @click="toggleLeftDrawer" />
        <q-btn v-else flat dense round icon="arrow_back" @click="navigateBack" />

        <q-toolbar-title>
          {{ $route.meta.title || "nkrypt.xyz" }}
        </q-toolbar-title>

        <q-btn flat dense round :icon="isDarkMode ? 'light_mode' : 'dark_mode'" @click="toggleDarkMode">
          <q-tooltip>{{ isDarkMode ? "Light Mode" : "Dark Mode" }}</q-tooltip>
        </q-btn>

        <q-btn flat dense round icon="person">
          <q-menu>
            <q-list style="min-width: 150px">
              <q-item clickable v-close-popup @click="profileClicked">
                <q-item-section>Profile</q-item-section>
              </q-item>
              <q-item clickable v-close-popup @click="settingsClicked">
                <q-item-section>Settings</q-item-section>
              </q-item>
              <q-separator />
              <q-item clickable v-close-popup @click="logoutClicked">
                <q-item-section>Logout</q-item-section>
              </q-item>
            </q-list>
          </q-menu>
        </q-btn>
      </q-toolbar>
    </q-header>

    <q-drawer v-model="leftDrawerOpen" show-if-above bordered :dark="isDarkMode" v-if="userStore.isUserLoggedIn">
      <div class="drawer-content">
        <q-list>
          <q-item-label header>ENCRYPTED BUCKETS</q-item-label>

          <template v-if="contentStore.bucketList.length === 0">
            <q-item>
              <q-item-section>
                <q-item-label caption> Create a bucket to get started </q-item-label>
              </q-item-section>
            </q-item>
          </template>

          <template v-else>
            <q-item v-for="bucket in contentStore.bucketList" :key="bucket._id" clickable :active="contentStore.activeBucket?._id === bucket._id" @click="bucketClicked(bucket)">
              <q-item-section avatar>
                <q-icon name="folder_zip" />
              </q-item-section>
              <q-item-section>
                <q-item-label>{{ bucket.name }}</q-item-label>
              </q-item-section>
            </q-item>
          </template>

          <q-item clickable @click="createBucketClicked">
            <q-item-section avatar>
              <q-icon name="add" color="primary" />
            </q-item-section>
            <q-item-section>
              <q-item-label>Create a bucket</q-item-label>
            </q-item-section>
          </q-item>
        </q-list>

        <q-separator />

        <q-list>
          <q-item-label header>NAVIGATION</q-item-label>

          <q-item clickable @click="navigateTo('dashboard')">
            <q-item-section avatar>
              <q-icon name="home" />
            </q-item-section>
            <q-item-section>Dashboard</q-item-section>
          </q-item>

          <q-item clickable @click="navigateTo('buckets')">
            <q-item-section avatar>
              <q-icon name="folder_zip" />
            </q-item-section>
            <q-item-section>Buckets</q-item-section>
          </q-item>

          <q-item v-if="isAdmin" clickable @click="navigateTo('users')">
            <q-item-section avatar>
              <q-icon name="group" />
            </q-item-section>
            <q-item-section>Users</q-item-section>
          </q-item>
        </q-list>

        <div style="flex: 1"></div>

        <div class="drawer-bottom">
          <div class="drawer-bottom-content">
            <div class="logo-container">
              <img src="/src/assets/logo-512-sqr.png" alt="nkrypt.xyz" class="logo-image" />
            </div>
            <div class="logged-in-as">
              Logged in as:<br />
              {{ userStore.displayName }}
            </div>
          </div>
        </div>
      </div>
    </q-drawer>

    <q-page-container>
      <router-view />
    </q-page-container>

    <q-dialog v-model="uiStore.isGloballyLoading" persistent>
      <q-card>
        <q-card-section class="text-center q-pa-lg">
          <q-spinner color="primary" size="50px" />
          <div class="q-mt-md text-grey-7">Please wait...</div>
        </q-card-section>
      </q-card>
    </q-dialog>
  </q-layout>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import { useQuasar } from "quasar";
import { useRouter } from "vue-router";
import { useUserStore } from "stores/user";
import { useSessionStore } from "stores/session";
import { useContentStore } from "stores/content";
import { useSettingsStore } from "stores/settings";
import { useUIStore } from "stores/ui";
import { usePasswordStore } from "stores/password";
import { callBucketListApi } from "integration/content-apis";
import { dialogService } from "services/dialog-service";
import { navigationService } from "services/navigation-service";
import { Bucket } from "models/common";

const router = useRouter();
const $q = useQuasar();
const userStore = useUserStore();
const sessionStore = useSessionStore();
const contentStore = useContentStore();
const settingsStore = useSettingsStore();
const uiStore = useUIStore();
const passwordStore = usePasswordStore();

const leftDrawerOpen = ref(false);

const isDarkMode = computed(() => $q.dark.isActive);
const isAdmin = computed(() => userStore.hasPermission("MANAGE_USERS"));

function toggleLeftDrawer() {
  leftDrawerOpen.value = !leftDrawerOpen.value;
}

function toggleDarkMode() {
  const newDarkMode = !$q.dark.isActive;
  $q.dark.set(newDarkMode);
  settingsStore.setDarkMode(newDarkMode);
}

function navigateBack() {
  navigationService.navigateToPreviousPageOrDashboard();
}

async function navigateTo(name: string) {
  await router.push({ name });
}

function bucketClicked(bucket: Bucket) {
  router.push(`/explore/${bucket._id}`);
}

function createBucketClicked() {
  router.push({ name: "bucket-create" });
}

function profileClicked() {
  router.push({ name: "profile" });
}

function settingsClicked() {
  router.push({ name: "settings" });
}

async function logoutClicked() {
  const skipConfirm = typeof sessionStorage !== "undefined" && sessionStorage.getItem("e2e-skip-logout-confirm");
  const confirmed = skipConfirm || (await dialogService.confirm("Logout", "Are you sure you want to logout?"));

  if (!confirmed) return;

  if (skipConfirm) sessionStorage.removeItem("e2e-skip-logout-confirm");
  userStore.clearUser();
  sessionStore.clearSession();
  contentStore.setBucketList([]);
  contentStore.setActiveBucket(null);
  passwordStore.clearAllPasswords();

  router.push({ name: "login" });
}

async function loadBucketList() {
  if (!sessionStore.isSessionActive) return;

  try {
    const response = await callBucketListApi({});
    contentStore.setBucketList(response.bucketList);
  } catch (error) {
    console.error("Failed to load bucket list:", error);
  }
}

onMounted(() => {
  navigationService.setRouter(router);
  loadBucketList();

  if (settingsStore.darkMode !== null) {
    $q.dark.set(settingsStore.darkMode);
  }
});
</script>

<style scoped lang="scss">
.drawer-content {
  display: flex;
  flex-direction: column;
  height: 100%;
  min-height: 100%;
}

.drawer-bottom {
  padding: 16px;
  background: rgb(233, 245, 237);
  border-top: 1px solid rgba(0, 0, 0, 0.1);
}

body.body--dark .drawer-bottom {
  background: #1e2538;
  color: #cbd5e1;
  border-top: 1px solid rgba(255, 255, 255, 0.05);
}

.drawer-bottom-content {
  display: flex;
  align-items: center;
  gap: 12px;
}

.logo-container {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.logo-image {
  height: 40px;
  width: 40px;
  object-fit: contain;
}

.logged-in-as {
  font-size: 12px;
  color: #555;
  flex: 1;
}

body.body--dark .logged-in-as {
  color: #cbd5e1;
}
</style>
