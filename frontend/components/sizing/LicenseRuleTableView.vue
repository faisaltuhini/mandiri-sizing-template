<template>
  <div class="license-view-container">
    <!-- TOPBAR WITH TOGGLE & ACTIONS -->
    <header class="sketch-topbar license-topbar">
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
          <h1 class="page-title">Aturan Lisensi Software &amp; Metrik</h1>
        </div>
        <p class="page-subtitle">Formula kalkulasi lisensi enterprise: per worker node, per 5 pods (AppDynamics), per DB agent, per shard (Redis min 3 quorum), per core.</p>
      </div>

      <div class="topbar-actions">
        <button
          type="button"
          class="btn-primary-action add-rule-btn"
          @click="openAddModal"
        >
          <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
            <line x1="12" y1="5" x2="12" y2="19"></line>
            <line x1="5" y1="12" x2="19" y2="12"></line>
          </svg>
          <span>Tambah Aturan Lisensi</span>
        </button>
      </div>
    </header>

    <!-- CONTENT WRAPPER -->
    <div class="license-content-scroll">
      <!-- SUMMARY KPI METRICS CARDS -->
      <section class="license-kpi-grid">
        <div class="kpi-card">
          <div class="kpi-icon-wrap icon-blue">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
              <polyline points="14 2 14 8 20 8"></polyline>
              <line x1="16" y1="13" x2="8" y2="13"></line>
              <line x1="16" y1="17" x2="8" y2="17"></line>
            </svg>
          </div>
          <div class="kpi-info">
            <span class="kpi-label">Total Aturan Lisensi</span>
            <div class="kpi-value-row">
              <span class="kpi-num">{{ rules.length }}</span>
              <span class="kpi-sub">Rule aktif: {{ activeRulesCount }}</span>
            </div>
          </div>
        </div>

        <div class="kpi-card highlight-card">
          <div class="kpi-icon-wrap icon-cyan">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <rect x="2" y="2" width="20" height="8" rx="2" ry="2"></rect>
              <rect x="2" y="14" width="20" height="8" rx="2" ry="2"></rect>
              <line x1="6" y1="6" x2="6.01" y2="6"></line>
              <line x1="6" y1="18" x2="6.01" y2="18"></line>
            </svg>
          </div>
          <div class="kpi-info">
            <span class="kpi-label">AppDynamics Monitoring</span>
            <div class="kpi-value-row">
              <span class="kpi-num">Pack 5 Pods</span>
              <span class="kpi-tag-sub">Rp 8.39M / 5 Pods · DB Agent Rp 7.51M</span>
            </div>
          </div>
        </div>

        <div class="kpi-card">
          <div class="kpi-icon-wrap icon-gold">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"></path>
            </svg>
          </div>
          <div class="kpi-info">
            <span class="kpi-label">Redis Quorum Rule</span>
            <div class="kpi-value-row">
              <span class="kpi-num">Min 3 Nodes</span>
              <span class="kpi-sub">HA Sentinel Cluster Quorum</span>
            </div>
          </div>
        </div>

        <div class="kpi-card">
          <div class="kpi-icon-wrap icon-purple">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <rect x="4" y="4" width="16" height="16" rx="2" ry="2"></rect>
              <rect x="9" y="9" width="6" height="6"></rect>
              <line x1="9" y1="1" x2="9" y2="4"></line>
              <line x1="15" y1="1" x2="15" y2="4"></line>
              <line x1="9" y1="20" x2="9" y2="23"></line>
              <line x1="15" y1="20" x2="15" y2="23"></line>
            </svg>
          </div>
          <div class="kpi-info">
            <span class="kpi-label">RH OpenShift (OCP)</span>
            <div class="kpi-value-row">
              <span class="kpi-num">Per Worker</span>
              <span class="kpi-sub">Master &amp; Router Free of Charge</span>
            </div>
          </div>
        </div>
      </section>

      <!-- SEARCH & FILTER TOOLBAR -->
      <section class="license-filter-bar">
        <div class="search-input-wrap">
          <svg class="search-icon" width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="11" cy="11" r="8"></circle>
            <line x1="21" y1="21" x2="16.65" y2="16.65"></line>
          </svg>
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Cari aturan lisensi, vendor (AppDynamics, Redis, Red Hat), SKU, atau layer..."
            class="search-input"
          />
          <button v-if="searchQuery" @click="searchQuery = ''" class="clear-search-btn" title="Hapus pencarian">
            &times;
          </button>
        </div>

        <div class="filter-controls">
          <div class="filter-group">
            <label class="filter-lbl">Kategori:</label>
            <select v-model="categoryFilter" class="filter-select">
              <option value="all">Semua Kategori</option>
              <option value="apm">APM / Monitoring</option>
              <option value="database">Database</option>
              <option value="cache">Distributed Cache</option>
              <option value="platform">Platform / K8s</option>
              <option value="ingress">Web Ingress</option>
              <option value="streaming">Event Streaming</option>
            </select>
          </div>

          <div class="filter-group">
            <label class="filter-lbl">Metrik:</label>
            <select v-model="metricFilter" class="filter-select">
              <option value="all">Semua Metrik</option>
              <option value="per_5_pods">Per 5 Pods (Pack 5)</option>
              <option value="per_node_db">Per Node DB Agent</option>
              <option value="per_shard">Per Shard (Min 3 Node)</option>
              <option value="per_worker">Per Worker Node</option>
              <option value="per_core">Per Core Fisik</option>
              <option value="per_instance">Per Instance / VM</option>
            </select>
          </div>

          <button
            v-if="searchQuery || categoryFilter !== 'all' || metricFilter !== 'all'"
            type="button"
            class="reset-filters-btn"
            @click="resetFilters"
            title="Reset semua filter"
          >
            Reset Filter
          </button>
        </div>
      </section>

      <!-- RULES TABLE SECTION -->
      <section class="rules-table-card">
        <div class="table-card-header">
          <div class="table-title-row">
            <h2 class="table-title">Daftar Aturan Lisensi ({{ filteredRules.length }})</h2>
            <span class="table-badge">{{ categoryFilter === 'all' ? 'Seluruh Komponen' : categoryFilter.toUpperCase() }}</span>
          </div>
          <button
            type="button"
            class="btn-refresh"
            @click="fetchRules"
            :disabled="isLoading"
            title="Muat ulang aturan dari database SQL Server"
          >
            <svg :class="{ spinning: isLoading }" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M23 4v6h-6"></path>
              <path d="M1 20v-6h6"></path>
              <path d="M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15"></path>
            </svg>
            <span>Refresh</span>
          </button>
        </div>

        <div v-if="isLoading" class="table-loading-state">
          <div class="spinner-circle"></div>
          <p>Memuat aturan lisensi dari Microsoft SQL Server...</p>
        </div>

        <div v-else-if="filteredRules.length === 0" class="table-empty-state">
          <div class="empty-icon-wrap">
            <svg width="36" height="36" viewBox="0 0 24 24" fill="none" stroke="#D4A017" stroke-width="1.5">
              <circle cx="12" cy="12" r="10"></circle>
              <line x1="8" y1="12" x2="16" y2="12"></line>
            </svg>
          </div>
          <h3>Tidak Ada Aturan Lisensi yang Cocok</h3>
          <p>Coba sesuaikan kata kunci pencarian atau ubah filter kategori/metrik.</p>
          <button type="button" class="btn-primary-action btn-sm" @click="resetFilters">
            Reset Filter
          </button>
        </div>

        <div v-else class="table-responsive">
          <table class="sketch-table license-table">
            <thead>
              <tr>
                <th style="width: 28%">Aturan &amp; Vendor</th>
                <th style="width: 14%">Kategori &amp; Layer</th>
                <th style="width: 18%">Metrik &amp; Pack</th>
                <th style="width: 14%">Syarat Cluster</th>
                <th style="width: 14%">SKU &amp; Harga Satuan</th>
                <th style="width: 12%; text-align: right">Aksi</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="rule in filteredRules"
                :key="rule.id"
                :class="{ 'inactive-row': !rule.is_active }"
              >
                <!-- Aturan & Vendor -->
                <td>
                  <div class="rule-name-wrap">
                    <span class="rule-title">{{ rule.rule_name }}</span>
                    <div class="rule-meta-badges">
                      <span class="vendor-badge">{{ rule.vendor }}</span>
                      <span v-if="rule.is_active" class="status-pill status-active">Aktif</span>
                      <span v-else class="status-pill status-inactive">Non-Aktif</span>
                    </div>
                    <p class="rule-desc" :title="rule.description">{{ rule.description }}</p>
                  </div>
                </td>

                <!-- Kategori & Layer -->
                <td>
                  <div class="layer-meta">
                    <span :class="['cat-pill', 'cat-' + rule.category]">
                      {{ formatCategory(rule.category) }}
                    </span>
                    <span class="target-layer-name">{{ rule.target_layer || 'Seluruh Layer' }}</span>
                  </div>
                </td>

                <!-- Metrik & Pack -->
                <td>
                  <div class="metric-info-block">
                    <span class="metric-type-pill">
                      {{ formatMetricType(rule.metric_type) }}
                    </span>
                    <div class="pack-size-row">
                      <span class="pack-size-lbl">Kelipatan / Pack:</span>
                      <strong class="pack-size-val">{{ rule.pack_size }} unit</strong>
                    </div>
                  </div>
                </td>

                <!-- Syarat Cluster -->
                <td>
                  <div v-if="rule.min_cluster_node > 1" class="cluster-quorum-badge">
                    <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
                      <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"></path>
                    </svg>
                    <span>Min {{ rule.min_cluster_node }} Nodes (Quorum)</span>
                  </div>
                  <div v-else class="cluster-single-badge">
                    <span>1 Node Standalone</span>
                  </div>
                </td>

                <!-- SKU & Harga Satuan -->
                <td>
                  <div class="pricing-info-cell">
                    <span v-if="rule.pricing_sku" class="sku-code">{{ rule.pricing_sku }}</span>
                    <div v-if="rule.unit_price_idr && rule.unit_price_idr > 0" class="price-val">
                      {{ formatCurrency(rule.unit_price_idr) }}
                      <small class="price-sub">/ Thn</small>
                    </div>
                    <div v-else class="unconfirmed-price">
                      <span class="badge-sta">Menunggu Konfirmasi STA</span>
                    </div>
                  </div>
                </td>

                <!-- Aksi -->
                <td style="text-align: right">
                  <div class="row-actions">
                    <button
                      type="button"
                      class="btn-action-icon btn-edit"
                      @click="openEditModal(rule)"
                      title="Ubah Aturan Lisensi"
                    >
                      <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path>
                        <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path>
                      </svg>
                    </button>
                    <button
                      type="button"
                      class="btn-action-icon btn-delete"
                      @click="confirmDelete(rule)"
                      title="Hapus Aturan"
                    >
                      <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <polyline points="3 6 5 6 21 6"></polyline>
                        <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
                      </svg>
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>
    </div>

    <!-- MODAL: ADD / EDIT RULE (EXECUTIVE & NON-TECH FRIENDLY) -->
    <transition name="modal-fade">
      <div v-if="showModal" class="sketch-modal-overlay" @click.self="showModal = false">
        <div class="sketch-modal-card rule-modal-card-friendly" @click.stop>
          <!-- HEADER -->
          <div class="modal-card-header">
            <div class="modal-title-wrap">
              <div class="modal-icon-shield">
                <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="#002060" stroke-width="2.2">
                  <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"></path>
                  <path d="M9 12l2 2 4-4"></path>
                </svg>
              </div>
              <div>
                <h3 class="modal-heading">{{ isEditing ? 'Ubah Aturan Lisensi Software' : 'Tambah Aturan Lisensi Baru' }}</h3>
                <p class="modal-subheading">Panduan konfigurasi cara perhitungan, syarat arsitektur cluster, dan tarif resmi Enterprise Architecture</p>
              </div>
            </div>
            <button type="button" class="modal-close-btn" @click="showModal = false" title="Tutup">&times;</button>
          </div>

          <form @submit.prevent="saveRule" class="modal-form-friendly">
            <!-- SECTION 1: IDENTITAS SOFTWARE -->
            <div class="form-section-card">
              <div class="section-card-title">
                <span class="section-step-num">1</span>
                <div>
                  <h4 class="section-title-text">Identitas Software &amp; Vendor</h4>
                  <p class="section-desc-text">Nama produk software resmi dan vendor penyedia solusi</p>
                </div>
              </div>

              <div class="section-grid-2">
                <div class="form-group">
                  <div class="form-label-row">
                    <label class="form-label">Nama Software / Produk Lisensi <span class="req">*</span></label>
                    <span v-if="matchedCatalogComponent" class="label-match-tag match-yes">
                      ✓ Master Arsitektur
                    </span>
                  </div>
                  <BaseAutocomplete
                    v-model="form.rule_name"
                    :options="softwareNameOptions"
                    placeholder="Ketik atau cari nama software (cth: Redis, AppDynamics, OCP)..."
                    required
                    @select="onSoftwareNameSelected"
                  />
                  <span class="field-hint">Memilih dari katalog akan otomatis melengkapi vendor, kategori, layer &amp; SKU harga</span>
                </div>

                <div class="form-group">
                  <label class="form-label">Vendor Resmi <span class="req">*</span></label>
                  <input
                    v-model="form.vendor"
                    type="text"
                    required
                    placeholder="Contoh: Redis Ltd, Cisco / AppDynamics, Red Hat, Oracle, F5"
                    class="form-input"
                  />
                </div>
              </div>

              <div class="section-grid-2" style="margin-top: 12px">
                <div class="form-group">
                  <label class="form-label">Kategori Komponen</label>
                  <select v-model="form.category" class="form-select">
                    <option value="apm">APM / Observability &amp; Monitoring</option>
                    <option value="database">Core Database Engine</option>
                    <option value="cache">Distributed Memory Cache</option>
                    <option value="platform">Platform Container &amp; Kubernetes</option>
                    <option value="ingress">Web Ingress &amp; Load Balancer</option>
                    <option value="streaming">Event Streaming &amp; Message Bus</option>
                    <option value="other">Komponen Lainnya</option>
                  </select>
                </div>

                <div class="form-group">
                  <label class="form-label">Target Layer Arsitektur</label>
                  <select v-model="form.target_layer" class="form-select">
                    <option value="App Container">Layer 2: App Container (Microservices)</option>
                    <option value="Core Database">Layer 4: Core Database (Oracle/PostgreSQL)</option>
                    <option value="Distributed Cache">Layer 6: Distributed Cache (Redis)</option>
                    <option value="Web Ingress">Layer 1: Web Ingress (NGINX/F5)</option>
                    <option value="Container Router">Layer 3: Container Router &amp; Bastion</option>
                    <option value="Message Bus">Layer 8: Message Bus (Kafka)</option>
                    <option value="Log Analytics">Layer 9: Log Analytics (ELK/OpenSearch)</option>
                    <option value="All Layers">Semua Layer Arsitektur</option>
                  </select>
                </div>
              </div>
            </div>

            <!-- SECTION 2: METRIK & CARA HITUNG (VISUAL CARD SELECTOR) -->
            <div class="form-section-card">
              <div class="section-card-title">
                <span class="section-step-num">2</span>
                <div>
                  <h4 class="section-title-text">Model &amp; Cara Perhitungan Lisensi</h4>
                  <p class="section-desc-text">Pilih salah satu metode perhitungan yang berlaku pada kontrak software ini</p>
                </div>
              </div>

              <!-- CARDS GRID -->
              <div class="metric-cards-selector-grid">
                <div
                  v-for="opt in metricCardsList"
                  :key="opt.id"
                  :class="['metric-choice-card', { selected: form.metric_type === opt.id }]"
                  @click="selectMetricCard(opt)"
                >
                  <div class="choice-card-top">
                    <span class="choice-icon">{{ opt.icon }}</span>
                    <span v-if="form.metric_type === opt.id" class="choice-check-badge">
                      <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
                        <polyline points="20 6 9 17 4 12"></polyline>
                      </svg>
                    </span>
                  </div>
                  <h5 class="choice-title">{{ opt.title }}</h5>
                  <span class="choice-sub-tag">{{ opt.badge }}</span>
                  <p class="choice-desc">{{ opt.desc }}</p>
                </div>
              </div>

              <!-- DYNAMIC SIMULATION & EXPLANATION CALLOUT -->
              <div class="metric-human-guide-box">
                <div class="guide-header">
                  <span class="guide-lamp-icon">💡</span>
                  <span class="guide-heading">Simulasi &amp; Ketentuan Perhitungan Otomatis:</span>
                </div>
                <p class="guide-body-text">
                  {{ getMetricExplanation(form.metric_type, form.pack_size, form.min_cluster_node) }}
                </p>
              </div>

              <!-- CONTEXTUAL ADJUSTERS (ONLY SHOWN IF RELEVANT) -->
              <div v-if="form.metric_type === 'per_5_pods'" class="contextual-adjuster-panel">
                <div class="adjuster-icon">📦</div>
                <div class="adjuster-content">
                  <label class="adjuster-title">Ukuran Kelipatan Pod per Paket Lisensi (Pack Size):</label>
                  <p class="adjuster-hint">Setiap kelipatan 5 pods dihitung 1 paket lisensi (contoh: 12 pods = 3 paket).</p>
                </div>
                <div class="adjuster-input-wrap">
                  <input
                    v-model.number="form.pack_size"
                    type="number"
                    min="1"
                    class="adjuster-num-input"
                  />
                  <span class="adjuster-unit">Pods / Paket</span>
                </div>
              </div>

              <div v-if="form.metric_type === 'per_shard'" class="contextual-adjuster-panel alert-quorum">
                <div class="adjuster-icon">🛡️</div>
                <div class="adjuster-content">
                  <label class="adjuster-title">Syarat Minimal Node Cluster (Quorum High Availability):</label>
                  <p class="adjuster-hint">Ketentuan arsitektur HA Enterprise Architecture mewajibkan minimal 3 node untuk auto-failover sentinel/raft.</p>
                </div>
                <div class="adjuster-input-wrap">
                  <input
                    v-model.number="form.min_cluster_node"
                    type="number"
                    min="1"
                    class="adjuster-num-input"
                  />
                  <span class="adjuster-unit">Nodes Min</span>
                </div>
              </div>
            </div>

            <!-- SECTION 3: HARGA & TAUTAN KATALOG SKU -->
            <div class="form-section-card">
              <div class="section-card-title">
                <span class="section-step-num">3</span>
                <div>
                  <h4 class="section-title-text">Tautan Katalog Harga &amp; Tarif Resmi</h4>
                  <p class="section-desc-text">Hubungkan ke SKU master data katalog dan tentukan tarif tahunan resmi</p>
                </div>
              </div>

              <div class="section-grid-2">
                <!-- SMART SKU AUTOCOMPLETE -->
                <div class="form-group">
                  <div class="form-label-row">
                    <label class="form-label">
                      Tautan SKU Katalog Enterprise Architecture
                      <span class="req">*</span>
                    </label>
                    <span v-if="matchedCatalogItem" class="label-match-tag match-yes">
                      ✓ Terverifikasi Katalog
                    </span>
                    <span v-else-if="form.pricing_sku" class="label-match-tag match-custom">
                      SKU Kustom
                    </span>
                  </div>

                  <!-- BaseAutocomplete component -->
                  <BaseAutocomplete
                    v-model="form.pricing_sku"
                    :options="pricingSkuOptions"
                    placeholder="Cari SKU atau nama katalog (cth: LIC_REDIS, APPDYN, OCP)..."
                    required
                    @select="onPricingSkuSelected"
                  />

                  <!-- Live Matched Card / Custom Info -->
                  <div v-if="matchedCatalogItem" class="sku-matched-card">
                    <div class="matched-top-line">
                      <span class="catalog-verified-badge">
                        <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
                          <polyline points="20 6 9 17 4 12"></polyline>
                        </svg>
                        Tersinkronisasi Master Katalog Enterprise
                      </span>
                      <span class="matched-env-badge">{{ matchedCatalogItem.environment || 'Production' }}</span>
                    </div>
                    <div class="matched-info-grid">
                      <div class="matched-info-left">
                        <strong class="matched-title">{{ matchedCatalogItem.name }}</strong>
                        <span class="matched-unit-text">Satuan: {{ matchedCatalogItem.unit || 'Unit / Thn' }}</span>
                      </div>
                      <div class="matched-price-tag">
                        {{ formatCurrency(matchedCatalogItem.unit_price_idr) }}
                      </div>
                    </div>
                    <p v-if="matchedCatalogItem.source" class="matched-source-note">
                      📋 Ref: {{ matchedCatalogItem.source }}
                    </p>
                  </div>
                  <div v-else-if="form.pricing_sku" class="sku-custom-card">
                    <span class="custom-sku-badge">ℹ️ SKU Kustom</span>
                    <span class="custom-sku-text">
                      Kode <strong>{{ form.pricing_sku }}</strong> tidak ada di katalog master pricing Enterprise Architecture.
                    </span>
                  </div>
                  <span v-else class="field-hint">Ketik untuk mencari SKU dari database katalog resmi Enterprise Architecture</span>
                </div>

                <!-- TARIF SATUAN RESMI (READ-ONLY, MURNI DARI KATALOG) -->
                <div class="form-group">
                  <div class="form-label-row">
                    <label class="form-label">Tarif Satuan Lisensi Resmi (IDR / Tahun)</label>
                    <span class="label-match-tag match-locked">
                      🔒 Murni Dari Master Katalog
                    </span>
                  </div>

                  <div class="catalog-price-readonly-box" :class="{ 'has-price': form.unit_price_idr && form.unit_price_idr > 0, 'pending-price': !form.unit_price_idr }">
                    <div class="price-val-row">
                      <span class="price-currency-label">Rp</span>
                      <span class="price-large-number">
                        {{ form.unit_price_idr ? Math.round(form.unit_price_idr).toLocaleString('id-ID') : '0' }}
                      </span>
                      <span class="price-unit-tag">
                        {{ matchedCatalogItem ? `/ ${matchedCatalogItem.unit}` : '/ Tahun' }}
                      </span>
                    </div>

                    <div v-if="form.unit_price_idr && form.unit_price_idr > 0" class="price-sync-status status-ok">
                      <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
                        <polyline points="20 6 9 17 4 12"></polyline>
                      </svg>
                      <span>Tarif resmi aktif otomatis terkunci dari katalog master Enterprise Architecture</span>
                    </div>
                    <div v-else class="price-sync-status status-pending">
                      <span>⚠️ Belum ada tarif kontrak (Menunggu konfirmasi Tim STA / Procurement)</span>
                    </div>
                  </div>

                  <span class="field-hint">Tarif tidak dapat diubah manual untuk menjaga kepatuhan dan integritas sizing terhadap katalog resmi</span>
                </div>
              </div>
            </div>

            <!-- SECTION 4: STATUS & CATATAN KONTRAK -->
            <div class="form-section-card">
              <div class="section-card-title">
                <span class="section-step-num">4</span>
                <div>
                  <h4 class="section-title-text">Status &amp; Dokumen Kontrak Pengadaan</h4>
                  <p class="section-desc-text">Status operasional dan catatan nomor kontrak Enterprise Architecture</p>
                </div>
              </div>

              <!-- STATUS TOGGLE CARD -->
              <div class="status-toggle-card" :class="{ 'is-active-card': form.is_active }">
                <label class="toggle-checkbox-row">
                  <input v-model="form.is_active" type="checkbox" class="modern-checkbox" />
                  <div class="toggle-texts">
                    <span class="toggle-heading">{{ form.is_active ? 'Aturan Lisensi Aktif' : 'Aturan Lisensi Non-Aktif (Diabaikan)' }}</span>
                    <span class="toggle-sub">
                      {{ form.is_active ? 'Aturan ini akan otomatis diterapkan pada kalkulasi biaya Sizing Arsitektur.' : 'Aturan ini disimpan namun dinonaktifkan sementara dari kalkulasi.' }}
                    </span>
                  </div>
                </label>
              </div>

              <div class="form-group" style="margin-top: 14px">
                <label class="form-label">Deskripsi, Klausul &amp; Nomor Kontrak Pengadaan</label>
                <textarea
                  v-model="form.description"
                  rows="3"
                  placeholder="Contoh: Kontrak Red Hat TIF.DDL/GH.1993/2025. Lisensi dihitung per Worker node. Master, Infra Router & Bastion tidak dikenakan biaya."
                  class="form-textarea"
                ></textarea>
                <span class="field-hint">Catatan ini akan otomatis muncul pada sheet 'Pricing' dan 'Formula' saat diexport ke Excel.</span>
              </div>
            </div>

            <!-- MODAL ACTIONS -->
            <div class="modal-form-actions-friendly">
              <button
                type="button"
                class="btn-secondary-friendly"
                @click="showModal = false"
              >
                Batal
              </button>
              <button
                type="submit"
                class="btn-primary-friendly"
                :disabled="isSaving"
              >
                <div v-if="isSaving" class="spinner-sm-white"></div>
                <template v-else>
                  <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
                    <polyline points="20 6 9 17 4 12"></polyline>
                  </svg>
                  <span>{{ isEditing ? 'Simpan Perubahan Aturan' : 'Tambah Aturan Lisensi' }}</span>
                </template>
              </button>
            </div>
          </form>
        </div>
      </div>
    </transition>

    <!-- MODAL: CONFIRM DELETE -->
    <transition name="modal-fade">
      <div v-if="showDeleteModal" class="sketch-modal-overlay" @click.self="showDeleteModal = false">
        <div class="sketch-modal-card delete-confirm-card" @click.stop>
          <div class="delete-icon-wrap">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#DC2626" stroke-width="2">
              <path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"></path>
              <line x1="12" y1="9" x2="12" y2="13"></line>
              <line x1="12" y1="17" x2="12.01" y2="17"></line>
            </svg>
          </div>
          <h3 class="delete-heading">Hapus Aturan Lisensi?</h3>
          <p class="delete-msg">
            Apakah Anda yakin ingin menghapus aturan <strong>{{ selectedRule?.rule_name }}</strong>?
            Tindakan ini akan menghapus aturan dari database Microsoft SQL Server.
          </p>
          <div class="delete-actions">
            <button type="button" class="btn-secondary" @click="showDeleteModal = false">
              Batal
            </button>
            <button
              type="button"
              class="btn-danger-confirm"
              @click="executeDelete"
              :disabled="isDeleting"
            >
              <div v-if="isDeleting" class="spinner-sm"></div>
              <span v-else>Ya, Hapus Aturan</span>
            </button>
          </div>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import BaseAutocomplete, { type AutocompleteOption } from '~/components/base/BaseAutocomplete.vue'
