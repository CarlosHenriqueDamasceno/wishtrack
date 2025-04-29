<script setup lang="ts">
import Cross from '@/components/icons/Cross.vue'

defineProps<{
  title: string
  show: boolean
  id: string
}>()

const emit = defineEmits(['close'])

function close(event: Event) {
  emit('close', event)
}
</script>

<template>
  <div
    :id="id"
    tabindex="-1"
    aria-hidden="true"
    :class="show === false ? 'opacity-0 pointer-events-none' : 'opacity-100 pointer-events-auto'"
    class="flex overflow-y-auto overflow-x-hidden fixed top-0 right-0 left-0 z-50 justify-center items-center w-full md:inset-0 h-[calc(100%-1rem)] max-h-full bg-slate-800/30 backdrop-blur-md transition-all duration-50"
  >
    <div class="relative p-4 w-full max-w-lg max-h-full">
      <div class="bg-zinc-950 border border-zinc-800 rounded-xl relative shadow-sm">
        <div
          class="flex items-center justify-between p-5 rounded-t dark:border-slate-600 border-slate-200"
        >
          <h2 class="text-xl font-semibold text-white">
            {{ title }}
          </h2>
          <button
            @click="close"
            type="button"
            class="end-2.5 text-slate-400 bg-transparent cursor-pointer hover:bg-white/10 rounded-lg text-sm w-8 h-8 ms-auto inline-flex justify-center items-center"
            :data-modal-toggle="id"
          >
            <Cross class="w-10 h-10" />
            <span class="sr-only">Close modal</span>
          </button>
        </div>
        <div class="px-5">
          <slot></slot>
        </div>
        <div class="p-5 dark:border-slate-600 border-slate-200 w-full">
          <slot name="footer"></slot>
        </div>
      </div>
    </div>
  </div>
</template>
