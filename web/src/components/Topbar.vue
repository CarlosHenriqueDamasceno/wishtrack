<script setup lang="ts">
import { useLoginStore } from '@/stores/login';
import { useRouter } from 'vue-router';
import CustomButton from './CustomButton.vue';

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
      <h1 class="text-lg font-bold hidden md:block">
        Wishtrack
      </h1>
    </div>
    <nav class="">
      <CustomButton type="danger" @click="logout">Sair</CustomButton>
    </nav>
  </header>
</template>
