import { RouteRecordRaw } from "vue-router";
import { GLOBAL_PERMISSION } from "constants/permissions";

const routes: RouteRecordRaw[] = [
  {
    path: "/",
    component: () => import("layouts/MainLayout.vue"),
    children: [
      {
        path: "",
        name: "dashboard",
        component: () => import("pages/DashboardPage.vue"),
        meta: { requiresAuthentication: true, title: "Dashboard" },
      },
      {
        path: "buckets",
        name: "buckets",
        component: () => import("pages/BucketsPage.vue"),
        meta: { requiresAuthentication: true, title: "Buckets" },
      },
      {
        path: "buckets/create",
        name: "bucket-create",
        component: () => import("pages/BucketCreatePage.vue"),
        meta: { requiresAuthentication: true, title: "Create Bucket", backButton: true, requiresPermission: GLOBAL_PERMISSION.CREATE_BUCKET },
      },
      {
        path: "buckets/:bucketId/edit",
        name: "bucket-edit",
        component: () => import("pages/BucketEditPage.vue"),
        meta: { requiresAuthentication: true, title: "Edit Bucket", backButton: true },
      },
      {
        path: "explore/:bucketId/:pathMatch(.*)*",
        name: "explore",
        component: () => import("pages/ExplorePage.vue"),
        meta: { requiresAuthentication: true, title: "Explore", backButton: true },
        props: true,
      },
      {
        path: "editor/text/:bucketId/:fileId",
        name: "plain-text-editor",
        component: () => import("pages/PlainTextEditorPage.vue"),
        meta: { requiresAuthentication: true, title: "Text Editor", backButton: true },
      },
      {
        path: "viewer/image/:bucketId/:fileId",
        name: "image-viewer",
        component: () => import("pages/ImageViewerPage.vue"),
        meta: { requiresAuthentication: true, title: "Image Viewer", backButton: true },
      },
      {
        path: "profile",
        name: "profile",
        component: () => import("pages/ProfilePage.vue"),
        meta: { requiresAuthentication: true, title: "Profile", backButton: true },
      },
      {
        path: "settings",
        name: "settings",
        component: () => import("pages/SettingsPage.vue"),
        meta: { requiresAuthentication: true, title: "Settings", backButton: true },
      },
      {
        path: "credentials",
        name: "credentials",
        component: () => import("pages/CredentialsPage.vue"),
        meta: { requiresAuthentication: true, title: "Credentials", backButton: true },
      },
      {
        path: "users",
        name: "users",
        component: () => import("pages/UsersPage.vue"),
        meta: { requiresAuthentication: true, title: "Users", requiresPermission: GLOBAL_PERMISSION.MANAGE_ALL_USER },
      },
      {
        path: "users/:userId",
        name: "user-save",
        component: () => import("pages/UserSavePage.vue"),
        meta: { requiresAuthentication: true, title: "User", backButton: true, requiresPermission: GLOBAL_PERMISSION.MANAGE_ALL_USER },
      },
    ],
  },
  {
    path: "/login",
    name: "login",
    component: () => import("pages/LoginPage.vue"),
    meta: { requiresAuthentication: false },
  },
  {
    path: "/:catchAll(.*)*",
    component: () => import("pages/ErrorNotFound.vue"),
  },
];

export default routes;
