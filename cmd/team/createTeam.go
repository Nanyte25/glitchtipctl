package team

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
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
		// Load environment variables from .env file at runtime
		if err := godotenv.Load(); err != nil {
			fmt.Println("Warning: No .env file found")
		}
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

	glitchtipURL := os.Getenv("GLITCHTIP_URL") // Retrieve the Glitchtip URL from environment variable
	if glitchtipURL == "" {
		fmt.Println("Error: GLITCHTIP_URL environment variable is not set.")
		return
	}

	url := fmt.Sprintf("%s/0/organizations/%s/teams/", glitchtipURL, orgName)

	payload := map[string]string{
		"name": teamName,
		"slug": generateSlug(teamName),
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
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Handle response
	if resp.StatusCode == http.StatusCreated {
		var teamResponse struct {
			DateCreated string   `json:"dateCreated"`
			ID          string   `json:"id"`
			IsMember    bool     `json:"isMember"`
			MemberCount int      `json:"memberCount"`
			Slug        string   `json:"slug"`
			Projects    []string `json:"projects"`
		}

		err := json.NewDecoder(resp.Body).Decode(&teamResponse)
		if err != nil {
			fmt.Printf("Error decoding response body: %v\n", err)
			return
		}

		fmt.Printf("Team created successfully:\n")
		fmt.Printf("- ID: %s\n- Slug: %s\n- Date Created: %s\n- Member Count: %d\n", teamResponse.ID, teamResponse.Slug, teamResponse.DateCreated, teamResponse.MemberCount)
	} else {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Error reading response body: %v\n", err)
			return
		}
		fmt.Printf("Failed to create team, status code: %d, details: %s\n", resp.StatusCode, string(bodyBytes))
	}
}

func generateSlug(name string) string {
	// Simple slug generation, adapt as needed for your application's requirements
	return strings.ToLower(strings.ReplaceAll(name, " ", "-"))
}
