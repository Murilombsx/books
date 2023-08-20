package api

import (
	"net/http"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestHealthHandler(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		want           string
		wantHttpStatus int
	}{
		{
			name:           "Test Successful Request",
			url:            "health",
			want:           "healthy",
			wantHttpStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := executeHttp(t, NewHeatlhHandler(), http.MethodGet, tt.url)
			if result.Result().StatusCode != tt.wantHttpStatus {
				t.Errorf("healthHandler.ServeHTTP() httpStatus = %v, wantHttpStatus = %v", result.Result().StatusCode, tt.wantHttpStatus)
				return
			}
			got := strings.TrimSpace(result.Body.String())
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("healthHandler.ServeHTTP() \n(+got) = %v \n(-want) = %v \n(+-diff) = %v", got, tt.want, diff)
			}
		})
	}
}
