<template>
  <div>
    <!-- MODAL 1: ADD / EDIT PRICING ITEM (CRUD) -->
    <div
      v-if="showPricingModal"
      class="sketch-modal-backdrop"
      @click.self="$emit('close-pricing-modal')"
    >
      <div class="sketch-modal-box modal-form-box">
        <div class="modal-hdr">
          <div class="modal-hdr-left">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="#002060" stroke-width="2">
              <rect x="3" y="3" width="7" height="7"/>
              <rect x="14" y="3" width="7" height="7"/>
              <rect x="14" y="14" width="7" height="7"/>
              <rect x="3" y="14" width="7" height="7"/>
            </svg>
            <div>
              <h3>{{ isEditingPricing ? 'Edit Item Harga' : 'Tambah Item SKU Baru' }}</h3>
              <p class="modal-hdr-sub">Master tarif infrastruktur &amp; lisensi resmi Enterprise Architecture</p>
            </div>
          </div>
          <button type="button" class="modal-close" @click="$emit('close-pricing-modal')">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18"></line>
              <line x1="6" y1="6" x2="18" y2="18"></line>
            </svg>
          </button>
        </div>

        <form @submit.prevent="$emit('save-pricing')" class="modal-content form-grid">
          <!-- Row 1: SKU & Kategori -->
          <div class="form-row">
            <div class="form-field">
              <label class="form-label">SKU Identifier <span class="req">*</span></label>
              <input
                type="text"
                v-model="pricingForm.sku"
                required
                class="form-input font-mono"
                placeholder="Contoh: VCPU_PROD, LIC_OCP_PROD"
              />
            </div>

            <div class="form-field">
              <label class="form-label">Kategori <span class="req">*</span></label>
              <select v-model="pricingForm.category" required class="form-input">
                <option value="compute">Compute (vCPU)</option>
                <option value="storage">Storage SAN</option>
                <option value="backup">Backup Netbackup</option>
                <option value="baremetal">Server Baremetal</option>
                <option value="license">Lisensi Software</option>
              </select>
            </div>
          </div>

          <!-- Row 2: Nama Item -->
          <div class="form-field">
            <label class="form-label">Nama Item / Software <span class="req">*</span></label>
            <input
              type="text"
              v-model="pricingForm.name"
              required
              class="form-input"
              placeholder="Contoh: Red Hat OpenShift Platform Plus Baremetal"
            />
          </div>

          <!-- Row 3: Environment & Satuan -->
          <div class="form-row">
            <div class="form-field">
              <label class="form-label">Environment <span class="req">*</span></label>
              <select v-model="pricingForm.environment" required class="form-input">
                <option value="Production">Production</option>
                <option value="Dev / Non-Prod">Dev / Non-Prod</option>
                <option value="All">All Environments</option>
              </select>
            </div>

            <div class="form-field">
              <label class="form-label">Satuan Tarif <span class="req">*</span></label>
              <input
                type="text"
                v-model="pricingForm.unit"
                required
                class="form-input"
                placeholder="vCPU, GB, Node, Unit, Instance, Shard"
              />
            </div>
          </div>

          <!-- Row 4: Status Konfirmasi STA -->
          <div class="form-field confirmation-toggle-field">
            <label class="checkbox-label">
              <input
                type="checkbox"
                v-model="pricingForm.is_confirmed"
                class="custom-chk"
              />
              <span class="chk-text">Harga Resmi Terkonfirmasi oleh Tim STA / Procurement Enterprise Architecture</span>
            </label>
            <span v-if="!pricingForm.is_confirmed" class="sta-warn-hint">
              Jika belum dikonfirmasi, harga akan disimpan sebagai pending dan muncul sebagai peringatan resmi pada Sizing BOQ.
            </span>
          </div>

          <!-- Row 5: Harga Satuan IDR -->
          <div class="form-field">
            <label class="form-label">
              Harga Satuan Per Tahun (IDR)
              <span v-if="pricingForm.is_confirmed" class="req">*</span>
            </label>
            <input
              type="number"
              v-model.number="pricingForm.unit_price_idr"
              :disabled="!pricingForm.is_confirmed"
              :required="pricingForm.is_confirmed"
              class="form-input font-bold"
              placeholder="Contoh: 8250000"
            />
            <span v-if="pricingForm.unit_price_idr" class="price-preview-label">
              Format: Rp {{ formatSimpleIDR(pricingForm.unit_price_idr) }}
            </span>
          </div>

          <!-- Row 6: Sumber Referensi -->
          <div class="form-field">
            <label class="form-label">Sumber Dokumen / Referensi Kontrak</label>
            <input
              type="text"
              v-model="pricingForm.source"
              class="form-input"
              placeholder="Contoh: Template Sizing.xlsx!Pricing, TIF.DDL/GH.9076/2025"
            />
          </div>

          <div class="modal-ftr">
            <button type="button" class="modal-cancel-btn" @click="$emit('close-pricing-modal')">Batal</button>
            <button type="submit" class="modal-submit-btn" :disabled="isSavingPricing">
              {{ isSavingPricing ? 'Menyimpan...' : (isEditingPricing ? 'Simpan Perubahan' : 'Tambah ke Katalog') }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- MODAL 2: DELETE PRICING ITEM CONFIRMATION -->
    <div
      v-if="showDeletePricingModal"
      class="sketch-modal-backdrop"
      @click.self="$emit('close-delete-modal')"
    >
      <div class="sketch-modal-box delete-confirm-box">
        <div class="modal-hdr">
          <div class="modal-hdr-left">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="#DC2626" stroke-width="2">
              <path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"/>
              <line x1="12" y1="9" x2="12" y2="13"/>
              <line x1="12" y1="17" x2="12.01" y2="17"/>
            </svg>
            <h3>Konfirmasi Hapus Item Katalog</h3>
          </div>
          <button type="button" class="modal-close" @click="$emit('close-delete-modal')">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18"></line>
              <line x1="6" y1="6" x2="18" y2="18"></line>
            </svg>
          </button>
        </div>

        <div class="modal-content">
          <p class="modal-desc">
            Apakah Anda yakin ingin menghapus item SKU <strong>{{ deletingPricingItem?.sku }}</strong> ({{ deletingPricingItem?.name }}) dari katalog database PostgreSQL?
          </p>
          <p class="delete-warn-note">
            Perubahan ini akan langsung mempengaruhi kalkulasi otomatis pada menu Sizing.
          </p>
        </div>

        <div class="modal-ftr">
          <button type="button" class="modal-cancel-btn" @click="$emit('close-delete-modal')">Batal</button>
          <button
            type="button"
            class="modal-delete-btn"
            :disabled="isDeletingPricing"
            @click="$emit('execute-delete')"
          >
            {{ isDeletingPricing ? 'Menghapus...' : 'Ya, Hapus Item' }}
          </button>
        </div>
      </div>
    </div>

    <!-- MODAL 3: UPLOAD DATA SIZING / TELEMETRI -->
    <div
      v-if="showUploadModal"
      class="sketch-modal-backdrop"
      @click.self="$emit('close-upload-modal')"
    >
      <div class="sketch-modal-box">
        <div class="modal-hdr">
          <div class="modal-hdr-left">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
              <polyline points="17 8 12 3 7 8"/>
              <line x1="12" y1="3" x2="12" y2="15"/>
            </svg>
            <h3>Upload Data Sizing / Telemetri</h3>
          </div>
          <button type="button" class="modal-close" @click="$emit('close-upload-modal')">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18"></line>
              <line x1="6" y1="6" x2="18" y2="18"></line>
            </svg>
          </button>
        </div>

        <div class="modal-content">
          <p class="modal-desc">
            Unggah berkas Excel template (*.xlsx) atau file telemetri SolarWinds (*.csv) untuk kalkulasi otomatis.
          </p>

          <div
            class="upload-dropzone"
            @dragover.prevent
            @drop.prevent="$emit('file-drop', $event)"
          >
            <input
              type="file"
              ref="fileInputRef"
              class="hidden-file-input"
              accept=".xlsx,.csv,.xls"
              @change="$emit('file-change', $event)"
            />
            <div class="drop-icon">
              <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"></path>
              </svg>
            </div>
            <div class="drop-text">
              <strong>Klik untuk memilih berkas</strong> atau seret file ke sini
            </div>
            <span class="drop-sub">Mendukung file .xlsx resmi Enterprise Architecture &amp; .csv SolarWinds</span>

            <button
              type="button"
              class="browse-btn"
              @click="fileInputRef?.click()"
            >
              Pilih Berkas Komputer
            </button>
          </div>

          <div v-if="uploadedFileName" class="uploaded-status">
            <span class="check-icon">
              <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
                <polyline points="20 6 9 17 4 12"></polyline>
              </svg>
            </span> Berkas terpilih: <strong>{{ uploadedFileName }}</strong>
          </div>
        </div>

        <div class="modal-ftr">
          <button type="button" class="modal-cancel-btn" @click="$emit('close-upload-modal')">Batal</button>
          <button type="button" class="modal-submit-btn" :disabled="!uploadedFileName" @click="$emit('process-upload')">
            Proses Data
          </button>
        </div>
      </div>
    </div>

    <!-- MODAL 4: GENERATE BILL OF QUANTITIES & EXCEL DOWNLOAD -->
    <div
      v-if="showGenerateModal"
      class="sketch-modal-backdrop"
      @click.self="$emit('close-generate-modal')"
    >
      <div class="sketch-modal-box generate-modal-large">
        <div class="modal-hdr">
          <div class="modal-hdr-left">
            <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="#D4A017" stroke-width="2">
              <polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2"/>
            </svg>
            <div>
              <h3>Hasil Generate Sizing &amp; Bill of Quantity (BOQ)</h3>
              <p class="modal-hdr-sub">Rekapitulasi 13 Layer Arsitektur &amp; Proyeksi Biaya Resmi Enterprise Architecture</p>
            </div>
          </div>
          <button type="button" class="modal-close" @click="$emit('close-generate-modal')">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18"></line>
              <line x1="6" y1="6" x2="18" y2="18"></line>
            </svg>
          </button>
        </div>

        <div class="modal-content modal-content-boq">
          <!-- Project Info Field & KPI Metrics -->
          <div class="boq-project-bar">
            <div class="project-input-wrap">
              <label class="project-input-lbl">Nama Project / Inisiatif Sistem:</label>
              <input
                type="text"
                :value="projectName"
                @input="$emit('update:projectName', ($event.target as HTMLInputElement).value)"
                class="project-input"
                placeholder="Contoh: Livin' Merchant Release 6 &amp; SuperApp Core"
              />
            </div>

            <div class="boq-meta-badges">
              <span class="boq-badge">Env: {{ environment }}</span>
              <span class="boq-badge">{{ includeDRC ? 'DC + DRC (2x Multiplier)' : 'Single Site (1x)' }}</span>
              <span class="boq-badge">13 Architecture Layers</span>
            </div>
          </div>

          <!-- Summary Metric Cards -->
          <div class="boq-summary-cards">
            <div class="boq-stat-card">
              <span class="stat-lbl">Total Node Alokasi</span>
              <span class="stat-num text-navy">{{ quoteResponse?.total_nodes }} Node</span>
              <span class="stat-sub">{{ includeDRC ? 'Termasuk replikasi DRC' : 'Alokasi DC Utama' }}</span>
            </div>
            <div class="boq-stat-card">
              <span class="stat-lbl">Kapasitas Compute Total</span>
              <span class="stat-num font-amber">{{ quoteResponse?.total_vcpu }} vCPU</span>
              <span class="stat-sub">{{ quoteResponse?.total_ram_gb }} GB Total RAM</span>
            </div>
            <div class="boq-stat-card">
              <span class="stat-lbl">Kapasitas Storage SAN</span>
              <span class="stat-num text-navy">{{ ((quoteResponse?.total_storage_gb || 0) / 1000).toFixed(1) }} TB</span>
              <span class="stat-sub">{{ quoteResponse?.total_storage_gb }} GB SAN Storage</span>
            </div>
            <div class="boq-stat-card card-grand-total">
              <span class="stat-lbl">TOTAL BIAYA SIZING (PERPETUAL)</span>
              <BaseMoney
                :amount="quoteResponse?.grand_total_annual"
                class="stat-grand-price"
              />
              <span class="stat-sub">3-Yr TCO: Rp {{ formatSimpleIDR(quoteResponse?.tco_3year) }}</span>
            </div>
          </div>

          <!-- Comprehensive 13-Layer BOQ Table -->
          <div class="boq-table-container">
            <table class="boq-table">
              <thead>
                <tr>
                  <th class="th-no">No</th>
                  <th>Layer Arsitektur</th>
                  <th>Komponen Sistem</th>
                  <th class="text-center">Nodes</th>
                  <th class="text-center">vCPU</th>
                  <th class="text-center">RAM</th>
                  <th class="text-center">Storage</th>
                  <th>Model Lisensi</th>
                  <th class="text-right">Biaya Compute</th>
                  <th class="text-right">Biaya Storage</th>
                  <th class="text-right">Biaya Lisensi</th>
                  <th class="text-right">Total Biaya (Perpetual)</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="(item, idx) in quoteResponse?.items" :key="item.component_id">
                  <td class="td-no">{{ idx + 1 }}</td>
                  <td class="font-medium text-navy">{{ item.layer_name }}</td>
                  <td>{{ item.component_name }}</td>
                  <td class="text-center">{{ item.node_count }}</td>
                  <td class="text-center">{{ item.total_vcpu }}</td>
                  <td class="text-center">{{ item.total_ram_gb }}G</td>
                  <td class="text-center">{{ item.total_storage_gb }}G</td>
                  <td class="license-td">
                    <span :class="['lic-badge-tag', item.license_kind]">
                      {{ item.license_product_name }}
                    </span>
                  </td>
                  <td class="text-right text-tabular">
                    <BaseMoney :amount="item.compute_cost_annual" />
                  </td>
                  <td class="text-right text-tabular">
                    <BaseMoney :amount="item.storage_cost_annual" />
                  </td>
                  <td class="text-right text-tabular">
                    <template v-if="item.license_kind === 'free'">
                      <span class="free-text">Free</span>
                    </template>
                    <template v-else-if="item.unconfirmed_prices && item.unconfirmed_prices.length > 0">
                      <span class="unconfirmed-pill" title="Pending STA">STA</span>
                    </template>
                    <template v-else>
                      <BaseMoney :amount="item.license_cost_annual" />
                    </template>
                  </td>
                  <td class="text-right text-tabular font-semibold text-navy">
                    <BaseMoney :amount="item.total_cost_annual" />
                  </td>
                </tr>
              </tbody>
              <tfoot>
                <tr class="tfoot-total-row">
                  <td colspan="3" class="font-bold text-navy">TOTAL KESELURUHAN (ANNUAL OPEX)</td>
                  <td class="text-center font-bold">{{ quoteResponse?.total_nodes }}</td>
                  <td class="text-center font-bold">{{ quoteResponse?.total_vcpu }}</td>
                  <td class="text-center font-bold">{{ quoteResponse?.total_ram_gb }}G</td>
                  <td class="text-center font-bold">{{ quoteResponse?.total_storage_gb }}G</td>
                  <td>-</td>
                  <td class="text-right font-bold">
                    <BaseMoney :amount="quoteResponse?.subtotal_compute_annual" />
                  </td>
                  <td class="text-right font-bold">
                    <BaseMoney :amount="quoteResponse?.subtotal_storage_annual" />
                  </td>
                  <td class="text-right font-bold">
                    <BaseMoney :amount="quoteResponse?.subtotal_license_annual" />
                  </td>
                  <td class="text-right font-bold grand-tfoot-price">
                    <BaseMoney :amount="quoteResponse?.grand_total_annual" />
                  </td>
                </tr>
              </tfoot>
            </table>
          </div>

          <!-- Unconfirmed Warning Alert if any -->
          <div v-if="quoteResponse?.warnings && quoteResponse.warnings.length > 0" class="boq-warning-notice">
            <div class="warning-notice-hdr">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="#D97706" stroke-width="2">
                <path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"/>
                <line x1="12" y1="9" x2="12" y2="13"/>
                <line x1="12" y1="17" x2="12.01" y2="17"/>
              </svg>
              <strong>Perhatian Tim Procurement / STA (Enterprise Architecture):</strong>
            </div>
            <ul class="warning-list">
              <li v-for="(warn, wIdx) in quoteResponse.warnings" :key="wIdx">{{ warn }}</li>
            </ul>
          </div>
        </div>

        <div class="modal-ftr boq-modal-ftr">
          <button type="button" class="modal-cancel-btn" @click="$emit('close-generate-modal')">Tutup</button>

          <div class="modal-ftr-actions">
            <button
              type="button"
              class="copy-summary-action"
              @click="$emit('copy-summary')"
            >
              {{ copyStatus || 'Salin Ringkasan Sizing' }}
            </button>

            <button
              type="button"
              class="download-excel-action"
              :disabled="isDownloadingExcel"
              @click="$emit('download-excel')"
            >
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
                <polyline points="7 10 12 15 17 10"/>
                <line x1="12" y1="15" x2="12" y2="3"/>
              </svg>
              <span>{{ isDownloadingExcel ? 'Membuat Spreadsheet...' : 'Unduh Template Sizing (.xlsx)' }}</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import type { PricingItemRecord, PricingFormState, QuoteResponse } from '~/types/sizing'
import BaseMoney from '~/components/base/BaseMoney.vue'

defineProps<{
  showPricingModal: boolean
  isEditingPricing: boolean
  isSavingPricing: boolean
  pricingForm: PricingFormState
  showDeletePricingModal: boolean
  isDeletingPricing: boolean
  deletingPricingItem: PricingItemRecord | null
  showUploadModal: boolean
  uploadedFileName: string
  showGenerateModal: boolean
  projectName: string
  environment: string
  includeDRC: boolean
  quoteResponse: QuoteResponse | null
  copyStatus: string
  isDownloadingExcel: boolean
}>()

defineEmits<{
  (e: 'close-pricing-modal'): void
  (e: 'save-pricing'): void
  (e: 'close-delete-modal'): void
  (e: 'execute-delete'): void
  (e: 'close-upload-modal'): void
  (e: 'file-change', ev: Event): void
  (e: 'file-drop', ev: DragEvent): void
  (e: 'process-upload'): void
  (e: 'close-generate-modal'): void
  (e: 'update:projectName', val: string): void
  (e: 'copy-summary'): void
  (e: 'download-excel'): void
}>()

const fileInputRef = ref<HTMLInputElement | null>(null)

function formatSimpleIDR(val?: number | null): string {
  if (val === null || val === undefined || isNaN(val)) return '0'
  return new Intl.NumberFormat('id-ID').format(Math.round(val))
}
</script>

<style scoped>
.sketch-modal-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(15, 23, 42, 0.6);
  backdrop-filter: blur(2px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 20px;
}

.sketch-modal-box {
  background: #FFFFFF;
  border-radius: 8px;
  width: 100%;
  max-width: 520px;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.2);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.modal-form-box {
  max-width: 600px;
}

.delete-confirm-box {
  max-width: 440px;
}

.generate-modal-large {
  max-width: 1100px;
  max-height: 90vh;
}

.modal-hdr {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  border-bottom: 1px solid #E2E8F0;
  background: #F8FAFC;
}

.modal-hdr-left {
  display: flex;
  align-items: center;
  gap: 10px;
}

.modal-hdr-left h3 {
  margin: 0;
  font-size: 15px;
  font-weight: 800;
  color: #002060;
}

.modal-hdr-sub {
  margin: 2px 0 0 0;
  font-size: 11px;
  color: #64748B;
}

.modal-close {
  background: transparent;
  border: none;
  font-size: 16px;
  color: #64748B;
  cursor: pointer;
}

.modal-content {
  padding: 20px;
  overflow-y: auto;
}

.form-grid {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}

.form-field {
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.form-label {
  font-size: 11px;
  font-weight: 700;
  color: #334155;
  text-transform: uppercase;
  letter-spacing: 0.03em;
}

.req {
  color: #DC2626;
}

.form-input {
  padding: 8px 10px;
  border: 1px solid #CBD5E1;
  border-radius: 5px;
  font-size: 12px;
  color: #0F172A;
  outline: none;
  transition: border-color 0.15s ease;
}

.form-input:focus {
  border-color: #002060;
}

.form-input:disabled {
  background: #F1F5F9;
  color: #94A3B8;
  cursor: not-allowed;
}

.font-mono {
  font-family: 'JetBrains Mono', monospace;
}

.confirmation-toggle-field {
  background: #F8FAFC;
  border: 1px solid #E2E8F0;
  padding: 10px 12px;
  border-radius: 6px;
  gap: 6px;
}

.checkbox-label {
  display: inline-flex;
  align-items: center;
  gap: 7px;
  font-size: 12px;
  font-weight: 600;
  color: #334155;
  cursor: pointer;
  user-select: none;
}

.custom-chk {
  width: 15px;
  height: 15px;
  cursor: pointer;
  accent-color: #002060;
}

.chk-text {
  font-size: 11px;
}

.sta-warn-hint {
  font-size: 11px;
  color: #B45309;
  line-height: 1.4;
}

.price-preview-label {
  font-size: 11px;
  font-weight: 600;
  color: #002060;
}

.delete-warn-note {
  font-size: 12px;
  color: #DC2626;
  background: #FEF2F2;
  border: 1px solid #FECACA;
  padding: 8px 12px;
  border-radius: 4px;
}

.modal-ftr {
  padding: 12px 20px;
  border-top: 1px solid #E2E8F0;
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 10px;
  background: #F8FAFC;
}

.modal-cancel-btn {
  padding: 7px 14px;
  background: transparent;
  border: 1px solid #CBD5E1;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 600;
  color: #475569;
  cursor: pointer;
}

.modal-submit-btn {
  padding: 7px 16px;
  background: #002060;
  color: #FFFFFF;
  border: none;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 700;
  cursor: pointer;
}

.modal-submit-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.modal-delete-btn {
  padding: 7px 16px;
  background: #DC2626;
  color: #FFFFFF;
  border: none;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 700;
  cursor: pointer;
}

.modal-delete-btn:hover:not(:disabled) {
  background: #B91C1C;
}

.modal-content-boq {
  padding: 18px 24px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.modal-desc {
  font-size: 13px;
  color: #475569;
  margin: 0 0 16px 0;
}

.upload-dropzone {
  border: 2px dashed #CBD5E1;
  border-radius: 8px;
  padding: 30px 20px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  text-align: center;
  background: #F8FAFC;
}

.hidden-file-input {
  display: none;
}

.drop-icon {
  font-size: 32px;
}

.drop-text {
  font-size: 13px;
  color: #334155;
}

.drop-sub {
  font-size: 11px;
  color: #64748B;
}

.browse-btn {
  margin-top: 8px;
  padding: 6px 14px;
  background: #FFFFFF;
  border: 1px solid #CBD5E1;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 600;
  cursor: pointer;
}

.uploaded-status {
  margin-top: 12px;
  font-size: 12px;
  color: #15803D;
  background: #DCFCE7;
  padding: 6px 12px;
  border-radius: 4px;
}

.boq-modal-ftr {
  justify-content: space-between;
}

.modal-ftr-actions {
  display: flex;
  align-items: center;
  gap: 10px;
}

.boq-project-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  padding-bottom: 12px;
  border-bottom: 1px solid #E2E8F0;
}

.project-input-wrap {
  display: flex;
  align-items: center;
  gap: 10px;
  flex: 1;
}

.project-input-lbl {
  font-size: 12px;
  font-weight: 700;
  color: #334155;
  white-space: nowrap;
}

.project-input {
  flex: 1;
  padding: 6px 10px;
  border: 1px solid #CBD5E1;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 600;
  color: #0F172A;
  outline: none;
}

.project-input:focus {
  border-color: #002060;
}

.boq-meta-badges {
  display: flex;
  gap: 6px;
}

.boq-badge {
  font-size: 11px;
  font-weight: 600;
  color: #002060;
  background: #EFF6FF;
  border: 1px solid #BFDBFE;
  padding: 3px 8px;
  border-radius: 4px;
}

.boq-summary-cards {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 12px;
}

.boq-stat-card {
  background: #F8FAFC;
  border: 1px solid #CBD5E1;
  border-radius: 6px;
  padding: 10px 14px;
  display: flex;
  flex-direction: column;
}

.card-grand-total {
  background: #EFF6FF;
  border-color: #BFDBFE;
}

.stat-lbl {
  font-size: 10px;
  font-weight: 700;
  color: #64748B;
  text-transform: uppercase;
}

.stat-num {
  font-size: 18px;
  font-weight: 800;
  margin: 2px 0;
}

.text-navy {
  color: #002060;
}

.font-amber {
  color: #B45309;
}

.stat-grand-price {
  font-size: 18px;
  font-weight: 800;
  color: #002060;
  margin: 2px 0;
}

.stat-sub {
  font-size: 10px;
  color: #64748B;
}

.boq-table-container {
  border: 1px solid #CBD5E1;
  border-radius: 6px;
  overflow-x: auto;
  max-height: 380px;
}

.boq-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 11px;
}

.boq-table th {
  background: #002060;
  color: #FFFFFF;
  padding: 8px 10px;
  font-weight: 700;
  text-align: left;
  white-space: nowrap;
}

.boq-table td {
  padding: 7px 10px;
  border-bottom: 1px solid #E2E8F0;
  color: #1E293B;
  vertical-align: middle;
}

.boq-table tr:hover td {
  background: #F8FAFC;
}

.th-no,
.td-no {
  width: 36px;
  text-align: center;
  color: #64748B;
}

.text-center {
  text-align: center;
}

.text-right {
  text-align: right;
}

.text-tabular {
  font-variant-numeric: tabular-nums;
}

.font-medium {
  font-weight: 600;
}

.font-bold {
  font-weight: 700;
}

.font-semibold {
  font-weight: 600;
}

.lic-badge-tag {
  display: inline-block;
  font-size: 10px;
  font-weight: 600;
  padding: 1px 6px;
  border-radius: 3px;
  max-width: 160px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.lic-badge-tag.enterprise {
  background: #DCFCE7;
  color: #15803D;
}

.lic-badge-tag.free {
  background: #F1F5F9;
  color: #475569;
}

.free-text {
  color: #15803D;
  font-weight: 600;
}

.unconfirmed-pill {
  font-size: 9px;
  font-weight: 700;
  color: #D97706;
  background: #FFFBEB;
  border: 1px solid #FDE68A;
  padding: 1px 4px;
  border-radius: 3px;
}

.tfoot-total-row td {
  background: #F1F5F9;
  border-top: 2px solid #002060;
  padding: 10px;
}

.grand-tfoot-price {
  font-size: 13px;
  color: #002060;
}

.boq-warning-notice {
  background: #FFFBEB;
  border: 1px solid #FDE68A;
  border-radius: 6px;
  padding: 10px 14px;
}

.warning-notice-hdr {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 11px;
  color: #B45309;
  margin-bottom: 4px;
}

.warning-list {
  margin: 0;
  padding-left: 20px;
  font-size: 11px;
  color: #92400E;
  line-height: 1.4;
}

.copy-summary-action {
  padding: 7px 14px;
  background: #FFFFFF;
  border: 1px solid #CBD5E1;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 600;
  color: #334155;
  cursor: pointer;
  transition: all 0.15s ease;
}

.copy-summary-action:hover {
  background: #F1F5F9;
  color: #002060;
}

.download-excel-action {
  display: inline-flex;
  align-items: center;
  gap: 7px;
  padding: 7px 16px;
  background: #15803D;
  color: #FFFFFF;
  border: none;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.15s ease;
}

.download-excel-action:hover:not(:disabled) {
  background: #166534;
}

.download-excel-action:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
</style>
