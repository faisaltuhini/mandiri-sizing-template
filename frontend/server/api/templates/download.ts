import path from 'path'
import fs from 'fs'
import { defineEventHandler, getQuery, createError, setHeader, sendStream } from 'h3'

export default defineEventHandler(async (event) => {
  const query = getQuery(event)
  const relPath = query.path as string

  if (!relPath) {
    throw createError({ statusCode: 400, statusMessage: 'Path parameter is required' })
  }

  // Prevent path traversal attacks
  const safeRelPath = path.normalize(relPath).replace(/^(\.\.[\/\\])+/, '')
  const absolutePath = path.resolve(process.cwd(), safeRelPath)

  if (!fs.existsSync(absolutePath) || !fs.statSync(absolutePath).isFile()) {
    throw createError({ statusCode: 404, statusMessage: 'File template tidak ditemukan' })
  }

  const filename = path.basename(absolutePath)
  const ext = path.extname(filename).toLowerCase()

  let contentType = 'application/octet-stream'
  if (ext === '.xlsx') {
    contentType = 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet'
  } else if (ext === '.docx') {
    contentType = 'application/vnd.openxmlformats-officedocument.wordprocessingml.document'
  } else if (ext === '.md') {
    contentType = 'text/markdown; charset=utf-8'
  }

  setHeader(event, 'Content-Type', contentType)
  setHeader(event, 'Content-Disposition', `attachment; filename="${encodeURIComponent(filename)}"`)
  setHeader(event, 'Content-Length', fs.statSync(absolutePath).size.toString())

  const stream = fs.createReadStream(absolutePath)
  return sendStream(event, stream)
})
