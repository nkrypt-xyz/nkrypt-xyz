import { usePasswordStore } from "stores/password";
import { dialogService } from "services/dialog-service";
import { decryptCryptoData } from "nkrypt-xyz-core-nodejs";
import { Bucket } from "models/common";

export async function getOrCollectPasswordForBucket(bucket: Bucket): Promise<string | null> {
  const passwordStore = usePasswordStore();

  // Check if password is already cached
  const cachedPassword = passwordStore.getPasswordForBucket(bucket._id);
  if (cachedPassword) {
    return cachedPassword;
  }

  // Prompt user for password
  const password = await dialogService.prompt("Bucket Password", `Enter the encryption password for "${bucket.name}":`, "");

  if (!password) return null;

  // Validate password
  const isValid = await decryptCryptoData(password, bucket.cryptData);
  if (!isValid) {
    await dialogService.alert("Invalid Password", "The password you entered is incorrect. Please try again.");
    return null;
  }

  //  Automatically cache the validated password
  passwordStore.setPasswordForBucket(bucket._id, password);

  return password;
}

export function validateBucketPassword(bucket: Bucket, password: string): Promise<boolean> {
  return decryptCryptoData(password, bucket.cryptData);
}
