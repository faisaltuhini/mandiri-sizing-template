<template>
  <div class="template-hub">
    <!-- Top Hero Banner -->
    <div class="hub-hero">
      <div class="hero-left">
        <div class="hero-tag">📁 Architecture Deliverable Hub</div>
        <h1 class="hero-title">Master Enterprise Templates & Capacity Tools</h1>
        <p class="hero-desc">
          Kumpulan template resmi (.xlsx, .docx, .md) arsitektur bank, panduan komprehensif ekstrapolasi Power-Law, dan kalkulator sizing server otomatis.
        </p>
      </div>

      <div class="hero-stats">
        <div class="stat-card">
          <span class="stat-num">7</span>
          <span class="stat-label">Template Berkas</span>
        </div>
        <div class="stat-card">
          <span class="stat-num">13 Layers</span>
          <span class="stat-label">Baseline Tech Stack</span>
        </div>
        <div class="stat-card">
          <span class="stat-num">100%</span>
          <span class="stat-label">EA Forum Verified</span>
        </div>
      </div>
    </div>

    <!-- Navigation Tabs -->
    <div class="hub-tabs">
      <button
        type="button"
        :class="['tab-btn', { active: activeTab === 'techstack' }]"
        @click="activeTab = 'techstack'"
      >
        🏛️ Baseline Tech Stack & Sizing Pattern
      </button>
      <button
        type="button"
        :class="['tab-btn', { active: activeTab === 'catalog' }]"
        @click="activeTab = 'catalog'"
      >
        📦 Deliverable Templates Catalog
      </button>
      <button
        type="button"
        :class="['tab-btn', { active: activeTab === 'calculator' }]"
        @click="activeTab = 'calculator'"
      >
        🧮 Interactive Capacity Calculator
      </button>
      <button
        type="button"
        :class="['tab-btn', { active: activeTab === 'solarwind' }]"
        @click="activeTab = 'solarwind'"
      >
        🛰️ SolarWinds Multi-Node DB Analytics
      </button>
      <button
        type="button"
        :class="['tab-btn', { active: activeTab === 'guide' }]"
        @click="activeTab = 'guide'"
      >
        📖 Comprehensive Guide & Formulas
      </button>
    </div>

    <!-- Tab Contents -->
    <div class="tab-content">
      <!-- TAB 1: Catalog -->
      <div v-if="activeTab === 'catalog'" class="catalog-view">
        <div v-if="pending" class="loading-state">
          ⏳ Memuat katalog template...
        </div>

        <div v-else-if="error" class="error-state">
          ⚠️ Gagal memuat daftar template: {{ error.message }}
        </div>

        <div v-else-if="data" class="categories-list">
          <div
            v-for="cat in data.categories"
            :key="cat.id"
            class="category-section"
          >
            <div class="cat-header">
              <span class="cat-icon">{{ cat.icon }}</span>
              <div>
                <h2>{{ cat.name }}</h2>
                <p>{{ cat.description }}</p>
              </div>
            </div>

            <div class="files-grid">
              <div
                v-for="file in cat.files"
                :key="file.id"
                class="file-card"
              >
                <div class="file-card-top">
                  <span class="file-icon">{{ file.icon }}</span>
                  <div class="file-badge" :class="`badge-${file.ext}`">
                    {{ file.ext.toUpperCase() }}
                  </div>
                </div>

                <h3 class="file-name">{{ file.name }}</h3>
                <p class="file-desc">{{ file.description }}</p>

                <div class="file-meta">
                  <span>📄 {{ file.filename }}</span>
                  <span>⚖️ {{ file.sizeFormatted }}</span>
                </div>

                <div class="file-actions">
                  <a
                    :href="`/api/templates/download?path=${encodeURIComponent(file.relativePath)}`"
                    class="action-btn download-btn"
                    title="Download Langsung Berkas"
                  >
                    📥 Download
                  </a>
                  <button
                    type="button"
                    class="action-btn agy-btn"
                    title="Buka di Antigravity IDE"
                    @click="openLocal(file.relativePath, 'antigravity')"
                  >
                    ✨ AGY
                  </button>
                  <button
                    type="button"
                    class="action-btn finder-btn"
                    title="Buka di System Finder"
                    @click="openLocal(file.relativePath, 'finder')"
                  >
                    📂 Finder
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- TAB: Tech Stack Baseline -->
      <div v-else-if="activeTab === 'techstack'">
        <TemplateTechStackBaselineView />
      </div>

      <!-- TAB 2: Calculator -->
      <div v-else-if="activeTab === 'calculator'">
        <TemplateCapacityCalculator />
      </div>

      <!-- TAB 3: SolarWinds Multi-Node Analytics -->
      <div v-else-if="activeTab === 'solarwind'">
        <TemplateSolarwindMultiNodeView />
      </div>

      <!-- TAB 4: Guide -->
      <div v-else-if="activeTab === 'guide'">
        <TemplateGuideReader />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useLocalOpen } from '~/composables/useLocalOpen'
import type { TemplateCategory } from '~/data/types'

const activeTab = ref<'techstack' | 'catalog' | 'calculator' | 'solarwind' | 'guide'>('techstack')
const { openPath } = useLocalOpen()

const { data, pending, error } = await useFetch<{ success: boolean; categories: TemplateCategory[] }>('/api/templates')

function openLocal(pathStr: string, app: 'antigravity' | 'vscode' | 'finder') {
  openPath(pathStr, app)
}
</script>

