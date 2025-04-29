import type Content from '@/types/content'
import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useModalStore = defineStore('modal', () => {
  const content = ref<Content | null>(null)
  const isVisible = ref(false)
  const type = ref<'write-down' | 'delete' | 'rating'>('write-down')

  function toggleModal(visibility?: boolean) {
    isVisible.value = visibility ?? !isVisible.value
    if (false === isVisible.value) {
      content.value = null
    }
  }

  return { content, isVisible, toggleModal, type }
})
