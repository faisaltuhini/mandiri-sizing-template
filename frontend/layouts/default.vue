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
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  background-color: #f8fafc;
  color: #1e293b;
}
</style>
