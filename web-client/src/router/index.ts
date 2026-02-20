import { route } from "quasar/wrappers";
import { createMemoryHistory, createRouter, createWebHashHistory, createWebHistory } from "vue-router";

import routes from "./routes";
import { useUserStore } from "stores/user";
import { useSessionStore } from "stores/session";

export default route(function (/* { store, ssrContext } */) {
  const createHistory = process.env.SERVER
    ? createMemoryHistory
    : process.env.VUE_ROUTER_MODE === "history"
      ? createWebHistory
      : createWebHashHistory;

  const Router = createRouter({
    scrollBehavior: () => ({ left: 0, top: 0 }),
    routes,
    history: createHistory(process.env.VUE_ROUTER_BASE),
  });

  Router.beforeEach((to, from, next) => {
    const userStore = useUserStore();
    const sessionStore = useSessionStore();

    const requiresAuth = to.meta.requiresAuthentication !== false;
    const isAuthenticated = userStore.isUserLoggedIn && sessionStore.isSessionActive;

    if (requiresAuth && !isAuthenticated) {
      next({
        name: "login",
        query: { next: to.fullPath },
      });
      return;
    }

    if (!requiresAuth && isAuthenticated && to.name === "login") {
      next({ name: "dashboard" });
      return;
    }

    const requiredPermission = to.meta.requiresPermission as string | undefined;
    if (requiredPermission && !userStore.hasPermission(requiredPermission)) {
      next({ name: "dashboard" });
      return;
    }

    next();
  });

  return Router;
});
