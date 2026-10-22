<template>
  <div class="capacity-eawg-view">
    <!-- Preset & Mode Toolbar -->
    <div class="eawg-toolbar">
      <div class="toolbar-left">
        <span class="doc-badge">🏛️ Enterprise Architecture EAWG Standard Format</span>
        <h3 class="doc-title-mini">Bab 10. Business Capacity Plan Deliverable</h3>
      </div>

      <!-- Controls Right -->
      <div class="toolbar-controls">
        <!-- Preset Selector -->
        <div class="preset-selector">
          <span class="preset-label">Proyek:</span>
          <div class="preset-pills">
            <button
              v-for="p in eawgPresets"
              :key="p.id"
              type="button"
              :class="['preset-pill', { active: selectedPresetId === p.id }]"
              @click="selectPreset(p.id)"
            >
              {{ p.icon }} {{ p.name }}
            </button>
          </div>
        </div>

        <!-- Mode Toggles -->
        <div class="mode-toggles">
          <button
            type="button"
            :class="['toggle-btn', { active: isEditMode }]"
            @click="isEditMode = !isEditMode"
            :title="isEditMode ? 'Kunci tampilan (Read-Only Preview)' : 'Buka input untuk ubah angka'"
          >
            {{ isEditMode ? '🔒 Selesai Edit' : '✏️ Mode Edit & Adjust' }}
          </button>
          <button
            type="button"
            :class="['toggle-btn grafana-btn', { active: showGrafanaModal || isGrafanaSynced }]"
            @click="showGrafanaModal = true"
            title="Tarik metrik utilisasi realtime dari Grafana API"
          >
            📡 Grafana API Sync
          </button>
          <button
            type="button"
            :class="['toggle-btn', { active: showGrafanaCols }]"
            @click="showGrafanaCols = !showGrafanaCols"
            title="Tampilkan / Sembunyikan kolom utilisasi Grafana di tabel"
          >
            📊 Kolom Utilisasi
          </button>
        </div>
      </div>
    </div>

    <!-- Live Grafana Sync Status Banner -->
    <div v-if="isGrafanaSynced" class="grafana-live-banner">
      <div class="banner-left">
        <span class="pulse-dot"></span>
        <span class="banner-title">📡 Terhubung ke Grafana API Monitoring:</span>
        <span class="banner-meta">Span: <strong>{{ lastSyncMeta.timeSpanLabel }}</strong></span>
        <span class="banner-meta">Agregasi: <strong>{{ lastSyncMeta.aggregation.toUpperCase() }}</strong></span>
        <span class="banner-meta">Host Tersinkron: <strong>{{ activeData.hostSizing.length }} Node</strong></span>
      </div>
      <div class="banner-right">
        <span class="banner-time">Terakhir update: {{ lastSyncTimeFormatted }}</span>
        <button type="button" class="btn-resync" @click="executeGrafanaSync">
          🔄 Sync Ulang
        </button>
      </div>
    </div>

    <!-- Grafana & APM Utilization Assistant Card (Expandable) -->
    <div v-if="showGrafanaCols || isGrafanaSynced" class="grafana-assistant-card">
      <div class="grafana-header">
        <div class="grafana-title-wrap">
          <span class="grafana-icon">📈</span>
          <div>
            <h4>Grafana & APM Monitoring Sizing Assistant</h4>
            <p>
              Tarik data utilisasi rata-rata / puncak CPU & RAM dari Grafana (span 1 tahun / kustom) untuk memproyeksikan kebutuhan sizing Year 2 & Year 3 secara presisi.
            </p>
          </div>
        </div>
        <div class="grafana-config-grid">
          <div class="grafana-field">
            <label>Batas Aman CPU (U<sub>max</sub>)</label>
            <select v-model.number="grafanaConfig.safeCpuUtil">
              <option :value="0.30">30% (Critical SLA)</option>
              <option :value="0.40">40% (Standar Enterprise Architecture)</option>
              <option :value="0.50">50% (Moderat)</option>
              <option :value="0.60">60% (Agresif)</option>
            </select>
          </div>
          <div class="grafana-field">
            <label>Batas Aman RAM (U<sub>ram</sub>)</label>
            <select v-model.number="grafanaConfig.safeRamUtil">
              <option :value="0.60">60% (Sangat Aman)</option>
              <option :value="0.70">70% (Standar Enterprise)</option>
              <option :value="0.80">80% (Agresif)</option>
            </select>
          </div>
          <div class="grafana-field">
            <label>Asumsi CAGR Pertumbuhan Bisnis</label>
            <div class="input-with-suffix">
              <input v-model.number="grafanaConfig.annualGrowthPercent" type="number" min="0" max="100" />
              <span>% / thn</span>
            </div>
          </div>
          <div class="grafana-action">
            <button
              type="button"
              class="grafana-apply-btn"
              @click="applyGrafanaExtrapolationAll"
            >
              ⚡ Auto-Scale Y2 & Y3 dari Grafana
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Official Document Container (Paper View) -->
    <div class="eawg-paper">
      <!-- Section 10 Header -->
      <div class="doc-section-header">
        <div class="section-title-wrap">
          <h2 class="section-num">10. Business Capacity Plan</h2>
          <span v-if="isEditMode" class="edit-badge">✏️ Mode Edit Aktif — Klik angka / teks untuk mengubah</span>
        </div>
        <p class="section-intro">
          Detail capacity plan berdasarkan proyeksi bisnis tertera pada EAWG, Asumsi bisnis terkait
          <span v-if="!isEditMode"><strong>{{ activeData.projectName }}</strong></span>
          <input v-else v-model="activeData.projectName" class="inline-text-input project-title-input" />
          adalah sebagai berikut:
        </p>
      </div>

      <!-- Table 1: Business Metrics Projection 5-Year -->
      <div class="table-container">
        <div class="table-title-row">
          <span class="table-caption">Tabel 10.1: Matriks Proyeksi Pertumbuhan Bisnis (2025–2029)</span>
          <div class="table-actions">
            <button v-if="isEditMode" type="button" class="mini-btn add-btn" @click="addBusinessRow">
              + Tambah Baris Metrik
            </button>
            <button type="button" class="mini-btn copy-btn" @click="copyBusinessTable">
              📋 Salin Tabel Bisnis
            </button>
          </div>
        </div>

        <div class="table-responsive">
          <table class="eawg-table business-table">
            <thead>
              <tr>
                <th class="th-channel">Channel / Modul</th>
                <th class="th-metrics">Metrics</th>
                <th>2025</th>
                <th>2026</th>
                <th>2027</th>
                <th>2028</th>
                <th>2029</th>
                <th class="th-growth">Growth</th>
                <th v-if="isEditMode" class="th-action">Aksi</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="(row, idx) in activeData.businessMetrics"
                :key="idx"
                :class="{ 'group-start': isNewChannel(idx) }"
              >
                <!-- Channel -->
                <td class="td-channel">
                  <span v-if="!isEditMode"><strong>{{ row.channel }}</strong></span>
                  <input v-else v-model="row.channel" class="inline-text-input channel-input" />
                </td>

                <!-- Metric Name -->
                <td class="td-metric">
                  <span v-if="!isEditMode">{{ row.metric }}</span>
                  <input v-else v-model="row.metric" class="inline-text-input metric-input" />
                </td>

                <!-- 2025 -->
                <td class="td-num">
                  <span v-if="!isEditMode">{{ formatNumber(row.y2025) }}</span>
                  <UiNumberInput
                    v-else
                    v-model="row.y2025"
                    :show-stepper="false"
                    class="table-number-input"
                  />
                </td>

                <!-- 2026 -->
                <td class="td-num">
                  <span v-if="!isEditMode">{{ formatNumber(row.y2026) }}</span>
                  <UiNumberInput
                    v-else
                    v-model="row.y2026"
                    :show-stepper="false"
                    class="table-number-input"
                  />
                </td>

                <!-- 2027 -->
                <td class="td-num">
                  <span v-if="!isEditMode">{{ formatNumber(row.y2027) }}</span>
                  <UiNumberInput
                    v-else
                    v-model="row.y2027"
                    :show-stepper="false"
                    class="table-number-input"
                  />
                </td>

                <!-- 2028 -->
                <td class="td-num">
                  <span v-if="!isEditMode">{{ formatNumber(row.y2028) }}</span>
                  <UiNumberInput
                    v-else
                    v-model="row.y2028"
                    :show-stepper="false"
                    class="table-number-input"
                  />
                </td>

                <!-- 2029 -->
                <td class="td-num">
                  <span v-if="!isEditMode">{{ formatNumber(row.y2029) }}</span>
                  <UiNumberInput
                    v-else
                    v-model="row.y2029"
                    :show-stepper="false"
                    class="table-number-input"
                  />
                </td>

                <!-- Growth -->
                <td class="td-growth">
                  <span v-if="!isEditMode" class="growth-badge">{{ row.growthPercent }}%</span>
                  <div v-else class="growth-edit-wrap">
                    <input
                      v-model.number="row.growthPercent"
                      type="number"
                      class="growth-input"
                    />
                    <span>%</span>
                    <button
                      type="button"
                      class="btn-calc-growth"
                      title="Hitung otomatis 2026-2029 dari 2025 & % Growth"
                      @click="autoProjectRow(idx)"
                    >
                      ⚡
                    </button>
                  </div>
                </td>

                <!-- Action Button in Edit Mode -->
                <td v-if="isEditMode" class="td-action">
                  <button
                    type="button"
                    class="btn-del-row"
                    title="Hapus baris"
                    @click="deleteBusinessRow(idx)"
                  >
                    🗑️
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Subtitle for Table 2 -->
      <div class="doc-sub-intro">
        <p>Capacity plan terkait <strong>{{ activeData.name }}</strong> adalah sebagai berikut:</p>
      </div>

      <!-- Table 2: Infrastructure Sizing BOM / Multi-Year Host Inventory -->
      <div class="table-container">
        <div class="table-title-row">
          <span class="table-caption">Tabel 10.2: Matriks Kebutuhan Infrastruktur & Host Sizing Multi-Year (Year 1, 2, 3)</span>
          <div class="table-actions">
            <button v-if="isEditMode" type="button" class="mini-btn add-btn" @click="addHostRow">
              + Tambah Host / VM
            </button>
            <button type="button" class="mini-btn copy-btn" @click="copyHostTable">
              📋 Salin Tabel Sizing
            </button>
          </div>
        </div>

        <div class="table-responsive">
          <table class="eawg-table host-table">
            <thead>
              <!-- Double Header Row 1 -->
              <tr>
                <th rowspan="2" class="th-envi">Envi</th>
                <th rowspan="2" class="th-usage">Penggunaan</th>
                <th rowspan="2" class="th-hostname">Host Name</th>
                <th rowspan="2" class="th-node">Node</th>

                <!-- Optional Grafana Header -->
                <th v-if="showGrafanaCols || isGrafanaSynced" colspan="3" class="th-grafana">
                  📊 Grafana Util ({{ lastSyncMeta.timeSpanLabel }})
                </th>

                <th colspan="3" class="th-year y1-header">Year 1</th>
                <th colspan="3" class="th-year y2-header">Year 2</th>
                <th colspan="3" class="th-year y3-header">Year 3</th>
                <th v-if="isEditMode" rowspan="2" class="th-action">Aksi</th>
              </tr>

              <!-- Double Header Row 2 -->
              <tr class="sub-header-row">
                <!-- Grafana Sub-Headers -->
                <th v-if="showGrafanaCols || isGrafanaSynced" class="sub-th th-graf-cpu">CPU %</th>
                <th v-if="showGrafanaCols || isGrafanaSynced" class="sub-th th-graf-ram">RAM %</th>
                <th v-if="showGrafanaCols || isGrafanaSynced" class="sub-th th-graf-status">Status</th>

                <!-- Year 1 -->
                <th class="sub-th y1-col">CPU</th>
                <th class="sub-th y1-col">Memory</th>
                <th class="sub-th y1-col">Disk</th>

                <!-- Year 2 -->
                <th class="sub-th y2-col">CPU</th>
                <th class="sub-th y2-col">Memory</th>
                <th class="sub-th y2-col">Disk</th>

                <!-- Year 3 -->
                <th class="sub-th y3-col">CPU</th>
                <th class="sub-th y3-col">Memory</th>
                <th class="sub-th y3-col">Disk</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="(row, hIdx) in activeData.hostSizing"
                :key="row.id"
                :class="getRowCategoryClass(row.usage)"
              >
                <!-- Envi -->
                <td class="td-envi">
                  <span v-if="!isEditMode" :class="['envi-tag', row.envi.toLowerCase()]">{{ row.envi }}</span>
                  <select v-else v-model="row.envi" class="envi-select">
                    <option value="DC">DC</option>
                    <option value="DRC">DRC</option>
                  </select>
                </td>

                <!-- Penggunaan / Role -->
                <td class="td-usage">
                  <span class="usage-icon">{{ getUsageIcon(row.usage) }}</span>
                  <span v-if="!isEditMode" class="usage-text">{{ row.usage }}</span>
                  <input v-else v-model="row.usage" class="inline-text-input usage-input" />
                </td>

                <!-- Host Name -->
                <td class="td-hostname">
                  <code v-if="!isEditMode">{{ row.hostName }}</code>
                  <input v-else v-model="row.hostName" class="inline-text-input host-input" />
                </td>

                <!-- Node Number -->
                <td class="td-node-val">
                  <span v-if="!isEditMode">{{ row.node }}</span>
                  <input v-else v-model.number="row.node" type="number" class="node-input" />
                </td>

                <!-- Optional Grafana Inputs -->
                <td v-if="showGrafanaCols || isGrafanaSynced" class="td-grafana-cell">
                  <div class="grafana-input-wrap">
                    <input
                      v-model.number="row.currentCpuUtil"
                      type="number"
                      min="1"
                      max="100"
                      placeholder="e.g. 45"
                      class="graf-input"
                      :disabled="!isEditMode"
                    />
                    <span>%</span>
                  </div>
                </td>
                <td v-if="showGrafanaCols || isGrafanaSynced" class="td-grafana-cell">
                  <div class="grafana-input-wrap">
                    <input
                      v-model.number="row.currentMemUtil"
                      type="number"
                      min="1"
                      max="100"
                      placeholder="e.g. 65"
                      class="graf-input"
                      :disabled="!isEditMode"
                    />
                    <span>%</span>
                  </div>
                </td>
                <td v-if="showGrafanaCols || isGrafanaSynced" class="td-grafana-status">
                  <span :class="['status-pill', getUtilStatus(row.currentCpuUtil, row.currentMemUtil).class]">
                    {{ getUtilStatus(row.currentCpuUtil, row.currentMemUtil).label }}
                  </span>
                </td>

                <!-- Year 1 Specs -->
                <td class="td-spec y1-val">
                  <span v-if="!isEditMode">{{ formatNumber(row.y1.cpu) }}</span>
                  <UiNumberInput
                    v-else
                    v-model="row.y1.cpu"
                    :show-stepper="false"
                    class="table-number-input"
                  />
                </td>
                <td class="td-spec y1-val">
                  <span v-if="!isEditMode">{{ formatNumber(row.y1.memory) }}</span>
                  <UiNumberInput
                    v-else
                    v-model="row.y1.memory"
                    :show-stepper="false"
                    class="table-number-input"
                  />
                </td>
                <td class="td-spec y1-val disk-val">
                  <span v-if="!isEditMode">{{ formatNumber(row.y1.disk) }}</span>
                  <UiNumberInput
                    v-else
                    v-model="row.y1.disk"
                    :show-stepper="false"
                    class="table-number-input"
                  />
                </td>

                <!-- Year 2 Specs -->
                <td class="td-spec y2-val">
                  <span v-if="!isEditMode">{{ formatNumber(row.y2.cpu) }}</span>
                  <UiNumberInput
                    v-else
                    v-model="row.y2.cpu"
                    :show-stepper="false"
                    class="table-number-input"
                  />
                </td>
                <td class="td-spec y2-val">
                  <span v-if="!isEditMode">{{ formatNumber(row.y2.memory) }}</span>
                  <UiNumberInput
                    v-else
                    v-model="row.y2.memory"
                    :show-stepper="false"
                    class="table-number-input"
                  />
                </td>
                <td class="td-spec y2-val disk-val">
                  <span v-if="!isEditMode">{{ formatNumber(row.y2.disk) }}</span>
                  <UiNumberInput
                    v-else
                    v-model="row.y2.disk"
                    :show-stepper="false"
                    class="table-number-input"
                  />
                </td>

                <!-- Year 3 Specs -->
                <td class="td-spec y3-val">
                  <span v-if="!isEditMode">{{ formatNumber(row.y3.cpu) }}</span>
                  <UiNumberInput
                    v-else
                    v-model="row.y3.cpu"
                    :show-stepper="false"
                    class="table-number-input"
                  />
                </td>
                <td class="td-spec y3-val">
                  <span v-if="!isEditMode">{{ formatNumber(row.y3.memory) }}</span>
                  <UiNumberInput
                    v-else
                    v-model="row.y3.memory"
                    :show-stepper="false"
                    class="table-number-input"
                  />
                </td>
                <td class="td-spec y3-val disk-val">
                  <span v-if="!isEditMode">{{ formatNumber(row.y3.disk) }}</span>
                  <UiNumberInput
                    v-else
                    v-model="row.y3.disk"
                    :show-stepper="false"
                    class="table-number-input"
                  />
                </td>

                <!-- Action Button in Edit Mode -->
                <td v-if="isEditMode" class="td-action">
                  <button
                    type="button"
                    class="btn-del-row"
                    title="Hapus baris host"
                    @click="deleteHostRow(hIdx)"
                  >
                    🗑️
                  </button>
                </td>
              </tr>
            </tbody>

            <!-- Summary Totals Footer -->
            <tfoot>
              <tr class="total-row">
                <td :colspan="(showGrafanaCols || isGrafanaSynced) ? 7 : 4" class="td-total-label">
                  <strong>Total Alokasi Resource (Semua Node)</strong>
                </td>
                <td class="td-total-val">{{ formatNumber(totals.y1.cpu) }}</td>
                <td class="td-total-val">{{ formatNumber(totals.y1.memory) }} GB</td>
                <td class="td-total-val">{{ formatNumber(totals.y1.disk) }} GB</td>

                <td class="td-total-val">{{ formatNumber(totals.y2.cpu) }}</td>
                <td class="td-total-val">{{ formatNumber(totals.y2.memory) }} GB</td>
                <td class="td-total-val">{{ formatNumber(totals.y2.disk) }} GB</td>

                <td class="td-total-val">{{ formatNumber(totals.y3.cpu) }}</td>
                <td class="td-total-val">{{ formatNumber(totals.y3.memory) }} GB</td>
                <td class="td-total-val">{{ formatNumber(totals.y3.disk) }} GB</td>
                <td v-if="isEditMode"></td>
              </tr>
            </tfoot>
          </table>
        </div>
      </div>

      <!-- Action Footer -->
      <div class="eawg-doc-actions">
        <div class="action-desc">
          <span>📄 Format resmi Bab 10 ini siap disalin ke Dokumen ASDF atau Slide Review EAWG / DRB.</span>
        </div>
        <div class="action-buttons">
          <button type="button" class="btn-reset" @click="resetToDefault">
            🔄 Reset ke Default Preset
          </button>
          <button type="button" class="btn-export-json" @click="downloadEawgJson">
            💾 Unduh JSON Spesifikasi
          </button>
          <button type="button" class="btn-copy-full" @click="copyFullAsdfDoc">
            📋 Salin Seluruh Format Bab 10 (Word / Markdown)
          </button>
        </div>
      </div>
    </div>

    <!-- Modal: Grafana Realtime API Sync Dialog -->
    <div v-if="showGrafanaModal" class="modal-backdrop" @click.self="showGrafanaModal = false">
      <div class="modal-card">
        <div class="modal-header">
          <div class="modal-title-wrap">
            <span class="modal-icon">📡</span>
            <div>
              <h3>Grafana API Real-Time Synchronization</h3>
              <p>Tarik data metrik CPU & RAM riil per host dari Grafana</p>
            </div>
          </div>
          <button type="button" class="close-btn" @click="showGrafanaModal = false">✕</button>
        </div>

        <div class="modal-body">
          <!-- Connection Settings -->
          <div class="form-section">
            <div class="form-row">
              <label>Grafana Endpoint URL</label>
              <input
                v-model="grafanaApiForm.url"
                type="text"
                placeholder="https://grafana.enterprise.co.id"
                class="form-input"
              />
            </div>

            <div class="form-row">
              <label>Service Account / API Bearer Token</label>
              <input
                v-model="grafanaApiForm.token"
                type="password"
                placeholder="glsa_xxxx... (kosongkan untuk Telemetry Simulator)"
                class="form-input"
              />
            </div>
          </div>

          <!-- Time Range Span & Aggregation -->
          <div class="form-section">
            <label class="section-title">Rentang Waktu Analisis (*Time Span Filter*)</label>
            <div class="timespan-grid">
              <button
                v-for="span in timeSpanOptions"
                :key="span.id"
                type="button"
                :class="['span-card-btn', { active: grafanaApiForm.timeSpan === span.id }]"
                @click="grafanaApiForm.timeSpan = span.id"
              >
                <span class="span-label">{{ span.label }}</span>
                <span class="span-desc">{{ span.desc }}</span>
              </button>
            </div>
          </div>

          <div class="form-section">
            <label class="section-title">Metode Agregasi Metrik (*Aggregation Strategy*)</label>
            <div class="aggr-selector">
              <label class="radio-label">
                <input v-model="grafanaApiForm.aggregation" type="radio" value="p95" />
                <span><strong>Percentile 95 (P95 SLA)</strong> — Rekomendasi Standar Enterprise Architecture</span>
              </label>
              <label class="radio-label">
                <input v-model="grafanaApiForm.aggregation" type="radio" value="max" />
                <span><strong>Peak Maximum (Max Spike)</strong> — Beban Puncak Tertinggi</span>
              </label>
              <label class="radio-label">
                <input v-model="grafanaApiForm.aggregation" type="radio" value="avg" />
                <span><strong>Average (Mean Baseline)</strong> — Rata-rata Beban Operasional</span>
              </label>
            </div>
          </div>

          <!-- Target Hosts Summary -->
          <div class="hosts-sync-preview">
            <span>Target Host Sinkronisasi: <strong>{{ activeData.hostSizing.length }} Node</strong> ({{ activeData.hostSizing[0]?.hostName }} .. {{ activeData.hostSizing[activeData.hostSizing.length - 1]?.hostName }})</span>
          </div>
        </div>

        <div class="modal-footer">
          <button type="button" class="btn-cancel" @click="showGrafanaModal = false">
            Batal
          </button>
          <button
            type="button"
            class="btn-sync-submit"
            :disabled="isSyncing"
            @click="executeGrafanaSync"
          >
            <span v-if="isSyncing">⏳ Menghubungkan & Menarik Metrik...</span>
            <span v-else>🚀 Tarik Data Realtime Sekarang</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, reactive } from 'vue'
