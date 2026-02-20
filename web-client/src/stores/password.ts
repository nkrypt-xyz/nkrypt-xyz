import { defineStore } from "pinia";

export const usePasswordStore = defineStore("password", {
  state: () => ({
    bucketPasswordCache: {} as Record<string, string>,
  }),

  getters: {
    getPasswordForBucket: (state) => (bucketId: string) => {
      return state.bucketPasswordCache[bucketId];
    },
  },

  actions: {
    setPasswordForBucket(bucketId: string, password: string) {
      this.bucketPasswordCache[bucketId] = password;
    },

    clearPasswordForBucket(bucketId: string) {
      delete this.bucketPasswordCache[bucketId];
    },

    clearAllPasswords() {
      this.bucketPasswordCache = {};
    },
  },
});
