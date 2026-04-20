package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	httpdelivery "github.com/org/itsm/internal/delivery/http"
	"github.com/org/itsm/internal/delivery/http/handler"
	"github.com/org/itsm/internal/infrastructure/webhook"
	notifinfra "github.com/org/itsm/internal/infrastructure/notification"
	"github.com/org/itsm/internal/repository/postgres"
	"github.com/org/itsm/internal/usecase"
	"github.com/org/itsm/pkg/config"
	jwtpkg "github.com/org/itsm/pkg/jwt"
	"github.com/org/itsm/pkg/logger"
)

func main() {
	cfg := config.Load()
	log := logger.New("itsm-server")

	// Database
	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Error("failed to connect to database", "error", err)
		panic(err)
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		log.Error("database ping failed", "error", err)
		panic(err)
	}
	log.Info("database connected")

	// Run migrations
	if err := runMigrations(context.Background(), pool, "migrations", log); err != nil {
		log.Error("migration failed", "error", err)
		panic(err)
	}

	// JWT
	jwtManager := jwtpkg.NewManager(cfg.JWTSecret, cfg.JWTRefreshSecret)

	// Repositories
	userRepo := postgres.NewUserRepository(pool)
	tokenRepo := postgres.NewRefreshTokenRepository(pool)
	ticketRepo := postgres.NewTicketRepository(pool)
	approvalRepo := postgres.NewApprovalRepository(pool)
	activityRepo := postgres.NewActivityLogRepository(pool)
	notifRepo := postgres.NewNotificationRepository(pool)
	webhookRepo := postgres.NewWebhookRepository(pool)
	deptRepo := postgres.NewDepartmentRepository(pool)
	divRepo := postgres.NewDivisionRepository(pool)
	teamRepo := postgres.NewTeamRepository(pool)
	projectRepo := postgres.NewProjectRepository(pool)
	projectColumnRepo := postgres.NewProjectColumnRepository(pool)
	projectRecordRepo := postgres.NewProjectRecordRepository(pool)
	projectActivityLogRepo := postgres.NewProjectActivityLogRepository(pool)
	projectMemberRepo := postgres.NewProjectMemberRepository(pool)

	// Infrastructure
	dispatcher := webhook.NewDispatcher(webhookRepo, userRepo, cfg.BaseURL)
	emailSvc := notifinfra.NewEmailService(cfg, log)
	if emailSvc.IsConfigured() {
		log.Info("email service configured", "host", cfg.SMTPHost)
	} else {
		log.Warn("email service not configured — email notifications disabled")
	}

	// Use cases
	authUC := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtManager)
	userUC := usecase.NewUserUseCase(userRepo, divRepo, teamRepo)
	webhookUC := usecase.NewWebhookUseCase(webhookRepo, dispatcher)
	ticketUC := usecase.NewTicketUseCase(ticketRepo, activityRepo, notifRepo, userRepo, webhookUC)
	approvalUC := usecase.NewApprovalUseCase(ticketRepo, approvalRepo, activityRepo, notifRepo, userRepo, webhookUC)
	assignmentUC := usecase.NewAssignmentUseCase(ticketRepo, userRepo, activityRepo, notifRepo, webhookUC, emailSvc)
	dashboardUC := usecase.NewDashboardUseCase(ticketRepo)
	notifUC := usecase.NewNotificationUseCase(notifRepo)
	orgUC := usecase.NewOrgUseCase(deptRepo, divRepo, teamRepo)
	projectBoardUC := usecase.NewProjectBoardUseCase(projectRepo, projectColumnRepo, projectRecordRepo, projectActivityLogRepo, projectMemberRepo)

	// Jira-like features repositories
	issueTypeRepo := postgres.NewIssueTypeRepository(pool)
	issueTypeSchemeRepo := postgres.NewIssueTypeSchemeRepository(pool)
	customFieldRepo := postgres.NewCustomFieldRepository(pool)
	customFieldOptionRepo := postgres.NewCustomFieldOptionRepository(pool)
	customFieldValueRepo := postgres.NewCustomFieldValueRepository(pool)
	workflowRepo := postgres.NewWorkflowRepository(pool)
	workflowStatusRepo := postgres.NewWorkflowStatusRepository(pool)
	workflowTransitionRepo := postgres.NewWorkflowTransitionRepository(pool)
	sprintRepo := postgres.NewSprintRepository(pool)
	sprintRecordRepo := postgres.NewSprintRecordRepository(pool)
	commentRepo := postgres.NewCommentRepository(pool)
	commentMentionRepo := postgres.NewCommentMentionRepository(pool)
	attachmentRepo := postgres.NewAttachmentRepository(pool)
	labelRepo := postgres.NewLabelRepository(pool)
	recordLabelRepo := postgres.NewRecordLabelRepository(pool)

	// Jira-like features use cases
	issueTypeUC := usecase.NewIssueTypeUseCase(issueTypeRepo, issueTypeSchemeRepo, projectRepo, projectActivityLogRepo, projectMemberRepo)
	customFieldUC := usecase.NewCustomFieldUseCase(customFieldRepo, customFieldOptionRepo, customFieldValueRepo, projectRepo, projectRecordRepo, projectActivityLogRepo, projectMemberRepo)
	workflowUC := usecase.NewWorkflowUseCase(workflowRepo, workflowStatusRepo, workflowTransitionRepo, projectRepo, projectRecordRepo, projectActivityLogRepo, projectMemberRepo)
	sprintUC := usecase.NewSprintUseCase(sprintRepo, sprintRecordRepo, projectRepo, projectRecordRepo, projectActivityLogRepo, projectMemberRepo)
	backlogUC := usecase.NewBacklogUseCase(sprintRecordRepo, projectRepo, projectRecordRepo, projectActivityLogRepo, projectMemberRepo)
	commentUC := usecase.NewCommentUseCase(commentRepo, commentMentionRepo, projectRecordRepo, userRepo, projectActivityLogRepo, projectMemberRepo)
	attachmentUC := usecase.NewAttachmentUseCase(attachmentRepo, projectRecordRepo, projectActivityLogRepo, projectMemberRepo)
	labelUC := usecase.NewLabelUseCase(labelRepo, recordLabelRepo, projectRepo, projectRecordRepo, projectActivityLogRepo, projectMemberRepo)
	bulkOpUC := usecase.NewBulkOperationUseCase(projectRecordRepo, workflowStatusRepo, labelRepo, recordLabelRepo, projectRepo, projectActivityLogRepo, projectMemberRepo)
	searchUC := usecase.NewSearchUseCase(projectRecordRepo, projectRepo, projectMemberRepo)

	// Project Board Features repositories and use cases
	releaseRepo := postgres.NewReleaseRepository(pool)
	releaseRecordRepo := postgres.NewReleaseRecordRepository(pool)
	componentRepo := postgres.NewComponentRepository(pool)
	reportsUC := usecase.NewReportsUseCase(projectRecordRepo, projectMemberRepo, sprintRepo, sprintRecordRepo)
	releaseUC := usecase.NewReleaseUseCase(releaseRepo, releaseRecordRepo, projectMemberRepo)
	componentUC := usecase.NewComponentUseCase(componentRepo, projectRecordRepo, projectMemberRepo)

	// Handlers
	handlers := &httpdelivery.Handlers{
		Auth:                 handler.NewAuthHandler(authUC),
		User:                 handler.NewUserHandler(userUC),
		Ticket:               handler.NewTicketHandler(ticketUC, approvalUC, assignmentUC, activityRepo),
		Approval:             handler.NewApprovalHandler(approvalUC),
		Dashboard:            handler.NewDashboardHandler(dashboardUC),
		Notification:         handler.NewNotificationHandler(notifUC),
		Webhook:              handler.NewWebhookHandler(webhookUC),
		Attachment:           handler.NewAttachmentHandler(pool),
		Org:                  handler.NewOrgHandler(orgUC),
		SSO:                  handler.NewSSOHandler(cfg, userRepo, jwtManager),
		Project:              handler.NewProjectHandler(projectBoardUC),
		Jira:                 handler.NewJiraHandler(issueTypeUC, customFieldUC, workflowUC, sprintUC, backlogUC, commentUC, attachmentUC, labelUC, bulkOpUC, searchUC),
		ProjectBoardFeatures: handler.NewProjectBoardFeaturesHandler(reportsUC, releaseUC, componentUC, projectRecordRepo, projectActivityLogRepo, projectMemberRepo),
	}

	// Router
	router := httpdelivery.NewRouter(handlers, jwtManager, &dbPinger{pool: pool})

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Info("server starting", "addr", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Error("server error", "error", err)
	}
}

