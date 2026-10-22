<template>
  <div class="sketch-workspace-grid">
    <!-- LEFT COLUMN: BULLET LIST (• ________) -->
    <section class="sketch-list-section">
      <div class="section-hdr">
        <h2 class="section-title">Daftar Komponen &amp; Layer</h2>
        <span class="item-count">{{ layerItems.length }} Layers</span>
      </div>

      <div class="bullet-list-container">
        <div
          v-for="item in layerItems"
          :key="item.id"
          :class="['bullet-row', { active: activeItemId === item.id }]"
          @click="$emit('select-layer', item.id)"
        >
          <span class="bullet-dot">•</span>
          <div class="bullet-content">
            <div class="bullet-title-row">
              <span class="bullet-label">{{ item.layer }}</span>
              <span v-if="isLayerCustomized(item.id)" class="custom-badge" title="Telah disesuaikan dari baseline">Kustom</span>
            </div>
            <span class="bullet-sub">{{ item.component }}</span>
          </div>
          <span class="bullet-chevron">→</span>
        </div>
      </div>
    </section>

    <!-- RIGHT COLUMN: SPECIFICATION, FRIENDLY STEP SECTIONS & COST PREVIEW -->
    <section class="sketch-card-section">
      <div class="pure-content-card">
        <!-- Header Banner of the Component -->
        <div class="component-header-banner">
          <div class="component-header-info">
            <div class="badge-row">
              <span class="badge-layer">Layer {{ activeItem.no }} · {{ activeItem.layer }}</span>
              <span v-if="isLayerCustomized(activeItem.id)" class="custom-status-tag">
                ⚠️ Ada Penyesuaian Kustom
              </span>
            </div>
            <h2 class="component-title">{{ activeItem.component }}</h2>
          </div>

          <div class="banner-actions-right">
            <div class="filter-status-tags">
              <span :class="['active-filter-tag', tierFilter]">
                Tier: {{ tierFilter.toUpperCase() }}
              </span>
              <span :class="['active-filter-tag', licenseFilter === 'free' ? 'free' : 'license']">
                Lisensi: {{ licenseFilter === 'free' ? 'FREE (OSS)' : 'LICENSE (ENTERPRISE)' }}
              </span>
            </div>
            <button
              type="button"
              class="btn-reset-baseline"
              @click="$emit('reset-active-layer')"
              title="Kembalikan konfigurasi layer ini ke standar resmi Bank Mandiri"
            >
              <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M3 12a9 9 0 1 0 9-9 9.75 9.75 0 0 0-6.74 2.74L3 8"/>
                <path d="M3 3v5h5"/>
              </svg>
              <span>Reset ke Standar</span>
            </button>
          </div>
        </div>

        <!-- ======================================================== -->
        <!-- STEP 1: UKURAN SERVER & KOMPUTASI                        -->
        <!-- ======================================================== -->
        <BaseSectionCard
          step-number="1"
          title="Kekuatan Server &amp; Komputasi"
          guide-question="Berapa banyak server dan seberapa tangguh otak pemroses yang dibutuhkan?"
          technical-label="vCPU &amp; RAM"
          :has-changes="isComputeCustomized"
        >
          <!-- Size Preset Selector (S / M / L / XL / Kustom) -->
          <div class="form-row-block">
            <div class="field-hdr-wrap">
              <label class="field-title-lbl">Ukuran Standar Server per Unit:</label>
              <BaseHelpTip
                title="Ukuran Server (CPU &amp; Memori)"
                subtitle="vCPU &amp; RAM"
                :content="glossary.vcpu.shortDesc"
                :business-impact="glossary.vcpu.businessImpact"
              />
            </div>
            <BaseSizePicker
              :vcpu="activeConfig.vcpuPerNode"
              :ram="activeConfig.ramGBPerNode"
              :recommended-vcpu="recommendedVcpu"
              :recommended-ram="recommendedRam"
              @update:vcpu="activeConfig.vcpuPerNode = $event"
              @update:ram="activeConfig.ramGBPerNode = $event"
            />
          </div>

          <!-- Stepper: Jumlah Server -->
          <div class="server-count-row">
            <div class="field-hdr-wrap">
              <div class="title-with-desc">
                <span class="field-title-lbl">Jumlah Server Aktif:</span>
                <span class="field-sub-hint">Berapa unit server mandiri yang dialokasikan</span>
              </div>
              <BaseHelpTip
                title="Jumlah Server (Nodes)"
                subtitle="High Availability Redundancy"
                :content="glossary.nodeCount.shortDesc"
                :business-impact="glossary.nodeCount.businessImpact"
              />
            </div>

            <div class="stepper-wrap-contained">
              <BaseStepper
                v-model="activeConfig.nodeCount"
                :min="1"
                :max="32"
                label="Server Aktif"
                unit="Server"
              />
            </div>
          </div>

          <!-- Advanced Toggle: Granular Steppers -->
          <BaseAdvancedToggle
            show-label="Sesuaikan angka vCPU dan RAM manual (Lanjutan)"
            hide-label="Sembunyikan pengaturan angka manual"
          >
            <div class="advanced-steppers-grid">
              <BaseStepper
                v-model="activeConfig.vcpuPerNode"
                :min="activeConfig.computeType === 'shared' ? 0 : 1"
                :max="128"
                :step="2"
                label="vCPU Pemroses"
                unit="Cores"
                inline-unit="vCPU"
                :disabled="activeConfig.computeType === 'shared'"
              />
              <BaseStepper
                v-model="activeConfig.ramGBPerNode"
                :min="activeConfig.computeType === 'shared' ? 0 : 1"
                :max="1024"
                :step="4"
                label="Memori RAM"
                unit="GB"
                inline-unit="GB"
                :disabled="activeConfig.computeType === 'shared'"
              />
            </div>
          </BaseAdvancedToggle>
        </BaseSectionCard>

        <!-- ======================================================== -->
        <!-- STEP 2: PENYIMPANAN DATA & BACKUP CADANGAN              -->
        <!-- ======================================================== -->
        <BaseSectionCard
          step-number="2"
          title="Penyimpanan Data &amp; Proteksi Cadangan"
          guide-question="Berapa kapasitas harddisk dan apakah membutuhkan pencadangan data otomatis?"
          technical-label="SAN Storage &amp; NBU"
          :has-changes="isStorageCustomized"
        >
          <!-- Kapasitas Storage Stepper -->
          <div class="storage-stepper-row">
            <div class="field-hdr-wrap">
              <div class="title-with-desc">
                <span class="field-title-lbl">Kapasitas Disk per Server:</span>
                <span class="field-sub-hint">Alokasi ruang simpan berkas dan database</span>
              </div>
              <BaseHelpTip
                title="SAN Storage (Kapasitas Disk)"
                subtitle="SAN Virtual Disk LUN"
                :content="glossary.storage.shortDesc"
                :business-impact="glossary.storage.businessImpact"
              />
            </div>

            <div class="stepper-wrap-contained">
              <BaseStepper
                v-model="activeConfig.storageGBPerNode"
                :min="10"
                :max="20000"
                :step="50"
                label="Storage"
                unit="GB"
                inline-unit="GB"
              />
            </div>
          </div>

          <!-- Pilihan Kecepatan Storage & Backup -->
          <div class="storage-options-grid">
            <!-- Pilihan Kelas Kecepatan Storage -->
            <div class="option-column">
              <div class="field-hdr-wrap">
                <span class="field-title-lbl">Kecepatan &amp; Kelas Disk:</span>
                <BaseHelpTip
                  title="Kelas Kecepatan Storage"
                  subtitle="Tier 2 vs Tier 3 ASM"
                  :content="glossary.storageTier2.shortDesc"
                  :business-impact="glossary.storageTier3.businessImpact"
                />
              </div>
              <BaseChoiceCard
                v-model="activeConfig.storageClass"
                :columns="2"
                :options="[
                  {
                    value: 'Tier 2',
                    title: 'Standar Andal',
                    technicalLabel: 'SAN Tier 2',
                    description: 'Kecepatan standar enterprise untuk web, worker, dan log (Rp 33.000/GB/thn).',
                    icon: '💾'
                  },
                  {
                    value: 'Tier 3',
                    title: 'Ultra Cepat',
                    technicalLabel: 'SAN Tier 3 ASM',
                    description: 'Penyimpanan premium IOPS tinggi khusus database Oracle ASM (Rp 48.000/GB/thn).',
                    icon: '⚡',
                    badge: 'Performa Tinggi',
                    badgeType: 'warning'
                  }
                ]"
              />
            </div>

            <!-- Toggle Proteksi Backup NetBackup -->
            <div class="backup-protection-box">
              <div class="field-hdr-wrap">
                <span class="field-title-lbl">Proteksi Cadangan Data:</span>
                <BaseHelpTip
                  title="NetBackup (NBU)"
                  subtitle="Enterprise Daily Backup"
                  :content="glossary.backupNBU.shortDesc"
                  :business-impact="glossary.backupNBU.businessImpact"
                />
              </div>

              <label class="backup-card-label" :class="{ checked: activeConfig.backupEnabled }">
                <input
                  type="checkbox"
                  v-model="activeConfig.backupEnabled"
                  class="custom-chk-large"
                />
                <div class="backup-text-wrap">
                  <span class="backup-main-title">Aktifkan Backup Otomatis Harian (NetBackup)</span>
                  <span class="backup-desc">Pencadangan data berkala terpusat Bank Mandiri (Tarif: Rp 27.000 / GB / tahun).</span>
                </div>
              </label>
            </div>
          </div>
        </BaseSectionCard>

        <!-- ======================================================== -->
        <!-- STEP 3: JENIS LINGKUNGAN INFRASTRUKTUR                  -->
        <!-- ======================================================== -->
        <BaseSectionCard
          step-number="3"
          title="Jenis Lingkungan Infrastruktur"
          guide-question="Di atas platform komputasi mana komponen ini akan di-deploy?"
          technical-label="VM / Baremetal / Shared"
          :has-changes="activeConfig.computeType !== recommendedComputeType"
        >
          <BaseChoiceCard
            v-model="activeConfig.computeType"
            :columns="3"
            :options="[
              {
                value: 'vm',
                title: 'Virtual Machine (VM)',
                technicalLabel: 'Server Virtual Fleksibel',
                description: 'Server virtual efisien dan hemat biaya. Paling fleksibel untuk sebagian besar workload.',
                icon: '🖥️'
              },
              {
                value: 'baremetal',
                title: 'Server Fisik Khusus',
                technicalLabel: 'Baremetal Dedicated',
                description: 'Mesin server fisik murni tanpa virtualisasi. Performa maksimal tanpa berbagi daya.',
                icon: '🏢',
                badge: 'Dedicated',
                badgeType: 'warning'
              },
              {
                value: 'shared',
                title: 'Penyimpanan Berbagi',
                technicalLabel: 'Shared / SAN LUN',
                description: 'Alokasi penyimpanan bersama terpusat tanpa server komputasi tambahan terpisah.',
                icon: '🔗'
              }
            ]"
          />
        </BaseSectionCard>

        <!-- ======================================================== -->
        <!-- STEP 4: MODEL LISENSI SOFTWARE                          -->
        <!-- ======================================================== -->
        <BaseSectionCard
          step-number="4"
          title="Pilihan Lisensi Software"
          guide-question="Gunakan versi open source gratis atau berlisensi resmi vendor dengan garansi support?"
          technical-label="OSS vs Enterprise"
          :has-changes="false"
        >
          <div class="license-choice-grid">
            <!-- Option 1: Community Free OSS -->
            <div
              class="license-choice-card"
              :class="{ active: activeConfig.licenseKind === 'free' }"
              @click="activeConfig.licenseKind = 'free'"
            >
              <div class="lic-card-hdr">
                <span class="lic-icon">🌱</span>
                <div class="lic-hdr-titles">
                  <span class="lic-name">Community / Free Open Source</span>
                  <span class="lic-price-tag">Biaya Lisensi: Rp 0 / thn</span>
                </div>
              </div>
              <p class="lic-product-name">{{ activeItem.license.free }}</p>
              <p class="lic-card-desc">Software gratis tanpa biaya langganan tahunan vendor. Pemeliharaan didukung oleh tim internal.</p>
            </div>

            <!-- Option 2: Enterprise Official License -->
            <div
              class="license-choice-card enterprise-card"
              :class="{ active: activeConfig.licenseKind === 'enterprise' }"
              @click="activeConfig.licenseKind = 'enterprise'"
            >
              <div class="lic-card-hdr">
                <span class="lic-icon">🏢</span>
                <div class="lic-hdr-titles">
                  <span class="lic-name">Enterprise / Berlisensi Resmi</span>
                  <span class="lic-price-tag font-emerald">Dukungan Resmi Vendor 24/7</span>
                </div>
              </div>
              <p class="lic-product-name">{{ activeItem.license.enterprise }}</p>
              <p class="lic-card-desc">Paket langganan resmi vendor dengan garansi SLA respon cepat, update patch keamanan, dan kepatuhan perbankan.</p>
            </div>
          </div>
        </BaseSectionCard>

        <!-- ======================================================== -->
        <!-- LIVE ESTIMATED COST BREAKDOWN FOR SELECTED LAYER         -->
        <!-- ======================================================== -->
        <div class="layer-cost-card">
          <div class="cost-card-header">
            <div class="cost-card-title-wrap">
              <span class="cost-card-title">Estimasi Biaya Tahunan (Layer {{ activeItem.no }})</span>
              <span class="cost-card-sub">Dihitung otomatis oleh Quote Engine berdasarkan katalog harga resmi</span>
            </div>
            <div class="cost-card-total-badge">
              <span class="total-lbl">Total Layer:</span>
              <BaseMoney
                :amount="activeLayerCost?.total_cost_annual"
                class="layer-grand-price"
                period="/ thn"
              />
            </div>
          </div>

          <div class="cost-metrics-row">
            <!-- Compute Cost -->
            <div class="cost-metric-box">
              <span class="metric-title">Biaya Compute</span>
              <BaseMoney :amount="activeLayerCost?.compute_cost_annual" period="/ thn" />
              <span class="metric-sub">{{ activeConfig.nodeCount }} Server · {{ activeConfig.nodeCount * activeConfig.vcpuPerNode }} vCPU</span>
            </div>

            <!-- Storage Cost -->
            <div class="cost-metric-box">
              <span class="metric-title">Biaya Storage</span>
              <BaseMoney :amount="activeLayerCost?.storage_cost_annual" period="/ thn" />
              <span class="metric-sub">{{ activeConfig.nodeCount * activeConfig.storageGBPerNode }} GB ({{ activeConfig.storageClass }})</span>
            </div>

            <!-- Backup Cost -->
            <div class="cost-metric-box">
              <span class="metric-title">Biaya Backup</span>
              <BaseMoney :amount="activeLayerCost?.backup_cost_annual" period="/ thn" />
              <span class="metric-sub">{{ activeConfig.backupEnabled ? 'NetBackup Aktif' : 'Tanpa NBU' }}</span>
            </div>

            <!-- License Cost -->
            <div class="cost-metric-box">
              <span class="metric-title">Biaya Lisensi</span>
              <template v-if="activeConfig.licenseKind === 'free'">
                <span class="free-badge">Rp 0 (Community OSS)</span>
              </template>
              <template v-else-if="activeLayerCost?.unconfirmed_prices && activeLayerCost.unconfirmed_prices.length > 0">
                <BaseMoney :amount="null" :is-unconfirmed="true" unconfirmed-text="Pending STA" />
              </template>
              <template v-else>
                <BaseMoney :amount="activeLayerCost?.license_cost_annual" period="/ thn" />
              </template>
              <span class="metric-sub text-truncate">{{ activeLayerCost?.license_product_name || activeItem.license.enterprise }}</span>
            </div>
          </div>

          <!-- Warning for unconfirmed prices if any -->
          <div v-if="activeLayerCost?.unconfirmed_prices && activeLayerCost.unconfirmed_prices.length > 0" class="layer-warning-strip">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="#D97706" stroke-width="2">
              <circle cx="12" cy="12" r="10"/>
              <line x1="12" y1="8" x2="12" y2="12"/>
              <line x1="12" y1="16" x2="12.01" y2="16"/>
            </svg>
            <span>Catatan STA: {{ activeLayerCost.unconfirmed_prices[0] }}</span>
          </div>
        </div>

        <!-- ======================================================== -->
        <!-- STEP 5: ACUAN BASELINE & CATATAN EA FORUM (Collapsible)  -->
        <!-- ======================================================== -->
        <BaseSectionCard
          step-number="5"
          title="Acuan Standar Baseline &amp; Catatan EA Forum"
          guide-question="Spesifikasi acuan resmi Bank Mandiri dan justifikasi Enterprise Architecture"
          technical-label="Baseline Standard"
          :collapsible="true"
          :default-collapsed="false"
        >
          <!-- Baseline Reference Card -->
          <div class="specs-layout-container">
            <!-- Tier Critical Box -->
            <div
              v-if="tierFilter === 'critical'"
              class="tier-detail-box box-critical"
            >
              <div class="tier-box-header">
                <div class="tier-indicator crit-dot"></div>
                <div class="tier-box-title-wrap">
                  <span class="tier-badge-title text-amber">Acuan Baseline: Tier Critical</span>
                  <span class="tier-box-subtitle">Livin' Merchant &amp; Super App · Standar High Availability</span>
                </div>
              </div>

              <div class="tier-specs-grid">
                <div class="spec-attr">
                  <span class="attr-label">Alokasi Node:</span>
                  <span class="attr-value">{{ activeItem.critical.nodes }}</span>
                </div>
                <div class="spec-attr">
                  <span class="attr-label">Tipe Compute:</span>
                  <span class="attr-value">{{ activeItem.critical.type }}</span>
                </div>
              </div>

              <div class="spec-config-strip">
                <span class="config-lbl">Spesifikasi Hardware Acuan:</span>
                <span class="config-code font-amber">{{ activeItem.critical.spec }}</span>
              </div>
            </div>

            <!-- Tier Moderate Box -->
            <div
              v-if="tierFilter === 'moderate'"
              class="tier-detail-box box-moderate"
            >
              <div class="tier-box-header">
                <div class="tier-indicator mod-dot"></div>
                <div class="tier-box-title-wrap">
                  <span class="tier-badge-title text-navy">Acuan Baseline: Tier Moderat</span>
                  <span class="tier-box-subtitle">IHSP &amp; RMT · Internal System &amp; Standard Workload</span>
                </div>
              </div>

              <div class="tier-specs-grid">
                <div class="spec-attr">
                  <span class="attr-label">Alokasi Node:</span>
                  <span class="attr-value">{{ activeItem.moderate.nodes }}</span>
                </div>
                <div class="spec-attr">
                  <span class="attr-label">Tipe Compute:</span>
                  <span class="attr-value">{{ activeItem.moderate.type }}</span>
                </div>
              </div>

              <div class="spec-config-strip">
                <span class="config-lbl">Spesifikasi Hardware Acuan:</span>
                <span class="config-code font-navy">{{ activeItem.moderate.spec }}</span>
              </div>
            </div>
          </div>

          <!-- Architecture Notes & Justification -->
          <div class="pure-notes-section">
            <h4 class="section-subheading">Catatan Arsitektur &amp; Justifikasi Sizing (EA Forum)</h4>
            <p class="notes-body">{{ activeItem.notes }}</p>
            <div v-if="activeItem.contractRef" class="contract-ref-badge">
              <span>Referensi Kontrak / Master: {{ activeItem.contractRef }}</span>
            </div>
          </div>
        </BaseSectionCard>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { TechStackItem } from '~/data/techStackBaselineData'
