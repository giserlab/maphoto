<template>
  <div class="places-list-page">
    <el-card v-loading="placesStore.loading">
      <template #header>
        <div class="card-header">
          <span>地点列表</span>
          <div>
            <el-button type="success" :disabled="!authStore.username" @click="shareGroup">
              <el-icon><Share /></el-icon>
              {{ filterGroup ? "分享分组" : "分享全部" }}
            </el-button>
            <el-button type="primary" @click="$router.push('/admin/places/create')">
              <el-icon><Plus /></el-icon>
              添加地点
            </el-button>
          </div>
        </div>
      </template>

      <!-- Filter bar -->
      <div class="filter-bar">
        <el-input v-model="searchQuery" placeholder="搜索地点名称" clearable style="width: 250px">
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>

        <el-select v-model="filterGroup" placeholder="全部分组" clearable>
          <el-option
            v-for="group in placesStore.groups"
            :key="group"
            :label="group"
            :value="group"
          />
        </el-select>
      </div>

      <!-- Places table -->
      <el-table :data="filteredPlaces" stripe>
        <el-table-column label="封面" width="100">
          <template #default="{ row }">
            <el-image
              v-if="row.cover"
              :src="row.cover"
              fit="cover"
              style="width: 60px; height: 60px; border-radius: 4px"
            />
            <div v-else class="no-cover">
              <el-icon><Picture /></el-icon>
            </div>
          </template>
        </el-table-column>

        <el-table-column prop="name" label="名称">
          <template #default="{ row }">
            {{ row.name || "未命名" }}
            <el-tag v-if="row.private" size="small" type="warning">私有</el-tag>
          </template>
        </el-table-column>

        <el-table-column prop="group" label="分组" width="120" />

        <el-table-column label="坐标" width="180">
          <template #default="{ row }">
            {{ formatCoordinate(row.lat) }}, {{ formatCoordinate(row.lon) }}
          </template>
        </el-table-column>

        <el-table-column label="照片数" width="100">
          <template #default="{ row }">
            {{ row.photos?.length || 0 }}
          </template>
        </el-table-column>

        <el-table-column label="操作" width="250" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="editPlace(row.id)"> 编辑 </el-button>
            <el-button type="primary" link @click="managePhotos(row.id)"> 照片 </el-button>
            <el-button type="danger" link @click="deletePlace(row)"> 删除 </el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-empty v-if="filteredPlaces.length === 0" description="暂无地点" />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import { useRouter } from "vue-router";
import { ElMessage, ElMessageBox } from "element-plus";
import { Plus, Search, Picture, Share } from "@element-plus/icons-vue";
import { usePlacesStore } from "@/stores/places";
import { useAuthStore } from "@/stores/auth";
import { formatCoordinate } from "@/utils/format";
import type { Place } from "@/types";

const router = useRouter();
const placesStore = usePlacesStore();
const authStore = useAuthStore();

const searchQuery = ref("");
const filterGroup = ref("");

const filteredPlaces = computed(() => {
  let result = placesStore.places;

  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase();
    result = result.filter(
      (p) => p.name?.toLowerCase().includes(query) || p.desc?.toLowerCase().includes(query),
    );
  }

  if (filterGroup.value) {
    result = result.filter((p) => p.group === filterGroup.value);
  }

  return result;
});

onMounted(() => {
  placesStore.fetchPlaces();
});

function editPlace(id: number) {
  router.push(`/admin/places/${id}/edit`);
}

function managePhotos(id: number) {
  router.push(`/admin/places/${id}/photos`);
}

async function deletePlace(place: Place) {
  try {
    await ElMessageBox.confirm(`确定要删除地点 "${place.name || "未命名"}" 吗？`, "确认删除", {
      type: "warning",
    });
    await placesStore.deletePlace(place.id);
    ElMessage.success("删除成功");
  } catch {
    // cancelled
  }
}

function shareGroup() {
  if (!authStore.username) {
    ElMessage.warning("未登录用户无法分享");
    return;
  }
  const baseUrl = window.location.origin;
  const query = filterGroup.value ? `?group=${encodeURIComponent(filterGroup.value)}` : "";
  const fullUrl = `${baseUrl}/#/?user=${authStore.username}&group=${query}`;
  navigator.clipboard
    .writeText(fullUrl)
    .then(() => {
      if (filterGroup.value) {
        ElMessage.success(`分组 "${filterGroup.value}" 的分享链接已复制到剪贴板`);
      } else {
        ElMessage.success("全部地点的分享链接已复制到剪贴板");
      }
    })
    .catch((err) => {
      console.error("复制失败:", err);
      ElMessage.error("复制失败，请手动复制链接");
    });
}
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.filter-bar {
  display: flex;
  gap: 16px;
  margin-bottom: 20px;
}

.no-cover {
  width: 60px;
  height: 60px;
  background: #f5f7fa;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #c0c4cc;
}

@media (max-width: 768px) {
  .filter-bar {
    flex-direction: column;
  }

  .filter-bar .el-input {
    width: 100% !important;
  }
}
</style>
