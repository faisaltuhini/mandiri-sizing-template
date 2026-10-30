<template>
  <div class="preset-view-container">
    <!-- TOPBAR -->
    <header class="sketch-topbar preset-topbar">
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
          <h1 class="page-title">Preset Arsitektur (Formula &amp; Standar Baseline)</h1>
        </div>
        <p class="page-subtitle">Kelola template formula sizing default untuk Tier Moderat &amp; Critical, serta tambahkan layer baru</p>
      </div>

      <div class="topbar-actions">
        <button
          v-if="(unregisteredCatalogItems?.length || 0) > 0 && isBannerDismissed"
          type="button"
          class="btn-unmapped-trigger"
          @click="openCatalogPickerModal"
          title="Lihat daftar komponen katalog yang belum memiliki preset arsitektur"
        >
          <span class="unmapped-dot"></span>
          <span>Katalog Belum Dipetakan</span>
          <span class="unmapped-count">{{ unregisteredCatalogItems.length }}</span>
        </button>

        <button
          type="button"
          class="btn-reset-defaults"
          @click="$emit('reset-to-official-defaults')"
          title="Reset semua preset ke standar baseline resmi Enterprise Architecture"
        >
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M3 12a9 9 0 1 0 9-9 9.75 9.75 0 0 0-6.74 2.74L3 8"/>
            <path d="M3 3v5h5"/>
          </svg>
          <span>Kembalikan Default Pabrik</span>
        </button>


        <button
          type="button"
          class="btn-add-layer"
          @click="openAddLayerModal"
        >
          <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
            <line x1="12" y1="5" x2="12" y2="19"/>
            <line x1="5" y1="12" x2="19" y2="12"/>
          </svg>
          <span>Tambah Layer Baru</span>
        </button>
      </div>
    </header>

    <!-- CONTENT WORKSPACE -->
    <div class="preset-workspace">
      <!-- COMPACT SINGLE-LINE CATALOG SYNC BAR -->
      <transition name="banner-slide">
        <div
          v-if="(unregisteredCatalogItems?.length || 0) > 0 && !isBannerDismissed"
          class="catalog-sync-strip"
        >
          <div class="sync-strip-left">
            <div class="sync-status-icon">
              <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <polygon points="12 2 2 7 12 12 22 7 12 2"></polygon>
                <polyline points="2 17 12 22 22 17"></polyline>
                <polyline points="2 12 12 17 22 12"></polyline>
              </svg>
            </div>
            <div class="sync-strip-copy">
              <span class="sync-strip-title">
                <strong>{{ unregisteredCatalogItems.length }} komponen katalog</strong> belum dipetakan ke layer preset
              </span>
              <span class="sync-strip-sep">·</span>
              <span class="sync-strip-preview">
                {{ unregisteredPreviewText }}
              </span>
            </div>
          </div>

          <div class="sync-strip-actions">
            <button
              type="button"
              class="btn-sync-review"
              @click="openCatalogPickerModal"
            >
              <span>Tinjau &amp; Pilih Katalog</span>
              <span class="review-count-pill">{{ unregisteredCatalogItems.length }}</span>
              <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round">
                <polyline points="9 18 15 12 9 6"></polyline>
              </svg>
            </button>
            <button
              type="button"
              class="btn-sync-close"
              @click="dismissNewCatalogBanner"
              title="Sembunyikan bar notifikasi ini"
              aria-label="Tutup notifikasi"
            >
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <line x1="18" y1="6" x2="6" y2="18"></line>
                <line x1="6" y1="6" x2="18" y2="18"></line>
              </svg>
            </button>
          </div>
        </div>
      </transition>

      <!-- EXECUTIVE GUIDANCE BANNER (MODERAT VS CRITICAL) -->
      <section class="preset-overview-section">
        <div class="overview-header-row">
          <div class="overview-title-group">
            <span class="overview-eyebrow">PANDUAN EKSEKUTIF</span>
            <h2 class="overview-heading">Standar Baseline Arsitektur Mandiri</h2>
          </div>
          <button
            type="button"
            class="btn-toggle-overview"
            @click="toggleOverviewBanner"
            :title="isOverviewBannerVisible ? 'Sembunyikan panduan ringkasan' : 'Tampilkan panduan ringkasan'"
          >
            <svg
              width="14"
              height="14"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              :class="{ 'rotate-180': !isOverviewBannerVisible }"
              style="transition: transform 0.2s ease;"
            >
              <polyline points="18 15 12 9 6 15"></polyline>
            </svg>
            <span>{{ isOverviewBannerVisible ? 'Sembunyikan Panduan' : 'Tampilkan Panduan' }}</span>
          </button>
        </div>

        <transition name="overview-collapse">
          <div v-if="isOverviewBannerVisible" class="overview-cards-row">
            <!-- MODERATE EXECUTIVE CARD -->
            <div class="overview-tier-card card-moderat">
              <div class="card-tier-top">
                <div class="tier-pill-group">
                  <span class="overview-tier-badge mod">MODERAT (N+1)</span>
                  <span class="overview-target-badge">Aplikasi Internal &amp; Operasional Reguler</span>
                </div>
              </div>
              <p class="overview-tier-desc">
                Dirancang untuk sistem non-kritis dengan efisiensi biaya optimal. Menggunakan klaster <strong>Server Virtual (VM)</strong> hemat daya, ketahanan failover <strong>N+1</strong> (1 server cadangan), dan opsi software <strong>Open Source (Gratis)</strong> tanpa lisensi tahunan berulang.
              </p>
              <div class="overview-totals-grid">
                <div class="overview-total-item">
                  <span class="tot-label">Total Unit Server</span>
                  <span class="tot-val text-navy">{{ clusterTotals.modNodes }} Node</span>
                </div>
                <div class="overview-total-item">
                  <span class="tot-label">Tenaga Komputasi</span>
                  <span class="tot-val text-navy">{{ clusterTotals.modVcpu }} vCPU</span>
                </div>
                <div class="overview-total-item">
                  <span class="tot-label">Total RAM</span>
                  <span class="tot-val text-navy">{{ clusterTotals.modRam }} GB</span>
                </div>
                <div class="overview-total-item">
                  <span class="tot-label">Kapasitas Disk</span>
                  <span class="tot-val text-navy">{{ formatDiskSize(clusterTotals.modStorage) }}</span>
                </div>
              </div>
            </div>

            <!-- CRITICAL EXECUTIVE CARD -->
            <div class="overview-tier-card card-critical">
              <div class="card-tier-top">
                <div class="tier-pill-group">
                  <span class="overview-tier-badge crit">CRITICAL (2N)</span>
                  <span class="overview-target-badge">Layanan Nasabah 24/7 &amp; Transaksi Finansial (Core)</span>
                </div>
              </div>
              <p class="overview-tier-desc">
                Dirancang tanpa toleransi downtime (High Availability &amp; zero single point of failure). Kapasitas <strong>digandakan penuh (2N Server Kembar)</strong>, prioritas <strong>Server Fisik Khusus (Baremetal)</strong> untuk performa puncak, serta <strong>Lisensi Resmi Vendor (Enterprise SLA 24/7)</strong>.
              </p>
              <div class="overview-totals-grid">
                <div class="overview-total-item">
                  <span class="tot-label">Total Unit Server</span>
                  <span class="tot-val text-amber">{{ clusterTotals.critNodes }} Node</span>
                </div>
                <div class="overview-total-item">
                  <span class="tot-label">Tenaga Komputasi</span>
                  <span class="tot-val text-amber">{{ clusterTotals.critVcpu }} vCPU</span>
                </div>
                <div class="overview-total-item">
                  <span class="tot-label">Total RAM</span>
                  <span class="tot-val text-amber">{{ clusterTotals.critRam }} GB</span>
                </div>
                <div class="overview-total-item">
                  <span class="tot-label">Kapasitas Disk</span>
                  <span class="tot-val text-amber">{{ formatDiskSize(clusterTotals.critStorage) }}</span>
                </div>
              </div>
            </div>
          </div>
        </transition>
      </section>

      <!-- SMART TOOLBAR: PENCARIAN + FILTER TIER + VIEW MODE + KATEGORI -->
      <div class="preset-smart-toolbar">
        <div class="toolbar-primary-row">
          <div class="search-input-wrap">
            <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="#64748B" stroke-width="2">
              <circle cx="11" cy="11" r="8"/>
              <line x1="21" y1="21" x2="16.65" y2="16.65"/>
            </svg>
            <input
              type="text"
              v-model="searchQuery"
              class="preset-search-input"
              placeholder="Cari berdasarkan nama layer, fungsi bisnis, komponen, atau lisensi..."
            />
            <button
              v-if="searchQuery"
              type="button"
              class="search-clear-btn"
              @click="searchQuery = ''"
            >
              ✕
            </button>
          </div>

          <div class="tier-view-segmented">
            <span class="seg-label">Filter Tier:</span>
            <div class="segmented-control">
              <button
                type="button"
                :class="['seg-btn', { active: activeTierTab === 'both' }]"
                @click="activeTierTab = 'both'"
              >
                Kedua Tier (Bandingkan)
              </button>
              <button
                type="button"
                :class="['seg-btn btn-mod', { active: activeTierTab === 'moderate' }]"
                @click="activeTierTab = 'moderate'"
              >
                Moderat Saja
              </button>
              <button
                type="button"
                :class="['seg-btn btn-crit', { active: activeTierTab === 'critical' }]"
                @click="activeTierTab = 'critical'"
              >
                Critical Saja
              </button>
            </div>
          </div>

          <!-- View Switcher (Cards vs Table vs Diagram) -->
          <div class="view-mode-switcher">
            <button
              type="button"
              :class="['view-switch-btn', { active: viewMode === 'cards' }]"
              @click="viewMode = 'cards'"
              title="Tampilan Kartu Visual Interaktif"
            >
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <rect x="3" y="3" width="7" height="7"></rect>
                <rect x="14" y="3" width="7" height="7"></rect>
                <rect x="14" y="14" width="7" height="7"></rect>
                <rect x="3" y="14" width="7" height="7"></rect>
              </svg>
              <span>Kartu Detail</span>
            </button>
            <button
              type="button"
              :class="['view-switch-btn', { active: viewMode === 'table' }]"
              @click="viewMode = 'table'"
              title="Tampilan Matriks Tabel 1 Layar"
            >
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="3" y1="6" x2="21" y2="6"></line>
                <line x1="3" y1="12" x2="21" y2="12"></line>
                <line x1="3" y1="18" x2="21" y2="18"></line>
              </svg>
              <span>Tabel Ringkas</span>
            </button>
            <button
              type="button"
              :class="['view-switch-btn', { active: viewMode === 'diagram' }]"
              @click="viewMode = 'diagram'"
              title="Tampilan High-Level Architecture Blueprint Diagram"
            >
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <rect x="9" y="2" width="6" height="5" rx="1"></rect>
                <rect x="2" y="17" width="6" height="5" rx="1"></rect>
                <rect x="9" y="17" width="6" height="5" rx="1"></rect>
                <rect x="16" y="17" width="6" height="5" rx="1"></rect>
                <path d="M12 7v5"></path>
                <path d="M5 17v-3a2 2 0 0 1 2-2h10a2 2 0 0 1 2 2v3"></path>
                <line x1="12" y1="12" x2="12" y2="17"></line>
              </svg>
              <span>Diagram</span>
            </button>
          </div>
        </div>

        <!-- CATEGORY PILL FILTER ROW -->
        <div class="category-pills-row">
          <button
            v-for="cat in availableCategories"
            :key="cat.id"
            type="button"
            :class="['category-pill-btn', { active: activeCategory === cat.id }]"
            @click="activeCategory = cat.id"
          >
            <span class="cat-pill-icon">{{ cat.icon }}</span>
            <span class="cat-pill-label">{{ cat.label }}</span>
            <span class="cat-pill-count">{{ categoryCounts[cat.id] || 0 }}</span>
          </button>
        </div>
      </div>

      <!-- VIEW MODE 1: CARDS LIST OF PRESET LAYERS -->
      <div v-if="viewMode === 'cards' && filteredPresetItems.length > 0" class="preset-cards-grid">
        <div
          v-for="item in filteredPresetItems"
          :key="item.id"
          class="preset-card"
        >
          <!-- Card Header -->
          <div class="preset-card-header">
            <div class="card-header-left">
              <!-- Component Icon Avatar -->
              <div class="card-avatar-box">
                <span class="avatar-emoji">{{ getLayerIcon(item) }}</span>
              </div>

              <div class="layer-title-meta-wrap">
                <div class="layer-badges-line">
                  <span class="layer-number-badge">Layer {{ String(item.no).padStart(2, '0') }}</span>
                  <span :class="['category-badge-chip', getCategoryBadge(item).theme]">
                    {{ getCategoryBadge(item).icon }} {{ getCategoryBadge(item).label }}
                  </span>
                  <span class="component-code-chip">{{ item.component }}</span>
                </div>
                <h3 class="layer-title">{{ item.layer }}</h3>
                <!-- Plain-Language Business Function Description -->
                <div v-if="item.notes" class="layer-business-role">
                  <span class="role-icon">💡</span>
                  <span class="role-text">{{ item.notes }}</span>
                </div>
              </div>
            </div>

            <div class="card-header-actions">
              <button
                type="button"
                class="btn-edit-layer"
                @click="openEditLayerModal(item)"
                title="Edit Formula & Baseline Layer"
              >
                <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
                  <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
                </svg>
                <span>Edit Preset</span>
              </button>

              <button
                type="button"
                class="btn-delete-layer"
                @click="confirmDeleteCustomLayer(item)"
                title="Hapus Layer Preset Ini"
              >
                <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="3 6 5 6 21 6"/>
                  <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
                </svg>
              </button>
            </div>
          </div>

          <!-- Formula Specs Comparison Grid -->
          <div :class="['preset-specs-row', { 'single-col': activeTierTab !== 'both' }]">
            <!-- MODERATE SPEC BOX -->
            <div
              v-if="activeTierTab === 'both' || activeTierTab === 'moderate'"
              class="spec-box box-mod"
            >
              <div class="spec-box-hdr">
                <div class="spec-tier-badge mod-theme">
                  <div class="tier-dot mod-dot"></div>
                  <span class="box-tier-title text-navy">Moderat (N+1)</span>
                  <span class="box-tier-sub">· Operasional Reguler</span>
                </div>
                <span :class="['tshirt-size-badge', getSizeBadge(item.moderate).tagClass]">
                  {{ getSizeBadge(item.moderate).label }}
                </span>
              </div>

              <!-- Shared Storage Special View vs Standard Compute Nodes -->
              <div v-if="item.moderate.computeType === 'shared'" class="shared-storage-callout mod-theme">
                <div class="shared-callout-icon">🗄️</div>
                <div class="shared-callout-info">
                  <span class="shared-callout-title">Penyimpanan Bersama Terpusat (SAN LUN Storage Pool)</span>
                  <span class="shared-callout-desc">Volume disk bersama terpusat tanpa unit server prosesor terpisah.</span>
                </div>
                <div class="shared-storage-val">
                  <span class="storage-num">{{ item.moderate.storageGBPerNode }} GB</span>
                  <span class="storage-class-tag">{{ item.moderate.storageClass }}</span>
                </div>
              </div>

              <div v-else class="spec-metrics-grid">
                <div class="metric-cell">
                  <span class="cell-lbl">Unit Server</span>
                  <span class="cell-val">{{ item.moderate.nodeCount }} Unit</span>
                  <span class="cell-sub">{{ formatComputeType(item.moderate.computeType) }}</span>
                </div>
                <div class="metric-cell">
                  <span class="cell-lbl">Tenaga Prosesor</span>
                  <span class="cell-val">{{ item.moderate.vcpuPerNode }} vCPU / node</span>
                  <span class="cell-sub">Total {{ item.moderate.nodeCount * item.moderate.vcpuPerNode }} vCPU</span>
                </div>
                <div class="metric-cell">
                  <span class="cell-lbl">Memori RAM</span>
                  <span class="cell-val">{{ item.moderate.ramGBPerNode }} GB / node</span>
                  <span class="cell-sub">Total {{ item.moderate.nodeCount * item.moderate.ramGBPerNode }} GB</span>
                </div>
                <div class="metric-cell">
                  <span class="cell-lbl">Ruang Disk</span>
                  <span class="cell-val">{{ item.moderate.storageGBPerNode }} GB</span>
                  <span class="cell-sub">{{ item.moderate.storageClass }} · Total {{ (item.moderate.nodeCount || 1) * item.moderate.storageGBPerNode }} GB</span>
                </div>
              </div>

              <div class="spec-box-tier-footer">
                <div class="tier-lic-indicator">
                  <span class="lic-icon">{{ item.moderate.licenseKind === 'enterprise' ? '🛡️' : '🆓' }}</span>
                  <span class="lic-name">
                    <strong>{{ item.moderate.licenseKind === 'enterprise' ? 'Lisensi Resmi' : 'Gratis (OSS)' }}:</strong>
                    {{ item.moderate.licenseKind === 'enterprise' ? (item.licenseEnterprise || 'Official Enterprise') : (item.licenseFree || 'Community OSS') }}
                  </span>
                </div>
                <span :class="['backup-badge', item.moderate.backupEnabled ? 'active' : 'inactive']">
                  {{ item.moderate.backupEnabled ? (item.component.toLowerCase().includes('nginx') ? '✓ Rubrik Backup' : '✓ NetBackup Aktif') : 'Tanpa Backup' }}
                </span>
              </div>
            </div>

            <!-- CRITICAL SPEC BOX -->
            <div
              v-if="activeTierTab === 'both' || activeTierTab === 'critical'"
              class="spec-box box-crit"
            >
              <div class="spec-box-hdr">
                <div class="spec-tier-badge crit-theme">
                  <div class="tier-dot crit-dot"></div>
                  <span class="box-tier-title text-amber">Critical (2N)</span>
                  <span class="box-tier-sub">· Layanan Inti 24/7</span>
                </div>
                <span :class="['tshirt-size-badge', getSizeBadge(item.critical).tagClass]">
                  {{ getSizeBadge(item.critical).label }}
                </span>
              </div>

              <!-- Delta Highlight Pills: Apa yang Meningkat di Critical? -->
              <div class="delta-pills-strip">
                <span class="delta-label">Perbedaan vs Moderat:</span>
                <div class="delta-chips-wrap">
                  <span
                    v-for="(pill, pIdx) in getDeltaPills(item.moderate, item.critical)"
                    :key="pIdx"
                    :class="['delta-chip', pill.type]"
                  >
                    {{ pill.text }}
                  </span>
                </div>
              </div>

              <!-- Shared Storage Special View vs Standard Compute Nodes -->
              <div v-if="item.critical.computeType === 'shared'" class="shared-storage-callout crit-theme">
                <div class="shared-callout-icon">🗄️</div>
                <div class="shared-callout-info">
                  <span class="shared-callout-title">Penyimpanan Bersama Terpusat (SAN LUN Storage Pool)</span>
                  <span class="shared-callout-desc">Volume disk berkecepatan tinggi tanpa unit server terpisah.</span>
                </div>
                <div class="shared-storage-val">
                  <span class="storage-num">{{ item.critical.storageGBPerNode }} GB</span>
                  <span class="storage-class-tag">{{ item.critical.storageClass }}</span>
                </div>
              </div>

              <div v-else class="spec-metrics-grid">
                <div class="metric-cell">
                  <span class="cell-lbl">Unit Server</span>
                  <span class="cell-val">{{ item.critical.nodeCount }} Unit</span>
                  <span class="cell-sub">{{ formatComputeType(item.critical.computeType) }}</span>
                </div>
                <div class="metric-cell">
                  <span class="cell-lbl">Tenaga Prosesor</span>
                  <span class="cell-val">{{ item.critical.vcpuPerNode }} vCPU / node</span>
                  <span class="cell-sub">Total {{ item.critical.nodeCount * item.critical.vcpuPerNode }} vCPU</span>
                </div>
                <div class="metric-cell">
                  <span class="cell-lbl">Memori RAM</span>
                  <span class="cell-val">{{ item.critical.ramGBPerNode }} GB / node</span>
                  <span class="cell-sub">Total {{ item.critical.nodeCount * item.critical.ramGBPerNode }} GB</span>
                </div>
                <div class="metric-cell">
                  <span class="cell-lbl">Ruang Disk</span>
                  <span class="cell-val">{{ item.critical.storageGBPerNode }} GB</span>
                  <span class="cell-sub">{{ item.critical.storageClass }} · Total {{ (item.critical.nodeCount || 1) * item.critical.storageGBPerNode }} GB</span>
                </div>
              </div>

              <div class="spec-box-tier-footer">
                <div class="tier-lic-indicator">
                  <span class="lic-icon">{{ item.critical.licenseKind === 'enterprise' ? '🛡️' : '🆓' }}</span>
                  <span class="lic-name">
                    <strong>{{ item.critical.licenseKind === 'enterprise' ? 'Lisensi Resmi' : 'Gratis (OSS)' }}:</strong>
                    {{ item.critical.licenseKind === 'enterprise' ? (item.licenseEnterprise || 'Official Enterprise') : (item.licenseFree || 'Community OSS') }}
                  </span>
                </div>
                <span :class="['backup-badge', item.critical.backupEnabled ? 'active' : 'inactive']">
                  {{ item.critical.backupEnabled ? (item.component.toLowerCase().includes('nginx') ? '✓ Rubrik Backup' : '✓ NetBackup Aktif') : 'Tanpa Backup' }}
                </span>
              </div>
            </div>
          </div>

          <!-- License & Meta Footer -->
          <div class="preset-card-ftr">
            <div class="ftr-left-refs">
              <span class="ftr-sku-item">
                <span class="sku-label">SKU OSS:</span>
                <span class="sku-value">{{ item.licenseFree || 'Community OSS' }}</span>
              </span>
              <span class="ftr-sku-sep">·</span>
              <span class="ftr-sku-item">
                <span class="sku-label">SKU Enterprise:</span>
                <span class="sku-value text-green">{{ item.licenseEnterprise || 'Enterprise Official' }}</span>
              </span>
            </div>
            <div v-if="item.contractRef" class="ftr-contract-ref" title="Nomor referensi kontrak pengadaan">
              <span class="ref-icon">📋</span>
              <span>Ref: {{ item.contractRef }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- VIEW MODE 2: MATRIX TABLE VIEW (1 SCREEN COMPACT COMPARISON) -->
      <div v-else-if="viewMode === 'table' && filteredPresetItems.length > 0" class="preset-matrix-table-wrap">
        <table class="preset-matrix-table">
          <thead>
            <tr>
              <th rowspan="2" class="th-layer">Layer Arsitektur</th>
              <th rowspan="2" class="th-cat">Kategori</th>
              <th colspan="4" class="th-group-mod">STANDAR BASELINE: MODERAT (N+1)</th>
              <th colspan="5" class="th-group-crit">STANDAR BASELINE: CRITICAL (2N)</th>
              <th rowspan="2" class="th-action">Aksi</th>
            </tr>
            <tr class="sub-headers">
              <!-- MODERAT SUB -->
              <th class="th-sub">Node &amp; Tipe</th>
              <th class="th-sub">vCPU &amp; RAM Total</th>
              <th class="th-sub">Storage Disk</th>
              <th class="th-sub">Lisensi &amp; NBU</th>
              <!-- CRITICAL SUB -->
              <th class="th-sub">Node &amp; Tipe</th>
              <th class="th-sub">vCPU &amp; RAM Total</th>
              <th class="th-sub">Storage Disk</th>
              <th class="th-sub">Lisensi &amp; NBU</th>
              <th class="th-sub th-delta">Upgrade Highlight</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in filteredPresetItems" :key="item.id" class="matrix-row">
              <!-- Layer Col -->
              <td class="td-layer-info">
                <div class="matrix-layer-cell">
                  <span class="matrix-emoji">{{ getLayerIcon(item) }}</span>
                  <div class="matrix-text-wrap">
                    <div class="matrix-layer-name">
                      <span class="matrix-no">L{{ String(item.no).padStart(2, '0') }}</span>
                      <strong>{{ item.layer }}</strong>
                    </div>
                    <div class="matrix-comp-sub">{{ item.component }}</div>
                    <div v-if="item.notes" class="matrix-role-sub" :title="item.notes">{{ item.notes }}</div>
                  </div>
                </div>
              </td>

              <!-- Kategori Col -->
              <td class="td-cat">
                <span :class="['matrix-cat-tag', getCategoryBadge(item).theme]">
                  {{ getCategoryBadge(item).icon }} {{ getCategoryBadge(item).label }}
                </span>
              </td>

              <!-- MODERAT COLS -->
              <td class="td-mod-node">
                <div v-if="item.moderate.computeType === 'shared'" class="matrix-shared-tag">Shared SAN</div>
                <div v-else>
                  <strong>{{ item.moderate.nodeCount }} Node</strong>
                  <div class="matrix-compute-type">{{ item.moderate.computeType === 'baremetal' ? 'Baremetal' : 'VM' }}</div>
                </div>
              </td>
              <td class="td-mod-compute">
                <div v-if="item.moderate.computeType === 'shared'" class="text-muted">-</div>
                <div v-else>
                  <div><strong>{{ item.moderate.nodeCount * item.moderate.vcpuPerNode }}</strong> vCPU</div>
                  <div class="text-muted">{{ item.moderate.nodeCount * item.moderate.ramGBPerNode }} GB RAM</div>
                </div>
              </td>
              <td class="td-mod-storage">
                <strong>{{ (item.moderate.nodeCount || 1) * item.moderate.storageGBPerNode }} GB</strong>
                <div class="text-muted">{{ item.moderate.storageClass }}</div>
              </td>
              <td class="td-mod-lic">
                <div class="matrix-lic-text" :title="item.moderate.licenseKind === 'enterprise' ? item.licenseEnterprise : item.licenseFree">
                  {{ item.moderate.licenseKind === 'enterprise' ? '🛡️ Enterprise' : '🆓 Free OSS' }}
                </div>
                <span :class="['matrix-nbu-dot', item.moderate.backupEnabled ? 'active' : 'inactive']" :title="item.moderate.backupEnabled ? (item.component.toLowerCase().includes('nginx') ? 'Rubrik Backup Aktif' : 'NetBackup Aktif') : 'Tanpa Backup'">
                  {{ item.moderate.backupEnabled ? (item.component.toLowerCase().includes('nginx') ? '✓ Rubrik' : '✓ NBU') : 'No Backup' }}
                </span>
              </td>

              <!-- CRITICAL COLS -->
              <td class="td-crit-node">
                <div v-if="item.critical.computeType === 'shared'" class="matrix-shared-tag">Shared SAN</div>
                <div v-else>
                  <strong>{{ item.critical.nodeCount }} Node</strong>
                  <div class="matrix-compute-type">{{ item.critical.computeType === 'baremetal' ? 'Baremetal' : 'VM' }}</div>
                </div>
              </td>
              <td class="td-crit-compute">
                <div v-if="item.critical.computeType === 'shared'" class="text-muted">-</div>
                <div v-else>
                  <div><strong>{{ item.critical.nodeCount * item.critical.vcpuPerNode }}</strong> vCPU</div>
                  <div class="text-muted">{{ item.critical.nodeCount * item.critical.ramGBPerNode }} GB RAM</div>
                </div>
              </td>
              <td class="td-crit-storage">
                <strong>{{ (item.critical.nodeCount || 1) * item.critical.storageGBPerNode }} GB</strong>
                <div class="text-muted">{{ item.critical.storageClass }}</div>
              </td>
              <td class="td-crit-lic">
                <div class="matrix-lic-text" :title="item.critical.licenseKind === 'enterprise' ? item.licenseEnterprise : item.licenseFree">
                  {{ item.critical.licenseKind === 'enterprise' ? '🛡️ Enterprise' : '🆓 Free OSS' }}
                </div>
                <span :class="['matrix-nbu-dot', item.critical.backupEnabled ? 'active' : 'inactive']" :title="item.critical.backupEnabled ? (item.component.toLowerCase().includes('nginx') ? 'Rubrik Backup Aktif' : 'NetBackup Aktif') : 'Tanpa Backup'">
                  {{ item.critical.backupEnabled ? (item.component.toLowerCase().includes('nginx') ? '✓ Rubrik' : '✓ NBU') : 'No Backup' }}
                </span>
              </td>
              <td class="td-crit-delta">
                <div class="matrix-delta-wrap">
                  <span
                    v-for="(pill, pIdx) in getDeltaPills(item.moderate, item.critical).slice(0, 2)"
                    :key="pIdx"
                    :class="['delta-chip-mini', pill.type]"
                  >
                    {{ pill.text }}
                  </span>
                </div>
              </td>

              <!-- ACTION COL -->
              <td class="td-actions">
                <div class="matrix-action-btns">
                  <button
                    type="button"
                    class="btn-matrix-action edit"
                    @click="openEditLayerModal(item)"
                    title="Edit Preset Layer"
                  >
                    <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
                      <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
                    </svg>
                  </button>
                  <button
                    type="button"
                    class="btn-matrix-action del"
                    @click="confirmDeleteCustomLayer(item)"
                    title="Hapus Layer Preset"
                  >
                    <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <polyline points="3 6 5 6 21 6"/>
                      <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
                    </svg>
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- VIEW MODE 3: HIGH-LEVEL ARCHITECTURE BLUEPRINT DIAGRAM -->
      <div v-else-if="viewMode === 'diagram' && filteredPresetItems.length > 0" class="preset-diagram-canvas">
        <!-- Blueprint Top Header & Visual Legend -->
        <div class="blueprint-header-bar">
          <div class="blueprint-title-block">
            <div class="bp-title-row">
              <span class="bp-badge-tag">REFERENCE ARCHITECTURE</span>
              <h3 class="bp-main-title">High-Level Enterprise Architecture Topology</h3>
              <span :class="['bp-lens-pill', activeTierTab]">
                {{ activeTierTab === 'both' ? '⚖️ Mode Perbandingan: Moderat (N+1) vs Critical (2N)' : (activeTierTab === 'moderate' ? '🔵 Fokus Blueprint: Tier Moderat (N+1)' : '🟠 Fokus Blueprint: Tier Critical (2N)') }}
              </span>
            </div>
            <p class="bp-subtitle">
              Visualisasi alur trafik berlapis (North-South Traffic Flow) beserta penempatan klaster compute, storage SAN, dan lisensi per zona arsitektur. Klik node untuk inspeksi detail atau edit formula.
            </p>
          </div>

          <div class="blueprint-legend-box">
            <div class="legend-item" title="Klaster Virtual Machine Standar">
              <span class="leg-dot vm"></span>
              <span>VM Cluster</span>
            </div>
            <div class="legend-item" title="Server Fisik Khusus Performa Tinggi">
              <span class="leg-dot bm"></span>
              <span>Baremetal Dedicated</span>
            </div>
            <div class="legend-item" title="Penyimpanan Terpusat SAN LUN">
              <span class="leg-dot san"></span>
              <span>Shared SAN LUN</span>
            </div>
            <div class="legend-sep"></div>
            <div class="legend-item" title="Jalur Trafik Sinkron Utama">
              <span class="leg-line sync"></span>
              <span>Sync Data Flow</span>
            </div>
            <div class="legend-item" title="Jalur Asinkron, Quorum, atau Log Telemetry">
              <span class="leg-line async"></span>
              <span>Async / Telemetry</span>
            </div>
          </div>
        </div>

        <!-- External Entry Traffic Node (Shown when viewing All or Access & App) -->
        <div v-if="activeCategory === 'all' || activeCategory === 'access_app'" class="bp-external-entry-wrap">
          <div class="bp-external-pill">
            <div class="ext-icon-pulse">🌍</div>
            <div class="ext-text-group">
              <span class="ext-eyebrow">EXTERNAL TRAFFIC &amp; CHANNELS</span>
              <span class="ext-title">Nasabah Mobile Banking · Web Portal · Partner Open API</span>
            </div>
            <div class="ext-protocol-tags">
              <span class="proto-chip">HTTPS / TLS 1.3</span>
              <span class="proto-chip">WAF &amp; Rate Limiting</span>
            </div>
          </div>
          <div class="bp-spine-connector entry-spine">
            <div class="spine-line"></div>
            <span class="spine-pill">⬇ Masuk ke Zona DMZ &amp; Ingress Layer</span>
          </div>
        </div>

        <!-- Architectural Swimlane Zones -->
        <div class="bp-zones-stack">
          <template v-for="(zone, zIdx) in diagramZones" :key="zone.id">
            <section :class="['bp-zone-lane', zone.themeClass]">
              <!-- Zone Header Bar -->
              <div class="bp-zone-header">
                <div class="zone-hdr-left">
                  <span class="zone-code-badge">{{ zone.code }}</span>
                  <span class="zone-icon-box">{{ zone.icon }}</span>
                  <div class="zone-title-wrap">
                    <div class="zone-title-line">
                      <h4 class="zone-name">{{ zone.title }}</h4>
                      <span class="zone-role-pill">{{ zone.architectureTier }}</span>
                    </div>
                    <p class="zone-desc">{{ zone.subtitle }}</p>
                  </div>
                </div>

                <div class="zone-hdr-metrics">
                  <div class="zone-stat-chip">
                    <span class="z-stat-lbl">Komponen</span>
                    <span class="z-stat-val">{{ zone.items.length }} Layer</span>
                  </div>
                  <div v-if="activeTierTab === 'both' || activeTierTab === 'moderate'" class="zone-stat-chip mod">
                    <span class="z-stat-lbl">Moderat (N+1)</span>
                    <span class="z-stat-val">{{ zone.modNodes }} Node · {{ zone.modVcpu }} vCPU</span>
                  </div>
                  <div v-if="activeTierTab === 'both' || activeTierTab === 'critical'" class="zone-stat-chip crit">
                    <span class="z-stat-lbl">Critical (2N)</span>
                    <span class="z-stat-val">{{ zone.critNodes }} Node · {{ zone.critVcpu }} vCPU</span>
                  </div>
                </div>
              </div>

              <!-- Zone Nodes Flow Grid -->
              <div :class="['bp-zone-nodes-flow', { 'single-node-hub': zone.items.length === 1 }]">
                <template v-for="(item, nIdx) in zone.items" :key="item.id">
                  <!-- Blueprint Node Card -->
                  <div
                    :class="[
                      'bp-node-card',
                      {
                        'is-selected': selectedDiagramItem?.id === item.id,
                        'has-baremetal': (activeTierTab !== 'moderate' && item.critical.computeType === 'baremetal') || (activeTierTab === 'moderate' && item.moderate.computeType === 'baremetal'),
                        'is-shared-san': item.moderate.computeType === 'shared' || item.critical.computeType === 'shared'
                      }
                    ]"
                    @click="selectDiagramNode(item)"
                  >
                    <!-- Node Top Meta Strip -->
                    <div class="bp-node-topbar">
                      <div class="bp-node-badges">
                        <span class="bp-layer-no">L{{ String(item.no).padStart(2, '0') }}</span>
                        <span class="bp-role-tag">{{ getNodeArchitectureRole(item) }}</span>
                      </div>
                      <div class="bp-node-actions" @click.stop>
                        <button
                          type="button"
                          class="bp-btn-edit"
                          @click="openEditLayerModal(item)"
                          title="Edit Formula Preset Layer Ini"
                        >
                          <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2">
                            <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
                            <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
                          </svg>
                          <span>Edit</span>
                        </button>
                      </div>
                    </div>

                    <!-- Node Identity (Icon + Layer + Software Component) -->
                    <div class="bp-node-identity">
                      <div class="bp-node-avatar">
                        <span>{{ getLayerIcon(item) }}</span>
                      </div>
                      <div class="bp-node-titles">
                        <h5 class="bp-node-layer-name">{{ item.layer }}</h5>
                        <div class="bp-node-comp-name">{{ item.component }}</div>
                      </div>
                    </div>

                    <p v-if="item.notes" class="bp-node-desc" :title="item.notes">
                      {{ item.notes }}
                    </p>

                    <!-- MODE A: BOTH TIERS COMPARISON INSIDE NODE -->
                    <div v-if="activeTierTab === 'both'" class="bp-node-dual-tiers">
                      <!-- Moderat Row -->
                      <div class="bp-tier-strip mod-strip">
                        <div class="strip-hdr">
                          <span class="strip-tier-dot mod"></span>
                          <span class="strip-tier-name">MODERAT (N+1)</span>
                          <span :class="['strip-compute-pill', item.moderate.computeType]">
                            {{ item.moderate.computeType === 'shared' ? '🗄️ Shared SAN' : `${item.moderate.nodeCount}x ${item.moderate.computeType === 'baremetal' ? 'Baremetal' : 'VM'}` }}
                          </span>
                        </div>
                        <div class="strip-specs-line">
                          <span v-if="item.moderate.computeType !== 'shared'">
                            <strong>{{ item.moderate.nodeCount * item.moderate.vcpuPerNode }}</strong> vCPU ·
                            <strong>{{ item.moderate.nodeCount * item.moderate.ramGBPerNode }}</strong> GB RAM ·
                          </span>
                          <span>
                            <strong>{{ formatDiskSize((item.moderate.nodeCount || 1) * item.moderate.storageGBPerNode) }}</strong> ({{ item.moderate.storageClass }})
                          </span>
                        </div>
                        <div class="strip-lic-line">
                          <span class="strip-lic-badge">
                            {{ item.moderate.licenseKind === 'enterprise' ? '🛡️' : '🆓' }}
                            {{ item.moderate.licenseKind === 'enterprise' ? (item.licenseEnterprise || 'Enterprise') : (item.licenseFree || 'Community OSS') }}
                          </span>
                        </div>
                      </div>

                      <!-- Critical Row -->
                      <div class="bp-tier-strip crit-strip">
                        <div class="strip-hdr">
                          <span class="strip-tier-dot crit"></span>
                          <span class="strip-tier-name">CRITICAL (2N)</span>
                          <span :class="['strip-compute-pill', item.critical.computeType]">
                            {{ item.critical.computeType === 'shared' ? '🗄️ Shared SAN' : `${item.critical.nodeCount}x ${item.critical.computeType === 'baremetal' ? '🔥 Baremetal' : 'VM'}` }}
                          </span>
                        </div>
                        <div class="strip-specs-line">
                          <span v-if="item.critical.computeType !== 'shared'">
                            <strong>{{ item.critical.nodeCount * item.critical.vcpuPerNode }}</strong> vCPU ·
                            <strong>{{ item.critical.nodeCount * item.critical.ramGBPerNode }}</strong> GB RAM ·
                          </span>
                          <span>
                            <strong>{{ formatDiskSize((item.critical.nodeCount || 1) * item.critical.storageGBPerNode) }}</strong> ({{ item.critical.storageClass }})
                          </span>
                        </div>
                        <div class="strip-lic-line">
                          <span class="strip-lic-badge crit-lic">
                            {{ item.critical.licenseKind === 'enterprise' ? '🛡️' : '🆓' }}
                            {{ item.critical.licenseKind === 'enterprise' ? (item.licenseEnterprise || 'Enterprise') : (item.licenseFree || 'Community OSS') }}
                          </span>
                        </div>
                      </div>

                      <!-- Upgrade Delta Footer -->
                      <div class="bp-node-delta-bar">
                        <span
                          v-for="(pill, pIdx) in getDeltaPills(item.moderate, item.critical).slice(0, 2)"
                          :key="pIdx"
                          :class="['bp-delta-tag', pill.type]"
                        >
                          {{ pill.text }}
                        </span>
                      </div>
                    </div>

                    <!-- MODE B: SINGLE TIER LENS (MODERATE ONLY OR CRITICAL ONLY) -->
                    <div v-else :class="['bp-node-single-tier', activeTierTab]">
                      <template v-for="activeSpec in [activeTierTab === 'moderate' ? item.moderate : item.critical]" :key="activeTierTab">
                        <div class="single-spec-badges">
                          <span :class="['single-compute-badge', activeSpec.computeType]">
                            <template v-if="activeSpec.computeType === 'shared'">
                              🗄️ Shared SAN LUN Pool
                            </template>
                            <template v-else>
                              {{ activeSpec.computeType === 'baremetal' ? '🔥' : '🖥️' }}
                              {{ activeSpec.nodeCount }} Node · {{ activeSpec.computeType === 'baremetal' ? 'Baremetal Dedicated' : 'Virtual Machine' }}
                            </template>
                          </span>
                          <span :class="['single-tshirt-pill', getSizeBadge(activeSpec).tagClass]">
                            {{ getSizeBadge(activeSpec).label }}
                          </span>
                        </div>

                        <div class="single-metrics-box">
                          <div class="s-metric">
                            <span class="s-lbl">Total vCPU</span>
                            <span class="s-val">{{ activeSpec.computeType === 'shared' ? '-' : `${activeSpec.nodeCount * activeSpec.vcpuPerNode} Core` }}</span>
                          </div>
                          <div class="s-metric">
                            <span class="s-lbl">Total RAM</span>
                            <span class="s-val">{{ activeSpec.computeType === 'shared' ? '-' : `${activeSpec.nodeCount * activeSpec.ramGBPerNode} GB` }}</span>
                          </div>
                          <div class="s-metric">
                            <span class="s-lbl">Total Disk ({{ activeSpec.storageClass }})</span>
                            <span class="s-val">{{ formatDiskSize((activeSpec.nodeCount || 1) * activeSpec.storageGBPerNode) }}</span>
                          </div>
                        </div>

                        <div class="single-lic-footer">
                          <span class="s-lic-text">
                            {{ activeSpec.licenseKind === 'enterprise' ? '🛡️' : '🆓' }}
                            {{ activeSpec.licenseKind === 'enterprise' ? (item.licenseEnterprise || 'Official Enterprise') : (item.licenseFree || 'Community OSS') }}
                          </span>
                          <span :class="['s-nbu-pill', { on: activeSpec.backupEnabled }]">
                            {{ activeSpec.backupEnabled ? (item.component.toLowerCase().includes('nginx') ? '✓ Rubrik' : '✓ NBU') : 'No Backup' }}
                          </span>
                        </div>
                      </template>
                    </div>
                  </div>

                  <!-- Intra-Zone Horizontal Flow Connector (If nodes in this zone have a direct relationship) -->
                  <div
                    v-if="nIdx < zone.items.length - 1 && getIntraZoneConnector(zone.id, item, zone.items[nIdx + 1])"
                    class="bp-intra-connector"
                  >
                    <span class="intra-label">{{ getIntraZoneConnector(zone.id, item, zone.items[nIdx + 1])?.label }}</span>
                    <div :class="['intra-arrow-line', getIntraZoneConnector(zone.id, item, zone.items[nIdx + 1])?.style]">
                      <span class="arrow-head">➔</span>
                    </div>
                  </div>
                </template>
              </div>
            </section>

            <!-- Inter-Zone North-South Spine Connector -->
            <div v-if="zIdx < diagramZones.length - 1" class="bp-spine-connector">
              <div :class="['spine-line', zone.spineStyle || 'sync']"></div>
              <span class="spine-pill">{{ zone.nextFlowLabel }}</span>
            </div>
          </template>
        </div>

        <!-- Selected Node Quick Inspector Bar (Appears when user clicks any node in the diagram) -->
        <transition name="banner-slide">
          <div v-if="selectedDiagramItem" class="bp-inspector-panel">
            <div class="insp-left">
              <div class="insp-avatar">{{ getLayerIcon(selectedDiagramItem) }}</div>
              <div class="insp-info">
                <div class="insp-title-line">
                  <span class="insp-no">LAYER {{ String(selectedDiagramItem.no).padStart(2, '0') }}</span>
                  <h4 class="insp-layer">{{ selectedDiagramItem.layer }}</h4>
                  <span class="insp-comp">{{ selectedDiagramItem.component }}</span>
                  <span v-if="selectedDiagramItem.contractRef" class="insp-ref">📋 Ref: {{ selectedDiagramItem.contractRef }}</span>
                </div>
                <p class="insp-notes">{{ selectedDiagramItem.notes || 'Komponen standar baseline arsitektur.' }}</p>
                <div class="insp-lic-pair">
                  <span><strong>SKU OSS:</strong> {{ selectedDiagramItem.licenseFree || 'Community OSS' }}</span>
                  <span class="insp-sep">·</span>
                  <span><strong>SKU Enterprise:</strong> <em class="text-green">{{ selectedDiagramItem.licenseEnterprise || 'Official Enterprise' }}</em></span>
                </div>
              </div>
            </div>

            <div class="insp-actions">
              <button
                type="button"
                class="btn-insp-edit"
                @click="openEditLayerModal(selectedDiagramItem)"
              >
                <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2">
                  <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
                  <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
                </svg>
                <span>Edit Formula Layer</span>
              </button>
              <button
                type="button"
                class="btn-insp-close"
                @click="selectedDiagramId = null"
                title="Tutup panel inspeksi"
              >
                ✕
              </button>
            </div>
          </div>
        </transition>
      </div>

      <!-- EMPTY STATE IF SEARCH / CATEGORY FINDS NOTHING -->
      <div v-if="filteredPresetItems.length === 0" class="preset-empty-state">
        <div class="empty-icon-box">🔍</div>
        <h3 class="empty-title">Tidak ada layer preset yang cocok</h3>
        <p class="empty-desc">Tidak ditemukan layer arsitektur yang sesuai dengan filter atau kata kunci pencarian Anda.</p>
        <button
          type="button"
          class="btn-reset-filter"
          @click="searchQuery = ''; activeCategory = 'all'"
        >
          Reset Filter &amp; Pencarian
        </button>
      </div>
    </div>

    <!-- MODAL: ADD / EDIT PRESET LAYER FORM (Ramah Non-Tech Savvy & Rapi Elegant) -->
    <div
      v-if="showModal"
      class="sketch-modal-backdrop"
      @click.self="showModal = false"
    >
      <div class="sketch-modal-box preset-edit-modal">
        <div class="modal-hdr">
          <div class="modal-hdr-left">
            <div class="hdr-icon-circle">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="#002060" stroke-width="2.2">
                <path d="M12 20h9"/>
                <path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z"/>
              </svg>
            </div>
            <div>
              <h3>{{ isEditing ? 'Edit Formula Preset Layer' : 'Tambah Layer Arsitektur Baru' }}</h3>
              <p class="modal-hdr-sub">Tetapkan formula baseline hardware &amp; lisensi resmi untuk Moderat dan Critical</p>
            </div>
          </div>
          <button type="button" class="modal-close" @click="showModal = false" title="Tutup">✕</button>
        </div>

        <!-- Wizard / Section Tabs Navigation -->
        <div class="modal-section-tabs">
          <button
            type="button"
            class="tab-step-btn"
            :class="{ active: modalActiveTab === 'general' }"
            @click="modalActiveTab = 'general'"
          >
            <span class="step-num">1</span>
            <span class="step-lbl">Identitas Komponen</span>
          </button>
          <button
            type="button"
            class="tab-step-btn"
            :class="{ active: modalActiveTab === 'moderate' }"
            @click="modalActiveTab = 'moderate'"
          >
            <span class="step-num step-mod">2</span>
            <span class="step-lbl">Tier Moderat (N+1)</span>
          </button>
          <button
            type="button"
            class="tab-step-btn"
            :class="{ active: modalActiveTab === 'critical' }"
            @click="modalActiveTab = 'critical'"
          >
            <span class="step-num step-crit">3</span>
            <span class="step-lbl">Tier Critical (2N)</span>
          </button>
        </div>

        <form novalidate @submit.prevent="savePresetLayer" class="modal-content preset-form-clean">
          <!-- ========================================== -->
          <!-- TAB 1: INFORMASI UMUM KOMPONEN             -->
          <!-- ========================================== -->
          <div v-show="modalActiveTab === 'general'" class="tab-pane-content">
            <div class="field-instruction-box">
              <span class="inst-icon">💡</span>
              <p class="inst-text">
                Tentukan nama layer arsitektur dan software yang digunakan. Nama lisensi free dan enterprise akan dijadikan acuan otomatis saat tim memilih model lisensi.
              </p>
            </div>

            <div class="form-grid-2col">
              <div class="clean-field">
                <label class="clean-label">Nama Layer Arsitektur <span class="req">*</span></label>
                <BaseAutocomplete
                  v-model="form.layer"
                  :options="layerOptions"
                  placeholder="Ketik atau pilih layer, cth: API Gateway, Security WAF"
                  required
                />
                <span class="clean-hint">Pilih dari katalog layer atau ketik kategori arsitektur baru.</span>
              </div>

              <div class="clean-field">
                <label class="clean-label">Komponen / Software Default <span class="req">*</span></label>
                <BaseAutocomplete
                  v-model="form.component"
                  :options="componentOptions"
                  placeholder="Ketik atau pilih software, cth: Kong Enterprise / NGINX"
                  required
                  @select="handleComponentSelect"
                />
                <span class="clean-hint">Memilih dari katalog akan otomatis menyarankan lisensi &amp; formula.</span>
              </div>
            </div>

            <div class="form-grid-2col">
              <div class="clean-field">
                <label class="clean-label">Nama Lisensi Versi Free (OSS)</label>
                <BaseAutocomplete
                  v-model="form.licenseFree"
                  :options="licenseFreeOptions"
                  placeholder="Ketik atau pilih lisensi free, cth: Community OSS"
                  required
                />
                <span class="clean-hint">Software open source tanpa biaya lisensi vendor (Rp 0).</span>
              </div>

              <div class="clean-field">
                <label class="clean-label">Nama Lisensi Versi Enterprise</label>
                <BaseAutocomplete
                  v-model="form.licenseEnterprise"
                  :options="licenseEnterpriseOptions"
                  placeholder="Ketik atau pilih SKU lisensi enterprise"
                  required
                />
                <span class="clean-hint">Lisensi resmi vendor bergaransi 24/7 dari katalog SKU.</span>
              </div>
            </div>

            <div class="form-grid-2col">
              <div class="clean-field full-span">
                <label class="clean-label">Catatan Arsitektur / Justifikasi EA Forum</label>
                <textarea
                  v-model="form.notes"
                  class="clean-textarea"
                  rows="2"
                  placeholder="Contoh: Standar Livin' Merchant Release 6 & SuperApp Core untuk High Throughput..."
                ></textarea>
              </div>

              <div class="clean-field">
                <label class="clean-label">No. Referensi Kontrak / Master SKU (Opsional)</label>
                <input
                  type="text"
                  v-model="form.contractRef"
                  class="clean-input font-mono"
                  placeholder="Contoh: TIF.DDL/GH.9076/2025"
                />
              </div>
            </div>
          </div>

          <!-- ========================================== -->
          <!-- TAB 2: FORMULA TIER MODERAT (N+1)          -->
          <!-- ========================================== -->
          <div v-show="modalActiveTab === 'moderate'" class="tab-pane-content">
            <div class="tier-banner-info banner-mod">
              <div class="tier-dot mod-dot"></div>
              <div class="banner-texts">
                <span class="banner-title text-navy">Formula Baseline: Tier Moderat (N+1)</span>
                <span class="banner-sub">Workload umum, backoffice, batch processing, dan service internal (HA N+1 redundancy).</span>
              </div>
            </div>

            <!-- Ukuran Server Size Picker -->
            <div class="clean-field">
              <div class="label-with-tip">
                <label class="clean-label">Ukuran Komputasi Server per Unit:</label>
                <BaseHelpTip
                  title="Ukuran Server Moderat"
                  subtitle="vCPU &amp; RAM"
                  :content="glossary.vcpu.shortDesc"
                  :business-impact="glossary.vcpu.businessImpact"
                />
              </div>
              <BaseSizePicker
                :vcpu="form.moderate.vcpuPerNode"
                :ram="form.moderate.ramGBPerNode"
                @update:vcpu="form.moderate.vcpuPerNode = $event"
                @update:ram="form.moderate.ramGBPerNode = $event"
              />
            </div>

            <!-- Stepper Jumlah Node & Storage -->
            <div class="steppers-card-row">
              <div class="stepper-item-box">
                <div class="label-with-tip">
                  <span class="box-lbl">Jumlah Server Aktif:</span>
                  <BaseHelpTip
                    title="Alokasi Node"
                    subtitle="Server Fisik/VM Aktif"
                    :content="glossary.nodeCount.shortDesc"
                    :business-impact="glossary.nodeCount.businessImpact"
                  />
                </div>
                <BaseStepper
                  v-model="form.moderate.nodeCount"
                  :min="1"
                  :max="32"
                  label="Server"
                  unit="Nodes"
                />
              </div>

              <div class="stepper-item-box">
                <div class="label-with-tip">
                  <span class="box-lbl">Storage / Server:</span>
                  <BaseHelpTip
                    title="Kapasitas Disk"
                    subtitle="SAN Storage"
                    :content="glossary.storage.shortDesc"
                    :business-impact="glossary.storage.businessImpact"
                    align="right"
                  />
                </div>
                <BaseStepper
                  v-model="form.moderate.storageGBPerNode"
                  :min="0"
                  :max="20000"
                  :step="50"
                  label="Storage"
                  unit="GB"
                  inline-unit="GB"
                />
              </div>
            </div>

            <!-- Pilihan Platform Compute -->
            <div class="clean-field">
              <label class="clean-label">Tipe Lingkungan Compute:</label>
              <BaseChoiceCard
                v-model="form.moderate.computeType"
                :columns="2"
                :options="[
                  {
                    value: 'vm',
                    title: 'Virtual Machine (VM)',
                    technicalLabel: 'Server Virtual Fleksibel',
                    description: 'Efisien & hemat biaya untuk aplikasi umum.',
                    icon: '🖥️'
                  },
                  {
                    value: 'baremetal',
                    title: 'Server Fisik Khusus',
                    technicalLabel: 'Baremetal Dedicated',
                    description: 'Mesin fisik dedicated tanpa lapisan virtual.',
                    icon: '🏢',
                    badge: 'Dedicated'
                  }
                ]"
              />
            </div>

            <!-- Storage Class & NetBackup / Rubrik -->
            <div class="storage-class-toggle-row">
              <div class="class-picker-half">
                <label class="clean-label">Kelas Kecepatan Storage:</label>
                <BaseChoiceCard
                  v-model="form.moderate.storageClass"
                  :columns="2"
                  :options="[
                    {
                      value: 'Tier 2',
                      title: 'Standar Andal (Tier 2)',
                      technicalLabel: 'SAN Tier 2 Standard',
                      description: 'Kecepatan standar untuk web & worker (Rp 33.000/GB).',
                      icon: '💾'
                    },
                    {
                      value: 'Tier 3',
                      title: 'Ultra Cepat (Tier 3)',
                      technicalLabel: 'High-Perf ASM',
                      description: 'IOPS tinggi khusus Oracle ASM (Rp 48.000/GB).',
                      icon: '⚡'
                    }
                  ]"
                />
              </div>

              <div class="backup-chk-half">
                <label class="clean-label">Proteksi Backup Cadangan:</label>
                <label class="clean-backup-chk-box" :class="{ active: form.moderate.backupEnabled }">
                  <input type="checkbox" v-model="form.moderate.backupEnabled" class="custom-chk-big" />
                  <div class="chk-desc-wrap">
                    <span class="chk-title">{{ form.component.toLowerCase().includes('nginx') ? 'Sertakan Rubrik Backup' : 'Sertakan NetBackup NBU' }}</span>
                    <span class="chk-sub">{{ form.component.toLowerCase().includes('nginx') ? 'Pencadangan Rubrik Backup Rp 35.000 / GB.' : 'Pencadangan NetBackup Rp 27.000 / GB.' }}</span>
                  </div>
                </label>
              </div>
            </div>
          </div>

          <!-- ========================================== -->
          <!-- TAB 3: FORMULA TIER CRITICAL (2N)          -->
          <!-- ========================================== -->
          <div v-show="modalActiveTab === 'critical'" class="tab-pane-content">
            <div class="tier-banner-info banner-crit">
              <div class="tier-dot crit-dot"></div>
              <div class="banner-texts">
                <div class="title-with-copy-btn">
                  <span class="banner-title text-amber">Formula Baseline: Tier Critical (2N)</span>
                  <button
                    type="button"
                    class="btn-copy-from-mod"
                    @click="copyModerateToCritical"
                    title="Salin parameter dari Moderat sebagai titik awal"
                  >
                    Salin dari Moderat 📋
                  </button>
                </div>
                <span class="banner-sub">Layanan Core Banking, transaksi nasabah intensif, SuperApp Livin' (2N Redundancy, HA 99.99%).</span>
              </div>
            </div>

            <!-- Ukuran Server Size Picker -->
            <div class="clean-field">
              <div class="label-with-tip">
                <label class="clean-label">Ukuran Komputasi Server per Unit:</label>
                <BaseHelpTip
                  title="Ukuran Server Critical"
                  subtitle="vCPU &amp; RAM"
                  :content="glossary.vcpu.shortDesc"
                  :business-impact="glossary.vcpu.businessImpact"
                />
              </div>
              <BaseSizePicker
                :vcpu="form.critical.vcpuPerNode"
                :ram="form.critical.ramGBPerNode"
                @update:vcpu="form.critical.vcpuPerNode = $event"
                @update:ram="form.critical.ramGBPerNode = $event"
              />
            </div>

            <!-- Stepper Jumlah Node & Storage -->
            <div class="steppers-card-row">
              <div class="stepper-item-box">
                <div class="label-with-tip">
                  <span class="box-lbl">Jumlah Server Aktif:</span>
                  <BaseHelpTip
                    title="Alokasi Node Critical"
                    subtitle="Minimal 2N (2 Server)"
                    :content="glossary.nodeCount.shortDesc"
                    :business-impact="glossary.nodeCount.businessImpact"
                  />
                </div>
                <BaseStepper
                  v-model="form.critical.nodeCount"
                  :min="1"
                  :max="32"
                  label="Server"
                  unit="Nodes"
                />
              </div>

              <div class="stepper-item-box">
                <div class="label-with-tip">
                  <span class="box-lbl">Storage / Server:</span>
                  <BaseHelpTip
                    title="Kapasitas Disk Critical"
                    subtitle="SAN Storage"
                    :content="glossary.storage.shortDesc"
                    :business-impact="glossary.storage.businessImpact"
                    align="right"
                  />
                </div>
                <BaseStepper
                  v-model="form.critical.storageGBPerNode"
                  :min="0"
                  :max="20000"
                  :step="50"
                  label="Storage"
                  unit="GB"
                  inline-unit="GB"
                />
              </div>
            </div>

            <!-- Pilihan Platform Compute -->
            <div class="clean-field">
              <label class="clean-label">Tipe Lingkungan Compute:</label>
              <BaseChoiceCard
                v-model="form.critical.computeType"
                :columns="2"
                :options="[
                  {
                    value: 'vm',
                    title: 'Virtual Machine (VM)',
                    technicalLabel: 'Server Virtual Fleksibel',
                    description: 'Efisien & mudah di-scale up.',
                    icon: '🖥️'
                  },
                  {
                    value: 'baremetal',
                    title: 'Server Fisik Khusus',
                    technicalLabel: 'Baremetal Dedicated',
                    description: 'Dedicated tanpa sharing daya untuk performa TPS maksimal.',
                    icon: '🏢',
                    badge: 'Direkomendasikan',
                    badgeType: 'warning'
                  }
                ]"
              />
            </div>

            <!-- Storage Class & NetBackup / Rubrik -->
            <div class="storage-class-toggle-row">
              <div class="class-picker-half">
                <label class="clean-label">Kelas Kecepatan Storage:</label>
                <BaseChoiceCard
                  v-model="form.critical.storageClass"
                  :columns="2"
                  :options="[
                    {
                      value: 'Tier 2',
                      title: 'Standar Andal (Tier 2)',
                      technicalLabel: 'SAN Tier 2 Standard',
                      description: 'Kecepatan standar untuk web & worker (Rp 33.000/GB).',
                      icon: '💾'
                    },
                    {
                      value: 'Tier 3',
                      title: 'Ultra Cepat (Tier 3)',
                      technicalLabel: 'High-Perf ASM',
                      description: 'IOPS tinggi khusus Oracle ASM (Rp 48.000/GB).',
                      icon: '⚡',
                      badge: 'High-IOPS',
                      badgeType: 'warning'
                    }
                  ]"
                />
              </div>

              <div class="backup-chk-half">
                <label class="clean-label">Proteksi Backup Cadangan:</label>
                <label class="clean-backup-chk-box" :class="{ active: form.critical.backupEnabled }">
                  <input type="checkbox" v-model="form.critical.backupEnabled" class="custom-chk-big" />
                  <div class="chk-desc-wrap">
                    <span class="chk-title">{{ form.component.toLowerCase().includes('nginx') ? 'Sertakan Rubrik Backup' : 'Sertakan NetBackup NBU' }}</span>
                    <span class="chk-sub">{{ form.component.toLowerCase().includes('nginx') ? 'Pencadangan Rubrik Backup Rp 35.000 / GB.' : 'Pencadangan NetBackup Rp 27.000 / GB.' }}</span>
                  </div>
                </label>
              </div>
            </div>
          </div>

          <!-- Bottom Action Buttons with Next/Prev and Submit -->
          <div class="modal-ftr-clean">
            <button type="button" class="clean-btn-cancel" @click="showModal = false">Batal</button>
            <div class="modal-ftr-right">
              <button
                v-if="modalActiveTab !== 'general'"
                type="button"
                class="clean-btn-nav"
                @click="modalActiveTab = modalActiveTab === 'critical' ? 'moderate' : 'general'"
              >
                ← Kembali
              </button>

              <button
                v-if="modalActiveTab !== 'critical'"
                type="button"
                class="clean-btn-nav primary-step"
                @click="modalActiveTab = modalActiveTab === 'general' ? 'moderate' : 'critical'"
              >
                Lanjut ke Formula {{ modalActiveTab === 'general' ? 'Moderat' : 'Critical' }} →
              </button>

              <button
                v-else
                type="submit"
                class="clean-btn-submit"
              >
                {{ isEditing ? 'Simpan Perubahan Preset' : 'Tambahkan ke Preset' }}
              </button>
            </div>
          </div>
        </form>
      </div>
    </div>

    <!-- MODAL: KATALOG BELUM DIPETAKAN (CLEAN PICKER FOR MANY ITEMS) -->
    <div v-if="showCatalogPickerModal" class="sketch-modal-backdrop" @click.self="showCatalogPickerModal = false">
      <div class="catalog-picker-modal">
        <div class="picker-modal-hdr">
          <div class="picker-hdr-left">
            <h2 class="picker-title">Komponen Katalog Belum Dipetakan</h2>
            <p class="picker-subtitle">
              Pilih komponen dari master katalog untuk dijadikan layer arsitektur baru, atau abaikan jika hanya berupa lisensi/tarif dasar.
            </p>
          </div>
          <button type="button" class="modal-close" @click="showCatalogPickerModal = false" title="Tutup">✕</button>
        </div>

        <!-- Search & Filter Bar inside Modal -->
        <div class="picker-toolbar">
          <div class="picker-search-box">
            <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="#64748B" stroke-width="2">
              <circle cx="11" cy="11" r="8"/>
              <line x1="21" y1="21" x2="16.65" y2="16.65"/>
            </svg>
            <input
              type="text"
              v-model="catalogPickerSearch"
              class="picker-search-input"
              placeholder="Cari nama komponen, kategori, atau kode SKU..."
            />
            <button
              v-if="catalogPickerSearch"
              type="button"
              class="picker-search-clear"
              @click="catalogPickerSearch = ''"
            >
              ✕
            </button>
          </div>

          <div class="picker-filter-tabs">
            <button
              type="button"
              :class="['picker-tab-btn', { active: catalogPickerFilter === 'all' }]"
              @click="catalogPickerFilter = 'all'"
            >
              Semua ({{ unregisteredCatalogItems.length }})
            </button>
            <button
              type="button"
              :class="['picker-tab-btn', { active: catalogPickerFilter === 'component' }]"
              @click="catalogPickerFilter = 'component'"
            >
              Komponen
            </button>
            <button
              type="button"
              :class="['picker-tab-btn', { active: catalogPickerFilter === 'pricing' }]"
              @click="catalogPickerFilter = 'pricing'"
            >
              SKU &amp; Lisensi
            </button>
          </div>
        </div>

        <!-- Clean Scrollable Table / List -->
        <div class="picker-list-body">
          <div v-if="filteredPickerCatalogItems.length === 0" class="picker-empty-state">
            <span>Tidak ada komponen katalog yang cocok dengan pencarian Anda.</span>
          </div>

          <div
            v-for="catItem in filteredPickerCatalogItems"
            :key="catItem.id + '_' + catItem.name"
            class="picker-row-item"
          >
            <div class="picker-row-main">
              <div class="picker-row-top">
                <span class="picker-item-name">{{ catItem.name }}</span>
                <span class="picker-item-cat">{{ catItem.category }}</span>
                <span v-if="catItem.sku" class="picker-item-sku">{{ catItem.sku }}</span>
              </div>
              <p v-if="catItem.description" class="picker-item-desc">{{ catItem.description }}</p>
            </div>

            <div class="picker-row-actions">
              <button
                type="button"
                class="btn-picker-ignore"
                @click="ignoreCatalogItem(catItem)"
                title="Sembunyikan item ini dari daftar rekomendasi layer"
              >
                Abaikan
              </button>
              <button
                type="button"
                class="btn-picker-select"
                @click="selectCatalogItemFromPicker(catItem)"
              >
                <span>Buat Layer</span>
                <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round">
                  <line x1="5" y1="12" x2="19" y2="12"></line>
                  <polyline points="12 5 19 12 12 19"></polyline>
                </svg>
              </button>
            </div>
          </div>
        </div>

        <div class="picker-modal-ftr">
          <div class="picker-ftr-left">
            <button
              v-if="ignoredCatalogKeys.size > 0"
              type="button"
              class="btn-restore-ignored"
              @click="restoreIgnoredCatalogItems"
            >
              Tampilkan kembali {{ ignoredCatalogKeys.size }} item yang diabaikan
            </button>
            <span v-else class="picker-ftr-hint">
              Menampilkan {{ filteredPickerCatalogItems.length }} dari {{ unregisteredCatalogItems.length }} komponen belum terpetakan
            </span>
          </div>
          <div class="picker-ftr-right">
            <button
              type="button"
              class="clean-btn-cancel"
              @click="showCatalogPickerModal = false"
            >
              Tutup
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, watch } from 'vue'
import type { SizingPresetItem, LayerPresetSpec, PricingItemRecord } from '~/types/sizing'
import BaseSizePicker from '~/components/base/BaseSizePicker.vue'
import BaseChoiceCard from '~/components/base/BaseChoiceCard.vue'
import BaseHelpTip from '~/components/base/BaseHelpTip.vue'
import BaseStepper from '~/components/base/BaseStepper.vue'
import BaseAutocomplete, { type AutocompleteOption } from '~/components/base/BaseAutocomplete.vue'
import { SIZING_GLOSSARY } from '~/data/sizingGlossary'
import { useBackendBaseUrl } from '~/composables/useSizingApi'

