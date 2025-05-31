package main

import (
	"agroport/handlers"
	"agroport/models"
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	// Get database URL from environment
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	// Connect to database
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	log.Println("Database connection established")

	// Initialize models with database connection
	models.InitDB(db)

	// Run migrations
	if err := models.RunMigrations(); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Initialize handlers
	h := handlers.NewHandler(db)

	// Setup routes
	r := mux.NewRouter()
	setupRoutes(r, h)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func setupRoutes(r *mux.Router, h *handlers.Handler) {
	// API prefix
	api := r.PathPrefix("/api/v1").Subrouter()

	// Workers endpoints
	api.HandleFunc("/workers", h.CreateWorker).Methods("POST")
	api.HandleFunc("/workers", h.GetWorkers).Methods("GET")
	api.HandleFunc("/workers/{id}", h.GetWorker).Methods("GET")
	api.HandleFunc("/workers/{id}", h.UpdateWorker).Methods("PUT")
	api.HandleFunc("/workers/{id}", h.DeleteWorker).Methods("DELETE")

	// Fields endpoints
	api.HandleFunc("/fields", h.CreateField).Methods("POST")
	api.HandleFunc("/fields", h.GetFields).Methods("GET")
	api.HandleFunc("/fields/{id}", h.GetField).Methods("GET")
	api.HandleFunc("/fields/{id}", h.UpdateField).Methods("PUT")
	api.HandleFunc("/fields/{id}", h.DeleteField).Methods("DELETE")

	// Schedules endpoints
	api.HandleFunc("/schedules", h.CreateSchedule).Methods("POST")
	api.HandleFunc("/schedules", h.GetSchedules).Methods("GET")
	api.HandleFunc("/schedules/{id}", h.GetSchedule).Methods("GET")
	api.HandleFunc("/schedules/{id}", h.UpdateSchedule).Methods("PUT")
	api.HandleFunc("/schedules/{id}", h.DeleteSchedule).Methods("DELETE")
	api.HandleFunc("/workers/{workerId}/schedules", h.GetWorkerSchedules).Methods("GET")

	// Operations endpoints
	api.HandleFunc("/operations", h.CreateOperation).Methods("POST")
	api.HandleFunc("/operations", h.GetOperations).Methods("GET")
	api.HandleFunc("/operations/{id}", h.GetOperation).Methods("GET")
	api.HandleFunc("/operations/{id}", h.UpdateOperation).Methods("PUT")
	api.HandleFunc("/operations/{id}", h.DeleteOperation).Methods("DELETE")
	api.HandleFunc("/operations/{id}/complete", h.CompleteOperation).Methods("POST")
	api.HandleFunc("/operations/{id}/start", h.StartOperation).Methods("POST")

	// Reports endpoints
	api.HandleFunc("/reports/daily", h.GetDailyReport).Methods("GET")
	api.HandleFunc("/reports/monthly", h.GetMonthlyReport).Methods("GET")
	api.HandleFunc("/reports/yearly", h.GetYearlyReport).Methods("GET")
	api.HandleFunc("/reports/daily/export", h.ExportDailyReport).Methods("GET")
	api.HandleFunc("/reports/monthly/export", h.ExportMonthlyReport).Methods("GET")
	api.HandleFunc("/reports/yearly/export", h.ExportYearlyReport).Methods("GET")

	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")
}
