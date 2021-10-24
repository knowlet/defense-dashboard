package quest

import (
	"context"
	"io"
	"net/http"
	"time"
)

func request(method, url, hostname string, body io.Reader) (*http.Response, error) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	req, _ := http.NewRequestWithContext(ctx, method, url, body)
	req.Host = hostname
	if body != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	client := &http.Client{
		// Timeout: 10 * time.Second,
	}
	go func() {
		time.Sleep(time.Second * 10)
		cancel()
	}()
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
