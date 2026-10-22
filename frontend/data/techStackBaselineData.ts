// ─────────────────────────────────────────────────────────────────────────────
// BASELINE TECHNOLOGY STACK & SIZING PATTERN (BANK MANDIRI)
// Standar Arsitektur: Tier Moderate (IHSP & RMT) vs Tier Critical (Livin' Merchant)
// ─────────────────────────────────────────────────────────────────────────────

export interface TechStackItem {
  id: string
  no: number
  layer: string
  component: string
  icon: string
  moderate: {
    nodes: string
    type: string
    spec: string
  }
  critical: {
    nodes: string
    type: string
    spec: string
  }
  license: {
    free: string
    enterprise: string
  }
  notes: string
  contractRef?: string
}

export interface ArchitectureForumReference {
  title: string
  tier: 'moderate' | 'critical'
  badge: string
  points: {
    label: string
    description: string
  }[]
}

export const BASELINE_KPI_SUMMARY = [
  {
    id: 'mod_compute',
    tier: 'moderate',
    title: 'Tier Moderate Compute',
    value: '14 VM Nodes',
    detail: '112 vCPU | 256 GB RAM (IHSP & RMT)',
    bg: '#002060',
    accent: '#38bdf8'
  },
  {
    id: 'mod_storage',
    tier: 'moderate',
    title: 'Tier Moderate Storage',
    value: '4.500 GB (4.5 TB)',
    detail: 'SAN Tier 2 LUN Virtual Disk',
    bg: '#1B365D',
    accent: '#60a5fa'
  },
  {
    id: 'crit_compute',
    tier: 'critical',
    title: 'Tier Critical Compute',
    value: '28 Dedicated Nodes',
    detail: '384 vCPU | 4.8 TB RAM (incl 4 Bare Metal)',
    bg: '#78350f',
    accent: '#fbbf24'
  },
  {
    id: 'crit_storage',
    tier: 'critical',
    title: 'Tier Critical Storage',
    value: '24.500 GB (24.5 TB)',
    detail: 'High-Perf ASM SAN Tier 2 Storage',
    bg: '#b45309',
    accent: '#fde047'
  }
]

