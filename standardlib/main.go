package main

import (
	"encoding/json"
	"net/http"
	"regexp"
	"example/jetbrains-tutorials/standardlib/recipes"

	"github.com/gosimple/slug"
)


var (
	RecipeRe = regexp.MustCompile(`^/recipes/$`)
	RecipeReWithId = regexp.MustCompile(`^/recipes/([a-z0-9]+(?:-[a-z0-9]+)+)$`)
)

func main() {
	store := recipes.NewMemStore()
	recipesHandler := NewRecipesHandler(store)

	// Create a new request multiplexer (I understand that this is a router)
	mux := http.NewServeMux()

	mux.Handle("/", &homeHandler{})
	mux.Handle("/recipes", recipesHandler)
	mux.Handle("/recipes/", recipesHandler)

	http.ListenAndServe(":8080", mux)
}


type homeHandler struct{}

func (h *homeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is my home page."))
}

func InternalErrorServerHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("500 Internal Server Error"))
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 Not Found"))
}

type RecipesHandler struct{
	store recipeStore
}

func NewRecipesHandler(s recipeStore) *RecipesHandler {
	return &RecipesHandler{
		store: s,
	}
}

func (h *RecipesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost && RecipeRe.MatchString(r.URL.Path):
		h.CreateRecipe(w, r)
		return
	case r.Method == http.MethodGet && RecipeRe.MatchString(r.URL.Path):
		h.ListRecipe(w, r)
		return
	case r.Method == http.MethodGet && RecipeReWithId.MatchString(r.URL.Path):
		h.ReadRecipe(w, r)
		return
	case r.Method == http.MethodPut && RecipeReWithId.MatchString(r.URL.Path):
		h.UpdateRecipe(w, r)
		return
	case r.Method == http.MethodDelete && RecipeReWithId.MatchString(r.URL.Path):
		h.DeleteRecipe(w, r)
		return
	}
}

func (h *RecipesHandler) CreateRecipe(w http.ResponseWriter, r *http.Request) {
	var recipe recipes.Recipe
	if err := json.NewDecoder(r.Body).Decode(&recipe); err != nil {
		InternalErrorServerHandler(w, r)
		return
	}

	resourceId := slug.Make(recipe.Name)

	if err := h.store.Add(resourceId, recipe); err != nil {
		InternalErrorServerHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *RecipesHandler) ListRecipe(w http.ResponseWriter, r *http.Request) {
	resources, err := h.store.List()

	jsonBytes, err := json.Marshal(resources)
	if err != nil {
		InternalErrorServerHandler(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *RecipesHandler) ReadRecipe(w http.ResponseWriter, r *http.Request) {
	matches := RecipeReWithId.FindStringSubmatch(r.URL.Path)

	if len(matches) < 2 {
		InternalErrorServerHandler(w, r)
		return
	}

	recipe, err := h.store.Get(matches[1])

	if err != nil {
		if err == recipes.NotFoundErr {
			NotFoundHandler(w, r)
			return
		}

		InternalErrorServerHandler(w, r)
		return
	}

	jsonByte, err := json.Marshal(recipe)
	if err != nil {
		InternalErrorServerHandler(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonByte)
}
func (h *RecipesHandler) UpdateRecipe(w http.ResponseWriter, r *http.Request) {}
func (h *RecipesHandler) DeleteRecipe(w http.ResponseWriter, r *http.Request) {}


type recipeStore interface {
	Add(name string, recipe recipes.Recipe) error
	Get(name string) (recipes.Recipe, error)
	Update(name string, recipe recipes.Recipe) error
	List() (map[string]recipes.Recipe, error)
	Remove(name string) error
}
