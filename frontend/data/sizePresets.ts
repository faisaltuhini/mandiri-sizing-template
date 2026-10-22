// ─────────────────────────────────────────────
// COMPUTE SIZE PRESETS: Preset Ukuran Server Siap Pakai
// ─────────────────────────────────────────────

export interface SizePresetOption {
  id: 'small' | 'medium' | 'standard' | 'large' | 'xlarge' | 'custom'
  name: string
  sublabel: string
  vcpu: number
  ram: number
  recommendedFor: string
}

export const COMPUTE_SIZE_PRESETS: SizePresetOption[] = [
  {
    id: 'small',
    name: 'Kecil',
    sublabel: '2 vCPU · 4 GB RAM',
    vcpu: 2,
    ram: 4,
    recommendedFor: 'Sentinel, log shipper ringan, tools internal'
  },
  {
    id: 'medium',
    name: 'Sedang',
    sublabel: '4 vCPU · 8 GB RAM',
    vcpu: 4,
    ram: 8,
    recommendedFor: 'Service penunjang, batch engine standar'
  },
  {
    id: 'standard',
    name: 'Standar',
    sublabel: '8 vCPU · 16 GB RAM',
    vcpu: 8,
    ram: 16,
    recommendedFor: 'Web Ingress, cache cluster, workload umum'
  },
  {
    id: 'large',
    name: 'Besar',
    sublabel: '16 vCPU · 32/64 GB RAM',
    vcpu: 16,
    ram: 32,
    recommendedFor: 'Kafka Message Bus, Elasticsearch Node'
  },
  {
    id: 'xlarge',
    name: 'Raksasa (Heavy)',
    sublabel: '64 vCPU · 1024 GB RAM',
    vcpu: 64,
    ram: 1024,
    recommendedFor: 'OpenShift Worker, Oracle Core Database'
  }
]
