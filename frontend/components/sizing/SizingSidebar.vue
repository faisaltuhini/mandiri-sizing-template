<template>
  <aside :class="['sketch-sidebar', { collapsed: isCollapsed }]">
    <div class="sidebar-brand">
      <div
        class="brand-shield"
        @click="$emit('toggle-collapse')"
        :title="isCollapsed ? 'Klik untuk memperluas sidebar' : 'Enterprise Architecture Sizing Suite'"
      >
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
          <path d="M12 2L3 7V12C3 17.5 7 21.5 12 22C17 21.5 21 17.5 21 12V7L12 2Z" fill="#D4A017" fill-opacity="0.25" stroke="#D4A017" stroke-width="2"/>
          <path d="M12 6L7 11H10V17H14V11H17L12 6Z" fill="#D4A017"/>
        </svg>
      </div>
      <div v-if="!isCollapsed" class="brand-text">
        <span class="bank-lbl">ENTERPRISE</span>
        <span class="app-lbl">Sizing Suite</span>
      </div>
    </div>

    <nav class="sidebar-menu">
      <!-- Menu 1: Price Catalog (Full Table CRUD) -->
      <button
        type="button"
        :class="['menu-btn', { active: currentMenu === 'catalog' }]"
        @click="$emit('select-menu', 'catalog')"
        :title="isCollapsed ? 'Katalog: Master Harga & CRUD' : 'Katalog Harga Resmi & Master Data SKU'"
      >
        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <rect x="3" y="3" width="7" height="7"/>
          <rect x="14" y="3" width="7" height="7"/>
          <rect x="14" y="14" width="7" height="7"/>
          <rect x="3" y="14" width="7" height="7"/>
        </svg>
        <div v-if="!isCollapsed" class="menu-btn-text">
          <span>Katalog</span>
          <small class="menu-btn-sub">Master Harga &amp; CRUD</small>
        </div>
      </button>

      <!-- Menu 2: Sizing Configurator -->
      <button
        type="button"
        :class="['menu-btn', { active: currentMenu === 'sizing' }]"
        @click="$emit('select-menu', 'sizing')"
        :title="isCollapsed ? 'Sizing: Configurator & BOQ' : 'Kalkulator Sizing Arsitektur & Penyesuaian Kapasitas'"
      >
        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <line x1="4" y1="21" x2="4" y2="14"/>
          <line x1="4" y1="10" x2="4" y2="3"/>
          <line x1="12" y1="21" x2="12" y2="12"/>
          <line x1="12" y1="8" x2="12" y2="3"/>
          <line x1="20" y1="21" x2="20" y2="16"/>
          <line x1="20" y1="12" x2="20" y2="3"/>
          <line x1="1" y1="14" x2="7" y2="14"/>
          <line x1="9" y1="8" x2="15" y2="8"/>
          <line x1="17" y1="16" x2="23" y2="16"/>
        </svg>
        <div v-if="!isCollapsed" class="menu-btn-text">
          <span>Sizing</span>
          <small class="menu-btn-sub">Configurator &amp; BOQ</small>
        </div>
      </button>

      <!-- Menu 3: Architecture Preset & Formulas -->
      <button
        type="button"
        :class="['menu-btn', { active: currentMenu === 'presets' }]"
        @click="$emit('select-menu', 'presets')"
        :title="isCollapsed ? (newCatalogCount && newCatalogCount > 0 ? `Preset: ${newCatalogCount} katalog baru` : 'Preset: Formula Moderat & Critical') : 'Preset Standar Arsitektur (Formula & Baseline)'"
      >
        <div class="menu-icon-wrap">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
            <polyline points="14 2 14 8 20 8"></polyline>
            <line x1="16" y1="13" x2="8" y2="13"></line>
            <line x1="16" y1="17" x2="8" y2="17"></line>
            <polyline points="10 9 9 9 8 9"></polyline>
          </svg>
          <span v-if="isCollapsed && newCatalogCount && newCatalogCount > 0" class="collapsed-dot-alert"></span>
        </div>
        <div v-if="!isCollapsed" class="menu-btn-text">
          <div class="preset-label-row">
            <span>Preset</span>
            <span v-if="newCatalogCount && newCatalogCount > 0" class="badge-new-catalog" title="Ada katalog baru yang belum masuk preset">
              {{ newCatalogCount }} Baru
            </span>
          </div>
          <small class="menu-btn-sub">Formula Baseline &amp; Tier</small>
        </div>
      </button>

      <!-- Menu 4: Software License Rules & Formulas -->
      <button
        type="button"
        :class="['menu-btn', { active: currentMenu === 'licenses' }]"
        @click="$emit('select-menu', 'licenses')"
        :title="isCollapsed ? 'Lisensi: Aturan & Metrik Software' : 'Aturan Lisensi Software (Pod, Shard, Worker, Core)'"
      >
        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"></path>
          <path d="M9 12l2 2 4-4"></path>
        </svg>
        <div v-if="!isCollapsed" class="menu-btn-text">
          <span>Lisensi</span>
          <small class="menu-btn-sub">Aturan &amp; Metrik CRUD</small>
        </div>
      </button>
    </nav>

    <!-- Sidebar Bottom Collapse Action -->
    <div class="sidebar-collapse-bar">
      <button
        type="button"
        class="collapse-action-btn"
        @click="$emit('toggle-collapse')"
        :title="isCollapsed ? 'Perluas Sidebar' : 'Kecilkan / Minimize Sidebar'"
      >
        <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <polyline v-if="!isCollapsed" points="11 17 6 12 11 7"></polyline>
          <polyline v-if="!isCollapsed" points="18 17 13 12 18 7"></polyline>
          <polyline v-if="isCollapsed" points="13 17 18 12 13 7"></polyline>
          <polyline v-if="isCollapsed" points="6 17 11 12 6 7"></polyline>
        </svg>
        <span v-if="!isCollapsed">Minimize Sidebar</span>
      </button>
    </div>

    <div class="sidebar-footer">
      <span class="sys-status" :title="'Engine Ready · :8089'">
        <span class="status-pulse"></span>
        <span v-if="!isCollapsed" class="status-text">Engine Ready · :8089</span>
      </span>
    </div>
  </aside>