import { eawgPresets, type EawgProjectPreset } from '~/data/capacityEawgData'
import { formatNumber } from '~/composables/useFormat'
import { useLocalOpen } from '~/composables/useLocalOpen'

const { showToast } = useLocalOpen()

const selectedPresetId = ref<string>('kopra')
const isEditMode = ref<boolean>(false)
const showGrafanaCols = ref<boolean>(false)
const showGrafanaModal = ref<boolean>(false)
const isSyncing = ref<boolean>(false)
const isGrafanaSynced = ref<boolean>(false)

// Time span options
const timeSpanOptions = [
  { id: '1y', label: '1 Tahun (1y)', desc: 'Full Fiscal Year Window (Paling Akurat untuk EAWG)' },
  { id: '6M', label: '6 Bulan (6M)', desc: 'Semi-Annual Peak Season' },
  { id: '90d', label: '3 Bulan (90d)', desc: 'Quarterly Planning Review' },
  { id: '30d', label: '30 Hari (30d)', desc: 'Monthly Payday Peak Period' },
  { id: '7d', label: '7 Hari (7d)', desc: 'Weekly SLA Window' },
]

// Grafana API Modal Form
const grafanaApiForm = reactive({
  url: 'https://grafana.enterprise.co.id',
  token: '',
  timeSpan: '1y' as '1y' | '6M' | '90d' | '30d' | '7d',
  aggregation: 'p95' as 'max' | 'p95' | 'avg',
})

