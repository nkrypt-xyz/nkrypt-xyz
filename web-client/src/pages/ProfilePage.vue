<template>
  <q-page padding class="page">
    <div class="row">
      <div class="col-12 col-md-6 offset-md-3">
        <q-card class="std-card">
          <q-card-section>
            <div class="text-h6">Profile</div>
          </q-card-section>

          <q-markup-table class="base-table">
            <tbody>
              <tr>
                <td class="text-bold">Username</td>
                <td>{{ userStore.userName }}</td>
              </tr>
              <tr>
                <td class="text-bold">Display Name</td>
                <td>{{ userStore.displayName }}</td>
              </tr>
            </tbody>
          </q-markup-table>

          <q-card-section>
            <div class="row q-gutter-sm">
              <q-btn label="Update Profile" color="primary" outline @click="updateProfileClicked" />
              <q-btn label="Change Password" color="primary" outline @click="changePasswordClicked" />
            </div>
          </q-card-section>
        </q-card>
      </div>
    </div>
  </q-page>
</template>

<script setup lang="ts">
import { useUserStore } from "stores/user";
import { useUIStore } from "stores/ui";
import { callUserUpdateProfileApi, callUserUpdatePasswordApi } from "integration/user-apis";
import { dialogService } from "services/dialog-service";
import { errorService } from "services/error-service";

const userStore = useUserStore();
const uiStore = useUIStore();

async function updateProfileClicked() {
  const newDisplayName = await dialogService.prompt("Update Profile", "Enter new display name:", userStore.displayName);

  if (!newDisplayName || newDisplayName === userStore.displayName) return;

  uiStore.incrementActiveGlobalObtrusiveTaskCount();
  try {
    await callUserUpdateProfileApi({ displayName: newDisplayName });
    userStore.updateDisplayName(newDisplayName);
    dialogService.notify("success", "Profile updated successfully!");
  } catch (error) {
    await errorService.handleUnexpectedError(error);
  } finally {
    uiStore.decrementActiveGlobalObtrusiveTaskCount();
  }
}

async function changePasswordClicked() {
  const currentPassword = await dialogService.prompt("Change Password", "Enter current password:", "");
  if (!currentPassword) return;

  const newPassword = await dialogService.prompt("Change Password", "Enter new password:", "");
  if (!newPassword) return;

  const confirmPassword = await dialogService.prompt("Change Password", "Confirm new password:", "");
  if (!confirmPassword) return;

  if (newPassword !== confirmPassword) {
    await dialogService.alert("Error", "Passwords do not match");
    return;
  }

  if (newPassword.length < 8) {
    await dialogService.alert("Error", "Password must be at least 8 characters");
    return;
  }

  uiStore.incrementActiveGlobalObtrusiveTaskCount();
  try {
    await callUserUpdatePasswordApi({ currentPassword, newPassword });
    dialogService.notify("success", "Password changed successfully!");
  } catch (error) {
    await errorService.handleUnexpectedError(error);
  } finally {
    uiStore.decrementActiveGlobalObtrusiveTaskCount();
  }
}
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