const glossary = SIZING_GLOSSARY
const modalActiveTab = ref<'general' | 'moderate' | 'critical'>('general')

const props = withDefaults(
  defineProps<{
    presetItems?: SizingPresetItem[]
    catalogItems?: PricingItemRecord[]
    isSidebarCollapsed?: boolean
  }>(),
  {
    presetItems: () => [],
    catalogItems: () => [],
    isSidebarCollapsed: false
  }
)

const emit = defineEmits<{
  (e: 'update:preset-items', items: SizingPresetItem[]): void
  (e: 'save-preset-item', item: SizingPresetItem): void
  (e: 'delete-preset-item', itemId: string): void
  (e: 'reset-to-official-defaults'): void
  (e: 'toggle-sidebar'): void
  (e: 'update:unregistered-count', count: number): void
}>()

// Backend catalog components if available
interface BackendCatalogComponent {
  id: string
  layer_no: number
  layer_name: string
  component_name: string
  description: string
  notes: string
  baselines?: Array<{
    tier: 'moderate' | 'critical'
    compute_type: 'vm' | 'baremetal' | 'shared'
    node_count: number
    vcpu_per_node: number
    ram_gb_per_node: number
    storage_gb_per_node: number
    storage_class: 'Tier 2' | 'Tier 3'
  }>
  licenses?: Array<{
    license_kind: 'free' | 'enterprise'
    product_name: string
    contract_ref?: string
  }>
}
const backendComponents = ref<BackendCatalogComponent[]>([])