import type { LayerAdjusterConfig, QuoteItemCost } from '~/types/sizing'
import { SIZING_GLOSSARY } from '~/data/sizingGlossary'
import BaseStepper from '~/components/base/BaseStepper.vue'
import BaseMoney from '~/components/base/BaseMoney.vue'
import BaseSectionCard from '~/components/base/BaseSectionCard.vue'
import BaseSizePicker from '~/components/base/BaseSizePicker.vue'
import BaseChoiceCard from '~/components/base/BaseChoiceCard.vue'
import BaseHelpTip from '~/components/base/BaseHelpTip.vue'
import BaseAdvancedToggle from '~/components/base/BaseAdvancedToggle.vue'

const props = defineProps<{
  layerItems: TechStackItem[]
  activeItemId: string
  activeItem: TechStackItem
  activeConfig: LayerAdjusterConfig
  activeLayerCost?: QuoteItemCost
  tierFilter: 'critical' | 'moderate'
  licenseFilter: 'free' | 'license'
  isLayerCustomized: (id: string) => boolean
}>()

defineEmits<{
  (e: 'select-layer', id: string): void
  (e: 'reset-active-layer'): void
}>()

const glossary = SIZING_GLOSSARY

// Check recommended baseline spec for sizing comparison
const recommendedVcpu = computed(() => {
  const isCrit = props.tierFilter === 'critical'
  if (props.activeItem.no === 2 || (isCrit && props.activeItem.no === 4)) return 64
  if (props.activeItem.no === 8 || props.activeItem.no === 9) return isCrit ? 16 : 8
  if ([1, 3, 6].includes(props.activeItem.no)) return 8
  return isCrit ? 8 : 4
})

