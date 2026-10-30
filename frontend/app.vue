<template>
  <div class="sizing-app-root">
    <NuxtLayout>
      <NuxtPage />
    </NuxtLayout>
  </div>
</template>

<style>
body {
  margin: 0;
  padding: 0;
  font-family: 'IBM Plex Sans', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  background-color: var(--ground, #f4f7f9);
  color: var(--text, #131416);
  -webkit-font-smoothing: antialiased;
}

* {
  box-sizing: border-box;
}
</style>
