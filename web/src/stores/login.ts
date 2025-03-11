import { defineStore } from 'pinia';

export const useLoginStore = defineStore('login', {
  state: () => ({
    token: localStorage.getItem('token')
  }),
  actions: {
    isLogged(): boolean {
      return this.token !== null
    },
    setToken(newToken: string) {
      this.token = newToken;
      localStorage.setItem('token', newToken);
    },
    clearToken() {
      this.token = '';
      localStorage.removeItem('token');
    }
  }
});
