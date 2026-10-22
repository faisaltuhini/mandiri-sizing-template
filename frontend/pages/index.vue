<template>
  <div class="sketch-sizing-app">
    <!-- LEFT SIDEBAR: [Sizing] and [Katalog] Navigation (Collapsible / Minimizable) -->
    <SizingSidebar
      :current-menu="currentMenu"
      :is-collapsed="isSidebarCollapsed"
      :new-catalog-count="unregisteredPresetCatalogCount"
      @select-menu="currentMenu = $event"
      @toggle-collapse="isSidebarCollapsed = !isSidebarCollapsed"
    />

    <!-- MAIN WORKSPACE -->
    <div class="sketch-main-container">
      <!-- VIEW 1: SIZING CONFIGURATOR -->
      <template v-if="currentMenu === 'sizing'">
        <!-- TOPBAR WITH FILTER POPOVER & ACTION BUTTONS -->
        <SizingTopBar
          :tier="tierFilter"
          :license="licenseFilter"
          :is-sidebar-collapsed="isSidebarCollapsed"
          @update:tier="setTierFilter"
          @update:license="setLicenseFilter"
          @toggle-sidebar="isSidebarCollapsed = !isSidebarCollapsed"
          @open-upload="showUploadModal = true"
          @trigger-generate="triggerGenerate"
        />

        <!-- TWO-COLUMN WORKSPACE: 13 LAYERS LIST + LIVE NUMERICAL ADJUSTER & SPEC REFERENCE -->
        <SizingConfigurator
          :layer-items="layerItems"
          :active-item-id="activeItemId"
          :active-item="activeItem"
          :active-config="activeConfig"
          :active-layer-cost="activeLayerCost"
          :tier-filter="tierFilter"
          :license-filter="licenseFilter"
          :is-layer-customized="isLayerCustomized"
          @select-layer="activeItemId = $event"
          @reset-active-layer="resetActiveLayerToBaseline"
        />

        <!-- STICKY ARCHITECTURE COST SUMMARY BAR (BOTTOM) -->
        <SizingBottomBar
          :environment="environment"
          :include-d-r-c="includeDRC"
          :quote-response="quoteResponse"
          @update:environment="environment = $event"
          @update:include-d-r-c="includeDRC = $event"
          @trigger-generate="triggerGenerate"
        />
      </template>

      <!-- VIEW 2: FULL TABLE PRICE CATALOG & CRUD -->
      <template v-else-if="currentMenu === 'catalog'">
        <CatalogTableView
          :items="filteredPricingList"
          :unconfirmed-count="unconfirmedCount"
          :search="catalogSearch"
          :category-filter="catalogCategoryFilter"
          :env-filter="catalogEnvFilter"
          :is-loading="isLoadingPricing"
          :is-sidebar-collapsed="isSidebarCollapsed"
          @update:search="catalogSearch = $event"
          @update:category-filter="catalogCategoryFilter = $event"
          @update:env-filter="catalogEnvFilter = $event"
          @refresh="loadPricingItems"
          @reset-all-filters="catalogSearch = ''; catalogCategoryFilter = 'all'; catalogEnvFilter = 'all'"
          @open-add-sku="openAddPricingModal"
          @edit-item="openEditPricingModal"
          @delete-item="confirmDeletePricing"
          @toggle-sidebar="isSidebarCollapsed = !isSidebarCollapsed"
        />
      </template>

      <!-- VIEW 3: ARCHITECTURE PRESET & FORMULA CRUD -->
      <template v-else-if="currentMenu === 'presets'">
        <PresetTableView
          :preset-items="presetItems"
          :catalog-items="pricingList"
          :is-sidebar-collapsed="isSidebarCollapsed"
          @save-preset-item="handleSavePresetItem"
          @delete-preset-item="handleDeletePresetItem"
          @reset-to-official-defaults="handleResetPresetDefaults"
          @toggle-sidebar="isSidebarCollapsed = !isSidebarCollapsed"
          @update:unregistered-count="unregisteredPresetCatalogCount = $event"
        />
      </template>

      <!-- VIEW 4: SOFTWARE LICENSE RULES & METRICS CRUD -->
      <template v-else-if="currentMenu === 'licenses'">
        <LicenseRuleTableView
          :is-sidebar-collapsed="isSidebarCollapsed"
          @toggle-sidebar="isSidebarCollapsed = !isSidebarCollapsed"
        />
      </template>
    </div>

    <!-- MODALS SECTION: ADD/EDIT SKU, DELETE CONFIRM, UPLOAD & BOQ SUMMARY -->
    <SizingModals
      :show-pricing-modal="showPricingModal"
      :is-editing-pricing="isEditingPricing"
      :is-saving-pricing="isSavingPricing"
      :pricing-form="pricingForm"
      :show-delete-pricing-modal="showDeletePricingModal"
      :is-deleting-pricing="isDeletingPricing"
      :deleting-pricing-item="deletingPricingItem"
      :show-upload-modal="showUploadModal"
      :uploaded-file-name="uploadedFileName"
      :show-generate-modal="showGenerateModal"
      :project-name="projectName"
      :environment="environment"
      :include-d-r-c="includeDRC"
      :quote-response="quoteResponse"
      :copy-status="copyStatus"
      :is-downloading-excel="isDownloadingExcel"
      @close-pricing-modal="showPricingModal = false"
      @save-pricing="savePricingItem"
      @close-delete-modal="showDeletePricingModal = false"
      @execute-delete="executeDeletePricing"
      @close-upload-modal="showUploadModal = false"
      @file-change="handleFileChange"
      @file-drop="handleFileDrop"
      @process-upload="processUpload"
      @close-generate-modal="showGenerateModal = false"
      @update:project-name="projectName = $event"
      @copy-summary="copyGeneratedSummary"
      @download-excel="downloadOfficialExcel"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch, onMounted } from 'vue'
