package aiview

import (
	"fmt"

	"github.com/jackwener/aiview/internal/config"
	"github.com/jackwener/aiview/internal/platform"
	_ "github.com/jackwener/aiview/internal/platform/bilibili"
	_ "github.com/jackwener/aiview/internal/platform/douyin"
	_ "github.com/jackwener/aiview/internal/platform/xiaohongshu"
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
