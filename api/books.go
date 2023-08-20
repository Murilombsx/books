package api

import (
	"books/dataprovider"
	"books/model"
	"encoding/json"
	"github.com/gorilla/schema"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

type IDataProvider interface {
	Fetch(params *dataprovider.BooksParams) ([]model.BookInformation, error)
}

var _ IDataProvider = dataprovider.DataProvider{}

type booksHandler struct {
	provider IDataProvider
}

func NewBooksHandler(provider IDataProvider) http.Handler {
	return booksHandler{
		provider: provider,
	}
}

type booksResponse struct {
	Items []model.BookInformation `json:"items"`
	Count uint8                   `json:"count"`
}

func (h booksHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	params := new(dataprovider.BooksParams)
	// Validate if the parameters comply with the schema defined
	if err := schema.NewDecoder().Decode(params, r.URL.Query()); err != nil {
		log.WithError(err).Error("parameters decoding failed")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else {
		// Validate if the values passed for each parameter are within the specifications
		if validation, ok := validateParams(*params); !ok {
			log.Errorf("parameters validation failed, reason: %s", validation)
			http.Error(w, validation, http.StatusBadRequest)
			return
		}
	}

	data, err := h.provider.Fetch(params)
	if err != nil {
		log.WithError(err).Error("data provider fetch failed")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Encode response as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(generateResponse(*params, data)); err != nil {
		log.WithError(err).Error("failed to encode response")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// validateParams will check:
//   - Title is greater than 1, less than 50
//   - Limit is greater than 1, less than 200
func validateParams(params dataprovider.BooksParams) (string, bool) {
	var validationErrors []string

	// Validate limit and store result
	if params.Limit <= 1 || params.Limit >= 200 {
		validationErrors = append(validationErrors, "invalid limit")
	}

	if len(params.Title) <= 1 || len(params.Title) >= 50 {
		validationErrors = append(validationErrors, "invalid title")
	}

	if len(validationErrors) > 0 {
		return strings.Join(validationErrors, ", "), false
	}
	return "", true
}

// generateResponse will execute the following:
//   - Limit results to be displayed based on the limit value defined by the caller
func generateResponse(params dataprovider.BooksParams, data []model.BookInformation) booksResponse {
	// Limit results
	if uint8(len(data)) > params.Limit {
		data = data[:params.Limit]
	}

	return booksResponse{
		Items: data,
		Count: uint8(len(data)),
	}
}
