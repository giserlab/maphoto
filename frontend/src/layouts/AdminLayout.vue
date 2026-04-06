<template>
  <el-container class="admin-layout">
    <!-- Sidebar for desktop -->
    <el-aside v-if="!isMobile" width="220px" class="sidebar">
      <div class="logo">
        <h2>Maphoto</h2>
      </div>
      <div class="user-info">
        <el-icon :size="18"><User /></el-icon>
        <span class="username">{{ username }}</span>
      </div>
      <el-menu
        :default-active="$route.path"
        router
        class="admin-menu"
        background-color="#304156"
        text-color="#bfcbd9"
        active-text-color="#409EFF"
      >
        <el-menu-item index="/admin">
          <el-icon><Odometer /></el-icon>
          <span>控制台</span>
        </el-menu-item>
        <el-menu-item index="/admin/places">
          <el-icon><Location /></el-icon>
          <span>地点管理</span>
        </el-menu-item>
        <el-sub-menu index="/admin/storage">
          <template #title>
            <el-icon><Picture /></el-icon>
            <span>图片存储</span>
          </template>
          <el-menu-item index="/admin/cos">腾讯云COS</el-menu-item>
          <el-menu-item index="/admin/storage">本地存储</el-menu-item>
        </el-sub-menu>
        <el-menu-item index="/admin/settings">
          <el-icon><Setting /></el-icon>
          <span>设置</span>
        </el-menu-item>
        <el-menu-item v-if="isAdmin" index="/admin/users">
          <el-icon><User /></el-icon>
          <span>用户管理</span>
        </el-menu-item>
      </el-menu>
      <div class="sidebar-actions">
        <el-button type="primary" @click="goHome">
          <el-icon><View /></el-icon>
          <span>查看首页</span>
        </el-button>
        <el-button type="danger" @click="handleLogout">
          <el-icon><SwitchButton /></el-icon>
          <span>退出</span>
        </el-button>
      </div>
    </el-aside>

    <el-container class="main-container">
      <!-- Main content -->
      <el-main class="main-content">
        <router-view />
      </el-main>

      <!-- Bottom navigation for mobile -->
      <div v-if="isMobile" class="mobile-nav">
        <div
          v-for="item in mobileNavItems"
          :key="item.path"
          class="nav-item"
          :class="{ active: $route.path === item.path }"
          @click="$router.push(item.path)"
        >
          <el-icon :size="20">
            <component :is="item.icon" />
          </el-icon>
          <span>{{ item.label }}</span>
        </div>
      </div>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { computed } from "vue";
import { useRouter } from "vue-router";
import { ElMessageBox } from "element-plus";
import {
  Odometer,
  Location,
  Setting,
  User,
  View,
  SwitchButton,
  Picture,
} from "@element-plus/icons-vue";
import { useAuthStore } from "@/stores/auth";
import { useMobile } from "@/composables/useMobile";

const router = useRouter();
const authStore = useAuthStore();
const { isMobile } = useMobile();

const username = computed(() => authStore.username);
const isAdmin = computed(() => authStore.isAdmin);

const mobileNavItems = computed(() => {
  const items = [
    { path: "/admin", icon: "Odometer", label: "控制台" },
    { path: "/admin/places", icon: "Location", label: "地点" },
    { path: "/admin/cos", icon: "Picture", label: "存储" },
    { path: "/admin/settings", icon: "Setting", label: "设置" },
  ];
  if (isAdmin.value) {
    items.push({ path: "/admin/users", icon: "User", label: "用户" });
  }
  return items;
});

function goHome() {
  const username = useAuthStore().username;
  window.open(`/#/?user=${username}`, "_blank");
}

async function handleLogout() {
  try {
    await ElMessageBox.confirm("确定要退出登录吗？", "提示", {
      confirmButtonText: "确定",
      cancelButtonText: "取消",
      type: "warning",
    });
    await authStore.logout();
    router.push("/login");
  } catch {
    // cancelled
  }
}
</script>

<style scoped>
.admin-layout {
  height: 100vh;
  overflow: hidden;
}

.sidebar {
  background-color: #304156;
  display: flex;
  flex-direction: column;
  height: 100vh;
  overflow: hidden;
}

.logo {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  border-bottom: 1px solid #1f2d3d;
}

.user-info {
  padding: 12px 20px;
  display: flex;
  align-items: center;
  gap: 8px;
  color: #bfcbd9;
  background-color: #263445;
  border-bottom: 1px solid #1f2d3d;
}

.user-info .username {
  font-size: 14px;
  font-weight: 500;
}

.admin-menu {
  border-right: none;
  flex: 1;
  overflow-y: auto;
}

.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background-color: #fff;
  box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
}

.header-left {
  font-size: 16px;
  font-weight: 500;
}

.header-right {
  display: flex;
  gap: 10px;
}

.sidebar-actions {
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 10px;
  border-top: 1px solid #1f2d3d;
}

.sidebar-actions .el-button {
  width: 100%;
  margin: 0em !important;
  justify-content: flex-start;
}

.main-container {
  height: 100vh;
  overflow: hidden;
}

.main-content {
  background-color: #f0f2f5;
  padding: 20px;
  height: 100vh;
  overflow-y: auto;
  box-sizing: border-box;
}

.admin-layout:has(.mobile-nav) .main-content {
  padding-bottom: 76px;
}

.mobile-nav {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  height: 56px;
  background-color: #fff;
  display: flex;
  justify-content: space-around;
  align-items: center;
  box-shadow: 0 -1px 4px rgba(0, 21, 41, 0.08);
  z-index: 100;
}

.nav-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  color: #909399;
  font-size: 12px;
  cursor: pointer;
}

.nav-item.active {
  color: #409eff;
}
</style>
