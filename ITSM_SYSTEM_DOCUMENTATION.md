# PCS ITSM - System Documentation
## IT Service Management Platform

---

## 1. System Overview

PCS ITSM adalah platform IT Service Management yang dibangun untuk mengelola ticketing, approval workflow, project management, dan organisasi perusahaan dalam satu sistem terintegrasi.

### Tech Stack
| Layer | Technology |
|-------|-----------|
| Backend | Go (Golang), Chi Router, pgx/pgxpool |
| Frontend | React, TypeScript, Tailwind CSS |
| Database | PostgreSQL 16 |
| Auth | JWT (Access + Refresh Token), Keycloak SSO |
| Deployment | Docker, Docker Compose |
| Storage | Google Drive API (attachments) |

---

## 2. Database Schema (43 Tables)

### 2.1 Core Tables

#### users
| Column | Type | Description |
|--------|------|-------------|
| id | UUID (PK) | Unique identifier |
| full_name | VARCHAR | Nama lengkap |
| email | VARCHAR (UNIQUE) | Email login |
| password | VARCHAR | Bcrypt hash |
| role | VARCHAR | user, agent, approver, admin |
| is_active | BOOLEAN | Status aktif |
| department_id | UUID (FK) | Departemen |
| division_id | UUID (FK) | Divisi |
| team_id | UUID (FK) | Tim |
| position | VARCHAR | manager, leader, staff |
| reports_to | UUID (FK → users) | Atasan langsung |
| created_at | TIMESTAMPTZ | Waktu dibuat |
| updated_at | TIMESTAMPTZ | Waktu diupdate |

#### tickets
| Column | Type | Description |
|--------|------|-------------|
| id | UUID (PK) | Unique identifier |
| ticket_number | VARCHAR | Nomor ticket (REQ-000001, INC-000001, CHG-000001) |
| title | VARCHAR | Judul ticket |
| description | TEXT | Deskripsi detail |
| type | VARCHAR | incident, request, change_request |
| category | VARCHAR | Kategori (Database, Network, dll) |
| priority | VARCHAR | low, medium, high, critical |
| status | VARCHAR | open, in_progress, pending_approval, approved, rejected, done |
| created_by | UUID (FK → users) | Pembuat ticket |
| assigned_to | UUID (FK → users) | Ditugaskan ke (individual) |
| assigned_team_id | UUID (FK → teams) | Ditugaskan ke tim |
| resolved_at | TIMESTAMPTZ | Waktu selesai |
| created_at | TIMESTAMPTZ | Waktu dibuat |
| updated_at | TIMESTAMPTZ | Waktu diupdate |

### 2.2 Organization Structure

#### divisions (Level 1 - Tertinggi)
| Column | Type | Description |
|--------|------|-------------|
| id | UUID (PK) | Unique identifier |
| name | VARCHAR (UNIQUE) | Nama divisi (IT, BIZ, DOS, OPS, dll) |
| code | VARCHAR (UNIQUE) | Kode singkat |

#### departments (Level 2 - Child of Division)
| Column | Type | Description |
|--------|------|-------------|
| id | UUID (PK) | Unique identifier |
| division_id | UUID (FK → divisions) | Parent divisi |
| name | VARCHAR | Nama departemen (R&D, ITOps, PSD, dll) |
| code | VARCHAR | Kode singkat |

#### teams (Level 3 - Child of Department)
| Column | Type | Description |
|--------|------|-------------|
| id | UUID (PK) | Unique identifier |
| department_id | UUID (FK → departments) | Parent departemen |
| name | VARCHAR | Nama tim (DevOps, BackEnd, FrontEnd, dll) |
| email | VARCHAR | Email tim (untuk notifikasi) |

### 2.3 Approval & Workflow

#### approvals
| Column | Type | Description |
|--------|------|-------------|
| id | UUID (PK) | Unique identifier |
| ticket_id | UUID (FK → tickets) | Ticket yang di-approve |
| approver_id | UUID (FK → users) | Siapa yang approve |
| level | INTEGER | Level approval (1, 2, 3...) |
| decision | VARCHAR | approved, rejected |
| comment | TEXT | Komentar approver |
| decided_at | TIMESTAMPTZ | Waktu keputusan |

