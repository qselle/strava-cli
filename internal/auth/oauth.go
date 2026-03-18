package auth

import (
	"bufio"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

const (
	AuthURL  = "https://www.strava.com/oauth/authorize"
	TokenURL = "https://www.strava.com/api/v3/oauth/token"
	Scopes   = "read,activity:read_all,profile:read_all"
)

type OAuthConfig struct {
	ClientID     string
	ClientSecret string
}

// LoginManual prints the auth URL and waits for the user to paste the code.
// Works on headless servers with no browser.
func LoginManual(ctx context.Context, cfg OAuthConfig) (*Token, error) {
	redirectURI := "http://localhost"
	state, err := randomState()
	if err != nil {
		return nil, fmt.Errorf("generating state: %w", err)
	}

	authURL := buildAuthURL(cfg.ClientID, redirectURI, state)

	fmt.Println("Open this URL in your browser:")
	fmt.Println()
	fmt.Println(authURL)
	fmt.Println()
	fmt.Println("Authorize the app, then copy the 'code' parameter from the redirect URL.")
	fmt.Println("The redirect URL will look like: http://localhost?state=...&code=THE_CODE")
	fmt.Println()
	fmt.Print("Paste the code here: ")

	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return nil, fmt.Errorf("no input received")
	}
	code := strings.TrimSpace(scanner.Text())
	if code == "" {
		return nil, fmt.Errorf("empty code")
	}

	return exchangeCode(ctx, cfg, code, redirectURI)
}

// LoginBrowser starts a local callback server and opens the browser.
// Works on machines with a browser.
func LoginBrowser(ctx context.Context, cfg OAuthConfig) (*Token, error) {
	state, err := randomState()
	if err != nil {
		return nil, fmt.Errorf("generating state: %w", err)
	}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, fmt.Errorf("starting callback server: %w", err)
	}
	defer listener.Close()

	port := listener.Addr().(*net.TCPAddr).Port
	redirectURI := fmt.Sprintf("http://127.0.0.1:%d/callback", port)

	authURL := buildAuthURL(cfg.ClientID, redirectURI, state)
	fmt.Printf("Opening browser for Strava authorization...\n")
	fmt.Printf("If it doesn't open, visit:\n%s\n\n", authURL)
	openBrowser(authURL)

	codeCh := make(chan string, 1)
	errCh := make(chan error, 1)

	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("state") != state {
			errCh <- fmt.Errorf("invalid state parameter")
			http.Error(w, "Invalid state", http.StatusBadRequest)
			return
		}
		if errMsg := r.URL.Query().Get("error"); errMsg != "" {
			errCh <- fmt.Errorf("authorization denied: %s", errMsg)
			fmt.Fprintf(w, "<html><body><h1>Authorization denied</h1><p>You can close this window.</p></body></html>")
			return
		}
		code := r.URL.Query().Get("code")
		if code == "" {
			errCh <- fmt.Errorf("no authorization code received")
			http.Error(w, "No code", http.StatusBadRequest)
			return
		}
		codeCh <- code
		fmt.Fprintf(w, "<html><body><h1>Authorized!</h1><p>You can close this window and return to the terminal.</p></body></html>")
	})

	server := &http.Server{Handler: mux}
	go func() { _ = server.Serve(listener) }()
	defer server.Shutdown(ctx)

	var code string
	select {
	case code = <-codeCh:
	case err := <-errCh:
		return nil, err
	case <-time.After(5 * time.Minute):
		return nil, fmt.Errorf("authorization timed out after 5 minutes")
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	return exchangeCode(ctx, cfg, code, redirectURI)
}

func RefreshAccessToken(ctx context.Context, cfg OAuthConfig, refreshToken string) (*Token, error) {
	data := url.Values{
		"client_id":     {cfg.ClientID},
		"client_secret": {cfg.ClientSecret},
		"grant_type":    {"refresh_token"},
		"refresh_token": {refreshToken},
	}

	req, err := http.NewRequestWithContext(ctx, "POST", TokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("refreshing token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token refresh failed with status %d", resp.StatusCode)
	}

	var token Token
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, fmt.Errorf("decoding token response: %w", err)
	}

	return &token, nil
}

func exchangeCode(ctx context.Context, cfg OAuthConfig, code, redirectURI string) (*Token, error) {
	data := url.Values{
		"client_id":     {cfg.ClientID},
		"client_secret": {cfg.ClientSecret},
		"code":          {code},
		"grant_type":    {"authorization_code"},
	}

	req, err := http.NewRequestWithContext(ctx, "POST", TokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("exchanging code: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token exchange failed with status %d", resp.StatusCode)
	}

	var token Token
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, fmt.Errorf("decoding token response: %w", err)
	}

	return &token, nil
}

func buildAuthURL(clientID, redirectURI, state string) string {
	params := url.Values{
		"client_id":       {clientID},
		"redirect_uri":    {redirectURI},
		"response_type":   {"code"},
		"scope":           {Scopes},
		"state":           {state},
		"approval_prompt": {"auto"},
	}
	return AuthURL + "?" + params.Encode()
}

func randomState() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func openBrowser(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	_ = cmd.Start()
}
