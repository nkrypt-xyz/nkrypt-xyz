import { callFileCreateApi } from "integration/content-apis";
import { encryptObject } from "nkrypt-xyz-core-nodejs";
import { MetaDataConstant } from "constants/meta-data-constants";
import { MiscConstant } from "constants/misc-constants";
import { CommonConstant } from "constants/common-constants";
import { Bucket, Directory } from "models/common";

export async function createTextFile(bucket: Bucket, parentDirectory: Directory, bucketPassword: string, fileName: string): Promise<void> {
  const encryptedMetaData = await encryptObject(
    {
      contentType: MiscConstant.TEXT_FILE_MIME,
    },
    bucketPassword
  );

  await callFileCreateApi({
    name: fileName,
    bucketId: bucket._id,
    parentDirectoryId: parentDirectory._id,
    metaData: {
      [MetaDataConstant.ORIGIN_GROUP_NAME]: {
        [MetaDataConstant.ORIGIN.CLIENT_NAME]: CommonConstant.CLIENT_NAME,
        [MetaDataConstant.ORIGIN.ORIGINATION_SOURCE]: MiscConstant.ORIGINATION_SOURCE_CREATE_TEXT_FILE,
        [MetaDataConstant.ORIGIN.ORIGINATION_DATE]: Date.now(),
      },
    },
    encryptedMetaData,
  });
}
