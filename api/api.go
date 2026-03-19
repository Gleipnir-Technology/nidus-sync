package api

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
)

func handleClientIos(w http.ResponseWriter, r *http.Request, u platform.User) {
	var sinceStr string
	err := r.ParseForm()
	if err != nil {
		render.Render(w, r, errRender(fmt.Errorf("Failed to parse GET form: %w", err)))
		return
	} else {
		sinceStr = r.FormValue("since")
	}

	var since *time.Time
	if sinceStr == "" {
		since = nil
	} else {
		since, err = parseTime(sinceStr)
		if err != nil {
			render.Render(w, r, errRender(fmt.Errorf("Failed to parse 'since' value: %w", err)))
			return
		}
	}

	csync, err := platform.ContentClientIos(r.Context(), u, since)
	if err != nil {
		render.Render(w, r, errRender(err))
		return
	}

	var since_used time.Time
	if since == nil {
		since_used = time.Unix(0, 0)
	} else {
		since_used = *since
	}
	response := ResponseClientIos{
		Fieldseeker: toResponseFieldseeker(csync.Fieldseeker),
		Since:       since_used,
	}
	if err := render.Render(w, r, response); err != nil {
		render.Render(w, r, errRender(err))
		return
	}
}

func apiMosquitoSource(w http.ResponseWriter, r *http.Request, u platform.User) {
	bounds, err := parseBounds(r)
	if err != nil {
		render.Render(w, r, errRender(err))
		return
	}

	query := db.NewGeoQuery()
	query.Bounds = *bounds
	query.Limit = 100
	sources, err := platform.MosquitoSourceQuery()
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

func apiTrapData(w http.ResponseWriter, r *http.Request, u platform.User) {
	bounds, err := parseBounds(r)
	if err != nil {
		render.Render(w, r, errRender(err))
		return
	}

	query := db.NewGeoQuery()
	query.Bounds = *bounds
	query.Limit = 100
	trap_data, err := platform.TrapDataQuery()
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

func apiServiceRequest(w http.ResponseWriter, r *http.Request, u platform.User) {
	bounds, err := parseBounds(r)
	if err != nil {
		render.Render(w, r, errRender(err))
		return
	}
	query := db.NewGeoQuery()
	query.Bounds = *bounds
	query.Limit = 100
	requests, err := platform.ServiceRequestQuery()
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

func parseTime(x string) (*time.Time, error) {
	created_epoch, err := strconv.ParseInt(x, 10, 64)
	if err != nil {
		return &time.Time{}, fmt.Errorf("Failed to parse time '%s': %w", x, err)
	}
	created := time.UnixMilli(created_epoch)
	return &created, nil
}
