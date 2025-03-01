#!/bin/bash

# movie_seeder.sh
# A script to seed movie data into a local API endpoint using curl

echo "Starting movie data seeding..."

# Movie 1: The Silent Echo
curl --request POST \
  --url http://localhost:4000/v1/movies \
  --header 'Accept: */*' \
  --header 'Accept-Encoding: gzip, deflate, br' \
  --header 'Connection: keep-alive' \
  --header 'Content-Type: application/json' \
  --header 'User-Agent: EchoapiRuntime/1.1.0' \
  --data '{
	"title": "The Silent Echo",
	"year": 2018,
	"genres": [
		"Drama",
		"Mystery"
	],
	"runtime": "115 mins"
}'
echo "Seeded: The Silent Echo"
sleep 1  # 1-second delay between requests

# Movie 2: Neon Horizon
curl --request POST \
  --url http://localhost:4000/v1/movies \
  --header 'Accept: */*' \
  --header 'Accept-Encoding: gzip, deflate, br' \
  --header 'Connection: keep-alive' \
  --header 'Content-Type: application/json' \
  --header 'User-Agent: EchoapiRuntime/1.1.0' \
  --data '{
	"title": "Neon Horizon",
	"year": 2023,
	"genres": [
		"Sci-Fi",
		"Action"
	],
	"runtime": "145 mins"
}'
echo "Seeded: Neon Horizon"
sleep 1

# Movie 3: Whispers in the Fog
curl --request POST \
  --url http://localhost:4000/v1/movies \
  --header 'Accept: */*' \
  --header 'Accept-Encoding: gzip, deflate, br' \
  --header 'Connection: keep-alive' \
  --header 'Content-Type: application/json' \
  --header 'User-Agent: EchoapiRuntime/1.1.0' \
  --data '{
	"title": "Whispers in the Fog",
	"year": 2015,
	"genres": [
		"Thriller"
	],
	"runtime": "98 mins"
}'
echo "Seeded: Whispers in the Fog"
sleep 1

# Movie 4: Desert Runner
curl --request POST \
  --url http://localhost:4000/v1/movies \
  --header 'Accept: */*' \
  --header 'Accept-Encoding: gzip, deflate, br' \
  --header 'Connection: keep-alive' \
  --header 'Content-Type: application/json' \
  --header 'User-Agent: EchoapiRuntime/1.1.0' \
  --data '{
	"title": "Desert Runner",
	"year": 2021,
	"genres": [
		"Adventure",
		"Drama"
	],
	"runtime": "127 mins"
}'
echo "Seeded: Desert Runner"
sleep 1

# Movie 5: Crimson Tide
curl --request POST \
  --url http://localhost:4000/v1/movies \
  --header 'Accept: */*' \
  --header 'Accept-Encoding: gzip, deflate, br' \
  --header 'Connection: keep-alive' \
  --header 'Content-Type: application/json' \
  --header 'User-Agent: EchoapiRuntime/1.1.0' \
  --data '{
	"title": "Crimson Tide",
	"year": 1995,
	"genres": [
		"Action",
		"Thriller"
	],
	"runtime": "116 mins"
}'
echo "Seeded: Crimson Tide"

echo "Seeding completed!"