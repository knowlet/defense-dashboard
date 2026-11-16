package quest

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMalformedRequestReturnsError(t *testing.T) {
	for _, helper := range []struct {
		name string
		call func(string, string, string, io.Reader) (*http.Response, error)
	}{{"form", request}, {"json", reqjson}} {
		for _, tc := range []struct{ name, method, url string }{
			{"invalid escape", http.MethodGet, "http://%zz/listing.php"},
			{"invalid port", http.MethodGet, "http://127.0.0.1:bad/listing.php"},
			{"invalid method", "bad method", "http://127.0.0.1/"},
		} {
			t.Run(helper.name+"/"+tc.name, func(t *testing.T) {
				defer func() {
					if p := recover(); p != nil {
						t.Errorf("malformed request panicked: %v", p)
					}
				}()
				resp, err := helper.call(tc.method, tc.url, "synthetic.invalid", nil)
				if err == nil || resp != nil {
					t.Fatalf("got response=%v error=%v, want nil response and error", resp, err)
				}
			})
		}
	}
}

func TestRequestHelpersPreserveValidRequests(t *testing.T) {
	for _, helper := range []struct {
		name, contentType string
		call              func(string, string, string, io.Reader) (*http.Response, error)
	}{
		{"form", "application/x-www-form-urlencoded", request},
		{"json", "application/json", reqjson},
	} {
		t.Run(helper.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				data, err := io.ReadAll(r.Body)
				if err != nil {
					t.Error(err)
				}
				if r.Host != "synthetic.invalid" || r.Method != http.MethodPost || r.Header.Get("Content-Type") != helper.contentType || string(data) != "payload" {
					t.Errorf("request changed: %s %s %s %s", r.Host, r.Method, r.Header.Get("Content-Type"), data)
				}
				w.Header().Set("Location", "/unexpected")
				w.WriteHeader(http.StatusFound)
				_, _ = w.Write([]byte("response"))
			}))
			defer server.Close()
			resp, err := helper.call(http.MethodPost, server.URL, "synthetic.invalid", strings.NewReader("payload"))
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()
			data, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}
			if resp.StatusCode != http.StatusFound || string(data) != "response" {
				t.Fatalf("response changed: %d %q", resp.StatusCode, data)
			}
		})
	}
}
