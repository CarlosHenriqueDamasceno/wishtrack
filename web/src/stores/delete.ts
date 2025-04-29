import type Content from '@/types/content'
import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useDeletingStore = defineStore('deleting', () => {
  const content = ref<Content | null>(null)
  const isModalVisible = ref(false)

  function toggleModal(visibility?: boolean) {
    isModalVisible.value = visibility ?? !isModalVisible.value
  }

  return { content, isModalVisible, toggleModal }
})
