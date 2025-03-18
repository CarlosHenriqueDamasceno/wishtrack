<script setup lang="ts">
import Button from '@/components/CustomButton.vue'
import ContentComponent from '@/components/Content.vue'
import Topbar from '@/components/Topbar.vue'
import Modal from '@/components/Modal.vue'
import { useLoginStore } from '@/stores/auth'
import { onMounted, reactive, ref, type Ref } from 'vue'
import type Content from '@/types/content'
import CustomButton from '@/components/CustomButton.vue'
import router from '@/router'

interface Errors {
  [index: string]: string[]
}

const store = useLoginStore()
const error: Ref<string> = ref('')
const errors = ref<Errors>({})
const contents = ref<[]>([])
const showWriteDownModal = ref<boolean>(false)
const showRateForm = ref<boolean>(false)
const writingTag = ref<string>('')
const selectedContent = ref<Content | null>(null)
const showDetailsModal = ref<boolean>(false)

interface Input {
  name: string
  category: string
  summary: string
  wishLevel: string
  genres: string[]
}

interface RateInput {
  rate: string
  comment: string
}

const writeDownInput: Input = reactive({
  name: '',
  category: '',
  summary: '',
  wishLevel: '4',
  genres: [],
})

const rateInput: RateInput = reactive({
  rate: '3',
  comment: '',
})

onMounted(feed)

async function feed() {
  const res = await fetch(import.meta.env.VITE_API_URL + '/contents/feed', {
    headers: {
      Authorization: 'Bearer ' + store.auth?.token,
    },
  })

  if (false === res.ok) {
    if (res.status === 401) {
      store.signOut()
      router.push({ name: 'login' })
    }

    return res.json().then(function (json) {
      error.value = json.error
      errors.value = json.errors
    })
  }

  contents.value = await res.json()
}

async function writeDown() {
  const res = await fetch(import.meta.env.VITE_API_URL + '/contents/write-down', {
    headers: {
      Authorization: 'Bearer ' + store.auth?.token,
      'Content-Type': 'application/json',
    },
    method: 'POST',
    body: JSON.stringify({
      name: writeDownInput.name,
      category: writeDownInput.category,
      summary: writeDownInput.summary,
      wish_level: parseInt(writeDownInput.wishLevel),
      genres: writeDownInput.genres,
    }),
  })

  if (false === res.ok) {
    return res.json().then(function (json) {
      error.value = json.error
      errors.value = json.errors
    })
  }

  toggleModal('write-down-modal')
  feed()
}

async function details(content: Content) {
  const res = await fetch(import.meta.env.VITE_API_URL + '/contents/' + content.id, {
    headers: {
      Authorization: 'Bearer ' + store.auth?.token,
    },
  })

  if (false === res.ok) {
    return res.json().then(function (json) {
      error.value = json.error
      errors.value = json.errors
    })
  }

  selectedContent.value = await res.json()
  toggleModal('details-modal')
}

async function remove() {
  const res = await fetch(import.meta.env.VITE_API_URL + '/contents/' + selectedContent.value?.id, {
    headers: {
      Authorization: 'Bearer ' + store.auth?.token,
    },
    method: 'DELETE',
  })

  if (false === res.ok) {
    return res.json().then(function (json) {
      error.value = json.error
      errors.value = json.errors
    })
  }

  toggleModal('details-modal')
  feed()
}

async function rate() {
  const res = await fetch(
    import.meta.env.VITE_API_URL + '/contents/' + selectedContent.value?.id + '/rate',
    {
      headers: {
        Authorization: 'Bearer ' + store.auth?.token,
        'Content-Type': 'application/json',
      },
      method: 'POST',
      body: JSON.stringify({
        rate: parseInt(rateInput.rate),
        comment: rateInput.comment,
      }),
    },
  )

  if (false === res.ok) {
    return res.json().then(function (json) {
      error.value = json.error
      errors.value = json.errors
    })
  }

  toggleModal('details-modal')
  feed()
}

function removeTag(tag: string) {
  writeDownInput.genres = writeDownInput.genres.filter((el: string) => el !== tag)
}

function addTag() {
  if (writeDownInput.genres.find((el: string) => el === writingTag.value)) {
    writingTag.value = ''
    return
  }
  writeDownInput.genres.push(writingTag.value)
  writingTag.value = ''
}

function toggleModal(target: string) {
  switch (target) {
    case 'write-down-modal':
      showWriteDownModal.value = !showWriteDownModal.value
      break
    case 'details-modal':
      showDetailsModal.value = !showDetailsModal.value
      break
  }
}
</script>