import {
  BASELINE_TECH_STACK_ITEMS,
  type TechStackItem
} from '~/data/techStackBaselineData'
import type {
  LayerAdjusterConfig,
  QuoteItemCost,
  QuoteResponse,
  PricingItemRecord,
  PricingFormState
} from '~/types/sizing'
import SizingSidebar from '~/components/sizing/SizingSidebar.vue'
import SizingTopBar from '~/components/sizing/SizingTopBar.vue'
import SizingConfigurator from '~/components/sizing/SizingConfigurator.vue'
import SizingBottomBar from '~/components/sizing/SizingBottomBar.vue'
import CatalogTableView from '~/components/sizing/CatalogTableView.vue'
import PresetTableView from '~/components/sizing/PresetTableView.vue'
import LicenseRuleTableView from '~/components/sizing/LicenseRuleTableView.vue'
import SizingModals from '~/components/sizing/SizingModals.vue'
import { useLocalOpen } from '~/composables/useLocalOpen'
import { useBackendBaseUrl } from '~/composables/useSizingApi'
import { useSizingIndexedDB } from '~/composables/useSizingIndexedDB'
import { DEFAULT_PRESET_ITEMS } from '~/data/defaultPresetItems'
import type { AppMenu, SizingPresetItem } from '~/types/sizing'

const { showToast } = useLocalOpen()
const { setItem: setDbItem, getItem: getDbItem } = useSizingIndexedDB()
const apiBase = useBackendBaseUrl()

// ==========================================
// ROUTER & QUERY PARAM PERSISTENCE
// ==========================================
const route = useRoute()
const router = useRouter()

// Initialize from URL Query Params (e.g. ?tab=presets | ?tab=catalog | ?tab=sizing | ?tab=licenses)
const validTabs: AppMenu[] = ['catalog', 'sizing', 'presets', 'licenses']
const initialTab = validTabs.includes(route.query.tab as AppMenu) ? (route.query.tab as AppMenu) : 'sizing'
const currentMenu = ref<AppMenu>(initialTab)
const isSidebarCollapsed = ref<boolean>(route.query.collapsed === '1')
const unregisteredPresetCatalogCount = ref<number>(0)

// Sync state changes to URL Query Params
function syncQueryParams() {
  const newQuery: Record<string, string | undefined> = {
    ...route.query,
    tab: currentMenu.value,
    collapsed: isSidebarCollapsed.value ? '1' : undefined,
  }

  // Remove undefined keys to keep URL clean
  Object.keys(newQuery).forEach(key => {
    if (newQuery[key] === undefined) delete newQuery[key]
  })

  router.replace({ query: newQuery })
}

watch([currentMenu, isSidebarCollapsed], () => {
  syncQueryParams()
})

