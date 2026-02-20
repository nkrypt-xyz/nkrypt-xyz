<template>
  <q-page padding class="page">
    <div class="row">
      <div class="col-12 col-md-6 offset-md-3">
        <q-card class="std-card" v-if="isLoading">
          <q-card-section class="text-center">
            <q-spinner color="primary" size="50px" />
          </q-card-section>
        </q-card>

        <q-card class="std-card" v-if="!isLoading">
          <q-card-section>
            <div class="text-h6">{{ isNewUser ? "Create User" : "Edit User" }}</div>
          </q-card-section>

          <q-card-section>
            <q-form ref="formRef" @submit="onSubmit" class="q-gutter-md">
              <q-input standout="bg-primary text-white" v-model="userData.userName" label="Username" :readonly="!isNewUser" lazy-rules :rules="[(val) => !!val || 'Username is required', (val) => val.length >= 4 || 'Username must be at least 4 characters']">
                <template v-slot:prepend>
                  <q-icon name="person" />
                </template>
              </q-input>

              <q-input standout="bg-primary text-white" v-model="userData.displayName" label="Display Name" lazy-rules :rules="[(val) => !!val || 'Display name is required']">
                <template v-slot:prepend>
                  <q-icon name="badge" />
                </template>
              </q-input>

              <q-input v-if="isNewUser" type="password" standout="bg-primary text-white" v-model="userData.password" label="Password" lazy-rules :rules="[(val) => !!val || 'Password is required', (val) => val.length >= 8 || 'Password must be at least 8 characters']">
                <template v-slot:prepend>
                  <q-icon name="lock" />
                </template>
              </q-input>

              <div class="q-mt-md">
                <div class="text-subtitle2 q-mb-sm">Global Permissions</div>
                <q-checkbox v-model="userData.globalPermissions[GLOBAL_PERMISSION.MANAGE_ALL_USER]" :label="GLOBAL_PERMISSION.MANAGE_ALL_USER" />
                <q-checkbox v-model="userData.globalPermissions[GLOBAL_PERMISSION.CREATE_USER]" :label="GLOBAL_PERMISSION.CREATE_USER" />
                <q-checkbox v-model="userData.globalPermissions[GLOBAL_PERMISSION.CREATE_BUCKET]" :label="GLOBAL_PERMISSION.CREATE_BUCKET" />
              </div>

              <div class="row justify-end q-gutter-sm q-mt-lg">
                <q-btn label="Cancel" outline color="grey" @click="cancelClicked" />
                <q-btn :label="isNewUser ? 'Create User' : 'Save Changes'" type="submit" color="primary" unelevated :loading="isSaving" />
              </div>
            </q-form>
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
import { useUIStore } from "stores/ui";
import { callUserFindApi, callAdminAddUserApi, callAdminSetGlobalPermissionsApi } from "integration/user-apis";
import { dialogService } from "services/dialog-service";
import { errorService } from "services/error-service";
import { MiscConstant } from "constants/misc-constants";
import { GLOBAL_PERMISSION } from "constants/permissions";

const router = useRouter();
const route = useRoute();
const uiStore = useUIStore();

const formRef = ref<QForm | null>(null);

const isLoading = ref(true);
const isSaving = ref(false);
const isNewUser = ref(true);

const userData = ref({
  userName: "",
  displayName: "",
  password: "",
  globalPermissions: {
    [GLOBAL_PERMISSION.MANAGE_ALL_USER]: false,
    [GLOBAL_PERMISSION.CREATE_USER]: false,
    [GLOBAL_PERMISSION.CREATE_BUCKET]: false,
  },
});

async function loadUser() {
  const userId = route.params.userId as string;

  if (userId === MiscConstant.NEW_ID_PLACEHOLDER) {
    isNewUser.value = true;
    isLoading.value = false;
    return;
  }

  isNewUser.value = false;

  try {
    const response = await callUserFindApi({
      filters: [{ by: "userId", userId }],
      includeGlobalPermissions: true,
    });
    const user = response.userList?.[0];

    if (!user) {
      await dialogService.alert("Error", "User not found");
      router.push({ name: "users" });
      return;
    }

    userData.value.userName = user.userName;
    userData.value.displayName = user.displayName;
    userData.value.globalPermissions = {
      ...{
        [GLOBAL_PERMISSION.MANAGE_ALL_USER]: false,
        [GLOBAL_PERMISSION.CREATE_USER]: false,
        [GLOBAL_PERMISSION.CREATE_BUCKET]: false,
      },
      ...user.globalPermissions,
    };
  } catch (error) {
    await errorService.handleUnexpectedError(error);
  } finally {
    isLoading.value = false;
  }
}

async function onSubmit() {
  const isValid = await formRef.value!.validate();
  if (!isValid) return;

  isSaving.value = true;
  uiStore.incrementActiveGlobalObtrusiveTaskCount();

  try {
    if (isNewUser.value) {
      const addResponse = await callAdminAddUserApi({
        userName: userData.value.userName,
        displayName: userData.value.displayName,
        password: userData.value.password,
      });
      const newUserId = addResponse.userId;

      await callAdminSetGlobalPermissionsApi({
        userId: newUserId,
        globalPermissions: userData.value.globalPermissions,
      });
      dialogService.notify("success", "User created successfully!");
    } else {
      const userId = route.params.userId as string;
      await callAdminSetGlobalPermissionsApi({
        userId,
        globalPermissions: userData.value.globalPermissions,
      });
      dialogService.notify("success", "User updated successfully!");
    }

    router.push({ name: "users" });
  } catch (error) {
    await errorService.handleUnexpectedError(error);
  } finally {
    isSaving.value = false;
    uiStore.decrementActiveGlobalObtrusiveTaskCount();
  }
}

function cancelClicked() {
  router.back();
}

onMounted(() => {
  loadUser();
});
</script>

<style scoped lang="scss">
.page {
  display: flex;
  flex-direction: column;
}
</style>
