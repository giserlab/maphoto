import { reactive, toRefs, watch, onBeforeUnmount } from "vue";
import { Map, View } from "ol";
import { ZoomSlider, Zoom } from "ol/control";
import { Vector as VectorSource, Cluster } from "ol/source";
import { fromLonLat } from "ol/proj";
import { Tile } from "ol/layer";
import { XYZ } from "ol/source";
import { Style, Stroke } from "ol/style";
import { GeoJSON } from "ol/format";
// @ts-ignore
import Photo from "ol-ext/style/Photo";
// @ts-ignore
import AnimatedCluster from "ol-ext/layer/AnimatedCluster";
import { ElMessage } from "element-plus";

export function initialMap() {
  // 初始化参数
  let state = reactive({
    popup: {
      open: false,
      srcs: [],
    },
    mapConfig: {
      autoCenter: true,
      iconSize: null,
      lat: null,
      link: null,
      lon: null,
      zoom: 2,
      maxZoom: 16,
      minZoom: 3,
      note: null,
      title: "maphoto",
      tolerance: 40.0,
    },
    loading: false,
  });

  let map: Map;

  let clusterLyr: AnimatedCluster;

  function genStyle(size, src) {
    if (size <= 1) {
      return new Style({
        image: new Photo({
          src: src,
          radius: 20,
          crop: true,
          kind: "square",
          shadow: true,
          stroke: new Stroke({
            width: 3,
            color: "#fff",
          }),
        }),
      });
    } else {
      return new Style({
        image: new Photo({
          src: src,
          radius: 20,
          crop: true,
          kind: "folio",
          shadow: true,
          stroke: new Stroke({
            width: 3,
            color: "#fff",
          }),
        }),
      });
    }
  }

  //   创建地图及图层
  function createMap() {
    map = new Map({
      target: "map-container",
      layers: [
        new Tile({
          zIndex: 99,
          source: new XYZ({
            url: "https://webrd02.is.autonavi.com/appmaptile?lang=zh_cn&size=1&scale=1&style=8&x={x}&y={y}&z={z}",
            attributions: "© 高德地图",
          }),
          opacity: 0.8,
        }),
      ],
      view: new View({
        center: fromLonLat([110, 39]),
        zoom: state.mapConfig.zoom,
        maxZoom: state.mapConfig.maxZoom || 16,
        minZoom: state.mapConfig.minZoom,
      }),
      controls: [new ZoomSlider(), new Zoom()],
    });
    map.on("loadstart", loadStartEvtHandler);
    map.on("rendercomplete", loadCompleteEvtHandler);
    window.map = map;
    bindClickEvt();
  }

  function loadCompleteEvtHandler(evt) {
    state.loading = false;
  }
  function loadStartEvtHandler(evt) {
    state.loading = true;
  }

  onBeforeUnmount(() => {
    // 解除监听
    map.un("rendercomplete", loadCompleteEvtHandler);
    map.un("loadstart", loadCompleteEvtHandler);
    map.un("click");
  });

  //   加载点资源
  function loadPhoto(data, forceAdjust = false) {
    try {
      // 创建矢量数据源
      let vecSource = new VectorSource({
        title: "poi",
        features: new GeoJSON().readFeatures(data),
        wrapX: false,
      });
      var newClusterSource = new Cluster({
        distance: state.mapConfig.tolerance,
        source: vecSource,
      });
      clusterLyr = new AnimatedCluster({
        name: "maphoto",
        zIndex: 999,
        source: newClusterSource,
        // maxResolution: 40,
        style: (feature, resolution) => {
          var clusterFeats = feature.get("features");
          const size = clusterFeats.length;
          const iconURL = clusterFeats[0].get("cover");
          return genStyle(size, iconURL);
        },
      });

      map.addLayer(clusterLyr);
      //   定位图层
      if (forceAdjust) {
        map.getView().fit(vecSource.getExtent(), map.getSize());
      } else {
        if (state.mapConfig.autoCenter) {
          map.getView().fit(vecSource.getExtent(), map.getSize());
        } else {
          let vw = map.getView();
          vw.setCenter(fromLonLat([state.mapConfig.lon, state.mapConfig.lat]));
          vw.setZoom(state.mapConfig.maxZoom);
        }
      }
    } catch (err) {
      ElMessage.warning("无照片");
      console.error(err);
    }
  }

  //   绑定点击事件
  function bindClickEvt() {
    map.on("click", (event) => {
      state.popup.srcs.length = 0;
      clusterLyr
        .getFeatures(event.pixel)
        .then((clusterFeat) => {
          if (clusterFeat.length > 0) {
            const features = clusterFeat[0].get("features");
            if (features.length > 0) {
              features.forEach((feat) => {
                const photos = feat.get("photos");
                const date = feat.get("date").slice(0, 10);
                if (photos.length > 0) {
                  state.popup.open = true;
                  photos.forEach((photo) => {
                    const src = photo.url;
                    state.popup.srcs.push({ date, src });
                  });
                }
              });
            }
          } else {
            state.popup.open = false;
          }
        })
        .catch((error) => {
          console.error("错误:" + error);
        });
    });
  }

  watch(
    () => state.mapConfig,
    (val) => {
      if (!map) return;
      let view = map.getView();
      if (val.zoom) view.setZoom(val.zoom);
      if (val.maxZoom) view.setMaxZoom(val.maxZoom);
      if (val.minZoom) view.setMinZoom(val.minZoom);
    },
    { deep: true, immediate: true },
  );

  return {
    createMap,
    loadPhoto,
    ...toRefs(state),
  };
}

export default initialMap;
