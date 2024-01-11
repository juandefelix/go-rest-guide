package main

import (
	"example/jetbraind-rest-tutorial/pkg/recipes"

	"testing"
	"os"
	"bytes"
	"net/http"
	"net/http/httptest"
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

	// hamAndCheeseWithButter := readTestData("han_and_cheese_with_butter_recipe.json")
	// hamAndCheeseWithButterReader := bytes.NewReader(hamAndCheeseWithButter)

	// Create the request
	req := httptest.NewRequest(http.MethodPost, "/recipes", hamAndCheeseReader)
	w := httptest.NewRecorder()
	recipesHandler.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()
	assert.Equal(t, 200, res.StatusCode)

	saved, _ := store.List()
	assert.Len(t, saved, 1)

}
