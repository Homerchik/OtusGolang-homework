package models

import "github.com/google/uuid"

type Event struct {
	ID           uuid.UUID `json:"id,omitempty"`
	Title        string    `json:"title"`
	UserID       uuid.UUID `json:"userId"`
	Description  string    `json:"description"`
	StartDate    int64     `json:"startDate"`
	EndDate      int64     `json:"endDate"`
	NotifyBefore int       `json:"notifyBefore"`
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
