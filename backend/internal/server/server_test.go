package server

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/wayne-stays/backend/internal/config"
	"github.com/wayne-stays/backend/internal/database"
	"github.com/wayne-stays/backend/internal/storage"
	"github.com/wayne-stays/backend/internal/testutil"
)

const staticURL = "http://test.local/static"

func TestMain(m *testing.M) {
	flag.Parse()
	if testing.Short() {
		os.Exit(m.Run())
	}

	stop, err := testutil.StartPostgres()
	if err != nil {
		log.Fatal(err)
	}
	code := m.Run()
	stop()
	os.Exit(code)
}

// newTestServer returns a router backed by a freshly migrated and seeded database
func newTestServer(t *testing.T) http.Handler {
	t.Helper()
	db := testutil.NewDatabase(t)
	if err := database.Migrate(db); err != nil {
		t.Fatal(err)
	}
	if err := database.Seed(db); err != nil {
		t.Fatal(err)
	}

	cfg := &config.Config{
		Server: config.ServerConfig{
			Mode:           "test",
			AllowedOrigins: []string{"http://localhost:5173"},
		},
		Storage: config.StorageConfig{
			Type:         "local",
			LocalPath:    t.TempDir(),
			LocalBaseURL: staticURL,
		},
	}
	store, err := storage.New(&cfg.Storage)
	if err != nil {
		t.Fatal(err)
	}

	return NewRouter(cfg, db, store)
}

// get performs a GET request and decodes the JSON body into out (if non-nil)
func get(t *testing.T, h http.Handler, path string, out interface{}) int {
	t.Helper()
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
	if out != nil {
		if err := json.Unmarshal(rec.Body.Bytes(), out); err != nil {
			t.Fatalf("GET %s: invalid JSON %q: %v", path, rec.Body.String(), err)
		}
	}
	return rec.Code
}

type accommodationJSON struct {
	ID         string  `json:"id"`
	Slug       string  `json:"slug"`
	Title      string  `json:"title"`
	PriceStart float64 `json:"priceStart"`
	Category   string  `json:"category"`
	Images     struct {
		Hero    string   `json:"hero"`
		Gallery []string `json:"gallery"`
	} `json:"images"`
}

func TestHealth(t *testing.T) {
	h := newTestServer(t)
	if code := get(t, h, "/health", nil); code != http.StatusOK {
		t.Errorf("status = %d, want 200", code)
	}
}

func TestAccommodationEndpoints(t *testing.T) {
	h := newTestServer(t)

	var list struct {
		Count int                 `json:"count"`
		Data  []accommodationJSON `json:"data"`
	}
	if code := get(t, h, "/api/v1/accommodations", &list); code != http.StatusOK {
		t.Fatalf("list status = %d, want 200", code)
	}
	if list.Count != 5 || len(list.Data) != 5 {
		t.Fatalf("list count = %d (%d items), want 5", list.Count, len(list.Data))
	}

	first := list.Data[0]
	if first.Slug != "wayne-tower-sky-suites" {
		t.Errorf("first slug = %q, want seed order (wayne-tower-sky-suites)", first.Slug)
	}
	if first.PriceStart != 1850 {
		t.Errorf("priceStart = %v, want 1850 (dollars, not cents)", first.PriceStart)
	}
	if want := staticURL + "/accommodations/wayne_tower_sky_suites/1.jpg"; first.Images.Hero != want {
		t.Errorf("hero = %q, want %q", first.Images.Hero, want)
	}

	t.Run("filter by category", func(t *testing.T) {
		var res struct {
			Data []accommodationJSON `json:"data"`
		}
		get(t, h, "/api/v1/accommodations?category=historic", &res)
		if len(res.Data) != 2 {
			t.Errorf("historic count = %d, want 2", len(res.Data))
		}
		for _, a := range res.Data {
			if a.Category != "historic" {
				t.Errorf("got category %q in historic filter", a.Category)
			}
		}
	})

	tests := []struct {
		name string
		path string
		want int
	}{
		{"invalid category", "/api/v1/accommodations?category=moon-base", http.StatusBadRequest},
		{"by id", "/api/v1/accommodations/" + first.ID, http.StatusOK},
		{"unknown id", "/api/v1/accommodations/00000000-0000-0000-0000-000000000000", http.StatusNotFound},
		{"malformed id", "/api/v1/accommodations/not-a-uuid", http.StatusNotFound},
		{"by slug", "/api/v1/accommodations/slug/clocktower-penthouse", http.StatusOK},
		{"unknown slug", "/api/v1/accommodations/slug/batcave", http.StatusNotFound},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if code := get(t, h, tt.path, nil); code != tt.want {
				t.Errorf("GET %s = %d, want %d", tt.path, code, tt.want)
			}
		})
	}
}

func TestMerchandiseEndpoints(t *testing.T) {
	h := newTestServer(t)

	var raw struct {
		Data []map[string]interface{} `json:"data"`
	}
	if code := get(t, h, "/api/v1/merchandise", &raw); code != http.StatusOK {
		t.Fatalf("status = %d, want 200", code)
	}
	if len(raw.Data) != 10 {
		t.Fatalf("count = %d, want 10", len(raw.Data))
	}

	// Field names must be camelCase, matching the accommodations API
	item := raw.Data[0]
	for _, key := range []string{"imageUrl", "isNew"} {
		if _, ok := item[key]; !ok {
			t.Errorf("missing key %q in %v", key, item)
		}
	}
	if url, _ := item["imageUrl"].(string); !strings.HasPrefix(url, staticURL+"/merch/") {
		t.Errorf("imageUrl = %q, want prefix %s/merch/", url, staticURL)
	}

	// Category filter is case-insensitive; "all" means no filter
	for path, want := range map[string]int{
		"/api/v1/merchandise?category=Apparel": 2,
		"/api/v1/merchandise?category=home":    3,
		"/api/v1/merchandise?category=all":     10,
	} {
		var res struct {
			Data []map[string]interface{} `json:"data"`
		}
		get(t, h, path, &res)
		if len(res.Data) != want {
			t.Errorf("GET %s returned %d items, want %d", path, len(res.Data), want)
		}
	}
}

func TestCORSAllowsConfiguredOrigin(t *testing.T) {
	h := newTestServer(t)

	req := httptest.NewRequest(http.MethodOptions, "/api/v1/accommodations", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	req.Header.Set("Access-Control-Request-Method", http.MethodGet)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:5173" {
		t.Errorf("Access-Control-Allow-Origin = %q, want the configured origin", got)
	}
}
