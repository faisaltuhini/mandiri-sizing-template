<template>
  <div class="base-autocomplete" ref="containerRef">
    <div class="input-wrapper" :class="{ 'is-focused': isOpen, 'is-disabled': disabled }">
      <input
        ref="inputRef"
        type="text"
        :value="modelValue"
        :placeholder="placeholder"
        :required="required"
        :disabled="disabled"
        class="clean-autocomplete-input"
        @input="handleInput"
        @focus="handleFocus"
        @keydown="handleKeydown"
      />

      <div class="suffix-actions">
        <!-- Clear (X) Button -->
        <button
          v-if="modelValue && !disabled"
          type="button"
          class="btn-clear-suffix"
          title="Hapus / Bersihkan"
          tabindex="-1"
          @mousedown.prevent="clearInput"
        >
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round">
            <line x1="18" y1="6" x2="6" y2="18"></line>
            <line x1="6" y1="6" x2="18" y2="18"></line>
          </svg>
        </button>

        <!-- Dropdown toggle arrow -->
        <button
          type="button"
          class="btn-toggle-suffix"
          :class="{ 'rotate-up': isOpen }"
          tabindex="-1"
          @mousedown.prevent="toggleDropdown"
          title="Tampilkan pilihan katalog"
        >
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
            <polyline points="6 9 12 15 18 9"></polyline>
          </svg>
        </button>
      </div>
    </div>

    <!-- Dropdown Popover Menu via Teleport to Body (Overlay di atas Modal) -->
    <Teleport to="body" v-if="isMounted">
      <transition name="dropdown-fade">
        <div
          v-if="isOpen && filteredOptions.length > 0"
          ref="dropdownRef"
          class="autocomplete-dropdown fixed-overlay-dropdown"
          :style="dropdownStyle"
          @mousedown.stop
        >
          <div class="dropdown-header">
            <span class="hdr-label">Pilihan dari Katalog ({{ filteredOptions.length }})</span>
            <span class="hdr-shortcut">Gunakan ↑ ↓ Enter</span>
          </div>

          <ul class="options-list" ref="listRef">
            <li
              v-for="(item, idx) in filteredOptions"
              :key="idx"
              class="option-item"
              :class="{ 'is-highlighted': highlightedIndex === idx, 'is-selected': item.value === modelValue }"
              @mouseenter="highlightedIndex = idx"
              @mousedown.prevent="selectOption(item)"
            >
              <div class="option-icon" v-if="item.icon">{{ item.icon }}</div>
              <div class="option-body">
                <div class="option-title-row">
                  <span class="option-label" v-html="highlightMatch(item.label)"></span>
                  <span v-if="item.badge" class="option-badge" :class="item.badgeType || 'default'">
                    {{ item.badge }}
                  </span>
                </div>
                <p v-if="item.description" class="option-desc">
                  {{ item.description }}
                </p>
              </div>
              <span v-if="item.category" class="option-category-tag">
                {{ item.category }}
              </span>
            </li>
          </ul>
        </div>

        <!-- Empty Suggestion State when typing something not in catalog -->
        <div
          v-else-if="isOpen && query.trim() && filteredOptions.length === 0"
          ref="dropdownRef"
          class="autocomplete-dropdown fixed-overlay-dropdown empty-dropdown"
          :style="dropdownStyle"
          @mousedown.stop
        >
          <div class="empty-state-content">
            <span class="empty-icon">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M12 20h9"></path>
                <path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z"></path>
              </svg>
            </span>
            <div class="empty-texts">
              <span class="empty-title">Input Kustom Baru</span>
              <span class="empty-sub">"{{ query }}" tidak ada di katalog standar, tetap akan disimpan sebagai nilai kustom.</span>
            </div>
          </div>
        </div>
      </transition>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted, nextTick } from 'vue'

export interface AutocompleteOption {
  label: string
  value: string
  icon?: string
  category?: string
  badge?: string
  badgeType?: 'default' | 'primary' | 'warning' | 'success'
  description?: string
  meta?: any
}

