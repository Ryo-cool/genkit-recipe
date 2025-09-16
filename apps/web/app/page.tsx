"use client";

import { FormEvent, useMemo, useState } from "react";

type Recipe = {
  title: string;
  description: string;
  prepTime: string;
  cookTime: string;
  servings: number;
  ingredients: string[];
  instructions: string[];
  tips?: string[];
};

function emptyRecipe(): Recipe {
  return {
    title: "",
    description: "",
    prepTime: "",
    cookTime: "",
    servings: 0,
    ingredients: [],
    instructions: []
  };
}

const defaultBase = "http://localhost:3400";

export default function Home() {
  const [ingredient, setIngredient] = useState("avocado");
  const [dietaryRestrictions, setDietaryRestrictions] = useState("vegetarian");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [recipe, setRecipe] = useState<Recipe | null>(null);

  const endpoint = useMemo(() => {
    const base = process.env.NEXT_PUBLIC_API_BASE ?? defaultBase;
    return `${base.replace(/\/$/, "")}/recipeGeneratorFlow`;
  }, []);

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setLoading(true);
    setError(null);
    setRecipe(null);

    try {
      const response = await fetch(endpoint, {
        method: "POST",
        headers: {
          "Content-Type": "application/json"
        },
        body: JSON.stringify({
          data: {
            ingredient,
            dietaryRestrictions: dietaryRestrictions.trim() || undefined
          }
        })
      });

      if (!response.ok) {
        throw new Error(`Request failed with status ${response.status}`);
      }

      const json = (await response.json()) as { data?: Recipe };
      const generated = json.data ?? emptyRecipe();
      if (!generated.title) {
        throw new Error("Response payload was empty");
      }

      setRecipe(generated);
    } catch (err) {
      const message = err instanceof Error ? err.message : "Unknown error";
      setError(message);
    } finally {
      setLoading(false);
    }
  }

  return (
    <main
      style={{
        maxWidth: "960px",
        margin: "0 auto",
        padding: "4rem 1.5rem",
        display: "grid",
        gap: "2rem"
      }}
    >
      <section
        style={{
          backgroundColor: "white",
          borderRadius: "16px",
          padding: "2rem",
          boxShadow: "0 20px 35px rgba(15, 23, 42, 0.08)"
        }}
      >
        <header style={{ marginBottom: "1.5rem" }}>
          <h1
            style={{
              margin: 0,
              fontSize: "2.25rem",
              lineHeight: 1.1
            }}
          >
            Genkit Recipe Studio
          </h1>
          <p style={{ marginTop: "0.5rem", color: "#475569" }}>
            Generate structured recipes with Gemini 2.5 Flash and inspect them in
            the Dev UI while you iterate on prompts.
          </p>
        </header>

        <form
          onSubmit={handleSubmit}
          style={{
            display: "grid",
            gap: "1rem"
          }}
        >
          <label style={{ display: "grid", gap: "0.35rem" }}>
            <span style={{ fontWeight: 600 }}>Main ingredient or cuisine</span>
            <input
              value={ingredient}
              onChange={(event) => setIngredient(event.target.value)}
              placeholder="tomato"
              required
              style={{
                padding: "0.75rem 1rem",
                borderRadius: "12px",
                border: "1px solid #CBD5F5",
                fontSize: "1rem"
              }}
            />
          </label>

          <label style={{ display: "grid", gap: "0.35rem" }}>
            <span style={{ fontWeight: 600 }}>Dietary restrictions (optional)</span>
            <input
              value={dietaryRestrictions}
              onChange={(event) => setDietaryRestrictions(event.target.value)}
              placeholder="gluten-free"
              style={{
                padding: "0.75rem 1rem",
                borderRadius: "12px",
                border: "1px solid #CBD5F5",
                fontSize: "1rem"
              }}
            />
          </label>

          <button
            type="submit"
            disabled={loading}
            style={{
              marginTop: "0.5rem",
              padding: "0.9rem 1.4rem",
              borderRadius: "999px",
              border: "none",
              background: loading ? "#94a3b8" : "linear-gradient(120deg, #2563eb, #7c3aed)",
              color: "white",
              fontWeight: 600,
              fontSize: "1rem",
              cursor: loading ? "not-allowed" : "pointer",
              transition: "opacity 0.2s ease"
            }}
          >
            {loading ? "Generating..." : "Generate recipe"}
          </button>
        </form>

        {error && (
          <p
            role="alert"
            style={{
              marginTop: "1rem",
              color: "#dc2626",
              fontWeight: 500
            }}
          >
            {error}
          </p>
        )}
      </section>

      <section
        aria-live="polite"
        style={{
          backgroundColor: "white",
          borderRadius: "16px",
          padding: "2rem",
          boxShadow: "0 20px 35px rgba(15, 23, 42, 0.08)",
          minHeight: "320px"
        }}
      >
        <h2
          style={{
            margin: 0,
            fontSize: "1.5rem"
          }}
        >
          Recipe output
        </h2>
        {!recipe && !error && !loading && (
          <p style={{ marginTop: "1rem", color: "#64748b" }}>
            Submit an ingredient to see a structured recipe response.
          </p>
        )}

        {recipe && (
          <article style={{ marginTop: "1.5rem", display: "grid", gap: "1.5rem" }}>
            <header>
              <h3 style={{ margin: 0, fontSize: "1.75rem" }}>{recipe.title}</h3>
              <p style={{ marginTop: "0.5rem", color: "#475569" }}>{recipe.description}</p>
              <dl
                style={{
                  marginTop: "1rem",
                  display: "grid",
                  gridTemplateColumns: "repeat(auto-fit, minmax(160px, 1fr))",
                  gap: "0.75rem",
                  color: "#334155",
                  fontSize: "0.95rem"
                }}
              >
                <div>
                  <dt style={{ fontWeight: 600 }}>Prep time</dt>
                  <dd>{recipe.prepTime}</dd>
                </div>
                <div>
                  <dt style={{ fontWeight: 600 }}>Cook time</dt>
                  <dd>{recipe.cookTime}</dd>
                </div>
                <div>
                  <dt style={{ fontWeight: 600 }}>Servings</dt>
                  <dd>{recipe.servings}</dd>
                </div>
              </dl>
            </header>

            <section>
              <h4 style={{ margin: 0, fontSize: "1.2rem" }}>Ingredients</h4>
              <ul
                style={{
                  marginTop: "0.75rem",
                  paddingLeft: "1.25rem",
                  display: "grid",
                  gap: "0.5rem"
                }}
              >
                {recipe.ingredients.map((item, index) => (
                  <li key={index}>{item}</li>
                ))}
              </ul>
            </section>

            <section>
              <h4 style={{ margin: 0, fontSize: "1.2rem" }}>Instructions</h4>
              <ol
                style={{
                  marginTop: "0.75rem",
                  paddingLeft: "1.25rem",
                  display: "grid",
                  gap: "0.75rem"
                }}
              >
                {recipe.instructions.map((item, index) => (
                  <li key={index}>{item}</li>
                ))}
              </ol>
            </section>

            {recipe.tips && recipe.tips.length > 0 && (
              <section>
                <h4 style={{ margin: 0, fontSize: "1.2rem" }}>Tips</h4>
                <ul
                  style={{
                    marginTop: "0.75rem",
                    paddingLeft: "1.25rem",
                    display: "grid",
                    gap: "0.5rem"
                  }}
                >
                  {recipe.tips.map((tip, index) => (
                    <li key={index}>{tip}</li>
                ))}
                </ul>
              </section>
            )}
          </article>
        )}
      </section>
    </main>
  );
}
