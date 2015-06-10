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
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gostack/ctxinfo"

	"golang.org/x/net/context"
)

func TestContextLifecycle(t *testing.T) {
	ctxHandler := func(ctx context.Context, w http.ResponseWriter, req *http.Request) {
		j, err := json.Marshal(ctxinfo.FromContext(ctx))
		if err != nil {
			panic(err)
		}
		w.Write(j)
	}

	handler := ContextHandlerAdapter("doximity.test", ContextHandlerFunc(ctxHandler))

	req, err := http.NewRequest("GET", "http://doximity.test", nil)
	if err != nil {
		log.Fatal(err)
	}

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	v := ctxinfo.Info{}
	if err := json.Unmarshal(w.Body.Bytes(), &v); err != nil {
		t.Fatal(err)
	}

	if v.Hostname == "" {
		t.Error("hostname is empty")
	}

	if v.TransactionID.String() == "" {
		t.Error("transaction id is empty")
	}

	if v.Application != "doximity.test" {
		t.Error("application is not doximity.test")
	}

	if v.Service != "webapp" {
		t.Error("service is not webapp")
	}
}
