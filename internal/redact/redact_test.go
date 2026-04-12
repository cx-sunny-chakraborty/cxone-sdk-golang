package redact

import (
	"strings"
	"testing"
)

func TestStringRedactsBearerToken(t *testing.T) {
	in := "Authorization: Bearer eyJhbGciOiJIUzI1NiJ9.payload.sig and more"
	out := String(in)
	if strings.Contains(out, "eyJhbGciOiJIUzI1NiJ9") {
		t.Fatalf("bearer token leaked: %q", out)
	}
	if !strings.Contains(out, "Bearer REDACTED") {
		t.Fatalf("expected 'Bearer REDACTED', got %q", out)
	}
}

func TestStringRedactsFormSecrets(t *testing.T) {
	in := "grant_type=client_credentials&client_id=ast-app&client_secret=supersecret&refresh_token=rt-12345"
	out := String(in)
	if strings.Contains(out, "supersecret") || strings.Contains(out, "rt-12345") {
		t.Fatalf("secret leaked: %q", out)
	}
	if !strings.Contains(out, "client_secret=REDACTED") || !strings.Contains(out, "refresh_token=REDACTED") {
		t.Fatalf("expected REDACTED placeholders: %q", out)
	}
	// Non-secret fields preserved.
	if !strings.Contains(out, "grant_type=client_credentials") {
		t.Fatalf("non-secret field stripped: %q", out)
	}
}

func TestStringRedactsJSONSecrets(t *testing.T) {
	in := `{"access_token":"abc123","client_secret":"shh","unrelated":"keep"}`
	out := String(in)
	if strings.Contains(out, "abc123") || strings.Contains(out, "shh") {
		t.Fatalf("json secret leaked: %q", out)
	}
	if !strings.Contains(out, `"unrelated":"keep"`) {
		t.Fatalf("non-secret field stripped: %q", out)
	}
}

func TestStringNoMatch(t *testing.T) {
	in := "no secrets here"
	if got := String(in); got != in {
		t.Fatalf("unexpected modification: %q", got)
	}
}

func TestURLStripsSensitiveQueryParams(t *testing.T) {
	in := "https://example.com/api?foo=bar&access_token=abcdef&q=hello"
	out := URL(in)
	if strings.Contains(out, "abcdef") {
		t.Fatalf("token leaked in URL: %q", out)
	}
	if !strings.Contains(out, "access_token=REDACTED") {
		t.Fatalf("expected redacted access_token: %q", out)
	}
	if !strings.Contains(out, "foo=bar") || !strings.Contains(out, "q=hello") {
		t.Fatalf("benign params lost: %q", out)
	}
}

func TestURLStripsUserInfo(t *testing.T) {
	in := "https://user:hunter2@example.com/path"
	out := URL(in)
	if strings.Contains(out, "hunter2") {
		t.Fatalf("password leaked: %q", out)
	}
}

func TestHeaderRedactsAuthorization(t *testing.T) {
	in := map[string][]string{
		"Authorization": {"Bearer abc.def.ghi"},
		"Content-Type":  {"application/json"},
		"X-Custom":      {"some Bearer inline.token here"},
	}
	out := Header(in)
	if got := out["Authorization"][0]; got != "REDACTED" {
		t.Fatalf("Authorization not fully redacted: %q", got)
	}
	if got := out["Content-Type"][0]; got != "application/json" {
		t.Fatalf("Content-Type modified: %q", got)
	}
	if strings.Contains(out["X-Custom"][0], "inline.token") {
		t.Fatalf("inline bearer not stripped: %q", out["X-Custom"][0])
	}
}
