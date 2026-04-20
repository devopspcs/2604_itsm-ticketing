# Jira-like Project Board - Quick Start Guide

## 🚀 5-Minute Setup

### Prerequisites
- PostgreSQL running
- Go 1.21+
- Node.js 16+
- Docker (optional)

### Step 1: Database Setup (2 min)

```bash
# Run migrations
migrate -path backend/migrations -database "postgresql://postgres:password@localhost:5432/itsm_db?sslmode=disable" up

# Or manually with psql
psql -U postgres -d itsm_db -f backend/migrations/000009_jira_features.up.sql
psql -U postgres -d itsm_db -f backend/migrations/000011_seed_jira_test_data.up.sql
```

### Step 2: Start Backend (1 min)

```bash
cd backend
go run cmd/main.go
# Backend running at http://localhost:8080
```

### Step 3: Start Frontend (1 min)

```bash
cd frontend
npm install
npm run dev
# Frontend running at http://localhost:5173
```

### Step 4: Verify Setup (1 min)

```bash
./verify-jira-setup.sh
```

## 📋 Common Commands

### Database

```bash
# Connect to database
psql -U postgres -d itsm_db

# Check issue types
SELECT * FROM issue_types;

# Check test project
SELECT * FROM projects WHERE name = 'Test Project - Jira Board';

# Check active sprint
SELECT * FROM sprints WHERE status = 'Active';

# Check custom fields
SELECT * FROM custom_fields;

# Check labels
SELECT * FROM labels;

# Backup database
pg_dump -U postgres -d itsm_db > backup.sql

# Restore database
psql -U postgres -d itsm_db < backup.sql
```

### Backend

```bash
# Build
cd backend && go build -o itsm-server

# Run
./itsm-server

# Run tests
go test ./...

# Run with debug logging
LOG_LEVEL=debug go run cmd/main.go

# Check API health
curl http://localhost:8080/api/v1/health
```

### Frontend

```bash
# Install dependencies
npm install

# Development server
npm run dev

# Build for production
npm run build

# Run tests
npm test

# Type check
npm run type-check

# Lint
npm run lint
```

### Docker

```bash
# Build backend image
docker build -f backend/Dockerfile -t itsm-backend:latest backend/

# Build frontend image
docker build -f frontend/Dockerfile -t itsm-frontend:latest frontend/

# Run with docker-compose
docker-compose up -d

# View logs
docker-compose logs -f backend
docker-compose logs -f frontend

# Stop services
docker-compose down
```

## 🔍 API Testing

### Get Issue Types

```bash
curl -X GET http://localhost:8080/api/v1/projects/{projectId}/issue-types \
  -H "Authorization: Bearer {token}"
```

### Get Active Sprint

```bash
curl -X GET http://localhost:8080/api/v1/projects/{projectId}/sprints/active \
  -H "Authorization: Bearer {token}"
```

### Get Sprint Records

```bash
curl -X GET http://localhost:8080/api/v1/sprints/{sprintId}/records \
  -H "Authorization: Bearer {token}"
```

### Get Backlog

```bash
curl -X GET http://localhost:8080/api/v1/projects/{projectId}/backlog \
  -H "Authorization: Bearer {token}"
```

### Get Custom Fields

```bash
curl -X GET http://localhost:8080/api/v1/projects/{projectId}/custom-fields \
  -H "Authorization: Bearer {token}"
```

### Get Labels

```bash
curl -X GET http://localhost:8080/api/v1/projects/{projectId}/labels \
  -H "Authorization: Bearer {token}"
```

### Add Comment

```bash
curl -X POST http://localhost:8080/api/v1/records/{recordId}/comments \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{"text": "This is a comment"}'
```

### Upload Attachment

```bash
curl -X POST http://localhost:8080/api/v1/records/{recordId}/attachments \
  -H "Authorization: Bearer {token}" \
  -F "file=@document.pdf"
```

### Transition Record

```bash
curl -X POST http://localhost:8080/api/v1/records/{recordId}/transition \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{"to_status_id": "{statusId}"}'
```

## 🧪 Testing Checklist

- [ ] Database connected
- [ ] Migrations applied
- [ ] Test data inserted
- [ ] Backend running
- [ ] Frontend running
- [ ] Can login
- [ ] Can view project board
- [ ] Can see sprint board
- [ ] Can drag records
- [ ] Can add comments
- [ ] Can upload files
- [ ] Can add labels
- [ ] Can view backlog

## 🐛 Troubleshooting

