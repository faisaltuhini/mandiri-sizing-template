<template>
  <div class="boq-estimator-view">
    <!-- Header Banner -->
    <div class="boq-header">
      <div>
        <div class="boq-badge">Standar Pengadaan & Anggaran IT Enterprise Architecture</div>
        <h2 class="boq-title">Estimasi Bill of Quantity (BOQ) & Anggaran Infrastruktur Tahunan</h2>
        <p class="boq-subtitle">
          Kalkulasi otomatis biaya pengadaan vCPU, SAN Storage, NetBackup, Server Bare Metal, dan Lisensi Enterprise
          berdasarkan standar tarif resmi Enterprise Architecture (Template Sizing.xlsx).
        </p>
      </div>

      <div class="boq-tier-toggle">
        <button
          type="button"
          :class="['tier-toggle-btn', { active: activeTier === 'moderate' }]"
          @click="changeTier('moderate')"
        >
          Tier Moderate (IHSP & RMT)
        </button>
        <button
          type="button"
          :class="['tier-toggle-btn', { active: activeTier === 'critical' }]"
          @click="changeTier('critical')"
        >
          Tier Critical (Livin' Merchant)
        </button>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="boq-loading">
      Memuat kalkulasi anggaran resmi dari Go Engine...
    </div>

    <div v-else-if="boqData" class="boq-content">
      <!-- Executive Cost Cards -->
      <div class="cost-summary-grid">
        <div class="cost-card card-total">
          <span class="cost-card-lbl">Grand Total Biaya Tahunan</span>
          <span class="cost-card-val">{{ formatRupiah(boqData.grand_total_annual) }}</span>
          <span class="cost-card-sub">Termasuk Production, DRC, & Non-Prod</span>
        </div>

        <div class="cost-card">
          <span class="cost-card-lbl">Biaya Hardware & Compute (vCPU + SAN)</span>
          <span class="cost-card-val">{{ formatRupiah(boqData.annual_hardware_cost) }}</span>
          <span class="cost-card-sub">{{ boqData.total_vcpu }} vCPU | {{ formatNumber(boqData.total_storage_gb) }} GB Storage</span>
        </div>

        <div class="cost-card">
          <span class="cost-card-lbl">Biaya Proteksi Backup Netbackup</span>
          <span class="cost-card-val">{{ formatRupiah(boqData.annual_backup_cost) }}</span>
          <span class="cost-card-sub">Formula: 50% data aktif × 80% dedup × Rp 27.000</span>
        </div>

        <div class="cost-card">
          <span class="cost-card-lbl">Pengadaan Server Bare Metal</span>
          <span class="cost-card-val">{{ formatRupiah(boqData.annual_bare_metal_cost) }}</span>
          <span class="cost-card-sub">{{ activeTier === 'critical' ? '4 Physical Server BM Oracle RAC & OCP' : '0 BM (100% Virtual Machine)' }}</span>
        </div>
      </div>

      <!-- Action Toolbar -->
      <div class="boq-tools">
        <div class="boq-stats-pills">
          <span class="stat-pill">Total Node: <strong>{{ boqData.total_vm }} Server/VM</strong></span>
          <span class="stat-pill">Total Core: <strong>{{ boqData.total_vcpu }} vCPU</strong></span>
          <span class="stat-pill">Total Memory: <strong>{{ formatNumber(boqData.total_ram_gb) }} GB RAM</strong></span>
          <span class="stat-pill">Total Storage: <strong>{{ formatNumber(boqData.total_storage_gb) }} GB SAN</strong></span>
        </div>

        <div class="boq-action-btns">
          <a
            :href="`${apiBase}/api/sizing/export/excel`"
            class="export-excel-btn"
            target="_blank"
          >
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="margin-right: 6px;">
              <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
              <polyline points="7 10 12 15 17 10"/>
              <line x1="12" y1="15" x2="12" y2="3"/>
            </svg>
            Unduh Excel BOQ Resmi (.xlsx)
          </a>
        </div>
      </div>

      <!-- Detailed BOQ Table -->
      <div class="boq-table-wrap">
        <table class="boq-table">
          <thead>
            <tr>
              <th class="th-center">No</th>
              <th>Nama Komponen / Host</th>
              <th class="th-center">Env</th>
              <th class="th-center">Jml Node</th>
              <th class="th-center">Total vCPU</th>
              <th class="th-center">Total RAM</th>
              <th class="th-center">Total Storage</th>
              <th class="th-right">Biaya vCPU / Thn</th>
              <th class="th-right">Biaya Storage / Thn</th>
              <th class="th-right">Biaya Backup / Thn</th>
              <th class="th-right">Total Biaya / Thn</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="item in boqData.items"
              :key="item.no"
              class="boq-row"
            >
              <td class="td-center font-bold">{{ item.no }}</td>
              <td class="font-bold text-navy">{{ item.item_name }}</td>
              <td class="td-center">
                <span :class="['env-badge', `env-${item.env.toLowerCase()}`]">{{ item.env }}</span>
              </td>
              <td class="td-center font-mono">{{ item.qty }}</td>
              <td class="td-center font-mono">{{ item.vcpu }}</td>
              <td class="td-center font-mono">{{ item.ram_gb }} GB</td>
              <td class="td-center font-mono">{{ formatNumber(item.storage_gb) }} GB</td>
              <td class="td-right font-mono">{{ formatRupiah(item.cost_vcpu) }}</td>
              <td class="td-right font-mono">{{ formatRupiah(item.cost_stg) }}</td>
              <td class="td-right font-mono">{{ formatRupiah(item.cost_nbu) }}</td>
              <td class="td-right font-bold font-mono text-dark">{{ formatRupiah(item.total_cost) }}</td>
            </tr>
          </tbody>
          <tfoot>
            <tr class="boq-total-row">
              <td colspan="3" class="font-bold">TOTAL RINGKASAN BIAYA ARSITEKTUR</td>
              <td class="td-center font-bold">{{ boqData.total_vm }}</td>
              <td class="td-center font-bold">{{ boqData.total_vcpu }}</td>
              <td class="td-center font-bold">{{ formatNumber(boqData.total_ram_gb) }} GB</td>
              <td class="td-center font-bold">{{ formatNumber(boqData.total_storage_gb) }} GB</td>
              <td colspan="3" class="td-right font-bold">GRAND TOTAL ANNUAL COST:</td>
              <td class="td-right font-bold text-gold text-lg">{{ formatRupiah(boqData.grand_total_annual) }}</td>
            </tr>
          </tfoot>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useSizingBackend } from '~/composables/useSizingBackend'
