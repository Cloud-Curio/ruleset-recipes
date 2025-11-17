/*
 * Data Ingestion TUI - Main TUI Interface
 * 
 * This package implements the main Terminal User Interface using Bubble Tea.
 * It provides a comprehensive interface for managing data ingestion operations.
 * 
 * Features:
 * - Multi-panel layout with navigation
 * - Real-time status updates
 * - Interactive forms and dialogs
 * - Keyboard shortcuts and mouse support
 * - Responsive design
 * 
 * Author: Codegen AI Assistant
 * Created: 2024
 * License: MIT
 */

package ui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/pkg/errors"

	"data-ingestion-tui/internal/app"
	"data-ingestion-tui/internal/config"
	"data-ingestion-tui/internal/models"
)

// @decorator: TUI
// @description: Main TUI application structure
// @version: 1.0.0
// @author: Codegen AI Assistant

// TUI represents the main Terminal User Interface
type TUI struct {
	app    *app.Application
	config *config.UIConfig
	logger *log.Logger
	
	// UI state
	width      int
	height     int
	activePane PaneType
	panes      map[PaneType]Pane
	
	// Styles
	styles *Styles
	
	// Status and notifications
	status       string
	notification string
	lastUpdate   time.Time
	
	// Context for cancellation
	ctx    context.Context
	cancel context.CancelFunc
}

// @decorator: PaneType
// @description: Enumeration of available UI panes
type PaneType int

const (
	DashboardPane PaneType = iota
	APIEndpointsPane
	IngestionJobsPane
	SchedulerPane
	AIAssistantPane
	ConfigurationPane
	LogsPane
	HelpPane
)

// @decorator: Pane
// @description: Interface for UI panes
type Pane interface {
	Init() tea.Cmd
	Update(tea.Msg) (Pane, tea.Cmd)
	View() string
	Title() string
	Help() []KeyBinding
	SetSize(width, height int)
	Focus()
	Blur()
	IsFocused() bool
}

// @decorator: KeyBinding
// @description: Keyboard shortcut definition
type KeyBinding struct {
	Key         string
	Description string
}

// @decorator: Styles
// @description: Lipgloss styles for the TUI
type Styles struct {
	// Layout styles
	Container    lipgloss.Style
	Header       lipgloss.Style
	Footer       lipgloss.Style
	Sidebar      lipgloss.Style
	Content      lipgloss.Style
	
	// Component styles
	Title        lipgloss.Style
	Subtitle     lipgloss.Style
	Button       lipgloss.Style
	ActiveButton lipgloss.Style
	Border       lipgloss.Style
	
	// Status styles
	StatusBar    lipgloss.Style
	Success      lipgloss.Style
	Warning      lipgloss.Style
	Error        lipgloss.Style
	Info         lipgloss.Style
	
	// Colors
	Primary      lipgloss.Color
	Secondary    lipgloss.Color
	Background   lipgloss.Color
	Foreground   lipgloss.Color
	BorderColor  lipgloss.Color
}

// @decorator: New
// @description: Create a new TUI instance
// @param app: Application instance
// @param config: UI configuration
// @return *TUI: New TUI instance
// @return error: Any error that occurred during creation
func New(app *app.Application, config config.UIConfig) (*TUI, error) {
	if app == nil {
		return nil, errors.New("application instance cannot be nil")
	}
	
	// Create context for cancellation
	ctx, cancel := context.WithCancel(context.Background())
	
	// Initialize styles
	styles := initializeStyles(&config)
	
	// Create TUI instance
	tui := &TUI{
		app:        app,
		config:     &config,
		logger:     app.Logger(),
		activePane: DashboardPane,
		panes:      make(map[PaneType]Pane),
		styles:     styles,
		ctx:        ctx,
		cancel:     cancel,
		lastUpdate: time.Now(),
	}
	
	// Initialize panes
	if err := tui.initializePanes(); err != nil {
		cancel()
		return nil, errors.Wrap(err, "failed to initialize panes")
	}
	
	return tui, nil
}

// @decorator: Run
// @description: Run the TUI application
// @param ctx: Context for cancellation
// @return error: Any error that occurred during execution
func (t *TUI) Run(ctx context.Context) error {
	// Create Bubble Tea program
	program := tea.NewProgram(
		t,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
		tea.WithContext(ctx),
	)
	
	// Set up graceful shutdown
	go func() {
		<-ctx.Done()
		program.Quit()
	}()
	
	// Run the program
	_, err := program.Run()
	if err != nil {
		return errors.Wrap(err, "TUI program execution failed")
	}
	
	return nil
}

