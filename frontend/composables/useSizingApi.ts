export function useBackendBaseUrl(): string {
  const config = useRuntimeConfig()
  if (import.meta.server && process.env.BACKEND_URL) {
    return process.env.BACKEND_URL.replace(/\/$/, '')
  }
  return ((config.public?.apiBase as string) ?? '').replace(/\/$/, '')
}

export function useSizingApi() {
  const backendBase = useBackendBaseUrl()

  async function checkHealth() {
    return await $fetch(`${backendBase}/api/health`)
  }

  async function fetchTechStackBaseline() {
    return await $fetch<{ total: number; data: any[] }>(`${backendBase}/api/v1/tech-stack/baseline`)
  }

  async function fetchKPISummaries() {
    return await $fetch<{ data: any[] }>(`${backendBase}/api/v1/tech-stack/kpis`)
  }

  async function calculatePowerLaw(payload: any) {
    return await $fetch(`${backendBase}/api/v1/sizing/calculate-power-law`, {
      method: 'POST',
      body: payload,
    })
  }

  async function calculateEAWG(payload: any) {
    return await $fetch(`${backendBase}/api/v1/sizing/calculate-eawg`, {
      method: 'POST',
      body: payload,
    })
  }

  async function uploadSolarwinds(file: File, metricName: string, unit: string) {
    const formData = new FormData()
    formData.append('file', file)
    formData.append('metric_name', metricName)
    formData.append('unit', unit)

    return await $fetch(`${backendBase}/api/v1/solarwinds/upload`, {
      method: 'POST',
      body: formData,
    })
  }

  async function fetchTemplates() {
    return await $fetch<{ total: number; data: any[] }>(`${backendBase}/api/v1/templates`)
  }

  function getTemplateDownloadUrl(filename: string) {
    return `${backendBase}/api/v1/templates/download/${encodeURIComponent(filename)}`
  }

  return {
    backendBase,
    checkHealth,
    fetchTechStackBaseline,
    fetchKPISummaries,
    calculatePowerLaw,
    calculateEAWG,
    uploadSolarwinds,
    fetchTemplates,
    getTemplateDownloadUrl,
  }
}
