package organization

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"
)

// Assuming you have an API endpoint to list organizations
const apiEndpoint = "http://localhost:8000/api/organizations"

// Organization represents the structure of your organization data
type Organization struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	// Add other fields as per your API response
}

// GetOrganizationsCmd represents the getOrganizations command
var GetOrganizationsCmd = &cobra.Command{
	Use:   "getOrganizations",
	Short: "List all organizations",
	Long:  `Retrieve and display a list of all organizations from the GlitchTip API.`,
	Run: func(cmd *cobra.Command, args []string) {
		getOrganizations()
	},
}

func getOrganizations() {
	// Make an HTTP GET request to the API endpoint
	resp, err := http.Get(apiEndpoint)
	if err != nil {
		fmt.Printf("Error fetching organizations: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Read and parse the response body
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

	// Print the fetched organizations
	for _, org := range organizations {
		fmt.Printf("ID: %s, Name: %s\n", org.ID, org.Name)
	}
}
