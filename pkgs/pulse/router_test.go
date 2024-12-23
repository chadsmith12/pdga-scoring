package pulse_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chadsmith12/pdga-scoring/pkgs/pulse"
	"github.com/stretchr/testify/assert"
)

func TestBasicRoute(t *testing.T) {
    app := pulse.Pulse()
    router := pulse.NewRouter(app)

    router.Get("/", func(req *http.Request) pulse.PuleHttpWriter {
        return pulse.OkResult()        
    })

    req := httptest.NewRequest("GET", "/", nil)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)
}

func TestJsonWriting(t *testing.T) {
    app := pulse.Pulse()
    router := pulse.NewRouter(app)
    var result = struct { Ok bool } { Ok: true }
    encoded, _ := json.Marshal(result)
    
    router.Get("/", func(req *http.Request) pulse.PuleHttpWriter {
        return pulse.JsonResult(result)
    })

    req := httptest.NewRequest("GET", "/", nil)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    encodedBody := string(encoded)
    assert.Equal(t, http.StatusOK, w.Code)
    assert.Equal(t, "application/json", w.Header().Get("content-type"))
    assert.Equal(t, "{\"Ok\":true}", encodedBody)
}