import { useBackendBaseUrl } from '~/composables/useSizingApi'

interface LicenseRule {
  id: string
  rule_name: string
  vendor: string
  category: string
  metric_type: string
  pack_size: number
  min_cluster_node: number
  target_layer: string
  pricing_sku: string
  unit_price_idr: number | null
  description: string
  is_active: boolean
  created_at?: string
  updated_at?: string
}

interface PricingItemOption {
  id?: string
  sku: string
  name: string
  category?: string
  unit_price_idr: number | null
  unit: string
  environment: string
  source?: string
}

const props = defineProps<{
  isSidebarCollapsed?: boolean
}>()

const emit = defineEmits<{
  (e: 'toggle-sidebar'): void
}>()

const apiBase = useBackendBaseUrl()

// State
const rules = ref<LicenseRule[]>([])
const availablePricingList = ref<PricingItemOption[]>([])
const catalogComponentsList = ref<any[]>([])
const isLoading = ref(false)
const isSaving = ref(false)
const isDeleting = ref(false)

const searchQuery = ref('')
const categoryFilter = ref('all')
const metricFilter = ref('all')

const showModal = ref(false)
const isEditing = ref(false)
const selectedRule = ref<LicenseRule | null>(null)
const showDeleteModal = ref(false)
const selectedSkuOption = ref('')

