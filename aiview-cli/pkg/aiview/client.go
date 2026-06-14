package aiview

import (
	"fmt"

	"github.com/jackwener/aiview/internal/config"
	"github.com/jackwener/aiview/internal/platform"
	_ "github.com/jackwener/aiview/internal/platform/bilibili"
	_ "github.com/jackwener/aiview/internal/platform/douyin"
	_ "github.com/jackwener/aiview/internal/platform/kuaishou"
	_ "github.com/jackwener/aiview/internal/platform/weibo"
	_ "github.com/jackwener/aiview/internal/platform/xiaohongshu"
	_ "github.com/jackwener/aiview/internal/platform/zhihu"
)

// Client provides a unified interface for accessing multiple platforms.
type Client struct {
	platform platform.Platform
	config   *config.Config
}

// New creates a new Client for the specified platform.
func New(platformName string) (*Client, error) {
	p, ok := platform.GetPlatform(platformName)
	if !ok {
		return nil, fmt.Errorf("platform %q not supported", platformName)
	}

	cfg := config.DefaultConfig()
	return &Client{
		platform: p,
		config:   cfg,
	}, nil
}

// PlatformName returns the current platform name.
func (c *Client) PlatformName() string {
	return c.platform.Name()
}

// BilibiliClient returns the underlying Bilibili client.
// Returns error if the platform is not bilibili.
func (c *Client) BilibiliClient() (*BilibiliClient, error) {
	if c.platform.Name() != "bilibili" {
		return nil, fmt.Errorf("not a bilibili client")
	}

	client, err := c.platform.NewClient(c.config)
	if err != nil {
		return nil, err
	}

	return &BilibiliClient{client: client}, nil
}

// DouyinClient returns the underlying Douyin client.
// Returns error if the platform is not douyin.
func (c *Client) DouyinClient() (*DouyinClient, error) {
	if c.platform.Name() != "douyin" {
		return nil, fmt.Errorf("not a douyin client")
	}

	client, err := c.platform.NewClient(c.config)
	if err != nil {
		return nil, err
	}

	return &DouyinClient{client: client}, nil
}

// XiaohongshuClient returns the underlying Xiaohongshu client.
// Returns error if the platform is not xiaohongshu.
func (c *Client) XiaohongshuClient() (*XiaohongshuClient, error) {
	if c.platform.Name() != "xiaohongshu" {
		return nil, fmt.Errorf("not a xiaohongshu client")
	}

	client, err := c.platform.NewClient(c.config)
	if err != nil {
		return nil, err
	}

	return &XiaohongshuClient{client: client}, nil
}

// WeiboClient returns the underlying Weibo client.
// Returns error if the platform is not weibo.
func (c *Client) WeiboClient() (*WeiboClient, error) {
	if c.platform.Name() != "weibo" {
		return nil, fmt.Errorf("not a weibo client")
	}

	client, err := c.platform.NewClient(c.config)
	if err != nil {
		return nil, err
	}

	return &WeiboClient{client: client}, nil
}

// KuaishouClient returns the underlying Kuaishou client.
// Returns error if the platform is not kuaishou.
func (c *Client) KuaishouClient() (*KuaishouClient, error) {
	if c.platform.Name() != "kuaishou" {
		return nil, fmt.Errorf("not a kuaishou client")
	}

	client, err := c.platform.NewClient(c.config)
	if err != nil {
		return nil, err
	}

	return &KuaishouClient{client: client}, nil
}

// ZhihuClient returns the underlying Zhihu client.
// Returns error if the platform is not zhihu.
func (c *Client) ZhihuClient() (*ZhihuClient, error) {
	if c.platform.Name() != "zhihu" {
		return nil, fmt.Errorf("not a zhihu client")
	}

	client, err := c.platform.NewClient(c.config)
	if err != nil {
		return nil, err
	}

	return &ZhihuClient{client: client}, nil
}
