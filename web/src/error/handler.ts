import { useLoginStore } from '@/stores/auth'
import { AxiosError } from 'axios'
import { useRouter } from 'vue-router'

const store = useLoginStore()
const router = useRouter()

export function useDefaultErrorHandlers(
  handler: (error: AxiosError) => void,
): (error: AxiosError) => void {
  return (error: Error) => {
    if (error instanceof AxiosError) {
      if (error.response?.status === 401) {
        store.signOut()
        router.push({ name: 'login' })
        return
      }
      handler(error)
    }
  }
}
