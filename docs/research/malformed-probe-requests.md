# Malformed service-probe requests

Base: 5773cd40de16cc5f020ae49f804f33a1bf8d2498.

Both HTTP helpers discarded the error returned by NewRequestWithContext and
immediately wrote req.Host. A malformed configured address such as `%zz` or a
non-numeric port makes request construction return a nil request. The resulting
nil-pointer panic happens in a probe goroutine, outside Gin's request recovery,
and can terminate the dashboard instead of recording the service as down.

Return the constructor error after cancelling the newly created context. Existing
probe callers already handle errors through their service-down path. This change
does not alter TLS verification, redirects, timeout duration, scoring, or valid
request semantics. Timeout/transport lifecycle cleanup is outside this patch.

Source: Go net/http NewRequestWithContext:
https://pkg.go.dev/net/http#NewRequestWithContext

Regression tests exercise malformed URL escapes, ports, and methods in both
helpers. All six cases panic on the unmodified base; the fixed helpers return an
error and nil response. Two local HTTP server controls verify POST body, explicit
Host, Content-Type, response-body readability, and preserved redirect refusal.
No configured competition endpoints or credentials are used.
