import { Router } from "vue-router";

let routerInstance: Router | null = null;

export const navigationService = {
  setRouter(router: Router) {
    routerInstance = router;
  },

  async navigateTo(name: string, params?: any, query?: any): Promise<void> {
    if (!routerInstance) {
      console.error("Router not initialized");
      return;
    }
    await routerInstance.push({ name, params, query });
  },

  async push(path: string): Promise<void> {
    if (!routerInstance) {
      console.error("Router not initialized");
      return;
    }
    await routerInstance.push(path);
  },

  async navigateReplace(name: string, params?: any, query?: any): Promise<void> {
    if (!routerInstance) {
      console.error("Router not initialized");
      return;
    }
    await routerInstance.replace({ name, params, query });
  },

  navigateBack(): void {
    if (!routerInstance) {
      console.error("Router not initialized");
      return;
    }
    routerInstance.back();
  },

  async navigateToPreviousPageOrDashboard(): Promise<void> {
    if (window.history.length > 2) {
      this.navigateBack();
    } else {
      await this.navigateTo("dashboard");
    }
  },
};
