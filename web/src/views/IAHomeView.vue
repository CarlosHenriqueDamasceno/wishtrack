<script setup lang="ts">
import AddContent from '@/components/AddContent.vue'
import ContentRow from '@/components/ContentRow.vue'
import DeleteContent from '@/components/DeleteContent.vue'
import RateContent from '@/components/RateContent.vue'
import Topbar from '@/components/Topbar.vue'
import httpClient from '@/http/client'
import { categoriesMap } from '@/types/categories'
import { useQuery } from '@tanstack/vue-query'

const { data: suggestions } = useQuery({
  queryKey: ['suggestions'],
  queryFn: async () => {
    const data = (await httpClient.get('/suggestions')).data

    return Object.values(data).map((el: any) => ({
      ...el,
      suggestions: el.suggestions.map((el: any) => {
        const category = categoriesMap[el.category as keyof typeof categoriesMap]
        if (!category) {
          throw new Error(`Category ${el.category} not found`)
        }
        return {
          ...el,
          category: category,
        }
      }),
    }))
  },
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
</script>

<template>
  <div class="min-h-screen bg-black text-white px-6 py-8 space-y-10">
    <Topbar />
    <div class="flex justify-end px-6 my-4">
      <AddContent />
      <RateContent />
      <DeleteContent />
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
