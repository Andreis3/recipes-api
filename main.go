// Recipes API
//
// This is a sample recipes API. You can find out more about the API at Building-Distributed-Applications-in-Gin.
//
//	Schemes: http
//  Host: localhost:3000
//	BasePath: /
//	Version: 1.0.0
//	Contact: Andrei Santos <andrei.as3@hotmail.com> https://github.com/Andreis3
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
// swagger:meta
package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
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

// swagger:parameters recipes newRecipe
type Recipe struct {
	//swagger:ignore
	RecipeID     string    `json:"recipe_id"`
	Name         string    `json:"name"`
	Tags         []string  `json:"tags"`
	Ingredients  []string  `json:"ingredients"`
	Instructions []string  `json:"instructions"`
	PublishedAt  time.Time `json:"published_at"`
}

// swagger:operation POST /recipes recipes newRecipe
// Create a new recipe
// ---
// produces:
// - application/json
// responses:
//     '200':
//         description: Successful operation
//     '400':
//         description: Invalid input
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

// swagger:operation GET /recipes recipes listRecipes
// Returns list of recipes
// ---
// produces:
// - application/json
// responses:
//     '200':
//         description: Successful operation
func ListRecipesHandler(c *gin.Context) {
	c.ShouldBindJSON(&recipes)
	c.JSON(http.StatusOK, recipes)
}

// swagger:operation GET /recipes/{id} recipes oneRecipe
// Get one recipe
// ---
// produces:
// - application/json
// parameters:
//   - name: id
//     in: path
//     description: ID of the recipe
//     required: true
//     type: string
// responses:
//     '200':
//         description: Successful operation
//     '404':
//         description: Invalid recipe ID
func GetRecipeIDHandler(c *gin.Context) {
	id := c.Param("id")

	var recipeResult Recipe

	for _, recipe := range recipes {
		if recipe.RecipeID == id {
			recipeResult = recipe
		}
	}

	if len(recipeResult.RecipeID) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "recipe not found"})
		return
	}
	c.JSON(http.StatusOK, recipeResult)
}

// swagger:operation PUT /recipes/{id} recipes updateRecipe
// Update an existing recipe
// ---
// parameters:
// - name: id
//   in: path
//   description: ID of the recipe
//   required: true
//   type: string
// produces:
// - application/json
// responses:
//     '200':
//         description: Successful operation
//     '400':
//         description: Invalid input
//     '404':
//         description: Invalid recipe ID
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
		return
	}
	recipe.RecipeID = recipes[index].RecipeID
	recipes[index] = recipe
	c.JSON(http.StatusOK, recipe)
}

// swagger:operation DELETE /recipes/{id} recipes deleteRecipe
// Delete an existing recipe
// ---
// produces:
// - application/json
// parameters:
//   - name: id
//     in: path
//     description: ID of the recipe
//     required: true
//     type: string
// responses:
//     '200':
//         description: Successful operation
//     '404':
//         description: Invalid recipe I
func DeleteRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	index := -1

	for key, recipe := range recipes {
		if recipe.RecipeID == id {
			index = key
		}
	}

	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{"error": "recipe not found"})
		return
	}
	recipes = append(recipes[:index], recipes[index+1:]...)
	c.JSON(http.StatusNoContent, nil)
}

// swagger:operation GET /recipes/search recipes findRecipe
// Search recipes based on tags
// ---
// produces:
// - application/json
// parameters:
//   - name: tag
//     in: query
//     description: recipe tag
//     required: true
//     type: string
// responses:
//     '200':
//         description: Successful operation
func SearchRecipesHandler(c *gin.Context) {
	tagQuery := c.Query("tag")
	listOfRecipes := make([]Recipe, 0)

	for _, recipe := range recipes {
		for _, tag := range recipe.Tags {
			if strings.EqualFold(tag, tagQuery) {
				listOfRecipes = append(listOfRecipes, recipe)
			}
		}
	}
	c.JSON(http.StatusOK, listOfRecipes)
}

func main() {
	router := gin.Default()
	router.POST("/recipes", NewRecipeHandler)
	router.GET("/recipes/:id", GetRecipeIDHandler)
	router.GET("/recipes", ListRecipesHandler)
	router.GET("/recipes/search", SearchRecipesHandler)
	router.PUT("/recipes/:id", UpdateRecipeHandler)
	router.DELETE("/recipes/:id", DeleteRecipeHandler)

	router.Run(":3000")
}
