<template>
  <el-form ref="formRef" :model="form" :rules="rules" label-position="top" class="place-form">
    <el-row :gutter="20">
      <el-col :xs="24" :lg="12">
        <el-form-item label="地点名称" prop="name">
          <el-input v-model="form.name" placeholder="输入地点名称" />
        </el-form-item>

        <el-form-item label="描述" prop="desc">
          <el-input v-model="form.desc" type="textarea" :rows="4" placeholder="输入地点描述" />
        </el-form-item>

        <el-form-item label="分组" prop="group">
          <el-select
            v-model="form.group"
            placeholder="选择或输入分组"
            filterable
            allow-create
            default-first-option
            style="width: 100%"
          >
            <el-option v-for="group in existingGroups" :key="group" :label="group" :value="group" />
          </el-select>
        </el-form-item>

        <el-form-item label="坐标">
          <div class="coordinate-inputs">
            <el-input-number
              v-model="form.lat"
              :precision="6"
              :step="0.001"
              placeholder="纬度"
              style="flex: 1"
            />
            <el-input-number
              v-model="form.lon"
              :precision="6"
              :step="0.001"
              placeholder="经度"
              style="flex: 1"
            />
            <el-button type="primary" plain @click="handleSelectPhoto">
              <el-icon><Upload /></el-icon>
              从照片读取
            </el-button>
            <input
              ref="fileInputRef"
              type="file"
              accept="image/*"
              style="display: none"
              @change="handleFileChange"
            />
          </div>
          <p class="form-hint">点击地图选择位置，或从本地照片 EXIF 信息自动读取坐标</p>
        </el-form-item>

        <el-form-item label="封面图片 URL">
          <el-input v-model="form.cover" placeholder="输入图片 URL 地址" clearable />
          <div v-if="form.cover" class="cover-preview">
            <el-image
              :src="form.cover"
              fit="cover"
              :preview-src-list="[form.cover]"
              style="width: 200px; height: 150px; border-radius: 8px"
            >
              <template #error>
                <div class="image-error">
                  <el-icon><Picture /></el-icon>
                  <span>加载失败</span>
                </div>
              </template>
            </el-image>
          </div>
          <p class="form-hint">输入图片的完整 URL 地址</p>
        </el-form-item>

        <el-form-item v-if="!isEdit" label="照片列表">
          <div class="photos-input">
            <el-input
              v-model="photoInput"
              placeholder="输入照片 URL 地址"
              clearable
              @keyup.enter="addPhoto"
            >
              <template #append>
                <el-button @click="addPhoto">
                  <el-icon><Plus /></el-icon>
                </el-button>
              </template>
            </el-input>
          </div>
          <div v-if="form.photos.length > 0" class="photos-list">
            <div v-for="(photo, index) in form.photos" :key="index" class="photo-tag">
              <el-image
                :src="photo"
                fit="cover"
                style="width: 60px; height: 60px; border-radius: 4px"
                :preview-src-list="[photo]"
              />
              <el-button type="danger" link size="small" @click="removePhoto(index)">
                <el-icon><Delete /></el-icon>
              </el-button>
            </div>
          </div>
          <p class="form-hint">已添加 {{ form.photos.length }} 张照片</p>
        </el-form-item>
      </el-col>

      <el-col :xs="24" :lg="12">
        <el-form-item label="在地图上选择位置">
          <div ref="mapPickerRef" class="map-picker" />
        </el-form-item>
      </el-col>
    </el-row>

    <el-form-item class="form-actions">
      <el-button type="primary" :loading="loading" @click="handleSubmit">
        {{ isEdit ? "保存" : "创建" }}
      </el-button>
      <el-button @click="$emit('cancel')">取消</el-button>
    </el-form-item>
  </el-form>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch, onMounted, onUnmounted } from "vue";
import { ElMessage } from "element-plus";
import { Map, View } from "ol";
import { Tile } from "ol/layer";
import { XYZ } from "ol/source";
import { fromLonLat, toLonLat } from "ol/proj";
import { Feature } from "ol";
import { Point } from "ol/geom";
import VectorLayer from "ol/layer/Vector";
import VectorSource from "ol/source/Vector";
import { Style, Icon } from "ol/style";
import { Picture, Plus, Delete, Upload } from "@element-plus/icons-vue";
import type { FormInstance, FormRules } from "element-plus";
import { usePlacesStore } from "@/stores/places";
import type { Place, PlaceAddForm, PlaceUpdateForm } from "@/types";
import exifr from "exifr";

