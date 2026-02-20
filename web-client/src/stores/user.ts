import { defineStore } from "pinia";
import { User } from "models/common";
import { getLocalStorageItem, setLocalStorageItem, removeLocalStorageItem } from "utils/store-utils";

const USER_LOCALSTORAGE_KEY = "--nkrypt-user";

export const useUserStore = defineStore("user", {
  state: () => ({
    user: getLocalStorageItem<User | null>(USER_LOCALSTORAGE_KEY, null),
  }),

  getters: {
    isUserLoggedIn(state): boolean {
      return !!state.user;
    },

    displayName(state): string {
      return state.user?.displayName || "";
    },

    userName(state): string {
      return state.user?.userName || "";
    },

    userId(state): string {
      return state.user?.userId || "";
    },

    globalPermissions(state): Record<string, boolean> {
      return state.user?.globalPermissions || {};
    },

    hasPermission: (state) => (permission: string) => {
      return state.user?.globalPermissions[permission] === true;
    },
  },

  actions: {
    setUser(user: User | null) {
      this.user = user;
      setLocalStorageItem(USER_LOCALSTORAGE_KEY, user);
    },

    updateDisplayName(displayName: string) {
      if (this.user) {
        this.user = { ...this.user, displayName };
        setLocalStorageItem(USER_LOCALSTORAGE_KEY, this.user);
      }
    },

    clearUser() {
      this.user = null;
      removeLocalStorageItem(USER_LOCALSTORAGE_KEY);
    },
  },
});
