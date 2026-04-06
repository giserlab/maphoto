<script setup lang="ts">
import { computed } from "vue";
import { useRoute } from "vue-router";
import DefaultLayout from "@/layouts/DefaultLayout.vue";
import AdminLayout from "@/layouts/AdminLayout.vue";
import AuthLayout from "@/layouts/AuthLayout.vue";
import zhCN from "element-plus/es/locale/lang/zh-cn";

const route = useRoute();

// 需要缓存的页面名称列表
const cachedViews = computed(() => {
  // 根据路由 meta.keepAlive 自动收集需要缓存的页面
  const routes = route.matched;
  const cached: string[] = [];
  routes.forEach((r) => {
    if (r.meta?.keepAlive && r.name) {
      cached.push(r.name as string);
    }
  });
  return cached;
});

const currentLang = computed(() => zhCN);
const layout = computed(() => {
  const layoutName = route.meta.layout as string;
  switch (layoutName) {
    case "admin":
      return AdminLayout;
    case "auth":
      return AuthLayout;
    case "default":
    default:
      return DefaultLayout;
  }
});
</script>

<template>
  <el-config-provider :locale="currentLang">
    <component :is="layout">
      <router-view v-slot="{ Component }">
        <keep-alive :include="cachedViews">
          <component :is="Component" />
        </keep-alive>
      </router-view>
    </component>
  </el-config-provider>
</template>

<style>
@import "./assets/styles/variables.css";
</style>