const props = defineProps<{
  place?: Place;
}>();

const emit = defineEmits<{
  submit: [data: PlaceAddForm | PlaceUpdateForm];
  cancel: [];
}>();

const placesStore = usePlacesStore();
const formRef = ref<FormInstance>();
const loading = ref(false);
const mapPickerRef = ref<HTMLElement>();
const map = ref<Map>();
const vectorSource = ref<VectorSource>();
const markerFeature = ref<Feature>();
const fileInputRef = ref<HTMLInputElement>();
const isCoordinateSet = ref(false);
const maxZoom = ref(14);

const isEdit = computed(() => !!props.place);
const existingGroups = computed(() => placesStore.groups);
const photoInput = ref("");

const form = reactive({
  name: "",
  desc: "",
  group: "",
  lat: 39.9042,
  lon: 116.4074,
  cover: "",
  photos: [] as string[],
});

const rules: FormRules = {
  name: [
    { required: true, message: "请输入地点名称", trigger: "blur" },
    { max: 50, message: "名称最多50个字符", trigger: "blur" },
  ],
};

// Initialize form if editing
watch(
  () => props.place,
  (place) => {
    if (place) {
      form.name = place.name || "";
      form.desc = place.desc;
      form.group = place.group;
      form.lat = place.lat;
      form.lon = place.lon;
      form.cover = place.cover || "";
      updateMarkerPosition();
      // 编辑模式：初始化时放大到 zoom
      updateMapCenter(true);
    }
  },
  { immediate: true },
);

// Watch for coordinate changes
watch(
  () => [form.lat, form.lon],
  () => {
    updateMarkerPosition();
  },
  { immediate: true },
);

onMounted(() => {
  if (!mapPickerRef.value) return;

  // Create vector source for marker
  vectorSource.value = new VectorSource();

  // Create marker feature
  markerFeature.value = new Feature({
    geometry: new Point(fromLonLat([form.lon, form.lat])),
  });

  markerFeature.value.setStyle(
    new Style({
      image: new Icon({
        src: "./static/location.png",
        width: 32,
        opacity: 1,
        anchor: [0.5, 1],
      }),
    }),
  );

  vectorSource.value.addFeature(markerFeature.value);

  // Create map
  map.value = new Map({
    target: mapPickerRef.value,
    layers: [
      new VectorLayer({
        source: vectorSource.value,
        zIndex: 100,
      }),
      new Tile({
        zIndex: 99,
        source: new XYZ({
          url: "https://webrd02.is.autonavi.com/appmaptile?lang=zh_cn&size=1&scale=1&style=8&x={x}&y={y}&z={z}",
        }),
        opacity: 0.8,
      }),
    ],
    view: new View({
      center: fromLonLat([form.lon, form.lat]),
      zoom: 4,
    }),
  });

  // Handle map click
  map.value.on("click", (event) => {
    const coordinate = toLonLat(event.coordinate);
    form.lon = parseFloat(coordinate[0].toFixed(6));
    form.lat = parseFloat(coordinate[1].toFixed(6));
    // 创建模式：用户点击地图选择坐标后，标记为已设置并放大
    if (!isEdit.value) {
      isCoordinateSet.value = true;
      updateMapCenter(true);
    }
  });

  // 编辑模式：地图初始化完成后放大到 zoom 14
  if (isEdit.value && props.place) {
    map.value.getView().setZoom(14);
  }
});

onUnmounted(() => {
  map.value?.dispose();
});

function updateMarkerPosition() {
  if (markerFeature.value) {
    const geometry = markerFeature.value.getGeometry() as Point;
    geometry.setCoordinates(fromLonLat([form.lon, form.lat]));
  }
}