// Listen to browser Back / Forward buttons
watch(() => route.query, (q) => {
  if (q.tab && validTabs.includes(q.tab as AppMenu)) {
    if (currentMenu.value !== q.tab) {
      currentMenu.value = q.tab as AppMenu
    }
  } else if (currentMenu.value !== 'sizing') {
    currentMenu.value = 'sizing'
  }

  const shouldBeCollapsed = q.collapsed === '1'
  if (isSidebarCollapsed.value !== shouldBeCollapsed) {
    isSidebarCollapsed.value = shouldBeCollapsed
  }
})



// ==========================================
// PRESETS & SIZING CONFIGURATOR STATE
// ==========================================
const presetItems = ref<SizingPresetItem[]>(JSON.parse(JSON.stringify(DEFAULT_PRESET_ITEMS)))

// Derive active layer items from presetItems so any newly added layer appears automatically
const layerItems = computed<TechStackItem[]>(() => {
  return presetItems.value.map(p => ({
    id: p.id,
    no: p.no,
    layer: p.layer,
    component: p.component,
    icon: p.icon || '📦',
    moderate: {
      nodes: `${p.moderate.nodeCount} Node`,
      type: p.moderate.computeType === 'baremetal' ? 'Baremetal Dedicated' : (p.moderate.computeType === 'shared' ? 'Shared SAN' : 'Virtual Machine (VM)'),
      spec: `${p.moderate.vcpuPerNode} vCPU / ${p.moderate.ramGBPerNode}GB RAM / ${p.moderate.storageGBPerNode}GB Storage`
    },
    critical: {
      nodes: `${p.critical.nodeCount} Node`,
      type: p.critical.computeType === 'baremetal' ? 'Baremetal Dedicated' : (p.critical.computeType === 'shared' ? 'Shared SAN' : 'Virtual Machine (VM)'),
      spec: `${p.critical.vcpuPerNode} vCPU / ${p.critical.ramGBPerNode}GB RAM / ${p.critical.storageGBPerNode}GB Storage`
    },
    license: {
      free: p.licenseFree,
      enterprise: p.licenseEnterprise
    },
    notes: p.notes || '',
    contractRef: p.contractRef || ''
  }))
})

const activeItemId = ref(DEFAULT_PRESET_ITEMS[0].id)

// Topbar filtering controls (Critical vs Moderat, Free vs License)
const tierFilter = ref<'critical' | 'moderate'>('critical')
const licenseFilter = ref<'free' | 'license'>('license')

// Architecture global scope
const environment = ref<'Production' | 'Dev / Non-Prod'>('Production')
const includeDRC = ref<boolean>(true)
const projectName = ref('Core Banking Architecture Livin\'')

// Modals
const showUploadModal = ref(false)
const showGenerateModal = ref(false)
const uploadedFileName = ref('')
const copyStatus = ref('')
const isDownloadingExcel = ref(false)

function parseBaselineSpec(item: TechStackItem | { id: string; no: number }, tier: 'moderate' | 'critical'): LayerAdjusterConfig {
  const isCrit = tier === 'critical'
  const preset = presetItems.value.find(p => p.id === item.id)

  if (preset) {
    const spec = isCrit ? preset.critical : preset.moderate
    return {
      nodeCount: spec.nodeCount,
      vcpuPerNode: spec.vcpuPerNode,
      ramGBPerNode: spec.ramGBPerNode,
      storageGBPerNode: spec.storageGBPerNode,
      storageClass: spec.storageClass,
      backupEnabled: spec.backupEnabled,
      licenseKind: spec.licenseKind,
      computeType: spec.computeType
    }
  }

  // Fallback defaults if preset not found
  return {
    nodeCount: isCrit ? 2 : 1,
    vcpuPerNode: isCrit ? 8 : 4,
    ramGBPerNode: isCrit ? 16 : 8,
    storageGBPerNode: isCrit ? 200 : 100,
    storageClass: 'Tier 2',
    backupEnabled: true,
    licenseKind: isCrit ? 'enterprise' : 'free',
    computeType: 'vm'
  }
}

const layerConfigs = reactive<Record<string, LayerAdjusterConfig>>({})

// Initialize config map from presets
for (const item of DEFAULT_PRESET_ITEMS) {
  layerConfigs[item.id] = parseBaselineSpec(item, 'critical')
}

// Active component
const activeItem = computed<TechStackItem>(() => {
  return layerItems.value.find(i => i.id === activeItemId.value) || layerItems.value[0]
})

const activeConfig = computed<LayerAdjusterConfig>({
  get() {
    if (!layerConfigs[activeItemId.value]) {
      layerConfigs[activeItemId.value] = parseBaselineSpec(activeItem.value, tierFilter.value)
    }
    return layerConfigs[activeItemId.value]
  },
  set(val) {
    layerConfigs[activeItemId.value] = val
  }
})

