<template>
  <div class="sketch-workspace-grid">
    <!-- LEFT COLUMN: LAYER LIST WITH NUMBER BADGE & COST PREVIEW -->
    <section class="sketch-list-section">
      <div class="section-hdr">
        <div class="section-hdr-top">
          <div class="section-hdr-title-wrap">
            <span class="hdr-accent-dot"></span>
            <h2 class="section-title">Daftar Komponen Arsitektur</h2>
          </div>
          <span class="item-count">{{ layerItems.length }} Layer</span>
        </div>
        <div class="list-hdr-actions">

          <button
            type="button"
            class="btn-add-catalog-layer"
            @click="openCatalogModal"
            title="Tambah komponen atau lisensi baru dari Katalog Harga"
          >
            <span class="btn-icon-circle">
              <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
                <line x1="12" y1="5" x2="12" y2="19"/>
                <line x1="5" y1="12" x2="19" y2="12"/>
              </svg>
            </span>
            <span class="btn-label-text">Tambah Komponen Katalog</span>
          </button>

          <button
            v-if="hasDeletedDefaultLayers"
            type="button"
            class="btn-restore-defaults"
            @click="$emit('reset-default-layers')"
            title="Pulihkan seluruh 13 layer standar Enterprise Architecture"
          >
            <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M3 12a9 9 0 1 0 9-9 9.75 9.75 0 0 0-6.74 2.74L3 8"/>
              <path d="M3 3v5h5"/>
            </svg>
            <span>Pulihkan 13 Layer Standar</span>
          </button>
        </div>
      </div>

      <div class="bullet-list-container">
        <div
          v-for="item in layerItems"
          :key="item.id"
          :class="['bullet-row', { active: activeItemId === item.id, customized: isLayerCustomized(item.id) }]"
          @click="$emit('select-layer', item.id)"
        >
          <span class="layer-num-badge">{{ item.no }}</span>

          <div class="bullet-content">
            <div class="bullet-title-row">
              <span class="bullet-label">{{ item.layer }}</span>
              <div class="bullet-badges">
                <span v-if="item.no > 13 || item.id.startsWith('layer_custom_')" class="catalog-tag" title="Komponen Katalog Tambahan">Katalog</span>
                <span v-if="isLayerCustomized(item.id)" class="custom-badge" title="Telah disesuaikan dari standar">Kustom</span>
              </div>
            </div>
            <div class="bullet-sub-row">
              <span class="bullet-sub">{{ item.component }}</span>
              <span v-if="getLayerShortCost(item)" class="bullet-cost-pill">
                {{ getLayerShortCost(item) }}
              </span>
            </div>
          </div>

          <button
            type="button"
            class="btn-delete-row"
            :title="`Hapus ${item.layer} (${item.component}) dari daftar Sizing`"
            @click.stop="confirmRemoveLayer(item)"
          >
            <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18"></line>
              <line x1="6" y1="6" x2="18" y2="18"></line>
            </svg>
          </button>
        </div>
      </div>
    </section>

    <!-- RIGHT COLUMN: CLEAN GUIDED CONFIGURATOR -->
    <section class="sketch-card-section">
      <div class="pure-content-card">
        <!-- 1. COMPONENT HEADER BANNER -->
        <div class="component-header-banner">
          <div class="component-header-info">
            <div class="badge-row">
              <span class="badge-layer">LAYER {{ activeItem.no }} OF {{ layerItems.length }} · {{ activeItem.layer }}</span>
              <span v-if="activeConfig.pricingSku" class="badge-sku" title="SKU Katalog Harga Terpilih">
                SKU: {{ activeConfig.pricingSku }}
              </span>
              <span v-if="isLayerCustomized(activeItem.id)" class="custom-status-tag">
                ● Disesuaikan dari Standar EA
              </span>
              <span v-else class="standard-status-tag">
                Standar EA
              </span>
            </div>
            <div class="component-title-row">
              <span v-if="activeItem.icon" class="component-emoji">{{ activeItem.icon }}</span>
              <h2 class="component-title">{{ activeItem.component }}</h2>
            </div>
          </div>

          <div class="layer-action-toolbar">
            <!-- 1. Expand / Collapse Steps -->
            <button
              type="button"
              class="toolbar-action-btn"
              @click="toggleAllSteps"
              :title="isAllCollapsed ? 'Buka semua 4 langkah konfigurasi' : 'Ringkas semua langkah konfigurasi'"
            >
              <svg v-if="isAllCollapsed" width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="15 3 21 3 21 9"></polyline>
                <polyline points="9 21 3 21 3 15"></polyline>
                <line x1="21" y1="3" x2="14" y2="10"></line>
                <line x1="3" y1="21" x2="10" y2="14"></line>
              </svg>
              <svg v-else width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="4 14 10 14 10 20"></polyline>
                <polyline points="20 10 14 10 14 4"></polyline>
                <line x1="14" y1="10" x2="21" y2="3"></line>
                <line x1="3" y1="21" x2="10" y2="14"></line>
              </svg>
              <span>{{ isAllCollapsed ? 'Buka Semua' : 'Ringkas Semua' }}</span>
            </button>

            <span class="toolbar-separator"></span>

            <!-- 2. Reset Layer Configuration to EA Baseline -->
            <button
              type="button"
              class="toolbar-action-btn btn-reset"
              :class="{ 'is-customized': isLayerCustomized(activeItem.id) }"
              @click="$emit('reset-active-layer')"
              title="Kembalikan konfigurasi layer ini ke standar resmi Enterprise Architecture"
            >
              <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2">
                <path d="M3 12a9 9 0 1 0 9-9 9.75 9.75 0 0 0-6.74 2.74L3 8"/>
                <path d="M3 3v5h5"/>
              </svg>
              <span>Reset Standar</span>
              <span v-if="isLayerCustomized(activeItem.id)" class="custom-dot" title="Layer ini telah disesuaikan dari standar EA"></span>
            </button>

            <span class="toolbar-separator"></span>

            <!-- 3. Delete Layer -->
            <button
              type="button"
              class="toolbar-action-btn btn-delete"
              title="Hapus layer ini dari arsitektur sizing"
              @click="confirmRemoveLayer(activeItem)"
            >
              <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M3 6h18"/>
                <path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6"/>
                <path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2"/>
              </svg>
              <span>Hapus</span>
            </button>
          </div>
        </div>

        <!-- 2. CLEAN & MODERN EA ARCHITECTURE BASELINE BANNER -->
        <div class="ea-context-banner">
          <div class="ea-context-left">
            <div class="ea-icon-badge" title="Informasi Komponen Arsitektur Enterprise">
              <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="#002060" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round">
                <circle cx="12" cy="12" r="10"/>
                <line x1="12" y1="16" x2="12" y2="12"/>
                <line x1="12" y1="8" x2="12.01" y2="8"/>
              </svg>
            </div>
            <div class="ea-context-body">
              <div class="ea-desc-row">
                <span class="ea-desc-lead">Fungsi Komponen:</span>
                <span class="ea-desc-text">{{ getFriendlyRoleDescription(activeItem) }}</span>
              </div>
              <div class="ea-meta-row">
                <span class="ea-baseline-chip">
                  <svg width="11" height="11" viewBox="0 0 24 24" fill="currentColor">
                    <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
                  </svg>
                  <span>Acuan Standar EA ({{ tierFilter === 'critical' ? 'Tier Critical' : 'Tier Moderat' }}):</span>
                </span>
                <span class="ea-baseline-spec">{{ cleanBaselineSummary }}</span>
                <span v-if="activeItem.contractRef" class="ea-contract-pill">
                  Ref: {{ activeItem.contractRef }}
                </span>
              </div>
            </div>
          </div>
        </div>

        <!-- ======================================================== -->
        <!-- STEP 1: UKURAN SERVER & KOMPUTASI                        -->
        <!-- ======================================================== -->
        <BaseSectionCard
          step-number="1"
          :title="activeConfig.computeType === 'none' ? 'Kuantitas Unit / Lisensi' : 'Kapasitas & Jumlah Server'"
          :guide-question="activeConfig.computeType === 'none' ? 'Berapa banyak kuantitas unit/node lisensi yang dibutuhkan?' : 'Seberapa besar ukuran per server dan berapa unit server yang diaktifkan?'"
          :technical-label="activeConfig.computeType === 'none' ? 'Unit / Kuantitas' : 'vCPU & RAM'"
          :has-changes="isComputeCustomized"
          :collapsible="true"
          :collapsed="collapsedSteps.step1"
          :collapsed-summary="step1Summary"
          @update:collapsed="setStepCollapsed('step1', $event)"
        >
          <!-- License-only banner notice -->
          <div v-if="activeConfig.computeType === 'none'" class="license-only-notice">
            <span class="notice-icon">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5"/>
              </svg>
            </span>
            <div class="notice-text">
              <span class="notice-title">Mode Lisensi / Add-on Saja (Tanpa Server Baru)</span>
              <span class="notice-desc">Komponen ini dihitung dari jumlah unit/node lisensi tanpa tambahan biaya server vCPU, RAM, atau Storage baru.</span>
            </div>
          </div>

          <!-- Size Preset Selector (Kecil / Sedang / Standar / Menengah+ / Besar / Raksasa / Kustom) -->
          <div v-if="activeConfig.computeType !== 'none'" class="form-row-block">
            <div class="field-hdr-wrap">
              <label class="field-title-lbl">Pilih Ukuran Kapasitas per Server:</label>
              <BaseHelpTip
                title="Ukuran Server (vCPU &amp; RAM)"
                subtitle="Otak Pemroses &amp; Memori Kerja"
                :content="glossary.vcpu.shortDesc"
                :business-impact="glossary.vcpu.businessImpact"
              />
            </div>
            <BaseSizePicker
              :vcpu="activeConfig.vcpuPerNode"
              :ram="activeConfig.ramGBPerNode"
              :recommended-vcpu="recommendedVcpu"
              :recommended-ram="recommendedRam"
              v-model:custom-expanded="isManualComputeOpen"
              @update:vcpu="activeConfig.vcpuPerNode = $event"
              @update:ram="activeConfig.ramGBPerNode = $event"
            />
          </div>

          <!-- Advanced Toggle: Granular Steppers - Automatically opens when "Kustom" is clicked -->
          <BaseAdvancedToggle
            v-if="activeConfig.computeType !== 'none'"
            v-model="isManualComputeOpen"
            show-label="Atur angka vCPU dan RAM secara manual"
            hide-label="Sembunyikan pengaturan angka manual vCPU & RAM"
          >
            <div class="advanced-steppers-grid">
              <BaseStepper
                v-model="activeConfig.vcpuPerNode"
                :min="activeConfig.computeType === 'shared' ? 0 : 1"
                :max="128"
                :step="2"
                label="Otak Pemroses (vCPU per Server)"
                unit="Cores"
                inline-unit="vCPU"
                :disabled="activeConfig.computeType === 'shared'"
              />
              <BaseStepper
                v-model="activeConfig.ramGBPerNode"
                :min="activeConfig.computeType === 'shared' ? 0 : 1"
                :max="1024"
                :step="4"
                label="Memori Kerja (RAM per Server)"
                unit="GB"
                inline-unit="GB"
                :disabled="activeConfig.computeType === 'shared'"
              />
            </div>
          </BaseAdvancedToggle>

          <!-- Stepper: Jumlah Server / Kuantitas Unit with Quick Chips -->
          <div class="server-count-row">
            <div class="field-hdr-wrap">
              <div class="title-with-desc">
                <span class="field-title-lbl">{{ activeConfig.computeType === 'none' ? 'Kuantitas Unit Lisensi:' : 'Jumlah Unit Server Aktif:' }}</span>
                <span class="field-sub-hint">{{ activeConfig.computeType === 'none' ? 'Berapa unit lisensi / agent yang dibutuhkan' : 'Semakin banyak unit server, semakin tahan terhadap gangguan (High Availability)' }}</span>
              </div>
              <BaseHelpTip
                :title="activeConfig.computeType === 'none' ? 'Kuantitas Lisensi (Unit)' : 'Jumlah Server (Nodes)'"
                :subtitle="activeConfig.computeType === 'none' ? 'Software License Allocation' : 'High Availability Redundancy'"
                :content="activeConfig.computeType === 'none' ? 'Jumlah lisensi atau agent software yang dibutuhkan sesuai skala arsitektur.' : glossary.nodeCount.shortDesc"
                :business-impact="activeConfig.computeType === 'none' ? 'Total biaya lisensi tahunan dihitung dari kuantitas unit dikalikan tarif katalog.' : glossary.nodeCount.businessImpact"
              />
            </div>

            <div class="quick-stepper-group">
              <div class="quick-preset-chips">
                <button
                  v-for="cnt in [1, 2, 3, 5, 7]"
                  :key="cnt"
                  type="button"
                  class="quick-chip-btn"
                  :class="{ active: activeConfig.nodeCount === cnt, recommended: recommendedNodeCount === cnt }"
                  @click="activeConfig.nodeCount = cnt"
                >
                  <span>{{ cnt }} {{ activeConfig.computeType === 'none' ? 'Unit' : 'Server' }}</span>
                  <span v-if="recommendedNodeCount === cnt" class="chip-star" title="Standar EA">Std</span>
                </button>
              </div>
              <div class="stepper-wrap-contained">
                <BaseStepper
                  v-model="activeConfig.nodeCount"
                  :min="1"
                  :max="activeConfig.computeType === 'none' ? 2000 : 32"
                  :inline-unit="activeConfig.computeType === 'none' ? 'Unit' : 'Server'"
                />
              </div>
            </div>
          </div>
        </BaseSectionCard>

        <!-- ======================================================== -->
        <!-- STEP 2: PENYIMPANAN DATA & BACKUP CADANGAN              -->
        <!-- ======================================================== -->
        <BaseSectionCard
          step-number="2"
          title="Penyimpanan Data &amp; Proteksi Backup"
          guide-question="Berapa kapasitas ruang simpan (harddisk) per server dan apakah perlu backup harian?"
          technical-label="SAN Storage &amp; NetBackup"
          :has-changes="isStorageCustomized"
          :collapsible="true"
          :collapsed="collapsedSteps.step2"
          :collapsed-summary="step2Summary"
          @update:collapsed="setStepCollapsed('step2', $event)"
        >
          <!-- Notice if computeType === 'none' -->
          <div v-if="activeConfig.computeType === 'none'" class="license-only-notice">
            <span class="notice-icon">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M22 12H2M5.45 5.11L2 12v6a2 2 0 0 0 2 2h16a2 2 0 0 0 2-2v-6l-3.45-6.89A2 2 0 0 0 16.76 4H7.24a2 2 0 0 0-1.79 1.11z"/>
              </svg>
            </span>
            <div class="notice-text">
              <span class="notice-title">Tidak Memerlukan Disk Storage / Backup Tambahan</span>
              <span class="notice-desc">Komponen berjenis lisensi / add-on menumpang pada server yang sudah ada, sehingga kapasitas disk diatur 0 GB (Biaya Storage &amp; Backup: Rp 0).</span>
            </div>
          </div>

          <template v-else>
            <!-- Kapasitas Storage Stepper + Quick Chips -->
            <div class="storage-stepper-row">
              <div class="field-hdr-wrap">
                <div class="title-with-desc">
                  <span class="field-title-lbl">Kapasitas Disk per Server:</span>
                  <span class="field-sub-hint">Total semua server: <strong>{{ (activeConfig.nodeCount * activeConfig.storageGBPerNode).toLocaleString('id-ID') }} GB</strong> ({{ ((activeConfig.nodeCount * activeConfig.storageGBPerNode) / 1000).toFixed(2) }} TB)</span>
                </div>
                <BaseHelpTip
                  title="SAN Storage (Kapasitas Disk)"
                  subtitle="SAN Virtual Disk LUN"
                  :content="glossary.storage.shortDesc"
                  :business-impact="glossary.storage.businessImpact"
                />
              </div>

              <div class="quick-stepper-group">
                <div class="quick-preset-chips">
                  <button
                    v-for="gb in storageQuickOptions"
                    :key="gb"
                    type="button"
                    class="quick-chip-btn"
                    :class="{ active: activeConfig.storageGBPerNode === gb, recommended: recommendedStorageGB === gb }"
                    @click="activeConfig.storageGBPerNode = gb"
                  >
                    <span>{{ gb >= 1000 ? `${gb / 1000} TB` : `${gb} GB` }}</span>
                    <span v-if="recommendedStorageGB === gb" class="chip-star" title="Standar EA">Std</span>
                  </button>
                </div>
                <div class="stepper-wrap-contained">
                  <BaseStepper
                    v-model="activeConfig.storageGBPerNode"
                    :min="10"
                    :max="20000"
                    :step="50"
                    inline-unit="GB"
                  />
                </div>
              </div>
            </div>

            <!-- Pilihan Kecepatan Storage & Backup -->
            <div class="storage-options-grid">
              <!-- Pilihan Kelas Kecepatan Storage -->
              <div class="option-column">
                <div class="field-hdr-wrap">
                  <span class="field-title-lbl">Kelas Kecepatan Disk:</span>
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
                      technicalLabel: 'SAN Tier 2 (Rp 33rb/GB)',
                      description: 'Kecepatan standar enterprise untuk web server, container worker, cache, dan log.',
                      icon: 'storage',
                      recommended: recommendedStorageClass === 'Tier 2'
                    },
                    {
                      value: 'Tier 3',
                      title: 'Ultra Cepat (Database)',
                      technicalLabel: 'SAN Tier 3 ASM (Rp 48rb/GB)',
                      description: 'Penyimpanan performa tinggi khusus database transaksi utama (Oracle ASM).',
                      icon: 'database',
                      badge: 'High IOPS',
                      badgeType: 'warning',
                      recommended: recommendedStorageClass === 'Tier 3'
                    }
                  ]"
                />
              </div>

              <!-- Toggle Proteksi Backup NetBackup / Rubrik -->
              <div class="backup-protection-box">
                <div class="field-hdr-wrap">
                  <span class="field-title-lbl">Proteksi Cadangan Data (Backup):</span>
                  <BaseHelpTip
                    :title="activeItem.component.toLowerCase().includes('nginx') ? 'Rubrik Backup' : 'NetBackup (NBU)'"
                    subtitle="Enterprise Daily Backup"
                    :content="glossary.backupNBU.shortDesc"
                    :business-impact="glossary.backupNBU.businessImpact"
                  />
                </div>

                <label class="backup-card-label" :class="{ checked: activeConfig.backupEnabled }">
                  <div class="backup-top-row">
                    <div class="backup-left-title">
                      <input
                        type="checkbox"
                        v-model="activeConfig.backupEnabled"
                        class="custom-chk-large"
                      />
                      <span class="backup-main-title">
                        {{ activeItem.component.toLowerCase().includes('nginx') ? 'Backup Otomatis (Rubrik Backup)' : 'Backup Otomatis Harian (NetBackup)' }}
                      </span>
                    </div>
                    <span :class="['backup-status-pill', activeConfig.backupEnabled ? 'on' : 'off']">
                      {{ activeConfig.backupEnabled ? 'Aktif' : 'Nonaktif' }}
                    </span>
                  </div>
                  <span class="backup-desc">
                    Mencadangkan otomatis <strong>{{ (activeConfig.nodeCount * activeConfig.storageGBPerNode).toLocaleString('id-ID') }} GB</strong> data setiap hari agar aman dari kehilangan data ({{ activeItem.component.toLowerCase().includes('nginx') ? 'Tarif Rubrik: Rp 35.000/GB' : 'Tarif NetBackup: Rp 27.000/GB' }}).
                  </span>
                </label>
              </div>
            </div>
          </template>
        </BaseSectionCard>

        <!-- ======================================================== -->
        <!-- STEP 3: JENIS LINGKUNGAN INFRASTRUKTUR                  -->
        <!-- ======================================================== -->
        <BaseSectionCard
          step-number="3"
          title="Jenis Lingkungan Infrastruktur"
          guide-question="Di atas platform mesin seperti apa komponen ini akan dijalankan?"
          technical-label="VM / Baremetal / Shared / Lisensi Saja"
          :has-changes="activeConfig.computeType !== recommendedComputeType"
          :collapsible="true"
          :collapsed="collapsedSteps.step3"
          :collapsed-summary="step3Summary"
          @update:collapsed="setStepCollapsed('step3', $event)"
        >
          <BaseChoiceCard
            v-model="activeConfig.computeType"
            :columns="2"
            :options="[
              {
                value: 'vm',
                title: 'Virtual Machine (VM)',
                technicalLabel: 'Server Virtual Fleksibel',
                description: 'Server virtual yang hemat biaya dan mudah diatur ukurannya. Cocok untuk sebagian besar aplikasi.',
                icon: 'server',
                recommended: recommendedComputeType === 'vm'
              },
              {
                value: 'baremetal',
                title: 'Server Fisik Khusus',
                technicalLabel: 'Baremetal Dedicated',
                description: 'Satu mesin server fisik penuh tanpa dibagi dengan aplikasi lain. Performa paling tinggi.',
                icon: 'baremetal',
                badge: 'Dedicated',
                badgeType: 'warning',
                recommended: recommendedComputeType === 'baremetal'
              }
            ]"
          />
        </BaseSectionCard>

        <!-- ======================================================== -->
        <!-- STEP 4: MODEL LISENSI SOFTWARE & SKU KATALOG             -->
        <!-- ======================================================== -->
        <BaseSectionCard
          step-number="4"
          title="Pilihan Lisensi Software &amp; Dukungan Vendor"
          guide-question="Gunakan software komunitas gratis atau versi berlisensi resmi dengan garansi dukungan vendor 24/7?"
          technical-label="OSS vs Enterprise / SKU"
          :has-changes="isLicenseCustomized"
          :collapsible="true"
          :collapsed="collapsedSteps.step4"
          :collapsed-summary="step4Summary"
          @update:collapsed="setStepCollapsed('step4', $event)"
        >
          <div class="license-choice-grid">
            <!-- Option 1: Community Free OSS -->
            <div
              class="license-choice-card"
              :class="{ active: activeConfig.licenseKind === 'free' }"
              @click="activeConfig.licenseKind = 'free'"
            >
              <div class="lic-card-hdr">
                <div class="lic-card-hdr-left">
                  <span class="lic-icon">
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <polyline points="16 18 22 12 16 6"></polyline>
                      <polyline points="8 6 2 12 8 18"></polyline>
                    </svg>
                  </span>
                  <div class="lic-hdr-titles">
                    <span class="lic-name">Community / Open Source (Gratis)</span>
                    <span class="lic-price-tag">Biaya Lisensi: Rp 0 / tahun</span>
                  </div>
                </div>
                <span class="lic-radio-indicator" :class="{ checked: activeConfig.licenseKind === 'free' }">
                  <svg v-if="activeConfig.licenseKind === 'free'" width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3.5" stroke-linecap="round" stroke-linejoin="round">
                    <polyline points="20 6 9 17 4 12"></polyline>
                  </svg>
                </span>
              </div>
              <p class="lic-product-name">{{ activeItem.license?.free || 'Community (Free OSS)' }}</p>
              <p class="lic-card-desc">Tanpa biaya langganan software tahunan. Pemeliharaan dan penanganan kendala dikelola mandiri oleh tim internal.</p>
            </div>

            <!-- Option 2: Enterprise Official License -->
            <div
              class="license-choice-card enterprise-card"
              :class="{ active: activeConfig.licenseKind === 'enterprise' }"
              @click="activeConfig.licenseKind = 'enterprise'"
            >
              <div class="lic-card-hdr">
                <div class="lic-card-hdr-left">
                  <span class="lic-icon">
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <rect x="4" y="2" width="16" height="20" rx="2" ry="2"></rect>
                      <line x1="9" y1="22" x2="9" y2="22.01"></line>
                      <line x1="15" y1="22" x2="15" y2="22.01"></line>
                    </svg>
                  </span>
                  <div class="lic-hdr-titles">
                    <span class="lic-name">Enterprise / Berlisensi Resmi</span>
                    <span class="lic-price-tag font-emerald">Garansi Support Vendor 24/7</span>
                  </div>
                </div>
                <span class="lic-radio-indicator emerald" :class="{ checked: activeConfig.licenseKind === 'enterprise' }">
                  <svg v-if="activeConfig.licenseKind === 'enterprise'" width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3.5" stroke-linecap="round" stroke-linejoin="round">
                    <polyline points="20 6 9 17 4 12"></polyline>
                  </svg>
                </span>
              </div>
              <p class="lic-product-name">{{ activeConfig.licenseProductName || activeItem.license?.enterprise || 'Enterprise License' }}</p>
              <p class="lic-card-desc">Langganan resmi vendor dengan jaminan SLA perbaikan cepat, pembaruan keamanan berkala, dan kepatuhan audit perbankan.</p>
            </div>
          </div>

        </BaseSectionCard>

        <!-- ======================================================== -->
        <!-- LIVE ESTIMATED COST BREAKDOWN WITH PROPORTION BAR        -->
        <!-- ======================================================== -->
        <div class="layer-cost-card">
          <div class="cost-card-header">
            <div class="cost-card-title-wrap">
              <span class="cost-card-title">Ringkasan Biaya Sizing (Perpetual) · Layer {{ activeItem.no }} ({{ activeItem.layer }})</span>
              <span class="cost-card-sub">Dihitung otomatis berdasarkan pilihan kapasitas di atas dan tarif katalog resmi</span>
            </div>
            <div class="cost-card-total-badge">
              <span class="total-lbl">Total Layer Ini:</span>
              <BaseMoney
                :amount="activeLayerCost?.total_cost_annual"
                class="layer-grand-price"
              />
            </div>
          </div>

          <!-- Visual Cost Proportion Bar -->
          <div v-if="(activeLayerCost?.total_cost_annual || 0) > 0" class="cost-proportion-bar" title="Proporsi komponen biaya pada layer ini">
            <div
              v-if="costPercentages.compute > 0"
              class="prop-segment seg-compute"
              :style="{ width: `${costPercentages.compute}%` }"
            ></div>
            <div
              v-if="costPercentages.storage > 0"
              class="prop-segment seg-storage"
              :style="{ width: `${costPercentages.storage}%` }"
            ></div>
            <div
              v-if="costPercentages.backup > 0"
              class="prop-segment seg-backup"
              :style="{ width: `${costPercentages.backup}%` }"
            ></div>
            <div
              v-if="costPercentages.license > 0"
              class="prop-segment seg-license"
              :style="{ width: `${costPercentages.license}%` }"
            ></div>
          </div>

          <div class="cost-metrics-row">
            <!-- Compute Cost -->
            <div class="cost-metric-box">
              <div class="metric-hdr">
                <span class="metric-dot dot-compute"></span>
                <span class="metric-title">Biaya Server (Compute)</span>
              </div>
              <BaseMoney :amount="activeLayerCost?.compute_cost_annual" />
              <span class="metric-sub">
                {{ activeConfig.computeType === 'none' ? 'Tanpa Server Baru' : `${activeConfig.nodeCount} Server · Total ${activeConfig.nodeCount * activeConfig.vcpuPerNode} vCPU` }}
              </span>
            </div>

            <!-- Storage Cost -->
            <div class="cost-metric-box">
              <div class="metric-hdr">
                <span class="metric-dot dot-storage"></span>
                <span class="metric-title">Biaya Penyimpanan</span>
              </div>
              <BaseMoney :amount="activeLayerCost?.storage_cost_annual" />
              <span class="metric-sub">
                {{ activeConfig.computeType === 'none' ? '0 GB' : `Total ${(activeConfig.nodeCount * activeConfig.storageGBPerNode).toLocaleString('id-ID')} GB (${activeConfig.storageClass})` }}
              </span>
            </div>

            <!-- Backup Cost -->
            <div class="cost-metric-box">
              <div class="metric-hdr">
                <span class="metric-dot dot-backup"></span>
                <span class="metric-title">{{ activeItem.component.toLowerCase().includes('nginx') ? 'Biaya Backup (Rubrik)' : 'Biaya Backup (NBU)' }}</span>
              </div>
              <BaseMoney :amount="activeLayerCost?.backup_cost_annual" />
              <span class="metric-sub">{{ activeConfig.backupEnabled ? (activeItem.component.toLowerCase().includes('nginx') ? 'Rubrik Backup Aktif' : 'NetBackup Harian Aktif') : 'Tanpa Backup' }}</span>
            </div>

            <!-- License Cost -->
            <div class="cost-metric-box">
              <div class="metric-hdr">
                <span class="metric-dot dot-license"></span>
                <span class="metric-title">Biaya Lisensi Software</span>
              </div>
              <template v-if="activeConfig.licenseKind === 'free'">
                <span class="free-badge">Rp 0 (Gratis / OSS)</span>
              </template>
              <template v-else-if="activeLayerCost?.unconfirmed_prices && activeLayerCost.unconfirmed_prices.length > 0">
                <BaseMoney :amount="null" :is-unconfirmed="true" unconfirmed-text="Pending STA" />
              </template>
              <template v-else>
                <BaseMoney :amount="activeLayerCost?.license_cost_annual" />
              </template>
              <span class="metric-sub text-truncate" :title="activeLayerCost?.license_product_name || activeItem.license.enterprise">
                {{ activeConfig.licenseKind === 'free' ? activeItem.license.free : (activeLayerCost?.license_product_name || activeItem.license.enterprise) }}
              </span>
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
        <!-- GUIDED WIZARD FOOTER: PREV / NEXT LAYER NAVIGATION       -->
        <!-- ======================================================== -->
        <div class="layer-wizard-nav">
          <button
            type="button"
            class="wizard-nav-btn prev"
            :disabled="!prevLayerItem"
            @click="prevLayerItem && $emit('select-layer', prevLayerItem.id)"
          >
            <span class="nav-arrow">←</span>
            <div v-if="prevLayerItem" class="nav-text-wrap">
              <span class="nav-hint">Layer Sebelumnya ({{ prevLayerItem.no }})</span>
              <span class="nav-name">{{ prevLayerItem.layer }} · {{ prevLayerItem.component }}</span>
            </div>
            <div v-else class="nav-text-wrap">
              <span class="nav-hint">Awal Daftar</span>
              <span class="nav-name">Anda berada di Layer 1</span>
            </div>
          </button>

          <div class="wizard-progress-indicator">
            <span class="prog-text">Layer <strong>{{ activeLayerIndex + 1 }}</strong> dari <strong>{{ layerItems.length }}</strong></span>
          </div>

          <button
            type="button"
            class="wizard-nav-btn next"
            :disabled="!nextLayerItem"
            @click="nextLayerItem && $emit('select-layer', nextLayerItem.id)"
          >
            <div v-if="nextLayerItem" class="nav-text-wrap align-right">
              <span class="nav-hint">Lanjut ke Layer {{ nextLayerItem.no }}</span>
              <span class="nav-name">{{ nextLayerItem.layer }} · {{ nextLayerItem.component }}</span>
            </div>
            <div v-else class="nav-text-wrap align-right">
              <span class="nav-hint">Selesai</span>
              <span class="nav-name">Layer Terakhir</span>
            </div>
            <span class="nav-arrow">→</span>
          </button>
        </div>
      </div>
    </section>

    <!-- MODAL: PILIH DARI KATALOG HARGA -->
    <div v-if="showCatalogModal" class="catalog-modal-backdrop" @click.self="showCatalogModal = false">
      <div class="catalog-modal-card">
        <div class="modal-hdr">
          <div class="modal-title-wrap">
            <div class="modal-icon-badge">
              <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="#002060" stroke-width="2">
                <polygon points="12 2 2 7 12 12 22 7 12 2"/>
                <polyline points="2 17 12 22 22 17"/>
                <polyline points="2 12 12 17 22 12"/>
              </svg>
            </div>
            <div>
              <h3 class="modal-title">Katalog Komponen &amp; Lisensi Software</h3>
              <p class="modal-subtitle">Pilih produk software, agent monitoring, atau add-on dari master katalog untuk dimasukkan ke lembar Sizing</p>
            </div>
          </div>
          <button type="button" class="modal-btn-close" @click="showCatalogModal = false" title="Tutup">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18"></line>
              <line x1="6" y1="6" x2="18" y2="18"></line>
            </svg>
          </button>
        </div>

        <!-- Search & Filter Controls -->
        <div class="modal-search-strip">
          <div class="modal-search-wrap">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="#64748B" stroke-width="2">
              <circle cx="11" cy="11" r="8"/>
              <line x1="21" y1="21" x2="16.65" y2="16.65"/>
            </svg>
            <input
              type="text"
              v-model="catalogModalSearch"
              class="modal-search-input"
              placeholder="Cari SKU atau nama produk (contoh: AppDynamics, APM, Redis, Oracle)..."
            />
            <button
              v-if="catalogModalSearch"
              type="button"
              class="modal-search-clear"
              @click="catalogModalSearch = ''"
            >
              <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="18" y1="6" x2="6" y2="18"></line>
                <line x1="6" y1="6" x2="18" y2="18"></line>
              </svg>
            </button>
          </div>

          <div class="modal-filter-pills">
            <button
              type="button"
              :class="['filter-pill', { active: catalogModalFilter === 'unmapped' }]"
              @click="catalogModalFilter = 'unmapped'"
            >
              <span class="pill-glow-dot"></span>
              <span>Siap Ditambahkan</span>
              <span class="pill-count">{{ unmappedCatalogItems.length }}</span>
            </button>
            <button
              type="button"
              :class="['filter-pill', { active: catalogModalFilter === 'license' }]"
              @click="catalogModalFilter = 'license'"
            >
              <span>Lisensi Software</span>
              <span class="pill-count">{{ licenseCatalogItems.length }}</span>
            </button>
            <button
              type="button"
              :class="['filter-pill', { active: catalogModalFilter === 'infra' }]"
              @click="catalogModalFilter = 'infra'"
            >
              <span>Tarif Dasar Infra</span>
              <span class="pill-count">{{ infraRatesCatalogItems.length }}</span>
            </button>
            <button
              type="button"
              :class="['filter-pill', { active: catalogModalFilter === 'all' }]"
              @click="catalogModalFilter = 'all'"
            >
              <span>Semua Master</span>
              <span class="pill-count">{{ availableCatalogItems.length }}</span>
            </button>
          </div>
        </div>

        <!-- Catalog List Container -->
        <div class="modal-items-scroll">
          <div v-if="filteredModalCatalogItems.length === 0" class="empty-catalog-state">
            <span class="empty-icon">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="11" cy="11" r="8"></circle>
                <line x1="21" y1="21" x2="16.65" y2="16.65"></line>
              </svg>
            </span>
            <p class="empty-text">
              <template v-if="catalogModalFilter === 'unmapped'">
                <strong>Semua komponen katalog sudah terpetakan di Sizing.</strong><br>
                Pilih tab <em>"Lisensi Software"</em> atau <em>"Semua Master"</em> jika ingin menambah instance layer tambahan.
              </template>
              <template v-else>
                Tidak ada item katalog yang cocok dengan kata kunci pencarian.
              </template>
            </p>
          </div>

          <div
            v-for="cat in filteredModalCatalogItems"
            :key="cat.sku"
            class="catalog-picker-card"
            :class="{ 'card-installed': getCatalogItemSizingInfo(cat).isAlreadyInSizing, 'card-infra': getCatalogItemSizingInfo(cat).isRawInfra }"
          >
            <!-- 1. Category Avatar Icon -->
            <div class="card-icon-avatar" :class="getCategoryTheme(cat)">
              <svg v-if="cat.sku.includes('APPDYN') || cat.name.toLowerCase().includes('agent') || cat.name.toLowerCase().includes('apm')" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M18 20V10M12 20V4M6 20v-6"/>
              </svg>
              <svg v-else-if="cat.category === 'database' || cat.sku.includes('ORACLE') || cat.sku.includes('REDIS')" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <ellipse cx="12" cy="5" rx="9" ry="3"/>
                <path d="M21 12c0 1.66-4 3-9 3s-9-1.34-9-3"/>
                <path d="M3 5v14c0 1.66 4 3 9 3s9-1.34 9-3V5"/>
              </svg>
              <svg v-else-if="cat.category === 'security' || cat.sku.includes('VAULT') || cat.sku.includes('NEUVECTOR')" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/>
              </svg>
              <svg v-else-if="cat.category === 'backup' || cat.sku.includes('RUBRIK') || cat.sku.includes('NBU')" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"/>
              </svg>
              <svg v-else-if="cat.category === 'compute' || cat.category === 'storage' || cat.category === 'baremetal'" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <rect x="2" y="2" width="20" height="8" rx="2" ry="2"/>
                <rect x="2" y="14" width="20" height="8" rx="2" ry="2"/>
                <line x1="6" y1="6" x2="6.01" y2="6"/>
                <line x1="6" y1="18" x2="6.01" y2="18"/>
              </svg>
              <svg v-else width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M21 2l-2 2m-1.5 1.5L14 9l-1.5-1.5-2 2 1.5 1.5L10 13H8v2H6v2H4v2h2v-2h2v-2h2l4-4 1.5 1.5 2-2-1.5-1.5L21 2z"/>
              </svg>
            </div>

            <!-- 2. Main Body Info -->
            <div class="card-main-info">
              <div class="card-meta-tags">
                <span class="sku-pill">{{ cat.sku }}</span>
                <span class="category-pill" :class="cat.category">{{ getCategoryLabel(cat) }}</span>
                <span class="env-pill" :class="cat.environment === 'Production' ? 'env-prod' : 'env-dev'">{{ cat.environment }}</span>
                <span v-if="getCatalogItemSizingInfo(cat).isAlreadyInSizing" class="status-pill installed">
                  Terpasang di {{ getCatalogItemSizingInfo(cat).matchedLayerText }}
                </span>
                <span v-else-if="getCatalogItemSizingInfo(cat).isRawInfra" class="status-pill infra">
                  Tarif Dasar Mesin &amp; Storage
                </span>
                <span v-else class="status-pill ready">
                  Siap Dimasukkan
                </span>
              </div>

              <h4 class="card-product-title">{{ cat.name }}</h4>

              <p class="card-description-text">
                {{ getItemExplanation(cat) }}
              </p>

              <div class="card-feature-hints">
                <div class="feature-item">
                  <span class="feature-label">Satuan Hitung:</span>
                  <span class="feature-val"><strong>{{ cat.unit || 'Unit / Thn' }}</strong></span>
                </div>
                <div class="feature-divider"></div>
                <div class="feature-item">
                  <span class="feature-label">Mode Sizing:</span>
                  <span class="feature-val highlight">{{ getSizingTypeHint(cat) }}</span>
                </div>
              </div>
            </div>

            <!-- 3. Pricing & Actions -->
            <div class="card-action-block">
              <div class="price-box">
                <span class="price-header">Tarif Tahunan</span>
                <div v-if="cat.unit_price_idr" class="price-val-row">
                  <span class="price-curr">Rp</span>
                  <span class="price-num">{{ formatSimpleIDR(cat.unit_price_idr) }}</span>
                </div>
                <div v-else class="price-pending-tag">Pending STA Tim</div>
                <span class="price-period">per {{ cat.unit || 'unit' }}</span>
              </div>

              <button
                v-if="!getCatalogItemSizingInfo(cat).isAlreadyInSizing && !getCatalogItemSizingInfo(cat).isRawInfra"
                type="button"
                class="btn-card-select"
                @click="pickCatalogItem(cat)"
              >
                <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
                  <line x1="12" y1="5" x2="12" y2="19"/>
                  <line x1="5" y1="12" x2="19" y2="12"/>
                </svg>
                <span>Masukkan ke Sizing</span>
              </button>
              <button
                v-else-if="getCatalogItemSizingInfo(cat).isAlreadyInSizing"
                type="button"
                class="btn-card-select btn-installed-action"
                title="Komponen ini sudah ada di daftar layer. Klik jika ingin menambahkan instance layer baru."
                @click="pickCatalogItem(cat)"
              >
                <span>+ Tambah Instance Baru</span>
              </button>
              <div v-else class="infra-auto-hint">
                Dihitung via Step 1-3
              </div>
            </div>
          </div>
        </div>

        <div class="modal-ftr">
          <span class="modal-ftr-count">Menampilkan {{ filteredModalCatalogItems.length }} item katalog</span>
          <button type="button" class="btn-modal-cancel" @click="showCatalogModal = false">Tutup</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import type { TechStackItem } from '~/data/techStackBaselineData'
