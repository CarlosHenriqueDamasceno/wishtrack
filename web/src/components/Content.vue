<script setup lang="ts">
import type Content from '@/types/content'
import { defineProps } from 'vue'

const props = defineProps<{
  content: Content
}>()
</script>
<template>
  <div
    class="cursor-pointer block max-w-md bg-white border border-slate-200 rounded-lg shadow-sm dark:bg-slate-800 dark:border-slate-700"
  >
    <div class="px-6 py-4">
      <h5 class="mb-2 text-2xl font-bold tracking-tight text-slate-900 dark:text-white">
        {{ props.content.name }}
      </h5>
      <span
        class="bg-slate-100 text-slate-800 text-xs font-medium px-2.5 py-0.5 rounded-sm dark:bg-slate-700 dark:text-slate-300"
      >
        {{ props.content.category }}
      </span>
      <p class="font-normal text-slate-700 dark:text-slate-400 mt-2">
        {{ props.content.summary }}
      </p>
    </div>
    <div class="px-6 pt-4 pb-2">
      <span
        v-for="genre in props.content.genres"
        class="inline-block bg-slate-200/20 rounded-full px-3 py-1 text-sm font-semibold text-white mr-2 mb-2"
      >
        #{{ genre }}
      </span>
    </div>
  </div>
</template>
