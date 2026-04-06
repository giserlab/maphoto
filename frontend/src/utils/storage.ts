const prefix = 'maphoto_'

export const storage = {
  get<T>(key: string): T | null {
    try {
      const item = localStorage.getItem(prefix + key)
      return item ? JSON.parse(item) : null
    } catch {
      return null
    }
  },

  set(key: string, value: unknown): void {
    localStorage.setItem(prefix + key, JSON.stringify(value))
  },

  remove(key: string): void {
    localStorage.removeItem(prefix + key)
  },

  clear(): void {
    Object.keys(localStorage)
      .filter(key => key.startsWith(prefix))
      .forEach(key => localStorage.removeItem(key))
  },
}