onMounted(async () => {
  try {
    const apiBase = useBackendBaseUrl()
    const res = await $fetch<{ data: BackendCatalogComponent[] }>(`${apiBase}/api/v1/catalog/components`)
    if (res && res.data) {
      backendComponents.value = res.data
    }
  } catch (e) {
    // Fail silently, fallback to preset and pricingList
  }
})

// ==========================================
// NEW CATALOG ALERT & LAYER AUTO-SUGGESTION
// ==========================================
export interface UnregisteredCatalogItem {
  id: string
  name: string
  category: string
  sku?: string
  description?: string
  source: 'component' | 'pricing'
  meta?: any
}

const isBannerDismissed = ref(false)
const showCatalogPickerModal = ref(false)
const catalogPickerSearch = ref('')
const catalogPickerFilter = ref<'all' | 'component' | 'pricing'>('all')
const ignoredCatalogKeys = ref<Set<string>>(new Set())

function dismissNewCatalogBanner() {
  isBannerDismissed.value = true
}

function openCatalogPickerModal() {
  catalogPickerSearch.value = ''
  catalogPickerFilter.value = 'all'
  showCatalogPickerModal.value = true
}

function ignoreCatalogItem(catItem: UnregisteredCatalogItem) {
  const next = new Set(ignoredCatalogKeys.value)
  next.add((catItem.name || '').trim().toLowerCase())
  ignoredCatalogKeys.value = next
}

