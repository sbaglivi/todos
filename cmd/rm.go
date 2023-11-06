/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/sbaglivi/todos/shared"
	"github.com/spf13/cobra"
)

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("rm expects 1 argument. received ", len(args))
			os.Exit(1)
		}
		re := regexp.MustCompile(`([^_])_([^_])`)
		output := re.ReplaceAllString(args[0], "$1 $2")

		// Replace all underscores prefixed by a backslash with an underscore
		re = regexp.MustCompile(`__`)
		output = re.ReplaceAllString(output, "_")
		parts := strings.Split(output, "/")
		forced, err := cmd.Flags().GetBool("force")
		if err != nil {
			panic(err)
		}
		if len(parts) > 3 {
			fmt.Println("/ is a reserved character, use it only as a separator between areas and tasks or tasks and subtasks")
			os.Exit(1)
		} else if len(parts) == 1 {
			areaName := parts[0]
			if !forced {
				var force string
				fmt.Println("Are you sure you want to delete an entire area? Enter 'yes' or 'y' to confirm")
				fmt.Scanln(&force)
				force = strings.ToLower(strings.TrimSpace(force))
				if !(force == "y" || force == "yes") {
					fmt.Println("did not receive confirmation. aborting")
					os.Exit(1)
				}
			}
			if err := store.DeleteArea(areaName); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else if len(parts) == 2 {
			areaName := parts[0]
			taskName := parts[1]
			if !forced {
				var force string
				fmt.Println("Are you sure you want to delete an entire task? Enter 'yes' or 'y' to confirm")
				fmt.Scanln(&force)
				force = strings.ToLower(strings.TrimSpace(force))
				if !(force == "y" || force == "yes") {
					fmt.Println("did not receive confirmation. aborting")
					os.Exit(1)
				}
			}
			if err := store.DeleteTask(areaName, taskName); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else {
			areaName := parts[0]
			taskName := parts[1]
			subtaskName := parts[2]
			if err := store.DeleteSubtask(areaName, taskName, subtaskName); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
		store.Save()

	},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		toComplete = strings.ReplaceAll(toComplete, "_", " ")
		possibleCompletions := make([]string, 0)
		for i := range store.Areas {
			area := &store.Areas[i]
			if shared.LongestContainsOtherIgnoreCase(area.Title, toComplete) {
				if len(toComplete) <= len(area.Title) {
					possibleCompletions = append(possibleCompletions, strings.ReplaceAll(area.Title, " ", "_"))
				}
				for j := range area.Tasks {
					task := &area.Tasks[j]
					rawFullTitle := area.Title + "/" + task.Title
					if shared.LongestContainsOtherIgnoreCase(rawFullTitle, toComplete) {
						if len(toComplete) < len(rawFullTitle) {
							possibleCompletions = append(possibleCompletions, strings.ReplaceAll(rawFullTitle, " ", "_"))
						}
						for k := range task.Subtasks {
							subtask := &task.Subtasks[k]
							rawFullSubtitle := area.Title + "/" + task.Title + "/" + subtask.Title
							if shared.StartsWithIgnoreCase(rawFullSubtitle, toComplete) {
								possibleCompletions = append(possibleCompletions, strings.ReplaceAll(rawFullSubtitle, " ", "_"))
							}
						}
					}
				}
			}
		}
		return possibleCompletions, cobra.ShellCompDirectiveNoSpace
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// rmCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	rmCmd.Flags().BoolP("force", "f", false, "force removal of whole areas / tasks")
	// rmCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