const recommendedRam = computed(() => {
  const isCrit = props.tierFilter === 'critical'
  if (props.activeItem.no === 2 || (isCrit && props.activeItem.no === 4)) return 1024
  if (props.activeItem.no === 8 || props.activeItem.no === 9) return isCrit ? 64 : 32
  if (props.activeItem.no === 3) return 32
  if (props.activeItem.no === 1) return 16
  return isCrit ? 16 : 8
})

const recommendedComputeType = computed(() => {
  const isCrit = props.tierFilter === 'critical'
  if (props.activeItem.no === 5) return 'shared'
  if (isCrit && (props.activeItem.no === 2 || props.activeItem.no === 4)) return 'baremetal'
  return 'vm'
})

const isComputeCustomized = computed(() => {
  return (
    props.activeConfig.vcpuPerNode !== recommendedVcpu.value ||
    props.activeConfig.ramGBPerNode !== recommendedRam.value
  )
})

const isStorageCustomized = computed(() => {
  const isCrit = props.tierFilter === 'critical'
  let defaultDisk = 200
  if (props.activeItem.no === 2) defaultDisk = isCrit ? 300 : 1200
  else if (props.activeItem.no === 4) defaultDisk = 1000
  else if (props.activeItem.no === 5) defaultDisk = isCrit ? 5500 : 1500
  else if (props.activeItem.no === 8) defaultDisk = isCrit ? 1000 : 500
  else if (props.activeItem.no === 9) defaultDisk = isCrit ? 2000 : 500
  else if (props.activeItem.no === 12) defaultDisk = isCrit ? 5000 : 500

  return props.activeConfig.storageGBPerNode !== defaultDisk
})
</script>

