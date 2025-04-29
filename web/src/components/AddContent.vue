<script setup lang="ts">
import { reactive, ref, watchEffect } from 'vue'
import httpClient from '@/http/client'
import Modal from './Modal.vue'
import CustomButton from './CustomButton.vue'
import Cross from './icons/Cross.vue'
import { useModalStore } from '@/stores/modal'
import { useQueryClient } from '@tanstack/vue-query'

interface Input {
  name: string
  category: string
  summary: string
  wishLevel: string
  genres: string[]
}

interface Errors {
  [index: string]: string[]
}

const modalStore = useModalStore()
const errors = ref<Errors>({})
const writingGenre = ref<string>('')
const writeDownInput: Input = reactive({
  name: '',
  category: '',
  summary: '',
  wishLevel: '3',
  genres: [],
})

watchEffect(() => {
  writeDownInput.name = modalStore.content?.name ?? ''
  writeDownInput.category = modalStore.content?.category ?? ''
  writeDownInput.summary = modalStore.content?.summary ?? ''
  writeDownInput.wishLevel = '3'
  writeDownInput.genres = modalStore.content?.genres ?? []
})

function addGenre() {
  if (writeDownInput.genres.find((el: string) => el === writingGenre.value)) {
    writingGenre.value = ''
    return
  }
  writeDownInput.genres.push(writingGenre.value)
  writingGenre.value = ''
}

function removeGenre(genre: string) {
  writeDownInput.genres = writeDownInput.genres.filter((el: string) => el !== genre)
}

const queryClient = useQueryClient()

async function writeDown() {
  httpClient
    .post('/contents/write-down', {
      name: writeDownInput.name,
      category: writeDownInput.category,
      summary: writeDownInput.summary,
      wish_level: parseInt(writeDownInput.wishLevel),
      genres: writeDownInput.genres,
    })
    .catch((error) => {
      error.value = error.response.error
      errors.value = error.response.errors
    })
    .then(function () {
      modalStore.toggleModal()
      queryClient.invalidateQueries({ queryKey: ['lastInserted'] })
    })
}
</script>

<template>
  <CustomButton theme="default" @click="modalStore.toggleModal()"> Anotar </CustomButton>
  <Modal
    id="write-down-modal"
    :show="modalStore.isVisible && modalStore.type === 'write-down'"
    :title="'Anotar novo conteúdo'"
    @close="modalStore.toggleModal()"
  >
    <div class="flex flex-col gap-2">
      <div class="flex gap-2 flex-col">
        <label class="block text-sm">Nome</label>
        <input
          v-model="writeDownInput.name"
          type="text"
          class="w-full p-2 rounded-md bg-zinc-900 border border-zinc-700"
        />
      </div>

      <div class="flex gap-2 flex-col">
        <label class="block text-sm mt-2">Sinopse</label>
        <textarea
          v-model="writeDownInput.summary"
          rows="3"
          class="w-full p-2 rounded-md bg-zinc-900 border border-zinc-700"
        ></textarea>
      </div>

      <div class="flex gap-2 flex-col">
        <label class="block text-sm mt-2">Category</label>
        <input
          v-model="writeDownInput.category"
          type="text"
          class="w-full p-2 rounded-md bg-zinc-900 border border-zinc-700"
        />
      </div>

      <div class="flex gap-2 flex-col">
        <label class="block text-sm mt-2">Quanto estou interessado (1-5)</label>
        <input
          v-model="writeDownInput.wishLevel"
          type="range"
          min="1"
          max="5"
          step="1"
          class="w-full cursor-pointer"
        />
      </div>

      <div class="flex gap-2 flex-col">
        <label class="block text-sm mt-2">Gêneros</label>
        <input
          type="text"
          name="genres"
          id="genres"
          class="w-full p-2 rounded-md bg-zinc-900 border border-zinc-700"
          placeholder="Pressione enter"
          @keydown.enter.prevent="addGenre"
          v-model="writingGenre"
        />
        <div class="flex flex-wrap mt-2 gap-1">
          <span
            v-for="genre in writeDownInput.genres"
            :key="genre"
            class="bg-zinc-800 text-sm text-zinc-300 pl-3 py-1 rounded-full pr-1"
          >
            {{ genre }}
            <button
              type="button"
              class="text-white cursor-pointer hover:bg-white/10 rounded-full p-2"
              @click="removeGenre(genre)"
            >
              <Cross class="w-5 h-5" />
            </button>
          </span>
        </div>
      </div>
    </div>

    <template #footer>
      <div class="flex justify-end">
        <CustomButton @click="writeDown" theme="default"> Salvar </CustomButton>
      </div>
    </template>
  </Modal>
</template>