// Last sync metadata
const lastSyncMeta = reactive({
  timeSpan: '1y',
  timeSpanLabel: '1 Tahun Terakhir (1y)',
  aggregation: 'p95',
  syncedAt: '',
})

// Config for Grafana Auto-Extrapolation
const grafanaConfig = reactive({
  safeCpuUtil: 0.40, // 40% target
  safeRamUtil: 0.70, // 70% target
  annualGrowthPercent: 20, // 20% growth
})

// Deep reactive copy of working data
const activeData = reactive<EawgProjectPreset>(JSON.parse(JSON.stringify(eawgPresets[0])))

const lastSyncTimeFormatted = computed(() => {
  if (!lastSyncMeta.syncedAt) return 'Baru saja'
  return new Date(lastSyncMeta.syncedAt).toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit', second: '2-digit' })
})

function selectPreset(id: string) {
  selectedPresetId.value = id
  const target = eawgPresets.find(p => p.id === id)
  if (target) {
    Object.assign(activeData, JSON.parse(JSON.stringify(target)))
    showToast(`📋 Memuat format EAWG ${target.name}`, 'info')
  }
}

function resetToDefault() {
  const target = eawgPresets.find(p => p.id === selectedPresetId.value)
  if (target) {
    Object.assign(activeData, JSON.parse(JSON.stringify(target)))
    isGrafanaSynced.value = false
    showToast(`🔄 Data preset ${target.name} dikembalikan ke default`, 'success')
  }
}

