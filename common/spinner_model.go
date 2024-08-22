package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/olekukonko/tablewriter"
)

// SpinnerModel represents the spinner model struct
type SpinnerModel struct {
	Spinner  spinner.Model
	Quitting bool
	ApiToken string
	OrgSlug  string
	Result   string
	Err      error
}

// NewSpinnerModel creates a new SpinnerModel instance
func NewSpinnerModel(apiToken string, orgSlug string) SpinnerModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	return SpinnerModel{Spinner: s, ApiToken: apiToken, OrgSlug: orgSlug}
}

// Init initializes the spinner and fetches data concurrently
func (m SpinnerModel) Init() tea.Cmd {
	// The fetchData function should be implemented in the specific context
	// that uses this spinner model, so this is just a placeholder.
	return tea.Batch(m.Spinner.Tick)
}

// Update handles spinner ticks and results
func (m SpinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.Spinner, cmd = m.Spinner.Update(msg)
		return m, cmd
	case string:
		m.Quitting = true
		m.Result = msg
		return m, tea.Quit
	case error:
		m.Quitting = true
		m.Err = msg
		return m, tea.Quit
	}
	return m, nil
}

// View displays the spinner or the result
func (m SpinnerModel) View() string {
	if m.Quitting {
		if m.Err != nil {
			return "Error: " + m.Err.Error() + "\n"
		}
		return m.Result + "\n"
	}
	return m.Spinner.View() + " Fetching data..."
}

// Example fetch function - can be replaced based on usage context
func fetchMembers(apiToken, orgSlug string) func() tea.Msg {
	return func() tea.Msg {
		url := fmt.Sprintf("http://localhost:8000/api/0/organizations/%s/members/", orgSlug)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return fmt.Sprintf("Error creating request: %v", err)
		}

		req.Header.Add("Authorization", "Bearer "+apiToken)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return fmt.Sprintf("Error fetching members: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Sprintf("Error: received status code %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Sprintf("Error reading response body: %v", err)
		}

		var members []map[string]interface{}
		err = json.Unmarshal(body, &members)
		if err != nil {
			return fmt.Sprintf("Error parsing JSON response: %v", err)
		}

		// Convert the members to a formatted table
		return formatMembersTable(members)
	}
}

// formatMembersTable converts the members data into a formatted table string
func formatMembersTable(members []map[string]interface{}) string {
	headers := []string{"ID", "Name", "Email"}

	// Create a buffer to capture the table output
	var buffer bytes.Buffer
	table := tablewriter.NewWriter(&buffer)
	table.SetHeader(headers)

	for _, member := range members {
		id := fmt.Sprintf("%v", member["id"])
		name := fmt.Sprintf("%v", member["name"])
		email := fmt.Sprintf("%v", member["email"])
		table.Append([]string{id, name, email})
	}

	// Render the table into the buffer
	table.Render()

	// Return the buffer's contents as a string
	return buffer.String()
}
