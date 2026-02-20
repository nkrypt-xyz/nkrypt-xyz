import { useSessionStore } from "stores/session";
import { callPostJsonApi } from "utils/api-utils";

export async function callUserLoginApi(serverUrl: string, data: { userName: string; password: string }) {
  return await callPostJsonApi(serverUrl, "", "/api/user/login", data);
}

export async function callUserUpdateProfileApi(data: { displayName: string }) {
  const sessionStore = useSessionStore();
  return await callPostJsonApi(sessionStore.session!.serverUrl, sessionStore.session!.apiKey, "/api/user/update-profile", data);
}

export async function callUserUpdatePasswordApi(data: { currentPassword: string; newPassword: string }) {
  const sessionStore = useSessionStore();
  return await callPostJsonApi(sessionStore.session!.serverUrl, sessionStore.session!.apiKey, "/api/user/update-password", data);
}

export async function callUserListApi(data: Record<string, never>) {
  const sessionStore = useSessionStore();
  return await callPostJsonApi(sessionStore.session!.serverUrl, sessionStore.session!.apiKey, "/api/user/list", data);
}

export async function callUserFindApi(data: { filters: Array<{ by: string; userId?: string; userName?: string }>; includeGlobalPermissions: boolean }) {
  const sessionStore = useSessionStore();
  return await callPostJsonApi(sessionStore.session!.serverUrl, sessionStore.session!.apiKey, "/api/user/find", data);
}

export async function callAdminAddUserApi(data: { userName: string; displayName: string; password: string }) {
  const sessionStore = useSessionStore();
  return await callPostJsonApi(sessionStore.session!.serverUrl, sessionStore.session!.apiKey, "/api/admin/iam/add-user", data);
}

export async function callAdminSetGlobalPermissionsApi(data: { userId: string; globalPermissions: Record<string, boolean> }) {
  const sessionStore = useSessionStore();
  return await callPostJsonApi(sessionStore.session!.serverUrl, sessionStore.session!.apiKey, "/api/admin/iam/set-global-permissions", data);
}
