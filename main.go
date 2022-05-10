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
	"context"
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/andreis3/recipes-api/handlers"
	"github.com/andreis3/recipes-api/models"
)

var recipesHandler *handlers.RecipesHandlers

const (
	MONGO_URI        = "mongodb://root:root@localhost:27017/test?authSource=admin"
	MONGO_DB         = "demo"
	MONGO_COLLECTION = "recipes"
)

func init() {
	recipes := make([]models.Recipe, 0)
	files, _ := ioutil.ReadFile("recipes.json")
	_ = json.Unmarshal(files, &recipes)

	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(MONGO_URI))
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")

	var listOfRecipes []any
	for _, recipe := range recipes {
		recipe.RecipeID = xid.New().String()
		listOfRecipes = append(listOfRecipes, recipe)
	}
	client.Database(MONGO_DB).Collection(MONGO_COLLECTION).Drop(ctx)
	collection := client.Database(MONGO_DB).Collection(MONGO_COLLECTION)
	recipesHandler = handlers.NewRecipesHandlers(ctx, collection)
	insertManyResult, err := collection.InsertMany(ctx, listOfRecipes)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Inserted recipes", len(insertManyResult.InsertedIDs))
}

func main() {
	router := gin.Default()

	router.POST("/recipes", recipesHandler.NewRecipeHandler)
	router.GET("/recipes/:id", recipesHandler.GetRecipeIDHandler)
	router.GET("/recipes", recipesHandler.ListRecipesHandler)
	router.GET("/recipes/search", recipesHandler.SearchRecipesHandler)
	router.PUT("/recipes/:id", recipesHandler.UpdateRecipeHandler)
	router.DELETE("/recipes/:id", recipesHandler.DeleteRecipeHandler)

	router.Run(":3000")
}
