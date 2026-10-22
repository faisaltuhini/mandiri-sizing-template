// ─────────────────────────────────────────────
// TypeScript Interfaces — Project Hub
// ─────────────────────────────────────────────

export interface TimelineItem {
  label: string
  percent: number
  color: string
  icon: string
}

export interface WorkflowStep {
  step: number
  name: string
  desc: string
  status: 'done' | 'progress' | 'pending'
  docLink?: string
  docLabel?: string
}

export interface Release {
  id: string
  name: string
  date: string
  drb: string
  forum: string
  status: 'progress' | 'done' | 'upcoming'
  isCurrent?: boolean
  features: string[]
  timeline: TimelineItem[]
  links?: LinkDoc[]
  workflow?: WorkflowStep[]
}

export interface LinkDoc {
  icon: string
  name: string
  desc: string
  href: string
  role?: 'sa' | 'sta' | 'dev' | 'official'
  roleBadge?: string
  author?: string
}

export interface TeamMember {
  name: string
  role: string
  wa?: string
  avatar: string
  avatarBg: string
}

export interface TeamGroup {
  title: string
  members: TeamMember[]
}

export interface Contact {
  name: string
  role: string
  wa: string
  note?: string
}

export interface ContactGroup {
  title: string
  contacts: Contact[]
}

export interface DirLink {
  path: string
  label: string
  desc: string
}

export interface TodoEntry {
  text: string
  status: 'done' | 'warn' | 'pending'
}

export interface Project {
  id: string
  name: string
  icon: string
  appId: string
  domain: string
  description: string
  modules?: string
  picSolar?: string
  picPt?: string
  priority?: number
  statusType?: 'active' | 'enterprise'
  workspacePath?: string
  colorA: string
  colorB: string
  heroAccent: string
  heroGlow: string
  releases: Release[]
  links: LinkDoc[]
  team: TeamGroup[]
  keyContacts?: ContactGroup[]
  dirs: DirLink[]
  todos: TodoEntry[]
}

export interface EawgPerson {
  floor: string
  name: string
  role: string
}

export interface GlobalContacts {
  eawg: EawgPerson[]
}

export interface TemplateFile {
  id: string
  name: string
  category: 'CAPLAN' | 'PENGADAAN' | 'TECH_STACK' | string
  categoryName: string
  filename: string
  ext: 'xlsx' | 'docx' | 'md'
  sizeBytes: number
  sizeFormatted: string
  relativePath: string
  description: string
  icon: string
}

export interface TemplateCategory {
  id: string
  name: string
  description: string
  icon: string
  files: TemplateFile[]
}