</template>

<script setup lang="ts">
import type { AppMenu } from '~/types/sizing'

defineProps<{
  currentMenu: AppMenu
  isCollapsed: boolean
  newCatalogCount?: number
}>()

defineEmits<{
  (e: 'select-menu', menu: AppMenu): void
  (e: 'toggle-collapse'): void
}>()
</script>

<style scoped>
.sketch-sidebar {
  width: 220px;
  background: #002060;
  color: #FFFFFF;
  display: flex;
  flex-direction: column;
  border-right: 1px solid rgba(255, 255, 255, 0.08);
  box-shadow: 2px 0 10px rgba(0, 32, 96, 0.1);
  flex-shrink: 0;
  transition: width 0.25s cubic-bezier(0.4, 0, 0.2, 1);
  overflow: hidden;
}

.sketch-sidebar.collapsed {
  width: 68px;
}

.sidebar-brand {
  padding: 18px 14px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  min-height: 72px;
  box-sizing: border-box;
}

.sketch-sidebar.collapsed .sidebar-brand {
  padding: 16px 8px;
  justify-content: center;
  flex-direction: column;
  gap: 8px;
}

.brand-shield {
  width: 36px;
  height: 36px;
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.08);
  border: 1px solid rgba(212, 160, 23, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  cursor: pointer;
  transition: transform 0.15s ease;
}

.brand-shield:hover {
  transform: scale(1.05);
}

.brand-text {
  display: flex;
  flex-direction: column;
  flex: 1;
  overflow: hidden;
  white-space: nowrap;
}

.bank-lbl {
  font-size: 13px;
  font-weight: 800;
  letter-spacing: 0.5px;
  color: #FFFFFF;
}