const form = ref<Partial<LicenseRule>>({
  rule_name: '',
  vendor: '',
  category: 'apm',
  metric_type: 'per_5_pods',
  pack_size: 5,
  min_cluster_node: 1,
  target_layer: 'App Container',
  pricing_sku: '',
  unit_price_idr: null,
  description: '',
  is_active: true
})

// Visual Cards for Non-Tech Users
const metricCardsList = [
  {
    id: 'per_5_pods',
    title: 'Paket Pods (Microservices)',
    badge: 'e.g. AppDynamics APM',
    desc: 'Dihitung per kelipatan 5 Pods container (dibulatkan ke atas)',
    icon: '📦',
    defaultPack: 5,
    defaultMinNode: 1,
    suggestedSku: 'APPDYN_POD_PROD',
    defaultPrice: 8392364
  },
  {
    id: 'per_node_db',
    title: 'Per Node Database Agent',
    badge: 'e.g. AppDynamics DB Agent',
    desc: '1 agent monitoring per instance database aktif',
    icon: '🗄️',
    defaultPack: 1,
    defaultMinNode: 1,
    suggestedSku: 'APPDYN_DB_PROD',
    defaultPrice: 7514100
  },
  {
    id: 'per_shard',
    title: 'Per Shard Data (Cluster HA)',
    badge: 'e.g. Redis Enterprise',
    desc: 'Dihitung per shard data (Wajib min. 3 node cluster quorum)',
    icon: '⚡',
    defaultPack: 1,
    defaultMinNode: 3,
    suggestedSku: 'LIC_REDIS_PROD',
    defaultPrice: 125000000
  },
  {
    id: 'per_worker',
    title: 'Per Worker Node',
    badge: 'e.g. Red Hat OpenShift',
    desc: 'Per server worker (Master, Router & Bastion gratis)',
    icon: '🏢',
    defaultPack: 1,
    defaultMinNode: 1,
    suggestedSku: 'LIC_OCP_PROD',
    defaultPrice: 240000000
  },
  {
    id: 'per_core',
    title: 'Per Core Processor',
    badge: 'e.g. Oracle Enterprise RAC',
    desc: 'Dihitung per core CPU fisik / vCPU dengan Core Factor',
    icon: '🧠',
    defaultPack: 1,
    defaultMinNode: 1,
    suggestedSku: 'LIC_ORACLE_PROD',
    defaultPrice: null
  },
  {
    id: 'per_instance',
    title: 'Per Instance Server / VM',
    badge: 'e.g. NGINX Plus',
    desc: 'Dihitung per unit VM load balancer / reverse proxy',
    icon: '🖥️',
    defaultPack: 1,
    defaultMinNode: 1,
    suggestedSku: 'LIC_NGINX_PROD',
    defaultPrice: null
  }
]

