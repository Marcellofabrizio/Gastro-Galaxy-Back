package server

import (
	"encoding/json"
	"fmt"
	"gastro-galaxy-back/internal/models"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", s.HelloWorldHandler)

	r.Get("/health", s.healthHandler)

	r.Post("/recipe", s.insertRecipeHandler)

	return r
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, _ := json.Marshal(s.db.Health())
	_, _ = w.Write(jsonResp)
}

func (s *Server) insertRecipeHandler(w http.ResponseWriter, r *http.Request) {

	var recipe models.Recipe

	err := json.NewDecoder(r.Body).Decode(&recipe)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	id, err := s.db.Insert(recipe.Name, recipe.Description, recipe.Url, recipe.CategoryId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Recipe id: %d", id)
}
