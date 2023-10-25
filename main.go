package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
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
type jsonArea struct {
	Id        int
	Title     string
	CreatedOn string
	Tasks     []jsonTask
}

type jsonTask struct {
	id          int
	title       string
	done        bool
	createdOn   time.Time
	completedOn time.Time
	subtasks    []jsonSubtask
}

type jsonSubtask struct {
	id          int
	title       string
	done        bool
	createdOn   time.Time
	completedOn time.Time
}

func (s *jsonStore) Init() {
	jsonFile, err := os.Open(s.fileName)
	if err != nil {
		// file does not exist yet
		s.areas = make([]jsonArea, 0)
	} else {
		defer jsonFile.Close()
		byteValue, _ := io.ReadAll(jsonFile)
		json.Unmarshal([]byte(byteValue), &s.areas)
		highestCount := 0
		for _, a := range s.areas {
			if a.Id > highestCount {
				highestCount = a.Id
			}
		}
		s.highestId = highestCount
	}
}

func (s *jsonStore) CreateArea(title string) jsonArea {
	s.areas = append(s.areas, jsonArea{
		Id:        s.highestId + 1,
		Title:     title,
		Tasks:     make([]jsonTask, 0),
		CreatedOn: time.Now().UTC().Format(time.RFC3339),
	})
	s.highestId++
	return s.areas[len(s.areas)-1]
}

func (s *jsonStore) Save() {
	for _, a := range s.areas {
		fmt.Printf("%+v\n", a)
	}
	areasJson, err := json.MarshalIndent(s.areas, "", "  ")
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(s.fileName, areasJson, 0644)
	if err != nil {
		panic(err)
	}
}

func main() {
	store := jsonStore{fileName: "todos.json"}
	store.Init()
	store.CreateArea("new area")
	store.CreateArea("another are")
	store.Save()
}
