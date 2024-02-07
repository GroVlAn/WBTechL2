package api

import (
	"dev11/tools/response"
	"net/http"
)

func (hh *HTTPHandler) createEvent(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		response.Resp(w, nil, nil, http.StatusNotFound)
	}
}

func (hh *HTTPHandler) updateEvent(w http.ResponseWriter, req *http.Request) {}

func (hh *HTTPHandler) deleteEvent(w http.ResponseWriter, req *http.Request) {}

func (hh *HTTPHandler) eventsForDay(w http.ResponseWriter, req *http.Request) {}

func (hh *HTTPHandler) eventsForWeek(w http.ResponseWriter, req *http.Request) {}

func (hh *HTTPHandler) eventsForMonth(w http.ResponseWriter, req *http.Request) {}