export const BASELINE_TECH_STACK_ITEMS: TechStackItem[] = [
  {
    id: 'web_ingress',
    no: 1,
    layer: 'Web Ingress',
    component: 'NGINX Web Server',
    icon: '🌐',
    moderate: {
      nodes: '2N (2 Nodes)',
      type: 'Virtual Machine (VM)',
      spec: '8 vCPU / 16GB RAM / 100-200GB Disk'
    },
    critical: {
      nodes: '2N (2 Nodes)',
      type: 'Virtual Machine (VM)',
      spec: '8-16 vCPU / 16-32GB RAM / 100-200GB Disk'
    },
    license: {
      free: 'NGINX Community (com)',
      enterprise: 'NGINX Plus (F5 Enterprise)'
    },
    notes: 'Moderate: NGINX VM 8/16/100 (IHSP) / 8/16/200 (RMT). Critical: NGINX Plus VM 8-16 vCPU / 16-32 GB (Kontrak F5: TIF.DDL/GH.9076/2025).',
    contractRef: 'TIF.DDL/GH.9076/2025'
  },
  {
    id: 'app_container',
    no: 2,
    layer: 'App Container',
    component: 'Red Hat OCP Worker',
    icon: '📦',
    moderate: {
      nodes: '1W (1 Worker)',
      type: 'Shared / VM Worker',
      spec: '64 vCPU / 1024GB RAM / 300GB Disk'
    },
    critical: {
      nodes: '2N (2 Nodes)',
      type: 'Dedicated Baremetal Worker',
      spec: '2x 64 vCPU / 1024GB RAM / 300GB Disk'
    },
    license: {
      free: 'OKD (Non-Prod)',
      enterprise: 'RH OCP Platform Plus BM Premium'
    },
    notes: 'Moderate: 1 OCP Worker VM (64/1024/300) per DC. Critical: 2 Dedicated Bare Metal Workers per DC (PPWRKR01/02 & RPWRKR01/02, Kontrak Red Hat: TIF.DDL/GH.1993/2025).',
    contractRef: 'TIF.DDL/GH.1993/2025'
  },
  {
    id: 'container_router',
    no: 3,
    layer: 'Container Router',
    component: 'OCP Infra Router & Bastion',
    icon: '🔀',
    moderate: {
      nodes: '3 Infra + 1 Bastion',
      type: 'VM Router / HA Proxy',
      spec: '8 vCPU / 32GB RAM / 150GB Disk'
    },
    critical: {
      nodes: 'Dedicated Router Pods',
      type: 'High-Throughput Router Pods',
      spec: '8 vCPU / 32GB RAM / 150GB Disk'
    },
    license: {
      free: 'Included in OCP',
      enterprise: 'Included in OCP Platform Plus'
    },
    notes: 'Moderate: 3 Master (8/32/150) + 3 Infra Router (8/32/150) + 1 Bastion (4/8/100). Critical: Dedicated high-throughput router pods cross-DC.',
    contractRef: 'Included in OCP'
  },
  {
    id: 'core_db',
    no: 4,
    layer: 'Core Database',
    component: 'Database Engine (Oracle)',
    icon: '🗄️',
    moderate: {
      nodes: 'Single Node, VM',
      type: 'Virtual Machine (VM)',
      spec: '16 vCPU / 32GB RAM / 1000GB Disk'
    },
    critical: {
      nodes: 'RAC, BM (Baremetal)',
      type: 'Physical Baremetal + ASM',
      spec: '2x 64 vCPU BM Server / 1024GB RAM'
    },
    license: {
      free: 'PostgreSQL / MySQL Com',
      enterprise: 'Oracle DB Enterprise (RAC + ADG)'
    },
    notes: 'Moderate: Single Node VM + Data Guard (16/32/1000) tanpa RAC Option untuk efisiensi lisensi core. Critical: Oracle RAC Bare Metal 2x64-Core (High TPS & Zero Data Loss, Kontrak: TIF.DDL/SCD.3585/2025).',
    contractRef: 'TIF.DDL/SCD.3585/2025'
  },
  {
    id: 'db_storage',
    no: 5,
    layer: 'Database Storage',
    component: 'Oracle ASM Shared SAN',
    icon: '💾',
    moderate: {
      nodes: 'SAN Virtual Disk LUN',
      type: 'Storage SAN Tier 2',
      spec: '1.000 GB - 1.500 GB SAN LUN'
    },
    critical: {
      nodes: 'Dedicated ASM SAN Cluster',
      type: 'High-Perf SAN Tier 2',
      spec: '5.500 GB (5.5 TB: Reco 500GB, Data 5TB)'
    },
    license: {
      free: 'SAN LUN Tier 2',
      enterprise: 'SAN LUN + Oracle ASM'
    },
    notes: 'Moderate: SAN Tier 2 Virtual LUN 1.0 - 1.5 TB. Critical: SAN Tier 2 ASM 5.5 TB (Reco 500GB, Data 5000GB) pada node bare metal.',
    contractRef: 'SAN Tier 2 ASM'
  },
  {
    id: 'dist_cache',
    no: 6,
    layer: 'Distributed Cache',
    component: 'Redis Caching Cluster',
    icon: '⚡',
    moderate: {
      nodes: '1-3 Nodes',
      type: 'Virtual Machine (VM)',
      spec: '8 vCPU / 16GB RAM / 200GB Disk'
    },
    critical: {
      nodes: '3 Nodes Multi-Shard',
      type: 'Virtual Machine (VM)',
      spec: '3x 8 vCPU / 32GB RAM / 200GB Disk'
    },
    license: {
      free: 'Redis Community (com)',
      enterprise: 'Redis Enterprise Software'
    },
    notes: 'Moderate: 1-3 VM Redis Community (8/16/200). Critical: 3 VM Redis Enterprise Multi-Shard (PVLMRD01-03 8/32/200, Kontrak: TIF.DDL/GH.9080/2024).',
    contractRef: 'TIF.DDL/GH.9080/2024'
  },
  {
    id: 'cache_sentinel',
    no: 7,
    layer: 'Cache Sentinel',
    component: 'Redis Sentinel Quorum',
    icon: '🛡️',
    moderate: {
      nodes: 'Co-located / Sentinel',
      type: 'Shared VM',
      spec: '4 vCPU / 8GB RAM / 100GB Disk'
    },
    critical: {
      nodes: 'Dedicated Sentinel DC2',
      type: 'Virtual Machine (VM)',
      spec: '4-8 vCPU / 16GB RAM / 100GB Disk'
    },
    license: {
      free: 'Redis Community (com)',
      enterprise: 'Redis Enterprise'
    },
    notes: 'Moderate: Co-located sentinel. Critical: Dedicated sentinel cross-DC (Wisdan & Rempoa) untuk quorum deteksi failover otomatis.',
    contractRef: 'TIF.DDL/GH.9080/2024'
  },
  {
    id: 'message_bus',
    no: 8,
    layer: 'Message Bus',
    component: 'Confluent / Apache Kafka',
    icon: '📨',
    moderate: {
      nodes: '3 Nodes (Shared/VM)',
      type: 'VM Broker Cluster',
      spec: '16 vCPU / 64GB / 600GB (IHSP)'
    },
    critical: {
      nodes: 'Dedicated Cluster (7 Nodes)',
      type: 'ZK + Broker + C3 Cluster',
      spec: '3 ZK (4/32) + 3 Broker (16/32/1000) + 1 C3'
    },
    license: {
      free: 'Apache / Confluent Com',
      enterprise: 'Confluent Kafka Enterprise'
    },
    notes: 'Moderate: 3 VM Kafka Community (16/64/600 atau 8/16/1000). Critical: 7 Dedicated VM Confluent Kafka Enterprise (3 ZK + 3 Broker + 1 C3, Kontrak: TIF.DDL/GH.9075/2024).',
    contractRef: 'TIF.DDL/GH.9075/2024'
  },
  {
    id: 'log_analytics',
    no: 9,
    layer: 'Log Analytics',
    component: 'Elasticsearch Analytics Node',
    icon: '🔍',
    moderate: {
      nodes: 'Combined dgn Logstash',
      type: 'Shared VM / Combined',
      spec: '3 Nodes (8 vCPU / 32GB RAM / 500GB)'
    },
    critical: {
      nodes: 'Dedicated Cluster (8 Nodes)',
      type: 'Master + Hot + Warm + Coord',
      spec: '3 Master (8/32/50) + 2 Hot (16/32/1900) + Warm + Coord'
    },
    license: {
      free: 'Elasticsearch / OpenSearch',
      enterprise: 'Elasticsearch Platinum'
    },
    notes: 'Moderate: ELK combined (1 Logstash + 3 Master/Data 8/32/500). Critical: Dedicated 8 VM (3 Master 8/32/50, 2 Hot Data 16/32/1900, 1 Warm, 2 Coord, Kontrak: TIF.DDL/GH.9074/2024).',
    contractRef: 'TIF.DDL/GH.9074/2024'
  },
  {
    id: 'log_shipper',
    no: 10,
    layer: 'Log Shipper',
    component: 'Logstash Pipeline Shipper',
    icon: '🪵',
    moderate: {
      nodes: 'Combined / Co-located',
      type: 'Shared VM',
      spec: '1 Node (8 vCPU / 32GB RAM / 100GB)'
    },
    critical: {
      nodes: 'Dedicated Shipper Nodes',
      type: 'Virtual Machine (VM)',
      spec: '2 Nodes (8 vCPU / 32GB RAM / 100GB)'
    },
    license: {
      free: 'Logstash Open Source',
      enterprise: 'Elastic Platinum'
    },
    notes: 'Moderate: Combined pada node ELK. Critical: Dedicated 2 VM Logstash Shipper (PVLMEL01/02 8/32/100 di DC & RVLMEL01/02 di DRC) terpisah dari data cluster demi kestabilan heap JVM.',
    contractRef: 'TIF.DDL/GH.9074/2024'
  },
  {
    id: 'log_dashboard',
    no: 11,
    layer: 'Log Dashboard',
    component: 'Kibana Observability Portal',
    icon: '📊',
    moderate: {
      nodes: 'Shared Portal',
      type: 'Shared VM Portal',
      spec: 'Resource sharing on ELK'
    },
    critical: {
      nodes: 'Dedicated Portal Nodes',
      type: 'Virtual Machine (VM)',
      spec: '2 Nodes (4 vCPU / 32GB RAM / 50GB)'
    },
    license: {
      free: 'Kibana Community (com)',
      enterprise: 'Elastic Platinum'
    },
    notes: 'Moderate: Shared portal. Critical: Terintegrasi pada coordinator node PVLMEC01/02 (4/32/50 di DC & RVLMEC01/02 di DRC).',
    contractRef: 'TIF.DDL/GH.9074/2024'
  },
  {
    id: 'doc_storage',
    no: 12,
    layer: 'Document Storage',
    component: 'SFTP / Object Storage / RDM',
    icon: '🗂️',
    moderate: {
      nodes: 'SFTP Server VM',
      type: 'Virtual Machine (VM)',
      spec: '8 vCPU / 16GB RAM / 1.500GB Disk'
    },
    critical: {
      nodes: 'SAN RDM / Dell ECS S3',
      type: 'Dedicated SAN LUN',
      spec: '8 vCPU / 16GB RAM / 365GB Disk'
    },
    license: {
      free: 'MinIO / OpenSFTP (com)',
      enterprise: 'Dell ECS / SAN Enterprise'
    },
    notes: 'Moderate: VM SFTP Server (8/16/1500) untuk pertukaran file benefit MFT. Critical: Storage Gateway SAN RDM / Dell ECS S3 untuk arsip transaksi & bukti settlement.',
    contractRef: 'Dell ECS / SAN RDM'
  },
  {
    id: 'batch_engine',
    no: 13,
    layer: 'Batch Engine',
    component: 'BMC Control-M / OCP Cron',
    icon: '⚙️',
    moderate: {
      nodes: 'Scheduled Cron on OCP',
      type: 'Shared OCP Worker',
      spec: 'Resource sharing on OCP Pods'
    },
    critical: {
      nodes: 'Dedicated Batch Control-M',
      type: 'Virtual Machine (VM)',
      spec: '4 vCPU / 8GB RAM / 100GB Disk'
    },
    license: {
      free: 'OCP Native Batch Cron',
      enterprise: 'BMC Control-M Enterprise'
    },
    notes: 'Moderate: Scheduled CronJob pada worker container OCP. Critical: Dedicated VM Agent BMC Control-M (PVLMBC01 / RVLMBC01 4/8/100) untuk batch settlement end-of-day.',
    contractRef: 'BMC Control-M'
  }
]

