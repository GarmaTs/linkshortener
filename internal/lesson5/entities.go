package lesson5

import (
	"log"
)

type Item struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
	ImageLink   string `json:"image_link"`
}

type User struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Config struct {
	Version string
	Port    int
}

type Application struct {
	Items  []Item
	Config Config
	Logger *log.Logger
}
