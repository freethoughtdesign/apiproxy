package main

import (
	"compress/gzip"
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func gzipHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()

		gzr := gzipResponseWriter{Writer: gz, ResponseWriter: w}
		next.ServeHTTP(gzr, r)
	})
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8787"
	}

	listen := os.Getenv("LISTEN")
	if listen == "" {
		listen = "0.0.0.0"
	}

	api_base := os.Getenv("API_BASE")
	if api_base == "" {
		log.Fatal("Error loading API_BASE env var.")
	}

	authn := os.Getenv("AUTHORIZATION")
	if authn == "" {
		log.Fatal("Error loading AUTHORIZATION env var.")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Create a custom transport with system root CAs
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{RootCAs: nil}, // Use system root CAs
		}
		client := &http.Client{Transport: tr}
		query := ""
		if r.URL.RawQuery != "" {
			query = "?" + r.URL.RawQuery
		}
		url := api_base + r.URL.Path + query
		req, err := http.NewRequest(r.Method, url, r.Body)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		req.Header.Set("Content-Type", r.Header.Get("Content-Type"))
		req.Header.Set("Authorization", authn)
		resp, err := client.Do(req)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("%s request for %s: %s", r.Method, r.URL.Path+query, resp.Status)

		// TODO: Consider caching responses here.

		w.Header().Set("X-Generator", "API Proxy for "+api_base)
		w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
		w.Header().Set("Date", resp.Header.Get("Date"))
		w.Header().Set("Etag", resp.Header.Get("Etag"))
		// w.Header().Set("Cache-Control", resp.Header.Get("Cache-Control"))

		// Set some relaxed CORS headers to allow use in browsers.
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Channel")

		w.WriteHeader(resp.StatusCode)
		w.Write(body)
	})

	log.Printf("Proxying %s on port %s", api_base, port)
	log.Fatal(http.ListenAndServe(listen+":"+port, gzipHandler(mux)))
}
