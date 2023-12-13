package precode

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	totalCount := 4

	req, _ := http.NewRequest("GET", "/cafe", nil)
	q := req.URL.Query()
	q.Add("count", "5")
	q.Add("city", "moscow")
	req.URL.RawQuery = q.Encode()

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	require.NotEmpty(t, responseRecorder.Body)
	body, _ := io.ReadAll(responseRecorder.Body)
	answer := string(body)
	cafes := strings.Split(answer, ",")
	assert.Len(t, cafes, totalCount)
}

func TestMainHandlerCorrectRequest(t *testing.T) {
	req, _ := http.NewRequest("GET", "/cafe", nil)
	q := req.URL.Query()
	q.Add("count", "1")
	q.Add("city", "moscow")
	req.URL.RawQuery = q.Encode()

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	require.NotEmpty(t, responseRecorder.Body)
	require.Equal(t, responseRecorder.Result().Status, "200 OK")
	body, _ := io.ReadAll(responseRecorder.Body)
	answer := string(body)
	cafes := strings.Split(answer, ",")
	assert.Len(t, cafes, 1)
}

func TestMainHandlerWrongCity(t *testing.T) {
	req, _ := http.NewRequest("GET", "/cafe", nil)
	q := req.URL.Query()
	q.Add("count", "1")
	q.Add("city", "spb")
	req.URL.RawQuery = q.Encode()

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	require.NotEmpty(t, responseRecorder.Body)
	require.Equal(t, responseRecorder.Result().Status, "400 Bad Request")

	body, _ := io.ReadAll(responseRecorder.Body)
	answer := string(body)
	assert.Equal(t, answer, "wrong city value")
}
