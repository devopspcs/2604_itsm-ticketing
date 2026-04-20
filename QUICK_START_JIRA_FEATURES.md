# Quick Start - Lihat Fitur Jira-like dalam 5 Menit

## ⚡ 5 Langkah Cepat

### Langkah 1: Deploy (2 menit)
```bash
# Build dan start services
docker compose build
docker compose up -d

# Tunggu 30 detik untuk services siap
sleep 30

# Verify
docker compose ps
```

### Langkah 2: Buka Browser (30 detik)
```
http://localhost:3000
```

### Langkah 3: Login
- Username: admin
- Password: admin

### Langkah 4: Pilih Project
- Klik project yang ada atau buat project baru

### Langkah 5: Explore Fitur (2 menit)

---

## 🎯 Fitur yang Bisa Langsung Dilihat

### 1️⃣ Sprint Board (Kanban-style)
```
URL: http://localhost:3000/projects/{projectId}/sprint

Klik "Sprint" di navigation bar
↓
Lihat:
- Sprint metrics (total, completed, %, days left)
- Kanban board dengan status columns
- Records dengan issue type icon
- Drag-and-drop antar status
```

### 2️⃣ Backlog (Prioritized List)
```
URL: http://localhost:3000/projects/{projectId}/backlog

Klik "Backlog" di navigation bar
↓
Lihat:
- Unassigned records ordered by priority
- Priority order (1, 2, 3, ...)
- Drag-and-drop untuk reorder
- Drag-and-drop ke sprint
```

### 3️⃣ Project Settings
```
URL: http://localhost:3000/projects/{projectId}/settings

Klik "Settings" di navigation bar
↓
Lihat 4 tabs:
1. Issue Types - 5 tipe (Bug, Task, Story, Epic, Sub-task)
2. Custom Fields - Create/edit/delete custom fields
3. Workflows - Workflow configuration
4. Labels - Create/edit/delete labels
```

### 4️⃣ Record Detail
```
Klik record di Sprint Board atau Backlog
↓
Lihat:
- Issue type
- Custom fields
- Comments section (add comment, @mention)
- Attachments section (upload file)
- Labels section (add label)
- Activity log
```

---

## 🎬 Demo Workflow (3 menit)

### Scenario: Create dan Manage Sprint

**Step 1: Create Sprint**
```
1. Buka Project Settings
2. Klik "Sprints" (atau di main menu)
3. Klik "Create Sprint"
4. Isi: Name, Start Date, End Date, Goal
5. Klik "Create"
```

**Step 2: Start Sprint**
```
1. Lihat sprint di list
2. Klik "Start Sprint"
3. Sprint status berubah ke "Active"
4. Sprint Board menjadi available
```

**Step 3: Assign Records to Sprint**
```
1. Buka Backlog
2. Drag record ke sprint
3. Record muncul di Sprint Board
```

**Step 4: Work on Sprint**
```
1. Buka Sprint Board
2. Drag record antar status columns
3. Lihat metrics update
4. Klik record untuk add comment/attachment
```

**Step 5: Complete Sprint**
```
1. Klik "Complete Sprint"
2. Sprint metrics calculated
3. Incomplete records moved to backlog
```

---

## 📸 Visual Guide

### Sprint Board View
```
┌─────────────────────────────────────────────────────────────┐
│ Sprint: Sprint 1 (Apr 19 - Apr 26)                          │
│ Goal: Implement user authentication                         │
│                                                             │
│ Metrics: Total: 10 | Completed: 3 | Progress: 30% | Days: 7│
│ ████░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░  │
├─────────────────────────────────────────────────────────────┤
│ Backlog    │ To Do      │ In Progress │ Done               │
├────────────┼────────────┼─────────────┼────────────────────┤
│ [Card]     │ [Card]     │ [Card]      │ [Card]             │
│ [Card]     │ [Card]     │             │ [Card]             │
│            │ [Card]     │ [Card]      │ [Card]             │
└────────────┴────────────┴─────────────┴────────────────────┘
```

