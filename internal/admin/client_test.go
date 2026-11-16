package admin

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestARequestsOwnTimeoutReplacesTheDefault(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
	}))
	defer server.Close()
	client := &Client{Origin: server.URL, HTTP: server.Client(), Token: func(context.Context) (string, error) { return "token", nil }}

	started := time.Now()
	_, err := client.Send(context.Background(), Request{Method: http.MethodPost, Path: Path("users", "u-ada", "restore-serving"), Timeout: 20 * time.Millisecond})
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected the request to time out, got %v", err)
	}
	if waited := time.Since(started); waited > requestTimeout/2 {
		t.Fatalf("waited %s, as long as the default allows", waited)
	}
}
