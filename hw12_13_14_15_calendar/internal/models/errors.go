package models

import "errors"

var (
	ErrStartTimeBeforeNow = errors.New("start time is before now")
	ErrEventIntersection  = errors.New("event intersection")
	ErrNoEventFound       = errors.New("no event found")
	ErrEventCantBeAdded   = errors.New("event can't be added")
	ErrEventCantBeUpdated = errors.New("event can't be updated")
	ErrEventCantBeDeleted = errors.New("event can't be deleted")
)
