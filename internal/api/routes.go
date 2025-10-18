package api

import (
	"net/http"
	"strings"
)

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/health", h.Health)

	mux.HandleFunc("/api/v1/history", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			if strings.HasPrefix(r.URL.Path, "/api/v1/history/") && r.URL.Path != "/api/v1/history/" {
				h.GetEntry(w, r)
			} else {
				h.GetHistory(w, r)
			}
		case http.MethodDelete:
			if strings.HasPrefix(r.URL.Path, "/api/v1/history/") && r.URL.Path != "/api/v1/history/" {
				h.DeleteEntry(w, r)
			} else {
				h.ClearHistory(w, r)
			}
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/v1/entries", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			h.CreateEntry(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/v1/search", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			h.Search(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/v1/stats", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			h.GetStats(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	return mux
}
