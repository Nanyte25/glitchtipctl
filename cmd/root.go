package cmd

import (
	"os"

	"github.com/nanyte25/glitchtipctl/cmd/organization"
	"github.com/nanyte25/glitchtipctl/cmd/project"
	"github.com/nanyte25/glitchtipctl/cmd/team"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "glitchtipctl",
	Short: "glitchtipctl is a commandline tool for Glitchtip error Tracking software.",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Add your commands
	rootCmd.AddCommand(project.GetProjectsCmd)
	rootCmd.AddCommand(team.GetTeamsCmd)
	rootCmd.AddCommand(team.CreateTeamCmd)
	rootCmd.AddCommand(organization.CreateOrganizationCmd)
	rootCmd.AddCommand(organization.GetOrganizationsCmd)

	// Additional features can be added here.
	rootCmd.Flags().BoolP("toggle", "t", false, "To toggle the debug mode")
}
