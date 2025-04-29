<script setup lang="ts">
import httpClient from '@/http/client'
import { useModalStore } from '@/stores/modal'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import CustomButton from './CustomButton.vue'
import Modal from './Modal.vue'

const store = useModalStore()
const queryClient = useQueryClient()

const { mutate: deleteContent, isPending } = useMutation({
  mutationFn: async () => {
    return httpClient
      .delete(`/contents/${store.content?.id}`)
      .catch((error) => {
        error.value = error.response.error
      })
      .then(function (data) {
        store.toggleModal()
        return data
      })
  },
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['lastInserted'] })
  },
})

function handleCancelClick() {
  store.type = 'rating'
}
</script>
<template>
  <Modal
    id="delete-modal"
    :show="store.isVisible && store.type === 'delete'"
    :title="'Remover conteúdo'"
    @close="store.toggleModal()"
  >
    <div class="mb-4 p-3 rounded-md border border-zinc-700">
      <p class="text-md font-bold leading-tight">{{ store.content?.name }}</p>
      <p class="mt-1 text-sm text-zinc-400 truncate">{{ store.content?.summary }}</p>
    </div>
    <template #footer>
      <div class="flex justify-between">
        <CustomButton theme="default" @click="handleCancelClick">Voltar</CustomButton>
        <CustomButton theme="danger" @click="deleteContent" :disabled="isPending">
          Remover
        </CustomButton>
      </div>
    </template>
  </Modal>
</template>
