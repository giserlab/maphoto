<template>
  <div ref="mapRef" class="map-container" :class="{ 'is-mobile': isMobile }">
    <!-- Map controls -->
    <MapControls @locate="handleLocate" @zoom-in="zoomIn" @zoom-out="zoomOut" />

    <!-- Mobile bottom sheet for place details -->
    <transition name="slide-up">
      <div v-if="isMobile && selectedPlace" class="mobile-place-sheet">
        <div class="sheet-header" @click="selectedPlace = null">
          <div class="drag-handle" />
        </div>
        <PlaceCard :place="selectedPlace" compact />
      </div>
    </transition>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted } from "vue";
import { Map, View } from "ol";
import { Tile } from "ol/layer";
import { XYZ } from "ol/source";
import { fromLonLat, toLonLat } from "ol/proj";
import { Feature } from "ol";
import { Point } from "ol/geom";
import VectorLayer from "ol/layer/Vector";
import VectorSource from "ol/source/Vector";
import { Style, Icon } from "ol/style";
import MapControls from "./MapControls.vue";
import PlaceCard from "@/components/place/PlaceCard.vue";
import { useMobile } from "@/composables/useMobile";
import type { Place, Config } from "@/types";

const props = defineProps<{
  places: Place[];
  config?: Config | null;
}>();

const emit = defineEmits<{
  select: [place: Place];
  "update:location": [lat: number, lon: number];
}>();

const { isMobile } = useMobile();
const mapRef = ref<HTMLElement>();
const map = ref<Map>();
const vectorSource = ref<VectorSource>();
const selectedPlace = ref<Place | null>(null);

const center = computed(() => {
  const lon = props.config?.lon || 116.4074;
  const lat = props.config?.lat || 39.9042;
  return fromLonLat([lon, lat]);
});

const zoom = computed(() => props.config?.zoom || 8);

// Watch for config changes
watch(
  () => props.config,
  (newConfig) => {
    if (newConfig && map.value) {
      const view = map.value.getView();
      view.setCenter(fromLonLat([newConfig.lon, newConfig.lat]));
      view.setZoom(newConfig.zoom);
    }
  },
  { deep: true },
);

// Watch for places changes
watch(
  () => props.places,
  (newPlaces) => {
    updateMarkers(newPlaces);
  },
  { deep: true },
);

onMounted(() => {
  if (!mapRef.value) return;

  // Create vector source for markers
  vectorSource.value = new VectorSource();

  // Create vector layer
  const vectorLayer = new VectorLayer({
    zIndex: 100,
    source: vectorSource.value,
  });

  // Create map
  map.value = new Map({
    target: mapRef.value,
    layers: [
      new Tile({
        source: new XYZ({
          url: "https://webrd02.is.autonavi.com/appmaptile?lang=zh_cn&size=1&scale=1&style=8&x={x}&y={y}&z={z}",
        }),
        zIndex: 99,
        opacity: 0.8,
      }),
      vectorLayer,
    ],
    view: new View({
      center: center.value,
      zoom: zoom.value,
      minZoom: props.config?.minZoom || 3,
      maxZoom: props.config?.maxZoom || 18,
    }),
  });

  // Handle map click
  map.value.on("click", (event) => {
    const coordinate = toLonLat(event.coordinate);
    emit("update:location", coordinate[1], coordinate[0]);

    // Check if clicked on a feature
    const feature = map.value?.forEachFeatureAtPixel(event.pixel, (f) => f);
    if (feature) {
      const placeId = feature.get("placeId");
      const place = props.places.find((p) => p.id === placeId);
      if (place) {
        selectPlace(place);
      }
    }
  });

  // Initial markers
  updateMarkers(props.places);
});

onUnmounted(() => {
  map.value?.dispose();
});

function updateMarkers(places: Place[]) {
  if (!vectorSource.value) return;

  vectorSource.value.clear();

  places.forEach((place) => {
    const feature = new Feature({
      geometry: new Point(fromLonLat([place.lon, place.lat])),
      placeId: place.id,
    });

    feature.setStyle(
      new Style({
        image: new Icon({
          src: "./static/location.png",
          width: 32,
          opacity: 1,
          anchor: [0.5, 1],
        }),
      }),
    );

    vectorSource.value?.addFeature(feature);
  });
}

function selectPlace(place: Place) {
  selectedPlace.value = place;
  const view = map.value?.getView();
  if (view) {
    view.setCenter(fromLonLat([place.lon, place.lat]));
  }
  emit("select", place);
}

function handleLocate() {
  if (navigator.geolocation) {
    navigator.geolocation.getCurrentPosition((position) => {
      const view = map.value?.getView();
      if (view) {
        view.setCenter(fromLonLat([position.coords.longitude, position.coords.latitude]));
        view.setZoom(15);
      }
    });
  }
}

function zoomIn() {
  const view = map.value?.getView();
  if (view) {
    const currentZoom = view.getZoom() || 10;
    view.setZoom(currentZoom + 1);
  }
}

function zoomOut() {
  const view = map.value?.getView();
  if (view) {
    const currentZoom = view.getZoom() || 10;
    view.setZoom(currentZoom - 1);
  }
}
</script>

<style scoped>
.map-container {
  width: 100%;
  height: 100%;
  position: relative;
}

.map-container :deep(.ol-map) {
  width: 100%;
  height: 100%;
}

.mobile-place-sheet {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  background: #fff;
  border-radius: 16px 16px 0 0;
  box-shadow: 0 -4px 20px rgba(0, 0, 0, 0.1);
  z-index: 1000;
  padding: 8px 16px 16px;
}

.sheet-header {
  display: flex;
  justify-content: center;
  padding: 8px 0;
  cursor: pointer;
}

.drag-handle {
  width: 40px;
  height: 4px;
  background: #dcdfe6;
  border-radius: 2px;
}

.slide-up-enter-active,
.slide-up-leave-active {
  transition: transform 0.3s ease;
}

.slide-up-enter-from,
.slide-up-leave-to {
  transform: translateY(100%);
}
</style>