// @decorator: Init
// @description: Initialize the TUI (Bubble Tea interface)
// @return tea.Cmd: Initial command
func (t *TUI) Init() tea.Cmd {
	t.logger.Info("Initializing TUI...")
	
	// Initialize all panes
	var cmds []tea.Cmd
	for _, pane := range t.panes {
		cmds = append(cmds, pane.Init())
	}
	
	// Add periodic update command
	cmds = append(cmds, t.tickCmd())
	
	return tea.Batch(cmds...)
}

// @decorator: Update
// @description: Handle TUI updates (Bubble Tea interface)
// @param msg: Bubble Tea message
// @return tea.Model: Updated model
// @return tea.Cmd: Command to execute
func (t *TUI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		t.width = msg.Width
		t.height = msg.Height
		t.updatePaneSizes()
		
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			t.logger.Info("Shutting down TUI...")
			t.cancel()
			return t, tea.Quit
			
		case "tab":
			t.nextPane()
			
		case "shift+tab":
			t.prevPane()
			
		case "f1":
			t.activePane = HelpPane
			
		case "f2":
			t.activePane = DashboardPane
			
		case "f3":
			t.activePane = APIEndpointsPane
			
		case "f4":
			t.activePane = IngestionJobsPane
			
		case "f5":
			t.activePane = SchedulerPane
			
		case "f6":
			t.activePane = AIAssistantPane
			
		case "f7":
			t.activePane = ConfigurationPane
			
		case "f8":
			t.activePane = LogsPane
			
		default:
			// Pass key to active pane
			if pane, exists := t.panes[t.activePane]; exists {
				updatedPane, cmd := pane.Update(msg)
				t.panes[t.activePane] = updatedPane
				cmds = append(cmds, cmd)
			}
		}
		
		// Update focus
		t.updatePaneFocus()
		
	case TickMsg:
		// Periodic update
		t.lastUpdate = time.Now()
		t.updateStatus()
		cmds = append(cmds, t.tickCmd())
		
	default:
		// Pass message to active pane
		if pane, exists := t.panes[t.activePane]; exists {
			updatedPane, cmd := pane.Update(msg)
			t.panes[t.activePane] = updatedPane
			cmds = append(cmds, cmd)
		}
	}
	
	return t, tea.Batch(cmds...)
}

// @decorator: View
// @description: Render the TUI (Bubble Tea interface)
// @return string: Rendered view
func (t *TUI) View() string {
	if t.width == 0 || t.height == 0 {
		return "Initializing..."
	}
	
	// Build the layout
	header := t.renderHeader()
	sidebar := t.renderSidebar()
	content := t.renderContent()
	footer := t.renderFooter()
	
	// Calculate dimensions
	headerHeight := lipgloss.Height(header)
	footerHeight := lipgloss.Height(footer)
	contentHeight := t.height - headerHeight - footerHeight
	
	sidebarWidth := 25
	contentWidth := t.width - sidebarWidth
	
	// Style the components
	styledSidebar := t.styles.Sidebar.
		Width(sidebarWidth).
		Height(contentHeight).
		Render(sidebar)
	
	styledContent := t.styles.Content.
		Width(contentWidth).
		Height(contentHeight).
		Render(content)
	
	// Combine layout
	body := lipgloss.JoinHorizontal(lipgloss.Top, styledSidebar, styledContent)
	
	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		body,
		footer,
	)
}

// @decorator: renderHeader
// @description: Render the header section
// @return string: Rendered header
func (t *TUI) renderHeader() string {
	title := t.styles.Title.Render("📊 Data Ingestion TUI")
	subtitle := t.styles.Subtitle.Render("Comprehensive Data Management Interface")
	
	// Status indicator
	statusColor := t.styles.Success
	statusText := "●"
	if t.app.HasErrors() {
		statusColor = t.styles.Error
	} else if t.app.HasWarnings() {
		statusColor = t.styles.Warning
	}
	
	status := statusColor.Render(statusText + " " + t.status)
	
	// Time display
	timeStr := t.lastUpdate.Format("15:04:05")
	
	header := lipgloss.JoinHorizontal(lipgloss.Center,
		title,
		lipgloss.NewStyle().Width(10).Render(""),
		subtitle,
		lipgloss.NewStyle().Flex(1).Render(""),
		status,
		lipgloss.NewStyle().Width(5).Render(""),
		timeStr,
	)
	
	return t.styles.Header.Width(t.width).Render(header)
}

