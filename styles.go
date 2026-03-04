package main

import "charm.land/lipgloss/v2"

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#C084FC")).
			Padding(0, 1)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6B7280"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F87171")).
			Bold(true)

	columnFocusedStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#A855F7")).
				Padding(0, 1)

	columnBlurredStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#374151")).
				Padding(0, 1)

	pendingTitleStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#6B7280")).
				Italic(true).
				Padding(0, 0, 0, 2)

	refChoiceStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9CA3AF")).
			PaddingLeft(2)

	refChoiceSelectedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#A855F7")).
				Bold(true).
				PaddingLeft(2)

	// Header banner
	headerBannerStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#5B21B6")).
				Foreground(lipgloss.Color("#F5F3FF")).
				Bold(true).
				Padding(0, 2)

	headerBadgeStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#4C1D95")).
				Background(lipgloss.Color("#DDD6FE")).
				Bold(true).
				Padding(0, 1)

	// Project panel title
	projectPanelTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#A855F7"))

	moveChooserHeaderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#6B7280")).
				Italic(true).
				PaddingLeft(2)
)

const (
	headerHeight = 2 // banner line + blank line
	footerHeight = 2 // help row + newline
)

// columnSize returns the outer visual width (for style.Width), inner content
// width (for list.SetSize), and content height (for list.SetSize).
//
// In lipgloss v2, Width() sets the outer width (border + padding + content).
// outerW is termW/3 so three columns exactly fill the terminal.
// innerW is outerW minus the frame so the list content fits inside the border.
func columnSize(termW, termH int) (outerW, innerW, h int) {
	frameW := columnBlurredStyle.GetHorizontalFrameSize()
	frameH := columnBlurredStyle.GetVerticalFrameSize()
	outerW = termW / 3
	innerW = outerW - frameW
	h = termH - headerHeight - footerHeight - frameH
	return
}
