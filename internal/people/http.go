package people

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) http.Handler {
	return &Handler{service: service}
}

func (h *Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path == "/health" {
		if request.Method != http.MethodGet {
			writeError(writer, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
			return
		}
		writeJSON(writer, http.StatusOK, map[string]string{"status": "ok"})
		return
	}
	if request.URL.Path == "/admin/people" {
		switch request.Method {
		case http.MethodGet:
			people, err := h.service.List(request.Context())
			if err != nil {
				writeServiceError(writer, err)
				return
			}
			writeJSON(writer, http.StatusOK, people)
		case http.MethodPost:
			var input CreateInput
			if !decodeJSON(writer, request, &input) {
				return
			}
			person, err := h.service.Create(request.Context(), input)
			if err != nil {
				writeServiceError(writer, err)
				return
			}
			writeJSON(writer, http.StatusCreated, person)
		default:
			writeError(writer, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
		}
		return
	}
	const prefix = "/admin/people/"
	if strings.HasPrefix(request.URL.Path, prefix) {
		id := strings.TrimPrefix(request.URL.Path, prefix)
		if id == "" || strings.Contains(id, "/") {
			writeError(writer, http.StatusNotFound, "person_not_found", "person not found")
			return
		}
		switch request.Method {
		case http.MethodPatch:
			var input StatusInput
			if !decodeJSON(writer, request, &input) {
				return
			}
			person, err := h.service.UpdateStatus(request.Context(), id, input.Status)
			if err != nil {
				writeServiceError(writer, err)
				return
			}
			writeJSON(writer, http.StatusOK, person)
		case http.MethodDelete:
			if err := h.service.Delete(request.Context(), id); err != nil {
				writeServiceError(writer, err)
				return
			}
			writer.WriteHeader(http.StatusNoContent)
		default:
			writeError(writer, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
		}
		return
	}
	writeError(writer, http.StatusNotFound, "not_found", "route not found")
}

func decodeJSON(writer http.ResponseWriter, request *http.Request, destination any) bool {
	decoder := json.NewDecoder(request.Body)
	if err := decoder.Decode(destination); err != nil {
		writeError(writer, http.StatusBadRequest, "invalid_request", "request body is invalid")
		return false
	}
	return true
}

func writeServiceError(writer http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrContactConflict):
		writeError(writer, http.StatusConflict, "contact_conflict", "phone is already assigned")
	case errors.Is(err, ErrNotFound):
		writeError(writer, http.StatusNotFound, "person_not_found", "person not found")
	case errors.Is(err, ErrInvalidRole):
		writeError(writer, http.StatusBadRequest, "invalid_role", "role is invalid")
	case errors.Is(err, ErrInvalidStatus):
		writeError(writer, http.StatusBadRequest, "invalid_status", "status is invalid")
	case errors.Is(err, ErrInvalidRequest):
		writeError(writer, http.StatusBadRequest, "invalid_request", "person data is invalid")
	default:
		writeError(writer, http.StatusInternalServerError, "internal_error", "internal server error")
	}
}

func writeError(writer http.ResponseWriter, status int, code string, message string) {
	writeJSON(writer, status, map[string]string{"code": code, "message": message})
}

func writeJSON(writer http.ResponseWriter, status int, value any) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	_ = json.NewEncoder(writer).Encode(value)
}
