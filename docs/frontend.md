# Frontend (Next.js) Overview

## Stack
- Next.js 14 App Router with TypeScript and React 18.
- Standalone output mode for deployability (`next.config.mjs`).
- Inline CSS for now; tailwind or design system can be layered on later.

## Directory Layout
```
frontend/
  app/
    layout.tsx      # Root layout + metadata
    page.tsx        # Main recipe form and renderer
    globals.css     # Base styles
  next.config.mjs
  package.json
  tsconfig.json
  next-env.d.ts
  .eslintrc.json
```

バックエンドとの分離を明確にするため、Next.js プロジェクトは `frontend/` 直下にまとまっています。Go 側は `backend/` に収容され、モジュール構成が分かりやすくなりました。

## Local Development
1. Install dependencies once: `cd frontend && npm install`.
2. Run the backend flow server (`cd backend && go run ./cmd/recipe` または `cd backend && genkit start -- go run ./cmd/recipe`).
3. Start the Next.js dev server: `npm run dev` (defaults to `http://localhost:3000`).
4. Update `.env.local` if you proxy through Next.js (optional; see below).

## Environment Variables
- `NEXT_PUBLIC_API_BASE` (optional): set to override the default flow endpoint (`http://localhost:3400`).
  - Example `.env.local`: `NEXT_PUBLIC_API_BASE="https://api.example.com"`

## Network Flow
1. Form submit posts to `${NEXT_PUBLIC_API_BASE ?? "http://localhost:3400"}/recipeGeneratorFlow`.
2. Response JSON is parsed and rendered into the UI sections (ingredients, instructions, tips).
3. Errors bubble up via inline alert state.

## Future Enhancements
- Move inline styles into CSS Modules or a design system.
- Add loading skeletons and persisted history per session.
- Use Next.js Route Handlers under `app/api` to proxy requests (and hide CORS differences in production).
- Integrate auth once backend flow is secured.

_Last updated: 2025-09-16_