// runMigrations executes all *.up.sql files in order, tracking applied ones.
func runMigrations(ctx context.Context, pool *pgxpool.Pool, dir string, log *logger.Logger) error {
	// Create migrations tracking table
	_, err := pool.Exec(ctx, `CREATE TABLE IF NOT EXISTS schema_migrations (
		version VARCHAR(255) PRIMARY KEY,
		applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	)`)
	if err != nil {
		return fmt.Errorf("create migrations table: %w", err)
	}

	// Read applied migrations
	rows, err := pool.Query(ctx, `SELECT version FROM schema_migrations`)
	if err != nil {
		return err
	}
	applied := map[string]bool{}
	for rows.Next() {
		var v string
		if err := rows.Scan(&v); err != nil {
			return err
		}
		applied[v] = true
	}
	rows.Close()

	// Find all .up.sql files
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("read migrations dir: %w", err)
	}
	var files []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".up.sql") {
			files = append(files, e.Name())
		}
	}
	sort.Strings(files)

	for _, f := range files {
		version := strings.TrimSuffix(f, ".up.sql")
		if applied[version] {
			continue
		}
		content, err := os.ReadFile(filepath.Join(dir, f))
		if err != nil {
			return fmt.Errorf("read %s: %w", f, err)
		}
		if _, err := pool.Exec(ctx, string(content)); err != nil {
			return fmt.Errorf("execute %s: %w", f, err)
		}
		if _, err := pool.Exec(ctx, `INSERT INTO schema_migrations (version) VALUES ($1)`, version); err != nil {
			return fmt.Errorf("record migration %s: %w", f, err)
		}
		log.Info("migration applied", "version", version)
	}
	return nil
}

type dbPinger struct {
	pool *pgxpool.Pool
}

func (d *dbPinger) Ping() error {
	return d.pool.Ping(context.Background())
}
