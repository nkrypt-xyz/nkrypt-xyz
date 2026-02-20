import { defineStore } from "pinia";

export const useUIStore = defineStore("ui", {
  state: () => ({
    activeGlobalObtrusiveTaskCount: 0,
  }),

  getters: {
    isGloballyLoading(state): boolean {
      return state.activeGlobalObtrusiveTaskCount > 0;
    },
  },

  actions: {
    incrementActiveGlobalObtrusiveTaskCount() {
      this.activeGlobalObtrusiveTaskCount++;
    },

    decrementActiveGlobalObtrusiveTaskCount() {
      this.activeGlobalObtrusiveTaskCount = Math.max(0, this.activeGlobalObtrusiveTaskCount - 1);
    },

    resetActiveGlobalObtrusiveTaskCount() {
      this.activeGlobalObtrusiveTaskCount = 0;
    },
  },
});
