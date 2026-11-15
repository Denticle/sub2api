package oauthprovider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

// 首页作品墙从视频工作台的灵感库取数据。工作台前面的 WAF 会拦浏览器直连，
// 所以由主站后端代取，并且只透出首页要用的字段。

const (
	inspirationFeedPath  = "/api/public/inspirations/feed"
	inspirationCacheTTL  = 10 * time.Minute
	inspirationFailTTL   = time.Minute
	inspirationMaxBody   = 4 << 20
	inspirationMaxItems  = 60
	inspirationFetchWait = 8 * time.Second
)

// Inspiration 是作品墙的一张卡片。
type Inspiration struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Image  string `json:"image"`
	Mode   string `json:"mode"`
	Credit string `json:"credit,omitempty"`
	Prompt string `json:"prompt,omitempty"`
}

type inspirationFeed struct {
	Items []Inspiration `json:"items"`
	Total int           `json:"total"`
}

type inspirationCache struct {
	mu      sync.Mutex
	base    string
	feed    inspirationFeed
	expires time.Time
	client  *http.Client
}

var feedCache = &inspirationCache{client: &http.Client{Timeout: inspirationFetchWait}}

func (c *inspirationCache) get(ctx context.Context, base string) inspirationFeed {
	c.mu.Lock()
	defer c.mu.Unlock()
	if base == c.base && time.Now().Before(c.expires) {
		return c.feed
	}
	feed, err := fetchInspirations(ctx, c.client, base)
	if err != nil {
		// 取不到就沿用上一次的结果（换了地址则清空），一分钟后再试，不让首页报错。
		if base != c.base {
			c.feed = inspirationFeed{Items: []Inspiration{}}
		}
		c.base = base
		c.expires = time.Now().Add(inspirationFailTTL)
		return c.feed
	}
	c.base = base
	c.feed = feed
	c.expires = time.Now().Add(inspirationCacheTTL)
	return feed
}

func fetchInspirations(ctx context.Context, client *http.Client, base string) (inspirationFeed, error) {
	u, err := url.Parse(base + inspirationFeedPath)
	if err != nil || (u.Scheme != "https" && u.Scheme != "http") {
		return inspirationFeed{}, fmt.Errorf("invalid studio url")
	}
	q := u.Query()
	q.Set("limit", "500")
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return inspirationFeed{}, err
	}
	req.Header.Set("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return inspirationFeed{}, err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return inspirationFeed{}, fmt.Errorf("studio feed status %d", resp.StatusCode)
	}

	var body struct {
		Data struct {
			Entries []struct {
				ID     string `json:"id"`
				Title  string `json:"title"`
				Image  string `json:"image"`
				Mode   string `json:"mode"`
				Credit any    `json:"credit"`
				Prompt string `json:"prompt"`
			} `json:"entries"`
		} `json:"data"`
	}
	if err := json.NewDecoder(io.LimitReader(resp.Body, inspirationMaxBody)).Decode(&body); err != nil {
		return inspirationFeed{}, err
	}

	items := make([]Inspiration, 0, inspirationMaxItems)
	for _, e := range body.Data.Entries {
		if !isHTTPSURL(e.Image) || strings.TrimSpace(e.Title) == "" {
			continue
		}
		mode := "image"
		if e.Mode == "video" {
			mode = "video"
		}
		items = append(items, Inspiration{
			ID:     e.ID,
			Title:  strings.TrimSpace(e.Title),
			Image:  e.Image,
			Mode:   mode,
			Credit: creditName(e.Credit),
			Prompt: truncateRunes(strings.TrimSpace(e.Prompt), 600),
		})
		if len(items) == inspirationMaxItems {
			break
		}
	}
	return inspirationFeed{Items: items, Total: len(body.Data.Entries)}, nil
}

func isHTTPSURL(raw string) bool {
	u, err := url.Parse(raw)
	return err == nil && u.Scheme == "https" && u.Host != ""
}

// credit 在上游有时是字符串，有时是 {name: ...} 对象。
func creditName(v any) string {
	switch c := v.(type) {
	case string:
		return truncateRunes(strings.TrimSpace(c), 40)
	case map[string]any:
		if name, ok := c["name"].(string); ok {
			return truncateRunes(strings.TrimSpace(name), 40)
		}
	}
	return ""
}

func truncateRunes(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n])
}

// Inspirations GET /api/v1/public/inspirations —— 首页作品墙
func (h *Handler) Inspirations(c *gin.Context) {
	cfg := h.svc.PublicConfig(c.Request.Context())
	if !cfg.Enabled || cfg.AppURL == "" {
		response.Success(c, inspirationFeed{Items: []Inspiration{}})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), inspirationFetchWait)
	defer cancel()
	c.Header("Cache-Control", "public, max-age=300")
	response.Success(c, feedCache.get(ctx, cfg.AppURL))
}
