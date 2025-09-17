package flows

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
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
			ai.WithSystem("You are a helpful cooking assistant. Respond only with JSON matching the given schema. Do not include any extra text."),
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

Respond with concise instructions and note any cooking tips when helpful.

Respond ONLY with valid JSON that matches this structure:
{
  "title": string,
  "description": string,
  "prepTime": string,
  "cookTime": string,
  "servings": number,
  "ingredients": [string],
  "instructions": [string],
  "tips": [string] (optional)
}
Do not include any extra text outside the JSON.`, ingredient, restriction)
	return prompt, restriction
}

// Make a defensive copy so GenerateData-internal option mutations don't leak into fallback.
var generateRecipe = func(ctx context.Context, app *genkit.Genkit, options ...ai.GenerateOption) (*models.Recipe, error) {
	// Text-first mode (for models where native structured output may be flaky)
	if os.Getenv("RECIPE_TEXT_ONLY") == "1" {
		text, err := genkit.GenerateText(ctx, app, options...)
		if err != nil {
			return nil, err
		}
		text = strings.TrimSpace(text)
		var r models.Recipe
		if uErr := json.Unmarshal([]byte(text), &r); uErr == nil {
			return &r, nil
		}
		if candidate, ok := extractFirstJSONObject(text); ok {
			if uErr2 := json.Unmarshal([]byte(candidate), &r); uErr2 == nil {
				return &r, nil
			}
		}
		return nil, fmt.Errorf("model returned non-JSON or malformed JSON")
	}
	optCopy := append([]ai.GenerateOption(nil), options...)
	if recipe, _, err := genkit.GenerateData[models.Recipe](ctx, app, optCopy...); err == nil {
		return recipe, nil
	} else {
		// Fallback: try plain text then unmarshal to the expected struct (with one retry)
		attempt := 0
		for attempt < 2 {
			text, err2 := genkit.GenerateText(ctx, app, options...)
			if err2 != nil {
				return nil, err
			}
			text = strings.TrimSpace(text)
			var r models.Recipe
			if uErr := json.Unmarshal([]byte(text), &r); uErr == nil {
				return &r, nil
			}
			if candidate, ok := extractFirstJSONObject(text); ok {
				if uErr2 := json.Unmarshal([]byte(candidate), &r); uErr2 == nil {
					return &r, nil
				}
			}
			attempt++
		}
		return nil, err
	}
}

// extractFirstJSONObject scans text and returns the first top-level JSON object by balancing braces.
func extractFirstJSONObject(text string) (string, bool) {
	start := -1
	depth := 0
	for i, ch := range text {
		if ch == '{' {
			if depth == 0 {
				start = i
			}
			depth++
		} else if ch == '}' {
			if depth > 0 {
				depth--
				if depth == 0 && start >= 0 {
					return text[start : i+1], true
				}
			}
		}
	}
	return "", false
}
