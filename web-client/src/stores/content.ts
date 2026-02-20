import { defineStore } from "pinia";
import { Bucket } from "models/common";

export const useContentStore = defineStore("content", {
  state: () => ({
    bucketList: [] as Bucket[],
    activeBucket: null as Bucket | null,
  }),

  getters: {
    getBucketById: (state) => (bucketId: string) => {
      return state.bucketList.find((b) => b._id === bucketId);
    },
  },

  actions: {
    setBucketList(buckets: Bucket[]) {
      this.bucketList = buckets;
    },

    setActiveBucket(bucket: Bucket | null) {
      this.activeBucket = bucket;
    },

    addBucket(bucket: Bucket) {
      this.bucketList.push(bucket);
    },

    updateBucket(bucketId: string, updates: Partial<Bucket>) {
      const index = this.bucketList.findIndex((b) => b._id === bucketId);
      if (index !== -1) {
        this.bucketList[index] = { ...this.bucketList[index], ...updates };
      }
    },

    removeBucket(bucketId: string) {
      this.bucketList = this.bucketList.filter((b) => b._id !== bucketId);
      if (this.activeBucket?._id === bucketId) {
        this.activeBucket = null;
      }
    },
  },
});
