<template>
  <div class="guide-reader">
    <!-- Header bar -->
    <div class="guide-header">
      <div class="header-left">
        <span class="guide-badge"> Enterprise Architecture Manual</span>
        <h2>Panduan Komprehensif Rumus, Teori & Metodologi Capacity Planning</h2>
        <p class="guide-subtitle">
          Dokumen standar metodologi matematis arsitektur IT Digital Channel Delivery — Enterprise Architecture
        </p>
      </div>
      <div class="header-actions">
        <button
          type="button"
          class="action-btn agy-action"
          title="Buka File MD Asli di Antigravity IDE"
          @click="openInAgy"
        >
           Open in Antigravity
        </button>
        <button
          type="button"
          class="action-btn finder-action"
          title="Buka Folder di Finder"
          @click="openInFinder"
        >
           Open in Finder
        </button>
      </div>
    </div>

    <!-- Loading / Error / Main View -->
    <div v-if="pending" class="loading-box">
       Memuat dokumen panduan...
    </div>

    <div v-else-if="error" class="error-box">
       Gagal memuat dokumen panduan: {{ error.message }}
    </div>

    <div v-else class="guide-body">
      <!-- Section Navigation Sidebar -->
      <div class="guide-toc">
        <div class="toc-title"> DAFTAR ISI & BAB</div>
        <a
          v-for="(heading, index) in headings"
          :key="index"
          :href="`#section-${index}`"
          :class="['toc-link', `toc-level-${heading.level}`]"
        >
          <span v-if="heading.level === 2" class="toc-bullet"></span>
          <span v-else class="toc-bullet sub"></span>
          {{ heading.text }}
        </a>
      </div>

      <!-- Main Rendered Content -->
      <div class="guide-content">
        <div
          v-for="(block, index) in parsedBlocks"
          :key="index"
          :id="block.isHeading ? `section-${block.headingIndex}` : undefined"
          :class="['guide-block', `block-${block.type}`]"
        >
          <!-- Headings -->
          <component
            :is="block.tag"
            v-if="block.isHeading"
            class="content-heading"
            v-html="renderInline(block.content)"
          />

          <!-- Block Math Formula -->
          <div v-else-if="block.type === 'formula'" class="math-display-card">
            <div class="math-card-badge"> FORMULA EKSTRAPOLASI</div>
            <div class="math-katex-container" v-html="renderBlockMath(block.content)"></div>
          </div>

          <!-- Code Block -->
          <div v-else-if="block.type === 'code'" class="code-card">
            <div class="code-card-header">
              <span> Snippet / Konfigurasi</span>
            </div>
            <pre><code>{{ block.content }}</code></pre>
          </div>

          <!-- Blockquote / Callout -->
          <div v-else-if="block.type === 'quote'" class="callout-card" v-html="renderInline(block.content)">
          </div>

          <!-- Paragraph -->
          <p v-else-if="block.type === 'paragraph'" class="content-p" v-html="renderInline(block.content)"></p>

          <!-- List Items -->
          <div v-else-if="block.type === 'list'" class="content-list-wrapper">
            <ul :class="['content-list', { 'numbered-list': block.isOrdered }]">
              <li
                v-for="(item, i) in block.items"
                :key="i"
                v-html="renderInline(item)"
              ></li>
            </ul>
          </div>

          <!-- Divider -->
          <hr v-else-if="block.type === 'hr'" class="content-hr" />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import katex from 'katex'
import { useLocalOpen } from '~/composables/useLocalOpen'

const { openPath } = useLocalOpen()

const { data, pending, error } = await useFetch<{ success: boolean; content: string }>('/api/templates/guide')

function openInAgy() {
  openPath('template/CAPLAN/CAPACITY_PLAN_GUIDE.md', 'antigravity')
}

function openInFinder() {
  openPath('template/CAPLAN/', 'finder')
}

interface HeadingItem {
  text: string
  level: number
}

interface ContentBlock {
  type: 'heading' | 'paragraph' | 'code' | 'formula' | 'quote' | 'list' | 'hr'
  tag?: string
  isHeading?: boolean
  headingIndex?: number
  content?: string
  items?: string[]
  isOrdered?: boolean
}

