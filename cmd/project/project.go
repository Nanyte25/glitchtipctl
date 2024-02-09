package project

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// GetProjectsCmd represents the getProjects command
var GetProjectsCmd = &cobra.Command{
	Use:   "getProjects",
	Short: "Get a list of projects from your organization",
	Long: `Get a list of projects from your organization. This command makes an HTTP GET request to the GlitchTip API
and prints out the list of projects.`,
	Run: func(cmd *cobra.Command, args []string) {
		apiToken := os.Getenv("GLITCHTIP_API_TOKEN")
		if apiToken == "" {
			fmt.Println("Error: GLITCHTIP_API_TOKEN environment variable is not set.")
			return
		}

		url := "http://localhost:8000/api/0/projects/"
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Printf("Error creating request: %v\n", err)
			return
		}

		req.Header.Add("Authorization", "Bearer "+apiToken)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Error making request: %v\n", err)
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Error reading response body: %v\n", err)
			return
		}

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Error: HTTP request failed with status %d\n", resp.StatusCode)
			return
		}

		var projects []map[string]interface{}
		if err := json.Unmarshal(body, &projects); err != nil {
			fmt.Printf("Error unmarshaling response: %v\n", err)
			return
		}

		// Initialize table writer
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Name", "Slug"})

		for _, project := range projects {
			id := fmt.Sprintf("%v", project["id"])
			name := fmt.Sprintf("%v", project["name"])
			slug := fmt.Sprintf("%v", project["slug"])
			table.Append([]string{id, name, slug})
		}

		fmt.Println("Projects:")
		table.Render() // Render the table to standard output
	},
}

func init() {
	// your init code here, if any
}
