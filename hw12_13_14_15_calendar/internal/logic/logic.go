package logic

import (
	"fmt"
	"sort"
	"time"

	"github.com/homerchik/OtusGolang-homework/hw12_13_14_15_calendar/internal/models"
)

func EventAfterNow(event models.Event) error {
	if event.StartDate < time.Now().UTC().Unix() {
		fmt.Println(event.StartDate)
		fmt.Println(time.Now().UTC().Unix())
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

func MergeEvents(basic, updated models.Event) models.Event {
	if updated.Title != "" {
		basic.Title = updated.Title
	}
	if updated.NotifyBefore != 0 {
		basic.NotifyBefore = updated.NotifyBefore
	}
	if updated.Description != "" {
		basic.Description = updated.Description
	}
	if updated.StartDate != 0 {
		basic.StartDate = updated.StartDate
	}
	if updated.EndDate != 0 {
		basic.EndDate = updated.EndDate
	}
	return basic
}

func BuildNotification(e models.Event) models.Notification {
	return models.Notification{ID: e.ID, UserID: e.UserID, Title: e.Title, Date: e.StartDate}
}
