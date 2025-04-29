<script setup lang="ts">
import { reactive, watchEffect } from 'vue'
import Modal from './Modal.vue'
import CustomButton from './CustomButton.vue'
import httpClient from '@/http/client'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { useModalStore } from '@/stores/modal'

interface RateInput {
  rate: string
  comment: string
}

const store = useModalStore()
const rateInput: RateInput = reactive({
  rate: '3',
  comment: '',
})

watchEffect(() => {
  rateInput.rate = store.content?.rate?.toString() ?? '3'
  rateInput.comment = store.content?.comment ?? ''
})

const queryClient = useQueryClient()

const { mutate: saveRating, isPending } = useMutation({
  mutationFn: async () => {
    return httpClient
      .post(`/contents/${store.content?.id}/rate`, {
        rate: parseInt(rateInput.rate),
        comment: rateInput.comment,
      })
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

function handleDeleteClick() {
  store.type = 'delete'
}
</script>

<template>
  <Modal
    id="rate-modal"
    :show="store.isVisible && store.type === 'rating'"
    :title="'Avaliar conteúdo'"
    @close="store.toggleModal()"
  >
    <div class="mb-4 p-3 rounded-md border border-zinc-700">
      <p class="text-md font-bold leading-tight">{{ store.content?.name }}</p>
      <p class="mt-1 text-sm text-zinc-400 truncate">{{ store.content?.summary }}</p>
    </div>
    <div class="flex gap-2 flex-col">
      <label class="block text-sm">Nota (1-5)</label>
      <input v-model="rateInput.rate" type="range" min="1" max="5" step="1" class="w-full" />
    </div>
    <div class="flex gap-2 flex-col">
      <label class="block text-sm mt-2">Comentário</label>
      <textarea
        v-model="rateInput.comment"
        rows="3"
        class="w-full p-2 rounded-md bg-zinc-900 border border-zinc-700"
      ></textarea>
    </div>
    <template #footer>
      <div class="flex justify-between">
        <CustomButton theme="danger" @click="handleDeleteClick">Remover</CustomButton>
        <CustomButton theme="default" @click="saveRating" :disabled="isPending"
          >Avaliar</CustomButton
        >
      </div>
    </template>
  </Modal>
</template>