function updateMapCenter(autoZoom = false) {
  const view = map.value?.getView();
  if (!view) return;

  view.setCenter(fromLonLat([form.lon, form.lat]));

  if (isEdit.value) {
    // 编辑模式：默认放大到 zoom
    view.setZoom(maxZoom.value);
  } else if (autoZoom && isCoordinateSet.value) {
    // 创建模式：用户选定坐标后放大到 zoom
    view.setZoom(maxZoom.value);
  }
  // 创建模式初始状态：不放大（保持默认 zoom 4）
}

function addPhoto() {
  const url = photoInput.value.trim();
  if (!url) {
    ElMessage.warning("请输入照片 URL");
    return;
  }
  if (!url.startsWith("http://") && !url.startsWith("https://")) {
    ElMessage.error("请输入有效的 HTTP/HTTPS 地址");
    return;
  }
  if (form.photos.includes(url)) {
    ElMessage.warning("该照片已添加");
    return;
  }
  form.photos.push(url);
  photoInput.value = "";
}

function removePhoto(index: number) {
  form.photos.splice(index, 1);
}

// Handle select photo button click
function handleSelectPhoto() {
  fileInputRef.value?.click();
}

// Handle file selection and extract EXIF GPS
async function handleFileChange(event: Event) {
  const target = event.target as HTMLInputElement;
  const file = target.files?.[0];

  if (!file) return;

  try {
    // Use exifr to parse GPS, it supports JPEG, TIFF, HEIC/HEIF formats
    const gps = await exifr.gps(file);

    if (!gps) {
      ElMessage.warning("该照片不包含 GPS 位置信息");
      return;
    }

    const { latitude, longitude } = gps;

    if (typeof latitude !== "number" || typeof longitude !== "number") {
      ElMessage.warning("GPS 坐标信息格式不正确");
      return;
    }

    // Update form coordinates
    form.lat = parseFloat(latitude.toFixed(6));
    form.lon = parseFloat(longitude.toFixed(6));

    // 创建模式：从照片读取坐标后，标记为已设置并放大
    if (!isEdit.value) {
      isCoordinateSet.value = true;
      updateMapCenter(true);
    } else {
      updateMapCenter();
    }

    ElMessage.success(`已从照片读取坐标：${form.lat.toFixed(4)}, ${form.lon.toFixed(4)}`);
  } catch (error) {
    console.error("EXIF parsing error:", error);
    ElMessage.error("读取照片信息失败，请检查文件格式");
  } finally {
    // Reset file input
    if (fileInputRef.value) {
      fileInputRef.value.value = "";
    }
  }
}

async function handleSubmit() {
  if (!formRef.value) return;

  await formRef.value.validate((valid) => {
    if (!valid) return;

    const submitData = isEdit.value
      ? {
          name: form.name,
          desc: form.desc,
          group: form.group,
          lat: form.lat,
          lon: form.lon,
          cover: form.cover,
        }
      : {
          name: form.name,
          desc: form.desc,
          group: form.group,
          lat: form.lat,
          lon: form.lon,
          cover: form.cover,
          photos: form.photos,
        };

    emit("submit", submitData as PlaceAddForm | PlaceUpdateForm);
  });
}
</script>

<style scoped>
.place-form {
  max-width: 1200px;
}

.coordinate-inputs {
  display: flex;
  gap: 12px;
}

.form-hint {
  font-size: 12px;
  color: #909399;
  margin: 4px 0 0;
}

.map-picker {
  width: 100%;
  height: 400px;
  border-radius: 8px;
  overflow: hidden;
  border: 1px solid #dcdfe6;
}

.form-actions {
  margin-top: 30px;
  padding-top: 20px;
  border-top: 1px solid #e4e7ed;
}

.cover-preview {
  margin-top: 12px;
}

.image-error {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
  background: #f5f7fa;
  color: #909399;
  font-size: 14px;
  gap: 8px;
}

.photos-input {
  margin-bottom: 12px;
}

.photos-list {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  margin-top: 12px;
}

.photo-tag {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px;
  background: #f5f7fa;
  border-radius: 8px;
  border: 1px solid #e4e7ed;
}

@media (max-width: 768px) {
  .coordinate-inputs {
    flex-direction: column;
  }

  .map-picker {
    height: 300px;
  }
}
</style>