const headings = computed<HeadingItem[]>(() => {
  if (!data.value?.content) return []
  const lines = data.value.content.split('\n')
  const list: HeadingItem[] = []

  for (const line of lines) {
    const trimmed = line.trim()
    if (trimmed.startsWith('## ')) {
      list.push({ text: cleanHeadingText(trimmed.replace(/^##\s+/, '')), level: 2 })
    } else if (trimmed.startsWith('### ')) {
      list.push({ text: cleanHeadingText(trimmed.replace(/^###\s+/, '')), level: 3 })
    }
  }
  return list
})

function cleanHeadingText(raw: string): string {
  return raw
    .replace(/`([^`]+)`/g, '$1')
    .replace(/\$([^$]+)\$/g, '$1')
    .replace(/\*\*(.*?)\*\*/g, '$1')
    .replace(/\\cdot/g, '·')
    .replace(/\\times/g, '×')
}

const parsedBlocks = computed<ContentBlock[]>(() => {
  if (!data.value?.content) return []
  const lines = data.value.content.split('\n')
  const blocks: ContentBlock[] = []

  let headingCounter = 0
  let inCode = false
  let inDisplayMath = false
  let codeBuffer: string[] = []
  let mathBuffer: string[] = []
  let listBuffer: string[] = []
  let isListOrdered = false

  function flushList() {
    if (listBuffer.length > 0) {
      blocks.push({
        type: 'list',
        items: [...listBuffer],
        isOrdered: isListOrdered,
      })
      listBuffer = []
      isListOrdered = false
    }
  }

  for (let i = 0; i < lines.length; i++) {
    const line = lines[i]
    const trimmed = line.trim()

    // 1. Handle Code Blocks
    if (trimmed.startsWith('```')) {
      if (inCode) {
        blocks.push({ type: 'code', content: codeBuffer.join('\n') })
        codeBuffer = []
        inCode = false
      } else {
        flushList()
        inCode = true
      }
      continue
    }

    if (inCode) {
      codeBuffer.push(line)
      continue
    }

    // 2. Handle Multi-line Display Math ($$ ... $$)
    if (trimmed.startsWith('$$') && trimmed.endsWith('$$') && trimmed.length > 4) {
      flushList()
      const mathContent = trimmed.slice(2, -2).trim()
      blocks.push({ type: 'formula', content: mathContent })
      continue
    } else if (trimmed === '$$') {
      if (inDisplayMath) {
        blocks.push({ type: 'formula', content: mathBuffer.join('\n') })
        mathBuffer = []
        inDisplayMath = false
      } else {
        flushList()
        inDisplayMath = true
      }
      continue
    }

    if (inDisplayMath) {
      mathBuffer.push(line)
      continue
    }

    // Single-line display math indented or wrapped
    if (trimmed.startsWith('$$') && trimmed.endsWith('$$')) {
      flushList()
      blocks.push({ type: 'formula', content: trimmed.replace(/^\$\$/g, '').replace(/\$\$$/g, '').trim() })
      continue
    }

    // 3. Handle Lists
    const isBullet = /^[-*]\s+/.test(trimmed)
    const isNumber = /^\d+\.\s+/.test(trimmed)

    if (isBullet || isNumber) {
      const itemContent = trimmed.replace(/^[-*]|\d+\.\s+/, '').trim()
      if (listBuffer.length === 0) {
        isListOrdered = isNumber
      }
      listBuffer.push(itemContent)
      continue
    } else if (listBuffer.length > 0 && trimmed === '') {
      flushList()
      continue
    }

    // 4. Horizontal Rule
    if (trimmed === '---') {
      flushList()
      blocks.push({ type: 'hr' })
      continue
    }

    // 5. Headings
    if (trimmed.startsWith('# ')) {
      flushList()
      blocks.push({ type: 'heading', tag: 'h1', isHeading: true, headingIndex: headingCounter++, content: trimmed.replace(/^#\s+/, '') })
      continue
    }
    if (trimmed.startsWith('## ')) {
      flushList()
      blocks.push({ type: 'heading', tag: 'h2', isHeading: true, headingIndex: headingCounter++, content: trimmed.replace(/^##\s+/, '') })
      continue
    }
    if (trimmed.startsWith('### ')) {
      flushList()
      blocks.push({ type: 'heading', tag: 'h3', isHeading: true, headingIndex: headingCounter++, content: trimmed.replace(/^###\s+/, '') })
      continue
    }
    if (trimmed.startsWith('#### ')) {
      flushList()
      blocks.push({ type: 'heading', tag: 'h4', isHeading: true, headingIndex: headingCounter++, content: trimmed.replace(/^####\s+/, '') })
      continue
    }

    // 6. Blockquote
    if (trimmed.startsWith('> ')) {
      flushList()
      blocks.push({ type: 'quote', content: trimmed.replace(/^>\s+/, '') })
      continue
    }

    // 7. Regular Paragraph
    if (trimmed !== '') {
      flushList()
      blocks.push({ type: 'paragraph', content: trimmed })
    }
  }

  flushList()
  return blocks
})

/**
 * KaTeX block rendering
 */
function renderBlockMath(mathExpr?: string): string {
  if (!mathExpr) return ''
  try {
    return katex.renderToString(mathExpr, {
      displayMode: true,
      throwOnError: false,
      strict: false,
    })
  } catch (err: any) {
    return `<div class="katex-error">${mathExpr}</div>`
  }
}

/**
 * Renders inline markdown + KaTeX inline math ($ ... $)
 */
function renderInline(text?: string): string {
  if (!text) return ''

  // Replace inline LaTeX math ($...$) with KaTeX rendered HTML
  let processed = text.replace(/\$([^$]+)\$/g, (_, math) => {
    try {
      return katex.renderToString(math.trim(), {
        displayMode: false,
        throwOnError: false,
        strict: false,
      })
    } catch {
      return `<code class="inline-math">${math}</code>`
    }
  })

  // Format standard markdown tags
  processed = processed
    .replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')
    .replace(/\*(.*?)\*/g, '<em>$1</em>')
    .replace(/`([^`]+)`/g, '<code class="inline-code">$1</code>')

  return processed
}
</script>

<style scoped>
.guide-reader {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.guide-header {
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius-lg);
  padding: 24px 28px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  box-shadow: var(--shadow-sm);
  gap: 20px;
}

.header-left h2 {
  font-size: 20px;
  font-weight: 800;
  color: var(--navy);
  margin-top: 6px;
  margin-bottom: 4px;
}

.guide-subtitle {
  font-size: 13px;
  color: var(--text3);
}

.guide-badge {
  display: inline-block;
  font-size: 11px;
  font-weight: 700;
  text-transform: uppercase;
  color: #2563eb;
  background: #eff6ff;
  border: 1px solid #bfdbfe;
  padding: 3px 10px;
  border-radius: 20px;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-shrink: 0;
}

.action-btn {
  font-family: 'Inter', sans-serif;
  font-size: 12px;
  font-weight: 600;
  padding: 8px 14px;
  border-radius: var(--radius-sm);
  border: 1px solid var(--border2);
  cursor: pointer;
  transition: all 0.15s;
  display: flex;
  align-items: center;
  gap: 6px;
}

.agy-action {
  background: #003b73;
  color: #ffffff;
  border-color: #003b73;
}

.agy-action:hover {
  background: #1b6c80;
}

.finder-action {
  background: var(--surface);
  color: var(--text2);
}

.finder-action:hover {
  background: var(--bg3);
}

.guide-body {
  display: grid;
  grid-template-columns: 260px 1fr;
  gap: 24px;
  align-items: start;
}

.guide-toc {
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius-lg);
  padding: 18px;
  position: sticky;
  top: 24px;
  display: flex;
  flex-direction: column;
  gap: 6px;
  max-height: 85vh;
  overflow-y: auto;
  box-shadow: var(--shadow-sm);
}

.toc-title {
  font-size: 11px;
  font-weight: 800;
  color: var(--navy);
  padding-bottom: 10px;
  border-bottom: 1px solid var(--border);
  letter-spacing: 0.8px;
}

.toc-link {
  font-size: 12px;
  color: var(--text2);
  line-height: 1.4;
  padding: 6px 8px;
  border-radius: 4px;
  transition: all 0.15s;
  display: flex;
  align-items: center;
  gap: 6px;
  word-break: break-word;
}

.toc-link:hover {
  color: var(--brand-blue);
  background: var(--bg3);
}

.toc-bullet {
  font-size: 10px;
  color: var(--brand-blue);
}

.toc-bullet.sub {
  color: var(--text3);
  font-size: 8px;
}

.toc-level-3 {
  padding-left: 18px;
  font-size: 11px;
  color: var(--text3);
}

.guide-content {
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius-lg);
  padding: 36px 40px;
  display: flex;
  flex-direction: column;
  gap: 20px;
  box-shadow: var(--shadow-sm);
  line-height: 1.8;
}

/* Headings Typography */
.content-heading {
  color: var(--navy);
}

h1.content-heading {
  font-size: 24px;
  font-weight: 800;
  border-bottom: 2px solid var(--border);
  padding-bottom: 12px;
  margin-top: 8px;
}

h2.content-heading {
  font-size: 18px;
  font-weight: 700;
  color: var(--brand-blue);
  border-bottom: 1px solid var(--border);
  padding-bottom: 8px;
  margin-top: 24px;
}

h3.content-heading {
  font-size: 15px;
  font-weight: 700;
  color: var(--text);
  margin-top: 14px;
}

h4.content-heading {
  font-size: 14px;
  font-weight: 600;
  color: var(--text2);
  margin-top: 10px;
}

.content-p {
  color: var(--text2);
  font-size: 14px;
}

/* Display Math Card */
.math-display-card {
  background: var(--surface, #ffffff);
  color: #f8fafc;
  padding: 24px 28px;
  border-radius: var(--radius-lg);
  border-left: 5px solid #3b82f6;
  box-shadow: 0 10px 25px -5px rgba(15, 23, 42, 0.2);
  margin: 12px 0;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.math-card-badge {
  font-size: 10px;
  font-weight: 800;
  letter-spacing: 0.8px;
  color: #94a3b8;
  text-transform: uppercase;
}

.math-katex-container {
  overflow-x: auto;
  padding: 6px 0;
  font-size: 1.15em;
}

/* Make KaTeX text sharp and clean */
:deep(.katex-display) {
  margin: 0 !important;
  color: #ffffff;
}

:deep(.katex) {
  font-size: 1.1em;
}

:deep(.katex .mord),
:deep(.katex .mbin),
:deep(.katex .mrel),
:deep(.katex .mopen),
:deep(.katex .mclose) {
  color: inherit;
}

/* Inline math on normal text */
:deep(p .katex),
:deep(li .katex),
:deep(span .katex) {
  color: #1d4ed8;
  background: #eff6ff;
  padding: 1px 6px;
  border-radius: 4px;
  font-size: 1em;
}

/* List Styling */
.content-list-wrapper {
  padding-left: 8px;
}

.content-list {
  padding-left: 20px;
  color: var(--text2);
  font-size: 13.5px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.content-list li {
  line-height: 1.6;
}

/* Code Block */
.code-card {
  background: #0f172a;
  border-radius: var(--radius);
  overflow: hidden;
  margin: 10px 0;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.code-card-header {
  background: #1e293b;
  padding: 8px 14px;
  font-size: 11px;
  color: #94a3b8;
  font-weight: 600;
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
}

.code-card pre {
  padding: 16px;
  color: #e2e8f0;
  font-family: 'JetBrains Mono', monospace;
  font-size: 12.5px;
  overflow-x: auto;
  line-height: 1.6;
}

/* Callout */
.callout-card {
  background: #fffbeb;
  border-left: 4px solid #f59e0b;
  border-radius: 4px;
  padding: 14px 18px;
  font-size: 13.5px;
  color: #92400e;
  line-height: 1.6;
}

.content-hr {
  border: none;
  border-top: 1px solid var(--border);
  margin: 14px 0;
}

.loading-box,
.error-box {
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius-lg);
  padding: 40px;
  text-align: center;
  font-size: 14px;
  color: var(--text3);
}

@media (max-width: 900px) {
  .guide-body {
    grid-template-columns: 1fr;
  }
  .guide-toc {
    position: static;
    max-height: none;
  }
  .guide-content {
    padding: 24px 20px;
  }
}
</style>
