# Cara Melihat Fitur-Fitur Jira-like

## 🚀 Langkah 1: Deploy Aplikasi

### Option A: Docker Deployment (Recommended)

```bash
# 1. Build Docker images
docker compose build

# 2. Start services
docker compose up -d

# 3. Wait for services to be ready (30 seconds)
sleep 30

# 4. Verify services
docker compose ps
```

### Option B: Manual Deployment

```bash
# Terminal 1: Start PostgreSQL
docker run --name itsm-postgres \
  -e POSTGRES_USER=itsm \
  -e POSTGRES_PASSWORD=itsm \
  -e POSTGRES_DB=itsm \
  -p 5432:5432 \
  postgres:16-alpine

# Terminal 2: Start Backend
cd backend
go run ./cmd/server

# Terminal 3: Start Frontend
cd frontend
npm run dev
```

---

## 🌐 Langkah 2: Akses Aplikasi

### Frontend URL
```
http://localhost:3000
```

### Backend API
```
http://localhost:8080/api/v1
```

### Database
```
Host: localhost
Port: 5432
User: itsm
Password: itsm
Database: itsm
```

---

## 📍 Dimana Melihat Fitur-Fitur Jira

### 1. 🎯 Issue Types
**Lokasi:** Project Settings → Issue Types Tab

```
URL: http://localhost:3000/projects/{projectId}/settings
```

**Yang Bisa Dilihat:**
- ✅ 5 tipe issue (Bug, Task, Story, Epic, Sub-task)
- ✅ Icon untuk setiap tipe
- ✅ Deskripsi tipe issue

**Cara Akses:**
1. Buka aplikasi di http://localhost:3000
2. Pilih project
3. Klik "Settings" di navigation bar
4. Lihat tab "Issue Types"

---

### 2. 📝 Custom Fields
**Lokasi:** Project Settings → Custom Fields Tab

```
URL: http://localhost:3000/projects/{projectId}/settings
```

**Yang Bisa Dilihat:**
- ✅ List custom fields yang sudah dibuat
- ✅ Tipe field (text, textarea, dropdown, date, number, checkbox)
- ✅ Tombol "Add Field" untuk membuat field baru
- ✅ Tombol delete untuk menghapus field

**Cara Akses:**
1. Buka aplikasi di http://localhost:3000
2. Pilih project
3. Klik "Settings" di navigation bar
4. Klik tab "Custom Fields"
5. Klik "Add Field" untuk membuat field baru

**Contoh Custom Fields:**
- Priority (dropdown)
- Estimated Hours (number)
- Due Date (date)
- Description (textarea)

---

### 3. ⚙️ Workflows
**Lokasi:** Project Settings → Workflows Tab

```
URL: http://localhost:3000/projects/{projectId}/settings
```

**Yang Bisa Dilihat:**
- ✅ Workflow configuration
- ✅ Initial status
- ✅ Workflow name

**Cara Akses:**
1. Buka aplikasi di http://localhost:3000
2. Pilih project
3. Klik "Settings" di navigation bar
4. Klik tab "Workflows"

---

### 4. 🏃 Sprint Board
**Lokasi:** Project → Sprint Tab

```
URL: http://localhost:3000/projects/{projectId}/sprint
```

**Yang Bisa Dilihat:**
- ✅ Active sprint dengan nama, dates, goal
- ✅ Sprint metrics (total, completed, %, days left)
- ✅ Kanban board dengan status columns
- ✅ Records organized by status
- ✅ Drag-and-drop antar status
- ✅ Issue type icon pada setiap card
- ✅ Assignee avatar
- ✅ Labels sebagai badges

**Cara Akses:**
1. Buka aplikasi di http://localhost:3000
2. Pilih project
3. Klik "Sprint" di navigation bar
4. Lihat sprint board dengan metrics

**Fitur yang Bisa Dicoba:**
- Drag record dari satu status ke status lain
- Lihat metrics update secara real-time
- Klik record untuk melihat detail

---

### 5. 📦 Backlog
**Lokasi:** Project → Backlog Tab

```
URL: http://localhost:3000/projects/{projectId}/backlog
```

**Yang Bisa Dilihat:**
- ✅ Unassigned records ordered by priority
- ✅ Priority order (1, 2, 3, ...)
- ✅ Drag-and-drop untuk reorder
- ✅ Drag-and-drop ke sprint
- ✅ Issue type icon pada setiap card
- ✅ Assignee avatar
- ✅ Labels sebagai badges

