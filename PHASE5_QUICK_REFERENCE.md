# Phase 5: Quick Reference Guide

## Using the Jira Service

```typescript
import { jiraService } from '@/services/jira.service'

// Issue Types
const issueTypes = await jiraService.listIssueTypes(projectId)
const scheme = await jiraService.getIssueTypeScheme(projectId)
await jiraService.createIssueTypeScheme(projectId, { name: 'Default', issue_type_ids: [...] })

// Custom Fields
const field = await jiraService.createCustomField(projectId, { name: 'Priority', field_type: 'dropdown', ... })
const fields = await jiraService.listCustomFields(projectId)
await jiraService.updateCustomField(projectId, fieldId, { name: 'New Name' })
await jiraService.deleteCustomField(projectId, fieldId)

// Workflows
const workflow = await jiraService.getWorkflow(projectId)
await jiraService.createWorkflow(projectId, { name: 'Default', initial_status: 'To Do', statuses: [...] })
await jiraService.updateWorkflow(projectId, { statuses: [...] })
await jiraService.transitionRecord(recordId, { to_status_id: statusId })

// Sprints
const sprint = await jiraService.createSprint(projectId, { name: 'Sprint 1', end_date: '2024-12-31' })
const sprints = await jiraService.listSprints(projectId)
const activeSprint = await jiraService.getActiveSprint(projectId)
await jiraService.startSprint(sprintId)
const metrics = await jiraService.completeSprint(sprintId)
const records = await jiraService.getSprintRecords(sprintId)

// Backlog
const backlog = await jiraService.getBacklog(projectId)
await jiraService.reorderBacklog(projectId, { record_ids: [...] })
await jiraService.bulkAssignToSprint(projectId, { sprint_id: sprintId, record_ids: [...] })

// Comments
const comment = await jiraService.addComment(recordId, { text: 'Hello @user' })
const comments = await jiraService.listComments(recordId)
await jiraService.updateComment(commentId, { text: 'Updated' })
await jiraService.deleteComment(commentId)

// Attachments
const attachment = await jiraService.uploadAttachment(recordId, file)
const attachments = await jiraService.listAttachments(recordId)
await jiraService.deleteAttachment(attachmentId)

// Labels
const label = await jiraService.createLabel(projectId, { name: 'Bug', color: '#FF0000' })
const labels = await jiraService.listLabels(projectId)
await jiraService.addLabelToRecord(recordId, labelId)
await jiraService.removeLabelFromRecord(recordId, labelId)
await jiraService.deleteLabel(labelId)

// Bulk Operations
await jiraService.bulkChangeStatus(projectId, { record_ids: [...], status_id: statusId })
await jiraService.bulkAssignTo(projectId, { record_ids: [...], assignee_id: userId })
await jiraService.bulkAddLabel(projectId, { record_ids: [...], label_id: labelId })
await jiraService.bulkDelete(projectId, { record_ids: [...] })

// Search & Filters
const results = await jiraService.searchRecords(projectId, 'query', { status: 'In Progress' })
await jiraService.saveFilter(projectId, { name: 'My Filter', filters: {...} })
const filters = await jiraService.listSavedFilters(projectId)
```

## Using the Hooks

### useJiraBoard - Sprint Board Management

```typescript
import { useJiraBoard } from '@/hooks/useJiraBoard'

function SprintBoard() {
  const {
    sprint,
    records,
    loading,
    error,
    transitioningRecordId,
    handleDragEnd,
    getRecordsByStatus,
    refresh,
    clearError,
  } = useJiraBoard(projectId, sprintId)

  if (loading) return <div>Loading...</div>
  if (error) return <div>{error} <button onClick={clearError}>Dismiss</button></div>

  const todoRecords = getRecordsByStatus('To Do')
  const inProgressRecords = getRecordsByStatus('In Progress')

  return (
    <div>
      <h1>{sprint?.name}</h1>
      <div onDragEnd={handleDragEnd}>
        {/* Render columns */}
      </div>
    </div>
  )
}
```

### useBacklog - Backlog Management

```typescript
import { useBacklog } from '@/hooks/useBacklog'

function BacklogView() {
  const {
    records,
    sprints,
    loading,
    error,
    selectedRecords,
    reordering,
    assigning,
    reorderRecords,
    assignToSprint,
    toggleRecordSelection,
    clearSelection,
    refresh,
    clearError,
  } = useBacklog(projectId)

  const handleAssignToSprint = async (sprintId: string) => {
    await assignToSprint(sprintId, Array.from(selectedRecords))
  }

  return (
    <div>
      {error && <div>{error}</div>}
      {records.map(record => (
        <div key={record.id}>
          <input
            type="checkbox"
            checked={selectedRecords.has(record.id)}
            onChange={() => toggleRecordSelection(record.id)}
          />
          {record.title}
        </div>
      ))}
      <button onClick={() => handleAssignToSprint(sprints[0].id)} disabled={assigning}>
        Assign to Sprint
      </button>
    </div>
  )
}
```

### useComments - Comment Management

```typescript
import { useComments } from '@/hooks/useComments'

function CommentSection() {
  const {
    comments,
    loading,
    error,
    addingComment,
    updatingCommentId,
    deletingCommentId,
    addComment,
    updateComment,
    deleteComment,
    parseMentions,
    refresh,
    clearError,
  } = useComments(recordId)

  const handleAddComment = async (text: string) => {
    const mentions = parseMentions(text)
    // Send notifications for mentions
    await addComment(text)
  }

  return (
    <div>
      {error && <div>{error}</div>}
      <textarea onBlur={(e) => handleAddComment(e.target.value)} />
      {comments.map(comment => (
        <div key={comment.id}>
          <p>{comment.text}</p>
          <button onClick={() => deleteComment(comment.id)} disabled={deletingCommentId === comment.id}>
            Delete
          </button>
        </div>
      ))}
    </div>
  )
}
```

