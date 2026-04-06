<!--
  This component is kept for compatibility but is no longer used.
  Markers are now rendered directly in MapContainer.vue using OpenLayers VectorLayer.
-->
<template>
  <!-- Empty template - markers rendered by parent -->
</template>

<script setup lang="ts">
import type { Place } from "@/types";

// Props are kept for API compatibility
defineProps<{
  place: Place;
  selected?: boolean;
}>();

defineEmits<{
  click: [place: Place];
}>();
</script>