#### approval_configs
| Column | Type | Description |
|--------|------|-------------|
| id | UUID (PK) | Unique identifier |
| ticket_type | VARCHAR | Tipe ticket yang butuh approval |
| level | INTEGER | Level approval |
| approver_id | UUID (FK → users) | Default approver |

### 2.4 Application Management (ACL)

#### applications
| Column | Type | Description |
|--------|------|-------------|
| id | UUID (PK) | Unique identifier |
| name | VARCHAR (UNIQUE) | Nama aplikasi |
| code | VARCHAR (UNIQUE) | Kode (ticketing, project-board) |
| description | TEXT | Deskripsi |
| icon | VARCHAR | Material icon name |
| color | VARCHAR | Hex color |
| is_active | BOOLEAN | Status aktif |

#### user_app_access
| Column | Type | Description |
|--------|------|-------------|
| id | UUID (PK) | Unique identifier |
| user_id | UUID (FK → users) | User |
| app_id | UUID (FK → applications) | Aplikasi |
| role | VARCHAR | Role di aplikasi tersebut |
| granted_at | TIMESTAMPTZ | Waktu diberikan akses |
| granted_by | UUID (FK → users) | Siapa yang memberikan |

### 2.5 Supporting Tables

#### activity_logs
| Column | Type | Description |
|--------|------|-------------|
| id | UUID (PK) | Unique identifier |
| ticket_id | UUID (FK) | Ticket terkait |
| actor_id | UUID (FK) | Siapa yang melakukan |
| action | VARCHAR | ticket_created, assigned, approval_requested, dll |
| old_value | TEXT | Nilai lama |
| new_value | TEXT | Nilai baru |
| created_at | TIMESTAMPTZ | Waktu |

#### notifications
| Column | Type | Description |
|--------|------|-------------|
| id | UUID (PK) | Unique identifier |
| user_id | UUID (FK) | Penerima notifikasi |
| ticket_id | UUID (FK) | Ticket terkait |
| message | TEXT | Isi notifikasi |
| is_read | BOOLEAN | Sudah dibaca? |
| created_at | TIMESTAMPTZ | Waktu |

---

## 3. Business Flow

### 3.1 Ticket Creation Flow

```
┌─────────────────────────────────────────────────────────┐
│                    TICKET CREATION                        │
├─────────────────────────────────────────────────────────┤
│                                                          │
│  User/Agent creates ticket                               │
│       │                                                  │
│       ├── User (role: user)                              │
│       │     → Status: OPEN                               │
│       │     → Can create: incident, request              │
│       │     → Must manually "Request Approval"           │
│       │                                                  │
│       └── Agent (role: agent)                            │
│             → Status: PENDING APPROVAL (auto)            │
│             → Can create: incident, request, change_req  │
│             → Auto-notifies direct manager (reports_to)  │
│                                                          │
└─────────────────────────────────────────────────────────┘
```

### 3.2 Ticket Assignment Flow

```
┌─────────────────────────────────────────────────────────┐
│                  TICKET ASSIGNMENT                        │
├─────────────────────────────────────────────────────────┤
│                                                          │
│  Assign to User (individual)                             │
│       → assigned_to = user_id                            │
│       → Status → IN PROGRESS                             │
│       → Notification to assignee                         │
│       → Email notification (if configured)               │
│                                                          │
│  Assign to Team                                          │
│       → assigned_team_id = team_id                       │
│       → assigned_to = NULL (all members handle)          │
│       → Status → IN PROGRESS                             │
│       → Notification to ALL team members                 │
│       → All team members can see ticket                  │
│       → "Assigned to" shows all member names             │
│                                                          │
└─────────────────────────────────────────────────────────┘
```

### 3.3 Approval Flow

```
┌─────────────────────────────────────────────────────────┐
│                   APPROVAL FLOW                          │
├─────────────────────────────────────────────────────────┤
│                                                          │
│  Ticket submitted for approval                           │
│       │                                                  │
│       ├── Auto-notify direct manager (reports_to)        │
│       │                                                  │
│       ├── Who can approve?                               │
│       │     • Admin (any ticket)                         │
│       │     • Approver (tickets in their department)     │
│       │     • Direct manager (reports_to of creator)     │
│       │     • Position: Manager (same division)          │
│       │                                                  │
│       ├── Decision: APPROVED                             │
│       │     → Status → APPROVED                          │
│       │     → Notify creator                             │
│       │                                                  │
│       └── Decision: REJECTED                             │
│             → Status → REJECTED                          │
│             → Notify creator                             │
│                                                          │
└─────────────────────────────────────────────────────────┘
```

