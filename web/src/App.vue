<script setup lang="ts">
import { RouterView, useRouter } from 'vue-router'
import { useLoginStore } from './stores/auth'
import { setUnauthorizedErrorHandler } from './http/client'

const store = useLoginStore()
const router = useRouter()

setUnauthorizedErrorHandler(() => {
  store.signOut()
  router.push({ name: 'login' })
})
</script>

<template>
  <Suspense>
    <RouterView />
    <template #fallback>
      <div class="loading">Loading...</div>
    </template>
  </Suspense>
</template>
