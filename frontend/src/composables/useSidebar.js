import { ref } from 'vue'

// Global reactive state untuk menyimpan status collapse sidebar
const isCollapsed = ref(false)

export function useSidebar() {
  const toggleSidebar = () => {
    isCollapsed.value = !isCollapsed.value
  }

  return {
    isCollapsed,
    toggleSidebar
  }
}
