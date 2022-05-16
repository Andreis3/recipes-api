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
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/gin-contrib/sessions"
	redisStore "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/xid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/andreis3/recipes-api/handlers"
	"github.com/andreis3/recipes-api/metrics"
	"github.com/andreis3/recipes-api/models"
	"github.com/andreis3/recipes-api/utils"
)

var authHandler *handlers.AuthHandler
var recipesHandler *handlers.RecipesHandlers
var err error
var config utils.Config

func init() {
	config, err = utils.LoadConfig()
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	ctx := context.Background()
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err = mongoClient.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")
	redisClient := redisClient(ctx)

	collectionRecipes := mongoRecipesCollection(config, ctx, mongoClient)
	collectionUsers := mongoUsersCollection(mongoClient, config, ctx)

	recipesHandler = handlers.NewRecipesHandlers(ctx, collectionRecipes, redisClient, config)
	authHandler = handlers.NewAuthHandler(ctx, collectionUsers)

	metrics.RegisterPrometheus()
}

func main() {
	router := gin.Default()

	store, _ := redisStore.NewStore(10, "tcp", config.RedisURI, config.RedisPassword, []byte("secret"))
	router.Use(sessions.Sessions("recipes_api", store))
	router.Use(metrics.PrometheusMiddleware([]string{"/metrics"}))

	authorized := router.Group("/")
	authorized.Use(authHandler.AuthMiddleware())
	{
		authorized.POST("/recipes", recipesHandler.NewRecipeHandler)
		authorized.PUT("/recipes/:id", recipesHandler.UpdateRecipeHandler)
		authorized.DELETE("/recipes/:id", recipesHandler.DeleteRecipeHandler)
	}

	router.GET("/recipes/:id", recipesHandler.GetRecipeIDHandler)
	router.GET("/recipes", recipesHandler.ListRecipesHandler)
	router.GET("/recipes/search", recipesHandler.SearchRecipesHandler)

	router.POST("/signin", authHandler.SignInHandler)
	router.POST("/signout", authHandler.SignOutHandler)
	router.POST("/refresh", authHandler.RefreshHandler)

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	port := fmt.Sprintf(":%s", config.Port)
	router.Run(port)
}

func mongoRecipesCollection(config utils.Config, ctx context.Context, client *mongo.Client) *mongo.Collection {
	recipes := make([]models.Recipe, 0)
	files, _ := ioutil.ReadFile("recipes.json")
	json.Unmarshal(files, &recipes)

	var listOfRecipes []any
	for _, recipe := range recipes {
		recipe.RecipeID = xid.New().String()
		listOfRecipes = append(listOfRecipes, recipe)
	}

	client.Database(config.MongoDB).Collection(config.CollectionRecipes).Drop(ctx)
	collectionRecipes := client.Database(config.MongoDB).Collection(config.CollectionRecipes)

	insertManyRecipesResult, err := collectionRecipes.InsertMany(ctx, listOfRecipes)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Inserted recipes", len(insertManyRecipesResult.InsertedIDs))

	return collectionRecipes
}

func redisClient(ctx context.Context) *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.RedisURI,
		Password: config.RedisPassword,
		DB:       config.RedisDB,
	})

	status := redisClient.Ping(ctx)
	log.Println(status)

	return redisClient
}

func mongoUsersCollection(client *mongo.Client, config utils.Config, ctx context.Context) *mongo.Collection {
	users := map[string]string{
		"andrei": "123456",
		"admin":  "123456",
		"santos": "123456",
	}
	h := sha256.New()

	collectionUsers := client.Database(config.MongoDB).Collection(config.CollectionUsers)

	client.Database(config.MongoDB).Collection(config.CollectionUsers).Drop(ctx)
	for username, password := range users {
		collectionUsers.InsertOne(ctx, bson.M{
			"username": username,
			"password": string(h.Sum([]byte(password))),
		})
	}
	return collectionUsers
}
