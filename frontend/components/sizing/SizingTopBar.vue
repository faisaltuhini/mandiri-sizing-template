<template>
  <header class="sketch-topbar">
    <div class="topbar-left">
      <div class="title-with-toggle">
        <button
          type="button"
          class="topbar-sidebar-toggle"
          @click="$emit('toggle-sidebar')"
          :title="isSidebarCollapsed ? 'Perluas Sidebar' : 'Minimize Sidebar'"
        >
          <svg width="17" height="17" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <rect x="3" y="3" width="18" height="18" rx="2" ry="2"></rect>
            <line x1="9" y1="3" x2="9" y2="21"></line>
          </svg>
        </button>
        <h1 class="page-title">Sizing Configuration &amp; Live Adjuster</h1>
      </div>
      <p class="page-subtitle">Sesuaikan alokasi compute, memory, storage &amp; lisensi per layer arsitektur</p>
    </div>

    <!-- TOP BAR ACTIONS: 1 MINIMALIST FILTER BUTTON WITH FLOATING POPUP -->
    <div class="topbar-filters">
      <div class="filter-popover-wrapper" ref="popoverRef">
        <button
          type="button"
          class="minimal-filter-btn"
          :class="{ active: isOpen }"
          @click="isOpen = !isOpen"
          title="Buka pengaturan filter Tier &amp; Lisensi"
        >
          <svg class="filter-svg-icon" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <polygon points="22 3 2 3 10 12.46 10 19 14 21 14 12.46 22 3"></polygon>
          </svg>
          <span class="filter-btn-title">Filter:</span>
          <span :class="['filter-pill', tier === 'critical' ? 'pill-crit' : 'pill-mod']">
            {{ tier === 'critical' ? 'Critical' : 'Moderat' }}
          </span>
          <span class="filter-pill-dot">·</span>
          <span :class="['filter-pill', license === 'license' ? 'pill-lic' : 'pill-free']">
            {{ license === 'license' ? 'License' : 'Free' }}
          </span>
          <svg class="chevron-svg" :class="{ rotated: isOpen }" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
            <polyline points="6 9 12 15 18 9"></polyline>
          </svg>
        </button>

        <!-- FLOATING POPUP / POPOVER -->
        <transition name="popover-fade">
          <div v-if="isOpen" class="filter-popup-modal" @click.stop>
            <div class="popup-header">
              <div class="popup-title-wrap">
                <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="#002060" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round">
                  <polygon points="22 3 2 3 10 12.46 10 19 14 21 14 12.46 22 3"></polygon>
                </svg>
                <span class="popup-heading">Filter Parameter Arsitektur</span>
              </div>
              <button type="button" class="popup-close-btn" @click="isOpen = false" title="Tutup">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <line x1="18" y1="6" x2="6" y2="18"></line>
                  <line x1="6" y1="6" x2="18" y2="18"></line>
                </svg>
              </button>
            </div>

            <div class="popup-body">
              <!-- Section 1: Tier Arsitektur -->
              <div class="popup-section">
                <div class="section-label-row">
                  <label class="section-field-lbl">Target Kesiapan (Tier):</label>
                  <span :class="['section-badge', tier === 'critical' ? 'badge-crit' : 'badge-mod']">
                    {{ tier === 'critical' ? 'Tier Critical (2N)' : 'Tier Moderat (N+1)' }}
                  </span>
                </div>
                <div class="segmented-control popup-segmented">
                  <button
                    type="button"
                    :class="['seg-btn', 'btn-crit', { active: tier === 'critical' }]"
                    @click="$emit('update:tier', 'critical')"
                  >
                    Critical
                  </button>
                  <button
                    type="button"
                    :class="['seg-btn', 'btn-mod', { active: tier === 'moderate' }]"
                    @click="$emit('update:tier', 'moderate')"
                  >
                    Moderat
                  </button>
                </div>
                <p class="section-hint">
                  {{ tier === 'critical' 
                    ? 'Direkomendasikan untuk Core Banking Livin\' & transaksi tinggi (2N redundancy, HA SLA 99.99%).' 
                    : 'Sesuai untuk backoffice staging, batch processing, dan service internal non-kritis.' }}
                </p>
              </div>

              <!-- Section 2: Model Lisensi -->
              <div class="popup-section">
                <div class="section-label-row">
                  <label class="section-field-lbl">Model Lisensi Software:</label>
                  <span :class="['section-badge', license === 'license' ? 'badge-lic' : 'badge-free']">
                    {{ license === 'license' ? 'Enterprise Support' : 'Community OSS' }}
                  </span>
                </div>
                <div class="segmented-control popup-segmented">
                  <button
                    type="button"
                    :class="['seg-btn', { active: license === 'free' }]"
                    @click="$emit('update:license', 'free')"
                  >
                    Free (OSS)
                  </button>
                  <button
                    type="button"
                    :class="['seg-btn', 'btn-lic', { active: license === 'license' }]"
                    @click="$emit('update:license', 'license')"
                  >
                    License (Enterprise)
                  </button>
                </div>
                <p class="section-hint">
                  {{ license === 'license'
                    ? 'Menyertakan estimasi biaya lisensi resmi vendor (Red Hat OCP, F5 NGINX Plus, Oracle DB).'
                    : 'Versi Community OSS tanpa lisensi tahunan (biaya software Rp 0).' }}
                </p>
              </div>
            </div>

            <div class="popup-footer">
              <span class="popup-active-note">
                Aktif: <strong>{{ tier === 'critical' ? 'Critical' : 'Moderat' }}</strong> · <strong>{{ license === 'license' ? 'License' : 'Free' }}</strong>
              </span>
              <button type="button" class="popup-apply-btn" @click="isOpen = false">
                Terapkan &amp; Tutup
              </button>
            </div>
          </div>
        </transition>
      </div>

      <!-- Action Buttons -->
      <div class="topbar-buttons">
        <button
          type="button"
          class="action-btn-upload"
          @click="$emit('open-upload')"
          title="Unggah template sizing atau data telemetri"
        >
          <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
            <polyline points="17 8 12 3 7 8"/>
            <line x1="12" y1="3" x2="12" y2="15"/>
          </svg>
          <span>Upload</span>
        </button>

        <button
          type="button"
          class="action-btn-generate"
          @click="$emit('trigger-generate')"
          title="Generate hasil perhitungan sizing dan BOQ"
        >
          <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2"/>
          </svg>
          <span>Generate</span>
        </button>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'

