package html

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

// fileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func fileServer(r chi.Router, path string, root http.FileSystem, embeddedFS embed.FS, embeddedPath string) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")

		// Determine the actual file path
		requestedPath := strings.TrimPrefix(r.URL.Path, pathPrefix)

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
			embeddedFilePath := filepath.Join(embeddedPath, requestedPath)
			embeddedFile, err := embeddedFS.Open(embeddedFilePath)

			if err != nil {
				http.NotFound(w, r)
				return
			}

			// Wrap the embedded file to implement http.File interface
			fileToServe = &embeddedFileWrapper{embeddedFile}
		}

		// Create a custom ResponseWriter that allows us to modify headers
		crw := &customResponseWriter{ResponseWriter: w}

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
