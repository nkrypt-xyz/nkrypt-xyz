import { defineStore } from "pinia";
import { getLocalStorageItem, setLocalStorageItem } from "utils/store-utils";

const SUGGESTED_SERVER_URL_KEY = "--nkrypt-suggested-server";

export const useCacheStore = defineStore("cache", {
  state: () => ({
    suggestedServerUrl: getLocalStorageItem(SUGGESTED_SERVER_URL_KEY, ""),
  }),

  actions: {
    setSuggestedServerUrl(url: string) {
      this.suggestedServerUrl = url;
      setLocalStorageItem(SUGGESTED_SERVER_URL_KEY, url);
    },
  },
});