// Computed
const activeRulesCount = computed(() => {
  return rules.value.filter(r => r.is_active).length
})

const filteredRules = computed(() => {
  return rules.value.filter(rule => {
    if (searchQuery.value) {
      const q = searchQuery.value.toLowerCase()
      const matchName = rule.rule_name.toLowerCase().includes(q)
      const matchVendor = rule.vendor.toLowerCase().includes(q)
      const matchSku = (rule.pricing_sku || '').toLowerCase().includes(q)
      const matchLayer = (rule.target_layer || '').toLowerCase().includes(q)
      const matchDesc = (rule.description || '').toLowerCase().includes(q)
      if (!matchName && !matchVendor && !matchSku && !matchLayer && !matchDesc) {
        return false
      }
    }

    if (categoryFilter.value !== 'all' && rule.category !== categoryFilter.value) {
      return false
    }

    if (metricFilter.value !== 'all' && rule.metric_type !== metricFilter.value) {
      return false
    }

    return true
  })
})

// Methods
const fetchRules = async () => {
  isLoading.value = true
  try {
    const res = await fetch(`${apiBase}/api/v1/license/rules`)
    if (res.ok) {
      const data = await res.json()
      rules.value = Array.isArray(data) ? data : []
    }
  } catch (err) {
    console.error('Failed to fetch license rules:', err)
  } finally {
    isLoading.value = false
  }
}

const fetchPricingItems = async () => {
  try {
    const res = await fetch(`${apiBase}/api/v1/pricing/items`)
    if (res.ok) {
      const data = await res.json()
      const items = Array.isArray(data.data) ? data.data : (Array.isArray(data) ? data : [])
      availablePricingList.value = items
    }
  } catch (err) {
    console.error('Failed to fetch pricing items:', err)
  }
}

const fetchCatalogComponents = async () => {
  try {
    const res = await fetch(`${apiBase}/api/v1/catalog/components`)
    if (res.ok) {
      const data = await res.json()
      const items = Array.isArray(data.data) ? data.data : (Array.isArray(data) ? data : [])
      catalogComponentsList.value = items
    }
  } catch (err) {
    console.error('Failed to fetch catalog components:', err)
  }
}



function getSkuIcon(sku: string): string {
  const s = (sku || '').toUpperCase()
  if (s.includes('APPDYN')) return '📈'
  if (s.includes('REDIS')) return '⚡'
  if (s.includes('OCP') || s.includes('OPENSHIFT')) return '🏢'
  if (s.includes('ORACLE')) return '🧠'
  if (s.includes('NGINX')) return '🌐'
  if (s.includes('KAFKA')) return '📬'
  if (s.includes('RUBRIK')) return '💾'
  if (s.includes('VCPU') || s.includes('BM_')) return '🖥️'
  if (s.includes('STG_')) return '💿'
  return '📦'
}

// Pricing SKU Autocomplete Options
const pricingSkuOptions = computed<AutocompleteOption[]>(() => {
  return availablePricingList.value.map(item => {
    const priceText = item.unit_price_idr != null
      ? formatCurrency(item.unit_price_idr)
      : 'Menunggu Konfirmasi STA'

    return {
      label: `[${item.sku}] ${item.name}`,
      value: item.sku,
      category: item.category ? item.category.toUpperCase() : 'SKU',
      badge: item.environment || (item.category === 'license' ? 'LISENSI RESMI' : 'KATALOG'),
      badgeType: item.category === 'license' ? 'primary' : 'default',
      description: `Tarif: ${priceText} / ${item.unit || 'Unit'} • Ref: ${item.source || 'Katalog Enterprise Architecture'}`,
      icon: getSkuIcon(item.sku),
      meta: item
    }
  })
})

// Software / Rule Name Autocomplete Options (Matching Catalog Components + Pricing)
const softwareNameOptions = computed<AutocompleteOption[]>(() => {
  const optionsMap = new Map<string, AutocompleteOption>()

  // From catalog components
  for (const comp of catalogComponentsList.value) {
    if (comp.component_name) {
      optionsMap.set(comp.component_name.toLowerCase(), {
        label: comp.component_name,
        value: comp.component_name,
        category: comp.layer_name || 'Komponen Arsitektur',
        badge: `Layer ${comp.layer_no || ''}`,
        description: comp.description || comp.notes || '',
        icon: '🏛️',
        meta: { type: 'component', ...comp }
      })
    }
    if (Array.isArray(comp.licenses)) {
      for (const lic of comp.licenses) {
        if (lic.product_name && !optionsMap.has(lic.product_name.toLowerCase())) {
          optionsMap.set(lic.product_name.toLowerCase(), {
            label: lic.product_name,
            value: lic.product_name,
            category: comp.layer_name || 'Lisensi Produk',
            badge: lic.vendor || 'Vendor',
            description: `SKU: ${lic.pricing_sku || 'Belum di-mapping'} • Metrik: ${formatMetricType(lic.metric)}`,
            icon: '📜',
            meta: { type: 'license', component: comp, license: lic }
          })
        }
      }
    }
  }

  // From pricing items
  for (const p of availablePricingList.value) {
    if (p.name && !optionsMap.has(p.name.toLowerCase())) {
      optionsMap.set(p.name.toLowerCase(), {
        label: p.name,
        value: p.name,
        category: p.category ? p.category.toUpperCase() : 'Katalog Pricing',
        badge: p.sku,
        description: `Tarif: ${p.unit_price_idr ? formatCurrency(p.unit_price_idr) : 'Belum ada tarif'} / ${p.unit || 'Unit'}`,
        icon: getSkuIcon(p.sku),
        meta: { type: 'pricing', item: p }
      })
    }
  }

  return Array.from(optionsMap.values())
})

