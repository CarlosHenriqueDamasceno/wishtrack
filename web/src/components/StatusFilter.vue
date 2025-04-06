<script setup lang="ts">
import { ref } from 'vue'
import { Status } from '@/types/status'

const selectedFilter = ref<Status>(Status.NONE)
const emit = defineEmits(['filter-selected'])

const emitFilter = () => {
  emit('filter-selected', selectedFilter.value)
}

emit('filter-selected', selectedFilter.value)
</script>
<template>
  <div class="flex items-center">
    <label for="filter" class="text-sm font-medium text-white mr-3">Status:</label>
    <select
      id="filter"
      class="border text-sm rounded-lg p-2.5 bg-gray-700 border-gray-600 placeholder-gray-400 text-white focus:ring-blue-500 focus:border-blue-500"
      v-model="selectedFilter"
      @change="emitFilter"
    >
      <option value="0">Todos</option>
      <option value="1">Avaliados</option>
      <option value="2">Anotados</option>
    </select>
  </div>
</template>
