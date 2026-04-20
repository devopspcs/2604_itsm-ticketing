package usecase

import (
	"context"
	"regexp"
	"time"

	"github.com/google/uuid"
	"github.com/org/itsm/internal/domain/entity"
	"github.com/org/itsm/internal/domain/repository"
	domainUC "github.com/org/itsm/internal/domain/usecase"
	"github.com/org/itsm/pkg/apperror"
)

type commentUseCase struct {
	commentRepo     repository.CommentRepository
	mentionRepo     repository.CommentMentionRepository
	recordRepo      repository.ProjectRecordRepository
	userRepo        repository.UserRepository
	activityRepo    repository.ProjectActivityLogRepository
	memberRepo      repository.ProjectMemberRepository
}

func NewCommentUseCase(
	commentRepo repository.CommentRepository,
	mentionRepo repository.CommentMentionRepository,
	recordRepo repository.ProjectRecordRepository,
	userRepo repository.UserRepository,
	activityRepo repository.ProjectActivityLogRepository,
	memberRepo repository.ProjectMemberRepository,
) domainUC.CommentUseCase {
	return &commentUseCase{
		commentRepo:     commentRepo,
		mentionRepo:     mentionRepo,
		recordRepo:      recordRepo,
		userRepo:        userRepo,
		activityRepo:    activityRepo,
		memberRepo:      memberRepo,
	}
}

