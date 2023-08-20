package app

import (
	"books/constants"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestApp(t *testing.T) {
	app := &App{}

	t.Run("Test App Init", func(t *testing.T) {
		app.Init()
		if diff := cmp.Diff(constants.SERVER_ADDR, app.server.Addr); diff != "" {
			t.Errorf("app.server.Addr \n(+got) = %v \n(-want) = %v \n(+-diff) = %v", app.server.Addr, constants.SERVER_ADDR, diff)
		}
		if diff := cmp.Diff(constants.DEFAULT_TIMEOUT, app.server.ReadTimeout); diff != "" {
			t.Errorf("app.server.ReadTimeout \n(+got) = %v \n(-want) = %v \n(+-diff) = %v", app.server.ReadTimeout, constants.DEFAULT_TIMEOUT, diff)
		}
	})

	t.Run("Test App Start", func(t *testing.T) {
		app.Start()
		resp, err := http.Get("http://localhost:8080/health")
		if err != nil {
			t.Errorf("app.server.healthHandler.ServerHTTP() error = %v, wantErr = %v", err, false)
			return
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("app.server.healthHandler.ServeHTTP() httpStatus = %v, wantHttpStatus = %v", resp.StatusCode, http.StatusOK)
			return
		}
		got, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("app.server.healthHandler.ServerHTTP() error = %v, wantErr = %v", err, false)
			return
		}
		if diff := cmp.Diff("healthy", string(got)); diff != "" {
			t.Errorf("app.server.healthHandler.ServeHTTP() \n(+got) = %v \n(-want) = %v \n(+-diff) = %v", got, "healthy", diff)
		}
	})
}