function isNewChannel(idx: number): boolean {
  if (idx === 0) return true
  const list = activeData.businessMetrics
  return list[idx].channel !== list[idx - 1].channel
}

function getUsageIcon(usage: string): string {
  const u = usage.toLowerCase()
  if (u.includes('web')) return '🌐'
  if (u.includes('master')) return '👑'
  if (u.includes('infra') || u.includes('router')) return '🔀'
  if (u.includes('worker')) return '⚙️'
  if (u.includes('redis')) return '⚡'
  if (u.includes('nfs')) return '📁'
  if (u.includes('oracle') || u.includes('postgres') || u.includes('db')) return '🗄️'
  return '🖥️'
}

function getRowCategoryClass(usage: string): string {
  const u = usage.toLowerCase()
  if (u.includes('worker')) return 'row-worker'
  if (u.includes('db') || u.includes('oracle') || u.includes('postgres')) return 'row-database'
  if (u.includes('redis')) return 'row-redis'
  return ''
}

function getUtilStatus(cpu?: number, ram?: number) {
  if (!cpu && !ram) return { label: 'No Data', class: 'status-gray' }
  const c = cpu || 0
  const r = ram || 0

  if (c >= 70 || r >= 85) return { label: '🔴 Critical', class: 'status-red' }
  if (c >= 45 || r >= 70) return { label: '🟡 Warning', class: 'status-yellow' }
  return { label: '🟢 Normal', class: 'status-green' }
}

// Execute Grafana API Sync
async function executeGrafanaSync() {
  isSyncing.value = true
  try {
    const payload = {
      grafanaUrl: grafanaApiForm.url,
      apiToken: grafanaApiForm.token,
      timeSpan: grafanaApiForm.timeSpan,
      aggregation: grafanaApiForm.aggregation,
      hosts: activeData.hostSizing.map(h => ({
        id: h.id,
        hostName: h.hostName,
        usage: h.usage,
        node: h.node,
      })),
      isDemoMode: !grafanaApiForm.token,
    }

    const res = await $fetch<{
      success: boolean
      source: string
      timeSpan: string
      timeSpanLabel: string
      aggregation: string
      syncedAt: string
      results: Record<string, { cpuUtil: number; memUtil: number; status: string }>
    }>('/api/grafana/sync', {
      method: 'POST',
      body: payload,
    })

    if (res && res.results) {
      // Map results back into activeData.hostSizing
      for (const row of activeData.hostSizing) {
        if (res.results[row.hostName]) {
          row.currentCpuUtil = res.results[row.hostName].cpuUtil
          row.currentMemUtil = res.results[row.hostName].memUtil
        }
      }

      lastSyncMeta.timeSpan = res.timeSpan
      lastSyncMeta.timeSpanLabel = res.timeSpanLabel
      lastSyncMeta.aggregation = res.aggregation
      lastSyncMeta.syncedAt = res.syncedAt
      isGrafanaSynced.value = true
      showGrafanaCols.value = true
      showGrafanaModal.value = false

      // Auto scale Year 2 and Year 3
      applyGrafanaExtrapolationAll(false)

      showToast(`📡 Berhasil sinkronisasi ${activeData.hostSizing.length} host dari Grafana API (Span: ${res.timeSpanLabel})!`, 'success')
    }
  } catch (err: any) {
    showToast(`⚠️ Gagal sinkronisasi Grafana: ${err.message}`, 'error')
  } finally {
    isSyncing.value = false
  }
}

// Auto-project business metrics row
function autoProjectRow(idx: number) {
  const row = activeData.businessMetrics[idx]
  const g = (row.growthPercent || 0) / 100
  const base = row.y2025 || 0

  row.y2026 = Math.round(base * (1 + g))
  row.y2027 = Math.round(base * Math.pow(1 + g, 2))
  row.y2028 = Math.round(base * Math.pow(1 + g, 3))
  row.y2029 = Math.round(base * Math.pow(1 + g, 4))
  showToast(`⚡ Proyeksi tahun 2026–2029 berhasil dihitung (${row.growthPercent}% CAGR)!`, 'info')
}

