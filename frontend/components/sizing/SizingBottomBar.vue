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
        <BaseMoney :amount="quoteResponse?.subtotal_compute_annual" period="/ thn" />
      </div>
      <div class="metric-pill">
        <span class="pill-lbl">Storage Subtotal</span>
        <BaseMoney :amount="quoteResponse?.subtotal_storage_annual" period="/ thn" />
      </div>
      <div class="metric-pill">
        <span class="pill-lbl">Lisensi Subtotal</span>
        <BaseMoney :amount="quoteResponse?.subtotal_license_annual" period="/ thn" />
      </div>
    </div>

    <div class="cost-bar-right">
      <div class="grand-total-box">
        <span class="grand-total-lbl">ESTIMASI TOTAL OPEX ARSITEKTUR</span>
        <BaseMoney
          :amount="quoteResponse?.grand_total_annual"
          class="grand-total-val"
          period="/ thn"
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
  background: #FFFFFF;
  border-top: 1px solid #CBD5E1;
  box-shadow: 0 -2px 8px rgba(0, 0, 0, 0.05);
  padding: 12px 24px;
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
  gap: 6px;
}

.env-lbl {
  font-size: 11px;
  font-weight: 700;
  color: #475569;
}

.drc-toggle-label {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  user-select: none;
}

.custom-chk {
  width: 15px;
  height: 15px;
  cursor: pointer;
  accent-color: #002060;
}

.drc-text {
  font-size: 11px;
  font-weight: 600;
  color: #334155;
}

.cost-bar-center {
  display: flex;
  align-items: center;
  gap: 12px;
}

.metric-pill {
  display: flex;
  flex-direction: column;
  background: #F8FAFC;
  border: 1px solid #E2E8F0;
  padding: 4px 10px;
  border-radius: 6px;
}

.pill-lbl {
  font-size: 9px;
  font-weight: 700;
  color: #64748B;
  text-transform: uppercase;
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
  font-size: 9px;
  font-weight: 800;
  letter-spacing: 0.05em;
  color: #64748B;
}

.grand-total-val {
  font-size: 18px;
  font-weight: 800;
  color: #002060;
}

.bar-generate-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  background: #002060;
  color: #FFFFFF;
  border: none;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.15s ease;
}

.bar-generate-btn:hover {
  background: #003399;
}
</style>
