package main

import (
	"encoding/json"
	"net/http"
	"regexp"
	"rest/internal/services/recipes"

	"github.com/gosimple/slug"
)

func main() {
	startStore()
}

func startStore() {
	store := recipes.NewMemStore()
	recipesHandler := NewRecipesHandler(store)

	mux := http.NewServeMux()

	mux.Handle("/", &HomeHandler{})
	mux.Handle("/recipes", recipesHandler)
	mux.Handle("/recipes/", recipesHandler)

	http.ListenAndServe("localhost:8080", mux)

}

type HomeHandler struct{}

func (hh *HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is my home page"))
}

type RecipesHandler struct {
	store recipeStore
}

func NewRecipesHandler(s recipeStore) *RecipesHandler {
	return &RecipesHandler{store: s}
}

var (
	RecipeRe     = regexp.MustCompile(`^/recipes/*$`)
	RecipeWithID = regexp.MustCompile(`^/recipes/(a-z0-9+(?:-[a-z0-9]+)+)$`)
)

func (rh *RecipesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost && RecipeRe.MatchString(r.URL.Path):
		rh.CreateRecipe(w, r)
	case r.Method == http.MethodGet && RecipeRe.MatchString(r.URL.Path):
		rh.ListRecipes(w, r)
	case r.Method == http.MethodGet && RecipeWithID.MatchString(r.URL.Path):
		rh.GetRecipe(w, r)
	case r.Method == http.MethodPut && RecipeWithID.MatchString(r.URL.Path):
		rh.UpdateRecipe(w, r)
	case r.Method == http.MethodDelete && RecipeWithID.MatchString(r.URL.Path):
		rh.DeleteRecipe(w, r)
		return
	default:
		return
	}
}

func (rh *RecipesHandler) CreateRecipe(w http.ResponseWriter, r *http.Request) {
	var recipe recipes.Recipe
	if err := json.NewDecoder(r.Body).Decode(&recipe); err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	resourceID := slug.Make(recipe.Name)
	if err := rh.store.Add(resourceID, recipe); err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)

}
func (rh *RecipesHandler) ListRecipes(w http.ResponseWriter, r *http.Request) {
	recipes, err := rh.store.List()
	if err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	jsonBytes, err := json.Marshal(recipes)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
func (rh *RecipesHandler) GetRecipe(w http.ResponseWriter, r *http.Request) {
	matches := RecipeWithID.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		InternalServerErrorHandler(w, r)
		return
	}

	recipe, err := rh.store.Get(matches[1])
	if err != nil {
		if err == recipes.NotFoundErr {
			NotFoundHandler(w, r)
			return
		}

		InternalServerErrorHandler(w, r)
		return
	}

	jsonBytes, err := json.Marshal(recipe)
	if err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)

}
func (rh *RecipesHandler) UpdateRecipe(w http.ResponseWriter, r *http.Request) {
	matches := RecipeWithID.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		InternalServerErrorHandler(w, r)
		return
	}

	var recipe recipes.Recipe
	if err := json.NewDecoder(r.Body).Decode(&recipe); err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	if err := rh.store.Update(matches[1], recipe); err != nil {
		if err == recipes.NotFoundErr {
			NotFoundHandler(w, r)
			return
		}

		InternalServerErrorHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
}
func (rh *RecipesHandler) DeleteRecipe(w http.ResponseWriter, r *http.Request) {
	matches := RecipeWithID.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		InternalServerErrorHandler(w, r)
		return
	}

	if err := rh.store.Remove(matches[1]); err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
}

type recipeStore interface {
	Add(name string, recipe recipes.Recipe) error
	Get(name string) (recipes.Recipe, error)
	Update(name string, recipe recipes.Recipe) error
	List() (map[string]recipes.Recipe, error)
	Remove(name string) error
}

func InternalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("500 Internal Server Error"))
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 Not Found"))
}
