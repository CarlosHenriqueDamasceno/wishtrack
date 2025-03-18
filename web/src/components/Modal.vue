<script setup lang="ts">
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
    <div class="relative p-4 w-full max-w-md max-h-full">
      <div class="relative bg-white rounded-lg shadow-sm dark:bg-slate-700">
        <div
          class="flex items-center justify-between p-4 md:p-5 rounded-t dark:border-slate-600 border-slate-200"
        >
          <h3 class="text-xl font-semibold text-slate-900 dark:text-white">
            {{ title }}
          </h3>
          <button
            @click="close"
            type="button"
            class="end-2.5 text-slate-400 bg-transparent cursor-pointer hover:bg-slate-200 hover:text-slate-900 rounded-lg text-sm w-8 h-8 ms-auto inline-flex justify-center items-center dark:hover:bg-slate-600 dark:hover:text-white"
            :data-modal-toggle="id"
          >
            <svg
              class="w-3 h-3"
              aria-hidden="true"
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 14 14"
            >
              <path
                stroke="currentColor"
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="m1 1 6 6m0 0 6 6M7 7l6-6M7 7l-6 6"
              />
            </svg>
            <span class="sr-only">Close modal</span>
          </button>
        </div>
        <div class="p-4">
          <slot></slot>
        </div>
        <div class="p-4 md:p-5 dark:border-slate-600 border-slate-200 w-full">
          <slot name="footer"></slot>
        </div>
      </div>
    </div>
  </div>
</template>
