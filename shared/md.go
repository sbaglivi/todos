package shared

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// func setup(cmd *cobra.Command, args []string) {
// 	fmt.Println("setup ran")
// 	sharedOpts.Store = markdownStore{filename: "todos.md"}
// }

type Area struct {
	Title string
	Tasks []Task
}

type Task struct {
	Title    string
	Subtasks []Subtask
	Done     bool
}

type Subtask struct {
	Title string
	Done  bool
}

func parseTask(s string) (string, bool) {
	re := regexp.MustCompile(`^ {2}- \[([x ])\] ([\w ]+)$`)
	match := re.FindStringSubmatch(s)
	if match != nil {
		return match[2], match[1] == "x"
	}
	return "", false
}

func parseSubtask(s string) (string, bool) {
	re := regexp.MustCompile(`^ {4}- \[([x ])\] ([\w ]+)$`)
	match := re.FindStringSubmatch(s)
	if match != nil {
		return match[2], match[1] == "x"
	}
	return "", false
}

func printAreas(areas []Area) {
	fmt.Print(areasToString)
}

func startsWithIgnoreCase(s, prefix string) bool {
	s = strings.ToLower(s)
	prefix = strings.ToLower(prefix)
	return strings.HasPrefix(s, prefix)
}

func areasToString(areas []Area) string {
	var builder strings.Builder
	for _, a := range areas {
		builder.WriteString(fmt.Sprintf("%s\n", a.Title))
	}
	return builder.String()
}

func (s *MarkdownStore) FindArea(prefix string) (error, Area) {
	matches := make([]Area, 0)
	for _, a := range s.areas {
		if !startsWithIgnoreCase(a.Title, prefix) {
			continue
		}
		matches = append(matches, a)
	}
	if len(matches) == 1 {
		return nil, matches[0]
	} else if len(matches) > 1 {
		return errors.New(fmt.Sprintf("Found %d matches for input \"%s\":\n%s", len(matches), prefix, areasToString(matches))), Area{}
	} else {
		return errors.New(fmt.Sprintf("Found no matches for input \"%s\"", prefix)), Area{}
	}
}

type MarkdownStore struct {
	filename string
	areas    []Area
}

func Setup() MarkdownStore {
	store := MarkdownStore{filename: "todos.md"}
	store.load()
	return store

}

func (s *MarkdownStore) load() {
	file, err := os.Open(s.filename)
	if err != nil {
		return
	}
	defer file.Close()

	var areas []Area
	var currentArea Area
	var currentTask Task
	areaPrefix := "- "
	taskPrefix := "  - ["
	subtaskPrefix := "    - ["
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, areaPrefix) {
			// This is an area
			if currentTask.Title != "" {
				currentArea.Tasks = append(currentArea.Tasks, currentTask)
				currentTask = Task{}
			}
			if currentArea.Title != "" {
				areas = append(areas, currentArea)
			}
			currentArea = Area{Title: strings.TrimPrefix(line, "- ")}
		} else if strings.HasPrefix(line, taskPrefix) {
			// This is a task
			if currentTask.Title != "" {
				currentArea.Tasks = append(currentArea.Tasks, currentTask)
			}
			if currentArea.Title != "" {
				title, done := parseTask(line)
				if title == "" {
					continue
				}
				currentTask = Task{Title: title, Done: done}
			}
		} else if strings.HasPrefix(line, subtaskPrefix) {
			if currentTask.Title != "" {
				title, done := parseSubtask(line)
				if title == "" {
					continue
				}
				currentTask.Subtasks = append(currentTask.Subtasks, Subtask{
					Title: title,
					Done:  done,
				})
			}
		}
	}

	// Append the last area and task if not already appended
	if currentTask.Title != "" {
		currentArea.Tasks = append(currentArea.Tasks, currentTask)
	}
	if currentArea.Title != "" {
		areas = append(areas, currentArea)
	}

	s.areas = areas
}

func (s *MarkdownStore) areasAsMDString() string {
	var builder strings.Builder
	done := " "
	for _, area := range s.areas {
		builder.WriteString(fmt.Sprintf("- %s\n", area.Title))
		for _, task := range area.Tasks {
			done = " "
			if task.Done {
				done = "x"
			}
			builder.WriteString(fmt.Sprintf("%*s- [%s] %s\n", 2, "", done, task.Title))
			for _, subtask := range task.Subtasks {
				done = " "
				if subtask.Done {
					done = "x"
				}
				builder.WriteString(fmt.Sprintf("%*s- [%s] %s\n", 4, "", done, subtask.Title))
			}
		}
	}
	return builder.String()
}

func (s *MarkdownStore) Print() {
	fmt.Println(s.areasAsMDString())
}

func (s *MarkdownStore) Save() {
	newFile, err := os.Create(s.filename)
	if err != nil {
		panic(err)
	}
	defer newFile.Close()
	toWrite := s.areasAsMDString()
	_, err = newFile.Write([]byte(toWrite))
	if err != nil {
		panic(err)
	}
}

func (s *MarkdownStore) ToggleTask(prefix string) {
	taskMatches := make([]*Task, 0)
	subtaskMatches := make([]*Subtask, 0)
	for i := range s.areas {
		area := &s.areas[i]
		for j := range area.Tasks {
			task := &area.Tasks[j]
			if startsWithIgnoreCase(task.Title, prefix) {
				taskMatches = append(taskMatches, task)
			}
			for k := range task.Subtasks {
				subtask := &task.Subtasks[k]
				if startsWithIgnoreCase(subtask.Title, prefix) {
					subtaskMatches = append(subtaskMatches, subtask)
				}
			}
		}
	}
	if len(taskMatches)+len(subtaskMatches) != 1 {
		fmt.Printf("search for task with prefix \"%s\" found %d results\n", prefix, len(taskMatches)+len(subtaskMatches))
		os.Exit(1)
	}
	if len(taskMatches) == 1 {
		taskMatches[0].Done = !taskMatches[0].Done
	} else {
		subtaskMatches[0].Done = !subtaskMatches[0].Done
	}
}

func oldoldmain() {
	s := MarkdownStore{filename: "notes.md"}
	s.load()
	input := "not exists"
	err, match := s.FindArea(input)
	if err != nil {
		panic(err)
	}
	fmt.Println(match.Title)

}

func saveToFile(filename, toWrite string) {
	newFile, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer newFile.Close()
	_, err = newFile.Write([]byte(toWrite))
	if err != nil {
		panic(err)
	}
}

func todosToString(areas []Area) string {
	var builder strings.Builder
	for _, area := range areas {
		builder.WriteString(fmt.Sprintf("- %s\n", area.Title))
		for _, task := range area.Tasks {
			done := " "
			if task.Done {
				done = "x"
			}
			builder.WriteString(fmt.Sprintf("%*s- [%s] %s\n", 2, "", done, task.Title))
			for _, subtask := range task.Subtasks {
				done := " "
				if subtask.Done {
					done = "x"
				}
				builder.WriteString(fmt.Sprintf("%*s- [%s] %s\n", 4, "", done, subtask.Title))
			}
		}
	}
	return builder.String()
}

// func main() {

// 	store := jsonStore{fileName: "todos.json"}
// 	store.Init()
// 	for _, a := range(store.Areas) {
// 		black
// 	}
// }