.app-lbl {
  font-size: 11px;
  color: #D4A017;
  font-weight: 600;
}

.sidebar-toggle-btn {
  width: 28px;
  height: 28px;
  border-radius: 5px;
  background: rgba(255, 255, 255, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.15);
  color: #CBD5E1;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.15s ease;
  flex-shrink: 0;
}

.sidebar-toggle-btn:hover {
  background: rgba(255, 255, 255, 0.18);
  color: #FFFFFF;
  border-color: rgba(255, 255, 255, 0.3);
}

.sketch-sidebar.collapsed .sidebar-toggle-btn {
  width: 24px;
  height: 24px;
}

.sidebar-menu {
  padding: 16px 10px;
  display: flex;
  flex-direction: column;
  gap: 6px;
  flex: 1;
}

.sketch-sidebar.collapsed .sidebar-menu {
  padding: 16px 6px;
}

.menu-btn {
  display: flex;
  align-items: center;
  gap: 12px;
  width: 100%;
  padding: 10px 14px;
  background: transparent;
  color: #CBD5E1;
  border: 1px solid transparent;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s ease;
  text-align: left;
}

.sketch-sidebar.collapsed .menu-btn {
  padding: 12px 0;
  justify-content: center;
  gap: 0;
}

.menu-btn:hover {
  background: rgba(255, 255, 255, 0.08);
  color: #FFFFFF;
}

.menu-btn.active {
  background: #FFFFFF;
  color: #002060;
  border-color: #FFFFFF;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.12);
}

.menu-btn-text {
  display: flex;
  flex-direction: column;
  overflow: hidden;
  white-space: nowrap;
}

.menu-btn-text span {
  font-size: 13px;
  font-weight: 700;
}

.menu-icon-wrap {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
}

.collapsed-dot-alert {
  position: absolute;
  top: -2px;
  right: -2px;
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: #F59E0B;
  box-shadow: 0 0 0 2px #002060;
}

.preset-label-row {
  display: flex;
  align-items: center;
  gap: 6px;
}

.badge-new-catalog {
  display: inline-flex;
  align-items: center;
  padding: 1px 6px;
  background: #D97706;
  color: #FFFFFF;
  border-radius: 10px;
  font-size: 10px;
  font-weight: 800;
  letter-spacing: 0.02em;
  animation: pulse-badge 2.4s infinite;
}

@keyframes pulse-badge {
  0%, 100% {
    transform: scale(1);
  }
  50% {
    transform: scale(1.05);
  }
}

.menu-btn-sub {
  font-size: 10px;
  color: #94A3B8;
  font-weight: 500;
}

.menu-btn.active .menu-btn-sub {
  color: #64748B;
}

.sidebar-collapse-bar {
  padding: 8px 10px;
  border-top: 1px solid rgba(255, 255, 255, 0.08);
}

.sketch-sidebar.collapsed .sidebar-collapse-bar {
  padding: 8px 6px;
}

.collapse-action-btn {
  display: flex;
  align-items: center;
  justify-content: flex-start;
  gap: 10px;
  width: 100%;
  padding: 8px 12px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  color: #94A3B8;
  border-radius: 6px;
  font-size: 11px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.15s ease;
}

.sketch-sidebar.collapsed .collapse-action-btn {
  justify-content: center;
  padding: 8px 0;
}

.collapse-action-btn:hover {
  background: rgba(255, 255, 255, 0.12);
  color: #FFFFFF;
}

.sidebar-footer {
  padding: 14px 18px;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}

.sketch-sidebar.collapsed .sidebar-footer {
  padding: 14px 8px;
  display: flex;
  justify-content: center;
}

.sys-status {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #94A3B8;
  font-size: 11px;
  font-weight: 600;
  white-space: nowrap;
}

.sketch-sidebar.collapsed .sys-status {
  justify-content: center;
}

.status-pulse {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #10B981;
  box-shadow: 0 0 6px #10B981;
}
</style>
