package models

import "time"

// swagger:parameters recipes newRecipe
type Recipe struct {
	//swagger:ignore
	ID           string    `json:"id" bson:"_id"`
	RecipeID     string    `json:"recipe_id" bson:"recipe_id"`
	Name         string    `json:"name" bson:"name"`
	Tags         []string  `json:"tags" bson:"tags"`
	Ingredients  []string  `json:"ingredients" bson:"ingredients"`
	Instructions []string  `json:"instructions" bson:"instructions"`
	PublishedAt  time.Time `json:"published_at" bson:"published_at"`
}
