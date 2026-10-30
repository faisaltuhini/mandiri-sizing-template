<template>
  <div
    ref="wrapperRef"
    class="base-help-tip-wrapper"
    @mouseenter="handleMouseEnter"
    @mouseleave="handleMouseLeave"
  >
    <button
      type="button"
      class="help-tip-btn"
      :class="{ active: show }"
      @click.stop="toggleShow"
      aria-label="Bantuan informasi parameter"
    >
      <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round">
        <circle cx="12" cy="12" r="10"></circle>
        <line x1="12" y1="16" x2="12" y2="12"></line>
        <line x1="12" y1="8" x2="12.01" y2="8"></line>
      </svg>
    </button>

    <transition name="tip-fade">
      <div
        v-if="show"
        class="help-popover"
        :class="[
          placement === 'top' ? 'pos-top' : 'pos-bottom',
          computedAlign === 'right' ? 'align-right' : 'align-left'
        ]"
        @click.stop
      >
        <div class="popover-hdr">
          <span class="popover-title">{{ title }}</span>
          <span v-if="subtitle" class="popover-sub">{{ subtitle }}</span>
        </div>
        <p class="popover-desc">{{ content }}</p>
        <div v-if="businessImpact" class="popover-impact">
          <span class="impact-label">Dampak Bisnis:</span>
          <span class="impact-text">{{ businessImpact }}</span>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const props = withDefaults(
  defineProps<{
    title: string
    subtitle?: string
    content: string
    businessImpact?: string
    placement?: 'bottom' | 'top'
    align?: 'left' | 'right' | 'auto'
  }>(),
  {
    placement: 'bottom',
    align: 'auto'
  }
)

const show = ref(false)
const wrapperRef = ref<HTMLElement | null>(null)
const computedAlign = ref<'left' | 'right'>('left')

function updateAlignment() {
  if (props.align === 'left' || props.align === 'right') {
    computedAlign.value = props.align
    return
  }
  if (wrapperRef.value && typeof window !== 'undefined') {
    const rect = wrapperRef.value.getBoundingClientRect()
    // If within 310px of the right window edge, or on right half of parent/screen
    if (window.innerWidth - rect.right < 310 || rect.left > window.innerWidth * 0.52) {
      computedAlign.value = 'right'
    } else {
      computedAlign.value = 'left'
    }
  }
}

function handleMouseEnter() {
  updateAlignment()
  show.value = true
}

function handleMouseLeave() {
  show.value = false
}

function toggleShow() {
  updateAlignment()
  show.value = !show.value
}
</script>

<style scoped>
.base-help-tip-wrapper {
  position: relative;
  display: inline-flex;
  align-items: center;
}

.help-tip-btn {
  background: none;
  border: none;
  padding: 2px 4px;
  color: #94A3B8;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  transition: all 0.15s ease;
}

.help-tip-btn:hover,
.help-tip-btn.active {
  color: #002060;
  background: #E2E8F0;
}

.help-popover {
  position: absolute;
  width: 290px;
  background: #0F172A;
  color: #FFFFFF;
  border-radius: 8px;
  padding: 12px 14px;
  box-shadow: 0 12px 30px rgba(0, 0, 0, 0.4), 0 0 0 1px rgba(255, 255, 255, 0.12);
  z-index: 99999;
  pointer-events: auto;
  font-size: 11px;
  line-height: 1.45;
}

/* Vertical placement */
.help-popover.pos-bottom {
  top: calc(100% + 8px);
}

.help-popover.pos-top {
  bottom: calc(100% + 8px);
}

/* Horizontal alignment: Left aligned (default) */
.help-popover.align-left {
  left: 0;
  right: auto;
}

.help-popover.align-left.pos-bottom::after {
  content: '';
  position: absolute;
  bottom: 100%;
  left: 12px;
  border-width: 5px;
  border-style: solid;
  border-color: transparent transparent #0F172A transparent;
}

.help-popover.align-left.pos-top::after {
  content: '';
  position: absolute;
  top: 100%;
  left: 12px;
  border-width: 5px;
  border-style: solid;
  border-color: #0F172A transparent transparent transparent;
}

/* Horizontal alignment: Right aligned (for elements near right edge) */
.help-popover.align-right {
  left: auto;
  right: 0;
}

.help-popover.align-right.pos-bottom::after {
  content: '';
  position: absolute;
  bottom: 100%;
  right: 12px;
  border-width: 5px;
  border-style: solid;
  border-color: transparent transparent #0F172A transparent;
}

.help-popover.align-right.pos-top::after {
  content: '';
  position: absolute;
  top: 100%;
  right: 12px;
  border-width: 5px;
  border-style: solid;
  border-color: #0F172A transparent transparent transparent;
}

.popover-hdr {
  display: flex;
  flex-direction: column;
  border-bottom: 1px solid rgba(255, 255, 255, 0.15);
  padding-bottom: 6px;
  margin-bottom: 6px;
}

.popover-title {
  font-weight: 800;
  color: #FCD34D;
  font-size: 12px;
}

.popover-sub {
  font-size: 10px;
  color: #94A3B8;
}

.popover-desc {
  margin: 0;
  color: #E2E8F0;
}

.popover-impact {
  margin-top: 8px;
  padding-top: 6px;
  border-top: 1px dashed rgba(255, 255, 255, 0.15);
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.impact-label {
  font-size: 9px;
  font-weight: 800;
  text-transform: uppercase;
  color: #38BDF8;
  letter-spacing: 0.05em;
}

.impact-text {
  font-size: 10px;
  color: #CBD5E1;
}

.tip-fade-enter-active,
.tip-fade-leave-active {
  transition: opacity 0.15s ease, transform 0.15s ease;
}

.help-popover.pos-bottom.tip-fade-enter-from,
.help-popover.pos-bottom.tip-fade-leave-to {
  opacity: 0;
  transform: translateY(-4px);
}

.help-popover.pos-top.tip-fade-enter-from,
.help-popover.pos-top.tip-fade-leave-to {
  opacity: 0;
  transform: translateY(4px);
}
</style>
