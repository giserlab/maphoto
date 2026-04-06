<template>
  <el-dialog
    v-model="visible"
    :title="title"
    width="90%"
    top="5vh"
    class="photo-viewer-dialog"
    destroy-on-close
  >
    <div v-if="photos.length === 0" class="no-photos">
      <el-empty description="暂无照片" />
    </div>

    <el-carousel v-else height="60vh" :interval="5000" arrow="always">
      <el-carousel-item v-for="photo in photos" :key="photo.id">
        <div class="carousel-image-container">
          <img :src="photo.url" class="carousel-image" />
        </div>
      </el-carousel-item>
    </el-carousel>
  </el-dialog>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Photo } from '@/types'

const props = defineProps<{
  modelValue: boolean
  photos: Photo[]
  title?: string
}>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
}>()

const visible = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value),
})
</script>

<style scoped>
.photo-viewer-dialog :deep(.el-dialog__body) {
  padding: 10px;
}

.no-photos {
  padding: 40px;
}

.carousel-image-container {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  background: #f5f7fa;
}

.carousel-image {
  max-width: 100%;
  max-height: 100%;
  object-fit: contain;
}
</style>
