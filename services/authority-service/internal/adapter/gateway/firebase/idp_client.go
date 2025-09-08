package firebase

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    outgateway "github.com/YukiOnishi1129/youtube-analytics/services/authority-service/internal/port/output/gateway"
)

type Client struct {
    APIKey string
    HTTP   *http.Client
}

func New(apiKey string) *Client { return &Client{APIKey: apiKey, HTTP: http.DefaultClient} }

var _ outgateway.IdentityProvider = (*Client)(nil)

func (c *Client) SignUp(ctx context.Context, email, password string) (outgateway.AuthTokens, error) {
    endpoint := fmt.Sprintf("https://identitytoolkit.googleapis.com/v1/accounts:signUp?key=%s", c.APIKey)
    body := map[string]any{"email": email, "password": password, "returnSecureToken": true}
    var resp struct{ IdToken, RefreshToken string; ExpiresIn string }
    if err := c.post(ctx, endpoint, body, &resp); err != nil { return outgateway.AuthTokens{}, err }
    return outgateway.AuthTokens{IDToken: resp.IdToken, RefreshToken: resp.RefreshToken}, nil
}

func (c *Client) SignIn(ctx context.Context, email, password string) (outgateway.AuthTokens, error) {
    endpoint := fmt.Sprintf("https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=%s", c.APIKey)
    body := map[string]any{"email": email, "password": password, "returnSecureToken": true}
    var resp struct{ IdToken, RefreshToken, ExpiresIn string }
    if err := c.post(ctx, endpoint, body, &resp); err != nil { return outgateway.AuthTokens{}, err }
    // ExpiresIn is string seconds per Firebase
    return outgateway.AuthTokens{IDToken: resp.IdToken, RefreshToken: resp.RefreshToken}, nil
}

func (c *Client) SignOut(ctx context.Context, refreshToken string) error {
    endpoint := fmt.Sprintf("https://securetoken.googleapis.com/v1/token:revoke?key=%s", c.APIKey)
    body := map[string]any{"token": refreshToken}
    var resp map[string]any
    return c.post(ctx, endpoint, body, &resp)
}

func (c *Client) ResetPassword(ctx context.Context, email string) error {
    endpoint := fmt.Sprintf("https://identitytoolkit.googleapis.com/v1/accounts:sendOobCode?key=%s", c.APIKey)
    body := map[string]any{"requestType": "PASSWORD_RESET", "email": email}
    var resp map[string]any
    return c.post(ctx, endpoint, body, &resp)
}

func (c *Client) post(ctx context.Context, url string, body map[string]any, out any) error {
    b, _ := json.Marshal(body)
    req, _ := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(b))
    req.Header.Set("Content-Type", "application/json")
    res, err := c.HTTP.Do(req)
    if err != nil { return err }
    defer res.Body.Close()
    if res.StatusCode/100 != 2 {
        var e map[string]any
        _ = json.NewDecoder(res.Body).Decode(&e)
        return fmt.Errorf("firebase api error: %v", e)
    }
    if out != nil {
        return json.NewDecoder(res.Body).Decode(out)
    }
    return nil
}