function restoreIgnoredCatalogItems() {
  ignoredCatalogKeys.value = new Set()
}

function selectCatalogItemFromPicker(catItem: UnregisteredCatalogItem) {
  showCatalogPickerModal.value = false
  openAddLayerModalWithCatalogItem(catItem)
}

const unregisteredCatalogItems = computed<UnregisteredCatalogItem[]>(() => {
  const items = props.presetItems || []
  const existingComponents = new Set(
    items.map(p => (p?.component || '').trim().toLowerCase()).filter(Boolean)
  )
  const existingLayers = new Set(
    items.map(p => (p?.layer || '').trim().toLowerCase()).filter(Boolean)
  )
  const existingLicenses = new Set([
    ...items.map(p => (p?.licenseEnterprise || '').trim().toLowerCase()).filter(Boolean),
    ...items.map(p => (p?.licenseFree || '').trim().toLowerCase()).filter(Boolean)
  ])

  const allExistingText = items
    .map(p => `${p?.layer || ''} ${p?.component || ''} ${p?.licenseEnterprise || ''} ${p?.licenseFree || ''}`.toLowerCase())
    .join(' | ')

  const productKeywordFamilies = [
    ['openshift', 'ocp'],
    ['confluent', 'kafka'],
    ['oracle database', 'oracle rac', 'oracle asm'],
    ['redis enterprise', 'redis caching', 'redis sentinel'],
    ['elasticsearch', 'elastic platinum', 'logstash', 'kibana'],
    ['nginx'],
    ['minio'],
    ['control-m']
  ]

  const isAlreadyCoveredByPreset = (nameLower: string, layerLower?: string) => {
    if (existingComponents.has(nameLower) || existingLayers.has(nameLower) || existingLicenses.has(nameLower)) {
      return true
    }
    if (layerLower && (existingLayers.has(layerLower) || existingComponents.has(layerLower))) {
      return true
    }
    for (const lic of existingLicenses) {
      if (lic.length > 5 && (lic.includes(nameLower) || nameLower.includes(lic))) {
        return true
      }
    }
    for (const family of productKeywordFamilies) {
      const itemMatchesFamily = family.some(kw => nameLower.includes(kw))
      const presetHasFamily = family.some(kw => allExistingText.includes(kw))
      if (itemMatchesFamily && presetHasFamily) {
        return true
      }
    }
    return false
  }

  // Raw infrastructure unit rate names that are not software/architectural layers
  const rawInfraKeywords = [
    'compute vcpu',
    'memory ram',
    'san storage',
    'netbackup',
    'baremetal server',
    'vcpu',
    'ram gb'
  ]

  const unreg: UnregisteredCatalogItem[] = []
  const seenNames = new Set<string>()

  // 1. Check backend catalog components
  for (const bc of backendComponents.value) {
    const compName = (bc.component_name || '').trim()
    const compLower = compName.toLowerCase()
    const layerLower = (bc.layer_name || '').trim().toLowerCase()
    if (!compName || seenNames.has(compLower) || ignoredCatalogKeys.value.has(compLower)) continue

    if (!isAlreadyCoveredByPreset(compLower, layerLower)) {
      seenNames.add(compLower)
      const freeLic = bc.licenses?.find(l => l.license_kind === 'free')?.product_name || ''
      const entLic = bc.licenses?.find(l => l.license_kind === 'enterprise')?.product_name || ''
      const contractRef = bc.licenses?.find(l => l.contract_ref)?.contract_ref || ''
      const modBaseline = bc.baselines?.find(b => b.tier === 'moderate')
      const critBaseline = bc.baselines?.find(b => b.tier === 'critical')

      unreg.push({
        id: bc.id,
        name: compName,
        category: bc.layer_name || 'Komponen Arsitektur',
        description: bc.description || bc.notes,
        source: 'component',
        meta: {
          layer: bc.layer_name,
          licenseFree: freeLic,
          licenseEnterprise: entLic,
          notes: bc.notes,
          contractRef,
          modBaseline,
          critBaseline
        }
      })
    }
  }

  // 2. Check catalog items from pricingList (SKUs)
  if (props.catalogItems) {
    for (const item of props.catalogItems) {
      const name = (item.name || '').trim()
      const nameLower = name.toLowerCase()
      if (!name || seenNames.has(nameLower) || ignoredCatalogKeys.value.has(nameLower)) continue

      const catUpper = (item.category || '').toUpperCase()
      // Exclude raw compute/storage/backup unit rates from layer suggestions
      const isRawInfraUnit =
        ['COMPUTE', 'STORAGE', 'BACKUP'].includes(catUpper) ||
        rawInfraKeywords.some(k => nameLower === k || nameLower.startsWith(k))
      if (isRawInfraUnit) continue

      const isRelevant = ['LICENSE', 'SOFTWARE', 'DATABASE', 'MIDDLEWARE', 'SECURITY'].includes(catUpper)
      if (isRelevant && !isAlreadyCoveredByPreset(nameLower)) {
        seenNames.add(nameLower)
        unreg.push({
          id: item.id || item.sku,
          name: item.name,
          category: item.category,
          sku: item.sku,
          description: `Katalog SKU: ${item.sku} · ${item.environment}`,
          source: 'pricing',
          meta: {
            layer: catUpper === 'DATABASE' ? 'Database Core OLTP' : (catUpper === 'SECURITY' ? 'Security WAF & API Gateway' : item.name),
            licenseEnterprise: item.name,
            contractRef: item.sku
          }
        })
      }
    }
  }

  return unreg
})

const unregisteredPreviewText = computed(() => {
  const list = unregisteredCatalogItems.value || []
  if (list.length === 0) return ''
  const firstTwo = list.slice(0, 2).map(i => i.name)
  if (list.length <= 2) {
    return firstTwo.join(', ')
  }
  return `${firstTwo.join(', ')}, dan ${list.length - 2} lainnya`
})

const filteredPickerCatalogItems = computed(() => {
  let list = unregisteredCatalogItems.value || []
  if (catalogPickerFilter.value !== 'all') {
    list = list.filter(item => item.source === catalogPickerFilter.value)
  }
  const q = catalogPickerSearch.value.trim().toLowerCase()
  if (!q) return list
  return list.filter(
    item =>
      item.name.toLowerCase().includes(q) ||
      (item.category || '').toLowerCase().includes(q) ||
      (item.sku || '').toLowerCase().includes(q) ||
      (item.description || '').toLowerCase().includes(q)
  )
})

function openAddLayerModalWithCatalogItem(catItem: UnregisteredCatalogItem) {
  modalActiveTab.value = 'general'
  isEditing.value = false
  editingId.value = null
  const nextNo = props.presetItems.length > 0 ? Math.max(...props.presetItems.map(i => i.no)) + 1 : 1
  form.id = `layer_custom_${Date.now()}`
  form.no = nextNo
  
  // Auto-fill form fields based on the selected catalog item
  const suggestedLayer = catItem.meta?.layer || (catItem.category !== 'LICENSE' && catItem.category !== 'COMPUTE' ? catItem.category : catItem.name)
  form.layer = suggestedLayer
  form.component = catItem.name
  form.licenseFree = catItem.meta?.licenseFree || 'Community OSS'
  form.licenseEnterprise = catItem.meta?.licenseEnterprise || (catItem.category === 'LICENSE' ? catItem.name : 'Official Enterprise')
  form.notes = catItem.description || catItem.meta?.notes || `Formula arsitektur baseline untuk ${catItem.name}`
  form.contractRef = catItem.meta?.contractRef || (catItem.sku ? `SKU: ${catItem.sku}` : '')

  if (catItem.meta?.modBaseline) {
    const mb = catItem.meta.modBaseline
    form.moderate.nodeCount = mb.node_count ?? mb.nodeCount ?? 2
    form.moderate.vcpuPerNode = mb.vcpu_per_node ?? mb.vcpuPerNode ?? 4
    form.moderate.ramGBPerNode = mb.ram_gb_per_node ?? mb.ramGBPerNode ?? 8
    form.moderate.storageGBPerNode = mb.storage_gb_per_node ?? mb.storageGBPerNode ?? 100
    if (mb.compute_type || mb.computeType) form.moderate.computeType = mb.compute_type || mb.computeType
    if (mb.storage_class || mb.storageClass) form.moderate.storageClass = mb.storage_class || mb.storageClass
  } else {
    form.moderate = {
      nodeCount: 2,
      vcpuPerNode: 4,
      ramGBPerNode: 8,
      storageGBPerNode: 100,
      storageClass: 'Tier 2',
      computeType: 'vm',
      backupEnabled: true,
      licenseKind: 'free'
    }
  }

  if (catItem.meta?.critBaseline) {
    const cb = catItem.meta.critBaseline
    form.critical.nodeCount = cb.node_count ?? cb.nodeCount ?? 2
    form.critical.vcpuPerNode = cb.vcpu_per_node ?? cb.vcpuPerNode ?? 8
    form.critical.ramGBPerNode = cb.ram_gb_per_node ?? cb.ramGBPerNode ?? 16
    form.critical.storageGBPerNode = cb.storage_gb_per_node ?? cb.storageGBPerNode ?? 200
    if (cb.compute_type || cb.computeType) form.critical.computeType = cb.compute_type || cb.computeType
    if (cb.storage_class || cb.storageClass) form.critical.storageClass = cb.storage_class || cb.storageClass
  } else {
    form.critical = {
      nodeCount: 2,
      vcpuPerNode: 8,
      ramGBPerNode: 16,
      storageGBPerNode: 200,
      storageClass: 'Tier 2',
      computeType: 'vm',
      backupEnabled: true,
      licenseKind: 'enterprise'
    }
  }

  showModal.value = true
}

watch(
  () => unregisteredCatalogItems.value?.length || 0,
  count => {
    emit('update:unregistered-count', count)
  },
  { immediate: true }
)

// ==========================================
// AUTOCOMPLETE OPTIONS FROM CATALOG & PRESET
// ==========================================
const layerOptions = computed<AutocompleteOption[]>(() => {
  const set = new Set<string>()
  const list: AutocompleteOption[] = []

  // Standard architectural layers at Enterprise Architecture
  const standardLayers = [
    { label: 'Web Ingress', icon: '🌐', desc: 'Reverse proxy, SSL termination, and rate limiting entry layer.' },
    { label: 'App Container', icon: '📦', desc: 'Microservices orchestration cluster on OpenShift / Kubernetes.' },
    { label: 'Message Broker & Event Streaming', icon: '📨', desc: 'Enterprise event streaming (Apache Kafka / RabbitMQ).' },
    { label: 'Database Core OLTP', icon: '🗄️', desc: 'Core transaction relational database (Oracle RAC / PostgreSQL).' },
    { label: 'Database Read Replica', icon: '📖', desc: 'Read-only replica for reporting and query scalability.' },
    { label: 'Caching / In-Memory Store', icon: '⚡', desc: 'Low-latency in-memory data store (Redis Enterprise / KeyDB).' },
    { label: 'Search Engine & Analytics', icon: '🔍', desc: 'High-speed distributed text search (OpenSearch / Elasticsearch).' },
    { label: 'Log Shipper & Ingest', icon: '📋', desc: 'Log pipeline, parsing, and normalization (Logstash / Fluentbit).' },
    { label: 'Security WAF & API Gateway', icon: '🛡️', desc: 'API gateway, WAF security, authentication, and traffic throttling.' },
    { label: 'Object Storage & File Archive', icon: '🪣', desc: 'S3-compatible scalable unstructured storage (MinIO / Ceph).' },
    { label: 'CI/CD & Worker Runner', icon: '⚙️', desc: 'Automated CI/CD build agents and batch workers.' },
    { label: 'Batch Processing / ETL', icon: '🔄', desc: 'Scheduled ETL processing and large-scale data transformation.' },
    { label: 'Observability & Monitoring', icon: '📊', desc: 'Telemetry, metrics, trace collection, and alerting.' }
  ]

  for (const sl of standardLayers) {
    if (!set.has(sl.label)) {
      set.add(sl.label)
      list.push({
        label: sl.label,
        value: sl.label,
        icon: sl.icon,
        category: 'Layer Standar',
        badge: 'Enterprise',
        badgeType: 'primary',
        description: sl.desc
      })
    }
  }

  // Include any extra layers in presetItems
  for (const p of props.presetItems) {
    if (!set.has(p.layer)) {
      set.add(p.layer)
      list.push({
        label: p.layer,
        value: p.layer,
        icon: p.icon || '🏷️',
        category: 'Preset Aktif',
        badge: 'Katalog',
        description: p.notes || `Komponen: ${p.component}`
      })
    }
  }

  return list
})

const componentOptions = computed<AutocompleteOption[]>(() => {
  const set = new Set<string>()
  const list: AutocompleteOption[] = []

  // 1. From backend catalog components
  for (const bc of backendComponents.value) {
    if (!set.has(bc.component_name)) {
      set.add(bc.component_name)
      const freeLic = bc.licenses?.find(l => l.license_kind === 'free')?.product_name || ''
      const entLic = bc.licenses?.find(l => l.license_kind === 'enterprise')?.product_name || ''
      const contractRef = bc.licenses?.find(l => l.contract_ref)?.contract_ref || ''
      const modBaseline = bc.baselines?.find(b => b.tier === 'moderate')
      const critBaseline = bc.baselines?.find(b => b.tier === 'critical')

      list.push({
        label: bc.component_name,
        value: bc.component_name,
        icon: '📦',
        category: bc.layer_name,
        badge: 'Katalog Resmi',
        badgeType: 'primary',
        description: bc.description || bc.notes,
        meta: {
          layer: bc.layer_name,
          licenseFree: freeLic,
          licenseEnterprise: entLic,
          notes: bc.notes,
          contractRef,
          modBaseline,
          critBaseline
        }
      })
    }
  }

  // 2. From presetItems
  for (const p of props.presetItems) {
    if (!set.has(p.component)) {
      set.add(p.component)
      list.push({
        label: p.component,
        value: p.component,
        icon: p.icon || '⚙️',
        category: p.layer,
        badge: 'Preset Enterprise Architecture',
        description: p.notes || `Layer: ${p.layer}`,
        meta: {
          layer: p.layer,
          licenseFree: p.licenseFree,
          licenseEnterprise: p.licenseEnterprise,
          notes: p.notes,
          contractRef: p.contractRef,
          modBaseline: p.moderate,
          critBaseline: p.critical
        }
      })
    }
  }

  // 3. From pricingList (catalog items)
  if (props.catalogItems) {
    for (const pi of props.catalogItems) {
      if (!set.has(pi.name) && (pi.category === 'LICENSE' || pi.category === 'COMPUTE' || pi.category === 'SOFTWARE')) {
        set.add(pi.name)
        list.push({
          label: pi.name,
          value: pi.name,
          icon: '🏷️',
          category: pi.category,
          badge: pi.sku || 'SKU',
          description: `Katalog SKU: ${pi.sku} (${pi.environment})`
        })
      }
    }
  }

  return list
})

const licenseFreeOptions = computed<AutocompleteOption[]>(() => {
  const set = new Set<string>()
  const list: AutocompleteOption[] = []
  const defaults = [
    'Community OSS',
    'NGINX Community (OSS)',
    'OKD (Community OSS)',
    'PostgreSQL Community OSS',
    'Strimzi Kafka Community (OSS)',
    'Redis Open Source (BSD)',
    'OpenSearch Community (OSS)',
    'Logstash OSS Pipeline',
    'MinIO Community OSS',
    'ModSecurity Open Source',
    'Prometheus & Grafana OSS'
  ]
  for (const d of defaults) {
    set.add(d)
    list.push({ label: d, value: d, icon: '🆓', badge: 'Rp 0 (OSS)', badgeType: 'success' })
  }
  for (const p of props.presetItems) {
    if (p.licenseFree && !set.has(p.licenseFree)) {
      set.add(p.licenseFree)
      list.push({ label: p.licenseFree, value: p.licenseFree, icon: '🆓', badge: 'Preset Free' })
    }
  }
  return list
})

const licenseEnterpriseOptions = computed<AutocompleteOption[]>(() => {
  const set = new Set<string>()
  const list: AutocompleteOption[] = []
  for (const p of props.presetItems) {
    if (p.licenseEnterprise && !set.has(p.licenseEnterprise)) {
      set.add(p.licenseEnterprise)
      list.push({
        label: p.licenseEnterprise,
        value: p.licenseEnterprise,
        icon: '💼',
        category: 'Preset Enterprise',
        badge: 'Official',
        badgeType: 'warning'
      })
    }
  }
  if (props.catalogItems) {
    for (const pi of props.catalogItems) {
      if (pi.category === 'LICENSE' && !set.has(pi.name)) {
        set.add(pi.name)
        list.push({
          label: pi.name,
          value: pi.name,
          icon: '💼',
          category: 'Katalog SKU',
          badge: pi.sku,
          badgeType: 'primary',
          description: `Lisensi Resmi Enterprise Architecture (${pi.sku})`
        })
      }
    }
  }
  return list
})

