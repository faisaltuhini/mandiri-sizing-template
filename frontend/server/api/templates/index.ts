import path from 'path'
import fs from 'fs'
import { defineEventHandler } from 'h3'
import type { TemplateCategory, TemplateFile } from '~/data/types'

function formatBytes(bytes: number): string {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return `${parseFloat((bytes / Math.pow(k, i)).toFixed(1))} ${sizes[i]}`
}

export default defineEventHandler(async () => {
  const baseDir = path.resolve(process.cwd(), 'template')

  const categoryConfigs: Record<string, { name: string; description: string; icon: string }> = {
    CAPLAN: {
      name: 'Master Capacity Planning Enterprise',
      description: 'Template Excel Superpower, Dokumen Word Resmi DAC, dan Panduan Teori Ekstrapolasi Power-Law.',
      icon: '📊',
    },
    TECH_STACK: {
      name: 'Baseline Tech Stack & Sizing Pattern',
      description: 'Template Excel Baseline Arsitektur Bank Mandiri, Komparasi Tier Moderate (IHSP/RMT) vs Critical (Livin Merchant), BOQ Sizing, dan Pricing Matrix.',
      icon: '🏛️',
    },
    PENGADAAN: {
      name: 'Dokumen Pengadaan & Lisensi',
      description: 'Draft Justifikasi & Kebutuhan Lisensi Middleware AppDynamics dan RHEL OS.',
      icon: '📝',
    },
  }

  const fileDescriptions: Record<string, string> = {
    'CAPACITY_PLAN_GUIDE.md': 'Panduan komprehensif rumus matematika, derivasi Power-Law, OpenShift Worker Node sizing, dan database storage 5 tahun.',
    'Capacity_Plan_Superpower_Document.docx': 'Format dokumen resmi pengadaan & persetujuan DAC (Deliverable Acceptance Certificate) arsitektur bank.',
    'Capacity_Plan_Superpower_Template.xlsx': 'Workbook 10-tab interaktif untuk kalkulasi otomatis ekstrapolasi utilisasi server & BOM hardware.',
    'Template_Tech_Stack_Architecture_Sizing.xlsx': 'Workbook resmi 6-sheet mencakup Cover, Baseline Tech Stack (13 layer), Sizing Moderate, Sizing Critical, Pricing Master, dan Kamus Formula.',
    'generate_template_tech_stack.py': 'Python automation script untuk generate seluruh sheet dan formula sizing secara terstandarisasi.',
    'License AppDyn.docx': 'Draft pengajuan lisensi AppDynamics APM & Infrastructure Monitoring.',
    'License RHEL.docx': 'Draft pengajuan lisensi Red Hat Enterprise Linux (RHEL) Core Subscriptions.',
  }

  const fileNamesReadable: Record<string, string> = {
    'CAPACITY_PLAN_GUIDE.md': 'Capacity Plan Theoretical & Practical Guide (MD)',
    'Capacity_Plan_Superpower_Document.docx': 'Master Capacity Plan Word Document Template (.docx)',
    'Capacity_Plan_Superpower_Template.xlsx': 'Master Capacity Plan Superpower Excel Workbook (.xlsx)',
    'Template_Tech_Stack_Architecture_Sizing.xlsx': 'Baseline Tech Stack & Architecture Sizing Template (.xlsx)',
    'generate_template_tech_stack.py': 'Tech Stack Sizing Generator Script (.py)',
    'License AppDyn.docx': 'Procurement License Document — AppDynamics (.docx)',
    'License RHEL.docx': 'Procurement License Document — Red Hat Enterprise Linux (.docx)',
  }

  const categoriesMap: Record<string, TemplateFile[]> = {
    CAPLAN: [],
    TECH_STACK: [],
    PENGADAAN: [],
  }

  if (fs.existsSync(baseDir)) {
    const subdirs = fs.readdirSync(baseDir, { withFileTypes: true })
    for (const dirent of subdirs) {
      if (dirent.isDirectory()) {
        const catKey = dirent.name
        const catPath = path.join(baseDir, dirent.name)
        const files = fs.readdirSync(catPath)

        for (const filename of files) {
          if (filename.startsWith('.')) continue
          const fullPath = path.join(catPath, filename)
          const stat = fs.statSync(fullPath)
          if (!stat.isFile()) continue

          const extName = path.extname(filename).replace('.', '').toLowerCase()
          let icon = '📄'
          if (extName === 'xlsx') icon = '🟢'
          else if (extName === 'docx') icon = '🔵'
          else if (extName === 'md') icon = '🟣'
          else if (extName === 'py') icon = '🐍'

          const relativePath = `template/${dirent.name}/${filename}`

          const item: TemplateFile = {
            id: `${catKey}_${filename.replace(/[^a-zA-Z0-9]/g, '_')}`,
            name: fileNamesReadable[filename] || filename,
            category: catKey,
            categoryName: categoryConfigs[catKey]?.name || catKey,
            filename,
            ext: extName,
            sizeBytes: stat.size,
            sizeFormatted: formatBytes(stat.size),
            relativePath,
            description: fileDescriptions[filename] || 'Template dokumen enterprise.',
            icon,
          }

          if (categoriesMap[catKey]) {
            categoriesMap[catKey].push(item)
          }
        }
      }
    }
  }

  const result: TemplateCategory[] = Object.keys(categoryConfigs).map((catKey) => ({
    id: catKey,
    name: categoryConfigs[catKey].name,
    description: categoryConfigs[catKey].description,
    icon: categoryConfigs[catKey].icon,
    files: categoriesMap[catKey] || [],
  }))

  return {
    success: true,
    categories: result,
  }
})
