package test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adnanahmady/go-rest-api-blog/internal"
	"github.com/go-chi/chi/v5"
)

type testServer struct {
	App    *internal.App
	routes *chi.Mux
}

func Setup() (*testServer, error) {
	app, err := internal.WireUpApp()
	if err != nil {
		return nil, err
	}
	app.Database.Migrate()
	routes := app.Server.GetEngine()

	return &testServer{
		App:    app,
		routes: routes,
	}, nil
}

func (s *testServer) Close() error {
	return s.App.Database.Close()
}

func (s *testServer) jsonRequest(
	t testing.TB,
	method, path string,
	body any,
	headers map[string]string,
) (*httptest.ResponseRecorder, any) {
	t.Helper()
	requestBody, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("failed to marshal body: %v", err)
		return nil, err
	}
	requestBodyReader := bytes.NewReader(requestBody)
	recorder, buffer := s.rawRequest(t, method, path, requestBodyReader, headers)

	var responseBody any
	_ = json.Unmarshal(buffer, &responseBody)

	return recorder, responseBody
}

func (s *testServer) rawRequest(
	t testing.TB,
	method, path string,
	body io.Reader,
	headers map[string]string,
) (*httptest.ResponseRecorder, []byte) {
	t.Helper()
	req, err := http.NewRequest(method, path, body)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
		return nil, nil
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	recorder := httptest.NewRecorder()
	s.routes.ServeHTTP(recorder, req)

	buffer, _ := io.ReadAll(recorder.Body)

	return recorder, buffer
}

func (s *testServer) Post(
	t testing.TB, path string, body any,
	headers map[string]string,
) (*httptest.ResponseRecorder, any) {
	t.Helper()
	return s.jsonRequest(t, http.MethodPost, path, body, headers)
}

func (s *testServer) Put(
	t testing.TB, path string, body any,
	headers map[string]string,
) (*httptest.ResponseRecorder, any) {
	t.Helper()
	return s.jsonRequest(t, http.MethodPut, path, body, headers)
}

func (s *testServer) Delete(
	t testing.TB, path string,
	headers map[string]string,
) (*httptest.ResponseRecorder, any) {
	t.Helper()
	return s.jsonRequest(t, http.MethodDelete, path, nil, headers)
}

func (s *testServer) Get(
	t testing.TB, path string,
	headers map[string]string,
) (*httptest.ResponseRecorder, any) {
	t.Helper()
	return s.jsonRequest(t, http.MethodGet, path, nil, headers)
}
