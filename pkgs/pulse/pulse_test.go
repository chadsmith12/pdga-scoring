package pulse_test

import (
	"net"
	"testing"
	"time"

	"github.com/chadsmith12/pdga-scoring/pkgs/pulse"
	"github.com/stretchr/testify/assert"
)

func TestPulseStart(t *testing.T) {
    app := pulse.Pulse()
  
    go func() {
        err := app.Start()
        if err != nil {
            t.Errorf("Server failed to start: %v", err)
        }
    }()
    
    // wait a little bit for server
    time.Sleep(time.Millisecond * 400)

    conn, err := net.Dial("tcp", ":4500")
    assert.NoErrorf(t, err, "error connecting to :4500: %v", err)
    conn.Close()
}

func TestPulseStartCustom(t *testing.T) {
    app := pulse.Pulse(pulse.WithAddr(":8080"))

    go func () {
	err := app.Start()
	if err != nil {
	    t.Errorf("Server failed to start: %v", err)
	}
    }()

    time.Sleep(time.Millisecond * 400)
    
    conn, err := net.Dial("tcp", ":8080")
    assert.NoErrorf(t, err, "error connecting to :8080: %v", err)
    conn.Close()
}
