<template>
  <div class="base-stepper" :class="{ 'is-disabled': disabled }">
    <div v-if="label" class="stepper-label-row">
      <span class="stepper-label">{{ label }}</span>
      <span v-if="unit" class="stepper-unit-badge">{{ unit }}</span>
    </div>
    <div class="stepper-controls">
      <button
        type="button"
        class="stepper-btn btn-dec"
        :disabled="disabled || modelValue <= min"
        @click="decrement"
        aria-label="Kurangi"
      >
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round">
          <line x1="5" y1="12" x2="19" y2="12"/>
        </svg>
      </button>

      <div class="stepper-value-wrapper">
        <input
          type="number"
          class="stepper-input"
          :value="modelValue"
          :min="min"
          :max="max"
          :step="step"
          :disabled="disabled"
          @input="onInput"
        />
        <span v-if="inlineUnit" class="inline-unit">{{ inlineUnit }}</span>
      </div>

      <button
        type="button"
        class="stepper-btn btn-inc"
        :disabled="disabled || modelValue >= max"
        @click="increment"
        aria-label="Tambah"
      >
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round">
          <line x1="12" y1="5" x2="12" y2="19"/>
          <line x1="5" y1="12" x2="19" y2="12"/>
        </svg>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
const props = withDefaults(
  defineProps<{
    modelValue: number
    min?: number
    max?: number
    step?: number
    label?: string
    unit?: string
    inlineUnit?: string
    disabled?: boolean
  }>(),
  {
    min: 0,
    max: 999999,
    step: 1,
    label: '',
    unit: '',
    inlineUnit: '',
    disabled: false,
  }
)

const emit = defineEmits<{
  (e: 'update:modelValue', val: number): void
}>()

function clamp(val: number): number {
  if (isNaN(val)) return props.min
  return Math.min(Math.max(val, props.min), props.max)
}

function increment() {
  if (props.disabled) return
  emit('update:modelValue', clamp(props.modelValue + props.step))
}

function decrement() {
  if (props.disabled) return
  emit('update:modelValue', clamp(props.modelValue - props.step))
}

function onInput(event: Event) {
  const target = event.target as HTMLInputElement
  const parsed = parseInt(target.value, 10)
  if (!isNaN(parsed)) {
    emit('update:modelValue', clamp(parsed))
  }
}
</script>

<style scoped>
.base-stepper {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.stepper-label-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.stepper-label {
  font-size: 11px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.04em;
  color: #475569;
}

.stepper-unit-badge {
  font-size: 11px;
  font-weight: 500;
  color: #64748B;
  background: #F1F5F9;
  padding: 1px 6px;
  border-radius: 4px;
}

.stepper-controls {
  display: inline-flex;
  align-items: center;
  background: #FFFFFF;
  border: 1px solid #CBD5E1;
  border-radius: 6px;
  overflow: hidden;
  transition: border-color 0.15s ease, box-shadow 0.15s ease;
  height: 36px;
}

.stepper-controls:focus-within {
  border-color: #002060;
  box-shadow: 0 0 0 2px rgba(0, 32, 96, 0.12);
}

.stepper-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 34px;
  height: 100%;
  background: #F8FAFC;
  border: none;
  color: #334155;
  cursor: pointer;
  transition: background 0.12s ease, color 0.12s ease;
  user-select: none;
}

.stepper-btn:hover:not(:disabled) {
  background: #E2E8F0;
  color: #002060;
}

.stepper-btn:active:not(:disabled) {
  background: #CBD5E1;
}

.stepper-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.btn-dec {
  border-right: 1px solid #E2E8F0;
}

.btn-inc {
  border-left: 1px solid #E2E8F0;
}

.stepper-value-wrapper {
  display: flex;
  align-items: center;
  justify-content: center;
  flex: 1;
  min-width: 60px;
  padding: 0 8px;
}

.stepper-input {
  width: 100%;
  border: none;
  background: transparent;
  text-align: center;
  font-size: 13px;
  font-weight: 600;
  color: #0F172A;
  font-family: inherit;
  font-variant-numeric: tabular-nums;
  outline: none;
  -moz-appearance: textfield;
}

.stepper-input::-webkit-outer-spin-button,
.stepper-input::-webkit-inner-spin-button {
  -webkit-appearance: none;
  margin: 0;
}

.inline-unit {
  font-size: 11px;
  color: #64748B;
  margin-left: 3px;
  user-select: none;
}

.is-disabled {
  opacity: 0.6;
  pointer-events: none;
}
</style>
