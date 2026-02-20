import { defineStore } from "pinia";
import { Session } from "models/common";
import { getLocalStorageItem, setLocalStorageItem, removeLocalStorageItem } from "utils/store-utils";

const SESSION_LOCALSTORAGE_KEY = "--nkrypt-session";

export const useSessionStore = defineStore("session", {
  state: () => ({
    session: getLocalStorageItem<Session | null>(SESSION_LOCALSTORAGE_KEY, null),
  }),

  getters: {
    isSessionActive(state): boolean {
      return !!state.session;
    },

    apiKey(state): string | null {
      return state.session?.apiKey || null;
    },

    serverUrl(state): string | null {
      return state.session?.serverUrl || null;
    },
  },

  actions: {
    setSession(session: Session | null) {
      this.session = session;
      setLocalStorageItem(SESSION_LOCALSTORAGE_KEY, session);
    },

    clearSession() {
      this.session = null;
      removeLocalStorageItem(SESSION_LOCALSTORAGE_KEY);
    },
  },
});
