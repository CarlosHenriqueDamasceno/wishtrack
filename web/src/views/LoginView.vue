<script setup lang="ts">
import { useLoginStore } from '@/stores/login';
import { ref, type Ref } from 'vue';
import { useRouter } from 'vue-router';

const email: Ref<string> = ref("")
const password: Ref<string> = ref("")
const error: Ref<string> = ref("")
const store = useLoginStore()
const router = useRouter()

async function login() {
  const res = await fetch('http://localhost:8080/api/v1/users/login', {
    method: 'POST',
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      email: email.value,
      password: password.value
    })
  })

  if (false === res.ok) {
    return res.json().then(function (json) { error.value = json.error })
  }

  const token = (await res.json()).token
  store.setToken(token)
  router.push({ name: 'home' })
}

</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-900 px-4 sm:px-0">
    <div class="bg-gray-800 p-8 rounded-lg shadow-lg w-full max-w-md">
      <h2 class="text-2xl font-bold text-white mb-6">Login</h2>
      <div v-if="error" class="mb-4 p-2 bg-red-400 text-white rounded">
        {{ error }}
      </div>
      <form @submit.prevent="login">
        <div class="mb-4">
          <label class="block text-gray-400 mb-2" for="email">Email</label>
          <input class="w-full p-2 rounded bg-gray-700 text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
            type="email" id="email" v-model="email" required>
        </div>
        <div class="mb-6">
          <label class="block text-gray-400 mb-2" for="password">Password</label>
          <input class="w-full p-2 rounded bg-gray-700 text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
            type="password" id="password" v-model="password" required>
        </div>
        <button
          class="cursor-pointer w-full bg-blue-500 text-white p-2 rounded hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500"
          type="submit">Login</button>
      </form>
    </div>
  </div>
</template>
