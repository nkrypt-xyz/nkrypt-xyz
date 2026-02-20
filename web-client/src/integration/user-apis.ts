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

export async function callUserCreateApi(data: { userName: string; displayName: string; password: string; globalPermissions: Record<string, boolean> }) {
  const sessionStore = useSessionStore();
  return await callPostJsonApi(sessionStore.session!.serverUrl, sessionStore.session!.apiKey, "/api/user/create", data);
}

export async function callUserUpdateApi(data: { userId: string; displayName: string; globalPermissions: Record<string, boolean> }) {
  const sessionStore = useSessionStore();
  return await callPostJsonApi(sessionStore.session!.serverUrl, sessionStore.session!.apiKey, "/api/user/update", data);
}