<style scoped>
.template-hub {
  padding: 32px;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.hub-hero {
  background: linear-gradient(135deg, #0f2b48 0%, #003b73 100%);
  color: #ffffff;
  padding: 32px 36px;
  border-radius: var(--radius-lg);

  display: flex;
  align-items: center;
  justify-content: space-between;
  box-shadow: var(--shadow-md);
}

.hero-tag {
  display: inline-block;
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.5px;
  text-transform: uppercase;
  padding: 4px 12px;
  background: rgba(255, 255, 255, 0.15);
  border: 1px solid rgba(255, 255, 255, 0.25);
  border-radius: 20px;
  margin-bottom: 12px;
}

.hero-title {
  font-size: 24px;
  font-weight: 800;
  margin-bottom: 8px;

  color: #ffffff;
}

.hero-desc {
  font-size: 14px;
  color: #cbd5e1;

  max-width: 700px;
  line-height: 1.6;
}

.hero-stats {
  display: flex;
  gap: 16px;
}

.stat-card {
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  padding: 16px 20px;
  border-radius: var(--radius);
  display: flex;
  flex-direction: column;

  align-items: center;
  min-width: 120px;
}

.stat-num {
  font-size: 24px;
  font-weight: 800;

  color: #ffffff;
}

.stat-label {
  font-size: 10px;

  color: #94a3b8;
  text-transform: uppercase;
  text-align: center;
}

.hub-tabs {
  display: flex;
  gap: 8px;
  border-bottom: 2px solid var(--border);

  padding-bottom: 2px;
}

.tab-btn {
  font-family: 'Inter', sans-serif;
  font-size: 13px;

  font-weight: 600;
  padding: 10px 18px;
  border: none;
  background: transparent;

  color: var(--text3);
  cursor: pointer;
  border-bottom: 3px solid transparent;

  transition: all 0.15s;
}

.tab-btn:hover {
  color: var(--navy);
}

.tab-btn.active {
  color: var(--brand-blue);
  border-bottom-color: var(--brand-blue);
}

.catalog-view {
  display: flex;
  flex-direction: column;
  gap: 32px;
}

.category-section {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.cat-header {
  display: flex;
  align-items: center;
  gap: 12px;
}

.cat-icon {
  font-size: 28px;
  width: 44px;
  height: 44px;
  display: flex;
  align-items: center;

  justify-content: center;
  background: var(--surface);
  border: 1px solid var(--border);

  border-radius: 10px;
}

.cat-header h2 {
  font-size: 17px;
  font-weight: 700;

  color: var(--navy);
}

.cat-header p {
  font-size: 12px;
  color: var(--text3);
}

.files-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 20px;
}

.file-card {
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius-lg);

  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 12px;

  box-shadow: var(--shadow-sm);
  transition: transform 0.15s, box-shadow 0.15s;
}

.file-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
}

.file-card-top {
  display: flex;
  align-items: center;

  justify-content: space-between;
}

.file-icon {
  font-size: 26px;
}

.file-badge {
  font-size: 10px;
  font-weight: 800;

  padding: 3px 8px;
  border-radius: 4px;
  text-transform: uppercase;
}

.badge-xlsx {
  background: #f0fdf4;
  color: #16a34a;

  border: 1px solid #bbf7d0;
}

.badge-docx {
  background: #eff6ff;

  color: #2563eb;
  border: 1px solid #bfdbfe;
}

.badge-md {
  background: #f5f3ff;
  color: #7c3aed;

  border: 1px solid #ddd6fe;
}

.file-name {
  font-size: 14px;

  font-weight: 700;
  color: var(--text);
  line-height: 1.4;
}

.file-desc {
  font-size: 12px;
  color: var(--text3);

  line-height: 1.5;
  flex: 1;
}

.file-meta {
  display: flex;

  align-items: center;
  justify-content: space-between;
  font-size: 11px;

  color: var(--text3);
  padding-top: 8px;
  border-top: 1px dashed var(--border);
}

.file-actions {
  display: flex;

  align-items: center;
  gap: 8px;
}

.action-btn {
  font-family: 'Inter', sans-serif;

  font-size: 11px;
  font-weight: 600;
  padding: 6px 12px;

  border-radius: var(--radius-sm);
  border: 1px solid var(--border2);
  cursor: pointer;

  transition: all 0.15s;
  display: inline-flex;
  align-items: center;

  justify-content: center;
  text-decoration: none;
}

.download-btn {
  background: var(--brand-blue);
  color: #ffffff;

  border-color: var(--brand-blue);
  flex: 1;
}

.download-btn:hover {
  background: var(--navy);
}

.agy-btn {
  background: rgba(99, 102, 241, 0.1);
  color: #6366f1;

  border-color: rgba(99, 102, 241, 0.3);
}

.agy-btn:hover {
  background: #6366f1;
  color: #ffffff;
}

.finder-btn {
  background: var(--bg3);
  color: var(--text2);

  border-color: var(--border2);
}

.finder-btn:hover {
  background: var(--border2);
}

.loading-state,
.error-state {
  padding: 40px;
  text-align: center;
  background: var(--surface);

  border: 1px solid var(--border);
  border-radius: var(--radius-lg);
  font-size: 14px;
}

@media (max-width: 800px) {
  .hub-hero {
    flex-direction: column;

    align-items: flex-start;
    gap: 16px;
  }
}
</style>
