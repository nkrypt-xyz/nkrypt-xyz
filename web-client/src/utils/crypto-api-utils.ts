import { BUCKET_CRYPTO_SPEC } from "constants/crypto-specs";
import { convertSmallStringToBuffer } from "./buffer-utils";

export const buildCryptoHeader = (iv: string, salt: string): string => {
  return `${BUCKET_CRYPTO_SPEC}|${iv}|${salt}`;
};

export const unbuildCryptoHeader = (cryptoHeader: string): [ArrayBuffer, ArrayBuffer] => {
  const [_, iv, salt] = cryptoHeader.split("|");
  const ivBuffer = convertSmallStringToBuffer(iv);
  const saltBuffer = convertSmallStringToBuffer(salt);
  return [ivBuffer, saltBuffer];
};
