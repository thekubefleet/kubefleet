package server

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
)

type HTTPServer struct {
	dataStore *DataStore
	router    *mux.Router
}

func NewHTTPServer(dataStore *DataStore) *HTTPServer {
	server := &HTTPServer{
		dataStore: dataStore,
		router:    mux.NewRouter(),
	}

	// API routes
	server.router.HandleFunc("/api/data", server.handleGetData).Methods("GET")
	server.router.HandleFunc("/api/data/latest", server.handleGetLatestData).Methods("GET")
	server.router.HandleFunc("/api/health", server.handleHealth).Methods("GET")

	// Serve React app
	server.router.PathPrefix("/").HandlerFunc(server.handleReactApp)

	return server
}

func (s *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *HTTPServer) handleGetData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	data := s.dataStore.GetAllData()
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":  data,
		"count": len(data),
	})
}

func (s *HTTPServer) handleGetLatestData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	data := s.dataStore.GetLatestData()
	if data == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "No data available"})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": data,
	})
}

func (s *HTTPServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":     "healthy",
		"dataPoints": s.dataStore.GetDataCount(),
	})
}

func (s *HTTPServer) handleReactApp(w http.ResponseWriter, r *http.Request) {
	// If the path is for an API endpoint, don't serve the React app
	if strings.HasPrefix(r.URL.Path, "/api/") {
		http.NotFound(w, r)
		return
	}

	// Check if we're in development mode (no build directory)
	if _, err := os.Stat("build"); os.IsNotExist(err) {
		// Development mode - serve a simple HTML page that loads from localhost:3001
		if r.URL.Path == "/" {
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(`
<!DOCTYPE html>
<html>
<head>
    <title>KubeFleet Dashboard</title>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <script src="http://localhost:3001/static/js/bundle.js"></script>
    <link rel="stylesheet" href="http://localhost:3001/static/css/main.css">
</head>
<body>
    <div id="root"></div>
</body>
</html>
			`))
			return
		}
		http.NotFound(w, r)
		return
	}

	// Production mode - serve from build directory
	path := r.URL.Path
	if path == "/" {
		path = "/index.html"
	}

	// Try to serve the file from build directory
	filePath := filepath.Join("build", path)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// If file doesn't exist, serve index.html for SPA routing
		filePath = filepath.Join("build", "index.html")
	}

	http.ServeFile(w, r, filePath)
}
