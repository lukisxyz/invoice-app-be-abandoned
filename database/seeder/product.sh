#!/bin/bash

# Function to generate a random string
generate_random_string() {
  head /dev/urandom | tr -dc A-Za-z0-9 | head -c $1
  echo
}

# Function to generate random data
generate_random_data() {
  echo "{"
  echo "\"Sku\": \"$(generate_random_string 8)\","
  echo "\"Barcode\": \"$(generate_random_string 12)\","
  echo "\"Name\": \"$(generate_random_string 10)\","
  echo "\"Description\": \"$(generate_random_string 20)\","
  echo "\"Image\": null,"
  echo "\"Amount\": $(awk -v min=10 -v max=1000 'BEGIN{srand(); print min+rand()*(max-min)}')"
  echo "}"
}

# Function to send an HTTP POST request
send_post_request() {
  local url=$1
  local data=$2
  curl -X POST -H "Content-Type: application/json" -d "$data" "$url"
}

# Generate 20 random data entries
for _ in {1..20}
do
  random_data=$(generate_random_data)
  send_post_request "http://127.0.0.1:8080/api/product" "$random_data"
done
