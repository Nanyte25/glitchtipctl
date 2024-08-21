package organization

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// Organization represents the structure of your organization data
type Organization struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	Slug              string `json:"slug"`
	DateCreated       string `json:"dateCreated"`
	Status            Status `json:"status"`
	Avatar            Avatar `json:"avatar"`
	IsEarlyAdopter    bool   `json:"isEarlyAdopter"`
	Require2FA        bool   `json:"require2FA"`
	IsAcceptingEvents bool   `json:"isAcceptingEvents"`
}

// Status represents the nested status structure
type Status struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Avatar represents the nested avatar structure
type Avatar struct {
	AvatarType string      `json:"avatarType"`
	AvatarUUID interface{} `json:"avatarUuid"`
}

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

	// Add the API token in the Authorization header
	req.Header.Add("Authorization", "Bearer "+apiToken)

	// Execute the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error fetching organizations: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Check for non-200 status code
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: received status code %d\n", resp.StatusCode)
		return
	}

	// Read and parse the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return
	}

	// Unmarshal JSON into a slice of organizations
	var organizations []Organization
	err = json.Unmarshal(body, &organizations)
	if err != nil {
		fmt.Printf("Error parsing JSON response: %v\n", err)
		return
	}

	// Print the fetched organizations in a table format
	printOrganizationsTable(organizations)
}

func printOrganizationsTable(organizations []Organization) {
	// Create a new tablewriter instance
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Name", "Slug", "Status", "Created", "2FA Required"})

	// Populate the table with data
	for _, org := range organizations {
		table.Append([]string{
			fmt.Sprintf("%d", org.ID),
			org.Name,
			org.Slug,
			org.Status.Name,
			org.DateCreated,
			fmt.Sprintf("%t", org.Require2FA),
		})
	}

	// Render the table
	table.Render()
}