### useAttachments - Attachment Management

```typescript
import { useAttachments } from '@/hooks/useAttachments'

function AttachmentSection() {
  const {
    attachments,
    loading,
    error,
    uploading,
    uploadProgress,
    deletingAttachmentId,
    uploadAttachment,
    deleteAttachment,
    isImageFile,
    refresh,
    clearError,
  } = useAttachments(recordId)

  const handleFileUpload = async (file: File) => {
    await uploadAttachment(file)
  }

  return (
    <div>
      {error && <div>{error}</div>}
      <input type="file" onChange={(e) => e.target.files && handleFileUpload(e.target.files[0])} />
      {uploading && <div>Uploading: {uploadProgress}%</div>}
      {attachments.map(attachment => (
        <div key={attachment.id}>
          {isImageFile(attachment.file_type) && <img src={attachment.file_path} alt={attachment.file_name} />}
          <p>{attachment.file_name}</p>
          <button onClick={() => deleteAttachment(attachment.id)} disabled={deletingAttachmentId === attachment.id}>
            Delete
          </button>
        </div>
      ))}
    </div>
  )
}
```

## Using Utility Functions

```typescript
import * as jiraUtils from '@/utils/jira.utils'

// Text Processing
const mentions = jiraUtils.parseMentions('@user1 @user2')
const highlighted = jiraUtils.highlightMentions('@user hello')
const truncated = jiraUtils.truncateText('Long text', 10)
const initials = jiraUtils.getInitials('John Doe')

// Date Formatting
const formatted = jiraUtils.formatDate('2024-12-31')
const withTime = jiraUtils.formatDateTime('2024-12-31T10:30:00')
const relative = jiraUtils.formatRelativeTime('2024-12-31T10:30:00')
const overdue = jiraUtils.isOverdue('2024-01-01')

// Status & Priority
const statusColor = jiraUtils.getStatusColor('In Progress')
const priorityColor = jiraUtils.getPriorityColor('High')
const icon = jiraUtils.getIssueTypeIcon('Bug')

// File Handling
const fileSize = jiraUtils.formatFileSize(1024000)

// Sprint Management
const progress = jiraUtils.calculateSprintProgress(5, 10)
const daysLeft = jiraUtils.getDaysRemaining('2024-12-31')
const isActive = jiraUtils.isSprintActive('Active')

// Record Operations
const filtered = jiraUtils.filterByStatus(records, 'In Progress')
const grouped = jiraUtils.groupByStatus(records)
const sorted = jiraUtils.sortByPriority(records)
const isCompleted = jiraUtils.isRecordCompleted(record)
```

## Error Handling Pattern

All hooks follow this error handling pattern:

```typescript
const { error, clearError } = useHook(...)

// Errors are automatically dismissed after 4 seconds
// Or manually clear with:
clearError()

// In components:
{error && (
  <div className="error-banner">
    {error}
    <button onClick={clearError}>Dismiss</button>
  </div>
)}
```

## Loading States Pattern

All hooks provide granular loading states:

```typescript
const {
  loading,           // Initial data load
  uploading,         // File upload in progress
  reordering,        // Reorder operation in progress
  assigning,         // Sprint assignment in progress
  addingComment,     // Comment addition in progress
  updatingCommentId, // Specific comment being updated
  deletingCommentId, // Specific comment being deleted
} = useHook(...)
```

## Type Safety

All types are exported from `frontend/src/types/jira.ts`:

```typescript
import type {
  IssueType,
  CustomField,
  Workflow,
  Sprint,
  Comment,
  Attachment,
  Label,
  JiraProjectRecord,
  SprintMetrics,
  SearchFilters,
} from '@/types/jira'
```

## Best Practices

1. **Always check loading state** before rendering data
2. **Handle errors gracefully** with user-friendly messages
3. **Use optimistic updates** for better UX
4. **Memoize callbacks** to prevent unnecessary re-renders
5. **Clean up subscriptions** in useEffect cleanup
6. **Validate user input** before API calls
7. **Use TypeScript** for type safety
8. **Test error scenarios** thoroughly

## Common Patterns

### Refresh Data
```typescript
const { refresh } = useHook(...)
// Manually refresh data
await refresh()
```

### Handle Async Operations
```typescript
const { addingComment } = useComments(recordId)
const handleAdd = async () => {
  await addComment('text')
  // Component will re-render automatically
}
```

### Batch Operations
```typescript
const { selectedRecords, assignToSprint } = useBacklog(projectId)
await assignToSprint(sprintId, Array.from(selectedRecords))
```

## Troubleshooting

**Issue**: Hook returns undefined data
- **Solution**: Check if recordId/projectId/sprintId is defined

**Issue**: Error not clearing
- **Solution**: Call `clearError()` manually or wait 4 seconds

**Issue**: Optimistic update fails
- **Solution**: Check error handling, data is reverted on failure

**Issue**: Upload progress not updating
- **Solution**: Check file size, ensure proper error handling

## Performance Tips

1. Use `useCallback` for event handlers
2. Memoize expensive computations
3. Avoid inline object/array creation in dependencies
4. Use `React.memo` for components receiving hooks
5. Batch multiple operations when possible
