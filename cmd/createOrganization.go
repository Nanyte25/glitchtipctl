/*
Mark Freer GlitchtipCTL Author 2024
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// createOrganizationCmd represents the createOrganization command
var createOrganizationCmd = &cobra.Command{
	Use:   "createOrganization",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("createOrganization called")
	},
}

func init() {
	rootCmd.AddCommand(createOrganizationCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createOrganizationCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createOrganizationCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
