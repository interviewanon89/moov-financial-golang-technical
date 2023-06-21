package customers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/moov-io/base/log"
)

func jsonResponse(w http.ResponseWriter, value interface{}) {
	jsonResponseStatus(w, http.StatusOK, value)
}

func jsonResponseStatus(w http.ResponseWriter, status int, value interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	e := json.NewEncoder(w)
	e.Encode(value)
}

func errorResponse(w http.ResponseWriter, err error, logger log.Logger) {
	validationErr := &validation.Errors{}

	switch true {
	case errors.Is(err, &json.InvalidUnmarshalError{}):
		w.WriteHeader(http.StatusBadRequest)
	case errors.As(err, validationErr):
		jsonResponseStatus(w, http.StatusUnprocessableEntity, validationErr)
	case errors.Is(err, sql.ErrNoRows):
		w.WriteHeader(http.StatusNotFound)
	default:
		logger.LogErrorf("unexpected: %w", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
