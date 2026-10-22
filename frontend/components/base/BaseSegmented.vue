<template>
  <div class="base-segmented">
    <button
      v-for="opt in options"
      :key="opt.value"
      type="button"
      :class="['segmented-item', { active: modelValue === opt.value }, opt.accent || '']"
      @click="$emit('update:modelValue', opt.value)"
    >
      {{ opt.label }}
    </button>
  </div>
</template>

<script setup lang="ts">
export interface SegmentedOption {
  value: string
  label: string
  accent?: string
}

defineProps<{
  modelValue: string
  options: SegmentedOption[]
}>()

defineEmits<{
  (e: 'update:modelValue', val: string): void
}>()
</script>

<style scoped>
.base-segmented {
  display: inline-flex;
  align-items: center;
  background: #F1F5F9;
  padding: 2px;
  border-radius: 6px;
  border: 1px solid #E2E8F0;
  gap: 2px;
}

.segmented-item {
  padding: 4px 11px;
  font-size: 11px;
  font-weight: 500;
  color: #64748B;
  background: transparent;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.15s ease;
  user-select: none;
  white-space: nowrap;
}

.segmented-item:hover:not(.active) {
  color: #0F172A;
  background: rgba(255, 255, 255, 0.6);
}

.segmented-item.active {
  background: #FFFFFF;
  color: #002060;
  font-weight: 600;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.06);
}

.segmented-item.crit.active {
  color: #B45309;
}

.segmented-item.mod.active {
  color: #0369A1;
}
</style>