const props = withDefaults(
  defineProps<{
    modelValue: string
    options: (string | AutocompleteOption)[]
    placeholder?: string
    required?: boolean
    disabled?: boolean
    maxItems?: number
  }>(),
  {
    placeholder: '',
    required: false,
    disabled: false,
    maxItems: 8
  }
)

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void
  (e: 'select', option: AutocompleteOption): void
  (e: 'clear'): void
}>()

const isMounted = ref(false)
const containerRef = ref<HTMLElement | null>(null)
const inputRef = ref<HTMLInputElement | null>(null)
const dropdownRef = ref<HTMLElement | null>(null)
const listRef = ref<HTMLElement | null>(null)

const isOpen = ref(false)
const query = ref('')
const highlightedIndex = ref(0)
const dropdownStyle = ref<Record<string, string>>({})

// Standardize options into unified structure
const normalizedOptions = computed<AutocompleteOption[]>(() => {
  return props.options.map(opt => {
    if (typeof opt === 'string') {
      return {
        label: opt,
        value: opt
      }
    }
    return opt
  })
})

// Filter options according to search query
const filteredOptions = computed(() => {
  const q = query.value.toLowerCase().trim()
  let list = normalizedOptions.value

  if (q) {
    list = list.filter(item => {
      const matchLabel = item.label.toLowerCase().includes(q)
      const matchVal = item.value.toLowerCase().includes(q)
      const matchCat = item.category ? item.category.toLowerCase().includes(q) : false
      const matchDesc = item.description ? item.description.toLowerCase().includes(q) : false
      return matchLabel || matchVal || matchCat || matchDesc
    })
  }

  return list.slice(0, props.maxItems)
})

function updateDropdownPosition() {
  if (!containerRef.value) return
  const rect = containerRef.value.getBoundingClientRect()
  const viewportHeight = window.innerHeight
  const dropdownHeight = 260
  const spaceBelow = viewportHeight - rect.bottom
  const spaceAbove = rect.top

  // Smart flip: show above if not enough space below
  const showAbove = spaceBelow < dropdownHeight && spaceAbove > spaceBelow

  if (showAbove) {
    dropdownStyle.value = {
      position: 'fixed',
      left: `${rect.left}px`,
      bottom: `${viewportHeight - rect.top + 4}px`,
      width: `${rect.width}px`,
      zIndex: '99999',
      maxHeight: `${Math.min(dropdownHeight, Math.max(spaceAbove - 16, 120))}px`
    }
  } else {
    dropdownStyle.value = {
      position: 'fixed',
      left: `${rect.left}px`,
      top: `${rect.bottom + 4}px`,
      width: `${rect.width}px`,
      zIndex: '99999',
      maxHeight: `${Math.min(dropdownHeight, Math.max(spaceBelow - 16, 120))}px`
    }
  }
}

watch(
  () => props.modelValue,
  newVal => {
    query.value = newVal || ''
  },
  { immediate: true }
)

watch(isOpen, val => {
  if (val) {
    nextTick(() => {
      updateDropdownPosition()
    })
  }
})

function handleInput(event: Event) {
  const target = event.target as HTMLInputElement
  query.value = target.value
  emit('update:modelValue', target.value)
  isOpen.value = true
  highlightedIndex.value = 0
  nextTick(() => updateDropdownPosition())
}

function handleFocus() {
  if (props.disabled) return
  query.value = props.modelValue || ''
  isOpen.value = true
  highlightedIndex.value = 0
  nextTick(() => updateDropdownPosition())
}

function toggleDropdown() {
  if (props.disabled) return
  if (isOpen.value) {
    isOpen.value = false
  } else {
    inputRef.value?.focus()
    isOpen.value = true
    nextTick(() => updateDropdownPosition())
  }
}

function clearInput() {
  query.value = ''
  emit('update:modelValue', '')
  emit('clear')
  nextTick(() => {
    inputRef.value?.focus()
    updateDropdownPosition()
  })
}

function selectOption(option: AutocompleteOption) {
  emit('update:modelValue', option.value)
  emit('select', option)
  query.value = option.value
  isOpen.value = false
}

