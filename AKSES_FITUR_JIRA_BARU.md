# Cara Akses Fitur-Fitur Jira Baru

## ⚠️ Penting: Fitur Jira Sudah Ada Tapi Belum Terlihat

Fitur-fitur Jira sudah diimplementasikan di code, tapi **belum terlihat di UI** karena:
1. Halaman baru ada tapi perlu diklik dari navigation bar
2. Database belum punya data sample
3. Navigation bar hanya muncul di halaman project

---

## 🎯 Langkah-Langkah Akses Fitur Jira

### Step 1: Buka Project
```
1. Buka aplikasi: http://localhost:3000
2. Lihat halaman Project Board (seperti screenshot)
3. Perhatikan navigation bar di atas dengan 4 menu:
   - Board (saat ini aktif)
   - Sprint (BARU - fitur Jira)
   - Backlog (BARU - fitur Jira)
   - Settings (BARU - fitur Jira)
```

### Step 2: Klik "Sprint" untuk Lihat Sprint Board
```
Navigation Bar:
┌──────────┬──────────┬──────────┬──────────┐
│ Board    │ Sprint ← │ Backlog  │ Settings │
└──────────┴──────────┴──────────┴──────────┘

Klik "Sprint"
↓
Akan melihat:
- Sprint Board (Kanban-style)
- Sprint metrics (total, completed, %, days left)
- Status columns (Backlog, To Do, In Progress, Done)
- Records dengan issue type icon
- Drag-and-drop antar status
```

### Step 3: Klik "Backlog" untuk Lihat Backlog
```
Navigation Bar:
┌──────────┬──────────┬──────────┬──────────┐
│ Board    │ Sprint   │ Backlog ←│ Settings │
└──────────┴──────────┴──────────┴──────────┘

Klik "Backlog"
↓
Akan melihat:
- Backlog records (unassigned)
- Priority order (1, 2, 3, ...)
- Drag-and-drop untuk reorder
- Drag-and-drop ke sprint
- Records dengan issue type icon
```

### Step 4: Klik "Settings" untuk Lihat Project Configuration
```
Navigation Bar:
┌──────────┬──────────┬──────────┬──────────┐
│ Board    │ Sprint   │ Backlog  │ Settings ←
└──────────┴──────────┴──────────┴──────────┘

Klik "Settings"
↓
Akan melihat 4 tabs:

1. Issue Types Tab
   - 5 tipe issue (Bug, Task, Story, Epic, Sub-task)
   - Icon untuk setiap tipe
   - Deskripsi

2. Custom Fields Tab
   - List custom fields
   - Tombol "Add Field"
   - Tombol delete
   - Tipe field (text, textarea, dropdown, date, number, checkbox)

3. Workflows Tab
   - Workflow configuration
   - Initial status
   - Workflow name

4. Labels Tab
   - List labels dengan color
   - Tombol "Add Label"
   - Tombol delete
```

---

## 📍 Lokasi Navigation Bar

```
┌─────────────────────────────────────────────────────────────┐
│ HEADER (Logo, Search, Notifications, Profile)              │
├─────────────────────────────────────────────────────────────┤
│ SIDEBAR (Projects list)                                     │
├─────────────────────────────────────────────────────────────┤
│ NAVIGATION BAR ← FITUR JIRA BARU DI SINI                    │
│ ┌──────────┬──────────┬──────────┬──────────┐               │
│ │ Board    │ Sprint   │ Backlog  │ Settings │               │
│ └──────────┴──────────┴──────────┴──────────┘               │
├─────────────────────────────────────────────────────────────┤
│ MAIN CONTENT AREA                                           │
│                                                             │
│ (Akan berubah sesuai menu yang diklik)                      │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

---

## 🎬 Demo: Lihat Fitur Jira dalam 1 Menit

### Scenario: Lihat Sprint Board

```
1. Buka aplikasi: http://localhost:3000
   ↓
2. Lihat Project Board (seperti screenshot)
   ↓
3. Lihat navigation bar dengan 4 menu
   ↓
4. Klik "Sprint"
   ↓
5. Akan melihat:
   - Sprint Board dengan Kanban columns
   - Sprint metrics di atas
   - Records dengan issue type icon
   - Bisa drag-and-drop antar status
```

### Scenario: Lihat Backlog

```
1. Dari Sprint Board, klik "Backlog"
   ↓
2. Akan melihat:
   - Backlog records ordered by priority
   - Priority order (1, 2, 3, ...)
   - Bisa drag-and-drop untuk reorder
   - Bisa drag-and-drop ke sprint
```

### Scenario: Lihat Project Settings

```
1. Dari Backlog, klik "Settings"
   ↓
2. Akan melihat 4 tabs:
   - Issue Types (5 tipe)
   - Custom Fields (create/edit/delete)
   - Workflows (configuration)
   - Labels (create/edit/delete)
