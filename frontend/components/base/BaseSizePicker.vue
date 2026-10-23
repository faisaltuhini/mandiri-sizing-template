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
          <span v-if="!forceCustomMode && currentPresetId === p.id" class="active-check">✓</span>
        </div>
        <span class="size-sub">{{ p.sublabel }}</span>
        <span v-if="p.id === recommendedPresetId" class="rec-badge" title="Rekomendasi Standar Arsitektur">★ Standar EA</span>
      </button>

      <button
        type="button"
        class="size-pill-btn btn-custom"
        :class="{ active: forceCustomMode || currentPresetId === 'custom' }"
        @click="selectCustom"
      >
        <div class="size-top-row">
          <span class="size-name">Kustom</span>
          <span v-if="forceCustomMode || currentPresetId === 'custom'" class="active-check">✓</span>
        </div>
        <span class="size-sub">Atur Angka Manual</span>
      </button>
    </div>

    <!-- Active size summary statement -->
    <div class="active-size-statement">
      <span class="statement-icon">💡</span>
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
  font-size: 9.5px;
  font-weight: 600;
  color: #64748B;
  margin-top: 3px;
  white-space: nowrap;
}

.rec-badge {
  position: absolute;
  top: -8px;
  right: 6px;
  background: #D97706;
  color: #FFFFFF;
  font-size: 8px;
  font-weight: 800;
  padding: 1.5px 5px;
  border-radius: 4px;
  letter-spacing: 0.02em;
  text-transform: uppercase;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
}

.active-size-statement {
  display: flex;
  align-items: center;
  gap: 10px;
  background: #F8FAFC;
  border: 1px solid #E2E8F0;
  padding: 8px 12px;
  border-radius: 6px;
  font-size: 11.5px;
  color: #334155;
}

.statement-icon {
  font-size: 14px;
  flex-shrink: 0;
}

.statement-body {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 6px;
}

.statement-hint {
  color: #64748B;
  font-size: 11px;
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
