# Phase 1: Analisis & Persiapan - Detail

**Status**: IN PROGRESS  
**Date**: April 19, 2026

---

## 📊 Analisis Struktur Saat Ini

### Backend Structure

#### Handlers (backend/internal/delivery/http/handler/)
**Project Board Related:**
- ✅ `project_handler.go` - Project CRUD (KEEP - digunakan Jira juga)
- ✅ `jira_handler.go` - Jira features (KEEP - sudah ada)

**Jira Features (KEEP):**
- ✅ `attachment_handler.go` - Attachments
- ✅ `dashboard_handler.go` - Dashboard

**Other Handlers (KEEP):**
- ✅ `approval_handler.go`
- ✅ `auth_handler.go`
- ✅ `notification_handler.go`
- ✅ `org_handler.go`
- ✅ `sso_handler.go`
- ✅ `ticket_handler.go`
- ✅ `user_handler.go`
- ✅ `webhook_handler.go`

#### Repositories (backend/internal/repository/postgres/)
**Project Board Related (KEEP - digunakan Jira):**
- ✅ `project_repository.go` - Projects
- ✅ `project_column_repository.go` - Columns (akan jadi workflow statuses)
- ✅ `project_record_repository.go` - Records (akan jadi issues)
- ✅ `project_activity_log_repository.go` - Activity logs
- ✅ `project_member_repository.go` - Project members

**Jira Features (KEEP):**
- ✅ `issue_type_repository.go`
- ✅ `issue_type_scheme_repository.go`
- ✅ `custom_field_repository.go`
- ✅ `custom_field_option_repository.go`
- ✅ `custom_field_value_repository.go`
- ✅ `workflow_repository.go`
- ✅ `workflow_status_repository.go`
- ✅ `workflow_transition_repository.go`
- ✅ `sprint_repository.go`
- ✅ `sprint_record_repository.go`
- ✅ `comment_repository.go`
- ✅ `comment_mention_repository.go`
- ✅ `attachment_repository.go`
- ✅ `label_repository.go`
- ✅ `record_label_repository.go`

**Other Repositories (KEEP):**
- ✅ `activity_log_repository.go`
- ✅ `approval_repository.go`
- ✅ `department_repository.go`
- ✅ `division_repository.go`
- ✅ `notification_repository.go`
- ✅ `refresh_token_repository.go`
- ✅ `team_repository.go`
- ✅ `ticket_repository.go`
- ✅ `user_repository.go`
- ✅ `webhook_repository.go`

#### Usecases (backend/internal/usecase/)
**Project Board Related (KEEP - digunakan Jira):**
- ✅ `project_usecase.go` - Project management

**Jira Features (KEEP):**
- ✅ `backlog_usecase.go` - Backlog management
- ✅ `sprint_usecase.go` - Sprint management
- ✅ `issue_type_usecase.go` - Issue types
- ✅ `custom_field_usecase.go` - Custom fields
- ✅ `workflow_usecase.go` - Workflows
- ✅ `comment_usecase.go` - Comments
- ✅ `attachment_usecase.go` - Attachments
- ✅ `label_usecase.go` - Labels
- ✅ `bulk_operation_usecase.go` - Bulk operations
- ✅ `assignment_usecase.go` - Assignments

**Other Usecases (KEEP):**
- ✅ `approval_usecase.go`
- ✅ `auth_usecase.go`
- ✅ `dashboard_usecase.go`
- ✅ `notification_usecase.go`
- ✅ `org_usecase.go`
- ✅ `search_usecase.go`
- ✅ `ticket_usecase.go`
- ✅ `user_usecase.go`
- ✅ `webhook_usecase.go`

### Frontend Structure

#### Pages (frontend/src/pages/)
**Project Board Related:**
- ✅ `ProjectBoardPage.tsx` - REPLACE dengan Jira board view
- ✅ `ProjectHomePage.tsx` - KEEP (Jira home)
- ✅ `ProjectCalendarPage.tsx` - KEEP (Jira calendar)
- ✅ `BacklogPage.tsx` - KEEP (Jira backlog)
- ✅ `SprintBoardPage.tsx` - KEEP (Jira sprint)
- ✅ `ProjectSettingsPage.tsx` - KEEP (Jira settings)

**New Pages (KEEP):**
- ✅ `ComponentsPage.tsx` - Components
- ✅ `IssuesPage.tsx` - Issues
- ✅ `ReleasesPage.tsx` - Releases
- ✅ `ReportsPage.tsx` - Reports
- ✅ `RepositoryPage.tsx` - Repository

