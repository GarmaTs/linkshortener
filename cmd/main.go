package main

import (
	"fmt"

	l1 "github.com/GarmaTs/linkshortener/internal/lesson1"
)

func main() {
	greetStr := l1.Greet("Garma")
	fmt.Println(greetStr)
}
