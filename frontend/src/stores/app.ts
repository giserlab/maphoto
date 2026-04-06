import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useAppStore = defineStore('app', () => {
  const sidebarCollapsed = ref(false)
  const isMobile = ref(false)

  function toggleSidebar() {
    sidebarCollapsed.value = !sidebarCollapsed.value
  }

  function setIsMobile(value: boolean) {
    isMobile.value = value
  }

  return {
    sidebarCollapsed,
    isMobile,
    toggleSidebar,
    setIsMobile,
  }
})
