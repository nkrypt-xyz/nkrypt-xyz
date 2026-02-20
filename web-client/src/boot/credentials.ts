import { boot } from "quasar/wrappers";
import { usePasswordStore } from "stores/password";

export default boot(({ app }) => {
  const passwordStore = usePasswordStore(app.config.globalProperties.$pinia);
  passwordStore.load();
});
