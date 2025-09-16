package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/firebase/genkit/go/core"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/googlegenai"
	"github.com/firebase/genkit/go/plugins/server"

	"example.com/genkit-recipe/internal/flows"
	"example.com/genkit-recipe/internal/models"
)

const (
	defaultBindAddress = "127.0.0.1:3400"
	defaultModel       = "googleai/gemini-2.5-flash"
)

func main() {
	ctx := context.Background()

	app := genkit.Init(ctx,
		genkit.WithPlugins(&googlegenai.GoogleAI{}),
		genkit.WithDefaultModel(defaultModel),
	)

	recipeFlow := flows.RegisterRecipeFlow(app, flows.FlowConfig{
		Temperature:     0.4,
		MaxOutputTokens: 600,
	})

	if err := runStartupSample(ctx, recipeFlow); err != nil {
		log.Fatalf("startup sample generation failed: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc(fmt.Sprintf("POST /%s", flows.RecipeFlowName), genkit.Handler(recipeFlow))

	addr := resolveBindAddress()
	log.Printf("Starting server on http://%s", addr)
	log.Printf("Flow available at: POST http://%s/%s", addr, flows.RecipeFlowName)

	if err := server.Start(ctx, addr, mux); err != nil {
		log.Fatalf("server exited: %v", err)
	}
}

func runStartupSample(ctx context.Context, flow *core.Flow[*models.RecipeInput, *models.Recipe, struct{}]) error {
	sample, err := flow.Run(ctx, &models.RecipeInput{
		Ingredient:          "avocado",
		DietaryRestrictions: "vegetarian",
	})
	if err != nil {
		return err
	}

	payload, err := json.MarshalIndent(sample, "", "  ")
	if err != nil {
		return err
	}
	log.Printf("Sample recipe generated at startup:\n%s", string(payload))
	return nil
}

func resolveBindAddress() string {
	if addr := strings.TrimSpace(os.Getenv("BIND_ADDR")); addr != "" {
		return addr
	}
	if port := strings.TrimSpace(os.Getenv("PORT")); port != "" {
		if strings.Contains(port, ":") {
			return port
		}
		return fmt.Sprintf("0.0.0.0:%s", port)
	}
	return defaultBindAddress
}
