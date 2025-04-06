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
    return 'bg-blue-500 text-white border-blue-500 hover:bg-blue-700 hover:border-blue-700 focus:ring-blue-100'
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
    class="border cursor-pointer font-bold rounded-lg text-sm px-5 py-2.5 focus:outline-none focus:ring-4 disabled:opacity-50 disabled:cursor-not-allowed"
    :class="buttonClass"
  >
    <slot></slot>
  </button>
</template>
