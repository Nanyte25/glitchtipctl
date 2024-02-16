/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
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
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
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
	// Create Glitchtip Error tracking organization
	rootCmd.AddCommand(project.GetProjectsCmd)
	rootCmd.AddCommand(team.GetTeamsCmd)
	rootCmd.AddCommand(team.CreateTeamCmd)
	rootCmd.AddCommand(organization.CreateOrganizationCmd) // N.B. make these global package once all are completed

	//add new teams feature
	// Initialize other commands and flags
	rootCmd.Flags().BoolP("toggle", "t", false, "To toggle the debug mode")
}
