package tui

import (
	"github.com/charmbracelet/bubbletea"
	"github.com/jackwener/aiview/internal/config"
	"github.com/jackwener/aiview/internal/platform"
	_ "github.com/jackwener/aiview/internal/platform/bilibili"
	_ "github.com/jackwener/aiview/internal/platform/douyin"
	_ "github.com/jackwener/aiview/internal/platform/kuaishou"
	_ "github.com/jackwener/aiview/internal/platform/weibo"
	_ "github.com/jackwener/aiview/internal/platform/xiaohongshu"
	_ "github.com/jackwener/aiview/internal/platform/zhihu"
)

// AppState represents the current state of the TUI
type AppState int

const (
	StatePlatformSelect AppState = iota
	StateHotList
	StateSearch
	StateDetail
)

// Item represents a generic list item (hot search, search result, etc.)
type Item struct {
	Title    string
	Subtitle string
	HotValue int
	URL      string
	ID       string
}

// Model is the main TUI model
type Model struct {
	state       AppState
	platforms   []string
	selected    int
	items       []Item
	searchQuery string
	detail      map[string]interface{}
	err         error
	width       int
	height      int
	cfg         *config.Config
	currentPlat string
	loading     bool
}

// InitialModel creates and initializes the TUI model
func InitialModel() Model {
	platforms := platform.List()
	return Model{
		state:     StatePlatformSelect,
		platforms: platforms,
		selected:  0,
		items:     []Item{},
		cfg:       config.DefaultConfig(),
		loading:   false,
	}
}

// Init implements tea.Model
func (m Model) Init() tea.Cmd {
	return nil
}

// loadHotSearchCmd returns a command that loads hot search data for a platform
func loadHotSearchCmd(platformName string) tea.Cmd {
	return func() tea.Msg {
		p, ok := platform.GetPlatform(platformName)
		if !ok {
			return errMsg{err: nil, msg: "platform not found"}
		}

		cfg := config.DefaultConfig()
		client, err := p.NewClient(cfg)
		if err != nil {
			return errMsg{err: err}
		}

		// Type assert to get the specific client
		switch c := client.(type) {
		case interface{ GetHotSearch(int) (map[string]interface{}, error) }:
			data, err := c.GetHotSearch(50)
			if err != nil {
				return errMsg{err: err}
			}
			return hotSearchLoadedMsg{data: data, platform: platformName}
		case interface{ GetHotSearch() (map[string]interface{}, error) }:
			data, err := c.GetHotSearch()
			if err != nil {
				return errMsg{err: err}
			}
			return hotSearchLoadedMsg{data: data, platform: platformName}
		case interface{ GetHotNotes() (map[string]interface{}, error) }:
			data, err := c.GetHotNotes()
			if err != nil {
				return errMsg{err: err}
			}
			return hotSearchLoadedMsg{data: data, platform: platformName}
		default:
			return errMsg{err: nil, msg: "platform does not support hot search"}
		}
	}
}

// Messages
type errMsg struct {
	err error
	msg string
}

type hotSearchLoadedMsg struct {
	data     map[string]interface{}
	platform string
}