// Check if currently typed SKU is matched with catalog
const matchedCatalogItem = computed(() => {
  if (!form.value.pricing_sku) return null
  return availablePricingList.value.find(
    p => p.sku.toLowerCase() === (form.value.pricing_sku || '').toLowerCase().trim()
  ) || null
})

// Check if currently typed rule name is matched with catalog component
const matchedCatalogComponent = computed(() => {
  if (!form.value.rule_name) return null
  const rn = form.value.rule_name.toLowerCase().trim()
  return catalogComponentsList.value.find(c =>
    c.component_name?.toLowerCase().includes(rn) ||
    rn.includes(c.component_name?.toLowerCase()) ||
    c.licenses?.some((l: any) => l.product_name?.toLowerCase().includes(rn))
  ) || null
})

function matchTargetLayer(layerName: string) {
  const ln = layerName.toLowerCase()
  if (ln.includes('container') || ln.includes('microservice') || ln.includes('app')) {
    form.value.target_layer = 'App Container'
  } else if (ln.includes('database') || ln.includes('db')) {
    form.value.target_layer = 'Core Database'
  } else if (ln.includes('cache') || ln.includes('redis')) {
    form.value.target_layer = 'Distributed Cache'
  } else if (ln.includes('ingress') || ln.includes('load balancer') || ln.includes('web')) {
    form.value.target_layer = 'Web Ingress'
  } else if (ln.includes('router') || ln.includes('bastion')) {
    form.value.target_layer = 'Container Router'
  } else if (ln.includes('message') || ln.includes('bus') || ln.includes('kafka') || ln.includes('streaming')) {
    form.value.target_layer = 'Message Bus'
  } else if (ln.includes('log') || ln.includes('analytics') || ln.includes('elk')) {
    form.value.target_layer = 'Log Analytics'
  }
}

function autoMatchRuleFromCatalog(item: any) {
  const sku = (item.sku || '').toUpperCase()
  const name = (item.name || '').toLowerCase()

  if (sku.includes('APPDYN') || name.includes('appdynamics')) {
    form.value.vendor = 'Cisco / AppDynamics'
    if (sku.includes('DB') || name.includes('database')) {
      form.value.category = 'database'
      form.value.target_layer = 'Core Database'
      form.value.metric_type = 'per_node_db'
      form.value.pack_size = 1
      form.value.min_cluster_node = 1
    } else {
      form.value.category = 'apm'
      form.value.target_layer = 'App Container'
      form.value.metric_type = 'per_5_pods'
      form.value.pack_size = 5
      form.value.min_cluster_node = 1
    }
  } else if (sku.includes('REDIS') || name.includes('redis')) {
    form.value.vendor = 'Redis Ltd'
    form.value.category = 'cache'
    form.value.target_layer = 'Distributed Cache'
    form.value.metric_type = 'per_shard'
    form.value.pack_size = 1
    form.value.min_cluster_node = 3
  } else if (sku.includes('OCP') || name.includes('openshift')) {
    form.value.vendor = 'Red Hat'
    form.value.category = 'platform'
    form.value.target_layer = 'App Container'
    form.value.metric_type = 'per_worker'
    form.value.pack_size = 1
    form.value.min_cluster_node = 1
  } else if (sku.includes('ORACLE') || name.includes('oracle')) {
    form.value.vendor = 'Oracle'
    form.value.category = 'database'
    form.value.target_layer = 'Core Database'
    form.value.metric_type = 'per_core'
    form.value.pack_size = 1
    form.value.min_cluster_node = 1
  } else if (sku.includes('NGINX') || name.includes('nginx') || name.includes('f5')) {
    form.value.vendor = 'F5'
    form.value.category = 'ingress'
    form.value.target_layer = 'Web Ingress'
    form.value.metric_type = 'per_instance'
    form.value.pack_size = 1
    form.value.min_cluster_node = 1
  } else if (sku.includes('KAFKA') || name.includes('confluent')) {
    form.value.vendor = 'Confluent'
    form.value.category = 'streaming'
    form.value.target_layer = 'Message Bus'
    form.value.metric_type = 'per_node'
    form.value.pack_size = 1
    form.value.min_cluster_node = 3
  }
}

function onPricingSkuSelected(opt: AutocompleteOption) {
  form.value.pricing_sku = opt.value
  const item = opt.meta
  if (!item) return

  if (item.unit_price_idr != null) {
    form.value.unit_price_idr = item.unit_price_idr
  }
  if (!form.value.rule_name) {
    form.value.rule_name = item.name
  }
  autoMatchRuleFromCatalog(item)
}

function onSoftwareNameSelected(opt: AutocompleteOption) {
  form.value.rule_name = opt.value
  const meta = opt.meta
  if (!meta) return

  if (meta.type === 'pricing' && meta.item) {
    onPricingSkuSelected({
      label: meta.item.name,
      value: meta.item.sku,
      meta: meta.item
    })
  } else if (meta.type === 'license') {
    if (meta.license.vendor) form.value.vendor = meta.license.vendor
    if (meta.license.pricing_sku) {
      form.value.pricing_sku = meta.license.pricing_sku
    }
    if (meta.component && meta.component.layer_name) {
      matchTargetLayer(meta.component.layer_name)
    }
  } else if (meta.type === 'component') {
    if (meta.layer_name) {
      matchTargetLayer(meta.layer_name)
    }
    const entLic = Array.isArray(meta.licenses) ? meta.licenses.find((l: any) => l.kind === 'enterprise') : null
    if (entLic) {
      if (entLic.vendor) form.value.vendor = entLic.vendor
      if (entLic.pricing_sku) {
        form.value.pricing_sku = entLic.pricing_sku
      }
    }
  }
}

function selectMetricCard(opt: any) {
  form.value.metric_type = opt.id
  form.value.pack_size = opt.defaultPack
  form.value.min_cluster_node = opt.defaultMinNode

  if (opt.suggestedSku && !form.value.pricing_sku) {
    form.value.pricing_sku = opt.suggestedSku
  }
}

// Watch pricing_sku to strictly sync price directly from master catalog
watch(
  () => form.value.pricing_sku,
  (newSku) => {
    if (!newSku) {
      form.value.unit_price_idr = null
      return
    }
    const found = availablePricingList.value.find(
      p => p.sku.toLowerCase() === newSku.toLowerCase().trim()
    )
    if (found) {
      form.value.unit_price_idr = found.unit_price_idr != null ? found.unit_price_idr : null
    } else {
      form.value.unit_price_idr = null
    }
  }
)

const resetFilters = () => {
  searchQuery.value = ''
  categoryFilter.value = 'all'
  metricFilter.value = 'all'
}

function getMetricExplanation(metric: string | undefined, packSize: number | undefined, minNode: number | undefined): string {
  switch (metric) {
    case 'per_5_pods':
      const ps = packSize || 5
      return `Dihitung per kelipatan ${ps} Pods aplikasi microservices. Jika kebutuhan sistem Anda 12 pods, sistem otomatis menghitung 12 ÷ ${ps} = 2.4 dibulatkan ke atas menjadi ${Math.ceil(12 / ps)} paket lisensi.`
    case 'per_node_db':
      return 'Dihitung per Node / Instance Database aktif. Setiap server database (Master maupun Standby / Replica) membutuhkan 1 Database Visibility Agent.'
    case 'per_shard':
      const mn = minNode || 3
      return `Dihitung berdasarkan jumlah Shard data memori. Ketentuan arsitektur High Availability (HA) Enterprise Architecture mewajibkan minimal ${mn} Node server untuk menjamin Sentinel / Quorum auto-failover cluster.`
    case 'per_worker':
      return 'Dihitung HANYA pada node Worker aktif yang menjalankan container aplikasi. Server Control Plane (Master Node), Infra Router, dan Bastion bersifat GRATIS (Free of Charge).'
    case 'per_core':
      return 'Dihitung berdasarkan jumlah Core Processor fisik atau vCPU pada server database (dengan Core Factor Intel x86 standar 0.5).'
    case 'per_instance':
      return 'Dihitung per instance software virtual machine (VM) load balancer atau web proxy aktif.'
    case 'per_node':
      return 'Dihitung secara sederhana per unit server fisik / virtual machine.'
    default:
      return 'Dihitung berdasarkan jumlah unit komponen yang digunakan pada arsitektur sistem.'
  }
}

const openAddModal = () => {
  isEditing.value = false
  selectedRule.value = null
  selectedSkuOption.value = ''
  form.value = {
    rule_name: '',
    vendor: '',
    category: 'apm',
    metric_type: 'per_5_pods',
    pack_size: 5,
    min_cluster_node: 1,
    target_layer: 'App Container',
    pricing_sku: '',
    unit_price_idr: null,
    description: '',
    is_active: true
  }
  showModal.value = true
}

const openEditModal = (rule: LicenseRule) => {
  isEditing.value = true
  selectedRule.value = rule
  selectedSkuOption.value = rule.pricing_sku || ''
  form.value = {
    ...rule
  }
  showModal.value = true
}