import type { LayerAdjusterConfig, QuoteItemCost, PricingItemRecord } from '~/types/sizing'
import { DEFAULT_PRESET_ITEMS } from '~/data/defaultPresetItems'
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
  allLayerCosts?: QuoteItemCost[]
  tierFilter: 'critical' | 'moderate'
  licenseFilter: 'free' | 'license'
  isLayerCustomized: (id: string) => boolean
  pricingList?: PricingItemRecord[]
}>()

const emit = defineEmits<{
  (e: 'select-layer', id: string): void
  (e: 'reset-active-layer'): void
  (e: 'reset-default-layers'): void
  (e: 'add-catalog-to-sizing', item: PricingItemRecord): void
  (e: 'remove-custom-layer', id: string): void
  (e: 'update:tier', val: 'critical' | 'moderate'): void
  (e: 'update:license', val: 'free' | 'license'): void
}>()

const glossary = SIZING_GLOSSARY

// UI states for non-tech friendly guided experience
const showBaselineDetails = ref(false)
const isManualComputeOpen = ref(false)
const storageQuickOptions = [50, 100, 150, 200, 500, 1000, 2000]

// Modal state
const showCatalogModal = ref(false)
const catalogModalSearch = ref('')
const catalogModalFilter = ref<'all' | 'license' | 'infra' | 'unmapped'>('unmapped')

