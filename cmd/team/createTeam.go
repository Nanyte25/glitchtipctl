package team

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var orgName string
var teamName string

// CreateTeamCmd represents the createTeam command
var CreateTeamCmd = &cobra.Command{
	Use:   "createTeam",
	Short: "Create a new team within an organization using the GlitchTip API",
	Long: `Create a new team within a specified organization. This command requires both the organization name 
and the team name.`,
	Run: func(cmd *cobra.Command, args []string) {
		createTeam(orgName, teamName)
	},
}

func init() {
	CreateTeamCmd.Flags().StringVarP(&orgName, "org", "o", "", "Name of the organization")
	CreateTeamCmd.Flags().StringVarP(&teamName, "name", "n", "", "Name of the team to create")
	CreateTeamCmd.MarkFlagRequired("org")
	CreateTeamCmd.MarkFlagRequired("name")
}

func createTeam(orgName, teamName string) {
	apiToken := os.Getenv("GLITCHTIP_API_TOKEN")
	if apiToken == "" {
		fmt.Println("Error: GLITCHTIP_API_TOKEN environment variable is not set.")
		return
	}

	// Replace "https://localhost:8000" with the actual Glitchtip API endpoint if not using localhost
	url := fmt.Sprintf("https://localhost:8000/api/0/organizations/%s/teams/", orgName)

	// Including both 'name' and 'slug' in the payload
	payload := map[string]string{
		"name": teamName,
		"slug": generateSlug(teamName), // This function should generate a slug from the team name
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Error marshaling payload: %v\n", err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}

	req.Header.Add("Authorization", "Bearer "+apiToken)
	req.Header.Add("Content-Type", "application/json") // Set Content-Type header

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusCreated { // Assuming 201 is the success status code
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Error reading response body: %v\n", err)
			return
		}
		fmt.Printf("Failed to create team, status code: %d, details: %s\n", resp.StatusCode, string(bodyBytes))
	} else {
		fmt.Println("Team created successfully")
	}
}

// GenerateSlug is a placeholder for your slug generation logic.
// You need to implement this function based on your requirements or Glitchtip's slug rules.
func generateSlug(name string) string {
	// Implement slug generation logic here. This is a simplistic approach.
	// For real use, consider edge cases and ensure uniqueness within the organization.
	return strings.ToLower(strings.ReplaceAll(name, " ", "-"))
}
