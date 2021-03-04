package domain

import "time"

type User struct {
	ID           int       `json:"-"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	RegisteredAt time.Time `json:"registeredAt"`
	LastVisitAt  time.Time `json:"lastVisitAt"`
}
