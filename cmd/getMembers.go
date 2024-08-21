package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

// GetMembersCmd represents the getMembers command
var GetMembersCmd = &cobra.Command{
	Use:   "getMembers",
	Short: "Get a list of members from your organization",
	Long:  `This command makes an HTTP GET request to the GlitchTip API and prints out the list of members from your organization.`,
	Run: func(cmd *cobra.Command, args []string) {
		apiToken := os.Getenv("GLITCHTIP_API_TOKEN")
		if apiToken == "" {
			fmt.Println("Error: GLITCHTIP_API_TOKEN environment variable is not set.")
			return
		}

		// Use shared spinner model and fetch function
		model := newModel(apiToken, fetchMembers(apiToken))
		program := tea.NewProgram(model)

		// Start the program
		if err := program.Start(); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	},
}

func fetchMembers(apiToken string) func() tea.Msg {
	return func() tea.Msg {
		url := "http://localhost:8000/api/0/organizations/0/members/" // Replace <org_slug> with your actual org slug
		headers := []string{"ID", "Name", "Email"}

		// Extract data function for formatting
		extractData := func(body []byte) [][]string {
			var members []map[string]interface{}
			json.Unmarshal(body, &members)

			data := [][]string{}
			for _, member := range members {
				id := fmt.Sprintf("%v", member["id"])
				name := fmt.Sprintf("%v", member["name"])
				email := fmt.Sprintf("%v", member["email"])
				data = append(data, []string{id, name, email})
			}
			return data
		}

		return fetchData(url, apiToken, headers, extractData)
	}
}

func init() {
	rootCmd.AddCommand(GetMembersCmd)
}
