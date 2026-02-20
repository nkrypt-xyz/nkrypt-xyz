import { boot } from "quasar/wrappers";
import { useSettingsStore } from "stores/settings";

export default boot(({ app }) => {
  const settingsStore = useSettingsStore(app.config.globalProperties.$pinia);

  settingsStore.load();

  if (settingsStore.darkMode !== null) {
    app.config.globalProperties.$q.dark.set(settingsStore.darkMode);
  }
});
