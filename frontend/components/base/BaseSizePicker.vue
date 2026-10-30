<template>
  <div class="base-size-picker">
    <div class="size-pill-group">
      <button
        v-for="p in presets"
        :key="p.id"
        type="button"
        class="size-pill-btn"
        :class="{
          active: !forceCustomMode && currentPresetId === p.id,
          recommended: p.id === recommendedPresetId
        }"
        @click="selectPreset(p)"
      >
        <div class="size-top-row">
          <span class="size-name">{{ p.name }}</span>
          <span v-if="p.id === recommendedPresetId" class="rec-badge-inline" title="Rekomendasi Standar Arsitektur">Standar EA</span>
          <span v-if="!forceCustomMode && currentPresetId === p.id" class="active-check">
            <svg width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
              <polyline points="20 6 9 17 4 12"></polyline>
            </svg>
          </span>
        </div>
        <span class="size-sub">{{ p.sublabel }}</span>
      </button>

      <button
        type="button"
        class="size-pill-btn btn-custom"
        :class="{ active: forceCustomMode || currentPresetId === 'custom' }"
        @click="selectCustom"
      >
        <div class="size-top-row">
          <span class="size-name">Kustom</span>
          <span v-if="forceCustomMode || currentPresetId === 'custom'" class="active-check">
            <svg width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
              <polyline points="20 6 9 17 4 12"></polyline>
            </svg>
          </span>
        </div>
        <span class="size-sub">Atur Angka Manual</span>
      </button>
    </div>

    <!-- Active size summary statement -->
    <div class="active-size-statement">
      <span class="statement-icon">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="12" cy="12" r="10"/>
          <line x1="12" y1="16" x2="12" y2="12"/>
          <line x1="12" y1="8" x2="12.01" y2="8"/>
        </svg>
      </span>
      <div class="statement-body">
        <span class="statement-text">
          Spesifikasi per 1 unit server: <strong>{{ vcpu }} vCPU</strong> (otak pemroses) &amp; <strong>{{ ram }} GB RAM</strong> (memori kerja).
        </span>
        <span v-if="activePresetObj && !forceCustomMode" class="statement-hint">
          Cocok untuk: {{ activePresetObj.recommendedFor }}
        </span>
        <span v-else class="statement-hint">
          Mode Kustom aktif — silakan sesuaikan angka vCPU dan RAM secara bebas di bawah.
        </span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { COMPUTE_SIZE_PRESETS, type SizePresetOption } from '~/data/sizePresets'

const props = defineProps<{
  vcpu: number
  ram: number
  recommendedVcpu?: number
  recommendedRam?: number
  customExpanded?: boolean
}>()

const emit = defineEmits<{
  (e: 'update:vcpu', val: number): void
  (e: 'update:ram', val: number): void
  (e: 'update:customExpanded', val: boolean): void
}>()

const presets = COMPUTE_SIZE_PRESETS
const forceCustomMode = ref(false)

function matchesPreset(p: SizePresetOption, vcpu: number, ram: number): boolean {
  return p.vcpu === vcpu && (p.ram === ram || (p.altRam !== undefined && p.altRam === ram))
}

const currentPresetId = computed(() => {
  const match = presets.find(p => matchesPreset(p, props.vcpu, props.ram))
  return match ? match.id : 'custom'
})

watch(
  () => [props.vcpu, props.ram],
  () => {
    if (currentPresetId.value === 'custom') {
      emit('update:customExpanded', true)
    }
  },
  { immediate: true }
)

const activePresetObj = computed(() => {
  return presets.find(p => p.id === currentPresetId.value)
})

const recommendedPresetId = computed(() => {
  if (!props.recommendedVcpu || !props.recommendedRam) return null
  const match = presets.find(p => matchesPreset(p, props.recommendedVcpu!, props.recommendedRam!))
  return match ? match.id : null
})

function selectPreset(preset: SizePresetOption) {
  forceCustomMode.value = false
  // If this preset is the recommended one and uses altRam, preserve recommendedRam
  if (preset.id === recommendedPresetId.value && props.recommendedRam) {
    emit('update:vcpu', preset.vcpu)
    emit('update:ram', props.recommendedRam)
  } else {
    emit('update:vcpu', preset.vcpu)
    emit('update:ram', preset.ram)
  }
  emit('update:customExpanded', false)
}

function selectCustom() {
  forceCustomMode.value = true
  emit('update:customExpanded', true)
}
</script>

<style scoped>
.base-size-picker {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.size-pill-group {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 8px;
}

.size-pill-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 10px 6px 8px 6px;
  background: #FFFFFF;
  border: 1px solid #CBD5E1;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.15s ease;
  position: relative;
  text-align: center;
}

.size-pill-btn:hover {
  border-color: #94A3B8;
  background: #F8FAFC;
  transform: translateY(-1px);
}

.size-pill-btn.active {
  background: #EFF6FF;
  border-color: #002060;
  box-shadow: 0 0 0 2px rgba(0, 32, 96, 0.12);
}

.size-pill-btn.recommended:not(.active) {
  border-color: #FCD34D;
  background: #FFFBEB;
}

.size-top-row {
  display: flex;
  align-items: center;
  gap: 4px;
}

.size-name {
  font-size: 12px;
  font-weight: 800;
  color: #002060;
}

.active-check {
  font-size: 10px;
  font-weight: 900;
  color: #002060;
  background: #DBEAFE;
  width: 15px;
  height: 15px;
  border-radius: 50%;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.size-sub {
  font-size: 11px;
  font-weight: 600;
  color: #475569;
  margin-top: 3px;
  white-space: nowrap;
}

.rec-badge-inline {
  background: #FEF3C7;
  color: #92400E;
  border: 1px solid #FDE68A;
  font-size: 11px;
  font-weight: 800;
  padding: 1px 6px;
  border-radius: 4px;
  letter-spacing: 0.02em;
  text-transform: uppercase;
}

.active-size-statement {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  background: transparent;
  border: none;
  border-top: 1px solid #F1F5F9;
  padding: 10px 2px 2px 2px;
  font-size: 11.5px;
  color: #334155;
  margin-top: 4px;
}

.statement-icon {
  font-size: 13px;
  flex-shrink: 0;
  margin-top: 1px;
}

.statement-body {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 6px;
}

.statement-hint {
  color: #475569;
  font-size: 11.5px;
}

@media (max-width: 1100px) {
  .size-pill-group {
    grid-template-columns: repeat(4, 1fr);
  }
}

@media (max-width: 700px) {
  .size-pill-group {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>
