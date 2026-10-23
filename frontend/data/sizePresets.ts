// ─────────────────────────────────────────────
// COMPUTE SIZE PRESETS: Preset Ukuran Server Siap Pakai
// ─────────────────────────────────────────────

export interface SizePresetOption {
  id: 'small' | 'medium' | 'standard' | 'standard_plus' | 'large' | 'xlarge' | 'custom'
  name: string
  sublabel: string
  vcpu: number
  ram: number
  altRam?: number
  recommendedFor: string
}

export const COMPUTE_SIZE_PRESETS: SizePresetOption[] = [
  {
    id: 'small',
    name: 'Kecil',
    sublabel: '2 vCPU · 4 GB',
    vcpu: 2,
    ram: 4,
    recommendedFor: 'Sentinel, log shipper ringan, tools internal'
  },
  {
    id: 'medium',
    name: 'Sedang',
    sublabel: '4 vCPU · 8 GB',
    vcpu: 4,
    ram: 8,
    recommendedFor: 'Service penunjang, batch engine standar'
  },
  {
    id: 'standard',
    name: 'Standar',
    sublabel: '8 vCPU · 16 GB',
    vcpu: 8,
    ram: 16,
    recommendedFor: 'Web Ingress, log portal, workload umum'
  },
  {
    id: 'standard_plus',
    name: 'Menengah+',
    sublabel: '8 vCPU · 32 GB',
    vcpu: 8,
    ram: 32,
    recommendedFor: 'OCP Router & Bastion, Redis Cluster, Memori Ekstra'
  },
  {
    id: 'large',
    name: 'Besar',
    sublabel: '16 vCPU · 64 GB',
    vcpu: 16,
    ram: 64,
    altRam: 32,
    recommendedFor: 'Kafka Message Bus, Elasticsearch Analytics'
  },
  {
    id: 'xlarge',
    name: 'Raksasa',
    sublabel: '64 vCPU · 1024 GB',
    vcpu: 64,
    ram: 1024,
    recommendedFor: 'OpenShift Worker, Oracle Core Database'
  }
]