<style scoped>
.sketch-workspace-grid {
  flex: 1;
  display: grid;
  grid-template-columns: 290px 1fr;
  overflow: hidden;
}

.sketch-list-section {
  background: #FFFFFF;
  border-right: 1px solid #E2E8F0;
  display: flex;
  flex-direction: column;
  overflow-y: auto;
}

.section-hdr {
  padding: 16px 18px 12px 18px;
  border-bottom: 1px solid #F1F5F9;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.section-title {
  font-size: 13px;
  font-weight: 800;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: #002060;
  margin: 0;
}

.item-count {
  font-size: 11px;
  font-weight: 700;
  color: #64748B;
  background: #F1F5F9;
  padding: 2px 7px;
  border-radius: 4px;
}

.bullet-list-container {
  display: flex;
  flex-direction: column;
  padding: 6px 0;
}

.bullet-row {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 11px 18px;
  cursor: pointer;
  border-left: 3px solid transparent;
  transition: all 0.15s ease;
  user-select: none;
  border-bottom: 1px solid #F8FAFC;
}

.bullet-row:hover {
  background: #F8FAFC;
}

.bullet-row.active {
  background: #F1F5F9;
  border-left-color: #002060;
}

.bullet-dot {
  font-size: 16px;
  color: #94A3B8;
  line-height: 1;
}

