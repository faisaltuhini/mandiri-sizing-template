<template>
  <div class="catalog-view-container">
    <!-- CATALOG TOPBAR WITH FILTERS & ADD BUTTON -->
    <header class="sketch-topbar catalog-topbar">
      <div class="topbar-left">
        <div class="title-with-toggle">
          <button
            type="button"
            class="topbar-sidebar-toggle"
            @click="$emit('toggle-sidebar')"
            :title="isSidebarCollapsed ? 'Perluas Sidebar' : 'Kecilkan / Minimize Sidebar'"
          >
            <svg width="17" height="17" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <rect x="3" y="3" width="18" height="18" rx="2" ry="2"></rect>
              <line x1="9" y1="3" x2="9" y2="21"></line>
            </svg>
          </button>
          <h1 class="page-title">Katalog Harga &amp; Master Data SKU</h1>
        </div>
        <p class="page-subtitle">Master tarif infrastruktur &amp; lisensi software Enterprise Architecture (Template Sizing.xlsx)</p>
      </div>

      <div class="topbar-filters">
        <!-- 1 MINIMALIST FILTER BUTTON WITH FLOATING POPUP (KATALOG) -->
        <div class="filter-popover-wrapper" ref="popoverRef">
          <button
            type="button"
            class="minimal-filter-btn"
            :class="{ active: isOpen || categoryFilter !== 'all' || envFilter !== 'all' }"
            @click="isOpen = !isOpen"
            title="Filter Kategori &amp; Environment Katalog"
          >
            <svg class="filter-svg-icon" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <polygon points="22 3 2 3 10 12.46 10 19 14 21 14 12.46 22 3"></polygon>
            </svg>
            <span class="filter-btn-title">Filter:</span>
            <span :class="['filter-pill', categoryFilter !== 'all' ? 'pill-cat' : 'pill-all']">
              {{ getCategoryLabel(categoryFilter) }}
            </span>
            <span class="filter-pill-dot">·</span>
            <span :class="['filter-pill', envFilter !== 'all' ? 'pill-env' : 'pill-all']">
              {{ getEnvLabel(envFilter) }}
            </span>
            <svg class="chevron-svg" :class="{ rotated: isOpen }" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
              <polyline points="6 9 12 15 18 9"></polyline>
            </svg>
          </button>

          <!-- FLOATING POPUP / POPOVER FOR KATALOG -->
          <transition name="popover-fade">
            <div v-if="isOpen" class="filter-popup-modal" @click.stop>
              <div class="popup-header">
                <div class="popup-title-wrap">
                  <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="#002060" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round">
                    <polygon points="22 3 2 3 10 12.46 10 19 14 21 14 12.46 22 3"></polygon>
                  </svg>
                  <span class="popup-heading">Filter Katalog Harga</span>
                </div>
                <button type="button" class="popup-close-btn" @click="isOpen = false" title="Tutup">
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <line x1="18" y1="6" x2="6" y2="18"></line>
                    <line x1="6" y1="6" x2="18" y2="18"></line>
                  </svg>
                </button>
              </div>

              <div class="popup-body">
                <!-- Section 1: Kategori Item -->
                <div class="popup-section">
                  <div class="section-label-row">
                    <label class="section-field-lbl">Kategori Infrastruktur:</label>
                    <span class="section-badge badge-cat">
                      {{ getCategoryLabel(categoryFilter) }}
                    </span>
                  </div>
                  <div class="segmented-control popup-segmented wrap-tabs">
                    <button
                      type="button"
                      :class="['seg-btn', { active: categoryFilter === 'all' }]"
                      @click="$emit('update:categoryFilter', 'all')"
                    >
                      Semua
                    </button>
                    <button
                      type="button"
                      :class="['seg-btn', { active: categoryFilter === 'compute' }]"
                      @click="$emit('update:categoryFilter', 'compute')"
                    >
                      Compute
                    </button>
                    <button
                      type="button"
                      :class="['seg-btn', { active: categoryFilter === 'storage' }]"
                      @click="$emit('update:categoryFilter', 'storage')"
                    >
                      Storage
                    </button>
                    <button
                      type="button"
                      :class="['seg-btn', { active: categoryFilter === 'backup' }]"
                      @click="$emit('update:categoryFilter', 'backup')"
                    >
                      Backup
                    </button>
                    <button
                      type="button"
                      :class="['seg-btn', { active: categoryFilter === 'license' }]"
                      @click="$emit('update:categoryFilter', 'license')"
                    >
                      Lisensi
                    </button>
                  </div>
                </div>

                <!-- Section 2: Environment Filter -->
                <div class="popup-section">
                  <div class="section-label-row">
                    <label class="section-field-lbl">Target Environment:</label>
                    <span class="section-badge badge-env">
                      {{ getEnvLabel(envFilter) }}
                    </span>
                  </div>
                  <div class="segmented-control popup-segmented">
                    <button
                      type="button"
                      :class="['seg-btn', { active: envFilter === 'all' }]"
                      @click="$emit('update:envFilter', 'all')"
                    >
                      Semua Env
                    </button>
                    <button
                      type="button"
                      :class="['seg-btn', { active: envFilter === 'Production' }]"
                      @click="$emit('update:envFilter', 'Production')"
                    >
                      Production
                    </button>
                    <button
                      type="button"
                      :class="['seg-btn', { active: envFilter === 'Dev / Non-Prod' }]"
                      @click="$emit('update:envFilter', 'Dev / Non-Prod')"
                    >
                      Dev / Non-Prod
                    </button>
                  </div>
                </div>
              </div>

              <div class="popup-footer">
                <button
                  v-if="categoryFilter !== 'all' || envFilter !== 'all'"
                  type="button"
                  class="popup-reset-link"
                  @click="$emit('update:categoryFilter', 'all'); $emit('update:envFilter', 'all')"
                >
                  Reset Filter
                </button>
                <span v-else class="popup-active-note">
                  Menampilkan semua SKU
                </span>
                <button type="button" class="popup-apply-btn" @click="isOpen = false">
                  Terapkan &amp; Tutup
                </button>
              </div>
            </div>
          </transition>
        </div>

        <!-- Add SKU Item Button -->
        <button
          type="button"
          class="action-btn-generate btn-add-pricing"
          @click="$emit('open-add-sku')"
        >
          <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
            <line x1="12" y1="5" x2="12" y2="19"/>
            <line x1="5" y1="12" x2="19" y2="12"/>
          </svg>
          <span>Tambah Item SKU</span>
        </button>
      </div>
    </header>

    <!-- CATALOG FULL TABLE WORKSPACE -->
    <div class="catalog-workspace">
      <!-- Search & Meta Strip -->
      <div class="catalog-toolbar">
        <div class="search-input-wrap">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="#64748B" stroke-width="2">
            <circle cx="11" cy="11" r="8"/>
            <line x1="21" y1="21" x2="16.65" y2="16.65"/>
          </svg>
          <input
            type="text"
            :value="search"
            @input="$emit('update:search', ($event.target as HTMLInputElement).value)"
            class="catalog-search-input"
            placeholder="Cari berdasarkan SKU, Nama Komponen, atau Sumber..."
          />
          <button
            v-if="search"
            type="button"
            class="search-clear-btn"
            @click="$emit('update:search', '')"
          >
            <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18"></line>
              <line x1="6" y1="6" x2="18" y2="18"></line>
            </svg>
          </button>
        </div>

        <div class="catalog-stats">
          <span class="stat-tag">
            <strong>{{ items.length }}</strong> Item Ditampilkan
          </span>
          <span class="stat-tag unconfirmed-tag">
            <strong>{{ unconfirmedCount }}</strong> Perlu Konfirmasi STA
          </span>
          <button
            type="button"
            class="refresh-catalog-btn"
            :disabled="isLoading"
            @click="$emit('refresh')"
            title="Muat ulang data dari PostgreSQL"
          >
            <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" :class="{ 'spin-icon': isLoading }">
              <path d="M23 4v6h-6"/>
              <path d="M1 20v-6h6"/>
              <path d="M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15"/>
            </svg>
            <span>Refresh</span>
          </button>
        </div>
      </div>

      <!-- FULL TABLE -->
      <div class="catalog-table-card">
        <table class="catalog-table">
          <thead>
            <tr>
              <th class="th-no">No</th>
              <th>SKU Identifier</th>
              <th>Kategori</th>
              <th>Nama Komponen / Lisensi</th>
              <th>Environment</th>
              <th>Satuan</th>
              <th class="text-right">Harga Satuan (IDR)</th>
              <th class="text-center">Status STA</th>
              <th>Sumber / Dokumen</th>
              <th class="text-center th-action">Aksi</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="items.length === 0">
              <td colspan="10" class="empty-table-cell">
                <div class="empty-table-state">
                  <span class="empty-icon">
                    <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <circle cx="11" cy="11" r="8"></circle>
                      <line x1="21" y1="21" x2="16.65" y2="16.65"></line>
                    </svg>
                  </span>
                  <p>Tidak ada item harga yang sesuai dengan filter atau pencarian.</p>
                  <button
                    type="button"
                    class="btn-reset-filter"
                    @click="$emit('reset-all-filters')"
                  >
                    Reset Semua Filter
                  </button>
                </div>
              </td>
            </tr>

            <tr v-for="(item, idx) in items" :key="item.id">
              <td class="td-no">{{ idx + 1 }}</td>
              <td>
                <code class="sku-badge">{{ item.sku }}</code>
              </td>
              <td>
                <span :class="['category-pill', item.category]">
                  {{ formatCategory(item.category) }}
                </span>
              </td>
              <td class="font-semibold text-navy">
                {{ item.name }}
              </td>
              <td>
                <span :class="['env-badge', item.environment === 'Production' ? 'prod' : 'dev']">
                  {{ item.environment }}
                </span>
              </td>
              <td class="text-slate-600 font-medium">
                {{ item.unit }}
              </td>
              <td class="text-right text-tabular">
                <template v-if="item.unit_price_idr !== null && item.unit_price_idr !== undefined">
                  <BaseMoney :amount="item.unit_price_idr" />
                </template>
                <template v-else>
                  <span class="unconfirmed-text-badge">Pending STA</span>
                </template>
              </td>
              <td class="text-center">
                <span v-if="item.is_confirmed" class="status-badge-confirmed" title="Harga resmi terkonfirmasi">
                  Resmi
                </span>
                <span v-else class="status-badge-pending" title="Estimasi harga belum dikonfirmasi STA">
                  Pending STA
                </span>
              </td>
              <td class="source-cell text-slate-500">
                {{ item.source || '-' }}
              </td>
              <td class="text-center action-cell">
                <button
                  type="button"
                  class="tbl-action-btn btn-edit"
                  @click="$emit('edit-item', item)"
                  title="Edit Item Harga"
                >
                  <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
                    <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
                  </svg>
                  <span>Edit</span>
                </button>
                <button
                  type="button"
                  class="tbl-action-btn btn-delete"
                  @click="$emit('delete-item', item)"
                  title="Hapus Item Harga"
                >
                  <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <polyline points="3 6 5 6 21 6"/>
                    <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
                  </svg>
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import type { PricingItemRecord } from '~/types/sizing'
import BaseMoney from '~/components/base/BaseMoney.vue'

