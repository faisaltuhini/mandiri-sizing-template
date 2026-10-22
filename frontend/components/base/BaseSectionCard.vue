<template>
  <div class="base-section-card" :class="{ 'has-changes': hasChanges, collapsed: isCollapsed }">
    <div class="section-card-hdr" @click="collapsible && (isCollapsed = !isCollapsed)">
      <div class="hdr-left">
        <span class="step-badge">{{ stepNumber }}</span>
        <div class="title-wrap">
          <div class="title-row">
            <h3 class="section-main-title">{{ title }}</h3>
            <span v-if="technicalLabel" class="technical-pill">{{ technicalLabel }}</span>
            <span v-if="hasChanges" class="change-indicator" title="Telah disesuaikan dari standar baseline">• Diubah</span>
          </div>
          <p class="section-guide-question">{{ guideQuestion }}</p>
        </div>
      </div>

      <div class="hdr-right" @click.stop>
        <slot name="header-action" />
        <button
          v-if="collapsible"
          type="button"
          class="collapse-btn"
          @click="isCollapsed = !isCollapsed"
          :title="isCollapsed ? 'Buka Section' : 'Tutup Section'"
        >
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
import { ref } from 'vue'

const props = withDefaults(
  defineProps<{
    stepNumber: string | number
    title: string
    guideQuestion: string
    technicalLabel?: string
    hasChanges?: boolean
    collapsible?: boolean
    defaultCollapsed?: boolean
  }>(),
  {
    collapsible: false,
    defaultCollapsed: false,
    hasChanges: false
  }
)

const isCollapsed = ref(props.defaultCollapsed)
</script>

<style scoped>
.base-section-card {
  background: #FFFFFF;
  border: 1px solid #E2E8F0;
  border-radius: 8px;
  overflow: hidden;
  transition: all 0.2s ease;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.02);
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
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  cursor: default;
}

.base-section-card.collapsed .section-card-hdr {
  border-bottom: none;
}

.hdr-left {
  display: flex;
  align-items: center;
  gap: 12px;
  flex: 1;
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
}

.title-row {
  display: flex;
  align-items: center;
  gap: 8px;
}

.section-main-title {
  margin: 0;
  font-size: 13px;
  font-weight: 800;
  color: #002060;
  letter-spacing: -0.01em;
}

.technical-pill {
  font-size: 10px;
  font-family: 'JetBrains Mono', monospace;
  color: #64748B;
  background: #F1F5F9;
  padding: 1px 6px;
  border-radius: 3px;
}

.change-indicator {
  font-size: 10px;
  font-weight: 700;
  color: #B45309;
  background: #FEF3C7;
  padding: 1px 6px;
  border-radius: 3px;
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
}

.collapse-btn {
  background: none;
  border: none;
  color: #64748B;
  cursor: pointer;
  padding: 4px;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.15s ease;
}

.collapse-btn:hover {
  background: #E2E8F0;
  color: #002060;
}

.collapse-btn svg.rotated {
  transform: rotate(180deg);
}

.section-card-body {
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.collapse-enter-active,
.collapse-leave-active {
  transition: max-height 0.25s ease, opacity 0.2s ease;
  overflow: hidden;
}
</style>