// @decorator: renderSidebar
// @description: Render the sidebar navigation
// @return string: Rendered sidebar
func (t *TUI) renderSidebar() string {
	var items []string
	
	paneNames := map[PaneType]string{
		DashboardPane:     "📊 Dashboard",
		APIEndpointsPane:  "🔗 API Endpoints",
		IngestionJobsPane: "⚡ Ingestion Jobs",
		SchedulerPane:     "⏰ Scheduler",
		AIAssistantPane:   "🤖 AI Assistant",
		ConfigurationPane: "⚙️  Configuration",
		LogsPane:          "📝 Logs",
		HelpPane:          "❓ Help",
	}
	
	for paneType := DashboardPane; paneType <= HelpPane; paneType++ {
		name := paneNames[paneType]
		if paneType == t.activePane {
			items = append(items, t.styles.ActiveButton.Render("▶ "+name))
		} else {
			items = append(items, t.styles.Button.Render("  "+name))
		}
	}
	
	// Add keyboard shortcuts
	items = append(items, "")
	items = append(items, t.styles.Subtitle.Render("Shortcuts:"))
	items = append(items, "Tab - Next pane")
	items = append(items, "F1-F8 - Direct navigation")
	items = append(items, "Ctrl+C/Q - Quit")
	
	return strings.Join(items, "\n")
}

// @decorator: renderContent
// @description: Render the main content area
// @return string: Rendered content
func (t *TUI) renderContent() string {
	if pane, exists := t.panes[t.activePane]; exists {
		return pane.View()
	}
	
	return t.styles.Error.Render("Pane not found")
}

// @decorator: renderFooter
// @description: Render the footer section
// @return string: Rendered footer
func (t *TUI) renderFooter() string {
	// Get help for active pane
	var helpItems []string
	if pane, exists := t.panes[t.activePane]; exists {
		for _, binding := range pane.Help() {
			helpItems = append(helpItems, fmt.Sprintf("%s: %s", binding.Key, binding.Description))
		}
	}
	
	helpText := strings.Join(helpItems, " • ")
	if len(helpText) > t.width-10 {
		helpText = helpText[:t.width-13] + "..."
	}
	
	footer := t.styles.StatusBar.Width(t.width).Render(helpText)
	
	// Add notification if present
	if t.notification != "" {
		notification := t.styles.Info.Render("💡 " + t.notification)
		footer = lipgloss.JoinVertical(lipgloss.Left, notification, footer)
	}
	
	return footer
}

// @decorator: initializePanes
// @description: Initialize all UI panes
// @return error: Any error that occurred during initialization
func (t *TUI) initializePanes() error {
	// Initialize dashboard pane
	dashboard, err := NewDashboardPane(t.app, t.styles)
	if err != nil {
		return errors.Wrap(err, "failed to create dashboard pane")
	}
	t.panes[DashboardPane] = dashboard
	
	// Initialize API endpoints pane
	apiPane, err := NewAPIEndpointsPane(t.app, t.styles)
	if err != nil {
		return errors.Wrap(err, "failed to create API endpoints pane")
	}
	t.panes[APIEndpointsPane] = apiPane
	
	// Initialize ingestion jobs pane
	jobsPane, err := NewIngestionJobsPane(t.app, t.styles)
	if err != nil {
		return errors.Wrap(err, "failed to create ingestion jobs pane")
	}
	t.panes[IngestionJobsPane] = jobsPane
	
	// Initialize scheduler pane
	schedulerPane, err := NewSchedulerPane(t.app, t.styles)
	if err != nil {
		return errors.Wrap(err, "failed to create scheduler pane")
	}
	t.panes[SchedulerPane] = schedulerPane
	
	// Initialize AI assistant pane
	aiPane, err := NewAIAssistantPane(t.app, t.styles)
	if err != nil {
		return errors.Wrap(err, "failed to create AI assistant pane")
	}
	t.panes[AIAssistantPane] = aiPane
	
	// Initialize configuration pane
	configPane, err := NewConfigurationPane(t.app, t.styles)
	if err != nil {
		return errors.Wrap(err, "failed to create configuration pane")
	}
	t.panes[ConfigurationPane] = configPane
	
	// Initialize logs pane
	logsPane, err := NewLogsPane(t.app, t.styles)
	if err != nil {
		return errors.Wrap(err, "failed to create logs pane")
	}
	t.panes[LogsPane] = logsPane
	
	// Initialize help pane
	helpPane, err := NewHelpPane(t.app, t.styles)
	if err != nil {
		return errors.Wrap(err, "failed to create help pane")
	}
	t.panes[HelpPane] = helpPane
	
	return nil
}