defineProps<{
  items: PricingItemRecord[]
  unconfirmedCount: number
  search: string
  categoryFilter: string
  envFilter: string
  isLoading: boolean
  isSidebarCollapsed: boolean
}>()

defineEmits<{
  (e: 'update:search', val: string): void
  (e: 'update:categoryFilter', val: string): void
  (e: 'update:envFilter', val: string): void
  (e: 'refresh'): void
  (e: 'reset-all-filters'): void
  (e: 'open-add-sku'): void
  (e: 'edit-item', item: PricingItemRecord): void
  (e: 'delete-item', item: PricingItemRecord): void
  (e: 'toggle-sidebar'): void
}>()

const isOpen = ref(false)
const popoverRef = ref<HTMLElement | null>(null)

function getCategoryLabel(cat: string): string {
  switch (cat) {
    case 'compute': return 'Compute'
    case 'storage': return 'Storage'
    case 'backup': return 'Backup'
    case 'license': return 'Lisensi'
    case 'baremetal': return 'Baremetal'
    default: return 'Semua Kategori'
  }
}

function getEnvLabel(env: string): string {
  switch (env) {
    case 'Production': return 'Production'
    case 'Dev / Non-Prod': return 'Dev / Non-Prod'
    default: return 'Semua Env'
  }
}

