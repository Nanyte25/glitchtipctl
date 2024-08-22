package organization

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/nanyte25/glitchtipctl/common" // Import the common package
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// GetMembersModel represents the spinner model struct for fetching members
type GetMembersModel struct {
	common.SpinnerModel
}

// NewGetMembersModel creates a new GetMembersModel instance
func NewGetMembersModel(apiToken string, orgSlug string) GetMembersModel {
	s := common.NewSpinnerModel(apiToken, orgSlug)
	return GetMembersModel{SpinnerModel: s}
}

// Init initializes the spinner and fetchMembers concurrently
func (m GetMembersModel) Init() tea.Cmd {
	return tea.Batch(m.SpinnerModel.Spinner.Tick, fetchMembers(m.SpinnerModel.ApiToken, m.SpinnerModel.OrgSlug))
}

// Update handles spinner ticks and results
func (m GetMembersModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m.SpinnerModel.Update(msg)
}

// View displays the spinner or result
func (m GetMembersModel) View() string {
	return m.SpinnerModel.View()
}

// GetMembersCmd represents the getMembers command
var GetMembersCmd = &cobra.Command{
	Use:   "getMembers [organization_slug]",
	Short: "Fetch the members of an organization by organizational slug",
	Long:  `Fetch and display the members of a specified organization by passing its slug.`,
	Args:  cobra.ExactArgs(1), // Ensure exactly one argument is passed (the org slug)
	Run: func(cmd *cobra.Command, args []string) {
		apiToken := os.Getenv("GLITCHTIP_API_TOKEN")
		if apiToken == "" {
			fmt.Println("Error: GLITCHTIP_API_TOKEN environment variable is not set.")
			return
		}

		// Assign the passed slug to orgSlug
		orgSlug := args[0]

		// Start the spinner model using the NewSpinnerModel from the common package
		model := NewGetMembersModel(apiToken, orgSlug)
		program := tea.NewProgram(model)

		// Run the program
		if err := program.Start(); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	},
}

// AddGetMembersCmd adds the getMembers command to the root command
func AddGetMembersCmd(rootCmd *cobra.Command) {
	rootCmd.AddCommand(GetMembersCmd)
}

// fetchMembers fetches members of the given organization using the API token
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