### "No Active Sprint"
```bash
# Check if sprint exists
psql -U postgres -d itsm_db -c "SELECT * FROM sprints WHERE status = 'Active';"

# Create sprint if missing
psql -U postgres -d itsm_db -c "
INSERT INTO sprints (id, project_id, name, goal, start_date, end_date, status, actual_start_date, created_at)
SELECT gen_random_uuid(), p.id, 'Sprint 1', 'Test sprint', CURRENT_DATE, CURRENT_DATE + INTERVAL '14 days', 'Active', CURRENT_DATE, NOW()
FROM projects p WHERE p.name = 'Test Project - Jira Board' LIMIT 1;
"
```

### "401 Unauthorized"
```bash
# Get new token
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "admin@example.com", "password": "password"}'

# Use token in requests
curl -H "Authorization: Bearer {token}" http://localhost:8080/api/v1/projects/{projectId}/issue-types
```

### "Database connection failed"
```bash
# Check PostgreSQL is running
sudo systemctl status postgresql

# Test connection
psql -h localhost -U postgres -d itsm_db

# Check credentials in .env
cat .env | grep DB_
```

### "Drag-and-drop not working"
```bash
# Check @dnd-kit is installed
cd frontend && npm list @dnd-kit/core

# Install if missing
npm install @dnd-kit/core @dnd-kit/utilities

# Clear cache and rebuild
rm -rf node_modules package-lock.json
npm install
npm run dev
```

## 📚 Documentation

- **Setup Guide**: JIRA_BOARD_SETUP_GUIDE.md
- **API Reference**: JIRA_API_REFERENCE.md
- **Troubleshooting**: JIRA_TROUBLESHOOTING.md
- **Implementation Summary**: JIRA_BOARD_IMPLEMENTATION_SUMMARY.md

## 🔗 Useful Links

- Frontend: http://localhost:5173
- Backend API: http://localhost:8080/api/v1
- PostgreSQL: localhost:5432
- Project Board: http://localhost:5173/projects/{projectId}
- Sprint Board: http://localhost:5173/projects/{projectId}/sprint
- Backlog: http://localhost:5173/projects/{projectId}/backlog

## 📊 Database Schema

```
issue_types
├── id (UUID)
├── name (VARCHAR)
├── icon (VARCHAR)
└── description (TEXT)

workflows
├── id (UUID)
├── project_id (UUID)
├── name (VARCHAR)
└── initial_status (VARCHAR)

workflow_statuses
├── id (UUID)
├── workflow_id (UUID)
├── status_name (VARCHAR)
└── status_order (INTEGER)

sprints
├── id (UUID)
├── project_id (UUID)
├── name (VARCHAR)
├── goal (TEXT)
├── start_date (DATE)
├── end_date (DATE)
├── status (VARCHAR)
└── actual_start_date (DATE)

custom_fields
├── id (UUID)
├── project_id (UUID)
├── name (VARCHAR)
├── field_type (VARCHAR)
└── is_required (BOOLEAN)

labels
├── id (UUID)
├── project_id (UUID)
├── name (VARCHAR)
└── color (VARCHAR)

comments
├── id (UUID)
├── record_id (UUID)
├── author_id (UUID)
├── text (TEXT)
└── created_at (TIMESTAMPTZ)

attachments
├── id (UUID)
├── record_id (UUID)
├── file_name (VARCHAR)
├── file_size (BIGINT)
├── file_path (VARCHAR)
└── uploader_id (UUID)
```

## 🎯 Next Steps

1. ✅ Run setup commands above
2. ✅ Verify with `./verify-jira-setup.sh`
3. ✅ Open http://localhost:5173 in browser
4. ✅ Navigate to project board
5. ✅ Test features (drag-drop, comments, attachments, labels)
6. ✅ Check logs for any errors
7. ✅ Read full documentation for advanced features

## 💡 Tips

- Use `LOG_LEVEL=debug` for detailed backend logs
- Check browser DevTools Console for frontend errors
- Use `psql` to verify database state
- Run `./verify-jira-setup.sh` to diagnose issues
- Check `JIRA_TROUBLESHOOTING.md` for common problems

## 🆘 Need Help?

1. Check **JIRA_TROUBLESHOOTING.md** for common issues
2. Run **verify-jira-setup.sh** to diagnose problems
3. Check backend logs: `tail -f backend.log`
4. Check frontend console: Browser DevTools (F12)
5. Check database: `psql -U postgres -d itsm_db`

---

**Ready to go!** 🎉

Start with the 5-minute setup above, then explore the full documentation for advanced features.
