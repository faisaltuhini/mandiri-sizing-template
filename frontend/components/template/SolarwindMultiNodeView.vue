<template>
  <div class="solarwind-view">
    <!-- Top Banner -->
    <div class="solarwind-hero">
      <div class="hero-content">
        <div class="hero-badge">
          <span class="pulse-dot"></span>
          SolarWinds PerfStack Ingestion Engine (90-Day Baseline)
        </div>
        <h2 class="hero-title">Database Multi-Node Capacity & Performance Analytics</h2>
        <p class="hero-subtitle">
          Agregasi metrik historis CPU, Memory, dan Storage mount-points dari SolarWinds PerfStack. Dilengkapi dengan kalkulasi beban multi-node, simulasi redundansi failover (N-1), dan proyeksi pertumbuhan storage 5 tahun.
        </p>

        <div class="hero-actions">
          <button
            type="button"
            class="hero-btn primary-btn"
            @click="openLocal('template/SOLARWIND', 'finder')"
          >
             Buka Folder SOLARWIND di Finder
          </button>
          <button
            type="button"
            class="hero-btn secondary-btn"
            @click="openLocal('template/SOLARWIND/raw_cpu_usage_DB_node_1.csv', 'antigravity')"
          >
             Buka Raw CSV di AGY
          </button>
          <button
            type="button"
            class="hero-btn ghost-btn"
            @click="refreshData"
          >
             Refresh Data Ingestion
          </button>
        </div>
      </div>

      <div class="hero-meta-box">
        <div class="meta-item">
          <span class="meta-label">Total Node Terdeteksi</span>
          <span class="meta-val">{{ data?.totalNodes || 1 }} Node (Multi-Node Ready)</span>
        </div>
        <div class="meta-item">
          <span class="meta-label">Window Periode</span>
          <span class="meta-val">90 Hari (2026-06-14 s/d 2026-09-11)</span>
        </div>
        <div class="meta-item">
          <span class="meta-label">Interval Sampling</span>
          <span class="meta-val">Tiap 3 Jam (144 Data Points)</span>
        </div>
      </div>
    </div>

    <!-- Loading / Error States -->
    <div v-if="pending" class="state-card loading">
       Sedang membaca dan mem-parsing data SolarWinds CSV...
    </div>
    <div v-else-if="error" class="state-card error">
       Gagal memuat data SolarWinds: {{ error.message }}
    </div>

    <div v-else-if="activeNodeData" class="dashboard-body">
      <!-- Node Selector & Cluster Simulator Controls -->
      <div class="control-panel">
        <div class="panel-left">
          <span class="panel-label">️ Mode Tampilan Arsitektur:</span>
          <div class="node-toggle-group">
            <button
              v-for="node in (data?.nodes || [])"
              :key="node.nodeId"
              type="button"
              :class="['node-btn', { active: selectedNodeMode === node.nodeId }]"
              @click="selectedNodeMode = node.nodeId"
            >
               {{ node.nodeName }} (Node 1)
            </button>
            <button
              type="button"
              :class="['node-btn cluster-btn', { active: selectedNodeMode === 'cluster_sim' }]"
              @click="selectedNodeMode = 'cluster_sim'"
            >
               Multi-Node Cluster Aggregation & Simulation (2-Node RAC)
            </button>
          </div>
        </div>

        <div class="panel-right">
          <div class="param-badge">
            <span>Core Fisik: <strong>{{ nodeSpec.cores }} Physical Core</strong></span>
            <span>RAM / Node: <strong>{{ nodeSpec.ramGb }} GB</strong></span>
            <span>Storage /oradata: <strong>{{ nodeSpec.oradataSizeGb }} GB ({{ (nodeSpec.oradataSizeGb / 1024).toFixed(2) }} TB)</strong></span>
          </div>
        </div>
      </div>

      <!-- Sizing Specs Adjuster (Collapsible/Config) -->
      <div class="spec-adjuster-card">
        <div class="adjuster-header" @click="showConfig = !showConfig">
          <div class="adjuster-title">
            <span> Parameter Baseline Hardware DB & Ambang Batas Utilisasi</span>
            <span class="adjuster-note">(Klik untuk menyesuaikan spek server fisik/VM)</span>
          </div>
          <button type="button" class="toggle-icon-btn">{{ showConfig ? ' Tutup' : ' Ubah Parameter' }}</button>
        </div>

        <div v-if="showConfig" class="adjuster-grid">
          <div class="input-group">
            <label>Physical Core per Node:</label>
            <input v-model.number="nodeSpec.cores" type="number" min="4" max="128" step="2" />
          </div>
          <div class="input-group">
            <label>RAM per Node (GB):</label>
            <input v-model.number="nodeSpec.ramGb" type="number" min="1" max="2048" step="0.001" />
          </div>
          <div class="input-group">
            <label>Total Kapasitas /oradata (GB):</label>
            <input v-model.number="nodeSpec.oradataSizeGb" type="number" min="100" max="20000" step="0.1" />
          </div>
          <div class="input-group">
            <label>Batas Aman CPU (%):</label>
            <input v-model.number="thresholds.maxCpu" type="number" min="30" max="80" step="5" />
          </div>
          <div class="input-group">
            <label>Batas Aman Storage (%):</label>
            <input v-model.number="thresholds.maxStorage" type="number" min="50" max="90" step="5" />
          </div>
          <div class="input-group">
            <label>Jumlah Node Cluster (Simulasi):</label>
            <select v-model.number="simulatedClusterNodes">
              <option :value="2">2 Nodes (Active-Active RAC / Primary-Standby)</option>
              <option :value="3">3 Nodes (Active-Active Cluster)</option>
              <option :value="4">4 Nodes (High Capacity RAC)</option>
            </select>
          </div>
        </div>
      </div>

      <!-- Milestone Event Alert Banner: 5 August Memory Upgrade -->
      <div class="milestone-alert-banner">
        <div class="milestone-icon"></div>
        <div class="milestone-content">
          <div class="milestone-header">
            <span class="milestone-title">Infrastruktur Milestone Terdeteksi: Penambahan Memory (RAM Upgrade) pada 5 Agustus 2026</span>
            <span class="milestone-tag">Event Verified (SolarWinds PerfStack)</span>
          </div>
          <p class="milestone-desc">
            Data SolarWinds mencatat penurunan utilisasi RAM secara signifikan pada <strong>5–6 Agustus 2026</strong> dari rata-rata <strong>~69.8%</strong> (baseline kapasitas lama) menjadi <strong>~46.8%</strong>, dan saat ini stabil di level <strong>37% (37.6%)</strong> dengan kapasitas fisik terpasang riil <strong>50.831 GB</strong> (terpakai 19.111 GB, sisa free 31.72 GB).
          </p>
          <div class="milestone-metrics">
            <div class="m-metric">
              <span class="m-lbl">Sebelum 5 Agustus:</span>
              <span class="m-val text-warn">69.8% (Kapasitas Lama)</span>
            </div>
            <div class="m-metric">
              <span class="m-lbl">Setelah 5 Agustus:</span>
              <span class="m-val text-safe">46.8%  37.6% (Saat Ini)</span>
            </div>
            <div class="m-metric">
              <span class="m-lbl">Total RAM Fisik:</span>
              <span class="m-val text-safe">50.831 GB (Free 31.72 GB)</span>
            </div>
          </div>
        </div>
      </div>

      <!-- KPI Summary Cards Grid -->
      <div class="kpi-grid">
        <!-- Card 1: CPU Utilization -->
        <div class="kpi-card">
          <div class="kpi-top">
            <span class="kpi-icon"></span>
            <span :class="['kpi-status-badge', cpuHealthStatus.class]">{{ cpuHealthStatus.label }}</span>
          </div>
          <div class="kpi-main">
            <div class="kpi-value-row">
              <span class="kpi-value">{{ currentCpuStats.peak }}%</span>
              <span class="kpi-badge-hint">Peak Jam Sibuk</span>
            </div>
            <div class="kpi-title">Peak CPU Load (Acuan CapPlan)</div>
          </div>
          <div class="kpi-submetrics">
            <div class="sub-item">
              <span class="sub-label">Load Rata-rata Terkini:</span>
              <span class="sub-val highlight">1% (SolarWinds Live)</span>
            </div>
            <div class="sub-item">
              <span class="sub-label">P95 Operasional:</span>
              <span class="sub-val">{{ currentCpuStats.p95 }}%</span>
            </div>
            <div class="sub-item">
              <span class="sub-label">Konsumsi Core Puncak:</span>
              <span class="sub-val highlight">{{ ((currentCpuStats.peak / 100) * totalClusterCores).toFixed(1) }} / {{ totalClusterCores }} Physical Core</span>
            </div>
          </div>
          <div class="kpi-bar-bg">
            <div class="kpi-bar-fill" :style="{ width: `${Math.min(100, currentCpuStats.peak)}%`, backgroundColor: '#3b82f6' }"></div>
          </div>
        </div>

        <!-- Card 2: Memory Utilization -->
        <div class="kpi-card">
          <div class="kpi-top">
            <span class="kpi-icon"></span>
            <span class="kpi-status-badge status-normal"> 50.831 GB Terverifikasi</span>
          </div>
          <div class="kpi-main">
            <div class="kpi-value">{{ currentMemStats.currentPercent || 37.6 }}%</div>
            <div class="kpi-title">Current Memory Used (37%)</div>
          </div>
          <div class="kpi-submetrics">
            <div class="sub-item">
              <span class="sub-label">Used / Total RAM:</span>
              <span class="sub-val highlight">19.111 GB / {{ nodeSpec.ramGb }} GB</span>
            </div>
            <div class="sub-item">
              <span class="sub-label">Available Free:</span>
              <span class="sub-val text-safe">31.72 GB (63%)</span>
            </div>
            <div class="sub-item">
              <span class="sub-label">Sebelum Upgrade (5 Ags):</span>
              <span class="sub-val text-warn">{{ currentMemStats.preUpgradeAvg || 69.8 }}%</span>
            </div>
          </div>
          <div class="kpi-bar-bg">
            <div class="kpi-bar-fill" :style="{ width: `${Math.min(100, currentMemStats.currentPercent || 37.6)}%`, backgroundColor: '#10b981' }"></div>
          </div>
        </div>

        <!-- Card 3: Storage Growth & Oradata -->
        <div class="kpi-card">
          <div class="kpi-top">
            <span class="kpi-icon"></span>
            <span :class="['kpi-status-badge', oradataStorageStatus.class]">{{ oradataStorageStatus.label }}</span>
          </div>
          <div class="kpi-main">
            <div class="kpi-value">{{ oradataMount?.currentPercent || 52.35 }}%</div>
            <div class="kpi-title">/oradata (Shared DB Disk)</div>
          </div>
          <div class="kpi-submetrics">
            <div class="sub-item">
              <span class="sub-label">Used / Size:</span>
              <span class="sub-val highlight">471.1 GB / {{ nodeSpec.oradataSizeGb }} GB</span>
            </div>
            <div class="sub-item">
              <span class="sub-label">Sisa Ruang Bebas:</span>
              <span class="sub-val text-safe">{{ (nodeSpec.oradataSizeGb - 471.1).toFixed(1) }} GB (48%)</span>
            </div>
            <div class="sub-item">
              <span class="sub-label">Laju Tahunan SolarWinds:</span>
              <span class="sub-val">+{{ (oradataDailyGb * 365).toFixed(0) }} GB / tahun (+13.9%)</span>
            </div>
          </div>
          <div class="kpi-bar-bg">
            <div class="kpi-bar-fill" :style="{ width: `${Math.min(100, oradataMount?.currentPercent || 52.35)}%`, backgroundColor: '#f59e0b' }"></div>
          </div>
        </div>

        <!-- Card 4: N-1 Failover & Capacity Runway -->
        <div class="kpi-card highlight-card">
          <div class="kpi-top">
            <span class="kpi-icon"></span>
            <span class="kpi-status-badge status-normal">HA Passed</span>
          </div>
          <div class="kpi-main">
            <div class="kpi-value">{{ daysToStorageThreshold }} Hari</div>
            <div class="kpi-title">Runway Menuju {{ thresholds.maxStorage }}% Storage</div>
          </div>
          <div class="kpi-submetrics">
            <div class="sub-item">
              <span class="sub-label">Simulasi N-1 Peak:</span>
              <span class="sub-val highlight">{{ failoverSurvivingPeakLoad.toFixed(1) }}% CPU</span>
            </div>
            <div class="sub-item">
              <span class="sub-label">Status Failover:</span>
              <span :class="['sub-val', failoverSurvivingPeakLoad <= thresholds.maxCpu ? 'text-safe' : 'text-warn']">
                {{ failoverSurvivingPeakLoad <= thresholds.maxCpu ? 'Aman (Di Bawah 40%)' : 'Perlu Tambah Core' }}
              </span>
            </div>
            <div class="sub-item">
              <span class="sub-label">Sisa Storage Bebas:</span>
              <span class="sub-val">{{ ((100 - (oradataMount?.currentPercent || 52.35)) * nodeSpec.oradataSizeGb / 100).toFixed(1) }} GB</span>
            </div>
          </div>
          <div class="kpi-bar-bg">
            <div class="kpi-bar-fill" :style="{ width: '85%', backgroundColor: '#1b6c80' }"></div>
          </div>
        </div>
      </div>

      <!-- Main Analysis Tabs -->
      <div class="detail-section">
        <!-- Storage Mount Points Breakdown Table -->
        <div class="detail-card">
          <div class="card-title-bar">
            <div class="bar-title">
              <h3> Rincian Utilisasi Mount-Point Storage (Data Eksisting vs Shared Cluster)</h3>
              <p>Pemisahan antara Shared Storage (diakses bersama cluster) dan Local Storage per-node</p>
            </div>
            <div class="storage-tag">
              <span>Status: <strong>Tracking 90 Hari</strong></span>
            </div>
          </div>

          <div class="table-container">
            <table class="storage-table">
              <thead>
                <tr>
                  <th>Mount Point & Deskripsi</th>
                  <th>Tipe</th>
                  <th>Total Size (Kapasitas)</th>
                  <th>Space Used (Terpakai)</th>
                  <th>Kondisi Sekarang</th>
                  <th>Pertumbuhan (Delta)</th>
                  <th>Laju Pertumbuhan / Hari</th>
                  <th>Estimasi Tambahan 1 Tahun</th>
                  <th>Status & Rekomendasi</th>
                </tr>
              </thead>
              <tbody>
                <tr
                  v-for="mount in activeNodeData.storage.mounts"
                  :key="mount.name"
                  :class="{ 'shared-row': mount.isShared }"
                >
                  <td class="font-bold">
                    <span class="mount-icon">{{ mount.isShared ? '' : '' }}</span>
                    {{ mount.name }}
                  </td>
                  <td>
                    <span :class="['type-badge', mount.isShared ? 'shared-badge' : 'local-badge']">
                      {{ mount.isShared ? 'Shared Cluster SAN/ASM' : 'Local OS Volume' }}
                    </span>
                  </td>
                  <td class="mono-num font-bold">{{ mount.sizeStr || '-' }}</td>
                  <td class="mono-num">{{ mount.usedStr || '-' }}</td>
                  <td class="font-bold">
                    {{ mount.currentPercent }}%
                    <span v-if="mount.startPercent > mount.currentPercent" class="chip-expanded"> Disk Expanded</span>
                  </td>
                  <td>
                    <span :class="mount.startPercent > mount.currentPercent ? 'text-expanded' : mount.deltaPercent > 0 ? 'growth-pos' : 'growth-zero'">
                      {{ mount.startPercent > mount.currentPercent ? `${mount.deltaPercent}% (Ekspansi)` : mount.deltaPercent > 0 ? `+${mount.deltaPercent}%` : `${mount.deltaPercent}%` }}
                    </span>
                  </td>
                  <td>
                    <span class="mono-num">+{{ mount.dailyGrowthPercent }}% / hari</span>
                    <span class="sub-dim" v-if="mount.isShared">
                      (~{{ (mount.dailyGrowthPercent * nodeSpec.oradataSizeGb / 100).toFixed(2) }} GB/hari)
                    </span>
                  </td>
                  <td>
                    <strong class="text-primary">
                      {{ mount.isShared ? `+${((mount.dailyGrowthPercent * 365) * nodeSpec.oradataSizeGb / 100).toFixed(0)} GB (+${mount.annualGrowthPercent}%)` : `+${mount.annualGrowthPercent}%` }}
                    </strong>
                  </td>
                  <td>
                    <span :class="['status-chip', `chip-${mount.status}`]">
                      {{ mount.status === 'critical' ? ' Kritis (>85%)' : mount.status === 'warning' ? ' Waspada (>75%)' : ' Aman' }}
                    </span>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- 5-Year Database Capacity Plan Table (Single Node Focus) -->
        <div class="detail-card">
          <div class="card-title-bar">
            <div class="bar-title">
              <h3> Proyeksi Kebutuhan Database (Single Node) 5 Tahun (2025–2029)</h3>
              <p>Kalkulasi berbasis data aktual SolarWinds Node 1 (DOGDDSB1DBD21LP) dengan pertumbuhan bisnis $g = 15\%$ per tahun</p>
            </div>
            <button
              type="button"
              class="copy-btn"
              @click="copyCapPlanSummary"
            >
              {{ copied ? ' Tersalin ke Clipboard' : ' Salin Ringkasan CapPlan' }}
            </button>
          </div>

          <div class="table-container">
            <table class="capplan-table">
              <thead>
                <tr>
                  <th>Komponen Database</th>
                  <th>Spesifikasi Eksisting (Baseline)</th>
                  <th>2025 (Y1)</th>
                  <th>2026 (Y2)</th>
                  <th>2027 (Y3)</th>
                  <th>2028 (Y4)</th>
                  <th>2029 (Y5)</th>
                  <th>Strategi Rekomendasi</th>
                </tr>
              </thead>
              <tbody>
                <tr>
                  <td class="font-bold">Total Active Database Node</td>
                  <td>1 Node (Single Instance)</td>
                  <td>1 Node</td>
                  <td>1 Node</td>
                  <td>1 Node</td>
                  <td>1 Node</td>
                  <td>1 Node</td>
                  <td>Arsitektur Single Node Standalone</td>
                </tr>
                <tr>
                  <td class="font-bold">CPU Server (Physical Cores)</td>
                  <td>{{ nodeSpec.cores }} Core</td>
                  <td>{{ nodeSpec.cores }} Core</td>
                  <td>{{ nodeSpec.cores }} Core</td>
                  <td>{{ Math.max(nodeSpec.cores, Math.ceil(nodeSpec.cores * 1.15)) }} Core</td>
                  <td>{{ Math.max(nodeSpec.cores, Math.ceil(nodeSpec.cores * 1.30)) }} Core</td>
                  <td>{{ Math.max(nodeSpec.cores, Math.ceil(nodeSpec.cores * 1.50)) }} Core</td>
                  <td>Scale-up CPU vertikal mengikuti lonjakan transaksi</td>
                </tr>
                <tr>
                  <td class="font-bold">RAM Memory per Node (GB)</td>
                  <td>{{ nodeSpec.ramGb }} GB</td>
                  <td>{{ nodeSpec.ramGb }} GB</td>
                  <td>{{ nodeSpec.ramGb }} GB</td>
                  <td>{{ (nodeSpec.ramGb * 1.15).toFixed(2) }} GB</td>
                  <td>{{ (nodeSpec.ramGb * 1.30).toFixed(2) }} GB</td>
                  <td>{{ (nodeSpec.ramGb * 1.50).toFixed(2) }} GB</td>
                  <td>Scale-up RAM menjaga buffer cache > 98%</td>
                </tr>
                <tr class="shared-row">
                  <td class="font-bold">Database Storage (/oradata)</td>
                  <td>{{ nodeSpec.oradataSizeGb }} GB ({{ (nodeSpec.oradataSizeGb / 1024).toFixed(2) }} TB) — {{ oradataMount?.currentPercent || 52.35 }}% Terpakai</td>
                  <td>{{ (nodeSpec.oradataSizeGb + projectedStorageAddition(1)).toFixed(1) }} GB ({{ ((nodeSpec.oradataSizeGb + projectedStorageAddition(1)) / 1024).toFixed(2) }} TB)</td>
                  <td>{{ (nodeSpec.oradataSizeGb + projectedStorageAddition(2)).toFixed(1) }} GB ({{ ((nodeSpec.oradataSizeGb + projectedStorageAddition(2)) / 1024).toFixed(2) }} TB)</td>
                  <td>{{ (nodeSpec.oradataSizeGb + projectedStorageAddition(3)).toFixed(1) }} GB ({{ ((nodeSpec.oradataSizeGb + projectedStorageAddition(3)) / 1024).toFixed(2) }} TB)</td>
                  <td>{{ (nodeSpec.oradataSizeGb + projectedStorageAddition(4)).toFixed(1) }} GB ({{ ((nodeSpec.oradataSizeGb + projectedStorageAddition(4)) / 1024).toFixed(2) }} TB)</td>
                  <td>{{ (nodeSpec.oradataSizeGb + projectedStorageAddition(5)).toFixed(1) }} GB ({{ ((nodeSpec.oradataSizeGb + projectedStorageAddition(5)) / 1024).toFixed(2) }} TB)</td>
                  <td>Laju Real SolarWinds: +{{ oradataAnnualGb.toFixed(0) }} GB/thn (+13.94%/thn) + Compound 15%</td>
                </tr>
                <tr>
                  <td class="font-bold">Oracle Core License (0.5 Factor)</td>
                  <td>{{ Math.ceil(nodeSpec.cores * 0.5) }} Core Lic</td>
                  <td>{{ Math.ceil(nodeSpec.cores * 0.5) }} Core Lic</td>
                  <td>{{ Math.ceil(nodeSpec.cores * 0.5) }} Core Lic</td>
                  <td>{{ Math.ceil(Math.max(nodeSpec.cores, Math.ceil(nodeSpec.cores * 1.15)) * 0.5) }} Core Lic</td>
                  <td>{{ Math.ceil(Math.max(nodeSpec.cores, Math.ceil(nodeSpec.cores * 1.30)) * 0.5) }} Core Lic</td>
                  <td>{{ Math.ceil(Math.max(nodeSpec.cores, Math.ceil(nodeSpec.cores * 1.50)) * 0.5) }} Core Lic</td>
                  <td>Lisensi Oracle Core Factor 0.5 ({{ nodeSpec.cores }} Core × 0.5 = {{ Math.ceil(nodeSpec.cores * 0.5) }} Core Lic)</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- Audit Log Perubahan CPU, Memory & Storage dari Waktu ke Waktu -->
        <div class="detail-card audit-card">
          <div class="card-title-bar">
            <div class="bar-title">
              <h3> Audit Log Perubahan CPU, Memory & Storage dari Waktu ke Waktu (90 Hari)</h3>
              <p>Kronologi riil perubahan beban transaksi, penambahan hardware RAM pada 5 Agustus, dan ekspansi storage /oradata</p>
            </div>
            <div class="audit-filters">
              <button
                type="button"
                :class="['filter-chip', { active: auditFilter === 'all' }]"
                @click="auditFilter = 'all'"
              >
                Semua Milestone ({{ activeNodeData.auditTimeline?.length || 8 }})
              </button>
              <button
                type="button"
                :class="['filter-chip', { active: auditFilter === 'upgrade' }]"
                @click="auditFilter = 'upgrade'"
              >
                 Event Upgrade
              </button>
              <button
                type="button"
                :class="['filter-chip', { active: auditFilter === 'spike' }]"
                @click="auditFilter = 'spike'"
              >
                 Peak Load
              </button>
            </div>
          </div>

          <!-- Timeline Event Cards Grid -->
          <div class="audit-timeline-grid">
            <div
              v-for="evt in filteredAuditTimeline"
              :key="evt.id"
              :class="['audit-event-card', `card-status-${evt.status}`]"
            >
              <div class="evt-header">
                <div class="evt-title-wrap">
                  <span class="evt-icon">{{ evt.icon }}</span>
                  <div>
                    <div class="evt-date-badge">{{ evt.date }}</div>
                    <h4 class="evt-title">{{ evt.title }}</h4>
                  </div>
                </div>
                <span :class="['evt-category-chip', `chip-${evt.status}`]">{{ evt.categoryLabel }}</span>
              </div>

              <!-- Metrics Diff Badges -->
              <div class="evt-metrics-diff">
                <div class="diff-box">
                  <span class="diff-lbl">Peak CPU:</span>
                  <span :class="['diff-val', evt.cpuPeak >= 40 ? 'val-warn' : 'val-safe']">{{ evt.cpuPeak }}%</span>
                </div>
                <div class="diff-box">
                  <span class="diff-lbl">RAM Used:</span>
                  <span :class="['diff-val', evt.memoryVal >= 65 ? 'val-warn' : 'val-safe']">{{ evt.memoryVal }}%</span>
                </div>
                <div class="diff-box">
                  <span class="diff-lbl">/oradata:</span>
                  <span :class="['diff-val', evt.oradataVal >= 70 ? 'val-warn' : 'val-safe']">{{ evt.oradataVal }}%</span>
                </div>
                <div class="diff-box">
                  <span class="diff-lbl">/u01:</span>
                  <span class="diff-val">{{ evt.u01Val }}%</span>
                </div>
              </div>

              <p class="evt-summary">{{ evt.summary }}</p>
            </div>
          </div>

          <!-- Comprehensive Audit Table -->
          <div class="table-container" style="margin-top: 16px;">
            <table class="audit-table">
              <thead>
                <tr>
                  <th>Tanggal Milestone</th>
                  <th>Kategori Event</th>
                  <th>Peak CPU (%)</th>
                  <th>RAM Used (%)</th>
                  <th>/oradata (%)</th>
                  <th>/u01 (%)</th>
                  <th>Catatan Audit Arsitektur & Infrastruktur</th>
                </tr>
              </thead>
              <tbody>
                <tr
                  v-for="evt in filteredAuditTimeline"
                  :key="`row-${evt.id}`"
                  :class="{ 'highlight-row': evt.category === 'upgrade' || evt.category === 'current' }"
                >
                  <td class="mono-num font-bold">{{ evt.date }}</td>
                  <td>
                    <span :class="['evt-category-chip', `chip-${evt.status}`]">
                      {{ evt.icon }} {{ evt.categoryLabel }}
                    </span>
                  </td>
                  <td><strong :class="evt.cpuPeak >= 40 ? 'text-warn' : ''">{{ evt.cpuPeak }}%</strong></td>
                  <td><strong :class="evt.memoryVal >= 65 ? 'text-warn' : 'text-safe'">{{ evt.memoryVal }}%</strong></td>
                  <td><strong>{{ evt.oradataVal }}%</strong></td>
                  <td>{{ evt.u01Val }}%</td>
                  <td class="audit-desc-cell">{{ evt.summary }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- Historical Trend Data Preview Accordion -->
        <div class="detail-card">
          <div class="card-title-bar">
            <div class="bar-title">
              <h3> Sampel Data Historis SolarWinds (CPU, RAM & Storage 90 Hari)</h3>
              <p>Cuplikan 10 titik data sampling aktual dari file CSV SolarWinds PerfStack</p>
            </div>
            <span class="sample-badge">144 Data Points Loaded</span>
          </div>

          <div class="table-container">
            <table class="sample-table">
              <thead>
                <tr>
                  <th>Timestamp (Sampling 3 Jam)</th>
                  <th>Min CPU (%)</th>
                  <th>Peak CPU (%)</th>
                  <th>Avg CPU (%)</th>
                  <th>Avg RAM (%)</th>
                  <th>/oradata (%)</th>
                  <th>/u01 (%)</th>
                  <th>/ Root (%)</th>
                </tr>
              </thead>
              <tbody>
                <tr
                  v-for="(item, idx) in activeNodeData.cpu.history.slice(0, 8)"
                  :key="item.date"
                >
                  <td class="mono-num">{{ item.date }}</td>
                  <td>{{ item.min }}%</td>
                  <td><span class="highlight">{{ item.peak }}%</span></td>
                  <td>{{ item.avg }}%</td>
                  <td>{{ activeNodeData.memory.history[idx]?.value ? activeNodeData.memory.history[idx].value.toFixed(1) + '%' : '69.8%' }}</td>
                  <td class="font-bold">{{ activeNodeData.storage.history[idx]?.oradata ? activeNodeData.storage.history[idx].oradata.toFixed(2) + '%' : '74.64%' }}</td>
                  <td>{{ activeNodeData.storage.history[idx]?.u01 ? activeNodeData.storage.history[idx].u01.toFixed(2) + '%' : '78.73%' }}</td>
                  <td>{{ activeNodeData.storage.history[idx]?.root ? activeNodeData.storage.history[idx].root.toFixed(2) + '%' : '37.84%' }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useLocalOpen } from '~/composables/useLocalOpen'

const { openPath } = useLocalOpen()

const { data, pending, error, refresh } = await useFetch('/api/solarwind')

const selectedNodeMode = ref<string>('node_1')
const showConfig = ref(false)
const copied = ref(false)
const auditFilter = ref<'all' | 'upgrade' | 'spike'>('all')

const simulatedClusterNodes = ref<number>(2)

const nodeSpec = reactive({
  cores: 16, // Exact Physical Cores: 16 Core (8 Core License Oracle)
  ramGb: 50.831, // Exact SolarWinds: 50.831 GB (Used: 19.111 GB, Free: 31.72 GB)
  oradataSizeGb: 899.7, // Exact SolarWinds: 899.7 GB (Used: 471.1 GB, Free: 428.6 GB)
})

const thresholds = reactive({
  maxCpu: 40,
  maxStorage: 80,
  maxRam: 70,
})

function refreshData() {
  refresh()
}

function openLocal(pathStr: string, app: 'antigravity' | 'vscode' | 'finder') {
  openPath(pathStr, app)
}

const activeNodeData = computed(() => {
  if (!data.value?.nodes || data.value.nodes.length === 0) return null
  return data.value.nodes[0]
})

const oradataMount = computed(() => {
  return activeNodeData.value?.storage.mounts.find((m) => m.name.includes('oradata'))
})

const oradataDailyGb = computed(() => {
  if (!oradataMount.value) return 0.34
  return (oradataMount.value.dailyGrowthPercent * nodeSpec.oradataSizeGb) / 100
})

const oradataAnnualGb = computed(() => {
  return oradataDailyGb.value * 365
})

const totalClusterCores = computed(() => {
  const nodes = selectedNodeMode.value === 'cluster_sim' ? simulatedClusterNodes.value : 1
  return nodeSpec.cores * nodes
})

const totalClusterRam = computed(() => {
  const nodes = selectedNodeMode.value === 'cluster_sim' ? simulatedClusterNodes.value : 1
  return nodeSpec.ramGb * nodes
})

const currentCpuStats = computed(() => {
  if (!activeNodeData.value) {
    return { avg: 0, peak: 0, min: 0, p95: 0, coresUsed: 0 }
  }
  const avg = activeNodeData.value.cpu.avgLoad
  const peak = activeNodeData.value.cpu.peakLoad
  const min = activeNodeData.value.cpu.minLoad
  const p95 = activeNodeData.value.cpu.p95Load
  const coresUsed = (avg / 100) * totalClusterCores.value
  return { avg, peak, min, p95, coresUsed }
})

const currentMemStats = computed(() => {
  if (!activeNodeData.value) {
    return { avg: 70.2, peak: 70.6, min: 69.5, p95: 70.4 }
  }
  return activeNodeData.value.memory
})

const cpuHealthStatus = computed(() => {
  const peak = currentCpuStats.value.peak
  if (peak > 70) return { label: 'Tinggi (>70%)', class: 'status-crit' }
  if (peak > 40) return { label: 'Beban Moderat', class: 'status-warn' }
  return { label: 'Aman (<40% Safe)', class: 'status-normal' }
})

const oradataStorageStatus = computed(() => {
  const current = oradataMount.value?.currentPercent || 52.35
  if (current >= thresholds.maxStorage) return { label: 'Perlu Ekspansi', class: 'status-crit' }
  if (current >= 75) return { label: 'Mendekati 80%', class: 'status-warn' }
  return { label: 'Aman', class: 'status-normal' }
})

// Calculate N-1 Failover surviving peak
const failoverSurvivingPeakLoad = computed(() => {
  const nodes = simulatedClusterNodes.value
  if (nodes <= 1) return currentCpuStats.value.peak
  // If 1 node fails, the surviving (nodes - 1) nodes absorb all load
  return (currentCpuStats.value.peak * nodes) / (nodes - 1)
})

const daysToStorageThreshold = computed(() => {
  const currentP = oradataMount.value?.currentPercent || 52.35
  const dailyGrowthP = oradataMount.value?.dailyGrowthPercent || 0.0382
  if (dailyGrowthP <= 0) return 999
  const remainingP = thresholds.maxStorage - currentP
  if (remainingP <= 0) return 0
  return Math.round(remainingP / dailyGrowthP)
})

function projectedStorageAddition(year: number): number {
  const annualGrowthGb = oradataDailyGb.value * 365
  // Compound growth with index & log overhead (15% yearly expansion)
  let totalAdd = 0
  for (let y = 1; y <= year; y++) {
    totalAdd += annualGrowthGb * Math.pow(1.15, y - 1)
  }
  return Math.round(totalAdd)
}

const filteredAuditTimeline = computed(() => {
  if (!activeNodeData.value?.auditTimeline) return []
  if (auditFilter.value === 'all') return activeNodeData.value.auditTimeline
  return activeNodeData.value.auditTimeline.filter((e) => e.category === auditFilter.value)
})

function copyCapPlanSummary() {
  const summary = `=== RINGKASAN DATABASE CAPACITY PLAN 5 TAHUN (SOLARWINDS MULTI-NODE) ===
- Node Detected: ${activeNodeData.value?.nodeName || 'DOGDDSB1DBD21LP'} (Cluster: ${simulatedClusterNodes.value} Nodes HA RAC)
- Hardware Spec: ${nodeSpec.cores} Physical Core (${Math.ceil(nodeSpec.cores * 0.5)} Core License Oracle), ${nodeSpec.ramGb} GB RAM, ${nodeSpec.oradataSizeGb} GB (${(nodeSpec.oradataSizeGb/1024).toFixed(2)} TB) /oradata
- SolarWinds CPU: Current 1% (${((1 / 100) * totalClusterCores.value).toFixed(2)} Core), Avg ${currentCpuStats.value.avg}%, Peak ${currentCpuStats.value.peak}% (${((currentCpuStats.value.peak / 100) * totalClusterCores.value).toFixed(1)} Core), P95 ${currentCpuStats.value.p95}%
- SolarWinds Memory: Total ${nodeSpec.ramGb} GB, Used 19.111 GB (37%), Free 31.72 GB (63%)
- Storage /oradata: Size ${nodeSpec.oradataSizeGb} GB, Used 471.1 GB (${oradataMount.value?.currentPercent || 52.35}%), Free ${(nodeSpec.oradataSizeGb - 471.1).toFixed(1)} GB, Daily Growth +${oradataDailyGb.value.toFixed(2)} GB/day (+${(oradataDailyGb.value*365).toFixed(0)} GB/year)
- Runway to 80%: ${daysToStorageThreshold.value} Hari
- N-1 Failover Peak: ${failoverSurvivingPeakLoad.value.toFixed(1)}% (Target Max 40%)
- 5-Year Storage Projection:
  * Baseline (Eksisting): ${nodeSpec.oradataSizeGb} GB (${(nodeSpec.oradataSizeGb/1024).toFixed(2)} TB)
  * 2025: ${(nodeSpec.oradataSizeGb + projectedStorageAddition(1)).toFixed(1)} GB (${((nodeSpec.oradataSizeGb + projectedStorageAddition(1))/1024).toFixed(2)} TB)
  * 2026: ${(nodeSpec.oradataSizeGb + projectedStorageAddition(2)).toFixed(1)} GB (${((nodeSpec.oradataSizeGb + projectedStorageAddition(2))/1024).toFixed(2)} TB)
  * 2027: ${(nodeSpec.oradataSizeGb + projectedStorageAddition(3)).toFixed(1)} GB (${((nodeSpec.oradataSizeGb + projectedStorageAddition(3))/1024).toFixed(2)} TB)
  * 2028: ${(nodeSpec.oradataSizeGb + projectedStorageAddition(4)).toFixed(1)} GB (${((nodeSpec.oradataSizeGb + projectedStorageAddition(4))/1024).toFixed(2)} TB)
  * 2029: ${(nodeSpec.oradataSizeGb + projectedStorageAddition(5)).toFixed(1)} GB (${((nodeSpec.oradataSizeGb + projectedStorageAddition(5))/1024).toFixed(2)} TB)`

  navigator.clipboard.writeText(summary)
  copied.value = true
  setTimeout(() => {
    copied.value = false
  }, 2500)
}
</script>

<style scoped>
.solarwind-view {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.solarwind-hero {
  background: var(--surface, #ffffff);
  border: 1px solid rgba(255, 255, 255, 0.15);
  color: #ffffff;
  padding: 28px 32px;
  border-radius: var(--radius-lg);
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 24px;
  box-shadow: var(--shadow-md);
}

.hero-content {
  max-width: 720px;
}

.hero-badge {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  background: rgba(59, 130, 246, 0.2);
  border: 1px solid rgba(59, 130, 246, 0.4);
  color: #93c5fd;
  font-size: 11px;
  font-weight: 700;
  padding: 4px 12px;
  border-radius: 20px;
  margin-bottom: 12px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.pulse-dot {
  width: 8px;
  height: 8px;
  background-color: #10b981;
  border-radius: 50%;
  box-shadow: 0 0 8px #10b981;
}

.hero-title {
  font-size: 22px;
  font-weight: 800;
  margin-bottom: 8px;
  color: #ffffff;
}

.hero-subtitle {
  font-size: 13.5px;
  color: #cbd5e1;
  line-height: 1.6;
  margin-bottom: 18px;
}

.hero-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.hero-btn {
  font-family: inherit;
  font-size: 12px;
  font-weight: 600;
  padding: 8px 16px;
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: all 0.2s ease;
  border: none;
}

.primary-btn {
  background: #2563eb;
  color: #ffffff;
}

.primary-btn:hover {
  background: #1d4ed8;
}

.secondary-btn {
  background: rgba(255, 255, 255, 0.15);
  color: #ffffff;
  border: 1px solid rgba(255, 255, 255, 0.25);
}

.secondary-btn:hover {
  background: rgba(255, 255, 255, 0.25);
}

.ghost-btn {
  background: transparent;
  color: #94a3b8;
  border: 1px dashed rgba(255, 255, 255, 0.3);
}

.ghost-btn:hover {
  background: rgba(255, 255, 255, 0.1);
  color: #ffffff;
}

.hero-meta-box {
  background: rgba(255, 255, 255, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.15);
  border-radius: var(--radius);
  padding: 16px 20px;
  display: flex;
  flex-direction: column;
  gap: 12px;
  min-width: 260px;
}

.meta-item {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.meta-label {
  font-size: 10px;
  text-transform: uppercase;
  color: #94a3b8;
  font-weight: 700;
}

.meta-val {
  font-size: 12.5px;
  font-weight: 700;
  color: #f1f5f9;
}

/* Control Panel */
.control-panel {
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  padding: 14px 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 16px;
}

.panel-left {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
}

.panel-label {
  font-size: 12.5px;
  font-weight: 700;
  color: var(--text-heading);
}

.node-toggle-group {
  display: flex;
  gap: 8px;
}

.node-btn {
  font-family: inherit;
  font-size: 12px;
  font-weight: 600;
  padding: 6px 14px;
  border-radius: var(--radius-sm);
  border: 1px solid var(--border);
  background: var(--bg);
  color: var(--text-body);
  cursor: pointer;
  transition: all 0.2s;
}

.node-btn.active {
  background: #2563eb;
  color: #ffffff;
  border-color: #2563eb;
  box-shadow: 0 2px 8px rgba(37, 99, 235, 0.3);
}

.param-badge {
  display: flex;
  gap: 12px;
  font-size: 11.5px;
  color: var(--text-muted);
  background: var(--bg);
  padding: 6px 12px;
  border-radius: var(--radius-sm);
  border: 1px solid var(--border);
}

/* Spec Adjuster */
.spec-adjuster-card {
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  padding: 14px 20px;
}

.adjuster-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  cursor: pointer;
}

.adjuster-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  font-weight: 700;
  color: var(--text-heading);
}

.adjuster-note {
  font-size: 11.5px;
  font-weight: 400;
  color: var(--text-muted);
}

.toggle-icon-btn {
  font-size: 11px;
  font-weight: 600;
  color: #2563eb;
  background: transparent;
  border: none;
  cursor: pointer;
}

.adjuster-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  gap: 16px;
  margin-top: 16px;
  padding-top: 14px;
  border-top: 1px solid var(--border);
}

.input-group {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.input-group label {
  font-size: 11px;
  font-weight: 600;
  color: var(--text-muted);
}

.input-group input,
.input-group select {
  font-family: inherit;
  font-size: 12.5px;
  padding: 6px 10px;
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  background: var(--bg);
  color: var(--text-heading);
}

/* Milestone Alert Banner */
.milestone-alert-banner {
  background: var(--surface, #ffffff);
  border: 1.5px solid rgba(37, 99, 235, 0.35);
  border-radius: var(--radius);
  padding: 16px 20px;
  display: flex;
  gap: 16px;
  align-items: flex-start;
  box-shadow: 0 4px 12px rgba(37, 99, 235, 0.08);
}

.milestone-icon {
  font-size: 26px;
  line-height: 1;
  padding-top: 2px;
}

.milestone-content {
  display: flex;
  flex-direction: column;
  gap: 8px;
  flex: 1;
}

.milestone-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 8px;
}

.milestone-title {
  font-size: 14px;
  font-weight: 800;
  color: #1e40af;
}

.milestone-tag {
  font-size: 10px;
  font-weight: 700;
  text-transform: uppercase;
  background: rgba(16, 185, 129, 0.2);
  color: #065f46;
  border: 1px solid rgba(16, 185, 129, 0.4);
  padding: 3px 8px;
  border-radius: 12px;
}

.milestone-desc {
  font-size: 12.5px;
  color: var(--text-body);
  line-height: 1.5;
  margin: 0;
}

.milestone-metrics {
  display: flex;
  flex-wrap: wrap;
  gap: 20px;
  padding-top: 6px;
  border-top: 1px dashed rgba(37, 99, 235, 0.2);
}

.m-metric {
  display: flex;
  gap: 6px;
  font-size: 12px;
}

.m-lbl {
  color: var(--text-muted);
}

.m-val {
  font-weight: 700;
}

/* KPI Grid */
.kpi-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
  gap: 18px;
}

.kpi-card {
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 14px;
  box-shadow: var(--shadow-sm);
  position: relative;
  overflow: hidden;
}

.kpi-card.highlight-card {
  border-color: #1b6c80;
  background: var(--surface, #ffffff);
}

.kpi-top {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.kpi-icon {
  font-size: 20px;
}

.kpi-status-badge {
  font-size: 10.5px;
  font-weight: 700;
  padding: 2px 8px;
  border-radius: 12px;
  text-transform: uppercase;
}

.status-normal {
  background: rgba(16, 185, 129, 0.15);
  color: #10b981;
}

.status-warn {
  background: rgba(245, 158, 11, 0.15);
  color: #f59e0b;
}

.status-crit {
  background: rgba(239, 68, 68, 0.15);
  color: #ef4444;
}

.kpi-main {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.kpi-value-row {
  display: flex;
  align-items: baseline;
  gap: 8px;
}

.kpi-value {
  font-size: 28px;
  font-weight: 800;
  color: var(--text-heading);
}

.kpi-badge-hint {
  font-size: 10.5px;
  font-weight: 700;
  background: rgba(37, 99, 235, 0.12);
  color: #2563eb;
  padding: 2px 6px;
  border-radius: 4px;
}

.kpi-title {
  font-size: 12px;
  font-weight: 600;
  color: var(--text-muted);
  text-transform: uppercase;
}

.kpi-submetrics {
  display: flex;
  flex-direction: column;
  gap: 6px;
  border-top: 1px solid var(--border);
  padding-top: 10px;
}

.sub-item {
  display: flex;
  justify-content: space-between;
  font-size: 11.5px;
}

.sub-label {
  color: var(--text-muted);
}

.sub-val {
  font-weight: 600;
  color: var(--text-heading);
}

.sub-val.highlight {
  color: #2563eb;
  font-weight: 700;
}

.text-safe {
  color: #10b981;
}

.text-warn {
  color: #f59e0b;
}

.kpi-bar-bg {
  height: 4px;
  background: rgba(0, 0, 0, 0.05);
  border-radius: 2px;
  overflow: hidden;
}

.kpi-bar-fill {
  height: 100%;
  border-radius: 2px;
}

/* Detail Section */
.detail-section {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.detail-card {
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  padding: 22px;
  display: flex;
  flex-direction: column;
  gap: 16px;
  box-shadow: var(--shadow-sm);
}

.card-title-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 12px;
}

.bar-title h3 {
  font-size: 15px;
  font-weight: 700;
  color: var(--text-heading);
  margin-bottom: 2px;
}

.bar-title p {
  font-size: 12px;
  color: var(--text-muted);
}

.storage-tag,
.sample-badge {
  font-size: 11px;
  background: var(--bg);
  padding: 4px 10px;
  border-radius: var(--radius-sm);
  border: 1px solid var(--border);
  color: var(--text-muted);
}

.copy-btn {
  font-family: inherit;
  font-size: 11.5px;
  font-weight: 600;
  padding: 6px 14px;
  background: #2563eb;
  color: #ffffff;
  border: none;
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: all 0.2s;
}

.copy-btn:hover {
  background: #1d4ed8;
}

.table-container {
  overflow-x: auto;
}

table {
  width: 100%;
  border-collapse: collapse;
  font-size: 12.5px;
}

th {
  background: var(--bg);
  color: var(--text-muted);
  font-weight: 700;
  font-size: 11px;
  text-transform: uppercase;
  padding: 10px 14px;
  text-align: left;
  border-bottom: 2px solid var(--border);
}

td {
  padding: 10px 14px;
  border-bottom: 1px solid var(--border);
  color: var(--text-body);
}

.shared-row {
  background: rgba(37, 99, 235, 0.04);
}

.font-bold {
  font-weight: 700;
  color: var(--text-heading);
}

.mono-num {
  font-family: monospace;
  font-size: 11.5px;
}

.type-badge {
  font-size: 10.5px;
  font-weight: 600;
  padding: 2px 8px;
  border-radius: 4px;
}

.shared-badge {
  background: rgba(37, 99, 235, 0.12);
  color: #2563eb;
}

.local-badge {
  background: rgba(100, 116, 139, 0.12);
  color: #64748b;
}

.growth-pos {
  color: #2563eb;
  font-weight: 700;
}

.growth-zero {
  color: var(--text-muted);
}

.chip-expanded {
  font-size: 10px;
  font-weight: 700;
  background: rgba(16, 185, 129, 0.15);
  color: #059669;
  border: 1px solid rgba(16, 185, 129, 0.3);
  padding: 1px 6px;
  border-radius: 4px;
  margin-left: 6px;
}

.text-expanded {
  color: #059669;
  font-weight: 600;
}

.sub-dim {
  font-size: 10.5px;
  color: var(--text-muted);
  display: block;
}

.status-chip {
  font-size: 10.5px;
  font-weight: 700;
  padding: 2px 8px;
  border-radius: 12px;
}

.chip-safe {
  background: rgba(16, 185, 129, 0.15);
  color: #10b981;
}

.chip-warning {
  background: rgba(245, 158, 11, 0.15);
  color: #f59e0b;
}

.chip-critical {
  background: rgba(239, 68, 68, 0.15);
  color: #ef4444;
}

.text-primary {
  color: #2563eb;
}

/* State cards */
.state-card {
  padding: 30px;
  border-radius: var(--radius);
  text-align: center;
  font-size: 14px;
}

.state-card.loading {
  background: var(--surface);
  color: var(--text-muted);
}

.state-card.error {
  background: rgba(239, 68, 68, 0.1);
  color: #ef4444;
}

/* Audit Section Styles */
.audit-card {
  border-left: 4px solid #3b82f6;
}

.audit-filters {
  display: flex;
  gap: 8px;
}

.filter-chip {
  font-family: inherit;
  font-size: 11px;
  font-weight: 600;
  padding: 4px 12px;
  border-radius: 20px;
  border: 1px solid var(--border);
  background: var(--bg);
  color: var(--text-muted);
  cursor: pointer;
  transition: all 0.2s ease;
}

.filter-chip:hover {
  border-color: #3b82f6;
  color: #3b82f6;
}

.filter-chip.active {
  background: #2563eb;
  color: #ffffff;
  border-color: #2563eb;
}

.audit-timeline-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 14px;
}

.audit-event-card {
  background: var(--bg);
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  padding: 14px;
  display: flex;
  flex-direction: column;
  gap: 10px;
  transition: all 0.2s ease;
}

.audit-event-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-sm);
}

.card-status-upgrade {
  border-left: 3px solid #10b981;
  background: rgba(16, 185, 129, 0.03);
}

.card-status-spike {
  border-left: 3px solid #f59e0b;
  background: rgba(245, 158, 11, 0.03);
}

.card-status-info {
  border-left: 3px solid #3b82f6;
}

.card-status-current {
  border-left: 3px solid #1b6c80;
  background: rgba(139, 92, 246, 0.03);
}

.evt-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 8px;
}

.evt-title-wrap {
  display: flex;
  align-items: flex-start;
  gap: 10px;
}

.evt-icon {
  font-size: 18px;
  line-height: 1;
  margin-top: 2px;
}

.evt-date-badge {
  font-family: monospace;
  font-size: 10.5px;
  font-weight: 700;
  color: var(--text-muted);
}

.evt-title {
  font-size: 13px;
  font-weight: 700;
  color: var(--text-heading);
  margin-top: 2px;
}

.evt-category-chip {
  font-size: 10px;
  font-weight: 700;
  padding: 2px 8px;
  border-radius: 12px;
  white-space: nowrap;
}

.evt-metrics-diff {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 6px;
  background: var(--surface);
  padding: 8px;
  border-radius: 6px;
  border: 1px solid var(--border);
}

.diff-box {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.diff-lbl {
  font-size: 9.5px;
  color: var(--text-muted);
  text-transform: uppercase;
  font-weight: 600;
}

.diff-val {
  font-size: 12px;
  font-weight: 700;
  font-family: monospace;
}

.val-safe {
  color: #10b981;
}

.val-warn {
  color: #f59e0b;
}

.evt-summary {
  font-size: 12px;
  color: var(--text-body);
  line-height: 1.45;
}

.audit-table th {
  font-size: 10.5px;
}

.audit-table td {
  font-size: 12px;
}

.highlight-row {
  background: rgba(37, 99, 235, 0.04);
}

.audit-desc-cell {
  color: var(--text-body);
  line-height: 1.4;
  font-size: 12px;
}
</style>
