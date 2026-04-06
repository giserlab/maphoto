<template>
  <div class="map-controls">
    <el-button
      circle
      :icon="Plus"
      @click="$emit('zoom-in')"
    />
    <el-button
      circle
      :icon="Minus"
      @click="$emit('zoom-out')"
    />
    <el-button
      circle
      :icon="Location"
      @click="$emit('locate')"
    />
  </div>
</template>

<script setup lang="ts">
import { Plus, Minus, Location } from '@element-plus/icons-vue'

defineEmits<{
  'zoom-in': []
  'zoom-out': []
  'locate': []
}>()
</script>

<style scoped>
.map-controls {
  position: absolute;
  right: 16px;
  bottom: 24px;
  display: flex;
  flex-direction: column;
  gap: 8px;
  z-index: 400;
}

.map-controls :deep(.el-button) {
  width: 36px;
  height: 36px;
  background: #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}
</style>
