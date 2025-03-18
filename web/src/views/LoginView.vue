<script setup lang="ts">
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
  <div class="min-h-screen flex items-center justify-center bg-slate-900 px-4 sm:px-0">
    <div class="bg-slate-800 p-8 rounded-lg shadow-lg w-full max-w-md">
      <h2 class="text-2xl font-bold text-white mb-6">Login</h2>
      <div v-if="error" class="mb-4 p-2 bg-red-400 text-white rounded">
        {{ error }}
      </div>
      <form @submit.prevent="login">
        <div class="mb-4">
          <label class="block mb-2 text-sm font-medium text-slate-900 dark:text-white" for="email"
            >Seu e-mail</label
          >
          <input
            class="border text-sm rounded-lg block w-full p-2.5 bg-slate-600 border-slate-500 placeholder-slate-400 text-white"
            type="email"
            id="email"
            v-model="email"
            placeholder="nome@wishtrack.com"
            required
          />
        </div>
        <div class="mb-6">
          <label
            class="block mb-2 text-sm font-medium text-slate-900 dark:text-white"
            for="password"
            >Sua senha</label
          >
          <input
            class="border text-sm rounded-lg block w-full p-2.5 bg-slate-600 border-slate-500 placeholder-slate-400 text-white"
            type="password"
            id="password"
            v-model="password"
            placeholder="••••••••"
            required
          />
        </div>
        <button
          class="cursor-pointer w-full bg-blue-500 text-white p-2 rounded hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500"
          type="submit"
        >
          Login
        </button>
      </form>
    </div>
  </div>
</template>
