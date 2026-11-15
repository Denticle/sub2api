package oauthprovider

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const feedFixture = `{"code":0,"data":{"entries":[
 {"id":"a","title":"仙侠短片","image":"https://cdn.example/a.jpg","mode":"video","credit":{"name":"Seedance 2.5"},"prompt":"御剑飞行"},
 {"id":"b","title":"产品海报","image":"https://cdn.example/b.jpg","mode":"image","credit":"蛋黄堡"},
 {"id":"c","title":"没有图","image":"","mode":"image"},
 {"id":"d","title":"不安全图","image":"http://cdn.example/d.jpg","mode":"image"},
 {"id":"e","title":"  ","image":"https://cdn.example/e.jpg","mode":"image"}
]}}`

func TestFetchInspirationsKeepsOnlyUsableEntries(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, inspirationFeedPath, r.URL.Path)
		_, _ = w.Write([]byte(feedFixture))
	}))
	defer srv.Close()

	feed, err := fetchInspirations(context.Background(), srv.Client(), srv.URL)
	require.NoError(t, err)
	require.Equal(t, 5, feed.Total)
	require.Len(t, feed.Items, 2)
	require.Equal(t, Inspiration{ID: "a", Title: "仙侠短片", Image: "https://cdn.example/a.jpg", Mode: "video", Credit: "Seedance 2.5", Prompt: "御剑飞行"}, feed.Items[0])
	require.Equal(t, "蛋黄堡", feed.Items[1].Credit)
	require.Equal(t, "image", feed.Items[1].Mode)
}

func TestInspirationCacheReusesAndFallsBack(t *testing.T) {
	var calls atomic.Int32
	fail := atomic.Bool{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		calls.Add(1)
		if fail.Load() {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		_, _ = w.Write([]byte(feedFixture))
	}))
	defer srv.Close()

	c := &inspirationCache{client: srv.Client()}
	first := c.get(context.Background(), srv.URL)
	require.Len(t, first.Items, 2)
	_ = c.get(context.Background(), srv.URL)
	require.EqualValues(t, 1, calls.Load(), "second call within TTL is served from cache")

	// 缓存过期后上游出错：沿用上一次的结果。
	fail.Store(true)
	c.expires = time.Now().Add(-time.Second)
	stale := c.get(context.Background(), srv.URL)
	require.Len(t, stale.Items, 2)
	require.EqualValues(t, 2, calls.Load())

	// 换了一个取不到的地址：不能把旧地址的数据当成新地址的。
	c.expires = time.Now().Add(-time.Second)
	other := c.get(context.Background(), srv.URL+"/other")
	require.Empty(t, other.Items)
}
