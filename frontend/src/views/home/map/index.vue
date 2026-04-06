<template>
  <div style="width: 100%; height: 100%">
    <Loading v-show="loading" top="1.5vw" right="2vw" />
    <div id="map-container"></div>
    <Transition>
      <Looker v-if="popup.open" :pics="popup.srcs" @close="popup.open = false" />
    </Transition>
  </div>
</template>

<script>
import { initialMap } from "./index.js";
import Looker from "./components/Looker.vue";
import Loading from "@/components/Loading.vue";
import * as shareApi from "@/api/share";
import { ElMessage } from "element-plus";

export default {
  name: "Map",
  components: { Looker, Loading },
  setup() {
    return initialMap();
  },
  mounted() {
    let name = this.$route.query.user;
    const group = this.$route.query.group || "";
    name = name ? name : "admin";
    this.createMap();
    shareApi
      .getSharedPlaces({ username: name, group: group || "" })
      .then((data) => {
        if (data.status) {
          Object.assign(this.mapConfig, data.data.config);
          this.loadPhoto(data.data.features, group ? true : false);
        } else throw data.msg;
      })
      .catch((err) => {
        ElMessage.warning(err);
      });
  },
};
</script>

<style scoped lang="scss">
.v-enter-active,
.v-leave-active {
  transition: opacity 0.7s ease;
}

.v-enter-from,
.v-leave-to {
  opacity: 0;
}

#map-container {
  width: 100%;
  height: 100%;
  &:active {
    cursor: pointer;
  }
  &:visited {
    cursor: pointer;
  }
}
</style>