function handleComponentSelect(option: AutocompleteOption) {
  if (option.meta) {
    if (option.meta.layer && (!form.layer || form.layer === '')) {
      form.layer = option.meta.layer
    }
    if (option.meta.licenseFree && (!form.licenseFree || form.licenseFree === 'Community OSS' || form.licenseFree === '')) {
      form.licenseFree = option.meta.licenseFree
    }
    if (option.meta.licenseEnterprise && (!form.licenseEnterprise || form.licenseEnterprise === 'Official Enterprise' || form.licenseEnterprise === '')) {
      form.licenseEnterprise = option.meta.licenseEnterprise
    }
    if (option.meta.notes && !form.notes) {
      form.notes = option.meta.notes
    }
    if (option.meta.contractRef && !form.contractRef) {
      form.contractRef = option.meta.contractRef
    }
    // If adding a new layer, automatically populate baselines from catalog if available
    if (!isEditing.value) {
      if (option.meta.modBaseline) {
        const mb = option.meta.modBaseline
        if (mb.node_count ?? mb.nodeCount) form.moderate.nodeCount = mb.node_count ?? mb.nodeCount
        if (mb.vcpu_per_node ?? mb.vcpuPerNode) form.moderate.vcpuPerNode = mb.vcpu_per_node ?? mb.vcpuPerNode
        if (mb.ram_gb_per_node ?? mb.ramGBPerNode) form.moderate.ramGBPerNode = mb.ram_gb_per_node ?? mb.ramGBPerNode
        if (mb.storage_gb_per_node ?? mb.storageGBPerNode) form.moderate.storageGBPerNode = mb.storage_gb_per_node ?? mb.storageGBPerNode
        if (mb.compute_type || mb.computeType) form.moderate.computeType = mb.compute_type || mb.computeType
        if (mb.storage_class || mb.storageClass) form.moderate.storageClass = mb.storage_class || mb.storageClass
      }
      if (option.meta.critBaseline) {
        const cb = option.meta.critBaseline
        if (cb.node_count ?? cb.nodeCount) form.critical.nodeCount = cb.node_count ?? cb.nodeCount
        if (cb.vcpu_per_node ?? cb.vcpuPerNode) form.critical.vcpuPerNode = cb.vcpu_per_node ?? cb.vcpuPerNode
        if (cb.ram_gb_per_node ?? cb.ramGBPerNode) form.critical.ramGBPerNode = cb.ram_gb_per_node ?? cb.ramGBPerNode
        if (cb.storage_gb_per_node ?? cb.storageGBPerNode) form.critical.storageGBPerNode = cb.storage_gb_per_node ?? cb.storageGBPerNode
        if (cb.compute_type || cb.computeType) form.critical.computeType = cb.compute_type || cb.computeType
        if (cb.storage_class || cb.storageClass) form.critical.storageClass = cb.storage_class || cb.storageClass
      }
    }
  }
}

const searchQuery = ref('')
const activeTierTab = ref<'both' | 'moderate' | 'critical'>('both')
const viewMode = ref<'cards' | 'table' | 'diagram'>('cards')
const selectedDiagramId = ref<string | null>(null)
const activeCategory = ref<string>('all')
const isOverviewBannerVisible = ref(true)
const showModal = ref(false)
const isEditing = ref(false)
const editingId = ref<string | null>(null)

// Restore saved preferences
if (typeof window !== 'undefined') {
  const savedMode = localStorage.getItem('mandiri_preset_view_mode')
  if (savedMode === 'cards' || savedMode === 'table' || savedMode === 'diagram') {
    viewMode.value = savedMode
  }
  const savedBanner = localStorage.getItem('mandiri_preset_overview_visible')
  if (savedBanner !== null) {
    isOverviewBannerVisible.value = savedBanner === 'true'
  }
}

watch(viewMode, (newVal) => {
  if (typeof window !== 'undefined') {
    localStorage.setItem('mandiri_preset_view_mode', newVal)
  }
})

function toggleOverviewBanner() {
  isOverviewBannerVisible.value = !isOverviewBannerVisible.value
  if (typeof window !== 'undefined') {
    localStorage.setItem('mandiri_preset_overview_visible', String(isOverviewBannerVisible.value))
  }
}

function formatDiskSize(gb: number): string {
  if (!gb || gb <= 0) return '0 GB'
  if (gb >= 1000) {
    const tb = (gb / 1000).toFixed(1)
    return `${tb.endsWith('.0') ? tb.slice(0, -2) : tb} TB`
  }
  return `${gb} GB`
}

// Cluster Totals (Moderat vs Critical) for Executive Overview
const clusterTotals = computed(() => {
  const items = props.presetItems || []
  let modNodes = 0
  let modVcpu = 0
  let modRam = 0
  let modStorage = 0

  let critNodes = 0
  let critVcpu = 0
  let critRam = 0
  let critStorage = 0

  for (const item of items) {
    if (item.moderate) {
      if (item.moderate.computeType !== 'shared') {
        const nodes = Number(item.moderate.nodeCount) || 0
        modNodes += nodes
        modVcpu += nodes * (Number(item.moderate.vcpuPerNode) || 0)
        modRam += nodes * (Number(item.moderate.ramGBPerNode) || 0)
      }
      modStorage += (Number(item.moderate.nodeCount) || 1) * (Number(item.moderate.storageGBPerNode) || 0)
    }
    if (item.critical) {
      if (item.critical.computeType !== 'shared') {
        const nodes = Number(item.critical.nodeCount) || 0
        critNodes += nodes
        critVcpu += nodes * (Number(item.critical.vcpuPerNode) || 0)
        critRam += nodes * (Number(item.critical.ramGBPerNode) || 0)
      }
      critStorage += (Number(item.critical.nodeCount) || 1) * (Number(item.critical.storageGBPerNode) || 0)
    }
  }

  return {
    modNodes,
    modVcpu,
    modRam,
    modStorage,
    critNodes,
    critVcpu,
    critRam,
    critStorage
  }
})

// Architecture Categories & Categorization
const CATEGORY_DEFINITIONS = [
  { id: 'all', label: 'Semua Layer', icon: '📑' },
  { id: 'access_app', label: 'Akses & Aplikasi', icon: '🌐' },
  { id: 'db_storage', label: 'Database & Penyimpanan', icon: '🗄️' },
  { id: 'cache_queue', label: 'Memori Cepat & Antrean', icon: '⚡' },
  { id: 'monitor_batch', label: 'Pemantauan & Penjadwalan', icon: '📊' },
  { id: 'other', label: 'Layer Tambahan', icon: '🔧' }
]

function getLayerCategory(item: SizingPresetItem): 'access_app' | 'db_storage' | 'cache_queue' | 'monitor_batch' | 'other' {
  const no = Number(item.no) || 0
  if (no === 1 || no === 2 || no === 3) return 'access_app'
  if (no === 4 || no === 5 || no === 12) return 'db_storage'
  if (no === 6 || no === 7 || no === 8) return 'cache_queue'
  if (no === 9 || no === 10 || no === 11 || no === 13) return 'monitor_batch'

  const text = `${item.layer} ${item.component} ${item.notes || ''}`.toLowerCase()
  if (text.includes('ingress') || text.includes('router') || text.includes('container') || text.includes('web') || text.includes('gateway') || text.includes('api')) {
    return 'access_app'
  }
  if (text.includes('database') || text.includes('storage') || text.includes('oracle') || text.includes('postgres') || text.includes('minio') || text.includes('s3') || text.includes('san')) {
    return 'db_storage'
  }
  if (text.includes('cache') || text.includes('redis') || text.includes('kafka') || text.includes('message') || text.includes('queue') || text.includes('rabbit')) {
    return 'cache_queue'
  }
  if (text.includes('log') || text.includes('monitor') || text.includes('batch') || text.includes('elastic') || text.includes('kibana') || text.includes('cron') || text.includes('audit')) {
    return 'monitor_batch'
  }
  return 'other'
}

const categoryCounts = computed(() => {
  const counts: Record<string, number> = {
    all: props.presetItems.length,
    access_app: 0,
    db_storage: 0,
    cache_queue: 0,
    monitor_batch: 0,
    other: 0
  }
  for (const item of props.presetItems) {
    const cat = getLayerCategory(item)
    counts[cat] = (counts[cat] || 0) + 1
  }
  return counts
})

const availableCategories = computed(() => {
  return CATEGORY_DEFINITIONS.filter(c => c.id === 'all' || (categoryCounts.value[c.id] || 0) > 0)
})

function getCategoryBadge(item: SizingPresetItem) {
  const cat = getLayerCategory(item)
  switch (cat) {
    case 'access_app': return { label: 'Akses & Aplikasi', icon: '🌐', theme: 'cat-app' }
    case 'db_storage': return { label: 'Database & Storage', icon: '🗄️', theme: 'cat-db' }
    case 'cache_queue': return { label: 'Memori & Antrean', icon: '⚡', theme: 'cat-cache' }
    case 'monitor_batch': return { label: 'Pemantauan & Batch', icon: '📊', theme: 'cat-monitor' }
    default: return { label: 'Layer Tambahan', icon: '🔧', theme: 'cat-other' }
  }
}

function getLayerIcon(item: SizingPresetItem): string {
  if (item.icon) return item.icon
  const defaultIcons: Record<number, string> = {
    1: '🌐',
    2: '📦',
    3: '🔀',
    4: '🗄️',
    5: '💾',
    6: '⚡',
    7: '🛡️',
    8: '📨',
    9: '🔍',
    10: '📋',
    11: '📊',
    12: '📁',
    13: '⏱️'
  }
  if (defaultIcons[item.no]) return defaultIcons[item.no]

  const text = `${item.layer} ${item.component}`.toLowerCase()
  if (text.includes('ingress') || text.includes('web') || text.includes('gateway')) return '🌐'
  if (text.includes('database') || text.includes('db') || text.includes('oracle') || text.includes('sql')) return '🗄️'
  if (text.includes('storage') || text.includes('san') || text.includes('minio') || text.includes('s3')) return '💾'
  if (text.includes('cache') || text.includes('redis')) return '⚡'
  if (text.includes('message') || text.includes('kafka') || text.includes('queue')) return '📨'
  if (text.includes('log') || text.includes('elastic') || text.includes('monitor')) return '📊'
  if (text.includes('container') || text.includes('app') || text.includes('pod')) return '📦'
  return '⚙️'
}

function getSizeBadge(spec?: LayerPresetSpec) {
  if (!spec) return { label: 'Standar', tagClass: 'badge-size-standard' }
  if (spec.computeType === 'shared') {
    return { label: 'Khusus Ruang Disk', tagClass: 'badge-size-storage' }
  }
  const vcpu = Number(spec.vcpuPerNode) || 0
  const ram = Number(spec.ramGBPerNode) || 0
  if (vcpu <= 2) {
    return { label: 'Ukuran Kecil', tagClass: 'badge-size-small' }
  }
  if (vcpu <= 4) {
    return { label: 'Ukuran Sedang', tagClass: 'badge-size-medium' }
  }
  if (vcpu <= 8 && ram <= 16) {
    return { label: 'Ukuran Standar', tagClass: 'badge-size-standard' }
  }
  if (vcpu <= 8 && ram > 16) {
    return { label: 'Ukuran Menengah+', tagClass: 'badge-size-plus' }
  }
  if (vcpu <= 32) {
    return { label: 'Ukuran Besar', tagClass: 'badge-size-large' }
  }
  return { label: 'Ukuran Raksasa', tagClass: 'badge-size-xlarge' }
}

function getDeltaPills(mod?: LayerPresetSpec, crit?: LayerPresetSpec): Array<{ text: string; type: 'upgrade' | 'scale' | 'license' | 'info' }> {
  if (!mod || !crit) return []
  const pills: Array<{ text: string; type: 'upgrade' | 'scale' | 'license' | 'info' }> = []

  const modNodes = Number(mod.nodeCount) || 0
  const critNodes = Number(crit.nodeCount) || 0
  if (critNodes > modNodes) {
    pills.push({ text: `↑ +${critNodes - modNodes} Unit Server`, type: 'scale' })
  } else if (critNodes < modNodes) {
    pills.push({ text: `↓ ${critNodes - modNodes} Unit Server`, type: 'info' })
  }

  if (mod.computeType !== crit.computeType) {
    if (crit.computeType === 'baremetal') {
      pills.push({ text: '↑ Server Fisik Dedicated (Baremetal)', type: 'upgrade' })
    } else if (crit.computeType === 'vm') {
      pills.push({ text: 'Virtual Machine (VM)', type: 'info' })
    }
  }

  if (mod.computeType !== 'shared' && crit.computeType !== 'shared') {
    const modTotalVcpu = modNodes * (Number(mod.vcpuPerNode) || 0)
    const critTotalVcpu = critNodes * (Number(crit.vcpuPerNode) || 0)
    if (modTotalVcpu > 0 && critTotalVcpu > modTotalVcpu) {
      const ratio = critTotalVcpu / modTotalVcpu
      if (ratio === 2) {
        pills.push({ text: '↑ Kapasitas 2x Lipat (2N)', type: 'scale' })
      } else if (ratio > 1) {
        pills.push({ text: `↑ Kapasitas +${Math.round((ratio - 1) * 100)}% vCPU`, type: 'scale' })
      }
    }

    const modTotalRam = modNodes * (Number(mod.ramGBPerNode) || 0)
    const critTotalRam = critNodes * (Number(crit.ramGBPerNode) || 0)
    if (modTotalRam > 0 && critTotalRam > modTotalRam) {
      const ramRatio = critTotalRam / modTotalRam
      const vcpuRatio = modTotalVcpu > 0 ? (critTotalVcpu / modTotalVcpu) : 1
      if (Math.abs(ramRatio - vcpuRatio) > 0.1) {
        pills.push({ text: `↑ RAM +${Math.round((ramRatio - 1) * 100)}%`, type: 'scale' })
      }
    }
  }

  const modTotalStorage = (modNodes || 1) * (Number(mod.storageGBPerNode) || 0)
  const critTotalStorage = (critNodes || 1) * (Number(crit.storageGBPerNode) || 0)
  if (critTotalStorage > modTotalStorage) {
    const diff = critTotalStorage - modTotalStorage
    pills.push({ text: `↑ Ruang Disk +${formatDiskSize(diff)}`, type: 'upgrade' })
  }
  if (mod.storageClass !== crit.storageClass) {
    pills.push({ text: `↑ Kecepatan Disk ${crit.storageClass}`, type: 'upgrade' })
  }

  if (mod.licenseKind !== crit.licenseKind && crit.licenseKind === 'enterprise') {
    pills.push({ text: '🛡️ Lisensi Resmi Vendor (SLA 24/7)', type: 'license' })
  }

  if (pills.length === 0) {
    pills.push({ text: 'Spesifikasi Hardware & Lisensi Identik', type: 'info' })
  }

  return pills
}

const form = reactive({
  id: '',
  no: 14,
  layer: '',
  component: '',
  licenseFree: '',
  licenseEnterprise: '',
  notes: '',
  contractRef: '',
  moderate: {
    nodeCount: 2,
    vcpuPerNode: 8,
    ramGBPerNode: 16,
    storageGBPerNode: 200,
    storageClass: 'Tier 2' as 'Tier 2' | 'Tier 3',
    computeType: 'vm' as 'vm' | 'baremetal' | 'shared',
    backupEnabled: true,
    licenseKind: 'free' as 'free' | 'enterprise'
  },
  critical: {
    nodeCount: 2,
    vcpuPerNode: 16,
    ramGBPerNode: 32,
    storageGBPerNode: 500,
    storageClass: 'Tier 2' as 'Tier 2' | 'Tier 3',
    computeType: 'vm' as 'vm' | 'baremetal' | 'shared',
    backupEnabled: true,
    licenseKind: 'enterprise' as 'free' | 'enterprise'
  }
})

const filteredPresetItems = computed(() => {
  let list = props.presetItems || []
  if (activeCategory.value !== 'all') {
    list = list.filter(item => getLayerCategory(item) === activeCategory.value)
  }
  if (!searchQuery.value.trim()) {
    return list
  }
  const q = searchQuery.value.toLowerCase().trim()
  return list.filter(item => {
    return (
      item.layer.toLowerCase().includes(q) ||
      item.component.toLowerCase().includes(q) ||
      item.licenseFree.toLowerCase().includes(q) ||
      item.licenseEnterprise.toLowerCase().includes(q) ||
      (item.notes || '').toLowerCase().includes(q)
    )
  })
})

// ==========================================
// HIGH-LEVEL ARCHITECTURE DIAGRAM HELPERS
// ==========================================
const selectedDiagramItem = computed<SizingPresetItem | null>(() => {
  if (!selectedDiagramId.value) return null
  return (props.presetItems || []).find(item => item.id === selectedDiagramId.value) || null
})

function selectDiagramNode(item: SizingPresetItem) {
  if (selectedDiagramId.value === item.id) {
    selectedDiagramId.value = null
  } else {
    selectedDiagramId.value = item.id
  }
}

function getNodeArchitectureRole(item: SizingPresetItem): string {
  const no = Number(item.no) || 0
  switch (no) {
    case 1: return 'DMZ Reverse Proxy'
    case 2: return 'Microservices Worker'
    case 3: return 'Ingress & Service Mesh'
    case 4: return 'Core OLTP Engine'
    case 5: return 'High-IOPS SAN LUN'
    case 6: return 'In-Memory Speed Layer'
    case 7: return 'HA Failover Quorum'
    case 8: return 'Event Streaming Bus'
    case 9: return 'Search & Log Indexer'
    case 10: return 'Log Ingest Pipeline'
    case 11: return 'Observability Portal'
    case 12: return 'Object & ECM Archive'
    case 13: return 'Workload Automation'
  }
  const cat = getLayerCategory(item)
  if (cat === 'access_app') return 'App & Access Tier'
  if (cat === 'db_storage') return 'Data Persistence'
  if (cat === 'cache_queue') return 'Middleware & Queue'
  if (cat === 'monitor_batch') return 'Telemetry & Ops'
  return 'Custom Extension'
}

function getIntraZoneConnector(
  zoneId: string,
  currentItem: SizingPresetItem,
  nextItem?: SizingPresetItem
): { label: string; style: 'sync' | 'async' | 'san' } | null {
  if (!nextItem) return null
  const curNo = Number(currentItem.no) || 0
  const nextNo = Number(nextItem.no) || 0

  if (zoneId === 'edge_ingress') {
    if (curNo === 1 && nextNo === 3) {
      return { label: 'HTTPS / Route Proxy', style: 'sync' }
    }
    return { label: 'Traffic Forward', style: 'sync' }
  }

  if (zoneId === 'middleware_mesh') {
    if (curNo === 6 && nextNo === 7) {
      return { label: 'HA Quorum Watch', style: 'async' }
    }
    if ((curNo === 6 || curNo === 7) && nextNo === 8) {
      return { label: 'Pub/Sub Stream', style: 'async' }
    }
  }

  if (zoneId === 'data_persistence') {
    if (curNo === 4 && nextNo === 5) {
      return { label: 'Fibre Channel ASM LUN', style: 'san' }
    }
  }

  if (zoneId === 'ops_observability') {
    if (curNo === 10 && nextNo === 9) {
      return { label: 'Parse & Index Stream', style: 'async' }
    }
    if (curNo === 9 && nextNo === 11) {
      return { label: 'Visual Query API', style: 'sync' }
    }
  }

  return null
}

interface DiagramZoneGroup {
  id: string
  code: string
  icon: string
  title: string
  architectureTier: string
  subtitle: string
  themeClass: string
  nextFlowLabel: string
  spineStyle: 'sync' | 'async'
  items: SizingPresetItem[]
  modNodes: number
  modVcpu: number
  critNodes: number
  critVcpu: number
}

const diagramZones = computed<DiagramZoneGroup[]>(() => {
  const items = filteredPresetItems.value || []

  const isEdgeIngress = (item: SizingPresetItem) => {
    const no = Number(item.no) || 0
    if (no === 1 || no === 3) return true
    if (no === 2) return false
    if (getLayerCategory(item) !== 'access_app') return false
    const text = `${item.layer} ${item.component}`.toLowerCase()
    return text.includes('ingress') || text.includes('router') || text.includes('gateway') || text.includes('waf') || text.includes('proxy') || text.includes('load balancer')
  }

  const isCoreCompute = (item: SizingPresetItem) => {
    const no = Number(item.no) || 0
    if (no === 2) return true
    if (no === 1 || no === 3) return false
    return getLayerCategory(item) === 'access_app' && !isEdgeIngress(item)
  }

  const edgeItems = items.filter(isEdgeIngress).sort((a, b) => a.no - b.no)
  const computeItems = items.filter(isCoreCompute).sort((a, b) => a.no - b.no)
  const middlewareItems = items.filter(i => getLayerCategory(i) === 'cache_queue').sort((a, b) => a.no - b.no)
  const dataItems = items.filter(i => getLayerCategory(i) === 'db_storage').sort((a, b) => a.no - b.no)

  // Order Observability logically: L10 (Log Shipper) -> L09 (Log Analytics) -> L11 (Log Dashboard) -> L13 (Batch Engine) -> others
  const opsOrderWeight = (no: number) => {
    if (no === 10) return 1
    if (no === 9) return 2
    if (no === 11) return 3
    if (no === 13) return 4
    return 10 + no
  }
  const opsItems = items
    .filter(i => getLayerCategory(i) === 'monitor_batch')
    .sort((a, b) => opsOrderWeight(a.no) - opsOrderWeight(b.no))

  const customOtherItems = items.filter(i => getLayerCategory(i) === 'other').sort((a, b) => a.no - b.no)

  const rawZones = [
    {
      id: 'edge_ingress',
      code: 'ZONE 01',
      icon: '🌐',
      title: 'Edge, DMZ & Traffic Ingress Layer',
      architectureTier: 'Network Entry & Routing',
      subtitle: 'Pintu gerbang utama lalu lintas nasabah, terminasi SSL/TLS, proteksi WAF, dan internal mesh router.',
      themeClass: 'zone-theme-edge',
      nextFlowLabel: '⬇ Routed HTTP/gRPC Request ke Klaster Aplikasi Inti',
      spineStyle: 'sync' as const,
      items: edgeItems
    },
    {
      id: 'core_compute',
      code: 'ZONE 02',
      icon: '📦',
      title: 'Core Application & Microservices Runtime',
      architectureTier: 'Compute Orchestration Hub',
      subtitle: 'Klaster eksekusi utama tempat seluruh container microservices dan logika bisnis transaksi berjalan.',
      themeClass: 'zone-theme-compute',
      nextFlowLabel: '⬇ Akses Sub-ms In-Memory Cache & Event Streaming Asinkron',
      spineStyle: 'sync' as const,
      items: computeItems
    },
    {
      id: 'middleware_mesh',
      code: 'ZONE 03',
      icon: '⚡',
      title: 'In-Memory Speed & Event Streaming Mesh',
      architectureTier: 'High-Speed Middleware',
      subtitle: 'Akselerasi baca/tulis instan (Redis Cluster + Sentinel Quorum) dan antrean transaksi berkapasitas tinggi (Kafka).',
      themeClass: 'zone-theme-middleware',
      nextFlowLabel: '⬇ Transaksi Finansial OLTP Read/Write & Penyimpanan Persisten',
      spineStyle: 'sync' as const,
      items: middlewareItems
    },
    {
      id: 'data_persistence',
      code: 'ZONE 04',
      icon: '🗄️',
      title: 'Core Database & High-IOPS Persistent Storage',
      architectureTier: 'Enterprise Data Persistence',
      subtitle: 'Penyimpanan transaksi inti berketahanan tinggi (Oracle RAC + ASM SAN LUN Tier 3) serta arsip dokumen objek.',
      themeClass: 'zone-theme-data',
      nextFlowLabel: '⇣ Pengiriman Log Telemetry, Audit Trail & Eksekusi Batch EOD',
      spineStyle: 'async' as const,
      items: dataItems
    },
    {
      id: 'ops_observability',
      code: 'ZONE 05',
      icon: '📊',
      title: 'Observability Pipeline & Batch Automation',
      architectureTier: 'Telemetry, Audit & Operations',
      subtitle: 'Pipa pengumpulan log (Logstash ➔ Elasticsearch ➔ Kibana) serta otomasi penjadwalan tutup buku harian (Control-M).',
      themeClass: 'zone-theme-ops',
      nextFlowLabel: '⇣ Integrasi Ekstensi Tambahan',
      spineStyle: 'async' as const,
      items: opsItems
    },
    {
      id: 'custom_other',
      code: 'ZONE 06',
      icon: '🔧',
      title: 'Extended & Custom Architecture Layers',
      architectureTier: 'Custom Extensions',
      subtitle: 'Komponen arsitektur tambahan sesuai kebutuhan spesifik proyek.',
      themeClass: 'zone-theme-other',
      nextFlowLabel: '',
      spineStyle: 'async' as const,
      items: customOtherItems
    }
  ]

  return rawZones
    .filter(z => z.items.length > 0)
    .map(z => {
      let modNodes = 0
      let modVcpu = 0
      let critNodes = 0
      let critVcpu = 0
      for (const it of z.items) {
        if (it.moderate && it.moderate.computeType !== 'shared') {
          const n = Number(it.moderate.nodeCount) || 0
          modNodes += n
          modVcpu += n * (Number(it.moderate.vcpuPerNode) || 0)
        }
        if (it.critical && it.critical.computeType !== 'shared') {
          const n = Number(it.critical.nodeCount) || 0
          critNodes += n
          critVcpu += n * (Number(it.critical.vcpuPerNode) || 0)
        }
      }
      return {
        ...z,
        modNodes,
        modVcpu,
        critNodes,
        critVcpu
      }
    })
})

