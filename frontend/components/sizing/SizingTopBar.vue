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

.filter-pill.pill-crit {
  background: #FEF3C7;
  color: #92400E;
  border: 1px solid #FCD34D;
}

.filter-pill.pill-mod {
  background: #E0F2FE;
  color: #0369A1;
  border: 1px solid #BAE6FD;
}

.filter-pill.pill-lic {
  background: #EFF6FF;
  color: #1D4ED8;
  border: 1px solid #BFDBFE;
}

.filter-pill.pill-free {
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

.section-badge.badge-crit {
  background: #FEF3C7;
  color: #92400E;
}

.section-badge.badge-mod {
  background: #E0F2FE;
  color: #0369A1;
}

.section-badge.badge-lic {
  background: #EFF6FF;
  color: #1D4ED8;
}

.section-badge.badge-free {
  background: #F1F5F9;
  color: #475569;
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

.seg-btn.btn-crit.active {
  color: #B45309;
}

.seg-btn.btn-mod.active {
  color: #0284C7;
}

.seg-btn.btn-lic.active {
  color: #1D4ED8;
}

.section-hint {
  font-size: 11px;
  color: #64748B;
  margin: 2px 0 0 0;
  line-height: 1.4;
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

.popup-active-note {
  font-size: 11px;
  color: #64748B;
}

.popup-active-note strong {
  color: #002060;
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

.topbar-buttons {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-left: 8px;
}

.action-btn-upload,
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
}

.action-btn-upload {
  background: #F8FAFC;
  color: #334155;
  border-color: #CBD5E1;
}

.action-btn-upload:hover {
  background: #F1F5F9;
  color: #002060;
  border-color: #94A3B8;
}

.action-btn-generate {
  background: #002060;
  color: #FFFFFF;
}

.action-btn-generate:hover {
  background: #003399;
  box-shadow: 0 2px 6px rgba(0, 32, 96, 0.25);
}
</style>
