"use client";

// =============================================================================
//
// The whole frontend is basically is in this one file. It is split into three parts:
//   1. Types and data - the shape of the data, and the list of serveries
//   2. API calls      - the two requests this page makes to the backend
//   3. The page       - what is shown on screen and what happens on clicks
//
// Every request goes to a URL that starts with /api/. next.config.ts forwards
// those requests, unchanged, to the backend at http://localhost:8080.
// =============================================================================

import { useEffect, useState } from "react";

// -----------------------------------------------------------------------------
// 1. Types and data
// -----------------------------------------------------------------------------

// A Servery. The id is what we send to the backend; the name is shown on screen.
type Servery = {
  id: string;
  name: string;
};

// The list of serveries, in the order we want to show them on screen.
const SERVERIES: Servery[] = [
  { id: "baker", name: "Baker Kitchen" },
  { id: "north", name: "North Servery" },
  { id: "seibel", name: "Seibel Servery" },
  { id: "south", name: "South Servery" },
  { id: "south-main", name: "South Main Servery" },
  { id: "west", name: "West Servery" },
];

// The value of the "+ New dish" option in the dish dropdown.
const NEW_DISH = "new-dish";

// A Review Object, We send this shape to the backend when posting a review,
// and the backend sends a list of them back when we load reviews.
type Review = {
  servery: string; // the servery's id, e.g. "north"
  dish: string; // the dish's name, e.g. "Pad Thai"
  rating: number; // 1 to 5
  comment: string;
};

// A dish and all of its reviews. We build these ourselves from the review list
// (see groupReviewsByDish below); the backend never sends a Dish.
type Dish = {
  name: string;
  reviews: Review[];
};

// -----------------------------------------------------------------------------
// 2. API calls
// -----------------------------------------------------------------------------

// GET /api/reviews?servery=north
// Returns every review for one servery.
async function getReviews(serveryId: string): Promise<Review[]> {
  const response = await fetch(`/api/reviews?servery=${serveryId}`, {
    method: "GET",
  });
  await throwIfFailed(response);

  const reviews = await response.json();
  return reviews ?? []; // treat a `null` response as an empty list
}

// POST /api/reviews
// Sends a new review to the backend as JSON, e.g.
//   { "servery": "north", "dish": "Pad Thai", "rating": 4, "comment": "Pretty good" }
async function postReview(review: Review): Promise<void> {
  const response = await fetch("/api/reviews", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(review),
  });
  await throwIfFailed(response);
}

// If the backend answered with an error status (anything other than 2xx),
// stop here and show the backend's message on the page.
async function throwIfFailed(response: Response): Promise<void> {
  if (!response.ok) {
    const message = await response.text();
    throw new Error(`${response.status} ${response.statusText}: ${message}`);
  }
}

// -----------------------------------------------------------------------------
// 3. The page
// -----------------------------------------------------------------------------

