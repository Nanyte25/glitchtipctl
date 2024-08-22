package cmd

import (
	"fmt"
	"os"

	"github.com/nanyte25/glitchtipctl/cmd/organization"
	"github.com/nanyte25/glitchtipctl/cmd/project"
	"github.com/nanyte25/glitchtipctl/cmd/team"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "glitchtipctl",
	Short: "glitchtipctl is a commandline tool for Glitchtip error tracking software.",
	Long: `glitchtipctl is a commandline tool for interacting with the GlitchTip error tracking software.

This tool provides various commands for managing organizations, projects, teams, users, and more. 
Use this CLI to automate and manage tasks within your GlitchTip account.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// RootCmd exposes the root command to other packages
func RootCmd() *cobra.Command {
	return rootCmd
}

func init() {
	// Add your commands here. These commands are added as subcommands of the root command.
	rootCmd.AddCommand(project.GetProjectsCmd)
	rootCmd.AddCommand(team.GetTeamsCmd)
	rootCmd.AddCommand(team.CreateTeamCmd)
	rootCmd.AddCommand(organization.CreateOrganizationCmd)
	rootCmd.AddCommand(organization.GetOrganizationsCmd)

	// Register the GetMembersCmd from the organization package
	organization.AddGetMembersCmd(rootCmd) // Add getMembers command to root

	// Additional commands can be added here.
	rootCmd.Flags().BoolP("toggle", "t", false, "To toggle the debug mode")
}