// Apply Grafana sizing extrapolation across all host rows
function applyGrafanaExtrapolationAll(notify = true) {
  const g = (grafanaConfig.annualGrowthPercent || 0) / 100
  const uCpuTarget = grafanaConfig.safeCpuUtil || 0.40
  const uRamTarget = grafanaConfig.safeRamUtil || 0.70

  let updatedCount = 0

  for (const row of activeData.hostSizing) {
    const currentCpuUtil = (row.currentCpuUtil || 40) / 100
    const currentMemUtil = (row.currentMemUtil || 50) / 100

    // Compute required base CPU/RAM at safe utilization threshold
    const baselineRequiredCpu = Math.ceil(row.y1.cpu * (currentCpuUtil / uCpuTarget))
    const baselineRequiredRam = Math.ceil(row.y1.memory * (currentMemUtil / uRamTarget))

    // Project Y2 and Y3 with growth
    row.y2.cpu = Math.max(row.y1.cpu, Math.ceil(baselineRequiredCpu * (1 + g)))
    row.y2.memory = Math.max(row.y1.memory, Math.ceil(baselineRequiredRam * (1 + g)))
    row.y2.disk = Math.ceil(row.y1.disk * (1 + g * 0.5))

    row.y3.cpu = Math.max(row.y2.cpu, Math.ceil(baselineRequiredCpu * Math.pow(1 + g, 2)))
    row.y3.memory = Math.max(row.y2.memory, Math.ceil(baselineRequiredRam * Math.pow(1 + g, 2)))
    row.y3.disk = Math.ceil(row.y1.disk * Math.pow(1 + g * 0.5, 2))

    updatedCount++
  }

  if (notify) {
    showToast(`⚡ Sizing Year 2 & 3 untuk ${updatedCount} host berhasil disesuaikan dengan utilisasi Grafana!`, 'success')
  }
}

// Add / Delete rows
function addBusinessRow() {
  activeData.businessMetrics.push({
    channel: 'Channel Baru',
    metric: 'Jumlah Transaksi',
    y2025: 10000,
    y2026: 12000,
    y2027: 14400,
    y2028: 17280,
    y2029: 20736,
    growthPercent: 20,
  })
  showToast('➕ Baris metrik bisnis ditambahkan', 'info')
}

function deleteBusinessRow(idx: number) {
  activeData.businessMetrics.splice(idx, 1)
  showToast('🗑️ Baris metrik bisnis dihapus', 'info')
}

function addHostRow() {
  const newId = `host-${Date.now()}`
  activeData.hostSizing.push({
    id: newId,
    envi: 'DC',
    usage: 'App Server Baru',
    hostName: 'PVAPP01WBD25LP',
    node: 1,
    currentCpuUtil: 35,
    currentMemUtil: 50,
    y1: { cpu: 8, memory: 16, disk: 100 },
    y2: { cpu: 12, memory: 24, disk: 120 },
    y3: { cpu: 16, memory: 32, disk: 150 },
  })
  showToast('➕ Baris Host / VM ditambahkan', 'info')
}

function deleteHostRow(idx: number) {
  activeData.hostSizing.splice(idx, 1)
  showToast('🗑️ Baris host dihapus', 'info')
}

// Totals calculation
const totals = computed(() => {
  const t = {
    y1: { cpu: 0, memory: 0, disk: 0 },
    y2: { cpu: 0, memory: 0, disk: 0 },
    y3: { cpu: 0, memory: 0, disk: 0 },
  }

  for (const row of activeData.hostSizing) {
    t.y1.cpu += row.y1?.cpu || 0
    t.y1.memory += row.y1?.memory || 0
    t.y1.disk += row.y1?.disk || 0

    t.y2.cpu += row.y2?.cpu || 0
    t.y2.memory += row.y2?.memory || 0
    t.y2.disk += row.y2?.disk || 0

    t.y3.cpu += row.y3?.cpu || 0
    t.y3.memory += row.y3?.memory || 0
    t.y3.disk += row.y3?.disk || 0
  }

  return t
})

function copyBusinessTable() {
  let text = `| Channel | Metrics | 2025 | 2026 | 2027 | 2028 | 2029 | Growth |\n`
  text += `|---|---|---|---|---|---|---|---|\n`
  for (const r of activeData.businessMetrics) {
    text += `| ${r.channel} | ${r.metric} | ${formatNumber(r.y2025)} | ${formatNumber(r.y2026)} | ${formatNumber(r.y2027)} | ${formatNumber(r.y2028)} | ${formatNumber(r.y2029)} | ${r.growthPercent}% |\n`
  }

  navigator.clipboard.writeText(text)
  showToast('✅ Tabel 10.1 Proyeksi Bisnis berhasil disalin ke clipboard!', 'success')
}

function copyHostTable() {
  let text = `| Envi | Penggunaan | Host Name | Node | Y1 CPU | Y1 RAM (GB) | Y1 Disk (GB) | Y2 CPU | Y2 RAM (GB) | Y2 Disk (GB) | Y3 CPU | Y3 RAM (GB) | Y3 Disk (GB) |\n`
  text += `|---|---|---|---|---|---|---|---|---|---|---|---|---|\n`
  for (const r of activeData.hostSizing) {
    text += `| ${r.envi} | ${r.usage} | ${r.hostName} | ${r.node} | ${formatNumber(r.y1.cpu)} | ${formatNumber(r.y1.memory)} | ${formatNumber(r.y1.disk)} | ${formatNumber(r.y2.cpu)} | ${formatNumber(r.y2.memory)} | ${formatNumber(r.y2.disk)} | ${formatNumber(r.y3.cpu)} | ${formatNumber(r.y3.memory)} | ${formatNumber(r.y3.disk)} |\n`
  }

  navigator.clipboard.writeText(text)
  showToast('✅ Tabel 10.2 Sizing Host berhasil disalin ke clipboard!', 'success')
}

function copyFullAsdfDoc() {
  const p = activeData
  let doc = `10. Business Capacity Plan\n`
  doc += `Detail capacity plan berdasarkan proyeksi bisnis tertera pada EAWG, Asumsi bisnis terkait ${p.projectName} adalah sebagai berikut:\n\n`

  doc += `Tabel 10.1: Asumsi Proyeksi Bisnis 5 Tahun\n`
  doc += `| Channel | Metrics | 2025 | 2026 | 2027 | 2028 | 2029 | Growth |\n`
  doc += `|---|---|---|---|---|---|---|---|\n`
  for (const r of p.businessMetrics) {
    doc += `| ${r.channel} | ${r.metric} | ${formatNumber(r.y2025)} | ${formatNumber(r.y2026)} | ${formatNumber(r.y2027)} | ${formatNumber(r.y2028)} | ${formatNumber(r.y2029)} | ${r.growthPercent}% |\n`
  }

  doc += `\nCapacity plan terkait ${p.name} adalah sebagai berikut:\n\n`
  doc += `Tabel 10.2: Capacity Plan Sizing BOM Multi-Year\n`
  doc += `| Envi | Penggunaan | Host Name | Node | Y1 CPU | Y1 RAM (GB) | Y1 Disk (GB) | Y2 CPU | Y2 RAM (GB) | Y2 Disk (GB) | Y3 CPU | Y3 RAM (GB) | Y3 Disk (GB) |\n`
  doc += `|---|---|---|---|---|---|---|---|---|---|---|---|---|\n`
  for (const r of p.hostSizing) {
    doc += `| ${r.envi} | ${r.usage} | ${r.hostName} | ${r.node} | ${formatNumber(r.y1.cpu)} | ${formatNumber(r.y1.memory)} | ${formatNumber(r.y1.disk)} | ${formatNumber(r.y2.cpu)} | ${formatNumber(r.y2.memory)} | ${formatNumber(r.y2.disk)} | ${formatNumber(r.y3.cpu)} | ${formatNumber(r.y3.memory)} | ${formatNumber(r.y3.disk)} |\n`
  }

  doc += `| **TOTAL** | - | - | - | **${formatNumber(totals.value.y1.cpu)}** | **${formatNumber(totals.value.y1.memory)}** | **${formatNumber(totals.value.y1.disk)}** | **${formatNumber(totals.value.y2.cpu)}** | **${formatNumber(totals.value.y2.memory)}** | **${formatNumber(totals.value.y2.disk)}** | **${formatNumber(totals.value.y3.cpu)}** | **${formatNumber(totals.value.y3.memory)}** | **${formatNumber(totals.value.y3.disk)}** |\n`

  navigator.clipboard.writeText(doc)
  showToast('📋 Format Lengkap Bab 10 ASDF berhasil disalin ke clipboard!', 'success')
}

