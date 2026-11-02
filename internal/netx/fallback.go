package netx

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// DialError is a proxy dial that was refused or timed out. A direct retry
// was attempted. Error is safe to show: it names the proxy endpoint and
// never includes the transport's "proxyconnect" text or credentials.
type DialError struct {
	Endpoint    string
	DirectTried bool
	Cause       error
}

func (e *DialError) Error() string {
	endpoint := e.Endpoint
	if endpoint == "" {
		endpoint = "proxy"
	}
	if e.DirectTried {
		return "local proxy " + endpoint + " is unreachable; a direct retry also failed"
	}
	return "local proxy " + endpoint + " is unreachable"
}

// Unwrap exposes the original transport error.
func (e *DialError) Unwrap() error { return e.Cause }

// Is matches ErrConnectionFailed so existing proxy handling still applies.
func (e *DialError) Is(target error) bool { return target == ErrConnectionFailed }

// FallbackNotice records the latest dial that was retried without the proxy.
type FallbackNotice struct {
	Endpoint string
	DirectOK bool
	At       time.Time
}

var (
	fallbackMu sync.Mutex
	fallback   FallbackNotice
	fallbackOK bool
)

// NoteFallback returns the latest direct retry, if one has happened in this
// process. The endpoint is host:port with no user info.
func NoteFallback() (FallbackNotice, bool) {
	fallbackMu.Lock()
	defer fallbackMu.Unlock()
	return fallback, fallbackOK
}

func recordFallback(n FallbackNotice) {
	if n.At.IsZero() {
		n.At = time.Now()
	}
	fallbackMu.Lock()
	fallback, fallbackOK = n, true
	fallbackMu.Unlock()
}

func (m *mappingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	proxyURL, proxyErr := m.proxyFor(req)
	resp, err := m.base.RoundTrip(req)
	if err == nil {
		if m.manual && resp.StatusCode == http.StatusProxyAuthRequired {
			_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 4<<10))
			_ = resp.Body.Close()
			return nil, &Error{Kind: ErrAuthenticationFailed}
		}
		return resp, nil
	}
	if m.canFallback(req.Context(), err, proxyURL, proxyErr) {
		endpoint := endpointOf(proxyURL)
		if endpoint == "" {
			endpoint = m.endpoint
		}
		direct, derr := m.roundTripDirect(req)
		if derr == nil {
			recordFallback(FallbackNotice{Endpoint: endpoint, DirectOK: true, At: time.Now()})
			return direct, nil
		}
		recordFallback(FallbackNotice{Endpoint: endpoint, DirectOK: false, At: time.Now()})
		return nil, &DialError{Endpoint: endpoint, DirectTried: true, Cause: err}
	}
	return nil, m.mapError(req.Context(), err)
}

func (m *mappingTransport) proxyFor(req *http.Request) (*url.URL, error) {
	if m.base == nil || m.base.Proxy == nil || req == nil {
		return nil, nil
	}
	return m.base.Proxy(req)
}

// canFallback is true only for a dial refusal or timeout through a proxy
// that was actually selected. Resolver errors, authentication failures, and
// a canceled request stay fail-closed.
func (m *mappingTransport) canFallback(ctx context.Context, err error, proxy *url.URL, proxyErr error) bool {
	if err == nil || proxyErr != nil || proxy == nil {
		return false
	}
	if ctx.Err() != nil || errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return false
	}
	if errors.Is(err, errProxyAuthRequired) || socksAuthenticationFailure(err) {
		return false
	}
	var proxyErrKind *Error
	if errors.As(err, &proxyErrKind) && errors.Is(proxyErrKind, ErrAuthenticationFailed) {
		return false
	}
	return dialFailure(err)
}

func (m *mappingTransport) roundTripDirect(req *http.Request) (*http.Response, error) {
	clone := req.Clone(req.Context())
	if req.Body != nil {
		if req.GetBody == nil {
			return nil, errors.New("request body cannot be retried")
		}
		body, err := req.GetBody()
		if err != nil {
			return nil, err
		}
		clone.Body = body
	}
	direct := m.base.Clone()
	direct.Proxy = nil
	return direct.RoundTrip(clone)
}

func dialFailure(err error) bool {
	var op *net.OpError
	if errors.As(err, &op) && (op.Op == "dial" || strings.Contains(op.Op, "proxyconnect")) {
		return true
	}
	msg := strings.ToLower(err.Error())
	if strings.Contains(msg, "proxyconnect") || strings.Contains(msg, "connection refused") ||
		strings.Contains(msg, "i/o timeout") || strings.Contains(msg, "no route to host") ||
		strings.Contains(msg, "network is unreachable") {
		return true
	}
	var ne net.Error
	return errors.As(err, &ne) && ne.Timeout()
}

func endpointOf(u *url.URL) string {
	if u == nil {
		return ""
	}
	host := u.Hostname()
	if host == "" {
		return ""
	}
	port := u.Port()
	if port == "" {
		switch strings.ToLower(u.Scheme) {
		case "https":
			port = "443"
		case "socks5", "socks5h":
			port = "1080"
		default:
			port = "80"
		}
	}
	return net.JoinHostPort(host, port)
}
