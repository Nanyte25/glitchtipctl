package cmd

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/olekukonko/tablewriter"
)

// Shared spinner and model structure
type model struct {
	spinner  spinner.Model
	quitting bool
	apiToken string
	result   string
	err      error
	fetch    func() tea.Msg
}

func newModel(apiToken string, fetchFunc func() tea.Msg) model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	return model{spinner: s, apiToken: apiToken, fetch: fetchFunc}
}

func (m model) Init() tea.Cmd {
	// Start spinner and trigger data fetching
	return tea.Batch(m.spinner.Tick, m.fetch)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case string:
		m.quitting = true
		m.result = msg
		return m, tea.Quit
	case error:
		m.quitting = true
		m.err = msg
		return m, tea.Quit
	}
	return m, nil
}

func (m model) View() string {
	if m.quitting {
		if m.err != nil {
			return fmt.Sprintf("Error: %v\n", m.err)
		}
		return m.result
	}
	return fmt.Sprintf("\n\n   %s Fetching data...\n\n", m.spinner.View())
}

// Helper function to fetch data and format it in a table
func fetchData(url string, apiToken string, headers []string, extractData func(body []byte) [][]string) tea.Msg {
	time.Sleep(2 * time.Second) // Simulate network latency

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+apiToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP request failed with status %d", resp.StatusCode)
	}

	// Extract data from response and format it in a table
	data := extractData(body)

	var tableString bytes.Buffer
	table := tablewriter.NewWriter(&tableString)
	table.SetHeader(headers)
	for _, row := range data {
		table.Append(row)
	}
	table.Render()

	return tableString.String()
}