### 3.4 Ticket Visibility (Who sees what)

```
┌─────────────────────────────────────────────────────────┐
│                 TICKET VISIBILITY                         │
├─────────────────────────────────────────────────────────┤
│                                                          │
│  Admin         → ALL tickets                             │
│  Agent         → ALL tickets                             │
│  Approver      → Tickets from teams in their department  │
│                  + unassigned tickets                     │
│                                                          │
│  User with subordinates (reports_to chain):              │
│       → Own tickets                                      │
│       → All subordinates' tickets (recursive)            │
│       → Tickets assigned to their team                   │
│                                                          │
│  User without subordinates:                              │
│       • Staff  → Own tickets + team-assigned tickets     │
│       • Leader → Team members' tickets                   │
│       • Manager → Division members' tickets              │
│                                                          │
│  Regular user  → Only own tickets                        │
│                                                          │
└─────────────────────────────────────────────────────────┘
```

### 3.5 Organization Hierarchy

```
┌─────────────────────────────────────────────────────────┐
│              ORGANIZATION HIERARCHY                       │
├─────────────────────────────────────────────────────────┤
│                                                          │
│  Division (Top Level)                                    │
│    └── Department (Child of Division)                    │
│          └── Team (Child of Department)                  │
│                └── Users (Members of Team)               │
│                                                          │
│  Example:                                                │
│  IT (Division)                                           │
│    ├── R&D (Department)                                  │
│    │     ├── BackEnd (Team)                              │
│    │     ├── FrontEnd (Team)                             │
│    │     └── Mobile (Team)                               │
│    ├── ITOps (Department)                                │
│    │     ├── DevOps (Team)                               │
│    │     ├── PPQA (Team)                                 │
│    │     └── Service Desk (Team)                         │
│    └── PSD (Department)                                  │
│          ├── BQA (Team)                                  │
│          ├── UI/UX (Team)                                │
│          └── Project Manager (Team)                      │
│                                                          │
└─────────────────────────────────────────────────────────┘
```

### 3.6 Reporting Hierarchy (Org Chart)

```
┌─────────────────────────────────────────────────────────┐
│              REPORTING HIERARCHY                          │
├─────────────────────────────────────────────────────────┤
│                                                          │
│  CEO (Daniel T Pesireron)                                │
│    ├── CTO (Erwin Fransiscus) - Division Head IT         │
│    │     ├── Manager ITOps (Ryan Kenedy)                 │
│    │     │     ├── Leader DevOps (Achya Syahputra)       │
│    │     │     │     └── Staff (Tioya, Alfarizi, ...)    │
│    │     │     ├── Leader PsOps (Aditya Warman)          │
│    │     │     │     └── Staff (...)                     │
│    │     │     └── Staff (Nazar, Erwin Aji, ...)         │
│    │     ├── Manager PSD (Vanesha Asyariza)              │
│    │     │     └── ...                                   │
│    │     └── Manager R&D (M Arif Yudhistira)             │
│    │           └── ...                                   │
│    ├── Head BIZ (Yonanda Syafriade)                      │
│    │     └── ...                                         │
│    └── Head OPS (Omar Saladdin)                          │
│          └── ...                                         │
│                                                          │
│  reports_to field determines parent-child relationship   │
│  Used for: ticket visibility, approval routing,          │
│            org chart display                             │
│                                                          │
└─────────────────────────────────────────────────────────┘
```

---

## 4. Application Access Control (ACL)

```
┌─────────────────────────────────────────────────────────┐
│              APPLICATION ACCESS CONTROL                   │
├─────────────────────────────────────────────────────────┤
│                                                          │
│  Applications:                                           │
│    • Ticketing System (code: ticketing)                   │
│    • Project Board (code: project-board)                  │
│                                                          │
│  Access Control:                                         │
│    • Admin → always has access to all apps               │
│    • Other users → must be granted access per app        │
│    • Revoked users → cannot see app in switcher          │
│    • Backend middleware blocks API access                 │
│    • Frontend guard redirects to /dashboard              │
│                                                          │
│  Management:                                             │
│    • Admin can grant/revoke access                       │
│    • Bulk grant supported                                │
│    • Per-app role assignment                             │
│                                                          │
└─────────────────────────────────────────────────────────┘
```

