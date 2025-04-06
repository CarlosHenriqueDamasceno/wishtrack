<script setup lang="ts">
import { ref, defineEmits } from 'vue'

const genreInput = ref('')
const emit = defineEmits(['update:genres'])

const updateGenres = () => {
  const genres = genreInput.value
    .split(',')
    .map((genre) => genre.trim())
    .filter((genre) => genre !== '')
  emit('update:genres', genres)
}
</script>

<template>
  <div class="genre-filter flex items-center justify-center">
    <input
      type="text"
      id="genres"
      v-model="genreInput"
      placeholder="Gêneros separados por vírgula"
      class="border text-sm rounded-lg min-w-70 p-2.5 block bg-slate-600 border-slate-500 placeholder-slate-400 text-white"
      @input="updateGenres"
    />
  </div>
</template>
