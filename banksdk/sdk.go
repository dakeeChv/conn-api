package jdbsdk

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

// APIError is JDB API error.
type APIError struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// Config is SDK configuration.
type Config struct {
	// JDB core banking backend
	BaseURL   string
	UserID    string
	SecretKey string

	// JDB card management system in malaysia
	CmsURL string
	CmsZpk string

	// hmacKey is a share key between JDB mobile backend and JDB core baning backend,
	// used to signed http request signature
	HMacKey []byte

	// timeout is a http client request timeout
	Timeout time.Duration
}

// SDK is a JDB SDK to interect with JDB backend written in Go.
type SDK struct {
	Hc *http.Client

	Cfg     *Config
	idToken string
}

// New create new sdk instance.
func New(ctx context.Context, cfg *Config) (*SDK, error) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		}}
	hc := &http.Client{
		Transport: transport,
		Timeout:   cfg.Timeout,
	}

	sdk := &SDK{Hc: hc, Cfg: cfg}
	if err := sdk.connect(ctx); err != nil {
		return nil, err
	}

	go sdk.notifyIDTokenExpired(ctx)
	return sdk, nil
}

// connect makes http call to perform authentication with JDB backend.
func (s *SDK) connect(ctx context.Context) error {
	body, _ := json.Marshal(map[string]string{
		"requestId": strconv.FormatInt(time.Now().UnixNano(), 10),
		"userId":    s.Cfg.UserID,
		"secretId":  s.Cfg.SecretKey,
	})

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, s.Cfg.BaseURL+"/authentication/login", bytes.NewBuffer(body))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	rsp, err := s.Hc.Do(req)
	if err != nil {
		return fmt.Errorf("connect: http.Do: %v", err)
	}
	defer rsp.Body.Close() //nolint

	content, err := io.ReadAll(rsp.Body)
	if err != nil {
		return fmt.Errorf("connect: read response body: %v", err)
	}

	if rsp.StatusCode != http.StatusOK {
		var apiErr APIError
		if err := json.Unmarshal(content, &apiErr); err != nil {
			return fmt.Errorf("connect: unmarshal json: %s", content)
		}
		return fmt.Errorf("connect: %s: %s", apiErr.Status, apiErr.Message)
	}

	var reply struct {
		Data struct {
			IDToken string `json:"token"`
		} `json:"data"`
	}
	if err := json.Unmarshal(content, &reply); err != nil {
		return fmt.Errorf("connect: unmarshal json: %w", err)
	}
	s.idToken = reply.Data.IDToken
	return nil
}

// notifyIDTokenExpired do infinit loop with period of time
// to perform auto renew token from JDB backend with
// exponential backoff strategy.
func (s *SDK) notifyIDTokenExpired(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Minute)
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return
		case <-ticker.C:
		}

		// reconnect with exponential backoff strategy
		backoff := time.Second
		for s.connect(ctx) != nil {
			select {
			case <-ctx.Done():
				return
			case <-time.After(backoff):
				backoff *= 2
			}
		}
	}
}

// newHTTPReq creates http request with auto signed signature.
func (s *SDK) NewHTTPReq(ctx context.Context, method, url string, body []byte) *http.Request {
	req, _ := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(body))
	req.Header.Set("Authorization", s.idToken)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("SignedHash", hmac256(s.Cfg.HMacKey, body))
	return req
}