import { formatNumber } from '~/composables/useFormat'

const { apiBase, fetchBOQ } = useSizingBackend()

const activeTier = ref<'moderate' | 'critical'>('moderate')
const boqData = ref<any>(null)
const loading = ref(false)

function formatRupiah(val: number): string {
  if (!val && val !== 0) return 'Rp 0'
  return 'Rp ' + Math.round(val).toLocaleString('id-ID')
}

async function loadData() {
  loading.value = true
  try {
    const res = await fetchBOQ(activeTier.value)
    if (res) {
      boqData.value = res
    }
  } finally {
    loading.value = false
  }
}

function changeTier(tier: 'moderate' | 'critical') {
  activeTier.value = tier
  loadData()
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.boq-estimator-view {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.boq-header {
  background: linear-gradient(135deg, #002060 0%, #001233 100%);
  color: #ffffff;
  padding: 24px 28px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 20px;
}

.boq-badge {
  display: inline-block;
  background: rgba(255, 192, 0, 0.2);
  color: #ffc000;
  border: 1px solid rgba(255, 192, 0, 0.4);
  font-size: 11px;
  font-weight: 700;
  text-transform: uppercase;
  padding: 3px 10px;
  border-radius: 20px;
  margin-bottom: 8px;
}

.boq-title {
  font-size: 20px;
  font-weight: 800;
  margin: 0 0 6px 0;
  color: #ffffff;
}

.boq-subtitle {
  font-size: 13px;
  color: #cbd5e1;
  margin: 0;
  max-width: 720px;
  line-height: 1.5;
}

.boq-tier-toggle {
  display: flex;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 8px;
  padding: 4px;
  gap: 4px;
}

.tier-toggle-btn {
  border: none;
  background: transparent;
  color: #e2e8f0;
  font-size: 12px;
  font-weight: 700;
  padding: 8px 14px;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s;
  white-space: nowrap;
}

.tier-toggle-btn.active {
  background: #ffc000;
  color: #002060;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
}

.boq-loading {
  padding: 40px;
  text-align: center;
  font-size: 14px;
  color: #64748b;
  background: #ffffff;
  border-radius: 10px;
}

.cost-summary-grid {
  display: grid;
  grid-template-columns: 1.3fr 1fr 1fr 1fr;
  gap: 16px;
}

.cost-card {
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 10px;
  padding: 16px 20px;
  display: flex;
  flex-direction: column;
  gap: 4px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

.cost-card.card-total {
  background: linear-gradient(135deg, #002060 0%, #003399 100%);
  color: #ffffff;
  border-color: #ffc000;
}

.cost-card-lbl {
  font-size: 11px;
  font-weight: 700;
  text-transform: uppercase;
  color: #64748b;
}

.card-total .cost-card-lbl {
  color: #ffc000;
}

.cost-card-val {
  font-size: 20px;
  font-weight: 800;
  letter-spacing: -0.5px;
  color: #0f172a;
}

.card-total .cost-card-val {
  color: #ffffff;
  font-size: 22px;
}

.cost-card-sub {
  font-size: 11px;
  color: #64748b;
}

.card-total .cost-card-sub {
  color: #cbd5e1;
}

.boq-tools {
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 10px;
  padding: 12px 18px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.boq-stats-pills {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.stat-pill {
  background: #f1f5f9;
  border: 1px solid #e2e8f0;
  padding: 4px 10px;
  border-radius: 6px;
  font-size: 11.5px;
  color: #334155;
}

.export-excel-btn {
  background: #107c41;
  color: #ffffff;
  font-size: 12px;
  font-weight: 700;
  padding: 8px 16px;
  border-radius: 6px;
  text-decoration: none;
  display: inline-flex;
  align-items: center;
  gap: 6px;
  transition: all 0.2s;
}

.export-excel-btn:hover {
  background: #0b5e30;
}

.boq-table-wrap {
  background: #ffffff;
  border: 1px solid #cbd5e1;
  border-radius: 10px;
  overflow-x: auto;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.05);
}

.boq-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 12px;
}

.boq-table th {
  background: #1b365d;
  color: #ffffff;
  font-weight: 700;
  padding: 10px 12px;
  border: 1px solid #cbd5e1;
  text-align: left;
}

.boq-table td {
  padding: 8px 12px;
  border: 1px solid #e2e8f0;
  vertical-align: middle;
}

.boq-row:nth-child(even) {
  background-color: #f8fafc;
}

.boq-row:hover {
  background-color: #eff6ff;
}

.th-center, .td-center {
  text-align: center;
}

.th-right, .td-right {
  text-align: right;
}

.font-bold {
  font-weight: 700;
}

.font-mono {
  font-family: ui-monospace, monospace;
}

.text-navy {
  color: #002060;
}

.text-dark {
  color: #0f172a;
}

.text-gold {
  color: #b45309;
}

.text-lg {
  font-size: 14px;
}

.env-badge {
  font-size: 10px;
  font-weight: 700;
  padding: 2px 6px;
  border-radius: 4px;
  text-transform: uppercase;
}

.env-prod {
  background: #dbeafe;
  color: #1e40af;
}

.env-drc {
  background: #fef3c7;
  color: #92400e;
}

.env-dev {
  background: #f1f5f9;
  color: #475569;
}

.boq-total-row {
  background: #f8fafc;
  border-top: 2px solid #002060;
}
</style>