function isLayerCustomized(id: string): boolean {
  const item = layerItems.value.find(i => i.id === id)
  if (!item || !layerConfigs[id]) return false
  const baseline = parseBaselineSpec(item, tierFilter.value)
  const current = layerConfigs[id]
  return (
    current.nodeCount !== baseline.nodeCount ||
    current.vcpuPerNode !== baseline.vcpuPerNode ||
    current.ramGBPerNode !== baseline.ramGBPerNode ||
    current.storageGBPerNode !== baseline.storageGBPerNode ||
    current.storageClass !== baseline.storageClass ||
    current.licenseKind !== baseline.licenseKind ||
    current.computeType !== baseline.computeType ||
    current.backupEnabled !== baseline.backupEnabled
  )
}

function resetActiveLayerToBaseline() {
  const targetTier = tierFilter.value
  layerConfigs[activeItemId.value] = parseBaselineSpec(activeItem.value, targetTier)
  layerConfigs[activeItemId.value].licenseKind = licenseFilter.value === 'license' ? 'enterprise' : 'free'
  showToast(`Konfigurasi Layer ${activeItem.value.no} dikembalikan ke baseline ${targetTier.toUpperCase()}`, 'info')
}

function setTierFilter(tier: 'critical' | 'moderate') {
  tierFilter.value = tier
  for (const item of layerItems.value) {
    const parsed = parseBaselineSpec(item, tier)
    layerConfigs[item.id].nodeCount = parsed.nodeCount
    layerConfigs[item.id].vcpuPerNode = parsed.vcpuPerNode
    layerConfigs[item.id].ramGBPerNode = parsed.ramGBPerNode
    layerConfigs[item.id].storageGBPerNode = parsed.storageGBPerNode
    layerConfigs[item.id].storageClass = parsed.storageClass
    layerConfigs[item.id].computeType = parsed.computeType
  }
}

// Preset CRUD Event Handlers
async function handleSavePresetItem(item: SizingPresetItem) {
  const idx = presetItems.value.findIndex(p => p.id === item.id)
  if (idx >= 0) {
    presetItems.value[idx] = item
    showToast(`Preset Layer "${item.layer}" berhasil diperbarui!`, 'success')
  } else {
    presetItems.value.push(item)
    showToast(`Preset Layer Baru "${item.layer}" berhasil ditambahkan!`, 'success')
  }

  // Update layerConfigs baseline
  layerConfigs[item.id] = parseBaselineSpec(item, tierFilter.value)

  // Persist presets to IndexedDB
  await setDbItem('sizing_presets', JSON.parse(JSON.stringify(presetItems.value)))
  triggerQuoteCalculation()
}

async function handleDeletePresetItem(itemId: string) {
  presetItems.value = presetItems.value.filter(p => p.id !== itemId)
  delete layerConfigs[itemId]
  if (activeItemId.value === itemId && presetItems.value.length > 0) {
    activeItemId.value = presetItems.value[0].id
  }
  showToast('Preset layer kustom berhasil dihapus!', 'info')
  await setDbItem('sizing_presets', JSON.parse(JSON.stringify(presetItems.value)))
  triggerQuoteCalculation()
}

async function handleResetPresetDefaults() {
  if (confirm('Kembalikan semua preset dan formula layer ke standar resmi pabrik Enterprise Architecture?')) {
    presetItems.value = JSON.parse(JSON.stringify(DEFAULT_PRESET_ITEMS))
    for (const item of DEFAULT_PRESET_ITEMS) {
      layerConfigs[item.id] = parseBaselineSpec(item, tierFilter.value)
    }
    await setDbItem('sizing_presets', JSON.parse(JSON.stringify(DEFAULT_PRESET_ITEMS)))
    showToast('Semua preset dikembalikan ke standar resmi Enterprise Architecture!', 'success')
    triggerQuoteCalculation()
  }
}

function setLicenseFilter(lic: 'free' | 'license') {
  licenseFilter.value = lic
  for (const item of layerItems.value) {
    if (layerConfigs[item.id]) {
      layerConfigs[item.id].licenseKind = lic === 'license' ? 'enterprise' : 'free'
    }
  }
}

const quoteResponse = ref<QuoteResponse | null>(null)

