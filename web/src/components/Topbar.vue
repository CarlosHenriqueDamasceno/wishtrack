<script setup lang="ts">
import { useLoginStore } from '@/stores/auth'
import { useRouter } from 'vue-router'
import { ref } from 'vue'

const store = useLoginStore()
const router = useRouter()
const isOpen = ref<boolean>(false)

function logout() {
  store.signOut()
  router.push({ name: 'login' })
}
</script>

<template>
  <header class="bg-slate-900 text-white py-3 px-6 flex justify-between items-center">
    <div class="flex items-center">
      <img class="w-10 mr-1" src="../assets/logo.png" alt="" />
      <h1 class="text-lg font-bold hidden md:block">Wishtrack</h1>
    </div>
    <nav class="flex items-center gap-1">
      <div class="relative">
        <div
          @click="isOpen = !isOpen"
          class="flex items-center justify-center w-12 h-12 rounded-full bg-slate-300 text-slate-700 font-semibold text-lg cursor-pointer"
        >
          {{ store.auth?.username.charAt(0).toUpperCase() }}
        </div>

        <div
          v-if="isOpen"
          class="absolute right-0 mt-2 w-32 text-white border rounded-lg bg-slate-800 border-slate-700"
        >
          <button
            @click="logout"
            class="block cursor-pointer w-full px-4 py-2 text-left rounded-lg hover:bg-slate-500"
          >
            Sair
          </button>
        </div>
      </div>
    </nav>
  </header>
</template>
