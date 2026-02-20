import { BLOB_CHUNK_SIZE_INCLUDING_TAG_BYTES } from "./crypto-specs";

export const CommonConstant = {
  APP_VERSION: "2.0.0",
  COPYRIGHT: "nkrypt.xyz",
  COPYRIGHT_HREF: "https://nkrypt.xyz",
  DEFAULT_SERVER_URL: import.meta.env.VITE_DEFAULT_SERVER_URL || "",
  PACKET_SIZE_FOR_QUANTIZED_STREAMS: 100 * BLOB_CHUNK_SIZE_INCLUDING_TAG_BYTES,
  CLIENT_NAME: import.meta.env.VITE_CLIENT_APPLICATION_NAME || "nkrypt-web-client",
};
