package pulse_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/chadsmith12/pdga-scoring/pkgs/pulse"
	"github.com/stretchr/testify/assert"
)

type Test struct {
	Name string
}

func TestReadJsonBody(t *testing.T) {
	var body = Test{Name: "Chad"}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("can't marshal body: %v", err)
	}

	reader := bytes.NewReader(bodyBytes)
	var actual Test
	err = pulse.Json(reader, &actual)

	assert.NoError(t, err)
	assert.Equal(t, body.Name, actual.Name)
}
