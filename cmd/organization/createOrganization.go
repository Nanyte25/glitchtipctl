package organization

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

// createOrganizationCmd represents the createOrganization command
var CreateOrganizationCmd = &cobra.Command{
	Use:   "createOrganization -n <name>",
	Short: "Create a new organization using the GlitchTip API",
	Long: `Create a new organization within GlitchTip. This command requires the organization name to be provided:

Example usage:
  glitchtipctl createOrganization -n "MyOrganization"
`,
	Run: func(cmd *cobra.Command, args []string) {
		createOrganization(orgName)
	},
}

func init() {
	CreateOrganizationCmd.Flags().StringVarP(&orgName, "name", "n", "", "Name of the organization to create")
	CreateOrganizationCmd.MarkFlagRequired("name")
}

// Create the organization and print the updated list of organizations
func createOrganization(orgName string) {
	apiToken := os.Getenv("GLITCHTIP_API_TOKEN")
	if apiToken == "" {
		fmt.Println("Error: GLITCHTIP_API_TOKEN environment variable is not set.")
		return
	}

	url := "http://localhost:8000/api/0/organizations/" // Replace with actual Glitchtip API endpoint

	// Define the payload for creating the organization
	payload := map[string]string{
		"name": orgName,
		"slug": generateSlug(orgName),
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

	// Handle non-200 status codes
	if resp.StatusCode >= 300 {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Error reading response body: %v\n", err)
			return
		}
		fmt.Printf("Failed to create organization, status code: %d, details: %s\n", resp.StatusCode, string(bodyBytes))
	} else {
		fmt.Println("Organization created successfully")
		// Fetch and print the updated list of organizations
		getAndPrintOrganizations(apiToken)
	}
}

// Generate slug from the organization name (simple version)
func generateSlug(name string) string {
	return strings.ToLower(strings.ReplaceAll(name, " ", "-"))
}

// Fetch the updated list of organizations and print them
func getAndPrintOrganizations(apiToken string) {
	url := "http://localhost:8000/api/0/organizations/" // Replace with the actual API endpoint

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}

	req.Header.Add("Authorization", "Bearer "+apiToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Error fetching organizations: %v\n", err)
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

	var organizations []Organization
	err = json.Unmarshal(body, &organizations)
	if err != nil {
		fmt.Printf("Error parsing JSON response: %v\n", err)
		return
	}

	printOrganizationsTable(organizations)
}