function openCatalogModal() {
  catalogModalSearch.value = ''
  catalogModalFilter.value = 'unmapped'
  showCatalogModal.value = true
}

const hasDeletedDefaultLayers = computed(() => {
  return DEFAULT_PRESET_ITEMS.some(def => !props.layerItems.some(item => item.id === def.id))
})

function confirmRemoveLayer(item: TechStackItem) {
  if (props.layerItems.length <= 1) {
    alert('Tidak dapat menghapus layer. Arsitektur sizing minimal harus memiliki 1 komponen.')
    return
  }
  if (
    confirm(
      `Hapus layer "${item.layer} - ${item.component}" dari daftar Sizing?\n\nKomponen ini beserta seluruh alokasi biaya server, storage, dan lisensinya akan dikeluarkan dari perhitungan arsitektur proyek ini.`
    )
  ) {
    emit('remove-custom-layer', item.id)
  }
}

function pickCatalogItem(cat: PricingItemRecord) {
  emit('add-catalog-to-sizing', cat)
  showCatalogModal.value = false
}

function formatSimpleIDR(val?: number | null): string {
  if (val === null || val === undefined || isNaN(val)) return '0'
  return new Intl.NumberFormat('id-ID').format(Math.round(val))
}

function getLayerShortCost(item: TechStackItem): string {
  if (!props.allLayerCosts) return ''
  const found =
    props.allLayerCosts.find(c => c.component_id === item.id) ||
    props.allLayerCosts.find(c => c.component_id === `layer-${item.no}`)
  if (!found || found.total_cost_annual === undefined || found.total_cost_annual === null) return ''
  const val = found.total_cost_annual
  if (val === 0) return 'Rp 0'
  if (val >= 1_000_000_000) {
    return `Rp ${(val / 1_000_000_000).toLocaleString('id-ID', { maximumFractionDigits: 1 })} M`
  }
  if (val >= 1_000_000) {
    return `Rp ${(val / 1_000_000).toLocaleString('id-ID', { maximumFractionDigits: 0 })} Jt`
  }
  return `Rp ${formatSimpleIDR(val)}`
}

