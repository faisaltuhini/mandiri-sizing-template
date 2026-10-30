<template>
  <footer class="sticky-cost-bar">
    <div class="cost-bar-left">
      <div class="env-selector-wrap">
        <span class="env-lbl">Environment:</span>
        <BaseSegmented
          :model-value="environment"
          @update:model-value="$emit('update:environment', $event)"
          :options="[
            { value: 'Production', label: 'Production' },
            { value: 'Dev / Non-Prod', label: 'Dev / Non-Prod' }
          ]"
        />
      </div>

      <label class="drc-toggle-label">
        <input
          type="checkbox"
          :checked="includeDRC"
          @change="$emit('update:includeDRC', ($event.target as HTMLInputElement).checked)"
          class="custom-chk"
        />
        <span class="drc-text">Topologi DRC (DC + DRC Active-Active 2x Multiplier)</span>
      </label>
    </div>

    <div class="cost-bar-center">
      <div class="metric-pill">
        <span class="pill-lbl">Compute Subtotal</span>
        <BaseMoney :amount="quoteResponse?.subtotal_compute_annual" />
      </div>
      <div class="metric-pill">
        <span class="pill-lbl">Storage Subtotal</span>
        <BaseMoney :amount="quoteResponse?.subtotal_storage_annual" />
      </div>
      <div class="metric-pill">
        <span class="pill-lbl">Lisensi Subtotal</span>
        <BaseMoney :amount="quoteResponse?.subtotal_license_annual" />
      </div>
    </div>

    <div class="cost-bar-right">
      <div class="grand-total-box">
        <span class="grand-total-lbl">ESTIMASI TOTAL BIAYA ARSITEKTUR (PERPETUAL)</span>
        <BaseMoney
          :amount="quoteResponse?.grand_total_annual"
          class="grand-total-val"
        />
      </div>

      <button
        type="button"
        class="bar-generate-btn"
        @click="$emit('trigger-generate')"
      >
        <span>Generate BOQ &amp; Excel</span>
        <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
          <path d="M5 12h14"/>
          <path d="M12 5l7 7-7 7"/>
        </svg>
      </button>
    </div>
  </footer>
</template>

<script setup lang="ts">
import type { QuoteResponse } from '~/types/sizing'
import BaseMoney from '~/components/base/BaseMoney.vue'
import BaseSegmented from '~/components/base/BaseSegmented.vue'

defineProps<{
  environment: 'Production' | 'Dev / Non-Prod'
  includeDRC: boolean
  quoteResponse: QuoteResponse | null
}>()

defineEmits<{
  (e: 'update:environment', val: 'Production' | 'Dev / Non-Prod'): void
  (e: 'update:includeDRC', val: boolean): void
  (e: 'trigger-generate'): void
}>()
</script>

<style scoped>
.sticky-cost-bar {
  background: var(--surface, #ffffff);
  border-top: 1px solid var(--rule, #edf1f4);
  box-shadow: none;
  padding: 10px 24px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  flex-shrink: 0;
  z-index: 20;
}

.cost-bar-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.env-selector-wrap {
  display: flex;
  align-items: center;
  gap: 8px;
}

.env-lbl {
  font-size: 12px;
  font-weight: 500;
  color: var(--muted, #5d6165);
}

.drc-toggle-label {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  user-select: none;
}

.custom-chk {
  width: 14px;
  height: 14px;
  cursor: pointer;
  accent-color: var(--accent, #003b73);
}

.drc-text {
  font-size: 12px;
  font-weight: 500;
  color: var(--text, #131416);
}

.cost-bar-center {
  display: flex;
  align-items: center;
  gap: 16px;
}

.metric-pill {
  display: flex;
  flex-direction: column;
  background: transparent;
  padding: 0 4px;
}

.pill-lbl {
  font-size: 11px;
  font-weight: 500;
  color: var(--muted, #5d6165);
}

.cost-bar-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.grand-total-box {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
}

.grand-total-lbl {
  font-size: 11px;
  font-weight: 500;
  letter-spacing: 0.02em;
  color: var(--muted, #5d6165);
}

.grand-total-val {
  font-size: var(--fs-title, 20px);
  font-weight: 600;
  color: var(--accent, #003b73);
  font-variant-numeric: tabular-nums;
}

.bar-generate-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 7px 14px;
  background: var(--accent, #003b73);
  color: #FFFFFF;
  border: none;
  border-radius: var(--radius-control, 6px);
  font-size: 12px;
  font-weight: 600;
  cursor: pointer;
  transition: opacity 0.15s ease;
}

.bar-generate-btn:hover {
  opacity: 0.9;
}
</style>
