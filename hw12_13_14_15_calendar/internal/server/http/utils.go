package internalhttp

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/homerchik/OtusGolang-homework/hw12_13_14_15_calendar/internal/models"
)

func ReadAndUnmarshalJSON(payload io.ReadCloser, v interface{}) error {
	data, err := io.ReadAll(payload)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, v); err != nil {
		return err
	}
	return nil
}

func WriteJSON(w io.Writer, v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	if _, err := w.Write(data); err != nil {
		return err
	}
	return nil
}

func SendError(w http.ResponseWriter, message string, err error, log Logger) {
	log.Error(message + ":" + err.Error())
	code := MatchHTTPCode(err)
	errorMsg := models.RestError{Message: message}
	w.WriteHeader(code)
	if err := WriteJSON(w, errorMsg); err != nil {
		log.Error("error sending error" + err.Error())
	}
}

func MatchHTTPCode(err error) int {
	if errors.Is(err, models.ErrInternalDBError) {
		return http.StatusInternalServerError
	}
	if errors.Is(err, models.ErrEventCantBeAdded) {
		return http.StatusBadRequest
	}
	if errors.Is(err, models.ErrNoEventFound) {
		return http.StatusNotFound
	}
	if errors.Is(err, models.ErrEventCantBeUpdated) {
		return http.StatusBadRequest
	}
	if errors.Is(err, models.ErrEventCantBeDeleted) {
		return http.StatusBadRequest
	}
	if err != nil {
		return http.StatusInternalServerError
	}
	return http.StatusOK
}
