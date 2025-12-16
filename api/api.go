package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/queue"
	"github.com/Gleipnir-Technology/nidus-sync/userfile"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func apiAudioPost(w http.ResponseWriter, r *http.Request, u *models.User) {
	id := chi.URLParam(r, "uuid")
	noteUUID, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "Failed to decode the uuid", http.StatusBadRequest)
		return
	}

	var payload NoteAudioPayload
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read the payload", http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		log.Error().Err(err).Msg("Audio note POST JSON decode error")
		output, err := os.OpenFile("/tmp/request.body", os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			log.Info().Msg("Failed to open temp request.bady")
		}
		defer output.Close()
		output.Write(body)
		log.Info().Msg("Wrote request to /tmp/request.body")

		http.Error(w, "Failed to decode the payload", http.StatusBadRequest)
		return
	}
	if err := db.NoteAudioCreate(context.Background(), noteUUID, db.NoteAudio{}, u.ID); err != nil {
		render.Render(w, r, errRender(err))
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func apiAudioContentPost(w http.ResponseWriter, r *http.Request, u *models.User) {
	u_str := chi.URLParam(r, "uuid")
	audioUUID, err := uuid.Parse(u_str)
	if err != nil {
		http.Error(w, "Failed to parse image UUID", http.StatusBadRequest)
		return
	}
	err = userfile.AudioFileContentWrite(audioUUID, r.Body)
	if err != nil {
		log.Printf("Failed to write content file: %v", err)
		http.Error(w, "failed to write content file", http.StatusInternalServerError)
	}

	queue.EnqueueAudioJob(queue.AudioJob{AudioUUID: audioUUID})
	w.WriteHeader(http.StatusOK)
}

func apiClientIos(w http.ResponseWriter, r *http.Request, u *models.User) {
	query := db.NewGeoQuery()
	query.Limit = 0
	sources, err := db.MosquitoSourceQuery(&query)
	if err != nil {
		render.Render(w, r, errRender(err))
		return
	}
	requests, err := db.ServiceRequestQuery(&query)
	if err != nil {
		render.Render(w, r, errRender(err))
		return
	}
	traps, err := db.TrapDataQuery(&query)
	if err != nil {
		render.Render(w, r, errRender(err))
		return
	}

	response := NewResponseClientIos(sources, requests, traps)
	if err := render.Render(w, r, response); err != nil {
		render.Render(w, r, errRender(err))
		return
	}
}

func apiClientIosNotePut(w http.ResponseWriter, r *http.Request, u *models.User) {
	id := chi.URLParam(r, "uuid")
	noteUUID, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "Failed to decode the uuid", http.StatusBadRequest)
		return
	}
	var payload NidusNotePayload
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read the payload", http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		log.Error().Err(err).Msg("Note PUT JSON decode error")
		output, err := os.OpenFile("/tmp/request.body", os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			log.Info().Msg("Failed to open temp request.bady")
		}
		defer output.Close()
		output.Write(body)
		log.Info().Msg("Wrote request to /tmp/request.body")

		http.Error(w, "Failed to decode the payload", http.StatusBadRequest)
		return
	}
	if err := db.NoteUpdate(context.Background(), noteUUID, db.NidusNotePayload{}); err != nil {
		render.Render(w, r, errRender(err))
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func apiImagePost(w http.ResponseWriter, r *http.Request, u *models.User) {
	id := chi.URLParam(r, "uuid")
	noteUUID, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "Failed to decode the uuid", http.StatusBadRequest)
		return
	}

	var payload NoteImagePayload
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read the payload", http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		log.Error().Err(err).Msg("Image note POST JSON decode error")
		output, err := os.OpenFile("/tmp/request.body", os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			log.Info().Msg("Failed to open temp request.bady")
		}
		defer output.Close()
		output.Write(body)
		log.Info().Msg("Wrote request to /tmp/request.body")

		http.Error(w, "Failed to decode the payload", http.StatusBadRequest)
		return
	}
	err = db.NoteImageCreate(context.Background(), noteUUID, db.NoteImage{}, u.ID)
	if err != nil {
		render.Render(w, r, errRender(err))
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func apiImageContentPost(w http.ResponseWriter, r *http.Request, u *models.User) {
	u_str := chi.URLParam(r, "uuid")
	imageUUID, err := uuid.Parse(u_str)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse image UUID")
		http.Error(w, "Failed to parse image UUID", http.StatusBadRequest)
	}
	// Read first 8 bytes to check PNG signature
	filepath := fmt.Sprintf("%s/%s.photo", userfile.UserFilesDirectory, imageUUID.String())

	// Create file in configured directory
	dst, err := os.Create(filepath)
	if err != nil {
		log.Printf("Failed to create image file %s: %v", filepath, err)
		http.Error(w, "Unable to create file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy rest of request body to file
	_, err = io.Copy(dst, r.Body)
	if err != nil {
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	log.Printf("Saved image file %s\n", imageUUID)
	fmt.Fprintf(w, "PNG uploaded successfully to %s", filepath)
}

func apiMosquitoSource(w http.ResponseWriter, r *http.Request, u *models.User) {
	bounds, err := parseBounds(r)
	if err != nil {
		render.Render(w, r, errRender(err))
		return
	}

	query := db.NewGeoQuery()
	query.Bounds = *bounds
	query.Limit = 100
	sources, err := db.MosquitoSourceQuery(&query)
	if err != nil {
		render.Render(w, r, errRender(err))
		return
	}

	data := []render.Renderer{}
	for _, s := range sources {
		data = append(data, NewResponseMosquitoSource(s))
	}
	if err := render.RenderList(w, r, data); err != nil {
		render.Render(w, r, errRender(err))
	}
}

func apiTrapData(w http.ResponseWriter, r *http.Request, u *models.User) {
	bounds, err := parseBounds(r)
	if err != nil {
		render.Render(w, r, errRender(err))
		return
	}

	query := db.NewGeoQuery()
	query.Bounds = *bounds
	query.Limit = 100
	trap_data, err := db.TrapDataQuery(&query)
	if err != nil {
		render.Render(w, r, errRender(err))
		return
	}

	data := []render.Renderer{}
	for _, td := range trap_data {
		data = append(data, NewResponseTrapDatum(td))
	}
	if err := render.RenderList(w, r, data); err != nil {
		render.Render(w, r, errRender(err))
	}
}

func apiServiceRequest(w http.ResponseWriter, r *http.Request, u *models.User) {
	bounds, err := parseBounds(r)
	if err != nil {
		render.Render(w, r, errRender(err))
		return
	}
	query := db.NewGeoQuery()
	query.Bounds = *bounds
	query.Limit = 100
	requests, err := db.ServiceRequestQuery(&query)
	if err != nil {
		render.Render(w, r, errRender(err))
		return
	}

	data := []render.Renderer{}
	for _, sr := range requests {
		data = append(data, NewResponseServiceRequest(sr))
	}
	if err := render.RenderList(w, r, data); err != nil {
		render.Render(w, r, errRender(err))
	}
}

func parseBounds(r *http.Request) (*db.GeoBounds, error) {
	err := r.ParseForm()
	if err != nil {
		return nil, err
	}

	east := r.FormValue("east")
	north := r.FormValue("north")
	south := r.FormValue("south")
	west := r.FormValue("west")

	bounds := db.GeoBounds{}

	var temp float64
	temp, err = strconv.ParseFloat(east, 64)
	if err != nil {
		return nil, err
	}
	bounds.East = temp
	temp, err = strconv.ParseFloat(north, 64)
	if err != nil {
		return nil, err
	}
	bounds.North = temp
	temp, err = strconv.ParseFloat(south, 64)
	if err != nil {
		return nil, err
	}
	bounds.South = temp
	temp, err = strconv.ParseFloat(west, 64)
	if err != nil {
		return nil, err
	}
	bounds.West = temp
	return &bounds, nil
}

func errRender(err error) render.Renderer {
	log.Error().Err(err).Msg("Rendering error")
	return &ResponseErr{
		Error:          err,
		HTTPStatusCode: 500,
		StatusText:     "Error rendering response",
		ErrorText:      err.Error(),
	}
}

func webhookFieldseeker(w http.ResponseWriter, r *http.Request) {
	// Create or open the log file
	file, err := os.OpenFile("webhook/request.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Printf("Error opening log file: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Write timestamp
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Fprintf(file, "\n=== Request logged at %s ===\n", timestamp)

	// Write request line
	fmt.Fprintf(file, "%s %s %s\n", r.Method, r.RequestURI, r.Proto)

	// Write all headers
	fmt.Fprintf(file, "\nHeaders:\n")
	for name, values := range r.Header {
		for _, value := range values {
			fmt.Fprintf(file, "%s: %s\n", name, value)
		}
	}

	// Write body
	fmt.Fprintf(file, "\nBody:\n")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		fmt.Fprintf(file, "Error reading body: %v\n", err)
	} else {
		file.Write(body)
		if len(body) == 0 {
			fmt.Fprintf(file, "(empty body)")
		}
	}

	fmt.Fprintf(file, "\n=== End of request ===\n\n")

	// Extract the crc_token value for the signature portion

	// Respond with 204 No Content
	w.WriteHeader(http.StatusNoContent)
}

func parseTime(x string) time.Time {
	created_epoch, err := strconv.ParseInt(x, 10, 64)
	if err != nil {
		log.Error().Err(err).Msg("Unable to convert inspection timestamp")
	}
	created := time.UnixMilli(created_epoch)
	return created
}