**Cara Akses:**
1. Buka aplikasi di http://localhost:3000
2. Pilih project
3. Klik "Backlog" di navigation bar
4. Lihat backlog dengan prioritization

**Fitur yang Bisa Dicoba:**
- Drag record untuk mengubah priority
- Drag record ke sprint untuk assign
- Lihat priority order berubah

---

### 6. 💬 Comments dengan @Mentions
**Lokasi:** Record Detail → Comments Section

```
URL: http://localhost:3000/projects/{projectId}/records/{recordId}
```

**Yang Bisa Dilihat:**
- ✅ Comments section di record detail
- ✅ Tombol "Add Comment"
- ✅ Comment text dengan author, timestamp
- ✅ @mention dropdown saat mengetik "@"
- ✅ Highlighted @mentions

**Cara Akses:**
1. Buka aplikasi di http://localhost:3000
2. Pilih project
3. Klik record untuk membuka detail
4. Scroll ke "Comments" section
5. Klik "Add Comment"
6. Ketik "@" untuk melihat dropdown mention

**Fitur yang Bisa Dicoba:**
- Ketik comment
- Ketik "@" untuk mention team member
- Submit comment
- Lihat comment dengan author dan timestamp

---

### 7. 📎 Attachments
**Lokasi:** Record Detail → Attachments Section

```
URL: http://localhost:3000/projects/{projectId}/records/{recordId}
```

**Yang Bisa Dilihat:**
- ✅ Attachments section di record detail
- ✅ Tombol "Add Attachment"
- ✅ List uploaded files
- ✅ File name, size, uploader, timestamp
- ✅ Delete button untuk setiap file

**Cara Akses:**
1. Buka aplikasi di http://localhost:3000
2. Pilih project
3. Klik record untuk membuka detail
4. Scroll ke "Attachments" section
5. Klik "Add Attachment"
6. Pilih file untuk upload

**Fitur yang Bisa Dicoba:**
- Upload file (image, document, dll)
- Lihat file muncul di list
- Delete file

---

### 8. 🏷️ Labels
**Lokasi:** Project Settings → Labels Tab

```
URL: http://localhost:3000/projects/{projectId}/settings
```

**Yang Bisa Dilihat:**
- ✅ Labels list dengan color
- ✅ Tombol "Add Label"
- ✅ Label name dan color
- ✅ Delete button untuk setiap label

**Cara Akses:**
1. Buka aplikasi di http://localhost:3000
2. Pilih project
3. Klik "Settings" di navigation bar
4. Klik tab "Labels"
5. Klik "Add Label" untuk membuat label baru

**Fitur yang Bisa Dicoba:**
- Buat label dengan name dan color
- Lihat label muncul di list
- Delete label

---

### 9. 🎨 Record Card Display
**Lokasi:** Sprint Board atau Backlog

```
URL: http://localhost:3000/projects/{projectId}/sprint
atau
http://localhost:3000/projects/{projectId}/backlog
```

**Yang Bisa Dilihat pada Card:**
- ✅ Issue type icon (🐛 Bug, ✓ Task, 📖 Story, 🎯 Epic, ↳ Sub-task)
- ✅ Record title
- ✅ Assignee avatar
- ✅ Labels sebagai colored badges
- ✅ Due date (jika ada)
- ✅ Attachment count
- ✅ Comment count

**Contoh Card:**
```
┌─────────────────────────────────────┐
│ 🐛 [BUG] Fix login error            │
│                                     │
│ 👤 John Doe                         │
│ 📅 Due: 2026-04-25                  │
│ 📎 2 attachments                    │
│ 💬 3 comments                       │
│                                     │
│ [🏷️ Bug] [🏷️ Critical]              │
└─────────────────────────────────────┘
```

---

### 10. 🔍 Filtering & Search
**Lokasi:** Sprint Board atau Backlog

```
URL: http://localhost:3000/projects/{projectId}/sprint
atau
http://localhost:3000/projects/{projectId}/backlog
```

**Yang Bisa Dilihat:**
- ✅ Search bar untuk mencari records
- ✅ Filter button untuk advanced filters
- ✅ Filter by issue type
- ✅ Filter by status
- ✅ Filter by label
- ✅ Filter by assignee
- ✅ Filter by due date

