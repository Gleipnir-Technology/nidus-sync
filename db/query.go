package db

import (
	"context"

	"github.com/google/uuid"
)
type NidusNotePayload struct {}
type NoteAudio struct {
	Transcription string
	Version int
	UUID uuid.UUID
}
type NoteImage struct {}
type MosquitoSource struct { }
type MosquitoTreatment struct { }
type ServiceRequest struct { }
type TrapData struct { }

func MosquitoSourceQuery() ([]MosquitoSource, error) {
	return make([]MosquitoSource, 0), nil
}
func NoteAudioCreate(ctx context.Context, noteUUID uuid.UUID, payload NoteAudio, userID int32) error {
	return nil
}
func NoteAudioGetLatest(ctx context.Context, uuid string) (*NoteAudio, error) {
	return nil, nil
}
func NoteAudioNormalized(uuid string) error {
	return nil
}
func NoteAudioTranscodedToOgg(uuid string) error {
	return nil
}
func NoteImageCreate(ctx context.Context, noteUUID uuid.UUID, payload NoteImage, userID int32) error {
	return nil
}
func NoteUpdate(ctx context.Context, noteUUID uuid.UUID, payload NidusNotePayload) error {
	return nil
}
func ServiceRequestQuery() ([]ServiceRequest, error) {
	return make([]ServiceRequest, 0), nil
}
func TrapDataQuery() ([]TrapData, error) {
	return make([]TrapData, 0), nil
}
