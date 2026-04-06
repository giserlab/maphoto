<template>
  <div class="share-page" :class="{ 'is-mobile': isMobile }">
    <!-- Map Section -->
    <div class="map-section">
      <MapContainer
        :places="places"
        :config="config"
        @select="handlePlaceSelect"
      />
    </div>

    <!-- Sidebar for desktop -->
    <div v-if="!isMobile" class="sidebar">
      <div class="sidebar-header">
        <h2 class="site-title">{{ config?.title || `${username} 的地图` }}</h2>
        <p v-if="config?.note" class="site-note">{{ config.note }}</p>

        <el-select
          v-model="selectedGroup"
          placeholder="全部分组"
          clearable
          class="group-filter"
          @change="handleGroupChange"
        >
          <el-option
            v-for="group in groups"
            :key="group"
            :label="group"
            :value="group"
          />
        </el-select>
      </div>

      <div class="place-list" v-loading="loading">
        <PlaceCard
          v-for="place in filteredPlaces"
          :key="place.id"
          :place="place"
          compact
          class="place-list-item"
          @click="handlePlaceSelect(place)"
        />

        <el-empty v-if="filteredPlaces.length === 0" description="暂无地点" />
      </div>
    </div>

    <!-- Mobile header -->
    <div v-else class="mobile-header">
      <h1 class="mobile-title">{{ config?.title || `${username} 的地图` }}</h1>
      <el-button circle :icon="Menu" @click="showMobileDrawer = true" />
    </div>

    <!-- Mobile drawer -->
    <el-drawer
      v-model="showMobileDrawer"
      title="地点列表"
      size="85%"
      direction="rtl"
    >
      <div class="mobile-drawer-content">
        <p v-if="config?.note" class="site-note">{{ config.note }}</p>

        <el-select
          v-model="selectedGroup"
          placeholder="全部分组"
          clearable
          class="group-filter"
          @change="handleGroupChange"
        >
          <el-option
            v-for="group in groups"
            :key="group"
            :label="group"
            :value="group"
          />
        </el-select>

        <div class="place-list-mobile" v-loading="loading">
          <PlaceCard
            v-for="place in filteredPlaces"
            :key="place.id"
            :place="place"
            class="place-list-item"
            @click="handlePlaceSelectMobile(place)"
          />

          <el-empty v-if="filteredPlaces.length === 0" description="暂无地点" />
        </div>
      </div>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { Menu } from '@element-plus/icons-vue'
import MapContainer from '@/components/map/MapContainer.vue'
import PlaceCard from '@/components/place/PlaceCard.vue'
import { useMobile } from '@/composables/useMobile'
import * as shareApi from '@/api/share'
import type { Place, Config } from '@/types'

const props = defineProps<{
  username: string
  group?: string
}>()

const { isMobile } = useMobile()

const loading = ref(false)
const places = ref<Place[]>([])
const config = ref<Config | null>(null)
const selectedGroup = ref(props.group || '')
const showMobileDrawer = ref(false)

const groups = computed(() => {
  const groupSet = new Set(places.value.map((p) => p.group).filter(Boolean) as string[])
  return Array.from(groupSet)
})

const filteredPlaces = computed(() => {
  if (!selectedGroup.value) return places.value
  return places.value.filter((p) => p.group === selectedGroup.value)
})

onMounted(() => {
  fetchData()
})

watch(() => props.group, (newGroup) => {
  selectedGroup.value = newGroup || ''
  fetchData()
})

async function fetchData() {
  loading.value = true
  try {
    const res = await shareApi.getSharedPlaces(props.username, selectedGroup.value || undefined)
    config.value = res.data.config

    // Convert GeoJSON to places
    const features = res.data.features?.features || []
    places.value = features.map((f: any) => ({
      id: f.properties.id,
      userid: f.properties.userid,
      name: f.properties.name,
      desc: f.properties.desc,
      cover: f.properties.cover,
      private: f.properties.private,
      group: f.properties.group || '',
      date: f.properties.date ? new Date(f.properties.date) : undefined,
      photos: f.properties.photos || [],
      lon: f.geometry.coordinates[0],
      lat: f.geometry.coordinates[1],
    }))
  } finally {
    loading.value = false
  }
}

function handlePlaceSelect(place: Place) {
  console.log('Selected place:', place)
}

function handlePlaceSelectMobile(place: Place) {
  showMobileDrawer.value = false
  handlePlaceSelect(place)
}

function handleGroupChange() {
  fetchData()
}
</script>

<style scoped>
.share-page {
  display: flex;
  height: 100vh;
  overflow: hidden;
}

.share-page.is-mobile {
  position: relative;
}

.map-section {
  flex: 1;
  height: 100%;
}

.is-mobile .map-section {
  height: 100vh;
}

.sidebar {
  width: 320px;
  background: #fff;
  border-left: 1px solid #e4e7ed;
  display: flex;
  flex-direction: column;
}

.sidebar-header {
  padding: 16px;
  border-bottom: 1px solid #e4e7ed;
}

.site-title {
  margin: 0 0 8px;
  font-size: 18px;
  color: #303133;
}

.site-note {
  margin: 0 0 12px;
  font-size: 13px;
  color: #606266;
  line-height: 1.4;
}

.group-filter {
  width: 100%;
}

.place-list {
  flex: 1;
  overflow-y: auto;
  padding: 12px;
}

.place-list-item {
  margin-bottom: 12px;
  cursor: pointer;
}

/* Mobile styles */
.mobile-header {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: 50px;
  background: #fff;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  z-index: 500;
}

.mobile-title {
  margin: 0;
  font-size: 18px;
  color: #303133;
}

.mobile-drawer-content {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.place-list-mobile {
  flex: 1;
  overflow-y: auto;
  padding: 12px 0;
}
</style>
