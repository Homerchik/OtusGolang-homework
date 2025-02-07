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
	var event *models.Event
	defer r.Body.Close()
	if err := ReadAndUnmarshalJSON(r.Body, &event); err != nil {
		SendError(w, "failed to unmarshal event"+err.Error(), models.ErrEventCantBeAdded, h.logger)
		return
	}
	event.ID = uuid.New()
	if err := h.storage.AddEvent(*event); err != nil {
		SendError(w, "error adding event", err, h.logger)
		return
	}
	w.Header().Add("Location", "/api/events/"+event.ID.String())
	w.WriteHeader(http.StatusCreated)
}

func (h *CalendarHandler) GetEvent(w http.ResponseWriter, r *http.Request) {
	var (
		id    uuid.UUID
		err   error
		event models.Event
	)

	if id, err = uuid.Parse(r.PathValue("id")); err != nil {
		w.WriteHeader(http.StatusNotFound)
		SendError(w, "wrong uuid has gotten", err, h.logger)
		return
	}
	if _, event, err = h.storage.GetEventByID(id); err != nil {
		SendError(w, "error sending error", err, h.logger)
		return
	}
	if err := WriteJSON(w, event); err != nil {
		h.logger.Error("error writing json" + err.Error())
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

	if id, err = uuid.Parse(r.PathValue("id")); err != nil {
		w.WriteHeader(http.StatusNotFound)
		SendError(w, "wrong uuid has gotten", err, h.logger)
		return
	}

	if err = ReadAndUnmarshalJSON(r.Body, &event); err != nil {
		h.logger.Error(message)
		SendError(w, "failed to unmarshal event"+err.Error(), models.ErrEventCantBeUpdated, h.logger)
		return
	}
	if event.ID != uuid.Nil && event.ID != id {
		SendError(w, "event ID in path and body are different", models.ErrEventCantBeUpdated, h.logger)
		return
	}

	event.ID = id

	if err = h.storage.UpdateEvent(*event); err != nil {
		SendError(w, "error updating event", err, h.logger)
		return
	}

	w.Header().Add("Location", "/api/events/"+id.String())
	w.WriteHeader(http.StatusNoContent)
}

func (h *CalendarHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	var (
		id  uuid.UUID
		err error
	)
	if id, err = uuid.Parse(r.PathValue("id")); err != nil {
		w.WriteHeader(http.StatusNotFound)
		SendError(w, "wrong uuid has gotten", err, h.logger)
		return
	}
	if err = h.storage.DeleteEvent(id); err != nil {
		SendError(w, "", err, h.logger)
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
		SendError(w, "", err, h.logger)
	} else {
		if err := WriteJSON(w, map[string]models.Schedule{"data": events}); err != nil {
			h.logger.Error("error writing json" + err.Error())
		}
	}
}
