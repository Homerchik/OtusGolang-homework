package logic

import (
	"sort"
	"time"

	"github.com/homerchik/OtusGolang-homework/hw12_13_14_15_calendar/internal/models"
)

func EventAfterNow(event models.Event) error {
	if event.StartDate < time.Now().Unix() {
		return models.ErrStartTimeBeforeNow
	}
	return nil
}

func EventFitsSchedule(event models.Event, schedule models.Schedule) error {
	if len(schedule) == 0 {
		return nil
	}
	sort.Sort(schedule)
	for _, scheduleEvent := range schedule {
		if scheduleEvent.StartDate < event.StartDate && scheduleEvent.EndDate < event.StartDate {
			continue
		}
		if scheduleEvent.StartDate > event.StartDate && scheduleEvent.StartDate > event.EndDate {
			return nil
		}
		return models.ErrEventIntersection
	}
	return nil
}

func CheckEvent(event models.Event, schedule models.Schedule) error {
	if err := EventAfterNow(event); err != nil {
		return err
	}
	if err := EventFitsSchedule(event, schedule); err != nil {
		return err
	}
	return nil
}
