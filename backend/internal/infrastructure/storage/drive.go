package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// DriveStorage handles file operations with Google Drive API
type DriveStorage struct {
	folderID     string
	credentials  *serviceAccountCredentials
	cachedToken  *accessToken
	httpClient   *http.Client
}

type serviceAccountCredentials struct {
	Type                    string `json:"type"`
	ProjectID               string `json:"project_id"`
	PrivateKeyID            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientID                string `json:"client_id"`
	AuthURI                 string `json:"auth_uri"`
	TokenURI                string `json:"token_uri"`
	AuthProviderX509CertURL string `json:"auth_provider_x509_cert_url"`
	ClientX509CertURL       string `json:"client_x509_cert_url"`
}

type accessToken struct {
	Token     string
	ExpiresAt time.Time
}

type DriveFile struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	MimeType string `json:"mimeType"`
	WebViewLink string `json:"webViewLink"`
}

// NewDriveStorage creates a new Google Drive storage instance
func NewDriveStorage() (*DriveStorage, error) {
	folderID := os.Getenv("GOOGLE_DRIVE_FOLDER_ID")
	if folderID == "" {
		return nil, fmt.Errorf("GOOGLE_DRIVE_FOLDER_ID is required")
	}

	credPath := os.Getenv("GOOGLE_SERVICE_ACCOUNT_KEY_PATH")
	credJSON := os.Getenv("GOOGLE_SERVICE_ACCOUNT_KEY_JSON")

	var creds serviceAccountCredentials

	if credJSON != "" {
		// Use inline JSON credentials
		if err := json.Unmarshal([]byte(credJSON), &creds); err != nil {
			return nil, fmt.Errorf("failed to parse GOOGLE_SERVICE_ACCOUNT_KEY_JSON: %w", err)
		}
	} else if credPath != "" {
		// Use file-based credentials
		data, err := os.ReadFile(credPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read service account key file: %w", err)
		}
		if err := json.Unmarshal(data, &creds); err != nil {
			return nil, fmt.Errorf("failed to parse service account key: %w", err)
		}
	} else {
		return nil, fmt.Errorf("either GOOGLE_SERVICE_ACCOUNT_KEY_PATH or GOOGLE_SERVICE_ACCOUNT_KEY_JSON is required")
	}

	return &DriveStorage{
		folderID:    folderID,
		credentials: &creds,
		httpClient:  &http.Client{Timeout: 60 * time.Second},
	}, nil
}

// Upload uploads a file to Google Drive and returns the file ID
func (d *DriveStorage) Upload(ctx context.Context, filename string, mimeType string, content []byte) (string, error) {
	token, err := d.getAccessToken(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get access token: %w", err)
	}

	// Create multipart upload request
	boundary := "----DriveUploadBoundary"
	metadata := map[string]interface{}{
		"name":    filename,
		"parents": []string{d.folderID},
	}
	metaJSON, _ := json.Marshal(metadata)

	var body bytes.Buffer
	body.WriteString("--" + boundary + "\r\n")
	body.WriteString("Content-Type: application/json; charset=UTF-8\r\n\r\n")
	body.Write(metaJSON)
	body.WriteString("\r\n--" + boundary + "\r\n")
	body.WriteString("Content-Type: " + mimeType + "\r\n\r\n")
	body.Write(content)
	body.WriteString("\r\n--" + boundary + "--\r\n")

	req, err := http.NewRequestWithContext(ctx, "POST",
		"https://www.googleapis.com/upload/drive/v3/files?uploadType=multipart&fields=id,name,mimeType,webViewLink",
		&body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "multipart/related; boundary="+boundary)

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("drive upload request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("drive upload failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	var driveFile DriveFile
	if err := json.NewDecoder(resp.Body).Decode(&driveFile); err != nil {
		return "", fmt.Errorf("failed to decode drive response: %w", err)
	}

	return driveFile.ID, nil
}

// Download downloads a file from Google Drive by file ID
func (d *DriveStorage) Download(ctx context.Context, fileID string) ([]byte, error) {
	token, err := d.getAccessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "GET",
		"https://www.googleapis.com/drive/v3/files/"+fileID+"?alt=media", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("drive download request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("drive download failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	return io.ReadAll(resp.Body)
}

// Delete removes a file from Google Drive
func (d *DriveStorage) Delete(ctx context.Context, fileID string) error {
	token, err := d.getAccessToken(ctx)
	if err != nil {
		return fmt.Errorf("failed to get access token: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "DELETE",
		"https://www.googleapis.com/drive/v3/files/"+fileID, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("drive delete request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("drive delete failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

// getAccessToken returns a valid access token, refreshing if needed
func (d *DriveStorage) getAccessToken(ctx context.Context) (string, error) {
	if d.cachedToken != nil && time.Now().Before(d.cachedToken.ExpiresAt) {
		return d.cachedToken.Token, nil
	}

	// Create JWT for service account
	now := time.Now()
	claims := map[string]interface{}{
		"iss":   d.credentials.ClientEmail,
		"scope": "https://www.googleapis.com/auth/drive.file",
		"aud":   d.credentials.TokenURI,
		"iat":   now.Unix(),
		"exp":   now.Add(time.Hour).Unix(),
	}

	jwt, err := d.signJWT(claims)
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT: %w", err)
	}

	// Exchange JWT for access token
	tokenReq := strings.NewReader("grant_type=urn%3Aietf%3Aparams%3Aoauth%3Agrant-type%3Ajwt-bearer&assertion=" + jwt)
	req, err := http.NewRequestWithContext(ctx, "POST", d.credentials.TokenURI, tokenReq)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("token request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("token exchange failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("failed to decode token response: %w", err)
	}

	d.cachedToken = &accessToken{
		Token:     tokenResp.AccessToken,
		ExpiresAt: now.Add(time.Duration(tokenResp.ExpiresIn-60) * time.Second), // refresh 60s before expiry
	}

	return d.cachedToken.Token, nil
}
