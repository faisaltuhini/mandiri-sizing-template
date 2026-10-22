<template>
  <span class="base-money" :class="{ 'is-unconfirmed': isUnconfirmed }">
    <template v-if="isUnconfirmed || amount === null || amount === undefined">
      <span class="unconfirmed-badge" :title="unconfirmedText">
        <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="12" cy="12" r="10"/>
          <line x1="12" y1="8" x2="12" y2="12"/>
          <line x1="12" y1="16" x2="12.01" y2="16"/>
        </svg>
        <span>{{ unconfirmedText }}</span>
      </span>
    </template>
    <template v-else>
      <span class="currency-prefix">{{ currency }}</span>
      <span class="currency-value">{{ formattedValue }}</span>
      <span v-if="period" class="currency-period">{{ period }}</span>
    </template>
  </span>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(
  defineProps<{
    amount?: number | null
    currency?: string
    period?: string
    isUnconfirmed?: boolean
    unconfirmedText?: string
  }>(),
  {
    amount: 0,
    currency: 'Rp',
    period: '',
    isUnconfirmed: false,
    unconfirmedText: 'Konfirmasi STA',
  }
)

const formattedValue = computed(() => {
  if (props.amount === null || props.amount === undefined) return '0'
  return new Intl.NumberFormat('id-ID').format(Math.round(props.amount))
})
</script>

<style scoped>
.base-money {
  display: inline-flex;
  align-items: baseline;
  gap: 3px;
  font-variant-numeric: tabular-nums;
  white-space: nowrap;
}

.currency-prefix {
  font-size: 0.85em;
  font-weight: 500;
  color: #64748B;
}

.currency-value {
  font-weight: 700;
  color: #0F172A;
}

.currency-period {
  font-size: 0.8em;
  font-weight: 400;
  color: #94A3B8;
  margin-left: 2px;
}

.unconfirmed-badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-size: 11px;
  font-weight: 500;
  color: #D97706;
  background: #FFFBEB;
  border: 1px solid #FDE68A;
  padding: 1px 7px;
  border-radius: 4px;
}
</style>
