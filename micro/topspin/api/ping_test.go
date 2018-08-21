package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestPing(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := NewAPIRouter()

	request, err := http.NewRequest(http.MethodGet, "/ping", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v\n", err)
	}

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)
	if recorder.Code != 200 {
		t.Fatalf("Expected to get status %d but instead got %d\n", 200, recorder.Code)
	}

	expected := "pong"
	received := recorder.Body.String()
	if expected != received {
		t.Fatalf("Expected to get '%s' in the response, but instead got '%s'\n", expected, received)
	}
}
