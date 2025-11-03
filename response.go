package main
import (
	"net/http"
)

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
