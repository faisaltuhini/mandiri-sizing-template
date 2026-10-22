<template>
  <div class="tech-stack-baseline-view">
    <!-- Top Executive Banner (Mandiri Deep Navy & Brushed Gold) -->
    <div class="baseline-header">
      <div class="header-content">
        <div class="header-badges">
          <span class="badge-mandiri">Bank Mandiri EA Forum Standard</span>
          <span class="badge-accent">Arsitektur Baseline Resmi</span>
          <span class="badge-count">13 Layer Arsitektur</span>
        </div>
        <h1 class="header-title">Standar Teknologi & Sizing Infrastruktur</h1>
        <p class="header-subtitle">
          Pedoman resmi penentuan spesifikasi server, lisensi, dan kapasitas komputasi Bank Mandiri untuk
          <strong>Tier Moderate (IHSP & RMT)</strong> vs <strong>Tier Critical (Livin' Merchant & Super App)</strong>.
        </p>
      </div>

      <!-- Quick Actions Right -->
      <div class="header-actions">
        <a
          href="/api/templates/download?path=template/TECH_STACK/Template_Tech_Stack_Architecture_Sizing.xlsx"
          class="action-btn-primary"
          title="Download berkas Excel resmi 6-sheet"
        >
          <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
            <polyline points="7 10 12 15 17 10"/>
            <line x1="12" y1="15" x2="12" y2="3"/>
          </svg>
          Unduh Template Excel Resmi
        </a>
        <button
          type="button"
          class="action-btn-secondary"
          @click="openInExcel"
          title="Buka langsung di aplikasi Microsoft Excel"
        >
          <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"/>
            <polyline points="15 3 21 3 21 9"/>
            <line x1="10" y1="14" x2="21" y2="3"/>
          </svg>
          Buka di Excel
        </button>
      </div>
    </div>

    <!-- Executive Comparison Snapshot (Human-Friendly for Non-Tech Savvy) -->
    <div class="executive-snapshot-grid">
      <!-- Tier Moderate Profile Card -->
      <div class="profile-card profile-moderate">
        <div class="profile-header">
          <div class="profile-badge-mod">TIER MODERATE (IHSP & RMT)</div>
          <span class="profile-target">Untuk Aplikasi Internal & Transaksi Menengah</span>
        </div>
        <div class="profile-body">
          <div class="profile-metric">
            <div class="metric-block">
              <span class="metric-val">14 Node</span>
              <span class="metric-lbl">Total Virtual Machine (VM)</span>
            </div>
            <div class="metric-divider"></div>
            <div class="metric-block">
              <span class="metric-val">112 Core</span>
              <span class="metric-lbl">vCPU & 256 GB RAM</span>
            </div>
            <div class="metric-divider"></div>
            <div class="metric-block">
              <span class="metric-val">4.5 TB</span>
              <span class="metric-lbl">SAN Tier 2 Storage</span>
            </div>
          </div>
          <div class="profile-features">
            <div class="feature-item">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="#2563eb" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"/></svg>
              <span>Efisiensi anggaran optimal dengan infrastruktur Shared VM & SAN LUN Tier 2</span>
            </div>
            <div class="feature-item">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="#2563eb" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"/></svg>
              <span>Arsitektur HA Dual-Node (2N) pada layer Web Ingress dan OpenShift Worker</span>
            </div>
            <div class="feature-item">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="#2563eb" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"/></svg>
              <span>Database Single-Node VM dengan replikasi Data Guard untuk DRC</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Tier Critical Profile Card -->
      <div class="profile-card profile-critical">
        <div class="profile-header">
          <div class="profile-badge-crit">TIER CRITICAL (LIVIN' MERCHANT & SUPER APP)</div>
          <span class="profile-target">Customer-Facing 24/7, Transaksi Tinggi & Nol Downtime</span>
        </div>
        <div class="profile-body">
          <div class="profile-metric">
            <div class="metric-block">
              <span class="metric-val text-amber">28 Node</span>
              <span class="metric-lbl">Dedicated Servers (incl. 4 Bare Metal)</span>
            </div>
            <div class="metric-divider"></div>
            <div class="metric-block">
              <span class="metric-val text-amber">384 Core</span>
              <span class="metric-lbl">vCPU & 4.8 TB RAM</span>
            </div>
            <div class="metric-divider"></div>
            <div class="metric-block">
              <span class="metric-val text-amber">24.5 TB</span>
              <span class="metric-lbl">High-Perf ASM SAN Storage</span>
            </div>
          </div>
          <div class="profile-features">
            <div class="feature-item">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="#d97706" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"/></svg>
              <span>Dedicated Physical Bare-Metal Servers untuk Database Oracle RAC 2-Node ASM</span>
            </div>
            <div class="feature-item">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="#d97706" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"/></svg>
              <span>Dedicated Bare-Metal OpenShift Workers (2x 64 vCPU) untuk penanganan spike transaksi masif</span>
            </div>
            <div class="feature-item">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="#d97706" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"/></svg>
              <span>High-speed Kafka 3-Node Cluster & Redis Multi-Shard 3-Node untuk caching real-time</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Controls Bar: Filter, View Mode, Search -->
    <div class="controls-bar">
      <!-- Search Box -->
      <div class="search-box">
        <svg class="search-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="#94a3b8" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <circle cx="11" cy="11" r="8"/>
          <line x1="21" y1="21" x2="16.65" y2="16.65"/>
        </svg>
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Cari layer, komponen, atau spesifikasi (misal: NGINX, Oracle RAC, Kafka, Redis)..."
          class="search-input"
        />
        <button
          v-if="searchQuery"
          type="button"
          class="clear-btn"
          @click="searchQuery = ''"
        >
          ✕
        </button>
      </div>

      <!-- View Presentation Switcher: Cards vs Table -->
      <div class="view-presentation-toggle">
        <button
          type="button"
          :class="['toggle-tab', { active: displayMode === 'cards' }]"
          @click="displayMode = 'cards'"
        >
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <rect x="3" y="3" width="7" height="7"/>
            <rect x="14" y="3" width="7" height="7"/>
            <rect x="14" y="14" width="7" height="7"/>
            <rect x="3" y="14" width="7" height="7"/>
          </svg>
          Ringkasan Kartu (Mudah Dipahami)
        </button>
        <button
          type="button"
          :class="['toggle-tab', { active: displayMode === 'table' }]"
          @click="displayMode = 'table'"
        >
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <line x1="3" y1="6" x2="21" y2="6"/>
            <line x1="3" y1="12" x2="21" y2="12"/>
            <line x1="3" y1="18" x2="21" y2="18"/>
          </svg>
          Tabel Matriks Lengkap (13 Layer)
        </button>
      </div>

      <!-- Tier Filter Buttons -->
      <div class="tier-filter-group">
        <button
          type="button"
          :class="['filter-btn', { active: tierFilter === 'all' }]"
          @click="tierFilter = 'all'"
        >
          Semua Tier
        </button>
        <button
          type="button"
          :class="['filter-btn', { active: tierFilter === 'moderate' }]"
          @click="tierFilter = 'moderate'"
        >
          Moderate Saja
        </button>
        <button
          type="button"
          :class="['filter-btn', { active: tierFilter === 'critical' }]"
          @click="tierFilter = 'critical'"
        >
          Critical Saja
        </button>
      </div>

      <!-- Copy Action -->
      <button
        type="button"
        class="copy-btn"
        @click="copyTableMarkdown"
        title="Salin ringkasan ke clipboard format Markdown"
      >
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <rect x="9" y="9" width="13" height="13" rx="2" ry="2"/>
          <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/>
        </svg>
        {{ copyStatus || 'Salin Ringkasan' }}
      </button>
    </div>

    <!-- PRESENTATION 1: EXECUTIVE CARDS VIEW (Clean & Non-Tech Friendly) -->
    <div v-if="displayMode === 'cards'" class="executive-cards-grid">
      <div
        v-for="item in filteredItems"
        :key="item.id"
        class="layer-card"
        @click="selectItem(item)"
      >
        <div class="layer-card-top">
          <div class="layer-index">Layer {{ item.no }}</div>
          <div class="layer-category">{{ item.layer }}</div>
        </div>

        <h3 class="layer-component-title">{{ item.component }}</h3>

        <div class="layer-specs-comparison">
          <!-- Moderate Box -->
          <div v-if="tierFilter === 'all' || tierFilter === 'moderate'" class="spec-box box-moderate">
            <span class="box-lbl">Tier Moderate</span>
            <span class="box-nodes">{{ item.moderate.nodes }} · {{ item.moderate.type }}</span>
            <span class="box-detail">{{ item.moderate.spec }}</span>
          </div>

          <!-- Critical Box -->
          <div v-if="tierFilter === 'all' || tierFilter === 'critical'" class="spec-box box-critical">
            <span class="box-lbl">Tier Critical</span>
            <span class="box-nodes">{{ item.critical.nodes }} · {{ item.critical.type }}</span>
            <span class="box-detail">{{ item.critical.spec }}</span>
          </div>
        </div>

        <!-- License & Vendor -->
        <div class="layer-license-strip">
          <span class="lic-tag">Lisensi: {{ item.license.enterprise }}</span>
        </div>

        <!-- Notes -->
        <p class="layer-notes">{{ item.notes }}</p>

        <div class="card-action-cue">
          <span>Klik untuk detail arsitektur lengkap →</span>
        </div>
      </div>
    </div>

    <!-- PRESENTATION 2: CLEAN ARCHITECTURE MATRIX TABLE -->
    <div v-else class="table-container">
      <table class="baseline-table">
        <thead>
          <tr class="header-row-primary">
            <th class="col-no">No</th>
            <th class="col-layer">Layer Arsitektur</th>
            <th class="col-comp">Komponen Utama</th>
            <th v-if="tierFilter === 'all' || tierFilter === 'moderate'" class="col-moderate">
              Standar Tier Moderate (IHSP & RMT)
            </th>
            <th v-if="tierFilter === 'all' || tierFilter === 'critical'" class="col-critical">
              Standar Tier Critical (Livin' Merchant)
            </th>
            <th class="col-lic">Lisensi Vendor</th>
            <th class="col-notes">Catatan & Justifikasi Sizing</th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="item in filteredItems"
            :key="item.id"
            class="data-row"
            @click="selectItem(item)"
          >
            <td class="cell-no">{{ item.no }}</td>
            <td class="cell-layer">
              <span class="layer-pill">{{ item.layer }}</span>
            </td>
            <td class="cell-comp">
              <strong>{{ item.component }}</strong>
            </td>

            <!-- Moderate Spec Cell -->
            <td v-if="tierFilter === 'all' || tierFilter === 'moderate'" class="cell-moderate">
              <div class="cell-tier-wrap">
                <span class="badge-soft-blue">{{ item.moderate.nodes }}</span>
                <span class="badge-soft-gray">{{ item.moderate.type }}</span>
                <span class="spec-text">{{ item.moderate.spec }}</span>
              </div>
            </td>

            <!-- Critical Spec Cell -->
            <td v-if="tierFilter === 'all' || tierFilter === 'critical'" class="cell-critical">
              <div class="cell-tier-wrap">
                <span class="badge-soft-amber">{{ item.critical.nodes }}</span>
                <span class="badge-soft-amber-light">{{ item.critical.type }}</span>
                <span class="spec-text font-bold">{{ item.critical.spec }}</span>
              </div>
            </td>

            <!-- License Cell -->
            <td class="cell-lic">
              <div class="lic-wrap">
                <span class="badge-soft-green">{{ item.license.enterprise }}</span>
              </div>
            </td>

            <!-- Notes Cell -->
            <td class="cell-notes">
              <div class="notes-content">
                {{ item.notes }}
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Polished Executive Modal Detail -->
    <div
      v-if="selectedItemData"
      class="modal-backdrop"
      @click.self="selectedItemData = null"
    >
      <div class="modal-card">
        <div class="modal-header">
          <div>
            <div class="modal-tag">Layer {{ selectedItemData.no }} · {{ selectedItemData.layer }}</div>
            <h2 class="modal-title">{{ selectedItemData.component }}</h2>
          </div>
          <button
            type="button"
            class="modal-close"
            @click="selectedItemData = null"
          >
            ✕
          </button>
        </div>

        <div class="modal-body">
          <div class="modal-compare-grid">
            <!-- Moderate Section -->
            <div class="tier-card tier-card-moderate">
              <div class="tier-card-hdr">
                <span>Tier Moderate (IHSP & RMT)</span>
              </div>
              <div class="tier-card-content">
                <div class="meta-row">
                  <span class="meta-label">Jumlah Node:</span>
                  <span class="meta-val">{{ selectedItemData.moderate.nodes }}</span>
                </div>
                <div class="meta-row">
                  <span class="meta-label">Tipe Mesin:</span>
                  <span class="meta-val">{{ selectedItemData.moderate.type }}</span>
                </div>
                <div class="meta-row">
                  <span class="meta-label">Spesifikasi Server:</span>
                  <span class="meta-val highlight-blue">{{ selectedItemData.moderate.spec }}</span>
                </div>
              </div>
            </div>

            <!-- Critical Section -->
            <div class="tier-card tier-card-critical">
              <div class="tier-card-hdr">
                <span>Tier Critical (Livin' Merchant & Super App)</span>
              </div>
              <div class="tier-card-content">
                <div class="meta-row">
                  <span class="meta-label">Jumlah Node:</span>
                  <span class="meta-val">{{ selectedItemData.critical.nodes }}</span>
                </div>
                <div class="meta-row">
                  <span class="meta-label">Tipe Mesin:</span>
                  <span class="meta-val">{{ selectedItemData.critical.type }}</span>
                </div>
                <div class="meta-row">
                  <span class="meta-label">Spesifikasi Server:</span>
                  <span class="meta-val highlight-amber">{{ selectedItemData.critical.spec }}</span>
                </div>
              </div>
            </div>
          </div>

          <!-- License & Justification Details -->
          <div class="detail-section">
            <h4 class="section-title">Kategori Lisensi & Penyedia Vendor</h4>
            <div class="lic-compare">
              <div class="lic-col">
                <span class="lic-sub-lbl">Komunitas / Free:</span>
                <span class="lic-badge lic-free">{{ selectedItemData.license.free }}</span>
              </div>
              <div class="lic-col">
                <span class="lic-sub-lbl">Lisensi Resmi Bank Mandiri:</span>
                <span class="lic-badge lic-ent">{{ selectedItemData.license.enterprise }}</span>
              </div>
            </div>
          </div>

          <div class="detail-section">
            <h4 class="section-title">Catatan Arsitektur & Rekomendasi Sizing</h4>
            <p class="justification-text">{{ selectedItemData.notes }}</p>
          </div>
        </div>

        <div class="modal-footer">
          <button
            type="button"
            class="modal-btn"
            @click="selectedItemData = null"
          >
            Tutup
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import {
  BASELINE_TECH_STACK_ITEMS,
  type TechStackItem
} from '~/data/techStackBaselineData'
import { useLocalOpen } from '~/composables/useLocalOpen'

const { openPath, showToast } = useLocalOpen()

const searchQuery = ref('')
const tierFilter = ref<'all' | 'moderate' | 'critical'>('all')
const displayMode = ref<'cards' | 'table'>('cards')
const selectedItemData = ref<TechStackItem | null>(null)
const copyStatus = ref('')

const filteredItems = computed(() => {
  const q = searchQuery.value.trim().toLowerCase()
  if (!q) return BASELINE_TECH_STACK_ITEMS

  return BASELINE_TECH_STACK_ITEMS.filter((item) => {
    return (
      item.component.toLowerCase().includes(q) ||
      item.layer.toLowerCase().includes(q) ||
      item.moderate.spec.toLowerCase().includes(q) ||
      item.critical.spec.toLowerCase().includes(q) ||
      item.license.free.toLowerCase().includes(q) ||
      item.license.enterprise.toLowerCase().includes(q) ||
      item.notes.toLowerCase().includes(q) ||
      (item.contractRef && item.contractRef.toLowerCase().includes(q))
    )
  })
})

function selectItem(item: TechStackItem) {
  selectedItemData.value = item
}

function openInExcel() {
  openPath('template/TECH_STACK/Template_Tech_Stack_Architecture_Sizing.xlsx', 'finder')
  showToast('Membuka Template_Tech_Stack_Architecture_Sizing.xlsx di Excel...')
}

async function copyTableMarkdown() {
  const headers = [
    'No',
    'Layer',
    'Komponen',
    'Moderate Nodes',
    'Moderate Type',
    'Moderate Spec',
    'Critical Nodes',
    'Critical Type',
    'Critical Spec',
    'Enterprise License',
    'Catatan Arsitektur'
  ]
  const rows = filteredItems.value.map((i) => [
    i.no,
    i.layer,
    i.component,
    i.moderate.nodes,
    i.moderate.type,
    i.moderate.spec,
    i.critical.nodes,
    i.critical.type,
    i.critical.spec,
    i.license.enterprise,
    i.notes
  ])
  const md = [
    `| ${headers.join(' | ')} |`,
    `| ${headers.map(() => '---').join(' | ')} |`,
    ...rows.map((r) => `| ${r.join(' | ')} |`)
  ].join('\n')

  try {
    await navigator.clipboard.writeText(md)
    copyStatus.value = 'Tersalin!'
    setTimeout(() => {
      copyStatus.value = ''
    }, 2500)
    showToast('Tabel berhasil disalin ke clipboard dalam format Markdown', 'success')
  } catch (err) {
    showToast('Gagal menyalin tabel ke clipboard', 'error')
  }
}
</script>

<style scoped>
.tech-stack-baseline-view {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

/* ─────────────────────────────────────────────
   Executive Corporate Header
   ───────────────────────────────────────────── */
.baseline-header {
  background: linear-gradient(135deg, #002060 0%, #00143d 100%);
  color: #ffffff;
  padding: 28px 36px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-left: 5px solid #d4a017;
  box-shadow: 0 4px 20px rgba(0, 32, 96, 0.15);
}

.header-badges {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 10px;
}

.badge-mandiri {
  font-size: 11px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  background: rgba(212, 160, 23, 0.2);
  color: #fde047;
  border: 1px solid rgba(212, 160, 23, 0.4);
  padding: 3px 10px;
  border-radius: 6px;
}

.badge-accent {
  font-size: 11px;
  font-weight: 600;
  background: rgba(255, 255, 255, 0.12);
  color: #ffffff;
  padding: 3px 10px;
  border-radius: 6px;
}

.badge-count {
  font-size: 11px;
  font-weight: 600;
  background: rgba(37, 99, 235, 0.3);
  color: #93c5fd;
  padding: 3px 10px;
  border-radius: 6px;
}

.header-title {
  font-size: 22px;
  font-weight: 800;
  letter-spacing: -0.01em;
  color: #ffffff;
  margin-bottom: 6px;
}

.header-subtitle {
  font-size: 13px;
  color: #cbd5e1;
  max-width: 800px;
  line-height: 1.5;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.action-btn-primary {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  background: #d4a017;
  color: #002060;
  font-size: 12px;
  font-weight: 700;
  padding: 10px 18px;
  border-radius: 8px;
  text-decoration: none;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
  transition: all 0.2s;
}

.action-btn-primary:hover {
  background: #e6b325;
  transform: translateY(-1px);
}

.action-btn-secondary {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  background: rgba(255, 255, 255, 0.12);
  color: #ffffff;
  border: 1px solid rgba(255, 255, 255, 0.25);
  font-size: 12px;
  font-weight: 600;
  padding: 10px 16px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
}

.action-btn-secondary:hover {
  background: rgba(255, 255, 255, 0.2);
}

/* ─────────────────────────────────────────────
   Executive Comparison Snapshot
   ───────────────────────────────────────────── */
.executive-snapshot-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
}

.profile-card {
  background: #ffffff;
  border-radius: 12px;
  border: 1px solid #e2e8f0;
  padding: 24px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
  display: flex;
  flex-direction: column;
  gap: 16px;
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.profile-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.08);
}

.profile-moderate {
  border-top: 4px solid #002060;
}

.profile-critical {
  border-top: 4px solid #d4a017;
}

.profile-header {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.profile-badge-mod {
  font-size: 12px;
  font-weight: 800;
  color: #002060;
  letter-spacing: 0.5px;
}

.profile-badge-crit {
  font-size: 12px;
  font-weight: 800;
  color: #b45309;
  letter-spacing: 0.5px;
}

.profile-target {
  font-size: 12px;
  color: #64748b;
}

.profile-metric {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: #f8fafc;
  padding: 16px;
  border-radius: 8px;
  border: 1px solid #f1f5f9;
}

.metric-block {
  display: flex;
  flex-direction: column;
}

.metric-val {
  font-size: 18px;
  font-weight: 800;
  color: #002060;
}

.metric-val.text-amber {
  color: #b45309;
}

.metric-lbl {
  font-size: 11px;
  color: #64748b;
  margin-top: 2px;
}

.metric-divider {
  width: 1px;
  height: 32px;
  background: #e2e8f0;
}

.profile-features {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.feature-item {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  font-size: 12px;
  color: #334155;
  line-height: 1.4;
}

/* ─────────────────────────────────────────────
   Controls Bar
   ───────────────────────────────────────────── */
.controls-bar {
  display: flex;
  align-items: center;
  gap: 12px;
  background: #ffffff;
  padding: 14px 20px;
  border-radius: 10px;
  border: 1px solid #e2e8f0;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.03);
  flex-wrap: wrap;
}

.search-box {
  flex: 1;
  min-width: 260px;
  position: relative;
  display: flex;
  align-items: center;
}

.search-icon {
  position: absolute;
  left: 12px;
}

.search-input {
  width: 100%;
  padding: 9px 36px 9px 36px;
  font-size: 13px;
  border: 1px solid #cbd5e1;
  border-radius: 8px;
  background: #f8fafc;
  color: #0f172a;
  outline: none;
  transition: all 0.2s;
}

.search-input:focus {
  border-color: #002060;
  background: #ffffff;
  box-shadow: 0 0 0 3px rgba(0, 32, 96, 0.08);
}

.clear-btn {
  position: absolute;
  right: 10px;
  background: none;
  border: none;
  color: #94a3b8;
  cursor: pointer;
  font-size: 12px;
}

.view-presentation-toggle {
  display: flex;
  background: #f1f5f9;
  padding: 3px;
  border-radius: 8px;
}

.toggle-tab {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  font-weight: 600;
  color: #64748b;
  padding: 6px 12px;
  border: none;
  background: transparent;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.15s ease;
}

.toggle-tab.active {
  background: #ffffff;
  color: #002060;
  font-weight: 700;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.08);
}

.tier-filter-group {
  display: flex;
  gap: 4px;
}

.filter-btn {
  font-size: 12px;
  font-weight: 600;
  color: #475569;
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  padding: 6px 12px;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.15s;
}

.filter-btn.active {
  background: #002060;
  color: #ffffff;
  border-color: #002060;
}

.copy-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  font-weight: 600;
  color: #475569;
  background: #ffffff;
  border: 1px solid #cbd5e1;
  padding: 6px 12px;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.15s;
}

.copy-btn:hover {
  background: #f8fafc;
  color: #002060;
}

/* ─────────────────────────────────────────────
   Executive Cards Grid Presentation
   ───────────────────────────────────────────── */
.executive-cards-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(380px, 1fr));
  gap: 18px;
}

.layer-card {
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.03);
  cursor: pointer;
  display: flex;
  flex-direction: column;
  gap: 12px;
  transition: all 0.2s ease;
}

.layer-card:hover {
  border-color: #002060;
  transform: translateY(-2px);
  box-shadow: 0 8px 18px rgba(0, 32, 96, 0.08);
}

.layer-card-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.layer-index {
  font-size: 11px;
  font-weight: 700;
  color: #94a3b8;
  text-transform: uppercase;
}

.layer-category {
  font-size: 11px;
  font-weight: 700;
  background: #f1f5f9;
  color: #475569;
  padding: 2px 8px;
  border-radius: 4px;
}

.layer-component-title {
  font-size: 16px;
  font-weight: 700;
  color: #0f172a;
}

.layer-specs-comparison {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.spec-box {
  padding: 10px 12px;
  border-radius: 8px;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.box-moderate {
  background: #f0f7ff;
  border-left: 3px solid #2563eb;
}

.box-critical {
  background: #fefce8;
  border-left: 3px solid #d97706;
}

.box-lbl {
  font-size: 10px;
  font-weight: 700;
  text-transform: uppercase;
  color: #64748b;
}

.box-nodes {
  font-size: 12px;
  font-weight: 700;
  color: #1e293b;
}

.box-detail {
  font-size: 11px;
  color: #475569;
  font-family: 'JetBrains Mono', monospace;
}

.layer-license-strip {
  display: flex;
}

.lic-tag {
  font-size: 11px;
  font-weight: 600;
  background: #ecfdf5;
  color: #065f46;
  border: 1px solid #a7f3d0;
  padding: 3px 8px;
  border-radius: 6px;
}

.layer-notes {
  font-size: 12px;
  color: #64748b;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.card-action-cue {
  margin-top: auto;
  font-size: 11px;
  font-weight: 600;
  color: #002060;
  text-align: right;
}

/* ─────────────────────────────────────────────
   Clean Table Presentation
   ───────────────────────────────────────────── */
.table-container {
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.03);
}

.baseline-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 12px;
  text-align: left;
}

.header-row-primary th {
  background: #002060;
  color: #ffffff;
  padding: 14px 16px;
  font-weight: 700;
  font-size: 12px;
  border-bottom: 2px solid #d4a017;
  letter-spacing: 0.3px;
}

.col-no { width: 48px; text-align: center; }
.col-layer { width: 140px; }
.col-comp { width: 180px; }
.col-moderate { width: 240px; }
.col-critical { width: 260px; }
.col-lic { width: 160px; }
.col-notes { min-width: 220px; }

.data-row {
  border-bottom: 1px solid #e2e8f0;
  cursor: pointer;
  transition: background 0.15s ease;
}

.data-row:hover {
  background: #f8fafc;
}

.data-row td {
  padding: 12px 16px;
  vertical-align: middle;
}

.cell-no {
  text-align: center;
  font-weight: 700;
  color: #94a3b8;
}

.layer-pill {
  font-size: 11px;
  font-weight: 600;
  background: #f1f5f9;
  color: #334155;
  padding: 3px 8px;
  border-radius: 4px;
}

.cell-tier-wrap {
  display: flex;
  flex-direction: column;
  gap: 3px;
}

.badge-soft-blue {
  display: inline-block;
  font-size: 11px;
  font-weight: 700;
  color: #1e40af;
}

.badge-soft-gray {
  font-size: 10px;
  color: #64748b;
}

.badge-soft-amber {
  display: inline-block;
  font-size: 11px;
  font-weight: 700;
  color: #b45309;
}

.badge-soft-amber-light {
  font-size: 10px;
  color: #92400e;
}

.badge-soft-green {
  display: inline-block;
  font-size: 11px;
  font-weight: 600;
  background: #ecfdf5;
  color: #065f46;
  border: 1px solid #a7f3d0;
  padding: 2px 6px;
  border-radius: 4px;
}

.spec-text {
  font-size: 11px;
  color: #334155;
  font-family: 'JetBrains Mono', monospace;
}

.notes-content {
  font-size: 11px;
  color: #64748b;
  line-height: 1.4;
}

/* ─────────────────────────────────────────────
   Modal Detail
   ───────────────────────────────────────────── */
.modal-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(15, 23, 42, 0.6);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 200;
  padding: 24px;
}

.modal-card {
  background: #ffffff;
  border-radius: 14px;
  width: 100%;
  max-width: 720px;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.25);
  display: flex;
  flex-direction: column;
}

