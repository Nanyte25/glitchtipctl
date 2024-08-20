package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// GetUsersCmd represents the getUsers command
var GetUsersCmd = &cobra.Command{
	Use:   "getUsers",
	Short: "Get a list of users from your organization",
	Long:  `Get a list of users from your organization by making an HTTP GET request to the GlitchTip API and printing out the list of users.`,
	Run: func(cmd *cobra.Command, args []string) {
		apiToken := os.Getenv("GLITCHTIP_API_TOKEN")
		if apiToken == "" {
			fmt.Println("Error: GLITCHTIP_API_TOKEN environment variable is not set.")
			return
		}

		// Start the spinner model
		model := newSpinnerModel(apiToken)
		program := tea.NewProgram(model)

		// Run the spinner and handle errors
		if err := program.Start(); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	},
}

// Spinner model to display the spinner while loading
type spinnerModel struct {
	spinner  spinner.Model
	quitting bool
	apiToken string
	result   string
	err      error
}

func newSpinnerModel(apiToken string) spinnerModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	return spinnerModel{spinner: s, apiToken: apiToken}
}

func (m spinnerModel) Init() tea.Cmd {
	// Start the spinner and fetch users concurrently
	return tea.Batch(m.spinner.Tick, m.fetchUsers)
}

func (m spinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m spinnerModel) View() string {
	if m.quitting {
		if m.err != nil {
			return fmt.Sprintf("Error: %v\n", m.err)
		}
		return fmt.Sprintf("%s\n", m.result)
	}
	return fmt.Sprintf("\n\n   %s Fetching users...\n\n", m.spinner.View())
}

func (m spinnerModel) fetchUsers() tea.Msg {
	// Introduce a delay to simulate network latency (for testing)
	time.Sleep(2 * time.Second) // Add artificial delay here for the spinner to display

	// Make the API call to get users
	url := "http://localhost:8000/api/0/users/"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+m.apiToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body) // Replaces ioutil.ReadAll
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP request failed with status %d", resp.StatusCode)
	}

	var users []map[string]interface{}
	if err := json.Unmarshal(body, &users); err != nil {
		return err
	}

	// Create a table of users
	var tableString bytes.Buffer
	table := tablewriter.NewWriter(&tableString)
	table.SetHeader([]string{"ID", "Name", "Email"})

	for _, user := range users {
		id := fmt.Sprintf("%v", user["id"])
		name := fmt.Sprintf("%v", user["name"])
		email := fmt.Sprintf("%v", user["email"])
		table.Append([]string{id, name, email})
	}
	table.Render()

	return tableString.String()
}

func init() {
	// Add the getUsers command to the root
	rootCmd.AddCommand(GetUsersCmd)
}
