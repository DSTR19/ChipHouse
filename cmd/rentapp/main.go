package main

import (
	"fmt"

	"ChipHouse/internal/app"
)

func main() {
	a, err := app.NewApp()
	if err != nil {
		panic(err)
	}
	defer a.DB.Pool.Close()

	fmt.Println("ChipHouse started with .env config")
}