### Backlog View
```
┌─────────────────────────────────────────────────────────────┐
│ Backlog                                                     │
├─────────────────────────────────────────────────────────────┤
│ Priority │ Record                                           │
│ ────────────────────────────────────────────────────────── │
│ 1        │ [🐛 Bug] Fix login error                        │
│ 2        │ [✓ Task] Update documentation                   │
│ 3        │ [📖 Story] Add user profile page                │
│ 4        │ [🎯 Epic] Implement payment system              │
│ 5        │ [↳ Sub-task] Design database schema             │
└─────────────────────────────────────────────────────────────┘
```

### Record Card
```
┌─────────────────────────────────────────┐
│ 🐛 [BUG] Fix login error                │
│                                         │
│ Description: Users can't login...       │
│                                         │
│ 👤 John Doe                             │
│ 📅 Due: 2026-04-25                      │
│ 📎 2 attachments                        │
│ 💬 3 comments                           │
│                                         │
│ [🏷️ Bug] [🏷️ Critical]                  │
└─────────────────────────────────────────┘
```

### Project Settings
```
┌─────────────────────────────────────────────────────────────┐
│ Project Settings                                            │
├─────────────────────────────────────────────────────────────┤
│ [Issue Types] [Custom Fields] [Workflows] [Labels]          │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│ Issue Types:                                                │
│ ✓ Bug                                                       │
│ ✓ Task                                                      │
│ ✓ Story                                                     │
│ ✓ Epic                                                      │
│ ✓ Sub-task                                                  │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

---

## 🔗 Direct Links

### Frontend Pages
- **Sprint Board:** http://localhost:3000/projects/{projectId}/sprint
- **Backlog:** http://localhost:3000/projects/{projectId}/backlog
- **Settings:** http://localhost:3000/projects/{projectId}/settings
- **Project Board:** http://localhost:3000/projects/{projectId}

### API Endpoints
- **Health Check:** http://localhost:8080/health
- **Issue Types:** http://localhost:8080/api/v1/projects/{projectId}/issue-types
- **Custom Fields:** http://localhost:8080/api/v1/projects/{projectId}/custom-fields
- **Sprints:** http://localhost:8080/api/v1/projects/{projectId}/sprints
- **Backlog:** http://localhost:8080/api/v1/projects/{projectId}/backlog

---

## ✅ Checklist: Fitur yang Harus Dilihat

- [ ] **Sprint Board** - Kanban view dengan metrics
- [ ] **Backlog** - Prioritized list dengan drag-and-drop
- [ ] **Issue Types** - 5 tipe (Bug, Task, Story, Epic, Sub-task)
- [ ] **Custom Fields** - Create/edit/delete fields
- [ ] **Workflows** - Status configuration
- [ ] **Labels** - Create/edit/delete labels
- [ ] **Comments** - Add comment dengan @mention
- [ ] **Attachments** - Upload file
- [ ] **Record Card** - Display issue type, labels, metrics
- [ ] **Filtering** - Filter by issue type, assignee, label
- [ ] **Drag-and-drop** - Move records antar status/sprint

---

## 🆘 Jika Ada Masalah

### Services tidak jalan
```bash
docker compose ps
docker compose logs -f
```

### Port sudah dipakai
```bash
lsof -i :3000
kill -9 <PID>
```

### Database error
```bash
docker compose exec postgres psql -U itsm -d itsm -c "\dt"
```

### Restart semua
```bash
docker compose down
docker compose up -d
```

---

## 📊 Fitur yang Sudah Ada

✅ Issue Types (5 tipe)  
✅ Custom Fields (7 tipe)  
✅ Workflows (Custom statuses)  
✅ Sprint Planning  
✅ Backlog Management  
✅ Comments dengan @Mentions  
✅ Attachments  
✅ Labels  
✅ Sprint Board (Kanban)  
✅ Backlog View  
✅ Filtering & Search  
✅ Bulk Operations  
✅ Activity Logging  
✅ Notifications  

---

**Status:** ✅ Ready to View  
**Time to Deploy:** 2 minutes  
**Time to Explore:** 3 minutes  
**Total:** 5 minutes  

**Sekarang buka browser dan lihat fitur-fitur Jira-like yang sudah diimplementasikan!** 🚀
