package server

import (
	"encoding/json"
	"fmt"
	"gastro-galaxy-back/internal/models"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", s.HelloWorldHandler)

	r.Get("/health", s.HealthHandler)

	r.Get("/recipes", s.GetRecipesHandler)

	r.Get("/recipe/{recipeId}", s.GetRecipeWithIngredientsHandler)

	r.Put("/recipe/{recipeId}", s.PutRecipeHandler)

	r.Post("/recipe", s.InsertRecipeHandler)

	r.Post("/ingredient", s.InsertIngredientHandler)

	r.Get("/ingredients", s.GetIngredientsHandler)

	return r
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Gastro Galaxy Back-End"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) HealthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, _ := json.Marshal(s.db.Health())
	_, _ = w.Write(jsonResp)
}

func (s *Server) InsertRecipeHandler(w http.ResponseWriter, r *http.Request) {

	var recipeDto models.RecipeInputDto

	if err := json.NewDecoder(r.Body).Decode(&recipeDto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	recipe := models.Recipe{
		Name:            recipeDto.Name,
		Url:             recipeDto.Url,
		CategoryId:      recipeDto.CategoryId,
		Description:     recipeDto.Description,
		LongDescription: recipeDto.LongDescription,
	}

	id, err := s.db.InsertRecipe(recipe.Name, recipe.Description, recipe.LongDescription, recipe.Url, recipe.CategoryId, recipeDto.IngredientIds)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Recipe id: %d", id)
}

func (s *Server) GetRecipesHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var data map[string]interface{}

	var category string

	if len(body) > 0 {
		if err := json.Unmarshal(body, &data); err != nil {
			fmt.Println(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		category = data["category"].(string)
	} else {
		category = ""
	}

	recipes, err := s.db.GetRecipes(category)

	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(recipes)
}

func (s *Server) GetRecipeWithIngredientsHandler(w http.ResponseWriter, r *http.Request) {

	recipeId, err := strconv.Atoi(r.PathValue("recipeId"))

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	recipe, err := s.db.GetRecipeWithIngredients(recipeId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if recipe == nil {
		http.Error(w, "Recipe not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(recipe)

}

func (s *Server) PutRecipeHandler(w http.ResponseWriter, r *http.Request) {

	recipeId, err := strconv.Atoi(r.PathValue("recipeId"))

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	recipe, err := s.db.GetRecipeWithIngredients(recipeId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if recipe == nil {
		http.Error(w, "Recipe not found", http.StatusNotFound)
		return
	}

	var recipeDto models.RecipeInputDto

	if err := json.NewDecoder(r.Body).Decode(&recipeDto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedRecipe := models.Recipe{
		Name:        recipeDto.Name,
		Url:         recipeDto.Url,
		Description: recipeDto.Description,
	}

	if err := s.db.UpdateRecipe(recipeId, updatedRecipe.Name, updatedRecipe.Description, updatedRecipe.Url); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := s.db.InsertRecipeIngredient(recipeId, recipeDto.IngredientIds); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return

	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Recipe UPDATED")

}

func (s *Server) InsertIngredientHandler(w http.ResponseWriter, r *http.Request) {

	var ingredient models.Ingedient

	err := json.NewDecoder(r.Body).Decode(&ingredient)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	id, err := s.db.InsertIngredient(ingredient.Name, ingredient.Amount, ingredient.Url, ingredient.IsAvailable)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Ingredient id: %d", id)
}

func (s *Server) GetIngredientsHandler(w http.ResponseWriter, r *http.Request) {

	ingredients, err := s.db.GetIngredients()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ingredients)
}
