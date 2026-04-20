# Rencana Migrasi ke Jira-Only Project Board

**Status**: PLANNING  
**Date**: April 19, 2026

---

## 📋 Analisis Situasi Saat Ini

Project board sekarang adalah **campuran dari 4 aplikasi**:

1. **ITSM Web App** - Ticketing system
2. **Kanban Board** - Basic kanban
3. **Project Board** - Project management (columns, records, drag-drop)
4. **Jira-like Project Board** - Advanced features (issue types, custom fields, workflows, sprints, backlog)

**Masalah**: Fitur tercampur, UI tidak konsisten, user bingung

**Solusi**: Hapus fitur project board lama, gunakan Jira sepenuhnya

---

## 🎯 Rencana Migrasi

### Phase 1: Analisis dan Persiapan (1-2 jam)

**Apa yang perlu dihapus:**
- [ ] Fitur project board lama (columns, records, drag-drop)
- [ ] Halaman ProjectBoardPage (ganti dengan Jira board)
- [ ] Komponen project board lama
- [ ] Database schema project board lama

**Apa yang perlu dipertahankan:**
- [ ] Fitur Jira (issue types, custom fields, workflows, sprints, backlog)
- [ ] Fitur Jira yang sudah ada
- [ ] Database schema Jira

**Apa yang perlu ditambahkan:**
- [ ] Migrasi data dari project board lama ke Jira
- [ ] Update UI untuk Jira only
- [ ] Update routes untuk Jira only

### Phase 2: Migrasi Data (1-2 jam)

**Data yang perlu dimigrasikan:**
- [ ] Projects → Projects (tetap sama)
- [ ] Columns → Workflow statuses
- [ ] Records → Issues dengan issue type Task
- [ ] Activity logs → Activity logs (tetap sama)

**Script migrasi:**
```sql
-- Migrate projects (no change needed)
-- Migrate columns to workflow statuses
-- Migrate records to issues
-- Create default workflows
-- Create default issue type schemes
```

### Phase 3: Update Backend (2-3 jam)

**Apa yang perlu diubah:**
- [ ] Hapus endpoint project board lama
- [ ] Tambahkan endpoint Jira yang belum ada
- [ ] Update database schema
- [ ] Update repository dan usecase

**Endpoint yang perlu ditambahkan:**
- [ ] Sprint management endpoints
- [ ] Backlog endpoints
- [ ] Custom fields endpoints
- [ ] Workflows endpoints
- [ ] Comments endpoints
- [ ] Attachments endpoints
- [ ] Labels endpoints

### Phase 4: Update Frontend (2-3 jam)

**Apa yang perlu diubah:**
- [ ] Hapus ProjectBoardPage lama
- [ ] Ganti dengan Jira board view
- [ ] Update sidebar navigation
- [ ] Update routes
- [ ] Update components

**Halaman yang perlu diubah:**
- [ ] ProjectBoardPage → Jira board view
- [ ] BacklogPage → Jira backlog view
- [ ] SprintBoardPage → Jira sprint board view
- [ ] ProjectSettingsPage → Jira settings view

### Phase 5: Testing dan Deployment (1-2 jam)

**Testing:**
- [ ] Test semua Jira features
- [ ] Test migrasi data
- [ ] Test backward compatibility
- [ ] Test UI/UX

**Deployment:**
- [ ] Build frontend
- [ ] Build backend
- [ ] Deploy ke production
- [ ] Monitor

---

## 📊 Perbandingan

### Project Board Lama vs Jira

| Fitur | Project Board Lama | Jira | Status |
|-------|-------------------|------|--------|
| Columns | ✅ | ❌ (Workflow statuses) | Replace |
| Records | ✅ | ✅ (Issues) | Migrate |
| Drag-drop | ✅ | ✅ | Keep |
| Issue Types | ❌ | ✅ | Add |
| Custom Fields | ❌ | ✅ | Add |
| Workflows | ❌ | ✅ | Add |
| Sprints | ❌ | ✅ | Add |
| Backlog | ❌ | ✅ | Add |
| Comments | ❌ | ✅ | Add |
| Attachments | ❌ | ✅ | Add |
| Labels | ❌ | ✅ | Add |

---

## 🗂️ File-File yang Perlu Diubah

### Backend

**Hapus:**
- [ ] `backend/internal/delivery/http/project_board_handler.go`
- [ ] `backend/internal/usecase/project_board_usecase.go`
- [ ] `backend/internal/repository/project_board_repository.go`

**Tambahkan/Update:**
- [ ] `backend/internal/delivery/http/jira_handler.go` (update)
- [ ] `backend/internal/usecase/jira_usecase.go` (update)
- [ ] `backend/internal/repository/jira_repository.go` (update)
- [ ] `backend/internal/delivery/http/sprint_handler.go` (new)
- [ ] `backend/internal/delivery/http/backlog_handler.go` (new)
- [ ] `backend/internal/delivery/http/custom_field_handler.go` (new)
- [ ] `backend/internal/delivery/http/workflow_handler.go` (new)
- [ ] `backend/internal/delivery/http/comment_handler.go` (new)
- [ ] `backend/internal/delivery/http/attachment_handler.go` (new)
- [ ] `backend/internal/delivery/http/label_handler.go` (new)

