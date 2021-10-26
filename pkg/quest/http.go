package quest

import (
	"context"
	"crypto/tls"
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
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
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

func reqjson(method, url, hostname string, body io.Reader) (*http.Response, error) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	req, _ := http.NewRequestWithContext(ctx, method, url, body)
	req.Host = hostname
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
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

func reqBaseJson(method, url, hostname, token string, body io.Reader) (*http.Response, error) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	req, _ := http.NewRequestWithContext(ctx, method, url, body)
	req.Host = hostname
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.163 Safari/537.36")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Authorization", "Basic "+token)
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
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
