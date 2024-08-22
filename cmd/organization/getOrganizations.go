package organization

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

// API endpoint to list organizations
const apiEndpoint = "http://localhost:8000/api/0/organizations/"

// GetOrganizationsCmd represents the getOrganizations command
var GetOrganizationsCmd = &cobra.Command{
	Use:   "getOrganizations",
	Short: "List all organizations",
	Long:  `Retrieve and display a list of all organizations from the GlitchTip API.`,
	Run: func(cmd *cobra.Command, args []string) {
		apiToken := os.Getenv("GLITCHTIP_API_TOKEN")
		if apiToken == "" {
			fmt.Println("Error: GLITCHTIP_API_TOKEN environment variable is not set.")
			return
		}
		getOrganizations(apiToken)
	},
}

func getOrganizations(apiToken string) {
	// Create an HTTP request
	req, err := http.NewRequest("GET", apiEndpoint, nil)
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
