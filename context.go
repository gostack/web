/*
Copyright 2015 Rodrigo Rafael Monti Kochenburger

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package web

import (
	"net/http"

	"github.com/gostack/ctxinfo"
	"github.com/zenazn/goji/web"
	"golang.org/x/net/context"
)

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
func ContextHandlerAdapter(ctx context.Context, ch ContextHandler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		ch.ServeHTTP(ctxinfo.TxContext(ctx), w, req)
	}

	return http.HandlerFunc(fn)
}

// ContextHandlerAdapter wraps a ContextHandler, returning an http.Handler that initializes a
// context allowing ContexHandler to be mounted on any net/http compatible library.
func GojiContextHandlerAdapter(ctx context.Context, ch ContextHandler) web.Handler {
	fn := func(c web.C, w http.ResponseWriter, req *http.Request) {
		ctx := ctxinfo.TxContext(ctx)
		ctx = context.WithValue(ctx, "github.com/gostack/web:goji", c.URLParams)
		ch.ServeHTTP(ctx, w, req)
	}

	return web.HandlerFunc(fn)
}

func GojiParam(ctx context.Context, key string) string {
	m := ctx.Value("github.com/gostack/web:goji").(map[string]string)
	return m[key]
}