function getFriendlyRoleDescription(item: TechStackItem): string {
  const map: Record<string, string> = {
    web_ingress: 'Pintu gerbang utama penerima trafik internet & pengaman koneksi aplikasi (Load Balancer / Reverse Proxy).',
    app_container: 'Mesin utama tempat seluruh aplikasi dan layanan transaksi (microservices) dijalankan.',
    container_router: 'Pengatur jalur lalu lintas jaringan internal antar-aplikasi di dalam cluster OpenShift.',
    core_db: 'Brankas penyimpanan data utama untuk seluruh catatan transaksi nasabah (Core Database).',
    db_storage: 'Ruang harddisk berkecepatan ekstra tinggi yang dipakai bersama oleh mesin database utama.',
    dist_cache: 'Memori penyimpanan super cepat sementara agar aplikasi tidak perlu berulang kali membaca database.',
    cache_sentinel: 'Penjaga otomatis yang langsung memindahkan tugas jika salah satu server cache mengalami gangguan.',
    message_bus: 'Jalur antrean pengiriman pesan transaksi berkecepatan tinggi antar-sistem (Event Streaming).',
    log_analytics: 'Pusat pencarian dan analisis catatan aktivitas sistem (log audit) secara real-time.',
    log_shipper: 'Kurir pengumpul dan perapi catatan log dari semua server sebelum disimpan ke pusat analitik.',
    log_dashboard: 'Layar pemantauan visual (dashboard) untuk melihat kesehatan sistem dan grafik log.',
    doc_storage: 'Gudang penyimpanan berkas dokumen elektronik seperti e-statement, PDF, dan foto KYC.',
    batch_engine: 'Pengatur jadwal pekerjaan otomatis harian/bulanan (seperti tutup buku & rekonsiliasi akhir hari).'
  }
  return map[item.id] || item.notes || 'Komponen infrastruktur pendukung arsitektur aplikasi.'
}

// Navigation between layers (Prev / Next)
const activeLayerIndex = computed(() => {
  const idx = props.layerItems.findIndex(i => i.id === props.activeItemId)
  return idx >= 0 ? idx : 0
})

const prevLayerItem = computed<TechStackItem | null>(() => {
  return activeLayerIndex.value > 0 ? props.layerItems[activeLayerIndex.value - 1] : null
})

const nextLayerItem = computed<TechStackItem | null>(() => {
  return activeLayerIndex.value < props.layerItems.length - 1 ? props.layerItems[activeLayerIndex.value + 1] : null
})

// Cost proportions for visual bar
const costPercentages = computed(() => {
  const total = props.activeLayerCost?.total_cost_annual || 0
  if (total <= 0) {
    return { compute: 0, storage: 0, backup: 0, license: 0 }
  }
  const c = ((props.activeLayerCost?.compute_cost_annual || 0) / total) * 100
  const s = ((props.activeLayerCost?.storage_cost_annual || 0) / total) * 100
  const b = ((props.activeLayerCost?.backup_cost_annual || 0) / total) * 100
  const l = ((props.activeLayerCost?.license_cost_annual || 0) / total) * 100
  return { compute: c, storage: s, backup: b, license: l }
})

const availableCatalogItems = computed<PricingItemRecord[]>(() => {
  return props.pricingList || []
})

interface CatalogSizingStatus {
  isAlreadyInSizing: boolean
  isRawInfra: boolean
  matchedLayerText?: string
}

function getCatalogItemSizingInfo(item: PricingItemRecord): CatalogSizingStatus {
  const catLower = (item.category || '').toLowerCase()
  const skuUpper = (item.sku || '').toUpperCase()
  const nameLower = (item.name || '').toLowerCase()
  const unitLower = (item.unit || '').toLowerCase()

  const isRawInfra =
    ['compute', 'storage', 'baremetal'].includes(catLower) ||
    ['VCPU_PROD', 'VCPU_DEV', 'BM_STD_PROD', 'BM_STD_DEV', 'STG_T2_PROD', 'STG_T2_DEV', 'STG_T3_PROD', 'NBU_ALL', 'RUBRIK_PROD'].includes(skuUpper) ||
    (catLower === 'backup' && (unitLower.includes('gb') || unitLower.includes('tb'))) ||
    nameLower.includes('compute vcpu') ||
    nameLower.includes('san storage') ||
    nameLower.includes('server baremetal') ||
    nameLower.includes('backup netbackup') ||
    nameLower.includes('rubrik backup')

  const directMatch = (props.layerItems || []).find(l => {
    return (l.contractRef && l.contractRef.toLowerCase() === item.sku.toLowerCase()) ||
      (l.component && l.component.toLowerCase() === nameLower)
  })
  if (directMatch) {
    return {
      isAlreadyInSizing: true,
      isRawInfra,
      matchedLayerText: `Layer ${directMatch.no} (${directMatch.layer})`
    }
  }

  const productKeywordFamilies: { key: string; keywords: string[]; layerKeywords: string[] }[] = [
    { key: 'nginx', keywords: ['nginx'], layerKeywords: ['nginx'] },
    { key: 'ocp', keywords: ['openshift', 'ocp'], layerKeywords: ['ocp', 'openshift'] },
    { key: 'kafka', keywords: ['kafka', 'confluent'], layerKeywords: ['kafka', 'confluent'] },
    { key: 'redis', keywords: ['redis'], layerKeywords: ['redis'] },
    { key: 'oracle', keywords: ['oracle'], layerKeywords: ['oracle'] },
    { key: 'minio', keywords: ['minio'], layerKeywords: ['minio'] },
    { key: 'elastic', keywords: ['elastic', 'elasticsearch', 'kibana', 'logstash'], layerKeywords: ['elastic', 'kibana', 'logstash'] },
    { key: 'controlm', keywords: ['control-m', 'controlm'], layerKeywords: ['control-m', 'controlm'] },
    { key: 'vault', keywords: ['vault'], layerKeywords: ['vault'] },
    { key: 'dynatrace', keywords: ['dynatrace'], layerKeywords: ['dynatrace'] },
    { key: 'neuvector', keywords: ['neuvector'], layerKeywords: ['neuvector'] }
  ]

  for (const fam of productKeywordFamilies) {
    const itemMatches = fam.keywords.some(kw => nameLower.includes(kw) || skuUpper.toLowerCase().includes(kw))
    if (itemMatches) {
      const existingLayer = (props.layerItems || []).find(l => {
        const layerText = `${l.layer} ${l.component} ${l.license?.enterprise || ''} ${l.license?.free || ''}`.toLowerCase()
        return fam.layerKeywords.some(kw => layerText.includes(kw))
      })
      if (existingLayer) {
        return {
          isAlreadyInSizing: true,
          isRawInfra,
          matchedLayerText: `Layer ${existingLayer.no} (${existingLayer.layer})`
        }
      }
    }
  }

  return {
    isAlreadyInSizing: false,
    isRawInfra
  }
}

