package main

import (
	"encoding/json"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/gosimple/slug"
	"example/jetbrain-rest-tutorial/pkg/recipes"
)

type homeHandler struct {}

func (h *homeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>This is my homepage</h1>"))
}

func InternalErrorServerHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("500 Internal Server Error"))
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 Not Found"))
}

func main() {
	store := recipes.NewMemStore()
	router := mux.NewRouter()

	recipesSubRouter := router.PathPrefix("/recipes").Subrouter()
	NewRecipesHandler(store, recipesSubRouter)

	http.ListenAndServe(":8010", router)
}

type recipeStore interface {
	Add(name string, recipe recipes.Recipe) error
	Get(name string) (recipes.Recipe, error)
	Update(name string, recipe recipes.Recipe) error
	List() (map[string]recipes.Recipe, error)
	Remove(name string) error
}

type RecipesHandler struct{
	store recipeStore
}

func NewRecipesHandler(store recipeStore, router *mux.Router) *RecipesHandler {
	handler :=  &RecipesHandler{
		store: store,
	}

	router.HandleFunc("/", handler.ListRecipes).Methods("GET")
	router.HandleFunc("/", handler.CreateRecipe).Methods("POST")
	router.HandleFunc("/{id}", handler.GetRecipe).Methods("GET")
	router.HandleFunc("/{id}", handler.UpdateRecipe).Methods("PUT")
	router.HandleFunc("/{id}", handler.DeleteRecipe).Methods("DELETE")

	return handler
}

func (h RecipesHandler) CreateRecipe(w http.ResponseWriter, r *http.Request) {
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

func (h RecipesHandler) ListRecipes(w http.ResponseWriter, r *http.Request) {
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

func (h RecipesHandler) GetRecipe(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	recipe, err := h.store.Get(id)

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

func (h RecipesHandler) UpdateRecipe(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

    var recipe recipes.Recipe
    if err := json.NewDecoder(r.Body).Decode(&recipe); err != nil {
        InternalErrorServerHandler(w, r)
        return
    }

    if err := h.store.Update(id, recipe); err != nil {
		if err == recipes.NotFoundErr {
			NotFoundHandler(w, r)
			return
		}

		InternalErrorServerHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
}
func (h RecipesHandler) DeleteRecipe(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

    if err :=h.store.Remove(id); err != nil {
    	InternalErrorServerHandler(w, r)
        return
    }

    w.WriteHeader(http.StatusOK)
}