function formatCategory(cat: string): string {
  switch (cat) {
    case 'compute': return 'Compute'
    case 'storage': return 'Storage'
    case 'backup': return 'Backup'
    case 'baremetal': return 'Baremetal'
    case 'license': return 'Lisensi'
    default: return cat
  }
}

function handleClickOutside(event: MouseEvent) {
  if (isOpen.value && popoverRef.value && !popoverRef.value.contains(event.target as Node)) {
    isOpen.value = false
  }
}

function handleKeydown(event: KeyboardEvent) {
  if (event.key === 'Escape' && isOpen.value) {
    isOpen.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
  document.addEventListener('keydown', handleKeydown)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
  document.removeEventListener('keydown', handleKeydown)
})
</script>

<style scoped>
.catalog-view-container {
  display: flex;
  flex-direction: column;
  flex: 1;
  overflow: hidden;
}

.catalog-topbar {
  background: #FFFFFF;
  border-bottom: 1px solid #E2E8F0;
  padding: 14px 24px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 20px;
  flex-shrink: 0;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.03);
}

.topbar-left {
  display: flex;
  flex-direction: column;
}

.title-with-toggle {
  display: flex;
  align-items: center;
  gap: 12px;
}

.topbar-sidebar-toggle {
  width: 32px;
  height: 32px;
  border-radius: 6px;
  background: #F1F5F9;
  border: 1px solid #CBD5E1;
  color: #475569;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.15s ease;
  flex-shrink: 0;
}

