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

func getAreaAutocompleteList(cmd *cobra.Command, args []string, toComplete string, store *shared.MarkdownStore) ([]string, cobra.ShellCompDirective) {
	// Retrieve the list of strings from your store and filter based on the substring toComplete.
	if len(args) != 0 {
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}
	options := make([]string, 0)
	for _, a := range store.Areas {
		if len(toComplete) < len(a.Title) && shared.StartsWithIgnoreCase(a.Title, toComplete) {
			// options = append(options, t.Title)
			options = append(options, strings.ReplaceAll(a.Title, " ", "\\ "))
			// options = append(options, strings.Split(t.Title, " ")[0])
		}
	}
	// to debug args len
	// options = append(options, fmt.Sprintf("%d", len(args)))
	return options, cobra.ShellCompDirectiveNoFileComp
}

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List available areas",
	Long:  ``,
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getAreaAutocompleteList(cmd, args, toComplete, &store)
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			shared.Ls(cmd, args, &store)
		} else {
			err, area := store.FindArea(strings.Join(args, " "))
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			fmt.Println(store.AreaAsMDString(area))
		}
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
}