// Selected layer item cost
const activeLayerCost = computed<QuoteItemCost | undefined>(() => {
  if (!quoteResponse.value || !quoteResponse.value.items) return undefined
  return (
    quoteResponse.value.items.find(i => i.component_id === activeItemId.value) ||
    quoteResponse.value.items.find(i => i.component_id === `layer-${activeItem.value.no}`)
  )
})

function buildQuoteRequest() {
  const items = layerItems.value.map(l => {
    const cfg = layerConfigs[l.id] || parseBaselineSpec(l, 'moderate')
    return {
      component_id: `layer-${l.no}`,
      tier: tierFilter.value === 'critical' ? 'critical' : 'moderate',
      compute_type: cfg.computeType,
      node_count: cfg.nodeCount,
      vcpu_per_node: cfg.vcpuPerNode,
      ram_gb_per_node: cfg.ramGBPerNode,
      storage_gb_per_node: cfg.storageGBPerNode,
      storage_class: cfg.storageClass,
      backup_enabled: cfg.backupEnabled,
      license_kind: cfg.licenseKind,
    }
  })

  return {
    environment: environment.value,
    include_drc: includeDRC.value,
    items,
  }
}

let quoteTimer: ReturnType<typeof setTimeout> | null = null

function triggerQuoteCalculation() {
  if (quoteTimer) clearTimeout(quoteTimer)
  quoteTimer = setTimeout(async () => {
    try {
      const payload = buildQuoteRequest()
      const res = await $fetch<QuoteResponse>(`${apiBase}/api/v1/quotes/calculate`, {
        method: 'POST',
        body: payload,
      })
      quoteResponse.value = res
    } catch (err) {
      console.warn('Backend quote calculation fallback:', err)
      runOptimisticQuoteCalculation()
    }
  }, 120)
}

function runOptimisticQuoteCalculation() {
  const isDev = environment.value.toLowerCase().includes('dev')
  const vcpuRate = isDev ? 7150000 : 8250000
  const stgT2Rate = isDev ? 28600 : 33000
  const stgT3Rate = 48000
  const nbuRate = 27000
  const bmRate = 794000000

  let totalNodes = 0
  let totalVCPU = 0
  let totalRAM = 0
  let totalStorage = 0
  let subCompute = 0
  let subStorage = 0
  let subBackup = 0
  let subLicense = 0

  const items: QuoteItemCost[] = []
  const warnings: string[] = []

  for (const l of layerItems.value) {
    const cfg = layerConfigs[l.id] || parseBaselineSpec(l, 'moderate')
    const itemNodes = cfg.nodeCount
    const itemVcpu = itemNodes * cfg.vcpuPerNode
    const itemRam = itemNodes * cfg.ramGBPerNode
    const itemStorage = itemNodes * cfg.storageGBPerNode

    let computeCost = 0
    if (cfg.computeType === 'baremetal') {
      computeCost = itemNodes * bmRate
    } else if (cfg.computeType === 'vm') {
      computeCost = itemVcpu * vcpuRate
    }

    const stgRate = cfg.storageClass === 'Tier 3' ? stgT3Rate : stgT2Rate
    const storageCost = itemStorage * stgRate
    const backupCost = cfg.backupEnabled ? itemStorage * nbuRate : 0

    let licenseCost = 0
    const unconfirmedPrices: string[] = []
    let licProductName = l.license.free

    if (cfg.licenseKind === 'enterprise') {
      licProductName = l.license.enterprise
      if (l.no === 2) {
        licenseCost = itemNodes * (isDev ? 134000000 : 240000000)
      } else if (l.no === 6) {
        licenseCost = itemNodes * (isDev ? 75000000 : 125000000)
      } else if (l.no === 1 || l.no === 4 || l.no === 8) {
        licenseCost = 0
        const msg = `${licProductName} (${l.component}): Harga belum dikonfirmasi STA`
        unconfirmedPrices.push(msg)
        if (!warnings.includes(msg)) warnings.push(msg)
      }
    }

    const itemTotal = computeCost + storageCost + backupCost + licenseCost

    totalNodes += itemNodes
    totalVCPU += itemVcpu
    totalRAM += itemRam
    totalStorage += itemStorage

    subCompute += computeCost
    subStorage += storageCost
    subBackup += backupCost
    subLicense += licenseCost

    items.push({
      component_id: l.id,
      component_name: l.component,
      layer_name: l.layer,
      node_count: itemNodes,
      total_vcpu: itemVcpu,
      total_ram_gb: itemRam,
      total_storage_gb: itemStorage,
      compute_cost_annual: computeCost,
      storage_cost_annual: storageCost,
      backup_cost_annual: backupCost,
      license_cost_annual: licenseCost,
      total_cost_annual: itemTotal,
      license_product_name: licProductName,
      license_kind: cfg.licenseKind,
      unconfirmed_prices: unconfirmedPrices,
    })
  }

  const drcMult = includeDRC.value ? 2 : 1
  subCompute *= drcMult
  subStorage *= drcMult
  subBackup *= drcMult

  const grandTotal = subCompute + subStorage + subBackup + subLicense

  quoteResponse.value = {
    environment: environment.value,
    include_drc: includeDRC.value,
    total_nodes: totalNodes * drcMult,
    total_vcpu: totalVCPU * drcMult,
    total_ram_gb: totalRAM * drcMult,
    total_storage_gb: totalStorage * drcMult,
    subtotal_compute_annual: subCompute,
    subtotal_storage_annual: subStorage,
    subtotal_backup_annual: subBackup,
    subtotal_license_annual: subLicense,
    grand_total_annual: grandTotal,
    tco_3year: grandTotal * 3,
    tco_5year: grandTotal * 5,
    items,
    warnings,
  }
}

