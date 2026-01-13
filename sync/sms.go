package sync

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

type SMSWebhookBody struct {
	Data SMSWebhookData `json:"data"`
}
type SMSWebhookData struct {
	ID         int64             `json:"id"`
	EventType  string            `json:"event_type"`
	RecordType string            `json:"record_type"`
	Payload    SMSMessagePayload `json:"payload"`
}

type SMSMessagePayload struct {
	ID         int64        `json:"id"`
	RecordType string       `json:"record_type"`
	From       SMSContact   `json:"from"`
	To         []SMSContact `json:"to"`
	Text       string       `json:"text"`
	ReceivedAt string       `json:"received_at"`
	Type       string       `json:"type"`
	Media      []MMSMedia   `json:"media"`
}

type MMSMedia struct {
	URL string `json:"url"`
}

// Contact represents a phone contact
type SMSContact struct {
	PhoneNumber string `json:"phone_number"`
	Status      string `json:"status,omitempty"`
}

func handleSMSMessage(data *SMSWebhookData) error {
	log.Info().Int64("ID", data.ID).Str("event_type", data.EventType).Str("record_type", data.RecordType).Str("from", data.Payload.From.PhoneNumber).Str("msg", data.Payload.Text).Str("receieved", data.Payload.ReceivedAt).Msg("Got SMS Message")

	for _, media := range data.Payload.Media {
		filePath, err := downloadMedia(media.URL)
		if err != nil {
			fmt.Errorf("Failed to download media from %s: %w", filePath, err)
			continue
		}
		fmt.Printf("Downloaded media to: %s\n", filePath)
	}
	return nil
}

// DownloadMedia downloads a media file from the given URL to a temporary location
// and returns the path to the downloaded file
func downloadMedia(mediaURL string) (string, error) {
	// Make GET request to the media URL
	resp, err := http.Get(mediaURL)
	if err != nil {
		return "", fmt.Errorf("failed to download media: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download media: status code %d", resp.StatusCode)
	}

	// Extract filename from URL or headers
	filename := getFilenameFromURL(mediaURL, resp)

	// Create temporary file with proper extension
	tmpDir := os.TempDir()
	timestamp := time.Now().UnixNano()
	tmpFilePath := filepath.Join(tmpDir, fmt.Sprintf("media_%d_%s", timestamp, filename))

	// Create the file
	out, err := os.Create(tmpFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer out.Close()

	// Write the response body to the file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to save media file: %w", err)
	}

	return tmpFilePath, nil
}

// getFilenameFromURL extracts filename from URL or Content-Disposition header
func getFilenameFromURL(mediaURL string, resp *http.Response) string {
	// First try Content-Disposition header
	contentDisp := resp.Header.Get("Content-Disposition")
	if contentDisp != "" {
		if strings.Contains(contentDisp, "filename=") {
			parts := strings.Split(contentDisp, "filename=")
			if len(parts) > 1 {
				filename := strings.Trim(parts[1], "\"' ")
				if filename != "" {
					return sanitizeFilename(filename)
				}
			}
		}
	}

	// Fall back to URL path
	urlPath := path.Base(mediaURL)
	if urlPath != "" && urlPath != "." && urlPath != "/" {
		return sanitizeFilename(urlPath)
	}

	// Default to generic name with extension based on Content-Type
	contentType := resp.Header.Get("Content-Type")
	ext := ".bin"

	switch {
	case strings.Contains(contentType, "image/jpeg"):
		ext = ".jpg"
	case strings.Contains(contentType, "image/png"):
		ext = ".png"
	case strings.Contains(contentType, "image/gif"):
		ext = ".gif"
	case strings.Contains(contentType, "video/mp4"):
		ext = ".mp4"
	case strings.Contains(contentType, "audio/mpeg"):
		ext = ".mp3"
	}

	return fmt.Sprintf("media%s", ext)
}

// sanitizeFilename removes potentially unsafe characters from filename
func sanitizeFilename(name string) string {
	// Replace unsafe characters with underscore
	unsafe := []string{"/", "\\", "?", "%", "*", ":", "|", "\"", "<", ">"}
	result := name
	for _, c := range unsafe {
		result = strings.ReplaceAll(result, c, "_")
	}
	return result
}
func postSMS(w http.ResponseWriter, r *http.Request) {
	// Log all request headers
	for name, values := range r.Header {
		for _, value := range values {
			log.Info().Str("name", name).Str("value", value).Msg("header")
		}
	}

	// Read the request body
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		//return nil, fmt.Errorf("failed to read request body: %w", err)
		respondError(w, "Failed to read request body", err, http.StatusInternalServerError)
		return
	}
	log.Info().Str("body", string(bodyBytes)).Msg("body")
	// Close the original body
	defer r.Body.Close()

	// Parse JSON into webhook struct
	var body SMSWebhookBody
	if err := json.Unmarshal(bodyBytes, &body); err != nil {
		respondError(w, "Failed to parse JSON", err, http.StatusBadRequest)
		return
	}

	if err := handleSMSMessage(&body.Data); err != nil {
		log.Error().Err(err).Msg("Failed to handle SMS Message")
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
func getSMS(w http.ResponseWriter, r *http.Request) {
	org := chi.URLParam(r, "org")

	to := r.URL.Query().Get("error")
	from := r.URL.Query().Get("error")
	message := r.URL.Query().Get("error")
	files := r.URL.Query().Get("error")
	id := r.URL.Query().Get("error")
	date := r.URL.Query().Get("error")

	log.Info().Str("org", org).Str("to", to).Str("from", from).Str("message", message).Str("files", files).Str("id", id).Str("date", date).Msg("Got SMS Message")
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-type", "text/plain")
	// Signifies to Voip.ms that the callback worked.
	fmt.Fprintf(w, "ok")
}
