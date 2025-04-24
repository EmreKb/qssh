package ui

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/EmreKb/qssh/pkg/config"
)

var (
	titleStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
	selectedItemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("39")).Bold(true)
	itemStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("246"))
	infoStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Italic(true)
)

type model struct {
	hosts    []config.SSHHost
	cursor   int
	selected int
	err      error
}

func initialModel() (model, error) {
	hosts, err := config.GetSSHHosts()
	if err != nil {
		return model{}, err
	}
	return model{
		hosts:    hosts,
		cursor:   0,
		selected: -1,
	}, nil
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.hosts)-1 {
				m.cursor++
			}

		case "enter", " ":
			if m.cursor >= 0 && m.cursor < len(m.hosts) {
				m.selected = m.cursor
				return m, connectToSSH(m.hosts[m.cursor])
			}
		}

	case sshConnectionMsg:
		// When we receive the SSH connection message, quit the program
		// The actual SSH connection will happen after the TUI closes
		return m, tea.Quit
	}

	return m, nil
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n\nPress any key to exit.", m.err)
	}

	if len(m.hosts) == 0 {
		return "No SSH hosts found in your config.\n\nPress q to quit."
	}

	s := titleStyle.Render("QSSH - SSH Host Selector") + "\n\n"

	for i, host := range m.hosts {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
			s += selectedItemStyle.Render(fmt.Sprintf("%s %s", cursor, host.Name))
		} else {
			s += itemStyle.Render(fmt.Sprintf("%s %s", cursor, host.Name))
		}

		// Display user and hostname if available
		details := []string{}
		if host.User != "" {
			details = append(details, host.User+"@")
		}
		if host.Hostname != "" {
			details = append(details, host.Hostname)
		}
		if host.Port != "22" {
			details = append(details, "port "+host.Port)
		}

		if len(details) > 0 {
			s += infoStyle.Render(" (" + strings.Join(details, ", ") + ")")
		}
		s += "\n"
	}

	s += "\n" + infoStyle.Render("↑/↓: Navigate • Enter: Connect • q: Quit")

	return s
}

func connectToSSH(host config.SSHHost) tea.Cmd {
	return func() tea.Msg {
		// Return a message that will signal the program to quit
		// and then connect to SSH
		return sshConnectionMsg{host: host}
	}
}

// Custom message type to signal SSH connection
type sshConnectionMsg struct {
	host config.SSHHost
}

// Start initializes and runs the TUI
func Start() error {
	m, err := initialModel()
	if err != nil {
		return err
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		return err
	}

	// Check if we need to connect to SSH
	if m, ok := finalModel.(model); ok && m.selected >= 0 && m.selected < len(m.hosts) {
		// Now that the TUI is closed, execute SSH in the normal terminal
		host := m.hosts[m.selected]

		// Give terminal time to reset
		time.Sleep(100 * time.Millisecond)

		// Construct the SSH command
		sshArgs := []string{}

		// Add user if specified
		if host.User != "" {
			sshArgs = append(sshArgs, "-l", host.User)
		}

		// Add port if not default
		if host.Port != "22" {
			sshArgs = append(sshArgs, "-p", host.Port)
		}

		// Add host name
		sshArgs = append(sshArgs, host.Name)

		// Create shell command to execute SSH through the shell
		shell := os.Getenv("SHELL")
		if shell == "" {
			shell = "sh"
		}
		sshCmd := "ssh " + strings.Join(sshArgs, " ")
		cmd := exec.Command(shell, "-ic", sshCmd)

		// Prepare environment without problematic DYLD_* vars (macOS SIP)
		env := []string{}
		for _, v := range os.Environ() {
			if strings.HasPrefix(v, "DYLD_") {
				continue
			}
			env = append(env, v)
		}
		cmd.Env = env

		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		// Execute the shell command
		return cmd.Run()
	}

	return nil
}