.modal-header {
  padding: 20px 24px;
  background: #002060;
  color: #ffffff;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.modal-tag {
  font-size: 11px;
  font-weight: 700;
  color: #d4a017;
  text-transform: uppercase;
}

.modal-title {
  font-size: 18px;
  font-weight: 800;
  color: #ffffff;
  margin-top: 2px;
}

.modal-close {
  background: rgba(255, 255, 255, 0.1);
  border: none;
  color: #ffffff;
  font-size: 14px;
  width: 32px;
  height: 32px;
  border-radius: 8px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}

.modal-body {
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.modal-compare-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.tier-card {
  border-radius: 10px;
  padding: 16px;
  border: 1px solid #e2e8f0;
}

.tier-card-moderate {
  background: #f0f7ff;
  border-color: #bfdbfe;
}

.tier-card-critical {
  background: #fefce8;
  border-color: #fde68a;
}

.tier-card-hdr {
  font-size: 12px;
  font-weight: 800;
  margin-bottom: 12px;
}

.tier-card-moderate .tier-card-hdr { color: #1e40af; }
.tier-card-critical .tier-card-hdr { color: #b45309; }

.tier-card-content {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.meta-row {
  display: flex;
  flex-direction: column;
  font-size: 11px;
}

.meta-label {
  color: #64748b;
  font-weight: 600;
}

.meta-val {
  color: #0f172a;
  font-weight: 700;
}

.highlight-blue { color: #1d4ed8; font-family: 'JetBrains Mono', monospace; }
.highlight-amber { color: #b45309; font-family: 'JetBrains Mono', monospace; }

.detail-section {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.section-title {
  font-size: 13px;
  font-weight: 700;
  color: #002060;
  border-bottom: 1px solid #f1f5f9;
  padding-bottom: 4px;
}

.lic-compare {
  display: flex;
  gap: 16px;
}

.lic-col {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.lic-sub-lbl {
  font-size: 11px;
  color: #64748b;
}

.lic-badge {
  display: inline-block;
  padding: 6px 10px;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 600;
}

.lic-free { background: #f1f5f9; color: #475569; }
.lic-ent { background: #ecfdf5; color: #065f46; border: 1px solid #a7f3d0; }

.justification-text {
  font-size: 12px;
  color: #334155;
  line-height: 1.6;
  background: #f8fafc;
  padding: 12px;
  border-radius: 8px;
  border: 1px solid #f1f5f9;
}

.modal-footer {
  padding: 16px 24px;
  background: #f8fafc;
  border-top: 1px solid #e2e8f0;
  display: flex;
  justify-content: flex-end;
}

.modal-btn {
  background: #002060;
  color: #ffffff;
  font-size: 12px;
  font-weight: 600;
  padding: 8px 18px;
  border-radius: 6px;
  border: none;
  cursor: pointer;
}
</style>
