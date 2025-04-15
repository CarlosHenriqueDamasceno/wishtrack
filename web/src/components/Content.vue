<script setup lang="ts">
import { useRatingStore } from '@/stores/rating'
import type Content from '@/types/content'
import { Status, statusName } from '@/types/status'
import { computed, defineProps } from 'vue'

const props = defineProps<{
  content: Content
}>()

const status = computed(() => {
  if (undefined === props.content.id) {
    return Status.UNSAVED
  }
  if (undefined === props.content.rate) {
    return Status.SAVED
  }

  return Status.RATED
})

const store = useRatingStore()

function handleClick() {
  store.content = props.content
  store.toggleModal(true)
}
</script>
<template>
  <div
    @click="handleClick"
    class="cursor-pointer min-w-[250px] aspect-[16/9] bg-zinc-900 rounded-2xl border border-zinc-800 shadow-md p-4 flex flex-col justify-between hover:scale-105 transition-transform"
  >
    <div>
      <div class="mb-0.5 rounded-sm text-white/50 text-xs">
        {{ props.content.category }}
        {{ statusName[status] }}
      </div>
      <h3 class="text-lg font-semibold text-white mb-2">{{ props.content.name }}</h3>
      <p class="text-sm text-zinc-400 mb-3 truncate-lines-3">
        {{ props.content.summary }}
      </p>
    </div>
    <div class="flex flex-wrap gap-1">
      <span
        v-for="genre in props.content.genres"
        class="bg-zinc-800 text-xs text-zinc-300 px-2 py-1 rounded-full"
        >{{ genre }}</span
      >
    </div>
  </div>
</template>

<style scoped>
.truncate-lines-3 {
  display: -webkit-box;
  -webkit-line-clamp: 3; /* Change this number for different line limits */
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
