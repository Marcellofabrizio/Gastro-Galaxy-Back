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

	r.Post("/ingredient", s.insertIngredientHandler)

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

	var recipeDto models.RecipeDto

	if err := json.NewDecoder(r.Body).Decode(&recipeDto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	recipe := models.Recipe{
		Name:        recipeDto.Name,
		Url:         recipeDto.Url,
		CategoryId:  recipeDto.CategoryId,
		Description: recipeDto.Description,
	}

	id, err := s.db.InsertRecipe(recipe.Name, recipe.Description, recipe.Url, recipe.CategoryId, recipeDto.IngedientIds)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Recipe id: %d", id)
}

func (s *Server) insertIngredientHandler(w http.ResponseWriter, r *http.Request) {

	var ingredient models.Ingedient

	err := json.NewDecoder(r.Body).Decode(&ingredient)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	id, err := s.db.InsertIngredient(ingredient.Name, ingredient.Amount, ingredient.Url, ingredient.IsAvailable)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Ingredient id: %d", id)
}
