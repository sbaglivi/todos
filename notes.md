- Area:
  - [ ] Task
  - [ ] Another task
    - [x] A subtask that is completed
    - [ ] one that is not
- Empty area
- New area
    - [ ] subtask misplaced
  - [x] correct task completed
    - [x] with a subtask
    - [ ] and another one
I'm confused because when working with json I don't get to just load part of the todos or of the file.
I retrieve the whole thing at once. 
Which makes it seem really inefficient not to just load it all in memory


How do I want commands to work?
Functionalities required:
- ls areas
`ls`
- ls tasks in a single area (for now with subtasks I think, I'll avoid creating a list task)
`ls areaname` 
- create area
`create -a areaname`
- create task in a given area (+ possibly a shortcut to also create some subtasks in the same line)
`create -a areaname -t taskname:subtask1;subtask2`
- create subtask in a given task
`create [-a areaname] -t taskname -s subtaskname`
- toggle completion of subtask (should it auto complete the task if all subtasks are completed?)
`toggle [-a areaname] [-t taskname] -s subtaskname`
`toggle areaname/taskName/subtaskName`
- toggle completion of task (? conditioned on one above maybe, but it still could be useful)
- edit task name
`rename -t oldname newname`
- edit subtask name
`rename -s oldname newname`
- edit area name 
`rename -a oldname newname`
- delete subtask
`delete -s name`
- delete task
`delete -t name`
- delete area
`delete -a name`

Improvements:
- printing to use colors
- add expiration time to tasks / subtasks?
- automatic removal of completed?