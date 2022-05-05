package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

var recipes []Recipe

func init() {
	recipes = make([]Recipe, 0)
	files, _ := ioutil.ReadFile("recipes.json")
	_ = json.Unmarshal([]byte(files), &recipes)
}

type Recipe struct {
	RecipeID     string    `json:"recipe_id"`
	Name         string    `json:"name"`
	Tags         []string  `json:"tags"`
	Ingredients  []string  `json:"ingredients"`
	Instructions []string  `json:"instructions"`
	PublishedAt  time.Time `json:"published_at"`
}

func NewRecipeHandler(c *gin.Context) {
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	recipe.RecipeID = xid.New().String()
	recipe.PublishedAt = time.Now()
	recipes = append(recipes, recipe)
	c.JSON(http.StatusCreated, recipe)
}

func ListRecipesHandler(c *gin.Context) {
	c.ShouldBindJSON(&recipes)
	c.JSON(http.StatusOK, recipes)
}

func UpdateRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	index := -1
	for i := 0; i < len(recipes); i++ {
		if recipes[i].RecipeID == id {
			index = i
		}
	}

	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{"error": "recipe not found"})
	}
	recipe.RecipeID = recipes[index].RecipeID
	recipes[index] = recipe
	c.JSON(http.StatusOK, recipe)
}

func main() {
	router := gin.Default()
	router.POST("/recipes", NewRecipeHandler)
	router.GET("/recipes", ListRecipesHandler)
	router.PUT("/recipes/:id", UpdateRecipeHandler)
	router.Run(":3000")
}