**Other Pages (KEEP):**
- ✅ `ActivityLogsPage.tsx`
- ✅ `ApprovalsPage.tsx`
- ✅ `DashboardPage.tsx`
- ✅ `KanbanBoardPage.tsx`
- ✅ `LoginPage.tsx`
- ✅ `NotificationsPage.tsx`
- ✅ `OrgStructurePage.tsx`
- ✅ `ProfilePage.tsx`
- ✅ `SSOCallbackPage.tsx`
- ✅ `TicketDetailPage.tsx`
- ✅ `TicketFormPage.tsx`
- ✅ `TicketListPage.tsx`
- ✅ `UserManagementPage.tsx`
- ✅ `WebhookConfigPage.tsx`

#### Components (frontend/src/components/project/)
**Project Board Related (HAPUS):**
- ❌ `ProjectBoardColumn.tsx` - Old column component
- ❌ `ProjectRecordCard.tsx` - Old record card component
- ❌ `ProjectFilterBar.tsx` - Old filter bar

**Jira Features (KEEP):**
- ✅ `BacklogView.tsx` - Backlog view
- ✅ `SprintBoard.tsx` - Sprint board
- ✅ `RecordCard.tsx` - Record card (Jira)
- ✅ `RecordDetailModal.tsx` - Record detail modal
- ✅ `SearchFilterBar.tsx` - Search filter
- ✅ `CommentSection.tsx` - Comments
- ✅ `AttachmentSection.tsx` - Attachments
- ✅ `LabelManager.tsx` - Labels
- ✅ `BulkOperationsBar.tsx` - Bulk operations
- ✅ `CalendarGrid.tsx` - Calendar
- ✅ `CreateProjectDialog.tsx` - Create project
- ✅ `MemberManagement.tsx` - Member management

---

## 🗂️ File-File yang Perlu Diubah

### Backend

**KEEP (Tidak perlu diubah):**
- ✅ Semua handler
- ✅ Semua repository
- ✅ Semua usecase

**UPDATE (Perlu diubah):**
- ⚠️ `backend/internal/delivery/http/router.go` - Update routes untuk Jira only

**HAPUS (Tidak perlu lagi):**
- ❌ Tidak ada file yang perlu dihapus (semua sudah Jira-compatible)

### Frontend

**KEEP (Tidak perlu diubah):**
- ✅ `BacklogPage.tsx`
- ✅ `SprintBoardPage.tsx`
- ✅ `ProjectSettingsPage.tsx`
- ✅ `ProjectHomePage.tsx`
- ✅ `ProjectCalendarPage.tsx`
- ✅ `ComponentsPage.tsx`
- ✅ `IssuesPage.tsx`
- ✅ `ReleasesPage.tsx`
- ✅ `ReportsPage.tsx`
- ✅ `RepositoryPage.tsx`

**REPLACE (Perlu diganti):**
- 🔄 `ProjectBoardPage.tsx` - Ganti dengan Jira board view

**HAPUS (Tidak perlu lagi):**
- ❌ `ProjectBoardColumn.tsx` - Old column component
- ❌ `ProjectRecordCard.tsx` - Old record card component
- ❌ `ProjectFilterBar.tsx` - Old filter bar

**UPDATE (Perlu diubah):**
- ⚠️ `frontend/src/App.tsx` - Update routes untuk Jira only

---

## 📊 Database Schema Analysis

### Current Tables (Project Board)
```
projects
├── id (UUID)
├── name (VARCHAR)
├── icon_color (VARCHAR)
├── created_by (UUID)
├── created_at (TIMESTAMPTZ)
└── updated_at (TIMESTAMPTZ)

project_columns
├── id (UUID)
├── project_id (UUID) → projects
├── name (VARCHAR)
├── position (INTEGER)
└── created_at (TIMESTAMPTZ)

project_records
├── id (UUID)
├── column_id (UUID) → project_columns
├── project_id (UUID) → projects
├── title (VARCHAR)
├── description (TEXT)
├── assigned_to (UUID) → users
├── due_date (DATE)
├── position (INTEGER)
├── is_completed (BOOLEAN)
├── completed_at (TIMESTAMPTZ)
├── created_by (UUID) → users
├── created_at (TIMESTAMPTZ)
└── updated_at (TIMESTAMPTZ)

project_activity_logs
├── id (UUID)
├── project_id (UUID) → projects
├── record_id (UUID) → project_records
├── actor_id (UUID) → users
├── action (VARCHAR)
├── detail (TEXT)
└── created_at (TIMESTAMPTZ)
```