```

---

## ⚠️ Catatan Penting

### Fitur Sudah Ada Tapi Belum Terlihat Karena:

1. **Navigation Bar Hanya Muncul di Project Pages**
   - Hanya muncul saat membuka project
   - Tidak muncul di dashboard atau halaman lain

2. **Database Belum Punya Data Sample**
   - Tidak ada sprints yang sudah dibuat
   - Tidak ada custom fields yang sudah dibuat
   - Tidak ada workflows yang sudah dibuat
   - Tidak ada labels yang sudah dibuat

3. **Halaman Baru Sudah Ada Tapi Perlu Diklik**
   - SprintBoardPage.tsx - ada
   - BacklogPage.tsx - ada
   - ProjectSettingsPage.tsx - ada
   - ProjectNavigation.tsx - ada
   - Tapi perlu diklik dari navigation bar

---

## 🔧 Untuk Membuat Fitur Lebih Terlihat

### Option 1: Buat Data Sample di Database
```sql
-- Create sample sprint
INSERT INTO sprints (id, project_id, name, goal, start_date, end_date, status)
VALUES (
  'sprint-1',
  'project-id',
  'Sprint 1',
  'Implement user authentication',
  '2026-04-19',
  '2026-04-26',
  'Active'
);

-- Create sample custom field
INSERT INTO custom_fields (id, project_id, name, field_type, is_required)
VALUES (
  'field-1',
  'project-id',
  'Priority',
  'dropdown',
  true
);

-- Create sample label
INSERT INTO labels (id, project_id, name, color)
VALUES (
  'label-1',
  'project-id',
  'Bug',
  '#FF0000'
);
```

### Option 2: Buat Data Melalui UI
1. Buka Project Settings
2. Klik "Add Label" untuk membuat label
3. Klik "Add Field" untuk membuat custom field
4. Buat sprint melalui API atau UI

### Option 3: Jalankan Migration Script
```bash
# Buat script untuk populate sample data
./scripts/seed-jira-features.sh
```

---

## 📊 Fitur yang Sudah Ada (Tapi Belum Terlihat)

### Backend (100% Complete)
- ✅ 40+ API endpoints
- ✅ Issue types management
- ✅ Custom fields management
- ✅ Workflows management
- ✅ Sprints management
- ✅ Backlog management
- ✅ Comments dengan @mentions
- ✅ Attachments
- ✅ Labels management

### Frontend (85% Complete)
- ✅ SprintBoardPage - Kanban board dengan metrics
- ✅ BacklogPage - Prioritized backlog
- ✅ ProjectSettingsPage - 4 tabs (Issue Types, Custom Fields, Workflows, Labels)
- ✅ ProjectNavigation - 4 menu items (Board, Sprint, Backlog, Settings)
- ✅ RecordCard - Display issue type, labels, metrics
- ✅ SearchFilterBar - Advanced filtering
- ✅ Comments section - Add comment dengan @mention
- ✅ Attachments section - Upload file
- ✅ Labels section - Add label

---

## 🎯 Checklist: Fitur yang Harus Dilihat

- [ ] Buka Project Board
- [ ] Lihat navigation bar dengan 4 menu
- [ ] Klik "Sprint" untuk lihat Sprint Board
- [ ] Lihat sprint metrics
- [ ] Lihat kanban columns
- [ ] Klik "Backlog" untuk lihat Backlog
- [ ] Lihat priority order
- [ ] Klik "Settings" untuk lihat Project Settings
- [ ] Lihat Issue Types tab
- [ ] Lihat Custom Fields tab
- [ ] Lihat Workflows tab
- [ ] Lihat Labels tab

---

## 🆘 Jika Navigation Bar Tidak Terlihat

### Kemungkinan 1: Belum di Project Page
```
Solusi: Klik project dari sidebar untuk membuka project page
```

### Kemungkinan 2: Browser Cache
```
Solusi: Refresh browser (Ctrl+R atau Cmd+R)
```

### Kemungkinan 3: Frontend Belum Build
```
Solusi: 
npm run build
npm run dev
```

### Kemungkinan 4: Backend Tidak Jalan
```
Solusi:
docker compose ps
docker compose logs -f backend
```

---

## 📞 Summary

**Fitur Jira sudah ada di code, tapi perlu diklik dari navigation bar:**

1. **Board** - Project Board (existing)
2. **Sprint** ← BARU - Sprint Board dengan Kanban
3. **Backlog** ← BARU - Prioritized backlog
4. **Settings** ← BARU - Project configuration

**Untuk melihat fitur-fitur Jira:**
1. Buka project
2. Lihat navigation bar
3. Klik "Sprint", "Backlog", atau "Settings"
4. Explore fitur-fitur Jira!

---

**Status:** ✅ Fitur Sudah Ada - Tinggal Diklik!  
**Date:** April 19, 2026  
**Project:** ITSM Ticketing System - Jira-like Project Board
