<template>
  <div class="base-size-picker">
    <div class="size-pill-group">
      <button
        v-for="p in presets"
        :key="p.id"
        type="button"
        class="size-pill-btn"
        :class="{ active: currentPresetId === p.id, recommended: p.id === recommendedPresetId }"
        @click="selectPreset(p)"
      >
        <span class="size-name">{{ p.name }}</span>
        <span class="size-sub">{{ p.sublabel }}</span>
        <span v-if="p.id === recommendedPresetId" class="rec-badge">Standar</span>
      </button>

      <button
        type="button"
        class="size-pill-btn btn-custom"
        :class="{ active: currentPresetId === 'custom' }"
        @click="selectCustom"
      >
        <span class="size-name">Kustom</span>
        <span class="size-sub">Atur Manual</span>
      </button>
    </div>

    <!-- Active size summary statement -->
    <div class="active-size-statement">
      <span class="statement-icon">💡</span>
      <span class="statement-text">
        Alokasi per server: <strong>{{ vcpu }} vCPU</strong> pemroses &amp; <strong>{{ ram }} GB</strong> memori RAM.
        <span v-if="activePresetObj" class="statement-hint">({{ activePresetObj.recommendedFor }})</span>
      </span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { COMPUTE_SIZE_PRESETS, type SizePresetOption } from '~/data/sizePresets'

const props = defineProps<{
  vcpu: number
  ram: number
  recommendedVcpu?: number
  recommendedRam?: number
}>()

const emit = defineEmits<{
  (e: 'update:vcpu', val: number): void
  (e: 'update:ram', val: number): void
}>()

const presets = COMPUTE_SIZE_PRESETS

const currentPresetId = computed(() => {
  const match = presets.find(p => p.vcpu === props.vcpu && p.ram === props.ram)
  return match ? match.id : 'custom'
})

const activePresetObj = computed(() => {
  return presets.find(p => p.id === currentPresetId.value)
})

const recommendedPresetId = computed(() => {
  if (!props.recommendedVcpu || !props.recommendedRam) return null
  const match = presets.find(p => p.vcpu === props.recommendedVcpu && p.ram === props.recommendedRam)
  return match ? match.id : null
})

function selectPreset(preset: SizePresetOption) {
  emit('update:vcpu', preset.vcpu)
  emit('update:ram', preset.ram)
}

function selectCustom() {
  // Keeps current numbers but switches to custom mode
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
  grid-template-columns: repeat(6, 1fr);
  gap: 8px;
}

.size-pill-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 8px 6px;
  background: #FFFFFF;
  border: 1px solid #CBD5E1;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.15s ease;
  position: relative;
  text-align: center;
}

.size-pill-btn:hover {
  border-color: #94A3B8;
  background: #F8FAFC;
}

.size-pill-btn.active {
  background: #EFF6FF;
  border-color: #002060;
  box-shadow: 0 0 0 2px rgba(0, 32, 96, 0.12);
}

.size-pill-btn.recommended {
  border-color: #FCD34D;
}

.size-name {
  font-size: 12px;
  font-weight: 800;
  color: #002060;
}

.size-sub {
  font-size: 9px;
  color: #64748B;
  margin-top: 2px;
  white-space: nowrap;
}

.rec-badge {
  position: absolute;
  top: -6px;
  right: -4px;
  background: #D4A017;
  color: #FFFFFF;
  font-size: 8px;
  font-weight: 800;
  padding: 1px 4px;
  border-radius: 3px;
  letter-spacing: 0.03em;
  text-transform: uppercase;
}

.active-size-statement {
  display: flex;
  align-items: center;
  gap: 8px;
  background: #F8FAFC;
  border: 1px solid #E2E8F0;
  padding: 7px 12px;
  border-radius: 5px;
  font-size: 11px;
  color: #334155;
}

.statement-icon {
  font-size: 13px;
}

.statement-hint {
  color: #64748B;
  margin-left: 4px;
}

@media (max-width: 900px) {
  .size-pill-group {
    grid-template-columns: repeat(3, 1fr);
  }
}
</style>
