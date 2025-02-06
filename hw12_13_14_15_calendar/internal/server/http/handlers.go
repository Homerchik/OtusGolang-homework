package internalhttp

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/homerchik/OtusGolang-homework/hw12_13_14_15_calendar/internal/app"
	"github.com/homerchik/OtusGolang-homework/hw12_13_14_15_calendar/internal/models"
)

type CalendarHandler struct {
	logger  app.Logger
	storage models.Storage
}

var TimeFormat = "2006-01-02T15:04:05Z"

func (h *CalendarHandler) Hello(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("Hello world!"))
	time.Sleep(time.Second)
}

func (h *CalendarHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var (
		event   *models.Event
		message string
		err     error
	)
	defer r.Body.Close()
	if err = ReadAndUnmarshalJSON(r.Body, &event); err != nil {
		message = "failed to unmarshal event" + err.Error()
		err = models.ErrEventCantBeAdded
		h.logger.Error(message)
	} else {
		event.ID = uuid.New()
		if err = h.storage.AddEvent(*event); err != nil {
			message = "error adding event"
			h.logger.Error(message + err.Error())
		} else {
			w.Header().Add("Location", "/api/events/"+event.ID.String())
			w.WriteHeader(http.StatusCreated)
		}
	}
	if err != nil {
		code := MatchHTTPCode(err)
		if err := SendError(w, message, code); err != nil {
			h.logger.Error("error sending error" + err.Error())
		}
	}
}

func (h *CalendarHandler) GetEvent(w http.ResponseWriter, r *http.Request) {
	var (
		id  uuid.UUID
		err error
	)
	if id, err = GetIDFromPath(r.URL.Path); err != nil {
		w.WriteHeader(http.StatusNotFound)
	}
	if _, event, err := h.storage.GetEventByID(id); err != nil {
		code := MatchHTTPCode(err)
		if err := SendError(w, err.Error(), code); err != nil {
			h.logger.Error("error sending error" + err.Error())
		}
	} else {
		if err := WriteJSON(w, event); err != nil {
			h.logger.Error("error writing json" + err.Error())
		}
	}
}

func (h *CalendarHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	var (
		id      uuid.UUID
		err     error
		event   *models.Event
		message string
	)
	defer r.Body.Close()

	if id, err = GetIDFromPath(r.URL.Path); err != nil {
		w.WriteHeader(http.StatusNotFound)
	}

	if err = ReadAndUnmarshalJSON(r.Body, &event); err != nil {
		message = "failed to unmarshal event" + err.Error()
		err = models.ErrEventCantBeUpdated
		h.logger.Error(message)
	} else {
		if event.ID != uuid.Nil && event.ID != id {
			err = models.ErrEventCantBeUpdated
			message = "event ID in path and body are different"
			goto ErrorLabel
		}
		event.ID = id

		if err = h.storage.UpdateEvent(*event); err != nil {
			message = "error updating event"
			h.logger.Error(message + err.Error())
			goto ErrorLabel
		}

		w.Header().Add("Location", "/api/events/"+id.String())
		w.WriteHeader(http.StatusNoContent)
	}

ErrorLabel:
	if err != nil {
		code := MatchHTTPCode(err)
		if err := SendError(w, message, code); err != nil {
			h.logger.Error("error sending error" + err.Error())
		}
	}
}

func (h *CalendarHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	var (
		id  uuid.UUID
		err error
	)
	if id, err = GetIDFromPath(r.URL.Path); err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	if err = h.storage.DeleteEvent(id); err != nil {
		code := MatchHTTPCode(err)
		if err = SendError(w, err.Error(), code); err != nil {
			h.logger.Error("error sending error" + err.Error())
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *CalendarHandler) GetEventsForRange(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	fromTS := params.Get("from")
	toTS := params.Get("to")
	if fromTS == "" || toTS == "" {
		w.Write([]byte("Please specify 'from' and 'to' parameters"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	from, err := time.Parse(TimeFormat, fromTS)
	if err != nil {
		w.Write([]byte("Unsupported time format, please use iso8601"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	to, err := time.Parse(TimeFormat, toTS)
	if err != nil {
		w.Write([]byte("Unsupported time format, please use iso8601"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if events, err := h.storage.GetEvents(from.Unix(), to.Unix()); err != nil {
		code := MatchHTTPCode(err)
		if err := SendError(w, err.Error(), code); err != nil {
			h.logger.Error("error sending error" + err.Error())
		}
	} else {
		if err := WriteJSON(w, map[string]models.Schedule{"data": events}); err != nil {
			h.logger.Error("error writing json" + err.Error())
		}
	}
}
