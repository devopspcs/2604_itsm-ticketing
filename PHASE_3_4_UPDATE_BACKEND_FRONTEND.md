# Phase 3 & 4: Update Backend & Frontend

**Status**: READY FOR EXECUTION  
**Date**: April 19, 2026

---

## Phase 3: Update Backend

### Status: ✅ ALREADY DONE

**Good News!** Backend sudah 100% Jira-compatible!

#### What's Already in Place:
- ✅ All Jira handlers implemented
- ✅ All Jira repositories implemented
- ✅ All Jira usecases implemented
- ✅ Router already has all Jira endpoints
- ✅ No project board specific code in backend

#### Backend Files Status:
```
backend/internal/delivery/http/handler/
├── jira_handler.go ✅ (Jira features)
├── project_handler.go ✅ (Project management)
└── (other handlers) ✅ (Keep as is)

backend/internal/repository/postgres/
├── (Jira repositories) ✅ (All implemented)
├── project_repository.go ✅ (Keep)
├── project_column_repository.go ✅ (Keep - for migration)
├── project_record_repository.go ✅ (Keep - will be used by Jira)
└── (other repositories) ✅ (Keep as is)

backend/internal/usecase/
├── (Jira usecases) ✅ (All implemented)
├── project_usecase.go ✅ (Keep)
└── (other usecases) ✅ (Keep as is)

backend/internal/delivery/http/router.go ✅ (Already has all Jira routes)
```

#### Conclusion:
**No backend changes needed!** All Jira features are already implemented and integrated.

---

## Phase 4: Update Frontend

### Status: ⏳ IN PROGRESS

#### What Needs to Be Done:

### 1. Delete Old Project Board Components

**Files to Delete:**
```
frontend/src/components/project/
├── ProjectBoardColumn.tsx ❌ (DELETE)
├── ProjectRecordCard.tsx ❌ (DELETE)
└── ProjectFilterBar.tsx ❌ (DELETE)
```

**Reason:** These are old project board components, replaced by Jira components

### 2. Replace ProjectBoardPage

**File to Replace:**
```
frontend/src/pages/ProjectBoardPage.tsx 🔄 (REPLACE)
```

**Current Content:** Old project board with columns and records

**New Content:** Jira board view with:
- Issue types
- Custom fields
- Workflows
- Sprint board
- Backlog integration
- Comments
- Attachments
- Labels

**New Implementation:**
```typescript
// frontend/src/pages/ProjectBoardPage.tsx

import { useState, useEffect } from 'react'
import { useParams } from 'react-router-dom'
import { jiraService } from '../services/jira.service'
import { SprintBoard } from '../components/project/SprintBoard'
import { BacklogView } from '../components/project/BacklogView'

export function ProjectBoardPage() {
  const { id: projectId } = useParams()
  const [sprint, setSprint] = useState(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const fetchSprint = async () => {
      if (!projectId) return
      try {
        const res = await jiraService.getActiveSprint(projectId)
        setSprint(res.data)
      } catch (err) {
        console.error('Failed to fetch sprint:', err)
      } finally {
        setLoading(false)
      }
    }

    fetchSprint()
  }, [projectId])

  if (loading) {
    return <div>Loading...</div>
  }

  if (!sprint) {
    return (
      <div className="flex flex-col h-[calc(100vh-64px)] items-center justify-center">
        <p className="text-on-surface-variant">No active sprint</p>
        <p className="text-sm text-on-surface-variant">Create or start a sprint to see the board</p>
      </div>
    )
  }

  return (
    <div className="flex flex-col h-[calc(100vh-64px)]">
      <div className="px-8 pt-8 pb-4">
        <h1 className="text-2xl font-extrabold text-on-surface">Sprint Board</h1>
      </div>

      <div className="flex-1 overflow-auto px-8 pb-8">
        {projectId && sprint.id && (
          <SprintBoard projectId={projectId} sprintId={sprint.id} statuses={[]} />
        )}
      </div>
    </div>
  )
}
```

### 3. Update App.tsx Routes

**File to Update:**
```
frontend/src/App.tsx ⚠️ (UPDATE)
```

**Current Routes:**
```typescript
<Route path="/projects/:id" element={<ProjectBoardPage />} />
```

**New Routes:** (Already correct, no change needed)
```typescript
<Route path="/projects/:id" element={<ProjectBoardPage />} />
<Route path="/projects/:id/sprint" element={<SprintBoardPage />} />
<Route path="/projects/:id/backlog" element={<BacklogPage />} />
<Route path="/projects/:id/settings" element={<ProjectSettingsPage />} />
```

