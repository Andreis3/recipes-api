package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/rs/xid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/andreis3/recipes-api/models"
	"github.com/andreis3/recipes-api/utils"
)

type RecipesHandlers struct {
	collection  *mongo.Collection
	ctx         context.Context
	redisClient *redis.Client
	config      utils.Config
}

func NewRecipesHandlers(ctx context.Context, collection *mongo.Collection, redisClient *redis.Client, config utils.Config) *RecipesHandlers {
	return &RecipesHandlers{collection: collection, ctx: ctx, redisClient: redisClient, config: config}
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
func (handler *RecipesHandlers) NewRecipeHandler(c *gin.Context) {
	var recipe models.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	recipe.ID = primitive.NewObjectID().Hex()
	recipe.RecipeID = xid.New().String()
	recipe.PublishedAt = time.Now()

	_, err := handler.collection.InsertOne(handler.ctx, recipe)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while inserting a new recipe"})
	}

	log.Println("Remove recipes data from redis")
	handler.redisClient.Del(handler.ctx, "recipes")

	iter := handler.redisClient.Scan(handler.ctx, 0, "recipes:search:*", 0).Iterator()
	for iter.Next(handler.ctx) {
		handler.redisClient.Del(handler.ctx, iter.Val())
	}

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
func (handler *RecipesHandlers) ListRecipesHandler(c *gin.Context) {
	val, err := handler.redisClient.Get(handler.ctx, "recipes").Result()
	if err == redis.Nil {
		log.Printf("Request to Mongo")

		cur, err := handler.collection.Find(handler.ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer cur.Close(handler.ctx)
		recipes := make([]models.Recipe, 0)
		for cur.Next(handler.ctx) {
			var recipe models.Recipe
			cur.Decode(&recipe)
			recipes = append(recipes, recipe)
		}
		data, _ := json.Marshal(recipes)
		handler.redisClient.Set(handler.ctx, "recipes", string(data), 0)
		c.JSON(http.StatusOK, recipes)
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		log.Printf("Request to Redis")
		recipes := make([]models.Recipe, 0)
		json.Unmarshal([]byte(val), &recipes)
		c.JSON(http.StatusOK, recipes)
	}
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
func (handler *RecipesHandlers) GetRecipeIDHandler(c *gin.Context) {
	id := c.Param("id")

	var recipeResult models.Recipe

	err := handler.collection.FindOne(handler.ctx, bson.M{"_id": id}).Decode(&recipeResult)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
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
func (handler *RecipesHandlers) UpdateRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	var recipe models.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//objectId, _ := primitive.ObjectIDFromHex(id)
	_, err := handler.collection.UpdateOne(handler.ctx, bson.M{
		"_id": id,
	}, bson.D{{"$set", bson.D{
		{"name", recipe.Name},
		{"instructions", recipe.Instructions},
		{"ingredients", recipe.Ingredients},
		{"tags", recipe.Tags},
	}}})
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while updating a recipe"})
		return
	}

	cur := handler.collection.FindOne(handler.ctx, bson.M{"_id": id})
	err = cur.Decode(&recipe)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Println("Remove recipes data from redis")
	handler.redisClient.Del(handler.ctx, "recipes")
	iter := handler.redisClient.Scan(handler.ctx, 0, "recipes:search:*", 0).Iterator()
	for iter.Next(handler.ctx) {
		handler.redisClient.Del(handler.ctx, iter.Val())
	}

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
func (handler *RecipesHandlers) DeleteRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	result, err := handler.collection.DeleteOne(handler.ctx, bson.M{"_id": id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
		return
	}

	log.Println("Remove recipes data from redis")
	handler.redisClient.Del(handler.ctx, "recipes")
	iter := handler.redisClient.Scan(handler.ctx, 0, "recipes:search:*", 0).Iterator()
	for iter.Next(handler.ctx) {
		handler.redisClient.Del(handler.ctx, iter.Val())
	}

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
func (handler *RecipesHandlers) SearchRecipesHandler(c *gin.Context) {
	tagQuery := c.Query("tag")
	key := fmt.Sprintf("recipes:search:%s", tagQuery)
	val, err := handler.redisClient.Get(handler.ctx, key).Result()
	if err == redis.Nil {
		log.Println("Request Mongo")

		recipes := make([]models.Recipe, 0)
		cur, err := handler.collection.Find(handler.ctx, bson.M{"tags": tagQuery})
		defer cur.Close(handler.ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		for cur.Next(handler.ctx) {
			var recipe models.Recipe
			err := cur.Decode(&recipe)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			recipes = append(recipes, recipe)
		}

		data, _ := json.Marshal(recipes)
		handler.redisClient.Set(handler.ctx, key, string(data), 0)

		c.JSON(http.StatusOK, recipes)
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		log.Println("Request Redis")
		recipes := make([]models.Recipe, 0)
		json.Unmarshal([]byte(val), &recipes)
		c.JSON(http.StatusOK, recipes)

	}

}