### Jira Tables (Already Exist)
```
issue_types
issue_type_schemes
issue_type_scheme_items
custom_fields
custom_field_options
custom_field_values
workflows
workflow_statuses
workflow_transitions
sprints
sprint_records
comments
comment_mentions
attachments
labels
record_labels
```

### Migration Strategy
```
projects → projects (no change)
project_columns → workflow_statuses (via migration)
project_records → project_records (add Jira fields)
project_activity_logs → project_activity_logs (no change)
```

---

## 🔄 Data Migration Plan

### Step 1: Backup
```sql
-- Backup existing data
CREATE TABLE projects_backup AS SELECT * FROM projects;
CREATE TABLE project_columns_backup AS SELECT * FROM project_columns;
CREATE TABLE project_records_backup AS SELECT * FROM project_records;
CREATE TABLE project_activity_logs_backup AS SELECT * FROM project_activity_logs;
```

### Step 2: Create Workflows
```sql
-- For each project, create a default workflow
INSERT INTO workflows (id, project_id, name, initial_status)
SELECT gen_random_uuid(), id, 'Default Workflow', 'Backlog'
FROM projects;

-- Create workflow statuses from project columns
INSERT INTO workflow_statuses (id, workflow_id, status_name, status_order)
SELECT gen_random_uuid(), w.id, pc.name, pc.position
FROM project_columns pc
JOIN workflows w ON pc.project_id = w.project_id;
```

### Step 3: Create Issue Type Schemes
```sql
-- For each project, create a default issue type scheme
INSERT INTO issue_type_schemes (id, project_id, name)
SELECT gen_random_uuid(), id, 'Default Scheme'
FROM projects;

-- Add all issue types to the scheme
INSERT INTO issue_type_scheme_items (id, scheme_id, issue_type_id)
SELECT gen_random_uuid(), iss.id, it.id
FROM issue_type_schemes iss
CROSS JOIN issue_types it;
```

### Step 4: Update Records
```sql
-- Add Jira fields to project_records
ALTER TABLE project_records
ADD COLUMN issue_type_id UUID REFERENCES issue_types(id),
ADD COLUMN status VARCHAR DEFAULT 'Backlog',
ADD COLUMN parent_record_id UUID REFERENCES project_records(id);

-- Set default issue type (Task) for all records
UPDATE project_records
SET issue_type_id = (SELECT id FROM issue_types WHERE name = 'Task')
WHERE issue_type_id IS NULL;

-- Set status from column name
UPDATE project_records pr
SET status = (
  SELECT ws.status_name
  FROM workflow_statuses ws
  JOIN workflows w ON ws.workflow_id = w.id
  JOIN project_columns pc ON w.project_id = pc.project_id
  WHERE pc.id = pr.column_id
  LIMIT 1
)
WHERE status = 'Backlog';
```

---

## ✅ Phase 1 Checklist

- [x] Identifikasi semua backend files
- [x] Identifikasi semua frontend files
- [x] Analisis database schema
- [x] Buat migration strategy
- [x] Buat backup plan
- [ ] Backup database production
- [ ] Verify backup integrity
- [ ] Create migration scripts
- [ ] Test migration scripts di development

---

## 📋 Summary

### Backend
- **KEEP**: Semua handler, repository, usecase (sudah Jira-compatible)
- **UPDATE**: router.go (update routes)
- **HAPUS**: Tidak ada

### Frontend
- **KEEP**: Semua Jira pages dan components
- **REPLACE**: ProjectBoardPage.tsx
- **HAPUS**: ProjectBoardColumn.tsx, ProjectRecordCard.tsx, ProjectFilterBar.tsx
- **UPDATE**: App.tsx (update routes)

### Database
- **KEEP**: Semua Jira tables
- **MIGRATE**: project_columns → workflow_statuses
- **UPDATE**: project_records (add Jira fields)
- **BACKUP**: Semua data sebelum migrasi

---

## 🎯 Next Steps

1. ✅ Phase 1 Complete - Analisis selesai
2. ⏳ Phase 2 - Migrasi Data
3. ⏳ Phase 3 - Update Backend
4. ⏳ Phase 4 - Update Frontend
5. ⏳ Phase 5 - Testing
6. ⏳ Phase 6 - Deployment

---

**Phase 1 Status**: ✅ COMPLETE

**Ready untuk Phase 2?** 🚀
