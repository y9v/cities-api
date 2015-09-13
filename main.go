package main

import (
	"fmt"
	"log"
)

func main() {
	fmt.Printf("* Listening on port %s\n\n", configuration.Port)
	log.Fatal(Server().ListenAndServe())
}

func init() {
	fmt.Println("* Booting cities service...")

	fmt.Println("* Loading configuration...")
	LoadConfiguration(&configuration)

	fmt.Println("* Connecting to the database...")
	InitDBSession()

	go ParseCities()
}
