<template>
  <div class="place-card" :class="{ 'is-compact': compact }">
    <div class="place-cover" @click="handleClick">
      <img v-if="place.cover" :src="place.cover" :alt="place.name" />
      <div v-else class="no-cover">
        <el-icon :size="40"><Picture /></el-icon>
      </div>
      <span v-if="place.group" class="place-group">{{ place.group }}</span>
    </div>
    <div class="place-info">
      <h3 class="place-name">{{ place.name || '未命名地点' }}</h3>
      <p v-if="!compact && place.desc" class="place-desc">{{ place.desc }}</p>
      <div class="place-meta">
        <span class="photo-count">
          <el-icon><Picture /></el-icon>
          {{ place.photos?.length || 0 }}
        </span>
        <span v-if="place.date" class="place-date">{{ formatDate(place.date) }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Picture } from '@element-plus/icons-vue'
import { formatDate } from '@/utils/format'
import type { Place } from '@/types'

const props = defineProps<{
  place: Place
  compact?: boolean
}>()

const emit = defineEmits<{
  click: [place: Place]
}>()

function handleClick() {
  emit('click', props.place)
}
</script>

<style scoped>
.place-card {
  background: #fff;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.05);
  transition: transform 0.2s, box-shadow 0.2s;
}

.place-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
}

.place-card.is-compact {
  display: flex;
  gap: 12px;
  padding: 12px;
}

.place-cover {
  position: relative;
  aspect-ratio: 16/10;
  overflow: hidden;
  cursor: pointer;
}

.is-compact .place-cover {
  width: 80px;
  height: 80px;
  flex-shrink: 0;
  border-radius: 4px;
  aspect-ratio: auto;
}

.place-cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s;
}

.place-cover:hover img {
  transform: scale(1.05);
}

.no-cover {
  width: 100%;
  height: 100%;
  background: #f5f7fa;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #c0c4cc;
}

.place-group {
  position: absolute;
  top: 8px;
  left: 8px;
  font-size: 12px;
  padding: 2px 8px;
  background: rgba(64, 158, 255, 0.9);
  color: #fff;
  border-radius: 4px;
}

.place-info {
  padding: 12px;
}

.is-compact .place-info {
  padding: 0;
  flex: 1;
}

.place-name {
  margin: 0 0 8px;
  font-size: 15px;
  font-weight: 500;
  color: #303133;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.place-desc {
  margin: 0 0 8px;
  font-size: 13px;
  color: #606266;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.place-meta {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 12px;
  color: #909399;
}

.photo-count {
  display: flex;
  align-items: center;
  gap: 4px;
}
</style>
