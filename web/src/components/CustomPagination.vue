<script setup lang="ts">
import { computed } from 'vue'
import CustomButton from './CustomButton.vue'

const props = defineProps({
  totalItems: {
    type: Number,
    required: true,
  },
  itemsPerPage: {
    type: Number,
    default: 10,
  },
  currentPage: {
    type: Number,
    default: 1,
  },
})

const totalPages = computed(() => {
  return Math.ceil(props.totalItems / props.itemsPerPage)
})

const pages = computed(() => {
  return Array.from({ length: totalPages.value }, (_, i) => i + 1)
})

const emit = defineEmits(['update:currentPage'])

const goToPage = (page: number) => {
  if (page < 1 || page > pages.value.length) {
    return
  }
  if (page !== props.currentPage) {
    emit('update:currentPage', page)
  }
}
</script>
<template>
  <nav class="flex items-center justify-center space-x-2" aria-label="Pagination">
    <CustomButton :disabled="currentPage === 1" @click="goToPage(currentPage - 1)">
      Anterior
    </CustomButton>
    <CustomButton
      v-for="page in pages"
      :key="page"
      class="px-4 py-2 text-sm font-medium border rounded-md"
      :theme="page === currentPage ? 'default' : 'secondary'"
      @click="goToPage(page)"
    >
      {{ page }}
    </CustomButton>
    <CustomButton :disabled="currentPage === totalPages" @click="goToPage(currentPage + 1)">
      Próximo
    </CustomButton>
  </nav>
</template>
