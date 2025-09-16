# Genkit レシピアプリケーション仕様書

## 1. 概要

- **目的**: Genkit Go 1.0 と Gemini 2.5 Flash を活用した型安全なレシピ生成 API を提供する。このサービスは HTTP フローを公開し、Genkit Dev UI を通じて観測可能性を提供する。
- **主要ユーザー**: AI 支援食品提案を探求する社内プロダクトチーム、デモ準備を行う開発者アドボケート、構造化出力を検証する QA エンジニア。
- **価値提案**: Go 構造体にマッピングされた構造化 JSON レスポンスを保証し、スキーマドリフトなしでダウンストリーム統合を可能にする。

## 2. 目標と非目標

### 目標

- `RecipeInput` → `Recipe` スキーマを通じて決定論的なリクエスト/レスポンス契約を提供する。
- 最小限のセットアップでローカルでエンドツーエンドで実行する（`cd backend && go run ./cmd/recipe`、`cd backend && genkit start -- go run ./cmd/recipe`）。
- Dev UI 内でデバッグのためのフロートレースとプロンプト履歴を表示する。
- HTTP クライアント（curl、Postman、BFF）が `/recipeGeneratorFlow` を利用できるようにする。

### 非目標

- 生成されたレシピの永続化（データベース/RAG ストレージはまだなし）。
- フローエンドポイントの認証/認可（外部展開後の将来作業）。
- マルチモデルオーケストレーション（Gemini 2.5 Flash がデフォルトプロバイダー）。

## 3. 機能要件

1. `ingredient` と オプションの `dietaryRestrictions` を含む JSON ペイロードを受け入れる。
2. タイトル、説明、準備/調理時間、人数、材料リスト、手順、オプションのコツを含む構造化された `Recipe` JSON ボディを生成する。
3. 起動時にスモークチェック用のサンプルレシピを stdout にログ出力する。
4. `127.0.0.1:3400` で `POST /recipeGeneratorFlow` を公開する HTTP サーバーをホストする。
5. Genkit Dev UI 経由でのインタラクティブ実行をサポートする（`cd backend && genkit start -- go run ./cmd/recipe`、UI は `http://localhost:4000`）。

## 4. システムアーキテクチャ

### コンポーネント

- **Go サービス (`main.go`)**: Genkit を初期化し、`recipeGeneratorFlow` を定義し、ブートストラップ実行を行い、HTTP ハンドラーをマウントする。
- **Genkit ランタイム**: フロー調整、構造化出力ヘルパー、Dev UI 計装を提供する。
- **Gemini プロバイダー**: `GEMINI_API_KEY` シークレットを使用して `googlegenai` プラグインを通じて LLM 推論を処理する。

### ランタイムフロー

1. クライアントが `RecipeInput` JSON を HTTP エンドポイントまたは Dev UI に送信する。
2. フローが入力フィールドからプロンプト文字列を構築する（`dietaryRestrictions` のデフォルトは `"none"`）。
3. `genkit.GenerateData[Recipe]` が Gemini に構造化出力をリクエストする。
4. レスポンスが `Recipe` 構造体にマーシャルされ、クライアントに返される。トレースが Dev UI に表示される。

### シーケンス図（テキスト形式）

```
クライアント -> フローハンドラー: POST /recipeGeneratorFlow (RecipeInput)
フローハンドラー -> Genkit ランタイム: recipeGeneratorFlow を実行
Genkit ランタイム -> Gemini: 構造化生成リクエスト
Gemini -> Genkit ランタイム: Recipe JSON ペイロード
Genkit ランタイム -> クライアント: HTTP 200 + Recipe
```

## 5. データ契約

```go
// リクエスト
 type RecipeInput struct {
     Ingredient          string `json:"ingredient"`
     DietaryRestrictions string `json:"dietaryRestrictions,omitempty"`
 }

// レスポンス
 type Recipe struct {
     Title        string   `json:"title"`
     Description  string   `json:"description"`
     PrepTime     string   `json:"prepTime"`
     CookTime     string   `json:"cookTime"`
     Servings     int      `json:"servings"`
     Ingredients  []string `json:"ingredients"`
     Instructions []string `json:"instructions"`
     Tips         []string `json:"tips,omitempty"`
 }
```

- **バリデーション**: フローは Genkit のスキーマ推論に依存する。将来的には、不足している `ingredient` 値の手動バリデーションを追加する可能性がある。

## 6. 設定とシークレット

- 環境変数: `GEMINI_API_KEY`（必須）。シェルエクスポートまたは `.env` ツールを通じて注入する。
- オプションの将来のシークレット: 追加プラグインを有効にする際の Firebase 認証情報。
- ポート: `3400`（HTTP フロー）、`4000`（`genkit start` を通じた Dev UI プロキシ）。

## 7. 運用上の考慮事項

- **バージョン固定**: プロジェクトが安定したら `github.com/firebase/genkit/go` をロックして API 変更から保護する。
- **ログ**: 標準ライブラリログ。ローカル環境を超えて展開する場合は構造化ログを検討する。
- **エラーハンドリング**: フローは生成エラーをコンテキストでラップし、HTTP 500 を返す。
- **テスト戦略（計画中）**: プロンプトフォーマット用のテーブル駆動単体テスト。CI がセットアップされたら `Recipe` JSON 形状用のゴールデンテスト。

## 8. ディレクトリレイアウト（提案）

| パス                        | 目的                                                                                    |
| --------------------------- | --------------------------------------------------------------------------------------- |
| `backend/cmd/recipe/main.go`        | 複数のバイナリをスケールする際にエントリーポイントをここに移動する。                    |
| `backend/internal/flows/recipe.go`  | フロー定義とプロンプト構築ヘルパー。                                                    |
| `backend/internal/models/recipe.go` | リクエスト/レスポンススキーマ用の共有構造体。                                           |
| `frontend/`                          | Next.js クライアントアプリ。                                                            |
| `docs/`                     | 仕様書、エージェントガイド（`AGENTS.md`）、アーキテクチャドキュメント（このファイル）。 |
| `testdata/`                 | 統合テスト用のゴールデン JSON レスポンス。                                              |

> **注意**: 2025-09-16 時点で上記レイアウトへ移行済み。追加サービスが増えた場合も `backend/` / `frontend/` 配下で整理する方針。

## 9. 開発ワークフロー

1. `cd backend && go mod tidy` – 依存関係が同期されていることを確認する。
2. `npm install -g genkit-cli` – マシンごとに一度実行。`genkit --version` で確認する。
3. `export GEMINI_API_KEY=...` – フロー実行前に設定する。
4. `go run ./cmd/recipe`（`backend/` 内） – クイックスモークテスト（stdout サンプル出力 + サーバー開始）。
5. `genkit start -- go run ./cmd/recipe`（`backend/` 内） – トレース検査のための Dev UI を有効にする。
6. `curl -X POST http://localhost:3400/recipeGeneratorFlow ...` – 契約を確認する。

## 10. 将来の機能拡張

- スキーマ進化戦略でオプションの `servingSize` と栄養フィールドを追加する。
- ローカル環境を超えて展開する際に認証（API キーヘッダー）を導入する。
- 設定フラグで複数プロバイダー（OpenAI 互換プラグイン）をサポートする。
- プロンプト調整をベンチマークするための評価（`genkit evaluate`）を統合する。

---

_最終更新日: 2025 年 9 月 16 日_
