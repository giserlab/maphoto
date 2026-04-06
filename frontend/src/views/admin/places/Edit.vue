<template>
  <div class="place-edit-page">
    <div class="page-actions">
      <el-button @click="$router.back()">
        <el-icon><ArrowLeft /></el-icon>
        返回
      </el-button>
    </div>

    <PlaceForm
      v-if="place"
      :place="place"
      @submit="handleSubmit"
      @cancel="$router.back()"
    />

    <el-skeleton v-else :rows="10" animated />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { ArrowLeft } from '@element-plus/icons-vue'
import PlaceForm from '@/components/place/PlaceForm.vue'
import { usePlacesStore } from '@/stores/places'
import type { Place, PlaceUpdateForm } from '@/types'

const props = defineProps<{
  id: string
}>()

const router = useRouter()
const placesStore = usePlacesStore()

const place = computed(() =>
  placesStore.places.find(p => p.id === parseInt(props.id))
)

onMounted(() => {
  if (placesStore.places.length === 0) {
    placesStore.fetchPlaces()
  }
})

async function handleSubmit(formData: PlaceUpdateForm) {
  try {
    await placesStore.updatePlace(parseInt(props.id), formData)
    ElMessage.success('更新成功')
    router.push('/admin/places')
  } catch (error) {
    console.error('Failed to update place:', error)
  }
}
</script>

<style scoped>
.page-actions {
  margin-bottom: 20px;
}
</style>
