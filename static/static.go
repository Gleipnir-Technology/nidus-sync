package static

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

//go:embed css gen file ico img js vendor
var embeddedStaticFS embed.FS

// fileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
var startedTime time.Time = time.Now()

var localFS http.Dir

func AddStaticRoute(r *mux.Router, path string) {
	if localFS == "" {
		localFS = http.Dir("./static")
		// Useful for debugging embedded file issues
		if config.IsProductionEnvironment() {
			fs.WalkDir(embeddedStaticFS, ".", func(path string, d fs.DirEntry, err error) error {
				log.Debug().Str("path", path).Send()
				return nil
			})
		}
	}
	fileServer(r, "/static", localFS, embeddedStaticFS)
}

func fileServer(r *mux.Router, path string, root http.FileSystem, embeddedFS embed.FS) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.HandleFunc(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		//rctx := chi.RouteContext(r.Context())

		//pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		pathPrefix := strings.TrimPrefix(r.URL.Path, "/static")

		// Determine the actual file path
		requestedPath := strings.TrimPrefix(r.URL.Path, pathPrefix+"/")

		var err error
		var fileToServe http.File
		found := false

		// For dev, try the current filesystem
		if !config.IsProductionEnvironment() {
			// Try to open from local filesystem for development
			fileToServe, err = root.Open(requestedPath)
			if err != nil {
				log.Warn().Str("path", requestedPath).Msg("Failed to read static file for dev")
				found = false
			} else {
				found = true
			}
		}
		// For production use the embedded filesystem
		if !found {
			// Requested paths start with
			embeddedFile, err := embeddedFS.Open(requestedPath)

			if err != nil {
				log.Debug().Err(err).Str("requested path", requestedPath).Msg("Failed to find resource")
				http.NotFound(w, r)
				return
			}

			// Wrap the embedded file to implement http.File interface
			fileToServe = &embeddedFileWrapper{embeddedFile}
		}

		// Create a custom ResponseWriter that allows us to modify headers
		crw := &customResponseWriter{ResponseWriter: w}

		// Add caching headers
		if config.IsProductionEnvironment() {
			ext := filepath.Ext(requestedPath)
			switch ext {
			case ".css", ".jpg", ".jpeg", ".png", ".gif", ".svg", ".woff", ".woff2", ".ttf":
				// Cache for 1 week (604800 seconds)
				crw.Header().Set("Cache-Control", "public, max-age=604800, stale-while-revalidate=86400")
			default:
				// If it's a generated file, cache it essentially forever (1 year)
				if strings.HasPrefix(requestedPath, "gen/") {
					crw.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
				} else {
					// Other files, 1 hour
					crw.Header().Set("Cache-Control", "public, max-age=3600")
				}
			}
		}
		// Serve the file
		http.ServeContent(crw, r, requestedPath, startedTime, fileToServe)

		// Close the file
		fileToServe.Close()
	})
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
