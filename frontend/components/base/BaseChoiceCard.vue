<template>
  <div class="choice-cards-grid" :style="{ gridTemplateColumns: `repeat(${columns || 3}, 1fr)` }">
    <div
      v-for="opt in options"
      :key="opt.value"
      class="choice-card"
      :class="{ active: modelValue === opt.value, disabled: opt.disabled }"
      @click="!opt.disabled && $emit('update:modelValue', opt.value)"
    >
      <div class="card-top">
        <span v-if="opt.icon" class="opt-icon">{{ opt.icon }}</span>
        <div class="opt-title-wrap">
          <span class="opt-title">{{ opt.title }}</span>
          <span v-if="opt.technicalLabel" class="opt-tech-pill">{{ opt.technicalLabel }}</span>
        </div>
      </div>
      <p class="opt-desc">{{ opt.description }}</p>
      <div v-if="opt.badge" class="opt-badge-strip">
        <span :class="['opt-badge', opt.badgeType || 'info']">{{ opt.badge }}</span>
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
  border: 1px solid #CBD5E1;
  border-radius: 8px;
  padding: 12px 14px;
  cursor: pointer;
  transition: all 0.15s ease;
  display: flex;
  flex-direction: column;
  gap: 6px;
  position: relative;
  text-align: left;
}

.choice-card:hover:not(.disabled) {
  border-color: #94A3B8;
  background: #F8FAFC;
}

.choice-card.active {
  background: #EFF6FF;
  border-color: #002060;
  box-shadow: 0 0 0 2px rgba(0, 32, 96, 0.12);
}

.choice-card.disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.card-top {
  display: flex;
  align-items: center;
  gap: 8px;
}

.opt-icon {
  font-size: 16px;
}

.opt-title-wrap {
  display: flex;
  flex-direction: column;
}

.opt-title {
  font-size: 12px;
  font-weight: 800;
  color: #0F172A;
}

.opt-tech-pill {
  font-size: 9px;
  font-family: 'JetBrains Mono', monospace;
  color: #64748B;
}

.opt-desc {
  margin: 0;
  font-size: 11px;
  color: #475569;
  line-height: 1.4;
}

.opt-badge-strip {
  margin-top: 2px;
}

.opt-badge {
  font-size: 9px;
  font-weight: 700;
  padding: 2px 6px;
  border-radius: 3px;
  text-transform: uppercase;
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