.topbar-sidebar-toggle:hover {
  background: #E2E8F0;
  color: #002060;
  border-color: #94A3B8;
}

.page-title {
  font-size: 18px;
  font-weight: 800;
  color: #002060;
  margin: 0;
  letter-spacing: -0.02em;
}

.page-subtitle {
  font-size: 12px;
  color: #64748B;
  margin: 2px 0 0 0;
}

.topbar-filters {
  display: flex;
  align-items: center;
  gap: 16px;
}

.filter-popover-wrapper {
  position: relative;
}

.minimal-filter-btn {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 6px 14px;
  background: #F8FAFC;
  border: 1px solid #CBD5E1;
  border-radius: 6px;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.15s ease;
  user-select: none;
}

.minimal-filter-btn:hover {
  background: #F1F5F9;
  border-color: #94A3B8;
}

.minimal-filter-btn.active {
  background: #FFFFFF;
  border-color: #002060;
  box-shadow: 0 0 0 3px rgba(0, 32, 96, 0.08);
}

.filter-svg-icon {
  color: #002060;
  flex-shrink: 0;
}

.filter-btn-title {
  font-weight: 700;
  color: #334155;
}

.filter-pill {
  padding: 2px 7px;
  border-radius: 4px;
  font-size: 11px;
  font-weight: 700;
}

.filter-pill.pill-cat {
  background: #E0E7FF;
  color: #3730A3;
  border: 1px solid #C7D2FE;
}

.filter-pill.pill-env {
  background: #DCFCE7;
  color: #166534;
  border: 1px solid #BBF7D0;
}

.filter-pill.pill-all {
  background: #F1F5F9;
  color: #475569;
  border: 1px solid #E2E8F0;
}

.filter-pill-dot {
  color: #94A3B8;
  font-weight: 800;
}

.chevron-svg {
  color: #64748B;
  transition: transform 0.2s cubic-bezier(0.16, 1, 0.3, 1);
}

.chevron-svg.rotated {
  transform: rotate(180deg);
}

.filter-popup-modal {
  position: absolute;
  top: calc(100% + 8px);
  right: 0;
  width: 380px;
  background: #FFFFFF;
  border-radius: 10px;
  border: 1px solid #CBD5E1;
  box-shadow: 0 16px 36px -6px rgba(0, 32, 96, 0.18), 0 6px 14px -3px rgba(0, 32, 96, 0.08);
  z-index: 200;
  overflow: hidden;
  animation: popoverDrop 0.2s cubic-bezier(0.16, 1, 0.3, 1);
}

