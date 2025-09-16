package flows

import (
	"context"
	"testing"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"

	"example.com/genkit-recipe/internal/models"
)

func TestBuildRecipePrompt(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name         string
		ingredient   string
		restriction  string
		wantPrompt   string
		wantSanitize string
	}{
		{
			name:        "empty restriction defaults to none",
			ingredient:  "tomato",
			restriction: "   ",
			wantPrompt: "Create a complete cooking recipe using the following requirements.\n\n" +
				"Main ingredient or cuisine focus: tomato\n" +
				"Dietary requirements: none\n\n" +
				"Respond with concise instructions and note any cooking tips when helpful.",
			wantSanitize: "none",
		},
		{
			name:        "passes through custom restriction",
			ingredient:  "miso",
			restriction: "gluten-free",
			wantPrompt: "Create a complete cooking recipe using the following requirements.\n\n" +
				"Main ingredient or cuisine focus: miso\n" +
				"Dietary requirements: gluten-free\n\n" +
				"Respond with concise instructions and note any cooking tips when helpful.",
			wantSanitize: "gluten-free",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			gotPrompt, gotRestriction := BuildRecipePrompt(tc.ingredient, tc.restriction)
			if gotPrompt != tc.wantPrompt {
				t.Errorf("prompt mismatch\nwant: %q\n got: %q", tc.wantPrompt, gotPrompt)
			}
			if gotRestriction != tc.wantSanitize {
				t.Errorf("restriction mismatch\nwant: %q\n got: %q", tc.wantSanitize, gotRestriction)
			}
		})
	}
}

func TestRegisterRecipeFlowHappyPath(t *testing.T) {
	t.Parallel()

	originalGenerator := generateRecipe
	defer func() { generateRecipe = originalGenerator }()

	want := &models.Recipe{
		Title:        "Tomato Delight",
		Description:  "A quick tomato-focused dish",
		PrepTime:     "10 minutes",
		CookTime:     "20 minutes",
		Servings:     2,
		Ingredients:  []string{"2 tomatoes", "1 tbsp olive oil"},
		Instructions: []string{"Slice tomatoes", "Saute with olive oil"},
	}
	called := false

	generateRecipe = func(ctx context.Context, app *genkit.Genkit, options ...ai.GenerateOption) (*models.Recipe, error) {
		called = true
		if len(options) == 0 {
			t.Fatalf("expected options to include prompt")
		}
		return want, nil
	}

	app := genkit.Init(context.Background())

	flow := RegisterRecipeFlow(app, FlowConfig{})
	got, err := flow.Run(context.Background(), &models.RecipeInput{Ingredient: "tomato"})
	if err != nil {
		t.Fatalf("flow run returned error: %v", err)
	}

	if !called {
		t.Fatalf("expected generator to be called")
	}

	if got != want {
		t.Fatalf("expected recipe pointer to match stub, got %#v", got)
	}
}
