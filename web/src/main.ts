import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'
import router from './router'
import axios from 'axios'

import App from './App.vue'

const httpClient = axios.create({
  baseURL: import.meta.env.VITE_API_URL,
  timeout: 1000,
  headers: {
    'Content-Type': 'application/json',
    Accept: 'application/json',
  },
})

const app = createApp(App)

const store = createPinia()
store.use(() => ({
  httpClient,
}))

app.use(router)
app.use(store)
app.provide('httpClient', httpClient)

app.mount('#app')
