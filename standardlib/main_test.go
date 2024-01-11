package main

import (
	"example/jetbraind-rest-tutorial/pkg/recipes"

	"testing"
	"os"
	"bytes"
	"net/http"
	"net/http/httptest"
	"io"
	"github.com/stretchr/testify/assert"
)

func readTestData(t *testing.T, name string) []byte {
	t.Helper()
	content, err := os.ReadFile("../testdata/" + name)
	if err != nil {
		t.Errorf("Could not read %v", name)
	}

	return content

}

func TestRecipesHandlerCRUD_Integration(t *testing.T) {
	store := recipes.NewMemStore()
	recipesHandler := NewRecipesHandler(store)

	hamAndCheese := readTestData(t, "ham_and_cheese_recipe.json")
	hamAndCheeseReader := bytes.NewReader(hamAndCheese)

	hamAndCheeseWithButter := readTestData(t, "ham_and_cheese_with_butter_recipe.json")
	hamAndCheeseWithButterReader := bytes.NewReader(hamAndCheeseWithButter)

	// CREATE
	req := httptest.NewRequest(http.MethodPost, "/recipes", hamAndCheeseReader)
	w := httptest.NewRecorder()
	recipesHandler.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()
	assert.Equal(t, 200, res.StatusCode)

	saved, _ := store.List()
	assert.Len(t, saved, 1)


	//GET
	req = httptest.NewRequest(http.MethodGet, "/recipes/ham-and-cheese-toasties", nil)
	w = httptest.NewRecorder()
	recipesHandler.ServeHTTP(w, req)

	res = w.Result()
	defer res.Body.Close()
	assert.Equal(t, 200, res.StatusCode)

	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	assert.JSONEq(t, string(hamAndCheese), string(data))


	//UPDATE
	req = httptest.NewRequest(http.MethodPut, "/recipes/ham-and-cheese-toasties", hamAndCheeseWithButterReader)
	w = httptest.NewRecorder()
	recipesHandler.ServeHTTP(w, req)

	res = w.Result()
	defer res.Body.Close()
	assert.Equal(t, 200, res.StatusCode)


	updatedHamAndCheese, err := store.Get("ham-and-cheese-toasties")
	assert.NoError(t, err)

	assert.Contains(t, updatedHamAndCheese.Ingredients, recipes.Ingredient{Name: "butter"})

	//DELETE
	req = httptest.NewRequest(http.MethodDelete, "/recipes/ham-and-cheese-toasties", nil)
	w = httptest.NewRecorder()
	recipesHandler.ServeHTTP(w, req)

	res = w.Result()
	defer res.Body.Close()
	assert.Equal(t, 200, res.StatusCode)

	saved, _ = store.List()
	assert.Len(t, saved, 0)
}