### 4. Keep Existing Jira Pages

**Pages to Keep:**
```
frontend/src/pages/
├── BacklogPage.tsx ✅ (Keep)
├── SprintBoardPage.tsx ✅ (Keep)
├── ProjectSettingsPage.tsx ✅ (Keep)
├── ProjectHomePage.tsx ✅ (Keep)
├── ProjectCalendarPage.tsx ✅ (Keep)
├── ComponentsPage.tsx ✅ (Keep)
├── IssuesPage.tsx ✅ (Keep)
├── ReleasesPage.tsx ✅ (Keep)
├── ReportsPage.tsx ✅ (Keep)
└── RepositoryPage.tsx ✅ (Keep)
```

### 5. Keep Existing Jira Components

**Components to Keep:**
```
frontend/src/components/project/
├── BacklogView.tsx ✅ (Keep)
├── SprintBoard.tsx ✅ (Keep)
├── RecordCard.tsx ✅ (Keep)
├── RecordDetailModal.tsx ✅ (Keep)
├── SearchFilterBar.tsx ✅ (Keep)
├── CommentSection.tsx ✅ (Keep)
├── AttachmentSection.tsx ✅ (Keep)
├── LabelManager.tsx ✅ (Keep)
├── BulkOperationsBar.tsx ✅ (Keep)
├── CalendarGrid.tsx ✅ (Keep)
├── CreateProjectDialog.tsx ✅ (Keep)
└── MemberManagement.tsx ✅ (Keep)
```

---

## 📋 Frontend Changes Summary

### Delete (3 files)
```
❌ frontend/src/components/project/ProjectBoardColumn.tsx
❌ frontend/src/components/project/ProjectRecordCard.tsx
❌ frontend/src/components/project/ProjectFilterBar.tsx
```

### Replace (1 file)
```
🔄 frontend/src/pages/ProjectBoardPage.tsx
   (Replace with Jira board view)
```

### Keep (20+ files)
```
✅ All Jira pages
✅ All Jira components
✅ All other pages
✅ All other components
```

### Update (1 file)
```
⚠️ frontend/src/App.tsx
   (Verify routes are correct - should already be correct)
```

---

## 🔄 Execution Steps

### Step 1: Delete Old Components
```bash
rm frontend/src/components/project/ProjectBoardColumn.tsx
rm frontend/src/components/project/ProjectRecordCard.tsx
rm frontend/src/components/project/ProjectFilterBar.tsx
```

### Step 2: Replace ProjectBoardPage
```bash
# Backup old file
cp frontend/src/pages/ProjectBoardPage.tsx frontend/src/pages/ProjectBoardPage.tsx.backup

# Replace with new Jira board view
# (See new implementation above)
```

### Step 3: Verify Routes
```bash
# Check frontend/src/App.tsx
# Routes should already be correct
```

### Step 4: Build Frontend
```bash
cd frontend
npm run build
```

### Step 5: Test
```bash
# Start dev server
npm run dev

# Test all pages:
# - /projects/{id} → Jira board
# - /projects/{id}/sprint → Sprint board
# - /projects/{id}/backlog → Backlog
# - /projects/{id}/settings → Settings
```

---

## ✅ Checklist

### Phase 3: Backend
- [x] Verify all Jira handlers implemented
- [x] Verify all Jira repositories implemented
- [x] Verify all Jira usecases implemented
- [x] Verify router has all Jira routes
- [x] No backend changes needed

### Phase 4: Frontend
- [ ] Delete ProjectBoardColumn.tsx
- [ ] Delete ProjectRecordCard.tsx
- [ ] Delete ProjectFilterBar.tsx
- [ ] Replace ProjectBoardPage.tsx
- [ ] Verify App.tsx routes
- [ ] Build frontend
- [ ] Test all pages

---

## 🎯 Next Steps

1. ✅ Phase 1 Complete - Analisis selesai
2. ⏳ Phase 2 - Migrasi Data (Ready)
3. ✅ Phase 3 Complete - Backend sudah Jira-ready
4. ⏳ Phase 4 - Update Frontend (Ready)
5. ⏳ Phase 5 - Testing
6. ⏳ Phase 6 - Deployment

---

## 📝 Summary

**Backend:** ✅ No changes needed - already Jira-compatible!

**Frontend:** 
- Delete 3 old components
- Replace ProjectBoardPage with Jira board view
- Keep all Jira pages and components
- Verify routes

**Result:** Project board will be 100% Jira-only!

---

**Ready untuk Phase 4 execution?** 🚀
