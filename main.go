package main

import (
	"fmt"
	"urlshortner/database"
	)

func main() {
	fmt.Println("Hello, world!")
	database.InitDB("urlshortener.db")
}