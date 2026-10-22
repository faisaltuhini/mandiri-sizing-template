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
          type="button"
          class="btn-reset-defaults"
          @click="$emit('reset-to-official-defaults')"
          title="Reset semua preset ke standar baseline resmi Bank Mandiri"
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
      <!-- ALERT BANNER: INFO KATALOG BARU TERDETEKSI -->
      <transition name="banner-slide">
        <div
          v-if="(unregisteredCatalogItems?.length || 0) > 0 && !isBannerDismissed"
          class="new-catalog-alert-banner"
        >
          <div class="banner-content-left">
            <div class="alert-icon-circle">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="#D97706" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round">
                <circle cx="12" cy="12" r="10"></circle>
                <line x1="12" y1="8" x2="12" y2="12"></line>
                <line x1="12" y1="16" x2="12.01" y2="16"></line>
              </svg>
            </div>
            <div class="alert-texts">
              <div class="alert-title-row">
                <span class="alert-badge">Info Katalog Baru</span>
                <h4 class="alert-title">
                  Ada {{ unregisteredCatalogItems?.length || 0 }} komponen katalog baru yang belum dimasukkan ke preset arsitektur
                </h4>
              </div>
              <p class="alert-sub">
                Komponen di bawah ini telah terdaftar di master data katalog. Apakah Anda ingin menginputkannya sebagai layer baru?
              </p>
              
              <!-- Chips of unregistered catalog items -->
              <div class="new-items-chips">
                <button
                  type="button"
                  v-for="(catItem, idx) in (unregisteredCatalogItems || []).slice(0, 4)"
                  :key="idx"
                  class="new-item-chip"
                  @click="openAddLayerModalWithCatalogItem(catItem)"
                  :title="'Klik untuk input layer baru: ' + catItem.name"
                >
                  <span class="chip-icon">📦</span>
                  <span class="chip-name">{{ catItem.name }}</span>
                  <span class="chip-cat">{{ catItem.category }}</span>
                  <span class="chip-action">+ Input</span>
                </button>
                <span v-if="(unregisteredCatalogItems?.length || 0) > 4" class="more-chips-pill">
                  +{{ (unregisteredCatalogItems?.length || 0) - 4 }} katalog lainnya
                </span>
              </div>
            </div>
          </div>

          <div class="banner-actions-right">
            <button
              type="button"
              class="btn-input-new-layer"
              @click="unregisteredCatalogItems?.[0] && openAddLayerModalWithCatalogItem(unregisteredCatalogItems[0])"
            >
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
                <line x1="12" y1="5" x2="12" y2="19"/>
                <line x1="5" y1="12" x2="19" y2="12"/>
              </svg>
              <span>Inputkan Layer Baru</span>
            </button>
            <button
              type="button"
              class="btn-dismiss-alert"
              @click="dismissNewCatalogBanner"
              title="Abaikan peringatan ini untuk sementara"
            >
              Nanti Saja
            </button>
          </div>
        </div>
      </transition>

      <!-- Toolbar Filter & Search -->
      <div class="preset-toolbar">
        <div class="search-input-wrap">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="#64748B" stroke-width="2">
            <circle cx="11" cy="11" r="8"/>
            <line x1="21" y1="21" x2="16.65" y2="16.65"/>
          </svg>
          <input
            type="text"
            v-model="searchQuery"
            class="preset-search-input"
            placeholder="Cari berdasarkan nama layer, komponen, atau lisensi..."
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
          <span class="seg-label">Tampilan Preset:</span>
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

        <div class="preset-stats">
          <span class="stat-badge">
            <strong>{{ filteredPresetItems.length }}</strong> Layer Terdaftar
          </span>
        </div>
      </div>

      <!-- CARDS LIST OF PRESET LAYERS -->
      <div class="preset-cards-grid">
        <div
          v-for="item in filteredPresetItems"
          :key="item.id"
          class="preset-card"
        >
          <!-- Card Header -->
          <div class="preset-card-header">
            <div class="card-header-left">
              <span class="layer-number-badge">Layer {{ item.no }}</span>
              <div class="layer-title-wrap">
                <h3 class="layer-title">{{ item.layer }}</h3>
                <span class="component-subtitle">{{ item.component }}</span>
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
                v-if="item.no > 13"
                type="button"
                class="btn-delete-layer"
                @click="confirmDeleteCustomLayer(item)"
                title="Hapus Layer Kustom Ini"
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
                <div class="tier-dot mod-dot"></div>
                <span class="box-tier-title text-navy">Formula Baseline: Moderat (N+1)</span>
              </div>

              <div class="spec-metrics-grid">
                <div class="metric-cell">
                  <span class="cell-lbl">Nodes:</span>
                  <span class="cell-val">{{ item.moderate.nodeCount }} Node</span>
                </div>
                <div class="metric-cell">
                  <span class="cell-lbl">Compute:</span>
                  <span class="cell-val">{{ item.moderate.vcpuPerNode }} vCPU / node</span>
                </div>
                <div class="metric-cell">
                  <span class="cell-lbl">RAM:</span>
                  <span class="cell-val">{{ item.moderate.ramGBPerNode }} GB / node</span>
                </div>
                <div class="metric-cell">
                  <span class="cell-lbl">Storage:</span>
                  <span class="cell-val">{{ item.moderate.storageGBPerNode }} GB ({{ item.moderate.storageClass }})</span>
                </div>
              </div>

              <div class="spec-type-strip">
                <span class="type-pill font-mono">{{ formatComputeType(item.moderate.computeType) }}</span>
                <span :class="['backup-badge', item.moderate.backupEnabled ? 'active' : 'inactive']">
                  {{ item.moderate.backupEnabled ? '✓ NetBackup NBU' : 'Tanpa NBU' }}
                </span>
                <span class="license-spec-pill">
                  Lisensi: <strong>{{ item.moderate.licenseKind === 'enterprise' ? 'Enterprise' : 'Free OSS' }}</strong>
                </span>
              </div>
            </div>

            <!-- CRITICAL SPEC BOX -->
            <div
              v-if="activeTierTab === 'both' || activeTierTab === 'critical'"
              class="spec-box box-crit"
            >
              <div class="spec-box-hdr">
                <div class="tier-dot crit-dot"></div>
                <span class="box-tier-title text-amber">Formula Baseline: Critical (2N)</span>
              </div>

              <div class="spec-metrics-grid">
                <div class="metric-cell">
                  <span class="cell-lbl">Nodes:</span>
                  <span class="cell-val">{{ item.critical.nodeCount }} Node</span>
                </div>
                <div class="metric-cell">
                  <span class="cell-lbl">Compute:</span>
                  <span class="cell-val">{{ item.critical.vcpuPerNode }} vCPU / node</span>
                </div>
                <div class="metric-cell">
                  <span class="cell-lbl">RAM:</span>
                  <span class="cell-val">{{ item.critical.ramGBPerNode }} GB / node</span>
                </div>
                <div class="metric-cell">
                  <span class="cell-lbl">Storage:</span>
                  <span class="cell-val">{{ item.critical.storageGBPerNode }} GB ({{ item.critical.storageClass }})</span>
                </div>
              </div>

              <div class="spec-type-strip">
                <span class="type-pill font-mono">{{ formatComputeType(item.critical.computeType) }}</span>
                <span :class="['backup-badge', item.critical.backupEnabled ? 'active' : 'inactive']">
                  {{ item.critical.backupEnabled ? '✓ NetBackup NBU' : 'Tanpa NBU' }}
                </span>
                <span class="license-spec-pill">
                  Lisensi: <strong>{{ item.critical.licenseKind === 'enterprise' ? 'Enterprise' : 'Free OSS' }}</strong>
                </span>
              </div>
            </div>
          </div>

          <!-- License & Meta Footer -->
          <div class="preset-card-ftr">
            <div class="ftr-lic-pair">
              <span class="lic-tag free" title="Community Free OSS SKU">
                OSS: <strong>{{ item.licenseFree }}</strong>
              </span>
              <span class="lic-tag enterprise" title="Enterprise Official License SKU">
                Enterprise: <strong>{{ item.licenseEnterprise }}</strong>
              </span>
            </div>
            <div v-if="item.contractRef" class="ftr-contract-ref">
              Ref: {{ item.contractRef }}
            </div>
          </div>
        </div>
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

        <form @submit.prevent="savePresetLayer" class="modal-content preset-form-clean">
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
                  />
                </div>
                <BaseStepper
                  v-model="form.moderate.storageGBPerNode"
                  :min="10"
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
                :columns="3"
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
                  },
                  {
                    value: 'shared',
                    title: 'Penyimpanan Berbagi',
                    technicalLabel: 'Shared / LUN',
                    description: 'Hanya ruang disk tanpa compute tambahan.',
                    icon: '🔗'
                  }
                ]"
              />
            </div>

            <!-- Storage Class & NetBackup -->
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
                    <span class="chk-title">Sertakan NetBackup NBU</span>
                    <span class="chk-sub">Pencadangan harian otomatis enterprise Bank Mandiri.</span>
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
                  />
                </div>
                <BaseStepper
                  v-model="form.critical.storageGBPerNode"
                  :min="10"
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
                :columns="3"
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
                  },
                  {
                    value: 'shared',
                    title: 'Penyimpanan Berbagi',
                    technicalLabel: 'Shared / LUN',
                    description: 'Hanya ruang disk tanpa compute tambahan.',
                    icon: '🔗'
                  }
                ]"
              />
            </div>

            <!-- Storage Class & NetBackup -->
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
                    <span class="chk-title">Sertakan NetBackup NBU</span>
                    <span class="chk-sub">Pencadangan harian otomatis enterprise Bank Mandiri.</span>
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

