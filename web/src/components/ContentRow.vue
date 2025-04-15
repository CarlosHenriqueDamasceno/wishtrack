<script setup lang="ts">
import ContentCard from '@/components/Content.vue'
import type Content from '@/types/content'
import { defineProps } from 'vue'

interface Provider {
  provider: string
  suggestions: Content[]
}

const props = defineProps<{
  provider: Provider
}>()
</script>

<template>
  <section>
    <h2 class="text-xl font-semibold mb-4">{{ props.provider.provider }}</h2>
    <div class="flex space-x-4 pb-2">
      <ContentCard
        v-for="(content, key) in props.provider.suggestions"
        :key="content.id ?? key"
        :content="content"
      />
    </div>
  </section>
</template>

<style scoped>
::-webkit-scrollbar {
  display: none;
}
</style>
