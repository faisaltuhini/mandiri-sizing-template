<template>
  <div class="base-advanced-toggle">
    <button
      type="button"
      class="toggle-link-btn"
      @click="toggle"
    >
      <svg
        class="toggle-chevron"
        :class="{ rotated: isOpen }"
        width="13"
        height="13"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="2.5"
        stroke-linecap="round"
        stroke-linejoin="round"
      >
        <polyline points="9 18 15 12 9 6"></polyline>
      </svg>
      <span class="btn-text">{{ displayLabel }}</span>
    </button>

    <div v-show="isOpen" class="advanced-slot-wrapper">
      <slot />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'

const props = withDefaults(
  defineProps<{
    showLabel?: string
    hideLabel?: string
    defaultOpen?: boolean
    modelValue?: boolean
  }>(),
  {
    defaultOpen: false,
    modelValue: undefined
  }
)

const emit = defineEmits<{
  (e: 'update:modelValue', val: boolean): void
}>()

const internalOpen = ref(props.modelValue !== undefined ? props.modelValue : props.defaultOpen)

watch(
  () => props.modelValue,
  (val) => {
    if (val !== undefined) {
      internalOpen.value = val
    }
  }
)

const isOpen = computed(() => {
  return props.modelValue !== undefined ? props.modelValue : internalOpen.value
})

function toggle() {
  const next = !isOpen.value
  internalOpen.value = next
  emit('update:modelValue', next)
}

const displayLabel = computed(() => {
  const raw = isOpen.value
    ? (props.hideLabel || 'Sembunyikan Pengaturan Detail')
    : (props.showLabel || 'Atur Angka Spesifik Manual (Lanjutan)')
  return raw.replace(/^[▸▾▶▼>v\s]+/, '').trim()
})
</script>

<style scoped>
.base-advanced-toggle {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.toggle-link-btn {
  background: none;
  border: none;
  padding: 4px 0;
  display: inline-flex;
  align-items: center;
  gap: 6px;
  color: #475569;
  font-size: 11.5px;
  font-weight: 700;
  cursor: pointer;
  align-self: flex-start;
  transition: all 0.15s ease;
  user-select: none;
}

.toggle-link-btn:hover {
  color: #002060;
}

.toggle-chevron {
  transition: transform 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  color: inherit;
  flex-shrink: 0;
}

.toggle-chevron.rotated {
  transform: rotate(90deg);
}

.btn-text {
  line-height: 1.2;
}

.advanced-slot-wrapper {
  background: #F8FAFC;
  border: 1px dashed #CBD5E1;
  border-radius: 8px;
  padding: 12px 14px;
}
</style>