<template>
  <main>
    <Topbar />
    <div class="flex justify-end px-6 my-4">
      <Button type="default" @click="toggleModal('write-down-modal')"> Anotar </Button>
    </div>
    <div class="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-3 mx-auto px-6">
      <ContentComponent v-for="content in contents" :content="content" @click="details(content)" />
    </div>
    <Modal
      id="write-down-modal"
      :show="showWriteDownModal"
      :title="'Anotar novo conteúdo'"
      @close="toggleModal('write-down-modal')"
    >
      <form class="space-y-4">
        <div v-if="error" class="mb-4 p-2 bg-red-400 text-white rounded">
          {{ error }}
        </div>
        <div v-for="error in errors" class="mb-4 p-2 bg-red-400 text-white rounded">
          {{ error.join('<br />') }}
        </div>
        <div>
          <label for="name" class="block mb-2 text-sm font-medium text-slate-900 dark:text-white"
            >Nome</label
          >
          <input
            type="text"
            name="name"
            id="name"
            class="border text-sm rounded-lg block w-full p-2.5 bg-slate-600 border-slate-500 placeholder-slate-400 text-white"
            placeholder="Lord of the rings"
            v-model="writeDownInput.name"
            required
          />
        </div>
        <div>
          <label
            for="category"
            class="block mb-2 text-sm font-medium text-slate-900 dark:text-white"
            >Categoria</label
          >
          <input
            type="text"
            name="category"
            id="category"
            class="border text-sm rounded-lg block w-full p-2.5 bg-slate-600 border-slate-500 placeholder-slate-400 text-white"
            placeholder="Filme"
            v-model="writeDownInput.category"
            required
          />
        </div>
        <div>
          <label for="summary" class="block mb-2 text-sm font-medium text-slate-900 dark:text-white"
            >Sinopse</label
          >
          <textarea
            name="summary"
            id="summary"
            class="border text-sm rounded-lg block w-full p-2.5 bg-slate-600 border-slate-500 placeholder-slate-400 text-white"
            placeholder="Description..."
            v-model="writeDownInput.summary"
          ></textarea>
        </div>
        <div>
          <label
            for="wishLevel"
            class="block mb-2 text-sm font-medium text-slate-900 dark:text-white"
            >O quanto tô afim</label
          >
          <input
            id="wishLevel"
            type="range"
            min="1"
            max="5"
            value="3"
            v-model="writeDownInput.wishLevel"
            class="w-full h-2 rounded-lg cursor-pointer bg-slate-500 accent-blue-600"
          />
        </div>
        <div>
          <label for="tags" class="block mb-2 text-sm font-medium text-slate-900 dark:text-white"
            >Gêneros</label
          >
          <input
            type="text"
            name="tags"
            id="tags"
            class="border text-sm rounded-lg block w-full p-2.5 bg-slate-600 border-slate-500 placeholder-slate-400 text-white"
            placeholder="Pressione enter"
            @keydown.enter.prevent="addTag"
            v-model="writingTag"
          />
          <div class="flex flex-wrap mt-2">
            <span
              v-for="tag in writeDownInput.genres"
              :key="tag"
              class="bg-blue-600 text-white text-sm font-medium mr-2 px-2.5 py-0.5 rounded dark:bg-blue-700"
            >
              {{ tag }}
              <button type="button" class="ml-1 text-white" @click="removeTag(tag)">x</button>
            </span>
          </div>
        </div>
      </form>
      <template #footer>
        <div class="flex justify-end">
          <CustomButton theme="default" @click="writeDown"> Anotar </CustomButton>
        </div>
      </template>
    </Modal>
    <Modal
      id="details-modal"
      :show="showDetailsModal"
      title="Detalhes"
      @close="
        () => {
          toggleModal('details-modal')
          showRateForm = false
        }
      "
    >
      <div v-if="showRateForm">
        <div v-if="error" class="mb-4 p-2 bg-red-400 text-white rounded">
          {{ error }}
        </div>
        <div v-for="error in errors" class="mb-4 p-2 bg-red-400 text-white rounded">
          {{ error.join('<br />') }}
        </div>
        <form>
          <div>
            <label for="rate" class="block mb-2 text-sm font-medium text-slate-900 dark:text-white"
              >O quanto curti</label
            >
            <input
              id="rate"
              type="range"
              min="1"
              max="5"
              value="3"
              class="w-full h-2 rounded-lg cursor-pointer bg-slate-500 accent-blue-600"
              v-model="rateInput.rate"
            />
          </div>
          <div class="mt-2">
            <label
              for="comment"
              class="block mb-2 text-sm font-medium text-slate-900 dark:text-white"
              >Comentário</label
            >
            <textarea
              name="comment"
              id="comment"
              class="border text-sm rounded-lg block w-full p-2.5 bg-slate-600 border-slate-500 placeholder-slate-400 text-white"
              placeholder="Obra prima..."
              v-model="rateInput.comment"
            ></textarea>
          </div>
        </form>
      </div>

      <div v-else class="flex justify-center">
        <ContentComponent class="w-full" v-if="selectedContent" :content="selectedContent" />
      </div>

      <template #footer>
        <div v-if="false === showRateForm" class="flex justify-between gap-3 items-center">
          <a href="#" class="text-red-400" @click="remove">Excluir</a>
          <CustomButton theme="default" @click="showRateForm = true">Avaliar</CustomButton>
        </div>
        <div v-else class="flex justify-between gap-3 items-center">
          <a href="#" class="text-white" @click="showRateForm = false">Voltar</a>
          <CustomButton theme="default" @click="rate">Salvar</CustomButton>
        </div>
      </template>
    </Modal>
  </main>
</template>
