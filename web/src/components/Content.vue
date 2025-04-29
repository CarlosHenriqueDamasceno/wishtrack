<script setup lang="ts">
import { useRatingStore } from '@/stores/rating'
import type Content from '@/types/content'
import { Status, statusName } from '@/types/status'
import { computed, defineProps } from 'vue'
import Plus from './icons/Plus.vue'
import StarOutline from './icons/StarOutline.vue'
import Star from './icons/Star.vue'
import { useDeletingStore } from '@/stores/delete'
import { useWriteDownStore } from '@/stores/writedown'
import { useModalStore } from '@/stores/modal'

const props = defineProps<{
  content: Content
}>()

const status = computed(() => {
  if (undefined === props.content.id) {
    return Status.UNSAVED
  }
  if (null === props.content.rate) {
    return Status.SAVED
  }

  return Status.RATED
})

const modalStore = useModalStore()

function handleClick() {
  modalStore.content = props.content

  if (status.value === Status.UNSAVED) {
    modalStore.type = 'write-down'
  }

  if (status.value === Status.SAVED || status.value === Status.RATED) {
    modalStore.type = 'rating'
  }

  modalStore.toggleModal()
}
</script>
<template>
  <div
    @click="handleClick"
    class="cursor-pointer min-w-[250px] aspect-[16/9] bg-zinc-900 rounded-2xl border border-zinc-800 shadow-md p-4 flex flex-col justify-between hover:scale-105 transition-transform"
  >
    <div>
      <div class="mb-0.5 rounded-sm text-white/50 text-xs flex justify-between">
        {{ props.content.category }}
        <button
          type="button"
          class="p-2 end-2.5 text-slate-400 cursor-pointer bg-white/10 rounded-full text-sm ms-auto inline-flex justify-center items-center"
        >
          <Plus v-if="status === Status.UNSAVED" class="w-4 h-4" />
          <StarOutline v-if="status === Status.SAVED" class="w-4 h-4" />
          <Star v-if="status === Status.RATED" class="w-4 h-4" />
        </button>
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