function dismissNewCatalogBanner() {
  isBannerDismissed.value = true
}

const unregisteredCatalogItems = computed<UnregisteredCatalogItem[]>(() => {
  const items = props.presetItems || []
  const existingComponents = new Set(
    items.map(p => (p?.component || '').trim().toLowerCase())
  )
  const existingLayers = new Set(
    items.map(p => (p?.layer || '').trim().toLowerCase())
  )

  const unreg: UnregisteredCatalogItem[] = []
  const seenNames = new Set<string>()

  // 1. Check backend catalog components
  for (const bc of backendComponents.value) {
    const compName = (bc.component_name || '').trim()
    const compLower = compName.toLowerCase()
    if (!compName || seenNames.has(compLower)) continue

    if (!existingComponents.has(compLower)) {
      seenNames.add(compLower)
      const freeLic = bc.licenses?.find(l => l.license_kind === 'free')?.product_name || ''
      const entLic = bc.licenses?.find(l => l.license_kind === 'enterprise')?.product_name || ''
      const contractRef = bc.licenses?.find(l => l.contract_ref)?.contract_ref || ''
      const modBaseline = bc.baselines?.find(b => b.tier === 'moderate')
      const critBaseline = bc.baselines?.find(b => b.tier === 'critical')

      unreg.push({
        id: bc.id,
        name: compName,
        category: bc.layer_name || 'Software Katalog',
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
      if (!name || seenNames.has(nameLower)) continue

      const isRelevant = ['LICENSE', 'SOFTWARE', 'COMPUTE', 'DATABASE', 'MIDDLEWARE', 'SECURITY', 'STORAGE'].includes(
        item.category?.toUpperCase() || ''
      )
      if (isRelevant && !existingComponents.has(nameLower) && !existingLayers.has(nameLower)) {
        seenNames.add(nameLower)
        unreg.push({
          id: item.id || item.sku,
          name: item.name,
          category: item.category,
          sku: item.sku,
          description: `Katalog SKU: ${item.sku} (${item.environment})`,
          source: 'pricing',
          meta: {
            layer: item.category === 'DATABASE' ? 'Database Core OLTP' : (item.category === 'SECURITY' ? 'Security WAF & API Gateway' : item.name),
            licenseEnterprise: item.name,
            contractRef: item.sku
          }
        })
      }
    }
  }

  return unreg
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

  // Standard architectural layers at Bank Mandiri
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
        badge: 'Preset Bank Mandiri',
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
          description: `Lisensi Resmi Bank Mandiri (${pi.sku})`
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
const showModal = ref(false)
const isEditing = ref(false)
const editingId = ref<string | null>(null)

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
  if (!searchQuery.value.trim()) {
    return props.presetItems
  }
  const q = searchQuery.value.toLowerCase().trim()
  return props.presetItems.filter(item => {
    return (
      item.layer.toLowerCase().includes(q) ||
      item.component.toLowerCase().includes(q) ||
      item.licenseFree.toLowerCase().includes(q) ||
      item.licenseEnterprise.toLowerCase().includes(q) ||
      (item.notes || '').toLowerCase().includes(q)
    )
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
  if (confirm(`Apakah Anda yakin ingin menghapus layer preset "${item.layer} (${item.component})"?`)) {
    emit('delete-preset-item', item.id)
  }
}
</script>

<style scoped>
.preset-view-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
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
  gap: 12px;
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
  display: flex;
  flex-direction: column;
  padding: 20px 24px;
  overflow-y: auto;
  gap: 16px;
}

/* ─────────────────────────────────────────────
   New Catalog Alert Banner
   ───────────────────────────────────────────── */
.new-catalog-alert-banner {
  background: linear-gradient(135deg, #FFFDF5 0%, #FEF3C7 100%);
  border: 1px solid #FCD34D;
  border-left: 5px solid #D97706;
  border-radius: 10px;
  padding: 14px 18px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  box-shadow: 0 4px 12px -2px rgba(217, 119, 6, 0.12);
  flex-shrink: 0;
}

.banner-content-left {
  display: flex;
  align-items: flex-start;
  gap: 14px;
  flex: 1;
  min-width: 0;
}

.alert-icon-circle {
  width: 38px;
  height: 38px;
  border-radius: 50%;
  background: #FEF3C7;
  border: 1.5px solid #FDE68A;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  margin-top: 2px;
}

.alert-texts {
  display: flex;
  flex-direction: column;
  gap: 4px;
  flex: 1;
  min-width: 0;
}

.alert-title-row {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.alert-badge {
  display: inline-block;
  padding: 2px 8px;
  background: #FDE68A;
  color: #92400E;
  font-size: 11px;
  font-weight: 800;
  border-radius: 4px;
  text-transform: uppercase;
  letter-spacing: 0.04em;
}

.alert-title {
  margin: 0;
  font-size: 14px;
  font-weight: 800;
  color: #78350F;
  letter-spacing: -0.01em;
}

.alert-sub {
  margin: 0;
  font-size: 12px;
  color: #92400E;
  line-height: 1.4;
}

.new-items-chips {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
  margin-top: 6px;
}

.new-item-chip {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 4px 10px;
  background: #FFFFFF;
  border: 1px solid #FCD34D;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.15s ease;
  font-family: inherit;
}

.new-item-chip:hover {
  background: #FEF3C7;
  border-color: #D97706;
  transform: translateY(-1px);
  box-shadow: 0 2px 4px rgba(217, 119, 6, 0.15);
}

.chip-icon {
  font-size: 13px;
}

.chip-name {
  font-size: 12px;
  font-weight: 700;
  color: #1E293B;
}

.chip-cat {
  font-size: 10px;
  color: #64748B;
  background: #F1F5F9;
  padding: 1px 5px;
  border-radius: 3px;
}

.chip-action {
  font-size: 11px;
  font-weight: 700;
  color: #D97706;
}

.more-chips-pill {
  font-size: 11px;
  color: #92400E;
  font-weight: 600;
}

.banner-actions-right {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-shrink: 0;
}

.btn-input-new-layer {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  background: #002060;
  color: #FFFFFF;
  border: none;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.15s ease;
  box-shadow: 0 2px 4px rgba(0, 32, 96, 0.2);
}

.btn-input-new-layer:hover {
  background: #003399;
  transform: translateY(-1px);
}

.btn-dismiss-alert {
  padding: 8px 12px;
  background: transparent;
  border: 1px solid #FCD34D;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 600;
  color: #92400E;
  cursor: pointer;
  transition: all 0.15s ease;
}

.btn-dismiss-alert:hover {
  background: #FEF3C7;
  color: #78350F;
}

/* Banner Slide Transition */
.banner-slide-enter-active,
.banner-slide-leave-active {
  transition: all 0.25s ease;
}

.banner-slide-enter-from,
.banner-slide-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}

.preset-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  background: #FFFFFF;
  border: 1px solid #E2E8F0;
  border-radius: 8px;
  padding: 10px 16px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.03);
}

.search-input-wrap {
  display: flex;
  align-items: center;
  background: #F8FAFC;
  border: 1px solid #CBD5E1;
  border-radius: 6px;
  padding: 0 12px;
  gap: 8px;
  width: 340px;
  height: 36px;
}

.preset-search-input {
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
  font-size: 12px;
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

.stat-badge {
  font-size: 12px;
  color: #475569;
  background: #F1F5F9;
  border: 1px solid #E2E8F0;
  padding: 4px 10px;
  border-radius: 4px;
}

/* ─────────────────────────────────────────────
   Cards Grid
   ───────────────────────────────────────────── */
.preset-cards-grid {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.preset-card {
  background: #FFFFFF;
  border: 1px solid #E2E8F0;
  border-radius: 8px;
  padding: 16px 18px;
  display: flex;
  flex-direction: column;
  gap: 12px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.02);
  transition: border-color 0.15s ease;
}

.preset-card:hover {
  border-color: #CBD5E1;
}

.preset-card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid #F1F5F9;
  padding-bottom: 10px;
}

.card-header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.layer-number-badge {
  font-size: 11px;
  font-weight: 800;
  color: #002060;
  background: #EFF6FF;
  border: 1px solid #DBEAFE;
  padding: 3px 8px;
  border-radius: 4px;
  text-transform: uppercase;
}

.layer-title-wrap {
  display: flex;
  align-items: baseline;
  gap: 8px;
}

.layer-title {
  font-size: 15px;
  font-weight: 800;
  color: #0F172A;
  margin: 0;
}

.component-subtitle {
  font-size: 12px;
  color: #64748B;
  font-weight: 600;
}

.card-header-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.btn-edit-layer {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  padding: 5px 10px;
  background: #F8FAFC;
  border: 1px solid #CBD5E1;
  border-radius: 4px;
  font-size: 11px;
  font-weight: 600;
  color: #334155;
  cursor: pointer;
  transition: all 0.15s ease;
}

.btn-edit-layer:hover {
  background: #EFF6FF;
  color: #002060;
  border-color: #002060;
}

.btn-delete-layer {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 5px 8px;
  background: #FFFFFF;
  border: 1px solid #CBD5E1;
  border-radius: 4px;
  color: #DC2626;
  cursor: pointer;
}

.btn-delete-layer:hover {
  background: #FEE2E2;
  border-color: #EF4444;
}

/* ─────────────────────────────────────────────
   Specs Comparison Row
   ───────────────────────────────────────────── */
.preset-specs-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}

.preset-specs-row.single-col {
  grid-template-columns: 1fr;
}

.spec-box {
  border-radius: 6px;
  padding: 12px 14px;
  border: 1px solid #E2E8F0;
  display: flex;
  flex-direction: column;
  gap: 10px;
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
  background: #D4A017;
}

.box-tier-title {
  font-size: 12px;
  font-weight: 800;
}

.text-navy {
  color: #002060;
}

.text-amber {
  color: #B45309;
}

.spec-metrics-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 8px;
  background: #FFFFFF;
  border: 1px solid #E2E8F0;
  border-radius: 4px;
  padding: 8px 10px;
}

.metric-cell {
  display: flex;
  flex-direction: column;
}

.cell-lbl {
  font-size: 9px;
  font-weight: 700;
  color: #64748B;
  text-transform: uppercase;
}

.cell-val {
  font-size: 11px;
  font-weight: 700;
  color: #0F172A;
}

.spec-type-strip {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 11px;
}

.type-pill {
  background: rgba(0, 0, 0, 0.04);
  padding: 2px 6px;
  border-radius: 3px;
  font-size: 10px;
  font-weight: 600;
}

.backup-badge {
  font-size: 10px;
  font-weight: 600;
  padding: 2px 6px;
  border-radius: 3px;
}

.backup-badge.active {
  background: #DCFCE7;
  color: #15803D;
}

.backup-badge.inactive {
  background: #F1F5F9;
  color: #64748B;
}

.license-spec-pill {
  font-size: 10px;
  color: #475569;
}

.preset-card-ftr {
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-top: 1px dashed #E2E8F0;
  padding-top: 8px;
}

.ftr-lic-pair {
  display: flex;
  align-items: center;
  gap: 12px;
}

.lic-tag {
  font-size: 11px;
}

.lic-tag.free {
  color: #475569;
}

.lic-tag.enterprise {
  color: #15803D;
}

.ftr-contract-ref {
  font-size: 11px;
  font-weight: 600;
  color: #0369A1;
  background: #E0F2FE;
  padding: 2px 7px;
  border-radius: 4px;
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
</style>