// IndexedDB Persistence helper: Auto-save sizing configuration state
let isRestoringFromDB = false

async function saveSizingStateToIndexedDB() {
  if (isRestoringFromDB) return
  const dataToSave = {
    tierFilter: tierFilter.value,
    licenseFilter: licenseFilter.value,
    environment: environment.value,
    includeDRC: includeDRC.value,
    activeItemId: activeItemId.value,
    projectName: projectName.value,
    layerConfigs: JSON.parse(JSON.stringify(layerConfigs)),
  }
  await setDbItem('sizing_state', dataToSave)
}

watch(
  [layerConfigs, environment, includeDRC, tierFilter, licenseFilter, activeItemId, projectName],
  () => {
    triggerQuoteCalculation()
    saveSizingStateToIndexedDB()
  },
  { deep: true }
)

// ==========================================
// PRICE CATALOG FULL TABLE & CRUD STATE
// ==========================================
const pricingList = ref<PricingItemRecord[]>([])
const isLoadingPricing = ref(false)
const catalogSearch = ref('')
const catalogCategoryFilter = ref<string>('all')
const catalogEnvFilter = ref<string>('all')

const showPricingModal = ref(false)
const isEditingPricing = ref(false)
const isSavingPricing = ref(false)
const editingPricingId = ref<string | null>(null)

const pricingForm = reactive<PricingFormState>({
  sku: '',
  name: '',
  category: 'compute',
  environment: 'Production',
  unit: 'vCPU / Thn',
  unit_price_idr: 0,
  is_confirmed: true,
  source: 'Template Sizing.xlsx!Pricing',
})

const showDeletePricingModal = ref(false)
const isDeletingPricing = ref(false)
const deletingPricingItem = ref<PricingItemRecord | null>(null)

async function loadPricingItems() {
  isLoadingPricing.value = true
  try {
    const res = await $fetch<{ total: number; data: PricingItemRecord[] }>(`${apiBase}/api/v1/pricing/items`)
    if (res && res.data) {
      pricingList.value = res.data
    }
  } catch (err) {
    console.error('Failed to load pricing items from database:', err)
  } finally {
    isLoadingPricing.value = false
  }
}

const filteredPricingList = computed(() => {
  return pricingList.value.filter(item => {
    if (catalogCategoryFilter.value !== 'all' && item.category !== catalogCategoryFilter.value) {
      return false
    }
    if (catalogEnvFilter.value !== 'all' && item.environment !== catalogEnvFilter.value) {
      return false
    }
    if (catalogSearch.value.trim()) {
      const q = catalogSearch.value.toLowerCase().trim()
      const matchSku = item.sku.toLowerCase().includes(q)
      const matchName = item.name.toLowerCase().includes(q)
      const matchSource = (item.source || '').toLowerCase().includes(q)
      if (!matchSku && !matchName && !matchSource) return false
    }
    return true
  })
})

const unconfirmedCount = computed(() => {
  return pricingList.value.filter(p => !p.is_confirmed || p.unit_price_idr === null).length
})

