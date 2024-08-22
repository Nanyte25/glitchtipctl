package organization

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
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

// Print organizations in a table format
func printOrganizationsTable(organizations []Organization) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Name", "Slug", "Status", "Created", "2FA Required"})

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

	table.Render()
}
