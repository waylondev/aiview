package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// View renders the TUI based on current state
func (m Model) View() string {
	if m.width == 0 {
		return "Initializing..."
	}

	var content string
	switch m.state {
	case StatePlatformSelect:
		content = m.viewPlatformSelect()
	case StateHotList:
		content = m.viewHotList()
	case StateSearch:
		content = m.viewSearch()
	case StateDetail:
		content = m.viewDetail()
	default:
		content = "Unknown state"
	}

	return content
}

// viewPlatformSelect renders the platform selection screen
func (m Model) viewPlatformSelect() string {
	var b strings.Builder

	b.WriteString(Title.Render("🎯 AIView - Platform Selection"))
	b.WriteString("\n\n")
	b.WriteString(Subtitle.Render("Select a platform to browse hot search:"))
	b.WriteString("\n\n")

	for i, p := range m.platforms {
		if i == m.selected {
			b.WriteString(SelectedItem.Render(p))
		} else {
			b.WriteString(NormalItem.Render(p))
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(Help.Render("↑/↓: Navigate • Enter: Select • q: Quit"))

	return b.String()
}

// viewHotList renders the hot search list
func (m Model) viewHotList() string {
	var b strings.Builder

	title := fmt.Sprintf("🔥 %s - Hot Search", strings.Title(m.currentPlat))
	b.WriteString(Title.Render(title))
	b.WriteString("\n")

	if m.loading {
		b.WriteString(Info.Render("Loading..."))
		return b.String()
	}

	if m.err != nil {
		b.WriteString(Error.Render(fmt.Sprintf("Error: %v", m.err)))
		b.WriteString("\n")
		b.WriteString(Help.Render("Esc: Back • q: Quit"))
		return b.String()
	}

	if len(m.items) == 0 {
		b.WriteString(Info.Render("No items found"))
		b.WriteString("\n")
		b.WriteString(Help.Render("Esc: Back • q: Quit"))
		return b.String()
	}

	// Calculate visible items based on terminal height
	maxVisible := m.height - 8 // Reserve space for header/footer
	if maxVisible < 5 {
		maxVisible = 5
	}
	startIdx := 0
	if m.selected >= maxVisible {
		startIdx = m.selected - maxVisible + 1
	}
	endIdx := startIdx + maxVisible
	if endIdx > len(m.items) {
		endIdx = len(m.items)
	}

	for i := startIdx; i < endIdx; i++ {
		item := m.items[i]
		prefix := fmt.Sprintf("%2d. ", i+1)

		if i == m.selected {
			b.WriteString(SelectedItem.Render(""))
			b.WriteString(lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#10B981")).Render(prefix + item.Title))
		} else {
			b.WriteString("    ")
			b.WriteString(prefix + item.Title)
		}

		if item.HotValue > 0 {
			b.WriteString("  ")
			b.WriteString(HotValue.Render(fmt.Sprintf("(热度: %d)", item.HotValue)))
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(Help.Render("↑/↓: Navigate • Enter: Detail • Esc: Back • q: Quit"))

	// Show scroll indicator if needed
	if len(m.items) > maxVisible {
		b.WriteString("\n")
		b.WriteString(Info.Render(fmt.Sprintf("Showing %d-%d of %d", startIdx+1, endIdx, len(m.items))))
	}

	return b.String()
}

// viewSearch renders the search screen
func (m Model) viewSearch() string {
	var b strings.Builder

	b.WriteString(Title.Render("🔍 Search"))
	b.WriteString("\n\n")
	b.WriteString(Subtitle.Render(fmt.Sprintf("Query: %s", m.searchQuery)))
	b.WriteString("\n\n")

	if m.loading {
		b.WriteString(Info.Render("Searching..."))
		return b.String()
	}

	if len(m.items) == 0 {
		b.WriteString(Info.Render("No results found"))
	} else {
		for i, item := range m.items {
			if i == m.selected {
				b.WriteString(SelectedItem.Render(item.Title))
			} else {
				b.WriteString(NormalItem.Render(item.Title))
			}
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")
	b.WriteString(Help.Render("Esc: Back • q: Quit"))

	return b.String()
}

// viewDetail renders the detail view
func (m Model) viewDetail() string {
	var b strings.Builder

	if m.selected >= len(m.items) {
		b.WriteString(Error.Render("No item selected"))
		return b.String()
	}

	item := m.items[m.selected]
	b.WriteString(Title.Render("📄 Detail"))
	b.WriteString("\n\n")

	b.WriteString(Subtitle.Render("Title:"))
	b.WriteString("\n")
	b.WriteString("  " + item.Title)
	b.WriteString("\n\n")

	if item.Subtitle != "" {
		b.WriteString(Subtitle.Render("Description:"))
		b.WriteString("\n")
		b.WriteString("  " + item.Subtitle)
		b.WriteString("\n\n")
	}

	if item.URL != "" {
		b.WriteString(Subtitle.Render("URL:"))
		b.WriteString("\n")
		b.WriteString("  " + item.URL)
		b.WriteString("\n\n")
	}

	if item.HotValue > 0 {
		b.WriteString(Subtitle.Render("Hot Value:"))
		b.WriteString("\n")
		b.WriteString(fmt.Sprintf("  %d", item.HotValue))
		b.WriteString("\n\n")
	}

	b.WriteString(Help.Render("Esc: Back • q: Quit"))

	return b.String()
}
