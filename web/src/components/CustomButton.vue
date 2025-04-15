<script setup lang="ts">
import { computed, defineProps } from 'vue'

const props = defineProps({
  theme: {
    type: String,
    default: 'default',
    validator: (value: string) => ['default', 'danger', 'secondary'].includes(value),
  },
  disabled: {
    type: Boolean,
    default: false,
  },
})

const buttonClass = computed(() => {
  if (props.theme === 'danger') {
    return 'bg-red-700 text-white border-red-700 hover:bg-red-800 hover:border-red-800 focus:ring-red-100'
  }

  if (props.theme === 'default') {
    return 'bg-white text-black hover:bg-zinc-200'
  }

  if (props.theme === 'secondary') {
    return 'bg-blue-100 text-gray-900 border-gray-200 hover:bg-gray-300 hover:border-gray-300 focus:ring-gray-100'
  }
})
</script>

<template>
  <button
    :disabled="disabled"
    type="button"
    class="px-4 py-2 cursor-pointer font-semibold rounded-md disabled:opacity-50 disabled:cursor-not-allowed"
    :class="buttonClass"
  >
    <slot></slot>
  </button>
</template>
