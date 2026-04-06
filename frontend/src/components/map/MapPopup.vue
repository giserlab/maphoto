<template>
  <div class="map-popup">
    <img v-if="place.cover" :src="place.cover" class="popup-image" @click="showGallery = true" />
    <div v-else class="popup-no-image">
      <el-icon :size="40"><Picture /></el-icon>
    </div>
    <div class="popup-content">
      <h4 class="popup-title">{{ place.name || "未命名地点" }}</h4>
      <p v-if="place.desc" class="popup-desc">{{ place.desc }}</p>
      <div class="popup-meta">
        <span v-if="place.group" class="popup-group">{{ place.group }}</span>
        <span class="photo-count">
          <el-icon><Picture /></el-icon>
          {{ place.photos?.length || 0 }}
        </span>
      </div>
    </div>

    <!-- Photo Gallery Dialog -->
    <PhotoViewer v-model="showGallery" :photos="place.photos || []" :title="place.name" />
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { Picture } from "@element-plus/icons-vue";
import PhotoViewer from "@/components/photo/PhotoViewer.vue";
import type { Place } from "@/types";

const props = defineProps<{
  place: Place;
}>();

const showGallery = ref(false);
</script>

<style scoped>
.map-popup {
  min-width: 200px;
  max-width: 280px;
}

.popup-image {
  width: 100%;
  height: 120px;
  object-fit: cover;
  border-radius: 4px;
  cursor: pointer;
  margin-bottom: 8px;
}

.popup-no-image {
  width: 100%;
  height: 120px;
  background: #f5f7fa;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #909399;
  border-radius: 4px;
  margin-bottom: 8px;
}

.popup-content {
  padding: 4px;
}

.popup-title {
  margin: 0 0 8px;
  font-size: 16px;
  font-weight: 500;
  color: #303133;
}

.popup-desc {
  margin: 0 0 8px;
  font-size: 13px;
  color: #606266;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.popup-meta {
  display: flex;
  align-items: center;
  gap: 10px;
}

.popup-group {
  font-size: 12px;
  padding: 2px 8px;
  background: #ecf5ff;
  color: #409eff;
  border-radius: 4px;
}

.photo-count {
  font-size: 12px;
  color: #909399;
  display: flex;
  align-items: center;
  gap: 4px;
}
</style>
