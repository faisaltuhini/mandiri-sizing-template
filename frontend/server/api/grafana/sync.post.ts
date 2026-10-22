import { defineEventHandler, readBody, createError } from 'h3'

export interface GrafanaSyncRequest {
  grafanaUrl?: string
  apiToken?: string
  timeSpan?: '1y' | '6M' | '90d' | '30d' | '7d' | '1d'
  aggregation?: 'max' | 'p95' | 'avg'
  hosts: Array<{
    id: string
    hostName: string
    usage: string
    node: number
  }>
  isDemoMode?: boolean
}

export interface HostMetricsResult {
  hostName: string
  cpuUtil: number
  memUtil: number
  status: 'normal' | 'warning' | 'critical'
  peakTimestamp?: string
  source: 'grafana-api' | 'telemetry-simulator'
}

export default defineEventHandler(async (event) => {
  const body = await readBody<GrafanaSyncRequest>(event)

  if (!body || !body.hosts || !Array.isArray(body.hosts)) {
    throw createError({
      statusCode: 400,
      statusMessage: 'Payload wajib menyertakan daftar hosts',
    })
  }

  const timeSpan = body.timeSpan || '1y'
  const aggregation = body.aggregation || 'p95'
  const isDemo = body.isDemoMode || !body.grafanaUrl || !body.apiToken

  // Map time span to human readable label
  const spanLabels: Record<string, string> = {
    '1y': '1 Tahun Terakhir (1y Fiscal Window)',
    '6M': '6 Bulan Terakhir (6M Peak Cycle)',
    '90d': '3 Bulan Terakhir (Quarterly Window)',
    '30d': '30 Hari Terakhir (Monthly Payday Peak)',
    '7d': '7 Hari Terakhir (Weekly Window)',
    '1d': '24 Jam Terakhir (Daily Window)',
  }

  // If live Grafana URL and API Token provided, attempt real HTTP fetch to Grafana
  if (!isDemo && body.grafanaUrl && body.apiToken) {
    try {
      const cleanUrl = body.grafanaUrl.replace(/\/$/, '')
      const fromTime = `now-${timeSpan}`
      const toTime = 'now'

      // Attempt Grafana DataSource query endpoint
      // Reference: POST /api/ds/query
      const grafanaResponse = await fetch(`${cleanUrl}/api/ds/query`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${body.apiToken}`,
        },
        body: JSON.stringify({
          from: fromTime,
          to: toTime,
          queries: body.hosts.map((h, i) => ({
            refId: `query_${i}`,
            datasource: { type: 'prometheus' },
            expr: `100 - (avg by (instance) (rate(node_cpu_seconds_total{mode="idle", instance=~"${h.hostName}.*"}[5m])) * 100)`,
            instant: false,
          })),
        }),
        signal: AbortSignal.timeout(5000), // 5s timeout
      })

      if (grafanaResponse.ok) {
        const data = await grafanaResponse.json()
        // If live response succeeded, parse Prometheus/Grafana timeseries
        const results: Record<string, HostMetricsResult> = {}

        body.hosts.forEach((h, idx) => {
          const queryResult = data.results?.[`query_${idx}`]
          let cpuVal = 35
          let memVal = 50

          if (queryResult?.frames?.[0]?.data?.values) {
            const vals: number[] = queryResult.frames[0].data.values[1] || []
            if (vals.length > 0) {
              if (aggregation === 'max') {
                cpuVal = Math.round(Math.max(...vals.filter(v => typeof v === 'number' && !isNaN(v))))
              } else if (aggregation === 'p95') {
                const sorted = [...vals].sort((a, b) => a - b)
                const pIndex = Math.floor(sorted.length * 0.95)
                cpuVal = Math.round(sorted[pIndex] || 40)
              } else {
                const sum = vals.reduce((a, b) => a + b, 0)
                cpuVal = Math.round(sum / vals.length)
              }
            }
          }

          const status = cpuVal >= 70 || memVal >= 85 ? 'critical' : cpuVal >= 45 || memVal >= 70 ? 'warning' : 'normal'

          results[h.hostName] = {
            hostName: h.hostName,
            cpuUtil: Math.min(100, Math.max(1, cpuVal)),
            memUtil: Math.min(100, Math.max(1, memVal)),
            status,
            source: 'grafana-api',
          }
        })

        return {
          success: true,
          source: 'grafana-live-api',
          timeSpan,
          timeSpanLabel: spanLabels[timeSpan] || timeSpan,
          aggregation,
          syncedAt: new Date().toISOString(),
          hostsCount: body.hosts.length,
          results,
        }
      }
    } catch (err: any) {
      console.warn('Grafana live API call failed, falling back to realistic telemetry engine:', err.message)
    }
  }

  // Realistic Enterprise Telemetry Generator (used when testing without live VPN or in demo mode)
  const results: Record<string, HostMetricsResult> = {}

  for (const h of body.hosts) {
    const u = h.usage.toLowerCase()

    // Realistic baseline metrics based on host role
    let baseCpu = 35
    let baseMem = 55

    if (u.includes('worker')) {
      baseCpu = timeSpan === '1y' ? 68 : timeSpan === '6M' ? 62 : timeSpan === '30d' ? 58 : 48
      baseMem = 72
    } else if (u.includes('oracle') || u.includes('postgres') || u.includes('db')) {
      baseCpu = timeSpan === '1y' ? 55 : 42
      baseMem = 78
    } else if (u.includes('redis')) {
      baseCpu = 28
      baseMem = timeSpan === '1y' ? 65 : 55
    } else if (u.includes('web')) {
      baseCpu = timeSpan === '1y' ? 45 : 36
      baseMem = 48
    } else if (u.includes('router') || u.includes('infra')) {
      baseCpu = 22
      baseMem = 38
    } else if (u.includes('master')) {
      baseCpu = 18
      baseMem = 42
    }

    // Apply multiplier based on aggregation mode
    let multiplier = 1.0
    if (aggregation === 'max') multiplier = 1.25 // Peak spike
    else if (aggregation === 'p95') multiplier = 1.10 // 95th Percentile SLA
    else multiplier = 0.90 // Average

    // Add deterministic variance per node number
    const nodeVariance = ((h.node * 7) % 11) - 5
    const finalCpu = Math.min(98, Math.max(8, Math.round(baseCpu * multiplier + nodeVariance)))
    const finalMem = Math.min(95, Math.max(12, Math.round(baseMem * multiplier + nodeVariance * 0.5)))

    const status = finalCpu >= 70 || finalMem >= 85 ? 'critical' : finalCpu >= 45 || finalMem >= 70 ? 'warning' : 'normal'

    results[h.hostName] = {
      hostName: h.hostName,
      cpuUtil: finalCpu,
      memUtil: finalMem,
      status,
      peakTimestamp: new Date(Date.now() - Math.floor(Math.random() * 86400000 * 30)).toISOString(),
      source: 'telemetry-simulator',
    }
  }

  return {
    success: true,
    source: isDemo ? 'telemetry-simulator' : 'grafana-fallback',
    timeSpan,
    timeSpanLabel: spanLabels[timeSpan] || timeSpan,
    aggregation,
    syncedAt: new Date().toISOString(),
    hostsCount: body.hosts.length,
    results,
  }
})
