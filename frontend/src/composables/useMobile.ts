import { ref, computed, onMounted, onUnmounted } from 'vue'

export function useMobile() {
  const windowWidth = ref(window.innerWidth)

  const isMobile = computed(() => windowWidth.value < 768)
  const isTablet = computed(() => windowWidth.value >= 768 && windowWidth.value < 1024)
  const isDesktop = computed(() => windowWidth.value >= 1024)

  const updateWidth = () => {
    windowWidth.value = window.innerWidth
  }

  onMounted(() => {
    window.addEventListener('resize', updateWidth)
  })

  onUnmounted(() => {
    window.removeEventListener('resize', updateWidth)
  })

  return {
    isMobile,
    isTablet,
    isDesktop,
    windowWidth,
  }
}
