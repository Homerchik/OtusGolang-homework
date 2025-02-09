package models

import (
	"github.com/google/uuid"
)

type Notification struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"userId"`
	Title  string    `json:"title"`
	Date   int64     `json:"date"`
}

func (n *Notification) String() string {
	return n.ID.String()
}
