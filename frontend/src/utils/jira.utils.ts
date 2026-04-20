/**
 * Parse @mentions from text
 * Returns array of mentioned usernames
 */
export function parseMentions(text: string): string[] {
  const mentionRegex = /@(\w+)/g
  const matches = text.match(mentionRegex) || []
  return matches.map(m => m.substring(1))
}

/**
 * Highlight mentions in text with HTML
 */
export function highlightMentions(text: string): string {
  return text.replace(/@(\w+)/g, '<span class="mention">@$1</span>')
}

/**
 * Format date for display
 */
export function formatDate(date: string | Date | undefined): string {
  if (!date) return ''
  const d = typeof date === 'string' ? new Date(date) : date
  return d.toLocaleDateString('id-ID', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  })
}

/**
 * Format date and time for display
 */
export function formatDateTime(date: string | Date | undefined): string {
  if (!date) return ''
  const d = typeof date === 'string' ? new Date(date) : date
  return d.toLocaleDateString('id-ID', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

/**
 * Format relative time (e.g., "2 hours ago")
 */
export function formatRelativeTime(date: string | Date | undefined): string {
  if (!date) return ''
  const d = typeof date === 'string' ? new Date(date) : date
  const now = new Date()
  const diffMs = now.getTime() - d.getTime()
  const diffMins = Math.floor(diffMs / 60000)
  const diffHours = Math.floor(diffMs / 3600000)
  const diffDays = Math.floor(diffMs / 86400000)

  if (diffMins < 1) return 'Baru saja'
  if (diffMins < 60) return `${diffMins} menit lalu`
  if (diffHours < 24) return `${diffHours} jam lalu`
  if (diffDays < 7) return `${diffDays} hari lalu`
  return formatDate(d)
}

/**
 * Check if date is overdue
 */
export function isOverdue(dueDate: string | undefined): boolean {
  if (!dueDate) return false
  return new Date(dueDate) < new Date()
}

/**
 * Check if date is today
 */
export function isToday(date: string | undefined): boolean {
  if (!date) return false
  const d = new Date(date)
  const today = new Date()
  return d.toDateString() === today.toDateString()
}

/**
 * Check if date is tomorrow
 */
export function isTomorrow(date: string | undefined): boolean {
  if (!date) return false
  const d = new Date(date)
  const tomorrow = new Date()
  tomorrow.setDate(tomorrow.getDate() + 1)
  return d.toDateString() === tomorrow.toDateString()
}

/**
 * Get status color for display
 */
export function getStatusColor(status: string): string {
  const colors: Record<string, string> = {
    'To Do': 'bg-gray-100 text-gray-800',
    'In Progress': 'bg-blue-100 text-blue-800',
    'In Review': 'bg-yellow-100 text-yellow-800',
    'Done': 'bg-green-100 text-green-800',
    'Backlog': 'bg-gray-100 text-gray-800',
  }
  return colors[status] || 'bg-gray-100 text-gray-800'
}

/**
 * Get priority color for display
 */
export function getPriorityColor(priority: string): string {
  const colors: Record<string, string> = {
    'Highest': 'bg-red-100 text-red-800',
    'High': 'bg-orange-100 text-orange-800',
    'Medium': 'bg-yellow-100 text-yellow-800',
    'Low': 'bg-blue-100 text-blue-800',
    'Lowest': 'bg-gray-100 text-gray-800',
  }
  return colors[priority] || 'bg-gray-100 text-gray-800'
}

/**
 * Get issue type icon
 */
export function getIssueTypeIcon(issueType: string): string {
  const icons: Record<string, string> = {
    'Bug': '🐛',
    'Task': '✓',
    'Story': '📖',
    'Epic': '🎯',
    'Sub-task': '↳',
  }
  return icons[issueType] || '📝'
}

/**
 * Format file size for display
 */
export function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 Bytes'
  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return Math.round((bytes / Math.pow(k, i)) * 100) / 100 + ' ' + sizes[i]
}

/**
 * Validate email format
 */
export function isValidEmail(email: string): boolean {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  return emailRegex.test(email)
}

/**
 * Truncate text to specified length
 */
export function truncateText(text: string, length: number): string {
  if (text.length <= length) return text
  return text.substring(0, length) + '...'
}

/**
 * Get initials from name
 */
export function getInitials(name: string): string {
  return name
    .split(' ')
    .map(n => n[0])
    .join('')
    .toUpperCase()
    .substring(0, 2)
}

/**
 * Generate random color for labels
 */
export function generateRandomColor(): string {
  const colors = [
    '#FF6B6B', '#4ECDC4', '#45B7D1', '#FFA07A', '#98D8C8',
    '#F7DC6F', '#BB8FCE', '#85C1E2', '#F8B88B', '#A9DFBF',
  ]
  return colors[Math.floor(Math.random() * colors.length)]
}

/**
 * Calculate sprint progress percentage
 */
export function calculateSprintProgress(completed: number, total: number): number {
  if (total === 0) return 0
  return Math.round((completed / total) * 100)
}

/**
 * Get days remaining in sprint
 */
export function getDaysRemaining(endDate: string | undefined): number {
  if (!endDate) return 0
  const end = new Date(endDate)
  const now = new Date()
  const diff = end.getTime() - now.getTime()
  return Math.ceil(diff / (1000 * 60 * 60 * 24))
}

/**
 * Check if sprint is active
 */
export function isSprintActive(status: string): boolean {
  return status === 'Active'
}

/**
 * Check if sprint is completed
 */
export function isSprintCompleted(status: string): boolean {
  return status === 'Completed'
}

/**
 * Check if sprint is planned
 */
export function isSprintPlanned(status: string): boolean {
  return status === 'Planned'
}

/**
 * Sort records by priority
 */
export function sortByPriority(records: any[]): any[] {
  const priorityOrder = { 'Highest': 0, 'High': 1, 'Medium': 2, 'Low': 3, 'Lowest': 4 }
  return [...records].sort((a, b) => {
    const aPriority = priorityOrder[a.priority as keyof typeof priorityOrder] ?? 5
    const bPriority = priorityOrder[b.priority as keyof typeof priorityOrder] ?? 5
    return aPriority - bPriority
  })
}

/**
 * Sort records by due date
 */
export function sortByDueDate(records: any[]): any[] {
  return [...records].sort((a, b) => {
    if (!a.due_date) return 1
    if (!b.due_date) return -1
    return new Date(a.due_date).getTime() - new Date(b.due_date).getTime()
  })
}

/**
 * Sort records by creation date (newest first)
 */
export function sortByCreatedDate(records: any[]): any[] {
  return [...records].sort((a, b) => {
    return new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
  })
}

/**
 * Filter records by status
 */
export function filterByStatus(records: any[], status: string): any[] {
  return records.filter(r => r.status === status)
}

/**
 * Filter records by assignee
 */
export function filterByAssignee(records: any[], assigneeId: string): any[] {
  return records.filter(r => r.assigned_to === assigneeId)
}

/**
 * Filter records by label
 */
export function filterByLabel(records: any[], labelId: string): any[] {
  return records.filter(r => r.labels?.some((l: any) => l.id === labelId))
}

/**
 * Filter records by due date range
 */
export function filterByDueDateRange(records: any[], startDate: Date, endDate: Date): any[] {
  return records.filter(r => {
    if (!r.due_date) return false
    const d = new Date(r.due_date)
    return d >= startDate && d <= endDate
  })
}

/**
 * Group records by status
 */
export function groupByStatus(records: any[]): Record<string, any[]> {
  return records.reduce((acc, record) => {
    const status = record.status || 'Unassigned'
    if (!acc[status]) acc[status] = []
    acc[status].push(record)
    return acc
  }, {} as Record<string, any[]>)
}

/**
 * Group records by assignee
 */
export function groupByAssignee(records: any[]): Record<string, any[]> {
  return records.reduce((acc, record) => {
    const assignee = record.assigned_to || 'Unassigned'
    if (!acc[assignee]) acc[assignee] = []
    acc[assignee].push(record)
    return acc
  }, {} as Record<string, any[]>)
}

/**
 * Check if record is overdue
 */
export function isRecordOverdue(record: any): boolean {
  return record.due_date ? isOverdue(record.due_date) : false
}

/**
 * Check if record is completed
 */
export function isRecordCompleted(record: any): boolean {
  return record.is_completed || record.status === 'Done'
}

/**
 * Get record status badge color
 */
export function getRecordStatusBadgeColor(record: any): string {
  if (isRecordCompleted(record)) return 'bg-green-100 text-green-800'
  if (isRecordOverdue(record)) return 'bg-red-100 text-red-800'
  return getStatusColor(record.status)
}
