import { defineStore } from "pinia";
import { Settings } from "models/common";
import { getLocalStorageItem, setLocalStorageItem } from "utils/store-utils";

const SETTINGS_KEY = "--nkrypt-settings";

export const useSettingsStore = defineStore("settings", {
  state: (): Settings => ({
    uploadMechanism: "stream",
    downloadMechanism: "stream",
    plainTextEditorNoRestrictions: false,
    darkMode: null,
  }),

  actions: {
    setUploadMechanism(mechanism: string) {
      this.uploadMechanism = mechanism;
      this.persist();
    },

    setDownloadMechanism(mechanism: string) {
      this.downloadMechanism = mechanism;
      this.persist();
    },

    setPlainTextEditorNoRestrictions(value: boolean) {
      this.plainTextEditorNoRestrictions = value;
      this.persist();
    },

    setDarkMode(isDark: boolean | null) {
      this.darkMode = isDark;
      this.persist();
    },

    persist() {
      setLocalStorageItem(SETTINGS_KEY, this.$state);
    },

    load() {
      const stored = getLocalStorageItem<Settings | null>(SETTINGS_KEY, null);
      if (stored) {
        this.$patch(stored);
      }
    },
  },
});