.bullet-row.active .bullet-dot {
  color: #002060;
}

.bullet-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
}

.bullet-title-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 6px;
}

.bullet-label {
  font-size: 13px;
  font-weight: 700;
  color: #1E293B;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.bullet-row.active .bullet-label {
  color: #002060;
}

.custom-badge {
  font-size: 9px;
  font-weight: 700;
  color: #B45309;
  background: #FEF3C7;
  padding: 1px 5px;
  border-radius: 3px;
  text-transform: uppercase;
}

.bullet-sub {
  font-size: 11px;
  color: #64748B;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  margin-top: 1px;
}

.bullet-chevron {
  font-size: 12px;
  color: #CBD5E1;
  font-weight: bold;
}

.bullet-row.active .bullet-chevron {
  color: #002060;
}

.sketch-card-section {
  padding: 20px 24px;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
}

.pure-content-card {
  background: #FFFFFF;
  border: 1px solid #E2E8F0;
  border-radius: 8px;
  padding: 22px 24px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.component-header-banner {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  border-bottom: 1px solid #E2E8F0;
  padding-bottom: 14px;
  gap: 16px;
}

.badge-row {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 4px;
}

.badge-layer {
  font-size: 11px;
  font-weight: 800;
  color: #002060;
  background: #EFF6FF;
  border: 1px solid #DBEAFE;
  padding: 3px 8px;
  border-radius: 4px;
  text-transform: uppercase;
}

.custom-status-tag {
  font-size: 10px;
  font-weight: 700;
  color: #B45309;
  background: #FEF3C7;
  border: 1px solid #FDE68A;
  padding: 2px 7px;
  border-radius: 4px;
}

.component-title {
  font-size: 20px;
  font-weight: 800;
  color: #0F172A;
  margin: 0;
  letter-spacing: -0.01em;
}

.banner-actions-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.filter-status-tags {
  display: flex;
  gap: 6px;
}

.active-filter-tag {
  font-size: 11px;
  font-weight: 700;
  padding: 3px 8px;
  border-radius: 4px;
}

.active-filter-tag.critical {
  background: #FEF3C7;
  color: #B45309;
  border: 1px solid #FDE68A;
}

.active-filter-tag.moderate {
  background: #E0F2FE;
  color: #0369A1;
  border: 1px solid #BAE6FD;
}

.active-filter-tag.license {
  background: #DCFCE7;
  color: #15803D;
  border: 1px solid #BBF7D0;
}

.active-filter-tag.free {
  background: #F1F5F9;
  color: #475569;
  border: 1px solid #CBD5E1;
}

.btn-reset-baseline {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  font-size: 11px;
  font-weight: 600;
  color: #475569;
  background: #FFFFFF;
  border: 1px solid #CBD5E1;
  border-radius: 4px;
  padding: 5px 10px;
  cursor: pointer;
  transition: all 0.15s ease;
  white-space: nowrap;
}

.btn-reset-baseline:hover {
  color: #002060;
  border-color: #002060;
  background: #F8FAFC;
}

/* ─────────────────────────────────────────────
   Step Section Form Elements
   ───────────────────────────────────────────── */
.form-row-block {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.field-hdr-wrap {
  display: flex;
  align-items: center;
  gap: 8px;
}

.title-with-desc {
  display: flex;
  flex-direction: column;
}

.field-title-lbl {
  font-size: 11px;
  font-weight: 800;
  color: #334155;
  text-transform: uppercase;
  letter-spacing: 0.03em;
}

.field-sub-hint {
  font-size: 10px;
  color: #64748B;
}

.server-count-row,
.storage-stepper-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  padding-top: 10px;
  border-top: 1px dashed #E2E8F0;
}

.stepper-wrap-contained {
  width: 200px;
}

.advanced-steppers-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}

