<template>
  <div class="choice-cards-grid" :style="{ gridTemplateColumns: `repeat(${columns || 3}, 1fr)` }">
    <div
      v-for="opt in options"
      :key="opt.value"
      class="choice-card"
      :class="{
        active: modelValue === opt.value,
        recommended: opt.recommended && modelValue !== opt.value,
        disabled: opt.disabled
      }"
      @click="!opt.disabled && $emit('update:modelValue', opt.value)"
    >
      <div class="card-top">
        <div class="card-top-left">
          <span v-if="opt.icon" class="opt-icon">{{ opt.icon }}</span>
          <div class="opt-title-wrap">
            <span class="opt-title">{{ opt.title }}</span>
            <span v-if="opt.technicalLabel" class="opt-tech-pill">{{ opt.technicalLabel }}</span>
          </div>
        </div>

        <!-- Clear Selection Radio / Checkmark Indicator for Non-Tech Users -->
        <span class="selection-indicator" :class="{ checked: modelValue === opt.value }">
          <svg v-if="modelValue === opt.value" width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3.5" stroke-linecap="round" stroke-linejoin="round">
            <polyline points="20 6 9 17 4 12"></polyline>
          </svg>
        </span>
      </div>

      <p class="opt-desc">{{ opt.description }}</p>

      <div v-if="opt.badge || opt.recommended || modelValue === opt.value" class="opt-badge-strip">
        <span v-if="modelValue === opt.value" class="opt-badge selected-pill">Terpilih</span>
        <span v-if="opt.recommended" class="opt-badge std-pill">★ Standar EA</span>
        <span v-if="opt.badge" :class="['opt-badge', opt.badgeType || 'info']">{{ opt.badge }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
export interface ChoiceOption {
  value: string
  title: string
  technicalLabel?: string
  description: string
  icon?: string
  badge?: string
  badgeType?: 'info' | 'success' | 'warning'
  recommended?: boolean
  disabled?: boolean
}

defineProps<{
  modelValue: string
  options: ChoiceOption[]
  columns?: number
}>()

defineEmits<{
  (e: 'update:modelValue', val: string): void
}>()
</script>

<style scoped>
.choice-cards-grid {
  display: grid;
  gap: 10px;
}

.choice-card {
  background: #FFFFFF;
  border: 1.5px solid #CBD5E1;
  border-radius: 9px;
  padding: 13px 14px;
  cursor: pointer;
  transition: all 0.16s ease;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  gap: 8px;
  position: relative;
  text-align: left;
}

.choice-card:hover:not(.disabled) {
  border-color: #64748B;
  background: #F8FAFC;
  transform: translateY(-1px);
}

.choice-card.recommended {
  border-color: #FCD34D;
  background: #FFFDF7;
}

.choice-card.active {
  background: #F0F6FF;
  border-color: #002060;
  box-shadow: 0 0 0 2px rgba(0, 32, 96, 0.12);
}

.choice-card.disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.card-top {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 8px;
}

.card-top-left {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
}

.opt-icon {
  font-size: 17px;
  flex-shrink: 0;
}

.opt-title-wrap {
  display: flex;
  flex-direction: column;
  min-width: 0;
}

.opt-title {
  font-size: 12.5px;
  font-weight: 800;
  color: #0F172A;
  line-height: 1.25;
}

.opt-tech-pill {
  font-size: 11px;
  font-family: 'JetBrains Mono', monospace;
  color: #475569;
  margin-top: 1px;
}

.selection-indicator {
  width: 18px;
  height: 18px;
  border-radius: 50%;
  border: 1.5px solid #CBD5E1;
  background: #FFFFFF;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  color: #FFFFFF;
  transition: all 0.15s ease;
}

.selection-indicator.checked {
  background: #002060;
  border-color: #002060;
}

.opt-desc {
  margin: 0;
  font-size: 11.5px;
  color: #475569;
  line-height: 1.45;
  flex: 1;
}

.opt-badge-strip {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 5px;
  margin-top: 2px;
}

.opt-badge {
  font-size: 11px;
  font-weight: 700;
  padding: 2px 7px;
  border-radius: 4px;
  text-transform: uppercase;
  letter-spacing: 0.02em;
}

.opt-badge.selected-pill {
  background: #002060;
  color: #FFFFFF;
}

.opt-badge.std-pill {
  background: #FEF3C7;
  color: #B45309;
  border: 1px solid #FDE68A;
}

.opt-badge.info {
  background: #EFF6FF;
  color: #1D4ED8;
}

.opt-badge.success {
  background: #DCFCE7;
  color: #15803D;
}

.opt-badge.warning {
  background: #FEF3C7;
  color: #92400E;
}

@media (max-width: 900px) {
  .choice-cards-grid {
    grid-template-columns: 1fr !important;
  }
}
</style>
