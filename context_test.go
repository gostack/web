package web

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/net/context"
)

func TestContextLifecycle(t *testing.T) {
	ctxHandler := func(ctx context.Context, w http.ResponseWriter, req *http.Request) {
		j, err := json.Marshal(InfoFromContext(ctx))
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

	v := Info{}
	if err := json.Unmarshal(w.Body.Bytes(), &v); err != nil {
		t.Fatal(err)
	}

	if v.Hostname == "" {
		t.Error("hostname is empty")
	}

	if v.TransactionID.String() == "" {
		t.Error("transaction id is empty")
	}

	if v.Service != "doximity.test" {
		t.Error("service is not doximity.test")
	}
}
