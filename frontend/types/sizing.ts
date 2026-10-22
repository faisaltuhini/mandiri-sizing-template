// ─────────────────────────────────────────────
// Sizing & Catalog Data Types (Bank Mandiri Sizing Suite)
// ─────────────────────────────────────────────

export interface LayerAdjusterConfig {
  nodeCount: number
  vcpuPerNode: number
  ramGBPerNode: number
  storageGBPerNode: number
  storageClass: 'Tier 2' | 'Tier 3'
  backupEnabled: boolean
  licenseKind: 'free' | 'enterprise'
  computeType: 'vm' | 'baremetal' | 'shared'
}

export interface QuoteItemCost {
  component_id: string
  component_name: string
  layer_name: string
  node_count: number
  total_vcpu: number
  total_ram_gb: number
  total_storage_gb: number
  compute_cost_annual: number
  storage_cost_annual: number
  backup_cost_annual: number
  license_cost_annual: number
  total_cost_annual: number
  license_product_name: string
  license_kind: string
  unconfirmed_prices?: string[]
}

export interface QuoteResponse {
  environment: string
  include_drc: boolean
  total_nodes: number
  total_vcpu: number
  total_ram_gb: number
  total_storage_gb: number
  subtotal_compute_annual: number
  subtotal_storage_annual: number
  subtotal_backup_annual: number
  subtotal_license_annual: number
  grand_total_annual: number
  tco_3year: number
  tco_5year: number
  items: QuoteItemCost[]
  warnings?: string[]
}

export interface PricingItemRecord {
  id: string
  price_list_version_id: string
  sku: string
  category: string
  name: string
  environment: string
  unit: string
  unit_price_idr: number | null
  is_confirmed: boolean
  source: string
  created_at?: string
  updated_at?: string
}

export interface PricingFormState {
  sku: string
  name: string
  category: string
  environment: string
  unit: string
  unit_price_idr: number | null
  is_confirmed: boolean
  source: string
}

export type AppMenu = 'catalog' | 'sizing' | 'presets' | 'licenses'

export interface LayerPresetSpec {
  nodeCount: number
  vcpuPerNode: number
  ramGBPerNode: number
  storageGBPerNode: number
  storageClass: 'Tier 2' | 'Tier 3'
  computeType: 'vm' | 'baremetal' | 'shared'
  backupEnabled: boolean
  licenseKind: 'free' | 'enterprise'
}

export interface SizingPresetItem {
  id: string
  no: number
  layer: string
  component: string
  icon?: string
  licenseFree: string
  licenseEnterprise: string
  notes?: string
  contractRef?: string
  moderate: LayerPresetSpec
  critical: LayerPresetSpec
}

