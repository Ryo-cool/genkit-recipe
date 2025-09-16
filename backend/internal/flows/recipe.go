package flows

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/core"
	"github.com/firebase/genkit/go/genkit"

	"example.com/genkit-recipe/internal/models"
)

const RecipeFlowName = "recipeGeneratorFlow"

// FlowConfig controls optional generation parameters.
type FlowConfig struct {
	Temperature     float32
	MaxOutputTokens int
	Logger          *slog.Logger
}

// RegisterRecipeFlow defines the recipe generator flow on the provided Genkit app instance.
func RegisterRecipeFlow(app *genkit.Genkit, cfg FlowConfig) *core.Flow[*models.RecipeInput, *models.Recipe, struct{}] {
	return genkit.DefineFlow(app, RecipeFlowName, func(ctx context.Context, input *models.RecipeInput) (*models.Recipe, error) {
		logger := cfg.Logger
		if logger == nil {
			logger = slog.Default()
		}

		if input == nil {
			return nil, fmt.Errorf("input payload is required")
		}

		ingredient := strings.TrimSpace(input.Ingredient)
		if ingredient == "" {
			return nil, fmt.Errorf("ingredient must be provided")
		}

		prompt, restriction := BuildRecipePrompt(ingredient, input.DietaryRestrictions)

		logger.Info("recipe flow invoked", "ingredient", ingredient, "dietaryRestrictions", restriction)

		options := []ai.GenerateOption{
			ai.WithPrompt(prompt),
		}
		if cfg.Temperature > 0 || cfg.MaxOutputTokens > 0 {
			config := map[string]any{}
			if cfg.Temperature > 0 {
				config["temperature"] = cfg.Temperature
			}
			if cfg.MaxOutputTokens > 0 {
				config["maxOutputTokens"] = cfg.MaxOutputTokens
			}
			if len(config) > 0 {
				options = append(options, ai.WithConfig(config))
			}
		}

		recipe, err := generateRecipe(ctx, app, options...)
		if err != nil {
			logger.Error("recipe generation failed", "error", err)
			return nil, fmt.Errorf("failed to generate recipe: %w", err)
		}

		logger.Info("recipe generated", "title", recipe.Title, "servings", recipe.Servings)
		return recipe, nil
	})
}

// BuildRecipePrompt assembles the prompt and returns the sanitized dietary restriction for logging/tests.
func BuildRecipePrompt(ingredient, restriction string) (string, string) {
	restriction = strings.TrimSpace(restriction)
	if restriction == "" {
		restriction = "none"
	}

	prompt := fmt.Sprintf(`Create a complete cooking recipe using the following requirements.

Main ingredient or cuisine focus: %s
Dietary requirements: %s

Respond with concise instructions and note any cooking tips when helpful.`, ingredient, restriction)
	return prompt, restriction
}

var generateRecipe = func(ctx context.Context, app *genkit.Genkit, options ...ai.GenerateOption) (*models.Recipe, error) {
	recipe, _, err := genkit.GenerateData[models.Recipe](ctx, app, options...)
	return recipe, err
}