function downloadEawgJson() {
  const payload = {
    title: 'Enterprise Architecture EAWG Business Capacity Plan (Bab 10)',
    timestamp: new Date().toISOString(),
    preset: activeData,
    grafanaSync: isGrafanaSynced.value ? lastSyncMeta : null,
    grafanaConfig,
    totals: totals.value,
  }

  const dataStr = 'data:text/json;charset=utf-8,' + encodeURIComponent(JSON.stringify(payload, null, 2))
  const downloadAnchor = document.createElement('a')
  downloadAnchor.setAttribute('href', dataStr)
  downloadAnchor.setAttribute('download', `eawg_capacity_plan_${activeData.id}_${new Date().toISOString().slice(0, 10)}.json`)
  document.body.appendChild(downloadAnchor)
  downloadAnchor.click()
  downloadAnchor.remove()

  showToast('💾 Berkas JSON Spesifikasi EAWG berhasil diunduh!', 'success')
}
</script>

<style scoped>
.capacity-eawg-view {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.eawg-toolbar {
  background: linear-gradient(135deg, #0f2b48 0%, #003b73 100%);
  color: #ffffff;
  padding: 16px 24px;
  border-radius: var(--radius-lg);
  display: flex;
  align-items: center;
  justify-content: space-between;
  box-shadow: var(--shadow-md);
  flex-wrap: wrap;
  gap: 16px;
}

.toolbar-left {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.doc-badge {
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.5px;
  text-transform: uppercase;
  color: #93c5fd;
}

.doc-title-mini {
  font-size: 16px;
  font-weight: 800;
  color: #ffffff;
}

.toolbar-controls {
  display: flex;
  align-items: center;
  gap: 16px;
  flex-wrap: wrap;
}

.preset-selector {
  display: flex;
  align-items: center;
  gap: 8px;
}

.preset-label {
  font-size: 12px;
  font-weight: 600;
  color: #cbd5e1;
}

.preset-pills {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-wrap: wrap;
}

.preset-pill {
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

.preset-pill:hover {
  background: rgba(255, 255, 255, 0.25);
}

.preset-pill.active {
  background: #ffffff;
  color: var(--navy);
  border-color: #ffffff;
  font-weight: 700;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
}

.mode-toggles {
  display: flex;
  align-items: center;
  gap: 8px;
}

.toggle-btn {
  font-family: 'Inter', sans-serif;
  font-size: 11px;
  font-weight: 700;
  padding: 6px 14px;
  border-radius: var(--radius-sm);
  border: 1px solid rgba(255, 255, 255, 0.3);
  background: rgba(255, 255, 255, 0.12);
  color: #ffffff;
  cursor: pointer;
  transition: all 0.15s;
}

.toggle-btn:hover {
  background: rgba(255, 255, 255, 0.25);
}

.toggle-btn.active {
  background: #f59e0b;
  color: #0f172a;
  border-color: #f59e0b;
  box-shadow: 0 2px 8px rgba(245, 158, 11, 0.3);
}

.toggle-btn.grafana-btn {
  background: rgba(59, 130, 246, 0.25);
  border-color: #60a5fa;
  color: #ffffff;
}

.toggle-btn.grafana-btn.active {
  background: #3b82f6;
  color: #ffffff;
  border-color: #93c5fd;
  box-shadow: 0 2px 8px rgba(59, 130, 246, 0.35);
}

/* Live Banner */
.grafana-live-banner {
  background: linear-gradient(90deg, #1e3a8a 0%, #1e40af 100%);
  color: #ffffff;
  padding: 10px 18px;
  border-radius: var(--radius);
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-size: 12px;
  border: 1px solid #3b82f6;
  box-shadow: 0 2px 6px rgba(30, 58, 138, 0.25);
  animation: fadeIn 0.2s ease-in-out;
  flex-wrap: wrap;
  gap: 8px;
}

.banner-left {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
}

.pulse-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #22c55e;
  box-shadow: 0 0 0 3px rgba(34, 197, 94, 0.3);
  animation: pulse 1.5s infinite;
}

.banner-title {
  font-weight: 700;
  color: #93c5fd;
}

.banner-meta {
  color: #e0e7ff;
}

.banner-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.banner-time {
  font-size: 11px;
  color: #93c5fd;
}

.btn-resync {
  font-family: 'Inter', sans-serif;
  font-size: 10px;
  font-weight: 700;
  padding: 3px 8px;
  background: rgba(255, 255, 255, 0.2);
  color: #ffffff;
  border: 1px solid rgba(255, 255, 255, 0.3);
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.15s;
}

.btn-resync:hover {
  background: #ffffff;
  color: #1e3a8a;
}

/* Grafana Assistant Card */
.grafana-assistant-card {
  background: #eff6ff;
  border: 1px solid #bfdbfe;
  border-radius: var(--radius-lg);
  padding: 18px 24px;
  box-shadow: var(--shadow-sm);
  animation: fadeIn 0.2s ease-in-out;
}

.grafana-header {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.grafana-title-wrap {
  display: flex;
  align-items: center;
  gap: 10px;
}

.grafana-icon {
  font-size: 24px;
}

.grafana-title-wrap h4 {
  font-size: 14px;
  font-weight: 700;
  color: #1e3a8a;
}

.grafana-title-wrap p {
  font-size: 12px;
  color: #3b82f6;
}

.grafana-config-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  gap: 12px;
  align-items: flex-end;
  background: #ffffff;
  padding: 12px 16px;
  border-radius: var(--radius);
  border: 1px solid #dbeafe;
}

.grafana-field {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.grafana-field label {
  font-size: 11px;
  font-weight: 600;
  color: var(--text2);
}

.grafana-field select,
.grafana-field input {
  font-family: 'Inter', sans-serif;
  font-size: 12px;
  padding: 6px 10px;
  border: 1px solid var(--border2);
  border-radius: var(--radius-sm);
  background: var(--bg2);
  outline: none;
}

.input-with-suffix {
  display: flex;
  align-items: center;
  background: var(--bg2);
  border: 1px solid var(--border2);
  border-radius: var(--radius-sm);
  padding-right: 8px;
}

.input-with-suffix input {
  border: none;
  background: transparent;
  width: 100%;
}

.input-with-suffix span {
  font-size: 11px;
  color: var(--text3);
  font-weight: 600;
}

.grafana-apply-btn {
  font-family: 'Inter', sans-serif;
  font-size: 12px;
  font-weight: 700;
  padding: 8px 16px;
  background: #2563eb;
  color: #ffffff;
  border: none;
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: background 0.15s;
  width: 100%;
}

.grafana-apply-btn:hover {
  background: #1d4ed8;
}

/* Paper Container */
.eawg-paper {
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius-lg);
  padding: 32px 36px;
  box-shadow: var(--shadow-sm);
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.doc-section-header {
  border-bottom: 2px solid var(--brand-blue, #003b73);
  padding-bottom: 12px;
}

.section-title-wrap {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 8px;
  flex-wrap: wrap;
}

.section-num {
  font-size: 20px;
  font-weight: 800;
  color: var(--navy);
}

.edit-badge {
  font-size: 11px;
  font-weight: 700;
  color: #d97706;
  background: #fef3c7;
  padding: 3px 8px;
  border-radius: 4px;
  border: 1px solid #fde68a;
}

.section-intro {
  font-size: 13px;
  color: var(--text);
  line-height: 1.6;
}

.project-title-input {
  font-size: 13px;
  font-weight: 700;
  color: var(--navy);
  border: 1px dashed var(--blue);
  background: #f0f9ff;
  padding: 2px 8px;
  border-radius: 4px;
  width: 300px;
}

.doc-sub-intro {
  font-size: 13px;
  color: var(--text);
  margin-top: 4px;
}

.table-container {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.table-title-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 8px;
}

.table-caption {
  font-size: 13px;
  font-weight: 700;
  color: var(--navy);
}

.table-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.mini-btn {
  font-family: 'Inter', sans-serif;
  font-size: 11px;
  font-weight: 600;
  padding: 4px 10px;
  border-radius: var(--radius-sm);
  border: 1px solid var(--border2);
  background: var(--bg3);
  color: var(--text);
  cursor: pointer;
  transition: all 0.15s;
}

.mini-btn:hover {
  background: var(--brand-blue);
  color: #ffffff;
  border-color: var(--brand-blue);
}

.add-btn {
  background: #f0fdf4;
  color: #16a34a;
  border-color: #bbf7d0;
  font-weight: 700;
}

.add-btn:hover {
  background: #16a34a;
  color: #ffffff;
}

.table-responsive {
  overflow-x: auto;
  border: 1px solid var(--border);
  border-radius: var(--radius);
  background: #ffffff;
}

.eawg-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 12px;
  text-align: left;
}

.eawg-table th,
.eawg-table td {
  padding: 8px 10px;
  border: 1px solid var(--border);
}

.eawg-table thead th {
  background: #f1f5f9;
  color: #1e293b;
  font-weight: 700;
  font-size: 11px;
  text-align: center;
}

.business-table th {
  background: #e2e8f0;
}

.business-table .th-channel,
.business-table .th-metrics {
  text-align: left;
}

.td-channel {
  background: #f8fafc;
  color: var(--navy);
}

.td-metric {
  color: var(--text);
}

.td-num {
  text-align: right;
  font-family: 'JetBrains Mono', monospace;
  font-weight: 600;
  color: #0f172a;
}

.td-growth {
  text-align: center;
}

.growth-badge {
  font-size: 11px;
  font-weight: 700;
  background: #f0fdf4;
  color: #16a34a;
  border: 1px solid #bbf7d0;
  padding: 2px 6px;
  border-radius: 4px;
}

.growth-edit-wrap {
  display: inline-flex;
  align-items: center;
  gap: 2px;
}

.growth-input {
  width: 48px;
  font-size: 11px;
  font-weight: 700;
  padding: 2px 4px;
  text-align: right;
  border: 1px solid var(--border2);
  border-radius: 3px;
}

.btn-calc-growth {
  background: #eff6ff;
  border: 1px solid #bfdbfe;
  border-radius: 3px;
  cursor: pointer;
  font-size: 10px;
  padding: 2px 4px;
}

.btn-calc-growth:hover {
  background: #2563eb;
  color: #ffffff;
}

.td-action {
  text-align: center;
  width: 40px;
}

.btn-del-row {
  background: transparent;
  border: none;
  cursor: pointer;
  font-size: 13px;
  opacity: 0.6;
  transition: opacity 0.15s;
}

.btn-del-row:hover {
  opacity: 1;
  transform: scale(1.1);
}

/* Host Table Styles */
.host-table thead tr:first-child th {
  background: #cbd5e1;
}

.th-year {
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.th-grafana {
  background: #dbeafe !important;
  color: #1e40af !important;
}

.y1-header {
  background: #e0e7ff !important;
  color: #3730a3 !important;
}

.y2-header {
  background: #f3e8ff !important;
  color: #6b21a8 !important;
}

.y3-header {
  background: #fae8ff !important;
  color: #86198f !important;
}

.sub-header-row th {
  font-size: 10px;
  padding: 4px 6px;
  background: #f8fafc;
}

.th-graf-cpu,
.th-graf-ram,
.th-graf-status {
  background: #eff6ff !important;
  color: #1e40af;
}

.y1-col {
  background: #eef2ff !important;
  color: #3730a3;
}

.y2-col {
  background: #faf5ff !important;
  color: #6b21a8;
}

.y3-col {
  background: #fdf4ff !important;
  color: #86198f;
}

.td-envi {
  text-align: center;
}

.envi-tag {
  font-size: 10px;
  font-weight: 800;
  padding: 2px 6px;
  border-radius: 4px;
}

.envi-tag.dc {
  background: #dbeafe;
  color: #1e40af;
  border: 1px solid #bfdbfe;
}

.envi-tag.drc {
  background: #fef3c7;
  color: #b45309;
  border: 1px solid #fde68a;
}

.envi-select {
  font-size: 10px;
  font-weight: 800;
  padding: 2px 4px;
  border-radius: 3px;
  border: 1px solid var(--border2);
}

.td-usage {
  display: flex;
  align-items: center;
  gap: 6px;
  font-weight: 600;
  color: var(--text);
}

.usage-icon {
  font-size: 13px;
}

.td-hostname code {
  font-family: 'JetBrains Mono', monospace;
  font-size: 11px;
  background: #f1f5f9;
  padding: 2px 6px;
  border-radius: 4px;
  border: 1px solid #e2e8f0;
  color: #0f172a;
}

.td-node-val {
  text-align: center;
  font-weight: 700;
  color: var(--text2);
}

.node-input {
  width: 38px;
  text-align: center;
  font-size: 11px;
  padding: 2px;
  border: 1px solid var(--border2);
  border-radius: 3px;
}

.td-spec {
  text-align: right;
  font-family: 'JetBrains Mono', monospace;
  font-size: 11px;
}

.y1-val {
  background: rgba(238, 242, 255, 0.4);
}

.y2-val {
  background: rgba(250, 245, 255, 0.4);
}

.y3-val {
  background: rgba(253, 244, 255, 0.4);
}

.disk-val {
  font-weight: 600;
  color: #0f2b48;
}

/* Grafana cells */
.td-grafana-cell {
  background: rgba(239, 246, 255, 0.5);
  text-align: center;
}

.grafana-input-wrap {
  display: inline-flex;
  align-items: center;
  gap: 2px;
}

.graf-input {
  width: 44px;
  font-family: 'JetBrains Mono', monospace;
  font-size: 11px;
  font-weight: 700;
  text-align: right;
  padding: 2px 4px;
  border: 1px solid #bfdbfe;
  background: #ffffff;
  border-radius: 3px;
}

.graf-input:disabled {
  background: transparent;
  border-color: transparent;
}

.td-grafana-status {
  text-align: center;
}

.status-pill {
  font-size: 10px;
  font-weight: 700;
  padding: 2px 6px;
  border-radius: 4px;
  white-space: nowrap;
}

.status-green {
  background: #f0fdf4;
  color: #16a34a;
  border: 1px solid #bbf7d0;
}

.status-yellow {
  background: #fffbeb;
  color: #d97706;
  border: 1px solid #fde68a;
}

.status-red {
  background: #fef2f2;
  color: #dc2626;
  border: 1px solid #fecaca;
}

.status-gray {
  background: #f1f5f9;
  color: #64748b;
  border: 1px solid #e2e8f0;
}

/* Inline Inputs in Table */
.inline-text-input {
  font-family: 'Inter', sans-serif;
  font-size: 12px;
  padding: 2px 6px;
  border: 1px solid var(--border2);
  border-radius: 3px;
  background: #ffffff;
  outline: none;
  width: 100%;
}

.inline-text-input:focus {
  border-color: var(--blue);
}

.channel-input {
  font-weight: 700;
  color: var(--navy);
  width: 110px;
}

.metric-input {
  width: 150px;
}

.usage-input {
  width: 140px;
}

.host-input {
  font-family: 'JetBrains Mono', monospace;
  font-size: 11px;
  width: 130px;
}

.table-number-input {
  min-width: 60px;
}

:deep(.table-number-input .formatted-number-input) {
  padding: 2px 4px;
  font-size: 11px;
  text-align: right;
  font-family: 'JetBrains Mono', monospace;
}

/* Category highlights */
.row-worker {
  background: #f0fdf4;
}

.row-worker td.td-usage {
  color: #15803d;
}

.row-database {
  background: #fffbeb;
}

.row-database td.td-usage {
  color: #b45309;
}

.row-redis {
  background: #fdf2f8;
}

.row-redis td.td-usage {
  color: #be185d;
}

/* Totals */
.total-row {
  background: #0f2b48 !important;
  color: #ffffff !important;
  font-weight: 700;
}

.total-row td {
  border-color: #1e3a5f;
}

.td-total-label {
  text-align: right;
  padding-right: 14px;
  color: #e2e8f0;
}

.td-total-val {
  text-align: right;
  font-family: 'JetBrains Mono', monospace;
  color: #ffffff;
  font-size: 11px;
}

/* Footer Actions */
.eawg-doc-actions {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-top: 16px;
  border-top: 1px solid var(--border);
  flex-wrap: wrap;
  gap: 12px;
}

.action-desc {
  font-size: 12px;
  color: var(--text3);
}

.action-buttons {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.btn-reset {
  font-family: 'Inter', sans-serif;
  font-size: 12px;
  font-weight: 600;
  padding: 8px 14px;
  background: var(--bg3);
  color: var(--text2);
  border: 1px solid var(--border2);
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: all 0.15s;
}

.btn-reset:hover {
  background: var(--border2);
}

.btn-copy-full {
  font-family: 'Inter', sans-serif;
  font-size: 12px;
  font-weight: 700;
  padding: 8px 16px;
  background: var(--brand-blue, #003b73);
  color: #ffffff;
  border: none;
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: background 0.15s;
}

.btn-copy-full:hover {
  background: var(--navy);
}

.btn-export-json {
  font-family: 'Inter', sans-serif;
  font-size: 12px;
  font-weight: 600;
  padding: 8px 14px;
  background: var(--bg3);
  color: var(--text);
  border: 1px solid var(--border2);
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: all 0.15s;
}

.btn-export-json:hover {
  background: var(--border2);
}

/* Modal Backdrop */
.modal-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(15, 23, 42, 0.65);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 20px;
  animation: fadeIn 0.15s ease-out;
}

.modal-card {
  background: #ffffff;
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-lg);
  width: 100%;
  max-width: 640px;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  border: 1px solid var(--border);
}

.modal-header {
  background: linear-gradient(135deg, #0f2b48 0%, #003b73 100%);
  color: #ffffff;
  padding: 20px 24px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.modal-title-wrap {
  display: flex;
  align-items: center;
  gap: 12px;
}

.modal-icon {
  font-size: 24px;
}

.modal-title-wrap h3 {
  font-size: 16px;
  font-weight: 700;
  color: #ffffff;
}

.modal-title-wrap p {
  font-size: 12px;
  color: #cbd5e1;
}

.close-btn {
  background: transparent;
  border: none;
  color: #ffffff;
  font-size: 18px;
  cursor: pointer;
  opacity: 0.8;
  transition: opacity 0.15s;
}

.close-btn:hover {
  opacity: 1;
}

.modal-body {
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 20px;
  max-height: 70vh;
  overflow-y: auto;
}

.form-section {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.section-title {
  font-size: 12px;
  font-weight: 700;
  color: var(--navy);
  text-transform: uppercase;
  letter-spacing: 0.3px;
}

.form-row {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.form-row label {
  font-size: 11px;
  font-weight: 600;
  color: var(--text2);
}

.form-input {
  font-family: 'Inter', sans-serif;
  font-size: 13px;
  padding: 8px 12px;
  border: 1px solid var(--border2);
  border-radius: var(--radius-sm);
  outline: none;
}

.form-input:focus {
  border-color: var(--blue);
}

.timespan-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(160px, 1fr));
  gap: 8px;
}

.span-card-btn {
  font-family: 'Inter', sans-serif;
  padding: 10px 12px;
  border-radius: var(--radius-sm);
  border: 1px solid var(--border);
  background: var(--bg2);
  cursor: pointer;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 3px;
  text-align: left;
  transition: all 0.15s;
}

.span-card-btn:hover {
  border-color: var(--blue);
  background: #eff6ff;
}

.span-card-btn.active {
  border-color: #2563eb;
  background: #dbeafe;
  box-shadow: 0 0 0 2px rgba(37, 99, 235, 0.2);
}

.span-label {
  font-size: 12px;
  font-weight: 700;
  color: var(--navy);
}

.span-desc {
  font-size: 10px;
  color: var(--text3);
}

.aggr-selector {
  display: flex;
  flex-direction: column;
  gap: 8px;
  background: var(--bg);
  padding: 12px 16px;
  border-radius: var(--radius-sm);
  border: 1px solid var(--border);
}

.radio-label {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
  color: var(--text);
  cursor: pointer;
}

.hosts-sync-preview {
  font-size: 11px;
  color: var(--text2);
  background: #f8fafc;
  padding: 8px 12px;
  border-radius: 4px;
  border: 1px solid #e2e8f0;
}

.modal-footer {
  padding: 16px 24px;
  background: var(--bg3);
  border-top: 1px solid var(--border);
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 10px;
}

.btn-cancel {
  font-family: 'Inter', sans-serif;
  font-size: 12px;
  font-weight: 600;
  padding: 8px 16px;
  border: 1px solid var(--border2);
  background: #ffffff;
  border-radius: var(--radius-sm);
  cursor: pointer;
}

.btn-sync-submit {
  font-family: 'Inter', sans-serif;
  font-size: 12px;
  font-weight: 700;
  padding: 8px 20px;
  background: #2563eb;
  color: #ffffff;
  border: none;
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: background 0.15s;
}

.btn-sync-submit:hover:not(:disabled) {
  background: #1d4ed8;
}

.btn-sync-submit:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

@keyframes pulse {
  0% {
    box-shadow: 0 0 0 0 rgba(34, 197, 94, 0.6);
  }
  70% {
    box-shadow: 0 0 0 8px rgba(34, 197, 94, 0);
  }
  100% {
    box-shadow: 0 0 0 0 rgba(34, 197, 94, 0);
  }
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(-6px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@media (max-width: 800px) {
  .eawg-paper {
    padding: 20px;
  }
  .eawg-toolbar {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>
