<template>
  <div class="number-input-wrapper" :class="{ disabled, focused: isFocused }">
    <div class="input-container">
      <span v-if="prefix" class="input-affix prefix">{{ prefix }}</span>
      <input
        ref="inputRef"
        type="text"
        inputmode="decimal"
        :value="displayValue"
        :placeholder="placeholder"
        :disabled="disabled"
        class="formatted-number-input"
        @focus="onFocus"
        @blur="onBlur"
        @input="onInput"
        @keydown="onKeyDown"
      />
      <span v-if="suffix" class="input-affix suffix">{{ suffix }}</span>
      
      <div v-if="showStepper && !disabled" class="stepper-controls">
        <button
          type="button"
          tabindex="-1"
          class="stepper-btn stepper-up"
          title="Tambah"
          @click.prevent="stepBy(1)"
        >
          ▲
        </button>
        <button
          type="button"
          tabindex="-1"
          class="stepper-btn stepper-down"
          title="Kurang"
          @click.prevent="stepBy(-1)"
        >
          ▼
        </button>
      </div>
    </div>
    <span v-if="unit" class="input-unit-label">{{ unit }}</span>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, nextTick } from 'vue'
import { formatNumber, parseFormattedNumber } from '~/composables/useFormat'

interface Props {
  modelValue?: number | null
  min?: number
  max?: number
  step?: number
  placeholder?: string
  unit?: string
  prefix?: string
  suffix?: string
  disabled?: boolean
  maxDecimals?: number
  minDecimals?: number
  showStepper?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: 0,
  step: 1,
  placeholder: '0',
  unit: '',
  prefix: '',
  suffix: '',
  disabled: false,
  maxDecimals: 4,
  minDecimals: 0,
  showStepper: true,
})

const emit = defineEmits<{
  (e: 'update:modelValue', value: number): void
  (e: 'change', value: number): void
}>()

const inputRef = ref<HTMLInputElement | null>(null)
const isFocused = ref(false)
const displayValue = ref('')

function formatCurrent(val: number | null | undefined): string {
  if (val === null || val === undefined || isNaN(val)) return ''
  return formatNumber(val, {
    minDecimals: props.minDecimals,
    maxDecimals: props.maxDecimals,
    locale: 'id-ID',
  })
}

// Format integer part with dots for Indonesian thousands
function formatIdInteger(val: number | string): string {
  const num = typeof val === 'string' ? parseFloat(val) : val
  if (isNaN(num)) return ''
  return new Intl.NumberFormat('id-ID', { useGrouping: true, maximumFractionDigits: 0 }).format(num)
}

// Initialize and sync with external prop changes
watch(
  () => props.modelValue,
  (newVal) => {
    if (!isFocused.value) {
      displayValue.value = formatCurrent(newVal)
    }
  },
  { immediate: true }
)

function onFocus(e: FocusEvent) {
  isFocused.value = true
  const target = e.target as HTMLInputElement
  if (target) {
    displayValue.value = formatCurrent(props.modelValue)
    nextTick(() => {
      target.select()
    })
  }
}

function onBlur() {
  isFocused.value = false
  let num = parseFormattedNumber(displayValue.value)

  if (props.min !== undefined && num < props.min) {
    num = props.min
  }
  if (props.max !== undefined && num > props.max) {
    num = props.max
  }

  emit('update:modelValue', num)
  emit('change', num)
  displayValue.value = formatCurrent(num)
}

