package static

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/lint"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

//go:embed css gen file ico img js vendor
var embeddedStaticFS embed.FS

// fileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
var startedTime time.Time = time.Now()

var localFS = http.Dir("./static")

func AddStaticRoute(r *mux.Router, path string) {
	fileServer(r, "/static/", localFS, embeddedStaticFS)
}

func SinglePageApp(gen_path string) http.Handler {
	// Accept the path as relative from project root, but
	// fix it to actually be relative to static filesystem root
	path := strings.TrimPrefix(gen_path, "static/")
	return spaHandler{
		genRoot: path,
	}

}

type spaHandler struct {
	genRoot string
}

func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	request_path := r.URL.Path
	path := h.genRoot + request_path
	fileToServe, err := fileFromFilesystem(path)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// default to index file
		fileToServe, err = fileFromFilesystem(h.genRoot + "/index.html")

		if err != nil {
			log.Error().Err(err).Msg("failed to open embedded index file")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	serveFileMaybeEmbedded(w, r, *fileToServe, path)
}

func fileServer(r *mux.Router, path string, root http.FileSystem, embeddedFS embed.FS) {
	log.Debug().Str("path", path).Msg("adding file server")
	r.PathPrefix(path).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/static/")
		fileToServe, err := fileFromFilesystem(path)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				http.NotFound(w, r)
			}
		}
		serveFileMaybeEmbedded(w, r, *fileToServe, path)
	})
}

func fileFromFilesystem(path string) (*http.File, error) {
	var err error
	var fileToServe http.File
	found := false

	// For dev, try the current filesystem
	if !config.IsProductionEnvironment() {
		// Try to open from local filesystem for development
		fileToServe, err = localFS.Open(path)
		if err != nil {
			//log.Warn().Err(err).Str("path", path).Msg("Failed to read static file for dev")
			found = false
		} else {
			found = true
		}
	}
	// For production use the embedded filesystem
	if !found {
		// Requested paths start with
		embeddedFile, err := embeddedStaticFS.Open(path)

		if err != nil {
			return nil, fmt.Errorf("open embedded file: %w", err)
		}

		// Wrap the embedded file to implement http.File interface
		fileToServe = &embeddedFileWrapper{embeddedFile}
	}
	return &fileToServe, nil
}

// Serve a file from the filesystem if we're in development mode or from the
// embedded filesystem if we aren't
func serveFileMaybeEmbedded(w http.ResponseWriter, r *http.Request, fileToServe http.File, path string) {
	// Create a custom ResponseWriter that allows us to modify headers
	crw := &customResponseWriter{ResponseWriter: w}

	// Add caching headers
	if config.IsProductionEnvironment() {
		ext := filepath.Ext(path)
		switch ext {
		case ".css", ".jpg", ".jpeg", ".png", ".gif", ".svg", ".woff", ".woff2", ".ttf":
			// Cache for 1 week (604800 seconds)
			crw.Header().Set("Cache-Control", "public, max-age=604800, stale-while-revalidate=86400")
		default:
			// If it's a generated file, cache it essentially forever (1 year)
			if strings.HasPrefix(path, "/static/gen/") {
				crw.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
			} else {
				// Other files, 1 hour
				crw.Header().Set("Cache-Control", "public, max-age=3600")
			}
		}
	}
	// Serve the file
	http.ServeContent(crw, r, path, startedTime, fileToServe)

	// Close the file
	lint.LogOnErr(fileToServe.Close, "close static file")
}

type embeddedFileWrapper struct {
	file fs.File
}

func (e *embeddedFileWrapper) Close() error {
	return e.file.Close()
}

func (e *embeddedFileWrapper) Read(p []byte) (n int, err error) {
	return e.file.Read(p)
}

type Seeker interface {
	Seek(offset int64, whence int) (int64, error)
}

func (e *embeddedFileWrapper) Seek(offset int64, whence int) (int64, error) {
	if seeker, ok := e.file.(Seeker); ok {
		return seeker.Seek(offset, whence)
	}
	return 0, fmt.Errorf("Seek not supported")
}

func (e *embeddedFileWrapper) Readdir(count int) ([]os.FileInfo, error) {
	// This is a bit tricky with embedded files
	if dirFile, ok := e.file.(fs.ReadDirFile); ok {
		entries, err := dirFile.ReadDir(count)
		if err != nil {
			return nil, err
		}

		fileInfos := make([]os.FileInfo, len(entries))
		for i, entry := range entries {
			fileInfos[i], err = entry.Info()
			if err != nil {
				return nil, err
			}
		}
		return fileInfos, nil
	}
	return nil, fmt.Errorf("Readdir not supported")
}

func (e *embeddedFileWrapper) Stat() (os.FileInfo, error) {
	return e.file.Stat()
}

// Custom ResponseWriter to track Content-Type
type customResponseWriter struct {
	http.ResponseWriter
	contentType string
	wroteHeader bool
}

func (crw *customResponseWriter) WriteHeader(code int) {
	crw.wroteHeader = true
	crw.ResponseWriter.WriteHeader(code)
}

func (crw *customResponseWriter) Header() http.Header {
	return crw.ResponseWriter.Header()
}

func (crw *customResponseWriter) Write(b []byte) (int, error) {
	if !crw.wroteHeader {
		if crw.contentType == "" {
			crw.contentType = http.DetectContentType(b)
			crw.ResponseWriter.Header().Set("Content-Type", crw.contentType)
		}
		crw.WriteHeader(http.StatusOK)
	}
	return crw.ResponseWriter.Write(b)
}
