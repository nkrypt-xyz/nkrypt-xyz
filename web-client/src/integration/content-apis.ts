import { useSessionStore } from "stores/session";
import { callPostJsonApi } from "utils/api-utils";

export async function callBucketListApi(data: Record<string, never>) {
  const sessionStore = useSessionStore();
  return await callPostJsonApi(sessionStore.session!.serverUrl, sessionStore.session!.apiKey, "/api/bucket/list", data);
}

export async function callBucketCreateApi(data: { name: string; cryptSpec: string; cryptData: string; metaData: Record<string, any> }) {
  const sessionStore = useSessionStore();
  return await callPostJsonApi(sessionStore.session!.serverUrl, sessionStore.session!.apiKey, "/api/bucket/create", data);
}

export async function callBucketUpdateApi(data: { bucketId: string; name: string }) {
  const sessionStore = useSessionStore();
  return await callPostJsonApi(sessionStore.session!.serverUrl, sessionStore.session!.apiKey, "/api/bucket/update", data);
}

export async function callBucketDestroyApi(data: { bucketId: string; name: string }) {
  const sessionStore = useSessionStore();
  return await callPostJsonApi(sessionStore.session!.serverUrl, sessionStore.session!.apiKey, "/api/bucket/destroy", data);
}

export async function callBucketSetAuthorizationApi(data: { bucketId: string; userId: string; isAuthorized: boolean }) {
  const sessionStore = useSessionStore();
  return await callPostJsonApi(sessionStore.session!.serverUrl, sessionStore.session!.apiKey, "/api/bucket/set-authorization", data);
}

export async function callDirectoryGetApi(data: { bucketId: string; directoryId: string }) {
  const sessionStore = useSessionStore();
  return await callPostJsonApi(sessionStore.session!.serverUrl, sessionStore.session!.apiKey, "/api/directory/get", data);
}

export async function callDirectoryCreateApi(data: { bucketId: string; name: string; parentDirectoryId: string; metaData: Record<string, any>; encryptedMetaData: string }) {
  const sessionStore = useSessionStore();
  return await callPostJsonApi(sessionStore.session!.serverUrl, sessionStore.session!.apiKey, "/api/directory/create", data);
}

export async function callDirectoryRenameApi(data: { bucketId: string; directoryId: string; name: string }) {
  const sessionStore = useSessionStore();
  return await callPostJsonApi(sessionStore.session!.serverUrl, sessionStore.session!.apiKey, "/api/directory/rename", data);
}

export async function callDirectoryMoveApi(data: { bucketId: string; directoryId: string; targetParentDirectoryId: string }) {
  const sessionStore = useSessionStore();
  return await callPostJsonApi(sessionStore.session!.serverUrl, sessionStore.session!.apiKey, "/api/directory/move", data);
}

export async function callDirectoryDeleteApi(data: { bucketId: string; directoryId: string }) {
  const sessionStore = useSessionStore();
  return await callPostJsonApi(sessionStore.session!.serverUrl, sessionStore.session!.apiKey, "/api/directory/delete", data);
}

export async function callFileGetApi(data: { bucketId: string; fileId: string }) {
  const sessionStore = useSessionStore();
  return await callPostJsonApi(sessionStore.session!.serverUrl, sessionStore.session!.apiKey, "/api/file/get", data);
}

export async function callFileCreateApi(data: { bucketId: string; name: string; parentDirectoryId: string; metaData: Record<string, any>; encryptedMetaData: string }) {
  const sessionStore = useSessionStore();
  return await callPostJsonApi(sessionStore.session!.serverUrl, sessionStore.session!.apiKey, "/api/file/create", data);
}

export async function callFileRenameApi(data: { bucketId: string; fileId: string; name: string }) {
  const sessionStore = useSessionStore();
  return await callPostJsonApi(sessionStore.session!.serverUrl, sessionStore.session!.apiKey, "/api/file/rename", data);
}

export async function callFileMoveApi(data: { bucketId: string; fileId: string; targetParentDirectoryId: string }) {
  const sessionStore = useSessionStore();
  return await callPostJsonApi(sessionStore.session!.serverUrl, sessionStore.session!.apiKey, "/api/file/move", data);
}

export async function callFileDeleteApi(data: { bucketId: string; fileId: string }) {
  const sessionStore = useSessionStore();
  return await callPostJsonApi(sessionStore.session!.serverUrl, sessionStore.session!.apiKey, "/api/file/delete", data);
}

export async function callFileSetMetaDataApi(data: { bucketId: string; fileId: string; metaData: Record<string, any> }) {
  const sessionStore = useSessionStore();
  return await callPostJsonApi(sessionStore.session!.serverUrl, sessionStore.session!.apiKey, "/api/file/set-meta-data", data);
}

export async function callFileSetEncryptedMetaDataApi(data: { bucketId: string; fileId: string; encryptedMetaData: string }) {
  const sessionStore = useSessionStore();
  return await callPostJsonApi(sessionStore.session!.serverUrl, sessionStore.session!.apiKey, "/api/file/set-encrypted-meta-data", data);
}

export class BlobApiError extends Error {
  constructor(
    message: string,
    public readonly status: number,
    public readonly code?: string
  ) {
    super(message);
    this.name = "BlobApiError";
  }
}

export async function callBlobGetApi(data: { bucketId: string; fileId: string }): Promise<{ blob: Blob; cryptoMeta: string }> {
  const sessionStore = useSessionStore();
  const response = await fetch(`${sessionStore.session!.serverUrl}/api/blob/read/${data.bucketId}/${data.fileId}`, {
    method: "POST",
    headers: {
      Authorization: `Bearer ${sessionStore.session!.apiKey}`,
    },
  });

  if (!response.ok) {
    let code: string | undefined;
    try {
      const errBody = await response.json();
      code = errBody?.error?.code;
    } catch {
      // Ignore JSON parse errors
    }
    throw new BlobApiError(`Failed to read blob: ${response.statusText}`, response.status, code);
  }

  const blob = await response.blob();
  const cryptoMeta = response.headers.get("nk-crypto-meta") || "";

  return { blob, cryptoMeta };
}

export async function callBlobSetApi(data: { bucketId: string; fileId: string; blob: Blob; cryptoMeta: string }): Promise<void> {
  const sessionStore = useSessionStore();
  const response = await fetch(`${sessionStore.session!.serverUrl}/api/blob/write/${data.bucketId}/${data.fileId}`, {
    method: "POST",
    headers: {
      Authorization: `Bearer ${sessionStore.session!.apiKey}`,
      "nk-crypto-meta": data.cryptoMeta,
      "Content-Type": "application/octet-stream",
    },
    body: data.blob,
  });

  if (!response.ok) {
    throw new Error(`Failed to write blob: ${response.statusText}`);
  }
}
