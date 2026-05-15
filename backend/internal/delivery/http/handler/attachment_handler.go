package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/org/itsm/internal/delivery/http/middleware"
	"github.com/org/itsm/internal/infrastructure/storage"
	"github.com/org/itsm/pkg/apperror"
)

type AttachmentHandler struct {
	db    *pgxpool.Pool
	drive *storage.DriveStorage
}

func NewAttachmentHandler(db *pgxpool.Pool, drive *storage.DriveStorage) *AttachmentHandler {
	return &AttachmentHandler{db: db, drive: drive}
}

type Attachment struct {
	ID          string    `json:"id"`
	TicketID    string    `json:"ticket_id"`
	UploadedBy  string    `json:"uploaded_by"`
	Filename    string    `json:"filename"`
	FileSize    int64     `json:"file_size"`
	MimeType    string    `json:"mime_type"`
	DriveFileID string    `json:"drive_file_id,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// Upload handles multipart file upload for a ticket — stores in Google Drive
func (h *AttachmentHandler) Upload(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r)
	if !ok {
		apperror.WriteError(w, apperror.ErrTokenInvalid)
		return
	}

	ticketID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}

	// Max 50MB
	if err := r.ParseMultipartForm(50 << 20); err != nil {
		apperror.WriteError(w, apperror.ErrValidation.WithDetails(map[string]interface{}{"error": "file too large, max 50MB"}))
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation.WithDetails(map[string]interface{}{"error": "no file provided"}))
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		apperror.WriteError(w, apperror.ErrInternal)
		return
	}

	mimeType := header.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	// Upload to Google Drive
	driveFileID := ""
	if h.drive != nil {
		driveFileName := ticketID.String() + "/" + header.Filename
		driveFileID, err = h.drive.Upload(r.Context(), driveFileName, mimeType, data)
		if err != nil {
			// Log error and fallback to DB storage
			fmt.Printf("[DRIVE ERROR] Upload failed for %s: %v\n", header.Filename, err)
			driveFileID = ""
		} else {
			fmt.Printf("[DRIVE OK] Uploaded %s -> fileID: %s\n", header.Filename, driveFileID)
		}
	} else {
		fmt.Println("[DRIVE WARN] Drive storage is nil, using DB fallback")
	}

	id := uuid.New()
	now := time.Now().UTC()

	if driveFileID != "" {
		// Store metadata only (file is in Google Drive)
		_, err = h.db.Exec(r.Context(),
			`INSERT INTO ticket_attachments (id, ticket_id, uploaded_by, filename, file_size, mime_type, drive_file_id, created_at)
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
			id, ticketID, claims.UserID, header.Filename, header.Size, mimeType, driveFileID, now,
		)
	} else {
		// Fallback: store file data in DB
		_, err = h.db.Exec(r.Context(),
			`INSERT INTO ticket_attachments (id, ticket_id, uploaded_by, filename, file_size, mime_type, file_data, created_at)
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
			id, ticketID, claims.UserID, header.Filename, header.Size, mimeType, data, now,
		)
	}
	if err != nil {
		apperror.WriteError(w, apperror.ErrInternal)
		return
	}

	apperror.WriteJSON(w, http.StatusCreated, Attachment{
		ID:          id.String(),
		TicketID:    ticketID.String(),
		Filename:    header.Filename,
		FileSize:    header.Size,
		MimeType:    mimeType,
		DriveFileID: driveFileID,
		CreatedAt:   now,
	})
}

// List returns all attachments for a ticket
func (h *AttachmentHandler) List(w http.ResponseWriter, r *http.Request) {
	ticketID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}

	rows, err := h.db.Query(r.Context(),
		`SELECT id, ticket_id, uploaded_by, filename, file_size, mime_type, COALESCE(drive_file_id, '') as drive_file_id, created_at
		 FROM ticket_attachments WHERE ticket_id = $1 ORDER BY created_at DESC`,
		ticketID,
	)
	if err != nil {
		apperror.WriteError(w, apperror.ErrInternal)
		return
	}
	defer rows.Close()

	var attachments []Attachment
	for rows.Next() {
		var a Attachment
		if err := rows.Scan(&a.ID, &a.TicketID, &a.UploadedBy, &a.Filename, &a.FileSize, &a.MimeType, &a.DriveFileID, &a.CreatedAt); err != nil {
			continue
		}
		attachments = append(attachments, a)
	}
	if attachments == nil {
		attachments = []Attachment{}
	}

	apperror.WriteJSON(w, http.StatusOK, attachments)
}

// Download serves the file content (from Drive or DB fallback)
func (h *AttachmentHandler) Download(w http.ResponseWriter, r *http.Request) {
	attachmentID, err := uuid.Parse(chi.URLParam(r, "attachmentId"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}

	var filename, mimeType, driveFileID string
	var data []byte
	err = h.db.QueryRow(r.Context(),
		`SELECT filename, mime_type, COALESCE(drive_file_id, ''), file_data FROM ticket_attachments WHERE id = $1`,
		attachmentID,
	).Scan(&filename, &mimeType, &driveFileID, &data)
	if err != nil {
		apperror.WriteError(w, apperror.ErrNotFound)
		return
	}

	// If stored in Google Drive, download from there
	if driveFileID != "" && h.drive != nil {
		driveData, err := h.drive.Download(r.Context(), driveFileID)
		if err == nil {
			data = driveData
		}
		// If Drive download fails, fall through to DB data if available
	}

	if data == nil {
		apperror.WriteError(w, apperror.ErrNotFound)
		return
	}

	w.Header().Set("Content-Type", mimeType)
	w.Header().Set("Content-Disposition", `attachment; filename="`+filename+`"`)
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// Delete removes an attachment (from Drive and DB)
func (h *AttachmentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r)
	if !ok {
		apperror.WriteError(w, apperror.ErrTokenInvalid)
		return
	}

	attachmentID, err := uuid.Parse(chi.URLParam(r, "attachmentId"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}

	// Only uploader or admin can delete
	var uploadedBy uuid.UUID
	var driveFileID string
	err = h.db.QueryRow(r.Context(),
		`SELECT uploaded_by, COALESCE(drive_file_id, '') FROM ticket_attachments WHERE id = $1`, attachmentID,
	).Scan(&uploadedBy, &driveFileID)
	if err != nil {
		apperror.WriteError(w, apperror.ErrNotFound)
		return
	}

	if claims.UserID != uploadedBy && claims.Role != "admin" {
		apperror.WriteError(w, apperror.ErrForbidden)
		return
	}

	// Delete from Google Drive if applicable
	if driveFileID != "" && h.drive != nil {
		_ = h.drive.Delete(r.Context(), driveFileID)
	}

	h.db.Exec(r.Context(), `DELETE FROM ticket_attachments WHERE id = $1`, attachmentID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "deleted"})
}

type Note struct {
	ID        string    `json:"id"`
	TicketID  string    `json:"ticket_id"`
	AuthorID  string    `json:"author_id"`
	Content   string    `json:"content"`
	ImageName *string   `json:"image_name,omitempty"`
	ImageMime *string   `json:"image_mime,omitempty"`
	HasImage  bool      `json:"has_image"`
	CreatedAt time.Time `json:"created_at"`
}

// AddNote handles multipart note with optional image
func (h *AttachmentHandler) AddNote(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r)
	if !ok {
		apperror.WriteError(w, apperror.ErrTokenInvalid)
		return
	}

	ticketID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}

	// Parse multipart (max 10MB)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		// Fallback to JSON body
		var body struct {
			Content string `json:"content"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Content == "" {
			apperror.WriteError(w, apperror.ErrValidation.WithDetails(map[string]interface{}{"error": "content is required"}))
			return
		}
		id := uuid.New()
		now := time.Now().UTC()
		_, err = h.db.Exec(r.Context(),
			`INSERT INTO ticket_notes (id, ticket_id, author_id, content, created_at) VALUES ($1,$2,$3,$4,$5)`,
			id, ticketID, claims.UserID, body.Content, now,
		)
		if err != nil {
			apperror.WriteError(w, apperror.ErrInternal)
			return
		}
		apperror.WriteJSON(w, http.StatusCreated, Note{
			ID: id.String(), TicketID: ticketID.String(), AuthorID: claims.UserID.String(),
			Content: body.Content, HasImage: false, CreatedAt: now,
		})
		return
	}

	content := r.FormValue("content")
	if content == "" {
		apperror.WriteError(w, apperror.ErrValidation.WithDetails(map[string]interface{}{"error": "content is required"}))
		return
	}

	id := uuid.New()
	now := time.Now().UTC()

	// Check for image
	file, header, fileErr := r.FormFile("image")
	if fileErr == nil {
		defer file.Close()
		imgData, _ := io.ReadAll(file)
		imgName := header.Filename
		imgMime := header.Header.Get("Content-Type")

		_, err = h.db.Exec(r.Context(),
			`INSERT INTO ticket_notes (id, ticket_id, author_id, content, image_data, image_name, image_mime, created_at)
			 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
			id, ticketID, claims.UserID, content, imgData, imgName, imgMime, now,
		)
		if err != nil {
			apperror.WriteError(w, apperror.ErrInternal)
			return
		}
		apperror.WriteJSON(w, http.StatusCreated, Note{
			ID: id.String(), TicketID: ticketID.String(), AuthorID: claims.UserID.String(),
			Content: content, ImageName: &imgName, ImageMime: &imgMime, HasImage: true, CreatedAt: now,
		})
	} else {
		_, err = h.db.Exec(r.Context(),
			`INSERT INTO ticket_notes (id, ticket_id, author_id, content, created_at) VALUES ($1,$2,$3,$4,$5)`,
			id, ticketID, claims.UserID, content, now,
		)
		if err != nil {
			apperror.WriteError(w, apperror.ErrInternal)
			return
		}
		apperror.WriteJSON(w, http.StatusCreated, Note{
			ID: id.String(), TicketID: ticketID.String(), AuthorID: claims.UserID.String(),
			Content: content, HasImage: false, CreatedAt: now,
		})
	}
}

// ListNotes returns all notes for a ticket
func (h *AttachmentHandler) ListNotes(w http.ResponseWriter, r *http.Request) {
	ticketID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}

	rows, err := h.db.Query(r.Context(),
		`SELECT id, ticket_id, author_id, content, image_name, image_mime, image_data IS NOT NULL as has_image, created_at
		 FROM ticket_notes WHERE ticket_id = $1 ORDER BY created_at DESC`, ticketID,
	)
	if err != nil {
		apperror.WriteError(w, apperror.ErrInternal)
		return
	}
	defer rows.Close()

	var notes []Note
	for rows.Next() {
		var n Note
		if err := rows.Scan(&n.ID, &n.TicketID, &n.AuthorID, &n.Content, &n.ImageName, &n.ImageMime, &n.HasImage, &n.CreatedAt); err != nil {
			continue
		}
		notes = append(notes, n)
	}
	if notes == nil {
		notes = []Note{}
	}
	apperror.WriteJSON(w, http.StatusOK, notes)
}

// GetNoteImage serves the image for a note
func (h *AttachmentHandler) GetNoteImage(w http.ResponseWriter, r *http.Request) {
	noteID, err := uuid.Parse(chi.URLParam(r, "noteId"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}

	var imgData []byte
	var imgMime string
	err = h.db.QueryRow(r.Context(),
		`SELECT image_data, image_mime FROM ticket_notes WHERE id = $1 AND image_data IS NOT NULL`, noteID,
	).Scan(&imgData, &imgMime)
	if err != nil {
		apperror.WriteError(w, apperror.ErrNotFound)
		return
	}

	w.Header().Set("Content-Type", imgMime)
	w.WriteHeader(http.StatusOK)
	w.Write(imgData)
}
