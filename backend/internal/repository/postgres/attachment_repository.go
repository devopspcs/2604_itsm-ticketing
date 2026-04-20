package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/org/itsm/internal/domain/entity"
	"github.com/org/itsm/internal/domain/repository"
	"github.com/org/itsm/pkg/apperror"
)

type attachmentRepository struct {
	db *pgxpool.Pool
}

func NewAttachmentRepository(db *pgxpool.Pool) repository.AttachmentRepository {
	return &attachmentRepository{db: db}
}

func (r *attachmentRepository) Create(ctx context.Context, attachment *entity.Attachment) error {
	query := `INSERT INTO attachments (id, record_id, file_name, file_size, file_type, file_path, uploader_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.Exec(ctx, query,
		attachment.ID, attachment.RecordID, attachment.FileName, attachment.FileSize,
		attachment.FileType, attachment.FilePath, attachment.UploaderID, attachment.CreatedAt)
	return err
}

func (r *attachmentRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Attachment, error) {
	query := `SELECT id, record_id, file_name, file_size, file_type, file_path, uploader_id, created_at FROM attachments WHERE id=$1`
	row := r.db.QueryRow(ctx, query, id)
	return scanAttachment(row)
}

func (r *attachmentRepository) ListByRecord(ctx context.Context, recordID uuid.UUID) ([]*entity.Attachment, error) {
	query := `SELECT a.id, a.record_id, a.file_name, a.file_size, a.file_type, a.file_path, a.uploader_id, COALESCE(u.full_name, u.email, '') as uploader_name, a.created_at
		FROM attachments a LEFT JOIN users u ON a.uploader_id = u.id
		WHERE a.record_id=$1 ORDER BY a.created_at DESC`
	rows, err := r.db.Query(ctx, query, recordID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attachments []*entity.Attachment
	for rows.Next() {
		attachment := &entity.Attachment{}
		if err := rows.Scan(&attachment.ID, &attachment.RecordID, &attachment.FileName, &attachment.FileSize,
			&attachment.FileType, &attachment.FilePath, &attachment.UploaderID, &attachment.UploaderName, &attachment.CreatedAt); err != nil {
			return nil, err
		}
		attachments = append(attachments, attachment)
	}
	return attachments, nil
}

func (r *attachmentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM attachments WHERE id=$1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *attachmentRepository) DeleteByRecord(ctx context.Context, recordID uuid.UUID) error {
	query := `DELETE FROM attachments WHERE record_id=$1`
	_, err := r.db.Exec(ctx, query, recordID)
	return err
}

func scanAttachment(row pgx.Row) (*entity.Attachment, error) {
	attachment := &entity.Attachment{}
	err := row.Scan(&attachment.ID, &attachment.RecordID, &attachment.FileName, &attachment.FileSize,
		&attachment.FileType, &attachment.FilePath, &attachment.UploaderID, &attachment.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrNotFound
		}
		return nil, err
	}
	return attachment, nil
}
