<script setup lang="ts">
import { computed, defineProps } from 'vue'

const props = defineProps({
  theme: {
    type: String,
    default: 'default',
    validator: (value: string) => ['default', 'danger'].includes(value),
  },
})

const buttonClass = computed(() => {
  return props.theme === 'danger'
    ? 'bg-red-700 text-white border-red-700 hover:bg-red-800 hover:border-red-800 focus:ring-red-100'
    : 'bg-blue-500 text-white border-blue-500 hover:bg-blue-700 hover:border-blue-700 focus:ring-blue-100'
})
</script>

<template>
  <button
    :class="buttonClass"
    type="button"
    class="border cursor-pointer font-bold rounded-lg text-sm px-5 py-2.5 focus:outline-none focus:ring-4"
  >
    <slot></slot>
  </button>
</template>
