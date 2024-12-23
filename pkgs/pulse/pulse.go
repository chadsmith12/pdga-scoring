package pulse

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const (
    GET = http.MethodGet
    HEAD = http.MethodHead
    POST = http.MethodPost
    PUT = http.MethodPut
    DELETE = http.MethodDelete
    CONNECT = http.MethodConnect
    OPTIONS = http.MethodOptions
    PATCH = http.MethodPatch
)

const DEFAULT_ADDR string = ":4500"

type EndpointHandler func(req *http.Request) PuleHttpWriter
type MiddlewareFunc func(EndpointHandler) EndpointHandler
type Option func(*PulseApp)

type PulseApp struct {
    server *http.Server
    router *PulseRouter
    logger *slog.Logger
    addr string
    shutdownTimeout time.Duration
}

func Pulse(options ...Option) *PulseApp {
    pulseApp := &PulseApp{ 
        addr: DEFAULT_ADDR,
        shutdownTimeout: time.Second * 30,
    }

    logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
    router := NewRouter(pulseApp)
    server := &http.Server { Handler: router,}

    pulseApp.logger = logger
    pulseApp.server = server
    pulseApp.router = router

    for _, option := range options {
        option(pulseApp)
    }
    pulseApp.server.Addr = pulseApp.addr

    return pulseApp
}

func WithAddr(addr string) Option {
    return func(pa *PulseApp) {
        pa.addr = addr
    }
}

func WithLogger(logger *slog.Logger) Option {
    return func(pa *PulseApp) {
        pa.logger = logger
    }
}

func WithShutdownTimeout(duration time.Duration) Option {
    return func(pa *PulseApp) {
        pa.shutdownTimeout = duration
    }
}

func (p *PulseApp) Start() error {
    ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
    defer cancel()

    done := make(chan struct{})
    go func() {
        err := p.server.ListenAndServe()
        if err != nil && !errors.Is(err, http.ErrServerClosed) {
            p.logger.Error("Failed to start server and listen", slog.Any("error", err))
        }
        close(done)
    }()

    p.logger.LogAttrs(context.TODO(), slog.LevelInfo, "server started", slog.String("addr", p.addr))

    select {
    case <-done:
        break;
    case <-ctx.Done():
        ctx, cancel := context.WithTimeout(context.Background(), p.shutdownTimeout)
        p.logger.LogAttrs(ctx, slog.LevelInfo, "Shutting down server...")
        p.server.Shutdown(ctx)
        cancel()
    }

    return nil
}

func (p *PulseApp) UseStaticFiles() {
    p.router.mux.Handle("/", http.FileServer(http.Dir("wwwroot")))
}

func (p *PulseApp) Get(pattern string, endpoint EndpointHandler) {
    p.router.Get(pattern, endpoint)
}

func (p *PulseApp) Post(pattern string, endpoint EndpointHandler) {
    p.router.Post(pattern, endpoint)
}

func (p *PulseApp) Group(prefix string) *Group {
    return p.router.Group(prefix)
}

func (p *PulseApp) Logger() *slog.Logger {
    return p.logger
}
