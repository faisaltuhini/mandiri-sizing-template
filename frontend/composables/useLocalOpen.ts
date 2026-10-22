export interface ToastNotification {
  message: string
  type: 'info' | 'success' | 'error'
  visible: boolean
  vscodeUrl?: string
}

const toastState = ref<ToastNotification>({
  message: '',
  type: 'info',
  visible: false,
})

let toastTimer: any = null

export function useLocalOpen() {
  function showToast(message: string, type: 'info' | 'success' | 'error' = 'info', vscodeUrl?: string) {
    if (toastTimer) clearTimeout(toastTimer)
    toastState.value = {
      message,
      type,
      visible: true,
      vscodeUrl,
    }
    toastTimer = setTimeout(() => {
      toastState.value.visible = false
    }, 4000)
  }

  function isLocalPath(pathStr: string): boolean {
    if (!pathStr || pathStr === '#') return false
    if (pathStr.startsWith('http://') || pathStr.startsWith('https://') || pathStr.startsWith('mailto:') || pathStr.startsWith('tel:')) {
      return false
    }
    return true
  }

  async function openPath(pathStr: string, app: 'antigravity' | 'vscode' | 'finder' = 'antigravity') {
    if (!pathStr || pathStr === '#') {
      showToast('Link belum dikonfigurasi (#)', 'info')
      return
    }

    if (!isLocalPath(pathStr)) {
      window.open(pathStr, '_blank', 'noopener,noreferrer')
      return
    }

    const appLabel = app === 'antigravity' ? 'Antigravity IDE' : app === 'vscode' ? 'VS Code' : 'Finder'
    showToast(`🚀 Membuka ${pathStr} di ${appLabel}...`, 'info')

    try {
      const res: any = await $fetch('/api/open', {
        params: { path: pathStr, app },
      })

      if (res.success) {
        showToast(`✅ ${res.message || 'Berhasil dibuka'} (${pathStr})`, 'success', res.vscodeUrl)
      } else {
        if (res.vscodeUrl && (app === 'vscode' || app === 'antigravity')) {
          window.location.href = res.vscodeUrl
          showToast(`💻 Membuka via URL Protocol...`, 'info', res.vscodeUrl)
        } else {
          showToast(`⚠️ ${res.message || 'Gagal membuka folder lokal'}`, 'error')
        }
      }
    } catch (err: any) {
      console.error('Error opening local path:', err)
      showToast(`⚠️ Gagal terhubung ke API lokal: ${err.message || err}`, 'error')
    }
  }

  return {
    toastState: readonly(toastState),
    showToast,
    isLocalPath,
    openPath,
  }
}
