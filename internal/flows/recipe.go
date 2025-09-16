package flows

import (
	"context"
	"fmt"
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
}

// RegisterRecipeFlow defines the recipe generator flow on the provided Genkit app instance.
func RegisterRecipeFlow(app *genkit.Genkit, cfg FlowConfig) *core.Flow[*models.RecipeInput, *models.Recipe, struct{}] {
	return genkit.DefineFlow(app, RecipeFlowName, func(ctx context.Context, input *models.RecipeInput) (*models.Recipe, error) {
		if input == nil {
			return nil, fmt.Errorf("input payload is required")
		}

		ingredient := strings.TrimSpace(input.Ingredient)
		if ingredient == "" {
			return nil, fmt.Errorf("ingredient must be provided")
		}

		restriction := strings.TrimSpace(input.DietaryRestrictions)
		if restriction == "" {
			restriction = "none"
		}

		prompt := fmt.Sprintf(`Create a complete cooking recipe using the following requirements.

Main ingredient or cuisine focus: %s
Dietary requirements: %s

Respond with concise instructions and note any cooking tips when helpful.`, ingredient, restriction)

		options := []ai.GenerateOption{
			ai.WithPrompt(prompt),
			ai.WithOutputType(models.Recipe{}),
		}
		if cfg.Temperature > 0 || cfg.MaxOutputTokens > 0 {
			options = append(options, ai.WithConfig(&ai.GenerationCommonConfig{
				Temperature:     float64(cfg.Temperature),
				MaxOutputTokens: cfg.MaxOutputTokens,
			}))
		}

		recipe, _, err := genkit.GenerateData[models.Recipe](ctx, app, options...)
		if err != nil {
			return nil, fmt.Errorf("failed to generate recipe: %w", err)
		}

		return recipe, nil
	})
}
