package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jackwener/aiview/internal/helper"
)

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	case hotSearchLoadedMsg:
		m.loading = false
		m.currentPlat = msg.platform
		m.items = parseHotSearchData(msg.data, msg.platform)
		m.state = StateHotList
		m.selected = 0
		m.err = nil
		return m, nil

	case errMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
		} else if msg.msg != "" {
			m.err = fmt.Errorf("%s", msg.msg)
		}
		return m, nil
	}

	return m, nil
}

// handleKeyPress processes keyboard input
func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit

	case "esc":
		return m.handleEscape()

	case "up", "k":
		if m.selected > 0 {
			m.selected--
		}
		return m, nil

	case "down", "j":
		if m.selected < len(m.items)-1 {
			m.selected++
		}
		return m, nil

	case "enter":
		return m.handleEnter()
	}

	return m, nil
}

// handleEscape handles the escape key based on current state
func (m Model) handleEscape() (tea.Model, tea.Cmd) {
	switch m.state {
	case StateHotList, StateSearch:
		m.state = StatePlatformSelect
		m.selected = 0
		m.items = []Item{}
		m.err = nil
		return m, nil

	case StateDetail:
		m.state = StateHotList
		m.detail = nil
		return m, nil
	}
	return m, nil
}

// handleEnter handles the enter key based on current state
func (m Model) handleEnter() (tea.Model, tea.Cmd) {
	switch m.state {
	case StatePlatformSelect:
		if m.selected < len(m.platforms) {
			platformName := m.platforms[m.selected]
			m.loading = true
			return m, loadHotSearchCmd(platformName)
		}

	case StateHotList:
		if m.selected < len(m.items) {
			m.state = StateDetail
			return m, nil
		}
	}
	return m, nil
}

// parseHotSearchData parses the hot search response into items
func parseHotSearchData(data map[string]interface{}, platform string) []Item {
	items := []Item{}

	switch platform {
	case "bilibili":
		// Bilibili returns data in "data" -> "trending" -> "list"
		d := helper.GetMap(data, "data")
		trending := helper.GetMap(d, "trending")
		list := helper.GetSlice(trending, "list")
		for i, item := range list {
			m, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			items = append(items, Item{
				Title:    helper.GetString(m, "keyword"),
				HotValue: helper.GetInt(m, "show_info"),
				URL:      helper.GetString(m, "url"),
				ID:       fmt.Sprintf("%d", i+1),
			})
		}

	case "douyin":
		// Douyin returns data in "data" -> "word_list"
		d := helper.GetMap(data, "data")
		wordList := helper.GetSlice(d, "word_list")
		for i, item := range wordList {
			m, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			items = append(items, Item{
				Title:    helper.GetString(m, "word"),
				HotValue: helper.GetInt(m, "hot_value"),
				ID:       fmt.Sprintf("%d", i+1),
			})
		}

	case "xiaohongshu":
		// Xiaohongshu returns data in "data"
		d := helper.GetMap(data, "data")
		for key, val := range d {
			if m, ok := val.(map[string]interface{}); ok {
				items = append(items, Item{
					Title: helper.GetString(m, "name"),
					ID:    key,
				})
			}
		}
	}

	return items
}

// truncateString truncates a string to maxLen and adds ellipsis if needed
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

// padRight pads a string with spaces to reach the desired width
func padRight(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-len(s))
}
