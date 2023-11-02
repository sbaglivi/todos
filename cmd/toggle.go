/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/sbaglivi/todos/shared"
	"github.com/spf13/cobra"
)

func getAutocompleteList(cmd *cobra.Command, args []string, toComplete string, store *shared.MarkdownStore) ([]string, cobra.ShellCompDirective) {
	// Retrieve the list of strings from your store and filter based on the substring toComplete.
	options := make([]string, 0)
	for _, a := range store.Areas {
		for _, t := range a.Tasks {
			// if toComplete == "" || (len(toComplete) <= len(item) && item[:len(toComplete)] == toComplete) {
			if len(toComplete) < len(t.Title) && shared.StartsWithIgnoreCase(t.Title, toComplete) {
				options = append(options, t.Title)
			}
		}
	}
	return options, cobra.ShellCompDirectiveNoFileComp
}

// toggleCmd represents the toggle command
var toggleCmd = &cobra.Command{
	Use:   "toggle",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			panic("no arguments given in call to toggle")
		}
		if len(args) == 1 {
			store.ToggleTask(args[0])
		}
		store.Print()
		store.Save()
		fmt.Println("toggle called")
	},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getAutocompleteList(cmd, args, toComplete, &store)
	},
}

func init() {
	rootCmd.AddCommand(toggleCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// toggleCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// toggleCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
