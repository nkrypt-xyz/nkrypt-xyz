<template>
  <q-layout view="lHh Lpr lFf">
    <q-page-container>
      <q-page class="row items-center justify-evenly">
        <q-card class="login-card">
      <div class="app-name q-pa-md">
        <img class="logo" src="~/assets/logo-512-sqr.png" alt="nkrypt.xyz" />
        <span>nkrypt.xyz</span>
      </div>

      <q-card-section>
        <div class="text-h6 text-center q-mb-md">Sign In</div>

        <q-form ref="loginForm" @submit="onSubmit" class="q-gutter-md">
          <q-input standout="bg-primary text-white" v-model="serverUrl" label="Server Address" lazy-rules :rules="validators.serverUrl">
            <template v-slot:prepend>
              <q-icon name="dns" />
            </template>
          </q-input>

          <q-input standout="bg-primary text-white" v-model="username" label="Username" lazy-rules :rules="validators.username">
            <template v-slot:prepend>
              <q-icon name="person" />
            </template>
          </q-input>

          <q-input type="password" standout="bg-primary text-white" v-model="password" label="Password" lazy-rules :rules="validators.password">
            <template v-slot:prepend>
              <q-icon name="lock" />
            </template>
          </q-input>

          <div class="row justify-center q-mt-lg">
            <q-btn label="Login" type="submit" color="primary" unelevated :loading="isLoading" padding="sm xl" icon="login" />
          </div>
        </q-form>
      </q-card-section>
        </q-card>
      </q-page>
    </q-page-container>
  </q-layout>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { QForm } from "quasar";
import { useRouter, useRoute } from "vue-router";
import { useUserStore } from "stores/user";
import { useSessionStore } from "stores/session";
import { useCacheStore } from "stores/cache";
import { callUserLoginApi } from "integration/user-apis";
import { dialogService } from "services/dialog-service";
import { errorService } from "services/error-service";
import { validators } from "utils/validators";
import { CommonConstant } from "constants/common-constants";

const router = useRouter();
const route = useRoute();

const userStore = useUserStore();
const sessionStore = useSessionStore();
const cacheStore = useCacheStore();

const loginForm = ref<QForm | null>(null);

  console.log("CommonConstant.DEFAULT_SERVER_URL", CommonConstant.DEFAULT_SERVER_URL);
  console.log("cacheStore.suggestedServerUrl", cacheStore.suggestedServerUrl);

const serverUrl = ref(cacheStore.suggestedServerUrl || CommonConstant.DEFAULT_SERVER_URL || "");
const username = ref("");
const password = ref("");

const isLoading = ref(false);

async function onSubmit() {
  const isValid = await loginForm.value!.validate();
  if (!isValid) return;

  isLoading.value = true;

  try {
    const response = await callUserLoginApi(serverUrl.value, {
      userName: username.value,
      password: password.value,
    });

    const { apiKey } = response;
    const { userName, displayName, _id: userId, globalPermissions } = response.user;

    userStore.setUser({ userName, displayName, userId, globalPermissions });
    sessionStore.setSession({ apiKey, serverUrl: serverUrl.value });
    cacheStore.setSuggestedServerUrl(serverUrl.value);

    dialogService.notify("success", "Login successful!");

    if (route.query.next) {
      await router.push(route.query.next as string);
    } else {
      await router.push({ name: "dashboard" });
    }
  } catch (error) {
    await errorService.handleUnexpectedError(error);
  } finally {
    isLoading.value = false;
  }
}
</script>

<style scoped lang="scss">
.login-card {
  min-width: 300px;
  max-width: 500px;
  width: 100%;
  margin: 12px;
}

.app-name {
  text-align: center;
  background-color: rgb(35, 35, 35);
  color: white;
  text-transform: uppercase;
  font-size: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;

  .logo {
    width: 40px;
    height: 40px;
  }
}

body.body--dark .app-name {
  background-color: #1e2538;
}
</style>
