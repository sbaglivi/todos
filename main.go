package main

import (
	"encoding/json"
	"io"
	"os"
	"time"

	"github.com/fatih/color"
)

type subtask struct {
	id          int
	title       string
	done        bool
	createdOn   time.Time
	completedOn time.Time
}

type task struct {
	id          int
	title       string
	done        bool
	createdOn   time.Time
	completedOn time.Time
}

type area struct {
	id        int
	title     string
	createdOn time.Time
}

type Store interface {
	GetTasksByArea(int) []task
	GetSubtaskByTask(int) []subtask
	GetAllAreas() []area
	ToggleTask(int)
	ToggleSubtask(int)
	CreateArea(string) area
	CreateTask(int, string) task
	CreateSubtask(int, string) subtask
	Init()
}

type jsonStore struct {
	highestId int
	fileName  string
	areas     []jsonArea
}

/*
if not capitalized they're not exported ?
When the first letter is capitalised, the identifier is public to any piece of code that you want to use.
When the first letter is lowercase, the identifier is private and could only be accessed within the package it was declared.
https://stackoverflow.com/questions/26327391/json-marshalstruct-returns
*/
type jsonWrapper struct {
	HighestId int
	Areas     []jsonArea
}
type jsonArea struct {
	Id        int
	Title     string
	CreatedOn string
	Tasks     []jsonTask
}

type jsonTask struct {
	Id          int
	Title       string
	Done        bool
	CreatedOn   string
	CompletedOn string
	Subtasks    []jsonSubtask
}

type jsonSubtask struct {
	Id          int
	Title       string
	Done        bool
	CreatedOn   string
	CompletedOn string
}

func (s *jsonStore) Init() {
	jsonFile, err := os.Open(s.fileName)
	if err != nil {
		// file does not exist yet
		s.areas = make([]jsonArea, 0)
	} else {
		var wrapper jsonWrapper
		defer jsonFile.Close()
		byteValue, _ := io.ReadAll(jsonFile)
		json.Unmarshal([]byte(byteValue), &wrapper)
		s.highestId = wrapper.HighestId
		s.areas = wrapper.Areas
	}
}

func (s *jsonStore) CreateArea(title string) jsonArea {
	s.areas = append(s.areas, jsonArea{
		Id:        s.highestId,
		Title:     title,
		Tasks:     make([]jsonTask, 0),
		CreatedOn: time.Now().UTC().Format(time.RFC3339),
	})
	s.highestId++
	return s.areas[len(s.areas)-1]
}

func subtasksFromTitles(highestId *int, titles []string) []jsonSubtask {
	subtasks := make([]jsonSubtask, len(titles))
	for i, t := range titles {
		subtasks[i] = jsonSubtask{
			Id:        *highestId,
			Title:     t,
			CreatedOn: time.Now().UTC().Format(time.RFC3339),
		}
		*highestId++
	}
	return subtasks
}

func (s *jsonStore) CreateTask(id int, title string, subtasks ...string) {
	for i := range s.areas {
		if s.areas[i].Id == id {
			taskId := s.highestId
			s.highestId++
			var jSubtasks []jsonSubtask
			if len(subtasks) > 0 {
				jSubtasks = subtasksFromTitles(&s.highestId, subtasks)
			} else {
				jSubtasks = make([]jsonSubtask, 0)
			}
			s.areas[i].Tasks = append(s.areas[i].Tasks, jsonTask{
				Title:     title,
				Id:        taskId,
				Subtasks:  jSubtasks,
				CreatedOn: time.Now().UTC().Format(time.RFC3339),
			})
		}
	}

}

func (s *jsonStore) Save(print bool) {
	if print {
		printTodos(s.areas)
	}
	wrapper := jsonWrapper{
		HighestId: s.highestId,
		Areas:     s.areas,
	}
	areasJson, err := json.MarshalIndent(wrapper, "", "  ")
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(s.fileName, areasJson, 0644)
	if err != nil {
		panic(err)
	}
}

func (s *jsonStore) CreateSubtask(id int, title string) {
	for i := range s.areas {
		for j := range s.areas[i].Tasks {
			if s.areas[i].Tasks[j].Id == id {
				s.areas[i].Tasks[j].Subtasks = append(s.areas[i].Tasks[j].Subtasks, jsonSubtask{
					Title:     title,
					Id:        s.highestId,
					CreatedOn: time.Now().UTC().Format(time.RFC3339),
				})
				s.highestId++
			}
		}
	}
}
func (s *jsonStore) createSomeShit() {
	s.CreateArea("new area")
	s.CreateTask(0, "my first task!")
	s.CreateSubtask(1, "this is a subtask!")
}

func printTodos(areas []jsonArea) {
	for _, a := range areas {
		// fmt.Printf("%d: %s\n", a.Id, a.Title)
		color.Cyan("%d: %s\n", a.Id, a.Title)
		for _, t := range a.Tasks {
			color.Yellow("%2s%d: %s\n", "", t.Id, t.Title)
			// fmt.Printf("%2s%d: %s\n", "", t.Id, t.Title)
			for _, s := range t.Subtasks {
				color.Red("%4s%d: %s\n", "", s.Id, s.Title)
				// fmt.Printf("%4s%d: %s\n", "", s.Id, s.Title)
			}
		}
	}
}

func main() {
	store := jsonStore{fileName: "todos.json"}
	store.Init()
	store.CreateTask(0, "new tech", "is it working?", "please tell me so")
	printTodos(store.areas)
	store.Save(false)
}
