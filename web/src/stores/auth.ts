import { defineStore } from 'pinia'
import { effect, ref } from 'vue'

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

  console.log(isLogged())

  function isLogged(): boolean {
    return auth.value?.token !== null
  }

  function signIn(user: Auth) {
    auth.value = Object.assign(auth.value ?? {}, user)
    localStorage.setItem('auth', JSON.stringify(user))
  }

  function signOut() {
    auth.value = defaultValue
    localStorage.removeItem('auth')
  }

  function restore(): Auth {
    const storedValue = localStorage.getItem('auth') ?? '{}'
    return Object.assign({}, defaultValue, JSON.parse(storedValue))
  }

  return { auth, isLogged, signIn, signOut }
})