function getCategoryLabel(cat: PricingItemRecord): string {
  const c = (cat.category || '').toLowerCase()
  const s = (cat.sku || '').toUpperCase()
  if (s.includes('APPDYN') || s.includes('DYNATRACE')) return 'Observability & APM'
  if (c === 'license') return 'Lisensi Software'
  if (c === 'database') return 'Database Engine'
  if (c === 'security') return 'Keamanan & Enkripsi'
  if (c === 'backup') return 'Tarif Backup Storage'
  if (c === 'compute') return 'Tarif Komputasi'
  if (c === 'storage') return 'Tarif SAN Storage'
  if (c === 'baremetal') return 'Tarif Server Baremetal'
  return cat.category.toUpperCase()
}

function getCategoryTheme(cat: PricingItemRecord): string {
  const c = (cat.category || '').toLowerCase()
  const s = (cat.sku || '').toUpperCase()
  if (s.includes('APPDYN') || s.includes('DYNATRACE')) return 'theme-apm'
  if (c === 'license') return 'theme-license'
  if (c === 'database') return 'theme-database'
  if (c === 'security') return 'theme-security'
  if (c === 'backup') return 'theme-backup'
  return 'theme-infra'
}

function getItemExplanation(cat: PricingItemRecord): string {
  const sku = (cat.sku || '').toUpperCase()
  if (sku.includes('APPDYN_DB')) {
    return 'Agent pemantauan database mendalam untuk diagnosa SQL query lambat, connection pool, deadlock, dan KPI database secara real-time.'
  }
  if (sku.includes('APPDYN_POD')) {
    return 'Paket lisensi APM per 5 worker pod untuk distributed request tracing dan deteksi bottleneck pada arsitektur microservices container.'
  }
  if (sku.includes('NGINX')) {
    return 'Software enterprise reverse proxy, SSL offloading, rate limiting, dan load balancing web ingress berkecepatan tinggi.'
  }
  if (sku.includes('OCP')) {
    return 'Platform orkestrasi container Kubernetes enterprise dari Red Hat untuk microservices dan DevOps workloads mission-critical.'
  }
  if (sku.includes('KAFKA')) {
    return 'Platform event streaming terdistribusi enterprise untuk pub/sub data pipeline skala jutaan TPS berlatensi rendah.'
  }
  if (sku.includes('REDIS')) {
    return 'Database in-memory enterprise dengan clustering multi-master untuk distributed caching, session store, dan low-latency lookup.'
  }
  if (sku.includes('ORACLE')) {
    return 'RDBMS enterprise mission-critical dengan ketahanan tinggi, dukungan RAC clustering, dan performa transaksi skala perbankan.'
  }
  if (sku.includes('RUBRIK')) {
    return 'Tarif penyimpanan proteksi data dan backup immutable terisolasi per Gigabyte / tahun.'
  }
  if (sku.includes('NBU')) {
    return 'Tarif proteksi data NetBackup per Gigabyte / tahun yang dihitung otomatis berdasarkan kapasitas storage aktif.'
  }
  if (sku.startsWith('VCPU')) {
    return 'Tarif satuan alokasi komputasi vCPU per tahun yang dikalikan otomatis dengan total core mesin virtual.'
  }
  if (sku.startsWith('STG_')) {
    return 'Tarif penyimpanan SAN Enterprise per Gigabyte / tahun berdasarkan tier performa storage yang dipilih.'
  }
  if (sku.startsWith('BM_STD')) {
    return 'Tarif sewa server fisik baremetal standard dedicated per unit per tahun untuk komputasi non-virtualisasi.'
  }
  return `Master katalog: ${cat.name} (${cat.category}).`
}

function getSizingTypeHint(cat: PricingItemRecord): string {
  const sku = (cat.sku || '').toUpperCase()
  const unit = (cat.unit || '').toLowerCase()
  if (sku.includes('APPDYN') || unit.includes('node') || unit.includes('pod') || unit.includes('agent')) {
    return 'Lisensi Add-on Saja (Tanpa Server Baru)'
  }
  if (['compute', 'storage', 'backup', 'baremetal'].includes((cat.category || '').toLowerCase())) {
    return 'Otomatis dihitung via Step 1-3'
  }
  return 'Dapat diset VM / Baremetal / Lisensi'
}

const unmappedCatalogItems = computed<PricingItemRecord[]>(() => {
  return availableCatalogItems.value.filter(item => {
    const info = getCatalogItemSizingInfo(item)
    return !info.isAlreadyInSizing && !info.isRawInfra
  })
})

const licenseCatalogItems = computed<PricingItemRecord[]>(() => {
  return availableCatalogItems.value.filter(item => item.category === 'license')
})

const infraRatesCatalogItems = computed<PricingItemRecord[]>(() => {
  return availableCatalogItems.value.filter(item => {
    return getCatalogItemSizingInfo(item).isRawInfra
  })
})

const filteredModalCatalogItems = computed<PricingItemRecord[]>(() => {
  let list = availableCatalogItems.value
  if (catalogModalFilter.value === 'unmapped') {
    list = unmappedCatalogItems.value
  } else if (catalogModalFilter.value === 'license') {
    list = licenseCatalogItems.value
  } else if (catalogModalFilter.value === 'infra') {
    list = infraRatesCatalogItems.value
  }

  const q = catalogModalSearch.value.trim().toLowerCase()
  if (!q) return list
  return list.filter(item => {
    return (
      item.sku.toLowerCase().includes(q) ||
      item.name.toLowerCase().includes(q) ||
      item.category.toLowerCase().includes(q) ||
      (item.unit || '').toLowerCase().includes(q)
    )
  })
})

// Active preset baseline reference for accurate customization detection
const activePresetBaseline = computed(() => {
  const isCrit = props.tierFilter === 'critical'
  const preset = DEFAULT_PRESET_ITEMS.find(p => p.id === props.activeItem.id)
  if (preset) {
    return isCrit ? preset.critical : preset.moderate
  }
  return null
})

const recommendedNodeCount = computed(() => {
  if (activePresetBaseline.value) return activePresetBaseline.value.nodeCount
  return props.tierFilter === 'critical' ? 2 : 1
})

const recommendedVcpu = computed(() => {
  if (props.activeConfig.computeType === 'none') return 0
  if (activePresetBaseline.value) return activePresetBaseline.value.vcpuPerNode
  const isCrit = props.tierFilter === 'critical'
  if (props.activeItem.no === 2 || (isCrit && props.activeItem.no === 4)) return 64
  if (props.activeItem.no === 8 || props.activeItem.no === 9) return isCrit ? 16 : 8
  if ([1, 3, 6].includes(props.activeItem.no)) return 8
  return isCrit ? 8 : 4
})

const recommendedRam = computed(() => {
  if (props.activeConfig.computeType === 'none') return 0
  if (activePresetBaseline.value) return activePresetBaseline.value.ramGBPerNode
  const isCrit = props.tierFilter === 'critical'
  if (props.activeItem.no === 2 || (isCrit && props.activeItem.no === 4)) return 1024
  if (props.activeItem.no === 8 || props.activeItem.no === 9) return isCrit ? 64 : 32
  if (props.activeItem.no === 3) return 32
  if (props.activeItem.no === 1) return 16
  return isCrit ? 16 : 8
})

const recommendedStorageGB = computed(() => {
  if (props.activeConfig.computeType === 'none') return 0
  if (activePresetBaseline.value) return activePresetBaseline.value.storageGBPerNode
  return 200
})

const recommendedStorageClass = computed(() => {
  if (activePresetBaseline.value) return activePresetBaseline.value.storageClass
  return 'Tier 2'
})

const recommendedComputeType = computed(() => {
  if (props.activeConfig.computeType === 'none' && props.activeItem.id.startsWith('layer_custom_')) return 'none'
  if (activePresetBaseline.value) return activePresetBaseline.value.computeType
  const isCrit = props.tierFilter === 'critical'
  if (props.activeItem.no === 5) return 'shared'
  if (isCrit && (props.activeItem.no === 2 || props.activeItem.no === 4)) return 'baremetal'
  return 'vm'
})

const isComputeCustomized = computed(() => {
  if (props.activeConfig.computeType === 'none') return false
  if (activePresetBaseline.value) {
    return (
      props.activeConfig.vcpuPerNode !== activePresetBaseline.value.vcpuPerNode ||
      props.activeConfig.ramGBPerNode !== activePresetBaseline.value.ramGBPerNode ||
      props.activeConfig.nodeCount !== activePresetBaseline.value.nodeCount
    )
  }
  return (
    props.activeConfig.vcpuPerNode !== recommendedVcpu.value ||
    props.activeConfig.ramGBPerNode !== recommendedRam.value
  )
})

const isStorageCustomized = computed(() => {
  if (props.activeConfig.computeType === 'none') return false
  if (activePresetBaseline.value) {
    return (
      props.activeConfig.storageGBPerNode !== activePresetBaseline.value.storageGBPerNode ||
      props.activeConfig.storageClass !== activePresetBaseline.value.storageClass ||
      props.activeConfig.backupEnabled !== activePresetBaseline.value.backupEnabled
    )
  }
  return false
})

const isLicenseCustomized = computed(() => {
  const expectedLic = props.licenseFilter === 'license' ? 'enterprise' : 'free'
  return props.activeConfig.licenseKind !== expectedLic || Boolean(props.activeConfig.pricingSku)
})

// Clean baseline summary without duplicate "Node Node"
const cleanBaselineSummary = computed(() => {
  const isCrit = props.tierFilter === 'critical'
  const spec = isCrit ? props.activeItem.critical : props.activeItem.moderate
  const nodesText = spec.nodes.replace(/\bNode\s+Node\b/gi, 'Node')
  return `${nodesText} · ${spec.spec} · ${spec.type}`
})

// ========================================================
// PERSISTENT COLLAPSIBLE STEPS (localStorage)
// ========================================================
const defaultCollapsedSteps: Record<string, boolean> = {
  step1: false,
  step2: false,
  step3: false,
  step4: false,
}

const collapsedSteps = ref<Record<string, boolean>>({ ...defaultCollapsedSteps })
const COLLAPSED_STORAGE_KEY = 'mandiri_sizing_collapsed_steps'

onMounted(() => {
  if (import.meta.client) {
    try {
      const saved = localStorage.getItem(COLLAPSED_STORAGE_KEY)
      if (saved) {
        const parsed = JSON.parse(saved)
        if (typeof parsed === 'object' && parsed !== null) {
          collapsedSteps.value = {
            step1: Boolean(parsed.step1),
            step2: Boolean(parsed.step2),
            step3: Boolean(parsed.step3),
            step4: Boolean(parsed.step4),
          }
        }
      }
    } catch (e) {
      console.warn('Failed to load collapsed steps from localStorage:', e)
    }
  }
})

