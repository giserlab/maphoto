<template>
  <div class="place-create-page">
    <div class="page-actions">
      <el-button @click="$router.back()">
        <el-icon><ArrowLeft /></el-icon>
        返回
      </el-button>
    </div>

    <PlaceForm @submit="handleSubmit" @cancel="$router.back()" />
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { ArrowLeft } from '@element-plus/icons-vue'
import PlaceForm from '@/components/place/PlaceForm.vue'
import { usePlacesStore } from '@/stores/places'
import type { PlaceAddForm } from '@/types'

const router = useRouter()
const placesStore = usePlacesStore()

async function handleSubmit(formData: PlaceAddForm) {
  try {
    await placesStore.createPlace(formData)
    ElMessage.success('添加成功')
    router.push('/admin/places')
  } catch (error) {
    console.error('Failed to create place:', error)
  }
}
</script>

<style scoped>
.page-actions {
  margin-bottom: 20px;
}
</style>
