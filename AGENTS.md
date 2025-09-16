# Genkit Resource Guide (AGENTS)

必ず日本語で返答します

## Core Documentation

- [Get started with Genkit using Go](https://firebase.google.com/docs/genkit-go/get-started-go) — Go 1.24+ setup, Gemini API key configuration, and first flow example.
- [Genkit developer tools (Go)](https://firebase.google.com/docs/genkit-go/devtools) — CLI install (`npm install -g genkit-cli`), Node.js ≥20 requirements, and Developer UI workflow.
- [Generating content with AI models](https://firebase.google.com/docs/genkit/models) — Covers structured output, `GenerateData[T]`, temperature/max token tuning, and provider switching.

## Flow & Server Patterns

- [Creating and securing flows](https://firebase.google.com/docs/genkit-go/flows) — Flow lifecycle, auth patterns, and deployment guidance.
- [Tool calling guide](https://firebase.google.com/docs/genkit/agents/tool-calling) — `DefineTool` usage, capability negotiation, and multi-step orchestration tips.

## Plugins & Integrations

- [Google AI (Gemini) Go plugin reference](https://pkg.go.dev/github.com/firebase/genkit/go/plugins/googlegenai) — Configuration options such as API key injection and model aliases.
- [Firebase plugin (Go)](https://genkit.dev/go/docs/plugins/firebase/) — Firestore vector store setup, credential env vars, and retriever definitions.
- [Plugin authoring (Go)](https://firebase.google.com/docs/genkit-go/plugin-authoring) — Create custom model/providers via `genkit.Plugin` implementations.
- [Community OpenAI plugin](https://thefireco.github.io/genkit-plugins/docs/plugins/genkitx-openai) — Alternative provider wiring, supported models, and install command.

## Tooling & Samples

- [Genkit CLI reference](https://genkit.dev/docs/devtools) — Command catalog (`genkit start`, `flow:run`, evaluations) and telemetry endpoints.
- [Official Genkit GitHub repo](https://github.com/firebase/genkit) — Release notes, stability tiers (Go = Beta), and sample code directories.
- [Go sample project](https://github.com/yukinagae/genkit-golang-sample) — Real-world wiring with Dev UI, Makefile, and `.env` setup for comparison.

## Operational Notes

- [Genkit for Go launch blog (Jul 17, 2024)](https://developers.googleblog.com/en/introducing-genkit-for-go-build-scalable-ai-powered-apps-in-go/) — Ecosystem overview, supported plugins, and observability considerations.
- [Genkit Go 1.0 announcement (Aug 2025)](https://developers.googleblog.com/en/announcing-genkit-go-10-and-enhanced-ai-assisted-development/) — Highlights release cadence, CLI bootstrap script, and pointers to examples/community.
- Monitor [Firebase Genkit releases on GitHub](https://github.com/firebase/genkit/releases) for changelogs and version pinning guidance.
- Environment prerequisites for this project: Go 1.24+, Node.js ≥20, exported `GEMINI_API_KEY`, optional Firebase credentials when enabling related plugins.
- Support: Firebase support portal, GitHub issues on `firebase/genkit`, and Discord linked from docs/blog posts above.

## Quick Checklist for New Agents

1. Install Go 1.24+ and Node.js 20+.
2. Install Genkit CLI: `npm install -g genkit-cli`.
3. Export `GEMINI_API_KEY` (or provider-specific secret).
4. Run `cd backend && go mod tidy` (ensures dependencies match docs examples).
5. Start Dev UI + flow: `cd backend && genkit start -- go run ./cmd/recipe` (server on :3400, UI on :4000 per docs).
6. Hit the flow: `curl -X POST http://localhost:3400/recipeGeneratorFlow -H 'Content-Type: application/json' -d '{"data":{"ingredient":"tomato","dietaryRestrictions":"vegan"}}'`.
7. Inspect traces in the Dev UI and compare against official samples when extending functionality.

_Revisit these links whenever Genkit announces SDK or CLI updates to stay aligned with upstream changes._
