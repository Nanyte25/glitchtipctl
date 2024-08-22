package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/nanyte25/glitchtipctl/common" // Updated to import the common package
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var orgSlug string

// GetUsersCmd represents the getUsers command
var GetUsersCmd = &cobra.Command{
	Use:   "getUsers [organization_slug]",
	Short: "Fetch the users of an organization",
	Long:  `Fetch and display the users of a specified organization by passing its slug.`,
	Args:  cobra.ExactArgs(1), // Ensure exactly one argument is passed (the org slug)
	Run: func(cmd *cobra.Command, args []string) {
		apiToken := os.Getenv("GLITCHTIP_API_TOKEN")
		if apiToken == "" {
			fmt.Println("Error: GLITCHTIP_API_TOKEN environment variable is not set.")
			return
		}

		// Assign the passed slug to orgSlug
		orgSlug = args[0]

		// Start the spinner model using the NewSpinnerModel from the common package
		model := common.NewSpinnerModel(apiToken, orgSlug)
		program := tea.NewProgram(model)

		// Run the program
		if err := program.Start(); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	},
}

func init() {
	// Register the GetUsersCmd
	RootCmd().AddCommand(GetUsersCmd)
}

// fetchData fetches the users of the given organization using the API token
func fetchData(apiToken, orgSlug string) func() tea.Msg {
	return func() tea.Msg {
		url := fmt.Sprintf("http://localhost:8000/api/0/organizations/%s/users/", orgSlug)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return fmt.Sprintf("Error creating request: %v", err)
		}

		req.Header.Add("Authorization", "Bearer "+apiToken)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return fmt.Sprintf("Error fetching users: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Sprintf("Error: received status code %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Sprintf("Error reading response body: %v", err)
		}

		var users []map[string]interface{}
		err = json.Unmarshal(body, &users)
		if err != nil {
			return fmt.Sprintf("Error parsing JSON response: %v", err)
		}

		// Convert the users to a formatted table
		return formatUsersTable(users)
	}
}

// formatUsersTable converts the users data into a formatted table string
func formatUsersTable(users []map[string]interface{}) string {
	headers := []string{"ID", "Name", "Email"}

	// Create a buffer to capture the table output
	var buffer bytes.Buffer
	table := tablewriter.NewWriter(&buffer)
	table.SetHeader(headers)

	for _, user := range users {
		id := fmt.Sprintf("%v", user["id"])
		name := fmt.Sprintf("%v", user["name"])
		email := fmt.Sprintf("%v", user["email"])
		table.Append([]string{id, name, email})
	}

	// Render the table into the buffer
	table.Render()

	// Return the buffer's contents as a string
	return buffer.String()
}
