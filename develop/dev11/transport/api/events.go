package api

import (
	"dev11/core"
	"dev11/tools/response"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func (hh *HTTPHandler) createEvent(w http.ResponseWriter, req *http.Request) {
	nameMethod := "create evet"
	if req.Method != http.MethodPost {
		response.Resp(w, nil, nil, http.StatusNotFound)
	}

	var newEvent core.Event
	body := req.Body
	defer func() {
		if err := body.Close(); err != nil {
			log.Printf("%s: can not close body: %s\n", nameMethod, err.Error())
		}
	}()

	errD := json.NewDecoder(body).Decode(&newEvent)

	if errD != nil {
		response.Resp(w,
			nil,
			fmt.Errorf("%s: bad event: %s\n", nameMethod, errD.Error()),
			http.StatusBadRequest,
		)
		return
	}

	errC := hh.sv.Create(newEvent)

	if errC != nil {
		response.Resp(w, nil, errC, http.StatusBadRequest)
		return
	}

	resp := core.SuccessResponse{Result: "ok"}
	log.Printf("%s: event %s is created\n", nameMethod, newEvent)

	response.Resp(w, resp, nil, http.StatusCreated)
}

func (hh *HTTPHandler) updateEvent(w http.ResponseWriter, req *http.Request) {
	nameMethod := "update event"
	if req.Method != http.MethodPut {
		response.Resp(w, nil, nil, http.StatusNotFound)
	}
	var newEvent core.Event
	body := req.Body
	defer func() {
		if err := body.Close(); err != nil {
			log.Printf("%s: can not close body: %s\n", nameMethod, err.Error())
		}
	}()

	errD := json.NewDecoder(body).Decode(&newEvent)

	if errD != nil {
		errResp := core.ErrorResponse{
			Error: fmt.Errorf("%s: bad event: %s\n", nameMethod, errD.Error()),
		}
		response.Resp(w,
			nil,
			errResp,
			http.StatusBadRequest,
		)
		return
	}

	errC := hh.sv.Update(newEvent)

	if errC != nil {
		errResp := core.ErrorResponse{
			Error: errC.Error(),
		}
		response.Resp(w, nil, errResp, http.StatusBadRequest)
		return
	}

	resp := core.SuccessResponse{Result: "ok"}
	log.Printf("%s: event %s is created\n", nameMethod, newEvent)

	response.Resp(w, resp, nil, http.StatusCreated)
}

func (hh *HTTPHandler) deleteEvent(w http.ResponseWriter, req *http.Request) {
	nameMethod := "delete event"
	if req.Method != http.MethodDelete {
		response.Resp(w, nil, nil, http.StatusNotFound)
	}

	var delEvent struct {
		ID     string `json:"id"`
		UserID string `json:"user_id"`
	}

	body := req.Body

	defer func() {
		if err := body.Close(); err != nil {
			log.Printf("%s: can not close body: %s\n", nameMethod, err.Error())
		}
	}()

	errD := json.NewDecoder(body).Decode(&delEvent)

	if errD != nil {
		errResp := core.ErrorResponse{
			Error: errD.Error(),
		}
		response.Resp(w, nil, errResp, http.StatusBadRequest)
		return
	}

	errDel := hh.sv.Delete(delEvent.ID, delEvent.UserID)

	if errDel != nil {
		errResp := core.ErrorResponse{
			Error: errDel.Error(),
		}
		response.Resp(w, nil, errResp, http.StatusBadRequest)
		return
	}

	log.Printf("%s: evetn by id %s is deleted", nameMethod, http.StatusOK)
	resp := core.SuccessResponse{Result: "ok"}

	response.Resp(w, resp, nil, http.StatusCreated)
}

func (hh *HTTPHandler) eventsForDay(w http.ResponseWriter, req *http.Request) {
	nameMethod := "events for day"
	if req.Method != http.MethodGet {
		response.Resp(w, nil, nil, http.StatusNotFound)
	}

	day := req.URL.Query().Get("day")
	userID := req.URL.Query().Get("user_id")

	if day == "" {
		errResp := core.ErrorResponse{
			Error: fmt.Errorf("%s: bad event: day %s not fond\n", nameMethod, day),
		}
		response.Resp(w,
			nil,
			errResp,
			http.StatusBadRequest,
		)

		return
	}

	if userID == "" {
		errResp := core.ErrorResponse{
			Error: fmt.Errorf("%s: bad event: user by id %s not fond\n", nameMethod, userID),
		}
		response.Resp(w,
			nil,
			errResp,
			http.StatusBadRequest,
		)

		return
	}

	dayTime, errP := time.Parse("2006-01-02", day)

	if errP != nil {
		errResp := core.ErrorResponse{
			Error: fmt.Errorf("%s: bad event: bad day data format: %s\n", nameMethod, errP),
		}
		response.Resp(w,
			nil,
			errResp,
			http.StatusBadRequest,
		)

		return
	}

	events, errEBD := hh.sv.EventByDay(dayTime, userID)

	if errEBD != nil {
		errResp := core.ErrorResponse{
			Error: fmt.Errorf("%s: bad event: bad day data format: %s\n", nameMethod, errP),
		}
		response.Resp(w,
			nil,
			errResp,
			http.StatusNotFound,
		)

		return
	}

	log.Printf("%s: return events by day %s", nameMethod, dayTime)

	resp := core.SuccessResponse{
		Result: events,
	}

	response.Resp(w, resp, nil, http.StatusCreated)
}

func (hh *HTTPHandler) eventsForWeek(w http.ResponseWriter, req *http.Request) {
	nameMethod := "events for week"
	if req.Method != http.MethodGet {
		response.Resp(w, nil, nil, http.StatusNotFound)
	}

	since := req.URL.Query().Get("since")
	userID := req.URL.Query().Get("user_id")

	if since == "" {
		errResp := core.ErrorResponse{
			Error: fmt.Errorf("%s: bad event: since %s not fond\n", nameMethod, since),
		}
		response.Resp(w,
			nil,
			errResp,
			http.StatusBadRequest,
		)

		return
	}

	if userID == "" {
		errResp := core.ErrorResponse{
			Error: fmt.Errorf("%s: bad event: user by id %s not fond\n", nameMethod, userID),
		}
		response.Resp(w,
			nil,
			errResp,
			http.StatusBadRequest,
		)

		return
	}

	sinceTime, errP := time.Parse("2006-01-02", since)

	if errP != nil {
		errResp := core.ErrorResponse{
			Error: fmt.Errorf("%s: bad event: bad since data format: %s\n", nameMethod, errP),
		}
		response.Resp(w,
			nil,
			errResp,
			http.StatusBadRequest,
		)

		return
	}

	events, errEBD := hh.sv.EventByWeek(sinceTime, userID)

	if errEBD != nil {
		errResp := core.ErrorResponse{
			Error: fmt.Errorf("%s: bad event: bad since data format: %s\n", nameMethod, errP),
		}
		response.Resp(w,
			nil,
			errResp,
			http.StatusNotFound,
		)

		return
	}

	log.Printf("%s: return events by week since %s", nameMethod, sinceTime)

	resp := core.SuccessResponse{
		Result: events,
	}

	response.Resp(w, resp, nil, http.StatusCreated)
}

func (hh *HTTPHandler) eventsForMonth(w http.ResponseWriter, req *http.Request) {
	nameMethod := "events for month"
	if req.Method != http.MethodGet {
		response.Resp(w, nil, nil, http.StatusNotFound)
	}

	since := req.URL.Query().Get("since")
	userID := req.URL.Query().Get("user_id")

	if since == "" {
		errResp := core.ErrorResponse{
			Error: fmt.Errorf("%s: bad event: since %s not fond\n", nameMethod, since),
		}
		response.Resp(w,
			nil,
			errResp,
			http.StatusBadRequest,
		)

		return
	}

	if userID == "" {
		errResp := core.ErrorResponse{
			Error: fmt.Errorf("%s: bad event: user by id %s not fond\n", nameMethod, userID),
		}
		response.Resp(w,
			nil,
			errResp,
			http.StatusBadRequest,
		)

		return
	}

	sinceTime, errP := time.Parse("2006-01-02", since)

	if errP != nil {
		errResp := core.ErrorResponse{
			Error: fmt.Errorf("%s: bad event: bad since data format: %s\n", nameMethod, errP),
		}
		response.Resp(w,
			nil,
			errResp,
			http.StatusBadRequest,
		)

		return
	}

	events, errEBD := hh.sv.EventByMonth(sinceTime, userID)

	if errEBD != nil {
		errResp := core.ErrorResponse{
			Error: fmt.Errorf("%s: bad event: bad since data format: %s\n", nameMethod, errP),
		}
		response.Resp(w,
			nil,
			errResp,
			http.StatusNotFound,
		)

		return
	}

	log.Printf("%s: return events by month since %s", nameMethod, sinceTime)

	resp := core.SuccessResponse{
		Result: events,
	}

	response.Resp(w, resp, nil, http.StatusCreated)
}
