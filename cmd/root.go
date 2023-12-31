/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/sbaglivi/todos/shared"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "todos",
	Short: "An application that keeps track of todos",
	Long: `Application that keeps track of your todos, organized in tasks which can be made of subtasks.
		Todos can be stored in markdown format so that you can interact with them even without the app.	
	`,
	ValidArgs: []string{"toggle", "ls"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			store.Print()
		} else if len(args) == 1 {
			store.LsArea(args[0])
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var store shared.MarkdownStore

func init() {
	store = shared.Setup()
	var areaName string
	rootCmd.PersistentFlags().StringVarP(&areaName, "area", "a", "", "Area name")
}
