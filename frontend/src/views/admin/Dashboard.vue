<template>
  <div class="dashboard-page">
    <el-row :gutter="20">
      <el-col :xs="24" :sm="12" :lg="6">
        <el-card class="stat-card">
          <div class="stat-value">{{ places.length }}</div>
          <div class="stat-label">地点总数</div>
        </el-card>
      </el-col>

      <el-col :xs="24" :sm="12" :lg="6">
        <el-card class="stat-card">
          <div class="stat-value">{{ totalPhotos }}</div>
          <div class="stat-label">照片总数</div>
        </el-card>
      </el-col>

      <el-col :xs="24" :sm="12" :lg="6">
        <el-card class="stat-card">
          <div class="stat-value">{{ groups.length }}</div>
          <div class="stat-label">分组数量</div>
        </el-card>
      </el-col>

      <el-col :xs="24" :sm="12" :lg="6">
        <el-card class="stat-card">
          <div class="stat-value">{{ publicPlaces }}</div>
          <div class="stat-label">公开地点</div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" class="mt-20">
      <el-col :xs="24" :lg="16">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>最近添加的地点</span>
              <el-button type="primary" @click="$router.push('/admin/places')">
                查看全部
              </el-button>
            </div>
          </template>

          <el-table :data="recentPlaces" v-loading="placesStore.loading">
            <el-table-column prop="name" label="名称">
              <template #default="{ row }">
                {{ row.name || '未命名' }}
              </template>
            </el-table-column>
            <el-table-column prop="group" label="分组" />
            <el-table-column label="照片数" width="100">
              <template #default="{ row }">
                {{ row.photos?.length || 0 }}
              </template>
            </el-table-column>
            <el-table-column label="操作" width="150">
              <template #default="{ row }">
                <el-button type="primary" link @click="editPlace(row.id)">
                  编辑
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>

      <el-col :xs="24" :lg="8" class="mt-20 lg:mt-0">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>快速操作</span>
            </div>
          </template>

          <div class="quick-actions">
            <el-button type="primary" @click="$router.push('/admin/places/create')">
              <el-icon><Plus /></el-icon>
              添加地点
            </el-button>

            <el-button @click="$router.push('/admin/settings')">
              <el-icon><Setting /></el-icon>
              个人设置
            </el-button>

            <el-button v-if="isAdmin" @click="$router.push('/admin/users')">
              <el-icon><User /></el-icon>
              用户管理
            </el-button>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Plus, Setting, User } from '@element-plus/icons-vue'
import { usePlacesStore } from '@/stores/places'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const placesStore = usePlacesStore()
const authStore = useAuthStore()

const isAdmin = computed(() => authStore.isAdmin)
const places = computed(() => placesStore.places)
const groups = computed(() => placesStore.groups)

const totalPhotos = computed(() => {
  return places.value.reduce((sum, p) => sum + (p.photos?.length || 0), 0)
})

const publicPlaces = computed(() => {
  return places.value.filter(p => !p.private).length
})

const recentPlaces = computed(() => {
  return [...places.value]
    .sort((a, b) => new Date(b.date || 0).getTime() - new Date(a.date || 0).getTime())
    .slice(0, 5)
})

onMounted(() => {
  placesStore.fetchPlaces()
})

function editPlace(id: number) {
  router.push(`/admin/places/${id}/edit`)
}
</script>

<style scoped>
.stat-card {
  text-align: center;
  margin-bottom: 20px;
}

.stat-value {
  font-size: 32px;
  font-weight: bold;
  color: #409EFF;
  line-height: 1.2;
}

.stat-label {
  font-size: 14px;
  color: #909399;
  margin-top: 8px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.quick-actions {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.quick-actions .el-button {
  justify-content: flex-start;
}

.mt-20 {
  margin-top: 20px;
}

@media (min-width: 1024px) {
  .lg\\:mt-0 {
    margin-top: 0;
  }
}
</style>
