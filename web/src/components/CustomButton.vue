<script setup lang="ts">
import { computed, defineProps } from 'vue';

const props = defineProps({
  type: {
    type: String,
    default: 'default',
    validator: (value:string) => ['default', 'danger'].includes(value)
  }
});

const buttonClass = computed(() => {
  return props.type === 'danger'
    ? 'bg-rose-600 text-white border-rose-300 hover:bg-rose-700 hover:border-rose-600 focus:ring-rose-700'
    : 'hover:bg-gray-100 focus:ring-gray-100 dark:bg-gray-800 dark:text-white dark:border-gray-600 dark:hover:bg-gray-700 dark:hover:border-gray-600 dark:focus:ring-gray-700';
});
</script>

<template>
  <button :class="buttonClass" type="button"
    class="border cursor-pointer font-medium rounded-lg text-sm px-5 py-2.5 focus:outline-none focus:ring-4">
    <slot></slot>
  </button>
</template>
