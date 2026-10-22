import { exec } from 'child_process'
import path from 'path'
import fs from 'fs'
import { defineEventHandler, getQuery, createError } from 'h3'

export default defineEventHandler(async (event) => {
  const query = getQuery(event)
  const targetPath = query.path as string
  const appMode = (query.app as string) || 'antigravity' // Default to Antigravity IDE or Finder

  if (!targetPath) {
    throw createError({ statusCode: 400, statusMessage: 'Path mandatory' })
  }

  // Resolve absolute path relative to workspace root (parent of 00_INIT)
  const absolutePath = path.resolve(process.cwd(), targetPath)
  const vscodeUrl = `vscode://file${absolutePath}`
  const antigravityUrl = `antigravity://file${absolutePath}`

  // Ensure directory exists if path ends with / or has no extension
  let pathExists = fs.existsSync(absolutePath)
  if (!pathExists) {
    try {
      const isDir = targetPath.endsWith('/') || !path.extname(targetPath)
      if (isDir) {
        fs.mkdirSync(absolutePath, { recursive: true })
      } else {
        const parentDir = path.dirname(absolutePath)
        if (!fs.existsSync(parentDir)) {
          fs.mkdirSync(parentDir, { recursive: true })
        }
        fs.writeFileSync(absolutePath, `# ${path.basename(absolutePath)}\n`, 'utf-8')
      }
      pathExists = true
    } catch (err: any) {
      console.error('Error auto-creating path:', err)
    }
  }

  let command = `open "${absolutePath}"`
  let appName = 'System Finder'

  if (appMode === 'antigravity' || appMode === 'agy') {
    command = `open -a "Antigravity IDE" "${absolutePath}"`
    appName = 'Antigravity IDE'
  } else if (appMode === 'vscode' || appMode === 'code') {
    command = `open -a "Visual Studio Code" "${absolutePath}" || code "${absolutePath}"`
    appName = 'VS Code'
  }

  return new Promise((resolve) => {
    exec(command, (error) => {
      if (error) {
        // Fallback to default open command
        exec(`open "${absolutePath}"`, (err2) => {
          if (err2) {
            resolve({
              success: false,
              message: `Gagal membuka: ${err2.message}`,
              path: absolutePath,
              vscodeUrl,
            })
          } else {
            resolve({
              success: true,
              message: `Dibuka di Finder (Fallback dari ${appName})`,
              path: absolutePath,
              vscodeUrl,
            })
          }
        })
      } else {
        resolve({
          success: true,
          message: `Dibuka di ${appName}`,
          path: absolutePath,
          vscodeUrl,
          antigravityUrl,
        })
      }
    })
  })
})