@keyframes popoverDrop {
  from {
    opacity: 0;
    transform: translateY(-8px) scale(0.98);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

.popover-fade-enter-active,
.popover-fade-leave-active {
  transition: opacity 0.15s ease, transform 0.15s ease;
}

.popover-fade-enter-from,
.popover-fade-leave-to {
  opacity: 0;
  transform: translateY(-6px);
}

.popup-header {
  padding: 12px 16px;
  background: #F8FAFC;
  border-bottom: 1px solid #E2E8F0;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.popup-title-wrap {
  display: flex;
  align-items: center;
  gap: 8px;
}

.popup-heading {
  font-size: 13px;
  font-weight: 800;
  color: #002060;
}

.popup-close-btn {
  background: transparent;
  border: none;
  color: #94A3B8;
  cursor: pointer;
  padding: 4px;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.15s ease;
}

.popup-close-btn:hover {
  background: #E2E8F0;
  color: #0F172A;
}

.popup-body {
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.popup-section {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.section-label-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.section-field-lbl {
  font-size: 12px;
  font-weight: 700;
  color: #1E293B;
}

.section-badge {
  font-size: 11px;
  font-weight: 700;
  padding: 2px 8px;
  border-radius: 12px;
}

.section-badge.badge-cat {
  background: #E0E7FF;
  color: #3730A3;
}

.section-badge.badge-env {
  background: #DCFCE7;
  color: #166534;
}

.popup-segmented {
  display: flex;
  width: 100%;
  background: #F1F5F9;
  padding: 3px;
  border-radius: 6px;
  border: 1px solid #CBD5E1;
  gap: 2px;
}

.popup-segmented .seg-btn {
  flex: 1;
  text-align: center;
  padding: 7px 0;
  font-size: 12px;
  font-weight: 600;
  color: #64748B;
  background: transparent;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.15s ease;
}

.popup-segmented .seg-btn.active {
  background: #FFFFFF;
  color: #002060;
  font-weight: 700;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.08);
}

.popup-segmented.wrap-tabs {
  display: flex;
  flex-wrap: wrap;
  gap: 2px;
}

.popup-segmented.wrap-tabs .seg-btn {
  flex: 1 1 28%;
}

.popup-footer {
  padding: 10px 16px;
  background: #F8FAFC;
  border-top: 1px solid #E2E8F0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.popup-reset-link {
  background: none;
  border: none;
  color: #EF4444;
  font-size: 11px;
  font-weight: 700;
  cursor: pointer;
  padding: 0;
  text-decoration: underline;
}

.popup-reset-link:hover {
  color: #B91C1C;
}

.popup-active-note {
  font-size: 11px;
  color: #64748B;
}

.popup-apply-btn {
  padding: 6px 14px;
  background: #002060;
  color: #FFFFFF;
  border: none;
  border-radius: 5px;
  font-size: 11px;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.15s ease;
}

.popup-apply-btn:hover {
  background: #003399;
}

.btn-add-pricing {
  margin-left: 8px;
}

.action-btn-generate {
  display: inline-flex;
  align-items: center;
  gap: 7px;
  padding: 7px 14px;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.15s ease;
  border: 1px solid transparent;
  background: #002060;
  color: #FFFFFF;
}

.action-btn-generate:hover {
  background: #003399;
  box-shadow: 0 2px 6px rgba(0, 32, 96, 0.25);
}

.catalog-workspace {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: 20px 24px;
  overflow-y: auto;
  gap: 16px;
}

.catalog-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.search-input-wrap {
  display: flex;
  align-items: center;
  background: #FFFFFF;
  border: 1px solid #CBD5E1;
  border-radius: 6px;
  padding: 0 12px;
  gap: 8px;
  width: 380px;
  height: 38px;
  transition: border-color 0.15s ease;
}

.search-input-wrap:focus-within {
  border-color: #002060;
  box-shadow: 0 0 0 2px rgba(0, 32, 96, 0.1);
}

.catalog-search-input {
  border: none;
  background: transparent;
  outline: none;
  font-size: 12px;
  color: #0F172A;
  width: 100%;
}

.search-clear-btn {
  background: transparent;
  border: none;
  color: #94A3B8;
  cursor: pointer;
  font-size: 12px;
  padding: 0;
}

.catalog-stats {
  display: flex;
  align-items: center;
  gap: 10px;
}

.stat-tag {
  font-size: 12px;
  color: #475569;
  background: #F1F5F9;
  border: 1px solid #E2E8F0;
  padding: 4px 10px;
  border-radius: 4px;
}

.unconfirmed-tag {
  background: #FEF3C7;
  color: #92400E;
  border-color: #FDE68A;
}

.refresh-catalog-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  background: #FFFFFF;
  border: 1px solid #CBD5E1;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 600;
  color: #334155;
  cursor: pointer;
  transition: all 0.15s ease;
}

.refresh-catalog-btn:hover:not(:disabled) {
  background: #F8FAFC;
  color: #002060;
}

.spin-icon {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.catalog-table-card {
  background: #FFFFFF;
  border: 1px solid #CBD5E1;
  border-radius: 8px;
  overflow-x: auto;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
}

.catalog-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 12px;
}

.catalog-table th {
  background: #002060;
  color: #FFFFFF;
  padding: 10px 14px;
  font-weight: 700;
  text-align: left;
  white-space: nowrap;
}

.catalog-table td {
  padding: 11px 14px;
  border-bottom: 1px solid #E2E8F0;
  color: #1E293B;
  vertical-align: middle;
}

.catalog-table tr:hover td {
  background: #F8FAFC;
}

.th-no,
.td-no {
  width: 36px;
  text-align: center;
  color: #64748B;
}

.text-right {
  text-align: right;
}

.text-center {
  text-align: center;
}

.text-tabular {
  font-variant-numeric: tabular-nums;
}

.font-semibold {
  font-weight: 600;
}

.text-navy {
  color: #002060;
}

.sku-badge {
  font-family: 'JetBrains Mono', monospace;
  font-size: 11px;
  font-weight: 700;
  color: #002060;
  background: #EFF6FF;
  border: 1px solid #BFDBFE;
  padding: 2px 6px;
  border-radius: 4px;
}

.category-pill {
  font-size: 10px;
  font-weight: 700;
  padding: 2px 8px;
  border-radius: 4px;
  text-transform: uppercase;
  letter-spacing: 0.03em;
}

.category-pill.compute {
  background: #E0F2FE;
  color: #0369A1;
}

.category-pill.storage {
  background: #F1F5F9;
  color: #475569;
}

.category-pill.backup {
  background: #FEF3C7;
  color: #B45309;
}

.category-pill.baremetal {
  background: var(--pill-neutral-bg, #ebedee);
  color: var(--accent-teal, #1b6c80);
}

.category-pill.license {
  background: #DCFCE7;
  color: #15803D;
}

.env-badge {
  font-size: 11px;
  font-weight: 600;
  padding: 2px 6px;
  border-radius: 3px;
}

.env-badge.prod {
  background: #EFF6FF;
  color: #1D4ED8;
}

.env-badge.dev {
  background: #F3F4F6;
  color: #4B5563;
}

.unconfirmed-text-badge {
  font-size: 11px;
  font-weight: 700;
  color: #B45309;
  background: #FEF3C7;
  border: 1px solid #FDE68A;
  padding: 2px 6px;
  border-radius: 4px;
}

.status-badge-confirmed {
  font-size: 11px;
  font-weight: 700;
  color: #15803D;
  background: #DCFCE7;
  padding: 2px 8px;
  border-radius: 4px;
}

.status-badge-pending {
  font-size: 11px;
  font-weight: 700;
  color: #B45309;
  background: #FEF3C7;
  padding: 2px 8px;
  border-radius: 4px;
}

.source-cell {
  max-width: 200px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  font-size: 11px;
}

.th-action {
  width: 120px;
}

.action-cell {
  white-space: nowrap;
}

.tbl-action-btn {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 11px;
  font-weight: 600;
  cursor: pointer;
  border: 1px solid #CBD5E1;
  background: #FFFFFF;
  color: #334155;
  margin: 0 2px;
  transition: all 0.15s ease;
}

.tbl-action-btn.btn-edit:hover {
  background: #EFF6FF;
  color: #002060;
  border-color: #002060;
}

.tbl-action-btn.btn-delete:hover {
  background: #FEE2E2;
  color: #DC2626;
  border-color: #F87171;
}

.empty-table-cell {
  text-align: center;
  padding: 40px !important;
}

.empty-table-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
}

.empty-icon {
  font-size: 32px;
}

.btn-reset-filter {
  padding: 6px 14px;
  background: #002060;
  color: #FFFFFF;
  border: none;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 600;
  cursor: pointer;
}
</style>
