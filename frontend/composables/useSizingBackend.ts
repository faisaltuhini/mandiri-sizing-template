// ─────────────────────────────────────────────────────────────────────────────
// Composable: useSizingBackend
// Connects frontend to the Go architecture sizing backend (default :8088)
// ─────────────────────────────────────────────────────────────────────────────

import { ref, onMounted } from 'vue'

export interface BackendStatus {
  connected: boolean
  service?: string
  version?: string
  engine?: string
  latencyMs?: number
  lastChecked?: string
}

export interface CalculateParams {
  daily_transactions: number
  peak_ratio: number
  peak_hours: number
  pod_tps_capacity: number
  cores_per_pod: number
  worker_cores: number
  avg_tx_kb: number
  days_retention: number
  annual_growth: number
}

export function useSizingBackend() {
  const config = useRuntimeConfig()
  const apiBase = (config.public?.apiBase as string) || 'http://localhost:8088'

  const backendStatus = ref<BackendStatus>({
    connected: false,
  })

  const isChecking = ref(false)

  async function checkHealth() {
    isChecking.value = true
    const start = performance.now()
    try {
      const res = await fetch(`${apiBase}/api/health`, { method: 'GET' })
      const latency = Math.round(performance.now() - start)
      if (res.ok) {
        const json = await res.json()
        backendStatus.value = {
          connected: true,
          service: json.service,
          version: json.version,
          engine: json.engine,
          latencyMs: latency,
          lastChecked: new Date().toLocaleTimeString(),
        }
      } else {
        backendStatus.value = { connected: false }
      }
    } catch (err) {
      backendStatus.value = { connected: false }
    } finally {
      isChecking.value = false
    }
  }

  async function fetchTechStackFromBackend() {
    try {
      const res = await fetch(`${apiBase}/api/sizing/tech-stack`)
      if (res.ok) {
        return await res.json()
      }
    } catch (e) {
      // Fallback
    }
    return null
  }

  async function calculateCapacity(params: CalculateParams) {
    try {
      const res = await fetch(`${apiBase}/api/sizing/calculate`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(params),
      })
      if (res.ok) {
        const json = await res.json()
        return json.result
      }
    } catch (e) {
      console.warn('Backend calculate offline, computing locally')
    }
    return null
  }

  async function fetchBOQ(tier: 'moderate' | 'critical' = 'moderate') {
    try {
      const res = await fetch(`${apiBase}/api/sizing/boq?tier=${tier}`)
      if (res.ok) {
        const json = await res.json()
        return json.boq
      }
    } catch (e) {
      console.warn('Backend BOQ offline')
    }
    return null
  }

  async function fetchPricingCatalog() {
    try {
      const res = await fetch(`${apiBase}/api/sizing/pricing`)
      if (res.ok) {
        const json = await res.json()
        return json.pricing
      }
    } catch (e) {
      console.warn('Backend pricing offline')
    }
    return null
  }

  onMounted(() => {
    checkHealth()
  })

  return {
    apiBase,
    backendStatus,
    isChecking,
    checkHealth,
    fetchTechStackFromBackend,
    calculateCapacity,
    fetchBOQ,
    fetchPricingCatalog,
  }
}
