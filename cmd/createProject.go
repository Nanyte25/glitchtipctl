package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

// Predefined list of valid platforms
var validPlatforms = []string{
	"python",
	"react",
	"django",
	"flutter",
	"react-native",
	"c",
	"javascript",
	"node",
}

// createProjectCmd represents the createProject command
var createProjectCmd = &cobra.Command{
	Use:   "createProject",
	Short: "Create a new project in GlitchTip",
	Long:  `Use this command to create a new project within a team and organization in GlitchTip by providing a name, slug, team slug, and platform.`,
	Run: func(cmd *cobra.Command, args []string) {
		apiToken := os.Getenv("GLITCHTIP_API_TOKEN")
		if apiToken == "" {
			fmt.Println("Error: GLITCHTIP_API_TOKEN environment variable is not set.")
			return
		}

		// Get flag values
		name, _ := cmd.Flags().GetString("name")
		slug, _ := cmd.Flags().GetString("slug")
		teamSlug, _ := cmd.Flags().GetString("team")
		orgSlug, _ := cmd.Flags().GetString("organization")
		platform, _ := cmd.Flags().GetString("platform")

		// Validate platform
		if !isValidPlatform(platform) {
			fmt.Printf("Error: '%s' is not a valid platform. Valid platforms are: %v\n", platform, validPlatforms)
			return
		}

		if name == "" || slug == "" || teamSlug == "" || orgSlug == "" || platform == "" {
			fmt.Println("Error: name, slug, team, organization, and platform must be provided.")
			return
		}

		// Create the project
		err := createProject(apiToken, name, slug, teamSlug, orgSlug, platform)
		if err != nil {
			fmt.Printf("Failed to create project: %v\n", err)
		} else {
			fmt.Println("Project created successfully!")
			// List the projects after creation
			listProjects(apiToken, orgSlug)
		}
	},
}

func init() {
	rootCmd.AddCommand(createProjectCmd)

	// Define the flags
	createProjectCmd.Flags().StringP("name", "n", "", "Name of the project (required)")
	createProjectCmd.Flags().StringP("slug", "s", "", "Slug for the project (required)")
	createProjectCmd.Flags().StringP("team", "t", "", "Slug of the team (required)")
	createProjectCmd.Flags().StringP("organization", "o", "", "Slug of the organization the project belongs to (required)")
	createProjectCmd.Flags().StringP("platform", "p", "", "Platform of the project e.g. python, React, Javascript, node, C#, or Flutter (required)")

	// Mark flags as required
	createProjectCmd.MarkFlagRequired("name")
	createProjectCmd.MarkFlagRequired("slug")
	createProjectCmd.MarkFlagRequired("team")
	createProjectCmd.MarkFlagRequired("organization")
	createProjectCmd.MarkFlagRequired("platform")
}

// createProject sends a POST request to the GlitchTip API to create a project
func createProject(apiToken, name, slug, teamSlug, orgSlug, platform string) error {
	url := fmt.Sprintf("http://localhost:8000/api/0/teams/%s/%s/projects/", orgSlug, teamSlug)

	// Define the payload
	payload := map[string]interface{}{
		"name":     name,
		"slug":     slug,
		"platform": platform,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshaling payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Add("Authorization", "Bearer "+apiToken)
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("error reading response body: %w", err)
		}
		return fmt.Errorf("error: received status code %d, details: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}

// isValidPlatform checks if the given platform is valid
func isValidPlatform(platform string) bool {
	for _, p := range validPlatforms {
		if p == platform {
			return true
		}
	}
	return false
}

// listProjects lists all projects for a given organization
func listProjects(apiToken, orgSlug string) {
	url := fmt.Sprintf("http://localhost:8000/api/0/organizations/%s/projects/", orgSlug)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}

	req.Header.Add("Authorization", "Bearer "+apiToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Error fetching projects: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: received status code %d\n", resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return
	}

	var projects []map[string]interface{}
	err = json.Unmarshal(body, &projects)
	if err != nil {
		fmt.Printf("Error parsing JSON response: %v\n", err)
		return
	}

	// Print the list of projects
	fmt.Println("+----+-------------+-------------+")
	fmt.Println("| ID |    NAME     |    SLUG     |")
	fmt.Println("+----+-------------+-------------+")
	for _, project := range projects {
		fmt.Printf("| %-2v | %-11v | %-11v |\n", project["id"], project["name"], project["slug"])
	}
	fmt.Println("+----+-------------+-------------+")
}