function setStepCollapsed(stepKey: string, val: boolean) {
  collapsedSteps.value[stepKey] = val
  if (import.meta.client) {
    try {
      localStorage.setItem(COLLAPSED_STORAGE_KEY, JSON.stringify(collapsedSteps.value))
    } catch (e) {
      console.warn('Failed to save collapsed steps to localStorage:', e)
    }
  }
}

const isAllCollapsed = computed(() => {
  return Object.values(collapsedSteps.value).every(val => val === true)
})

function toggleAllSteps() {
  const nextVal = !isAllCollapsed.value
  const updated: Record<string, boolean> = {
    step1: nextVal,
    step2: nextVal,
    step3: nextVal,
    step4: nextVal,
  }
  collapsedSteps.value = updated
  if (import.meta.client) {
    try {
      localStorage.setItem(COLLAPSED_STORAGE_KEY, JSON.stringify(updated))
    } catch (e) {
      console.warn('Failed to save collapsed steps to localStorage:', e)
    }
  }
}

// Collapsed Step Summaries (Clean & Concise)
const step1Summary = computed(() => {
  if (props.activeConfig.computeType === 'none') {
    return `${props.activeConfig.nodeCount} Unit Lisensi`
  }
  if (props.activeConfig.computeType === 'shared') {
    return `${props.activeConfig.nodeCount} Unit Shared SAN (0 vCPU)`
  }
  const vcpu = props.activeConfig.vcpuPerNode
  const ram = props.activeConfig.ramGBPerNode
  return `${props.activeConfig.nodeCount} Server · ${vcpu} vCPU / ${ram} GB RAM`
})

const step2Summary = computed(() => {
  if (props.activeConfig.computeType === 'none') {
    return '0 GB (Menumpang server eksisting)'
  }
  const cls = props.activeConfig.storageClass || 'Tier 2'
  const backup = props.activeConfig.backupEnabled ? 'Backup Harian' : 'Tanpa Backup'
  return `${props.activeConfig.storageGBPerNode} GB/server (${cls}) · ${backup}`
})

const step3Summary = computed(() => {
  const map: Record<string, string> = {
    vm: 'Virtual Machine (VM)',
    baremetal: 'Server Fisik Khusus (Baremetal)',
    shared: 'Penyimpanan Berbagi (Shared SAN)',
    none: 'Lisensi / Add-on Saja'
  }
  return map[props.activeConfig.computeType] || 'Virtual Machine (VM)'
})

const step4Summary = computed(() => {
  if (props.activeConfig.licenseKind === 'free') {
    return 'Gratis (Community OSS)'
  }
  const prodName = props.activeConfig.licenseProductName || props.activeLayerCost?.license_product_name || props.activeItem.license?.enterprise
  return prodName ? `Enterprise · ${prodName}` : 'Lisensi Resmi Enterprise'
})
</script>

<style scoped>
.sketch-workspace-grid {
  flex: 1;
  display: grid;
  grid-template-columns: 325px 1fr;
  overflow: hidden;
  min-height: 0;
  background: #F8FAFC;
}

.sketch-list-section {
  background: #FFFFFF;
  border-right: 1px solid #E2E8F0;
  display: flex;
  flex-direction: column;
  overflow-y: auto;
}

.section-hdr {
  padding: 14px 16px 12px 16px;
  border-bottom: 1px solid #E2E8F0;
  background: #FAFCFF;
  display: flex;
  flex-direction: column;
  gap: 10px;
  position: sticky;
  top: 0;
  z-index: 5;
}

.section-hdr-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
}

.section-hdr-title-wrap {
  display: flex;
  align-items: center;
  gap: 7px;
}

.hdr-accent-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: #002060;
  flex-shrink: 0;
}

.section-title {
  font-size: 11.5px;
  font-weight: 800;
  text-transform: uppercase;
  letter-spacing: 0.04em;
  color: #002060;
  margin: 0;
  white-space: nowrap;
}

.item-count {
  font-size: 11px;
  font-weight: 700;
  color: #475569;
  background: #EDF2F7;
  padding: 2px 8px;
  border-radius: 12px;
  border: 1px solid #E2E8F0;
  white-space: nowrap;
}

.bullet-list-container {
  display: flex;
  flex-direction: column;
  padding: 6px 8px;
  gap: 3px;
}

.bullet-row {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
  cursor: pointer;
  border-radius: 8px;
  border: 1px solid transparent;
  transition: all 0.15s ease;
  user-select: none;
}

.bullet-row:hover {
  background: #F8FAFC;
  border-color: #E2E8F0;
}

.bullet-row.active {
  background: #EFF6FF;
  border-color: #BFDBFE;
  box-shadow: 0 1px 2px rgba(0, 32, 96, 0.05);
}

.layer-num-badge {
  width: 24px;
  height: 24px;
  border-radius: 6px;
  background: #E2E8F0;
  color: #334155;
  font-size: 11.5px;
  font-weight: 800;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  font-variant-numeric: tabular-nums;
  transition: all 0.15s ease;
}

.bullet-row.active .layer-num-badge {
  background: #002060;
  color: #FFFFFF;
}

.bullet-row.customized:not(.active) .layer-num-badge {
  background: #FEF3C7;
  color: #B45309;
}

.bullet-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
  gap: 2px;
}

.bullet-title-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 6px;
}

.bullet-label {
  font-size: 12.5px;
  font-weight: 700;
  color: #1E293B;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.bullet-row.active .bullet-label {
  color: #002060;
  font-weight: 800;
}

.custom-badge {
  font-size: 11px;
  font-weight: 700;
  color: #B45309;
  background: #FEF3C7;
  padding: 1px 5px;
  border-radius: 4px;
  text-transform: uppercase;
}

.bullet-sub-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 6px;
}

.bullet-sub {
  font-size: 11px;
  color: #475569;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.bullet-cost-pill {
  font-size: 11px;
  font-weight: 700;
  color: #002060;
  background: #F8FAFC;
  border: 1px solid #E2E8F0;
  padding: 1.5px 6px;
  border-radius: 4px;
  white-space: nowrap;
  flex-shrink: 0;
  font-variant-numeric: tabular-nums;
}

.bullet-row.active .bullet-cost-pill {
  background: #FFFFFF;
  border-color: #BFDBFE;
}

.sketch-card-section {
  padding: 20px 24px 110px 24px;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
}

.pure-content-card {
  background: #FFFFFF;
  border: 1px solid #E2E8F0;
  border-radius: 12px;
  padding: 22px 24px;
  box-shadow: 0 2px 6px rgba(15, 23, 42, 0.04);
  display: flex;
  flex-direction: column;
  gap: 16px;
}

/* 1. Component Header Banner */
.component-header-banner {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  padding-bottom: 14px;
  border-bottom: 1px solid #E2E8F0;
  gap: 16px;
  flex-wrap: wrap;
}

.badge-row {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 6px;
  flex-wrap: wrap;
}

.badge-layer {
  font-size: 11px;
  font-weight: 800;
  color: #002060;
  background: #EFF6FF;
  border: 1px solid #DBEAFE;
  padding: 2.5px 8px;
  border-radius: 5px;
  text-transform: uppercase;
  letter-spacing: 0.02em;
}

.custom-status-tag {
  font-size: 11px;
  font-weight: 700;
  color: #B45309;
  background: #FEF3C7;
  border: 1px solid #FDE68A;
  padding: 2px 8px;
  border-radius: 5px;
}

.standard-status-tag {
  font-size: 11px;
  font-weight: 700;
  color: #15803D;
  background: #DCFCE7;
  border: 1px solid #BBF7D0;
  padding: 2px 8px;
  border-radius: 5px;
}

.component-title-row {
  display: flex;
  align-items: center;
  gap: 8px;
}

.component-emoji {
  font-size: 20px;
  line-height: 1;
}

.component-title {
  font-size: 20px;
  font-weight: 800;
  color: #0F172A;
  margin: 0;
  letter-spacing: -0.01em;
}

/* Unified Layer Action Toolbar (Clean, Compact, Rapi) */
.layer-action-toolbar {
  display: inline-flex;
  align-items: center;
  background: #FFFFFF;
  border: 1px solid #CBD5E1;
  border-radius: 8px;
  padding: 3px 5px;
  box-shadow: 0 1px 2px rgba(15, 23, 42, 0.04);
  gap: 2px;
}

.toolbar-action-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 11.5px;
  font-weight: 600;
  color: #475569;
  background: transparent;
  border: none;
  border-radius: 6px;
  padding: 5px 10px;
  cursor: pointer;
  white-space: nowrap;
  transition: all 0.15s ease;
  line-height: 1;
}

.toolbar-action-btn:hover {
  color: #002060;
  background: #F1F5F9;
}

.toolbar-action-btn.btn-reset.is-customized {
  color: #B45309;
  background: #FFFBEB;
}

.toolbar-action-btn.btn-reset.is-customized:hover {
  background: #FEF3C7;
  color: #92400E;
}

.custom-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: #F59E0B;
  box-shadow: 0 0 0 2px rgba(245, 158, 11, 0.2);
}

.toolbar-action-btn.btn-delete {
  color: #64748B;
}

.toolbar-action-btn.btn-delete:hover {
  color: #DC2626;
  background: #FEF2F2;
}

.toolbar-separator {
  width: 1px;
  height: 16px;
  background: #E2E8F0;
  margin: 0 3px;
}

/* 2. EA Architecture Baseline Banner (Simple, Clean, Rapi) */
.ea-context-banner {
  background: #F8FAFC;
  border: 1px solid #E2E8F0;
  border-radius: 8px;
  padding: 12px 16px;
}

.ea-context-left {
  display: flex;
  align-items: flex-start;
  gap: 12px;
}

.ea-icon-badge {
  width: 28px;
  height: 28px;
  border-radius: 6px;
  background: #EFF6FF;
  border: 1px solid #DBEAFE;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  margin-top: 1px;
}

.ea-context-body {
  display: flex;
  flex-direction: column;
  gap: 7px;
  flex: 1;
  min-width: 0;
}

.ea-desc-row {
  display: flex;
  align-items: baseline;
  gap: 6px;
  flex-wrap: wrap;
}

.ea-desc-lead {
  font-size: 12.5px;
  font-weight: 700;
  color: #002060;
}

.ea-desc-text {
  font-size: 12.5px;
  color: #334155;
  line-height: 1.45;
}