### Frontend

**Hapus:**
- [ ] `frontend/src/pages/ProjectBoardPage.tsx` (old version)
- [ ] `frontend/src/components/project/ProjectBoardColumn.tsx`
- [ ] `frontend/src/components/project/ProjectRecordCard.tsx`
- [ ] `frontend/src/components/project/ProjectFilterBar.tsx`

**Tambahkan/Update:**
- [ ] `frontend/src/pages/ProjectBoardPage.tsx` (Jira board view)
- [ ] `frontend/src/pages/BacklogPage.tsx` (update)
- [ ] `frontend/src/pages/SprintBoardPage.tsx` (update)
- [ ] `frontend/src/pages/ProjectSettingsPage.tsx` (update)
- [ ] `frontend/src/components/jira/JiraBoard.tsx` (new)
- [ ] `frontend/src/components/jira/JiraBacklog.tsx` (new)
- [ ] `frontend/src/components/jira/JiraSprint.tsx` (new)
- [ ] `frontend/src/components/jira/JiraSettings.tsx` (new)

### Database

**Migrasi:**
- [ ] Migrate projects
- [ ] Migrate columns to workflow statuses
- [ ] Migrate records to issues
- [ ] Create default workflows
- [ ] Create default issue type schemes

---

## 🔄 Migrasi Data

### Projects
```
projects → projects (no change)
```

### Columns → Workflow Statuses
```
project_columns → workflow_statuses
- id → id
- project_id → workflow_id
- name → status_name
- position → status_order
```

### Records → Issues
```
project_records → issues (new table)
- id → id
- column_id → status_id (from workflow_statuses)
- project_id → project_id
- title → title
- description → description
- assigned_to → assigned_to
- due_date → due_date
- position → priority
- is_completed → is_completed
- completed_at → completed_at
- created_by → created_by
- created_at → created_at
- updated_at → updated_at
- (new) issue_type_id → Task (default)
- (new) status → (from column name)
- (new) parent_record_id → null
```

---

## ✅ Checklist Migrasi

### Phase 1: Analisis
- [ ] Identifikasi semua fitur project board lama
- [ ] Identifikasi semua fitur Jira yang sudah ada
- [ ] Identifikasi data yang perlu dimigrasikan
- [ ] Buat backup database

### Phase 2: Migrasi Data
- [ ] Buat script migrasi
- [ ] Test script migrasi di development
- [ ] Backup data production
- [ ] Jalankan migrasi di production

### Phase 3: Backend
- [ ] Hapus endpoint project board lama
- [ ] Tambahkan endpoint Jira yang belum ada
- [ ] Update database schema
- [ ] Test semua endpoint
- [ ] Build backend

### Phase 4: Frontend
- [ ] Hapus komponen project board lama
- [ ] Update ProjectBoardPage ke Jira board
- [ ] Update BacklogPage
- [ ] Update SprintBoardPage
- [ ] Update ProjectSettingsPage
- [ ] Update routes
- [ ] Test semua halaman
- [ ] Build frontend

### Phase 5: Testing
- [ ] Test migrasi data
- [ ] Test semua Jira features
- [ ] Test UI/UX
- [ ] Test backward compatibility
- [ ] Test performance

### Phase 6: Deployment
- [ ] Deploy backend
- [ ] Deploy frontend
- [ ] Monitor production
- [ ] Gather user feedback

---

## 📈 Timeline

| Phase | Durasi | Status |
|-------|--------|--------|
| Phase 1: Analisis | 1-2 jam | ⏳ Pending |
| Phase 2: Migrasi Data | 1-2 jam | ⏳ Pending |
| Phase 3: Backend | 2-3 jam | ⏳ Pending |
| Phase 4: Frontend | 2-3 jam | ⏳ Pending |
| Phase 5: Testing | 1-2 jam | ⏳ Pending |
| Phase 6: Deployment | 1 jam | ⏳ Pending |
| **Total** | **8-13 jam** | ⏳ Pending |

---

## 🎯 Hasil Akhir

Setelah migrasi selesai:

✅ **Project board hanya menggunakan Jira features**
✅ **Tidak ada lagi fitur project board lama**
✅ **UI konsisten dan clean**
✅ **Semua data termigrasi dengan baik**
✅ **Backward compatible dengan data lama**
✅ **Production ready**

---

## ⚠️ Risiko dan Mitigasi

| Risiko | Mitigasi |
|--------|----------|
| Data loss | Backup database sebelum migrasi |
| User confusion | Dokumentasi dan training |
| Performance issue | Load testing sebelum deployment |
| Backward compatibility | Test dengan data lama |

---

## 📞 Next Steps

1. **Approve rencana migrasi** - Apakah Anda setuju dengan rencana ini?
2. **Mulai Phase 1** - Analisis dan persiapan
3. **Jalankan migrasi** - Phase 2-6
4. **Deploy ke production** - Phase 6
5. **Monitor dan gather feedback** - Post-deployment

---

**Apakah Anda ingin saya mulai dengan Phase 1?** 🚀
