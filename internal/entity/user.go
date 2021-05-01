package entity

import "time"

type User struct {
	ID           int       `json:"id"`
	FirstName    string    `json:"firstName"`
	LastName     string    `json:"lastName"`
	BirthDate    string    `json:"birthDate"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	InSearch     bool      `json:"inSearch"`
	RegisteredAt time.Time `json:"registeredAt"`
	ImageURL     string    `json:"imageURL"`
}
