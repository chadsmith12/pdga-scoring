package pulse

import (
	"encoding/json"
	"io"
)

func Json(body io.Reader, data any) error {
    return json.NewDecoder(body).Decode(&data)
}
