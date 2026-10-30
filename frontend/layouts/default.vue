<template>
  <div class="sizing-layout-root">
    <slot />
  </div>
</template>

<script setup lang="ts">
// Clean wrapper layout for Enterprise Sizing Suite
</script>

<style>
.sizing-layout-root {
  min-height: 100dvh;
  display: flex;
  flex-direction: column;
  background-color: var(--ground, #f4f7f9);
  color: var(--text, #131416);
}
</style>