defineProps<{
  tier: 'critical' | 'moderate'
  license: 'free' | 'license'
  isSidebarCollapsed: boolean
}>()

defineEmits<{
  (e: 'update:tier', val: 'critical' | 'moderate'): void
  (e: 'update:license', val: 'free' | 'license'): void
  (e: 'toggle-sidebar'): void
  (e: 'open-upload'): void
  (e: 'trigger-generate'): void
}>()

const isOpen = ref(false)
const popoverRef = ref<HTMLElement | null>(null)

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
.sketch-topbar {
  background: var(--surface, #ffffff);
  border-bottom: 1px solid var(--rule, #edf1f4);
  padding: 12px 24px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 20px;
  flex-shrink: 0;
  box-shadow: none;
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
  width: 30px;
  height: 30px;
  border-radius: var(--radius-control, 6px);
  background: var(--raised, #eef2f5);
  border: 1px solid var(--border, #e2e7eb);
  color: var(--muted, #5d6165);
  display: inline-flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: background-color 0.15s ease, color 0.15s ease, border-color 0.15s ease;
  flex-shrink: 0;
}

.topbar-sidebar-toggle:hover {
  background: var(--active, #e3e8ec);
  color: var(--text, #131416);
  border-color: var(--border2, #d1d6db);
}

.page-title {
  font-size: var(--fs-title, 20px);
  font-weight: 600;
  color: var(--text, #131416);
  margin: 0;
  line-height: 1.25;
}

.page-subtitle {
  font-size: 13px;
  color: var(--muted, #5d6165);
  margin: 2px 0 0 0;
}

.topbar-filters {
  display: flex;
  align-items: center;
  gap: 14px;
}

.filter-popover-wrapper {
  position: relative;
}

.minimal-filter-btn {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 6px 12px;
  background: var(--surface, #ffffff);
  border: 1px solid var(--border, #e2e7eb);
  border-radius: var(--radius-control, 6px);
  font-size: 12px;
  cursor: pointer;
  transition: border-color 0.15s ease, background-color 0.15s ease;
  user-select: none;
  color: var(--text, #131416);
}

.minimal-filter-btn:hover {
  background: var(--raised, #eef2f5);
  border-color: var(--border2, #d1d6db);
}

.minimal-filter-btn.active {
  background: var(--surface, #ffffff);
  border-color: var(--accent, #003b73);
}

.filter-svg-icon {
  color: var(--accent, #003b73);
  flex-shrink: 0;
}

.filter-btn-title {
  font-weight: 500;
  color: var(--muted, #5d6165);
}

.filter-pill {
  padding: 2px 7px;
  border-radius: var(--radius-sm, 4px);
  font-size: 11px;
  font-weight: 600;
}

.filter-pill.pill-crit {
  background: var(--pill-warning-bg, #fffbeb);
  color: var(--pill-warning-ink, #92400e);
  border: 1px solid rgba(146, 64, 14, 0.2);
}

.filter-pill.pill-mod {
  background: var(--pill-neutral-bg, #ebedee);
  color: var(--pill-neutral-ink, #353736);
  border: 1px solid var(--border, #e2e7eb);
}

.filter-pill.pill-lic {
  background: var(--pill-positive-bg, #dcfae5);
  color: var(--pill-positive-ink, #04250e);
  border: 1px solid rgba(30, 123, 79, 0.2);
}

.filter-pill.pill-free {
  background: var(--pill-neutral-bg, #ebedee);
  color: var(--pill-neutral-ink, #353736);
  border: 1px solid var(--border, #e2e7eb);
}

.filter-pill-dot {
  color: var(--border2, #d1d6db);
  font-weight: 600;
}

.chevron-svg {
  color: var(--muted, #5d6165);
  transition: transform 0.2s ease;
}

.chevron-svg.rotated {
  transform: rotate(180deg);
}

.filter-popup-modal {
  position: absolute;
  top: calc(100% + 6px);
  right: 0;
  width: 360px;
  background: var(--surface, #ffffff);
  border-radius: var(--radius-card, 8px);
  border: 1px solid var(--border, #e2e7eb);
  box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.08);
  z-index: 200;
  overflow: hidden;
}

.popup-header {
  padding: 12px 16px;
  background: var(--raised, #eef2f5);
  border-bottom: 1px solid var(--border, #e2e7eb);
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
  font-weight: 600;
  color: var(--text, #131416);
}

.popup-close-btn {
  background: transparent;
  border: none;
  color: var(--muted, #5d6165);
  cursor: pointer;
  padding: 4px;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: color 0.15s ease, background-color 0.15s ease;
}

.popup-close-btn:hover {
  background: var(--border, #e2e7eb);
  color: var(--text, #131416);
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
  font-weight: 600;
  color: var(--text, #131416);
}

.section-badge {
  font-size: 11px;
  font-weight: 600;
  padding: 2px 8px;
  border-radius: 12px;
}

.section-badge.badge-crit {
  background: var(--pill-warning-bg, #fffbeb);
  color: var(--pill-warning-ink, #92400e);
  border: 1px solid rgba(146, 64, 14, 0.2);
}

.section-badge.badge-mod {
  background: var(--pill-neutral-bg, #ebedee);
  color: var(--pill-neutral-ink, #353736);
  border: 1px solid var(--border, #e2e7eb);
}

.section-badge.badge-lic {
  background: var(--pill-positive-bg, #dcfae5);
  color: var(--pill-positive-ink, #04250e);
  border: 1px solid rgba(30, 123, 79, 0.2);
}

.section-badge.badge-free {
  background: var(--pill-neutral-bg, #ebedee);
  color: var(--pill-neutral-ink, #353736);
  border: 1px solid var(--border, #e2e7eb);
}

.popup-segmented {
  display: flex;
  width: 100%;
  background: var(--raised, #eef2f5);
  padding: 2px;
  border-radius: var(--radius-control, 6px);
  border: 1px solid var(--border, #e2e7eb);
  gap: 2px;
}

.popup-segmented .seg-btn {
  flex: 1;
  text-align: center;
  padding: 6px 0;
  font-size: 12px;
  font-weight: 500;
  color: var(--muted, #5d6165);
  background: transparent;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.15s ease, color 0.15s ease;
}

.popup-segmented .seg-btn.active {
  background: var(--surface, #ffffff);
  color: var(--text, #131416);
  font-weight: 600;
  border: 1px solid var(--border, #e2e7eb);
}

.section-hint {
  font-size: 11px;
  color: var(--muted, #5d6165);
  margin: 2px 0 0 0;
  line-height: 1.4;
}

.popup-footer {
  padding: 10px 16px;
  background: var(--raised, #eef2f5);
  border-top: 1px solid var(--border, #e2e7eb);
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.popup-active-note {
  font-size: 11px;
  color: var(--muted, #5d6165);
}

.popup-active-note strong {
  color: var(--text, #131416);
}

.popup-apply-btn {
  padding: 6px 14px;
  background: var(--accent, #003b73);
  color: #FFFFFF;
  border: none;
  border-radius: var(--radius-control, 6px);
  font-size: 12px;
  font-weight: 600;
  cursor: pointer;
  transition: opacity 0.15s ease;
}

.popup-apply-btn:hover {
  opacity: 0.9;
}

.topbar-buttons {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-left: 8px;
}

.action-btn-upload,
.action-btn-generate {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  border-radius: var(--radius-control, 6px);
  font-size: 12px;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.15s ease, border-color 0.15s ease, opacity 0.15s ease;
  border: 1px solid transparent;
}

.action-btn-upload {
  background: var(--surface, #ffffff);
  color: var(--text, #131416);
  border-color: var(--border, #e2e7eb);
}

.action-btn-upload:hover {
  background: var(--raised, #eef2f5);
  border-color: var(--border2, #d1d6db);
}

.action-btn-generate {
  background: var(--accent, #003b73);
  color: #FFFFFF;
}

.action-btn-generate:hover {
  opacity: 0.9;
  box-shadow: none;
}
</style>