function formatComputeType(type: string): string {
  if (type === 'baremetal') return 'Dedicated Baremetal'
  if (type === 'shared') return 'Shared / SAN LUN'
  return 'Virtual Machine (VM)'
}

function copyModerateToCritical() {
  form.critical.nodeCount = Math.max(form.moderate.nodeCount, 2)
  form.critical.vcpuPerNode = form.moderate.vcpuPerNode
  form.critical.ramGBPerNode = form.moderate.ramGBPerNode
  form.critical.storageGBPerNode = form.moderate.storageGBPerNode
  form.critical.storageClass = form.moderate.storageClass
  form.critical.computeType = form.moderate.computeType
  form.critical.backupEnabled = form.moderate.backupEnabled
}


function openAddLayerModal() {
  modalActiveTab.value = 'general'
  isEditing.value = false
  editingId.value = null
  const nextNo = props.presetItems.length > 0 ? Math.max(...props.presetItems.map(i => i.no)) + 1 : 1
  form.id = `layer_custom_${Date.now()}`
  form.no = nextNo
  form.layer = ''
  form.component = ''
  form.licenseFree = 'Community OSS'
  form.licenseEnterprise = 'Official Enterprise'
  form.notes = ''
  form.contractRef = ''
  form.moderate = {
    nodeCount: 2,
    vcpuPerNode: 4,
    ramGBPerNode: 8,
    storageGBPerNode: 100,
    storageClass: 'Tier 2',
    computeType: 'vm',
    backupEnabled: true,
    licenseKind: 'free'
  }
  form.critical = {
    nodeCount: 2,
    vcpuPerNode: 8,
    ramGBPerNode: 16,
    storageGBPerNode: 200,
    storageClass: 'Tier 2',
    computeType: 'vm',
    backupEnabled: true,
    licenseKind: 'enterprise'
  }
  showModal.value = true
}

function openEditLayerModal(item: SizingPresetItem) {
  modalActiveTab.value = 'general'
  isEditing.value = true
  editingId.value = item.id
  form.id = item.id
  form.no = item.no
  form.layer = item.layer
  form.component = item.component
  form.licenseFree = item.licenseFree
  form.licenseEnterprise = item.licenseEnterprise
  form.notes = item.notes || ''
  form.contractRef = item.contractRef || ''
  form.moderate = { ...item.moderate }
  form.critical = { ...item.critical }
  showModal.value = true
}

function savePresetLayer() {
  const itemToSave: SizingPresetItem = {
    id: form.id,
    no: form.no,
    layer: form.layer.trim(),
    component: form.component.trim(),
    licenseFree: form.licenseFree.trim(),
    licenseEnterprise: form.licenseEnterprise.trim(),
    notes: form.notes.trim(),
    contractRef: form.contractRef.trim(),
    moderate: { ...form.moderate },
    critical: { ...form.critical }
  }

  emit('save-preset-item', itemToSave)
  showModal.value = false
}

function confirmDeleteCustomLayer(item: SizingPresetItem) {
  const isBuiltIn = item.no <= 13
  const msg = isBuiltIn
    ? `Apakah Anda yakin ingin menghapus layer preset bawaan "${item.layer} (${item.component})"?\n\nCatatan: Anda dapat mengembalikannya kapan saja melalui tombol "Kembalikan Default Pabrik".`
    : `Apakah Anda yakin ingin menghapus layer preset kustom "${item.layer} (${item.component})"?`
  if (confirm(msg)) {
    emit('delete-preset-item', item.id)
  }
}
</script>

<style scoped>
.preset-view-container {
  display: flex;
  flex-direction: column;
  flex: 1;
  height: 100%;
  min-height: 0;
  overflow: hidden;
  background: #F8FAFC;
}

.preset-topbar {
  background: #FFFFFF;
  border-bottom: 1px solid #E2E8F0;
  padding: 14px 24px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-shrink: 0;
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

.topbar-actions {
  display: flex;
  align-items: center;
  gap: 10px;
}

.btn-unmapped-trigger {
  display: inline-flex;
  align-items: center;
  gap: 7px;
  padding: 6px 12px;
  background: #F8FAFC;
  border: 1px solid #CBD5E1;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 600;
  color: #334155;
  cursor: pointer;
  transition: all 0.15s ease;
}

.btn-unmapped-trigger:hover {
  background: #F1F5F9;
  border-color: #94A3B8;
  color: #0F172A;
}

.unmapped-dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: #D97706;
  flex-shrink: 0;
}

.unmapped-count {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 18px;
  height: 18px;
  padding: 0 5px;
  border-radius: 4px;
  background: #E2E8F0;
  color: #0F172A;
  font-size: 11px;
  font-weight: 700;
}

.btn-reset-defaults {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 7px 14px;
  background: #FFFFFF;
  border: 1px solid #CBD5E1;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 600;
  color: #475569;
  cursor: pointer;
  transition: all 0.15s ease;
}

.btn-reset-defaults:hover {
  background: #F1F5F9;
  color: #002060;
  border-color: #002060;
}

.btn-add-layer {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 7px 16px;
  background: #002060;
  color: #FFFFFF;
  border: none;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.15s ease;
}

.btn-add-layer:hover {
  background: #003399;
}

/* ─────────────────────────────────────────────
   Workspace & Toolbar
   ───────────────────────────────────────────── */
.preset-workspace {
  flex: 1;
  min-height: 0;
  display: flex;
  flex-direction: column;
  padding: 20px 24px 80px 24px;
  overflow-y: auto;
  overflow-x: hidden;
  gap: 14px;
}

/* Clear visible scrollbars for workspace */
.preset-workspace::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

.preset-workspace::-webkit-scrollbar-track {
  background: #F1F5F9;
  border-radius: 4px;
}

.preset-workspace::-webkit-scrollbar-thumb {
  background: #CBD5E1;
  border-radius: 4px;
}

.preset-workspace::-webkit-scrollbar-thumb:hover {
  background: #94A3B8;
}

/* ─────────────────────────────────────────────
   Compact Single-Line Catalog Sync Bar
   ───────────────────────────────────────────── */
.catalog-sync-strip {
  background: #FFFFFF;
  border: 1px solid #E2E8F0;
  border-left: 3px solid #002060;
  border-radius: 8px;
  padding: 9px 14px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  box-shadow: 0 1px 2px rgba(15, 23, 42, 0.03);
  flex-shrink: 0;
}

.sync-strip-left {
  display: flex;
  align-items: center;
  gap: 10px;
  min-width: 0;
  flex: 1;
}

.sync-status-icon {
  width: 28px;
  height: 28px;
  border-radius: 6px;
  background: #F1F5F9;
  border: 1px solid #E2E8F0;
  color: #002060;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.sync-strip-copy {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
  flex: 1;
  overflow: hidden;
  white-space: nowrap;
}

.sync-strip-title {
  font-size: 12.5px;
  color: #1E293B;
  flex-shrink: 0;
}

.sync-strip-title strong {
  font-weight: 700;
  color: #0F172A;
}

.sync-strip-sep {
  color: #CBD5E1;
  flex-shrink: 0;
}

.sync-strip-preview {
  font-size: 12px;
  color: #64748B;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.sync-strip-actions {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-shrink: 0;
}

.btn-sync-review {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 5px 11px;
  background: #F8FAFC;
  border: 1px solid #CBD5E1;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 600;
  color: #0F172A;
  cursor: pointer;
  transition: all 0.15s ease;
}

.btn-sync-review:hover {
  background: #002060;
  border-color: #002060;
  color: #FFFFFF;
}

.btn-sync-review:hover .review-count-pill {
  background: rgba(255, 255, 255, 0.2);
  color: #FFFFFF;
}

.review-count-pill {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 1px 6px;
  border-radius: 4px;
  background: #E2E8F0;
  color: #334155;
  font-size: 11px;
  font-weight: 700;
  transition: all 0.15s ease;
}

.btn-sync-close {
  width: 28px;
  height: 28px;
  border-radius: 6px;
  background: transparent;
  border: 1px solid transparent;
  color: #94A3B8;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.15s ease;
}

.btn-sync-close:hover {
  background: #F1F5F9;
  color: #475569;
}

/* ─────────────────────────────────────────────
   Clean Catalog Picker Modal
   ───────────────────────────────────────────── */