export default function ServeryReviewsPage() {
  // ---- Reviews from the backend ----
  const [reviews, setReviews] = useState<Review[]>([]);

  // ---- Which servery the user is looking at (its id, e.g. "baker") ----
  const [selectedServeryId, setSelectedServeryId] = useState(SERVERIES[0].id);

  // ---- What the user has entered in the review form ----
  // dishChoice is the dish picked in the dropdown, or NEW_DISH if they are adding
  // a dish. In that case, newDishName is the name they typed.
  const [dishChoice, setDishChoice] = useState("");
  const [newDishName, setNewDishName] = useState("");
  const [rating, setRating] = useState(5);
  const [comment, setComment] = useState("");

  // ---- An error message to show at the top of the page ("" means no error) ----
  const [errorMessage, setErrorMessage] = useState("");

  // When the page opens, and whenever the selected servery changes,
  // load that servery's reviews.
  useEffect(() => {
    loadReviews(selectedServeryId);
  }, [selectedServeryId]);

  async function loadReviews(serveryId: string) {
    try {
      const serveryReviews = await getReviews(serveryId);
      setReviews(serveryReviews);
      setErrorMessage("");
    } catch (error) {
      showError(error);
    }
  }

  // Runs when the user clicks a servery button.
  function selectServery(serveryId: string) {
    setSelectedServeryId(serveryId);
    setDishChoice(""); // each servery has different dishes, so clear the old choice
  }

  // Runs when the user clicks "Post review".
  async function handleSubmit(event: React.SubmitEvent<HTMLFormElement>) {
    event.preventDefault(); // stop the browser from reloading the page

    // Use the dish picked in the dropdown, or the typed name if it's a new dish.
    let dishName = dishChoice;
    if (dishChoice === NEW_DISH) {
      dishName = newDishName;
    }

    const newReview: Review = {
      servery: selectedServeryId,
      dish: dishName,
      rating: rating,
      comment: comment,
    };

    try {
      await postReview(newReview);

      // It worked: clear the form and reload the reviews so the new one shows up.
      setDishChoice("");
      setNewDishName("");
      setRating(5);
      setComment("");
      await loadReviews(selectedServeryId);
    } catch (error) {
      showError(error);
    }
  }

  function showError(error: unknown) {
    setErrorMessage(String(error));
  }

  // The reviews come back as one flat list. Group them by dish so each dish
  // gets its own card on the page.
  const dishes = groupReviewsByDish(reviews);

  return (
    <main>
      <h1>Servery Reviews</h1>
      <p className="subtitle">Rate the food at Rice&apos;s serveries.</p>

      {errorMessage !== "" && (
        <div className="error">
          <strong>Something went wrong talking to the backend:</strong>
          <pre>{errorMessage}</pre>
        </div>
      )}

      {/* ---- One button per servery ---- */}
      <nav className="servery-buttons">
        {SERVERIES.map((servery) => (
          <button
            key={servery.id}
            className={servery.id === selectedServeryId ? "selected" : ""}
            onClick={() => selectServery(servery.id)}
          >
            {servery.name}
          </button>
        ))}
      </nav>

      {/* ---- The review form ---- */}
      <form className="card" onSubmit={handleSubmit}>
        <h2>Write a review</h2>

        {/* Pick a dish someone already reviewed at this servery, or "+ New dish". */}
        <label>
          Dish
          <select value={dishChoice} onChange={(event) => setDishChoice(event.target.value)}>
            <option value="" disabled>
              Choose a dish
            </option>
            {dishes.map((dish) => (
              <option key={dish.name} value={dish.name}>
                {dish.name}
              </option>
            ))}
            <option value={NEW_DISH}>+ New dish</option>
          </select>
        </label>

        {/* Only shown when adding a new dish. The new dish shows up on the page
            once it has its first review. */}
        {dishChoice === NEW_DISH && (
          <label>
            New dish name
            <input
              value={newDishName}
              onChange={(event) => setNewDishName(event.target.value)}
              placeholder="e.g. Pad Thai"
            />
          </label>
        )}

        <label>
          Rating
          <select value={rating} onChange={(event) => setRating(Number(event.target.value))}>
            <option value={5}>5 - Amazing</option>
            <option value={4}>4 - Good</option>
            <option value={3}>3 - Okay</option>
            <option value={2}>2 - Bad</option>
            <option value={1}>1 - Awful</option>
          </select>
        </label>

        <label>
          Comment
          <textarea
            value={comment}
            onChange={(event) => setComment(event.target.value)}
            placeholder="What did you think?"
          />
        </label>

        <button type="submit">Post review</button>
      </form>

      {/* ---- The reviews, grouped by dish ---- */}
      {dishes.length === 0 && <p className="empty">No reviews here yet. Be the first!</p>}

      {dishes.map((dish) => (
        <section className="card" key={dish.name}>
          <h2>{dish.name}</h2>
          <p className="average">
            {stars(averageRating(dish.reviews))} {averageRating(dish.reviews).toFixed(1)} from{" "}
            {dish.reviews.length} review(s)
          </p>

          <ul>
            {dish.reviews.map((review, index) => (
              <li key={index}>
                <span className="stars">{stars(review.rating)}</span> {review.comment}
              </li>
            ))}
          </ul>
        </section>
      ))}
    </main>
  );
}

// -----------------------------------------------------------------------------
// Small helpers used by the page
// -----------------------------------------------------------------------------

// Turns a flat list of reviews into a list of dishes, each with its reviews.
function groupReviewsByDish(reviews: Review[]): Dish[] {
  const dishes: Dish[] = [];

  for (const review of reviews) {
    let dish = dishes.find((existingDish) => existingDish.name === review.dish);

    // First review we've seen for this dish: make a new Dish for it.
    if (dish === undefined) {
      dish = { name: review.dish, reviews: [] };
      dishes.push(dish);
    }

    dish.reviews.push(review);
  }

  return dishes;
}

// The average of all the ratings, e.g. ratings 4 and 5 -> 4.5
function averageRating(reviews: Review[]): number {
  let total = 0;
  for (const review of reviews) {
    total += review.rating;
  }
  return total / reviews.length;
}

// Shows a rating as stars, e.g. 3 -> "★★★☆☆"
function stars(rating: number): string {
  const filled = Math.min(5, Math.max(0, Math.round(rating))); // keep it between 0 and 5
  return "★".repeat(filled) + "☆".repeat(5 - filled);
}
