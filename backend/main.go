package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	dbPath := os.Getenv("HIVE_DB_PATH")
	if dbPath == "" {
		dbPath = filepath.Join("..", "hive.db")
	}
	addr := os.Getenv("HIVE_API_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	if err := initDB(dbPath); err != nil {
		log.Fatalf("Database init failed: %v", err)
	}

	mux := http.NewServeMux()

	// Runs
	mux.HandleFunc("GET /api/runs", listRuns)
	mux.HandleFunc("POST /api/runs", createRun)
	mux.HandleFunc("GET /api/runs/{id}", getRun)

	// Probes
	mux.HandleFunc("GET /api/probes", listProbes)
	mux.HandleFunc("POST /api/probes", createProbe)

	// Gap matrices
	mux.HandleFunc("GET /api/gap-matrices", listGapMatrices)
	mux.HandleFunc("GET /api/gap-matrices/{id}", getGapMatrix)
	mux.HandleFunc("POST /api/gap-matrices", createGapMatrix)

	// Comparisons
	mux.HandleFunc("GET /api/comparisons", listComparisons)
	mux.HandleFunc("POST /api/comparisons", createComparison)

	// Stats
	mux.HandleFunc("GET /api/stats", getStats)

	// Import from error_ledger
	mux.HandleFunc("POST /api/import", handleImport)

	// Health
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	handler := enableCORS(mux)
	log.Printf("Hive API server starting on %s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
