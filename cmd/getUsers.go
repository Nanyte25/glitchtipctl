package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

// GetUsersCmd represents the getUsers command
var GetUsersCmd = &cobra.Command{
	Use:   "getUsers",
	Short: "Get a list of users from your organization",
	Long:  `This command makes an HTTP GET request to the GlitchTip API and prints out the list of users from your organization.`,
	Run: func(cmd *cobra.Command, args []string) {
		apiToken := os.Getenv("GLITCHTIP_API_TOKEN")
		if apiToken == "" {
			fmt.Println("Error: GLITCHTIP_API_TOKEN environment variable is not set.")
			return
		}

		// Use shared spinner model and fetch function
		model := newModel(apiToken, fetchUsers(apiToken))
		program := tea.NewProgram(model)

		// Start the program
		if err := program.Start(); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	},
}

func fetchUsers(apiToken string) func() tea.Msg {
	return func() tea.Msg {
		url := "http://localhost:8000/api/0/users/"
		headers := []string{"ID", "Name", "Email"}

		// Extract data function for formatting
		extractData := func(body []byte) [][]string {
			var users []map[string]interface{}
			json.Unmarshal(body, &users)

			data := [][]string{}
			for _, user := range users {
				id := fmt.Sprintf("%v", user["id"])
				name := fmt.Sprintf("%v", user["name"])
				email := fmt.Sprintf("%v", user["email"])
				data = append(data, []string{id, name, email})
			}
			return data
		}

		return fetchData(url, apiToken, headers, extractData)
	}
}

func init() {
	rootCmd.AddCommand(GetUsersCmd)
}
