<template>
  <div class="photo-manage-page">
    <el-card v-if="place">
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <el-button @click="$router.back()">
              <el-icon><ArrowLeft /></el-icon>
              返回
            </el-button>
            <span class="place-name">{{ place.name || '未命名地点' }}</span>
          </div>
          <el-button type="primary" @click="showUpload = true">
            <el-icon><Plus /></el-icon>
            添加照片
          </el-button>
        </div>
      </template>

      <!-- Photos grid -->
      <div class="photos-grid">
        <div
          v-for="photo in place.photos"
          :key="photo.id"
          class="photo-item"
        >
          <div class="photo-wrapper">
            <img :src="photo.url" @click="previewPhoto(photo)" />

            <div class="photo-actions">
              <el-button
                v-if="place.cover !== photo.url"
                type="primary"
                link
                @click="setAsCover(photo.url)"
              >
                设为封面
              </el-button>
              <el-tag v-else size="small" type="success">封面</el-tag>

              <el-button type="danger" link @click="deletePhoto(photo)">
                删除
              </el-button>
            </div>
          </div>
        </div>

        <el-empty v-if="place.photos?.length === 0" description="暂无照片" />
      </div>
    </el-card>

    <!-- Add photo dialog -->
    <el-dialog
      v-model="showUpload"
      title="添加照片"
      width="500px"
      destroy-on-close
    >
      <PhotoUpload @success="handlePhotoAdded" />
    </el-dialog>

    <!-- Preview dialog -->
    <el-dialog
      v-model="showPreview"
      title="图片预览"
      width="80%"
      top="5vh"
      destroy-on-close
    >
      <img :src="previewUrl" style="width: 100%" />
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, ArrowLeft } from '@element-plus/icons-vue'
import { usePlacesStore } from '@/stores/places'
import PhotoUpload from '@/components/photo/PhotoUpload.vue'
import type { Place, Photo } from '@/types'

const props = defineProps<{
  id: string
}>()

const router = useRouter()
const placesStore = usePlacesStore()

const showUpload = ref(false)
const showPreview = ref(false)
const previewUrl = ref('')

const place = computed(() =>
  placesStore.places.find(p => p.id === parseInt(props.id))
)

onMounted(() => {
  if (placesStore.places.length === 0) {
    placesStore.fetchPlaces()
  }
})

function handlePhotoAdded(url: string) {
  placesStore.addPlacePhoto(parseInt(props.id), url)
  showUpload.value = false
  ElMessage.success('添加成功')
}

async function setAsCover(url: string) {
  try {
    await placesStore.updatePlaceCover(parseInt(props.id), url)
    ElMessage.success('封面设置成功')
  } catch {
    ElMessage.error('设置失败')
  }
}

async function deletePhoto(photo: Photo) {
  try {
    await ElMessageBox.confirm('确定要删除这张照片吗？', '确认删除', {
      type: 'warning',
    })
    await placesStore.removePlacePhoto(parseInt(props.id), photo.url)
    ElMessage.success('删除成功')
  } catch {
    // cancelled
  }
}

function previewPhoto(photo: Photo) {
  previewUrl.value = photo.url
  showPreview.value = true
}
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-header .header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.card-header .place-name {
  font-size: 16px;
  font-weight: 500;
}

.photos-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 16px;
}

.photo-item {
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.photo-wrapper {
  position: relative;
}

.photo-wrapper img {
  width: 100%;
  height: 150px;
  object-fit: cover;
  display: block;
  cursor: pointer;
}

.photo-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 12px;
  background: #fff;
  border-top: 1px solid #e4e7ed;
}

@media (max-width: 768px) {
  .photos-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>
