import type { SizingPresetItem } from '~/types/sizing'

export const DEFAULT_PRESET_ITEMS: SizingPresetItem[] = [
  {
    id: 'web_ingress',
    no: 1,
    layer: 'Web Ingress',
    component: 'NGINX Web Server',
    icon: '🌐',
    licenseFree: 'NGINX Community (OSS)',
    licenseEnterprise: 'NGINX Plus (F5 Enterprise)',
    notes: 'Reverse proxy, SSL termination, and rate limiting entry layer.',
    contractRef: 'TIF.DDL/GH.9076/2025',
    moderate: {
      nodeCount: 2,
      vcpuPerNode: 8,
      ramGBPerNode: 16,
      storageGBPerNode: 200,
      storageClass: 'Tier 2',
      computeType: 'vm',
      backupEnabled: true,
      licenseKind: 'free'
    },
    critical: {
      nodeCount: 2,
      vcpuPerNode: 8,
      ramGBPerNode: 16,
      storageGBPerNode: 200,
      storageClass: 'Tier 2',
      computeType: 'vm',
      backupEnabled: true,
      licenseKind: 'enterprise'
    }
  },
  {
    id: 'app_container',
    no: 2,
    layer: 'App Container',
    component: 'Red Hat OCP Worker',
    icon: '📦',
    licenseFree: 'OKD (Community OSS)',
    licenseEnterprise: 'Red Hat OpenShift Platform Plus Baremetal',
    notes: 'Microservices orchestration cluster on OpenShift.',
    contractRef: 'TIF.DDL/GH.1993/2025',
    moderate: {
      nodeCount: 1,
      vcpuPerNode: 64,
      ramGBPerNode: 1024,
      storageGBPerNode: 1200,
      storageClass: 'Tier 2',
      computeType: 'vm',
      backupEnabled: true,
      licenseKind: 'free'
    },
    critical: {
      nodeCount: 2,
      vcpuPerNode: 64,
      ramGBPerNode: 1024,
      storageGBPerNode: 300,
      storageClass: 'Tier 2',
      computeType: 'baremetal',
      backupEnabled: true,
      licenseKind: 'enterprise'
    }
  },
  {
    id: 'container_router',
    no: 3,
    layer: 'Container Router',
    component: 'OCP Infra Router & Bastion',
    icon: '🔀',
    licenseFree: 'Included in OCP',
    licenseEnterprise: 'Included in OCP Platform Plus',
    notes: 'Ingress routing & HA proxy pods for container workloads.',
    contractRef: 'Included in OCP',
    moderate: {
      nodeCount: 7,
      vcpuPerNode: 8,
      ramGBPerNode: 32,
      storageGBPerNode: 150,
      storageClass: 'Tier 2',
      computeType: 'vm',
      backupEnabled: true,
      licenseKind: 'free'
    },
    critical: {
      nodeCount: 7,
      vcpuPerNode: 8,
      ramGBPerNode: 32,
      storageGBPerNode: 150,
      storageClass: 'Tier 2',
      computeType: 'vm',
      backupEnabled: true,
      licenseKind: 'enterprise'
    }
  },
  {
    id: 'core_db',
    no: 4,
    layer: 'Core Database',
    component: 'Database Engine (Oracle)',
    icon: '🗄️',
    licenseFree: 'PostgreSQL / MySQL Community',
    licenseEnterprise: 'Oracle Database Enterprise Edition (RAC + ADG)',
    notes: 'High-availability transaction database engine.',
    contractRef: 'TIF.DDL/SCD.3585/2025',
    moderate: {
      nodeCount: 1,
      vcpuPerNode: 16,
      ramGBPerNode: 32,
      storageGBPerNode: 1000,
      storageClass: 'Tier 2',
      computeType: 'vm',
      backupEnabled: true,
      licenseKind: 'free'
    },
    critical: {
      nodeCount: 2,
      vcpuPerNode: 64,
      ramGBPerNode: 1024,
      storageGBPerNode: 1000,
      storageClass: 'Tier 3',
      computeType: 'baremetal',
      backupEnabled: true,
      licenseKind: 'enterprise'
    }
  },
  {
    id: 'db_storage',
    no: 5,
    layer: 'Database Storage',
    component: 'Oracle ASM Shared SAN',
    icon: '💾',
    licenseFree: 'LVM / Linux Software RAID',
    licenseEnterprise: 'Oracle ASM (Enterprise Storage Management)',
    notes: 'Low-latency SAN block storage with ASM redundancy.',
    contractRef: 'SAN Master Contract',
    moderate: {
      nodeCount: 1,
      vcpuPerNode: 0,
      ramGBPerNode: 0,
      storageGBPerNode: 1500,
      storageClass: 'Tier 2',
      computeType: 'shared',
      backupEnabled: true,
      licenseKind: 'free'
    },
    critical: {
      nodeCount: 1,
      vcpuPerNode: 0,
      ramGBPerNode: 0,
      storageGBPerNode: 5500,
      storageClass: 'Tier 3',
      computeType: 'shared',
      backupEnabled: true,
      licenseKind: 'enterprise'
    }
  },
  {
    id: 'dist_cache',
    no: 6,
    layer: 'Distributed Cache',
    component: 'Redis In-Memory Cluster',
    icon: '⚡',
    licenseFree: 'Redis Community OSS',
    licenseEnterprise: 'Redis Enterprise Software (Active-Active)',
    notes: 'High-speed session storage & API response caching.',
    contractRef: 'TIF.DDL/GH.2984/2025',
    moderate: {
      nodeCount: 2,
      vcpuPerNode: 8,
      ramGBPerNode: 16,
      storageGBPerNode: 200,
      storageClass: 'Tier 2',
      computeType: 'vm',
      backupEnabled: true,
      licenseKind: 'free'
    },
    critical: {
      nodeCount: 3,
      vcpuPerNode: 8,
      ramGBPerNode: 32,
      storageGBPerNode: 200,
      storageClass: 'Tier 2',
      computeType: 'vm',
      backupEnabled: true,
      licenseKind: 'enterprise'
    }
  },
  {
    id: 'cache_sentinel',
    no: 7,
    layer: 'Cache Sentinel',
    component: 'Redis Sentinel Quorum',
    icon: '🛡️',
    licenseFree: 'Redis Sentinel OSS',
    licenseEnterprise: 'Included in Redis Enterprise',
    notes: 'Automated failover & master election for Redis.',
    contractRef: 'Included in Redis',
    moderate: {
      nodeCount: 3,
      vcpuPerNode: 2,
      ramGBPerNode: 4,
      storageGBPerNode: 50,
      storageClass: 'Tier 2',
      computeType: 'vm',
      backupEnabled: true,
      licenseKind: 'free'
    },
    critical: {
      nodeCount: 3,
      vcpuPerNode: 4,
      ramGBPerNode: 8,
      storageGBPerNode: 100,
      storageClass: 'Tier 2',
      computeType: 'vm',
      backupEnabled: true,
      licenseKind: 'enterprise'
    }
  },
  {
    id: 'message_bus',
    no: 8,
    layer: 'Message Bus',
    component: 'Confluent / Apache Kafka',
    icon: '📨',
    licenseFree: 'Apache Kafka Community OSS',
    licenseEnterprise: 'Confluent Enterprise Platform (Kora Engine)',
    notes: 'Event-driven backbone for transaction streaming.',
    contractRef: 'TIF.DDL/GH.4412/2025',
    moderate: {
      nodeCount: 3,
      vcpuPerNode: 8,
      ramGBPerNode: 32,
      storageGBPerNode: 500,
      storageClass: 'Tier 2',
      computeType: 'vm',
      backupEnabled: true,
      licenseKind: 'free'
    },
    critical: {
      nodeCount: 3,
      vcpuPerNode: 16,
      ramGBPerNode: 64,
      storageGBPerNode: 1000,
      storageClass: 'Tier 2',
      computeType: 'vm',
      backupEnabled: true,
      licenseKind: 'enterprise'
    }
  },
  {
    id: 'log_analytics',
    no: 9,
    layer: 'Log Analytics',
    component: 'Elasticsearch Analytics Node',
    icon: '🔍',
    licenseFree: 'OpenSearch / Elastic Basic',
    licenseEnterprise: 'Elasticsearch Enterprise Platinum Subscription',
    notes: 'Hot-warm audit trail search & real-time analytics.',
    contractRef: 'TIF.DDL/GH.8821/2025',
    moderate: {
      nodeCount: 2,
      vcpuPerNode: 8,
      ramGBPerNode: 32,
      storageGBPerNode: 500,
      storageClass: 'Tier 2',
      computeType: 'vm',
      backupEnabled: true,
      licenseKind: 'free'
    },
    critical: {
      nodeCount: 3,
      vcpuPerNode: 16,
      ramGBPerNode: 64,
      storageGBPerNode: 2000,
      storageClass: 'Tier 2',
      computeType: 'vm',
      backupEnabled: true,
      licenseKind: 'enterprise'
    }
  },
  {
    id: 'log_shipper',
    no: 10,
    layer: 'Log Shipper',
    component: 'Logstash Pipeline Shipper',
    icon: '🚚',
    licenseFree: 'Logstash OSS',
    licenseEnterprise: 'Included in Elastic Platinum',
    notes: 'Log ingest pipeline, parsing, and normalization buffer.',
    contractRef: 'Included in Elastic',
    moderate: {
      nodeCount: 1,
      vcpuPerNode: 4,
      ramGBPerNode: 8,
      storageGBPerNode: 100,
      storageClass: 'Tier 2',
      computeType: 'vm',
      backupEnabled: true,
      licenseKind: 'free'
    },
    critical: {
      nodeCount: 2,
      vcpuPerNode: 8,
      ramGBPerNode: 16,
      storageGBPerNode: 200,
      storageClass: 'Tier 2',
      computeType: 'vm',
      backupEnabled: true,
      licenseKind: 'enterprise'
    }
  },
  {
    id: 'log_dashboard',
    no: 11,
    layer: 'Log Dashboard',
    component: 'Kibana Observability Portal',
    icon: '📊',
    licenseFree: 'OpenSearch Dashboards / Kibana Basic',
    licenseEnterprise: 'Included in Elastic Platinum',
    notes: 'Visualization dashboards and alerting gateway.',
    contractRef: 'Included in Elastic',
    moderate: {
      nodeCount: 1,
      vcpuPerNode: 4,
      ramGBPerNode: 8,
      storageGBPerNode: 100,
      storageClass: 'Tier 2',
      computeType: 'vm',
      backupEnabled: true,
      licenseKind: 'free'
    },
    critical: {
      nodeCount: 2,
      vcpuPerNode: 8,
      ramGBPerNode: 16,
      storageGBPerNode: 200,
      storageClass: 'Tier 2',
      computeType: 'vm',
      backupEnabled: true,
      licenseKind: 'enterprise'
    }
  },
  {
    id: 'doc_storage',
    no: 12,
    layer: 'Document Storage',
    component: 'SFTP / Object Storage / BDM',
    icon: '📁',
    licenseFree: 'MinIO Community (OSS)',
    licenseEnterprise: 'MinIO Enterprise / Bank Mandiri Doc Storage',
    notes: 'Long-term electronic customer statement and KYC store.',
    contractRef: 'ECM Storage Master',
    moderate: {
      nodeCount: 1,
      vcpuPerNode: 4,
      ramGBPerNode: 8,
      storageGBPerNode: 500,
      storageClass: 'Tier 2',
      computeType: 'vm',
      backupEnabled: true,
      licenseKind: 'free'
    },
    critical: {
      nodeCount: 2,
      vcpuPerNode: 8,
      ramGBPerNode: 16,
      storageGBPerNode: 5000,
      storageClass: 'Tier 2',
      computeType: 'vm',
      backupEnabled: true,
      licenseKind: 'enterprise'
    }
  },
  {
    id: 'batch_engine',
    no: 13,
    layer: 'Batch Engine',
    component: 'BMC Control-M / OCP Cron',
    icon: '⏱️',
    licenseFree: 'Kubernetes CronJob (OSS)',
    licenseEnterprise: 'BMC Control-M Workload Automation',
    notes: 'End-of-day batch settlement and clearing scheduler.',
    contractRef: 'TIF.DDL/GH.5510/2025',
    moderate: {
      nodeCount: 1,
      vcpuPerNode: 4,
      ramGBPerNode: 8,
      storageGBPerNode: 100,
      storageClass: 'Tier 2',
      computeType: 'vm',
      backupEnabled: true,
      licenseKind: 'free'
    },
    critical: {
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
]
