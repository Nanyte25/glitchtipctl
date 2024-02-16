package organization

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var orgName string

// createOrganizationCmd represents the createOrganization command
var CreateOrganizationCmd = &cobra.Command{
	Use:   "createOrganization -n <name>",
	Short: "Create a new organization using the GlitchTip API",
	Long: `Create a new organization within GlitchTip. This command requires the organization name to be provided:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		createOrganization(orgName)
	},
}

func init() {
	CreateOrganizationCmd.Flags().StringVarP(&orgName, "name", "n", "", "Name of the organization to create")
	CreateOrganizationCmd.MarkFlagRequired("name")
}

func createOrganization(orgName string) {
	apiToken := os.Getenv("GLITCHTIP_API_TOKEN")
	if apiToken == "" {
		fmt.Println("Error: GLITCHTIP_API_TOKEN environment variable is not set.")
		return
	}

	url := "http://localhost:8000/api/0/organizations/" // Use the actual Glitchtip API endpoint

	// Define the payload for creating the organization
	payload := map[string]string{
		"name": orgName,
		"slug": generateSlug(orgName), // Assuming you have a function to generate a slug from the org name
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

	if resp.StatusCode >= 300 {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {

			fmt.Fprintf(os.Stdout, "Error reading response body: %v\n", []any{err}...)
			return
		}
		fmt.Printf("Failed to create organization, status code: %d, details: %s\n", resp.StatusCode, string(bodyBytes))
	} else {
		fmt.Println("Organization created successfully")
	}
}

// This is a placeholder for your implementation.
func generateSlug(name string) string {
	// Implement slug generation logic here. This is a simplistic approach.

	return name // Placeholder, replace with actual slug generation logic.
}
