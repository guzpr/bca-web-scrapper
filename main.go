package main

import (
	"webscrapping/scrapper"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	scrapper.Login()
}
