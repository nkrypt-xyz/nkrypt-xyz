<template>
  <q-page padding class="page">
    <div class="row">
      <div class="col-12 col-lg-10 offset-lg-1">
        <div class="flex items-center justify-between q-mb-md">
          <div class="text-h5">Users</div>
          <q-btn v-if="canCreateUser()" icon="add" label="Create User" color="primary" @click="createUserClicked" />
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
                <th class="text-left">Username</th>
                <th class="text-left">Display Name</th>
                <th class="text-left">Permissions</th>
                <th class="text-right">Actions</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="user in userList" :key="user._id">
                <td>{{ user.userName }}</td>
                <td>{{ user.displayName }}</td>
                <td>
                  <q-chip v-for="(value, key) in user.globalPermissions" :key="key" v-show="value" dense color="primary" text-color="white" size="sm">
                    {{ key }}
                  </q-chip>
                </td>
                <td class="text-right">
                  <q-btn icon="edit" flat dense @click="editUser(user)">
                    <q-tooltip>Edit</q-tooltip>
                  </q-btn>
                </td>
              </tr>
            </tbody>
          </q-markup-table>

          <q-card-section v-if="userList.length === 0">
            <div class="text-center text-grey-7">No users found.</div>
          </q-card-section>
        </q-card>
      </div>
    </div>
  </q-page>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useRouter } from "vue-router";
import { useUserStore } from "stores/user";
import { callUserListApi, callUserFindApi } from "integration/user-apis";
import { errorService } from "services/error-service";
import { MiscConstant } from "constants/misc-constants";
import { GLOBAL_PERMISSION } from "constants/permissions";

const router = useRouter();
const userStore = useUserStore();

const isLoading = ref(true);
const userList = ref<any[]>([]);
const canCreateUser = () => userStore.hasPermission(GLOBAL_PERMISSION.CREATE_USER);

async function loadUsers() {
  try {
    isLoading.value = true;
    const listResponse = await callUserListApi({});
    const list = listResponse.userList || [];

    if (list.length === 0) {
      userList.value = [];
      return;
    }

    const findResponse = await callUserFindApi({
      filters: list.map((u: { _id: string }) => ({ by: "userId" as const, userId: u._id })),
      includeGlobalPermissions: true,
    });
    const findMap = new Map((findResponse.userList || []).map((u: { _id: string }) => [u._id, u]));
    userList.value = list.map((u: { _id: string }) => findMap.get(u._id) || u);
  } catch (error) {
    await errorService.handleUnexpectedError(error);
  } finally {
    isLoading.value = false;
  }
}

function createUserClicked() {
  router.push({ name: "user-save", params: { userId: MiscConstant.NEW_ID_PLACEHOLDER } });
}

function editUser(user: any) {
  router.push({ name: "user-save", params: { userId: user._id } });
}

onMounted(() => {
  loadUsers();
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
