<template>
  <div class="capacity-calculator">
    <!-- Header -->
    <div class="calc-header">
      <div class="calc-header-main">
        <div class="header-titles">
          <div class="calc-badge">Kalkulator Kapasitas & Ekstrapolasi Sizing</div>
          <h2 class="calc-title">Enterprise Capacity Planning & Sizing Engine</h2>
          <p class="calc-subtitle">
            Kalkulasi otomatis proyeksi Peak TPS, model scaling Power-Law (a · TPSᵇ), sizing OpenShift Worker Node, estimasi DB 5 Tahun, dan format tabel resmi dokumen EAWG.
          </p>
        </div>

        <!-- Mode Switcher -->
        <div class="view-mode-switch">
          <button
            type="button"
            :class="['mode-tab-btn', { active: activeViewMode === 'engine' }]"
            @click="activeViewMode = 'engine'"
          >
            Simulasi Sizing Terintegrasi
          </button>
          <button
            type="button"
            :class="['mode-tab-btn', { active: activeViewMode === 'eawg' }]"
            @click="activeViewMode = 'eawg'"
          >
            Format Resmi Dokumen EAWG
          </button>
        </div>

        <!-- Presets Selector (Engine mode) -->
        <div v-if="activeViewMode === 'engine'" class="preset-selector">
          <span class="preset-label">Pilihan Preset Layanan:</span>
          <div class="preset-buttons">
            <button
              v-for="p in presets"
              :key="p.id"
              type="button"
              :class="['preset-btn', { active: activePreset === p.id }]"
              @click="applyPreset(p.id)"
            >
              {{ p.icon }} {{ p.name }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- VIEW MODE 1: Power-Law Engine -->
    <div v-if="activeViewMode === 'engine'" class="engine-view-container">
      <!-- Main Grid -->
      <div class="calc-grid">
        <!-- Section 1: Business Growth & Peak TPS -->
        <div class="calc-card">
        <div class="card-header">
          <span class="card-icon">📈</span>
          <div>
            <h3>1. Proyeksi Transaksi & Peak TPS</h3>
            <p>Estimasi beban puncak di jam sibuk operasional</p>
          </div>
        </div>

        <div class="form-grid">
          <div class="form-group">
            <label>Daily Transactions (Trx<sub>daily</sub>)</label>
            <UiNumberInput
              v-model="params.trxDaily"
              :step="100000"
              :min="1000"
              unit="transaksi / hari"
            />
          </div>

          <div class="form-group">
            <label>Peak Hours Window</label>
            <UiNumberInput
              v-model="params.peakHours"
              :step="1"
              :min="1"
              :max="24"
              unit="jam operasional sibuk"
            />
          </div>

          <div class="form-group">
            <label>Concurrency Spike Factor</label>
            <UiNumberInput
              v-model="params.concurrencyFactor"
              :step="0.1"
              :min="1.0"
              :max="5.0"
              :max-decimals="2"
              unit="rasio lonjakan peak (1.5× – 2.0×)"
            />
          </div>

          <div class="form-group">
            <label>Tingkat Pertumbuhan Tahunan (g)</label>
            <UiNumberInput
              v-model="params.growthRatePercent"
              :step="1"
              :min="0"
              :max="100"
              unit="% CAGR per tahun"
            />
          </div>
        </div>

        <div class="result-box result-highlight">
          <div class="result-label">Proyeksi Peak TPS Puncak</div>
          <div class="result-value">{{ formatNumber(peakTps, { maxDecimals: 1 }) }} <span class="unit">TPS</span></div>
          <div class="result-formula">Formula: (Trx<sub>daily</sub> × Factor<sub>concurrency</sub>) / (PeakHours × 3600)</div>
        </div>
      </div>

      <!-- Section 2: Power-Law Extrapolation (PT Baseline) -->
      <div class="calc-card span-2">
        <div class="card-header">
          <span class="card-icon">🧪</span>
          <div>
            <h3>2. Baseline Uji Beban (PT) & Model Power-Law (a · TPSᵇ)</h3>
            <p>Input 2 titik data hasil stress test untuk menurunkan eksponen scaling non-linear</p>
          </div>
        </div>

        <div class="pt-points-grid">
          <!-- PT Point 1 -->
          <div class="pt-box">
            <h4>Titik Uji Beban 1 (PT Low Load)</h4>
            <div class="form-group">
              <label>TPS Titik 1 (TPS₁)</label>
              <UiNumberInput
                v-model="params.pt1Tps"
                :step="100"
                :min="1"
                unit="TPS"
              />
            </div>
            <div class="form-group">
              <label>Node Spec (Cores/Node)</label>
              <UiNumberInput
                v-model="params.pt1NodeCores"
                :step="2"
                :min="1"
                unit="Cores"
              />
            </div>
            <div class="form-group">
              <label>Jumlah Node PT 1</label>
              <UiNumberInput
                v-model="params.pt1NodeQty"
                :step="1"
                :min="1"
                unit="Node"
              />
            </div>
            <div class="form-group">
              <label>Rata-rata % CPU PT 1</label>
              <UiNumberInput
                v-model="params.pt1CpuPercent"
                :step="5"
                :min="1"
                :max="100"
                unit="%"
              />
            </div>
            <div class="pt-summary">
              Total Consumed Cores (C₁): <strong>{{ formatNumber(pt1Cores, { decimals: 2 }) }} cores</strong>
            </div>
          </div>

          <!-- PT Point 2 -->
          <div class="pt-box">
            <h4>Titik Uji Beban 2 (PT High Load)</h4>
            <div class="form-group">
              <label>TPS Titik 2 (TPS₂)</label>
              <UiNumberInput
                v-model="params.pt2Tps"
                :step="100"
                :min="1"
                unit="TPS"
              />
            </div>
            <div class="form-group">
              <label>Node Spec (Cores/Node)</label>
              <UiNumberInput
                v-model="params.pt2NodeCores"
                :step="2"
                :min="1"
                unit="Cores"
              />
            </div>
            <div class="form-group">
              <label>Jumlah Node PT 2</label>
              <UiNumberInput
                v-model="params.pt2NodeQty"
                :step="1"
                :min="1"
                unit="Node"
              />
            </div>
            <div class="form-group">
              <label>Rata-rata % CPU PT 2</label>
              <UiNumberInput
                v-model="params.pt2CpuPercent"
                :step="5"
                :min="1"
                :max="100"
                unit="%"
              />
            </div>
            <div class="pt-summary">
              Total Consumed Cores (C₂): <strong>{{ formatNumber(pt2Cores, { decimals: 2 }) }} cores</strong>
            </div>
          </div>
        </div>

        <!-- Derived Model Constants -->
        <div class="derived-constants">
          <div class="constant-chip">
            <span class="chip-title">Eksponen Scaling (b)</span>
            <span class="chip-val">{{ formatNumber(scalingB, { decimals: 4 }) }}</span>
            <span class="chip-desc">b = ln(C₂/C₁) / ln(TPS₂/TPS₁)</span>
            <span :class="['scaling-tag', scalingB > 1.1 ? 'warn' : 'good']">
              {{ scalingB > 1.1 ? '⚠️ Non-linear Overhead' : '✅ Efficient Scaling' }}
            </span>
          </div>

          <div class="constant-chip">
            <span class="chip-title">Konstanta Scaling (a)</span>
            <span class="chip-val">{{ scalingA.toExponential(4) }}</span>
            <span class="chip-desc">a = C₁ / (TPS₁)ᵇ</span>
          </div>

          <div class="constant-chip highlight">
            <span class="chip-title">Kebutuhan Raw Cores</span>
            <span class="chip-val">{{ formatNumber(predictedCores, { decimals: 2 }) }} Cores</span>
            <span class="chip-desc">Cores(Predicted) = a · (TPS_peak)ᵇ</span>
          </div>
        </div>
      </div>

      <!-- Section 3: OpenShift Worker Node & Container Sizing -->
      <div class="calc-card">
        <div class="card-header">
          <span class="card-icon">☸️</span>
          <div>
            <h3>3. OpenShift Worker Node Sizing</h3>
            <p>Penentuan vCPU aman & jumlah fisik node</p>
          </div>
        </div>

        <div class="form-grid">
          <div class="form-group">
            <label>Max Safe CPU Utilization (U<sub>max</sub>)</label>
            <select v-model.number="params.maxCpuUtil">
              <option :value="0.30">30% (Sangat Aman / Critical SLA)</option>
              <option :value="0.40">40% (Standar Enterprise Enterprise Architecture)</option>
              <option :value="0.50">50% (Moderat)</option>
              <option :value="0.60">60% (Agresif)</option>
            </select>
          </div>

          <div class="form-group">
            <label>vCPU per Worker Node</label>
            <UiNumberInput
              v-model="params.workerNodeCores"
              :step="8"
              :min="1"
              unit="vCPU per baremetal / VM"
            />
          </div>
        </div>

        <!-- Utilization Progress Gauge -->
        <div class="util-gauge-container">
          <div class="gauge-label-row">
            <span>Batas Utilisasi Aman CPU: <strong>{{ (params.maxCpuUtil * 100).toFixed(0) }}%</strong></span>
            <span>Status Safety: <strong class="text-green">BOM Buffer Integrated</strong></span>
          </div>
          <div class="gauge-bar-bg">
            <div
              class="gauge-bar-fill"
              :style="{ width: `${params.maxCpuUtil * 100}%` }"
            ></div>
          </div>
        </div>

        <div class="results-grid-2">
          <div class="result-box">
            <div class="result-label">Total Cores Dibutuhkan (Buffer)</div>
            <div class="result-value text-blue">{{ formatNumber(requiredCores) }} <span class="unit">vCPU</span></div>
            <div class="result-sub">Required = ⌈ Predicted / U_max ⌉</div>
          </div>

          <div class="result-box result-success">
            <div class="result-label">Jumlah Worker Node (N + 1 HA)</div>
            <div class="result-value text-green">{{ formatNumber(workerNodeQty) }} <span class="unit">Nodes</span></div>
            <div class="result-sub">Nodes = ⌈ Cores / WorkerCores ⌉ + 1 HA</div>
          </div>
        </div>
      </div>

      <!-- Section 4: Database Storage & Redis RAM Sizing -->
      <div class="calc-card">
        <div class="card-header">
          <span class="card-icon">🗄️</span>
          <div>
            <h3>4. Database & Redis Cache Sizing</h3>
            <p>Proyeksi storage 5 tahun & RAM Redis Enterprise</p>
          </div>
        </div>

        <div class="form-grid">
          <div class="form-group">
            <label>Avg Row Size (KB)</label>
            <UiNumberInput
              v-model="params.avgRowSizeKb"
              :step="0.5"
              :min="0.1"
              :max-decimals="2"
              unit="KB per record"
            />
          </div>

          <div class="form-group">
            <label>Active Redis Keys (Millions)</label>
            <UiNumberInput
              v-model="params.redisKeysMillions"
              :step="0.5"
              :min="0.1"
              :max-decimals="2"
              unit="Juta keys"
            />
          </div>

          <div class="form-group">
            <label>Avg Redis Key Payload (KB)</label>
            <UiNumberInput
              v-model="params.redisPayloadKb"
              :step="0.5"
              :min="0.1"
              :max-decimals="2"
              unit="KB per key"
            />
          </div>
        </div>

        <div class="results-grid-2">
          <div class="result-box">
            <div class="result-label">Proyeksi Storage DB (5 Tahun)</div>
            <div class="result-value">{{ formatNumber(dbStorage5YearsGb, { maxDecimals: 1 }) }} <span class="unit">GB</span></div>
            <div class="result-sub">+30% Index +20% WAL/REDO Logs</div>
          </div>

          <div class="result-box">
            <div class="result-label">Redis RAM Cluster (1:1 Repl)</div>
            <div class="result-value">{{ formatNumber(redisRamGb, { maxDecimals: 1 }) }} <span class="unit">GB</span></div>
            <div class="result-sub">Dict Overhead 20%, Target Max 70%</div>
          </div>
        </div>
      </div>

      <!-- Section 5: 5-Year Executive Projection Table -->
      <div class="calc-card span-2">
        <div class="card-header">
          <span class="card-icon">📊</span>
          <div>
            <h3>5. Matriks Proyeksi Kapasitas 5 Tahun (2025–2029)</h3>
            <p>Tabel ekstrapolasi CAGR {{ params.growthRatePercent }}% per tahun untuk pengajuan EAWG & DAC</p>
          </div>
        </div>

        <div class="table-responsive">
          <table class="projection-table">
            <thead>
              <tr>
                <th>Tahun</th>
                <th>Daily Trx</th>
                <th>Peak TPS</th>
                <th>Raw vCPU</th>
                <th>vCPU (Safe Buffer)</th>
                <th>Worker Nodes</th>
                <th>DB Storage Cumulative</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="row in projectionRows" :key="row.year" :class="{ 'row-highlight': row.year === 2025 }">
                <td class="td-year">
                  <strong>{{ row.year }}</strong>
                  <span v-if="row.year === 2025" class="year-badge">Base</span>
                </td>
                <td>{{ formatNumber(row.dailyTrx) }}</td>
                <td class="td-tps">{{ formatNumber(row.peakTps, { maxDecimals: 1 }) }} TPS</td>
                <td>{{ formatNumber(row.predictedCores, { maxDecimals: 1 }) }} cores</td>
                <td class="td-cores"><strong>{{ formatNumber(row.requiredCores) }} vCPU</strong></td>
                <td class="td-nodes"><span class="node-pill">{{ formatNumber(row.workerNodes) }} Nodes</span></td>
                <td class="td-storage">{{ formatNumber(row.dbStorageGb, { maxDecimals: 0 }) }} GB</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- Footer Summary Actions -->
    <div class="calc-summary-bar">
      <div class="summary-info">
        <span>📌 Target Peak Load: <strong>{{ formatNumber(peakTps, { maxDecimals: 1 }) }} TPS</strong></span>
        <span>🖥️ OpenShift Workers: <strong>{{ formatNumber(workerNodeQty) }} Nodes ({{ formatNumber(requiredCores) }} vCPU)</strong></span>
        <span>💾 Storage 5-Yr: <strong>{{ formatNumber(dbStorage5YearsGb, { maxDecimals: 0 }) }} GB</strong></span>
        <span>⚡ Power-Law (b): <strong>{{ formatNumber(scalingB, { decimals: 4 }) }}</strong></span>
      </div>

      <div class="summary-actions">
        <button type="button" class="copy-btn secondary-btn" @click="downloadJson">
          💾 Download JSON Spec
        </button>
        <button type="button" class="copy-btn" @click="copySummary">
          📋 Copy Ringkasan BOM (EAWG Format)
        </button>
      </div>
    </div>
  </div>

    <!-- VIEW MODE 2: Official EAWG Bab 10 Format -->
    <div v-else-if="activeViewMode === 'eawg'" class="eawg-view-container">
      <TemplateCapacityEawgView />
    </div>
  </div>
</template>

<script setup lang="ts">
import { useLocalOpen } from '~/composables/useLocalOpen'
import { formatNumber } from '~/composables/useFormat'

const { showToast } = useLocalOpen()

const activeViewMode = ref<'engine' | 'eawg'>('eawg')
const activePreset = ref<string>('lm')

const presets = [
  { id: 'lm', name: "Livin' Merchant (LM)", icon: '🛍️' },
  { id: 'kopra', name: 'Kopra Cash Management', icon: '🏦' },
  { id: 'rmt', name: 'RM Tools (RMT)', icon: '💼' },
  { id: 'custom', name: 'Custom Parameters', icon: '⚙️' },
]

const params = reactive({
  trxDaily: 15000000,
  peakHours: 8,
  concurrencyFactor: 1.8,
  growthRatePercent: 15,

  // PT Point 1
  pt1Tps: 2000,
  pt1NodeCores: 16,
  pt1NodeQty: 4,
  pt1CpuPercent: 35,

  // PT Point 2
  pt2Tps: 4500,
  pt2NodeCores: 16,
  pt2NodeQty: 4,
  pt2CpuPercent: 78,

  // Safety & Infrastructure Specs
  maxCpuUtil: 0.40,
  workerNodeCores: 32,
  avgRowSizeKb: 2.5,

  // Redis
  redisKeysMillions: 5.0,
  redisPayloadKb: 4.0,
})

function applyPreset(presetId: string) {
  activePreset.value = presetId
  if (presetId === 'lm') {
    params.trxDaily = 15000000
    params.peakHours = 8
    params.concurrencyFactor = 1.8
    params.growthRatePercent = 15
    params.pt1Tps = 2000
    params.pt1NodeCores = 16
    params.pt1NodeQty = 4
    params.pt1CpuPercent = 35
    params.pt2Tps = 4500
    params.pt2NodeCores = 16
    params.pt2NodeQty = 4
    params.pt2CpuPercent = 78
    params.maxCpuUtil = 0.40
    params.workerNodeCores = 32
    params.avgRowSizeKb = 2.5
    params.redisKeysMillions = 5.0
    params.redisPayloadKb = 4.0
    showToast("🛍️ Preset Livin' Merchant dimuat!", 'info')
  } else if (presetId === 'kopra') {
    params.trxDaily = 8500000
    params.peakHours = 7
    params.concurrencyFactor = 2.0
    params.growthRatePercent = 20
    params.pt1Tps = 1200
    params.pt1NodeCores = 16
    params.pt1NodeQty = 4
    params.pt1CpuPercent = 28
    params.pt2Tps = 3000
    params.pt2NodeCores = 16
    params.pt2NodeQty = 4
    params.pt2CpuPercent = 68
    params.maxCpuUtil = 0.40
    params.workerNodeCores = 32
    params.avgRowSizeKb = 3.5
    params.redisKeysMillions = 8.0
    params.redisPayloadKb = 5.0
    showToast('🏦 Preset Kopra Cash Management dimuat!', 'info')
  } else if (presetId === 'rmt') {
    params.trxDaily = 3200000
    params.peakHours = 9
    params.concurrencyFactor = 1.5
    params.growthRatePercent = 12
    params.pt1Tps = 800
    params.pt1NodeCores = 8
    params.pt1NodeQty = 4
    params.pt1CpuPercent = 22
    params.pt2Tps = 1800
    params.pt2NodeCores = 8
    params.pt2NodeQty = 4
    params.pt2CpuPercent = 55
    params.maxCpuUtil = 0.40
    params.workerNodeCores = 16
    params.avgRowSizeKb = 1.8
    params.redisKeysMillions = 2.5
    params.redisPayloadKb = 3.0
    showToast('💼 Preset RM Tools (RMT) dimuat!', 'info')
  }
}

// 1. Peak TPS calculation
const peakTps = computed(() => {
  const seconds = params.peakHours * 3600
  if (seconds <= 0) return 0
  return (params.trxDaily * params.concurrencyFactor) / seconds
})

// 2. Consumed cores in PT
const pt1Cores = computed(() => {
  return (params.pt1CpuPercent / 100) * params.pt1NodeCores * params.pt1NodeQty
})

const pt2Cores = computed(() => {
  return (params.pt2CpuPercent / 100) * params.pt2NodeCores * params.pt2NodeQty
})

// 3. Power-Law Scaling Exponent (b) & Constant (a)
const scalingB = computed(() => {
  const c1 = pt1Cores.value
  const c2 = pt2Cores.value
  const tps1 = params.pt1Tps
  const tps2 = params.pt2Tps

  if (c1 <= 0 || c2 <= 0 || tps1 <= 0 || tps2 <= 0 || tps1 === tps2) return 1.0
  return Math.log(c2 / c1) / Math.log(tps2 / tps1)
})

const scalingA = computed(() => {
  const c1 = pt1Cores.value
  const tps1 = params.pt1Tps
  const b = scalingB.value

  if (tps1 <= 0) return 0
  return c1 / Math.pow(tps1, b)
})

// 4. Predicted Raw Cores for Target Peak TPS
const predictedCores = computed(() => {
  const a = scalingA.value
  const b = scalingB.value
  const tps = peakTps.value

  if (tps <= 0) return 0
  return a * Math.pow(tps, b)
})

// 5. Required Cores with Max CPU Utilization Buffer (e.g. 40%)
const requiredCores = computed(() => {
  if (params.maxCpuUtil <= 0) return 0
  return Math.ceil(predictedCores.value / params.maxCpuUtil)
})

// 6. Worker Node Quantity (with +1 HA Node)
const workerNodeQty = computed(() => {
  if (params.workerNodeCores <= 0) return 1
  return Math.ceil(requiredCores.value / params.workerNodeCores) + 1
})

// 7. DB Storage 5 Years
const dbStorage5YearsGb = computed(() => {
  const dailyGb = (params.trxDaily * params.avgRowSizeKb * 1.3) / (1024 * 1024)
  const annualGb = dailyGb * 365 * 1.2
  return annualGb * 5
})

// 8. Redis RAM Sizing
const redisRamGb = computed(() => {
  const totalKeys = params.redisKeysMillions * 1000000
  const rawMb = (totalKeys * params.redisPayloadKb * 1.2) / 1024
  const rawGb = rawMb / 1024
  return (rawGb * 2) / 0.7
})

// 9. 5-Year Projections Rows
const projectionRows = computed(() => {
  const rows = []
  const startYear = 2025
  const g = (params.growthRatePercent || 0) / 100

  let cumulativeStorage = 0

  for (let i = 0; i < 5; i++) {
    const year = startYear + i
    const dailyTrx = Math.round(params.trxDaily * Math.pow(1 + g, i))
    const seconds = params.peakHours * 3600
    const yearTps = (dailyTrx * params.concurrencyFactor) / seconds
    const predCores = scalingA.value * Math.pow(yearTps, scalingB.value)
    const reqCores = Math.ceil(predCores / params.maxCpuUtil)
    const workerNodes = Math.ceil(reqCores / params.workerNodeCores) + 1

    const yearDailyGb = (dailyTrx * params.avgRowSizeKb * 1.3) / (1024 * 1024)
    const yearAnnualGb = yearDailyGb * 365 * 1.2
    cumulativeStorage += yearAnnualGb

    rows.push({
      year,
      dailyTrx,
      peakTps: yearTps,
      predictedCores: predCores,
      requiredCores: reqCores,
      workerNodes,
      dbStorageGb: cumulativeStorage,
    })
  }

  return rows
})

function copySummary() {
  const text = `=== RINGKASAN KAPASITAS ENTERPRISE (EAWG FORMAT) ===
- Presets/Project: ${presets.find(p => p.id === activePreset.value)?.name || 'Custom'}
- Daily Transactions (Base): ${formatNumber(params.trxDaily)} trx/day @ ${params.peakHours}h (${params.concurrencyFactor}x spike)
- Target Peak TPS: ${formatNumber(peakTps.value, { maxDecimals: 1 })} TPS
- Power-Law Model: a=${scalingA.value.toExponential(4)}, b=${formatNumber(scalingB.value, { decimals: 4 })}
- Predicted Raw vCPU: ${formatNumber(predictedCores.value, { decimals: 2 })} Cores
- Required vCPU (Safe Buffer ${params.maxCpuUtil * 100}%): ${formatNumber(requiredCores.value)} Cores
- OpenShift Worker Nodes (${params.workerNodeCores} vCPU/node): ${formatNumber(workerNodeQty.value)} Nodes (including +1 HA)
- DB Storage Projection (5-Year): ${formatNumber(dbStorage5YearsGb.value, { maxDecimals: 1 })} GB
- Redis Enterprise RAM: ${formatNumber(redisRamGb.value, { maxDecimals: 1 })} GB

=== PROYEKSI 5 TAHUN (2025-2029) ===
${projectionRows.value
  .map(
    r =>
      `${r.year}: ${formatNumber(r.dailyTrx)} trx/day -> ${formatNumber(r.peakTps, { maxDecimals: 1 })} TPS | ${formatNumber(r.requiredCores)} vCPU (${formatNumber(r.workerNodes)} Nodes) | DB ${formatNumber(r.dbStorageGb, { maxDecimals: 0 })} GB`
  )
  .join('\n')}`

  navigator.clipboard.writeText(text)
  showToast('✅ Ringkasan Kalkulasi Kapasitas berhasil disalin ke clipboard!', 'success')
}

function downloadJson() {
  const payload = {
    timestamp: new Date().toISOString(),
    preset: activePreset.value,
    parameters: params,
    results: {
      peakTps: peakTps.value,
      scalingA: scalingA.value,
      scalingB: scalingB.value,
      predictedCores: predictedCores.value,
      requiredCores: requiredCores.value,
      workerNodeQty: workerNodeQty.value,
      dbStorage5YearsGb: dbStorage5YearsGb.value,
      redisRamGb: redisRamGb.value,
      projectionRows: projectionRows.value,
    },
  }

  const dataStr = 'data:text/json;charset=utf-8,' + encodeURIComponent(JSON.stringify(payload, null, 2))
  const downloadAnchor = document.createElement('a')
  downloadAnchor.setAttribute('href', dataStr)
  downloadAnchor.setAttribute('download', `capacity_plan_spec_${activePreset.value}_${new Date().toISOString().slice(0, 10)}.json`)
  document.body.appendChild(downloadAnchor)
  downloadAnchor.click()
  downloadAnchor.remove()

  showToast('💾 Berkas JSON Spesifikasi Kapasitas berhasil diunduh!', 'success')
}
</script>

<style scoped>
.capacity-calculator {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.calc-header {
  background: linear-gradient(135deg, #0f2b48 0%, #003b73 100%);
  color: #ffffff;
  padding: 24px 28px;
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-md);
}

.calc-header-main {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.calc-badge {
  display: inline-block;
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.5px;
  text-transform: uppercase;
  padding: 4px 10px;
  background: rgba(255, 255, 255, 0.15);
  border: 1px solid rgba(255, 255, 255, 0.25);
  border-radius: 20px;
  margin-bottom: 8px;
}

.calc-title {
  font-size: 20px;
  font-weight: 700;
  margin-bottom: 6px;
  color: #ffffff;
}

.calc-subtitle {
  font-size: 13px;
  color: #cbd5e1;
  max-width: 800px;
  line-height: 1.5;
}

.view-mode-switch {
  display: flex;
  align-items: center;
  gap: 8px;
  padding-top: 12px;
  border-top: 1px solid rgba(255, 255, 255, 0.15);
  flex-wrap: wrap;
}

.mode-tab-btn {
  font-family: 'Inter', sans-serif;
  font-size: 12px;
  font-weight: 700;
  padding: 8px 16px;
  border-radius: var(--radius-sm);
  border: 1px solid rgba(255, 255, 255, 0.25);
  background: rgba(255, 255, 255, 0.1);
  color: #ffffff;
  cursor: pointer;
  transition: all 0.15s;
}

.mode-tab-btn:hover {
  background: rgba(255, 255, 255, 0.2);
}

.mode-tab-btn.active {
  background: #ffffff;
  color: var(--navy);
  border-color: #ffffff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
}

.engine-view-container,
.eawg-view-container {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.preset-selector {
  display: flex;
  align-items: center;
  gap: 12px;

  padding-top: 12px;
  border-top: 1px solid rgba(255, 255, 255, 0.15);
  flex-wrap: wrap;
}

.preset-label {
  font-size: 12px;
  font-weight: 600;
  color: #cbd5e1;
}

.preset-buttons {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.preset-btn {
  font-family: 'Inter', sans-serif;
  font-size: 11px;
  font-weight: 600;
  padding: 6px 12px;
  border-radius: 20px;
  border: 1px solid rgba(255, 255, 255, 0.25);
  background: rgba(255, 255, 255, 0.1);
  color: #ffffff;
  cursor: pointer;
  transition: all 0.15s;
}

.preset-btn:hover {
  background: rgba(255, 255, 255, 0.25);
}

.preset-btn.active {
  background: #ffffff;
  color: var(--navy);
  border-color: #ffffff;
  font-weight: 700;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
}

.calc-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 20px;
}

.span-2 {
  grid-column: span 2;
}

.calc-card {
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius-lg);
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 16px;
  box-shadow: var(--shadow-sm);
}

.card-header {
  display: flex;
  align-items: center;
  gap: 12px;
  padding-bottom: 12px;
  border-bottom: 1px solid var(--border);
}

.card-icon {
  font-size: 24px;
  width: 42px;
  height: 42px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--bg3);
  border-radius: 8px;
}

.card-header h3 {
  font-size: 15px;
  font-weight: 700;
  color: var(--text);
}

.card-header p {
  font-size: 12px;
  color: var(--text3);
}

.form-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  gap: 12px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.form-group label {
  font-size: 11px;
  font-weight: 600;
  color: var(--text2);
}

.form-group label sub {
  font-size: 9px;
  color: var(--text3);
}

.form-group input,
.form-group select {
  font-family: 'Inter', sans-serif;
  font-size: 13px;
  padding: 8px 10px;
  border: 1px solid var(--border2);
  border-radius: var(--radius-sm);
  background: var(--bg2);
  color: var(--text);
  outline: none;
  transition: border-color 0.15s;
}

.form-group input:focus,
.form-group select:focus {
  border-color: var(--blue);
}

.input-unit {
  font-size: 10px;
  color: var(--text3);
}

.result-box {
  background: var(--bg3);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  padding: 14px;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.result-highlight {
  background: #eff6ff;
  border-color: #bfdbfe;
}

.result-success {
  background: #f0fdf4;
  border-color: #bbf7d0;
}

.result-label {
  font-size: 11px;
  font-weight: 600;
  color: var(--text3);
  text-transform: uppercase;
  letter-spacing: 0.3px;
}

.result-value {
  font-size: 24px;
  font-weight: 800;
  color: var(--navy);
}

.result-value .unit {
  font-size: 14px;
  font-weight: 600;
  color: var(--text3);
}

.text-blue {
  color: #2563eb;
}

.text-green {
  color: #16a34a;
}

.result-formula,
.result-sub {
  font-size: 11px;
  color: var(--text3);
  font-family: 'JetBrains Mono', monospace;
}

.pt-points-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.pt-box {
  background: var(--bg);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  padding: 14px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.pt-box h4 {
  font-size: 13px;
  font-weight: 700;
  color: var(--navy);
  padding-bottom: 6px;
  border-bottom: 1px solid var(--border);
}

.pt-summary {
  font-size: 12px;
  color: var(--text2);
  background: var(--bg2);
  padding: 6px 10px;
  border-radius: 4px;
  border: 1px solid var(--border);
}

.derived-constants {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
  margin-top: 8px;
}

.constant-chip {
  background: var(--bg2);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  padding: 12px;
  display: flex;
  flex-direction: column;
  gap: 4px;
  position: relative;
}

.constant-chip.highlight {
  background: #fef3c7;
  border-color: #fde68a;
}

.chip-title {
  font-size: 11px;
  font-weight: 600;
  color: var(--text3);
}

.chip-val {
  font-size: 18px;
  font-weight: 700;
  font-family: 'JetBrains Mono', monospace;
  color: var(--navy);
}

.chip-desc {
  font-size: 10px;
  color: var(--text3);
}

.scaling-tag {
  font-size: 9px;
  font-weight: 700;
  padding: 2px 6px;
  border-radius: 4px;
  width: fit-content;
}

.scaling-tag.good {
  background: #f0fdf4;
  color: #16a34a;
  border: 1px solid #bbf7d0;
}

.scaling-tag.warn {
  background: #fffbeb;
  color: #d97706;
  border: 1px solid #fde68a;
}

.util-gauge-container {
  display: flex;
  flex-direction: column;
  gap: 6px;
  background: var(--bg3);
  padding: 12px;
  border-radius: var(--radius);
  border: 1px solid var(--border);
}

.gauge-label-row {
  display: flex;
  justify-content: space-between;
  font-size: 11px;
  color: var(--text2);
}

.gauge-bar-bg {
  height: 8px;
  background: var(--border2);
  border-radius: 4px;
  overflow: hidden;
}

.gauge-bar-fill {
  height: 100%;
  background: linear-gradient(90deg, #16a34a 0%, #2563eb 100%);
  border-radius: 4px;
  transition: width 0.3s ease;
}

.results-grid-2 {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}

.table-responsive {
  overflow-x: auto;
}

.projection-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 12px;
}

.projection-table th,
.projection-table td {
  padding: 10px 12px;
  text-align: left;
  border-bottom: 1px solid var(--border);
}

.projection-table th {
  background: var(--bg3);
  color: var(--text2);
  font-weight: 700;
  font-size: 11px;
  text-transform: uppercase;
}

.projection-table tr:hover {
  background: var(--bg);
}

.row-highlight {
  background: #eff6ff;
}

.td-year {
  display: flex;
  align-items: center;
  gap: 6px;
}

.year-badge {
  font-size: 9px;
  font-weight: 700;
  background: #2563eb;
  color: #ffffff;
  padding: 1px 5px;
  border-radius: 3px;
}

.td-tps {
  font-family: 'JetBrains Mono', monospace;
  font-weight: 600;
  color: var(--navy);
}

.td-cores {
  color: #2563eb;
}

.node-pill {
  background: #f0fdf4;
  color: #16a34a;
  border: 1px solid #bbf7d0;
  padding: 2px 8px;
  border-radius: 12px;
  font-weight: 700;
  font-size: 11px;
}

.td-storage {
  font-family: 'JetBrains Mono', monospace;
  color: var(--text2);
}

.calc-summary-bar {
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius-lg);
  padding: 16px 20px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  box-shadow: var(--shadow-sm);
  flex-wrap: wrap;
  gap: 12px;
}

.summary-info {
  display: flex;
  align-items: center;
  gap: 20px;
  font-size: 12px;
  color: var(--text2);
  flex-wrap: wrap;
}

.summary-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.copy-btn {
  font-family: 'Inter', sans-serif;
  font-size: 12px;
  font-weight: 600;
  padding: 8px 16px;
  background: var(--brand-blue);
  color: #ffffff;
  border: none;
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: background 0.15s;
}

.copy-btn:hover {
  background: var(--navy);
}

.secondary-btn {
  background: var(--bg3);
  color: var(--text);
  border: 1px solid var(--border2);
}

.secondary-btn:hover {
  background: var(--border2);
}

@media (max-width: 900px) {
  .calc-grid {
    grid-template-columns: 1fr;
  }
  .span-2 {
    grid-column: span 1;
  }
  .pt-points-grid {
    grid-template-columns: 1fr;
  }
  .derived-constants {
    grid-template-columns: 1fr;
  }
  .calc-summary-bar {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>
