<script setup lang="ts">
import CustomButton from '@/components/CustomButton.vue'
</script>

<template>
  <div class="flex flex-col items-center justify-center min-h-screen bg-slate-800 text-slate-200">
    <h1 class="text-7xl font-bold text-white">404</h1>
    <p class="text-xl mt-2 text-slate-300">Oops! Página não encontrada.</p>
    <router-link to="/" class="mt-4">
      <CustomButton :type="'default'">Voltar para home</CustomButton>
    </router-link>
  </div>
</template>
