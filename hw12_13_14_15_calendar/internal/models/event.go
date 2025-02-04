package models

import "github.com/google/uuid"

type Event struct {
	ID           uuid.UUID
	Title        string
	UserID       uuid.UUID
	Description  string
	StartDate    int64
	EndDate      int64
	NotifyBefore int
}

func NewEvent(userID uuid.UUID, title, description string, start, end int64, notifyBefore int) Event {
	return Event{
		ID:           uuid.New(),
		UserID:       userID,
		Title:        title,
		Description:  description,
		StartDate:    start,
		EndDate:      end,
		NotifyBefore: notifyBefore,
	}
}

func (e *Event) HasDifferentDate(event Event) bool {
	return e.StartDate != event.StartDate || e.EndDate != event.EndDate
}

type Schedule []Event

func (s Schedule) Len() int {
	return len(s)
}

func (s Schedule) Less(i, j int) bool {
	return s[i].StartDate < s[j].StartDate
}

func (s Schedule) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
