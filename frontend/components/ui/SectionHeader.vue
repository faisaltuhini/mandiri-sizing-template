<template>
  <div class="section-header">
    <span class="section-icon">{{ icon }}</span>
    <span class="section-title">{{ title }}</span>
    <div class="section-line" />
  </div>
</template>

<script setup lang="ts">
defineProps<{ icon: string; title: string }>()
</script>