function handleKeydown(e: KeyboardEvent) {
  if (!isOpen.value) {
    if (e.key === 'ArrowDown' || e.key === 'ArrowUp') {
      isOpen.value = true
      nextTick(() => updateDropdownPosition())
      e.preventDefault()
    }
    return
  }

  if (e.key === 'ArrowDown') {
    e.preventDefault()
    if (filteredOptions.value.length > 0) {
      highlightedIndex.value = (highlightedIndex.value + 1) % filteredOptions.value.length
      scrollToHighlighted()
    }
  } else if (e.key === 'ArrowUp') {
    e.preventDefault()
    if (filteredOptions.value.length > 0) {
      highlightedIndex.value =
        (highlightedIndex.value - 1 + filteredOptions.value.length) % filteredOptions.value.length
      scrollToHighlighted()
    }
  } else if (e.key === 'Enter') {
    if (filteredOptions.value.length > 0 && highlightedIndex.value >= 0) {
      e.preventDefault()
      selectOption(filteredOptions.value[highlightedIndex.value])
    }
  } else if (e.key === 'Escape') {
    isOpen.value = false
  }
}

function scrollToHighlighted() {
  nextTick(() => {
    if (!listRef.value) return
    const el = listRef.value.children[highlightedIndex.value] as HTMLElement
    if (el) {
      el.scrollIntoView({ block: 'nearest' })
    }
  })
}

