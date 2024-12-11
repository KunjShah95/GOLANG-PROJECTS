package middlewares

import (
	"compress/gzip"
	"log"
	"net/http"
	"strings"
	"sync"
)

var gzipPool = sync.Pool{
	New: func() interface{} {
		// Initialize gzip writer with a placeholder writer; it will be reset per request
		writer, err := gzip.NewWriterLevel(nil, gzip.BestSpeed)
		if err != nil {
			log.Fatalf("Failed to initialize gzip writer: %v", err)
		}
		return writer
	},
}

// gzipResponseWriter wraps http.ResponseWriter to provide gzip compression
type gzipResponseWriter struct {
	http.ResponseWriter
	Writer *gzip.Writer
}

// Header returns the header map that will be sent by WriteHeader.
func (w *gzipResponseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

// WriteHeader sends an HTTP response header with the provided status code.
func (w *gzipResponseWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
}

// Write compresses the data before writing it to the client.
func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// Flush ensures that all compressed data is sent to the client.
func (w *gzipResponseWriter) Flush() {
	w.Writer.Flush()
	if flusher, ok := w.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

// GzipCompression middleware compresses responses if the client supports it
func GzipCompression(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the client accepts gzip encoding
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		// Remove Content-Length header and set Content-Encoding to gzip
		w.Header().Del("Content-Length")
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Content-Length", "0")

		// Acquire a gzip writer from the pool
		gzipWriter := gzipPool.Get().(*gzip.Writer)
		// Initialize the writer with the response writer
		gzipWriter.Reset(w)

		defer func() {
			if err := gzipWriter.Close(); err != nil {
				log.Printf("Failed to close gzip writer: %v", err)
			}
			gzipPool.Put(gzipWriter)
		}()

		// Wrap the ResponseWriter with gzipResponseWriter
		gzWriter := &gzipResponseWriter{
			ResponseWriter: w,
			Writer:         gzipWriter,
		}

		// Serve the next handler with the gzipResponseWriter
		next.ServeHTTP(gzWriter, r)
	})
}
