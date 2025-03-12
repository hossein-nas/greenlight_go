#!/bin/bash

# movie_seeder_250.sh
# Seeds 250 movies into a local API endpoint using TMDb API

echo "Starting movie data seeding with TMDb API..."

# TMDb API key (replace with your own)
TMDB_API_KEY="a24b06047bd34ca165459534f85574b8"

# Check for required tools
command -v curl >/dev/null 2>&1 || { echo "curl is required but not installed. Aborting."; exit 1; }
command -v jq >/dev/null 2>&1 || { echo "jq is required but not installed. Aborting."; exit 1; }

# Base URL for TMDb popular movies endpoint
BASE_URL="https://api.themoviedb.org/3/movie/popular?api_key=$TMDB_API_KEY&language=en-US"

# Counter for total movies seeded
TOTAL_SEEDED=0
MAX_ITEMS=250
PAGE=1

# Loop until we have 250 items
while [ $TOTAL_SEEDED -lt $MAX_ITEMS ]; do
    # Fetch movie data for the current page
    RESPONSE=$(curl -s "$BASE_URL&page=$PAGE")
    
    # Check if API call was successful
    if [ $? -ne 0 ]; then
        echo "Error fetching data from TMDb API on page $PAGE. Exiting."
        exit 1
    fi

    # Extract movie list from response
    MOVIES=$(echo "$RESPONSE" | jq -c '.results[]')

    # Process each movie in the current page
    echo "$MOVIES" | while IFS= read -r MOVIE; do
        if [ $TOTAL_SEEDED -ge $MAX_ITEMS ]; then
            break
        fi

        # Extract fields and map to your schema
        TITLE=$(echo "$MOVIE" | jq -r '.title')
        YEAR=$(echo "$MOVIE" | jq -r '.release_date | split("-")[0]')
        RUNTIME=$(curl -s "https://api.themoviedb.org/3/movie/$(echo "$MOVIE" | jq -r '.id')?api_key=$TMDB_API_KEY" | jq -r '.runtime // 120') # Fallback to 120 if null
        GENRES=$(echo "$MOVIE" | jq -r '[.genre_ids[] | if . == 28 then "Action" elif . == 12 then "Adventure" elif . == 18 then "Drama" elif . == 35 then "Comedy" elif . == 878 then "Sci-Fi" elif . == 53 then "Thriller" else "Other" end] | join(",")')

        # Construct JSON payload
        PAYLOAD=$(jq -n \
            --arg title "$TITLE" \
            --arg year "$YEAR" \
            --arg runtime "$RUNTIME mins" \
            --arg genres "$GENRES" \
            '{title: $title, year: ($year | tonumber), runtime: $runtime, genres: ($genres | split(","))}')

        # POST to local API
        curl --request POST \
          --url http://localhost:4000/v1/movies \
          --header 'Accept: */*' \
          --header 'Accept-Encoding: gzip, deflate, br' \
          --header 'Connection: keep-alive' \
          --header 'Content-Type: application/json' \
          --header 'User-Agent: EchoapiRuntime/1.1.0' \
          --data "$PAYLOAD" >/dev/null 2>&1

        if [ $? -eq 0 ]; then
            echo "Seeded #$((TOTAL_SEEDED + 1)): $TITLE"
            TOTAL_SEEDED=$((TOTAL_SEEDED + 1))
        else
            echo "Failed to seed: $TITLE"
        fi

        # Small delay to avoid overwhelming the server
        sleep 0.05
    done

    # Increment page for next batch
    PAGE=$((PAGE + 1))

    # Check if we've exhausted available pages (TMDb has a max of 500 pages)
    if [ $PAGE -gt 500 ]; then
        echo "Reached TMDb page limit. Seeded $TOTAL_SEEDED movies."
        break
    fi
done

echo "Seeding completed! Total movies seeded: $TOTAL_SEEDED"