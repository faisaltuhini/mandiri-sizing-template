<template>
  <div class="base-section-card" :class="{ 'has-changes': hasChanges, collapsed: isCollapsed }">
    <div
      class="section-card-hdr"
      :class="{ 'is-clickable': collapsible }"
      @click="collapsible && toggleCollapse()"
    >
      <div class="hdr-left">
        <span class="step-badge">{{ stepNumber }}</span>
        <div class="title-wrap">
          <div class="title-row">
            <h3 class="section-main-title">{{ title }}</h3>
            <span v-if="technicalLabel" class="technical-pill">{{ technicalLabel }}</span>
            <span v-if="hasChanges" class="change-indicator" title="Telah disesuaikan dari standar baseline">• Diubah</span>
            <span v-if="isCollapsed && collapsedSummary" class="collapsed-summary-chip" :title="collapsedSummary">
              {{ collapsedSummary }}
            </span>
          </div>
          <p v-if="!isCollapsed || !collapsedSummary" class="section-guide-question">{{ guideQuestion }}</p>
        </div>
      </div>

      <div class="hdr-right" @click.stop>
        <slot name="header-action" />
        <button
          v-if="collapsible"
          type="button"
          class="collapse-btn"
          @click="toggleCollapse"
          :title="isCollapsed ? 'Buka Section (Klik untuk memperluas)' : 'Ciutkan Section (Klik untuk memperkecil)'"
        >
          <span class="collapse-action-hint">{{ isCollapsed ? 'Buka' : 'Tutup' }}</span>
          <svg
            width="14"
            height="14"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2.5"
            stroke-linecap="round"
            stroke-linejoin="round"
            :class="{ rotated: !isCollapsed }"
          >
            <polyline points="6 9 12 15 18 9"></polyline>
          </svg>
        </button>
      </div>
    </div>

    <transition name="collapse">
      <div v-show="!isCollapsed" class="section-card-body">
        <slot />
      </div>
    </transition>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'

const props = withDefaults(
  defineProps<{
    stepNumber: string | number
    title: string
    guideQuestion: string
    technicalLabel?: string
    hasChanges?: boolean
    collapsible?: boolean
    collapsed?: boolean
    defaultCollapsed?: boolean
    collapsedSummary?: string
  }>(),
  {
    collapsible: false,
    collapsed: undefined,
    defaultCollapsed: false,
    hasChanges: false
  }
)

const emit = defineEmits<{
  (e: 'update:collapsed', value: boolean): void
}>()

const internalCollapsed = ref(props.collapsed !== undefined ? props.collapsed : props.defaultCollapsed)

watch(
  () => props.collapsed,
  (val) => {
    if (val !== undefined) {
      internalCollapsed.value = val
    }
  }
)

const isCollapsed = computed(() => {
  return props.collapsed !== undefined ? props.collapsed : internalCollapsed.value
})

function toggleCollapse() {
  if (!props.collapsible) return
  const nextVal = !isCollapsed.value
  internalCollapsed.value = nextVal
  emit('update:collapsed', nextVal)
}
</script>

<style scoped>
.base-section-card {
  background: #FFFFFF;
  border: 1px solid #E2E8F0;
  border-radius: 8px;
  /* Removed overflow: hidden so popovers aren't clipped */
  transition: all 0.2s ease;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.02);
  position: relative;
}

.base-section-card:hover {
  border-color: #CBD5E1;
}

.base-section-card.has-changes {
  border-left: 3px solid #D4A017;
}

.section-card-hdr {
  padding: 12px 16px;
  background: #F8FAFC;
  border-bottom: 1px solid #E2E8F0;
  border-top-left-radius: 7px;
  border-top-right-radius: 7px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  user-select: none;
  transition: background-color 0.15s ease;
}

.section-card-hdr.is-clickable {
  cursor: pointer;
}

.section-card-hdr.is-clickable:hover {
  background: #F1F5F9;
}

.base-section-card.collapsed .section-card-hdr {
  border-bottom: none;
  border-bottom-left-radius: 7px;
  border-bottom-right-radius: 7px;
}

.hdr-left {
  display: flex;
  align-items: center;
  gap: 12px;
  flex: 1;
  min-width: 0;
}

.step-badge {
  width: 26px;
  height: 26px;
  border-radius: 50%;
  background: #002060;
  color: #FFFFFF;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 800;
  flex-shrink: 0;
}

.title-wrap {
  display: flex;
  flex-direction: column;
  min-width: 0;
}

.title-row {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.section-main-title {
  margin: 0;
  font-size: 13px;
  font-weight: 800;
  color: #002060;
  letter-spacing: -0.01em;
  white-space: nowrap;
}

.technical-pill {
  font-size: 10px;
  font-family: 'JetBrains Mono', monospace;
  color: #64748B;
  background: #F1F5F9;
  padding: 1px 6px;
  border-radius: 3px;
  white-space: nowrap;
}

.change-indicator {
  font-size: 10px;
  font-weight: 700;
  color: #B45309;
  background: #FEF3C7;
  padding: 1px 6px;
  border-radius: 3px;
  white-space: nowrap;
}

.collapsed-summary-chip {
  font-size: 11px;
  font-weight: 600;
  color: #002060;
  background: #EFF6FF;
  border: 1px solid #BFDBFE;
  padding: 1px 8px;
  border-radius: 4px;
  display: inline-flex;
  align-items: center;
  max-width: 480px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.section-guide-question {
  margin: 2px 0 0 0;
  font-size: 11px;
  color: #64748B;
}

.hdr-right {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
}

.collapse-btn {
  background: #FFFFFF;
  border: 1px solid #E2E8F0;
  color: #475569;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 4px;
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 11px;
  font-weight: 600;
  transition: all 0.15s ease;
}

.collapse-btn:hover {
  background: #E2E8F0;
  color: #002060;
  border-color: #CBD5E1;
}

.collapse-btn svg {
  transition: transform 0.2s ease;
}

.collapse-btn svg.rotated {
  transform: rotate(180deg);
}

.collapse-action-hint {
  font-size: 11px;
  color: #64748B;
}

.section-card-body {
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.collapse-enter-active,
.collapse-leave-active {
  transition: opacity 0.2s ease;
}

.collapse-enter-from,
.collapse-leave-to {
  opacity: 0;
}
</style>
