package api

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
	"github.com/Gleipnir-Technology/nidus-sync/resource"
	//"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

/*
type renderer struct {
}
func (ren *renderer) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
*/
// In the best case scenario, the excellent github.com/pkg/errors package
// helps reveal information on the error, setting it on Err, and in the Render()
// method, using it to set the application-specific error code in AppCode.
type ResponseErr struct {
	Error          error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

func (e *ResponseErr) Render(w http.ResponseWriter, r *http.Request) error {
	http.Error(w, e.StatusText, e.HTTPStatusCode)
	return nil
}

func errRender(err error) *ResponseErr {
	log.Error().Err(err).Msg("Rendering error")
	return &ResponseErr{
		Error:          err,
		HTTPStatusCode: 500,
		StatusText:     "Error rendering response",
		ErrorText:      err.Error(),
	}
}

type Renderable interface {
	Render(http.ResponseWriter, *http.Request) error
}

func renderShim(w http.ResponseWriter, r *http.Request, renderer Renderable) error {
	return renderer.Render(w, r)
}
func renderList(w http.ResponseWriter, r *http.Request, data []Renderable) error {
	return nil
}
func handleClientIos(w http.ResponseWriter, r *http.Request, u platform.User) {
	var sinceStr string
	err := r.ParseForm()
	if err != nil {
		renderShim(w, r, errRender(fmt.Errorf("Failed to parse GET form: %w", err)))
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
			renderShim(w, r, errRender(fmt.Errorf("Failed to parse 'since' value: %w", err)))
			return
		}
	}

	csync, err := platform.ContentClientIos(r.Context(), u, since)
	if err != nil {
		renderShim(w, r, errRender(err))
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
	if err := renderShim(w, r, response); err != nil {
		renderShim(w, r, errRender(err))
		return
	}
}

func apiMosquitoSource(w http.ResponseWriter, r *http.Request, u platform.User) {
	bounds, err := parseBounds(r)
	if err != nil {
		renderShim(w, r, errRender(err))
		return
	}

	query := db.NewGeoQuery()
	query.Bounds = *bounds
	query.Limit = 100
	sources, err := platform.MosquitoSourceQuery()
	if err != nil {
		renderShim(w, r, errRender(err))
		return
	}

	data := []Renderable{}
	for _, s := range sources {
		data = append(data, NewResponseMosquitoSource(s))
	}
	if err := renderList(w, r, data); err != nil {
		renderShim(w, r, errRender(err))
	}
}

func apiTrapData(w http.ResponseWriter, r *http.Request, u platform.User) {
	bounds, err := parseBounds(r)
	if err != nil {
		renderShim(w, r, errRender(err))
		return
	}

	query := db.NewGeoQuery()
	query.Bounds = *bounds
	query.Limit = 100
	trap_data, err := platform.TrapDataQuery()
	if err != nil {
		renderShim(w, r, errRender(err))
		return
	}

	data := []Renderable{}
	for _, td := range trap_data {
		data = append(data, NewResponseTrapDatum(td))
	}
	if err := renderList(w, r, data); err != nil {
		renderShim(w, r, errRender(err))
	}
}

func apiServiceRequest(w http.ResponseWriter, r *http.Request, u platform.User) {
	bounds, err := parseBounds(r)
	if err != nil {
		renderShim(w, r, errRender(err))
		return
	}
	query := db.NewGeoQuery()
	query.Bounds = *bounds
	query.Limit = 100
	requests, err := platform.ServiceRequestQuery()
	if err != nil {
		renderShim(w, r, errRender(err))
		return
	}

	data := []Renderable{}
	for _, sr := range requests {
		data = append(data, types.ServiceRequestFromModel(sr))
	}
	if err := renderList(w, r, data); err != nil {
		renderShim(w, r, errRender(err))
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

type about struct {
	Environment string `json:"environment"`
	SentryDSN   string `json:"sentry_dsn"`
	Version     string `json:"version"`
}

func getRoot(ctx context.Context, r *http.Request, q resource.QueryParams) (*about, *nhttp.ErrorWithStatus) {
	return &about{
		Environment: config.Environment,
		SentryDSN:   config.SentryDSNFrontend,
		Version:     version,
	}, nil
}
