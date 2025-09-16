# Genkit Recipe App

Gemini 2.5 Flash を使って構造化されたレシピを生成する Go + Genkit の最小アプリケーションです。フローは HTTP と Genkit Dev UI の両方から呼び出せます。

## セットアップ

1. 前提: Go 1.24 以上、Node.js 20 以上、Genkit CLI、Gemini API キー（`GEMINI_API_KEY`）。
2. 依存関係の取得:
   ```bash
   cd backend && go mod tidy
   cd ../frontend && npm install
   ```
3. 環境変数を設定:
   ```bash
   export GEMINI_API_KEY="your-key"
   ```

## 実行

### バックエンド（Go）
```bash
cd backend
go run ./cmd/recipe
```
- 起動時にサンプルレシピを標準出力へ表示します。
- `POST http://localhost:3400/recipeGeneratorFlow` が利用可能になります。

### Dev UI（推奨）
```bash
cd backend
genkit start -- go run ./cmd/recipe
```
- Dev UI は `http://localhost:4000` で開けます。

### フロントエンド（Next.js）
```bash
cd frontend
npm run dev
```
- ブラウザで `http://localhost:3000` を開くとフォームからレシピ生成を試せます。
- API エンドポイントを変更する場合は `.env.local` に `NEXT_PUBLIC_API_BASE` を設定します。

## テスト（今後）
- Go: `cd backend && go test ./...`
- Web: `cd frontend && npm run lint`

詳しい仕様や開発計画は `docs/` を参照してください。
