// ─────────────────────────────────────────────
// Composable: IndexedDB Persistence for Sizing Suite
// Stores layer adjustments, tier, license filter, environment, and DRC topology
// ─────────────────────────────────────────────

const DB_NAME = 'EnterpriseSizingSuiteDB'
const DB_VERSION = 1
const STORE_NAME = 'sizing_configs'

function openDB(): Promise<IDBDatabase> {
  return new Promise((resolve, reject) => {
    if (typeof window === 'undefined' || !window.indexedDB) {
      reject(new Error('IndexedDB not available'))
      return
    }

    const request = window.indexedDB.open(DB_NAME, DB_VERSION)

    request.onupgradeneeded = (event) => {
      const db = (event.target as IDBOpenDBRequest).result
      if (!db.objectStoreNames.contains(STORE_NAME)) {
        db.createObjectStore(STORE_NAME, { keyPath: 'key' })
      }
    }

    request.onsuccess = () => resolve(request.result)
    request.onerror = () => reject(request.error)
  })
}

export function useSizingIndexedDB() {
  async function setItem<T>(key: string, value: T): Promise<void> {
    if (typeof window === 'undefined') return
    try {
      const db = await openDB()
      return new Promise((resolve, reject) => {
        const tx = db.transaction(STORE_NAME, 'readwrite')
        const store = tx.objectStore(STORE_NAME)
        const req = store.put({ key, value: JSON.parse(JSON.stringify(value)), updatedAt: Date.now() })

        req.onsuccess = () => resolve()
        req.onerror = () => reject(req.error)
      })
    } catch (err) {
      console.warn(`[IndexedDB] Failed to save key "${key}":`, err)
    }
  }

  async function getItem<T>(key: string): Promise<T | null> {
    if (typeof window === 'undefined') return null
    try {
      const db = await openDB()
      return new Promise((resolve, reject) => {
        const tx = db.transaction(STORE_NAME, 'readonly')
        const store = tx.objectStore(STORE_NAME)
        const req = store.get(key)

        req.onsuccess = () => {
          if (req.result && req.result.value !== undefined) {
            resolve(req.result.value as T)
          } else {
            resolve(null)
          }
        }
        req.onerror = () => reject(req.error)
      })
    } catch (err) {
      console.warn(`[IndexedDB] Failed to load key "${key}":`, err)
      return null
    }
  }

  async function removeItem(key: string): Promise<void> {
    if (typeof window === 'undefined') return
    try {
      const db = await openDB()
      return new Promise((resolve, reject) => {
        const tx = db.transaction(STORE_NAME, 'readwrite')
        const store = tx.objectStore(STORE_NAME)
        const req = store.delete(key)

        req.onsuccess = () => resolve()
        req.onerror = () => reject(req.error)
      })
    } catch (err) {
      console.warn(`[IndexedDB] Failed to delete key "${key}":`, err)
    }
  }

  return {
    setItem,
    getItem,
    removeItem,
  }
}
