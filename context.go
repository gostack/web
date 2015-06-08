package web

import (
	"net/http"
	"os"

	"github.com/satori/go.uuid"
	"golang.org/x/net/context"
)

// ctxKey is a string identifier to be used when storing Info into a context,
// keeping all package data under the same key avoids collision from other packages.
const ctxKey = "github.com/gostack/web"

// caches the hostname for this system
var hostname string

// init the package
func init() {
	var err error

	hostname, err = os.Hostname()
	if err != nil {
		panic(err)
	}
}

// Info is a struct that stores information about the system and the request currently
// being processed.
type Info struct {
	Service       string
	Hostname      string
	TransactionID uuid.UUID
}

// NewInfo returns a new initialized Info instance.
func NewInfo(service string) Info {
	return Info{Service: service, Hostname: hostname, TransactionID: uuid.NewV4()}
}

// ContextHandler is a extension of a http.Handler that also includes a context.Context object.
type ContextHandler interface {
	ServeHTTP(c context.Context, w http.ResponseWriter, req *http.Request)
}

// ContextHandlerFunc implements ServeHTTP for a function
type ContextHandlerFunc func(c context.Context, w http.ResponseWriter, req *http.Request)

// ServeHTTP implements the http.Handler interface for a function, calling itself
func (ch ContextHandlerFunc) ServeHTTP(c context.Context, w http.ResponseWriter, req *http.Request) {
	ch(c, w, req)
}

// ContextHandlerAdapter wraps a ContextHandler, returning an http.Handler that initializes a
// context allowing ContexHandler to be mounted on any net/http compatible library.
func ContextHandlerAdapter(service string, ch ContextHandler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		ctx := ContextWithInfo(context.Background(), NewInfo(service))
		ch.ServeHTTP(ctx, w, req)
	}

	return http.HandlerFunc(fn)
}

// ContextWithInfo creates a new context containing which contains a instance of Info
func ContextWithInfo(ctx context.Context, info Info) context.Context {
	return context.WithValue(ctx, ctxKey, info)
}

// InfoFromContext retrieves an Info stored within a context.
func InfoFromContext(ctx context.Context) Info {
	return ctx.Value(ctxKey).(Info)
}
