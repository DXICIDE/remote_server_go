package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/DXICIDE/remote_server_go/internal/database"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
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
	mux.HandleFunc("POST /api/validate_chirp", handlerReq)
	mux.HandleFunc("POST /api/users", apiCfg.handlerUser)

	s := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	s.ListenAndServe()
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func handlerHealthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) handlerUser(w http.ResponseWriter, r *http.Request) {
	type email struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := email{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), params.Email)
	if err != nil {
		log.Printf("Error creating user: %s", err)
		w.WriteHeader(500)
		return
	}
	userMap := User{}
	userMap.ID = user.ID
	userMap.Email = user.Email
	userMap.CreatedAt = user.CreatedAt
	userMap.UpdatedAt = user.UpdatedAt
	respondWithJSON(w, 201, userMap)
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	hits := cfg.fileserverHits.Load()
	w.Write([]byte(fmt.Sprintf(`<!DOCTYPE html>
<html>
  <body>
	<h1>Welcome, Chirpy Admin</h1>
	<p>Chirpy has been visited %d times!</p>
  </body>
</html>`, hits)))
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, 403, "Error, not dev environment")
		return
	}

	if err := cfg.db.Resetu(r.Context()); err != nil {
		respondWithError(w, 500, "failed to reset database")
		return
	}

	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	cfg.fileserverHits.Store(0)
}

func handlerReq(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		// an error will be thrown if the JSON is invalid or has the wrong types
		// any missing fields will simply have their values in the struct set to their zero value
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}
	type cleaned_body struct {
		Cleaned_body string `json:"cleaned_body"`
	}

	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	bodysplit := strings.Split(params.Body, " ")

	for i := 0; i < len(bodysplit); i++ {
		body := strings.ToLower(bodysplit[i])
		if body == "kerfuffle" || body == "sharbert" || body == "fornax" {
			bodysplit[i] = "****"
		}
	}
	body := strings.Join(bodysplit, " ")
	cleanBody := &cleaned_body{}
	cleanBody.Cleaned_body = body
	respondWithJSON(w, 200, cleanBody)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type resErr struct {
		Error string `json:"error"`
	}
	respBody := resErr{
		Error: msg,
	}

	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}
