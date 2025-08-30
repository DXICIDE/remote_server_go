package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/DXICIDE/remote_server_go/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	platform := os.Getenv("PLATFORM")

	if err != nil {
		log.Fatal("couldnt connect to the database")
	}

	dbQueries := database.New(db)

	filepathRoot := "app"
	mux := http.NewServeMux()
	apiCfg := &apiConfig{}
	apiCfg.db = dbQueries
	apiCfg.platform = platform
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(filepathRoot)))))

	mux.HandleFunc("GET /api/healthz", handlerHealthz)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerChirps)
	mux.HandleFunc("POST /api/users", apiCfg.handlerUser)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetChirps)

	s := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	s.ListenAndServe()
}
