package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/firebase/genkit/go/core"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/googlegenai"
	"github.com/firebase/genkit/go/plugins/server"

	"example.com/genkit-recipe/internal/flows"
	"example.com/genkit-recipe/internal/models"
)

const (
	defaultBindAddress      = "127.0.0.1:3400"
	defaultModelEnvironment = "GENKIT_DEFAULT_MODEL"
	defaultModel            = "googleai/gemini-2.5-flash"
	temperatureEnvironment  = "RECIPE_TEMPERATURE"
	outputTokensEnvironment = "RECIPE_MAX_OUTPUT_TOKENS"
)

func main() {
	ctx := context.Background()
	slog.SetDefault(newLogger())
	logger := slog.Default()

	app := genkit.Init(ctx,
		genkit.WithPlugins(&googlegenai.GoogleAI{}),
		genkit.WithDefaultModel(resolveModel()),
	)

	genCfg := resolveGenerationConfig(logger)
	recipeFlow := flows.RegisterRecipeFlow(app, flows.FlowConfig{
		Temperature:     genCfg.Temperature,
		MaxOutputTokens: genCfg.MaxOutputTokens,
		Logger:          logger,
	})

	if os.Getenv("DISABLE_STARTUP_SAMPLE") != "1" {
		if err := runStartupSample(ctx, recipeFlow, logger); err != nil {
			logger.Error("startup sample generation failed", "error", err)
			os.Exit(1)
		}
	} else {
		logger.Info("startup sample disabled via DISABLE_STARTUP_SAMPLE=1")
	}

	mux := http.NewServeMux()
	mux.HandleFunc(fmt.Sprintf("POST /%s", flows.RecipeFlowName), genkit.Handler(recipeFlow))

	addr := resolveBindAddress()
	logger.Info("server starting", "address", fmt.Sprintf("http://%s", addr))
	logger.Info("flow available", "method", "POST", "url", fmt.Sprintf("http://%s/%s", addr, flows.RecipeFlowName))

	if err := server.Start(ctx, addr, mux); err != nil {
		logger.Error("server exited", "error", err)
		os.Exit(1)
	}
}

func runStartupSample(ctx context.Context, flow *core.Flow[*models.RecipeInput, *models.Recipe, struct{}], logger *slog.Logger) error {
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
	logger.Info("startup recipe generated", "payload", string(payload))
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

func resolveModel() string {
	if value := strings.TrimSpace(os.Getenv(defaultModelEnvironment)); value != "" {
		return value
	}
	return defaultModel
}

type generationConfig struct {
	Temperature     float32
	MaxOutputTokens int
}

func resolveGenerationConfig(logger *slog.Logger) generationConfig {
	cfg := generationConfig{
		Temperature:     0.4,
		MaxOutputTokens: 600,
	}

	if raw := strings.TrimSpace(os.Getenv(temperatureEnvironment)); raw != "" {
		value, err := strconv.ParseFloat(raw, 32)
		if err != nil {
			logger.Warn("invalid temperature override, falling back to default", "value", raw, "error", err)
		} else if value < 0 || value > 2 {
			logger.Warn("temperature override out of range, expecting 0-2", "value", value)
		} else {
			cfg.Temperature = float32(value)
		}
	}

	if raw := strings.TrimSpace(os.Getenv(outputTokensEnvironment)); raw != "" {
		value, err := strconv.Atoi(raw)
		if err != nil {
			logger.Warn("invalid max output tokens override, falling back to default", "value", raw, "error", err)
		} else if value <= 0 {
			logger.Warn("max output tokens must be positive", "value", value)
		} else {
			cfg.MaxOutputTokens = value
		}
	}

	logger.Info("generation config resolved", "temperature", cfg.Temperature, "maxOutputTokens", cfg.MaxOutputTokens)
	return cfg
}

func newLogger() *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	return slog.New(handler)
}