.ea-meta-row {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.ea-baseline-chip {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  font-size: 11px;
  font-weight: 700;
  color: #92400E;
  background: #FEF3C7;
  border: 1px solid #FDE68A;
  padding: 2px 8px;
  border-radius: 4px;
}

.ea-baseline-spec {
  font-size: 12px;
  font-weight: 600;
  color: #0F172A;
  letter-spacing: -0.01em;
}

.ea-contract-pill {
  font-size: 11px;
  font-weight: 600;
  color: #0369A1;
  background: #E0F2FE;
  border: 1px solid #BAE6FD;
  padding: 2px 7px;
  border-radius: 4px;
}

/* ─────────────────────────────────────────────
   Step Section Form Elements & Quick Chips
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
  gap: 2px;
}

.field-title-lbl {
  font-size: 11.5px;
  font-weight: 800;
  color: #1E293B;
  text-transform: uppercase;
  letter-spacing: 0.02em;
}

.field-sub-hint {
  font-size: 11.5px;
  color: #475569;
}

.server-count-row,
.storage-stepper-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  padding-top: 12px;
  border-top: 1px dashed #E2E8F0;
  flex-wrap: wrap;
}

.quick-stepper-group {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.quick-preset-chips {
  display: flex;
  align-items: center;
  gap: 5px;
  flex-wrap: wrap;
}

.quick-chip-btn {
  background: #FFFFFF;
  border: 1px solid #CBD5E1;
  color: #334155;
  font-size: 11px;
  font-weight: 700;
  padding: 6px 10px;
  border-radius: 6px;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  gap: 4px;
  transition: all 0.14s ease;
}

.quick-chip-btn:hover {
  border-color: #64748B;
  background: #F8FAFC;
}

.quick-chip-btn.recommended:not(.active) {
  border-color: #FCD34D;
  background: #FFFBEB;
}

.quick-chip-btn.active {
  background: #002060;
  border-color: #002060;
  color: #FFFFFF;
}

.chip-star {
  color: #92400E;
  font-size: 11px;
}

.quick-chip-btn.active .chip-star {
  color: #FCD34D;
}

.stepper-wrap-contained {
  width: 170px;
}

.advanced-steppers-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 14px;
}

.storage-options-grid {
  display: grid;
  grid-template-columns: 1.25fr 1fr;
  gap: 14px;
  padding-top: 12px;
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
  border: 1.5px solid #CBD5E1;
  border-radius: 9px;
  padding: 13px 14px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  gap: 8px;
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
  background: #F0F6FF;
  border-color: #002060;
  box-shadow: 0 0 0 2px rgba(0, 32, 96, 0.12);
}

.backup-top-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.backup-left-title {
  display: flex;
  align-items: center;
  gap: 8px;
}

.custom-chk-large {
  width: 16px;
  height: 16px;
  cursor: pointer;
  accent-color: #002060;
  flex-shrink: 0;
}

.backup-main-title {
  font-size: 12.5px;
  font-weight: 800;
  color: #0F172A;
}

.backup-status-pill {
  font-size: 11px;
  font-weight: 800;
  padding: 2px 7px;
  border-radius: 12px;
  text-transform: uppercase;
}

.backup-status-pill.on {
  background: #DCFCE7;
  color: #15803D;
}

.backup-status-pill.off {
  background: #E2E8F0;
  color: #475569;
}

.backup-desc {
  font-size: 11.5px;
  color: #475569;
  line-height: 1.45;
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
  border: 1.5px solid #CBD5E1;
  border-radius: 9px;
  padding: 14px;
  cursor: pointer;
  transition: all 0.15s ease;
  display: flex;
  flex-direction: column;
  gap: 8px;
  text-align: left;
}

.license-choice-card:hover {
  border-color: #94A3B8;
  background: #F8FAFC;
}

.license-choice-card.active {
  background: #F0F6FF;
  border-color: #002060;
  box-shadow: 0 0 0 2px rgba(0, 32, 96, 0.12);
}

.license-choice-card.enterprise-card.active {
  background: #F0FDF4;
  border-color: #15803D;
  box-shadow: 0 0 0 2px rgba(21, 128, 61, 0.12);
}

.lic-card-hdr {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 10px;
}

.lic-card-hdr-left {
  display: flex;
  align-items: center;
  gap: 10px;
}

.lic-icon {
  font-size: 20px;
}

.lic-hdr-titles {
  display: flex;
  flex-direction: column;
}

.lic-name {
  font-size: 12.5px;
  font-weight: 800;
  color: #0F172A;
}

.lic-price-tag {
  font-size: 11px;
  font-weight: 700;
  color: #475569;
}

.font-emerald {
  color: #15803D;
}

.lic-radio-indicator {
  width: 18px;
  height: 18px;
  border-radius: 50%;
  border: 1.5px solid #CBD5E1;
  background: #FFFFFF;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  color: #FFFFFF;
  flex-shrink: 0;
}

.lic-radio-indicator.checked {
  background: #002060;
  border-color: #002060;
}

.lic-radio-indicator.emerald.checked {
  background: #15803D;
  border-color: #15803D;
}

.lic-product-name {
  margin: 0;
  font-size: 12px;
  font-weight: 700;
  color: #002060;
  background: transparent;
  padding: 0;
  border: none;
}

.lic-card-desc {
  margin: 0;
  font-size: 11px;
  color: #475569;
  line-height: 1.45;
}

/* ─────────────────────────────────────────────
   Layer Cost Breakdown Card + Visual Proportion Bar
   ───────────────────────────────────────────── */
.layer-cost-card {
  background: #FFFFFF;
  border: 1px solid #E2E8F0;
  border-radius: 10px;
  padding: 16px 18px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.cost-card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  flex-wrap: wrap;
}

.cost-card-title-wrap {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.cost-card-title {
  font-size: 13.5px;
  font-weight: 800;
  color: #002060;
}

.cost-card-sub {
  font-size: 11.5px;
  color: #475569;
}

.cost-card-total-badge {
  display: flex;
  align-items: baseline;
  gap: 8px;
  background: transparent;
  border: none;
  padding: 0;
}

.total-lbl {
  font-size: 11.5px;
  font-weight: 700;
  color: #475569;
}

.layer-grand-price {
  font-size: 18px;
  font-weight: 800;
  color: #002060;
}

.cost-proportion-bar {
  height: 7px;
  width: 100%;
  background: #E2E8F0;
  border-radius: 99px;
  overflow: hidden;
  display: flex;
  gap: 2px;
}

.prop-segment {
  height: 100%;
  transition: width 0.25s ease;
}

.seg-compute { background: #002060; }
.seg-storage { background: #0284C7; }
.seg-backup { background: #D97706; }
.seg-license { background: #16A34A; }

.cost-metrics-row {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 12px;
}

.cost-metric-box {
  background: transparent;
  border: none;
  border-right: 1px solid #F1F5F9;
  border-radius: 0;
  padding: 2px 10px 2px 0;
  display: flex;
  flex-direction: column;
  gap: 3px;
}

.cost-metric-box:last-child {
  border-right: none;
}

.metric-hdr {
  display: flex;
  align-items: center;
  gap: 6px;
}

.metric-dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  flex-shrink: 0;
}

.dot-compute { background: #002060; }
.dot-storage { background: #0284C7; }
.dot-backup { background: #D97706; }
.dot-license { background: #16A34A; }

.metric-title {
  font-size: 11px;
  font-weight: 700;
  color: #334155;
  text-transform: uppercase;
}

.metric-sub {
  font-size: 11.5px;
  color: #475569;
  margin-top: 2px;
}

.text-truncate {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.free-badge {
  font-size: 12.5px;
  font-weight: 800;
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
  border-radius: 6px;
}

/* ─────────────────────────────────────────────
   Guided Wizard Prev / Next Navigation Footer
   ───────────────────────────────────────────── */
.layer-wizard-nav {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding-top: 12px;
  border-top: 1px solid #E2E8F0;
}

.wizard-nav-btn {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  background: #FFFFFF;
  border: 1px solid #CBD5E1;
  border-radius: 8px;
  padding: 9px 14px;
  cursor: pointer;
  transition: all 0.15s ease;
  text-align: left;
  min-width: 210px;
}

.wizard-nav-btn.next {
  justify-content: flex-end;
  background: #F8FAFC;
  border-color: #93C5FD;
}

.wizard-nav-btn:hover:not(:disabled) {
  border-color: #002060;
  background: #EFF6FF;
  transform: translateY(-1px);
}

.wizard-nav-btn:disabled {
  opacity: 0.45;
  cursor: not-allowed;
}

.nav-arrow {
  font-size: 15px;
  font-weight: 800;
  color: #002060;
}

.nav-text-wrap {
  display: flex;
  flex-direction: column;
  gap: 1px;
}

.nav-text-wrap.align-right {
  text-align: right;
}

.nav-hint {
  font-size: 11px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.02em;
  color: #475569;
}

.nav-name {
  font-size: 12px;
  font-weight: 800;
  color: #002060;
}

.wizard-progress-indicator {
  font-size: 11.5px;
  color: #64748B;
  background: #F1F5F9;
  padding: 5px 12px;
  border-radius: 20px;
}


.btn-add-catalog-layer {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  width: 100%;
  background: #F0F6FF;
  border: 1px dashed #93C5FD;
  color: #002060;
  font-size: 11.5px;
  font-weight: 700;
  padding: 8px 12px;
  border-radius: 7px;
  cursor: pointer;
  transition: all 0.18s ease;
  box-shadow: 0 1px 2px rgba(0, 32, 96, 0.03);
}

.btn-icon-circle {
  width: 18px;
  height: 18px;
  border-radius: 50%;
  background: #DBEAFE;
  color: #1D4ED8;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  transition: all 0.18s ease;
}

.btn-label-text {
  letter-spacing: -0.01em;
}

.btn-add-catalog-layer:hover {
  background: #002060;
  border: 1px solid #002060;
  color: #FFFFFF;
  box-shadow: 0 3px 8px rgba(0, 32, 96, 0.18);
  transform: translateY(-1px);
}

.btn-add-catalog-layer:hover .btn-icon-circle {
  background: rgba(255, 255, 255, 0.2);
  color: #FFFFFF;
}

.bullet-badges {
  display: flex;
  align-items: center;
  gap: 4px;
}

.catalog-tag {
  font-size: 9px;
  font-weight: 700;
  color: #0284C7;
  background: #E0F2FE;
  padding: 1px 5px;
  border-radius: 3px;
  text-transform: uppercase;
}

.list-hdr-actions {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.btn-restore-defaults {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  width: 100%;
  padding: 6px 12px;
  background: #EFF6FF;
  border: 1px dashed #93C5FD;
  border-radius: 6px;
  color: #1D4ED8;
  font-size: 11px;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.15s ease;
}

.btn-restore-defaults:hover {
  background: #DBEAFE;
  border-color: #3B82F6;
  color: #1E40AF;
}

.btn-delete-row {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 22px;
  height: 22px;
  border-radius: 4px;
  border: none;
  background: transparent;
  color: #94A3B8;
  cursor: pointer;
  transition: all 0.15s ease;
  margin-left: 4px;
  opacity: 0;
}

.bullet-row:hover .btn-delete-row,
.bullet-row.active .btn-delete-row {
  opacity: 1;
}

.btn-delete-row:hover {
  background: #FEE2E2;
  color: #DC2626;
}

.btn-delete-active-layer {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  font-size: 11px;
  font-weight: 600;
  color: #DC2626;
  background: #FEF2F2;
  border: 1px solid #FECACA;
  border-radius: 4px;
  padding: 5px 10px;
  cursor: pointer;
  transition: all 0.15s ease;
  white-space: nowrap;
}

.btn-delete-active-layer:hover {
  background: #DC2626;
  color: #FFFFFF;
  border-color: #DC2626;
}

.badge-sku {
  font-size: 10px;
  font-weight: 800;
  color: #1D4ED8;
  background: #EFF6FF;
  border: 1px solid #DBEAFE;
  padding: 2px 7px;
  border-radius: 4px;
  font-family: 'JetBrains Mono', monospace;
}

/* License only banner notice */
.license-only-notice {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  background: #F0F9FF;
  border: 1px solid #BAE6FD;
  border-radius: 8px;
  padding: 12px 14px;
}

.notice-icon {
  font-size: 20px;
  line-height: 1;
}

.notice-text {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.notice-title {
  font-size: 13px;
  font-weight: 700;
  color: #0369A1;
}

.notice-desc {
  font-size: 11.5px;
  color: #0284C7;
  line-height: 1.4;
}

/* Enterprise SKU picker box */
.enterprise-sku-picker-box {
  margin-top: 8px;
  padding: 14px 16px;
  background: #F8FAFC;
  border: 1px solid #E2E8F0;
  border-radius: 8px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.sku-picker-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
}

.sku-picker-title-wrap {
  display: flex;
  flex-direction: column;
  gap: 1px;
}

.sku-picker-lbl {
  font-size: 12px;
  font-weight: 700;
  color: #1E293B;
}

.sku-picker-sub {
  font-size: 11px;
  color: #64748B;
}

.sku-selected-tag {
  font-size: 10px;
  font-weight: 700;
  color: #047857;
  background: #D1FAE5;
  padding: 2px 8px;
  border-radius: 4px;
  font-family: 'JetBrains Mono', monospace;
}

.sku-picker-control-row {
  display: flex;
  width: 100%;
}

.sku-dropdown-select {
  width: 100%;
  padding: 8px 12px;
  font-size: 12px;
  font-weight: 600;
  color: #1E293B;
  background: #FFFFFF;
  border: 1px solid #CBD5E1;
  border-radius: 6px;
  outline: none;
  cursor: pointer;
  transition: border-color 0.15s ease;
}

.sku-dropdown-select:focus {
  border-color: #002060;
  box-shadow: 0 0 0 2px rgba(0, 32, 96, 0.1);
}

.active-sku-detail-card {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: #FFFFFF;
  border: 1px solid #E2E8F0;
  border-radius: 6px;
  padding: 10px 14px;
  gap: 12px;
}

.sku-card-main {
  display: flex;
  flex-direction: column;
  gap: 3px;
}

.sku-card-name {
  font-size: 12px;
  font-weight: 700;
  color: #0F172A;
}

.sku-card-env-badge {
  font-size: 10px;
  color: #64748B;
}

.sku-card-price-wrap {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 2px;
}

.sku-price-title {
  font-size: 10px;
  color: #64748B;
  text-transform: uppercase;
}

.sku-price-val {
  font-size: 13px;
  font-weight: 800;
  color: #059669;
}

.sku-price-unit {
  font-size: 11px;
  font-weight: 600;
  color: #64748B;
}

.sku-price-pending {
  font-size: 11px;
  font-weight: 700;
  color: #D97706;
}

/* ===================================================
   MODAL: KATALOG HARGA (ELEGANT ENTERPRISE AESTHETICS)
   =================================================== */
.catalog-modal-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(15, 23, 42, 0.55);
  backdrop-filter: blur(5px);
  z-index: 1000;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 10vh 10vw; /* 20vh vertical & 20vw horizontal total padding from viewport */
  box-sizing: border-box;
  animation: fadeInModal 0.2s cubic-bezier(0.16, 1, 0.3, 1);
}

@keyframes fadeInModal {
  from { opacity: 0; transform: scale(0.98); }
  to { opacity: 1; transform: scale(1); }
}

.catalog-modal-card {
  background: #FFFFFF;
  border-radius: 14px;
  width: 100%;
  height: 100%;
  max-width: 1560px;
  display: flex;
  flex-direction: column;
  box-shadow: 0 25px 50px -12px rgba(15, 23, 42, 0.25), 0 0 0 1px rgba(15, 23, 42, 0.08);
  overflow: hidden;
}

.modal-hdr {
  padding: 18px 28px;
  border-bottom: 1px solid #E2E8F0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: #FFFFFF;
  flex-shrink: 0;
}

.modal-title-wrap {
  display: flex;
  align-items: center;
  gap: 14px;
}

.modal-icon-badge {
  width: 42px;
  height: 42px;
  border-radius: 10px;
  background: #EFF6FF;
  border: 1px solid #BFDBFE;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.modal-title {
  font-size: 16px;
  font-weight: 800;
  color: #002060;
  margin: 0;
  letter-spacing: -0.01em;
}

.modal-subtitle {
  font-size: 12.5px;
  color: #64748B;
  margin: 3px 0 0 0;
  line-height: 1.4;
}

.modal-btn-close {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #F8FAFC;
  border: 1px solid #E2E8F0;
  font-size: 14px;
  color: #64748B;
  cursor: pointer;
  border-radius: 8px;
  transition: all 0.15s ease;
}

.modal-btn-close:hover {
  color: #0F172A;
  background: #F1F5F9;
  border-color: #CBD5E1;
}

.modal-search-strip {
  padding: 14px 28px;
  background: #F8FAFC;
  border-bottom: 1px solid #E2E8F0;
  display: flex;
  flex-direction: column;
  gap: 12px;
  flex-shrink: 0;
}

.modal-search-wrap {
  display: flex;
  align-items: center;
  gap: 10px;
  background: #FFFFFF;
  border: 1px solid #CBD5E1;
  border-radius: 8px;
  padding: 9px 14px;
  box-shadow: 0 1px 2px rgba(15, 23, 42, 0.04);
  transition: all 0.15s ease;
}

.modal-search-wrap:focus-within {
  border-color: #002060;
  box-shadow: 0 0 0 3px rgba(0, 32, 96, 0.08);
}

.modal-search-input {
  flex: 1;
  border: none;
  outline: none;
  font-size: 13px;
  color: #0F172A;
  background: transparent;
}

.modal-search-clear {
  border: none;
  background: none;
  color: #94A3B8;
  cursor: pointer;
  font-size: 12px;
  padding: 2px 6px;
  border-radius: 4px;
}

.modal-search-clear:hover {
  background: #F1F5F9;
  color: #334155;
}

.modal-filter-pills {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.filter-pill {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 11.5px;
  font-weight: 600;
  padding: 5px 12px;
  border-radius: 20px;
  border: 1px solid #E2E8F0;
  background: #FFFFFF;
  color: #475569;
  cursor: pointer;
  transition: all 0.15s ease;
}

.filter-pill:hover {
  background: #F1F5F9;
  color: #0F172A;
}

.filter-pill.active {
  background: #002060;
  color: #FFFFFF;
  border-color: #002060;
  box-shadow: 0 2px 4px rgba(0, 32, 96, 0.2);
}

.pill-count {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 1px 6px;
  border-radius: 10px;
  font-size: 10.5px;
  font-weight: 700;
  background: #F1F5F9;
  color: #475569;
}

.filter-pill.active .pill-count {
  background: rgba(255, 255, 255, 0.2);
  color: #FFFFFF;
}

.pill-glow-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: #10B981;
}

.modal-items-scroll {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
  padding: 20px 28px;
  display: flex;
  flex-direction: column;
  gap: 14px;
  background: #F8FAFC;
}

.empty-catalog-state {
  text-align: center;
  padding: 40px 20px;
  color: #64748B;
}

.empty-icon {
  font-size: 32px;
  display: block;
  margin-bottom: 8px;
}

.empty-text {
  font-size: 13px;
  margin: 0;
  line-height: 1.5;
}

.catalog-picker-card {
  display: flex;
  align-items: flex-start;
  gap: 16px;
  border: 1px solid #E2E8F0;
  border-radius: 10px;
  padding: 16px;
  background: #FFFFFF;
  transition: all 0.18s cubic-bezier(0.16, 1, 0.3, 1);
  box-shadow: 0 1px 3px rgba(15, 23, 42, 0.03);
}

.catalog-picker-card:hover {
  border-color: #93C5FD;
  box-shadow: 0 6px 16px -2px rgba(0, 32, 96, 0.08);
  transform: translateY(-1px);
}

.catalog-picker-card.card-installed {
  background: #FAFCFE;
  border-color: #E2E8F0;
}

.catalog-picker-card.card-infra {
  background: #F8FAFC;
  opacity: 0.9;
}

/* Category Avatars */
.card-icon-avatar {
  width: 42px;
  height: 42px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  transition: transform 0.15s ease;
}

.catalog-picker-card:hover .card-icon-avatar {
  transform: scale(1.05);
}

.card-icon-avatar.theme-apm {
  background: var(--pill-neutral-bg, #ebedee);
  color: var(--accent-teal, #1b6c80);
  border: 1px solid var(--border, #e2e7eb);
}

.card-icon-avatar.theme-database {
  background: var(--pill-positive-bg, #dcfae5);
  color: var(--positive, #1e7b4f);
  border: 1px solid rgba(30, 123, 79, 0.2);
}

.card-icon-avatar.theme-security {
  background: var(--pill-negative-bg, #ffedea);
  color: var(--negative, #c4242a);
  border: 1px solid rgba(196, 36, 42, 0.2);
}

.card-icon-avatar.theme-backup {
  background: var(--pill-warning-bg, #fffbeb);
  color: var(--warning, #92400e);
  border: 1px solid rgba(146, 64, 14, 0.2);
}

.card-icon-avatar.theme-license {
  background: var(--pill-positive-bg, #dcfae5);
  color: var(--positive, #1e7b4f);
  border: 1px solid rgba(30, 123, 79, 0.2);
}

.card-icon-avatar.theme-infra {
  background: var(--raised, #eef2f5);
  color: var(--muted, #5d6165);
  border: 1px solid var(--border, #e2e7eb);
}

/* Card Body */
.card-main-info {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
}

.card-meta-tags {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-wrap: wrap;
  margin-bottom: 5px;
}

.sku-pill {
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 11px;
  font-weight: 700;
  color: #0F172A;
  background: #F1F5F9;
  border: 1px solid #CBD5E1;
  padding: 1px 7px;
  border-radius: 4px;
}

.category-pill {
  font-size: 10px;
  font-weight: 700;
  color: #0369A1;
  background: #E0F2FE;
  padding: 1px 7px;
  border-radius: 4px;
  text-transform: uppercase;
  letter-spacing: 0.02em;
}

.env-pill {
  font-size: 10px;
  font-weight: 600;
  padding: 1px 6px;
  border-radius: 4px;
}

.env-pill.env-prod {
  background: #FEF3C7;
  color: #92400E;
}

.env-pill.env-dev {
  background: #F1F5F9;
  color: #475569;
}

.status-pill {
  font-size: 10px;
  font-weight: 700;
  padding: 1px 7px;
  border-radius: 4px;
}

.status-pill.ready {
  background: #D1FAE5;
  color: #065F46;
}

.status-pill.installed {
  background: #EFF6FF;
  color: #1E40AF;
  border: 1px solid #BFDBFE;
}

.status-pill.infra {
  background: #F1F5F9;
  color: #64748B;
}

.card-product-title {
  font-size: 14px;
  font-weight: 700;
  color: #0F172A;
  margin: 0 0 5px 0;
  line-height: 1.35;
}

.card-description-text {
  font-size: 12px;
  color: #475569;
  line-height: 1.45;
  margin: 0 0 10px 0;
}

.card-feature-hints {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 11px;
}

.feature-item {
  display: flex;
  align-items: center;
  gap: 5px;
}

.feature-label {
  color: #64748B;
}

.feature-val {
  color: #1E293B;
  font-weight: 600;
}

.feature-val.highlight {
  color: #002060;
  background: #F0F6FF;
  padding: 1px 6px;
  border-radius: 4px;
  border: 1px solid #DBEAFE;
}

.feature-divider {
  width: 1px;
  height: 12px;
  background: #CBD5E1;
}

/* Card Action Column */
.card-action-block {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  justify-content: space-between;
  gap: 12px;
  min-width: 175px;
  flex-shrink: 0;
  padding-left: 14px;
  border-left: 1px solid #F1F5F9;
}

.price-box {
  text-align: right;
  display: flex;
  flex-direction: column;
  align-items: flex-end;
}

.price-header {
  font-size: 10px;
  font-weight: 700;
  color: #64748B;
  text-transform: uppercase;
  letter-spacing: 0.04em;
}

.price-val-row {
  display: flex;
  align-items: baseline;
  gap: 3px;
  margin-top: 2px;
}

.price-curr {
  font-size: 11px;
  font-weight: 700;
  color: #059669;
}

.price-num {
  font-size: 16px;
  font-weight: 800;
  color: #002060;
  letter-spacing: -0.02em;
}

.price-pending-tag {
  font-size: 11px;
  font-weight: 700;
  color: #D97706;
  background: #FEF3C7;
  padding: 2px 7px;
  border-radius: 4px;
  margin-top: 3px;
}

.price-period {
  font-size: 10.5px;
  color: #64748B;
  margin-top: 1px;
}

.btn-card-select {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  background: #002060;
  color: #FFFFFF;
  border: none;
  font-size: 11.5px;
  font-weight: 700;
  padding: 8px 14px;
  border-radius: 6px;
  cursor: pointer;
  white-space: nowrap;
  transition: all 0.15s ease;
  box-shadow: 0 1px 3px rgba(0, 32, 96, 0.2);
}

.btn-card-select:hover {
  background: #003399;
  transform: translateY(-1px);
  box-shadow: 0 4px 10px rgba(0, 32, 96, 0.25);
}

.btn-card-select.btn-installed-action {
  background: #F8FAFC;
  color: #334155;
  border: 1px solid #CBD5E1;
  box-shadow: none;
}

.btn-card-select.btn-installed-action:hover {
  background: #002060;
  color: #FFFFFF;
  border-color: #002060;
}

.infra-auto-hint {
  font-size: 11px;
  color: #94A3B8;
  font-style: italic;
  text-align: right;
}

.modal-ftr {
  padding: 14px 28px;
  border-top: 1px solid #E2E8F0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: #FFFFFF;
  flex-shrink: 0;
}

.modal-ftr-count {
  font-size: 12px;
  color: #64748B;
}

.btn-modal-cancel {
  background: #FFFFFF;
  border: 1px solid #CBD5E1;
  color: #334155;
  font-size: 12px;
  font-weight: 700;
  padding: 7px 16px;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.15s ease;
}

.btn-modal-cancel:hover {
  background: #F1F5F9;
  color: #0F172A;
}

@media (max-width: 1024px) {
  .catalog-modal-backdrop {
    padding: 6vh 5vw;
  }
}

@media (max-width: 640px) {
  .catalog-modal-backdrop {
    padding: 12px;
  }
  .modal-hdr,
  .modal-search-strip,
  .modal-items-scroll,
  .modal-ftr {
    padding-left: 16px;
    padding-right: 16px;
  }
}

@media (max-width: 1024px) {
  .context-detail-grid {
    grid-template-columns: 1fr;
  }
  .cost-metrics-row {
    grid-template-columns: repeat(2, 1fr);
  }
  .storage-options-grid {
    grid-template-columns: 1fr;
  }
}
</style>
