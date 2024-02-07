package service

import (
	"dev11/core"
	"fmt"
	"time"
)

var Events = make([]core.Event, 0)
var Users = []core.User{
	core.User{
		ID:       "3",
		UserName: "vlad",
	},
}

type EventsService struct{}

func NewEventsService() *EventsService {
	return &EventsService{}
}

func (es *EventsService) Create(event core.Event) {
	Events = append(Events, event)
}

func (es *EventsService) Update(event core.Event) error {
	eventKey := es.foundEvent(event.ID, event.UserID)

	if eventKey == -1 {
		return fmt.Errorf("event: event by id %s and user id %s is not found", event.ID, event.UserID)
	}

	Events[eventKey] = event

	return nil
}

func (es *EventsService) Delete(evID, userID string) error {
	eventKey := es.foundEvent(evID, userID)

	if eventKey == -1 {
		return fmt.Errorf("event: event by id %s and user id %s is not found", evID, userID)
	}

	Events = append(Events[:eventKey], Events[eventKey+1:]...)

	return nil
}

func (es *EventsService) EventByDay(day time.Time, userID string) ([]core.Event, error) {
	eventsByDay := make([]core.Event, 0)
	for _, ev := range Events {
		if ev.Date == day && ev.UserID == userID {
			eventsByDay = append(eventsByDay, ev)
		}
	}

	if len(eventsByDay) == 0 {
		return nil, fmt.Errorf("event: have not events by day %s", day)
	}

	return eventsByDay, nil
}

func (es *EventsService) EventByWeek(since time.Time, userID string) ([]core.Event, error) {
	eventsByWeek := make([]core.Event, 0)
	forDay := since.AddDate(0, 0, since.Day()+7)

	for _, ev := range Events {
		if (ev.Date.After(since) && ev.Date.Before(forDay) || ev.Date == since) && ev.UserID == userID {
			eventsByWeek = append(eventsByWeek, ev)
		}
	}
	if len(eventsByWeek) == 0 {
		return nil, fmt.Errorf("event: have not events by week since %s", since)
	}

	return eventsByWeek, nil
}

func (es *EventsService) EventByMonth(since time.Time, userID string) ([]core.Event, error) {
	eventsByWeek := make([]core.Event, 0)
	forDay := since.AddDate(0, int(since.Month())+30, 0)

	for _, ev := range Events {
		if (ev.Date.After(since) && ev.Date.Before(forDay) || ev.Date == since) && ev.UserID == userID {
			eventsByWeek = append(eventsByWeek, ev)
		}
	}
	if len(eventsByWeek) == 0 {
		return nil, fmt.Errorf("event: have not events by week since %s", since)
	}

	return eventsByWeek, nil
}

func (es *EventsService) foundEvent(id, userID string) int {
	eventKey := -1

	for i, ev := range Events {
		if ev.ID == id && ev.UserID == userID {
			eventKey = i
		}
	}

	return eventKey
}
