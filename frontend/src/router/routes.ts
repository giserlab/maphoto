import type { RouteRecordRaw } from "vue-router";

export const routes: RouteRecordRaw[] = [
  // Public routes
  {
    path: "/",
    name: "Home",
    component: () => import("@/views/home/Index.vue"),
    meta: { public: true, layout: "default", keepAlive: true },
  },
  {
    path: "/share/:username",
    name: "Share",
    component: () => import("@/views/share/Index.vue"),
    meta: { public: true, layout: "default", keepAlive: true },
    props: (route) => ({
      username: route.params.username,
      group: route.query.group,
    }),
  },
  {
    path: "/login",
    name: "Login",
    component: () => import("@/views/login/Index.vue"),
    meta: { public: true, guestOnly: true, layout: "auth" },
  },

  // Admin routes (requires auth)
  {
    path: "/admin",
    component: () => import("@/layouts/AdminLayout.vue"),
    meta: { requiresAuth: true },
    children: [
      {
        path: "",
        name: "AdminDashboard",
        component: () => import("@/views/admin/Dashboard.vue"),
        meta: { keepAlive: true },
      },
      {
        path: "places",
        name: "PlaceList",
        component: () => import("@/views/admin/places/List.vue"),
        meta: { keepAlive: true },
      },
      {
        path: "places/create",
        name: "PlaceCreate",
        component: () => import("@/views/admin/places/Create.vue"),
      },
      {
        path: "places/:id/edit",
        name: "PlaceEdit",
        component: () => import("@/views/admin/places/Edit.vue"),
        props: true,
      },
      {
        path: "places/:id/photos",
        name: "PlacePhotos",
        component: () => import("@/views/admin/photos/Manage.vue"),
        props: true,
      },
      {
        path: "settings",
        name: "UserSettings",
        component: () => import("@/views/admin/settings/Index.vue"),
      },
      {
        path: "cos",
        name: "CosManager",
        component: () => import("@/views/admin/cos/Index.vue"),
      },
      {
        path: "storage",
        name: "LocalStorage",
        component: () => import("@/views/admin/storage/Index.vue"),
      },
      {
        path: "users",
        name: "UserManagement",
        component: () => import("@/views/admin/users/List.vue"),
        meta: { requiresAdmin: true },
      },
    ],
  },

  // 404
  {
    path: "/:pathMatch(.*)*",
    name: "NotFound",
    component: () => import("@/views/error/404.vue"),
  },
];