**Cara Akses:**
1. Buka Sprint Board atau Backlog
2. Lihat search bar di atas
3. Klik "Filter" button
4. Pilih filter criteria

---

## 📊 API Endpoints yang Bisa Dicoba

### Issue Types
```bash
GET /api/v1/projects/{projectId}/issue-types
```

### Custom Fields
```bash
GET /api/v1/projects/{projectId}/custom-fields
POST /api/v1/projects/{projectId}/custom-fields
PATCH /api/v1/projects/{projectId}/custom-fields/{fieldId}
DELETE /api/v1/projects/{projectId}/custom-fields/{fieldId}
```

### Workflows
```bash
GET /api/v1/projects/{projectId}/workflow
GET /api/v1/workflows/{workflowId}/statuses
```

### Sprints
```bash
GET /api/v1/projects/{projectId}/sprints
POST /api/v1/projects/{projectId}/sprints
PATCH /api/v1/projects/{projectId}/sprints/{sprintId}
GET /api/v1/projects/{projectId}/sprints/{sprintId}/records
```

### Backlog
```bash
GET /api/v1/projects/{projectId}/backlog
PATCH /api/v1/projects/{projectId}/backlog/reorder
POST /api/v1/projects/{projectId}/backlog/assign-sprint
```

### Comments
```bash
POST /api/v1/projects/{projectId}/records/{recordId}/comments
GET /api/v1/projects/{projectId}/records/{recordId}/comments
PATCH /api/v1/projects/{projectId}/records/{recordId}/comments/{commentId}
DELETE /api/v1/projects/{projectId}/records/{recordId}/comments/{commentId}
```

### Attachments
```bash
POST /api/v1/projects/{projectId}/records/{recordId}/attachments
GET /api/v1/projects/{projectId}/records/{recordId}/attachments
DELETE /api/v1/projects/{projectId}/records/{recordId}/attachments/{attachmentId}
```

### Labels
```bash
GET /api/v1/projects/{projectId}/labels
POST /api/v1/projects/{projectId}/labels
PATCH /api/v1/projects/{projectId}/records/{recordId}/labels
```

---

## 🧪 Testing dengan Postman/cURL

### Example: Get Issue Types
```bash
curl -X GET http://localhost:8080/api/v1/projects/{projectId}/issue-types \
  -H "Authorization: Bearer {token}"
```

### Example: Create Custom Field
```bash
curl -X POST http://localhost:8080/api/v1/projects/{projectId}/custom-fields \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Priority",
    "field_type": "dropdown",
    "is_required": true
  }'
```

### Example: Add Comment
```bash
curl -X POST http://localhost:8080/api/v1/projects/{projectId}/records/{recordId}/comments \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "text": "This is a comment with @mention"
  }'
```

---

## 📋 Checklist: Fitur yang Bisa Dicoba

- [ ] Lihat Issue Types di Settings
- [ ] Buat Custom Field baru
- [ ] Lihat Workflow configuration
- [ ] Buka Sprint Board
- [ ] Lihat sprint metrics
- [ ] Drag record antar status
- [ ] Buka Backlog
- [ ] Drag record untuk reorder priority
- [ ] Drag record ke sprint
- [ ] Buka record detail
- [ ] Tambah comment dengan @mention
- [ ] Upload attachment
- [ ] Tambah label ke record
- [ ] Lihat record card dengan semua info
- [ ] Coba filter records
- [ ] Coba search records

---

## 🆘 Troubleshooting

### Aplikasi tidak bisa diakses
```bash
# Check if services are running
docker compose ps

# Check logs
docker compose logs -f frontend
docker compose logs -f backend
docker compose logs -f postgres
```

### Database connection error
```bash
# Check PostgreSQL is running
docker compose ps postgres

# Check database exists
docker compose exec postgres psql -U itsm -d itsm -c "\dt"
```

### Port already in use
```bash
# Kill process using port
lsof -i :3000
kill -9 <PID>

# Or change port in docker-compose.yml
```

---

## 📞 Support

Jika ada masalah, cek:
1. Semua services running: `docker compose ps`
2. Logs: `docker compose logs -f`
3. Database: `docker compose exec postgres psql -U itsm -d itsm`
4. API: `curl http://localhost:8080/health`

---

**Status:** ✅ Ready to View  
**Date:** April 19, 2026  
**Project:** ITSM Ticketing System - Jira-like Project Board
