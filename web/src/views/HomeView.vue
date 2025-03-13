<script setup lang="ts">
import Button from '@/components/CustomButton.vue'
import Content from '@/components/Content.vue'
import Topbar from '@/components/Topbar.vue'
import Modal from '@/components/Modal.vue'
import { useLoginStore } from '@/stores/login'
import { onMounted, ref } from 'vue'

const store = useLoginStore()
const error = ref<boolean>(false)
const contents = ref<[]>([])
const showWriteDownModal = ref<boolean>(false)
const writingTag = ref<string>("")

const input = {
  name: ref<string>(),
  category: ref<string>(),
  summary: ref<string>(),
  wishLevel: ref<number>(),
  genres: ref<string[]>([])
}

onMounted(feed)

async function feed() {
  const res = await fetch(import.meta.env.VITE_API_URL + '/contents/feed', {
    headers: {
      'Authorization': 'Bearer '+store.token
    }
  })

  if (false === res.ok){
    error.value = true
  }

  contents.value = await res.json()
}

async function writeDown(){

}

function removeTag(tag: string){
  input.genres.value = input.genres.value.filter((el:string) => el !== tag)
}

function addTag(){
  if(input.genres.value.find((el:string) => el === writingTag.value)){
    writingTag.value = ""
    return
  }
  input.genres.value.push(writingTag.value)
  writingTag.value = ""
}

function toggleModal(){
  showWriteDownModal.value = !showWriteDownModal.value
}

</script>

<template>
  <main>
    <Topbar />
    <div class="flex justify-end px-6 my-4">
      <Button type="default" @click="toggleModal">Anotar</Button>
    </div>
    <div class="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-3 mx-auto px-6">
      <Content v-for="content in contents" :content="content"/>
    </div>
    <Modal :show="showWriteDownModal" :title="'Anotar novo conteúdo'" @close="toggleModal">
      <form class="space-y-4" @submit.prevent="writeDown">
        <div>
          <label for="name" class="block mb-2 text-sm font-medium text-gray-900 dark:text-white">Nome</label>
          <input type="text" name="name" id="name" class="border text-sm rounded-lg block w-full p-2.5 bg-gray-600 border-gray-500 placeholder-gray-400 text-white" placeholder="Lord of the rings" v-model="input.name" required />
        </div>
        <div>
          <label for="category" class="block mb-2 text-sm font-medium text-gray-900 dark:text-white">Categoria</label>
          <input type="text" name="category" id="category" class="border text-sm rounded-lg block w-full p-2.5 bg-gray-600 border-gray-500 placeholder-gray-400 text-white" placeholder="Filme" v-model="input.category" required />
        </div>
        <div>
          <label for="summary" class="block mb-2 text-sm font-medium text-gray-900 dark:text-white">Sinopse</label>
          <textarea name="summary" id="summary" class="border text-sm rounded-lg block w-full p-2.5 bg-gray-600 border-gray-500 placeholder-gray-400 text-white" placeholder="Description..." v-model="input.summary" required></textarea>
        </div>
        <div>
          <label for="wishLevel" class="block mb-2 text-sm font-medium text-gray-900 dark:text-white">O quanto tô afim</label>
          <input id="wishLevel" type="range" min="1" max="5" value="3" class="w-full h-2 rounded-lg cursor-pointer  bg-gray-500 accent-blue-600">
        </div>
        <div>
          <label for="tags" class="block mb-2 text-sm font-medium text-gray-900 dark:text-white">Gêneros</label>
          <input type="text" name="tags" id="tags" class="border text-sm rounded-lg block w-full p-2.5 bg-gray-600 border-gray-500 placeholder-gray-400 text-white" placeholder="Pressione enter" @keydown.enter.prevent="addTag" v-model="writingTag"/>
          <div class="flex flex-wrap mt-2">
            <span v-for="tag in input.genres.value" :key="tag" class="bg-blue-600 text-white text-sm font-medium mr-2 px-2.5 py-0.5 rounded dark:bg-blue-700">
              {{ tag }}
              <button type="button" class="ml-1 text-white" @click="removeTag(tag)">x</button>
            </span>
          </div>
        </div>
        <button type="submit" class="cursor-pointer w-full text-white focus:ring-4 focus:outline-none font-medium rounded-lg text-sm px-5 py-2.5 text-center dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800">
          Anotar
        </button>
      </form>
    </Modal>
  </main>
</template>
