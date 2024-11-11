package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "print the version number of ccrctl",
	Long:  `print the version number of ccrctl`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("===========================================================")
		log.Printf("|  Build Time  : %-40s |", BuildTime)
		log.Printf("|  Version     : %-39s |", Version)
		log.Printf("===========================================================")
	},
}
