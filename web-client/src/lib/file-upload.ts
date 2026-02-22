import { callFileCreateApi, callBlobSetApi } from "integration/content-apis";
import { encryptObject, createEncryptionKeyFromPassword, makeRandomIv, makeRandomSalt, buildCryptoHeader, convertSmallUint8ArrayToString } from "nkrypt-xyz-core-nodejs";
import { MetaDataConstant } from "constants/meta-data-constants";
import { MiscConstant } from "constants/misc-constants";
import { CommonConstant } from "constants/common-constants";
import { Bucket, Directory } from "models/common";

export interface UploadProgress {
  bytesUploaded: number;
  totalBytes: number;
  percentage: number;
  status: "preparing" | "uploading" | "complete" | "error";
}

export async function uploadFileWithEncryption(
  file: globalThis.File,
  bucket: Bucket,
  parentDirectory: Directory,
  bucketPassword: string,
  onProgress?: (progress: UploadProgress) => void
): Promise<void> {
  const totalBytes = file.size;
  let bytesUploaded = 0;

  const updateProgress = (status: UploadProgress["status"]) => {
    if (onProgress) {
      onProgress({
        bytesUploaded,
        totalBytes,
        percentage: totalBytes > 0 ? Math.round((bytesUploaded / totalBytes) * 100) : 0,
        status,
      });
    }
  };

  updateProgress("preparing");

  const encryptedMetaData = await encryptObject(
    {
      contentType: file.type || "application/octet-stream",
    },
    bucketPassword
  );

  const fileRecord = await callFileCreateApi({
    name: file.name,
    bucketId: bucket._id,
    parentDirectoryId: parentDirectory._id,
    metaData: {
      [MetaDataConstant.ORIGIN_GROUP_NAME]: {
        [MetaDataConstant.ORIGIN.CLIENT_NAME]: CommonConstant.CLIENT_NAME,
        [MetaDataConstant.ORIGIN.ORIGINATION_SOURCE]: MiscConstant.ORIGINATION_SOURCE_UPLOAD,
        [MetaDataConstant.ORIGIN.ORIGINATION_DATE]: Date.now(),
      },
    },
    encryptedMetaData,
  });

  updateProgress("uploading");

  // Generate random IV and salt for this encryption
  const { iv } = await makeRandomIv();
  const { salt } = await makeRandomSalt();

  // Create encryption key from password and salt
  const { key: encryptionKey } = await createEncryptionKeyFromPassword(bucketPassword, salt);

  // Encrypt the file data
  const fileArrayBuffer = await file.arrayBuffer();
  const plainData = new Uint8Array(fileArrayBuffer);

  const encryptedData = await crypto.subtle.encrypt(
    {
      name: "AES-GCM",
      iv,
    },
    encryptionKey,
    plainData
  );

  // Build crypto header for HTTP header (not embedded in blob)
  const ivStr = convertSmallUint8ArrayToString(iv);
  const saltStr = convertSmallUint8ArrayToString(salt);
  const cryptoMetaHeader = buildCryptoHeader(ivStr, saltStr);

  const blob = new Blob([encryptedData]);

  await callBlobSetApi({
    bucketId: bucket._id,
    fileId: fileRecord.fileId,
    blob,
    cryptoMeta: cryptoMetaHeader,
  });

  bytesUploaded = totalBytes;
  updateProgress("complete");
}
