package api

import (
	"books/dataprovider"
	"books/model"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/mock"
)

type DataProviderMock struct {
	mock.Mock
}

func (m DataProviderMock) Fetch(params *dataprovider.BooksParams) ([]model.BookInformation, error) {
	args := m.Called()
	return args.Get(0).([]model.BookInformation), args.Error(1)
}

func TestBooksHandler(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		want           string
		wantHttpStatus int
		mockProvider   func() DataProviderMock
	}{
		{
			name:           "Test Invalid Title",
			url:            "/books?title=flowersflowersflowersflowersflowersflowersflowersflowers&limit=20",
			want:           "invalid title",
			wantHttpStatus: http.StatusBadRequest,
			mockProvider: func() DataProviderMock {
				provider := DataProviderMock{}
				provider.On("Fetch").Return([]model.BookInformation{}, nil)
				return provider
			},
		},
		{
			name:           "Test Invalid Limit",
			url:            "/books?title=flowers&limit=11",
			want:           "invalid limit",
			wantHttpStatus: http.StatusBadRequest,
			mockProvider: func() DataProviderMock {
				provider := DataProviderMock{}
				provider.On("Fetch").Return([]model.BookInformation{}, nil)
				return provider
			},
		},
		{
			name:           "Test Invalid Limit & Title",
			url:            "/books?title=flowersflowersflowersflowersflowersflowersflowersflowers&limit=200",
			want:           "invalid limit, invalid title",
			wantHttpStatus: http.StatusBadRequest,
			mockProvider: func() DataProviderMock {
				provider := DataProviderMock{}
				provider.On("Fetch").Return([]model.BookInformation{}, nil)
				return provider
			},
		},
		{
			name:           "Test Empty Parameter",
			url:            "/books?limit=200",
			want:           "title is empty",
			wantHttpStatus: http.StatusBadRequest,
			mockProvider: func() DataProviderMock {
				provider := DataProviderMock{}
				provider.On("Fetch").Return([]model.BookInformation{}, nil)
				return provider
			},
		},
		{
			name:           "Test DataProvider Fetch Error",
			url:            "/books?title=flowers&limit=100",
			want:           "error msg",
			wantHttpStatus: http.StatusInternalServerError,
			mockProvider: func() DataProviderMock {
				provider := DataProviderMock{}
				provider.On("Fetch").Return([]model.BookInformation{}, errors.New("error msg"))
				return provider
			},
		},
		{
			name:           "Test",
			url:            "/books?title=flowers&limit=2",
			want:           `{"items":[{"id":"CVYKAQAAMAAJ","volumeInfo":{"title":"","language":""}},{"id":"K0lQAAAAMAAJ","volumeInfo":{"title":"","language":""}}],"count":2}`,
			wantHttpStatus: http.StatusOK,
			mockProvider: func() DataProviderMock {
				provider := DataProviderMock{}
				provider.On("Fetch").
					Return([]model.BookInformation{
						{Id: "CVYKAQAAMAAJ"},
						{Id: "K0lQAAAAMAAJ"},
					}, nil)
				return provider
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := executeHttp(t, NewBooksHandler(tt.mockProvider()), http.MethodGet, tt.url)
			if result.Result().StatusCode != tt.wantHttpStatus {
				t.Errorf("pageRankHandler.ServeHTTP() httpStatusCode = %v, wantHttpStatus = %v", result.Result().StatusCode, tt.wantHttpStatus)
				return
			}
			got := strings.TrimSpace(result.Body.String())
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("pageRankHandler.ServeHTTP() \n(+got) = %v \n(-want) = %v \n(+-diff) = %v", got, tt.want, diff)
			}
		})
	}
}
