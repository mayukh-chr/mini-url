curl -X POST http://localhost:8080/shorten -H "Content-Type: application/json" -d '{"url":"https://example.com"}'

randomly generates a short url at http://localhost:8080/u/example

curl -X POST http://localhost:8080/shorten -H "Content-Type: application/json" -d '{"url":"https://example.com", "short_code":"exmp"}'

redirects http://localhost:8080/u/exmp to example.com

curl -v http://localhost:8080/u/exmp

shows where the short link redirects to, returns 404 if no redirects available

curl -X PUT http://localhost:8080/u/{old_code}   -H "Content-Type: application/json"   -d '{
    "url": "https://example.com",
    "short_code": "new_code"
}'

changes the short code

curl -X DELETE http://localhost:8080/u/exmp

deletes the redirect; returns 404 now

curl http://localhost:8080/stats/googl

gives the access count of the url