.catalog-picker-modal {
  background: #FFFFFF;
  width: 680px;
  max-width: 94vw;
  max-height: 82vh;
  border-radius: 10px;
  border: 1px solid #E2E8F0;
  box-shadow: 0 20px 25px -5px rgba(15, 23, 42, 0.12), 0 8px 10px -6px rgba(15, 23, 42, 0.06);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.picker-modal-hdr {
  padding: 16px 20px;
  border-bottom: 1px solid #E2E8F0;
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
}

.picker-title {
  margin: 0;
  font-size: 15px;
  font-weight: 700;
  color: #0F172A;
  letter-spacing: -0.01em;
}

.picker-subtitle {
  margin: 3px 0 0 0;
  font-size: 12px;
  color: #64748B;
  line-height: 1.4;
}

.picker-toolbar {
  padding: 10px 20px;
  background: #F8FAFC;
  border-bottom: 1px solid #E2E8F0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.picker-search-box {
  display: flex;
  align-items: center;
  gap: 8px;
  background: #FFFFFF;
  border: 1px solid #CBD5E1;
  border-radius: 6px;
  padding: 0 10px;
  height: 32px;
  flex: 1;
}

.picker-search-input {
  border: none;
  outline: none;
  background: transparent;
  font-size: 12px;
  color: #0F172A;
  width: 100%;
}

.picker-search-clear {
  border: none;
  background: transparent;
  color: #94A3B8;
  font-size: 11px;
  cursor: pointer;
  padding: 0;
}

.picker-filter-tabs {
  display: inline-flex;
  background: #E2E8F0;
  padding: 2px;
  border-radius: 6px;
  gap: 2px;
}

.picker-tab-btn {
  padding: 4px 10px;
  border: none;
  background: transparent;
  border-radius: 4px;
  font-size: 11.5px;
  font-weight: 600;
  color: #475569;
  cursor: pointer;
  transition: all 0.15s ease;
}

.picker-tab-btn.active {
  background: #FFFFFF;
  color: #0F172A;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.06);
}

.picker-list-body {
  flex: 1;
  overflow-y: auto;
  divide-y: 1px solid #F1F5F9;
  max-height: 420px;
}

.picker-empty-state {
  padding: 40px 20px;
  text-align: center;
  font-size: 12.5px;
  color: #64748B;
}

.picker-row-item {
  padding: 12px 20px;
  border-bottom: 1px solid #F1F5F9;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  transition: background 0.12s ease;
}

.picker-row-item:hover {
  background: #F8FAFC;
}

.picker-row-main {
  display: flex;
  flex-direction: column;
  gap: 3px;
  min-width: 0;
  flex: 1;
}

.picker-row-top {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.picker-item-name {
  font-size: 13px;
  font-weight: 600;
  color: #0F172A;
}

.picker-item-cat {
  font-size: 11px;
  font-weight: 500;
  color: #475569;
  background: #F1F5F9;
  border: 1px solid #E2E8F0;
  padding: 1px 6px;
  border-radius: 4px;
}

.picker-item-sku {
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 10.5px;
  color: #64748B;
  background: #F8FAFC;
  border: 1px solid #E2E8F0;
  padding: 1px 6px;
  border-radius: 4px;
}

.picker-item-desc {
  margin: 0;
  font-size: 11.5px;
  color: #64748B;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.picker-row-actions {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
}

.btn-picker-ignore {
  padding: 5px 10px;
  background: transparent;
  border: 1px solid transparent;
  border-radius: 5px;
  font-size: 11.5px;
  font-weight: 500;
  color: #64748B;
  cursor: pointer;
  transition: all 0.15s ease;
}

.btn-picker-ignore:hover {
  background: #F1F5F9;
  color: #0F172A;
}

.btn-picker-select {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  padding: 5px 11px;
  background: #FFFFFF;
  border: 1px solid #CBD5E1;
  border-radius: 5px;
  font-size: 11.5px;
  font-weight: 600;
  color: #002060;
  cursor: pointer;
  transition: all 0.15s ease;
}

.btn-picker-select:hover {
  background: #002060;
  border-color: #002060;
  color: #FFFFFF;
}

.picker-modal-ftr {
  padding: 12px 20px;
  background: #F8FAFC;
  border-top: 1px solid #E2E8F0;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.picker-ftr-hint {
  font-size: 11.5px;
  color: #64748B;
}

.btn-restore-ignored {
  border: none;
  background: transparent;
  color: #002060;
  font-size: 11.5px;
  font-weight: 600;
  cursor: pointer;
  padding: 0;
  text-decoration: underline;
}

/* Banner Slide Transition */
.banner-slide-enter-active,
.banner-slide-leave-active {
  transition: all 0.2s ease;
}

.banner-slide-enter-from,
.banner-slide-leave-to {
  opacity: 0;
  transform: translateY(-6px);
}

/* ─────────────────────────────────────────────
   Executive Guidance Banner (Moderat vs Critical)
   ───────────────────────────────────────────── */
.preset-overview-section {
  background: #FFFFFF;
  border: 1px solid #E2E8F0;
  border-radius: 10px;
  padding: 14px 18px;
  box-shadow: 0 1px 3px rgba(15, 23, 42, 0.04);
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.overview-header-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.overview-title-group {
  display: flex;
  align-items: baseline;
  gap: 10px;
}

.overview-eyebrow {
  font-size: 10px;
  font-weight: 800;
  color: #002060;
  background: #EFF6FF;
  border: 1px solid #BFDBFE;
  padding: 2px 7px;
  border-radius: 4px;
  letter-spacing: 0.05em;
  text-transform: uppercase;
}

.overview-heading {
  margin: 0;
  font-size: 14px;
  font-weight: 800;
  color: #0F172A;
  letter-spacing: -0.01em;
}

.btn-toggle-overview {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  background: #F8FAFC;
  border: 1px solid #CBD5E1;
  border-radius: 6px;
  padding: 4px 10px;
  font-size: 11px;
  font-weight: 600;
  color: #475569;
  cursor: pointer;
  transition: all 0.15s ease;
}

.btn-toggle-overview:hover {
  background: #F1F5F9;
  color: #002060;
  border-color: #94A3B8;
}

.rotate-180 {
  transform: rotate(180deg);
}

.overview-cards-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 14px;
}

.overview-tier-card {
  border-radius: 8px;
  padding: 14px 16px;
  border: 1px solid #E2E8F0;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.overview-tier-card.card-moderat {
  background: linear-gradient(180deg, #F0F7FF 0%, #F8FAFC 100%);
  border-color: #BAE6FD;
}

.overview-tier-card.card-critical {
  background: linear-gradient(180deg, #FFFBEB 0%, #FFFDF7 100%);
  border-color: #FDE68A;
}

.card-tier-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.tier-pill-group {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.overview-tier-badge {
  font-size: 11px;
  font-weight: 800;
  padding: 3px 8px;
  border-radius: 4px;
  letter-spacing: 0.02em;
}

.overview-tier-badge.mod {
  background: #0284C7;
  color: #FFFFFF;
}

.overview-tier-badge.crit {
  background: #D97706;
  color: #FFFFFF;
}

.overview-target-badge {
  font-size: 11px;
  font-weight: 600;
  color: #475569;
}

.overview-tier-desc {
  margin: 0;
  font-size: 12px;
  color: #334155;
  line-height: 1.5;
}

.overview-totals-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 8px;
  background: #FFFFFF;
  border: 1px solid rgba(0, 0, 0, 0.06);
  border-radius: 6px;
  padding: 8px 10px;
}

.overview-total-item {
  display: flex;
  flex-direction: column;
}

.tot-label {
  font-size: 9.5px;
  font-weight: 700;
  color: #64748B;
  text-transform: uppercase;
  letter-spacing: 0.02em;
}

.tot-val {
  font-size: 13px;
  font-weight: 800;
}

.overview-collapse-enter-active,
.overview-collapse-leave-active {
  transition: all 0.2s ease;
}

.overview-collapse-enter-from,
.overview-collapse-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}

/* ─────────────────────────────────────────────
   Smart Toolbar & Filters
   ───────────────────────────────────────────── */
.preset-smart-toolbar {
  background: #FFFFFF;
  border: 1px solid #E2E8F0;
  border-radius: 10px;
  padding: 12px 16px;
  box-shadow: 0 1px 3px rgba(15, 23, 42, 0.03);
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.toolbar-primary-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 14px;
  flex-wrap: wrap;
}

.search-input-wrap {
  display: flex;
  align-items: center;
  background: #F8FAFC;
  border: 1px solid #CBD5E1;
  border-radius: 6px;
  padding: 0 12px;
  gap: 8px;
  flex: 1;
  min-width: 280px;
  max-width: 440px;
  height: 36px;
}

.preset-search-input {
  border: none;
  background: transparent;
  outline: none;
  font-size: 12.5px;
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

.tier-view-segmented {
  display: flex;
  align-items: center;
  gap: 8px;
}

.seg-label {
  font-size: 12px;
  font-weight: 700;
  color: #475569;
}

.segmented-control {
  display: inline-flex;
  align-items: center;
  background: #F1F5F9;
  padding: 3px;
  border-radius: 6px;
  border: 1px solid #CBD5E1;
  gap: 2px;
}

.seg-btn {
  padding: 5px 12px;
  font-size: 11.5px;
  font-weight: 600;
  color: #64748B;
  background: transparent;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.15s ease;
}

.seg-btn.active {
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

/* View Mode Switcher */
.view-mode-switcher {
  display: inline-flex;
  align-items: center;
  background: #F1F5F9;
  padding: 3px;
  border-radius: 6px;
  border: 1px solid #CBD5E1;
  gap: 2px;
}

.view-switch-btn {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  padding: 5px 11px;
  font-size: 11.5px;
  font-weight: 600;
  color: #64748B;
  background: transparent;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.15s ease;
}

.view-switch-btn.active {
  background: #002060;
  color: #FFFFFF;
  font-weight: 700;
  box-shadow: 0 1px 3px rgba(0, 32, 96, 0.25);
}

/* Category Pills Row */
.category-pills-row {
  display: flex;
  align-items: center;
  gap: 8px;
  overflow-x: auto;
  padding-top: 4px;
  border-top: 1px solid #F1F5F9;
}

.category-pill-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 5px 12px;
  background: #F8FAFC;
  border: 1px solid #E2E8F0;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 600;
  color: #475569;
  cursor: pointer;
  white-space: nowrap;
  transition: all 0.15s ease;
}

.category-pill-btn:hover {
  background: #F1F5F9;
  border-color: #CBD5E1;
  color: #0F172A;
}

.category-pill-btn.active {
  background: #002060;
  border-color: #002060;
  color: #FFFFFF;
  font-weight: 700;
  box-shadow: 0 2px 4px rgba(0, 32, 96, 0.15);
}

.cat-pill-count {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 18px;
  height: 18px;
  padding: 0 5px;
  border-radius: 9px;
  background: rgba(0, 0, 0, 0.06);
  font-size: 10.5px;
  font-weight: 700;
}

.category-pill-btn.active .cat-pill-count {
  background: rgba(255, 255, 255, 0.22);
  color: #FFFFFF;
}

/* ─────────────────────────────────────────────
   Visual Cards Grid
   ───────────────────────────────────────────── */
.preset-cards-grid {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.preset-card {
  background: #FFFFFF;
  border: 1px solid #E2E8F0;
  border-radius: 10px;
  padding: 16px 20px;
  display: flex;
  flex-direction: column;
  gap: 14px;
  box-shadow: 0 1px 3px rgba(15, 23, 42, 0.03);
  transition: all 0.2s ease;
}

.preset-card:hover {
  border-color: #CBD5E1;
  box-shadow: 0 4px 12px -2px rgba(15, 23, 42, 0.06);
}

.preset-card-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  border-bottom: 1px solid #F1F5F9;
  padding-bottom: 12px;
  gap: 16px;
}

.card-header-left {
  display: flex;
  align-items: flex-start;
  gap: 14px;
  flex: 1;
}

.card-avatar-box {
  width: 44px;
  height: 44px;
  border-radius: 10px;
  background: #F1F5F9;
  border: 1px solid #E2E8F0;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  font-size: 22px;
  box-shadow: inset 0 1px 2px rgba(0, 0, 0, 0.04);
}

.layer-title-meta-wrap {
  display: flex;
  flex-direction: column;
  gap: 4px;
  flex: 1;
}

.layer-badges-line {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.layer-number-badge {
  font-size: 10.5px;
  font-weight: 800;
  color: #002060;
  background: #EFF6FF;
  border: 1px solid #DBEAFE;
  padding: 2px 7px;
  border-radius: 4px;
  text-transform: uppercase;
  letter-spacing: 0.03em;
}

.category-badge-chip {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-size: 11px;
  font-weight: 700;
  padding: 2px 8px;
  border-radius: 4px;
}

.category-badge-chip.cat-app {
  background: #EFF6FF;
  color: #1D4ED8;
  border: 1px solid #BFDBFE;
}

.category-badge-chip.cat-db {
  background: #F5F3FF;
  color: #6D28D9;
  border: 1px solid #DDD6FE;
}

.category-badge-chip.cat-cache {
  background: #FEF3C7;
  color: #B45309;
  border: 1px solid #FDE68A;
}

.category-badge-chip.cat-monitor {
  background: #ECFDF5;
  color: #047857;
  border: 1px solid #A7F3D0;
}

.category-badge-chip.cat-other {
  background: #F1F5F9;
  color: #475569;
  border: 1px solid #CBD5E1;
}

.component-code-chip {
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 11px;
  font-weight: 600;
  color: #64748B;
  background: #F8FAFC;
  border: 1px solid #E2E8F0;
  padding: 2px 7px;
  border-radius: 4px;
}

.layer-title {
  font-size: 16px;
  font-weight: 800;
  color: #0F172A;
  margin: 2px 0 0 0;
  letter-spacing: -0.01em;
}

/* Plain-Language Business Function Description */
.layer-business-role {
  display: flex;
  align-items: flex-start;
  gap: 7px;
  background: #F8FAFC;
  border-left: 3px solid #002060;
  border-radius: 0 4px 4px 0;
  padding: 6px 10px;
  margin-top: 3px;
}

.role-icon {
  font-size: 12px;
  line-height: 1.4;
}

.role-text {
  font-size: 12px;
  color: #334155;
  line-height: 1.45;
  font-weight: 500;
}

.card-header-actions {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
}

.btn-edit-layer {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  padding: 6px 12px;
  background: #F8FAFC;
  border: 1px solid #CBD5E1;
  border-radius: 6px;
  font-size: 11.5px;
  font-weight: 700;
  color: #002060;
  cursor: pointer;
  transition: all 0.15s ease;
}

.btn-edit-layer:hover {
  background: #002060;
  color: #FFFFFF;
  border-color: #002060;
}

.btn-delete-layer {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 6px 9px;
  background: #FFFFFF;
  border: 1px solid #CBD5E1;
  border-radius: 6px;
  color: #DC2626;
  cursor: pointer;
  transition: all 0.15s ease;
}

.btn-delete-layer:hover {
  background: #FEE2E2;
  border-color: #EF4444;
}

/* ─────────────────────────────────────────────
   Specs Comparison Row (Moderat vs Critical)
   ───────────────────────────────────────────── */
.preset-specs-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 14px;
}

.preset-specs-row.single-col {
  grid-template-columns: 1fr;
}

.spec-box {
  border-radius: 8px;
  padding: 14px 16px;
  border: 1px solid #E2E8F0;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.spec-box.box-mod {
  background: #F8FAFC;
  border-color: #BAE6FD;
}

.spec-box.box-crit {
  background: #FFFDF7;
  border-color: #FDE68A;
}

.spec-box-hdr {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.spec-tier-badge {
  display: flex;
  align-items: center;
  gap: 8px;
}

.tier-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
}

.mod-dot {
  background: #0284C7;
}

.crit-dot {
  background: #D97706;
}

.box-tier-title {
  font-size: 13px;
  font-weight: 800;
}

.box-tier-sub {
  font-size: 11.5px;
  color: #64748B;
  font-weight: 600;
}

.text-navy {
  color: #002060;
}

.text-amber {
  color: #B45309;
}

/* T-Shirt Size Badges */
.tshirt-size-badge {
  font-size: 11px;
  font-weight: 800;
  padding: 3px 8px;
  border-radius: 4px;
  letter-spacing: 0.02em;
}

.badge-size-small {
  background: #EFF6FF;
  color: #1D4ED8;
  border: 1px solid #BFDBFE;
}

.badge-size-medium {
  background: #F0FDFA;
  color: #0F766E;
  border: 1px solid #99F6E4;
}

.badge-size-standard {
  background: #F0FDF4;
  color: #15803D;
  border: 1px solid #BBF7D0;
}

.badge-size-plus {
  background: #FEF9C3;
  color: #A16207;
  border: 1px solid #FDE047;
}

.badge-size-large {
  background: #FFF7ED;
  color: #C2410C;
  border: 1px solid #FED7AA;
}

.badge-size-xlarge {
  background: #FEF2F2;
  color: #B91C1C;
  border: 1px solid #FECACA;
}

.badge-size-storage {
  background: #F5F3FF;
  color: #6D28D9;
  border: 1px solid #DDD6FE;
}

/* Delta Pills Strip */
.delta-pills-strip {
  display: flex;
  align-items: center;
  gap: 8px;
  background: #FEF3C7;
  border: 1px solid #FDE68A;
  border-radius: 6px;
  padding: 6px 10px;
  flex-wrap: wrap;
}

.delta-label {
  font-size: 10px;
  font-weight: 800;
  color: #92400E;
  text-transform: uppercase;
  letter-spacing: 0.03em;
}

.delta-chips-wrap {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-wrap: wrap;
}

.delta-chip {
  font-size: 11px;
  font-weight: 700;
  padding: 2px 7px;
  border-radius: 4px;
}

.delta-chip.scale {
  background: #FFFFFF;
  color: #B45309;
  border: 1px solid #FCD34D;
}

.delta-chip.upgrade {
  background: #FFFFFF;
  color: #047857;
  border: 1px solid #A7F3D0;
}

.delta-chip.license {
  background: #FFFFFF;
  color: #1D4ED8;
  border: 1px solid #BFDBFE;
}

.delta-chip.info {
  background: #FFFFFF;
  color: #475569;
  border: 1px solid #CBD5E1;
}

/* Shared Storage Special Callout */
.shared-storage-callout {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  border-radius: 6px;
  padding: 12px 14px;
  border: 1px dashed #CBD5E1;
}

.shared-storage-callout.mod-theme {
  background: #FFFFFF;
  border-color: #BAE6FD;
}

.shared-storage-callout.crit-theme {
  background: #FFFFFF;
  border-color: #FDE68A;
}

.shared-callout-icon {
  font-size: 24px;
  flex-shrink: 0;
}

.shared-callout-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
  flex: 1;
}

.shared-callout-title {
  font-size: 12.5px;
  font-weight: 800;
  color: #0F172A;
}

.shared-callout-desc {
  font-size: 11px;
  color: #64748B;
}

.shared-storage-val {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
}

.storage-num {
  font-size: 14px;
  font-weight: 800;
  color: #0F172A;
}

.storage-class-tag {
  font-size: 10px;
  font-weight: 700;
  color: #6D28D9;
  background: #F5F3FF;
  padding: 1px 6px;
  border-radius: 3px;
}

/* Spec Metrics Grid */
.spec-metrics-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 8px;
  background: #FFFFFF;
  border: 1px solid #E2E8F0;
  border-radius: 6px;
  padding: 10px 12px;
}

.metric-cell {
  display: flex;
  flex-direction: column;
}

.cell-lbl {
  font-size: 9.5px;
  font-weight: 700;
  color: #64748B;
  text-transform: uppercase;
  letter-spacing: 0.02em;
}

.cell-val {
  font-size: 12.5px;
  font-weight: 800;
  color: #0F172A;
  margin-top: 2px;
}

.cell-sub {
  font-size: 10.5px;
  color: #64748B;
  font-weight: 500;
  margin-top: 1px;
}

/* Tier Footer inside box */
.spec-box-tier-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  padding-top: 8px;
  border-top: 1px dashed rgba(0, 0, 0, 0.08);
  font-size: 11.5px;
}

.tier-lic-indicator {
  display: flex;
  align-items: center;
  gap: 6px;
  color: #334155;
  min-width: 0;
  flex: 1;
}

.lic-icon {
  font-size: 13px;
  flex-shrink: 0;
}

.lic-name {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.backup-badge {
  font-size: 10.5px;
  font-weight: 700;
  padding: 2px 7px;
  border-radius: 4px;
  flex-shrink: 0;
}

.backup-badge.active {
  background: #DCFCE7;
  color: #15803D;
  border: 1px solid #BBF7D0;
}

.backup-badge.inactive {
  background: #F1F5F9;
  color: #64748B;
  border: 1px solid #E2E8F0;
}

/* Card Bottom Footer */
.preset-card-ftr {
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-top: 1px dashed #E2E8F0;
  padding-top: 10px;
  gap: 12px;
  flex-wrap: wrap;
}

.ftr-left-refs {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 11.5px;
}

.sku-label {
  color: #64748B;
  font-weight: 600;
}

.sku-value {
  color: #1E293B;
  font-weight: 700;
}

.text-green {
  color: #15803D;
}

.ftr-sku-sep {
  color: #CBD5E1;
}

.ftr-contract-ref {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  font-size: 11px;
  font-weight: 600;
  color: #0369A1;
  background: #E0F2FE;
  padding: 3px 8px;
  border-radius: 4px;
}

.ref-icon {
  font-size: 11px;
}

/* ─────────────────────────────────────────────
   Matrix Table View (1 Screen Dense Comparison)
   ───────────────────────────────────────────── */
.preset-matrix-table-wrap {
  background: #FFFFFF;
  border: 1px solid #E2E8F0;
  border-radius: 10px;
  overflow-x: auto;
  box-shadow: 0 1px 3px rgba(15, 23, 42, 0.04);
  flex-shrink: 0;
  margin-bottom: 24px;
}

.preset-matrix-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 12px;
  text-align: left;
}

.preset-matrix-table thead tr {
  background: #F8FAFC;
  border-bottom: 1px solid #CBD5E1;
}

.preset-matrix-table thead th {
  padding: 8px 10px;
  font-weight: 700;
  color: #475569;
  border-right: 1px solid #E2E8F0;
}

.th-layer {
  width: 260px;
  min-width: 240px;
}

.th-cat {
  width: 130px;
}

.th-group-mod {
  background: #EFF6FF;
  color: #002060;
  text-align: center;
  font-size: 11px;
  letter-spacing: 0.03em;
  border-bottom: 1px solid #BFDBFE;
}

.th-group-crit {
  background: #FEF3C7;
  color: #92400E;
  text-align: center;
  font-size: 11px;
  letter-spacing: 0.03em;
  border-bottom: 1px solid #FDE68A;
}

.th-sub {
  font-size: 10px;
  font-weight: 700;
  color: #64748B;
  background: #F8FAFC;
  text-transform: uppercase;
}

.th-action {
  width: 70px;
  text-align: center;
}

.preset-matrix-table tbody tr.matrix-row {
  border-bottom: 1px solid #F1F5F9;
  transition: background 0.1s ease;
}

.preset-matrix-table tbody tr.matrix-row:hover {
  background: #F8FAFC;
}

.preset-matrix-table td {
  padding: 8px 10px;
  border-right: 1px solid #F1F5F9;
  vertical-align: middle;
}

.matrix-layer-cell {
  display: flex;
  align-items: flex-start;
  gap: 8px;
}

.matrix-emoji {
  font-size: 18px;
  flex-shrink: 0;
  margin-top: 1px;
}

.matrix-text-wrap {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.matrix-layer-name {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12.5px;
  color: #0F172A;
}

.matrix-no {
  font-size: 9.5px;
  font-weight: 800;
  color: #002060;
  background: #EFF6FF;
  padding: 1px 4px;
  border-radius: 3px;
}

.matrix-comp-sub {
  font-size: 11px;
  color: #64748B;
}

.matrix-role-sub {
  font-size: 10.5px;
  color: #475569;
  max-width: 220px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.matrix-cat-tag {
  display: inline-flex;
  align-items: center;
  gap: 3px;
  font-size: 10.5px;
  font-weight: 700;
  padding: 2px 6px;
  border-radius: 4px;
  white-space: nowrap;
}

.matrix-shared-tag {
  font-size: 10.5px;
  font-weight: 700;
  color: #6D28D9;
  background: #F5F3FF;
  padding: 2px 6px;
  border-radius: 4px;
  display: inline-block;
}

.matrix-compute-type {
  font-size: 10px;
  color: #64748B;
}

.text-muted {
  color: #94A3B8;
  font-size: 11px;
}

.matrix-lic-text {
  font-weight: 600;
  font-size: 11px;
  color: #334155;
  white-space: nowrap;
}

.matrix-nbu-dot {
  font-size: 9.5px;
  font-weight: 700;
  padding: 1px 5px;
  border-radius: 3px;
  display: inline-block;
  margin-top: 2px;
}

.matrix-nbu-dot.active {
  background: #DCFCE7;
  color: #15803D;
}

.matrix-nbu-dot.inactive {
  background: #F1F5F9;
  color: #94A3B8;
}

.matrix-delta-wrap {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.delta-chip-mini {
  font-size: 10px;
  font-weight: 700;
  padding: 1px 5px;
  border-radius: 3px;
  white-space: nowrap;
}

.delta-chip-mini.scale {
  background: #FEF3C7;
  color: #B45309;
}

.delta-chip-mini.upgrade {
  background: #DCFCE7;
  color: #15803D;
}

.delta-chip-mini.license {
  background: #EFF6FF;
  color: #1D4ED8;
}

.delta-chip-mini.info {
  background: #F1F5F9;
  color: #64748B;
}

.matrix-action-btns {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
}

.btn-matrix-action {
  width: 24px;
  height: 24px;
  border-radius: 4px;
  border: 1px solid #CBD5E1;
  background: #FFFFFF;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  color: #475569;
  transition: all 0.15s ease;
}

.btn-matrix-action.edit:hover {
  background: #EFF6FF;
  color: #002060;
  border-color: #002060;
}

.btn-matrix-action.del:hover {
  background: #FEE2E2;
  color: #DC2626;
  border-color: #EF4444;
}

/* ─────────────────────────────────────────────
   Empty State
   ───────────────────────────────────────────── */
.preset-empty-state {
  background: #FFFFFF;
  border: 1px dashed #CBD5E1;
  border-radius: 10px;
  padding: 40px 24px;
  text-align: center;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
}

.empty-icon-box {
  width: 52px;
  height: 52px;
  border-radius: 26px;
  background: #F8FAFC;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
}

.empty-title {
  margin: 0;
  font-size: 15px;
  font-weight: 800;
  color: #0F172A;
}

.empty-desc {
  margin: 0;
  font-size: 12.5px;
  color: #64748B;
  max-width: 400px;
}

.btn-reset-filter {
  margin-top: 6px;
  padding: 6px 14px;
  background: #002060;
  color: #FFFFFF;
  border: none;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 700;
  cursor: pointer;
}

/* ─────────────────────────────────────────────
   Modal Styling
   ───────────────────────────────────────────── */
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

.preset-edit-modal {
  max-width: 780px;
  border-radius: 12px;
  box-shadow: 0 20px 45px -10px rgba(0, 32, 96, 0.28);
}

.modal-hdr {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 22px;
  border-bottom: 1px solid #E2E8F0;
  background: #FFFFFF;
}

.modal-hdr-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.hdr-icon-circle {
  width: 38px;
  height: 38px;
  border-radius: 50%;
  background: #EFF6FF;
  border: 1px solid #BFDBFE;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.modal-hdr-left h3 {
  margin: 0;
  font-size: 16px;
  font-weight: 800;
  color: #002060;
  letter-spacing: -0.01em;
}

.modal-hdr-sub {
  margin: 2px 0 0 0;
  font-size: 12px;
  color: #64748B;
}

.modal-close {
  background: transparent;
  border: none;
  font-size: 18px;
  color: #64748B;
  cursor: pointer;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 6px;
  transition: all 0.15s;
}

.modal-close:hover {
  background: #F1F5F9;
  color: #0F172A;
}

/* Modal Wizard Tabs */
.modal-section-tabs {
  display: flex;
  background: #F8FAFC;
  padding: 6px 12px;
  gap: 8px;
  border-bottom: 1px solid #E2E8F0;
}

.tab-step-btn {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 9px 14px;
  border-radius: 8px;
  border: 1px solid transparent;
  background: transparent;
  font-size: 13px;
  font-weight: 600;
  color: #64748B;
  cursor: pointer;
  transition: all 0.15s ease;
}

.tab-step-btn:hover {
  background: #F1F5F9;
  color: #0F172A;
}

.tab-step-btn.active {
  background: #FFFFFF;
  color: #002060;
  font-weight: 700;
  border-color: #E2E8F0;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

.step-num {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 22px;
  height: 22px;
  border-radius: 50%;
  background: #E2E8F0;
  color: #475569;
  font-size: 11px;
  font-weight: 700;
  transition: all 0.15s;
}

.tab-step-btn.active .step-num {
  background: #002060;
  color: #FFFFFF;
}

.tab-step-btn.active .step-num.step-mod {
  background: #0284C7;
  color: #FFFFFF;
}

.tab-step-btn.active .step-num.step-crit {
  background: #D97706;
  color: #FFFFFF;
}

/* Modal Content Clean Form */
.preset-form-clean {
  display: flex;
  flex-direction: column;
  padding: 0;
  background: #FFFFFF;
}

.tab-pane-content {
  padding: 20px 24px;
  display: flex;
  flex-direction: column;
  gap: 16px;
  overflow-y: auto;
  max-height: 56vh;
}

.field-instruction-box {
  display: flex;
  gap: 10px;
  align-items: flex-start;
  background: #F0FDF4;
  border: 1px solid #BBF7D0;
  border-radius: 8px;
  padding: 10px 14px;
}

.inst-icon {
  font-size: 16px;
  line-height: 1;
}

.inst-text {
  font-size: 12px;
  color: #166534;
  line-height: 1.45;
  margin: 0;
}

.form-grid-2col {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 14px;
}

.clean-field {
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.clean-field.full-span {
  grid-column: 1 / -1;
}

.clean-label {
  font-size: 12px;
  font-weight: 700;
  color: #1E293B;
  display: flex;
  align-items: center;
  gap: 4px;
}

.clean-input {
  padding: 9px 12px;
  border: 1px solid #CBD5E1;
  border-radius: 6px;
  font-size: 13px;
  color: #0F172A;
  background: #FFFFFF;
  transition: all 0.15s;
  outline: none;
}

.clean-input:focus {
  border-color: #002060;
  box-shadow: 0 0 0 3px rgba(0, 32, 96, 0.08);
}

.clean-textarea {
  padding: 9px 12px;
  border: 1px solid #CBD5E1;
  border-radius: 6px;
  font-size: 12px;
  color: #0F172A;
  background: #FFFFFF;
  resize: vertical;
  outline: none;
  font-family: inherit;
  line-height: 1.4;
  transition: all 0.15s;
}

.clean-textarea:focus {
  border-color: #002060;
  box-shadow: 0 0 0 3px rgba(0, 32, 96, 0.08);
}

.clean-hint {
  font-size: 11px;
  color: #64748B;
}

.req {
  color: #DC2626;
}

.font-mono {
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
}

/* Tier Banners */
.tier-banner-info {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 14px;
  border-radius: 8px;
  border: 1px solid #E2E8F0;
}

.banner-mod {
  background: linear-gradient(135deg, #F0F9FF 0%, #E0F2FE 100%);
  border-color: #BAE6FD;
}

.banner-crit {
  background: linear-gradient(135deg, #FFFBEB 0%, #FEF3C7 100%);
  border-color: #FDE68A;
}

.tier-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  flex-shrink: 0;
}

.mod-dot {
  background: #0284C7;
  box-shadow: 0 0 0 3px rgba(2, 132, 199, 0.2);
}

.crit-dot {
  background: #D97706;
  box-shadow: 0 0 0 3px rgba(217, 119, 6, 0.2);
}

.banner-texts {
  display: flex;
  flex-direction: column;
  gap: 2px;
  flex: 1;
}

.banner-title {
  font-size: 13px;
  font-weight: 800;
}

.banner-sub {
  font-size: 11px;
  color: #475569;
}

.text-navy {
  color: #002060;
}

.text-amber {
  color: #92400E;
}

.title-with-copy-btn {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.btn-copy-from-mod {
  padding: 4px 10px;
  background: #FFFFFF;
  border: 1px solid #FCD34D;
  border-radius: 6px;
  font-size: 11px;
  font-weight: 700;
  color: #B45309;
  cursor: pointer;
  transition: all 0.15s;
}

.btn-copy-from-mod:hover {
  background: #FEF3C7;
  border-color: #F59E0B;
}

.label-with-tip {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 6px;
  margin-bottom: 2px;
}

/* Stepper Row */
.steppers-card-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 14px;
}

.stepper-item-box {
  background: #F8FAFC;
  border: 1px solid #E2E8F0;
  border-radius: 8px;
  padding: 12px 14px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.box-lbl {
  font-size: 12px;
  font-weight: 700;
  color: #334155;
}

/* Storage & Backup Row */
.storage-class-toggle-row {
  display: grid;
  grid-template-columns: 1.4fr 1fr;
  gap: 14px;
  align-items: stretch;
}

.class-picker-half,
.backup-chk-half {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.clean-backup-chk-box {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  padding: 12px;
  background: #F8FAFC;
  border: 1.5px solid #E2E8F0;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.15s ease;
  flex: 1;
}

.clean-backup-chk-box:hover {
  background: #F1F5F9;
  border-color: #CBD5E1;
}

.clean-backup-chk-box.active {
  background: #EFF6FF;
  border-color: #3B82F6;
}

.custom-chk-big {
  width: 18px;
  height: 18px;
  margin-top: 2px;
  cursor: pointer;
  accent-color: #002060;
}

.chk-desc-wrap {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.chk-title {
  font-size: 12px;
  font-weight: 700;
  color: #0F172A;
}

.chk-sub {
  font-size: 11px;
  color: #64748B;
  line-height: 1.35;
}

/* Modal Bottom Actions */
.modal-ftr-clean {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 24px;
  border-top: 1px solid #E2E8F0;
  background: #F8FAFC;
}

.modal-ftr-right {
  display: flex;
  align-items: center;
  gap: 10px;
}

.clean-btn-cancel {
  padding: 8px 16px;
  background: transparent;
  border: 1px solid #CBD5E1;
  border-radius: 6px;
  font-size: 13px;
  font-weight: 600;
  color: #475569;
  cursor: pointer;
  transition: all 0.15s;
}

.clean-btn-cancel:hover {
  background: #E2E8F0;
  color: #0F172A;
}

.clean-btn-nav {
  padding: 8px 16px;
  background: #FFFFFF;
  border: 1px solid #CBD5E1;
  border-radius: 6px;
  font-size: 13px;
  font-weight: 600;
  color: #334155;
  cursor: pointer;
  transition: all 0.15s;
}

.clean-btn-nav:hover {
  background: #F1F5F9;
  border-color: #94A3B8;
}

.clean-btn-nav.primary-step {
  background: #002060;
  border-color: #002060;
  color: #FFFFFF;
  font-weight: 700;
}

.clean-btn-nav.primary-step:hover {
  background: #003399;
}

.clean-btn-submit {
  padding: 8px 20px;
  background: #002060;
  border: none;
  border-radius: 6px;
  font-size: 13px;
  font-weight: 700;
  color: #FFFFFF;
  cursor: pointer;
  transition: all 0.15s;
  box-shadow: 0 2px 4px rgba(0, 32, 96, 0.2);
}

.clean-btn-submit:hover {
  background: #003399;
  box-shadow: 0 4px 8px rgba(0, 32, 96, 0.25);
}

/* ─────────────────────────────────────────────────────────────
   VIEW MODE 3: HIGH-LEVEL ARCHITECTURE BLUEPRINT DIAGRAM
   ───────────────────────────────────────────────────────────── */
.preset-diagram-canvas {
  background-color: #F8FAFC;
  background-image: radial-gradient(#CBD5E1 1px, transparent 1px);
  background-size: 20px 20px;
  border: 1px solid #CBD5E1;
  border-radius: 12px;
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 16px;
  box-shadow: inset 0 2px 6px rgba(15, 23, 42, 0.03);
  position: relative;
}

/* Blueprint Header & Legend */
.blueprint-header-bar {
  background: linear-gradient(135deg, #002060 0%, #0F296B 60%, #1E3A8A 100%);
  color: #FFFFFF;
  border-radius: 10px;
  padding: 16px 20px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 20px;
  flex-wrap: wrap;
  box-shadow: 0 4px 12px rgba(0, 32, 96, 0.15);
}

.blueprint-title-block {
  display: flex;
  flex-direction: column;
  gap: 6px;
  flex: 1;
  min-width: 300px;
}

.bp-title-row {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.bp-badge-tag {
  font-size: 10px;
  font-weight: 800;
  letter-spacing: 0.08em;
  padding: 3px 8px;
  border-radius: 4px;
  background: rgba(255, 255, 255, 0.16);
  border: 1px solid rgba(255, 255, 255, 0.28);
  color: #FDE68A;
}

.bp-main-title {
  margin: 0;
  font-size: 16px;
  font-weight: 800;
  color: #FFFFFF;
  letter-spacing: -0.01em;
}

.bp-lens-pill {
  font-size: 11px;
  font-weight: 700;
  padding: 3px 10px;
  border-radius: 20px;
  background: rgba(255, 255, 255, 0.14);
  color: #E2E8F0;
  border: 1px solid rgba(255, 255, 255, 0.22);
}

.bp-lens-pill.moderate {
  background: rgba(56, 189, 248, 0.22);
  border-color: #7DD3FC;
  color: #E0F2FE;
}

.bp-lens-pill.critical {
  background: rgba(245, 158, 11, 0.28);
  border-color: #FCD34D;
  color: #FEF3C7;
}

.bp-subtitle {
  margin: 0;
  font-size: 12px;
  color: #CBD5E1;
  line-height: 1.45;
}

.blueprint-legend-box {
  display: flex;
  align-items: center;
  gap: 12px;
  background: rgba(15, 23, 42, 0.35);
  border: 1px solid rgba(255, 255, 255, 0.15);
  border-radius: 8px;
  padding: 8px 14px;
  flex-wrap: wrap;
}

.legend-item {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 11px;
  font-weight: 600;
  color: #F1F5F9;
  white-space: nowrap;
}

.leg-dot {
  width: 9px;
  height: 9px;
  border-radius: 50%;
  display: inline-block;
}

.leg-dot.vm {
  background: #38BDF8;
  box-shadow: 0 0 0 2px rgba(56, 189, 248, 0.3);
}

.leg-dot.bm {
  background: #F59E0B;
  box-shadow: 0 0 0 2px rgba(245, 158, 11, 0.35);
}

.leg-dot.san {
  background: #A855F7;
  box-shadow: 0 0 0 2px rgba(168, 85, 247, 0.35);
}

.legend-sep {
  width: 1px;
  height: 16px;
  background: rgba(255, 255, 255, 0.2);
}

.leg-line {
  width: 18px;
  height: 0;
  display: inline-block;
}

.leg-line.sync {
  border-top: 2px solid #38BDF8;
}

.leg-line.async {
  border-top: 2px dashed #FCD34D;
}

/* External Entry Pill */
.bp-external-entry-wrap {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-top: 2px;
}

.bp-external-pill {
  display: inline-flex;
  align-items: center;
  gap: 14px;
  background: #FFFFFF;
  border: 2px solid #002060;
  border-radius: 40px;
  padding: 10px 22px;
  box-shadow: 0 4px 12px rgba(0, 32, 96, 0.08);
  flex-wrap: wrap;
}

.ext-icon-pulse {
  width: 34px;
  height: 34px;
  border-radius: 50%;
  background: #EFF6FF;
  border: 1px solid #BFDBFE;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
}

.ext-text-group {
  display: flex;
  flex-direction: column;
}

.ext-eyebrow {
  font-size: 9.5px;
  font-weight: 800;
  color: #0284C7;
  letter-spacing: 0.06em;
}

.ext-title {
  font-size: 13px;
  font-weight: 800;
  color: #0F172A;
}

.ext-protocol-tags {
  display: flex;
  align-items: center;
  gap: 6px;
}

.proto-chip {
  font-size: 10.5px;
  font-weight: 700;
  color: #002060;
  background: #F1F5F9;
  border: 1px solid #CBD5E1;
  padding: 3px 9px;
  border-radius: 12px;
}

/* North-South Spine Connector */
.bp-spine-connector {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  position: relative;
  padding: 6px 0;
}

.bp-spine-connector .spine-line {
  width: 2px;
  height: 16px;
  background: #0284C7;
}

.bp-spine-connector .spine-line.async {
  background: transparent;
  border-left: 2px dashed #D97706;
}

.bp-spine-connector .spine-pill {
  font-size: 11px;
  font-weight: 700;
  color: #002060;
  background: #EFF6FF;
  border: 1px solid #93C5FD;
  padding: 3px 12px;
  border-radius: 14px;
  box-shadow: 0 2px 4px rgba(15, 23, 42, 0.05);
  z-index: 1;
}

/* Zones Stack & Swimlane Cards */
.bp-zones-stack {
  display: flex;
  flex-direction: column;
}

.bp-zone-lane {
  background: rgba(255, 255, 255, 0.92);
  backdrop-filter: blur(6px);
  border: 1px solid #CBD5E1;
  border-left-width: 5px;
  border-radius: 12px;
  padding: 16px 18px;
  display: flex;
  flex-direction: column;
  gap: 14px;
  box-shadow: 0 2px 8px rgba(15, 23, 42, 0.04);
  transition: all 0.2s ease;
}

.bp-zone-lane.zone-theme-edge {
  border-left-color: #0284C7;
}

.bp-zone-lane.zone-theme-compute {
  border-left-color: #002060;
  background: linear-gradient(180deg, #F8FAFF 0%, #FFFFFF 100%);
}

.bp-zone-lane.zone-theme-middleware {
  border-left-color: #D97706;
}

.bp-zone-lane.zone-theme-data {
  border-left-color: #7C3AED;
}

.bp-zone-lane.zone-theme-ops {
  border-left-color: #059669;
}

.bp-zone-lane.zone-theme-other {
  border-left-color: #475569;
}

.bp-zone-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  padding-bottom: 12px;
  border-bottom: 1px dashed #E2E8F0;
  flex-wrap: wrap;
}

.zone-hdr-left {
  display: flex;
  align-items: center;
  gap: 12px;
  flex: 1;
  min-width: 260px;
}

.zone-code-badge {
  font-size: 10.5px;
  font-weight: 800;
  letter-spacing: 0.06em;
  padding: 4px 8px;
  border-radius: 5px;
  background: #0F172A;
  color: #FFFFFF;
  flex-shrink: 0;
}

.zone-theme-edge .zone-code-badge { background: #0284C7; }
.zone-theme-compute .zone-code-badge { background: #002060; }
.zone-theme-middleware .zone-code-badge { background: #D97706; }
.zone-theme-data .zone-code-badge { background: #6D28D9; }
.zone-theme-ops .zone-code-badge { background: #047857; }

.zone-icon-box {
  width: 36px;
  height: 36px;
  border-radius: 8px;
  background: #F1F5F9;
  border: 1px solid #E2E8F0;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  flex-shrink: 0;
}

.zone-title-wrap {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.zone-title-line {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.zone-name {
  margin: 0;
  font-size: 14.5px;
  font-weight: 800;
  color: #0F172A;
}

.zone-role-pill {
  font-size: 10.5px;
  font-weight: 700;
  color: #475569;
  background: #F1F5F9;
  border: 1px solid #CBD5E1;
  padding: 2px 8px;
  border-radius: 12px;
}

.zone-desc {
  margin: 0;
  font-size: 11.5px;
  color: #64748B;
}

.zone-hdr-metrics {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.zone-stat-chip {
  display: flex;
  flex-direction: column;
  background: #F8FAFC;
  border: 1px solid #E2E8F0;
  border-radius: 6px;
  padding: 4px 10px;
}

.zone-stat-chip.mod {
  background: #F0F9FF;
  border-color: #BAE6FD;
}

.zone-stat-chip.crit {
  background: #FFFBEB;
  border-color: #FDE68A;
}

.z-stat-lbl {
  font-size: 9.5px;
  font-weight: 700;
  color: #64748B;
  text-transform: uppercase;
}

.z-stat-val {
  font-size: 11.5px;
  font-weight: 800;
  color: #0F172A;
}

.zone-stat-chip.mod .z-stat-val { color: #0369A1; }
.zone-stat-chip.crit .z-stat-val { color: #B45309; }

/* Zone Nodes Flow Layout */
.bp-zone-nodes-flow {
  display: flex;
  align-items: stretch;
  gap: 12px;
  flex-wrap: wrap;
}

.bp-zone-nodes-flow.single-node-hub {
  justify-content: center;
}

.bp-zone-nodes-flow.single-node-hub .bp-node-card {
  max-width: 640px;
  flex: 1 1 480px;
}

/* Blueprint Node Card */
.bp-node-card {
  flex: 1 1 280px;
  min-width: 265px;
  background: #FFFFFF;
  border: 1.5px solid #CBD5E1;
  border-radius: 10px;
  padding: 13px 14px;
  display: flex;
  flex-direction: column;
  gap: 10px;
  cursor: pointer;
  transition: all 0.18s ease;
  box-shadow: 0 1px 3px rgba(15, 23, 42, 0.04);
  position: relative;
}

.bp-node-card:hover {
  border-color: #002060;
  transform: translateY(-2px);
  box-shadow: 0 8px 18px -4px rgba(0, 32, 96, 0.12);
}

.bp-node-card.is-selected {
  border-color: #002060;
  box-shadow: 0 0 0 3px rgba(0, 32, 96, 0.16), 0 8px 20px -4px rgba(0, 32, 96, 0.15);
}

.bp-node-card.has-baremetal {
  border-top: 3px solid #D97706;
}

.bp-node-card.is-shared-san {
  border-top: 3px solid #7C3AED;
}

.bp-node-topbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.bp-node-badges {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-wrap: wrap;
}

.bp-layer-no {
  font-size: 10px;
  font-weight: 800;
  color: #002060;
  background: #EFF6FF;
  border: 1px solid #BFDBFE;
  padding: 2px 6px;
  border-radius: 4px;
}

.bp-role-tag {
  font-size: 10px;
  font-weight: 700;
  color: #475569;
  background: #F1F5F9;
  border: 1px solid #E2E8F0;
  padding: 2px 7px;
  border-radius: 4px;
}

.bp-btn-edit {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 3px 8px;
  font-size: 10.5px;
  font-weight: 700;
  color: #002060;
  background: #F8FAFC;
  border: 1px solid #CBD5E1;
  border-radius: 5px;
  cursor: pointer;
  transition: all 0.15s ease;
}

.bp-btn-edit:hover {
  background: #002060;
  color: #FFFFFF;
  border-color: #002060;
}

.bp-node-identity {
  display: flex;
  align-items: center;
  gap: 10px;
}

.bp-node-avatar {
  width: 38px;
  height: 38px;
  border-radius: 8px;
  background: #F8FAFC;
  border: 1px solid #E2E8F0;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  flex-shrink: 0;
}

.bp-node-titles {
  display: flex;
  flex-direction: column;
  min-width: 0;
}

.bp-node-layer-name {
  margin: 0;
  font-size: 14px;
  font-weight: 800;
  color: #0F172A;
  line-height: 1.25;
}

.bp-node-comp-name {
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 11px;
  font-weight: 600;
  color: #0284C7;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.bp-node-desc {
  margin: 0;
  font-size: 11.5px;
  color: #64748B;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

/* Dual Tiers Comparison inside Node */
.bp-node-dual-tiers {
  display: flex;
  flex-direction: column;
  gap: 6px;
  margin-top: auto;
}

.bp-tier-strip {
  border-radius: 7px;
  padding: 7px 9px;
  display: flex;
  flex-direction: column;
  gap: 3px;
  border: 1px solid #E2E8F0;
}

.bp-tier-strip.mod-strip {
  background: #F0F9FF;
  border-color: #BAE6FD;
}

.bp-tier-strip.crit-strip {
  background: #FFFBEB;
  border-color: #FDE68A;
}

.strip-hdr {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 6px;
}

.strip-tier-dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  flex-shrink: 0;
}

.strip-tier-dot.mod { background: #0284C7; }
.strip-tier-dot.crit { background: #D97706; }

.strip-tier-name {
  font-size: 10px;
  font-weight: 800;
  color: #0F172A;
  margin-right: auto;
}

.strip-compute-pill {
  font-size: 10px;
  font-weight: 800;
  padding: 1px 6px;
  border-radius: 4px;
  background: #FFFFFF;
  border: 1px solid #CBD5E1;
  color: #0F172A;
}

.strip-compute-pill.baremetal {
  background: #FEF3C7;
  border-color: #F59E0B;
  color: #92400E;
}

.strip-compute-pill.shared {
  background: #F3E8FF;
  border-color: #C084FC;
  color: #6B21A8;
}

.strip-specs-line {
  font-size: 11px;
  color: #334155;
}

.strip-specs-line strong {
  color: #0F172A;
  font-weight: 800;
}

.strip-lic-line {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.strip-lic-badge {
  font-size: 10.5px;
  font-weight: 600;
  color: #475569;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.strip-lic-badge.crit-lic {
  color: #92400E;
  font-weight: 700;
}

.bp-node-delta-bar {
  display: flex;
  align-items: center;
  gap: 5px;
  flex-wrap: wrap;
  padding-top: 2px;
}

.bp-delta-tag {
  font-size: 10px;
  font-weight: 700;
  padding: 2px 7px;
  border-radius: 4px;
  background: #F1F5F9;
  color: #475569;
  border: 1px solid #CBD5E1;
}

.bp-delta-tag.upgrade {
  background: #FEF3C7;
  color: #B45309;
  border-color: #FDE68A;
}

.bp-delta-tag.scale {
  background: #EFF6FF;
  color: #1D4ED8;
  border-color: #BFDBFE;
}

.bp-delta-tag.license {
  background: #ECFDF5;
  color: #047857;
  border-color: #A7F3D0;
}

/* Single Tier Lens inside Node */
.bp-node-single-tier {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 9px 10px;
  border-radius: 8px;
  border: 1px solid #E2E8F0;
  margin-top: auto;
}

.bp-node-single-tier.moderate {
  background: #F0F9FF;
  border-color: #BAE6FD;
}

.bp-node-single-tier.critical {
  background: #FFFBEB;
  border-color: #FDE68A;
}

.single-spec-badges {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 6px;
  flex-wrap: wrap;
}

.single-compute-badge {
  font-size: 11px;
  font-weight: 800;
  color: #0F172A;
  background: #FFFFFF;
  border: 1px solid #CBD5E1;
  padding: 2px 8px;
  border-radius: 5px;
}

.single-compute-badge.baremetal {
  background: #FEF3C7;
  border-color: #F59E0B;
  color: #92400E;
}

.single-compute-badge.shared {
  background: #F3E8FF;
  border-color: #C084FC;
  color: #6B21A8;
}

.single-tshirt-pill {
  font-size: 10px;
  font-weight: 700;
  padding: 2px 7px;
  border-radius: 4px;
}

.single-metrics-box {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 6px;
  background: #FFFFFF;
  border: 1px solid rgba(15, 23, 42, 0.07);
  border-radius: 6px;
  padding: 6px 8px;
}

.s-metric {
  display: flex;
  flex-direction: column;
}

.s-lbl {
  font-size: 9px;
  font-weight: 700;
  color: #64748B;
  text-transform: uppercase;
}

.s-val {
  font-size: 11.5px;
  font-weight: 800;
  color: #0F172A;
}

.single-lic-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.s-lic-text {
  font-size: 11px;
  font-weight: 700;
  color: #334155;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.s-nbu-pill {
  font-size: 10px;
  font-weight: 700;
  padding: 1px 6px;
  border-radius: 4px;
  background: #F1F5F9;
  color: #64748B;
  flex-shrink: 0;
}

.s-nbu-pill.on {
  background: #DCFCE7;
  color: #15803D;
}

/* Intra-Zone Horizontal Connector */
.bp-intra-connector {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 4px;
  padding: 0 2px;
  min-width: 84px;
  flex-shrink: 0;
}

.intra-label {
  font-size: 9.5px;
  font-weight: 700;
  color: #475569;
  background: #F1F5F9;
  border: 1px solid #CBD5E1;
  padding: 2px 7px;
  border-radius: 10px;
  white-space: nowrap;
}

.intra-arrow-line {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
}

.intra-arrow-line::before {
  content: '';
  flex: 1;
  height: 0;
  border-top: 2px solid #0284C7;
}

.intra-arrow-line.async::before {
  border-top: 2px dashed #D97706;
}

.intra-arrow-line.san::before {
  border-top: 3px double #7C3AED;
}

.arrow-head {
  font-size: 12px;
  font-weight: 800;
  color: #0284C7;
  margin-left: 2px;
}

.intra-arrow-line.async .arrow-head { color: #D97706; }
.intra-arrow-line.san .arrow-head { color: #7C3AED; }

/* Selected Node Inspector Panel */
.bp-inspector-panel {
  position: sticky;
  bottom: 12px;
  background: #FFFFFF;
  border: 2px solid #002060;
  border-radius: 12px;
  padding: 14px 18px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  box-shadow: 0 12px 28px -6px rgba(0, 32, 96, 0.25);
  z-index: 10;
  flex-wrap: wrap;
}

.insp-left {
  display: flex;
  align-items: center;
  gap: 14px;
  flex: 1;
  min-width: 280px;
}

.insp-avatar {
  width: 44px;
  height: 44px;
  border-radius: 10px;
  background: #EFF6FF;
  border: 1px solid #BFDBFE;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 22px;
  flex-shrink: 0;
}

.insp-info {
  display: flex;
  flex-direction: column;
  gap: 3px;
}

.insp-title-line {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.insp-no {
  font-size: 10.5px;
  font-weight: 800;
  color: #FFFFFF;
  background: #002060;
  padding: 2px 7px;
  border-radius: 4px;
}

.insp-layer {
  margin: 0;
  font-size: 15px;
  font-weight: 800;
  color: #0F172A;
}

.insp-comp {
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 11.5px;
  font-weight: 600;
  color: #0284C7;
  background: #F0F9FF;
  border: 1px solid #BAE6FD;
  padding: 1px 7px;
  border-radius: 4px;
}

.insp-ref {
  font-size: 11px;
  font-weight: 600;
  color: #475569;
  background: #F8FAFC;
  border: 1px solid #E2E8F0;
  padding: 1px 7px;
  border-radius: 4px;
}

.insp-notes {
  margin: 0;
  font-size: 12px;
  color: #334155;
}

.insp-lic-pair {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 11.5px;
  color: #475569;
  flex-wrap: wrap;
}

.insp-sep {
  color: #CBD5E1;
}

.insp-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.btn-insp-edit {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 8px 14px;
  background: #002060;
  color: #FFFFFF;
  border: none;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.15s ease;
}

.btn-insp-edit:hover {
  background: #003399;
}

.btn-insp-close {
  width: 32px;
  height: 32px;
  border-radius: 6px;
  border: 1px solid #CBD5E1;
  background: #F8FAFC;
  color: #475569;
  cursor: pointer;
  font-size: 13px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.btn-insp-close:hover {
  background: #E2E8F0;
  color: #0F172A;
}

@media (max-width: 960px) {
  .bp-intra-connector {
    display: none;
  }
}
</style>
