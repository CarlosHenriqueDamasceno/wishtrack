<script setup lang="ts">
import Content from '@/components/Content.vue';
import Header from '@/components/Header.vue';
</script>

<template>
  <main>
    <Header />
    <div class="xl:px-0 xl:max-w-[1400px] w-full max-w-full mt-5">
      <div class="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-3 mx-auto px-4 sm:px-6 lg:px-8">
        <Content />
        <Content />
        <Content />
      </div>
    </div>
  </main>
</template>