function openAddPricingModal() {
  isEditingPricing.value = false
  editingPricingId.value = null
  pricingForm.sku = ''
  pricingForm.name = ''
  pricingForm.category = 'compute'
  pricingForm.environment = 'Production'
  pricingForm.unit = 'vCPU / Thn'
  pricingForm.unit_price_idr = 0
  pricingForm.is_confirmed = true
  pricingForm.source = 'Template Sizing.xlsx!Pricing'
  showPricingModal.value = true
}

function openEditPricingModal(item: PricingItemRecord) {
  isEditingPricing.value = true
  editingPricingId.value = item.id
  pricingForm.sku = item.sku
  pricingForm.name = item.name
  pricingForm.category = item.category
  pricingForm.environment = item.environment
  pricingForm.unit = item.unit
  pricingForm.unit_price_idr = item.unit_price_idr
  pricingForm.is_confirmed = item.is_confirmed
  pricingForm.source = item.source || ''
  showPricingModal.value = true
}

async function savePricingItem() {
  if (!pricingForm.sku || !pricingForm.name) {
    showToast('SKU dan Nama Item wajib diisi', 'error')
    return
  }

  isSavingPricing.value = true
  try {
    const payload = {
      sku: pricingForm.sku.trim(),
      name: pricingForm.name.trim(),
      category: pricingForm.category,
      environment: pricingForm.environment,
      unit: pricingForm.unit.trim(),
      unit_price_idr: pricingForm.is_confirmed ? pricingForm.unit_price_idr : null,
      is_confirmed: pricingForm.is_confirmed,
      source: pricingForm.source.trim(),
    }

    if (isEditingPricing.value && editingPricingId.value) {
      await $fetch(`${apiBase}/api/v1/pricing/items/${editingPricingId.value}`, {
        method: 'PUT',
        body: payload,
      })
      showToast(`SKU ${payload.sku} berhasil diperbarui di database!`, 'success')
    } else {
      await $fetch(`${apiBase}/api/v1/pricing/items`, {
        method: 'POST',
        body: payload,
      })
      showToast(`SKU ${payload.sku} berhasil ditambahkan ke database!`, 'success')
    }

    showPricingModal.value = false
    await loadPricingItems()
    triggerQuoteCalculation()
  } catch (err: any) {
    console.error('Error saving pricing item:', err)
    showToast(`Gagal menyimpan: ${err.message || 'Server error'}`, 'error')
  } finally {
    isSavingPricing.value = false
  }
}

function confirmDeletePricing(item: PricingItemRecord) {
  deletingPricingItem.value = item
  showDeletePricingModal.value = true
}

async function executeDeletePricing() {
  if (!deletingPricingItem.value) return
  isDeletingPricing.value = true
  try {
    await $fetch(`${apiBase}/api/v1/pricing/items/${deletingPricingItem.value.id}`, {
      method: 'DELETE',
    })
    showToast(`SKU ${deletingPricingItem.value.sku} berhasil dihapus dari database!`, 'info')
    showDeletePricingModal.value = false
    await loadPricingItems()
    triggerQuoteCalculation()
  } catch (err: any) {
    console.error('Error deleting pricing item:', err)
    showToast(`Gagal menghapus: ${err.message || 'Server error'}`, 'error')
  } finally {
    isDeletingPricing.value = false
  }
}

// Lifecycle onMounted: Restore from IndexedDB, load initial data & ensure query param exists
onMounted(async () => {
  if (!route.query.tab) {
    syncQueryParams()
  }

  // Restore sizing configurator choices & presets from IndexedDB
  if (import.meta.client) {
    try {
      // 1. Restore Custom Presets if any
      const savedPresets = await getDbItem<SizingPresetItem[]>('sizing_presets')
      if (savedPresets && Array.isArray(savedPresets) && savedPresets.length > 0) {
        presetItems.value = savedPresets
        for (const item of savedPresets) {
          if (!layerConfigs[item.id]) {
            layerConfigs[item.id] = parseBaselineSpec(item, tierFilter.value)
          }
        }
      }

      // 2. Restore Live Sizing Configuration
      const savedState = await getDbItem<{
        tierFilter?: 'critical' | 'moderate'
        licenseFilter?: 'free' | 'license'
        environment?: 'Production' | 'Dev / Non-Prod'
        includeDRC?: boolean
        activeItemId?: string
        projectName?: string
        layerConfigs?: Record<string, LayerAdjusterConfig>
      }>('sizing_state')

      if (savedState) {
        isRestoringFromDB = true
        if (savedState.tierFilter) tierFilter.value = savedState.tierFilter
        if (savedState.licenseFilter) licenseFilter.value = savedState.licenseFilter
        if (savedState.environment) environment.value = savedState.environment
        if (typeof savedState.includeDRC === 'boolean') includeDRC.value = savedState.includeDRC
        if (savedState.activeItemId) activeItemId.value = savedState.activeItemId
        if (savedState.projectName) projectName.value = savedState.projectName
        if (savedState.layerConfigs) {
          Object.assign(layerConfigs, savedState.layerConfigs)
        }
        setTimeout(() => {
          isRestoringFromDB = false
        }, 100)
      }
    } catch (e) {
      console.warn('Failed to restore from IndexedDB:', e)
      isRestoringFromDB = false
    }
  }

  await loadPricingItems()
  triggerQuoteCalculation()
})