const saveRule = async () => {
  isSaving.value = true
  try {
    const payload = {
      rule_name: form.value.rule_name,
      vendor: form.value.vendor,
      category: form.value.category,
      metric_type: form.value.metric_type,
      pack_size: Number(form.value.pack_size) || 1,
      min_cluster_node: Number(form.value.min_cluster_node) || 1,
      target_layer: form.value.target_layer,
      pricing_sku: form.value.pricing_sku,
      unit_price_idr: form.value.unit_price_idr ? Number(form.value.unit_price_idr) : null,
      description: form.value.description,
      is_active: form.value.is_active
    }

    let url = `${apiBase}/api/v1/license/rules`
    let method = 'POST'

    if (isEditing.value && selectedRule.value) {
      url = `${apiBase}/api/v1/license/rules/${selectedRule.value.id}`
      method = 'PUT'
    }

    const res = await fetch(url, {
      method,
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload)
    })

    if (res.ok) {
      showModal.value = false
      await fetchRules()
    } else {
      const errData = await res.json()
      alert('Gagal menyimpan aturan: ' + (errData.error || res.statusText))
    }
  } catch (err) {
    console.error('Error saving rule:', err)
    alert('Terjadi kesalahan koneksi saat menyimpan aturan lisensi.')
  } finally {
    isSaving.value = false
  }
}

const confirmDelete = (rule: LicenseRule) => {
  selectedRule.value = rule
  showDeleteModal.value = true
}

const executeDelete = async () => {
  if (!selectedRule.value) return
  isDeleting.value = true
  try {
    const res = await fetch(`${apiBase}/api/v1/license/rules/${selectedRule.value.id}`, {
      method: 'DELETE'
    })
    if (res.ok) {
      showDeleteModal.value = false
      await fetchRules()
    } else {
      alert('Gagal menghapus aturan lisensi')
    }
  } catch (err) {
    console.error('Error deleting rule:', err)
  } finally {
    isDeleting.value = false
  }
}

// Formatters
const formatCurrency = (val: number | null | undefined): string => {
  if (!val) return 'Rp 0'
  return 'Rp ' + Math.round(val).toLocaleString('id-ID')
}

const formatCategory = (cat: string): string => {
  switch (cat) {
    case 'apm': return 'APM / Monitoring'
    case 'database': return 'Database'
    case 'cache': return 'Dist. Cache'
    case 'platform': return 'Platform / K8s'
    case 'ingress': return 'Web Ingress'
    case 'streaming': return 'Event Stream'
    default: return cat.toUpperCase()
  }
}

const formatMetricType = (metric: string): string => {
  switch (metric) {
    case 'per_5_pods': return 'Per 5 Pods (5-Pack)'
    case 'per_node_db': return 'Per DB Agent Node'
    case 'per_shard': return 'Per Shard Data'
    case 'per_worker': return 'Per Worker Node'
    case 'per_core': return 'Per Processor Core'
    case 'per_instance': return 'Per Instance / VM'
    case 'per_node': return 'Per Node Server'
    default: return metric
  }
}

onMounted(() => {
  fetchRules()
  fetchPricingItems()
  fetchCatalogComponents()
})
</script>

<style scoped>
.license-view-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
  overflow: hidden;
  background: #F8FAFC;
  color: #1E293B;
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
}

/* TOPBAR */
.license-topbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 24px;
  background: #FFFFFF;
  border-bottom: 1px solid #CBD5E1;
  flex-shrink: 0;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.03);
}

