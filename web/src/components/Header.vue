<script setup lang="ts">
import { useLoginStore } from '@/stores/login';
import { useRouter } from 'vue-router';

const store = useLoginStore()
const router = useRouter()

function logout() {
  store.clearToken()
  router.push({ name: 'login' })
}

</script>
<template>
  <header class="bg-gray-900 text-white py-3 px-6 flex justify-between items-center">
    <div class="flex items-center">
      <img class="w-10 mr-1" src="../assets/logo.png" alt="">
      <h1 class="text-lg font-bold">
        Wishtrack
      </h1>
    </div>
    <nav class="flex gap-4">
      <button type="button"
        class="border cursor-pointer focus:outline-none hover:bg-gray-100 focus:ring-4 focus:ring-gray-100 font-medium rounded-lg text-sm px-5 py-2.5 me-2 mb-2 dark:bg-gray-800 dark:text-white dark:border-gray-600 dark:hover:bg-gray-700 dark:hover:border-gray-600 dark:focus:ring-gray-700">
        Anotar
      </button>
      <button type="button" @click="logout"
        class="border cursor-pointer font-medium rounded-lg focus:ring-4 text-sm px-5 py-2.5 me-2 mb-2 bg-rose-600 text-white border-rose-300 hover:bg-rose-700 hover:border-rose-600 focus:ring-rose-700">
        Sair
      </button>
    </nav>
  </header>
</template>
