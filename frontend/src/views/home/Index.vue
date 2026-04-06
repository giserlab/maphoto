<template>
  <div class="main" :style="mainStyle">
    <div class="map-out">
      <Map />
    </div>
  </div>
</template>

<script setup>
import Map from "./map/index.vue";
import { computed } from "vue";
import { reactive } from "vue";
import { onMounted } from "vue";

const mainStyle = computed(() => {
  return {
    width: "100vw",
    height: state.screenHeight + "px",
  };
});

const state = reactive({
  screenHeight: 0,
});

onMounted(() => (state.screenHeight = window.innerHeight));
</script>

<style scoped lang="scss">
.main {
  display: flex;
  flex-direction: column;
  width: 100%;
  margin: 0px auto;
  background: rgb(35, 38, 55);
  .header-out {
    width: 100%;
    height: 3em;
  }
  .map-out {
    flex: 1;
    width: 100%;
  }
  .footer-out {
    width: 100%;
    height: 6em;
  }
}
</style>
