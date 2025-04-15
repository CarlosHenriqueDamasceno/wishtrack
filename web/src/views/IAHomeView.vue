<script setup lang="ts">
import AddContent from '@/components/AddContent.vue'
import ContentRow from '@/components/ContentRow.vue'
import RateContent from '@/components/RateContent.vue'
import Topbar from '@/components/Topbar.vue'
import httpClient from '@/http/client'
import { useQuery } from '@tanstack/vue-query'
import { ref } from 'vue'

const { data: suggestions } = useQuery({
  queryKey: ['suggestions'],
  queryFn: async () => (await httpClient.get('/suggestions')).data,
})

const { data: lastInserted } = useQuery({
  queryKey: ['lastInserted'],
  queryFn: async () =>
    (
      await httpClient.get('/contents', {
        params: {
          limit: 5,
        },
      })
    ).data.data,
})

const isRatingModalVisible = ref(false)
</script>

<template>
  <div class="min-h-screen bg-black text-white px-6 py-8 space-y-10">
    <Topbar />
    <div class="flex justify-end px-6 my-4">
      <AddContent />
      <RateContent />
    </div>

    <ContentRow
      :provider="{
        provider: 'Últimos adicionados',
        suggestions: lastInserted,
      }"
    />

    <ContentRow
      v-for="(provider, key) in suggestions"
      :key="provider.name ?? key"
      :provider="provider"
    />
  </div>
</template>
