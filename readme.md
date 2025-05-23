url -X POST http://localhost:8080/shorten -H "Content-Type: application/json" -d "{\"url\":\"https://example.com\"}" 

randomly generates a short url at http://localhost:8080/u/example

curl -X POST http://localhost:8080/shorten -H "Content-Type: application/json" -d "{\"url\":\"https://example.com\", \"short_code\":\"exmp\"}"

redirects http://localhost:8080/u/exmp to example.com