function highlightMatch(text: string): string {
  if (!query.value.trim()) return escapeHtml(text)
  const escapedQuery = query.value.replace(/[-[\]{}()*+?.,\\^$|#\s]/g, '\\$&')
  const regex = new RegExp(`(${escapedQuery})`, 'gi')
  return escapeHtml(text).replace(regex, '<mark class="hl-match">$1</mark>')
}

function escapeHtml(text: string): string {
  return text
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#039;')
}

function handleClickOutside(event: MouseEvent) {
  const target = event.target as Node
  const clickedInContainer = containerRef.value && containerRef.value.contains(target)
  const clickedInDropdown = dropdownRef.value && dropdownRef.value.contains(target)

  if (!clickedInContainer && !clickedInDropdown) {
    isOpen.value = false
  }
}

function onScrollOrResize() {
  if (isOpen.value) {
    updateDropdownPosition()
  }
}

onMounted(() => {
  isMounted.value = true
  document.addEventListener('click', handleClickOutside)
  window.addEventListener('resize', onScrollOrResize)
  window.addEventListener('scroll', onScrollOrResize, true)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
  window.removeEventListener('resize', onScrollOrResize)
  window.removeEventListener('scroll', onScrollOrResize, true)
})
</script>

<style scoped>
.base-autocomplete {
  position: relative;
  width: 100%;
}

.input-wrapper {
  position: relative;
  display: flex;
  align-items: center;
  width: 100%;
  background: #FFFFFF;
  border: 1px solid #CBD5E1;
  border-radius: 6px;
  transition: all 0.15s ease;
}

.input-wrapper:hover {
  border-color: #94A3B8;
}

.input-wrapper.is-focused {
  border-color: #002060;
  box-shadow: 0 0 0 3px rgba(0, 32, 96, 0.08);
}

.input-wrapper.is-disabled {
  background: #F1F5F9;
  cursor: not-allowed;
  opacity: 0.7;
}

.clean-autocomplete-input {
  flex: 1;
  width: 100%;
  padding: 9px 56px 9px 12px;
  border: none;
  background: transparent;
  font-size: 13px;
  color: #0F172A;
  outline: none;
  box-sizing: border-box;
}

.clean-autocomplete-input::placeholder {
  color: #94A3B8;
}

.suffix-actions {
  position: absolute;
  right: 6px;
  top: 50%;
  transform: translateY(-50%);
  display: flex;
  align-items: center;
  gap: 2px;
}

.btn-clear-suffix {
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #E2E8F0;
  border: none;
  border-radius: 50%;
  color: #475569;
  cursor: pointer;
  transition: all 0.12s ease;
  padding: 0;
}

.btn-clear-suffix:hover {
  background: #CBD5E1;
  color: #0F172A;
}

.btn-toggle-suffix {
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: transparent;
  border: none;
  color: #94A3B8;
  cursor: pointer;
  transition: transform 0.18s ease;
  padding: 0;
}

.btn-toggle-suffix.rotate-up {
  transform: rotate(180deg);
}

.btn-toggle-suffix:hover {
  color: #002060;
}
</style>

<style>
/* ===================================================
   TELEPORTED DROPDOWN STYLING (GLOBAL OVERLAY ON TOP)
   =================================================== */
.autocomplete-dropdown.fixed-overlay-dropdown {
  background: #FFFFFF;
  border: 1px solid #CBD5E1;
  border-radius: 8px;
  box-shadow: 0 20px 40px -8px rgba(0, 32, 96, 0.25), 0 10px 15px -3px rgba(0, 0, 0, 0.1);
  z-index: 99999 !important;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.autocomplete-dropdown .dropdown-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 12px;
  background: #F8FAFC;
  border-bottom: 1px solid #E2E8F0;
  font-size: 11px;
}

.autocomplete-dropdown .hdr-label {
  font-weight: 700;
  color: #002060;
  text-transform: uppercase;
  letter-spacing: 0.03em;
}

.autocomplete-dropdown .hdr-shortcut {
  color: #64748B;
  font-size: 10px;
}

.autocomplete-dropdown .options-list {
  list-style: none;
  margin: 0;
  padding: 4px;
  overflow-y: auto;
  max-height: 240px;
}

.autocomplete-dropdown .option-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 10px;
  border-radius: 6px;
  cursor: pointer;
  transition: background 0.12s ease;
}

.autocomplete-dropdown .option-item:hover,
.autocomplete-dropdown .option-item.is-highlighted {
  background: #EFF6FF;
}

.autocomplete-dropdown .option-item.is-selected {
  background: #E0E7FF;
}

.autocomplete-dropdown .option-icon {
  font-size: 16px;
  line-height: 1;
  flex-shrink: 0;
}

.autocomplete-dropdown .option-body {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.autocomplete-dropdown .option-title-row {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.autocomplete-dropdown .option-label {
  font-size: 13px;
  font-weight: 600;
  color: #0F172A;
  word-break: break-word;
}

.autocomplete-dropdown .option-desc {
  margin: 0;
  font-size: 11px;
  color: #64748B;
  line-height: 1.3;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.autocomplete-dropdown .option-badge {
  display: inline-block;
  padding: 1px 6px;
  border-radius: 4px;
  font-size: 10px;
  font-weight: 700;
  letter-spacing: 0.02em;
}

.autocomplete-dropdown .option-badge.default {
  background: #F1F5F9;
  color: #475569;
}

.autocomplete-dropdown .option-badge.primary {
  background: #DBEAFE;
  color: #1E40AF;
}

.autocomplete-dropdown .option-badge.warning {
  background: #FEF3C7;
  color: #92400E;
}

.autocomplete-dropdown .option-badge.success {
  background: #DCFCE7;
  color: #166534;
}

.autocomplete-dropdown .option-category-tag {
  font-size: 10px;
  font-weight: 600;
  color: #64748B;
  background: #F8FAFC;
  border: 1px solid #E2E8F0;
  padding: 2px 6px;
  border-radius: 4px;
  white-space: nowrap;
}

/* Empty state */
.autocomplete-dropdown.empty-dropdown {
  padding: 12px 14px;
  background: #FAFAFA;
}

.autocomplete-dropdown .empty-state-content {
  display: flex;
  align-items: center;
  gap: 10px;
}

.autocomplete-dropdown .empty-icon {
  font-size: 18px;
}

.autocomplete-dropdown .empty-texts {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.autocomplete-dropdown .empty-title {
  font-size: 12px;
  font-weight: 700;
  color: #334155;
}

.autocomplete-dropdown .empty-sub {
  font-size: 11px;
  color: #64748B;
}

/* Global mark highlight inside option label */
mark.hl-match {
  background: #FEF08A;
  color: #002060;
  font-weight: 800;
  padding: 0 1px;
  border-radius: 2px;
}

/* Transitions */
.dropdown-fade-enter-active,
.dropdown-fade-leave-active {
  transition: opacity 0.15s ease, transform 0.15s ease;
}

.dropdown-fade-enter-from,
.dropdown-fade-leave-to {
  opacity: 0;
  transform: translateY(-4px);
}
</style>
