package main

import (
	"compress/gzip"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

// enableCORS adds CORS headers to the response
func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Allow requests from any origin
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	}
}

type ORTBPayload struct {
	ID     string `json:"id"`
	At     int    `json:"at"`
	Tmax   int    `json:"tmax"`
	Imp    []struct {
		ID       string `json:"id"`
		Banner   map[string]interface{} `json:"banner"`
		Tagid    string  `json:"tagid"`
		Bidfloor float64 `json:"bidfloor"`
		Bidfloorcur string `json:"bidfloorcur"`
	} `json:"imp"`
	Site   map[string]interface{} `json:"site"`
	Device map[string]interface{} `json:"device"`
	User   map[string]string      `json:"user"`
}

type ResponsePayload struct {
	Success     bool                   `json:"success"`
	Message     string                 `json:"message"`
	RequestID   string                 `json:"request_id,omitempty"`
	QueryParams map[string]string      `json:"query_params,omitempty"`
	Data        map[string]interface{} `json:"data,omitempty"`
}

func main() {
	http.HandleFunc("/api", enableCORS(handleRequest))
	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read query parameters
	queryParams := make(map[string]string)
	for key, values := range r.URL.Query() {
		queryParams[key] = values[0]
	}

	// Read body
	var body []byte
	var err error

	// Check if the request is gzipped
	if r.URL.Query().Get("gzip") == "1" {
		// Handle gzipped request body
		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			http.Error(w, "Error creating gzip reader: "+err.Error(), http.StatusBadRequest)
			return
		}
		defer gz.Close()

		body, err = io.ReadAll(gz)
		if err != nil {
			http.Error(w, "Error reading gzipped request body: "+err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		// Handle regular request body
		body, err = io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}
	}

	// Parse JSON from the body (gzipped or not)
	var ortbPayload ORTBPayload
	if err := json.Unmarshal(body, &ortbPayload); err != nil {
		http.Error(w, "Error parsing ORTB payload: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Prepare response
	response := ResponsePayload{
		Success:     true,
		Message:     "Request processed successfully",
		RequestID:   ortbPayload.ID,
		QueryParams: queryParams,
		Data: map[string]interface{}{
			"imp_count": len(ortbPayload.Imp),
			"site":      ortbPayload.Site,
			"device":    ortbPayload.Device,
		},
	}

	// Set headers for gzipped JSON response
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Encoding", "gzip")

	// Create gzip writer
	gz := gzip.NewWriter(w)
	defer gz.Close()

	// Write JSON response through gzip writer
	if err := json.NewEncoder(gz).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}
