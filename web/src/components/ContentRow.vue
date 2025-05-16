<script setup lang="ts">
import { ref, onMounted } from 'vue'
import ContentCard from '@/components/Content.vue'
import type Provider from '@/types/provider'
import { defineProps } from 'vue'

const props = defineProps<{
  provider: Provider
}>()

const scrollRef = ref<HTMLElement | null>(null)
const canScrollLeft = ref(false)
const canScrollRight = ref(false)

function updateArrows() {
  const el = scrollRef.value
  if (!el) return
  canScrollLeft.value = el.scrollLeft > 0
  canScrollRight.value = el.scrollLeft + el.clientWidth < el.scrollWidth - 1
}

function scrollByAmount(amount: number) {
  const el = scrollRef.value
  if (!el) return
  el.scrollBy({ left: amount, behavior: 'smooth' })
}

function scrollLeft() {
  const el = scrollRef.value
  if (!el) return
  scrollByAmount(-el.clientWidth * 0.8)
}

function scrollRight() {
  const el = scrollRef.value
  if (!el) return
  scrollByAmount(el.clientWidth * 0.8)
}

onMounted(() => {
  updateArrows()
  if (scrollRef.value) {
    scrollRef.value.addEventListener('scroll', updateArrows)
    window.addEventListener('resize', updateArrows)
  }
})
</script>

<template>
  <section>
    <h2 class="text-xl font-semibold mb-4">{{ props.provider.provider }}</h2>
    <div class="relative group">
      <button
        v-if="canScrollLeft"
        @click="scrollLeft"
        class="absolute left-0 top-1/2 -translate-y-1/2 z-20 bg-black/60 hover:bg-black/80 text-white rounded-full w-10 h-10 flex items-center justify-center transition-opacity opacity-0 group-hover:opacity-100"
        aria-label="Rolar para a esquerda"
      >
        <svg width="24" height="24" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><path d="M15 19l-7-7 7-7"/></svg>
      </button>
      <div
        ref="scrollRef"
        class="overflow-x-auto -mx-4 scroll-smooth scrollbar-hide"
        style="scrollbar-width: none;"
      >
        <div class="flex flex-nowrap space-x-4 pb-2 px-4" style="overflow: visible;">
          <ContentCard
            v-for="(content, key) in props.provider.suggestions"
            :key="content.id ?? key"
            :content="content"
          />
        </div>
      </div>
      <button
        v-if="canScrollRight"
        @click="scrollRight"
        class="absolute right-0 top-1/2 -translate-y-1/2 z-20 bg-black/60 hover:bg-black/80 text-white rounded-full w-10 h-10 flex items-center justify-center transition-opacity opacity-0 group-hover:opacity-100"
        aria-label="Rolar para a direita"
      >
        <svg width="24" height="24" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><path d="M9 5l7 7-7 7"/></svg>
      </button>
    </div>
  </section>
</template>

<style scoped>
.scrollbar-hide::-webkit-scrollbar {
  display: none;
}
.scrollbar-hide {
  -ms-overflow-style: none;
  scrollbar-width: none;
}
</style>
