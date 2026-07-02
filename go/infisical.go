package infisical

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type SDK struct {
	siteURL     string
	accessToken string
	httpClient  *http.Client
}

type Options struct {
	SiteURL string
}

type LoginResponse struct {
	AccessToken string `json:"accessToken"`
}

type Secret struct {
	ID             string `json:"id"`
	WorkspaceID    string `json:"workspace"`
	Environment    string `json:"environment"`
	SecretKey      string `json:"secretKey"`
	SecretValue    string `json:"secretValue"`
	SecretComment  string `json:"secretComment"`
	SecretPath     string `json:"secretPath"`
	Version        int    `json:"version"`
	Type           string `json:"type"`
	CreatedAt      string `json:"createdAt"`
	UpdatedAt      string `json:"updatedAt"`
}

type SecretsResponse struct {
	Secrets []Secret `json:"secrets"`
}

func New(opts Options) *SDK {
	siteURL := opts.SiteURL
	if siteURL == "" {
		siteURL = "https://app.infisical.com"
	}
	siteURL = strings.TrimRight(siteURL, "/")

	return &SDK{
		siteURL:    siteURL,
		httpClient: &http.Client{},
	}
}

func (s *SDK) request(method, path string, params map[string]string, body any) ([]byte, error) {
	u, err := url.Parse(s.siteURL + path)
	if err != nil {
		return nil, fmt.Errorf("invalid url: %w", err)
	}

	if params != nil {
		q := u.Query()
		for k, v := range params {
			q.Set(k, v)
		}
		u.RawQuery = q.Encode()
	}

	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal body: %w", err)
		}
		reqBody = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, u.String(), reqBody)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	if s.accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+s.accessToken)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	res, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http do: %w", err)
	}
	defer res.Body.Close()

	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	if res.StatusCode >= 400 {
		var errResp struct {
			Message string `json:"message"`
			Error   string `json:"error"`
		}
		json.Unmarshal(respBody, &errResp)
		msg := errResp.Message
		if msg == "" {
			msg = errResp.Error
		}
		if msg == "" {
			msg = fmt.Sprintf("HTTP %d", res.StatusCode)
		}
		return nil, fmt.Errorf("%s", msg)
	}

	return respBody, nil
}

func (s *SDK) Login(id, secret string) error {
	body, err := s.request("POST", "/api/v1/auth/universal-auth/login", nil, map[string]string{
		"clientId":     id,
		"clientSecret": secret,
	})
	if err != nil {
		return err
	}

	var resp LoginResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return fmt.Errorf("parse login response: %w", err)
	}

	s.accessToken = resp.AccessToken
	return nil
}

func (s *SDK) Secrets(environment, projectID string) (*SecretsResponse, error) {
	if s.accessToken == "" {
		return nil, fmt.Errorf("not authenticated. call Login() first")
	}

	body, err := s.request("GET", "/api/v3/secrets/raw", map[string]string{
		"environment": environment,
		"workspaceId": projectID,
	}, nil)
	if err != nil {
		return nil, err
	}

	var resp SecretsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse secrets response: %w", err)
	}

	return &resp, nil
}