// @decorator: initializeStyles
// @description: Initialize Lipgloss styles
// @param config: UI configuration
// @return *Styles: Initialized styles
func initializeStyles(config *config.UIConfig) *Styles {
	// Parse colors
	primary := lipgloss.Color(config.Colors.Primary)
	secondary := lipgloss.Color(config.Colors.Secondary)
	background := lipgloss.Color(config.Colors.Background)
	foreground := lipgloss.Color(config.Colors.Foreground)
	borderColor := lipgloss.Color(config.Colors.Border)
	
	return &Styles{
		// Layout styles
		Container: lipgloss.NewStyle().
			Padding(1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(borderColor),
		
		Header: lipgloss.NewStyle().
			Padding(0, 1).
			Background(primary).
			Foreground(foreground).
			Bold(true),
		
		Footer: lipgloss.NewStyle().
			Padding(0, 1).
			Background(secondary).
			Foreground(foreground),
		
		Sidebar: lipgloss.NewStyle().
			Padding(1).
			Border(lipgloss.NormalBorder(), false, true, false, false).
			BorderForeground(borderColor),
		
		Content: lipgloss.NewStyle().
			Padding(1),
		
		// Component styles
		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(foreground),
		
		Subtitle: lipgloss.NewStyle().
			Italic(true).
			Foreground(secondary),
		
		Button: lipgloss.NewStyle().
			Padding(0, 1).
			Margin(0, 0, 1, 0),
		
		ActiveButton: lipgloss.NewStyle().
			Padding(0, 1).
			Margin(0, 0, 1, 0).
			Background(primary).
			Foreground(foreground).
			Bold(true),
		
		Border: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(borderColor),
		
		// Status styles
		StatusBar: lipgloss.NewStyle().
			Padding(0, 1).
			Background(background).
			Foreground(foreground),
		
		Success: lipgloss.NewStyle().
			Foreground(lipgloss.Color(config.Colors.Success)),
		
		Warning: lipgloss.NewStyle().
			Foreground(lipgloss.Color(config.Colors.Warning)),
		
		Error: lipgloss.NewStyle().
			Foreground(lipgloss.Color(config.Colors.Error)),
		
		Info: lipgloss.NewStyle().
			Foreground(secondary),
		
		// Colors
		Primary:     primary,
		Secondary:   secondary,
		Background:  background,
		Foreground:  foreground,
		BorderColor: borderColor,
	}
}

// @decorator: updatePaneSizes
// @description: Update sizes for all panes
func (t *TUI) updatePaneSizes() {
	sidebarWidth := 25
	contentWidth := t.width - sidebarWidth
	contentHeight := t.height - 4 // Account for header and footer
	
	for _, pane := range t.panes {
		pane.SetSize(contentWidth, contentHeight)
	}
}

// @decorator: updatePaneFocus
// @description: Update focus state for all panes
func (t *TUI) updatePaneFocus() {
	for paneType, pane := range t.panes {
		if paneType == t.activePane {
			pane.Focus()
		} else {
			pane.Blur()
		}
	}
}

// @decorator: nextPane
// @description: Navigate to the next pane
func (t *TUI) nextPane() {
	t.activePane = (t.activePane + 1) % (HelpPane + 1)
	t.updatePaneFocus()
}

// @decorator: prevPane
// @description: Navigate to the previous pane
func (t *TUI) prevPane() {
	if t.activePane == 0 {
		t.activePane = HelpPane
	} else {
		t.activePane--
	}
	t.updatePaneFocus()
}

// @decorator: updateStatus
// @description: Update the status display
func (t *TUI) updateStatus() {
	// Get application status
	status := t.app.GetStatus()
	t.status = fmt.Sprintf("%s | Jobs: %d | APIs: %d",
		status.State,
		status.ActiveJobs,
		status.ConnectedAPIs,
	)
}

// @decorator: TickMsg
// @description: Periodic update message
type TickMsg time.Time

// @decorator: tickCmd
// @description: Create a periodic tick command
// @return tea.Cmd: Tick command
func (t *TUI) tickCmd() tea.Cmd {
	return tea.Tick(t.config.RefreshInterval, func(time.Time) tea.Msg {
		return TickMsg{}
	})
}

// @decorator: SetNotification
// @description: Set a notification message
// @param message: Notification message
func (t *TUI) SetNotification(message string) {
	t.notification = message
	
	// Clear notification after 5 seconds
	go func() {
		time.Sleep(5 * time.Second)
		t.notification = ""
	}()
}
