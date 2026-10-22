<template>
  <Transition name="toast">
    <div
      v-if="toastState.visible"
      :class="['toast-notification', `toast-${toastState.type}`]"
    >
      <div class="toast-content">
        <span class="toast-icon">
          {{ toastState.type === 'success' ? '✅' : toastState.type === 'error' ? '⚠️' : 'ℹ️' }}
        </span>
        <span class="toast-message">{{ toastState.message }}</span>
      </div>
      <a
        v-if="toastState.vscodeUrl"
        :href="toastState.vscodeUrl"
        class="toast-vscode-link"
        title="Paksa Buka di VS Code"
      >
        💻 VS Code
      </a>
    </div>
  </Transition>
</template>

<script setup lang="ts">
import { useLocalOpen } from '~/composables/useLocalOpen'

const { toastState } = useLocalOpen()
</script>

<style scoped>
.toast-notification {
  position: fixed;
  bottom: 24px;
  right: 24px;
  z-index: 9999;
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 16px;
  border-radius: 8px;
  background: #1e293b;
  color: #f8fafc;
  box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.25), 0 8px 10px -6px rgba(0, 0, 0, 0.2);
  border: 1px solid rgba(255, 255, 255, 0.1);
  font-size: 13px;
  font-weight: 500;
  max-width: 420px;
}

.toast-success {
  border-left: 4px solid #22c55e;
}

.toast-error {
  border-left: 4px solid #ef4444;
}

.toast-info {
  border-left: 4px solid #3b82f6;
}

.toast-content {
  display: flex;
  align-items: center;
  gap: 8px;
  flex: 1;
}

.toast-icon {
  font-size: 15px;
}

.toast-message {
  word-break: break-word;
}

.toast-vscode-link {
  font-size: 11px;
  font-weight: 600;
  padding: 3px 8px;
  border-radius: 4px;
  background: #007acc;
  color: #ffffff;
  text-decoration: none;
  transition: opacity 0.15s;
  white-space: nowrap;
}

.toast-vscode-link:hover {
  opacity: 0.9;
}

.toast-enter-active,
.toast-leave-active {
  transition: all 0.25s ease;
}

.toast-enter-from,
.toast-leave-to {
  opacity: 0;
  transform: translateY(12px) scale(0.96);
}
</style>
