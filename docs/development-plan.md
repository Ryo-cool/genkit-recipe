# Development Plan

_Last updated: 2025-09-16_

## Phase 0 – Baseline Validation
- [ ] Confirm Go server runs (`cd backend && go run ./cmd/recipe`) and returns structured JSON for sample requests.
  - Blocked 2025-09-16: `go run` 呼び出し時に DNS 解決がサンドボックスで拒否され、Gemini API へ到達できず起動失敗。`GEMINI_API_KEY` は設定済みだがネットワーク制限が原因。
- [ ] Verify Genkit Dev UI traces appear via `cd backend && genkit start -- go run ./cmd/recipe`.
  - Blocked 2025-09-16: 上記と同様にバックエンド起動が失敗するため未検証。
- [ ] Smoke test Next.js client (`cd frontend && npm run dev`) hitting local flow endpoint.
  - Pending 2025-09-16: バックエンドが未起動のため実施できず。ネットワーク制限解除後に再試行。
- [x] Document any setup discrepancies in `AGENTS.md`.
  - 実際のネットワーク制限事項を本ドキュメントに記載。

## Phase 1 – Backend Hardening
- [x] Refactor into `backend/cmd/recipe`, `backend/internal/flows`, `backend/internal/models` per spec.
- [x] Add input validation (missing ingredient, max string length).
- [x] Introduce configurable model parameters (temperature, max tokens).
- [x] Implement structured logging for flow invocations.
- [x] Create unit tests for prompt builder and flow happy-path.

## Phase 2 – Frontend Enhancements
- [x] Add client-side validation (empty ingredient, loading states with skeletons).
  - 完了 2025-09-16: `frontend/app/page.tsx` で必須入力チェック・最大文字数チェック・ローディングスケルトンを実装。
- [ ] Persist recent recipes locally (session storage) for quick comparisons.
- [ ] Create responsive layout with reusable UI primitives (CSS modules or design system).
- [ ] Add status indicator for backend connectivity.

## Phase 3 – Observability & QA
- [ ] Integrate `genkit evaluate` scenarios for regression testing prompts.
- [ ] Establish Playwright or Cypress smoke tests targeting the Next.js UI.
- [ ] Configure GitHub Actions (or preferred CI) for `go test` + `npm test`/lint.
- [ ] Capture baseline metrics (latency, token usage) for recipe flow.

## Phase 4 – Production Readiness
- [ ] Secure endpoint (API key header or OAuth) and document usage.
- [ ] Containerize Go service and Next.js app (Dockerfiles + compose).
- [ ] Externalize configuration (12-factor `.env`, secrets manager integration).
- [ ] Prepare deployment guide (e.g., Cloud Run + Vercel) with environment matrix.

## Phase 5 – Extensions (Backlog)
- [ ] Add nutrition breakdown via tool-calling to external APIs.
- [ ] Implement RAG module with Firebase vector store for curated recipes.
- [ ] Multi-provider toggle (OpenAI-compatible plugin) with feature flag.
- [ ] Internationalization (UI + prompt localization).

## Milestones & Checkpoints
- **M1**: Baseline validation complete, plan reviewed (ETA Week 1).
- **M2**: Backend & frontend hardened, CI running (ETA Week 3).
- **M3**: Production readiness tasks complete, ready for pilot users (ETA Week 5).
- **M4**: Select backlog features prioritized post-pilot (ETA Week 7).

---
Use this checklist to track progress in stand-ups; update statuses weekly.
