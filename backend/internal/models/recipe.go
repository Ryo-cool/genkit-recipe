package models

// RecipeInput represents the prompt parameters provided by HTTP clients.
type RecipeInput struct {
	Ingredient          string `json:"ingredient" jsonschema:"description=Main ingredient or cuisine focus"`
	DietaryRestrictions string `json:"dietaryRestrictions,omitempty" jsonschema:"description=Dietary requirements, omit for none"`
}

// Recipe captures the structured response returned from the LLM.
type Recipe struct {
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	PrepTime     string   `json:"prepTime"`
	CookTime     string   `json:"cookTime"`
	Servings     int      `json:"servings"`
	Ingredients  []string `json:"ingredients"`
	Instructions []string `json:"instructions"`
	Tips         []string `json:"tips,omitempty"`
}
