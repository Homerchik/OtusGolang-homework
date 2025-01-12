package storage

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrStartTimeBeforeNow = errors.New("start time is before now")
	ErrEventIntersection   = errors.New("event intersection")
	ErrNoEventFound 	  = errors.New("no event found")
	ErrEventCantBeAdded   = errors.New("event can't be added")
	ErrEventCantBeUpdated = errors.New("event can't be updated")
	ErrEventCantBeDeleted = errors.New("event can't be deleted")
)

type Event struct {
	ID    uuid.UUID
	Title string
	UserId uuid.UUID
	Description string
	StartDate time.Time
	EndDate time.Time
	NotifyBefore int
}

func NewEvent(userId uuid.UUID, title, description string, start, end time.Time, notifyBefore time.Duration) Event {
	return Event{
		ID: uuid.New(),
		UserId: userId,
		Title: title,
		Description: description,
		StartDate: start.UTC(),
		EndDate: end.UTC(),
		NotifyBefore: int(notifyBefore.Seconds()),
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
	return s[i].StartDate.Before(s[j].StartDate)
}

func (s Schedule) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}