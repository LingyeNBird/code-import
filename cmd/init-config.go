package cmd

import (
	"ccrctl/pkg/config"
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initConfig)
}

var initConfig = &cobra.Command{
	Use:   "init-config",
	Short: "generate config.yaml.default file",
	Long:  `init config  will generate config.yaml.default file, you can copy it to config.yaml and edit it`,
	Run: func(cmd *cobra.Command, args []string) {
		err := config.InitConfig()
		if err != nil {
			panic(err)
		}
		fmt.Println("generate config.yaml.default file success, you can copy it to config.yaml and edit it")
	},
}
