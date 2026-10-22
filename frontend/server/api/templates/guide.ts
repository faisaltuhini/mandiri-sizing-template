import path from 'path'
import fs from 'fs'
import { defineEventHandler, createError } from 'h3'

export default defineEventHandler(async () => {
  const guidePath = path.resolve(process.cwd(), 'template/CAPLAN/CAPACITY_PLAN_GUIDE.md')

  if (!fs.existsSync(guidePath)) {
    throw createError({ statusCode: 404, statusMessage: 'Guide file not found' })
  }

  const content = fs.readFileSync(guidePath, 'utf-8')
  return {
    success: true,
    content,
  }
})