function onInput(e: Event) {
  const target = e.target as HTMLInputElement
  const rawInput = target.value

  // Handle empty input
  if (!rawInput || rawInput.trim() === '') {
    displayValue.value = ''
    emit('update:modelValue', 0)
    return
  }

  // In Indonesian typing: dots are thousands, commas (or dots) can be typed as decimal
  // Clean dots (thousands)
  let rawWithoutDots = rawInput.replace(/\./g, '')

  // Allow either ',' or '.' as decimal separator
  let decimalChar = ','
  if (rawWithoutDots.includes(',')) {
    decimalChar = ','
  } else if (rawWithoutDots.includes('.')) {
    decimalChar = '.'
  }

  // Split integer and decimal parts
  const parts = rawWithoutDots.split(decimalChar === ',' ? ',' : '.')
  let intPartClean = parts[0].replace(/[^0-9-]/g, '')
  let decPartClean = parts.length > 1 ? parts.slice(1).join('').replace(/[^0-9]/g, '') : null

  let numericVal = 0
  if (decPartClean !== null) {
    numericVal = parseFloat(`${intPartClean || '0'}.${decPartClean}`)
  } else {
    numericVal = parseFloat(intPartClean || '0')
  }

  if (!isNaN(numericVal)) {
    emit('update:modelValue', numericVal)
  }

  // Format display value with dots for thousands and comma for decimal
  const selectionStart = target.selectionStart || 0
  const prevLength = rawInput.length

  const formattedInt = intPartClean ? formatIdInteger(intPartClean) : ''

  if (decPartClean !== null) {
    // If decimal typed
    const truncatedDec = decPartClean.slice(0, props.maxDecimals)
    displayValue.value = `${formattedInt},${truncatedDec}`
  } else if (rawInput.endsWith(',') || rawInput.endsWith('.')) {
    // User just typed decimal comma/dot
    displayValue.value = `${formattedInt},`
  } else {
    displayValue.value = formattedInt
  }

  // Adjust cursor position
  nextTick(() => {
    if (inputRef.value) {
      const newLength = displayValue.value.length
      const diff = newLength - prevLength
      const newCursor = Math.max(0, Math.min(selectionStart + diff, newLength))
      inputRef.value.setSelectionRange(newCursor, newCursor)
    }
  })
}

function stepBy(direction: number) {
  const current = props.modelValue ?? 0
  const stepAmount = props.step || 1
  let updated = current + direction * stepAmount

  // Round to prevent floating point inaccuracies
  const precision = props.maxDecimals ?? 2
  updated = Math.round(updated * Math.pow(10, precision)) / Math.pow(10, precision)

  if (props.min !== undefined && updated < props.min) {
    updated = props.min
  }
  if (props.max !== undefined && updated > props.max) {
    updated = props.max
  }

  emit('update:modelValue', updated)
  emit('change', updated)
  displayValue.value = formatCurrent(updated)
}

function onKeyDown(e: KeyboardEvent) {
  if (e.key === 'ArrowUp') {
    e.preventDefault()
    stepBy(1)
  } else if (e.key === 'ArrowDown') {
    e.preventDefault()
    stepBy(-1)
  }
}
</script>

<style scoped>
.number-input-wrapper {
  display: flex;
  flex-direction: column;
  gap: 4px;
  width: 100%;
}

.input-container {
  display: flex;
  align-items: center;
  position: relative;
  background: var(--bg2);
  border: 1px solid var(--border2);
  border-radius: var(--radius-sm);
  transition: all 0.15s ease-in-out;
  overflow: hidden;
}

.input-container:hover {
  border-color: var(--border);
}

.number-input-wrapper.focused .input-container {
  border-color: var(--blue, #2563eb);
  box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.12);
  background: #ffffff;
}

.formatted-number-input {
  font-family: 'Inter', sans-serif;
  font-size: 13px;
  font-weight: 500;
  color: var(--text);
  background: transparent;
  border: none;
  outline: none;
  padding: 8px 10px;
  width: 100%;
  box-sizing: border-box;
}

.input-affix {
  font-size: 12px;
  color: var(--text3);
  font-weight: 600;
  user-select: none;
}

.prefix {
  padding-left: 10px;
}

.suffix {
  padding-right: 10px;
}

.stepper-controls {
  display: flex;
  flex-direction: column;
  height: 100%;
  border-left: 1px solid var(--border2);
  background: var(--bg3);
  opacity: 0.7;
  transition: opacity 0.15s;
}

.input-container:hover .stepper-controls,
.number-input-wrapper.focused .stepper-controls {
  opacity: 1;
}

.stepper-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  background: transparent;
  border: none;
  font-size: 7px;
  color: var(--text2);
  cursor: pointer;
  padding: 2px 7px;
  line-height: 1;
  transition: background 0.1s, color 0.1s;
}

.stepper-up {
  border-bottom: 1px solid var(--border2);
  padding-top: 4px;
}

.stepper-down {
  padding-bottom: 4px;
}

.stepper-btn:hover {
  background: var(--border);
  color: var(--navy);
}

.stepper-btn:active {
  background: var(--brand-blue, #003b73);
  color: #ffffff;
}

.input-unit-label {
  font-size: 10px;
  color: var(--text3);
  line-height: 1.2;
}

.number-input-wrapper.disabled {
  opacity: 0.6;
  pointer-events: none;
}
</style>
