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

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("zero or multiple arguments given to add. expected areaName[/taskName][/subtaskName]. Substitute spaces for _, if you need a _ prefix it with a \\")
			os.Exit(1)
		}
		creationSeparator := "::"
		joinedArgs := strings.Join(args, " ")
		re := regexp.MustCompile(`([^_])_([^_])`)
		output := re.ReplaceAllString(joinedArgs, "$1 $2")

		// Replace all underscores prefixed by a backslash with an underscore
		re = regexp.MustCompile(`__`)
		output = re.ReplaceAllString(output, "_")
		parts := strings.Split(output, "/")
		if len(parts) == 1 {
			newAreaName := parts[0]
			if strings.Contains(newAreaName, creationSeparator) {
				names := strings.Split(newAreaName, creationSeparator)
				for _, name := range names {
					store.Areas = append(store.Areas, shared.Area{
						Title: name,
					})
				}
			}
			store.Areas = append(store.Areas, shared.Area{
				Title: newAreaName,
			})
		} else if len(parts) == 2 {
			areaName := parts[0]
			newTaskName := parts[1]
			for i := range store.Areas {
				area := &store.Areas[i]
				if area.Title == areaName {
					if strings.Contains(newTaskName, creationSeparator) {
						names := strings.Split(newTaskName, creationSeparator)
						for _, name := range names {
							area.Tasks = append(area.Tasks, shared.Task{
								Title: name,
							})
						}
					} else {
						area.Tasks = append(area.Tasks, shared.Task{
							Title: newTaskName,
						})
					}
				}
			}
		} else if len(parts) == 3 {
			areaName := parts[0]
			taskName := parts[1]
			newSubtaskName := parts[2]
			for i := range store.Areas {
				area := &store.Areas[i]
				if area.Title != areaName {
					continue
				}
				for j := range area.Tasks {
					task := &area.Tasks[j]
					if task.Title == taskName {
						if strings.Contains(newSubtaskName, creationSeparator) {
							names := strings.Split(newSubtaskName, creationSeparator)
							for _, name := range names {
								task.Subtasks = append(task.Subtasks, shared.Subtask{
									Title: name,
								})
							}
						} else {
							task.Subtasks = append(task.Subtasks, shared.Subtask{
								Title: newSubtaskName,
							})
						}
					}
				}
			}
		} else {
			fmt.Println("add argument had more than 2 / in it, '/' is a reserved character")
			os.Exit(1)

		}
		store.Save()
	},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		// what's currenty in toComplete does not show in args
		// using __complete, the space after clean does now show in toComplete, it gets trimmed
		// "TestArae/clean<TAB>" shows both
		// "TestArae/clean <TAB>" shows nothing
		// "TestArae/clean\<TAB>" completes to clean kitchen
		// "TestArae/clean/<TAB>" completes to clean/SPACE
		// toComplete = strings.Join(append(args, toComplete), " ")
		toComplete = strings.ReplaceAll(toComplete, "_", " ")
		possibleCompletions := make([]string, 0)
		for i := range store.Areas {
			area := &store.Areas[i]
			if shared.LongestContainsOtherIgnoreCase(area.Title, toComplete) {
				if len(toComplete) <= len(area.Title) {
					possibleCompletions = append(possibleCompletions, strings.ReplaceAll(area.Title, " ", "_")+"/")
				}
				for j := range area.Tasks {
					task := &area.Tasks[j]
					rawFullTitle := area.Title + "/" + task.Title + "/"
					if shared.StartsWithIgnoreCase(rawFullTitle, toComplete) {
						possibleCompletions = append(possibleCompletions, strings.ReplaceAll(rawFullTitle, " ", "_"))
					}
				}
			}
		}
		// possibleCompletions = append(possibleCompletions, "args is "+strings.Join(args, "::"))
		// possibleCompletions = append(possibleCompletions, "=="+toComplete+"==")
		return possibleCompletions, cobra.ShellCompDirectiveNoSpace
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
