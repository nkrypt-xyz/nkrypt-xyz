/* eslint-disable */

declare namespace NodeJS {
  interface ProcessEnv {
    NODE_ENV: string;
    VUE_ROUTER_MODE: "hash" | "history" | "abstract" | undefined;
    VUE_ROUTER_BASE: string | undefined;
  }
}

interface ImportMetaEnv {
  readonly VITE_API_URL?: string;
  readonly VITE_TEST_API_URL?: string;
  readonly VITE_TEST_MODE?: string;
  readonly VITE_DEFAULT_SERVER_URL?: string;
  readonly VITE_CLIENT_APPLICATION_NAME?: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
