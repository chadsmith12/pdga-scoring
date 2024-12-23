package pulse

import (
	"encoding/json"
	"net/http"
)

type PuleHttpWriter interface {
    Write(w http.ResponseWriter, req *http.Request)
}

type HttpStatusCodeWriter struct {
    StatusCode int
}

func (scw HttpStatusCodeWriter) Write(w http.ResponseWriter, req *http.Request) {
    w.WriteHeader(scw.StatusCode)
}


type JsonResultWriter struct {
    HttpStatusCodeWriter
    Data interface{}
}

func (jw JsonResultWriter) Write(w http.ResponseWriter, req *http.Request) {
    w.Header().Set("content-type", "application/json")
    w.WriteHeader(jw.StatusCode)
    json.NewEncoder(w).Encode(jw.Data)
}


func JsonResult(data interface{}) JsonResultWriter {
    return JsonResultWriter{
        HttpStatusCodeWriter: HttpStatusCodeWriter{ StatusCode: 200 },
        Data: data,
    }
}

func ErrorJson(statusCode int, data interface{}) JsonResultWriter {
    return JsonResultWriter{
        HttpStatusCodeWriter: HttpStatusCodeWriter{ StatusCode: statusCode },
        Data: data,
    }
}

func InternalErrorJson(data interface{}) JsonResultWriter {
    return ErrorJson(500, data)
}

func OkResult() HttpStatusCodeWriter {
    return HttpStatusCodeWriter{ StatusCode: 200 }
}
