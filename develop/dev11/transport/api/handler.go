package api

import (
	"dev11/core"
	"dev11/middleware"
	"net/http"
	"time"
)

type HTTPHandler struct {
	mux *http.ServeMux
	sv  EventServ
}

func NewHTTPHandler(sv EventServ) *HTTPHandler {
	return &HTTPHandler{
		mux: http.NewServeMux(),
		sv:  sv,
	}
}

func (hh *HTTPHandler) Handler() http.Handler {
	hh.handle("/create_event/", hh.createEvent)
	hh.handle("/update_event/", hh.updateEvent)
	hh.handle("/delete_event/", hh.deleteEvent)
	hh.handle("/events_for_day/", hh.eventsForDay)
	hh.handle("/events_for_week/", hh.eventsForWeek)
	hh.handle("/events_for_month/", hh.eventsForMonth)

	return hh.mux
}

func (hh *HTTPHandler) handle(path string, hf http.HandlerFunc) {
	hh.mux.Handle(path, middleware.Middleware(hf))
}

type EventServ interface {
	Create(event core.Event)
	Update(event core.Event) error
	Delete(evID, userID string) error
	EventByDay(day time.Time, userID string) ([]core.Event, error)
	EventByWeek(since time.Time, userID string) ([]core.Event, error)
	EventByMonth(since time.Time, userID string) ([]core.Event, error)
}
