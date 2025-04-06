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
    <select
      id="filter"
      class="border text-sm rounded-lg p-2.5 bg-slate-600 border-slate-500 text-white focus:ring-blue-500 focus:border-blue-500"
      v-model="selectedFilter"
      @change="emitFilter"
    >
      <option value="0">Todos</option>
      <option value="1">Avaliados</option>
      <option value="2">Anotados</option>
    </select>
  </div>
</template>
