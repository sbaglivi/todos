package main

import (
	"encoding/json"
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
	getTasksByArea(int) (error, []task)
	getSubtaskByTask(int) (error, []subtask)
	getAllAreas() []area
	toggleTask(int)
	toggleSubtask(int)
	createArea(string) area
	createTask(int, string) task
	createSubtask(int, string) subtask
}

type jsonStore struct {
	fileName string
}

func (f jsonStore) getAllAreas() []area {
	jsonFile, err := os.Open(f.fileName)
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()
	byteValue, _ := io.ReadAll(jsonFile)
	var result map[int]interface{}
	json.Unmarshal([]byte(byteValue), &result)
	areas := make([]area, len(result))
	i := 0
	for k, v := range result {
		areas[i] = area{id: k, title: v.title, createdOn: v.createdOn}
	}
	return areas
}

func main() {

}
