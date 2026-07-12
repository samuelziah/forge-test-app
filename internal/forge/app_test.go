package forge

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestVersion(t *testing.T) {
	a := New()
	srv := httptest.NewServer(a.Handler())
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/version")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("status: got %d want 200", resp.StatusCode)
	}
	var v versionResponse
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		t.Fatal(err)
	}
	if v.Name != "forge-test-app" {
		t.Errorf("name: got %q want forge-test-app", v.Name)
	}
	if v.Version != "0.1.0" {
		t.Errorf("version: got %q want 0.1.0", v.Version)
	}
	if v.StartedAt == "" {
		t.Error("started_at should be set")
	}
}

func TestHealth(t *testing.T) {
	a := New()
	srv := httptest.NewServer(a.Handler())
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/health")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	var h healthResponse
	if err := json.NewDecoder(resp.Body).Decode(&h); err != nil {
		t.Fatal(err)
	}
	if h.Status != "ok" {
		t.Errorf("status: got %q want ok", h.Status)
	}
}

func TestEcho(t *testing.T) {
	a := New()
	srv := httptest.NewServer(a.Handler())
	defer srv.Close()

	resp, err := http.Post(srv.URL+"/echo", "application/json",
		strings.NewReader(`{"hello":"world","n":42}`))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("status: got %d want 200", resp.StatusCode)
	}
	var got map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}
	recv, ok := got["received"].(map[string]any)
	if !ok {
		t.Fatalf("received: not a map: %#v", got["received"])
	}
	if recv["hello"] != "world" {
		t.Errorf("hello: got %v want world", recv["hello"])
	}
}

func TestEchoBadJSON(t *testing.T) {
	a := New()
	srv := httptest.NewServer(a.Handler())
	defer srv.Close()

	resp, err := http.Post(srv.URL+"/echo", "application/json",
		strings.NewReader(`{not json`))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 400 {
		t.Errorf("status: got %d want 400", resp.StatusCode)
	}
}
