/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "generates completion scripts for the app",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// switch args[0] {
		// case "bash":
		// 	cmd.Root().GenBashCompletion(os.Stdout)
		// case "zsh":
		// 	cmd.Root().GenZshCompletion(os.Stdout)
		// case "fish":
		// 	cmd.Root().GenFishCompletion(os.Stdout, true)
		// case "powershell":
		// 	cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
		// }
		if err := cmd.Root().GenBashCompletionFileV2("complete.bash", true); err != nil {
			fmt.Printf("error while generating bash completion script: %s\n", err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
