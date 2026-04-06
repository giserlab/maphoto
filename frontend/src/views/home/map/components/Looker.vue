<template>
  <div class="looker-wrapper">
    <div class="looker">
      <div class="header">
        <div class="date"></div>
        <div class="close btn" @click="this.$emit('close')">
          <span class="icon">✖️</span>
        </div>
      </div>
      <div class="looker-inner">
        <img
          :src="activeSrc"
          v-loading="loading"
          @load="() => (loading = false)"
          @error="() => (loading = false)"
        />
      </div>
      <div class="img-handle">
        <div class="tip">
          <div
            :class="['tip-item', k === currentSrcIndex ? 'active' : '']"
            v-for="(item, k) in pics"
            :key="k.src"
          ></div>
        </div>
        <div class="btn-group" v-if="btnSeen">
          <div class="btn left" @click="backwardLook">
            <span class="icon">◀</span>
          </div>
          <div class="btn right" @click="towardLook">
            <span class="icon">▶</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "Looker",
  props: {
    pics: {
      type: Array,
      default: () => {
        return [
          {
            date: "2024-01-01",
            src: "https://md-1301600412.cos.ap-nanjing.myqcloud.com/pic/image-20220328152305216.png",
          },
        ];
      },
    },
  },
  data() {
    return {
      currentSrcIndex: 0,
      loading: true,
    };
  },
  computed: {
    btnSeen() {
      return this.pics.length > 1;
    },
    activeSrc() {
      return this.pics[this.currentSrcIndex].src;
    },
  },
  methods: {
    towardLook() {
      this.loading = true;
      if (this.currentSrcIndex < this.pics.length - 1) {
        this.currentSrcIndex += 1;
      } else {
        this.currentSrcIndex = 0;
      }
    },
    backwardLook() {
      this.loading = true;
      if (this.currentSrcIndex > 0) {
        this.currentSrcIndex -= 1;
      } else {
        this.currentSrcIndex = this.pics.length - 1;
      }
    },
  },
};
</script>

<style scoped lang="scss">
$--btn-size: 0.2em;
$--btn-color: rgb(18, 150, 219);
.looker-wrapper {
  position: absolute;
  left: 0em;
  top: 0em;
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100vw;
  height: 100dvh;
  pointer-events: none;
}
.looker {
  pointer-events: all;
  display: flex;
  flex-direction: column;
  z-index: 1000;
  position: relative;
  width: var(--looker-width);
  min-height: var(--looker-height);
  max-height: 80%;
  background: #ffffffe9;
  box-shadow: 1px 1px 8px 4px #454343;
  .header {
    position: relative;
    display: flex;
    justify-content: space-between;
    width: 100%;
    padding: 0.5em;
    .date {
      display: flex;
      justify-content: center;
      align-items: center;
      width: 100%;
      font-weight: bold;
    }
    .close {
      color: $--btn-color;
      cursor: pointer;
    }
  }
  .looker-inner {
    flex: 1;
    display: flex;
    flex-direction: row;
    justify-content: center;
    align-items: flex-start;
    padding: 0.5em;
    overflow-y: scroll;
    img {
      flex: 1;
      display: block;
      width: 95%;
      height: 60%;
      max-height: 80%;
      margin: 0 auto;
      object-fit: cover;
      border-radius: 0.2em;
      background-color: #fff;
    }
  }
  .btn {
    width: 2em;
    z-index: 999;
    padding: 0.5em;
    cursor: pointer;
    display: flex;
    flex-direction: row;
    justify-content: center;
    align-items: center;
    .icon {
      font-size: 1em;
      color: $--btn-color !important;
      &:active,
      &:hover {
        transform: scale(1.2);
      }
    }
  }
  .img-handle {
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    align-items: center;
    gap: 0.5em;
    width: 100%;
    padding: 0.5em 1em;
    box-sizing: border-box;
    .btn-group {
      width: 100%;
      display: flex;
      justify-content: space-around;
      .icon {
        font-size: 2em;
      }
    }

    .tip {
      flex: 1;
      display: flex;
      flex-direction: row;
      justify-content: center;
      align-items: center;
      flex-wrap: wrap;
      .tip-item {
        width: 0.5em;
        height: 0.5em;
        margin: 0.5em;
        cursor: pointer;
        border-radius: 0.25em;
        background: rgba(210, 210, 211, 0.805);
        &.active {
          background: $--btn-color;
        }
        &:hover {
          opacity: 0.5;
        }
      }
    }
  }
}
</style>
