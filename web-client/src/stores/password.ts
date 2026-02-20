import { defineStore } from "pinia";
import { getLocalStorageItem, setLocalStorageItem, removeLocalStorageItem } from "utils/store-utils";

const BUCKET_PASSWORD_CACHE_KEY = "--nkrypt-bucket-password-cache";

export const usePasswordStore = defineStore("password", {
  state: () => ({
    bucketPasswordCache: getLocalStorageItem<Record<string, string>>(BUCKET_PASSWORD_CACHE_KEY, {}),
  }),

  getters: {
    getPasswordForBucket: (state) => (bucketId: string) => {
      return state.bucketPasswordCache[bucketId];
    },

    cachedBucketIds: (state) => Object.keys(state.bucketPasswordCache),
  },

  actions: {
    setPasswordForBucket(bucketId: string, password: string) {
      this.bucketPasswordCache[bucketId] = password;
      this.persist();
    },

    clearPasswordForBucket(bucketId: string) {
      delete this.bucketPasswordCache[bucketId];
      this.persist();
    },

    clearAllPasswords() {
      this.bucketPasswordCache = {};
      removeLocalStorageItem(BUCKET_PASSWORD_CACHE_KEY);
    },

    load() {
      const stored = getLocalStorageItem<Record<string, string>>(BUCKET_PASSWORD_CACHE_KEY, {});
      if (stored && typeof stored === "object" && Object.keys(stored).length > 0) {
        this.bucketPasswordCache = { ...stored };
      }
    },

    persist() {
      setLocalStorageItem(BUCKET_PASSWORD_CACHE_KEY, this.bucketPasswordCache);
    },
  },
});
