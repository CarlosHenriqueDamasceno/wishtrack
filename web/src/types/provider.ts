import type Content from './content'

export default interface Provider {
  provider: string
  suggestions: Content[]
}
