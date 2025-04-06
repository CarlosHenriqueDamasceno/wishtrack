import { defineStore } from 'pinia'
import { ref } from 'vue'
import httpClient from '../http/client'

export interface Auth {
  token: string | null
  username: string
}

const defaultValue = {
  token: null,
  username: '',
}

export const useLoginStore = defineStore('login', () => {
  const auth = ref<Auth>(restore())

  setHttpClientToken()

  function isLogged(): boolean {
    return auth.value?.token !== null
  }

  function signIn(user: Auth) {
    auth.value = Object.assign(auth.value ?? {}, user)
    localStorage.setItem('auth', JSON.stringify(user))
    setHttpClientToken()
  }

  function signOut() {
    auth.value = defaultValue
    localStorage.removeItem('auth')
    setHttpClientToken()
  }

  function restore(): Auth {
    const storedValue = localStorage.getItem('auth') ?? '{}'
    return Object.assign({}, defaultValue, JSON.parse(storedValue))
  }

  function setHttpClientToken() {
    httpClient.defaults.headers.common.Authorization = `Bearer ${auth.value.token}`
  }

  return { auth, isLogged, signIn, signOut }
})
