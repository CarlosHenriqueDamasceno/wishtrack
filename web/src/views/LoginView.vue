<script setup lang="ts">
import Topbar from '@/components/Topbar.vue'
import { useLoginStore } from '@/stores/auth'
import { ref, type Ref } from 'vue'
import { useRouter } from 'vue-router'

const email: Ref<string> = ref('')
const password: Ref<string> = ref('')
const error: Ref<string> = ref('')
const store = useLoginStore()
const router = useRouter()

async function login() {
  const res = await fetch(import.meta.env.VITE_API_URL + '/users/login', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      email: email.value,
      password: password.value,
    }),
  })

  if (false === res.ok) {
    return res.json().then(function (json) {
      error.value = json.error
    })
  }

  const auth = await res.json()
  store.signIn(auth)
  router.push({ name: 'home' })
}
</script>

<template>
  <div class="min-h-screen bg-black text-white px-6 py-8 flex flex-col space-y-10">
    <div class="flex flex-1 items-center justify-center">
      <div class="bg-zinc-950 border border-zinc-800 rounded-xl shadow-sm w-full max-w-md p-8">
        <h2 class="text-2xl font-bold text-white mb-6">Login</h2>
        <div v-if="error" class="mb-4 p-2 bg-red-400 text-white rounded">
          {{ error }}
        </div>
        <form @submit.prevent="login" class="space-y-6">
          <div>
            <label class="block mb-2 text-sm font-medium text-white" for="email">Seu e-mail</label>
            <input
              class="w-full p-2 rounded-md bg-zinc-900 border border-zinc-700 text-white placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
              type="email"
              id="email"
              v-model="email"
              placeholder="nome@wishtrack.com"
              required
            />
          </div>
          <div>
            <label class="block mb-2 text-sm font-medium text-white" for="password">Sua senha</label>
            <input
              class="w-full p-2 rounded-md bg-zinc-900 border border-zinc-700 text-white placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
              type="password"
              id="password"
              v-model="password"
              placeholder="••••••••"
              required
            />
          </div>
          <button
            class="cursor-pointer w-full bg-blue-500 text-white p-2 rounded-md hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500"
            type="submit"
          >
            Login
          </button>
        </form>
      </div>
    </div>
  </div>
</template>
