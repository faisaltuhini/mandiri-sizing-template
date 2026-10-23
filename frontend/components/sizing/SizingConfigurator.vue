<template>
  <div class="sketch-workspace-grid">
    <!-- LEFT COLUMN: BULLET LIST (• ________) -->
    <section class="sketch-list-section">
      <div class="section-hdr">
        <div class="section-hdr-top">
          <div class="section-hdr-title-wrap">
            <span class="hdr-accent-dot"></span>
            <h2 class="section-title">Komponen &amp; Layer</h2>
          </div>
          <span class="item-count">{{ layerItems.length }} Layers</span>
        </div>
        <button
          type="button"
          class="btn-add-catalog-layer"
          @click="openCatalogModal"
          title="Tambah layer baru langsung dari Katalog Harga"
        >
          <span class="btn-icon-circle">
            <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
              <line x1="12" y1="5" x2="12" y2="19"/>
              <line x1="5" y1="12" x2="19" y2="12"/>
            </svg>
          </span>
          <span class="btn-label-text">Tambah Komponen dari Katalog</span>
        </button>
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
              <div class="bullet-badges">
                <span v-if="item.no > 13 || item.id.startsWith('layer_custom_')" class="catalog-tag" title="Komponen Katalog Kustom">Katalog</span>
                <span v-if="isLayerCustomized(item.id)" class="custom-badge" title="Telah disesuaikan dari baseline">Kustom</span>
              </div>
            </div>
            <span class="bullet-sub">{{ item.component }}</span>
          </div>
          <button
            v-if="item.no > 13 || item.id.startsWith('layer_custom_')"
            type="button"
            class="btn-delete-row"
            title="Hapus layer kustom ini dari Sizing"
            @click.stop="confirmRemoveLayer(item)"
          >
            <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18"></line>
              <line x1="6" y1="6" x2="18" y2="18"></line>
            </svg>
          </button>
          <span v-else class="bullet-chevron">→</span>
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
              <span v-if="activeConfig.pricingSku" class="badge-sku" title="SKU Katalog Harga Terpilih">
                SKU: {{ activeConfig.pricingSku }}
              </span>
              <span v-if="isLayerCustomized(activeItem.id)" class="custom-status-tag">
                ⚠️ Ada Penyesuaian Kustom
              </span>
            </div>
            <h2 class="component-title">{{ activeItem.component }}</h2>
          </div>

          <div class="banner-actions-right">
            <button
              type="button"
              class="btn-toggle-all-steps"
              @click="toggleAllSteps"
              :title="isAllCollapsed ? 'Buka semua point konfigurasi' : 'Ciutkan semua point konfigurasi'"
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
              <span>{{ isAllCollapsed ? 'Buka Semua Point' : 'Ciutkan Semua Point' }}</span>
            </button>

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
              title="Kembalikan konfigurasi layer ini ke standar resmi Enterprise Architecture"
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
          :title="activeConfig.computeType === 'none' ? 'Kuantitas Unit / Lisensi' : 'Kekuatan Server & Komputasi'"
          :guide-question="activeConfig.computeType === 'none' ? 'Berapa banyak kuantitas unit/node lisensi yang dialokasikan?' : 'Berapa banyak server dan seberapa tangguh otak pemroses yang dibutuhkan?'"
          :technical-label="activeConfig.computeType === 'none' ? 'Unit / Kuantitas' : 'vCPU & RAM'"
          :has-changes="isComputeCustomized"
          :collapsible="true"
          :collapsed="collapsedSteps.step1"
          :collapsed-summary="step1Summary"
          @update:collapsed="setStepCollapsed('step1', $event)"
        >
          <!-- License-only banner notice -->
          <div v-if="activeConfig.computeType === 'none'" class="license-only-notice">
            <span class="notice-icon">📦</span>
            <div class="notice-text">
              <span class="notice-title">Mode Lisensi / Add-on Saja (Tanpa Server VM Baru)</span>
              <span class="notice-desc">Komponen ini dihitung berdasarkan kuantitas lisensi unit/node tanpa biaya alokasi komputasi vCPU, RAM, atau SAN Storage baru.</span>
            </div>
          </div>

          <!-- Size Preset Selector (S / M / L / XL / Kustom) - Only when computeType !== 'none' -->
          <div v-if="activeConfig.computeType !== 'none'" class="form-row-block">
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

          <!-- Stepper: Jumlah Server / Kuantitas Unit -->
          <div class="server-count-row">
            <div class="field-hdr-wrap">
              <div class="title-with-desc">
                <span class="field-title-lbl">{{ activeConfig.computeType === 'none' ? 'Kuantitas Unit Lisensi:' : 'Jumlah Server Aktif:' }}</span>
                <span class="field-sub-hint">{{ activeConfig.computeType === 'none' ? 'Berapa unit lisensi / agent yang dibutuhkan' : 'Berapa unit server dedicated yang dialokasikan' }}</span>
              </div>
              <BaseHelpTip
                :title="activeConfig.computeType === 'none' ? 'Kuantitas Lisensi (Unit)' : 'Jumlah Server (Nodes)'"
                :subtitle="activeConfig.computeType === 'none' ? 'Software License Allocation' : 'High Availability Redundancy'"
                :content="activeConfig.computeType === 'none' ? 'Jumlah lisensi atau agent software yang dibutuhkan sesuai skala arsitektur.' : glossary.nodeCount.shortDesc"
                :business-impact="activeConfig.computeType === 'none' ? 'Total biaya lisensi tahunan dihitung dari kuantitas unit dikalikan tarif katalog.' : glossary.nodeCount.businessImpact"
              />
            </div>

            <div class="stepper-wrap-contained">
              <BaseStepper
                v-model="activeConfig.nodeCount"
                :min="1"
                :max="activeConfig.computeType === 'none' ? 2000 : 32"
                :label="activeConfig.computeType === 'none' ? 'Unit Lisensi' : 'Server Aktif'"
                :unit="activeConfig.computeType === 'none' ? 'Unit' : 'Server'"
              />
            </div>
          </div>

          <!-- Advanced Toggle: Granular Steppers - Only when computeType !== 'none' -->
          <BaseAdvancedToggle
            v-if="activeConfig.computeType !== 'none'"
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
          :collapsible="true"
          :collapsed="collapsedSteps.step2"
          :collapsed-summary="step2Summary"
          @update:collapsed="setStepCollapsed('step2', $event)"
        >
          <!-- Notice if computeType === 'none' -->
          <div v-if="activeConfig.computeType === 'none'" class="license-only-notice">
            <span class="notice-icon">💾</span>
            <div class="notice-text">
              <span class="notice-title">Tidak Memerlukan Storage / NetBackup Tambahan</span>
              <span class="notice-desc">Komponen berjenis lisensi / add-on menumpang pada cluster atau server eksisting, sehingga alokasi disk diatur 0 GB (Biaya Storage & Backup: Rp 0).</span>
            </div>
          </div>

          <template v-else>
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
                    <span class="backup-desc">Pencadangan data berkala terpusat Enterprise Architecture (Tarif: Rp 27.000 / GB / tahun).</span>
                  </div>
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
          guide-question="Di atas platform komputasi mana komponen ini akan di-deploy?"
          technical-label="VM / Baremetal / Shared / Lisensi Saja"
          :has-changes="activeConfig.computeType !== recommendedComputeType"
          :collapsible="true"
          :collapsed="collapsedSteps.step3"
          :collapsed-summary="step3Summary"
          @update:collapsed="setStepCollapsed('step3', $event)"
        >
          <BaseChoiceCard
            v-model="activeConfig.computeType"
            :columns="4"
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
              },
              {
                value: 'none',
                title: 'Lisensi / Add-on Saja',
                technicalLabel: 'Tanpa Server Baru',
                description: 'Lisensi software atau agent yang menumpang pada server eksisting (Biaya VM & Storage Rp 0).',
                icon: '📦',
                badge: 'Add-on',
                badgeType: 'info'
              }
            ]"
          />
        </BaseSectionCard>

        <!-- ======================================================== -->
        <!-- STEP 4: MODEL LISENSI SOFTWARE & SKU KATALOG             -->
        <!-- ======================================================== -->
        <BaseSectionCard
          step-number="4"
          title="Pilihan Lisensi Software &amp; Katalog SKU"
          guide-question="Gunakan versi open source gratis atau berlisensi resmi vendor dengan garansi support?"
          technical-label="OSS vs Enterprise / SKU"
          :has-changes="false"
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
                <span class="lic-icon">🌱</span>
                <div class="lic-hdr-titles">
                  <span class="lic-name">Community / Free Open Source</span>
                  <span class="lic-price-tag">Biaya Lisensi: Rp 0 / thn</span>
                </div>
              </div>
              <p class="lic-product-name">{{ activeItem.license?.free || 'Community (Free OSS)' }}</p>
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
              <p class="lic-product-name">{{ activeConfig.licenseProductName || activeItem.license?.enterprise || 'Enterprise License' }}</p>
              <p class="lic-card-desc">Paket langganan resmi vendor dengan garansi SLA respon cepat, update patch keamanan, dan kepatuhan perbankan.</p>
            </div>
          </div>

          <!-- SKU Selector & Rate Display when Enterprise is selected -->
          <div v-if="activeConfig.licenseKind === 'enterprise'" class="enterprise-sku-picker-box">
            <div class="sku-picker-header">
              <div class="sku-picker-title-wrap">
                <span class="sku-picker-lbl">Pilihan SKU Lisensi Katalog Harga:</span>
                <span class="sku-picker-sub">Tentukan atau ganti item katalog harga resmi untuk perhitungan biaya otomatis</span>
              </div>
              <span v-if="activeConfig.pricingSku" class="sku-selected-tag">SKU: {{ activeConfig.pricingSku }}</span>
            </div>

            <div class="sku-picker-control-row">
              <select
                :value="activeConfig.pricingSku || ''"
                class="sku-dropdown-select"
                @change="onSkuSelectChange(($event.target as HTMLSelectElement).value)"
              >
                <option value="">-- Standar Enterprise Layer --</option>
                <optgroup label="Katalog Lisensi Software">
                  <option
                    v-for="item in licenseCatalogItems"
                    :key="item.sku"
                    :value="item.sku"
                  >
                    {{ item.sku }} — {{ item.name }} ({{ formatSimpleIDR(item.unit_price_idr) ? 'Rp ' + formatSimpleIDR(item.unit_price_idr) : 'Pending STA' }} / {{ item.unit }})
                  </option>
                </optgroup>
                <optgroup v-if="otherCatalogItems.length > 0" label="Item Katalog Lainnya">
                  <option
                    v-for="item in otherCatalogItems"
                    :key="item.sku"
                    :value="item.sku"
                  >
                    {{ item.sku }} — {{ item.name }} ({{ formatSimpleIDR(item.unit_price_idr) ? 'Rp ' + formatSimpleIDR(item.unit_price_idr) : 'Pending STA' }} / {{ item.unit }})
                  </option>
                </optgroup>
              </select>
            </div>

            <!-- Active SKU Detail Pill -->
            <div v-if="activeCatalogPricingItem" class="active-sku-detail-card">
              <div class="sku-card-main">
                <span class="sku-card-name">{{ activeCatalogPricingItem.name }}</span>
                <span class="sku-card-env-badge">{{ activeCatalogPricingItem.environment }} · {{ activeCatalogPricingItem.unit }}</span>
              </div>
              <div class="sku-card-price-wrap">
                <span class="sku-price-title">Tarif Satuan Katalog:</span>
                <span v-if="activeCatalogPricingItem.unit_price_idr" class="sku-price-val">
                  Rp {{ formatSimpleIDR(activeCatalogPricingItem.unit_price_idr) }} <span class="sku-price-unit">/ {{ activeCatalogPricingItem.unit }}</span>
                </span>
                <span v-else class="sku-price-pending">Pending Konfirmasi STA</span>
              </div>
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
          guide-question="Spesifikasi acuan resmi Enterprise Architecture dan justifikasi Enterprise Architecture"
          technical-label="Baseline Standard"
          :collapsible="true"
          :collapsed="collapsedSteps.step5"
          :collapsed-summary="step5Summary"
          @update:collapsed="setStepCollapsed('step5', $event)"
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
          <button type="button" class="modal-btn-close" @click="showCatalogModal = false" title="Tutup">✕</button>
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
              ✕
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
            <span class="empty-icon">{{ catalogModalFilter === 'unmapped' ? '🎉' : '🔍' }}</span>
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
                  ✓ Terpasang di {{ getCatalogItemSizingInfo(cat).matchedLayerText }}
                </span>
                <span v-else-if="getCatalogItemSizingInfo(cat).isRawInfra" class="status-pill infra">
                  ⚙️ Tarif Dasar Mesin &amp; Storage
                </span>
                <span v-else class="status-pill ready">
                  ⚡ Siap Dimasukkan
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
                <span class="price-period">per {{ cat.unit || 'unit' }} / tahun</span>
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
  tierFilter: 'critical' | 'moderate'
  licenseFilter: 'free' | 'license'
  isLayerCustomized: (id: string) => boolean
  pricingList?: PricingItemRecord[]
}>()

