/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/sbaglivi/todos/shared"
	"github.com/spf13/cobra"
)

func getAutocompleteList(cmd *cobra.Command, args []string, toComplete string, store *shared.MarkdownStore) ([]string, cobra.ShellCompDirective) {
	// Retrieve the list of strings from your store and filter based on the substring toComplete.
	if len(args) != 0 {
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}
	options := make([]string, 0)
	for _, a := range store.Areas {
		for _, t := range a.Tasks {
			// if toComplete == "" || (len(toComplete) <= len(item) && item[:len(toComplete)] == toComplete) {
			if len(toComplete) < len(t.Title) && shared.StartsWithIgnoreCase(t.Title, toComplete) {
				// options = append(options, t.Title)
				options = append(options, strings.ReplaceAll(t.Title, " ", "\\ "))
				// options = append(options, strings.Split(t.Title, " ")[0])
			}
		}
	}
	// to debug args len
	// options = append(options, fmt.Sprintf("%d", len(args)))
	return options, cobra.ShellCompDirectiveNoFileComp
}

// toggleCmd represents the toggle command
var toggleCmd = &cobra.Command{
	Use:   "toggle",
	Short: "Toggle completion of a task",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		area, err := cmd.Flags().GetString("area")
		if err != nil {
			panic(err)
		}
		fmt.Println(area)
		if len(args) == 0 {
			fmt.Println("toggle expects 1 or more arguments to identify the task to toggle, none were given")
			os.Exit(1)
		}
		store.ToggleTask(strings.Join(args, " "))
		store.Print()
		store.Save()
	},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getAutocompleteList(cmd, args, toComplete, &store)
	},
}

func init() {
	rootCmd.AddCommand(toggleCmd)
	toggleCmd.MarkPersistentFlagRequired("area") // Make the "area" flag mandatory
}
