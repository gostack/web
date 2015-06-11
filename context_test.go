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
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gostack/ctxinfo"
	"golang.org/x/net/context"
)

func TestContextLifecycle(t *testing.T) {
	ctxHandler := func(ctx context.Context, w http.ResponseWriter, req *http.Request) {
		env, tx := ctxinfo.EnvFromContext(ctx), ctxinfo.TxFromContext(ctx)

		if env.Application != "myapp" || tx.TransactionID.String() == "" {
			t.Error("context not initialized properly")
		}
	}

	handler := ContextHandlerAdapter(
		ctxinfo.EnvContext(context.Background(), "myapp"),
		ContextHandlerFunc(ctxHandler),
	)

	req, err := http.NewRequest("GET", "http://doximity.test", nil)
	if err != nil {
		log.Fatal(err)
	}

	handler.ServeHTTP(httptest.NewRecorder(), req)
}