function triggerGenerate() {
  showGenerateModal.value = true
}

function handleFileChange(event: Event) {
  const target = event.target as HTMLInputElement
  if (target.files && target.files.length > 0) {
    uploadedFileName.value = target.files[0].name
  }
}

function handleFileDrop(event: DragEvent) {
  if (event.dataTransfer && event.dataTransfer.files.length > 0) {
    uploadedFileName.value = event.dataTransfer.files[0].name
  }
}

function processUpload() {
  showToast(`Berkas ${uploadedFileName.value} berhasil diimpor & diproses!`, 'success')
  showUploadModal.value = false
}

function formatSimpleIDR(val?: number | null): string {
  if (val === null || val === undefined || isNaN(val)) return '0'
  return new Intl.NumberFormat('id-ID').format(Math.round(val))
}

async function copyGeneratedSummary() {
  if (!quoteResponse.value) return
  const q = quoteResponse.value
  const summaryText = `ENTERPRISE - ESTIMASI SIZING & BOQ ARSITEKTUR
Project: ${projectName.value}
Environment: ${q.environment} | DRC: ${q.include_drc ? 'Active-Active (2x)' : 'Single DC'}
Total Nodes: ${q.total_nodes} Node | Total vCPU: ${q.total_vcpu} Core | Total Storage: ${(q.total_storage_gb / 1000).toFixed(1)} TB

RINCIAN BIAYA TAHUNAN:
- Compute Subtotal: Rp ${formatSimpleIDR(q.subtotal_compute_annual)} / thn
- Storage Subtotal: Rp ${formatSimpleIDR(q.subtotal_storage_annual)} / thn
- Lisensi Subtotal: Rp ${formatSimpleIDR(q.subtotal_license_annual)} / thn
- TOTAL ANNUAL OPEX: Rp ${formatSimpleIDR(q.grand_total_annual)} / thn

ESTIMASI TCO:
- 3-Year TCO: Rp ${formatSimpleIDR(q.tco_3year)}
- 5-Year TCO: Rp ${formatSimpleIDR(q.tco_5year)}
`

  try {
    await navigator.clipboard.writeText(summaryText)
    copyStatus.value = 'Tersalin ke Clipboard!'
    setTimeout(() => {
      copyStatus.value = ''
    }, 2500)
    showToast('Ringkasan BOQ arsitektur berhasil disalin!', 'success')
  } catch (e) {
    showToast('Gagal menyalin ringkasan', 'error')
  }
}

async function downloadOfficialExcel() {
  isDownloadingExcel.value = true
  try {
    const payload = {
      project_name: projectName.value,
      quote: buildQuoteRequest(),
    }

    const response = await fetch(`${apiBase}/api/v1/quotes/export-excel`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(payload),
    })

    if (!response.ok) {
      throw new Error(`Server returned HTTP ${response.status}`)
    }

    const blob = await response.blob()
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `Enterprise_Sizing_BOQ_${new Date().toISOString().slice(0, 10)}.xlsx`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)

    showToast('Template Sizing (.xlsx) resmi Enterprise Architecture berhasil diunduh!', 'success')
  } catch (err: any) {
    console.error('Download Excel error:', err)
    showToast(`Gagal mengunduh Excel: ${err.message || 'Koneksi error'}`, 'error')
  } finally {
    isDownloadingExcel.value = false
  }
}
</script>

<style scoped>
.sketch-sizing-app {
  display: flex;
  min-height: 100vh;
  background-color: #F8FAFC;
  color: #0F172A;
  font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
}

.sketch-main-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  height: 100vh;
  overflow: hidden;
}
</style>