// AddComment creates a new comment on a record
func (uc *commentUseCase) AddComment(ctx context.Context, recordID uuid.UUID, req domainUC.AddCommentRequest, requester domainUC.UserClaims) (*entity.Comment, error) {
	// Get record
	record, err := uc.recordRepo.FindByID(ctx, recordID)
	if err != nil {
		return nil, err
	}
	if record == nil {
		return nil, apperror.ErrNotFound
	}

	// Verify user is a member of the project
	isMember, _ := uc.memberRepo.IsMember(ctx, record.ProjectID, requester.UserID)
	if !isMember {
		return nil, apperror.ErrForbidden
	}

	// Create comment
	comment := &entity.Comment{
		ID:        uuid.New(),
		RecordID:  recordID,
		AuthorID:  requester.UserID,
		Text:      req.Text,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	if err := uc.commentRepo.Create(ctx, comment); err != nil {
		return nil, err
	}

	// Parse mentions and create mention records
	mentions := uc.ParseMentions(req.Text)
	for _, username := range mentions {
		// Find user by username
		user, err := uc.userRepo.FindByEmail(ctx, username+"@example.com") // Simplified - would need proper lookup
		if err == nil && user != nil {
			mention := &entity.CommentMention{
				ID:              uuid.New(),
				CommentID:       comment.ID,
				MentionedUserID: user.ID,
				CreatedAt:       time.Now().UTC(),
			}
			_ = uc.mentionRepo.Create(ctx, mention)

			// TODO: Send notification to mentioned user
		}
	}

	// Log activity
	uc.logActivity(ctx, record.ProjectID, &recordID, requester.UserID, "comment_added", "Comment added to record")

	return comment, nil
}

// ListComments returns all comments for a record
func (uc *commentUseCase) ListComments(ctx context.Context, recordID uuid.UUID, requester domainUC.UserClaims) ([]*entity.Comment, error) {
	// Get record
	record, err := uc.recordRepo.FindByID(ctx, recordID)
	if err != nil {
		return nil, err
	}
	if record == nil {
		return nil, apperror.ErrNotFound
	}

	// Verify user is a member
	isMember, _ := uc.memberRepo.IsMember(ctx, record.ProjectID, requester.UserID)
	if !isMember {
		return nil, apperror.ErrForbidden
	}

	comments, err := uc.commentRepo.ListByRecord(ctx, recordID)
	if err != nil {
		return nil, err
	}
	if comments == nil {
		comments = []*entity.Comment{}
	}

	return comments, nil
}

// UpdateComment updates a comment's text
func (uc *commentUseCase) UpdateComment(ctx context.Context, commentID uuid.UUID, req domainUC.UpdateCommentRequest, requester domainUC.UserClaims) (*entity.Comment, error) {
	// Get comment
	comment, err := uc.commentRepo.FindByID(ctx, commentID)
	if err != nil {
		return nil, err
	}
	if comment == nil {
		return nil, apperror.ErrNotFound
	}

	// Verify user is the author
	if comment.AuthorID != requester.UserID {
		return nil, apperror.ErrForbidden
	}

	// Get record to verify project membership
	record, err := uc.recordRepo.FindByID(ctx, comment.RecordID)
	if err != nil {
		return nil, err
	}
	if record == nil {
		return nil, apperror.ErrNotFound
	}

	// Update comment
	comment.Text = req.Text
	comment.UpdatedAt = time.Now().UTC()

	if err := uc.commentRepo.Update(ctx, comment); err != nil {
		return nil, err
	}

	// Delete old mentions
	_ = uc.mentionRepo.DeleteByComment(ctx, commentID)

	// Parse new mentions and create mention records
	mentions := uc.ParseMentions(req.Text)
	for _, username := range mentions {
		// Find user by username
		user, err := uc.userRepo.FindByEmail(ctx, username+"@example.com") // Simplified
		if err == nil && user != nil {
			mention := &entity.CommentMention{
				ID:              uuid.New(),
				CommentID:       comment.ID,
				MentionedUserID: user.ID,
				CreatedAt:       time.Now().UTC(),
			}
			_ = uc.mentionRepo.Create(ctx, mention)
		}
	}

	// Log activity
	uc.logActivity(ctx, record.ProjectID, &comment.RecordID, requester.UserID, "comment_updated", "Comment updated")

	return comment, nil
}

// DeleteComment deletes a comment
func (uc *commentUseCase) DeleteComment(ctx context.Context, commentID uuid.UUID, requester domainUC.UserClaims) error {
	// Get comment
	comment, err := uc.commentRepo.FindByID(ctx, commentID)
	if err != nil {
		return err
	}
	if comment == nil {
		return apperror.ErrNotFound
	}

	// Verify user is the author
	if comment.AuthorID != requester.UserID {
		return apperror.ErrForbidden
	}

	// Get record
	record, err := uc.recordRepo.FindByID(ctx, comment.RecordID)
	if err != nil {
		return err
	}
	if record == nil {
		return apperror.ErrNotFound
	}

	// Delete mentions
	_ = uc.mentionRepo.DeleteByComment(ctx, commentID)

	// Delete comment
	if err := uc.commentRepo.Delete(ctx, commentID); err != nil {
		return err
	}

	// Log activity
	uc.logActivity(ctx, record.ProjectID, &comment.RecordID, requester.UserID, "comment_deleted", "Comment deleted")

	return nil
}

// ParseMentions extracts @mentions from text
func (uc *commentUseCase) ParseMentions(text string) []string {
	// Regular expression to match @username patterns
	// Matches @word where word contains letters, numbers, underscores, and hyphens
	re := regexp.MustCompile(`@([a-zA-Z0-9_-]+)`)
	matches := re.FindAllStringSubmatch(text, -1)

	var mentions []string
	seen := make(map[string]bool)

	for _, match := range matches {
		if len(match) > 1 {
			username := match[1]
			// Avoid duplicates
			if !seen[username] {
				mentions = append(mentions, username)
				seen[username] = true
			}
		}
	}

	return mentions
}

// logActivity logs an activity to the project activity log
func (uc *commentUseCase) logActivity(ctx context.Context, projectID uuid.UUID, recordID *uuid.UUID, actorID uuid.UUID, action, detail string) {
	activity := &entity.ProjectActivityLog{
		ID:        uuid.New(),
		ProjectID: projectID,
		RecordID:  recordID,
		ActorID:   actorID,
		Action:    action,
		Detail:    detail,
		CreatedAt: time.Now().UTC(),
	}
	_ = uc.activityRepo.Append(ctx, activity)
}