.storage-options-grid {
  display: grid;
  grid-template-columns: 1.2fr 1fr;
  gap: 14px;
  padding-top: 10px;
  border-top: 1px dashed #E2E8F0;
}

.option-column {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.backup-protection-box {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.backup-card-label {
  background: #FFFFFF;
  border: 1px solid #CBD5E1;
  border-radius: 8px;
  padding: 12px 14px;
  display: flex;
  align-items: flex-start;
  gap: 10px;
  cursor: pointer;
  transition: all 0.15s ease;
  height: 100%;
  box-sizing: border-box;
}

.backup-card-label:hover {
  background: #F8FAFC;
  border-color: #94A3B8;
}

.backup-card-label.checked {
  background: #EFF6FF;
  border-color: #002060;
  box-shadow: 0 0 0 2px rgba(0, 32, 96, 0.12);
}

.custom-chk-large {
  width: 16px;
  height: 16px;
  cursor: pointer;
  accent-color: #002060;
  margin-top: 2px;
}

.backup-text-wrap {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.backup-main-title {
  font-size: 12px;
  font-weight: 800;
  color: #0F172A;
}

.backup-desc {
  font-size: 10px;
  color: #64748B;
  line-height: 1.35;
}

/* ─────────────────────────────────────────────
   Step 4: License Cards
   ───────────────────────────────────────────── */
.license-choice-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}

.license-choice-card {
  background: #FFFFFF;
  border: 1px solid #CBD5E1;
  border-radius: 8px;
  padding: 14px;
  cursor: pointer;
  transition: all 0.15s ease;
  display: flex;
  flex-direction: column;
  gap: 6px;
  text-align: left;
}

.license-choice-card:hover {
  border-color: #94A3B8;
  background: #F8FAFC;
}

.license-choice-card.active {
  background: #FFFFFF;
  border-color: #002060;
  box-shadow: 0 0 0 2px rgba(0, 32, 96, 0.12);
}

.license-choice-card.enterprise-card.active {
  border-color: #15803D;
  box-shadow: 0 0 0 2px rgba(21, 128, 61, 0.12);
}

.lic-card-hdr {
  display: flex;
  align-items: center;
  gap: 10px;
}

.lic-icon {
  font-size: 18px;
}

.lic-hdr-titles {
  display: flex;
  flex-direction: column;
}

.lic-name {
  font-size: 12px;
  font-weight: 800;
  color: #0F172A;
}

.lic-price-tag {
  font-size: 10px;
  font-weight: 700;
  color: #64748B;
}

.font-emerald {
  color: #15803D;
}

.lic-product-name {
  margin: 0;
  font-size: 12px;
  font-weight: 700;
  color: #002060;
  background: #F8FAFC;
  padding: 4px 8px;
  border-radius: 4px;
  border: 1px solid #E2E8F0;
}

.lic-card-desc {
  margin: 0;
  font-size: 11px;
  color: #64748B;
  line-height: 1.4;
}

/* ─────────────────────────────────────────────
   Layer Cost Breakdown Card
   ───────────────────────────────────────────── */
.layer-cost-card {
  background: #F1F5F9;
  border: 1px solid #CBD5E1;
  border-radius: 8px;
  padding: 16px 18px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.cost-card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid #E2E8F0;
  padding-bottom: 8px;
}

.cost-card-title-wrap {
  display: flex;
  flex-direction: column;
}

.cost-card-title {
  font-size: 13px;
  font-weight: 800;
  color: #002060;
}

.cost-card-sub {
  font-size: 11px;
  color: #64748B;
}

.cost-card-total-badge {
  display: flex;
  align-items: baseline;
  gap: 8px;
}

.total-lbl {
  font-size: 12px;
  font-weight: 700;
  color: #475569;
}

.layer-grand-price {
  font-size: 18px;
  font-weight: 800;
  color: #002060;
}

.cost-metrics-row {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 12px;
}

.cost-metric-box {
  background: #FFFFFF;
  border: 1px solid #E2E8F0;
  border-radius: 6px;
  padding: 10px 12px;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.metric-title {
  font-size: 11px;
  font-weight: 700;
  color: #64748B;
  text-transform: uppercase;
}

.metric-sub {
  font-size: 10px;
  color: #94A3B8;
  margin-top: 3px;
}

.text-truncate {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.free-badge {
  font-size: 12px;
  font-weight: 700;
  color: #15803D;
}

.layer-warning-strip {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 11px;
  font-weight: 600;
  color: #B45309;
  background: #FEF3C7;
  border: 1px solid #FDE68A;
  padding: 6px 10px;
  border-radius: 4px;
}

/* ─────────────────────────────────────────────
   Step 5: Baseline Reference & Notes
   ───────────────────────────────────────────── */
.specs-layout-container {
  display: grid;
  grid-template-columns: 1fr;
  gap: 14px;
}

.tier-detail-box {
  border-radius: 6px;
  padding: 14px;
  border: 1px solid #E2E8F0;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.tier-detail-box.box-critical {
  background: #FFFDF7;
  border-color: #FDE68A;
}

.tier-detail-box.box-moderate {
  background: #F8FAFC;
  border-color: #BAE6FD;
}

.tier-box-header {
  display: flex;
  align-items: center;
  gap: 8px;
}

.tier-indicator {
  width: 8px;
  height: 8px;
  border-radius: 50%;
}

.crit-dot {
  background: #D4A017;
}

.mod-dot {
  background: #0284C7;
}

.tier-badge-title {
  font-size: 12px;
  font-weight: 800;
}

.tier-box-subtitle {
  font-size: 10px;
  color: #64748B;
  display: block;
}

.text-amber {
  color: #B45309;
}

.font-amber {
  color: #B45309;
}

.text-navy {
  color: #002060;
}

.font-navy {
  color: #002060;
}

.tier-specs-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 6px;
}

.spec-attr {
  display: flex;
  flex-direction: column;
}

.attr-label {
  font-size: 10px;
  color: #64748B;
  text-transform: uppercase;
}

.attr-value {
  font-size: 12px;
  font-weight: 700;
  color: #0F172A;
}

.spec-config-strip {
  padding: 6px 8px;
  background: rgba(0, 0, 0, 0.02);
  border-radius: 4px;
  display: flex;
  flex-direction: column;
  gap: 1px;
}

.config-lbl {
  font-size: 10px;
  color: #64748B;
  text-transform: uppercase;
}

.config-code {
  font-size: 12px;
  font-weight: 700;
}

.pure-notes-section {
  border-top: 1px solid #E2E8F0;
  padding-top: 14px;
}

.section-subheading {
  font-size: 12px;
  font-weight: 800;
  color: #002060;
  text-transform: uppercase;
  letter-spacing: 0.04em;
  margin: 0 0 6px 0;
}

.notes-body {
  font-size: 12px;
  line-height: 1.5;
  color: #475569;
  margin: 0;
}

.contract-ref-badge {
  display: inline-block;
  margin-top: 8px;
  font-size: 11px;
  font-weight: 600;
  color: #0369A1;
  background: #E0F2FE;
  padding: 2px 8px;
  border-radius: 4px;
}
</style>
