<script setup lang="ts">
import ContentComponent from '@/components/Content.vue'
import Topbar from '@/components/Topbar.vue'
import Modal from '@/components/Modal.vue'
import { useLoginStore } from '@/stores/auth'
import { onMounted, reactive, ref, type Ref, watch } from 'vue'
import type Content from '@/types/content'
import CustomButton from '@/components/CustomButton.vue'
import router from '@/router'
import Pagination from '@/components/CustomPagination.vue'
import StatusFilter from '@/components/StatusFilter.vue'
import GenresFilter from '@/components/GenresFilter.vue'
import { Status } from '@/types/status'
import httpClient from '@/http/client'

interface Errors {
  [index: string]: string[]
}

const store = useLoginStore()

const error: Ref<string> = ref('')
const errors = ref<Errors>({})
const contents = ref<[]>([])

const pagination = ref({
  currentPage: 1,
  totalItems: 1,
  limit: 10,
})

const filters = ref<{
  watched: boolean | null
  search: string
}>({
  watched: null,
  search: '',
})

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

const genres = ref<string[]>([])

onMounted(fetchItems)

let searchTimeout: number | null = null

watch(
  () => filters.value.search,
  (newValue) => {
    if (searchTimeout) {
      clearTimeout(searchTimeout)
    }
    searchTimeout = setTimeout(() => {
      fetchItems()
    }, 300)
  },
)

async function fetchItems() {
  httpClient
    .get('/contents', {
      params: {
        page: pagination.value.currentPage,
        limit: pagination.value.limit,
        watched: filters.value.watched,
        search: filters.value.search,
        genres: genres.value.join(','),
      },
    })
    .catch((error) => {
      error.value = error.response.data.error
      errors.value = error.response.data.errors
    })
    .then(function (res) {
      if (res === undefined) {
        return
      }

      contents.value = res.data.data
      pagination.value.currentPage = res.data.page
      pagination.value.totalItems = res.data.total
      pagination.value.limit = res.data.limit
    })
}

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
      toggleModal('write-down-modal')
      fetchItems()
    })
}

async function details(content: Content) {
  httpClient
    .get('/contents/' + content.id)
    .catch((error) => {
      if (error.response.status === 401) {
        store.signOut()
        router.push({ name: 'login' })
      }

      error.value = error.response.data.error
      errors.value = error.response.data.errors
    })
    .then(async function (res) {
      if (res === undefined) {
        return
      }
      selectedContent.value = res.data.data
      toggleModal('details-modal')
    })
}

async function remove() {
  httpClient
    .delete('/contents/' + selectedContent.value?.id)
    .catch((error) => {
      error.value = error.response.data.error
      errors.value = error.response.data.errors
    })
    .then(function () {
      toggleModal('details-modal')
      fetchItems()
    })
}

async function rate() {
  httpClient
    .post('/contents/' + selectedContent.value?.id + '/rate', {
      rate: parseInt(rateInput.rate),
      comment: rateInput.comment,
    })
    .catch((error) => {
      error.value = error.response.data.error
      errors.value = error.response.data.errors
    })
    .then(function () {
      toggleModal('details-modal')
      fetchItems()
    })
}

function handleStatusFilter(status: Status) {
  if (status === Status.NONE) {
    filters.value.watched = null
    return
  }

  filters.value.watched = status == Status.WATCHED

  fetchItems()
}

function handleGenreFilter(newGenres: string[]) {
  genres.value = newGenres
  fetchItems()
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
      <CustomButton theme="default" @click="toggleModal('write-down-modal')"> Anotar </CustomButton>
    </div>
    <div class="flex gap-10 px-6 my-4 items-center">
      <input
        type="text"
        placeholder="Pesquisar por nome, categoria, descrição ou sinopse"
        v-model="filters.search"
        class="border text-sm rounded-lg block w-100 p-2.5 bg-slate-600 border-slate-500 placeholder-slate-400 text-white"
      />
      <GenresFilter @update:genres="handleGenreFilter" />
      <StatusFilter @filter-selected="handleStatusFilter" />
    </div>
    <div class="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-3 mx-auto px-6">
      <ContentComponent v-for="content in contents" :content="content" @click="details(content)" />
    </div>
    <div class="flex justify-center fixed bottom-10 w-full py-4">
      <Pagination
        :totalItems="pagination.totalItems"
        :currentPage="pagination.currentPage"
        :itemsPerPage="pagination.limit"
        @update:currentPage="
          (page: number) => {
            pagination.currentPage = page
            fetchItems()
          }
        "
      />
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
              class="bg-blue-600 text-white text-sm font-medium mr-2 px-2.5 py-0.5 rounded-md dark:bg-blue-700"
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
