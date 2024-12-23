package pulse_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chadsmith12/pdga-scoring/pkgs/pulse"
	"github.com/stretchr/testify/assert"
)

func TestGroup(t *testing.T) {
    app := pulse.Pulse()
    router := pulse.NewRouter(app)
    group := router.Group("/api")

    group.Get("/hello", func(req *http.Request) pulse.PuleHttpWriter {
        return pulse.OkResult()
    })

    req := httptest.NewRequest("GET", "/api/hello", nil)
    w := httptest.NewRecorder()

    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)
}