---

## 5. API Endpoints

### 5.1 Authentication
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /api/v1/auth/login | Login (email + password) |
| POST | /api/v1/auth/refresh | Refresh token |
| POST | /api/v1/auth/logout | Logout |
| GET | /api/v1/auth/sso/login-url | SSO login URL |
| GET | /api/v1/auth/sso/callback | SSO callback |

### 5.2 Tickets
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/v1/tickets | List tickets (filtered by visibility) |
| POST | /api/v1/tickets | Create ticket |
| GET | /api/v1/tickets/{id} | Get ticket detail |
| PATCH | /api/v1/tickets/{id} | Update ticket |
| DELETE | /api/v1/tickets/{id} | Delete ticket |
| POST | /api/v1/tickets/{id}/submit | Submit for approval |
| POST | /api/v1/tickets/{id}/assign | Assign ticket |

### 5.3 Organization
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/v1/divisions | List divisions |
| GET | /api/v1/departments | List departments |
| GET | /api/v1/teams | List teams |
| GET | /api/v1/org-chart | Get org chart tree |

### 5.4 Applications
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/v1/me/apps | Get my accessible apps |
| GET | /api/v1/applications | List all apps |
| POST | /api/v1/applications | Create app (admin) |
| GET | /api/v1/applications/{id}/users | List app users |
| POST | /api/v1/applications/{id}/access | Grant access |
| DELETE | /api/v1/applications/{id}/access/{userId} | Revoke access |

### 5.5 Dashboard & Notifications
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/v1/dashboard | Dashboard stats |
| GET | /api/v1/notifications | List notifications |
| PATCH | /api/v1/notifications/{id}/read | Mark as read |

---

## 6. Role & Permission Matrix

| Feature | Admin | Approver | Agent | User |
|---------|-------|----------|-------|------|
| See all tickets | ✅ | Department scope | ✅ | Own + subordinates |
| Create ticket (all types) | ✅ | ✅ | ✅ | incident, request only |
| Assign ticket | ✅ | ✅ | ✅ | ❌ (staff) |
| Approve ticket | ✅ | ✅ | ❌ | If direct manager |
| Manage users | ✅ | ❌ | ❌ | ❌ |
| Manage org structure | ✅ | ❌ | ❌ | ❌ |
| View org chart | ✅ | ✅ | ❌ | ❌ |
| Manage applications | ✅ | ❌ | ❌ | ❌ |
| Manage webhooks | ✅ | ❌ | ❌ | ❌ |
| Kanban board | ✅ | ✅ | ✅ | ❌ |
| Activity logs | ✅ | ✅ | ✅ | ❌ |

---

## 7. Ticket Status Flow

```
OPEN → IN_PROGRESS → PENDING_APPROVAL → APPROVED → DONE
                                       → REJECTED
```

| Status | Description |
|--------|-------------|
| open | Ticket baru dibuat |
| in_progress | Sedang dikerjakan (assigned) |
| pending_approval | Menunggu approval dari atasan |
| approved | Disetujui |
| rejected | Ditolak |
| done | Selesai |

---

## 8. SLA Targets

| Priority | Target Resolution Time |
|----------|----------------------|
| Critical | 4 hours |
| High | 8 hours |
| Medium | 24 hours |
| Low | 72 hours |

---

## 9. Notification System

Notifikasi dikirim otomatis pada event:
1. **Ticket Created** → Semua agent & admin
2. **Ticket Assigned (User)** → Assignee + email
3. **Ticket Assigned (Team)** → Semua member team
4. **Approval Requested** → Direct manager (reports_to)
5. **Approval Decided** → Ticket creator
6. **Agent Creates Ticket** → Auto-request approval → notify manager

---

## 10. Integration Points

| System | Integration |
|--------|-------------|
| Keycloak SSO | Single Sign-On authentication |
| Google Drive | Ticket attachment storage |
| SMTP Email | Email notifications |
| Webhooks | External system notifications |

---

*Document generated: May 21, 2026*
*Version: 1.0*
*System: PCS ITSM Platform*
