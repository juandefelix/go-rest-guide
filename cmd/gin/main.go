package main


import (
	"github.com/gin-gonic/gin"
	"net/http"
	"example/jetbrain-rest-tutorial/pkg/recipes"
	"github.com/gosimple/slug"

)

type recipeStore interface {
    Add(name string, recipe recipes.Recipe) error
    Get(name string) (recipes.Recipe, error)
    List() (map[string]recipes.Recipe, error)
    Update(name string, recipe recipes.Recipe) error
    Remove(name string) error
}

type RecipesHandler struct{
	store recipeStore
}

func NewRecipesHandler(store recipeStore) *RecipesHandler {
    return &RecipesHandler{
        store: store,
    }
}

func (h RecipesHandler) CreateRecipe(c *gin.Context) {
	var recipe recipes.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := slug.Make(recipe.Name)
	h.store.Add(id, recipe)
	c.JSON(http.StatusOK, gin.H{"status": "success"})

}
func (h RecipesHandler) ListRecipes(c *gin.Context)  {
	list, err := h.store.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, list)
}
func (h RecipesHandler) GetRecipe(c *gin.Context)    {
	id := c.Param("id")
	recipe, err := h.store.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})

	}

	c.JSON(http.StatusOK, recipe)
}
func (h RecipesHandler) UpdateRecipe(c *gin.Context) {
	var recipe recipes.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		if err == recipes.NotFoundErr {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id := c.Param("id")
	err := h.store.Update(id, recipe)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
func (h RecipesHandler) DeleteRecipe(c *gin.Context) {
	id := c.Param("id")

	err := h.store.Remove(id)

	if err != nil {
		if err == recipes.NotFoundErr {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func main() {
	store := recipes.NewMemStore()
	recipesHandler := NewRecipesHandler(store)

	router := gin.Default()

	router.GET("/", homePage)


	router.GET("/recipes", recipesHandler.ListRecipes)
	router.POST("/recipes", recipesHandler.CreateRecipe)
	router.GET("/recipes/:id", recipesHandler.GetRecipe)
	router.PUT("/recipes/:id", recipesHandler.UpdateRecipe)
	router.DELETE("/recipes/:id", recipesHandler.DeleteRecipe)

	router.Run()
}

func homePage(c *gin.Context) {
	c.String(http.StatusOK, "Hola Manola!")
}