.topbar-left {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.title-with-toggle {
  display: flex;
  align-items: center;
  gap: 10px;
}

.topbar-sidebar-toggle {
  background: #F8FAFC;
  border: 1px solid #CBD5E1;
  border-radius: 6px;
  color: #334155;
  padding: 5px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;
}

.topbar-sidebar-toggle:hover {
  background: #F1F5F9;
  border-color: #94A3B8;
  color: #002060;
}

.page-title {
  margin: 0;
  font-size: 1.15rem;
  font-weight: 800;
  color: #002060;
  letter-spacing: -0.01em;
}

.page-subtitle {
  margin: 0;
  font-size: 0.78rem;
  color: #64748B;
}

.btn-primary-action {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  background: #002060;
  border: 1px solid transparent;
  color: #FFFFFF;
  font-size: 0.80rem;
  font-weight: 700;
  padding: 8px 16px;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s ease;
  box-shadow: 0 1px 3px rgba(0, 32, 96, 0.2);
}

.btn-primary-action:hover {
  background: #003399;
  box-shadow: 0 4px 10px rgba(0, 32, 96, 0.25);
  transform: translateY(-1px);
}

/* CONTENT SCROLL */
.license-content-scroll {
  flex: 1;
  overflow-y: auto;
  padding: 20px 24px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

/* KPI GRID */
.license-kpi-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 14px;
}

.kpi-card {
  background: #FFFFFF;
  border: 1px solid #CBD5E1;
  border-radius: 8px;
  padding: 14px 16px;
  display: flex;
  align-items: center;
  gap: 14px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
  transition: all 0.2s ease;
}

.kpi-card:hover {
  border-color: #002060;
  box-shadow: 0 4px 12px rgba(0, 32, 96, 0.08);
}

.highlight-card {
  border-color: #BAE6FD;
  background: linear-gradient(135deg, #F0F9FF 0%, #E0F2FE 100%);
}

.kpi-icon-wrap {
  width: 42px;
  height: 42px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.icon-blue { background: #EFF6FF; color: #1D4ED8; }
.icon-cyan { background: #E0F2FE; color: #0284C7; }
.icon-gold { background: #FFFBEB; color: #D97706; }
.icon-purple { background: #FAF5FF; color: #7E22CE; }

.kpi-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.kpi-label {
  font-size: 0.70rem;
  color: #64748B;
  text-transform: uppercase;
  font-weight: 700;
  letter-spacing: 0.04em;
}

.kpi-value-row {
  display: flex;
  flex-direction: column;
}

.kpi-num {
  font-size: 1.15rem;
  font-weight: 800;
  color: #002060;
}

.kpi-sub {
  font-size: 0.72rem;
  color: #64748B;
}

.kpi-tag-sub {
  font-size: 0.70rem;
  color: #0369A1;
  font-weight: 700;
}

/* FILTER BAR */
.license-filter-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  background: #FFFFFF;
  border: 1px solid #CBD5E1;
  border-radius: 8px;
  padding: 10px 16px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
}

.search-input-wrap {
  position: relative;
  flex: 1;
  max-width: 480px;
}

.search-icon {
  position: absolute;
  left: 12px;
  top: 50%;
  transform: translateY(-50%);
  color: #94A3B8;
}

.search-input {
  width: 100%;
  background: #F8FAFC;
  border: 1px solid #CBD5E1;
  border-radius: 6px;
  padding: 7px 32px 7px 36px;
  font-size: 0.82rem;
  color: #0F172A;
  outline: none;
  transition: all 0.2s ease;
}

.search-input:focus {
  background: #FFFFFF;
  border-color: #002060;
  box-shadow: 0 0 0 2px rgba(0, 32, 96, 0.1);
}

.clear-search-btn {
  position: absolute;
  right: 10px;
  top: 50%;
  transform: translateY(-50%);
  background: transparent;
  border: none;
  color: #94A3B8;
  font-size: 1rem;
  cursor: pointer;
}

.filter-controls {
  display: flex;
  align-items: center;
  gap: 12px;
}

.filter-group {
  display: flex;
  align-items: center;
  gap: 6px;
}

.filter-lbl {
  font-size: 0.75rem;
  color: #475569;
  font-weight: 600;
}

.filter-select {
  background: #F8FAFC;
  border: 1px solid #CBD5E1;
  border-radius: 6px;
  padding: 6px 10px;
  font-size: 0.78rem;
  color: #1E293B;
  outline: none;
  cursor: pointer;
  font-weight: 600;
}

.filter-select:focus {
  border-color: #002060;
  background: #FFFFFF;
}

.reset-filters-btn {
  background: #F1F5F9;
  border: 1px solid #CBD5E1;
  color: #475569;
  font-size: 0.75rem;
  font-weight: 600;
  padding: 6px 10px;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.reset-filters-btn:hover {
  background: #E2E8F0;
  color: #0F172A;
}

/* RULES TABLE CARD */
.rules-table-card {
  background: #FFFFFF;
  border: 1px solid #CBD5E1;
  border-radius: 8px;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
}

.table-card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 18px;
  background: #F8FAFC;
  border-bottom: 1px solid #E2E8F0;
}

.table-title-row {
  display: flex;
  align-items: center;
  gap: 10px;
}

.table-title {
  margin: 0;
  font-size: 0.92rem;
  font-weight: 800;
  color: #002060;
}

.table-badge {
  background: #E0E7FF;
  color: #3730A3;
  font-size: 0.68rem;
  font-weight: 700;
  padding: 2px 8px;
  border-radius: 12px;
  letter-spacing: 0.04em;
}

.btn-refresh {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  background: #FFFFFF;
  border: 1px solid #CBD5E1;
  color: #334155;
  font-size: 0.75rem;
  font-weight: 600;
  padding: 5px 10px;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-refresh:hover:not(:disabled) {
  background: #F8FAFC;
  color: #002060;
}

.spinning {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.table-responsive {
  overflow-x: auto;
}

.sketch-table {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
  font-size: 0.8rem;
}

.sketch-table th {
  background: #002060;
  color: #FFFFFF;
  font-weight: 700;
  text-transform: uppercase;
  font-size: 0.68rem;
  letter-spacing: 0.05em;
  padding: 10px 16px;
  border-bottom: 1px solid #00173D;
  white-space: nowrap;
}

.sketch-table td {
  padding: 12px 16px;
  border-bottom: 1px solid #E2E8F0;
  color: #1E293B;
  vertical-align: middle;
}

.sketch-table tbody tr:hover td {
  background: #F8FAFC;
}

.inactive-row td {
  opacity: 0.6;
  background: #F8FAFC;
}

/* CELL CONTENT STYLES */
.rule-name-wrap {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.rule-title {
  font-size: 0.86rem;
  font-weight: 700;
  color: #0F172A;
}

.rule-meta-badges {
  display: flex;
  align-items: center;
  gap: 6px;
}

.vendor-badge {
  background: #EFF6FF;
  color: #1D4ED8;
  border: 1px solid #BFDBFE;
  font-size: 0.68rem;
  font-weight: 700;
  padding: 2px 7px;
  border-radius: 4px;
}

.status-pill {
  font-size: 0.65rem;
  font-weight: 700;
  padding: 2px 6px;
  border-radius: 4px;
}

.status-active {
  background: #DCFCE7;
  color: #15803D;
  border: 1px solid #BBF7D0;
}

.status-inactive {
  background: #FEE2E2;
  color: #B91C1C;
  border: 1px solid #FECACA;
}

.rule-desc {
  margin: 0;
  font-size: 0.72rem;
  color: #64748B;
  line-height: 1.35;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.layer-meta {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.cat-pill {
  display: inline-block;
  align-self: flex-start;
  font-size: 0.68rem;
  font-weight: 700;
  padding: 2px 7px;
  border-radius: 4px;
}

.cat-apm { background: #E0F2FE; color: #0369A1; border: 1px solid #BAE6FD; }
.cat-database { background: #FEF3C7; color: #92400E; border: 1px solid #FDE68A; }
.cat-cache { background: #FEE2E2; color: #B91C1C; border: 1px solid #FECACA; }
.cat-platform { background: #F3E8FF; color: #6B21A8; border: 1px solid #E9D5FF; }
.cat-ingress { background: #DCFCE7; color: #15803D; border: 1px solid #BBF7D0; }
.cat-streaming { background: #FCE7F3; color: #9D174D; border: 1px solid #FBCFE8; }
.cat-other { background: #F1F5F9; color: #475569; border: 1px solid #CBD5E1; }

.target-layer-name {
  font-size: 0.74rem;
  color: #334155;
  font-weight: 600;
}

.metric-info-block {
  display: flex;
  flex-direction: column;
  gap: 3px;
}

.metric-type-pill {
  font-family: monospace;
  font-size: 0.74rem;
  color: #002060;
  font-weight: 700;
  background: #F1F5F9;
  border: 1px solid #CBD5E1;
  padding: 2px 6px;
  border-radius: 4px;
  align-self: flex-start;
}

.pack-size-row {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 0.70rem;
}

.pack-size-lbl {
  color: #64748B;
}

.pack-size-val {
  color: #0284C7;
  font-weight: 700;
}

.cluster-quorum-badge {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  background: #FFFBEB;
  color: #B45309;
  border: 1px solid #FCD34D;
  padding: 4px 8px;
  border-radius: 5px;
  font-size: 0.72rem;
  font-weight: 700;
}

.cluster-single-badge {
  color: #64748B;
  font-size: 0.74rem;
}

.pricing-info-cell {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.sku-code {
  font-family: monospace;
  font-size: 0.68rem;
  color: #64748B;
  font-weight: 600;
}

.price-val {
  font-size: 0.86rem;
  font-weight: 800;
  color: #002060;
}

.price-sub {
  font-size: 0.70rem;
  color: #64748B;
  font-weight: normal;
}

.badge-sta {
  background: #FEF3C7;
  color: #92400E;
  font-size: 0.68rem;
  font-weight: 600;
  padding: 2px 6px;
  border-radius: 4px;
  border: 1px solid #FDE68A;
}

.row-actions {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 6px;
}

.btn-action-icon {
  background: #F8FAFC;
  border: 1px solid #CBD5E1;
  color: #475569;
  padding: 5px;
  border-radius: 5px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;
}

.btn-action-icon:hover {
  background: #FFFFFF;
  color: #002060;
}

.btn-edit:hover {
  background: #EFF6FF;
  border-color: #3B82F6;
  color: #1D4ED8;
}

.btn-delete:hover {
  background: #FEE2E2;
  border-color: #EF4444;
  color: #B91C1C;
}

/* EMPTY / LOADING STATES */
.table-loading-state,
.table-empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 48px 24px;
  text-align: center;
  color: #64748B;
}

.spinner-circle {
  width: 32px;
  height: 32px;
  border: 3px solid #E2E8F0;
  border-top-color: #002060;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
  margin-bottom: 12px;
}

.empty-icon-wrap {
  margin-bottom: 12px;
}

.btn-sm {
  font-size: 0.75rem;
  padding: 6px 14px;
  margin-top: 10px;
}

/* ============================================================
   MODAL STYLES (EXECUTIVE & NON-TECH FRIENDLY)
============================================================ */
.sketch-modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(15, 23, 42, 0.55);
  backdrop-filter: blur(5px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 20px;
}

.rule-modal-card-friendly {
  width: 100%;
  max-width: 820px;
  max-height: 90vh;
  background: #FFFFFF;
  border: 1px solid #CBD5E1;
  border-radius: 12px;
  box-shadow: 0 24px 50px -10px rgba(0, 32, 96, 0.25);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.modal-card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 24px;
  background: #F8FAFC;
  border-bottom: 1px solid #E2E8F0;
  flex-shrink: 0;
}

.modal-title-wrap {
  display: flex;
  align-items: center;
  gap: 12px;
}

.modal-icon-shield {
  width: 40px;
  height: 40px;
  background: #E0E7FF;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.modal-heading {
  margin: 0;
  font-size: 1.05rem;
  font-weight: 800;
  color: #002060;
}

.modal-subheading {
  margin: 2px 0 0 0;
  font-size: 0.74rem;
  color: #64748B;
}

.modal-close-btn {
  background: transparent;
  border: none;
  color: #94A3B8;
  font-size: 1.5rem;
  line-height: 1;
  cursor: pointer;
  padding: 4px;
}

.modal-close-btn:hover {
  color: #1E293B;
}

.modal-form-friendly {
  padding: 20px 24px;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 16px;
  background: #F8FAFC;
}

/* SECTION CARDS INSIDE FORM */
.form-section-card {
  background: #FFFFFF;
  border: 1px solid #E2E8F0;
  border-radius: 10px;
  padding: 16px 20px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.03);
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.section-card-title {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 4px;
  padding-bottom: 8px;
  border-bottom: 1px solid #F1F5F9;
}

.section-step-num {
  width: 24px;
  height: 24px;
  background: #002060;
  color: #FFFFFF;
  border-radius: 50%;
  font-size: 0.72rem;
  font-weight: 800;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.section-title-text {
  margin: 0;
  font-size: 0.88rem;
  font-weight: 800;
  color: #002060;
}

.section-desc-text {
  margin: 0;
  font-size: 0.72rem;
  color: #64748B;
}

.section-grid-2 {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 14px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.form-label {
  font-size: 0.74rem;
  font-weight: 700;
  color: #1E293B;
}

.req {
  color: #DC2626;
}

.field-hint {
  font-size: 0.68rem;
  color: #64748B;
}

.form-input,
.form-select,
.form-textarea {
  background: #FFFFFF;
  border: 1px solid #CBD5E1;
  border-radius: 6px;
  padding: 8px 12px;
  font-size: 0.82rem;
  color: #0F172A;
  outline: none;
  transition: all 0.2s ease;
}

.form-input:focus,
.form-select:focus,
.form-textarea:focus {
  border-color: #002060;
  box-shadow: 0 0 0 2px rgba(0, 32, 96, 0.1);
}

.font-bold {
  font-weight: 700;
}

/* METRIC CARDS SELECTOR GRID */
.metric-cards-selector-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 10px;
}

.metric-choice-card {
  background: #F8FAFC;
  border: 2px solid #E2E8F0;
  border-radius: 8px;
  padding: 10px 12px;
  cursor: pointer;
  display: flex;
  flex-direction: column;
  gap: 4px;
  transition: all 0.2s ease;
  position: relative;
}

.metric-choice-card:hover {
  background: #FFFFFF;
  border-color: #94A3B8;
  transform: translateY(-1px);
}

.metric-choice-card.selected {
  background: #FFFFFF;
  border-color: #002060;
  box-shadow: 0 4px 12px rgba(0, 32, 96, 0.12);
}

.choice-card-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.choice-icon {
  font-size: 1.25rem;
}

.choice-check-badge {
  width: 18px;
  height: 18px;
  background: #002060;
  color: #FFFFFF;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.choice-title {
  margin: 0;
  font-size: 0.78rem;
  font-weight: 800;
  color: #0F172A;
}

.choice-sub-tag {
  font-size: 0.65rem;
  font-weight: 700;
  color: #0369A1;
}

.choice-desc {
  margin: 0;
  font-size: 0.68rem;
  color: #64748B;
  line-height: 1.35;
}

/* HUMAN EXPLANATION CALLOUT */
.metric-human-guide-box {
  background: #EFF6FF;
  border: 1px solid #BFDBFE;
  border-radius: 8px;
  padding: 10px 14px;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.guide-header {
  display: flex;
  align-items: center;
  gap: 6px;
}

.guide-lamp-icon {
  font-size: 0.95rem;
}

.guide-heading {
  font-size: 0.74rem;
  font-weight: 800;
  color: #1E3A8A;
}

.guide-body-text {
  margin: 0;
  font-size: 0.74rem;
  color: #1E40AF;
  line-height: 1.45;
}

/* CONTEXTUAL ADJUSTERS */
.contextual-adjuster-panel {
  background: #F8FAFC;
  border: 1px solid #CBD5E1;
  border-radius: 8px;
  padding: 10px 14px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 14px;
}

.contextual-adjuster-panel.alert-quorum {
  background: #FFFBEB;
  border-color: #FDE68A;
}

.adjuster-icon {
  font-size: 1.35rem;
  flex-shrink: 0;
}

.adjuster-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.adjuster-title {
  font-size: 0.75rem;
  font-weight: 800;
  color: #1E293B;
}

.adjuster-hint {
  margin: 0;
  font-size: 0.70rem;
  color: #64748B;
}

.adjuster-input-wrap {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-shrink: 0;
}

.adjuster-num-input {
  width: 60px;
  padding: 6px 8px;
  background: #FFFFFF;
  border: 1px solid #002060;
  border-radius: 6px;
  font-size: 0.85rem;
  font-weight: 800;
  text-align: center;
  color: #002060;
}

.adjuster-unit {
  font-size: 0.72rem;
  font-weight: 700;
  color: #475569;
}

/* CATALOG PRICE READONLY DISPLAY (STRICT LOCKED TO MASTER CATALOG) */
.catalog-price-readonly-box {
  background: #F8FAFC;
  border: 1.5px solid #CBD5E1;
  border-radius: 8px;
  padding: 10px 14px;
  display: flex;
  flex-direction: column;
  gap: 8px;
  transition: all 0.2s ease;
}

.catalog-price-readonly-box.has-price {
  background: #F8FAFC;
  border-color: #94A3B8;
}

.catalog-price-readonly-box.pending-price {
  background: #FFFBEB;
  border-color: #FDE68A;
}

.price-val-row {
  display: flex;
  align-items: baseline;
  gap: 6px;
}

.price-currency-label {
  font-size: 0.95rem;
  font-weight: 800;
  color: #64748B;
}

.price-large-number {
  font-size: 1.4rem;
  font-weight: 800;
  color: #002060;
  letter-spacing: -0.02em;
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
}

.price-unit-tag {
  font-size: 0.78rem;
  font-weight: 600;
  color: #64748B;
}

.price-sync-status {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 0.72rem;
  font-weight: 600;
  padding: 4px 9px;
  border-radius: 6px;
  line-height: 1.3;
}

.price-sync-status.status-ok {
  background: #DCFCE7;
  color: #15803D;
  border: 1px solid #BBF7D0;
}

.price-sync-status.status-pending {
  background: #FEF3C7;
  color: #B45309;
  border: 1px solid #FDE68A;
}

/* FORM LABEL ROW & MATCH TAGS */
.form-label-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.label-match-tag {
  font-size: 0.65rem;
  font-weight: 700;
  padding: 2px 7px;
  border-radius: 4px;
  letter-spacing: 0.02em;
}

.label-match-tag.match-yes {
  background: #DCFCE7;
  color: #15803D;
  border: 1px solid #BBF7D0;
}

.label-match-tag.match-custom {
  background: #FEF3C7;
  color: #B45309;
  border: 1px solid #FDE68A;
}

.label-match-tag.match-locked {
  background: #EFF6FF;
  color: #1D4ED8;
  border: 1px solid #BFDBFE;
}

/* SKU MATCHED CARD (SUCCESS) */
.sku-matched-card {
  background: #F0FDF4;
  border: 1px solid #BBF7D0;
  border-radius: 8px;
  padding: 8px 12px;
  margin-top: 8px;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.matched-top-line {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.catalog-verified-badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-size: 0.66rem;
  font-weight: 700;
  color: #166534;
  background: #DCFCE7;
  padding: 1px 7px;
  border-radius: 9999px;
}

.matched-env-badge {
  font-size: 0.65rem;
  font-weight: 600;
  color: #002060;
  background: #E0E7FF;
  padding: 1px 6px;
  border-radius: 4px;
}

.matched-info-grid {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
}

.matched-info-left {
  display: flex;
  flex-direction: column;
}

.matched-title {
  font-size: 0.78rem;
  font-weight: 700;
  color: #0F172A;
}

.matched-unit-text {
  font-size: 0.68rem;
  color: #64748B;
}

.matched-price-tag {
  font-size: 0.85rem;
  font-weight: 800;
  color: #002060;
  white-space: nowrap;
}

.matched-source-note {
  font-size: 0.65rem;
  color: #64748B;
  margin: 0;
  border-top: 1px dashed #DCFCE7;
  padding-top: 4px;
}

/* SKU CUSTOM CARD (INFO) */
.sku-custom-card {
  background: #FFFBEB;
  border: 1px solid #FDE68A;
  border-radius: 8px;
  padding: 8px 12px;
  margin-top: 8px;
  display: flex;
  align-items: flex-start;
  gap: 8px;
}

.custom-sku-badge {
  font-size: 0.66rem;
  font-weight: 700;
  color: #92400E;
  background: #FEF3C7;
  padding: 2px 6px;
  border-radius: 4px;
  white-space: nowrap;
}

.custom-sku-text {
  font-size: 0.70rem;
  color: #78350F;
  line-height: 1.35;
}

/* STATUS TOGGLE CARD */
.status-toggle-card {
  background: #F1F5F9;
  border: 1px solid #CBD5E1;
  border-radius: 8px;
  padding: 10px 14px;
  transition: all 0.2s ease;
}

.status-toggle-card.is-active-card {
  background: #F0FDF4;
  border-color: #86EFAC;
}

.toggle-checkbox-row {
  display: flex;
  align-items: center;
  gap: 12px;
  cursor: pointer;
}

.modern-checkbox {
  width: 20px;
  height: 20px;
  accent-color: #002060;
  cursor: pointer;
}

.toggle-texts {
  display: flex;
  flex-direction: column;
}

.toggle-heading {
  font-size: 0.78rem;
  font-weight: 800;
  color: #0F172A;
}

.toggle-sub {
  font-size: 0.70rem;
  color: #64748B;
}

/* MODAL FORM ACTIONS */
.modal-form-actions-friendly {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 12px;
  padding-top: 8px;
}

.btn-secondary-friendly {
  background: #FFFFFF;
  border: 1px solid #CBD5E1;
  color: #475569;
  font-size: 0.82rem;
  font-weight: 700;
  padding: 9px 18px;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-secondary-friendly:hover {
  background: #F1F5F9;
  color: #0F172A;
}

.btn-primary-friendly {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  background: #002060;
  border: none;
  color: #FFFFFF;
  font-size: 0.82rem;
  font-weight: 800;
  padding: 9px 22px;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s ease;
  box-shadow: 0 2px 6px rgba(0, 32, 96, 0.25);
}

.btn-primary-friendly:hover:not(:disabled) {
  background: #003399;
  box-shadow: 0 4px 14px rgba(0, 32, 96, 0.35);
  transform: translateY(-1px);
}

.spinner-sm-white {
  width: 15px;
  height: 15px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: #FFFFFF;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

/* DELETE CONFIRM CARD */
.delete-confirm-card {
  width: 90%;
  max-width: 440px;
  padding: 24px;
  text-align: center;
}

.delete-icon-wrap {
  width: 48px;
  height: 48px;
  background: #FEE2E2;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 14px auto;
}

.delete-heading {
  margin: 0 0 8px 0;
  font-size: 1.05rem;
  font-weight: 800;
  color: #0F172A;
}

.delete-msg {
  font-size: 0.82rem;
  color: #64748B;
  line-height: 1.45;
  margin: 0 0 20px 0;
}

.delete-actions {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
}

.btn-danger-confirm {
  background: #DC2626;
  border: 1px solid #B91C1C;
  color: #FFFFFF;
  font-size: 0.82rem;
  font-weight: 700;
  padding: 8px 18px;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-danger-confirm:hover:not(:disabled) {
  background: #B91C1C;
}

/* TRANSITIONS */
.modal-fade-enter-active,
.modal-fade-leave-active {
  transition: opacity 0.2s ease;
}

.modal-fade-enter-from,
.modal-fade-leave-to {
  opacity: 0;
}
</style>