export const EA_FORUM_REFERENCES: ArchitectureForumReference[] = [
  {
    title: "TIER MODERATE: IHSP (STM2601010) & RMT (RM TOOLS SATU R1)",
    tier: 'moderate',
    badge: 'Forum C / SLA 99.5% - 99.9%',
    points: [
      {
        label: 'IHSP (Integrated HC Suite Platform)',
        description: 'Acuan Dokumen: EA Forum C STM2601010 R1.1. Sizing per site: NGINX 8/16/100, OCP Master 8/32/150, Infra Router 8/32/150, 1 OCP Worker 64/1024/300, DB VM 16/32/1000, Kafka 16/64/600, SFTP 8/16/1500, ELK 8/32/500.'
      },
      {
        label: 'RMT (RM Tools SATU R1) & Pola Migrasi RAC',
        description: 'Acuan Dokumen: Capacity Plan RM Tools SATU R1 & evaluasi migrasi database (20260824). NGINX 8/16/200, Redis 8/16/200, Kafka 8/16/1000. Database dialihkan dari RAC ke Single Node Data Guard (Wisdan Primary, Rempoa Standby) untuk optimasi lisensi core.'
      },
      {
        label: 'Ketentuan Arsitektur Tier Moderate',
        description: 'Kategori EA Forum: Forum C. Target SLA Ketersediaan: 99.5% - 99.9% (RTO < 2 jam, RPO < 15 menit). Seluruh server berbasis Virtual Machine (VM). Database Single Node Data Guard. Server logging (ELK & Logstash) diperbolehkan shared/co-located.'
      }
    ]
  },
  {
    title: "TIER CRITICAL: LIVIN' MERCHANT (DRB2502010 R6 SESI 2)",
    tier: 'critical',
    badge: 'Forum B (Very Critical) / SLA 99.99%',
    points: [
      {
        label: "Livin' Merchant (Release 6 Sesi 2)",
        description: "Acuan Dokumen: EA Forum B DRB2502010 R6 Sesi 2. Klasifikasi: Very Critical (Forum B) melayani transaksi merchant nasional 24/7 (SLA 99.99%). Wajib 2 Dedicated OCP Bare Metal Workers (64/1024/300) dan Oracle RAC Bare Metal 64-Core (ASM SAN 5.5 TB)."
      },
      {
        label: 'Pemisahan Node Logging (ELK & Logstash)',
        description: 'Acuan Server List LM (Tabel 108): Elasticsearch dipisah menjadi Master (PVLMEM01-03 8/32/50), Hot Data (PVLMDH01-02 16/32/1900), Warm Data (PVLMDW01), Coordinator (PVLMEC01-02), dan Logstash Shipper (PVLMEL01-02 8/32/100) di server tersendiri demi stabilitas indexing.'
      },
      {
        label: "Lisensi Resmi Enterprise di Livin' Merchant",
        description: 'NGINX Plus (TIF.DDL/GH.9076/2025), Red Hat OCP Platform Plus Bare Metal (TIF.DDL/GH.1993/2025), Oracle RAC Enterprise (TIF.DDL/SCD.3585/2025), Confluent Kafka (TIF.DDL/GH.9075/2024), Elasticsearch Platinum (TIF.DDL/GH.9074/2024), Redis Enterprise (TIF.DDL/GH.9080/2024).'
      }
    ]
  }
]
