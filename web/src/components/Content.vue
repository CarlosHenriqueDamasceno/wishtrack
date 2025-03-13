<script setup lang="ts">
import { defineProps } from 'vue';

type Content = {
  id: string,
  name: string,
  genres: string[],
  summary: string,
  category: string,
}

const props = defineProps<{
  content:Content
}>()

</script>
<template>
  <div
    class="cursor-pointer block max-w-md bg-white border border-gray-200 rounded-lg shadow-sm  dark:bg-gray-800 dark:border-gray-700">
    <div class="px-6 py-4">
      <h5 class="mb-2 text-2xl font-bold tracking-tight text-gray-900 dark:text-white">{{props.content.name}}</h5>
      <span
        class="bg-gray-100 text-gray-800 text-xs font-medium px-2.5 py-0.5 rounded-sm dark:bg-gray-700 dark:text-gray-300">
        {{props.content.category}}
      </span>
      <p class="font-normal text-gray-700 dark:text-gray-400 mt-2">
        {{props.content.summary}}
      </p>
    </div>
    <div class="px-6 pt-4 pb-2">
      <span v-for="genre in props.content.genres" class="inline-block bg-gray-200 rounded-full px-3 py-1 text-sm font-semibold text-gray-700 mr-2 mb-2">
        #{{genre}}
      </span>
    </div>
  </div>
</template>