const emit = defineEmits<{
  (e: 'select-layer', id: string): void
  (e: 'reset-active-layer'): void
  (e: 'add-catalog-to-sizing', item: PricingItemRecord): void
  (e: 'remove-custom-layer', id: string): void
}>()

const glossary = SIZING_GLOSSARY

// Modal state
const showCatalogModal = ref(false)
const catalogModalSearch = ref('')
const catalogModalFilter = ref<'all' | 'license' | 'infra' | 'unmapped'>('unmapped')

function openCatalogModal() {
  catalogModalSearch.value = ''
  catalogModalFilter.value = 'unmapped'
  showCatalogModal.value = true
}

function confirmRemoveLayer(item: TechStackItem) {
  if (confirm(`Hapus layer "${item.layer} - ${item.component}" dari daftar Sizing?`)) {
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

  // 1. Raw infrastructure pricing unit (hardware chassis, vCPU, SAN storage, backup storage per GB)
  const isRawInfra =
    ['compute', 'storage', 'baremetal'].includes(catLower) ||
    ['VCPU_PROD', 'VCPU_DEV', 'BM_STD_PROD', 'BM_STD_DEV', 'STG_T2_PROD', 'STG_T2_DEV', 'STG_T3_PROD', 'NBU_ALL', 'RUBRIK_PROD'].includes(skuUpper) ||
    (catLower === 'backup' && (unitLower.includes('gb') || unitLower.includes('tb'))) ||
    nameLower.includes('compute vcpu') ||
    nameLower.includes('san storage') ||
    nameLower.includes('server baremetal') ||
    nameLower.includes('backup netbackup') ||
    nameLower.includes('rubrik backup')

  // 2. Direct match with custom layer in sizing (via contractRef or component name)
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

  // 3. Match known preset software keyword families
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

function isItemAlreadyInSizing(sku: string): boolean {
  if (!sku) return false
  const target = (props.pricingList || []).find(p => p.sku === sku)
  if (target) {
    return getCatalogItemSizingInfo(target).isAlreadyInSizing
  }
  return false
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

const otherCatalogItems = computed<PricingItemRecord[]>(() => {
  return (props.pricingList || []).filter(item => item.category !== 'license')
})

const activeCatalogPricingItem = computed<PricingItemRecord | undefined>(() => {
  if (!props.activeConfig?.pricingSku || !props.pricingList) return undefined
  return props.pricingList.find(p => p.sku === props.activeConfig.pricingSku)
})

function onSkuSelectChange(sku: string) {
  props.activeConfig.pricingSku = sku
  if (sku && props.pricingList) {
    const found = props.pricingList.find(p => p.sku === sku)
    if (found) {
      props.activeConfig.licenseProductName = found.name
      props.activeConfig.licenseMetric = found.unit
    }
  }
}

// Active preset baseline reference for accurate customization detection
const activePresetBaseline = computed(() => {
  const isCrit = props.tierFilter === 'critical'
  const preset = DEFAULT_PRESET_ITEMS.find(p => p.id === props.activeItem.id)
  if (preset) {
    return isCrit ? preset.critical : preset.moderate
  }
  return null
})

// Check recommended baseline spec for sizing comparison
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

// ========================================================
// PERSISTENT COLLAPSIBLE STEPS (localStorage)
// ========================================================
const defaultCollapsedSteps: Record<string, boolean> = {
  step1: false,
  step2: false,
  step3: false,
  step4: false,
  step5: false,
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
            step5: Boolean(parsed.step5),
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
    step5: nextVal,
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

// Collapsed Step Summaries
const step1Summary = computed(() => {
  if (props.activeConfig.computeType === 'none') {
    return `${props.activeConfig.nodeCount} Unit Lisensi / Add-on`
  }
  const vcpu = props.activeConfig.vcpuPerNode
  const ram = props.activeConfig.ramGBPerNode
  let presetName = 'Kustom'
  if (vcpu === 2 && ram === 4) presetName = 'Kecil'
  else if (vcpu === 4 && ram === 8) presetName = 'Sedang'
  else if (vcpu === 8 && ram === 16) presetName = 'Standar'
  else if (vcpu === 16 && ram === 32) presetName = 'Besar'
  else if (vcpu === 64 && ram === 1024) presetName = 'Raksasa (Heavy)'

  return `${props.activeConfig.nodeCount} Server Aktif · ${presetName} (${vcpu} vCPU / ${ram} GB RAM)`
})

const step2Summary = computed(() => {
  if (props.activeConfig.computeType === 'none') {
    return 'Disk 0 GB (Menumpang cluster eksisting)'
  }
  const cls = props.activeConfig.storageClass || 'Tier 2'
  const backup = props.activeConfig.backupEnabled ? 'Backup Harian Aktif' : 'Tanpa Backup'
  return `${props.activeConfig.storageGBPerNode} GB / server (${cls}) · ${backup}`
})

const step3Summary = computed(() => {
  const map: Record<string, string> = {
    vm: 'Virtual Machine (VM)',
    baremetal: 'Server Fisik Khusus (Baremetal Dedicated)',
    shared: 'Penyimpanan Berbagi (Shared / SAN LUN)',
    none: 'Lisensi / Add-on Saja (Tanpa Server Baru)'
  }
  return map[props.activeConfig.computeType] || 'Virtual Machine (VM)'
})

const step4Summary = computed(() => {
  if (props.activeConfig.licenseKind === 'free') {
    return 'Community Free (OSS) · Tanpa Biaya Lisensi'
  }
  const sku = props.activeConfig.pricingSku
  const prodName = props.activeConfig.licenseProductName || props.activeLayerCost?.license_product_name
  if (sku) {
    return `Lisensi Enterprise · SKU: ${sku}${prodName ? ' (' + prodName + ')' : ''}`
  }
  return prodName ? `Lisensi Enterprise · ${prodName}` : 'Lisensi Resmi Vendor (Enterprise)'
})

const step5Summary = computed(() => {
  const isCrit = props.tierFilter === 'critical'
  const spec = isCrit ? props.activeItem.critical : props.activeItem.moderate
  const tLabel = isCrit ? 'Critical' : 'Moderate'
  return `Acuan EA (Tier ${tLabel}): ${spec.nodes} Node · ${spec.spec} (${spec.type})`
})
</script>

<style scoped>
.sketch-workspace-grid {
  flex: 1;
  display: grid;
  grid-template-columns: 315px 1fr;
  overflow: hidden;
  min-height: 0;
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
  font-size: 12px;
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
  padding: 20px 24px 110px 24px;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
}

.btn-toggle-all-steps {
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

.btn-toggle-all-steps:hover {
  color: #002060;
  border-color: #002060;
  background: #EFF6FF;
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

/* Add Catalog button and left list enhancements */
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
}

.btn-delete-row:hover {
  background: #FEE2E2;
  color: #DC2626;
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
  border-radius: 6px;
  padding: 12px 14px;
  margin-bottom: 12px;
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
  font-size: 11px;
  color: #0284C7;
  line-height: 1.4;
}

/* Enterprise SKU picker box */
.enterprise-sku-picker-box {
  margin-top: 14px;
  padding: 14px 16px;
  background: #F8FAFC;
  border: 1px solid #E2E8F0;
  border-radius: 6px;
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
  padding: 20px;
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
  max-width: 860px;
  max-height: 88vh;
  display: flex;
  flex-direction: column;
  box-shadow: 0 25px 50px -12px rgba(15, 23, 42, 0.25), 0 0 0 1px rgba(15, 23, 42, 0.08);
  overflow: hidden;
}

.modal-hdr {
  padding: 20px 24px;
  border-bottom: 1px solid #E2E8F0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: #FFFFFF;
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
  padding: 14px 24px;
  background: #F8FAFC;
  border-bottom: 1px solid #E2E8F0;
  display: flex;
  flex-direction: column;
  gap: 12px;
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
  overflow-y: auto;
  padding: 16px 24px;
  display: flex;
  flex-direction: column;
  gap: 12px;
  max-height: 480px;
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
  background: #EEF2FF;
  color: #4F46E5;
  border: 1px solid #C7D2FE;
}

.card-icon-avatar.theme-database {
  background: #ECFDF5;
  color: #059669;
  border: 1px solid #A7F3D0;
}

.card-icon-avatar.theme-security {
  background: #FEF2F2;
  color: #DC2626;
  border: 1px solid #FECACA;
}

.card-icon-avatar.theme-backup {
  background: #FAF5FF;
  color: #7C3AED;
  border: 1px solid #E9D5FF;
}

.card-icon-avatar.theme-license {
  background: #F0FDF4;
  color: #16A34A;
  border: 1px solid #BBF7D0;
}

.card-icon-avatar.theme-infra {
  background: #F1F5F9;
  color: #475569;
  border: 1px solid #CBD5E1;
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
  padding: 14px 24px;
  border-top: 1px solid #E2E8F0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: #FFFFFF;
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
</style>
