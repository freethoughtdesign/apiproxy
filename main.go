package main

import (
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8787"
	}

	listen := os.Getenv("LISTEN")
	if listen == "" {
		listen = "0.0.0.0"
	}

	hostname := os.Getenv("API_HOSTNAME")
	if hostname == "" {
		log.Fatal("Error loading API_HOSTNAME env var.")
	}

	token := os.Getenv("ACCESS_TOKEN")
	if token == "" {
		log.Fatal("Error loading ACCESS_TOKEN env var.")
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Create a custom transport with system root CAs
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{RootCAs: nil}, // Use system root CAs
		}
		client := &http.Client{Transport: tr}
		url := "https://" + hostname + r.URL.Path
		req, err := http.NewRequest(r.Method, url, r.Body)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		req.Header.Set("Content-Type", r.Header.Get("Content-Type"))
		req.Header.Set("Authorization", "Bearer "+token)
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

		log.Printf("%s request for %s: %s", r.Method, r.URL.Path, resp.Status)

		// TODO: Consider caching responses here.

		w.Header().Set("X-Generator", "API Proxy for "+hostname)
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

	log.Printf("Proxying %s on port %s", hostname, port)
	log.Fatal(http.ListenAndServe(listen+":"+port, nil))
}
