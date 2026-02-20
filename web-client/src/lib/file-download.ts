import { callBlobGetApi } from "integration/content-apis";
import { createEncryptionKeyFromPassword } from "utils/crypto-utils";
import { unbuildCryptoHeader } from "utils/crypto-api-utils";
import { File } from "models/common";
import streamSaver from "streamsaver";

export interface DownloadProgress {
  bytesDownloaded: number;
  totalBytes: number;
  percentage: number;
  status: "downloading" | "decrypting" | "complete" | "error";
}

export async function downloadFileWithDecryption(file: File, bucketPassword: string, onProgress?: (progress: DownloadProgress) => void): Promise<void> {
  const updateProgress = (status: DownloadProgress["status"], bytesDownloaded = 0, totalBytes = 0) => {
    if (onProgress) {
      onProgress({
        bytesDownloaded,
        totalBytes,
        percentage: totalBytes > 0 ? Math.round((bytesDownloaded / totalBytes) * 100) : 0,
        status,
      });
    }
  };

  updateProgress("downloading");

  const { blob, cryptoMeta } = await callBlobGetApi({
    bucketId: file.bucketId,
    fileId: file._id,
  });

  if (!cryptoMeta) {
    throw new Error("Missing crypto metadata from server");
  }

  const arrayBuffer = await blob.arrayBuffer();
  const encryptedData = new Uint8Array(arrayBuffer);

  updateProgress("decrypting", 0, encryptedData.length);

  // Parse crypto header to get IV and salt
  const [ivBuffer, saltBuffer] = unbuildCryptoHeader(cryptoMeta);
  const iv = new Uint8Array(ivBuffer);
  const salt = new Uint8Array(saltBuffer);

  // Create encryption key from password and salt
  const { key: encryptionKey } = await createEncryptionKeyFromPassword(bucketPassword, salt);

  // Decrypt the blob data (no embedded header, just encrypted content)
  const decryptedData = await crypto.subtle.decrypt(
    {
      name: "AES-GCM",
      iv,
    },
    encryptionKey,
    arrayBuffer
  );

  updateProgress("complete", encryptedData.length, encryptedData.length);

  const fileStream = streamSaver.createWriteStream(file.name, {
    size: decryptedData.byteLength,
  });

  const writer = fileStream.getWriter();
  await writer.write(new Uint8Array(decryptedData));
  await writer.close();
}
